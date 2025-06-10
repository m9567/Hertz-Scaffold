package common

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

var (
	proxyMap = make(map[string]*httputil.ReverseProxy)
)

func NewProxy(targetHost string) *httputil.ReverseProxy {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
	}
	return proxy
}

func ProxyUrl(w http.ResponseWriter, r *http.Request, uri string) {
	reverseProxy := proxyMap[uri]
	if reverseProxy != nil {
		reverseProxy.ServeHTTP(w, r)
		return
	}
	for i := 0; i < 10; i++ {
		tryLock := Locker.TryLock(uri)
		if tryLock {
			defer Locker.Unlock(uri)
			reverseProxy = NewProxy(uri)
			proxyMap[uri] = reverseProxy
			reverseProxy.ServeHTTP(w, r)
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
	reverseProxy = NewProxy(uri)
	reverseProxy.ServeHTTP(w, r)
}
