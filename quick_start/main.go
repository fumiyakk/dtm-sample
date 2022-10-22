package main

/*
Before start this demo, You have to run dtm server.
``` go
git clone https://github.com/dtm-labs/dtm && cd dtm
go run main.go
```
*/

import (
	"fmt"
	"log"
	"time"

	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid/v3"
)

// business address
const businessAPI = "/api/business_start"
const businessPort = 8082

var business = fmt.Sprintf("http://localhost:%d%s", businessPort, businessAPI)

// main starts demo
func main() {
	StartSvr()

	sr := SuccessRequest()
	log.Println(sr)

	fr := FailRequest()
	log.Println(fr)

	select {}
}

// StartSvr quick start: start server
func StartSvr() {
	app := gin.New()
	qsAddRoute(app)
	log.Printf("quick start examples listening at %d", businessPort)
	go func() {
		_ = app.Run(fmt.Sprintf(":%d", businessPort))
	}()
	time.Sleep(100 * time.Millisecond)
}

func qsAddRoute(app *gin.Engine) {
	app.POST(businessAPI+"/TransIn", func(c *gin.Context) {
		log.Printf("TransIn")
		c.JSON(200, "")
	})
	app.POST(businessAPI+"/TransInFail", func(c *gin.Context) {
		log.Printf("TransIn")
		c.JSON(409, "")
	})
	app.POST(businessAPI+"/TransInCompensate", func(c *gin.Context) {
		log.Printf("TransInCompensate")
		c.JSON(200, "")
	})
	app.POST(businessAPI+"/TransOut", func(c *gin.Context) {
		log.Printf("TransOut")
		c.JSON(200, "")
	})
	app.POST(businessAPI+"/TransOutCompensate", func(c *gin.Context) {
		log.Printf("TransOutCompensate")
		c.JSON(200, "")
	})
}

const dtmServer = "http://localhost:36789/api/dtmsvr"

func SuccessRequest() string {
	req := &gin.H{"amount": 30}

	saga := dtmcli.NewSaga(dtmServer, shortuuid.New()).
		Add(business+"/TransOut", business+"/TransOutCompensate", req).
		Add(business+"/TransIn", business+"/TransInCompensate", req)
	err := saga.Submit()

	if err != nil {
		panic(err)
	}

	return saga.Gid
}

func FailRequest() string {
	req := &gin.H{"amount": 50}

	saga := dtmcli.NewSaga(dtmServer, shortuuid.New()).
		Add(business+"/TransOut", business+"/TransOutCompensate", req).
		Add(business+"/TransInFail", business+"/TransInCompensate", req)
	err := saga.Submit()

	if err != nil {
		panic(err)
	}

	return saga.Gid
}
