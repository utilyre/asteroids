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

	player := &Player{}

	srv.Handle("player.move_forward", func(ctx context.Context, body []byte) error {
		player.MoveForward(0.16)
		slog.Debug("player rotated left",
			"position", player.Position,
			"rotation", player.Rotation,
		)
		return nil
	})
	srv.Handle("player.move_backward", func(ctx context.Context, body []byte) error {
		player.MoveBackward(0.16)
		slog.Debug("player rotated left",
			"position", player.Position,
			"rotation", player.Rotation,
		)
		return nil
	})
	srv.Handle("player.rotate_left", func(ctx context.Context, body []byte) error {
		player.RotateLeft(0.16)
		slog.Debug("player rotated left",
			"position", player.Position,
			"rotation", player.Rotation,
		)
		return nil
	})
	srv.Handle("player.rotate_right", func(ctx context.Context, body []byte) error {
		player.RotateRight(0.16)
		slog.Debug("player rotated right",
			"position", player.Position,
			"rotation", player.Rotation,
		)
		return nil
	})

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", "error", err)
	}
}
