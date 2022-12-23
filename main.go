package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kbinani/screenshot"
)

type Config struct {
	Url      string
	Interval int
}

func main() {
	cfg := readConfig()

	dir, err := os.MkdirTemp("", "capscr_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir) // clean up
	fmt.Println(dir)
	for {
		postFile(screen(dir), cfg.Url)
		time.Sleep(time.Second * time.Duration(cfg.Interval))
	}
}

func readConfig() Config {
	cfgPath := path.Join(getExeDir(), "config.json")
	fmt.Println(cfgPath)

	filePtr, err := os.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer filePtr.Close()
	var cfg Config
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func screen(dir string) string {
	//使用 GetDisplayBounds获取指定屏幕显示范围，全屏截图
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	//拼接图片名
	t := time.Now().Unix()
	name := strconv.Itoa(int(t)) + ".png"
	fullName := path.Join(dir, name)
	save(img, fullName)
	return fullName
}

// save *image.RGBA to filePath with PNG format.
func save(img *image.RGBA, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func postFile(fullName string, url string) {
	fileName := path.Base((fullName))

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, _ := bodyWriter.CreateFormFile("file", fileName)

	file, err := os.Open(fullName)
	if err != nil {
		return
	}
	defer file.Close()

	io.Copy(fileWriter, file)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuffer)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Print(string(resp_body))
}

func getExeDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}
