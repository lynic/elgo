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
	binary.BigEndian.PutUint32(dest[n:n+4], uint32(lensrc))
	return dest[:n+4], nil
}

func LZ4DeCompress(source []byte) ([]byte, error) {
	if len(source) < 5 {
		return nil, fmt.Errorf("Invalid length")
	}
	lensrc := len(source) - 4
	nn := binary.BigEndian.Uint32(source[lensrc:])
	dest := make([]byte, nn)
	n, err := lz4.DecompressSafe(source[:lensrc], dest)
	if err != nil {
		fmt.Printf("Failed to decompress source: %v\n", source)
		return nil, err
	}
	if uint32(n) != nn {
		fmt.Printf("Size after decompress %d is not the same as %s.\n", n, nn)
		return nil, fmt.Errorf("Error on decompress packet size")
	}
	return dest[:], nil
}
