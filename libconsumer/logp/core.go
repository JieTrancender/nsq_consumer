package logp

import (
	"io/ioutil"
	golog "log"
	"os"
	"sync/atomic"
	"unsafe"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_log unsafe.Pointer // Pointer to a coreLogger, Access via atomic.LoadPointer.
)

func init() {
	storeLogger(&coreLogger{
		selectors:  map[string]struct{}{},
		rootLogger: zap.NewNop(),
		logger:     newLogger(zap.NewNop(), ""),
	})
}

type coreLogger struct {
	selectors  map[string]struct{}
	rootLogger *zap.Logger
	logger     *Logger
}

func Configure(cfg Config) error {
	return ConfigureWithOutputs(cfg)
}

func ConfigureWithOutputs(cfg Config, outputs ...zapcore.Core) error {
	var (
		sink zapcore.Core
		err  error
	)

	sink, err = createLogOutput(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to build log output")
	}

	// Default logger is always discard, debug level below will posibly re-enable it
	golog.SetOutput(ioutil.Discard)

	// Enabled selectors when debug is enabled.
	selectors := make(map[string]struct{}, len(cfg.Selectors))
	// if cfg.Level.Enabled(DebugLevel) && len(cfg.Selectors) > 0 {
	// 	for _, sel := range cfg.Selectors {
	// 		selectors[strings.TrimSpace(sel)] = struct{}{}

	// 		if len(selectors) == 0 {
	// 			selectors["*"] = struct{}{}
	// 		}

	// 		// Re-enable the default go logger output when either stdlog or all selector is enabled.
	// 		_, stdlogEnabled := selectors["stdlog"]
	// 		_, allEnabled := selectors["*"]
	// 		if stdlogEnabled || allEnabled {
	// 			golog.SetOutput(_defaultGoLog)
	// 		}

	// 		sink = selectiveWrapper(sink, selectors)
	// 	}
	// }

	sink = newMultiCore(append(outputs, sink)...)
	root := zap.New(sink, makeOptions(cfg)...)
	storeLogger(&coreLogger{
		selectors:  selectors,
		rootLogger: root,
		logger:     newLogger(root, ""),
	})
	return nil
}

func createLogOutput(cfg Config) (zapcore.Core, error) {
	switch {
	case cfg.toIODiscard:
		return makeDiscardOutput(cfg)
	case cfg.ToStderr:
		return makeStderrOutput(cfg)
		// case cfg.ToFiles:
		// return makeFileOutput(cfg)
	}

	switch cfg.environment {
	case SystemdEnvironment, ContainerEnvironment:
		return makeStderrOutput(cfg)
	case MacOSServiceEnvironment, WindowsServiceEnvironment:
		fallthrough
	default:
		// return makeFileOutput(cfg)
		return makeStderrOutput(cfg)
	}
}

// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return loadLogger().rootLogger.Sync()
}

func makeOptions(cfg Config) []zap.Option {
	var options []zap.Option
	if cfg.addCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.development {
		options = append(options, zap.Development())
	}
	return options
}

func makeStderrOutput(cfg Config) (zapcore.Core, error) {
	stderr := zapcore.Lock(os.Stderr)
	return newCore(cfg, buildEncoder(cfg), stderr, cfg.Level.ZapLevel()), nil
}

func makeDiscardOutput(cfg Config) (zapcore.Core, error) {
	discard := zapcore.AddSync(ioutil.Discard)
	return newCore(cfg, buildEncoder(cfg), discard, cfg.Level.ZapLevel()), nil
}

func newCore(cfg Config, enc zapcore.Encoder, ws zapcore.WriteSyncer, enab zapcore.LevelEnabler) zapcore.Core {
	return wrappedCore(cfg, zapcore.NewCore(enc, ws, enab))
}

func wrappedCore(cfg Config, core zapcore.Core) zapcore.Core {
	// if cfg.ECSEnabled {
	// 	return ecszap.WrapCore(core)
	// }

	return core
}

func loadLogger() *coreLogger {
	p := atomic.LoadPointer(&_log)
	return (*coreLogger)(p)
}

func storeLogger(l *coreLogger) {
	if old := loadLogger(); old != nil {
		_ = old.rootLogger.Sync()
	}
	atomic.StorePointer(&_log, unsafe.Pointer(l))
}

// newMultiCore creates a sink that sends to multiple cores.
func newMultiCore(cores ...zapcore.Core) zapcore.Core {
	return &multiCore{cores: cores}
}

// multiCore allows multiple cores to be used for logging.
type multiCore struct {
	cores []zapcore.Core
}

// Enabled returns true if the level is enabled in any one of the cores.
func (m multiCore) Enabled(level zapcore.Level) bool {
	for _, core := range m.cores {
		if core.Enabled(level) {
			return true
		}
	}
	return false
}

// With creates a new multiCore with each core set with the given fields
func (m multiCore) With(fields []zapcore.Field) zapcore.Core {
	cores := make([]zapcore.Core, len(m.cores))
	for i, core := range m.cores {
		cores[i] = core.With(fields)
	}
	return &multiCore{cores: cores}
}

// Check will place each core that checks for that entry.
func (m multiCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	for _, core := range m.cores {
		checked = core.Check(entry, checked)
	}
	return checked
}

// Write writes the entry to each core.
func (m multiCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	var errs error
	for _, core := range m.cores {
		if err := core.Write(entry, fields); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// Sync syncs each core.
func (m multiCore) Sync() error {
	var errs error
	for _, core := range m.cores {
		if err := core.Sync(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}
