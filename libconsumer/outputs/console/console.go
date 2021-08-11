package console

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/consumer"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/JieTrancender/nsq_consumer/libconsumer/outputs"
)

type console struct {
	logger *logp.Logger
	out    *os.File
	writer *bufio.Writer
}

func init() {
	outputs.RegisterType("console", makeConsole)
}

func makeConsole(
	consumerInfo consumer.Info,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig
	err := cfg.Unpack(&config)
	if err != nil {
		return outputs.Group{}, err
	}

	c, err := newConsole()
	if err != nil {
		return outputs.Group{}, fmt.Errorf("console output initialization failed with: %v", err)
	}

	// check stdout actually being available
	if runtime.GOOS != "windows" {
		if _, err = c.out.Stat(); err != nil {
			err = fmt.Errorf("console output initialization failed with: %v", err)
			return outputs.Fail(err)
		}
	}

	return outputs.Success(config.BatchSize, 0, c)
}

func newConsole() (*console, error) {
	c := &console{
		logger: logp.NewLogger("console"),
		out:    os.Stdout,
	}
	c.writer = bufio.NewWriterSize(c.out, 8*1024)
	return c, nil
}

func (c *console) Close() error { return nil }

var nl = []byte("\n")

func (c *console) Publish(_ context.Context, m consumer.Message) error {
	if err := c.writeBuffer(m.Body()); err != nil {
		c.logger.Errorf("Unable to publish message to console: %+v", err)
		return err
	}

	if err := c.writeBuffer(nl); err != nil {
		c.logger.Errorf("Error when appending newline to console: %+v", err)
		return err
	}
	c.writer.Flush()

	return nil
}

func (c *console) writeBuffer(buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := c.writer.Write(buf[written:])
		if err != nil {
			return err
		}

		written += n
	}
	return nil
}

func (c *console) String() string {
	return "console"
}
