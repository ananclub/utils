package utils

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	ActionUrl string `json:"actionURL"` //点击按钮触发的URL
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
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"`
}
type DingMSG struct {
	MsgType    string            `json:"msgtype"`
	Text       DingMSGText       `json:"text"`
	Link       DingMSGLink       `json:"link"`
	Markdown   DingMSGMarkdown   `json:"markdown"`
	ActionCard DingMSGActionCard `json:"actionCard"`
	FeedCard   DingMSGFeedCard   `json:"feedCard"`
	At         DingMSGAt         `json:"at"`
}

type DingTalkRobot struct {
	AccessToken string
	Secret      string
	Url         string
	Vals        url.Values
	ChanSend    chan DingMSG
	LastActive  int64
}

func NewDingTalkRobot(token, secret, apiurl string) (r *DingTalkRobot, err error) {
	vals := url.Values{}
	if secret != "" {
		var sign string
		timestamp := strconv.Itoa(int(time.Now().Unix()) * 1000)
		plainText := timestamp + "\n" + secret
		sign, err = HmacSha256(plainText, secret)
		if err != nil {
			return
		}
		vals.Add("timestamp", timestamp)
		vals.Add("sign", sign)
	}
	vals.Add("access_token", token)
	if apiurl == "" {
		apiurl = "https://oapi.dingtalk.com/robot/send"
	}
	r = &DingTalkRobot{
		AccessToken: token,
		Secret:      secret,
		Url:         apiurl,
		Vals:        vals,
		LastActive:  time.Now().Unix(),
	}

	r.ChanSend = make(chan DingMSG, 1024)

	go func(r *DingTalkRobot) {
		for {
			msg, ok := <-r.ChanSend
			if ok {
				if content, err := json.Marshal(msg); err != nil {
					log.Println("json.Marshal err:", err.Error(), content)
				} else if _, b, err := HttpDo(r.Url, "POST", nil, nil, map[string]string{"Content-Type": "application/json; charset=utf-8"}, r.Vals, content); err != nil {
					log.Println("HttpDo err:", err.Error())
				} else {
					log.Println("HttpDo:", string(b))
				}
				r.LastActive = time.Now().Unix()
				time.Sleep(3 * time.Second)
			} else {
				break
			}
		}
	}(r)

	return
}

func (r *DingTalkRobot) SendText(content string, atmobiles []string, atuserids []string, isatall bool) {
	/*
		发送文本信息
		content: 文本内容
		atmobiles: 需要@的手机号列表
		atuserid: 需要@的用户id列表
		isatall: 是否需要@全体成员
	*/

	msg := DingMSG{
		MsgType: DINGMSG_TYPE_TEXT,
		Text:    DingMSGText{Content: content},
		At: DingMSGAt{
			AtMobiles: atmobiles,
			AtUserIds: atuserids,
			IsAtAll:   isatall,
		},
	}
	if len(msg.At.AtMobiles) > 0 {
		for _, p := range msg.At.AtMobiles {
			if len(p) > 0 && !strings.Contains(msg.Text.Content, "@"+p) {
				msg.Text.Content += " @" + p
			}
		}
	}
	if len(msg.At.AtUserIds) > 0 {
		for _, p := range msg.At.AtUserIds {
			if len(p) > 0 && !strings.Contains(msg.Text.Content, "@"+p) {
				msg.Text.Content += " @" + p
			}
		}
	}
	r.ChanSend <- msg
	return
}

func (r *DingTalkRobot) SendLink(title string, text string, messageurl string, picurl string) {
	/*
		发送链接信息
		title: 标题
		text: 文本内容
		messageurl: 链接URL
		picurl: 图片地址
	*/

	msg := DingMSG{
		MsgType: DINGMSG_TYPE_LINK,
		Link: DingMSGLink{
			MessageUrl: messageurl,
			PicUrl:     picurl,
			Title:      title,
			Text:       text,
		},
	}
	r.ChanSend <- msg
	return
}

func (r *DingTalkRobot) SendMarkdown(title string, text string, atmobiles []string, atuserids []string, isatall bool) {

	msg := DingMSG{
		MsgType: DINGMSG_TYPE_MARKDOWN,
		Markdown: DingMSGMarkdown{
			Title: title,
			Text:  text,
		},
		At: DingMSGAt{
			AtMobiles: atmobiles,
			AtUserIds: atuserids,
			IsAtAll:   isatall,
		},
	}
	if len(msg.At.AtMobiles) > 0 {
		for _, p := range msg.At.AtMobiles {
			if len(p) > 0 && !strings.Contains(msg.Text.Content, "@"+p) {
				msg.Text.Content += " @" + p
			}
		}
	}
	if len(msg.At.AtUserIds) > 0 {
		for _, p := range msg.At.AtUserIds {
			if len(p) > 0 && !strings.Contains(msg.Text.Content, "@"+p) {
				msg.Text.Content += " @" + p
			}
		}
	}
	r.ChanSend <- msg

	return

}

func (r *DingTalkRobot) SendFeedcard(title []string, messageurl []string, picurl []string) {
	/*
		发送feedcard类型消息
		title:单条信息文本 ["时代的火车向前开1",时代的火车向前开2]
		messageURL:点击单条信息到跳转链接 ["https://www.dingtalk.com/","https://www.dingtalk.com/"]
		picURL:单条信息后面图片的URL ["https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png","https://img.alicdn.com/tfs/TB1NwmBEL9TBuNjy1zbXXXpepXa-2400-1218.png"]

	*/
	msg := DingMSG{
		MsgType: DINGMSG_TYPE_FEEDCARD,
		FeedCard: DingMSGFeedCard{
			Links: []DingMSGFeedCardLink{},
		},
	}
	if len(title) == len(messageurl) && len(messageurl) == len(picurl) {
		length := len(title)
		for i := 0; i < length; i++ {

			var link DingMSGFeedCardLink
			link.Title = title[i]
			link.MessageURL = messageurl[i]
			link.PicURL = picurl[i]
			msg.FeedCard.Links = append(msg.FeedCard.Links, link)
		}
	} else {
		log.Println("标题与文章链接数量不一致")
	}
	r.ChanSend <- msg
	return
}

func (r *DingTalkRobot) SendActioncard(title string, text string, titlechild []string, actionurl []string, btnorientation string) {
	/*
		发送整体卡片
		title: 首屏会话透出的展示内容。
		text:markdown格式的消息。
		btns:按钮
		title_child:按钮标题
		actionURL:点击按钮触发的URL
		btnOrientation:0：按钮竖直排列  1：按钮横向排列 默认为0
	*/

	if btnorientation != "1" {
		btnorientation = "0"
	}
	msg := DingMSG{
		MsgType: DINGMSG_TYPE_ACTIONCARD,
		ActionCard: DingMSGActionCard{
			Title:          title,
			Text:           text,
			BtnOrientation: btnorientation,
			Btns:           []DingMSGActionCardBtn{},
		},
	}
	if len(titlechild) == len(actionurl) {
		length := len(titlechild)
		for i := 0; i < length; i++ {
			var btn DingMSGActionCardBtn
			btn.Title = titlechild[i]
			btn.ActionUrl = actionurl[i]
			msg.ActionCard.Btns = append(msg.ActionCard.Btns, btn)
		}
	} else {
		log.Println("标题与文章链接数量不一致")
	}
	r.ChanSend <- msg
	return
}

func (r *DingTalkRobot) SendWholeActioncard(title string, text string, singtitle string, singleurl string, btnorientation string) {
	/*
		发送独立卡片
		title:首屏会话透出的展示内容。
		text:markdown格式的消息。
		singleTitle:单个按钮的标题。注意 设置此项和singleURL后，btns无效。
		singleURL:点击singleTitle按钮触发的URL。
		btnOrientation:0：按钮竖直排列  1：按钮横向排列 默认为0
	*/
	if btnorientation != "1" {
		btnorientation = "0"
	}
	msg := DingMSG{
		MsgType: DINGMSG_TYPE_ACTIONCARD,
		ActionCard: DingMSGActionCard{
			Title:          title,
			Text:           text,
			SingleTitle:    singtitle,
			SingleURL:      singleurl,
			BtnOrientation: btnorientation,
		},
	}
	r.ChanSend <- msg
	return
}

/*
func SendRequest(webhook string, params []byte) {
	reader := bytes.NewReader(params)
	request, err := http.NewRequest("POST", webhook, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 设置请求头及代理
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	// 发送请求
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
}
*/
