package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

// structure of a torrent file
type TorrentFile struct {
	Announce string `bencode:"announce"`
	Info InfoDict `bencode:"info"`
}

// details of the file
type InfoDict struct {
	PieceLength int `bencode:"piece length"`
	Pieces string `bencode:"pieces"`
	Length int `bencode:"length"`
	Name string `bencode:"name"`
}

func ParseTorrentFile(filename string) (*TorrentFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var torrent TorrentFile
	err = bencode.Unmarshal(file, &torrent)
	if err != nil {
		return nil, err
	}
	return &torrent, nil
}

func (tf *TorrentFile) DecodePieces() ([][20]byte, error) {
	buf := bytes.NewBuffer(([]byte)(tf.Info.Pieces))
	numPieces := len(buf.Bytes()) / 20

	if numPieces * 20 != len(tf.Info.Pieces) {
		return nil, fmt.Errorf("invalid piece hash length")
	}

	hashes := make([][20]byte, numPieces)
	for i := 0; i < numPieces; i++ {
		binary.Read(buf, binary.BigEndian, &hashes[i])
	}
	return hashes, nil
}	