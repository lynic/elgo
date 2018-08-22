package elgo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ReadHTTPBody(r *http.Request) ([]byte, error) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	encoding := r.Header.Get("Content-Encoding")
	return ReadHTTP(content, encoding)
}

func ReadHTTP(content []byte, encoding string) ([]byte, error) {
	if encoding == "" {
		return content, nil
	}
	var payload []byte
	var err error
	payload = content[:]
	switch encoding {
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
	acceptEncoding := requestHeader.Get("Accept-Encoding")
	if acceptEncoding == "" {
		return body, nil
	}
	compressed := false
	var payload []byte
	var err error
	payload = body
	for _, coding := range TrimSplit(acceptEncoding, ",") {
		switch coding {
		case "lz4":
			payload, err = LZ4Compress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "lz4")
			compressed = true
			break
		case "snappy":
			payload, err = SnappyCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "snappy")
			compressed = true
			break
		case "gzip":
			payload, err = GZipCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "gzip")
			compressed = true
			break
		case "deflate":
			payload, err = FlateCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "deflate")
			compressed = true
		case "zlib":
			payload, err = ZlibCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "zlib")
			compressed = true
			break
		case "xz":
			payload, err = XzCompress(body[:])
			if err != nil {
				return nil, err
			}
			respondHeader.Set("Content-Encoding", "xz")
			compressed = true
			break
		}
	}
	if compressed {
		respondHeader.Del("Content-Type")
		respondHeader.Del("Content-Length")
	}
	return payload, nil
}

func DoRequest(method, url string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(fmt.Sprintf("Response code is %d %s", resp.StatusCode, string(data)))
	}
	cdata := data
	if resp.Header.Get("Content-Encoding") != "" {
		cdata, err = ReadHTTP(data, resp.Header.Get("Content-Encoding"))
		if err != nil {
			return nil, err
		}
	}
	return cdata, nil
}

const AuthContextKey = "AuthContext"

func HttpAuthWraper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authURL := os.Getenv("AUTHURL")
		if authURL == "" {
			next.ServeHTTP(w, r)
			return
		}
		authCtx, err := AuthRequest(r.Header.Get("Authorization"))
		if err != nil {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(MakeError(err.Error()))
			return
		}
		// authCtx := &AuthContext{}
		// authCtx.UserID = authUser.ID
		// authCtx.UserName = authUser.Name
		ctx := context.WithValue(r.Context(), AuthContextKey, authCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
		// next.ServeHTTP(w, r)
	})
}
