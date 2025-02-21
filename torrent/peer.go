package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Peer struct {
	IP   string
	Port int
}

// 68-byte handshake structure
type HandshakeMessage struct {
	Pstrlen  byte
	Pstr     [19]byte
	Reserved [8]byte
	InfoHash [20]byte
	PeerID   [20]byte
}

// func to build the handshake message
func CreateHandshake(infoHash [20]byte, peerID string) ([]byte, error) {
	if len(peerID) != 20 {
		return nil, fmt.Errorf("peerID must be exactly 20 bytes")
	}

	handshake := HandshakeMessage{
		Pstrlen:  19,
		Pstr:     [19]byte{'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't', ' ', 'p', 'r', 'o', 't', 'o', 'c', 'o', 'l'},
		Reserved: [8]byte{},
	}

	copy(handshake.InfoHash[:], infoHash[:])
	copy(handshake.PeerID[:], []byte(peerID))

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, handshake)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// establishes a TCP connection and performs a handshake and handles message exhange
func ConnectToPeer(peer Peer, infoHash [20]byte, peerID string) error {
    address := fmt.Sprintf("%s:%d", peer.IP, peer.Port)
    conn, err := net.DialTimeout("tcp", address, 10*time.Second)
    if err != nil {
        return fmt.Errorf("failed to connect to peer %s: %v", address, err)
    }
    defer conn.Close()

    // Send handshake
    handshake, err := CreateHandshake(infoHash, peerID)
    if err != nil {
        return err
    }

    _, err = conn.Write(handshake)
    if err != nil {
        return fmt.Errorf("failed to send handshake: %v", err)
    }

    // Read handshake response
    resp := make([]byte, 68)
    _, err = conn.Read(resp)
    if err != nil {
        return fmt.Errorf("failed to read handshake response: %v", err)
    }

    fmt.Printf("Handshake successful with peer %s\n", address)

    // Read peer's first response message (usually Bitfield)
	msg, err := ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to read message from peer: %v", err)
	}

	if msg == nil {
		fmt.Println("Error: Received a nil message from", address)
		return err
	}
	if msg.ID == MsgBitfield {
		fmt.Println("Received Bitfield message from", address)
	} else {
		fmt.Println("No Bitfield message received from", address)
	}

	// Send Interested message
	fmt.Println("Sending Interested message...")
	err = SendMessage(conn, MsgInterested, nil)
	if err != nil {
		return fmt.Errorf("failed to send interested message: %v", err)
	}

	// Read peer response (should be Unchoke)
	msg, err = ReadMessage(conn)
	if err != nil {
		return fmt.Errorf("failed to read message from peer: %v", err)
	}

	if msg.ID == MsgUnchoke {
		fmt.Println("Peer unchoked us! Requesting first piece...")

    // Select the first piece (later, optimize piece selection)
    pieceIndex := 0
    blockOffset := 0
    blockSize := 16384 // 16 KB blocks

    // Send a request for the first block of the first piece
    requestPayload := make([]byte, 12)
    binary.BigEndian.PutUint32(requestPayload[0:4], uint32(pieceIndex))
    binary.BigEndian.PutUint32(requestPayload[4:8], uint32(blockOffset))
    binary.BigEndian.PutUint32(requestPayload[8:12], uint32(blockSize))

    err = SendMessage(conn, MsgRequest, requestPayload)
    if err != nil {
        return fmt.Errorf("failed to send request: %v", err)
    }

    fmt.Println("Request sent for piece", pieceIndex, "block", blockOffset)
	} else {
		fmt.Println("Peer did not unchoke us.")
	}

    return nil
}