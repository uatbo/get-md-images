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
	matches := regexpObject.FindAllString(text, -1)
	// 匹配图像标签和图像下标
	//matches_index := regexpObject.FindAllStringIndex(text, -1)

	// 查找图像标签中的url链接
	pattern = `\(http(.*?)\)`
	regexpObject = regexp.MustCompile(pattern)

	for i := 0; i < len(matches); i++ {
		imageUrl := regexpObject.FindAllString(matches[i], -1)
		if len(imageUrl) == 0 {
			fmt.Println("")
			continue
		} else {
			fmt.Println(imageUrl[0][1 : len(imageUrl[0])-1])
		}
	}
}

// 替换字符串中的链接

// 下载图片存储到当前目录下images目录
