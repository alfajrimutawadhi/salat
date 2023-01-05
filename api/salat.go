package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/alfajrimutawadhi/salat/common"
	"github.com/alfajrimutawadhi/salat/constant"
	"github.com/alfajrimutawadhi/salat/domain"
)

// returns response API in byte
func RequestAPI(cl *domain.Location) []byte {
	url := fmt.Sprintf("https://%s/timingsByCity?country=%s&city=%s&method=1", constant.HeaderHost, cl.Country, cl.City)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		common.HandleError(err)
	}
	req.Header.Add("X-RapidAPI-Key", constant.HeaderKey)
	req.Header.Add("X-RapidAPI-Host", constant.HeaderHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			err = errors.New("unable connect to internet")
		}
		common.HandleError(err)
	}

	if res.StatusCode != 200 {
		fmt.Printf("Status code : %d\n", res.StatusCode)
		fmt.Println("make sure your location settings are correct")
		os.Exit(os.SEEK_END)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		common.HandleError(err)
	}
	return body
}
