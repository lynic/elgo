package elgo

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"runtime"
	"time"
)

func CheckOS() string {
	return runtime.GOOS
}

func ToJson(v interface{}) ([]byte, error) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return out, nil
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

func LoadStruct(v interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
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
