package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	url := "https://github.com/NigelWu95/qiniu-suits-java/releases/download/v7.72/qsuits-7.72.jar"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Range", "bytes=0-100")
	if err != nil {
		log.Println(err.Error())
	}
	resp, err := client.Do(req)
	contentLength := int(resp.ContentLength)
	contentRange := resp.Header.Get("Content-Range")
	contentRange = strings.Split(contentRange, "/")[1]
	fmt.Println(contentRange)
	mediaType, mediaParams, _ := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	log.Printf("Get %s MediaType:%s, Filename:%s, Length %d.\n", url, mediaType, mediaParams["filename"], contentLength)
}
