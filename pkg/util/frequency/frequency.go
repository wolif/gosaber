package frequency

import (
	"fmt"
	goCache "github.com/patrickmn/go-cache"
	"github.com/wolif/gosaber/pkg/util/strs"
	"sync"
	"time"
)

type frequency struct {
	CtrlMutex sync.Mutex
	Times     int64
	TimeSpan  time.Duration
	Cache     *goCache.Cache
	Key       string
}

func (f *frequency) Ctrl() {
	f.CtrlMutex.Lock()
	defer f.CtrlMutex.Unlock()
	timesData, expiration, found := f.Cache.GetWithExpiration(f.Key)
	if !found {
		f.Cache.Set(f.Key, int64(1), f.TimeSpan)
		return
	}

	times := int64(strs.StrToIntWithFallback(fmt.Sprint(timesData), 1))
	if times >= f.Times {
		now := time.Now()
		if now.UnixNano() < expiration.UnixNano() {
			time.Sleep(time.Duration(expiration.UnixNano() - now.UnixNano()))
		}
		f.Cache.Set(f.Key, int64(1), f.TimeSpan)
		return
	}
	_ = f.Cache.Increment(f.Key, 1)
}

func New(times int64, timeSpan time.Duration) *frequency {
	return &frequency{
		Times:    times,
		TimeSpan: timeSpan,
		Cache:    goCache.New(goCache.NoExpiration, timeSpan+time.Second*1),
		Key:      fmt.Sprint(time.Now().UnixNano()),
	}
}
