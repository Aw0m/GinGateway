package constant

import "errors"

const (
	BucketLockType_CAS   = 1
	BucketLockType_Mutex = 2
)

var (
	RouterErr_ServiceNotFound = errors.New("service not found")
	HttpType_HTTP             = "http"
	HttpType_HTTPS            = "https"
)
