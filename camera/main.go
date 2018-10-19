package camera

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client  Client for IP camera communication
type Client struct {
	imageURL string
	savePath string
}

// NewClient  Create new client
func NewClient(imageURL string, savePath string) *Client {
	return &Client{imageURL: imageURL, savePath: savePath}
}

func (c *Client) getImage() ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	response, err := client.Get(c.imageURL)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("couldn't get the image from the camera, please check your URL and credentials")
	}
	return ioutil.ReadAll(response.Body)
}

// SaveImage  Get image from IP camera and save it on disk
func (c *Client) SaveImage() (fileName string, err error) {
	output, err := c.getImage()
	if err != nil {
		return
	}
	fileName = fmt.Sprintf("%s/%s_%s", c.savePath, time.Now().Format(time.RFC3339), "image.jpg")
	err = ioutil.WriteFile(fileName, output, 0655)
	if err != nil {
		return
	}
	return
}
