package jsonwire_test

import (
	"github.com/mattn/go-jsonwire"
	"time"
	"errors"
	"encoding/json"
	"testing"
)

func TestStartStopServer(t *testing.T) {
	jsonwire.StartServer("localhost:9999")
	time.Sleep(1e9)
	jsonwire.StopServer()
}

func newClient(addr string) (c *jsonwire.Client, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = errors.New("Failed to create client")
		}
	}()
	return jsonwire.NewClient(addr), nil
}

func TestJsonWire(t *testing.T) {
	jsonwire.StartServer("localhost:9998")
	defer jsonwire.StopServer()

	c, err := newClient("localhost:9998")
	if err != nil {
		t.Fatal(`Failed to NewClient:`, err)
	}

	_, err = c.Post("/url", map[string]string {"url": "http://www.google.com/"})
	if err != nil {
		t.Fatal(`Failed to Post /url:`, err)
	}
	c.Get("/title")

	res, err := c.Get("/title")
	if err != nil {
		t.Fatal(`Failed to Get /title:`, err)
	}
	var result map[string]interface{}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		t.Fatal(`Failed to decode JSON:`, err)
	}
	title, ok := result["value"]
	if !ok {
		t.Fatal(`Failed to get title of http://www.google.com/:`, err)
	}
	if title != "Google" {
		t.Fatalf("Expected %s, but %s:", "Google", title)
	}
}
