package main

import (
	"fmt"
	"github.com/inconshreveable/go-update"
	"net/http"
)

func main() {

	url := "https://github.com/NigelWu95/qsuits-exec-go/raw/master/bin/qsuits_darwin_amd64"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		fmt.Println(err.Error())
	}
}

