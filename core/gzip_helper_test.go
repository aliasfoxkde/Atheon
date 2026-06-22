package core

import (
	"bytes"
	"compress/gzip"
)

// gzipBytes compresses the input using gzip.
func gzipBytes(in []byte) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(in); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
