package elgo

import (
	"io/ioutil"
	"net/http"
)

func ReadHTTPBody(r *http.Request) ([]byte, error) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var payload []byte
	payload = content[:]
	compress := r.Header.Get("Content-Encoding")
	switch compress {
	case "lz4":
		payload, err = LZ4DeCompress(content[:])
		if err != nil {
			return nil, err
		}
	case "gzip":
		payload, err = GZipDeCompress(content[:])
		if err != nil {
			return nil, err
		}
	case "deflate":
		payload, err = ZlibDeCompress(content[:])
		if err != nil {
			return nil, err
		}
	}
	return payload, nil
}

func CompressHTTPBody(r *http.Request, header http.Header, body []byte) ([]byte, error) {
	compress := r.Header.Get("Accept-Encoding")
	var payload []byte
	var err error
	switch compress {
	case "lz4":
		payload, err = LZ4Compress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "lz4")
	case "gzip":
		payload, err = GZipCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "gzip")
	case "deflate":
		payload, err = ZlibCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "deflate")
	default:
		payload = body
	}
	return payload, nil
}
