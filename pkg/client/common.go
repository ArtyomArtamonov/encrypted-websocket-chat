package client

type FrameType uint8

const (
	SystemFrameType FrameType = 0
	MessageFrameType FrameType = 1
)

type Frame struct {
	From string `json:"from"`
	Data interface{} `json:"data"`
	Type FrameType `json:"type"`
}

func NewFrame(from string, data interface{}, _type FrameType) *Frame {
	return &Frame{
		From: from,
		Data: data,
		Type: _type,
	}
}
