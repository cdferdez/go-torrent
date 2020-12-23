package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/cdferdez/go-torrent/handshake"
	"github.com/cdferdez/go-torrent/peers"
)

// A Client is a  TCP connection with a peer
type Client struct {
	Conn     net.Conn
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func completeHandshake(conn net.Conn, infohash, peerID [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{})

	req := handshake.New(infohash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], infohash[:]) {
		return nil, fmt.Errorf("Expected infohash %x but got %x", res.InfoHash, infohash)
	}

	return res, nil
}

// New creates a new connection, completes handshake
func New(peer peers.Peer, infohash, peerID [20]byte) (*Client, error) {
	fmt.Println("Establishing Connection to peer:", peer.String())
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("------Completing Handshake------")
	_, err = completeHandshake(conn, infohash, peerID)
	if err != nil {
		conn.Close()
		fmt.Print(err)
		return nil, err
	}
	fmt.Println("Success!")

	return &Client{
		Conn:     conn,
		peer:     peer,
		infoHash: infohash,
		peerID:   peerID,
	}, nil
}
