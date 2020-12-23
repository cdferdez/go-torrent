package message

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	// MsgChoke chokes the receiver
	MsgChoke messageID = iota
	// MsgUnchoke unchokes the receiver
	MsgUnchoke messageID = iota
	// MsgInterested expresses interest in receiving data
	MsgInterested messageID = iota
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested messageID = iota
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave messageID = iota
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield messageID = iota
	// MsgRequest requests a block of data from the receiver
	MsgRequest messageID = iota
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece messageID = iota
	// MsgCancel cancels a request
	MsgCancel messageID = iota
)

// Message stores ID and paylaod of a message
type Message struct {
	ID      messageID
	Payload []byte
}

// Serialize serializes a message into a buffer of the form
// <length prefix><message ID><payload>
// Interprets `nil` as a keep-alive message
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload[:])
	return buf
}

// Read parses a Message from a stream. Returns `nil` on keep-alive message
func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep alive
	if length == 0 {
		return nil, nil
	}

	// read id and payload
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}
