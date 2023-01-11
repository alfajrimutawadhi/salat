package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alfajrimutawadhi/salat/common"
	"github.com/alfajrimutawadhi/salat/constant"
	"github.com/alfajrimutawadhi/salat/salat"
	"github.com/jedib0t/go-pretty/v6/table"
)

var version string = "v1.0.1"

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
			var loc salat.Location
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
		case constant.TimeMode:
			fmt.Println("Choose time mode format")
			var m string
			fmt.Println("1. 12 Hours (AM/PM)")
			fmt.Println("2. 24 Hours (23:59)")
			fmt.Print("input time mode (1/2) : ")
			fmt.Scan(&m)
			setTimeMode(path, m)
		default:
			fmt.Println(constant.Help)
		}
	}

}

func scheduleNow(path string) {
	ti := time.Now()
	h, m, _ := ti.Clock()
	c := salat.ReadConfig(path)
	if c.TimeMode == 2 {
		fmt.Printf("Time now = %d:%d\n", h, m)
	} else {
		fmt.Printf("Time now = %s\n", common.ConvTime24To12(fmt.Sprintf("%d:%d", h, m)))
	}

	var sc salat.Salat
	res := salat.RequestAPI(&c.Location, "s", ti)

	var tmp salat.Response
	if err := json.Unmarshal(res, &tmp); err != nil {
		common.HandleError(err)
	}
	sc = tmp.Data.Timings

	sc.SetTimeMode(c.TimeMode)

	fmt.Println("Location =", c.Location.City)
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
	c := salat.ReadConfig(path)
	l := fmt.Sprintf("Country = %s\nCity = %s", c.Location.Country, c.Location.City)
	fmt.Println(l)
}

func setLocation(path string, req *salat.Location) {
	defer fmt.Println("Successfully changed location")
	c := salat.ReadConfig(path)

	d := salat.Config{
		Location: salat.Location{Country: req.Country, City: req.City},
		TimeMode: c.TimeMode,
	}
	cf, err := json.Marshal(d)
	if err != nil {
		common.HandleError(err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s/.salat/config.json", path), []byte(cf), os.ModeAppend.Perm()); err != nil {
		common.HandleError(err)
	}
}

func dateHijri(path string) {
	ti := time.Now()
	var d salat.Date

	c := salat.ReadConfig(path)
	// easier with use API prayer schedule
	res := salat.RequestAPI(&c.Location, "s", ti)

	var tmp salat.Response
	json.Unmarshal(res, &tmp)
	d.Hijri = tmp.Data.Date.Hijri

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Current date"})
	t.AppendRows([]table.Row{{fmt.Sprintf("%s %s %sH", d.Hijri.Day, d.Hijri.Month.En, d.Hijri.Year)}})
	t.Render()
}

func calendarHijri(path string) {
	var d []salat.Date

	ti := time.Now()
	c := salat.ReadConfig(path)
	res := salat.RequestAPI(&c.Location, "c", ti)

	var tmp salat.ResponseDate
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

func handleDate(d []salat.Date) []table.Row {
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

func setTimeMode(path string, req string) {
	defer fmt.Println("Successfully changed time mode")
	// validate
	if req != "1" && req != "2" {
		err := errors.New("invalid input")
		common.HandleError(err)
	}
	tm, err := strconv.Atoi(req)
	if err != nil {
		common.HandleError(err)
	}

	c := salat.ReadConfig(path)
	d := salat.Config{
		Location: c.Location,
		TimeMode: int8(tm),
	}
	cf, err := json.Marshal(d)
	if err != nil {
		common.HandleError(err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s/.salat/config.json", path), []byte(cf), os.ModeAppend.Perm()); err != nil {
		common.HandleError(err)
	}
}
