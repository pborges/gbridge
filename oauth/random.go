package oauth

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	rand "math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var src cryptoSource

// GenerateRandomString generates an cryptographically secure string
func GenerateRandomString(length int) string {
	rnd := rand.New(src)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return string(b)
}

// cryptoSource is used to create a secure pseudo-random number generator
type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

// Uint64 gets a unsigned integer from /dev/urandom (Unix) or CryptAcquireContext (Win)
func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		// The crand.Reader returns an error if the underlying system call fails. For instance if it can't read /dev/urandom on a Unix system, or if CryptAcquireContext fails on a Windows system.
		log.Fatal(err)
	}
	return v
}
