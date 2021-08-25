package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	jww "github.com/spf13/jwalterweatherman"
	base "goEs/controller"
	"goEs/utils"
	"reflect"
	"strconv"
	"time"
)

//数据结构
type Employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

var now = time.Now()

//Get方法带参数
func GetElasticById(c *gin.Context) {
	id := c.Params.ByName("id")
	var bb Employee
	//通过id查找
	get1, err := utils.Client.Get().Index("megacorp").Id(id).Do(context.Background())
	if err != nil {
		c.JSON(200, gin.H{
			"data":       err,
			"code":       200,
			"serverTime": now,
		})
		return
	}
	if get1.Found {
		err := json.Unmarshal(get1.Source, &bb)
		if err != nil {
			fmt.Println(err)
		}
	}

	c.JSON(200, gin.H{
		"data":       bb,
		"code":       200,
		"serverTime": now,
	})

}

//Post方法,创建
func Create(c *gin.Context) {
	//使用结构体
	e1 := Employee{"Jane", "Smith", 32, "I like to collect rock albums", []string{"music"}}
	_, err := utils.Client.Index().
		Index("megacorp").
		Id("1").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	utils.Client.Index().Index("megacorp").WaitForActiveShards("2").Do(context.Background())
	//使用字符串
	e2 := `{"first_name":"John","last_name":"Smith","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	_, err = utils.Client.Index().
		Index("megacorp").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	e3 := `{"first_name":"Douglas","last_name":"Fir","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
	_, err = utils.Client.Index().
		Index("megacorp").
		Id("3").
		BodyJson(e3).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	c.JSON(200, base.SuccessRep("新增成功"))
}

//Delete方法
func Destroy(c *gin.Context) {
	id := c.Params.ByName("id")
	res, err := utils.Client.Delete().Index("megacorp").
		Id(id).
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}

	c.JSON(200, base.SuccessRep(fmt.Sprintf("delete result %s\n", res.Result)))
}

//查找
func Query(c *gin.Context) {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = utils.Client.Search("megacorp").Do(context.Background())
	re1 := printEmployee(res, err)

	//字段相等
	q := elastic.NewQueryStringQuery("last_name:Smith")
	res, err = utils.Client.Search("megacorp").Query(q).Do(context.Background())
	if err != nil {
		jww.ERROR.Printf("Search error=%s", err.Error())
	}
	re2 := printEmployee(res, err)

	//条件查询
	//年龄大于30岁的
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
	res, err = utils.Client.Search("megacorp").Query(q).Do(context.Background())
	re3 := printEmployee(res, err)


	//短语搜索 搜索about字段中有 rock climbing
	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = utils.Client.Search("megacorp").Query(matchPhraseQuery).Do(context.Background())
	re4 := printEmployee(res, err)

	//分析 interests
	aggs := elastic.NewTermsAggregation().Field("interests")
	res, err = utils.Client.Search("megacorp").Aggregation("all_interests", aggs).Do(context.Background())
	re5 := printEmployee(res, err)
	re1 = append(re1, re2...)
	re1 = append(re1, re3...)
	re1 = append(re1, re4...)
	re1 = append(re1, re5...)
	c.JSON(200, base.SuccessRep(re1))
}

////简单分页
func List(c *gin.Context) {
	size, _ := strconv.Atoi(c.PostForm("size"))
	page, _ := strconv.Atoi(c.PostForm("page"))
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res, err := utils.Client.Search("megacorp").
		Size(size).
		From((page - 1) * size).
		Do(context.Background())
	re1 := printEmployee(res, err)
	c.JSON(200, base.SuccessRep(re1))

}

//打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) (result []Employee) {
	if err != nil {
		jww.ERROR.Printf("Search error=%s", err.Error())
		return result
	}

	var typ Employee
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(Employee)
		result = append(result, t)
	}
	return result
}

//Post方法,创建
func WaitForActiveShards(c *gin.Context) {
	fmt.Println(1111)
	ca,err := utils.Client.CatShards().Index("megacorp").Do(context.Background())
	fmt.Println("ca===>",ca)
	fmt.Println("err===>",err)
	//_, err := utils.Client.Index().Index("megacorp").WaitForActiveShards("2").Do(context.Background())
	//fmt.Println("err====>",err)
	//if err != nil {
	//	jww.ERROR.Printf("Search error=%s", err.Error())
	//	c.JSON(200, base.ErrorRep(500,"修改失败"))
	//}
	c.JSON(200, base.SuccessRep(nil))
}
