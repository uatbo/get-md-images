package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	// 指定文本文件路径
	filePath := "/Users/uatbo/Files/大兵/笔记/增量表与全量表的区别.md"
	// 读取文本文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		os.Exit(1)
	}
	// 将文件内容转换为字符串
	text := string(content)
	// 定义正则表达式模式
	pattern := `!\[(.*?)\]\((.*?)\)`
	// 编译正则表达式
	regexpObject := regexp.MustCompile(pattern)
	// 查找匹配
	// 匹配图像标签
	matches := regexpObject.FindAllString(text, -1)
	//匹配图像标签的下标

	// 安装匹配中的内容查找并下载图片
	for i := 0; i < len(matches); i++ {
		fmt.Println(matches[i])
	}

}
