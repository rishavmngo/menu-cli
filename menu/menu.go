package menu

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

type Node struct {
	name      string
	parent    *Node
	childrens []*Node
}

type Menu struct {
	Main *Node
}

var currentItem = 0

func (menu *Menu) Test() {

	fmt.Println(menu.Main.name)

	for _, item := range menu.Main.childrens {
		fmt.Println("\t" + item.name + "\t\t" + item.parent.name)

		for _, ineerItem := range item.childrens {
			fmt.Println("\t\t" + ineerItem.name + "\t\t" + ineerItem.parent.name)
			for _, superInnerItem := range ineerItem.childrens {
				fmt.Println("\t\t\t" + superInnerItem.name + "\t\t" + superInnerItem.parent.name)
			}
		}
	}

}
func (menu *Menu) Display() {
	head := menu.Main
	var buffer bytes.Buffer

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inpChan := make(chan byte)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()
	go func() {
		defer close(inpChan)

		inpCh := make([]byte, 1)
		var inpBuffer []byte

		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := os.Stdin.Read(inpCh)
				if err != nil {
					fmt.Println("Error reading input:", err)
					close(inpChan)
				}
				if n > 0 {
					if inpCh[0] == 27 { // ESC
						inpBuffer = append(inpBuffer, inpCh[0])
						continue
					}
					if len(inpBuffer) == 1 && inpCh[0] == 91 { // ESC [
						inpBuffer = append(inpBuffer, inpCh[0])
						continue
					}
					if len(inpBuffer) == 2 {
						// Check for escape sequences
						inpBuffer = append(inpBuffer, inpCh[0])
						switch inpBuffer[2] {
						case 'A':
							inpChan <- 'u' // Up arrow
						case 'B':
							inpChan <- 'd' // Down arrow
						case 'C':
							inpChan <- 'r' // right arrow
						case 'D':
							inpChan <- 'l' // left arrow
						}
						inpBuffer = inpBuffer[:0]
					} else {
						inpChan <- inpCh[0]
					}
				}
			}
		}

	}()
	ui := time.NewTicker(50 * time.Millisecond)
	defer ui.Stop()

mainLoop:
	for {
		select {
		case inp := <-inpChan:
			switch inp {
			case '\n', 13:
				expand(&head)
			case 3:
				cancel()
				break mainLoop
			case 127: //backspace
			case 23:
			case 32: //space
			case 'u':
				currentItem--
			case 'd':
				currentItem++
			case 'l':
				if canCollapse(head) {
					head = head.parent
					currentItem = 0
				}
			case 'r':
				expand(&head)
			default:
			}

		case <-ui.C:

			ClearScreenStandalone()
			headingOfList(head, &buffer)
			getListItems(head, &buffer)

			_, err := buffer.WriteTo(os.Stdout)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing buffer to stdout: %v\n", err)
			}
		}
	}

}

var Yellow = "\033[33m"

var reset = "\033[0m"
var bold = "\033[1m"

func headingOfList(head *Node, buffer *bytes.Buffer) {

	buffer.WriteString(bold)
	buffer.WriteString(fmt.Sprintf("%s(%d)\r\n", head.name, getActiveItemIndex(len(head.childrens))))
	buffer.WriteString(reset)
}

func getActiveItemIndex(length int) int {
	if currentItem >= 0 {
		return int(math.Abs(float64(currentItem % length)))
	} else {
		currentItem = length - 1
	}
	return currentItem
}

func getListItems(head *Node, buffer *bytes.Buffer) {

	for index, item := range head.childrens {
		if index == getActiveItemIndex(len(head.childrens)) {
			buffer.WriteString(Yellow)
			buffer.WriteString(fmt.Sprintf("> "))

		} else {
			buffer.WriteString("  ")
		}
		buffer.WriteString(fmt.Sprintf("%s\r\n", item.name))
		buffer.WriteString(reset)
	}

}

func NewMenu(name string) *Menu {

	// menu := &Menu{Main: &Node{Name: name, FirstChild: nil, NextSibiling: nil}}
	head := &Node{name: name}

	return &Menu{Main: head}

}

func (node *Node) Add(name string) *Node {

	newNode := &Node{name: name}
	newNode.parent = node
	node.childrens = append(node.childrens, newNode)
	return newNode

}
func expand(head **Node) {

	if canExpand(*head) {
		*head = (*head).childrens[getActiveItemIndex(len((*head).childrens))]
		currentItem = 0
	}
}

func canExpand(head *Node) bool {

	if len(head.childrens[getActiveItemIndex(len(head.childrens))].childrens) > 0 {
		return true
	}
	return false
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
