package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func insertLog(c *gin.Context) {
	var (
		headers string
		body    []byte
	)
	for k, v := range c.Request.Header {
		headerV := ""
		for _, s := range v {
			headerV += s + "; "
		}
		headers += k + ": " + headerV + "\n"
	}
	body, _ = ioutil.ReadAll(c.Request.Body)
	fmt.Println(body)
	/*	mongodb.TestLog(mongodb.Logs{
		Host:       c.Request.Host,
		RequestIP:  c.ClientIP(),
		Time:       time.Now().Unix(),
		RequestUri: c.Request.RequestURI,
		State:      200,
		Header:     headers,
		Body:       string(body),
	})*/
}
