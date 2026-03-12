package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var downloadDir = os.Getenv("HOME") + "/Downloads"

// 확장자 → 폴더 타입
var fileTypes = map[string]string{
	// images
	"png":  "images",
	"jpg":  "images",
	"jpeg": "images",
	"svg":  "images",
	"webp": "images",
	"gif":  "images",

	// videos
	"mp4": "videos",
	"mov": "videos",
	"mkv": "videos",

	// documents
	"pdf":  "documents",
	"csv":  "documents",
	"doc":  "documents",
	"docx": "documents",
	"txt":  "documents",

	// archives
	"zip": "archives",
	"rar": "archives",
	"tar": "archives",
	"gz":  "archives",
	"dmg": "archives",
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main () {
	fmt.Printf("downloadDir: %v\n", downloadDir)

	files, err := os.ReadDir(downloadDir)
	checkErr(err)

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		// 확장자 확인
		ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")

		if ext == "" {
			continue
		}

		// target 확장자 확인
		if fileTypes[ext] == "" {
			continue
		}

	// 확장자 폴더 생성
		folderPath := filepath.Join(downloadDir, fileTypes[ext])

		err := os.MkdirAll(folderPath, 0755)
		checkErr(err)

		oldPath := filepath.Join(downloadDir, file.Name())
		newPath := filepath.Join(folderPath, file.Name())

		err = os.Rename(oldPath, newPath)
		checkErr(err)

		fmt.Println(file.Name(), "→", fileTypes[ext])
	}
}