package main

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type Backend struct {
	Name string
	Ip   string
	Port string
}

func GetBackends(client *etcd.Client, followServicesString, backendName string) (map[string][]Backend, error) {

	resp, err := client.Get("sp/", false, true)
	if err != nil {
		log.Println("Error when reading etcd: ", err)
		return nil, err
	} else {
		backends := make(map[string][]Backend)
		followServices := strings.Split(followServicesString, ",")

		for _, element := range resp.Node.Nodes {
			key := (*element).Key // key format is: /sp/<SP_GROUP>/<IP:PORT>
			splited := strings.Split(key, "/")

			serviceFlag := false
			for _, followService := range followServices {
				if followService == splited[2] {
					serviceFlag = true
					break
				}
			}
			if serviceFlag == false {
				continue
			}

			for index2, element2 := range (*element).Nodes {
				key := (*element2).Key // key format is: /sp/<SP_GROUP>/<IP:PORT>
				splited2 := strings.Split(key, "/")
				service := strings.Split(splited2[3], ":")

				serviceType := splited2[2]

				backend := Backend{Name: fmt.Sprintf("back-%v", index2), Ip: service[0], Port: service[1]}

				backends[serviceType] = append(backends[serviceType], backend)
			}
		}
		return backends, nil
	}

}
