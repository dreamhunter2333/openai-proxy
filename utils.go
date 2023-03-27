package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/yalp/jsonpath"
)

func getTotalTokens(respBody []byte) (int, error) {
	var data interface{}
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		fmt.Println("getTotalTokens err: ", err)
		return 0, err
	}
	result, err := jsonpath.Read(data, "$.usage.total_tokens")
	if err != nil {
		fmt.Println("getTotalTokens jsonpath err: ", err)
		return 0, err
	}
	return int(result.(float64)), nil
}

func checkSelfApiKey(selfApiKey string, config ApiConfig) bool {
	limit, ok := config.ApiKeyToTokens[selfApiKey]
	if !ok {
		fmt.Println("no config for: ", selfApiKey)
		return false
	}
	used, err := getSelfApiKeyTokens(selfApiKey)
	if err != nil {
		fmt.Println("getSelfApiKeyTokens err: ", err)
		return false
	}
	return used < limit
}

func increaseSelfApiKeyTokens(respBody []byte, selfApiKey string) error {
	totalTokens, err := getTotalTokens(respBody)
	if err != nil {
		return err
	}
	redisHost := "localhost:6379"
	if redisHost, ok := os.LookupEnv("REDIS_HOST"); ok {
		fmt.Println("redisHost: " + redisHost)
	}

	// Connect to Redis server
	conn, err := redis.Dial("tcp", redisHost)
	if err != nil {
		return err
	}
	if redisPassword, ok := os.LookupEnv("REDIS_PASS"); ok {
		if _, err := conn.Do("AUTH", redisPassword); err != nil {
			conn.Close()
			return err
		}
	}

	defer conn.Close()

	_, err = conn.Do("INCRBY", selfApiKey, totalTokens)
	if err != nil {
		return err
	}

	return nil
}

func getSelfApiKeyTokens(selfApiKey string) (int, error) {
	redisHost := "localhost:6379"
	if redisHost, ok := os.LookupEnv("REDIS_HOST"); ok {
		fmt.Println("redisHost: " + redisHost)
	}

	// Connect to Redis server
	conn, err := redis.Dial("tcp", redisHost)
	if err != nil {
		return 0, err
	}
	if redisPassword, ok := os.LookupEnv("REDIS_PASS"); ok {
		if _, err := conn.Do("AUTH", redisPassword); err != nil {
			conn.Close()
			return 0, err
		}
	}

	defer conn.Close()

	reply, err := conn.Do("GET", selfApiKey)
	if err != nil {
		return 0, err
	}
	if reply == nil {
		return 0, nil
	}
	counter, err := redis.Int(reply, err)
	if err != nil {
		return 0, err
	}
	return counter, nil
}
