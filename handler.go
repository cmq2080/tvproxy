package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func tsProxyHandler(c *gin.Context) {
	remoteURL := c.Query("url")
	resp, err := newHTTPClient().Get(remoteURL)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer resp.Body.Close()
	c.DataFromReader(200, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func fhzwHandler(c *gin.Context) {
	commonHandler("凤凰中文", c)
}

func fhzxHandler(c *gin.Context) {
	commonHandler("凤凰资讯", c)
}

func wxxwtHandler(c *gin.Context) {
	commonHandler("无线新闻台", c)
}

func msxwHandler(c *gin.Context) {
	commonHandler("民视新闻", c)
}

func tsxwHandler(c *gin.Context) {
	commonHandler("台视新闻", c)
}

func commonHandler(channelName string, c *gin.Context) {
	channel, ok := FindInChannels(channelName)
	if !ok {
		c.AbortWithStatus(500)
	}

	m3u8ProxyHandler(channel, c)
}

func m3u8ProxyHandler(channel Channel, c *gin.Context) {
	m3u8Url := channel.Url
	resp, err := newHTTPClient().Get(m3u8Url)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	bodyString := string(bodyBytes)
	processedBody := m3u8Proc(bodyString, baseURL+"play.ts?url=", channel.PlayPrefix)
	c.Data(200, resp.Header.Get("Content-Type"), []byte(processedBody))
}

func m3u8Proc(data string, prefixUrl string, playPrefix string) string {
	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") {
			sb.WriteString(l)
		} else {
			sb.WriteString(prefixUrl)
			l = playPrefix + l
			fmt.Println(l)
			l = url.QueryEscape(l)
			fmt.Println(l)
			sb.WriteString(l)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
