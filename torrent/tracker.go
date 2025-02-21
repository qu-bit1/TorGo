package torrent

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jackpal/bencode-go"
)

type TrackerResp struct {
	Interval int   `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (tf *TorrentFile) ComputeInfoHash() ([20]byte, error) {
	var buf strings.Builder
	err := bencode.Marshal(&buf, tf.Info)
	if err != nil {
		return [20]byte{}, err
	}
	hash := sha1.Sum([]byte(buf.String()))
	return hash, nil
}

func (tf *TorrentFile) GetPeers() ([]string, error) {
	infoHash, err := tf.ComputeInfoHash()
	if err != nil {
		return nil, err
	}
	params := url.Values{
		"info_hash": {string(infoHash[:])},
		"peer_id": {"-MY1000-123456789012"}, // a unique 20-byte client ID
		"port": {"6881"},
		"uploaded": {"0"},
		"downloaded": {"0"},
		"left": {fmt.Sprintf("%d", tf.Info.Length)},
		"compact": {"1"}, // requests compact peer list
	}

	trackerURL := fmt.Sprintf("%s?%s", tf.Announce, params.Encode())
	resp, err := http.Get(trackerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// now decode the tracker resp
	var tracker TrackerResp
	err = bencode.Unmarshal(resp.Body, &tracker)
	if err != nil {
		return nil, err
	}

	// now parse the compact peer list (6 byte per peer: 4 bytes ip + 2 byte port)
	peerData := []byte(tracker.Peers)
	var peers []string
	for i := 0; i < len(peerData); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", peerData[i], peerData[i+1], peerData[i+2], peerData[i+3])
		port := binary.BigEndian.Uint16(peerData[i+4 : i+6])
		peers = append(peers, fmt.Sprintf("%s:%d", ip, port))
	}
	return peers, nil
}
