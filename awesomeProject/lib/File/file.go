package File

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 从文件中读取IP段信息并返回 IP 段列表
func ReadIPRangesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		println("error")
		return nil, err
	}
	defer file.Close()

	ipRanges := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//println(line)
		// 忽略空行或注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// 假设每行只有一个IP段，可以根据实际需求进行适当修改
		ipRanges = append(ipRanges, line)
	}

	if err := scanner.Err(); err != nil {

		return nil, err
	}
	return ipRanges, nil
}

func CreateFile(filename string) error {
	_, err := os.Create(filename)
	if err != nil {
		println("创建文件失败~ 请检查文件路径")
		return err
	}
	return nil
}

func WriteFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		//CreateFile(filename)
		fmt.Println("打开文件失败:", err)
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return err
	}

	fmt.Println("追加写入文件成功")
	return nil
}
