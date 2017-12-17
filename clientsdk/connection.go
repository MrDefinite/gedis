package clientsdk

import (
	"errors"
	"time"
)

func (gc *Gclient) heartbeat()  {

}

func (gc *Gclient) sendRequest(encodedRequest string) (string, error) {
	if encodedRequest == "" {
		return "", errors.New("cannot send empty request to server")
	}

	if gc.isCommunicatingToServer {
		return "", errors.New("there is another request being processed")
	}

	gc.commLock.Lock()

	// Send it to server now
	gc.conn.Write([]byte(encodedRequest))

	// Init buffer
	buff := make([]byte, 1024)

	gc.conn.SetReadDeadline(time.Now().Add(time.Duration(gc.requestTimeout)))
	// Wait for response from server
	n, err := gc.conn.Read(buff)
	if err != nil {
		return "", err
	}
	gc.commLock.Unlock()

	return string(buff[:n]), nil
}

