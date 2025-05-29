package taz

import (
	"io"
	"os"
)

var DEFAULT_KEY byte = 0x2A // Clé par défaut pour le taz

// SetKey permet de changer la clé par défaut utilisée pour le taz
func SetKey(key byte) {
	DEFAULT_KEY = key
}

// ReadTazzedContent lit le contenu d'un fichier taz et le déchiffre
func ReadUntazzedContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	for i := range data {
		data[i] ^= DEFAULT_KEY
	}

	return string(data), nil
}

// TazFile chiffre le contenu d'un fichier en utilisant la clé par défaut
func TazFile(path string) error {
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
				buf[i] ^= DEFAULT_KEY
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
