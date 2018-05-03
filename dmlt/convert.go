package dmlt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lymslive/cfgame/fsop"
)

// 将一组指定的文件转表，文件名相对数据根目录
func ConvertExcel(files []string) error {
	log.Println("clean work dir ...")
	if err := _cleanDir(); err != nil {
		log.Println(err)
		return err
	}

	log.Println("copy excel to tmp dir ...")
	for _, file := range files {
		name := filepath.Base(file)
		tmpfile := filepath.Join(subdir.midput, name)
		_, err := fsop.CopyFile(file, tmpfile)
		if err != nil {
			log.Printf("copy file[%s] failed", name)
			return err
		}
	}

	log.Println("convert excel to plain csv ...")
	pCmd := exec.Command(xls2csv, subdir.midput, subdir.sheet, x2cconf)
	if pCmd == nil {
		return fmt.Errorf("cannot build execute cmd: xls2csv ...")
	}
	if err := pCmd.Run(); err != nil {
		log.Println(err)
		return err
	}

	log.Println("convert plain csv to bin data ...")
	nFile, err := _convTDR()
	if err != nil {
		log.Printf("fails convsion: %s", err)
	}
	log.Printf("converted %d tables", nFile)

	return err
}

// 清理目录，并将当前路径切换到数据根目录
func _cleanDir() error {
	var err error

	err = os.Chdir(datadir)
	if err != nil {
		return err
	}

	err = os.RemoveAll(subdir.output)
	if err != nil {
		return err
	}
	err = os.RemoveAll(subdir.midput)
	if err != nil {
		return err
	}
	err = os.RemoveAll(subdir.sheet)
	if err != nil {
		return err
	}

	err = os.Mkdir(subdir.output, os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Mkdir(subdir.midput, os.ModeDir)
	if err != nil {
		return err
	}
	err = os.Mkdir(subdir.sheet, os.ModeDir)
	if err != nil {
		return err
	}

	return nil
}

// 将 csv 目录内所有文件转为 bin 文件
// 调用 resconv.exe 程序
func _convTDR() (nFile int, err error) {
	// 读取目录，保存文件名至一个 map
	dir, err := os.Open(subdir.sheet)
	if err != nil {
		return
	}

	files, err := dir.Readdirnames(0)
	if err != nil {
		return
	}

	var hasFile = make(map[string]bool)
	for _, file := range files {
		file = fsop.TrimExt(file)
		hasFile[file] = true
	}

	// 读取公共转表批处理文件，逐一解析执行
	// 略过每行开始的两个字段 call do_conv.bat
	hf, err := os.Open(convfile)
	if err != nil {
		log.Printf("cannot open file[%s]", convfile)
		return
	}
	defer hf.Close()

	input := bufio.NewScanner(hf)
	for input.Scan() {
		line := input.Text()
		token := strings.Split(line, " ")
		nt := len(token)
		if nt < 4 {
			return 0, fmt.Errorf("bad format of convfile[%s]", convfile)
		}

		if _, ok := hasFile[token[nt-1]]; !ok {
			continue
		}

		tdrname := token[2]
		tblname := token[3]
		csvname := filepath.Join(subdir.sheet, tblname) + ".csv"

		// 拼接 resconv.exe 的参数
		var args = []string{resarg1, tdrfile, tdrname, csvname, subdir.output}
		for i := 4; i < nt; i++ {
			csvname := filepath.Join(subdir.sheet, token[i]) + ".csv"
			args = append(args, csvname)
		}

		pCmd := exec.Command(resconv, args...)
		if pCmd == nil {
			return nFile, fmt.Errorf("cannot build execute cmd: resconv ...")
		}
		if err = pCmd.Run(); err != nil {
			log.Println(err)
			return
		}

		nFile++
		log.Println("convert %d table[%s] from [%s] etc", nFile, tdrname, tblname)
	}

	return
}
