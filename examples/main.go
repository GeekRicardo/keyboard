package main

import (
	"fmt"

	"github.com/GeekRicardo/keyboardhook"
)

func main() {
	keyboardhook.SetKeyDownCallback(func(event keyboardhook.KeyEvent) {
		fmt.Printf("Key down: %s (%d)\n", event.Name, event.VKCode)
	})

	keyboardhook.SetKeyUpCallback(func(event keyboardhook.KeyEvent) {
		fmt.Printf("Key up: %s (%d)\n", event.Name, event.VKCode)
	})

	fmt.Println("Press Ctrl-c to exit...")
	keyboardhook.Start()
}
