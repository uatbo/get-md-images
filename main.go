package main

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	matches_index := regexpObject.FindAllStringIndex(text, -1)

	// 查找图像标签中的url链接
	pattern = `\(http(.*?)\)`
	regexpObject = regexp.MustCompile(pattern)

	var newText *string = &text
	// 下载修改文档中的图像标签，从最后一个标签开始更改
	for i := len(matches) - 1; i >= 0; i-- {
		imageUrl := regexpObject.FindAllString(matches[i], -1)
		if len(imageUrl) == 0 { // 当前图像标签里面没有网络图像
			fmt.Println("")
			continue
		} else { // 当前图像标签是网络图像
			// 下载图片存储到当前目录下images目录，返回相对路径
			tmpPath := DownloadImage(imageUrl[0][1 : len(imageUrl[0])-1])
			if tmpPath == "" {
				panic("Downloading image failed")
			}
			// 生成新的子串
			tmpSubStr := NewSubString(matches[i], tmpPath)
			if tmpSubStr == "" {
				panic("New sub string error")
			}
			// 修改字符串
			newText = AlterString(newText, tmpSubStr, matches_index[i])
			if newText == nil {
				panic("Alter string error")
			}
		}
	}

	// 写入文件
	fmt.Println(*newText)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const prefix = "images/"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// DownloadImage 下载图片存储到当前目录下images目录，返回相对路径
func DownloadImage(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	mt, err := mimetype.DetectReader(response.Body)
	if err != nil {
		return ""
	}
	filename := prefix + randomString(5) + mt.String()

	file, err := os.Create(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return ""
	}

	return filename
}

// NewSubString 生成新的子串
func NewSubString(subString string, url string) string {

	index := strings.Index(subString, "(")
	if index == -1 {
		return ""
	}
	return subString[0:index+1] + url + ")"
}

// AlterString 替换字符串中的链接
func AlterString(str *string, subString string, subStringIndex []int) *string {

	newText := (*str)[0:subStringIndex[0]] + subString + (*str)[subStringIndex[1]:len(*str)-1]
	return &newText
}
