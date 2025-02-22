package route

import (
	"HelloGolang/internal/oss"
	"HelloGolang/internal/service"
	"encoding/json"
	"net/http"
)

func Article(w http.ResponseWriter, r *http.Request) {
	codes := []string{
		"202501193301333114",
		"202501193301288337",
		"202501193301321338",
	}
	result, err := service.EastmoneyArticleService(codes)
	if err != nil {
		http.Error(w, "服务器错误", 500)
		return
	}
	fileHelp := oss.FileHelper{}
	go fileHelp.WriteToFile("article-2025-01-19.json", result)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
