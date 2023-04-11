package salat

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
)

type Salat struct {
	Fajr    string `json:"Fajr"`
	Dhuhr   string `json:"Dhuhr"`
	Asr     string `json:"Asr"`
	Maghrib string `json:"Maghrib"`
	Isha    string `json:"Isha"`
}

type Response struct {
	Code   int      `json:"code"`
	Status string   `json:"status"`
	Data   MetaData `json:"data"`
}

type ResponseDate struct {
	Code   int        `json:"code"`
	Status string     `json:"status"`
	Data   []MetaData `json:"data"`
}

type MetaData struct {
	Timings Salat       `json:"timings"`
	Date    Date        `json:"date"`
	Meta    interface{} `json:"meta"` // ignore response meta
}

type Date struct {
	Gregorian struct {
		Day     string `json:"day"`
		Weekday struct {
			En string `json:"en"`
		} `json:"weekday"`
	} `json:"gregorian"`

	Hijri struct {
		Day   string `json:"day"`
		Month struct {
			Number int    `json:"number"`
			En     string `json:"en"`
		} `json:"month"`
		Year string `json:"year"`
	} `json:"hijri"`
}

// returns response API in byte.
// nec c for calendar, s for salat
func RequestAPI(cl *Location, nec string, t time.Time) []byte {
	var url string

	if nec == "c" {
		url = fmt.Sprintf("https://%s/calendarByCity?country=%s&city=%s&method=3&month=%d&year=%d", constant.HeaderHost, cl.Country, cl.City, int(t.Month()), t.Year())
	} else {
		url = fmt.Sprintf("https://%s/timingsByCity?country=%s&city=%s&method=3", constant.HeaderHost, cl.Country, cl.City)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		common.HandleError(err)
	}

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
