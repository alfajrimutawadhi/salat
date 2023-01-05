package common

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alfajrimutawadhi/salat/domain"
)

func ReadLocation(path string) *domain.Location {
	var loc domain.Location
	cf, err := os.ReadFile(fmt.Sprintf("%s/location.toml", path))
	if err != nil {
		HandleError(err)
	}
	lines := strings.Split(string(cf), "\n")
	if len(lines) < 2 {
		HandleError(errors.New("file location.toml is out of format"))
	}

	loc.Country = strings.Split(lines[0], " = ")[1]
	loc.City = strings.Split(lines[1], " = ")[1]

	return &loc
}
