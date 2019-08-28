package main

import (
	"fmt"
	"os"
)

func main(){

	fileName := "/Users/wubingheng/Downloads/8.jpeg"
	desFileName := "/Users/wubingheng/Downloads/8-bak.jpeg"
	file, err := os.Open(fileName)
	if err != nil{
		fmt.Println(nil)
	}

	info, _ := os.Stat(fileName)
	size := info.Size()
	var count int64 = 1
	//这里切分原意为通过协程来分段读取
	if size % 2 == 0 {
		count *= 2
	} else if size % 3 == 0 {
		count *= 3
	} else {
		count *= 1
	}
	si := size / count
	fmt.Printf("文件总大小：%v, 分片数：%v, 每个分片大小：%v\n", size, count, si)

	desFile, err := os.OpenFile(desFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i <= int(count); i++ {
		//申明一个byte
		b := make([]byte, si)
		//从哪个位置开始读
		file.Seek(int64(i) * si, 0)
		//读到byte数组里边
		file.Read(b)
		//从哪个位置开始写
		desFile.Seek(int64(i) * si, 0)
		//写入
		desFile.Write(b)
	}

	defer desFile.Close()
	defer file.Close()
}
