package domain

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
