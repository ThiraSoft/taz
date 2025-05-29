package main

import (
	"flag"
	"fmt"
	"os"

	taz "github.com/ThiraSoft/taz/pkg/taz"
)

func main() {
	key := flag.Uint("k", uint(taz.DEFAULT_KEY), "Clé (ex: 0x2A)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage : taz [-k clé] <fichier>")
		os.Exit(1)
	}

	if *key > 255 {
		fmt.Fprintln(os.Stderr, "Erreur : la clé doit être entre 0 et 255")
		os.Exit(1)
	}

	path := flag.Arg(0)
	taz.SetKey(uint8(*key))
	err := taz.TazFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
	}
}
