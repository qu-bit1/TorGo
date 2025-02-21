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

// establishes a TCP connection and performs a handshake
func ConnectToPeer(peer Peer, infoHash [20]byte, peerID string) error {
    address := fmt.Sprintf("%s:%d", peer.IP, peer.Port) // Ensure no extra colons
    conn, err := net.DialTimeout("tcp", address, 5*time.Second)
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

    // Read response
    resp := make([]byte, 68)
    _, err = conn.Read(resp)
    if err != nil {
        return fmt.Errorf("failed to read handshake response: %v", err)
    }

    fmt.Printf("Handshake successful with peer %s\n", address)
    return nil
}