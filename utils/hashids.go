package utils

import (
	"log"
	"os"
	"github.com/speps/go-hashids"
)

var h *hashids.HashID

func InitHashID() {
	salt := os.Getenv("HASHID_SALT")
	if salt == "" {
		log.Fatal("HASHID_SALT not set")
	}
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 6
	var err error
	h, err = hashids.NewWithData(hd)
	if err != nil {
		log.Fatalf("Failed to initialize Hashids: %v", err)
	}
}

func EncodeID(id int) (string, error) {
	encodedId, err := h.Encode([]int{id})
	return encodedId, err 
}

func DecodeID(hash string) (int, error) {
	decoded, err := h.DecodeWithError(hash)
	if err != nil || len(decoded) == 0 {
		return -1, err
	}
	return decoded[0], nil
}
