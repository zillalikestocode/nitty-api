package api

import (
	"net/http"

	application "github.com/zillalikestocode/community-api/app"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	r := application.LoadRoutes()
	r.ServeHTTP(w, req)
}
