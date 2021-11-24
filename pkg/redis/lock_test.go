package redis

import (
	"github.com/wolif/gosaber/pkg/util/strs"
	"testing"
	"time"
)

func initRedis() {
	Init(&Config{
		Addrs: []string{"10.20.1.20:6379"},
		Pwd:   "mqj123.com",
		DB:    0,
	})
}

func testMutex(t *testing.T) {
	mutex, e := Mutex(7, 100*time.Millisecond, 3*time.Minute)
	if e != nil {
		t.Error(e)
	}

	for i := 0; i < 1; i++ {
		go func(i int) {
			for j := 0; j < 100; j++ {
				func() {
					b, err := mutex.Lock("abc")
					if err != nil {
						t.Error(err)
						return
					}
					if b {
						defer time.Sleep(1 * time.Second)
						defer mutex.Unlock("abc")

						n, e := Client.Get("test_redis_mutex").Result()
						if e != nil {
							t.Error(e)
							return
						}
						num := strs.StrToIntWithDefaultZero(n)
						t.Logf("G%d - %d: %d", i, j+1, num)

						Client.Set("test_redis_mutex", num+1, 10*time.Second)
					} else {
						t.Logf("G%d - %d: %s", i, j+1, "get lock failed")
					}
				}()
			}
		}(i)
	}

	select {}
}

func TestSetCache(t *testing.T) {
	initRedis()
	Client.Set("test_redis_mutex", 0, 10*time.Second)
}

func BenchmarkGet(b *testing.B) {
	initRedis()
	mutex, e := Mutex(7, 100*time.Millisecond, 3*time.Minute)
	if e != nil {
		b.Error(e)
		return
	}
	for i:=0;i<b.N;i++ {
		mutex.Lock("abc")
		mutex.Unlock("abc")
	}
}

func BenchmarkMutex(b *testing.B) {
	initRedis()
	mutex, e := Mutex(7, 100*time.Millisecond, 3*time.Minute)
	if e != nil {
		b.Error(e)
	}

	for j := 0; j < b.N; j++ {
		bl, err := mutex.Lock("abc")
		if err != nil {
			b.Error(err)
			return
		}
		if bl {
			n, e := Client.Get("test_redis_mutex").Result()
			if e != nil {
				b.Error(e)
				return
			}
			num := strs.StrToIntWithDefaultZero(n)
			b.Logf("%d: %d", j+1, num)

			Client.Set("test_redis_mutex", num+1, 10*time.Second)
			mutex.Unlock("abc")
		} else {
			b.Logf("%d: %s", j+1, "get lock failed")
		}

	}
}

func TestMutex1(t *testing.T) {
	initRedis()
	Client.Set("test_redis_mutex", 0, 10*time.Second)
	testMutex(t)
}
func TestMutex2(t *testing.T) {
	initRedis()
	testMutex(t)
}
func TestMutex3(t *testing.T) {
	initRedis()
	testMutex(t)
}
func TestMutex4(t *testing.T) {
	initRedis()
	testMutex(t)
}
func TestMutex5(t *testing.T) {
	initRedis()
	testMutex(t)
}
