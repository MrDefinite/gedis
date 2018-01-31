package clientsdk

import (
	"errors"
	"time"
	"github.com/MrDefinite/gedis/common/protocol/resp"
)

func (gc *Gclient) heartbeat()  {

}

func (gc *Gclient) sendRequestAndGetResponse(encodedRequest []byte) ([]byte, error) {
	if encodedRequest == nil {
		return nil, errors.New("cannot send empty request to server")
	}

	if gc.isCommunicatingToServer {
		return nil, errors.New("there is another request being processed")
	}

	gc.commLock.Lock()

	// Send it to server now
	gc.conn.Write(encodedRequest)

	// Init buffer
	buff := make([]byte, resp.MaxDataSizeReadPerTime)

	gc.conn.SetReadDeadline(time.Now().Add(time.Duration(gc.requestTimeout) * time.Second))
	// Wait for response from server

	var n int
	var err error
	//var parser resp.
	for n, err = gc.conn.Read(buff);; {
		if err != nil {
			gc.commLock.Unlock()
			return nil, err
		}


	}

	gc.commLock.Unlock()

	return string(buff[:n]), nil
}

func parseBuff()  {

}
