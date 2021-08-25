package index

import (
	//"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	base "goEs/controller"
	"goEs/utils"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
)

func Create(c *gin.Context) {
	mapping := `{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"tags":{
				"type":"keyword"
			},
			"location":{
				"type":"geo_point"
			},
			"suggest_field":{
				"type":"completion"
			}
		}
	}
}`

	ctx := context.Background()
	createIndex, err := utils.Client.CreateIndex("twitter").BodyString(mapping).Do(ctx)
	if err != nil {
		c.JSON(500, base.ErrorRep(500, err.Error()))
		return
	}
	if createIndex == nil {
		c.JSON(500, base.ErrorRep(500, "createIndex is nil"))
		return
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
		c.JSON(500, base.ErrorRep(500, "Not acknowledged"))
	}
	c.JSON(200, base.SuccessRep("新增索引成功"))
}

func IndexExists(c *gin.Context) {
	index, _ := c.Params.Get("index")
	exists, err := utils.Client.IndexExists(index).Do(context.Background())
	if err != nil {
		c.JSON(500, base.ErrorRep(500, "createIndex is nil"))
		return
	}
	if !exists {
		c.JSON(500, base.SuccessRep("该索引不存在"))
		return
	}
	c.JSON(200, base.SuccessRep("该索引存在"))
}

func DeleteIndex(c *gin.Context) {
	index, _ := c.Params.Get("index")
	ctx := context.Background()
	deleteIndex, err := utils.Client.DeleteIndex(index).Do(ctx)
	if err != nil {
		c.JSON(500, base.ErrorRep(500, err.Error()))
		return
	}
	if deleteIndex == nil {
		c.JSON(500, base.ErrorRep(500, "deleteIndex is nil"))
		return
	}
	if !deleteIndex.Acknowledged {
		c.JSON(500, base.ErrorRep(500, "deleteIndex.Acknowledged is false"))
	}
	c.JSON(200, base.SuccessRep("索引已删除"))
}

type Product struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

func ScrollParallel(c *gin.Context) {
	index, _ := c.Params.Get("index")
	ctx := context.Background()
	total, err := utils.Client.Count(index).Do(ctx)
	if err != nil {
		c.JSON(500, base.ErrorRep(500, err.Error()))
		return
	}
	bar := pb.StartNew(int(total))

	// This example illustrates how to use goroutines to iterate
	// through a result set via ScrollService.
	//
	// It uses the excellent golang.org/x/sync/errgroup package to do so.
	//
	// The first goroutine will Scroll through the result set and send
	// individual documents to a channel.
	//
	// The second cluster of goroutines will receive documents from the channel and
	// deserialize them.
	//
	// Feel free to add a third goroutine to do something with the
	// deserialized results.
	//
	// Let's go.

	// 1st goroutine sends individual hits to channel.
	hits := make(chan json.RawMessage)
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		defer close(hits)
		// Initialize scroller. Just don't call Do yet.
		scroll := utils.Client.Scroll(index).Size(100)
		for {
			results, err := scroll.Do(ctx)
			if err == io.EOF {
				return nil // all results retrieved
			}
			if err != nil {
				return err // something went wrong
			}

			// Send the hits to the hits channel
			for _, hit := range results.Hits.Hits {
				select {
				case hits <- hit.Source:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	})

	// 2nd goroutine receives hits and deserializes them.
	//
	// If you want, setup a number of goroutines handling deserialization in parallel.
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for hit := range hits {
				// Deserialize
				var p Product
				err := json.Unmarshal(hit, &p)
				fmt.Println(p.Index)
				if err != nil {
					return err
				}

				// Do something with the product here, e.g. send it to another channel
				// for further processing.
				_ = p

				bar.Increment()

				// Terminate early?
				select {
				default:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Check whether any goroutines failed.
	if err := g.Wait(); err != nil {
		panic(err)
	}

	// Done.
	bar.FinishPrint("Done")

}

func Flush(c *gin.Context) {
	index, _ := c.Params.Get("index")
	ctx := context.Background()
	_, err := utils.Client.Flush().Index(index).Do(ctx)
	if err != nil {
		c.JSON(500, base.ErrorRep(500, err.Error()))
		return
	}
	c.JSON(200, base.SuccessRep("索引已刷新"))
}
