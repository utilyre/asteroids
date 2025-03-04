package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("listening on tcp localhost:3000")
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		slog.Error("failed to listen", "error", err)
	}

	for {
		slog.Info("waiting to accept a connection")
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err)
			continue
		}
		slog.Info("accepted connection", "remote", conn.RemoteAddr())

		monitorBuffer := &bytes.Buffer{}
		connReader := io.TeeReader(conn, monitorBuffer)

		go monitorConn(monitorBuffer, conn.RemoteAddr().String())
		go handleConn(connReader, conn)
	}
}

func handleConn(r io.Reader, conn net.Conn) {
	defer conn.Close()

	for {
		slog.Info("receiving string from network", "remote", conn.RemoteAddr())

		msg, err := recvString(r)
		if errors.Is(err, io.EOF) {
			slog.Info("connection closed", "remote", conn.RemoteAddr())
			break
		}
		if err != nil {
			slog.Error("failed to read from connection",
				"remote", conn.RemoteAddr(),
				"error", err,
			)
			return
		}

		slog.Info("received message",
			"remote", conn.RemoteAddr(),
			"message", msg,
		)

		slog.Info("sending message over the network", "remote", conn.RemoteAddr())
		if err := sendString(conn, msg); err != nil {
			slog.Error("failed to send string over network", "remote", conn.RemoteAddr())
			return
		}
	}
}

func monitorConn(r io.Reader, addr string) {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		slog.Error("failed to make logs directory")
		return
	}

	name := fmt.Sprintf("traffic_%s.log", addr)
	f, err := os.OpenFile(
		name,
		os.O_APPEND|os.O_WRONLY|os.O_CREATE,
		0o644,
	)
	if err != nil {
		slog.Error("failed to open logs file", "name", name, "error", err)
		return
	}

	_, err = io.Copy(f, r)
	if err != nil {
		slog.Error("failed to copy connection to traffic log file", "error", err)
	}
}
