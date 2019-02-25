package elgo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type SCNotify struct {
	SCKey string
}

func (s *SCNotify) Init() error {
	scKey := os.Getenv("SCKEY")
	if scKey == "" {
		return fmt.Errorf("no SCKey in env")
	}
	s.SCKey = scKey
	return nil
}

func (s *SCNotify) Send(title, content string) error {
	urlStr := fmt.Sprintf("https://sc.ftqq.com/%s.send", s.SCKey)
	urlStr += fmt.Sprintf("?text=%s", url.QueryEscape(title))
	urlStr += fmt.Sprintf("&desp=%s", url.QueryEscape(content))
	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("respond code is %d %s", resp.StatusCode, string(body))
	}
	return nil
}

func NotifySend(title, content string) error {
	if os.Getenv("SCKEY") != "" {
		sender := &SCNotify{}
		sender.Init()
		err := sender.Send(title, content)
		if err != nil {
			log.Printf("failed to send through SCNotify: %s", err.Error())
		}
	}
	return nil
}
