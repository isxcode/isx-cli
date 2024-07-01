/*
Copyright © 2024 jamie HERE <EMAIL ADDRESS>
*/
package common

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func textWithColor(color, text string) string {
	return color + text + Reset
}

func RedText(text string) string {
	return textWithColor(Red, text)
}

func GreenText(text string) string {
	return textWithColor(Green, text)
}

func YellowText(text string) string {
	return textWithColor(Yellow, text)
}

func BlueText(text string) string {
	return textWithColor(Blue, text)
}

func PurpleText(text string) string {
	return textWithColor(Purple, text)
}

func CyanText(text string) string {
	return textWithColor(Cyan, text)
}

func WhiteText(text string) string {
	return textWithColor(White, text)
}
