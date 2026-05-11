package app

import (
	"net/http"
	"os"
	"time"

	"github.com/envcrypts/envcrypt-cli/internal/client"
)

type App struct {
	HttpClient *client.Client
}

const defaultHTTPTimeout = 30 * time.Second

func NewApp(baseUrl string) *App {
	timeout := defaultHTTPTimeout
	if envTimeout := os.Getenv("ENVCRYPT_HTTP_TIMEOUT"); envTimeout != "" {
		if parsed, err := time.ParseDuration(envTimeout); err == nil {
			timeout = parsed
		}
	}

	httpClient := client.NewClient(baseUrl, &http.Client{Timeout: timeout})
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
