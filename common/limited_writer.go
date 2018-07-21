

package common

import (
	"errors"
	"io"
)

var ErrWriteExceedLimitedCount = errors.New("writer exceed limited count")

type LimitedWriter struct {
	count  uint64
	max    uint64
	writer io.Writer
}

func NewLimitedWriter(w io.Writer, max uint64) *LimitedWriter {
	return &LimitedWriter{
		writer: w,
		max:    max,
	}
}

func (self *LimitedWriter) Write(buf []byte) (int, error) {
	if self.count+uint64(len(buf)) > self.max {
		return 0, ErrWriteExceedLimitedCount
	}
	n, err := self.writer.Write(buf)
	self.count += uint64(n)
	return n, err
}

// Count function return counted bytes
func (self *LimitedWriter) Count() uint64 {
	return self.count
}
