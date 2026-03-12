package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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

	fmt.Println("Watching:", downloadDir)

	// create new watcher
	watcher, err := fsnotify.NewWatcher()
	checkErr(err)
	defer watcher.Close()

	go watchFiles(watcher)

	err = watcher.Add(downloadDir)
	checkErr(err)

	select {} // 프로그램 계속 실행
}

func watchFiles(watcher *fsnotify.Watcher) {

	for {
		select {
		case event, ok := <- watcher.Events:
			if !ok {
				return
			}

			// 파일 생성 감지
			if event.Has(fsnotify.Create) {
				moveFile(event.Name)
			}

		case err, ok := <- watcher.Errors:
			if !ok {
				return
			}

			log.Panicln("watch err:",err)
		}
	}
}

func moveFile(path string) {
	fileName := filepath.Base(path)

	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
	ext = strings.ToLower(ext)

	if ext == "" {
		return
	}

	folderType, ok := fileTypes[ext]
	if !ok {
		return
	}

	folderPath := filepath.Join(downloadDir, folderType)

	err := os.MkdirAll(folderPath, 0755)
	checkErr(err)

	newPath := filepath.Join(folderPath, fileName)

	err = os.Rename(path, newPath)
	if err != nil {
		log.Println("move failed:", err)
		return
	}

	fmt.Println("Moved:", fileName, "->", folderType)

}