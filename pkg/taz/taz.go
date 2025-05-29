package taz

import (
	"io"
	"os"
)

func Taz(path string, key byte) error {
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
