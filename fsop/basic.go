package fsop

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

// 复制文件，要求参数全路径或相对当前目录的路径
func CopyFile(src, dst string) (n int64, err error) {
	if src == "" || dst == "" || src == dst {
		return 0, errors.New("invalid copy src or dst")
	}

	hfSrc, err := os.Open(src)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer hfSrc.Close()

	hfDst, err := os.Create(dst)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer hfDst.Close()

	return io.Copy(hfDst, hfSrc)
}

// 判断目录是否存在
func DirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// 判断文件或目录不存在
func NotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// 返回不包括后缀的文件名
func TrimExt(file string) string {
	dot := strings.LastIndex(file, ".")
	if dot <= 0 {
		return file
	}

	return file[:dot]
}
