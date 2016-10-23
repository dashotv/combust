package main

import (
	"os"
	"os/signal"
	"fmt"
	"github.com/dashotv/rabbit"
	"github.com/dashotv/flame"
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
)

var (
	flameUrl string
	rabbitUrl string
	rabbitExchange string
	rabbitType string
	rabbitQueue string
)

func init() {
	flameUrl = os.Getenv("FLAME_URL")
	rabbitUrl = os.Getenv("RABBIT_URL")
	rabbitExchange = os.Getenv("RABBIT_EXCHANGE")
	rabbitType = os.Getenv("RABBIT_TYPE")
	rabbitQueue = os.Getenv("RABBIT_QUEUE")
}

func main() {
	var client *rabbit.Client
	var err error
	var consuming chan []byte

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	if client, err = rabbit.NewClient(rabbitUrl); err != nil {
		fmt.Println("error: ", err)
		return
	}

	if consuming, err = client.Consumer(rabbitExchange, rabbitType, rabbitQueue); err != nil {
		fmt.Println("error: ", err)
		return
	}

	// Block until a signal is received.
	// This means we will run until killed / interrupted
	for {
		select {
		case m := <-consuming:
			//fmt.Println("Got message: ", string(m))
			d := &flame.Response{}
			if err = json.Unmarshal(m, d); err != nil {
				fmt.Println("error unmarshaling: ", err)
			}
			//fmt.Println("decoded: ", d)
			for _, t := range d.Torrents {
				fmt.Printf("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
			}
		case s := <-c:
			fmt.Println("Got signal:", s)
			return
		}
	}
}
