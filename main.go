package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

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

	peers, err := torrentFile.GetPeers()
	if err != nil {
		log.Fatalf("Error fetching peers: %v", err)
	}

	fmt.Println("Peers List:")
	for _, peer := range peers {
		fmt.Println(peer)
	}

	fmt.Println("Connecting to peers...")

	infoHash, _ := torrentFile.ComputeInfoHash()
	peerID := "-MY1000-123456789012"

	// Connect to first 5 peers
	for i, peerAddr := range peers {
		if i >= 20 {
			break
		}

		host, port, err := net.SplitHostPort(peerAddr)
		if err != nil {
			fmt.Printf("Invalid peer address %s: %v\n", peerAddr, err)
			continue
		}

		portInt, err := strconv.Atoi(port)
		if err != nil {
			fmt.Printf("Invalid port number %s: %v\n", port, err)
			continue
		}

		peer := torrent.Peer{
			IP:   host,
			Port: portInt,
		}

		err = torrent.ConnectToPeer(peer, infoHash, peerID)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v\n", peerAddr, err)
		}
	}
}