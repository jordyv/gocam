package gocamhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gobuffalo/packr"
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

type alertApiItem struct {
	Timestamp time.Time `json:"timestamp"`
	ImageName string    `json:"image_name"`
}

func (h *GocamHttp) Listen() {
	box := packr.NewBox("./webui")

	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(box))
	handler.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir(h.Options.config.ImagePath))))
	handler.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		imageFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/alert", h.Options.config.ImagePath))
		if err != nil {
			log.Errorln("error reading image directory:", err)
			return
		}
		files := make([]alertApiItem, 0)
		for _, image := range imageFiles {
			if !image.IsDir() {
				files = append(files, alertApiItem{image.ModTime(), image.Name()})
			}
		}
		json, err := json.Marshal(files)
		if err != nil {
			log.Errorln("error while serializing to JSON", err)
			return
		}
		fmt.Fprint(w, string(json))
	})
	log.Fatalln(http.ListenAndServe(h.Options.config.HTTPAddr, handler))
}
