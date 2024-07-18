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
	// menu.Test()

	// mainMenu.addItem("Play")
	// settings := mainMenu.addItem("Settings")
	// mainMenu.addItem("Exit")
	//
	//  mode := settings.addItem("mode")
	//  mode.addItem("easy")
	//  mode.addItem("paragraph")
	//  mode.addItem("advance")
	//
	//  duration:= settings.addItem("duration")
	//
	//  duration.addItem("10")
	//  duration.addItem("20")
	//  duration.addItem("30")
	//  duration.addItem("50")

}
