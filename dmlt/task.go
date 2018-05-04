package dmlt

import (
	"log"
)

// 从本地的配置列表，指定转表
func FromLocal(partFile string) error {
	mapPart, err := ParseOper(partFile)
	if err != nil {
		log.Print(err)
		return err
	}

	listFile := make([]string, 0)
	if commFile, ok := mapPart[commPrefix]; ok {
		listFile = append(listFile, commFile...)
	}

	operName := ""
	for key, operFile := range mapPart {
		if key == commPrefix {
			continue
		}

		operName = key
		listFile = append(listFile, operFile...)
		err = ConvertExcel(listFile)
		if err != nil {
			log.Print(err)
			return err
		}

		// 先只转一个运营分表
		break
	}

	// 只有公共配置表部分
	if operName == "" && len(listFile) > 0 {
		err = ConvertExcel(listFile)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}
