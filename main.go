package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logJSONHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logJSONHandler)
	slog.SetDefault(logger)
	http.HandleFunc("/")
	logger.Info("Proxy server started on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("can not start server", slog.String("error", err.Error()))
		os.Exit(1)
	}

}
