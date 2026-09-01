package bootstrap

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearch() (*elasticsearch.Client, error) {
	oClient, oErr := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: CONFIG.ELASTICSEARCH.HOSTS,
		Username:  CONFIG.ELASTICSEARCH.USERNAME,
		Password:  CONFIG.ELASTICSEARCH.PASSWORD,
	})

	log.Info("[INFO] ELASTICSEARCH 連線完成.", "addr", strings.Join(CONFIG.ELASTICSEARCH.HOSTS, ","))

	return oClient, oErr
}

/**
ES 相當把 MYSQL 資料的 “條件” 又區分2類，
1. SQL： WHERE MATCH -> ES : query.must
2. SQL: WHERE -> ES: query.filter


SELECT *
FROM games
WHERE
    name LIKE '%poker%'
    AND status='online'
    AND price < 100
ORDER BY relevance;


{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "name": "poker"
          }
        }
      ],
      "filter": [
        {
          "term": {
            "status": "online"
          }
        },
        {
          "range": {
            "price": {
              "lt": 100
            }
          }
        }
      ]
    }
  }
}

*/
