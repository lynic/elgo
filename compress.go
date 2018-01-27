package elgo

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"

	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
	"github.com/ulikunitz/xz"
)

func LZ4Compress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w := lz4.NewWriter(&in)
	defer w.Close()
	_, err := w.Write(source)
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func LZ4DeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader := lz4.NewReader(buf)
	var out bytes.Buffer
	_, err := io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func ZlibCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := zlib.NewWriterLevel(&in, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	_, err = w.Write(source)
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func ZlibDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func GZipCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := gzip.NewWriterLevel(&in, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	_, err = w.Write(source)
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func GZipDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func FlateCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := flate.NewWriter(&in, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	_, err = w.Write(source)
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func FlateDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader := flate.NewReader(buf)
	defer reader.Close()
	var out bytes.Buffer
	_, err := io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func SnappyCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w := snappy.NewWriter(&in)
	defer w.Close()
	_, err := w.Write(source)
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func SnappyDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader := snappy.NewReader(buf)
	var out bytes.Buffer
	_, err := io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func XzCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, err := xz.NewWriter(&in)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	_, err = w.Write(source)
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func XzDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader, err := xz.NewReader(buf)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	_, err = io.Copy(&out, reader)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
