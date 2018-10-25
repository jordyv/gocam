package gocamhttp

import (
	"encoding/json"
	"errors"
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

func getAllItems(imagePath string) ([]alertApiItem, error) {
	imageFiles, err := ioutil.ReadDir(fmt.Sprintf("%s/alert", imagePath))
	if err != nil {
		log.Errorln("error reading image directory:", err)
		return nil, errors.New("could not read images in directory")
	}
	files := make([]alertApiItem, 0)
	for _, image := range imageFiles {
		if !image.IsDir() {
			files = append(files, alertApiItem{image.ModTime(), image.Name()})
		}
	}
	return files, nil
}

func getAllGroupedItems(imagePath string) (map[time.Time]([]alertApiItem), error) {
	items, err := getAllItems(imagePath)
	if err != nil {
		return nil, err
	}
	m := make(map[time.Time]([]alertApiItem), 0)
	for _, i := range items {
		groupTimestamp := time.Date(i.Timestamp.Year(), i.Timestamp.Month(), i.Timestamp.Day(), int(0), int(0), int(0), int(0), time.Local)
		m[groupTimestamp] = append(m[groupTimestamp], i)
	}

	return m, nil
}

func (h *GocamHttp) Listen() {
	box := packr.NewBox("./webui")

	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(box))
	handler.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir(h.Options.config.ImagePath))))
	handler.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		files, _ := getAllItems(h.Options.config.ImagePath)
		json, _ := json.Marshal(files)
		fmt.Fprint(w, string(json))
	})
	handler.HandleFunc("/api/alerts/grouped", func(w http.ResponseWriter, r *http.Request) {
		files, _ := getAllGroupedItems(h.Options.config.ImagePath)
		json, _ := json.Marshal(files)
		fmt.Fprint(w, string(json))
	})
	log.Fatalln(http.ListenAndServe(h.Options.config.HTTPAddr, handler))
}
