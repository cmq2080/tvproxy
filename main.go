package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CODE_SUCCESS = 200
)

var (
	debug      bool   = false
	currentDir string = "."
	baseURL    string

	listen  string
	timeout int = 10
)

func main() {
	initialize()
	fmt.Println("Proxy Started!!!\nListen: http://" + listen + "/\nTimeout: " + strconv.Itoa(timeout) + " Second(s)\nBaseURL: " + baseURL)

	if debug {
		fmt.Println("[channels]")
		for _, channel := range channels {
			fmt.Println("名称: " + channel.Name)
			fmt.Println("描述: " + channel.Desc)
			fmt.Println("链接: " + channel.Url)
			fmt.Println("前缀: " + channel.PlayPrefix + "\n")
		}
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Handle(http.MethodGet, "/*action", func(c *gin.Context) {
		action := c.Param("action")
		action = strings.TrimLeft(action, "/")

		if action == "play.ts" {
			tsProxyHandler(c)
		}
		if action == "play.m3u8" {
			m3u8ProxyHandler(c)
		}

		length := len(action)
		if strings.HasSuffix(action, ".m3u8") {
			channelName := action[:length-5]
			if debug {
				fmt.Println(channelName)
			}

			channel, ok := FindInChannels(channelName)
			if !ok {
				c.AbortWithStatus(500)
			}
			url := channel.Url
			c = SetQuery(c, "url", url)

			m3u8ProxyHandler(c)
		}

	})

	r.Run(listen)
}

func initialize() {
	thisExe, _ := os.Executable()
	currentDir = filepath.Dir(thisExe)

	configFilepath := currentDir + "/config.json"
	buff, err := os.ReadFile(configFilepath)
	if err != nil {
		panic(err)
	}

	var config map[string]interface{}

	json.Unmarshal(buff, &config)

	buff_debug, ok := config["debug"].(bool)
	if !ok {
		debug = false
	} else {
		debug = buff_debug
	}

	buff_listen, ok := config["listen"].(string)
	if !ok {
		panic("Property Error: " + "listen")
	}
	if !strings.Contains(buff_listen, ":") {
		buff_listen += ":80"
	}
	if strings.HasSuffix(buff_listen, ":") {
		buff_listen += "80"
	}

	buff2_array := strings.Split(buff_listen, ":")
	if len(buff2_array) != 2 {
		panic("Property Error: " + "listen" + "-Error Format")
	}
	listen = buff_listen

	serverUrl := buff2_array[0]
	serverPort := buff2_array[1]

	if serverUrl == "" {
		serverUrl = "0.0.0.0"
	}
	if serverPort == "" {
		serverPort = "80"
	}

	baseURL = "http://" + serverUrl + ":" + serverPort + "/"

	buff_timeout, ok := config["timeout"].(float64)
	if !ok {
		panic("Property Error: " + "timeout")
	}
	timeout = int(buff_timeout)

	buff_channels, ok := config["channels"].([]interface{})
	if !ok {
		panic("Property Error: " + "channels")
	}

	for _, v := range buff_channels {
		buff_channel := v.(map[string]interface{})
		channel := Channel{
			Name: buff_channel["name"].(string),
			Desc: buff_channel["desc"].(string),
			Url:  buff_channel["url"].(string),
		}
		channel.PlayPrefix, ok = buff_channel["play_prefix"].(string)
		if !ok || channel.PlayPrefix == "" {
			channel.PlayPrefix = GetPlayPrefix(channel.Url)
		}

		channels = append(channels, channel)
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func SetQuery(c *gin.Context, key, value string) *gin.Context {
	params, _ := url.ParseQuery(c.Request.URL.RawQuery)
	params.Set(key, value)
	c.Request.URL.RawQuery = params.Encode()

	return c
}
