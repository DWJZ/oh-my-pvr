package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	uploadDir = "./uploads" // 文件保存目录
)

func main() {
	// 创建上传目录
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create upload directory:", err)
		return
	}

	// 文件上传接口
	http.HandleFunc("/upload", enableCORS(uploadHandler))

	// 启动服务器
	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

// 启用 CORS 的中间件
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 头
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有域名访问
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 如果是预检请求（OPTIONS），直接返回成功
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next(w, r)
	}
}

// 文件上传处理函数
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		http.Error(w, "Request is nil", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 28); err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusInternalServerError)
		return
	}

	// 获取所有文件
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// 使用 WaitGroup 等待所有文件上传完成
	var wg sync.WaitGroup
	wg.Add(len(files))

	// 用于存储上传结果
	results := make([]string, len(files))
	errors := make([]error, len(files))

	// 并行处理每个文件
	for i, fileHeader := range files {
		go func(index int, header *multipart.FileHeader) {
			defer wg.Done()

			// 打开文件
			file, err := header.Open()
			if err != nil {
				errors[index] = err
				return
			}
			defer file.Close()

			// 创建目标文件
			filePath := filepath.Join(uploadDir, header.Filename)
			dst, err := os.Create(filePath)
			if err != nil {
				errors[index] = err
				return
			}
			defer dst.Close()

			// 将上传的文件内容复制到目标文件
			if _, err := io.Copy(dst, file); err != nil {
				errors[index] = err
				return
			}

			// 记录成功信息
			results[index] = fmt.Sprintf("File uploaded successfully: %s", header.Filename)
		}(i, fileHeader)
	}

	// 等待所有文件处理完成
	wg.Wait()

	// 检查是否有错误
	for _, err := range errors {
		if err != nil {
			http.Error(w, "Failed to upload some files", http.StatusInternalServerError)
			return
		}
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	for _, result := range results {
		fmt.Fprintln(w, result)
	}
}
