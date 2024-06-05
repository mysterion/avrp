package thumbnails

import (
	"encoding/hex"
	"io"
	"log"
	"os"
)

func logg(args ...string) {

}

func GenerateHash(file string) (string, error) {
	s, err := os.Stat(file)
	if err != nil {
		log.Println(err)
		return "", err
	}

	sz := s.Size()

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer f.Close()

	var i int64
	hash := make([]byte, 62)
	for i = 1; 1<<i < sz && i < 62; i++ {
		_, err := f.Seek(1<<i, io.SeekStart)
		if err != nil {
			log.Println(err)
			return "", err
		}

		_, err = f.Read(hash[i : i+1])
		if err != nil && err != io.EOF {
			log.Println(err)
			return "", err
		}

	}

	return hex.EncodeToString(hash[:i]), nil
}
