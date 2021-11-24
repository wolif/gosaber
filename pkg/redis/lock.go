package redis

import (
	"fmt"
	"sync"
	"time"
)

type mutex struct {
	Retries    int
	RetryDelay time.Duration
	Expire     time.Duration
	MT         sync.Mutex
}

func Mutex(tries int, interval time.Duration, expire time.Duration) (*mutex, error) {
	if tries < 1 {
		return nil, fmt.Errorf("pararm tries must greater than 0")
	}
	return &mutex{
		Retries:    tries,
		RetryDelay: interval,
		Expire:     expire,
	}, nil
}

func (m *mutex) Lock(key string) (bool, error) {
	m.MT.Lock()
	defer m.MT.Unlock()

	for i := 0; i < m.Retries; i++ {
		b, e := Client.SetNX(key, 1, m.Expire).Result()
		if e != nil {
			return false, e
		}
		if b {
			return true, nil
		}
		time.Sleep(m.RetryDelay)
	}
	return false, nil
}

func (m *mutex) Unlock(key string) (bool, error) {
	n, e := Client.Del(key).Result()
	if e != nil {
		return false, e
	}
	return n == 1, nil
}
