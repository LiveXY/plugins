package plugin

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/livexy/plugin/smser"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/template"
)

// NewSMS 创建一个新的麦讯通短信服务实例
func NewSMS(cfg smser.SMSConfig) smser.SMSer {
	sms := &SMS{cfg: cfg}
	return sms
}

type SMS struct {
	cfg smser.SMSConfig
}

// Send 发送短信
// 会将模板内容和数据合并，并调用麦讯通接口发送
func (o *SMS) Send(temp, mobile string, data map[string]any) error {
	if len(mobile) < 5 {
		return errors.New("参数错误！")
	}
	subject := "【" + o.cfg.SignName + "】"
	var body string
	var err error
	if data == nil {
		body = temp
	} else {
		body, err = template.FastTemplate(temp, data)
	}
	if err != nil {
		return err
	}
	body = url.QueryEscape(subject + body)
	url := "http://www.mxtong.cn:8080/GateWay/Services.asmx/DirectSend?" + o.cfg.ExtendData + "&Account=" + o.cfg.AccessID + "&Password=" + o.cfg.AccessSecret + "&Phones=%s&Content=%s&SendTime=&SendType=1&PostFixNumber="
	api := fmt.Sprintf(url, mobile, body)
	result, err := getXml(api)
	if err != nil {
		return err
	}
	if !strings.Contains(result, "<RetCode>Sucess</RetCode>") {
		return errors.New("发送短信失败：" + mobile + subject + body)
	}
	return nil
}

func getXml(api string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Host", "www.mxtong.cn:8080")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	data := bytex.ToStr(body)
	if !strings.HasPrefix(data, "<?xml") {
		return "", errors.New("非XML结果" + data)
	}
	return data, nil
}
