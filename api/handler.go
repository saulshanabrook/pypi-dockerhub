package api

import (
	"encoding/json"
	"net/http"

	"github.com/saulshanabrook/pypi-dockerhub/db"
)

func CreateHandler(c *db.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rels, err := c.GetReleases()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bs, err := json.Marshal(&rels)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(bs)
	}
}
