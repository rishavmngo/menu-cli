package menu

import (
	"bytes"
	"fmt"
	"math"

	"github.com/rishavmngo/menu-go/models"
)

func expand(head **Node, currentItem *models.CurrentItem) {

	if canExpand(*head, currentItem) {
		*head = (*head).childrens[getActiveItemIndex(len((*head).childrens), currentItem)]
		currentItem.Reset()
	} else {
		if (*head).childrens[getActiveItemIndex(len((*head).childrens), currentItem)].action != nil {

			(*head).childrens[getActiveItemIndex(len((*head).childrens), currentItem)].action()
		}
	}
}
func canExpand(head *Node, currentItem *models.CurrentItem) bool {

	if len(head.childrens[getActiveItemIndex(len(head.childrens), currentItem)].childrens) > 0 {
		return true
	}
	return false
}

func collapse(head **Node, currentItem *models.CurrentItem) {

	if canCollapse(*head) {
		*head = (*head).parent
		currentItem.Reset()
	}
}
func canCollapse(head *Node) bool {

	if head.parent != nil {
		return true
	}
	return false

}

func ClearScreenStandalone() {

	fmt.Print("\033[2J") // Clear the screen
	fmt.Print("\033[H")  // Move the cursor to the top-left corner
}

func getListItems(head *Node, buffer *bytes.Buffer, currentItem *models.CurrentItem) {

	for index, item := range head.childrens {
		if index == getActiveItemIndex(len(head.childrens), currentItem) {
			buffer.WriteString(bold)
			buffer.WriteString(Yellow)
			buffer.WriteString(fmt.Sprintf("> "))

		} else {
			buffer.WriteString("  ")
		}
		buffer.WriteString(fmt.Sprintf("%s\r\n", item.name))
		buffer.WriteString(reset)
	}

}

var Yellow = "\033[33m"

var reset = "\033[0m"
var bold = "\033[1m"

func headingOfList(head *Node, buffer *bytes.Buffer, currentItem *models.CurrentItem) {

	buffer.WriteString(bold)
	buffer.WriteString(fmt.Sprintf("%s\r\n", head.name))
	buffer.WriteString(reset)
}
func getActiveItemIndex(length int, currentItem *models.CurrentItem) int {
	if currentItem.Get() >= 0 {
		return int(math.Abs(float64(currentItem.Get() % length)))
	} else {
		// currentItem = length - 1
		currentItem.Set(length - 1)
	}
	return currentItem.Get()
}
