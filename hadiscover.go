package main

import (
	"flag"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"os"
)

func getEnvOrDefault(envName, defaultValue string) string {
	val := os.Getenv(envName)
	if val == "" {
		return defaultValue
	}
	return val
}

var filename = flag.String(
	"config",
	getEnvOrDefault("HADISCOVER_TPL", "./haproxy.cfg.tpl"),
	"Template config file used for HAproxy")
var etcdHost = flag.String(
	"etcd",
	getEnvOrDefault("HADISCOVER_ETCD", "http://localhost:4001"),
	"etcd server(s)")
var etcdKey = flag.String(
	"key",
	getEnvOrDefault("HADISCOVER_KEY", "services"),
	"etcd root key to look for")

var configFile = ".haproxy.cfg"
var name = "back"
var haproxyExec = "haproxy"

func reloadConf(haproxyExec string, etcdClient *etcd.Client) error {
	backends, _ := GetBackends(etcdClient, *etcdKey, name)

	err := createConfigFile(backends, *filename, configFile)
	if err != nil {
		log.Println("Cannot generate haproxy configuration: ", err)
		return err
	}
	return reloadHAproxy(haproxyExec, configFile)
}

func main() {
	flag.Parse()

	var etcdClient = etcd.NewClient([]string{*etcdHost})
	err := reloadConf(haproxyExec, etcdClient)
	if err != nil {
		log.Println("Cannot reload haproxy: ", err)
	}

	changeChan := make(chan *etcd.Response)
	stopChan := make(chan bool)

	go func() {
		for msg := range changeChan {
			reload := (msg.PrevNode == nil) || (msg.PrevNode.Key != msg.Node.Key) || (msg.Action != "set")
			if reload {
				err := reloadConf(haproxyExec, etcdClient)
				if err != nil {
					log.Println("Cannot reload haproxy: ", err, msg)
				}
			}
		}
	}()

	log.Println("Start watching changes in etcd")
	if _, err := etcdClient.Watch(*etcdKey, 0, true, changeChan, stopChan); err != nil {
		log.Println("Cannot register watcher for changes in etcd: ", err)
	}

}
