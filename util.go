package main

import (
	"os"
	"os/signal"
	"syscall"
)

func RegisterExitProcess() {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
		os.Kill)

	go func() {
		for range c {
			SweepFiles()
			os.Exit(1)
		}
	}()
}
