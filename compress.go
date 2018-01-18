package elgo

import (
	"encoding/binary"
	"fmt"

	"github.com/hungys/go-lz4"
)

func LZ4Compress(source []byte) ([]byte, error) {
	lensrc := len(source)
	dest := make([]byte, 2*lensrc)
	n, err := lz4.CompressFast(source, dest, 7)
	if err != nil {
		fmt.Printf("Failed to compress source: %v\n", source)
		return nil, err
	}
	binary.BigEndian.PutUint16(dest[n:n+2], uint16(lensrc))
	return dest[:n+2], nil
}

func LZ4DeCompress(source []byte) ([]byte, error) {
	lensrc := len(source) - 2
	nn := binary.BigEndian.Uint16(source[lensrc:])
	dest := make([]byte, nn)
	n, err := lz4.DecompressSafe(source[:lensrc], dest)
	if err != nil {
		fmt.Printf("Failed to decompress source: %v\n", source)
		return nil, err
	}
	if uint16(n) != nn {
		fmt.Printf("Size after decompress %d is not the same as %s.\n", n, nn)
		return nil, fmt.Errorf("Error on decompress packet size")
	}
	return dest[:], nil
}
