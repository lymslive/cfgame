package dmlt

import (
	"bufio"
	"fmt"
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
	pattern := operPrefix + "_(.+)/"
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
	log.Printf("will parse file[%s]\n", path)
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
		if strings.TrimSpace(line) == "" {
			continue
		}

		idx := strings.Index(line, subdir.excel)
		if idx < 0 {
			log.Printf("line[%s] has no substr[%s]? idx[%d]", line, subdir.excel, idx)
			continue
			return nil, fmt.Errorf("may not excel file")
		}
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
	file = filepath.ToSlash(file)
	lsMatch := operRegexp.FindStringSubmatch(file)
	if lsMatch != nil && len(lsMatch) > 1 {
		return lsMatch[1]
	}

	return commPrefix
}

/*
/product/data_hf/xls/运营_49/活动开启配置表.xlsx
/product/data_hf/xls/运营_JZ/活动开启配置表.xlsx
/product/data_hf/xls/运营_fh/活动开启配置表.xlsx
/product/data_hf/xls/运营_xw/活动开启配置表.xlsx
*/
