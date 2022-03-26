package main

import (
	"fmt"
	"os"
	"bufio"
)
// go语言圣经练习题: https://cloud.tencent.com/developer/article/1501697

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		counts[input.Text()]++
	}
	for k, v := range counts {
		fmt.Printf("%s == %d\n", k, v)
	}
}