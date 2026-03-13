package main

import (
	"context"
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

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}


// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Watch() {
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
				NewApp().MoveFile(event.Name)
			}
			
		case error, ok := <- watcher.Errors:
			if !ok {
				return
			}
			log.Panicln("watch err:",error)
		}
	}
}

func (a *App) MoveFile(path string) {
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

func (a *App) ClassifyFolder() {
	files, err := os.ReadDir(downloadDir)
	checkErr(err)

	fmt.Println(files)

	for _, file := range files {
		fmt.Println(downloadDir, file.Name())

		if file.IsDir() {
			continue
		}


			ext := strings.TrimPrefix(filepath.Ext(file.Name()), ".")
	ext = strings.ToLower(ext)

		fmt.Println(ext)

	if ext == "" {
		continue
	}

	 folderType, ok := fileTypes[ext]

	 fmt.Println(folderType, ok)

	 if !ok {
		continue
	 }

	 	folderPath := filepath.Join(downloadDir, folderType)

		prevPath := filepath.Join(downloadDir, file.Name())

		err := os.MkdirAll(folderPath, 0755)
	checkErr(err)

	newPath := filepath.Join(folderPath, file.Name())

	err = os.Rename(prevPath, newPath)

	if err != nil {
		log.Println("move failed:", err)
		return
	}

	fmt.Println("Moved:", file.Name(), "->", folderType)
	}	
}

// func classify(files os.DirEntry) {
// }



func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}