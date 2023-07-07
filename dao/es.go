package dao

import (
	"sync"

	"github.com/olivere/elastic"
)

var esClient *elastic.Client
var once sync.Once

func NewEs() *elastic.Client {
	once.Do(func() {
		client, err := elastic.NewClient(
			elastic.SetURL("127.0.0.1:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			panic(err)
		}
		esClient = client
	})

	return esClient
}
