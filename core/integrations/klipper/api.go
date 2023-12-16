package klipper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *KipplerPrinter) serverInfo() (*Result, error) {
	res, err := http.Get(fmt.Sprintf("%s/server/info", p.Address))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r MoonRakerResponse
	err = decoder.Decode(&r)
	if err != nil {
		return nil, err
	}

	return r.Result, nil
}
