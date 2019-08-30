package main

import (
	"fmt"
	"os"
	"time"
)

func main(){

	fileName := "/Users/wubingheng/Downloads/8.jpeg"
	desFileName := "/Users/wubingheng/Downloads/8-bak.jpeg"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(nil)
	}

	info, _ := os.Stat(fileName)
	size := info.Size()
	var count int64 = 1
	if size % 2 == 0 {
		count *= 2
	} else if size % 3 == 0 {
		count *= 3
	} else{
		count *= 1
	}

	si := size / count
	fmt.Printf("文件总大小：%v, 分片数：%v, 每个分片大小：%v\n", size, count, si)

	desF,err := os.OpenFile(desFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < int(count); i++ {
		go func(vs int) {
			//申明一个byte
			b := make([]byte, si)
			//从指定位置开始读
			file.ReadAt(b, int64(vs) * si)
			//从指定位置开始写
			desF.WriteAt(b, int64(vs) * si)

		}(i)
	}
	time.Sleep(time.Second*5)
	defer desF.Close()
	defer file.Close()
}
