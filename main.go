package main

import (
	"bytes"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 指定文本文件路径
var filePath string

func main() {
	// 读取文件路径
	//filePath = "/Users/uatbo/Files/大兵/笔记/增量表与全量表的区别.md"

	args := os.Args
	filePath = args[1]

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

	filePath = filepath.Dir(filePath) + "/tmp.md"
	// 写入文件
	WriteStringToFile(*newText, filePath)
}

const charset = "abcdefghijklmnopqrstuvwxyz"
const prefix = "images/"

// 生成长度为n的随机字符串
func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func getTime() string {
	currentTime := time.Now()
	// 自定义时间格式化字符串
	timeFormat := "200601021504"
	// 使用时间格式化字符串格式化时间
	formattedTime := currentTime.Format(timeFormat)
	return formattedTime
}

// DownloadImage 下载图片存储到当前目录下images目录，返回相对路径
func DownloadImage(url string) string {
	// 请求图片，获取http报文
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()

	// 获取图片类型
	res, _ := io.ReadAll(response.Body) // 将body内容转为字节串
	mt := mimetype.Detect(res)
	if err != nil {
		return ""
	}
	// 字节串写回body，后续写入图片文件需要
	response.Body = io.NopCloser(bytes.NewReader(res))

	// 定义图片存储路径
	filename := prefix + randomString(3) + getTime() + mt.Extension()
	filedir := filepath.Dir(filePath) + "/" + filename
	fmt.Println(filedir)
	dirPath := filepath.Dir(filedir) // 提取目录路径

	// 递归创建目录
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("创建目录时出错:", err)
		return ""
	}

	// 写入文件
	file, err := os.Create(filedir)
	if err != nil {
		fmt.Println("创建文件时出错:", err)
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

	newText := (*str)[0:subStringIndex[0]] + subString + (*str)[subStringIndex[1]:len(*str)]
	return &newText
}

// WriteStringToFile 将字符串写入文本文件函数
func WriteStringToFile(text string, path string) {
	// 打开文件以进行写入，如果文件不存在会创建文件
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 将字符串写入文件
	_, err = file.WriteString(text)
	if err != nil {
		panic(err)
	}
}
