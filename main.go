package main

import (
	"fmt"
	"github.com/go-pg/pg"
	router "goEs/route"
	"net/http"
	"time"

	"goEs/utils"

)


//初始化
func init() {
	utils.InitConfig()
	utils.InitElastic()
	router.InitRoute()
	utils.InitLooger()
}



func main()  {

	db := utils.NewDataBase()
	//测试查库
	var n string
	_, err := db.GetPg().QueryOne(pg.Scan(&n), "select now() ")
	if err != nil {
		panic(err)
	}
	fmt.Println(n)

	router := &http.Server{
		Addr:           utils.Conf.Server.Host+utils.Conf.Server.Port,
		Handler:        router.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	router.ListenAndServe()
}
