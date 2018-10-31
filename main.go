package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jordyv/gocam/alerting"
	"github.com/jordyv/gocam/camera"
	cfg "github.com/jordyv/gocam/config"
	"github.com/jordyv/gocam/hasher"
	"github.com/jordyv/gocam/http"
	"github.com/jordyv/gocam/metrics"
	sys "golang.org/x/sys/unix"
)

var (
	config *cfg.Config

	client         *camera.Client
	hashCalculator *hasher.Hasher
	alertManager   *alerting.AlertManager

	previousHash      *hasher.ImageHash
	previousImagePath string
)

func cleanUp() {
	files, _ := ioutil.ReadDir(config.ImagePath)
	if len(files) > config.MaxKeepedImageFiles {
		for _, file := range files[0 : len(files)-config.MaxKeepedImageFiles] {
			sys.Unlink(fmt.Sprintf("%s/%s", config.ImagePath, file.Name()))
		}
	}
}

func alertNotify(imagePath string, previousImagePath string, distance int) {
	newFileName := fmt.Sprintf("%s.jpg", time.Now().Format(time.RFC3339))
	newImagePath := fmt.Sprintf("%s/%s", config.AlertImagePath, newFileName)
	err := sys.Link(imagePath, newImagePath)
	if err != nil {
		log.Panicf("cannot copy %s to %s: %s", imagePath, newImagePath, err)
		return
	}
	go alertManager.Notify(alerting.NewAlert(previousImagePath, imagePath, distance))
}

func buildAlertHandlers() []alerting.AlertHandler {
	handlers := make([]alerting.AlertHandler, 0)
	if config.HasAlertHandler("console") {
		handlers = append(handlers, &alerting.LogAlertHandler{Options: config.GetLogAlertHandlerOptions()})
	}
	if config.HasAlertHandler("diff") {
		handlers = append(handlers, &alerting.DiffImageAlertHandler{Options: config.GetDiffImageAlertHandlerOptions()})
	}
	if config.HasAlertHandler("telegram") {
		handlers = append(handlers, &alerting.TelegramAlertHandler{Options: config.GetTelegramAlertHandlerOptions()})
	}
	return handlers
}

func initConfig() {
	config = &cfg.Config{}
	config.Build()
}

func runCycle() {
	filePath, err := client.SaveImage()
	if err != nil {
		log.Panicln("could not save image from camera;", err)
		os.Exit(1)
	}
	log.Debugln("Saved image from camera at:", filePath)

	currentHash, err := hashCalculator.CalculateHash(filePath)
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
	log.Debugln("Calculated hash for image:", currentHash)

	if previousHash != nil {
		distance, _ := hashCalculator.CalculateDistance(previousHash, currentHash)
		log.Debugln("Distance with previous image:", distance)
		if distance >= config.Treshold {
			alertNotify(filePath, previousImagePath, distance)
		}
	}
	previousHash = currentHash
	previousImagePath = filePath

	cleanUp()
}

func main() {
	initConfig()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(config.LogLevel)

	client = camera.NewClient(camera.NewIPCameraCommunicator(config.CameraURL, config.ImagePath))
	hashCalculator = hasher.New()
	alertManager = alerting.New(buildAlertHandlers())

	if config.HTTPEnabled {
		httpListener := gocamhttp.NewGocamHttp(config)
		log.Infoln("Start HTTP server at", config.HTTPAddr)
		go httpListener.Listen()
	}
	if config.MetricsEnabled {
		m := metrics.New(config)
		alertManager.AddAlertHandler(m)
		log.Infoln("Start metrics endpoint", config.MetricsAddr)
		go m.Listen()
	}

	go func() {
		for {
			runCycle()
			time.Sleep(config.Interval)
		}
	}()
	select {}
}
