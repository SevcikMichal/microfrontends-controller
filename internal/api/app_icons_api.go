package api

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SevcikMichal/microfrontends-controller/internal/provider"
)

type AppIconsApi struct {
	MicroFrontendProvider *provider.MicroFrontendProvider
}

func (api *AppIconsApi) GetAppIcon(w http.ResponseWriter, r *http.Request) {
	navigationPath := parseNavigationPath(r.URL.Path)
	appIcon := api.MicroFrontendProvider.GetMicrofrontendAppIcon(navigationPath)

	if appIcon == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(appIcon.Data) == 0 {
		response, err := http.Get(appIcon.Url)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer response.Body.Close()

		var data []byte
		_, readErr := response.Body.Read(data)
		if readErr != nil && readErr != io.EOF {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appIcon.Data = base64.StdEncoding.EncodeToString(data)

		copyHeaders(w.Header(), response.Header)
		w.Header().Set("Content-Type", appIcon.Mime)
		w.WriteHeader(response.StatusCode)
		io.Copy(w, response.Body)

		return
	}

	w.Header().Set("Content-Type", appIcon.Mime)

	decoded, err := base64.StdEncoding.DecodeString(appIcon.Data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(decoded)
}

func parseNavigationPath(path string) string {
	segments := strings.Split(path, "/")
	return strings.Join(segments[(len(segments)-1):], "/")
}
