package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// https://github.com/golang/go/issues/8456
	lock, err := os.OpenFile("lockfile", os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer lock.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
		os.Remove("lockfile")
		log.Println("lockfile removed by signal handler")
		os.Exit(1)
	}()

	//--

	log.Println("want lock")
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("acquired lock")

	defer func() {
		syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
		os.Remove("lockfile")
		log.Println("lockfile removed")
	}()

	time.Sleep(10 * time.Second)
	log.Println("done")
}
