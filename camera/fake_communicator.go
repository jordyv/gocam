package camera

import (
	"fmt"
	"os"
)

type FakeCommunicator struct{}

// NewFakeCommunicator  Create new fake camera communicator
func NewFakeCommunicator() FakeCommunicator {
	return FakeCommunicator{}
}

// SaveImage  Get image from IP camera and save it on disk
func (c FakeCommunicator) SaveImage() (string, error) {
	dir, _ := os.Getwd()
	return fmt.Sprintf("%s/test/images/same/1.jpg", dir), nil
}
