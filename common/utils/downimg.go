package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/hkyangyi/newe/common/base"
	"github.com/hkyangyi/newe/common/file"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func DownImg(url string) (string, error) {
	//url := "https://example.com/image.jpg" // 替换为要下载的图片的URL

	// response, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("图片下载失败:", err)
	// 	return "", err
	// }
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return "", err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 发送请求并获取响应
	client := &http.Client{
		Transport: tr,
	}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("图片下载失败:")
		fmt.Println("Request error:", err)
		return "", err
	}
	defer response.Body.Close()

	//生产图片保存地址
	now := time.Now().Format("20060102")
	path := base.Conf.IMG_SavePath + "/hikimg" + "/" + now + "/"
	filename := GetUUID() + ".jpg"
	imgurl := base.Conf.IMG_PrefixUrl + "/hikimg" + "/" + now + "/" + filename

	dir, err := os.Getwd()
	if err != nil {
		return "", err

	}
	err = file.IsNotExistMkDir(dir + "/" + path)
	if err != nil {
		return "", err
	}

	file, err := os.Create(path + filename) // 保存图片的文件名
	if err != nil {
		fmt.Println("图片保存失败:", err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("图片保存失败:", err)
		return "", err
	}
	return imgurl, nil
}

func WriteImg(src multipart.File) (string, error) {

	//生产图片保存地址
	now := time.Now().Format("20060102")
	path := base.Conf.IMG_SavePath + "/hongmo" + "/" + now + "/"
	filename := GetUUID() + ".jpg"
	imgurl := base.Conf.IMG_PrefixUrl + "/hongmo" + "/" + now + "/" + filename

	dir, err := os.Getwd()
	if err != nil {
		return "", err

	}
	err = file.IsNotExistMkDir(dir + "/" + path)
	if err != nil {
		return "", err
	}

	file, err := os.Create(path + filename) // 保存图片的文件名
	if err != nil {
		fmt.Println("图片保存失败:", err)
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, src)
	if err != nil {
		fmt.Println("图片保存失败:", err)
		return "", err
	}
	return imgurl, nil
}
