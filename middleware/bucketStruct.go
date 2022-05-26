package middleware

import (
	"sync"
	"time"
)

type TokenBucket struct {
	//token []int
	limit  int
	num    int
	ticker *time.Ticker
	mu     *sync.Mutex
}

var (
	bucket     *TokenBucket
	singleLock sync.Mutex
)

func GetBucket(limit, num int, ticker *time.Ticker) *TokenBucket {
	if bucket != nil {
		return bucket
	}
	singleLock.Lock()
	defer singleLock.Unlock()
	if bucket == nil {
		if limit < num {
			limit = num
		}
		bucket = &TokenBucket{
			limit:  limit,
			num:    num,
			ticker: ticker,
			mu:     &sync.Mutex{},
		}

		// 开始定时向令牌桶中添加token
		go func() {
			for {
				select {
				case <-bucket.ticker.C:
					bucket.addToken()
				}
			}
		}()
	}
	return bucket
}

func (bucket *TokenBucket) GetToken() bool {
	if bucket.num > 0 {
		bucket.mu.Lock()
		defer bucket.mu.Unlock()
		if bucket.num > 0 {
			bucket.num--
			return true
		}
	}
	return false
}

func (bucket *TokenBucket) addToken() bool {
	if bucket.num < bucket.limit {
		bucket.mu.Lock()
		defer bucket.mu.Unlock()
		if bucket.num < bucket.limit {
			bucket.num++
			return true
		}
	}
	return false
}
