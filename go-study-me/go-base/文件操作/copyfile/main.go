package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

// CopyFile dt -> src
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Printf("Open %s failed, err:%v.\n", srcName, err)
		return
	}

	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Open %s failed, err:%v.\n", srcName, err)
		return
	}
	defer dst.Close()
	return io.Copy(dst, src) //
}

func main() {
	wd, _ := os.Getwd()
	src := path.Join(wd, "./dst.txt")
	dt := path.Join(wd, "./src.txt")
	_, err := CopyFile(dt, src)
	if err != nil {
		fmt.Println("copy file failed, err: ", err)
		return
	}
	fmt.Println("copy done!")

}
