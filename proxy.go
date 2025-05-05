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

func m3u8ProxyHandler(c *gin.Context) {
	var m3u8Url string = ""
	url, ok := c.GetQuery("url")
	if ok {
		m3u8Url = url
	} else {
		panic("Param Not Found: url")
	}
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
	processedBody := m3u8Proc(bodyString, baseURL+"play.ts?url=", m3u8Url)
	c.Data(200, resp.Header.Get("Content-Type"), []byte(processedBody))
}

func m3u8Proc(data string, prefixUrl string, m3u8Url string) string {
	channel, ok := FindInChannels2("url", m3u8Url)

	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "#") { // Property
			sb.WriteString(l)
		} else { // Value
			sb.WriteString(prefixUrl)

			if !strings.HasPrefix(l, "http") { // Need Play Prefix
				playPrefix := GetPlayPrefix(m3u8Url)
				if ok {
					playPrefix = channel.PlayPrefix
				}
				if !strings.HasSuffix(playPrefix, "/") {
					playPrefix += "/"
				}

				l = playPrefix + l
			}
			if debug {
				fmt.Println(l)
			}
			l = url.PathEscape(l)
			if debug {
				fmt.Println(l)
			}
			sb.WriteString(l)
		}
		sb.WriteString("\n")
	}
	if debug {
		fmt.Println("<<<<<<<<<<<<Response: ")
		fmt.Println(sb.String())
		fmt.Println(">>>>>>>>>>>>\n")
	}

	return sb.String()
}

func GetPlayPrefix(url string) string {
	prefix := ""
	i := strings.LastIndex(url[:len(url)-1], "/")
	if i > 0 {
		prefix = url[:i] + "/"
	}

	return prefix
}
