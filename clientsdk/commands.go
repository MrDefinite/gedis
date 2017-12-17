package clientsdk

import (
	"strings"
	"github.com/MrDefinite/gedis/common/protocol/resp"
)

func (gc *Gclient) Set(params []string) (string, error) {

	plainCmd := "set "
	for _, param := range params {
		plainCmd += param + " "
	}

	plainCmd = strings.TrimSpace(plainCmd)

	encodedCmd := resp.EncodeCmd(plainCmd)
	response, err := gc.sendRequest(encodedCmd)

	if err != nil {
		return "", err
	}
	response = resp.ParseResponse(response)

	return response, nil
}

func (gc *Gclient) Get(params []string) string {
	return ""
}
