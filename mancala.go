package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"mancala/kalah"
	"os"
	"time"
)

func main() {
	toLog := flag.Bool("log", false, "whether to log gameplay to a file")
	flag.Parse()
	var w io.Writer
	if *toLog {
		l, err := os.Create("log" + fmt.Sprintf("%v", time.Now().Unix()))
		if err != nil {
			log.Fatal(err)
		}
		w = io.MultiWriter(os.Stdout, l)
	} else {
		w = os.Stdout
	}
	b := kalah.NewBoard(6, w, true)

	var r = bufio.NewReader(os.Stdin)
	for {
		b.Print()
		fmt.Fprintf(os.Stdout, ">> ")
		in, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if exit := b.Handle(in); exit {
			break
		}
	}
	b.Print()
	if b.GameOver() {
		b.PrintWinner()
	}
}
