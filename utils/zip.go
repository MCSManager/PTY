package utils

import (
	"archive/zip"
	"bufio"
	"compress/flate"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 示例 zip.Zip("MCSManager 9.4.5_win64_x86", "./test.zip") 可使用相对路径和绝对路径
func Zip(filePath []string, zipPath string) error {
	if len(filePath) == 0 {
		return errors.New("file is nil")
	}
	var err error
	filePath[0], err = filepath.Abs(filePath[0])
	if err != nil {
		return err
	}
	var baseDir = filepath.Dir(filePath[0])
	if len(filePath) == 1 {
		fi, err := os.Stat(filePath[0])
		if err != nil {
			return err
		}
		if fi.IsDir() {
			baseDir = filePath[0]
		}
	}
	for k, v := range filePath[1:] {
		filePath[k+1], err = filepath.Abs(v)
		if err != nil {
			return err
		}
		if filepath.Dir(filePath[k+1]) != baseDir {
			return errors.New("base dir err")
		}
	}
	zipPath, err = filepath.Abs(zipPath)
	if err != nil {
		return err
	}
	if strings.ToLower(filepath.Ext(zipPath)) != ".zip" {
		zipPath += ".zip"
	}
	err = os.MkdirAll(filepath.Dir(zipPath), os.ModePerm)
	if err != nil {
		return err
	}
	zipFileNamePrefix := strings.TrimPrefix(strings.TrimPrefix(zipPath, baseDir), string(os.PathSeparator))
	zipfile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	buf := bufio.NewWriterSize(zipfile, 4*bufSize)
	defer buf.Flush()
	zw := zip.NewWriter(buf)
	defer zw.Close()
	zw.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
	for _, fPath := range filePath {
		err = filepath.Walk(fPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			filePath := strings.TrimPrefix(strings.TrimPrefix(path, baseDir), string(os.PathSeparator))
			if info.IsDir() {
				_, err = zw.Create(filePath + `/`)
				return err
			}
			if filePath == zipFileNamePrefix {
				return nil
			}
			zipfile, err := zw.Create(filePath)
			if err != nil {
				return err
			}
			f1, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f1.Close()
			_, err = io.CopyBuffer(zipfile, f1, make([]byte, bufSize))
			return err
		})
		if err != nil {
			return err
		}
	}
	return err
}
