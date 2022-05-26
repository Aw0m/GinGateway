package middleware

import (
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucketCAS struct {
	//token []int
	limit  int32
	num    int32
	ticker *time.Ticker
}

var (
	bucketCAS     *TokenBucketCAS
	singleLockCAS sync.Mutex
)

func GetBucketCAS(limit, num int32, ticker *time.Ticker) *TokenBucketCAS {
	if bucketCAS != nil {
		return bucketCAS
	}
	singleLockCAS.Lock()
	defer singleLockCAS.Unlock()
	if bucketCAS == nil {
		if limit < num {
			limit = num
		}
		bucketCAS = &TokenBucketCAS{
			limit:  limit,
			num:    num,
			ticker: ticker,
		}
		go func() {
			for {
				select {
				case <-bucketCAS.ticker.C:
					bucketCAS.addToken()
				}
			}
		}()
	}
	return bucketCAS
}

func (bucketCAS *TokenBucketCAS) addToken() {
	if atomic.LoadInt32(&(bucketCAS.num)) >= bucketCAS.limit {
		return
	}
	atomic.AddInt32(&(bucketCAS.num), 1)
}

func (bucketCAS *TokenBucketCAS) GetToken() bool {
	for {
		oldVal := atomic.LoadInt32(&(bucketCAS.num))
		if oldVal <= 0 {
			return false
		}
		if atomic.CompareAndSwapInt32(&(bucketCAS.num), oldVal, oldVal+1) {
			return true
		}
	}
}
