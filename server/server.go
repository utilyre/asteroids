package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

type Server struct {
	Addr   string
	LogDir string
}

func (srv *Server) ListenAndServe() error {
	addr := "localhost:80"
	if len(srv.Addr) > 0 {
		addr = srv.Addr
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	slog.Info("started listening", "address", addr)

	for {
		slog.Info("waiting to accept a new connection", "address", addr)
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("failed to establish connection",
				"address", addr,
				"error", err,
			)
			continue
		}
		slog.Info("accepted connection", "remote", conn.RemoteAddr())

		monitorBuffer := &bytes.Buffer{}
		connReader := io.TeeReader(conn, monitorBuffer)

		go srv.monitorConn(monitorBuffer, conn.RemoteAddr())
		go srv.handleConn(connReader, conn)
	}
}

func (srv *Server) handleConn(r io.Reader, conn net.Conn) {
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

func (srv *Server) monitorConn(r io.Reader, remote net.Addr) {
	logDir := "logs"
	if len(srv.LogDir) > 0 {
		logDir = srv.LogDir
	}

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		slog.Error("failed to make log directory", "remote", remote, "log_dir", logDir, "error", err)
		return
	}

	name := fmt.Sprintf("logs/traffic_%s.log", remote)
	f, err := os.Create(name)
	if err != nil {
		slog.Error("failed to open log file", "remote", remote, "name", name, "error", err)
		return
	}
	defer f.Close() // TODO: return error or error group if function already failed

	slog.Info("copying traffic to log file", "remote", remote, "name", name)
	_, err = io.Copy(f, r)
	if err != nil {
		slog.Error("failed to copy connection to traffic log file", "remote", remote, "name", name, "error", err)
	}
}

func readUInt64(r io.Reader) (uint64, error) {
	var value uint64
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}

	return value, nil
}

func writeUInt64(w io.Writer, value uint64) error {
	if err := binary.Write(w, binary.BigEndian, value); err != nil {
		return err
	}

	return nil
}

func readBytes(r io.Reader) ([]byte, error) {
	size, err := readUInt64(r)
	if err != nil {
		return nil, err
	}
	slog.Debug("read a byte slice size", "size", size)
	if size > 64 {
		return nil, errors.New("receiving over-sized bytes")
	}

	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	if uint64(n) < size {
		panic("read bytes and pre-known size mismatch")
	}

	return buf, nil
}

func writeBytes(w io.Writer, p []byte) error {
	size := uint64(len(p))
	if err := writeUInt64(w, size); err != nil {
		return err
	}

	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if uint64(n) < size {
		panic("written bytes and the size of bytes mismatch")
	}

	return nil
}
