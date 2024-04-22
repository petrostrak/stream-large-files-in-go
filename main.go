package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
}

func (fs *FileServer) Start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.read(conn)
	}
}

func (fs *FileServer) read(conn net.Conn) {
	buffer := new(bytes.Buffer)
	for {
		n, err := io.CopyN(buffer, conn, 4096)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buffer.Bytes())
		fmt.Printf("received %d bytes over the network.\n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("received %d bytes over the network.\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		err := sendFile(4096)
		if err != nil {
			return
		}
	}()
	server := &FileServer{}
	server.Start()
}
