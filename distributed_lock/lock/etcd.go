package lock

import (
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdClient struct {
	*clientv3.Client
}

var eClient *etcdClient
var once sync.Once

func NewEtcdClient() *etcdClient {
	once.Do(func() {
		c, err := clientv3.New(clientv3.Config{
			Endpoints: []string{"http://123.57.167.85:2379"},
			// Username:  "root",
			// Password:  "ficates",
		})
		if err != nil {
			panic(err)
		}
		eClient = &etcdClient{c}
	})
	return eClient
}
