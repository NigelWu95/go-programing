package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HttpGet struct {
	Url           string
	HttpClient    *http.Client
	MediaType     string
	MediaParams   map[string]string
	ContentLength int64
	DownloadBlock int64
	DownloadRange [][]int64
	Count         int
	FilePath      string // 包括路径和文件名
	TempFiles     []*os.File
	File          *os.File
	WG            sync.WaitGroup
}

func main() {

	get := new(HttpGet)
	get.FilePath = "./qsuits-7.73-jar-with-dependencies.jar"
	get.HttpClient = new(http.Client)
	get.Url = "https://search.maven.org/remotecontent?filepath=com/qiniu/qsuits/7.73/qsuits-7.73-jar-with-dependencies.jar"
	get.DownloadBlock = 1048576
	downloadStart := time.Now()

	req, err := http.NewRequest("GET", get.Url, nil)
	req.Header.Set("Range", "bytes=0-100")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	resp, err := get.HttpClient.Do(req)
	if err != nil {
		log.Panicf("Get %s error %v.\n", get.Url, err)
	}
	get.MediaType, get.MediaParams, _ = mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	contentRange := strings.Split(resp.Header.Get("Content-Range"), "/")
	if len(contentRange) < 2 {
		err = errors.New("can not get content-range")
		fmt.Println(err.Error())
		panic(err)
	}
	get.ContentLength, _ = strconv.ParseInt(contentRange[1], 10, 64)
	get.Count = int(math.Ceil(float64(get.ContentLength / get.DownloadBlock)))
	get.File, err = os.Create(get.FilePath)
	if err != nil {
		log.Panicf("Create file %s error %v.\n", get.FilePath, err)
	}
	var rangeStart int64 = 0
	for i := 0; i < get.Count; i++ {
		if i != get.Count - 1 {
			get.DownloadRange = append(get.DownloadRange, []int64{rangeStart, rangeStart + get.DownloadBlock - 1})
		} else {
			// 最后一块
			get.DownloadRange = append(get.DownloadRange, []int64{rangeStart, get.ContentLength - 1})
		}
		rangeStart += get.DownloadBlock
	}
	// Check if the download has paused.
	for i := 0; i < len(get.DownloadRange); i++ {
		rangeI := fmt.Sprintf("%d-%d", get.DownloadRange[i][0], get.DownloadRange[i][1])
		tempFile, err := os.OpenFile(get.FilePath + "." + rangeI, os.O_RDONLY|os.O_APPEND, 0)
		if err != nil {
			tempFile, _ = os.Create(get.FilePath + "." + rangeI)
		} else {
			fi, err := tempFile.Stat()
			if err == nil {
				get.DownloadRange[i][0] += fi.Size()
			}
		}
		get.TempFiles = append(get.TempFiles, tempFile)
	}

	for i, _ := range get.DownloadRange {
		get.WG.Add(1)
		go get.Download(i)
	}

	get.WG.Wait()

	for i := 0; i < len(get.TempFiles); i++ {
		tempFile, _ := os.Open(get.TempFiles[i].Name())
		cnt, err := io.Copy(get.File, tempFile)
		if cnt <= 0 || err != nil {
			log.Printf("Download #%d error %v.\n", i, err)
		}
		tempFile.Close()
	}
	get.File.Close()
	log.Printf("Download complete and store file %s with %v.\n", get.FilePath, time.Now().Sub(downloadStart))
	defer func() {
		for i := 0; i < len(get.TempFiles); i++ {
			err := os.Remove(get.TempFiles[i].Name())
			if err != nil {
				log.Printf("Remove temp file %s error %v.\n", get.TempFiles[i].Name(), err)
			} else {
				log.Printf("Remove temp file %s.\n", get.TempFiles[i].Name())
			}
		}
	}()
}

func (get *HttpGet) Download(i int) {
	defer get.WG.Done()
	if get.DownloadRange[i][0] > get.DownloadRange[i][1] {
		return
	}
	rangeI := fmt.Sprintf("%d-%d", get.DownloadRange[i][0], get.DownloadRange[i][1])
	log.Printf("Download #%d bytes %s.\n", i, rangeI)

	defer get.TempFiles[i].Close()

	req, err := http.NewRequest("GET", get.Url, nil)
	req.Header.Set("Range", "bytes=" + rangeI)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	resp, err := get.HttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Panicf("Download #%d error %v.\n", i, err)
	} else {
		cnt, err := io.Copy(get.TempFiles[i], resp.Body)
		if cnt == int64(get.DownloadRange[i][1] - get.DownloadRange[i][0] + 1) {
			log.Printf("Download #%d complete.\n", i)
		} else {
			reqDump, _ := httputil.DumpRequest(req, false)
			respDump, _ := httputil.DumpResponse(resp, true)
			log.Panicf("Download error %d %v, expect %d-%d, but got %d.\nRequest: %s\nResponse: %s\n", resp.StatusCode, err, get.DownloadRange[i][0], get.DownloadRange[i][1], cnt, string(reqDump), string(respDump))
		}
	}
}
