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
					inpChan <- inpCh[0]
				}
			}
		}

	}()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	bs := make([]byte, 3)
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
			case 27:
				bs[0] = 27
			case 91:
				bs[1] = 91
			case 65:
				if bs[0] == 27 && bs[1] == 91 {
					// fmt.Println("up")
					currentItem--

					bs[0] = 0
					bs[1] = 0

				}
			case 66:

				currentItem++
				// fmt.Println("down")
				if bs[0] == 27 && bs[1] == 91 {

					bs[0] = 0
					bs[1] = 0
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

func getListItems(head *Node, buffer *bytes.Buffer) {

	for index, item := range head.childrens {
		if index == int(math.Abs(float64(currentItem%len(head.childrens)))) {
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
