package main

import (
	"context"
	"log/slog"
	"time"
)

// TODO: the server should spawn the asteroids

// client says asteroid/spawn(position, velocity)
// now, the server is ought to update its state
// so the server dispatches the message to the corresponding method (manually?)
// and the method updates state (maybe a response?) (but definitely log)

func TODO() {
	// client gotta keep track of last snapshot and current snapshot
	// then interpolate from frameTime-interpTime to frameTime
	// (instead of interpolating from the last snapshot's time to frameTime)

	// game loop
	for {
		// process user commands

		// run physical simulation step

		// check game rules

		// broadcast update
	}
}

type UserCommand string

const (
	PlayerMoveForward  UserCommand = "player.move_forward"
	PlayerMoveBackward UserCommand = "player.move_backward"
	PlayerRotateLeft   UserCommand = "player.rotate_left"
	PlayerRotateRight  UserCommand = "player.rotate_right"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &Server{Addr: ":3000"}

	var userCommandQueue []UserCommand

	slog.Info("starting to sample user command queue")
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			slog.Debug("user command queue sampled", "queue", userCommandQueue)

			select {
			case <-ticker.C:
			case <-ctx.Done():
				slog.Info("stopped sampling user command queue")
				break
			}
		}
	}()

	slog.Info("starting simulation loop")
	go func() {
		ticker := time.NewTicker(time.Second / 60)
		defer ticker.Stop()

		for {
			if len(userCommandQueue) == 0 {
				continue
			}

			var uc UserCommand
			uc, userCommandQueue = userCommandQueue[0], userCommandQueue[1:]

			slog.Debug("user command consumed", "user_command", uc)

			select {
			case <-ticker.C:
			case <-ctx.Done():
				slog.Info("stopped simulation loop")
				break
			}
		}
	}()

	srv.Handle("player.move_forward", func(ctx context.Context, body []byte) error {
		userCommandQueue = append(userCommandQueue, PlayerMoveForward)
		return nil
	})
	srv.Handle("player.move_backward", func(ctx context.Context, body []byte) error {
		userCommandQueue = append(userCommandQueue, PlayerMoveBackward)
		return nil
	})
	srv.Handle("player.rotate_left", func(ctx context.Context, body []byte) error {
		userCommandQueue = append(userCommandQueue, PlayerRotateLeft)
		return nil
	})
	srv.Handle("player.rotate_right", func(ctx context.Context, body []byte) error {
		userCommandQueue = append(userCommandQueue, PlayerRotateRight)
		return nil
	})

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", "error", err)
	}
}
