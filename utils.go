package main

import "fmt"

func GetPromptMessage() string {
	return "\r\033[K" + ToBold("Pokedex > ")
}

func ToBold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func StartFromClearLine() {
	fmt.Print("\r\033[K")
}
