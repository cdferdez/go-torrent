package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces       string `bencode:"pieces"`
	PiecesLength int    `bencode:"piece length"`
	Length       int    `bencode:"length"`
	Name         string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// TorrentFile is ...
type TorrentFile struct {
	Announce     string
	InfoHash     [20]byte
	PiecesHashes [][20]byte
	PieceLength  int
	Length       int
	Name         string
}

// Open parses a torrent file
func Open(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}

	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)

	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("received malformed pieces of length %d", len(buf))
		return nil, err
	}

	numHashes := len(buf)
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}

	pieceHashes, err := bto.Info.splitPieceHashes()

	t := TorrentFile{
		Announce:     bto.Announce,
		InfoHash:     infoHash,
		PiecesHashes: pieceHashes,
		PieceLength:  bto.Info.PiecesLength,
		Length:       bto.Info.Length,
		Name:         bto.Info.Name,
	}

	return t, nil
}
