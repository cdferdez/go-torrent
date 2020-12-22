package main

import (
	"log"
	"os"

	"github.com/cdferdez/go-torrent/torrentfile"
)

func main() {
	inPath := os.Args[1]
	//outPath := os.Args[2]

	tf, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	err = tf.DownloadToFile("test")
}
