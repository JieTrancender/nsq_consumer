module github.com/JieTrancender/nsq_consumer

go 1.14

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/elastic/go-ucfg v0.8.3
	github.com/google/uuid v1.2.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869
	github.com/nsqio/go-nsq v1.0.8
	github.com/olivere/elastic/v7 v7.0.27
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.14.0
	google.golang.org/grpc v1.40.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace google.golang.org/grpc v1.40.0 => google.golang.org/grpc v1.26.0
