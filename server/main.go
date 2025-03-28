package main

import (
	"context"
	"log/slog"
	"os"
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
	// WARN: this configuration is for development only
	slog.SetDefault(slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "time" {
					return slog.Attr{}
				}
				return a
			},
		}),
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* slog.Info("starting to sample user command queue")
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			userCommandMu.RLock()
			slog.Debug("user command queue sampled", "queue", userCommandQueue)
			userCommandMu.RUnlock()

			select {
			case <-ticker.C:
			case <-ctx.Done():
				slog.Info("stopped sampling user command queue")
				break
			}
		}
	}() */

	userCommandQueue := make(chan UserCommand, 1)
	defer close(userCommandQueue)

	srv := setupServer(userCommandQueue)
	slog.Info("starting server")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("failed to listen and serve", "error", err)
		}
	}()

	slog.Info("starting simulation loop")
	simulate(ctx, userCommandQueue)
}

func setupServer(userCommandQueue chan<- UserCommand) *Server {
	srv := &Server{Addr: ":3000"}

	enqueueCmd := func(cmd UserCommand) {
		// TODO: timeout only if it becomes a bottle neck
		userCommandQueue <- cmd
	}

	srv.Handle("player.move_forward", func(ctx context.Context, body []byte) error {
		enqueueCmd(PlayerMoveForward)
		return nil
	})
	srv.Handle("player.move_backward", func(ctx context.Context, body []byte) error {
		enqueueCmd(PlayerMoveBackward)
		return nil
	})
	srv.Handle("player.rotate_left", func(ctx context.Context, body []byte) error {
		enqueueCmd(PlayerRotateLeft)
		return nil
	})
	srv.Handle("player.rotate_right", func(ctx context.Context, body []byte) error {
		enqueueCmd(PlayerRotateRight)
		return nil
	})

	return srv
}

func simulate(ctx context.Context, userCommandQueue <-chan UserCommand) {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	for {
		select {
		case uc := <-userCommandQueue:
			slog.Debug("user command consumed", "user_command", uc)
		default:
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			slog.Info("stopped simulation loop")
			break
		}
	}
}
