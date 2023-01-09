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
	path := os.Getenv("HOME")

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
		case constant.Calendar:
			calendarHijri(path)
		default:
			fmt.Println(constant.Help)
		}
	}

}

func scheduleNow(path string) {
	ti := time.Now()
	h, m, s := ti.Clock()
	fmt.Printf("Time now = %d:%d:%d\n", h, m, s)

	var sc domain.Salat

	cl := common.ReadLocation(path)
	res := api.RequestAPI(cl, "s", ti)

	var tmp domain.Response
	if err := json.Unmarshal(res, &tmp); err != nil {
		common.HandleError(err)
	}
	sc = tmp.Data.Timings

	fmt.Println("Location =", cl.City)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Salat", "Schedule"})
	t.AppendRows([]table.Row{
		{constant.Fajr, sc.Fajr},
		{constant.Dhuhr, sc.Dhuhr},
		{constant.Asr, sc.Asr},
		{constant.Maghrib, sc.Maghrib},
		{constant.Isha, sc.Isha},
	})
	t.Render()

}

func showLocation(path string) {
	fmt.Println("Your current location")
	lc := common.ReadLocation(path)
	l := fmt.Sprintf("Country = %s\nCity = %s", lc.Country, lc.City)
	fmt.Println(l)
}

func setLocation(path string, req *domain.Location) {
	defer fmt.Println("Successfully changed location")
	sl := []string{
		"country = " + req.Country,
		"city = " + req.City,
	}

	nl := strings.Join(sl, "\n")
	
	if err := os.WriteFile(fmt.Sprintf("%s/.salat/location.toml", path), []byte(nl), os.ModeAppend.Perm()); err != nil {
		common.HandleError(err)
	}
}

func dateHijri(path string) {
	ti := time.Now()
	var d domain.Date

	cl := common.ReadLocation(path)
	// easier with use API sholat
	res := api.RequestAPI(cl, "s", ti)

	var tmp domain.Response
	json.Unmarshal(res, &tmp)
	d.Hijri = tmp.Data.Date.Hijri

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Current date"})
	t.AppendRows([]table.Row{{fmt.Sprintf("%s %s %sH", d.Hijri.Day, d.Hijri.Month.En, d.Hijri.Year)}})
	t.Render()
}

func calendarHijri(path string) {
	var d []domain.Date

	ti := time.Now()
	cl := common.ReadLocation(path)
	res := api.RequestAPI(cl, "d", ti)

	var tmp domain.ResponseDate
	if err := json.Unmarshal(res, &tmp); err != nil {
		common.HandleError(err)
	}
	for i := range tmp.Data {
		d = append(d, tmp.Data[i].Date)
	}

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRow(table.Row{
		"",
		ti.Year(),
		ti.Year(),
		ti.Year(),
		d[0].Hijri.Year + "H",
		d[0].Hijri.Year + "H",
		"",
	}, rowConfigAutoMerge)

	// check hijri month name
	me := d[0].Hijri.Month.En
	for i := range d {
		if d[i].Hijri.Month.En != me {
			me = fmt.Sprintf("%s - %s", me, d[i].Hijri.Month.En)
			break
		}
	}
	t.AppendRow(table.Row{
		"",
		ti.Month().String(),
		ti.Month().String(),
		ti.Month().String(),
		me,
		me,
		"",
	}, rowConfigAutoMerge)

	t.AppendSeparator()
	t.AppendRow(table.Row{
		constant.Sun,
		constant.Mon,
		constant.Tue,
		constant.Wed,
		constant.Thu,
		constant.Fri,
		constant.Sat,
	})
	t.AppendSeparator()

	dt := handleDate(d)

	t.AppendRows(dt)
	t.Render()
}

func handleDate(d []domain.Date) []table.Row {
	var dt []table.Row
	sd := []string{
		constant.Sunday,
		constant.Monday,
		constant.Tuesday,
		constant.Wednesday,
		constant.Thursday,
		constant.Friday,
		constant.Saturday,
	}

	for i := 0; i < len(d); i++ {
		var row []interface{}
		for j := 0; j < len(sd); j++ {
			if i < len(d) {
				if d[i].Gregorian.Weekday.En == sd[j] {
					row = append(row, fmt.Sprintf("%s (%s)", d[i].Gregorian.Day, d[i].Hijri.Day))
					i++
				} else {
					row = append(row, "")
				}
			} else {
				break
			}
		}
		i--
		dt = append(dt, row)
	}

	return dt
}
