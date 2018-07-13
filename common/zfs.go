package common

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"github.com/ipfs/go-ipfs-api"
)

func GetSysInfo(gateway string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v0/diag/sys", gateway))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func GetID(gateway string) (*shell.IdOutput, error) {
	var sh = shell.NewShell(gateway)
	return sh.ID()
}
