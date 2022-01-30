package client

type FrameType uint8

const (
	SystemFrameType FrameType = 0
	MessageFrameType FrameType = 1
	RsaHandshakeFrameType FrameType = 2
	PartnerDisconnectedFrameType FrameType = 3
	AesHandshakeFrameType FrameType = 4
)

type Frame struct {
	From string `json:"from"`
	Data []byte `json:"data"`
	Type FrameType `json:"type"`
}

func NewFrame(from string, data []byte, _type FrameType) *Frame {
	return &Frame{
		From: from,
		Data: data,
		Type: _type,
	}
}
