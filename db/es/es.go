package es

import (
	"sync"

	"github.com/olivere/elastic"
)

var esClient *elastic.Client

var once sync.Once

func NewEs() *elastic.Client {
	once.Do(func() {
		client, _ := elastic.NewClient(
			elastic.SetURL("http://123.57.167.85:9200"),
			elastic.SetSniff(false),
		)
		esClient = client
	})

	return esClient
}
