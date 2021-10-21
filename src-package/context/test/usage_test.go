/*
在goroutine中调用ctx.Done()时就会去初始化channel, ctx.Done()返回的是一个阻塞的只读通道，阻塞等待 closechann
当外层函数调用cannel函数时就会close channel, 此时通道为空, 通道可读, 读取到的值为通道类型对应的零值

*/

package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

// 创建临时压缩文件返回路径, 待上传oss后, 删除文件

func zipFileExecute(ctx context.Context) (string, error) {
	tmpDir, err := ioutil.TempDir("", "test-*")
	if err != nil {
		return tmpDir, fmt.Errorf("create tmp dir failed, error: %w", err)
	}

	fmt.Println("generate tmp dir: ", tmpDir)

	defer func() {
		go func() {
			// select {
			// case <-ctx.Done():
			// 	if err := os.RemoveAll(tmpDir); err != nil {
			// 		log.Fatalf("delete tmp dir %s failed", tmpDir)
			// 	}
			// }
			<-ctx.Done()
			if err := os.RemoveAll(tmpDir); err != nil {
				log.Fatalf("delete tmp dir %s failed", tmpDir)
			}

			fmt.Printf("delete tmp dir %s done\n", tmpDir)
		}()
	}()
	// do zip.....
	// return dir path
	return tmpDir, nil
}

func uploadOss(filePath string) error {
	return nil
}

func TestUsageWithCancel(t *testing.T) {
	subCtx, cancel := context.WithCancel(context.Background())

	tmpDir, err := zipFileExecute(subCtx)
	if err != nil {
		cancel()
		log.Fatal(err)
	}

	if err := uploadOss(tmpDir); err != nil {
		fmt.Println("upload failed")
	}
	// 上传成功后发送通知删除临时文件夹
	cancel()
	time.Sleep(10 * time.Second)
	fmt.Println("done")
}
