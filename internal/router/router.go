package router

import (
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
)

type RouterProvider struct {
	FrontendConfigApi *api.MicroFrontendConfigApi
	WebComponentApi   *api.WebComponentApi
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
