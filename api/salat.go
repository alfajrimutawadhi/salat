package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alfajrimutawadhi/salat/common"
	"github.com/alfajrimutawadhi/salat/constant"
	"github.com/alfajrimutawadhi/salat/domain"
)

// returns response API in byte
func RequestAPI(cl *domain.Location, nec string, t time.Time) []byte {
	var url string

	if nec == "d" {
		url = fmt.Sprintf("https://%s/calendarByCity?country=%s&city=%s&method=1&month=%d&year=%d", constant.HeaderHost, cl.Country, cl.City, int(t.Month()), t.Year())
	} else {
		url = fmt.Sprintf("https://%s/timingsByCity?country=%s&city=%s&method=1", constant.HeaderHost, cl.Country, cl.City)
	}

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
