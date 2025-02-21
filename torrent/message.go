package torrent

import (
	"encoding/binary"
	"fmt"
	"io"
)

// BitTorrent Message IDs
const (
	MsgChoke         = 0
	MsgUnchoke       = 1
	MsgInterested    = 2
	MsgNotInterested = 3
	MsgHave          = 4
	MsgBitfield      = 5
	MsgRequest       = 6
	MsgPiece         = 7
	MsgCancel        = 8
)

// Message represents a BitTorrent protocol message
type Message struct {
	ID   byte
	Data []byte
}

// ReadMessage reads a message from the peer
func ReadMessage(conn io.Reader) (*Message, error) {
	// Read message length (first 4 bytes)
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read message length: %v", err)
	}

	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, nil // Keep-alive message
	}

	// Read message ID (1 byte)
	msgID := make([]byte, 1)
	_, err = io.ReadFull(conn, msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to read message ID: %v", err)
	}

	// Read payload (length - 1 bytes)
	payload := make([]byte, length-1)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload: %v", err)
	}

	return &Message{
		ID:   msgID[0],
		Data: payload,
	}, nil
}

// SendMessage sends a BitTorrent message to a peer
func SendMessage(conn io.Writer, msgID byte, payload []byte) error {
	length := uint32(len(payload) + 1)
	buf := make([]byte, 4+1+len(payload))

	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = msgID
	copy(buf[5:], payload)

	_, err := conn.Write(buf)
	return err
}