package main

import "log/slog"

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	srv := &Server{
		Addr: ":3000",
	}
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", "error", err)
	}
}
