package router

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/SevcikMichal/microfrontends-controller/internal/api"
	"github.com/SevcikMichal/microfrontends-controller/internal/configuration"
	"github.com/gorilla/mux"
)

const (
	lookupWebComponentKeyWord = "lookup-web-component"
	contentHashKeyWord        = "content-hash"
)

type RouterProvider struct {
	FrontendConfigApi *api.MicroFrontendConfigApi
	WebComponentApi   *api.WebComponentApi
	AppIconsApi       *api.AppIconsApi
}

func (routerProvider *RouterProvider) CreateRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	basePathRouter := router.PathPrefix(configuration.GetBaseURL()).Subrouter()

	// Frontend config handlers
	feConfigHandleFunc := http.HandlerFunc(routerProvider.FrontendConfigApi.GetMicroFrontendConfigs)
	feConfigJsHandleFunc := http.HandlerFunc(routerProvider.FrontendConfigApi.GetMicroFrontendConfigsAsJavaScritp)
	basePathRouter.Handle("/fe-config", routerProvider.cache("30", routerProvider.FrontendConfigApi.MicroFrontendProvider.GetETag, feConfigHandleFunc)).Methods("GET")
	basePathRouter.Handle("/fe-config.mjs", routerProvider.cache("30", routerProvider.FrontendConfigApi.MicroFrontendProvider.GetETag, feConfigJsHandleFunc)).Methods("GET")

	// Health check handlers
	router.HandleFunc("/healthz", api.GetHealthInfo).Methods("GET")
	basePathRouter.HandleFunc("/healthz", api.GetHealthInfo).Methods("GET")

	// Web component handlers
	webComponentHandleFunc := http.HandlerFunc(routerProvider.WebComponentApi.GetWebComponent)
	basePathRouter.PathPrefix("/web-components").Handler(routerProvider.cache("3600", func() string {
		return lookupWebComponentKeyWord // Dirty hack to be able to fetch it here so that we don't need to duplicate the code in api
	}, webComponentHandleFunc)).Methods("GET")

	// App icon handlers
	appIconHandleFunc := http.HandlerFunc(routerProvider.AppIconsApi.GetAppIcon)
	basePathRouter.PathPrefix("/app-icons").Handler(routerProvider.cache("604800", func() string {
		return contentHashKeyWord // Dirty hack to calculate hash from content
	}, appIconHandleFunc)).Methods("GET")

	basePathRouter.PathPrefix("/").HandlerFunc(passThrough).Methods("GET")

	return router
}

func (routerProvider *RouterProvider) cache(durationInSeconds string, eTagGetter func() string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eTag := eTagGetter()

		// Dirty hack to be able to fetch it here so that we don't need to duplicate the code in api
		if eTag == lookupWebComponentKeyWord {
			segments := strings.Split(r.URL.Path, "/")
			namespace, name := segments[2], segments[3]
			eTag = routerProvider.FrontendConfigApi.MicroFrontendProvider.GetMicrofrontendHashSuffix(namespace, name)
		}

		// Dirty hack to calculate hash from content
		if eTag == contentHashKeyWord {
			segments := strings.Split(r.URL.Path, "/")
			navigationPath := strings.Join(segments[(len(segments)-1):], "/")
			appIcon := routerProvider.FrontendConfigApi.MicroFrontendProvider.GetMicrofrontendAppIcon(navigationPath)
			hash := md5.Sum([]byte(appIcon.Data))
			hashStr := hex.EncodeToString(hash[:])
			eTag = hashStr
		}

		if r.Header.Get("Cache-Control") != "no-cache" && r.Header.Get("If-None-Match") == eTag {
			w.Header().Set("Cache-Control", "max-age="+durationInSeconds)
			w.Header().Set("Last-Modified", r.Header.Get("Last-Modified"))
			w.Header().Set("ETag", eTag)
			w.WriteHeader(http.StatusNotModified)
			return
		}

		c := httptest.NewRecorder()
		handler(c, r)

		for k, v := range c.Header() {
			w.Header().Set(k, strings.Join(v, ", "))
		}

		w.Header().Set("Cache-Control", "max-age="+durationInSeconds)
		w.Header().Set("Last-Modified", time.Now().Format(time.RFC1123))
		w.Header().Set("ETag", eTag)

		w.WriteHeader(c.Code)
		content := c.Body.Bytes()

		w.Write(content)
	})
}

func passThrough(w http.ResponseWriter, r *http.Request) {
	log.Println("Request did not match any endpoint passing to web ui.")
	path := r.URL.Path

	req, err := http.NewRequest("GET", "http://localhost:8082"+path, r.Body) // TODO: Read from config
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Proxying request to resolved URL.", "Resolved URL:", "http://localhost:8082"+path) // TODO: Read from config
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		dst.Set(key, strings.Join(values, ", "))
	}
}
