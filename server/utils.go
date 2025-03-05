package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
)

func recvUInt64(r io.Reader) (uint64, error) {
	var value uint64
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}

	return value, nil
}

func sendUInt64(w io.Writer, value uint64) error {
	if err := binary.Write(w, binary.BigEndian, value); err != nil {
		return err
	}

	return nil
}

func recvBytes(r io.Reader) ([]byte, error) {
	size, err := recvUInt64(r)
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

func sendBytes(w io.Writer, p []byte) error {
	size := uint64(len(p))
	if err := sendUInt64(w, size); err != nil {
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
