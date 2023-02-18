package utils

import (
	"compress/flate"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/klauspost/compress/zip"

	archiver "github.com/mholt/archiver/v4"
)

var initZipCompressor = sync.Once{}

func _initZipCompressor() {
	initZipCompressor.Do(func() {
		zip.RegisterCompressor(flate.BestCompression, func(w io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(w, flate.BestCompression)
		})
		zip.RegisterCompressor(flate.BestSpeed, func(w io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(w, flate.BestSpeed)
		})
		zip.RegisterDecompressor(flate.BestCompression, flate.NewReader)
		zip.RegisterDecompressor(flate.BestSpeed, flate.NewReader)
	})
}

type ZipCfg struct {
	Ctx        context.Context
	FilePath   []string
	Exhaustive bool
}

func Zip(ZipPath string, cfg ZipCfg) error {
	_initZipCompressor()
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}
	if len(cfg.FilePath) == 0 {
		return errors.New("file is nil")
	}
	var err error
	cfg.FilePath[0], err = filepath.Abs(cfg.FilePath[0])
	if err != nil {
		return err
	}
	var baseDir = filepath.Dir(cfg.FilePath[0])
	if len(cfg.FilePath) == 1 {
		fi, err := os.Stat(cfg.FilePath[0])
		if err != nil {
			return err
		}
		if fi.IsDir() {
			baseDir = cfg.FilePath[0]
		}
	}
	for k, v := range cfg.FilePath[1:] {
		cfg.FilePath[k+1], err = filepath.Abs(v)
		if err != nil {
			return err
		}
		if filepath.Dir(cfg.FilePath[k+1]) != baseDir {
			return errors.New("base dir err")
		}
	}
	ZipPath, err = filepath.Abs(ZipPath)
	if err != nil {
		return err
	}
	zipExi := strings.ToLower(filepath.Ext(ZipPath))
	var format archiver.CompressedArchive
	switch zipExi {
	case "":
		ZipPath += ".zip"
		format = archiver.CompressedArchive{
			Archival: archiver.Zip{Compression: zip.Deflate, SelectiveCompression: true},
		}
	case ".tar":
		format = archiver.CompressedArchive{
			Archival: archiver.Tar{},
		}
	case ".gz", ".tgz":
		format = archiver.CompressedArchive{
			Compression: archiver.Gz{CompressionLevel: flate.BestCompression, Multithreaded: true},
			Archival:    archiver.Tar{},
		}
	case ".zip":
		format = archiver.CompressedArchive{
			Archival: archiver.Zip{Compression: zip.Deflate, SelectiveCompression: true},
		}
	}
	fileMap := make(map[string]string)
	for _, fPath := range cfg.FilePath {
		select {
		case <-cfg.Ctx.Done():
			return cfg.Ctx.Err()
		default:
			if cfg.Exhaustive {
				fmt.Println(fPath)
			}
			fileMap[fPath] = strings.TrimPrefix(strings.TrimPrefix(fPath, baseDir), string(os.PathSeparator))
		}
	}
	files, err := archiver.FilesFromDisk(nil, fileMap)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(ZipPath), os.ModePerm)
	if err != nil {
		return err
	}
	zipfile, err := os.Create(ZipPath)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	fmt.Println("Archiving, please wait...")
	err = format.Archive(cfg.Ctx, zipfile, files)
	if err != nil {
		return err
	}
	return err
}
