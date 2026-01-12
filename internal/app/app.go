package app

import (
	"net/http"

	"github.com/envcrypts/envcrypt-cli/internal/client"
)

type App struct {
	HttpClient *client.Client
}

func NewApp(baseUrl string) *App {
	httpClient := client.NewClient(baseUrl, &http.Client{})
	return &App{
		HttpClient: httpClient,
	}
}
