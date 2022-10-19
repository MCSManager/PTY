package utils

import (
	"archive/zip"
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 示例 zip.Zip("MCSManager 9.4.5_win64_x86", "./test.zip") 可使用相对路径和绝对路径
func Zip(filePath, zipPath string) error {
	var err error
	if filePath, err = filepath.Abs(filePath); err != nil {
		return err
	}
	if zipPath, err = filepath.Abs(zipPath); err != nil {
		return err
	}
	if strings.ToLower(filepath.Ext(zipPath)) != ".zip" {
		zipPath = filepath.Base(zipPath) + ".zip"
	}
	zipfile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	buf := bufio.NewWriter(zipfile)
	zw := zip.NewWriter(buf)
	defer zw.Close()
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(filePath)
	}
	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filePath == path {
			return nil
		}
		var zipfile io.Writer
		if info.IsDir() {
			_, err = zw.Create(baseDir + strings.TrimPrefix(path, filePath) + `/`)
			return err
		} else {
			zipfile, err = zw.Create(baseDir + strings.TrimPrefix(path, filePath))
			if err != nil {
				return err
			}
		}
		f1, err := os.Open(path)
		if err != nil {
			return err
		}
		io.Copy(zipfile, f1)
		f1.Close()
		return nil
	})
	buf.Flush()
	return err
}
