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

func (app *App) SetBaseURL(url string) {
	app.HttpClient.SetBaseURL(url)
}

func (app *App) BaseURL() string {
	return app.HttpClient.BaseURL()
}
