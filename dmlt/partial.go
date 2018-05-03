package dmlt

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	dftPartial string = "partialExcel.txt"
	operPrefix string = "运营"
	commPrefix string = "comm"
)

var operRegexp *regexp.Regexp

func init() {
	pattern := operPrefix + "_(.+)"
	operRegexp = regexp.MustCompile(pattern)
}

// 各种类的源配置表，key => 列表
type PartExcel map[string][]string

// 分析运营的部分转表需求
// 每个运营渠道映射一个文件名列表
// 公用的键名 comm
func ParseOper(file string) (PartExcel, error) {
	if file == "" {
		file = dftPartial
	}

	path := filepath.Join(datadir, file)
	hf, err := os.Open(path)
	if err != nil {
		log.Printf("cannot open file[%s]", path)
		return nil, err
	}
	defer hf.Close()

	var output = PartExcel{}
	input := bufio.NewScanner(hf)
	for input.Scan() {
		line := input.Text()
		idx := strings.Index(line, subdir.excel)
		relate := line[idx:]

		oper := _whichOper(relate)
		if _, ok := output[oper]; !ok {
			output[oper] = make([]string, 0)
		}
		output[oper] = append(output[oper], relate)
	}

	return output, nil
}

// 分析一个文件路径中是否含有(运营_)段
// 返回属性哪个运营标识，没有的话返回 comm
func _whichOper(file string) (oper string) {
	oper = commPrefix

	if !operRegexp.MatchString(file) {
		return
	}

	lsPath := strings.Split(file, string(filepath.Separator))
	for _, part := range lsPath {
		lsMatch := operRegexp.FindStringSubmatch(part)
		if lsMatch != nil && len(lsMatch) > 2 {
			oper = lsMatch[1]
			break
		}
	}

	return
}

/*
/product/data_hf/xls/运营_49/活动开启配置表.xlsx
/product/data_hf/xls/运营_JZ/活动开启配置表.xlsx
/product/data_hf/xls/运营_fh/活动开启配置表.xlsx
/product/data_hf/xls/运营_xw/活动开启配置表.xlsx
*/
