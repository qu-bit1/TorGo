package main

import (
	"fmt"
	"log"

	"github.com/qu-bit1/TorGo/torrent"
)

func main() {
	torrentFile, err := torrent.ParseTorrentFile("debian-mac-12.9.0-amd64-netinst.iso.torrent")
	if err != nil {
		log.Fatalf("error parsing torrent file: %v", err)
	}

	fmt.Println("Tracker URL: ", torrentFile.Announce)
	fmt.Println("File name: ", torrentFile.Info.Name)
	fmt.Println("File length: ", torrentFile.Info.Length)
	fmt.Println("Piece length: ", torrentFile.Info.PieceLength)
	// fmt.Println("Number of pieces: ", len(torrentFile.Info.Pieces) / 20)

	pieces, err := torrentFile.DecodePieces()
	if err != nil {
		log.Fatalf("error decoding pieces: %v", err)
	}
	fmt.Println("Number of Pieces: ", len(pieces))
}
