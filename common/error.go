package common

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	fmt.Print("Error : ", err.Error())
	os.Exit(os.SEEK_END)
}