package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

	logger := slog.With("address", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info("started listening")

	// serve
	for {
		logger.Info("waiting to accept a new connection")
		conn, err := ln.Accept()
		if err != nil {
			logger.Error("failed to establish connection", "error", err)
			continue
		}
		logger.Info("connection accepted", "remote", conn.RemoteAddr())

		monitorBuffer := &bytes.Buffer{}
		connReader := io.TeeReader(conn, monitorBuffer)

		logger.Info("monitoring connection")
		go srv.monitorConn(monitorBuffer, conn.RemoteAddr())

		logger.Info("handling connection")
		go func() {
			defer func() {
				if err := conn.Close(); err != nil {
					logger.Error("failed to close connection", "error", err)
				}
			}()

			srv.handleConn(connReader, conn.RemoteAddr())
		}()
	}
}

func (srv *Server) handleConn(r io.Reader, remote net.Addr) {
	logger := slog.With("remote", remote)

	for {
		logger.Info("reading message from network")

		msg, err := ReadMessage(r)
		if errors.Is(err, io.EOF) {
			logger.Info("connection closed")
			break
		}
		if err != nil {
			logger.Error("failed to read message from connection", "error", err)
			continue
		}

		slog.Info("received message", "message", msg)
		var m map[string]any
		if err := json.Unmarshal(msg.Body, &m); err != nil {
			logger.Error("failed to unmarshal message body", "error", err)
		}
		logger.Info("unmarshaled message body", "body", m)

		// TODO: dispatch message
		// client says asteroid/spawn(position, velocity)
		// now, the server is ought to update its state
		// so the server dispatches the message to the corresponding method (manually?)
		// and the method updates state (maybe a response?) (but definitely log)
	}
}

func (srv *Server) monitorConn(r io.Reader, remote net.Addr) {
	logger := slog.With("remote", remote)

	logDir := "logs"
	if len(srv.LogDir) > 0 {
		logDir = srv.LogDir
	}

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		logger.Error("failed to make log directory", "path", logDir, "error", err)
		return
	}

	name := fmt.Sprintf("logs/traffic_%s.log", remote)
	logger = logger.With("filename", name)

	f, err := os.Create(name)
	if err != nil {
		logger.Error("failed to create log file", "error", err)
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	logger.Info("copying traffic to log file")
	_, err = io.Copy(f, r)
	if err != nil {
		logger.Error("failed to copy traffic to log file", "error", err)
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
	if size > 1024 {
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
