package main

import (
	"fmt"
	"os"
	"os/signal"
	//"time"
	"encoding/json"
	"github.com/dashotv/flame"
	"github.com/dashotv/rabbit"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron"
	"log"
)

var (
	flameUrl       string
	rabbitUrl      string
	rabbitExchange string
	rabbitType     string
)

func init() {
	flameUrl = os.Getenv("FLAME_URL")
	rabbitUrl = os.Getenv("RABBIT_URL")
	rabbitExchange = os.Getenv("RABBIT_EXCHANGE")
	rabbitType = os.Getenv("RABBIT_TYPE")
}

func main() {
	var err error
	var r *rabbit.Client
	var f *flame.Client
	var p chan []byte

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	f = flame.NewClient(flameUrl)

	if r, err = rabbit.NewClient(rabbitUrl); err != nil {
		log.Fatal("error: ", err)
	}

	if p, err = r.Producer(rabbitExchange, rabbitType); err != nil {
		log.Fatal("error: ", err)
	}

	cron := cron.New()
	cron.AddFunc("* * * * * *", func() {
		var r *flame.Response
		var d []byte

		//fmt.Printf("timer: %s\n", time.Now())

		if r, err = f.List(); err != nil {
			log.Fatal("error: ", err)
		}

		if d, err = json.Marshal(r); err != nil {
			log.Fatal("error: ", err)
		}

		p <- d

		for _, t := range r.Torrents {
			fmt.Printf("%3.0f %6.2f%% %10.2fmb %8.8s %s\n", t.Queue, t.Progress, t.SizeMb(), t.State, t.Name)
		}
	})
	cron.Start()

	// Block until a signal is received.
	// This means we will run until killed / interrupted
	select {
	case s := <-sig:
		fmt.Println("Got signal:", s)
		return
	}
}
