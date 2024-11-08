package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	orderID := "12345"
	uploadURL := fmt.Sprintf("http://localhost:8080/upload/%s", orderID)
	imagesURL := fmt.Sprintf("http://localhost:8080/images/%s", orderID)
	demander := "张三" // 示例 demander

	// 测试上传图片
	err := uploadFile(uploadURL, "path/to/your/image.jpg", demander)
	if err != nil {
		fmt.Println("上传图片失败:", err)
	} else {
		fmt.Println("图片上传成功")
	}

	// 测试获取图片
	resp, err := http.Get(imagesURL)
	if err != nil {
		fmt.Println("获取图片失败:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	fmt.Println("获取图片成功:", string(body))
}

func uploadFile(url, filePath, demander string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// 添加 'demander' 字段
	err = writer.WriteField("demander", demander)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("上传失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
