package urlshort

import (
	"net/http"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the path from the incoming URL
		path := r.URL.Path
		// check if the path exist in our mapping
		if dest, ok := pathsToUrls[path]; ok {
			// if the path exists, we redirect the user to that URL (dest)
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}
