package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

/**
@files：需要压缩的文件
@compreFile：压缩之后的文件
*/
func main() {
	compreFile, err := os.Create("/Users/wubingheng/Downloads/200c387e89644e689aff2c06889be245-1.zip")
	if err != nil {
		panic(err)
	}
	zw := zip.NewWriter(compreFile)
	defer zw.Close()
	var files []*os.File
	err = filepath.Walk("/Users/wubingheng/Downloads/200c387e89644e689aff2c06889be245", func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		err := compress_zip(file, zw)
		if err != nil {
			panic(err)
		}
		file.Close()
	}
}

/**
功能：压缩文件
@file:压缩文件
@prefix：压缩文件内部的路径
@tw：写入压缩文件的流
*/
func compress_zip(file *os.File, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		fmt.Println("压缩文件失败：", err.Error())
		return err
	}
	// 获取压缩头信息
	head, err := zip.FileInfoHeader(info)
	if err != nil {
		fmt.Println("压缩文件失败：", err.Error())
		return err
	}
	// 指定文件压缩方式 默认为 Store 方式 该方式不压缩文件 只是转换为zip保存
	head.Method = zip.Deflate
	fw, err := zw.CreateHeader(head)
	if err != nil {
		fmt.Println("压缩文件失败：", err.Error())
		return err
	}
	// 写入文件到压缩包中
	_, err = io.Copy(fw, file)
	file.Close()
	if err != nil {
		fmt.Println("压缩文件失败：", err.Error())
		return err
	}
	return nil
}