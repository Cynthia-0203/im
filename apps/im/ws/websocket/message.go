package websocket

import "time"


type FrameType uint8

const (
	FrameData      FrameType = 0x0
	FramePing      FrameType = 0x1
	FrameAck       FrameType = 0x2
	FrameNoAck     FrameType = 0x3
	FrameErr       FrameType = 0x9
	FrameTranspond FrameType = 0x6
)

type Message struct{
	FrameType `json:"frameType"`
	Id        string      `json:"id"`
	Method string `json:"method,omitempty"`
	// UserId string `json:"userId,omitempty"`
	FormId string `json:"formId,omitempty"`
	Data interface{} `json:"data,omitempty"`
	errCount int `json:"-"`
	ackTime time.Time `json:"-"`
	AckSeq int `json:"ackSeq,omitempty"`
	TranspondUid string      `json:"transpondUid"`
}

func NewMessage(formId string,data interface{})*Message{

	return &Message{
		FrameType: FrameData,
		FormId: formId,
		Data: data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
