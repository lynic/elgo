package elgo

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"

	"github.com/pierrec/lz4"
)

// func LZ4Compress(source []byte) ([]byte, error) {
// 	lensrc := len(source)
// 	dest := make([]byte, 2*lensrc)
// 	n, err := lz4.CompressFast(source, dest, 7)
// 	if err != nil {
// 		fmt.Printf("Failed to compress source: %v\n", source)
// 		return nil, err
// 	}
// 	binary.BigEndian.PutUint32(dest[n:n+4], uint32(lensrc))
// 	return dest[:n+4], nil
// }

func LZ4Compress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, _ := lz4.NewWriterLevel(&in, lz4.BestCompression)
	_, err := w.Write(source)
	defer w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func LZ4DeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader, err := lz4.NewReader(buf)
	defer reader.Close()
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

// func LZ4DeCompress(source []byte) ([]byte, error) {
// 	if len(source) < 5 {
// 		return nil, fmt.Errorf("Invalid length")
// 	}
// 	lensrc := len(source) - 4
// 	nn := binary.BigEndian.Uint32(source[lensrc:])
// 	dest := make([]byte, nn)
// 	n, err := lz4.DecompressSafe(source[:lensrc], dest)
// 	if err != nil {
// 		fmt.Printf("Failed to decompress source: %v\n", source)
// 		return nil, err
// 	}
// 	if uint32(n) != nn {
// 		fmt.Printf("Size after decompress %d is not the same as %s.\n", n, nn)
// 		return nil, fmt.Errorf("Error on decompress packet size")
// 	}
// 	return dest[:], nil
// }

func ZlibCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevel(&in, zlib.BestCompression)
	_, err := w.Write(source)
	defer w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func ZlibDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader, err := zlib.NewReader(buf)
	defer reader.Close()
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

func GZipCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, _ := gzip.NewWriterLevel(&in, gzip.BestCompression)
	_, err := w.Write(source)
	defer w.Close()
	if err != nil {
		return nil, err
	}
	return in.Bytes(), nil
}

func GZipDeCompress(source []byte) ([]byte, error) {
	buf := bytes.NewBuffer(source[:])
	reader, err := gzip.NewReader(buf)
	defer reader.Close()
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

func FlateCompress(source []byte) ([]byte, error) {
	var in bytes.Buffer
	w, _ := flate.NewWriter(&in, flate.BestCompression)
	_, err := w.Write(source)
	defer w.Close()
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
