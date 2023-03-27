package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/yalp/jsonpath"

	"github.com/gin-gonic/gin"
)

var targetURL = "https://api.openai.com"
var authPrefix = "Bearer "

func getProxyUrl() *http.Transport {
	// Check if HTTP_PROXY environment variable is set
	if proxyUrl, ok := os.LookupEnv("http_proxy"); ok {
		fmt.Println("HTTP_PROXY: " + proxyUrl)
		proxyURL, _ := url.Parse(proxyUrl)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		return transport
	}
	return nil
}

func getApiKey() string {
	// use self API_KEY
	if apiKey, ok := os.LookupEnv("API_KEY"); ok {
		fmt.Println("API_KEY: ******")
		return apiKey
	}
	return ""
}

func proxy(c *gin.Context, transport *http.Transport) (*http.Response, error) {
	var client *http.Client

	remoteURL := targetURL + c.Param("path")
	req, _ := http.NewRequest(c.Request.Method, remoteURL, c.Request.Body)

	// 设置请求Header
	req.Header = c.Request.Header

	// new http Client
	if transport != nil {
		client = &http.Client{
			Transport: transport,
		}
	} else {
		client = &http.Client{}
	}

	return client.Do(req)
}

func getTotalTokens(respBody []byte) (int, error) {
	var data interface{}
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		return 0, err
	}
	result, err := jsonpath.Read(data, "$.usage.total_tokens")
	if err != nil {
		return 0, err
	}
	return int(result.(float64)), nil
}

func checkSelfApiKey(selfApiKey string) error {
	// TODO: check selfApiKey and balance
	return nil
}

func increaseSelfApiKeyTokens(respBody []byte, selfApiKey string) error {
	totalTokens, err := getTotalTokens(respBody)
	if err != nil {
		return err
	}
	// TODO: count totalTokens in db
	fmt.Println(selfApiKey, totalTokens)
	return nil
}

func handler(c *gin.Context, transport *http.Transport, apiKey string) {
	var selfApiKey string = ""
	// use self API_KEY
	if apiKey != "" {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, authPrefix) {
			selfApiKey = strings.TrimPrefix(authHeader, authPrefix)
			c.Request.Header.Set("Authorization", authPrefix+apiKey)
		}
	}
	resp, err := proxy(c, transport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	// 读取响应Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 计费
	if selfApiKey != "" {
		err := increaseSelfApiKeyTokens(respBody, selfApiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 设置响应Header
	for k, v := range resp.Header {
		c.Header(k, v[0])
	}
	// 设置响应状态码
	c.Status(resp.StatusCode)
	// 将响应体复制到原始响应中
	c.Writer.Write(respBody)
}

func main() {
	r := gin.Default()
	transport := getProxyUrl()
	apiKey := getApiKey()

	r.Any("/*path", func(c *gin.Context) {
		handler(c, transport, apiKey)
	})

	r.Run(":8080")
}
