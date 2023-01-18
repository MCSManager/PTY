package utils

import (
	"compress/flate"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"

	archiver "github.com/mholt/archiver/v4"
)

func init() {
	zip.RegisterCompressor(flate.BestCompression, func(w io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(w, flate.BestCompression)
	})
	zip.RegisterCompressor(flate.BestSpeed, func(w io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(w, flate.BestSpeed)
	})
	zip.RegisterDecompressor(flate.BestCompression, flate.NewReader)
	zip.RegisterDecompressor(flate.BestSpeed, flate.NewReader)
}

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
	zipExi := strings.ToLower(filepath.Ext(zipPath))
	var format archiver.CompressedArchive
	switch zipExi {
	case "":
		zipPath += ".zip"
		format = archiver.CompressedArchive{
			Archival: archiver.Zip{Compression: zip.Deflate, SelectiveCompression: true},
		}
	case ".tar":
		format = archiver.CompressedArchive{
			Archival: archiver.Tar{},
		}
	case ".gz", ".tgz":
		format = archiver.CompressedArchive{
			Compression: archiver.Gz{CompressionLevel: flate.DefaultCompression, Multithreaded: true},
			Archival:    archiver.Tar{},
		}
	case ".zip":
		format = archiver.CompressedArchive{
			Archival: archiver.Zip{Compression: zip.Deflate, SelectiveCompression: true},
		}
	}
	fileMap := make(map[string]string)
	for _, fPath := range filePath {
		fileMap[fPath] = strings.TrimPrefix(strings.TrimPrefix(fPath, baseDir), string(os.PathSeparator))
	}
	files, err := archiver.FilesFromDisk(nil, fileMap)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(zipPath), os.ModePerm)
	if err != nil {
		return err
	}
	zipfile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	err = format.Archive(context.Background(), zipfile, files)
	if err != nil {
		return err
	}
	return err
}
