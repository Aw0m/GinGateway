package router

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"

	"GinGateway/constant"
	"GinGateway/service/discovery"
)

func RouteForward(serviceName string, HttpType string) gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ForwardHandler(context.Writer, context.Request, serviceName, HttpType)
		if err == nil {
			context.Next()
			return
		}

		switch err {
		case constant.RouterErr_ServiceNotFound:
			context.Status(http.StatusNotFound)
		default:
			context.Status(http.StatusServiceUnavailable)
		}
		context.Abort()
	}

}

func ForwardHandler(writer http.ResponseWriter, request *http.Request, serviceName, HttpType string) error {
	host, exist := discovery.GetService(serviceName)
	if !exist {
		return constant.RouterErr_ServiceNotFound
	}

	rawURL := fmt.Sprintf("%s://%s%s?%s", HttpType, host, request.URL.Path, request.URL.RawQuery)
	u, err := url.Parse(rawURL)
	if nil != err {
		log.Println(err)
		return err
	}

	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = u
			r.Header = request.Header
		},
	}
	proxy.ServeHTTP(writer, request)
	return nil
}
