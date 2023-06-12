package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SevcikMichal/microfrontends-controller/internal/api"
	"github.com/SevcikMichal/microfrontends-controller/internal/configuration"
	"github.com/gorilla/mux"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

type RouterProvider struct {
	FrontendConfigApi *api.MicroFrontendConfigApi
}

func (routerProvider *RouterProvider) CreateRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	basePathRouter := router.PathPrefix(configuration.GetBaseURL()).Subrouter()

	cacheClient5sec := createCaheClient(5)
	feConfigHandleFunc := http.HandlerFunc(routerProvider.FrontendConfigApi.GetMicroFrontendConfigs)
	feConfigJsHandleFunc := http.HandlerFunc(routerProvider.FrontendConfigApi.GetMicroFrontendConfigsAsJavaScritp)

	basePathRouter.Handle("/fe-config", cacheClient5sec.Middleware(feConfigHandleFunc)).Methods("GET")
	basePathRouter.Handle("/fe-config.mjs", cacheClient5sec.Middleware(feConfigJsHandleFunc)).Methods("GET")
	router.HandleFunc("/healthz", api.GetHealthInfo).Methods("GET")
	basePathRouter.HandleFunc("/healthz", api.GetHealthInfo).Methods("GET")

	return router
}

func createCaheClient(timeToLiveInSeconds int) *cache.Client {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		panic(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(time.Duration(timeToLiveInSeconds)*time.Second),
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return cacheClient
}
