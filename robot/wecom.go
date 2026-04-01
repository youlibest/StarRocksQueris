/*
 *@author  chengkenli
 *@project StarRocksQueris
 *@package robot
 *@file    wecom
 *@date    2026/03/30
 *@desc    企业微信应用告警功能（使用CorpID+AgentId+Secret）
 */

package robot

import (
	"StarRocksQueris/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	WeComGetTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	WeComSendMsgURL  = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

// WeComAccessToken 企业微信访问令牌响应
type WeComAccessToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// WeComMessage 企业微信消息结构体
type WeComMessage struct {
	ToUser   string          `json:"touser"`
	MsgType  string          `json:"msgtype"`
	AgentID  int             `json:"agentid"`
	Text     *WeComText      `json:"text,omitempty"`
	Markdown *WeComMarkdown  `json:"markdown,omitempty"`
}

// WeComText 文本消息类型
type WeComText struct {
	Content string `json:"content"`
}

// WeComMarkdown Markdown消息类型
type WeComMarkdown struct {
	Content string `json:"content"`
}

// WeComResponse 企业微信响应
type WeComResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var (
	wecomTokenCache *cache.Cache
)

func init() {
	wecomTokenCache = cache.New(7000*time.Second, 100*time.Second)
}

// GetWeComAccessToken 获取企业微信访问令牌
func GetWeComAccessToken(corpID, secret string) (string, error) {
	// 检查缓存
	if token, found := wecomTokenCache.Get("access_token"); found {
		return token.(string), nil
	}

	url := fmt.Sprintf("%s?corpid=%s&corpsecret=%s", WeComGetTokenURL, corpID, secret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp WeComAccessToken
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("获取access_token失败: %s", tokenResp.ErrMsg)
	}

	// 缓存token（有效期7200秒，这里设置7000秒提前刷新）
	wecomTokenCache.Set("access_token", tokenResp.AccessToken, cache.DefaultExpiration)

	return tokenResp.AccessToken, nil
}

// SendWeComAppMessage 发送企业微信应用消息
func SendWeComAppMessage(msg *WeComMessage, corpID, secret string) error {
	accessToken, err := GetWeComAccessToken(corpID, secret)
	if err != nil {
		return fmt.Errorf("获取access_token失败: %v", err)
	}

	url := fmt.Sprintf("%s?access_token=%s", WeComSendMsgURL, accessToken)

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result WeComResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("发送消息失败: %s", result.ErrMsg)
	}

	return nil
}

// SendWeComMarkdown 发送Markdown消息到企业微信
func SendWeComMarkdown(title, message, logURL string) {
	// 构建markdown内容
	var content strings.Builder
	content.WriteString(fmt.Sprintf("## %s\n\n", title))
	content.WriteString(message)
	if logURL != "" {
		content.WriteString(fmt.Sprintf("\n\n[查看日志](%s)", logURL))
	}

	// 构建接收人列表
	var toUser string
	if util.GlobalRuleConfig.WeComApp.MentionAll {
		toUser = "@all"
	} else if len(util.GlobalRuleConfig.WeComApp.MentionedUserList) > 0 {
		toUser = strings.Join(util.GlobalRuleConfig.WeComApp.MentionedUserList, "|")
	} else {
		toUser = "@all"
	}

	// 构建消息
	agentID := 0
	fmt.Sscanf(util.GlobalRuleConfig.WeComApp.AgentID, "%d", &agentID)

	msg := &WeComMessage{
		ToUser:  toUser,
		MsgType: "markdown",
		AgentID: agentID,
		Markdown: &WeComMarkdown{
			Content: content.String(),
		},
	}

	// 发送消息
	err := SendWeComAppMessage(msg, util.GlobalRuleConfig.WeComApp.CorpID, util.GlobalRuleConfig.WeComApp.Secret)
	if err != nil {
		util.Loggrs.Error(uid, fmt.Sprintf("企业微信发送失败: %v", err))
	} else {
		util.Loggrs.Info(uid, "企业微信发送成功")
	}
}

// SendWeComText 发送文本消息到企业微信
func SendWeComText(title, message, logURL string) {
	// 构建文本内容
	var content strings.Builder
	content.WriteString(fmt.Sprintf("【%s】\n\n", title))
	content.WriteString(message)
	if logURL != "" {
		content.WriteString(fmt.Sprintf("\n\n日志链接: %s", logURL))
	}

	// 构建接收人列表
	var toUser string
	if util.GlobalRuleConfig.WeComApp.MentionAll {
		toUser = "@all"
	} else if len(util.GlobalRuleConfig.WeComApp.MentionedUserList) > 0 {
		toUser = strings.Join(util.GlobalRuleConfig.WeComApp.MentionedUserList, "|")
	} else {
		toUser = "@all"
	}

	// 构建消息
	agentID := 0
	fmt.Sscanf(util.GlobalRuleConfig.WeComApp.AgentID, "%d", &agentID)

	msg := &WeComMessage{
		ToUser:  toUser,
		MsgType: "text",
		AgentID: agentID,
		Text: &WeComText{
			Content: content.String(),
		},
	}

	// 发送消息
	err := SendWeComAppMessage(msg, util.GlobalRuleConfig.WeComApp.CorpID, util.GlobalRuleConfig.WeComApp.Secret)
	if err != nil {
		util.Loggrs.Error(uid, fmt.Sprintf("企业微信发送失败: %v", err))
	} else {
		util.Loggrs.Info(uid, "企业微信发送成功")
	}
}

// SendWeComQueris 发送慢查询告警到企业微信
func SendWeComQueris(i *util.InQue, queris bool) error {
	cid := fmt.Sprintf("wecom_%d_%s", i.Action, i.Item.Id)
	_, ok := i.Larkcache.Get(cid)
	util.Loggrs.Info(uid, fmt.Sprintf("企业微信最后关卡：识别缓存中的%s，缓存状态：%t", cid, ok))
	if ok {
		return nil
	}

	// 构建消息内容
	var msgs []string
	msgs = append(msgs, fmt.Sprintf("> **告警类型**: %s\n", getActionName(i.Action)))

	if i.App != "" {
		msgs = append(msgs, fmt.Sprintf("> **集群**: %s\n", i.App))
	}
	if i.Fe != "" {
		msgs = append(msgs, fmt.Sprintf("> **FE节点**: %s\n", i.Fe))
	}
	if i.Item.Host != "" {
		msgs = append(msgs, fmt.Sprintf("> **客户端IP**: %s\n", i.Item.Host))
	}
	if i.Item.User != "" {
		msgs = append(msgs, fmt.Sprintf("> **用户**: %s\n", i.Item.User))
	}
	if i.Item.Id != "" {
		msgs = append(msgs, fmt.Sprintf("> **连接ID**: %s\n", i.Item.Id))
	}
	if i.Item.Time != "" {
		msgs = append(msgs, fmt.Sprintf("> **执行时间**: %s秒\n", i.Item.Time))
	}
	if i.Item.State != "" {
		msgs = append(msgs, fmt.Sprintf("> **状态**: %s\n", i.Item.State))
	}

	// SQL语句处理
	var sessionSQL string
	if len(i.Item.Info) >= 300 {
		sessionSQL = i.Item.Info[0:280] + " ..."
	} else {
		sessionSQL = i.Item.Info
	}
	sessionSQL = strings.NewReplacer("\n", "", `"`, `"`).Replace(sessionSQL)
	if sessionSQL != "" {
		msgs = append(msgs, fmt.Sprintf("> **SQL**: `%s`\n", sessionSQL))
	}

	// 构建标题
	title := fmt.Sprintf("StarRocks慢查询告警 - %s", getActionName(i.Action))

	// 构建日志链接
	logURL := ""
	if i.Logfile != "" {
		logURL = fmt.Sprintf("http://%s:7890/log%s", util.H.Ip, i.Logfile)
	}

	// 根据配置选择消息类型
	msgType := util.GlobalRuleConfig.WeComApp.MsgType
	if msgType == "text" {
		SendWeComText(title, strings.Join(msgs, ""), logURL)
	} else {
		// 默认使用markdown
		SendWeComMarkdown(title, strings.Join(msgs, ""), logURL)
	}

	// 写入缓存
	i.Larkcache.Set(cid, i.Item.Id, cache.DefaultExpiration)
	return nil
}

// getActionName 获取动作名称
func getActionName(action int) string {
	names := map[int]string{
		0: "状态异常停留清退",
		1: "异常违规参数查杀",
		2: "10分钟慢查询提醒",
		3: "30分钟慢查询查杀",
		4: "全表扫描亿级查杀",
		5: "TB级扫描字节查杀",
		6: "百亿扫描行数查杀",
		7: "CATALOG违规查杀",
		8: "GB级消耗内存查杀",
		9: "其他",
	}
	if name, ok := names[action]; ok {
		return name
	}
	return "未知告警"
}
