package main

import (
	"flag"
	"fmt"
	"os"

	taz "github.com/ThiraSoft/taz/pkg/taz"
)

func main() {
	key := flag.Uint("k", 0xAA, "Clé (ex: 0xAA)")
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
	err := taz.TazFile(path, byte(*key))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
	}
}
