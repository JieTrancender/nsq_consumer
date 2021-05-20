package instance

import (
	"flag"
	"fmt"

	"github.com/JieTrancender/nsq_to_consumer/internal/app"
)

var etcdEndpoints = app.StringArray{}

func init() {
	fs := flag.CommandLine
	// etcdEndpoints := app.StringArray{}
	fs.Var(&etcdEndpoints, "etcd-endpoints", "etcd endpoint, may be given multi times")
}

func Run(settings Settings) error {
	etcdEndpoints, _ := settings.RunFlags.GetStringArray("etcd-endpoints")
	etcdPath, _ := settings.RunFlags.GetString("etcd-path")
	etcdUsername, _ := settings.RunFlags.GetString("etcd-username")
	etcdPassword, _ := settings.RunFlags.GetString("etcd-password")
	fmt.Printf("etcd(%v):%s %s:%s\n", etcdEndpoints, etcdPath, etcdUsername, etcdPassword)

	return nil
}
