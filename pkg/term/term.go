package term

import "fmt"

func ClearCurrentLine() {
	fmt.Print("\033[K")
}

func MoveCursorUp(lines int) {
	fmt.Printf("\033[%dA", lines)
}

func MoveCursorDown(lines int) {
	fmt.Printf("\033[%dB", lines)
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}
