package salat

import "github.com/alfajrimutawadhi/salat/common"

func (s *Salat) SetTimeMode(ts int8) Salat {
	if ts == 2 {
		return *s
	}
	s.Fajr = common.ConvTime24To12(s.Fajr)
	s.Dhuhr = common.ConvTime24To12("12:01")
	s.Asr = common.ConvTime24To12(s.Asr)
	s.Maghrib = common.ConvTime24To12(s.Maghrib)
	s.Isha = common.ConvTime24To12(s.Isha)
	return *s
}
