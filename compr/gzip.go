package compr

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GzipDecompress(in []byte) (out []byte, err error) {
	gr, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		return
	}
	out, err = io.ReadAll(gr)
	if err != nil {
		return
	}
	return
}

func GzipCompress(in []byte) (out []byte, err error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err = w.Write(in)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	return buf.Bytes(), err
}
