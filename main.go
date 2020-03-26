package main

import (
	"bufio"
	"fmt"
	"os"
)

func inputLoop(ch chan string) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		ch <- s.Text()
	}

}

func main() {
	userInput := make(chan string)
	go inputLoop(userInput)
	fmt.Println("start...")
	for {
		fmt.Println(<-userInput)
		// TODO: mainloop
	}
}
