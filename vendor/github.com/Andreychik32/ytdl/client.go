package ytdl

import (
	"net/http"

	"github.com/rs/zerolog"
)

// Client export
type Client struct {
	Logger     zerolog.Logger
	HTTPClient *http.Client
}
// DefaultClient check
var DefaultClient = &Client{
	HTTPClient: http.DefaultClient,
	Logger:     zerolog.Nop(),
}
