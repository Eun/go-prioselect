package prioselect_test

import (
	"fmt"
	"log"
	"os"
	"time"

	prioselect "github.com/Eun/go-prioselect"
)

func ExampleSelect() {
	timeoutChan := time.After(time.Second * 5)
	textChan := make(chan string)

	go func() {
		fmt.Printf("Enter any text: ")
		var text string
		fmt.Scan(&text)
		textChan <- text
	}()

	// if textChan and timeoutChan trigger at the same time
	// textChan will be prioritized
	value, channel := prioselect.Select(textChan, timeoutChan)
	switch channel {
	case timeoutChan:
		log.Fatalln("timeout")
	case textChan:
		fmt.Printf("Text is %s", value.(string))
		os.Exit(0)
	case nil:
		log.Fatalln("unexpected error")
	}
}
