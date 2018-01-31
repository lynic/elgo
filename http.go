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
		payload, err = FlateDeCompress(content[:])
		if err != nil {
			return nil, err
		}
	case "zlib":
		payload, err = ZlibDeCompress(content[:])
		if err != nil {
			return nil, err
		}
	case "snappy":
		payload, err = SnappyDeCompress(content[:])
		if err != nil {
			return nil, err
		}
	case "xz":
		payload, err = XzDeCompress(content[:])
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
	case "snappy":
		payload, err = SnappyCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "snappy")
	case "gzip":
		payload, err = GZipCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "gzip")
	case "deflate":
		payload, err = FlateCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "deflate")
	case "zlib":
		payload, err = ZlibCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "zlib")
	case "xz":
		payload, err = XzCompress(body[:])
		if err != nil {
			return nil, err
		}
		header.Set("Content-Encoding", "xz")
	default:
		payload = body
	}
	return payload, nil
}

func CompressHTTP(requestHeader http.Header, respondHeader http.Header, body []byte) ([]byte, error) {
	contentEncoding := respondHeader.Get("Content-Encoding")
	if contentEncoding != "" {
		return body, nil
	}
	acceptEncoding := requestHeader["Accept-Encoding"]
	if len(acceptEncoding) == 0 {
		return body, nil
	}
	var payload []byte
	var err error
	for _, coding := range acceptEncoding {
		switch coding {
		case "lz4":
			payload, err = LZ4Compress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "lz4")
			break
		case "snappy":
			payload, err = SnappyCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "snappy")
			break
		case "gzip":
			payload, err = GZipCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "gzip")
			break
		case "deflate":
			payload, err = FlateCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "deflate")
		case "zlib":
			payload, err = ZlibCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "zlib")
			break
		case "xz":
			payload, err = XzCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "xz")
			break
		}
	}
	respondHeader.Del("Content-Type")
	respondHeader.Del("Content-Length")
	return payload, nil
}
