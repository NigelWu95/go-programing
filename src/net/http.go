package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	//resp, err := client.Do(req)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	time.Sleep(900 * time.Millisecond)
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println(resp)
		fmt.Println(err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(body))
}
