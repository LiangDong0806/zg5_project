package models

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"log"
	"zg5/work01/server/common/global"
)

// TODO: 初始化Elastic
func InitElastic() {
	var err error
	dsn := fmt.Sprintf("http://%s:%d", global.ServerConfig.ElasticSearch.Host, global.ServerConfig.ElasticSearch.Port)
	global.Client, err = elastic.NewClient(elastic.SetURL(dsn), elastic.SetSniff(false))
	if err != nil {
		log.Println("Error creating elastic")
	}
	fmt.Println(global.Client, dsn, "[p[[[")
}
func EscIns(body map[string]interface{}) error {
	fmt.Println(global.Client, "[[[")
	_, err := global.Client.Index().Index(global.ServerConfig.ElasticSearch.Index).BodyJson(body).Do(context.Background())
	if err != nil {
		return errors.New("err creating")
	}
	return nil
}
