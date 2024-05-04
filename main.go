package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CODE_SUCCESS = 200
)

var (
	currentDir string = "."
	baseURL    string

	listen string
)

func main() {
	initialize()
	fmt.Println("Proxy Started!!!\nListen:" + listen)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("play.ts", tsProxyHandler)

	r.GET("/fhzw.m3u8", fhzwHandler)
	r.GET("/fhzx.m3u8", fhzxHandler)
	r.GET("/wxxwt.m3u8", wxxwtHandler)
	r.GET("/msxw.m3u8", msxwHandler)
	r.GET("/tsxw.m3u8", tsxwHandler)

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

	buff2, ok := config["listen"].(string)
	if !ok {
		panic("Property Error: " + "listen")
	}
	listen = buff2
	if strings.Index(listen, ":") == 0 {
		listen = "127.0.0.1" + listen
	}
	baseURL = "http://" + listen + "/"

	buff3, ok := config["channels"].([]interface{})
	if !ok {
		panic("Property Error: " + "channels")
	}

	for _, v := range buff3 {
		buff3_1 := v.(map[string]interface{})
		channel := Channel{
			Name:       buff3_1["name"].(string),
			Url:        buff3_1["url"].(string),
			PlayPrefix: buff3_1["play_prefix"].(string),
		}

		channels = append(channels, channel)
	}
}

type Response struct {
	ctx *gin.Context
}

func newResponse(c *gin.Context) *Response {
	return &Response{ctx: c}
}

func (resp *Response) Success(msg string, data interface{}) {
	resp.ctx.JSON(CODE_SUCCESS, map[string]interface{}{
		"stat": 0,
		"msg":  msg,
		"data": data,
	})
}

func (resp *Response) SuccessWithoutData(msg string) {
	resp.ctx.JSON(CODE_SUCCESS, map[string]interface{}{
		"stat": 0,
		"msg":  msg,
	})
}

func (resp *Response) Error(msg string) {
	resp.ctx.JSON(CODE_SUCCESS, map[string]interface{}{
		"stat": 1,
		"msg":  msg,
	})
}

func (resp *Response) ErrorWithData(msg string, data interface{}) {
	resp.ctx.JSON(CODE_SUCCESS, map[string]interface{}{
		"stat": 1,
		"msg":  msg,
		"data": data,
	})
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
