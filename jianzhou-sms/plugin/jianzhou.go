package plugin

import (
	"errors"
	"net/url"

	"github.com/livexy/pkg/request"

	"github.com/livexy/plugins/plugin/smser"

	"github.com/livexy/pkg/strx"
	"github.com/livexy/pkg/template"
)

func NewSMS(cfg smser.SMSConfig) smser.SMSer {
	sms := &SMS{cfg: cfg}
	return sms
}

type SMS struct {
	cfg smser.SMSConfig
}

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
	api := "http://www.jianzhou.sh.cn/JianzhouSMSWSServer/http/sendBatchMessage"
	params := "account=" + o.cfg.AccessID + "&password=" + o.cfg.AccessSecret + "&destmobile=" + mobile + "&msgText=" + body + "&sendDateTime="
	result, err := request.HttpPost(api, params, 10, request.Header{Name: "Content-Type", Value: "application/x-www-form-urlencoded"})
	if err != nil {
		return err
	}
	iret := strx.ToInt(result)
	if iret > 0 {
		return nil
	}
	switch iret {
	case -1:
		return errors.New("短信余额不足")
	case -2:
		return errors.New("短信帐号或密码错误")
	case -3:
		return errors.New("连接服务商失败")
	case -4:
		return errors.New("短信发送超时")
	case -5:
		return errors.New("其他错误，一般为网络问题，IP受限等")
	case -6:
		return errors.New("短信内容为空")
	case -7:
		return errors.New("目标号码为空")
	case -11:
		return errors.New("超过最大定时时间限制")
	case -12:
		return errors.New("目标号码在黑名单里")
	case -13:
		return errors.New("没有权限使用该网关")
	case -22:
		return errors.New("Ip被封停")
	default:
		return errors.New("错误编号：" + result)
	}
}
