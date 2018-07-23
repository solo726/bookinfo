package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)

	// Place whatever shutdown handling you want here
	time.Sleep(10 * time.Second)
	fmt.Println("service shutdown...")

	errc <- terminateError
}
