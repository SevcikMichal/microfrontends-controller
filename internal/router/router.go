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

	cacheClient := createCaheClient(120)
	feConfigHandleFunc := http.HandlerFunc(routerProvider.FrontendConfigApi.GetMicroFrontendConfigs)

	basePathRouter.Handle("/fe-config", cacheClient.Middleware(feConfigHandleFunc)).Methods("GET")
	basePathRouter.HandleFunc("/fe-config.mjs", routerProvider.FrontendConfigApi.GetMicroFrontendConfigsAsJavaScritp).Methods("GET")
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
