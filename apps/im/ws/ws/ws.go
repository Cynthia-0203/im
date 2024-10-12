package ws

import (
	"gim/pkg/constants"

)



type (
	Msg struct {
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
		MsgId           string            `mapstructure:"msgId"`
		ReadRecords     map[string]string `mapstructure:"readRecords"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}

	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		RecvIds            []string `mapstructure:"recvIds"`
		SendTime           int64  `mapstructure:"sendTime"`
		ReadRecords map[string]string     `mapstructure:"readRecords"`
		constants.MType `mapstructure:"mType"`
		Content         string `mapstructure:"content"`
		MsgId       string                `mapstructure:"msgId"`
		ContentType constants.ContentType `mapstructure:"contentType"`
	}

	MarkRead struct{
		constants.ChatType `mapstructure:"chatType"`
		RecvId string `mapstructure:"recvId"`
		ConversationId string `mapstructure:"conversationId"`
		MsgIds []string `mapstructure:"msgIds"`
	}
)
