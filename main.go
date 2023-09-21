package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/robfig/cron"
)

var (
	beeper sync.WaitGroup
	db     *pg.DB
)

func main() {
	r := gin.Default()

	err := getError()
	if err != nil {
		fmt.Println("Okay Error")
	}
	db, err = new("localhost:4848", "postgres", "pass123", "xrp_db")
	if err != nil {
		fmt.Println("Connection Error")
	}

	// XRPLedgerService(db)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})
	r.GET("/start", func(c *gin.Context) {

		cron := cron.New()
		cron.AddFunc("@every 2s", XRPLedgerService)
		cron.Run()
		cron.Stop()
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	//User Auth
	r.POST("/register", RegisterUserHandler)
	r.POST("/signin", SignInHandler)
	r.POST("/create-stream", CreateStreamHandler)

	//update
	//pass streamID and Address in body and token in headers
	r.POST("/add-address-to-stream", registerAddressToStreamHandler)

	//pass streamID  and url in body and token
	r.POST("/update-stream-url", updateStreamURLHandler)

	//Read operations
	//
	r.POST("/view-addresses", viewAddressesHandler)
	r.GET("/view-stream", viewStreamHandler)

	//Delete operation
	r.DELETE("/delete-address", deleteAddressHandler)

	r.Run()

}

func getError() error {
	return nil
}
