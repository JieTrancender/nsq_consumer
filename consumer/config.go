package consumer

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/JieTrancender/nsq_consumer/internal/version"
	"github.com/JieTrancender/nsq_consumer/libconsumer/common"
	"github.com/JieTrancender/nsq_consumer/libconsumer/logp"
	"github.com/nsqio/go-nsq"
)

type config struct {
	DialTimeout time.Duration `config:"dial_timeout"`

	ReadTimeout  time.Duration `config:"read_timeout"`
	WriteTimeout time.Duration `config:"write_timeout"`

	MaxInFlight int `config:"max_in_flight"`

	Test int `config:"test"`
}

func newNSQConfig(rawConfig *common.Config, consumerType string) *nsq.Config {
	cfg := nsq.NewConfig()
	cfg.UserAgent = fmt.Sprintf("nsq_consumer/%s go-nsq/%s", version.GetDefaultVersion(), nsq.VERSION)

	configName := strings.ToLower(consumerType)
	if !rawConfig.HasField(configName) {
		return cfg
	}

	sub, err := rawConfig.Child(configName, -1)
	if err != nil {
		return cfg
	}

	c := &config{
		DialTimeout:  6 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 1 * time.Second,
		MaxInFlight:  200,
	}
	if err := sub.Unpack(&c); err != nil {
		logp.L().Errorf("newNSQConfig unpack fail: %v", err)
		return cfg
	}

	val := reflect.ValueOf(c).Elem()
	typ := val.Type()
	cfgVal := reflect.ValueOf(cfg).Elem()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := cfgVal.FieldByName(field.Name)
		if fieldVal.IsValid() {
			fieldVal.Set(reflect.ValueOf(val.Field(i).Interface()))
		}
	}

	return cfg
}
