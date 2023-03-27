package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Any("/*path", func(c *gin.Context) {
		remoteURL := "https://api.openai.com" + c.Param("path")
		req, _ := http.NewRequest(c.Request.Method, remoteURL, c.Request.Body)
		for k, v := range c.Request.Header {
			req.Header.Set(k, v[0])
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			c.Header(k, v[0])
		}
		c.Status(resp.StatusCode)
		c.Writer.Write([]byte{})
	})

	r.Run(":8080")
}
