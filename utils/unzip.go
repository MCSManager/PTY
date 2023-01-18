package utils

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	archiver "github.com/mholt/archiver/v4"

	"golang.org/x/text/transform"
)

const bufSize = 512 * 1024

// 示例: zip.Unzip("./mcsm.zip", "./", "auto") 可使用相对路径和绝对路径
func Unzip(zipPath, targetPath string, coderTypes CoderType) error {
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
	zipFile, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	return UnzipWithFile(zipFile, targetPath, coderTypes)
}

func UnzipWithFile(file io.Reader, targetPath string, coderTypes CoderType) error {
	seek, ok := file.(io.Seeker)
	if !ok {
		return errors.New("seek file error")
	}
	var err error
	if targetPath, err = filepath.Abs(targetPath); err != nil {
		return err
	}
	err = os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return err
	}
	format, _, err := archiver.Identify("", file)
	if err != nil {
		return err
	}
	if coderTypes == T_Auto {
		m := zipEncode(format, file, isUtf8, isGBK)
		_, err = seek.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		if m[T_UTF8] || !m[T_GBK] {
			err = decode(format, file, targetPath, T_UTF8)
		} else {
			err = decode(format, file, targetPath, T_GBK)
		}
	} else {
		err = decode(format, file, targetPath, coderTypes)
	}
	return err
}

func zipEncode(format archiver.Format, r io.Reader, fun ...func(data []byte) (bool, CoderType)) (res map[CoderType]bool) {
	res = make(map[CoderType]bool)
	if ex, ok := format.(archiver.Extractor); ok {
		ex.Extract(context.Background(), r, nil, func(ctx context.Context, f archiver.File) error {
			for _, fn := range fun {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					ok, name := fn([]byte(f.Name()))
					if b, o := res[name]; o {
						if !b {
							continue
						} else {
							res[name] = ok
						}
					} else {
						res[name] = ok
					}
				}
			}
			return nil
		})
	}
	return
}

func decode(format archiver.Format, r io.Reader, targetPath string, coderTypes CoderType) error {
	decoder := newDeCoder(coderTypes)
	if ex, ok := format.(archiver.Extractor); ok {
		buffer := make([]byte, bufSize)
		return ex.Extract(context.Background(), r, nil, func(ctx context.Context, f archiver.File) error {
			if result, _, err := transform.String(decoder, f.NameInArchive); err != nil {
				return err
			} else {
				fpath := filepath.Join(targetPath, result)
				if f.IsDir() {
					return os.MkdirAll(fpath, f.Mode())
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
					if f.Size() > bufSize {
						buf := bufio.NewWriterSize(file, 4*bufSize)
						outFile = buf
						defer buf.Flush()
					} else {
						outFile = file
					}
					_, err = io.CopyBuffer(outFile, inFile, buffer)
					return err
				}
			}
		})
	}
	return errors.New("format.(archiver.Extractor) err")
}
