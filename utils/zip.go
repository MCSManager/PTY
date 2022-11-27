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
func Zip(filePath []string, zipPath string) error {
	zipPath, err := filepath.Abs(zipPath)
	if err != nil {
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
	defer buf.Flush()
	zw := zip.NewWriter(buf)
	defer zw.Close()
	for _, fPath := range filePath {
		fPath, err = filepath.Abs(fPath)
		if err != nil {
			return err
		}
		err = filepath.Walk(fPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				_, err = zw.Create(strings.TrimPrefix(strings.TrimPrefix(path, filepath.Dir(fPath)), string(os.PathSeparator)) + `/`)
				return err
			}
			zipfile, err := zw.Create(strings.TrimPrefix(strings.TrimPrefix(path, filepath.Dir(fPath)), string(os.PathSeparator)))
			if err != nil {
				return err
			}
			f1, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f1.Close()
			_, err = io.Copy(zipfile, f1)
			return err
		})
		if err != nil {
			return err
		}
	}
	return err
}
