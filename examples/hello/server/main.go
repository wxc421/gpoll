package main

import (
	"gpoll"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	gPoll := gpoll.New(":8080")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	logger.Info("start")
	if err := gPoll.Start(); err != nil {
		panic(err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("close...")
	gPoll.Close()

	select {}
}
