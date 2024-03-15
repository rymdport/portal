package main

import (
	"fmt"
	"os"

	"github.com/rymdport/portal/trash"
)

func main() {

	file, _ := os.OpenFile("test.html", os.O_RDWR, 0)

	fmt.Println(trash.TrashFile(file.Fd()))
}
