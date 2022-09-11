package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

const (
	copyChunkSize int64 = 100
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetLimitNegative   = errors.New("offset and limit must not be negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return ErrOffsetLimitNegative
	}

	sourceFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer destFile.Close()

	fileSize, err := getFileSize(sourceFile)
	if err != nil {
		return err
	}
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	bytesToCopyCount := fileSize - offset
	if limit == 0 || limit > bytesToCopyCount {
		limit = bytesToCopyCount
	}

	err = doCopy(sourceFile, destFile, offset, limit)
	if err != nil {
		return err
	}

	return nil
}

func doCopy(sourceFile *os.File, destFile *os.File, offset int64, limit int64) error {
	sourceFile.Seek(offset, 0)
	sourceReader := io.LimitReader(sourceFile, limit)

	bar := pb.StartNew(int(limit))
	for {
		n, err := io.CopyN(destFile, sourceReader, copyChunkSize)
		bar.Add64(n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}
	bar.Finish()

	return nil
}

func getFileSize(sourceFile *os.File) (int64, error) {
	fileStat, err := sourceFile.Stat()
	if err != nil {
		return 0, err
	}

	fileSize := fileStat.Size()
	if fileSize <= 0 {
		return 0, ErrUnsupportedFile
	}

	return fileSize, nil
}
