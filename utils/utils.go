package utils

import (
	"archive/zip"
	"bufio"
	"io"
	"os"

	"github.com/Axope/JOJ/common/log"
)

// 解压
func Unzip(path string, dir string) error {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	CreateDir(dir)
	for _, f := range zipReader.File {
		// 打开文件
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		// 创建目标文件
		dst, err := os.Create(dir + f.Name)
		if err != nil {
			return err
		}
		defer dst.Close()
		// 将文件内容复制到目标文件中
		_, err = io.Copy(dst, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Logger.Error("remove file error", log.Any("err", err))
	}
}

func CreateDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			log.Logger.Error("Failed to create directory:", log.Any("err", err))
			return
		}
	}
}

func GetStringSlice(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		log.Logger.Error("Error opening file:", log.Any("err", err))
		return nil
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// 检查是否有读取文件的错误
	if err := scanner.Err(); err != nil {
		log.Logger.Error("Error reading file:", log.Any("err", err))
		return nil
	}

	return lines
}
