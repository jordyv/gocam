package gocamhttp

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/jordyv/gocam/config"
)

type GocamHttp struct {
	Options *GocamHttpOptions
}

type GocamHttpOptions struct {
	config *config.Config
}

func NewGocamHttp(config *config.Config) *GocamHttp {
	return &GocamHttp{Options: &GocamHttpOptions{config: config}}
}

func (h *GocamHttp) Listen() {
	logrus.Fatalln(http.ListenAndServe(h.Options.config.HTTPAddr, http.FileServer(http.Dir(h.Options.config.ImagePath))))
}
