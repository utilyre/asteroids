package main

import (
	"bytes"
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
		go HandleConn(connReader, conn)
	}
}

func monitorConn(r io.Reader, addr string) {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		slog.Error("failed to make logs directory")
		return
	}

	name := fmt.Sprintf("logs/traffic_%s.log", addr)
	f, err := os.Create(name)
	if err != nil {
		slog.Error("failed to open logs file", "name", name, "error", err)
		return
	}
	defer f.Close()

	slog.Info("traffic logs file opened", "name", f.Name())

	slog.Info("copying traffic to log file", "name", name)
	_, err = io.Copy(f, r)
	if err != nil {
		slog.Error("failed to copy connection to traffic log file", "error", err)
	}
}
