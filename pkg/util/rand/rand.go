package rand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))

	// NewSource returns a new pseudo-random Source seeded with the given value.
	// Unlike the default Source used by top-level functions, this source is not
	// safe for concurrent use by multiple goroutines. Hence the need for a mutex.
	seededIDLock sync.Mutex
)

func Int63() int64 {
	seededIDLock.Lock()
	defer seededIDLock.Unlock()
	return seededIDGen.Int63()
}

func Int63Range(min, max int64) int64 {
	return Int63()%(max-min) + min
}

func Int31() int32 {
	seededIDLock.Lock()
	defer seededIDLock.Unlock()
	return seededIDGen.Int31()
}

func Int31Range(min, max int32) int32 {
	return Int31()%(max-min) + min
}

func Bytes(n int) ([]byte, error) {
	seededIDLock.Lock()
	defer seededIDLock.Unlock()

	b := make([]byte, n)
	_, err := seededIDGen.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
