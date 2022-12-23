package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/kbinani/screenshot"
)

func main() {
	dir, err := os.MkdirTemp("", "capscr_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir) // clean up
	fmt.Println(dir)
	for {
		postFile(screen(dir))
		time.Sleep(time.Second * 3)
	}
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

func postFile(fullName string) {
	fileName := path.Base((fullName))
	//这是一个Post 参数会被返回的地址
	uri := "http://192.168.3.131:8888/" + fileName
	byte, err := os.ReadFile(fullName)
	if err != nil {
		fmt.Println("readFile err=", err)
	}
	res, err := http.Post(uri, "image/png", bytes.NewReader(byte))
	if err != nil {
		fmt.Println("post err=", err)
	}
	//http返回的response的body必须close,否则就会有内存泄露
	defer func() {
		res.Body.Close()
	}()
	//读取body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("res err=", err)
	}
	fmt.Print(string(body))
}
