package main

import (
	"awesomeProject/lib/File"
	"time"
)

func main() {

	fp := "C:\\Users\\simple.chen\\Desktop\\陈新宇工作文档\\网段隔离扫描结果\\扫描2022\\ip.txt"

	ips, _ := File.ReadIPRangesFromFile(fp)

	File.WriteFile("output.txt", "段存活ip数量为：12")

	for _, ip := range ips {
		println(ip)
	}

	time.Sleep(time.Second)

}
