package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("unsupported file error case", func(t *testing.T) {
		from := "/dev/urandom"
		to := "copy_test.txt"
		limit := int64(0)
		offset := int64(0)

		err := Copy(from, to, offset, limit)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)

		os.Remove(to)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		from := "test_offset_exceeds_file_size_source"
		to := "test_offset_exceeds_file_size_dest"
		limit := int64(0)
		offset := int64(4)

		fileFrom, _ := os.Create(from)
		fileFrom.WriteString("123")
		os.Create(to)

		err := Copy(from, to, offset, limit)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)

		os.Remove(from)
		os.Remove(to)
	})

	t.Run("valid case", func(t *testing.T) {
		from := "test_offset_exceeds_file_size_source"
		to := "test_offset_exceeds_file_size_dest"
		limit := int64(4)
		offset := int64(5)

		fileFrom, _ := os.Create(from)
		fileFrom.WriteString("some content")
		fileTo, _ := os.Create(to)

		err := Copy(from, to, offset, limit)
		require.Equal(t, nil, err)

		destContent := make([]byte, 4)
		fileTo.Read(destContent)
		require.Equal(t, "cont", string(destContent))

		os.Remove(from)
		os.Remove(to)
	})
}
