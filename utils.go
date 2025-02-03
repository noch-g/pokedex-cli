package main

func ToBold(s string) string {
	return "\033[1m" + s + "\033[0m"
}
