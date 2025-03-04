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

		msg, err := ReadMessage(conn)
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
	Entity  string
	Action  string
}

func ReadMessage(r io.Reader) (*Message, error) {
	version, err := recvUInt64(r)
	if err != nil {
		return nil, err
	}

	entity, err := recvString(r)
	if err != nil {
		return nil, err
	}

	action, err := recvString(r)
	if err != nil {
		return nil, err
	}

	return &Message{
		Version: version,
		Entity:  entity,
		Action:  action,
	}, nil
}
