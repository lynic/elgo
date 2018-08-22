package elgo

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/json-iterator/go"
)

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	eauth := base64.StdEncoding.EncodeToString([]byte(auth))
	return fmt.Sprintf("Basic %s", eauth)
}

func DigestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func StrIn(s []string, ss string) bool {
	for _, v := range s {
		if v == ss {
			return true
		}
	}
	return false
}

func StrRemove(s []string, ss string) []string {
	for i, v := range s {
		if v == ss {
			s[i] = s[len(s)-1]
			return s[:len(s)-1]
		}
	}
	return s
}

func CheckOS() string {
	return runtime.GOOS
}

func ToJson(v interface{}) ([]byte, error) {
	out, err := jsoniter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func ToPrettyJson(v interface{}) ([]byte, error) {
	out, err := jsoniter.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return out, nil
}

func FromJson(data []byte, v interface{}) error {
	err := jsoniter.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func TrimSplit(s, seq string) []string {
	splited := strings.Split(s, seq)
	ret := make([]string, len(splited))
	for i, v := range splited {
		ret[i] = strings.TrimSpace(v)
	}
	return ret
}

func SaveStruct(v interface{}, filename string) error {
	out, err := ToJson(v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveStructCompress(v interface{}, filename string) error {
	out, err := ToJson(v)
	if err != nil {
		return err
	}
	cout, err := GZipCompress(out)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, cout, 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveStructPretty(v interface{}, filename string) error {
	out, err := ToPrettyJson(v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadStruct(v interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return err
}

func LoadStructCompress(v interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	cdata, err := GZipDeCompress(data)
	if err != nil {
		return err
	}
	err = jsoniter.Unmarshal(cdata, v)
	if err != nil {
		return err
	}
	return err
}

func GenString(n uint, letters string) string {
	ss := ""
	llen := len(letters)
	var i uint = 0
	rand.Seed(time.Now().UnixNano())
	for ; i < n; i++ {
		ss += string(letters[rand.Intn(llen)])
	}
	return ss
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func EnsurePath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func CIDRToPool(cidr string) map[string]string {
	incip := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}
	ips := make(map[string]string)
	ip, ipnet, _ := net.ParseCIDR(cidr)
	ipa := ip.Mask(ipnet.Mask)
	// skip 0 and 1
	incip(ipa)
	ips[ipa.String()] = "router"
	incip(ipa)
	for ; ipnet.Contains(ipa); incip(ipa) {
		ips[ipa.String()] = ""
	}
	// skip last one
	delete(ips, ipa.String())
	return ips
}

func DecodeAuthHeader(authorization string) (string, string, error) {
	auth := strings.SplitN(authorization, " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", "", fmt.Errorf("Invalid header")
	}
	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return "", "", err
	}
	pairs := strings.SplitN(string(payload), ":", 2)
	if len(pairs) != 2 {
		return "", "", fmt.Errorf("Invalid header")
	}
	return pairs[0], pairs[1], nil
}

func GenAuthorization(name, password string) string {
	mstr := fmt.Sprintf("%s:%s", name, password)
	return base64.StdEncoding.EncodeToString([]byte(mstr))
}

// func AuthRequest(authorization string) (*AuthContext, error) {
// 	username, password, err := DecodeAuthHeader(authorization)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return AuthUserRequest(username, password)
// }

func MakeError(msg string) []byte {
	ret := fmt.Sprintf(`{
		"errors": [{
			"message": "%s"
		}]
	}`, msg)
	return []byte(ret)
}
