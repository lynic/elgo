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
	SCKey    string
	UseHttps bool
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
	urlStr := fmt.Sprintf("sc.ftqq.com/%s.send", s.SCKey)
	if s.UseHttps {
		urlStr = "https://" + urlStr
	} else {
		urlStr = "http://" + urlStr
	}
	data, err := url.ParseQuery("text=" + url.QueryEscape(title) + "&desp=" + url.QueryEscape(content))
	if err != nil {
		return err
	}
	resp, err := http.PostForm(urlStr, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
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
