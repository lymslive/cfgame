package dmlt

import (
	"log"
	"path/filepath"

	"github.com/lymslive/cfgame/fsop"
)

// 数据基础目录
var datadir string

func SetDataDir(path string) (old string) {
	if path == "" {
		return datadir
	}

	old = datadir
	if !fsop.DirExists(path) {
		log.Printf("may not a directory[%s]", path)
		return
	}

	datadir = path
	return
}

// 子目录名
var subdir struct {
	excel  string
	sheet  string
	output string
	midput string
}

// xls2csv 转化工具路径及配置
var xls2csv string
var x2cconf string

// tdr 转表工具参数
var resconv = "resconv.exe"
var resarg1 = "conv"
var tdrfile = "ResMeta.tdr"
var convfile = "策划转表_公共.bat"

func init() {
	subdir.excel = "xls"
	subdir.sheet = "csv"
	subdir.output = "bin"
	subdir.midput = "xls_tmp"

	xls2csv = filepath.Join("x2c", "xls2csv.exe")
	x2cconf = "x2c.x2c"
}
