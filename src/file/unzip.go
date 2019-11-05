package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/**
@zipFile：压缩文件
@dest：解压之后文件保存路径
*/
func DeCompress(srcFile *os.File, dest string) error {
	zipFile, err := zip.OpenReader(srcFile.Name())
	if err != nil {
		fmt.Println("Unzip File Error：", err.Error())
		return err
	}
	defer zipFile.Close()
	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				fmt.Println("Unzip File Error : " + err.Error())
				return err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			fmt.Println("Unzip File Error : " + err.Error())
			continue
		}
		defer srcFile.Close()
		if strings.Contains(innerFile.Name, string(filepath.Separator)) {
			paths := strings.Split(innerFile.Name, string(filepath.Separator))
			err = os.MkdirAll(filepath.Join(dest, strings.Join(paths[0:len(paths) - 1], string(filepath.Separator))), os.ModePerm)
			if err != nil {
				fmt.Println("Unzip File Error : " + err.Error())
				return err
			}
		}
		newFile, err := os.Create(filepath.Join(dest, innerFile.Name))
		if err != nil {
			fmt.Println("Unzip File Error : " + err.Error())
			continue
		}
		io.Copy(newFile, srcFile)
		newFile.Close()
	}
	return nil
}

func main() {

	srcFile, err := os.Open("/Users/wubingheng/Downloads/200c387e89644e689aff2c06889be245.zip")
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	err = DeCompress(srcFile, "/Users/wubingheng/Downloads/200c387e89644e689aff2c06889be245")
	if err != nil {
		panic(err)
	}
}
