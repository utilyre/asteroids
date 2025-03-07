package main

import (
	"context"
	"log/slog"
)

// TODO: the server should spawn the asteroids

// client says asteroid/spawn(position, velocity)
// now, the server is ought to update its state
// so the server dispatches the message to the corresponding method (manually?)
// and the method updates state (maybe a response?) (but definitely log)

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
