package main

import (
	"fmt"
	"os"
)

func write_file(str string) {
	f, err := os.OpenFile("/data/data.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		fmt.Printf("file open error: %v\n", err)
	}
	if _, err = f.WriteString(str); err != nil {
		fmt.Printf("file write error: %v\n", err)
	}
}
