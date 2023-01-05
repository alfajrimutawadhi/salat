package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alfajrimutawadhi/salat/api"
	"github.com/alfajrimutawadhi/salat/common"
	"github.com/alfajrimutawadhi/salat/constant"
	"github.com/alfajrimutawadhi/salat/domain"
	"github.com/jedib0t/go-pretty/v6/table"
)

var version string = "v1.0.0"

func main() {
	// path := os.Getenv("HOME")
	path, _ := os.Getwd()

	args := os.Args
	if len(args) == 1 {
		fmt.Print(constant.Help)
	} else {
		switch args[1] {
		case constant.Version:
			fmt.Printf("Salat %s", version)
		case constant.ScheduleNow:
			scheduleNow(path)
		case constant.ShowLocation:
			showLocation(path)
		case constant.SetLocation:
			var err error
			var loc domain.Location
			fmt.Print("input your country : ")
			in := bufio.NewReader(os.Stdin)
			loc.Country, err = in.ReadString('\n')
			if err != nil {
				common.HandleError(err)
			}
			fmt.Print("input your city : ")
			loc.City, err = in.ReadString('\n')
			if err != nil {
				common.HandleError(err)
			}
			loc.Country = strings.ReplaceAll(loc.Country, "\n", "")
			loc.City = strings.ReplaceAll(loc.City, "\n", "")

			setLocation(path, &loc)
		case constant.Date:
			dateHijri(path)
		default:
			fmt.Println(constant.Help)

		}
	}

}

func scheduleNow(path string) {
	h, m, s := time.Now().Clock()
	fmt.Printf("Time now = %d:%d:%d\n", h, m, s)

	var sc domain.Salat

	cl := common.ReadLocation(path)
	res := api.RequestAPI(cl)

	var tmp domain.Response
	json.Unmarshal(res, &tmp)
	sc = tmp.Data.Timings

	fmt.Println("Location =", cl.City)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Salat", "Schedule"})
	t.AppendRows([]table.Row{
		{"Fajr", sc.Fajr},
		{"Dhuhr", sc.Dhuhr},
		{"Asr", sc.Asr},
		{"Maghrib", sc.Maghrib},
		{"Isha", sc.Isha},
	})
	t.Render()
}

func showLocation(path string) {
	fmt.Println("Your current location")
	b, err := os.ReadFile(fmt.Sprintf("%s/location.toml", path))
	if err != nil {
		common.HandleError(err)
	}
	fmt.Println(string(b))
}

func setLocation(path string, req *domain.Location) {
	defer fmt.Println("Successfully changed location")
	sl := []string{
		"country = " + req.Country,
		"city = " + req.City,
	}

	nl := strings.Join(sl, "\n")
	if err := os.WriteFile(fmt.Sprintf("%s/location.toml", path), []byte(nl), os.ModeAppend); err != nil {
		common.HandleError(err)
	}
}

func dateHijri(path string) {
	var d domain.Date

	cl := common.ReadLocation(path)
	res := api.RequestAPI(cl)

	var tmp domain.Response
	json.Unmarshal(res, &tmp)
	d.Hijri.Day = tmp.Data.Date.Hijri.Day
	d.Hijri.Month = tmp.Data.Date.Hijri.Month
	d.Hijri.Year = tmp.Data.Date.Hijri.Year

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Current date"})
	t.AppendRows([]table.Row{{fmt.Sprintf("%s %s %sH", d.Hijri.Day, d.Hijri.Month.En, d.Hijri.Year)}})
	t.Render()
}
