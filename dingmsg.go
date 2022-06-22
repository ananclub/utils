package utils

import (
	"net/url"
	"strconv"
	"time"
)

const (
	DINGAPIURL = "https://oapi.dingtalk.com/robot/send"
)

const (
	DINGMSG_TYPE_TEXT       = "text"       //text
	DINGMSG_TYPE_MARKDOWN   = "markdown"   //markdown
	DINGMSG_TYPE_LINK       = "link"       //link类型
	DINGMSG_TYPE_FEEDCARD   = "feedCard"   //FeedCard
	DINGMSG_TYPE_ACTIONCARD = "actionCard" //整体跳转&独立跳转
)

type DingMSGText struct {
	Content string `json:"content"`
}
type DingMSGMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingMSGActionCardBtn struct {
	Title     string `json:"title"`
	ActionUrl string `json:"actionUrl"` //点击按钮触发的URL
}
type DingMSGLink struct {
	Title      string `json:"title"`      //消息标题
	Text       string `json:"text"`       //消息内容。如果太长只会部分展示
	MessageUrl string `json:"messageUrl"` //点击消息跳转的URL
	PicUrl     string `json:"picUrl"`     //图片URL
}
type DingMSGActionCard struct {
	Title          string                 `json:"title"`          //首屏会话透出的展示内容
	Text           string                 `json:"text"`           //markdown格式的消息
	SingleTitle    string                 `json:"singleTitle"`    //单个按钮的方案。(设置此项和singleURL后btns无效)
	SingleURL      string                 `json:"singleURL"`      //点击singleTitle按钮触发的URL
	BtnOrientation string                 `json:"btnOrientation"` //0-按钮竖直排列，1-按钮横向排列
	HideAvatar     string                 `json:"hideAvatar"`     //0-正常发消息者头像，1-隐藏发消息者头像
	Btns           []DingMSGActionCardBtn `json:"btns"`           //按钮的信息：title-按钮方案，actionURL-点击按钮触发的URL
}
type DingMSGFeedCardLink struct {
	Title      string `json:"title"`      //单条信息文本
	MessageURL string `json:"messageURL"` //点击单条信息到跳转链接
	PicURL     string `json:"picURL"`     //单条信息后面图片的URL
}
type DingMSGFeedCard struct {
	Links []DingMSGFeedCardLink `json:"links"`
}
type DingMSGAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}
type DingMSG struct {
	MsgType    string            `json:"msgtype"`
	Text       DingMSGText       `json:"text"`
	Markdown   DingMSGMarkdown   `json:"markdown"`
	ActionCard DingMSGActionCard `json:"actionCard"`
	FeedCard   DingMSGFeedCard   `json:"feedCard"`
	At         DingMSGAt         `json:"at"`
}

func genSign(token string, sec string) (v url.Values, err error) {
	timestamp := strconv.Itoa(int(time.Now().Unix()) * 1000)
	plainText := timestamp + "\n" + sec
	sign, err := HmacSha256(plainText, sec)
	if err != nil {
		return
	}
	v = url.Values{}
	v.Add("access_token", token)
	v.Add("timestamp", timestamp)
	v.Add("sign", sign)
	return
}
func SendDingMsg(token string, sec string, msg *DingMSG) (result string, err error) {

	if len(msg.At.AtMobiles) > 0 {
		for _, p := range msg.At.AtMobiles {
			if msg.MsgType == DINGMSG_TYPE_MARKDOWN {
				msg.Markdown.Text += "@" + p + " "
			} else if msg.MsgType == DINGMSG_TYPE_TEXT {
				msg.Text.Content += "@" + p + " "
			}
		}
	}
	if v, err := genSign(token, sec); err != nil {
		return "", err
	} else if content, err := json.Marshal(msg); err != nil {
		return "", err
	} else if _, b, err := HttpDo(DINGAPIURL, "POST", nil, nil, map[string]string{"Content-Type": "application/json; charset=utf-8"}, v, content); err != nil {
		return "", err
	} else {
		result = string(b)
	}
	return
}

func SendDingMsgText(token string, sec string, text string) (result string, err error) {

	msg := DingMSG{MsgType: DINGMSG_TYPE_TEXT,
		Text: DingMSGText{Content: text}}
	if v, err := genSign(token, sec); err != nil {
		return "", err
	} else if content, err := json.Marshal(msg); err != nil {
		return "", err
	} else if _, b, err := HttpDo(DINGAPIURL, "POST", nil, nil, map[string]string{"Content-Type": "application/json; charset=utf-8"}, v, content); err != nil {
		return "", err
	} else {
		result = string(b)
	}
	return
}
func SendDingMsgStr(token string, sec string, msg []byte) (result string, err error) {

	if v, err := genSign(token, sec); err != nil {
		return "", err
	} else if _, b, err := HttpDo(DINGAPIURL, "POST", nil, nil, map[string]string{"Content-Type": "application/json; charset=utf-8"}, v, msg); err != nil {
		return "", err
	} else {
		result = string(b)
	}
	return
}
