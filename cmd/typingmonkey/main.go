package main

import (
	"fmt"
	"time"

	"typingmonkey/internal/bot"
	"typingmonkey/internal/driver"
	"typingmonkey/internal/killswitch"
)

const (
	/*
	 port for geckodriver to run on
	*/
	GECKO_PORT = 4444
	/*
	 delay too little may cause missed keystrokes
	 8 is fine results in about ~780 WPM ≧◡≦
	*/
	TYPING_DELAY = 8 * time.Millisecond
)

func main() {
	go killswitch.Monitor()

	gecko, err := driver.New(GECKO_PORT)
	if err != nil {
		fmt.Printf("oopsies!: %v\n", err)
		return
	}

	fmt.Println("starting geckodriver...")
	if err := gecko.Start(); err != nil {
		fmt.Printf("oopsies!: %v\n", err)
		return
	}
	defer gecko.Stop()

	fmt.Println("talking to Firefox...")
	if err := gecko.Connect(); err != nil {
		fmt.Printf("oopsies!: %v\n", err)
		return
	}

	monkeybot := bot.New(gecko.WebDriver(), TYPING_DELAY)

	fmt.Println("open: MonkeyType...")
	if err := monkeybot.OpenMonkeyType(); err != nil {
		fmt.Printf("oopsies!: %v\n", err)
		return
	}

	for {
		words, err := monkeybot.RunTest()
		if err != nil {
			fmt.Printf("oopsies! during typing: %v\n", err)
			break
		}

		fmt.Printf("\nfinished! (ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ typed %d words\n\n", words)
		fmt.Println("'B' to go again! or 'ENTER' to exit")

		var choice string
		fmt.Scanln(&choice)

		if choice != "b" && choice != "B" {
			break
		}

		fmt.Println("\nrestarting test...")
		if err := monkeybot.Restart(); err != nil {
			fmt.Printf("oopsies! restarting: %v\n", err)
			break
		}
	}

	fmt.Println("\nturning off...")
}
