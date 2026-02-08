package main

import (
	"io"
)

func ioReadAllLimit(r io.Reader, limit int64) ([]byte, error) {
	lr := &io.LimitedReader{R: r, N: limit}
	return io.ReadAll(lr)
}
