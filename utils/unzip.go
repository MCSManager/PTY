package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	archiver "github.com/mholt/archiver/v4"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

const bufSize = 512 * 1024

type UnzipCfg struct {
	BufferSize                int
	Ctx                       context.Context
	CoderTypes                CoderType
	SkipExistFile, Exhaustive bool
}

func Unzip(zipPath string, TargetPath string, cfg UnzipCfg) (err error) {
	if zipPath, err = filepath.Abs(zipPath); err != nil {
		return
	}
	zipFile, err := os.Open(zipPath)
	if err != nil {
		return
	}
	defer zipFile.Close()
	return UnzipWithFile(zipFile, TargetPath, cfg)
}

func UnzipWithFile(zipFile io.Reader, TargetPath string, cfg UnzipCfg) error {
	_initZipCompressor()
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}
	if cfg.BufferSize == 0 {
		cfg.BufferSize = bufSize
	}
	seek, ok := zipFile.(io.Seeker)
	if !ok {
		return errors.New("seek file error")
	}
	var err error
	if TargetPath, err = filepath.Abs(TargetPath); err != nil {
		return err
	}
	err = os.MkdirAll(TargetPath, os.ModePerm)
	if err != nil {
		return err
	}
	format, _, err := archiver.Identify("", zipFile)
	if err != nil {
		return err
	}
	if cfg.CoderTypes == T_Auto {
		m := zipEncode(cfg.Ctx, format, zipFile, isUtf8, isGBK)
		_, err = seek.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		if m[T_GBK] && !m[T_UTF8] {
			cfg.CoderTypes = T_GBK
			err = decode(format, zipFile, TargetPath, cfg)
		} else {
			err = decode(format, zipFile, TargetPath, cfg)
		}
	} else {
		err = decode(format, zipFile, TargetPath, cfg)
	}
	return err
}

func zipEncode(ctx context.Context, format archiver.Format, r io.Reader, fun ...func(data []byte) (bool, CoderType)) (res map[CoderType]bool) {
	res = make(map[CoderType]bool)
	if ex, ok := format.(archiver.Extractor); ok {
		ex.Extract(ctx, r, nil, func(ctx context.Context, f archiver.File) error {
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

func decode(format archiver.Format, r io.Reader, TargetPath string, cfg UnzipCfg) error {
	var decoder *encoding.Decoder
	if cfg.CoderTypes != T_Auto {
		decoder = newDeCoder(cfg.CoderTypes)
	}
	if ex, ok := format.(archiver.Extractor); ok {
		buffer := make([]byte, cfg.BufferSize)
		return ex.Extract(cfg.Ctx, r, nil, func(ctx context.Context, f archiver.File) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				var result string
				if decoder == nil {
					result = f.NameInArchive
				} else {
					var err error
					result, _, err = transform.String(decoder, f.NameInArchive)
					if err != nil {
						fmt.Printf("File %s err: %v", f.NameInArchive, err)
						return err
					}
				}
				if cfg.Exhaustive {
					fmt.Println(result)
				}
				fpath := filepath.Join(TargetPath, result)
				if f.IsDir() {
					return os.MkdirAll(fpath, f.Mode())
				}
				if cfg.SkipExistFile {
					_, err := os.Stat(fpath)
					if err == nil {
						return err
					}
				}
				inFile, err := f.Open()
				if err != nil {
					return err
				}
				defer inFile.Close()

				if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
					return err
				}
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
		})
	}
	return errors.New("format.(archiver.Extractor) err")
}
