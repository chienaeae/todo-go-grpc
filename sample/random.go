package sample

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)


const alphabet = "abcdefghijklimnopqrstuvwxyz"

func init () {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
}

func randomID() string {
	return uuid.New().String()
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0 ; i < n ; i ++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}