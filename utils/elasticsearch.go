package utils

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	jww "github.com/spf13/jwalterweatherman"
)

var Client *elastic.Client

func InitElastic() {
	//errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	dbConf := fmt.Sprintf("http://%s:%s", Conf.Databases.Host, Conf.Databases.Port)
	var err error
	//这个地方有个小坑 不加上elastic.SetSniff(false) 会连接不上
	Client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(dbConf))
	if err != nil {
		jww.ERROR.Printf("NewClient error=%s", err.Error())
		return
	}
	_, _, err = Client.Ping(dbConf).Do(context.Background())
	if err != nil {
		jww.ERROR.Printf("Ping error=%s", err.Error())
		return
	}

	vsesion, err := Client.ElasticsearchVersion(dbConf)
	fmt.Printf("ElasticsearchVersion = %+v\n",vsesion)
	if err != nil {
		jww.ERROR.Printf("ElasticsearchVersion error=%s", err.Error())
		return
	}
}
