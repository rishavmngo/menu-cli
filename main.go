package main

import (
	"fmt"

	m "github.com/rishavmngo/menu-go/menu"
)

func main() {

	menu := m.NewMenu("Main Menu")

	menu.Main.Add("Play", func() {
		fmt.Println("play")
	})
	settings := menu.Main.Add("Settings", nil)
	menu.Main.Add("Exit", func() {
		menu.Exit()
	})

	Mode := settings.Add("Mode", nil)
	Mode.Add("Easy", nil)
	Mode.Add("Advance", nil)
	Mode.Add("Paragraph", nil)

	Duration := settings.Add("Duration", nil)

	Duration.Add("10", nil)
	Duration.Add("30", nil)
	Duration.Add("60", nil)
	Duration.Add("120", nil)

	menu.Display()
	fmt.Println("hello")
}
