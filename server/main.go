package main

import (
	"context"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	srv := &Server{Addr: ":3000"}

	srv.Handle("player.move_forward", func(ctx context.Context, body []byte) error {
		slog.Debug("⚠️ player moved forward")
		return nil
	})
	srv.Handle("asteroid.spawn", func(ctx context.Context, body []byte) error {
		slog.Debug("⚠️ ASTEROID SPAWNED")
		return nil
	})

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", "error", err)
	}
}
