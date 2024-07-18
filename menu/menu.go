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
					inpBuffer = append(inpBuffer, inpCh[0])
					if len(inpBuffer) >= 3 {
						// Check for escape sequences
						if inpBuffer[0] == 27 && inpBuffer[1] == 91 { // ESC [
							switch inpBuffer[2] {
							case 'A':
								inpChan <- 'u' // Up arrow
								inpBuffer = inpBuffer[:0]
							case 'B':
								inpChan <- 'd' // Down arrow
								inpBuffer = inpBuffer[:0]
							case 'C':
								inpChan <- 'r' // right arrow
								inpBuffer = inpBuffer[:0]
							case 'D':
								inpChan <- 'l' // left arrow
								inpBuffer = inpBuffer[:0]
							default:
								inpBuffer = inpBuffer[1:]
							}
						} else {
							// Send non-escape sequence characters to the channel
							for _, b := range inpBuffer {
								inpChan <- b
							}
							inpBuffer = inpBuffer[:0]
						}
					}
				}
			}
		}

	}()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

mainLoop:
	for {
		select {
		case inp := <-inpChan:
			switch inp {
			case '\n':
			case '\r':
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
				if head.parent != nil {
					head = head.parent
					currentItem = 0
				}
			case 'r':
				if len(head.childrens) > 0 {
					head = head.childrens[currentItem]
					currentItem = 0
				}
			default:
			}

		case <-ticker.C:

			ClearScreenStandalone()
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

func getActiveItem(length int) int {
	if currentItem >= 0 {
		return int(math.Abs(float64(currentItem % length)))
	} else {
		currentItem = length - 1
	}
	return currentItem
}
func cyclicAccess(arrayLength int) int {
	currentIndex := currentItem
	// Handle negative or out-of-bounds index
	if currentIndex < 0 {
		currentIndex = (currentIndex + arrayLength) % arrayLength
	} else if currentIndex >= arrayLength {
		currentIndex = currentIndex % arrayLength
	}
	// Access the element using the adjusted index
	return currentIndex
}

func getListItems(head *Node, buffer *bytes.Buffer) {

	for index, item := range head.childrens {
		if index == getActiveItem(len(head.childrens)) {
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

func ClearScreenStandalone() {

	fmt.Print("\033[2J") // Clear the screen
	fmt.Print("\033[H")  // Move the cursor to the top-left corner
}
