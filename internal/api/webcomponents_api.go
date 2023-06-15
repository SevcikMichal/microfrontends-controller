package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/SevcikMichal/microfrontends-controller/internal/provider"
)

type WebComponentApi struct {
	MicroFrontendProvider *provider.MicroFrontendProvider
	Client                *http.Client
}

func (api *WebComponentApi) GetWebComponent(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get web component started.")
	namespace, name, resource := parseWebComponentPath(r.URL.Path)

	requestModuleUri := api.MicroFrontendProvider.GetMicrofrontendRequestModuleUri(namespace, name)
	realModuleUri := api.MicroFrontendProvider.GetMicrofrontendModuleUri(namespace, name)

	if requestModuleUri == "" || realModuleUri == "" {
		http.NotFound(w, r)
		return
	}

	proxyUrl := realModuleUri
	if requestModuleUri != r.URL.Path {
		proxyUrl = addResourceToLastUrlSegment(realModuleUri, resource)
	}

	req, err := http.NewRequest("GET", proxyUrl, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Proxying request to the module.", "Resolved URL:", proxyUrl)
	resp, err := api.Client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.Header().Set("Content-Type", "application/javascript")

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func parseWebComponentPath(path string) (namespace, name, resource string) {
	segments := strings.Split(path, "/")
	return segments[2], segments[3], strings.Join(segments[4:], "/")
}

func addResourceToLastUrlSegment(urlPath, resource string) string {
	u, _ := url.Parse(urlPath)
	u.Path = path.Join(path.Dir(u.Path), resource)
	return u.String()
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		dst.Set(key, strings.Join(values, ", "))
	}
}
