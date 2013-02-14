package jsonwire

import (
	"fmt"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"time"
	"syscall"
)

type Client struct {
	base string
}

var phantomjs *exec.Cmd

func StartServer(addr string, options ...string) {
	if phantomjs != nil {
		return
	}

	if addr == "" {
		addr = "localhost:8910"
	}

	go func() {
		path, err := exec.LookPath("phantomjs")
		if err != nil {
			panic("phantomjs is not installed")
		}
		phantomjs = exec.Command(path, append([]string{"--webdriver=" + addr}, options...)...)
		//phantomjs.Stdout = os.Stdout
		phantomjs.Stderr = os.Stderr
		sig := make(chan os.Signal, 10)
		go func() {
			sig <- syscall.SIGINT
			if phantomjs.Process != nil {
				phantomjs.Process.Kill()
			}
			phantomjs = nil
			err = nil
		}()
		err = phantomjs.Run()
		if err != nil && phantomjs != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			panic(err.Error())
		}
	}()
}

func StopServer() (err error) {
	if phantomjs != nil && phantomjs.Process != nil {
		phantomjs.Process.Signal(syscall.SIGINT)
	}
	return
}

func NewClient(addr string) *Client {
	if addr == "" {
		addr = "localhost:8910"
	}

	client := &http.Client{
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			return errors.New("no need to redirect")
		},
	}

	for i := 0; i < 3; i++ {
		res, err := client.Post("http://" + addr + "/session", "application/json", bytes.NewBufferString(`{"desiredCapabilities":{}}`))
		if err != nil && res == nil || res.StatusCode != 303 {
			time.Sleep(1e9)
			continue
		}
		base := res.Header.Get("Location")
		return &Client{base}
	}
	panic("error")
}

func (c *Client) Get(path string) (*http.Response, error) {
	return http.Get(c.base + path)
}

func (c *Client) Post(path string, arg interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(arg)
	return http.Post(c.base+path, "application/json", &buf)
}
