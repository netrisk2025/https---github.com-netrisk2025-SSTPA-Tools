package apihttp

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func healthHandler(version string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(healthResponse{
			Status:  "ok",
			Service: "sstpa-api",
			Version: version,
		})
	}
}
