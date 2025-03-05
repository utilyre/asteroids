package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
)

func HandleConn(r io.Reader, conn net.Conn) {
	defer conn.Close()

	for {
		slog.Info("reading message from network", "remote", conn.RemoteAddr())

		msg, err := ReadMessage(r)
		if errors.Is(err, io.EOF) {
			slog.Info("connection closed", "remote", conn.RemoteAddr())
			break
		}
		if err != nil {
			slog.Error("failed to read message from connection",
				"remote", conn.RemoteAddr(),
				"error", err,
			)
			return
		}

		slog.Info("received message",
			"remote", conn.RemoteAddr(),
			"message", msg,
		)
	}
}

type Message struct {
	Version uint64
	Scope   string
	Body    []byte
}

func ReadMessage(r io.Reader) (*Message, error) {
	version, err := recvUInt64(r)
	if err != nil {
		return nil, err
	}

	scope, err := recvBytes(r)
	if err != nil {
		return nil, err
	}

	body, err := recvBytes(r)
	if err != nil {
		return nil, err
	}

	return &Message{
		Version: version,
		Scope:   string(scope),
		Body:    body,
	}, nil
}
