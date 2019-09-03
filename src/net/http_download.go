package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func httpClientDo(resultDir string, version string, req *http.Request) (qsuitsFilePath string, err error) {

	var jarFile string
	err = os.MkdirAll(filepath.Join(resultDir, ".qsuits"), os.ModePerm)
	if err != nil {
		return jarFile, err
	}
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return jarFile, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jarFile, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		jarFile = filepath.Join(resultDir, ".qsuits", "qsuits-" + version + ".jar")
		err = ioutil.WriteFile(jarFile, body, 0755)
		if err != nil {
			return jarFile, err
		}
		return jarFile, nil
	} else {
		return jarFile, errors.New(resp.Status)
	}
}

func main() {

	resultDir := ".."
	version := "7.72"
	req, err := http.NewRequest("GET", "https://github.com/NigelWu95/qiniu-suits-java/releases/download/v" +
		version + "/qsuits-" + version + ".jar", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	qsuitsFilePath, err := httpClientDo(resultDir, version, req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(qsuitsFilePath)
}