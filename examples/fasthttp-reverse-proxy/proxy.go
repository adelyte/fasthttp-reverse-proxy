package main

import (
	"strings"
	"time"

	"github.com/yeqown/log"

	"github.com/valyala/fasthttp"
	proxy "github.com/adelyte/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer  = proxy.NewReverseProxy("localhost:8080", proxy.WithTimeout(5*time.Second))
	proxyServer2 = proxy.NewReverseProxy("api-js.mixpanel.com")
	proxyServer3 = proxy.NewReverseProxy("baidu.com")
)

// ProxyHandler ... fasthttp.RequestHandler func
func ProxyHandler(ctx *fasthttp.RequestCtx) {
	requestURI := string(ctx.RequestURI())
	log.Info("requestURI=", requestURI)

	if strings.HasPrefix(requestURI, "/local") {
		// "/local" path proxy to localhost
		arr := strings.Split(requestURI, "?")
		if len(arr) > 1 {
			arr = append([]string{"/foo"}, arr[1:]...)
			requestURI = strings.Join(arr, "?")
		}

		ctx.Request.SetRequestURI(requestURI)
		proxyServer.ServeHTTP(ctx)
	} else if strings.HasPrefix(requestURI, "/baidu") {
		proxyServer3.ServeHTTP(ctx)
	} else {
		proxyServer2.ServeHTTP(ctx)
	}
}

func main() {
	if err := fasthttp.ListenAndServe(":8081", ProxyHandler); err != nil {
		log.Fatal(err)
	}
}
