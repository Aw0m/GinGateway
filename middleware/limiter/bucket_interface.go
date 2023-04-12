package limiter

type TokenBucket interface {
	GetToken() bool
	addToken()
}
