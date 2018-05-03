package fsop

import (
	"log"
	"testing"

	"github.com/lymslive/gotoo/assert"
)

func TestBasic(t *testing.T) {
	assert.BeginTest(t)
	defer assert.EndTest()

	// 去文件后缀名
	log.Println("test TrimExt ...")
	var file = "basic.go"
	assert.Equal(TrimExt(file), "basic", "TrimExt fail")
	file = ".basic.go"
	assert.Equal(TrimExt(file), ".basic", "TrimExt fail")
	file = ".basic"
	assert.Equal(TrimExt(file), ".basic", "TrimExt fail")
	file = "basic_test.go"
	assert.Equal(TrimExt(file), "basic_test", "TrimExt fail")
	file = "basic.test.go"
	assert.Equal(TrimExt(file), "basic.test", "TrimExt fail")

	// 目录文件存在性
	log.Println("test DirExists ...")
	assert.True(DirExists("/usr"), "DirExists fails")
	assert.True(DirExists("/home/lymslive/"), "DirExists fails")
	assert.True(!DirExists("/home/lymslive/bin/vex"), "DirExists fails")
	assert.True(!NotExists("/home/lymslive/bin/vex"), "NotExists fails")
	assert.True(DirExists("../dmlt"), "DirExists fails")
	assert.True(!NotExists("../dmlt/convert.go"), "NotExists fails")
	assert.True(NotExists("../dmlt/convert.gox"), "NotExists fails")

	log.Println("test CopyFile ...")
	n, err := CopyFile("basic.go", "basic.bak")
	assert.True(err == nil, "CopyFile fails")
	log.Printf("copy [%d] byte", n)
}
