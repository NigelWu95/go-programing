package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	//DefaultDownloadBlock int64 = 4194304
	DefaultDownloadBlock int64 = 1048576
)

type GoGet struct {
	Url           string
	Cnt           int
	DownloadBlock int64
	CustomCnt     int
	Latch         int
	Header        http.Header
	MediaType     string
	MediaParams   map[string]string
	FilePath      string // 包括路径和文件名
	GetClient     *http.Client
	ContentLength int64
	DownloadRange [][]int64
	File          *os.File
	TempFiles     []*os.File
	WG            sync.WaitGroup
}

func NewGoGet() *GoGet {
	get := new(GoGet)
	get.FilePath = "./"
	get.GetClient = new(http.Client)

	flag.Parse()
	get.Url = *urlFlag
	get.DownloadBlock = DefaultDownloadBlock

	return get
}

var urlFlag = flag.String("u", "https://search.maven.org/remotecontent?filepath=com/qiniu/qsuits/7.73/qsuits-7.73-jar-with-dependencies.jar", "Fetch file url")

// var cntFlag = flag.Int("c", 1, "Fetch concurrently counts")

func main() {
	get := NewGoGet()

	downloadStart := time.Now()

	req, err := http.NewRequest("HEAD", get.Url, nil)
	resp, err := get.GetClient.Do(req)
	get.Header = resp.Header
	if err != nil {
		log.Panicf("Get %s error %v.\n", get.Url, err)
	}
	get.MediaType, get.MediaParams, _ = mime.ParseMediaType(get.Header.Get("Content-Disposition"))
	get.ContentLength = resp.ContentLength
	get.Cnt = int(math.Ceil(float64(get.ContentLength / get.DownloadBlock)))
	if strings.HasSuffix(get.FilePath, "/") {
		get.FilePath += get.MediaParams["filename"]
	}
	get.File, err = os.Create(get.FilePath)
	if err != nil {
		log.Panicf("Create file %s error %v.\n", get.FilePath, err)
	}
	log.Printf("Get %s MediaType:%s, Filename:%s, Size %d.\n", get.Url, get.MediaType, get.MediaParams["filename"], get.ContentLength)
	if get.Header.Get("Accept-Ranges") != "" {
		log.Printf("Server %s support Range by %s.\n", get.Header.Get("Server"), get.Header.Get("Accept-Ranges"))
	} else {
		log.Printf("Server %s doesn't support Range.\n", get.Header.Get("Server"))
	}

	log.Printf("Start to download %s with %d thread.\n", get.MediaParams["filename"], get.Cnt)
	var rangeStart int64 = 0
	for i := 0; i < get.Cnt; i++ {
		if i != get.Cnt - 1 {
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

	go get.Watch()
	get.Latch = get.Cnt
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

func (get *GoGet) Download(i int) {
	defer get.WG.Done()
	if get.DownloadRange[i][0] > get.DownloadRange[i][1] {
		return
	}
	rangeI := fmt.Sprintf("%d-%d", get.DownloadRange[i][0], get.DownloadRange[i][1])
	log.Printf("Download #%d bytes %s.\n", i, rangeI)

	defer get.TempFiles[i].Close()

	req, err := http.NewRequest("GET", get.Url, nil)
	req.Header.Set("Range", "bytes=" + rangeI)
	resp, err := get.GetClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Download #%d error %v.\n", i, err)
	} else {
		cnt, err := io.Copy(get.TempFiles[i], resp.Body)
		if cnt == int64(get.DownloadRange[i][1] - get.DownloadRange[i][0]+1) {
			log.Printf("Download #%d complete.\n", i)
		} else {
			reqDump, _ := httputil.DumpRequest(req, false)
			respDump, _ := httputil.DumpResponse(resp, true)
			log.Panicf("Download error %d %v, expect %d-%d, but got %d.\nRequest: %s\nResponse: %s\n", resp.StatusCode, err, get.DownloadRange[i][0], get.DownloadRange[i][1], cnt, string(reqDump), string(respDump))
		}
	}
}

func (get *GoGet) Watch() {
	fmt.Printf("[=================>]\n")
}

