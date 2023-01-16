package utils

import (
	"archive/zip"
	"bufio"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/text/transform"
)

const bufSize = 1024 * 1024 * 4

// 示例: zip.Unzip("./mcsm.zip", "./", "auto") 可使用相对路径和绝对路径
func Unzip(zipPath, targetPath, coder string) error {
	var err error
	if targetPath, err = filepath.Abs(targetPath); err != nil {
		return err
	}
	err = os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return err
	}
	if zipPath, err = filepath.Abs(zipPath); err != nil {
		return err
	}
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	if coder == "auto" {
		if zipEncode(zipReader.File, isUtf8) || !zipEncode(zipReader.File, isGBK) {
			err = decode(zipReader.File, targetPath, "utf8")
		} else {
			err = decode(zipReader.File, targetPath, "gbk")
		}
	} else {
		err = decode(zipReader.File, targetPath, coder)
	}
	return err
}

func zipEncode(f []*zip.File, fun func(data []byte) bool) bool {
	for _, v := range f {
		if fun([]byte(v.Name)) {
			continue
		}
		return false
	}
	return true
}

func decode(files []*zip.File, targetPath string, types string) error {
	decoder := newDeCoder(types)
	for _, f := range files {
		if result, _, err := transform.String(decoder, f.Name); err != nil {
			return err
		} else if err = handleFile(f, targetPath, result); err != nil {
			return err
		}
	}
	return nil
}

func handleFile(f *zip.File, targetPath, decodeName string) error {
	fpath := filepath.Join(targetPath, decodeName)
	if f.FileInfo().IsDir() {
		return os.MkdirAll(fpath, os.ModePerm)
	} else {
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()
		file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		var outFile io.Writer
		if f.UncompressedSize64 > bufSize {
			buf := bufio.NewWriterSize(file, 4*bufSize)
			outFile = buf
			defer buf.Flush()
		} else {
			outFile = file
		}
		_, err = io.CopyBuffer(outFile, inFile, make([]byte, bufSize))
		return err
	}
}
