package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	lock, err := os.Create("lockfile")
	if err != nil {
		log.Println(err)
	}
	defer lock.Close()

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
		os.Remove("lockfile")
		done <- true
	}()

	//--

	fmt.Println("want lock")
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("acquired lock")
	defer func() {
		syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
		os.Remove("lockfile")
	}()

	<-done
}
