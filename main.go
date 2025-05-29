package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
}

func xorFile(path string, key byte) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// Vérifie s'il y a un \n final
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	if size == 0 {
		return nil // fichier vide
	}

	lastByte := make([]byte, 1)
	_, err = f.ReadAt(lastByte, size-1)
	if err != nil {
		return err
	}

	hasTrailingNewline := lastByte[0] == '\n'

	if hasTrailingNewline {
		// Tronque le \n
		err = f.Truncate(size - 1)
		if err != nil {
			return err
		}
	}

	// Rewind
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	buf := make([]byte, 8192)

	for {
		pos, err := f.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		n, err := f.Read(buf)
		if n > 0 {
			for i := 0; i < n; i++ {
				buf[i] ^= key
			}
			_, err = f.Seek(pos, io.SeekStart)
			if err != nil {
				return err
			}
			_, err = f.Write(buf[:n])
			if err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	// Ajoute à nouveau \n si besoin
	if hasTrailingNewline {
		_, err = f.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte{'\n'})
		if err != nil {
			return err
		}
	}

	return nil
}

func Taz(path string, key byte) error {
	return xorFile(path, key)
}

func main() {
	key := flag.Uint("k", 200, "Clé (ex: 170)")
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
	check(xorFile(path, byte(*key)))

	fmt.Println("File", path, "has been tazed !")
}
