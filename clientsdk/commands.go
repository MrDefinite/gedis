package clientsdk

import (
	"strings"
	"github.com/MrDefinite/gedis/common/protocol/resp"
	"errors"
)

const (
	setCommandName = "set"
	getCommandName = "get"
)


func (gc *Gclient) ProcessCmd(cmd string) (string, error) {
	requestCmd, err := gc.ParseCmd(cmd)
	if err != nil {
		return "", err
	}

	switch requestCmd.Name {
	case getCommandName:
		return gc.Get(requestCmd.Params)
	case setCommandName:
		return gc.Set(requestCmd.Params)
	}

	return "", errors.New("no corresponding command found")
}

func (gc *Gclient) ParseCmd(cmd string) (*resp.RequestCmd, error) {
	if cmd == "" {
		return nil, errors.New("illegal command")
	}

	cmd = strings.TrimSpace(cmd)
	cmdArray := strings.Split(cmd, " ")

	if len(cmdArray) < 1 {
		return nil, errors.New("illegal command")
	}

	cmdName := strings.ToLower(cmdArray[0])
	var cmdParams []string
	for i := 1; i < len(cmdArray); i++ {
		cmdParams = append(cmdParams, cmdArray[i])
	}

	requestCmd := resp.RequestCmd{
		Name: cmdName,
		Params: cmdParams,
	}

	return &requestCmd, nil
}

func buildAndEncodeCmd(cmdName string, params []string) string {
	plainCmd := cmdName + " "
	for _, param := range params {
		plainCmd += param + " "
	}
	plainCmd = strings.TrimSpace(plainCmd)

	encodedCmd := resp.EncodeCmd(plainCmd)

	return encodedCmd
}

func processRequest(gc *Gclient, encodedCmd string) (string, error) {
	response, err := gc.sendRequestAndGetResponse(encodedCmd)
	if err != nil {
		return "", err
	}
	response = resp.ParseResponse(response)

	return response, nil
}

func (gc *Gclient) Set(params []string) (string, error) {
	if params == nil || len(params) == 0 {
		return "", errors.New("illegal command")
	}

	encodedCmd := buildAndEncodeCmd("set", params)
	return processRequest(gc, encodedCmd)
}

func (gc *Gclient) Get(params []string) (string, error) {
	if params == nil || len(params) == 0 {
		return "", errors.New("illegal command")
	}

	encodedCmd := buildAndEncodeCmd("get", params)
	return processRequest(gc, encodedCmd)
}
