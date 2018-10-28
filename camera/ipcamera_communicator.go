package camera

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type IPCameraCommunicator struct {
	imageURL string
	savePath string
}

// NewIPCameraCommunicator  Create new IP camera communicator
func NewIPCameraCommunicator(imageURL string, savePath string) IPCameraCommunicator {
	return IPCameraCommunicator{imageURL: imageURL, savePath: savePath}
}

func (c *IPCameraCommunicator) getImage() ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Get(c.imageURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("couldn't get the image from the camera, please check your URL and credentials")
	}
	return ioutil.ReadAll(response.Body)
}

// SaveImage  Get image from IP camera and save it on disk
func (c IPCameraCommunicator) SaveImage() (filePath string, err error) {
	output, err := c.getImage()
	if err != nil {
		return
	}
	filePath = fmt.Sprintf("%s/%s_%s", c.savePath, time.Now().Format(time.RFC3339), "image.jpg")
	err = ioutil.WriteFile(filePath, output, 0655)
	if err != nil {
		return
	}
	return
}
