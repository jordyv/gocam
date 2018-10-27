package camera

// Client  Client for camera communication
type Client struct {
	Communicator
}

// Communicator  Communicators actually communicating with the IP camera implementation, eg. IP cams, RPi's cams, etc.
type Communicator interface {
	SaveImage() (string, error)
}

// NewClient  Create a new client
func NewClient(communicator Communicator) *Client {
	return &Client{communicator}
}
