package dao

import (
	"github.com/olivere/elastic"
)

func NewEs() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL("127.0.0.1:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	return client
}
