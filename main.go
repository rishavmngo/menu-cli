package main

import m "github.com/rishavmngo/menu-cli/menu"

func main() {

	menu := m.NewMenu("Main Menu")
	menu.Main.Add("Play")
	settings := menu.Main.Add("Settings")
	menu.Main.Add("Exit")

	Mode := settings.Add("Mode")
	Mode.Add("Easy")
	Mode.Add("Advance")
	Mode.Add("Paragraph")

	Duration := settings.Add("Duration")

	Duration.Add("10")
	Duration.Add("30")
	Duration.Add("60")
	Duration.Add("120")

	menu.Display()
}
