# Download Organazier

## Keyword

`fsnotify`, `os.Getenv`, `os.ReadDir`, `file.IsDir`, `TrimPrefix`, `filepath.Join`, `os.MkdirAll`, `os.Rename`, `wails`

## Libarary

- fsnotify (파일 감시 라이브러리)
  > go get github.com/fsnotify/fsnotify

---

### 다운로드 경로 파악

```go
var downloadDir = os.Getenv("HOME") + "/Downloads"
```

### 폴더 탐색

```go
files, err := os.ReadDir(downloadDir)
```

### slice 대신 map

#### 기존

```go
var targetExts = []string{"png","zip"}
```

검색
`O(n)`

#### 개선

```go
var targetExts = map[string]bool
```

검색
`O(1)`

### filepath.Join 사용

#### 기존

```go
downloadDir + "/" + file.Name()
```

문제

> OS별 path separator 다름

#### 개선

```go
filepath.Join(downloadDir, file.Name())
```

### Mkdir → MkdirAll

#### 기존

```go
os.Mkdir(folderPath, 0755)
```

문제

> 이미 존재하면 에러

#### 개선

```go
os.MkdirAll(folderPath, 0755)
```

Mkdir, MkdirAll 의 두번째 인자는 파일 모드인데

```go
os.MkdirAll(path string, perm os.FileMode)
```

```
0755
 │││
 ││└ others (다른 사용자)
 │└ group  (같은 그룹)
 └ owner  (소유자)
```

| 숫자 | 권한        |
| ---- | ----------- |
| 4    | read (r)    |
| 2    | write (w)   |
| 1    | execute (x) |

> 0755 해석

```
owner  = 7 = rwx
group  = 5 = r-x
others = 5 = r-x
```

즉

```
owner   → 읽기 / 쓰기 / 실행
group   → 읽기 / 실행
others  → 읽기 / 실행
```

앞의 0은 `8진수(octal)`라는 의미

#### 자주 쓰는 권한

| 권한 | 의미              |
| ---- | ----------------- |
| 0755 | 일반적인 디렉토리 |
| 0700 | 나만 접근         |
| 0644 | 일반 파일         |
| 0600 | 비밀 파일         |

### TrimPrefix로 확장자 제거

기존

```go
ext = ext[1:]
```

개선

```go
strings.TrimPrefix(ext, ".")
```

더 안전함

### 폴더/파일 감시

`for { select { ... } }` 는 채널을 계속 감시하는 이벤트 루프를 만들 때 쓰는 패턴

## 전체 코드

```go
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
```

## 데스크탑 앱

`Go + Wails`

설치

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```
