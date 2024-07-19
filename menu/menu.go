package menu

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rishavmngo/menu-cli/models"
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

func (menu *Menu) Display() {
	head := menu.Main
	var buffer bytes.Buffer

	var currentItem = models.Get()

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
							inpChan <- '>' // right arrow
						case 'D':
							inpChan <- '<' // left arrow
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
				expand(&head, currentItem)
			case 3:
				cancel()
				break mainLoop
			case 127: //backspace
			case 23:
			case 32: //space
			case 'j':
				currentItem.Increment()
			case 'k':
				currentItem.Decrement()
			case 'u':
				currentItem.Decrement()
			case 'd':
				currentItem.Increment()
			case '<':
				collapse(&head, currentItem)
			case 'l':
				expand(&head, currentItem)
			case 'h':
				collapse(&head, currentItem)
			case '>':
				expand(&head, currentItem)
			default:
			}

		case <-ui.C:

			ClearScreenStandalone()
			headingOfList(head, &buffer, currentItem)
			getListItems(head, &buffer, currentItem)

			_, err := buffer.WriteTo(os.Stdout)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing buffer to stdout: %v\n", err)
			}
		}
	}

}

func NewMenu(name string) *Menu {

	head := &Node{name: name}

	return &Menu{Main: head}

}

func (node *Node) Add(name string) *Node {

	newNode := &Node{name: name}
	newNode.parent = node
	node.childrens = append(node.childrens, newNode)
	return newNode

}
