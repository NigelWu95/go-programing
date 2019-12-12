package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {

	file, err := os.Open("/Users/wubingheng/Downloads/redis-work.md")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n') // 以'\n'为结束符读入一行
		if io.EOF == err {
			//bytes, isPrefix, _ := r.ReadLine()
			//fmt.Println(isPrefix)
			fmt.Println(line)
			break
		} else if err != nil{
			panic(err)
		} else {
			fmt.Println(line)
		}
	}
}
