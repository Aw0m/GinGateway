package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"wxprojectApiGateway/service/discovery"
)

var (
	RouteMap map[string]string
)

func init() {
	RouteMap = make(map[string]string)
	yamlFile, err := ioutil.ReadFile("route.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	if err = yaml.Unmarshal(yamlFile, &RouteMap); err != nil {
		// error handling
		log.Fatalln(err)
	}
}

func RouteForward(serviceName string) gin.HandlerFunc {
	return func(context *gin.Context) {
		// 先通过Token获得用户id
		openid := GetOpenid(context.GetHeader("Authorization"))
		context.Request.Header.Set("userID", openid)
		if err := ForwardHandler(context.Writer, context.Request, serviceName); err != nil {
			if err.Error() == "not exist" {
				context.Status(http.StatusNotFound)
			} else {
				context.Status(http.StatusServiceUnavailable)
			}
			context.Abort()
			return
		}
		context.Next()
	}
}

func ForwardHandler(writer http.ResponseWriter, request *http.Request, serviceName string) error {
	host, exist := discovery.GetService(serviceName)
	if !exist {
		return errors.New("not exist")
	}

	path, exist := RouteMap[request.URL.Path]
	if !exist {
		path = request.URL.Path
	}

	var rawURL strings.Builder
	rawURL.WriteString("http://")
	rawURL.WriteString(host)
	rawURL.WriteString(path)
	rawURL.WriteString("?")
	rawURL.WriteString(request.URL.RawQuery)

	u, err := url.Parse(rawURL.String())
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
