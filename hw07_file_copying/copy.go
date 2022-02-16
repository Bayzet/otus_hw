package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrLimitCannotBeNegative = errors.New("limit cannot be negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if limit < 0 {
		return ErrLimitCannotBeNegative
	}

	fileTo, _ := os.Create(toPath)
	defer fileTo.Close()
	fwriter := io.Writer(fileTo)

	fileFrom, err := os.Open(fromPath)
	defer fileFrom.Close()
	if err != nil {
		return err
	}

	fStat, err := fileFrom.Stat()
	if !fStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if fStat.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	fileFrom.Seek(offset, io.SeekStart)
	freader := io.Reader(fileFrom)

	bar := pb.Start64(limit)

	barReader := bar.NewProxyReader(freader)

	if limit == 0 {
		limit = fStat.Size()
	}

	if _, err := io.CopyN(fwriter, barReader, limit); err != nil {
		if err != io.EOF {
			return err
		}
	}
	bar.Finish()

	return nil
}
