package main

import (
	"log"
	"os"

	"github.com/taimats/wc/internal"
)

func main() {
	err := internal.CmdWC(os.Args)
	if err != nil {
		log.Fatalf("wc stops running: (error: %s)", err)
	}
}
