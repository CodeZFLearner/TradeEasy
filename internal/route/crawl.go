package route

import (
	"HelloGolang/internal/service"
	"encoding/json"
	"net/http"
)

func Crawl(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	result, err := service.CrawlByCategory(query)
	if err != nil {
		http.Error(w, "服务器错误", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
