package main

import (
	"fmt"
	"github.com/mattn/go-jsonwire"
	"encoding/json"
)

func main() {
	jsonwire.StartServer("localhost:8910")
	defer jsonwire.StopServer()

	c := jsonwire.NewClient("localhost:8910")
	c.Post("/url", map[string]string {"url": "http://mixi.jp"})
	res, _ := c.Get("/title")
	var result map[string]interface{}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&result)
	fmt.Println(result["value"])
}
