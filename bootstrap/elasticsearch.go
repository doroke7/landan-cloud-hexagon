package bootstrap

import (
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearch() (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: CONFIG.ELASTICSEARCH.HOSTS,
		Username:  CONFIG.ELASTICSEARCH.USER,
		Password:  CONFIG.ELASTICSEARCH.PASS,
	})
}
