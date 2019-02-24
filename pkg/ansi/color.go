package ansi

const END string = "\u001b[0m"

// styles

func Bold(s string) string {
	return "\u001b[1m" + s + END
}

func Grey(s string) string {
	return "\u001b[30;1m" + s + END
}

func Underline(s string) string {
	return "\u001b[4m" + s + END
}

// colors

func Blue(s string) string {
	return "\u001b[34m" + s + END
}

func Yellow(s string) string {
	return "\u001b[33m" + s + END
}

func Red(s string) string {
	return "\u001b[31m" + s + END
}

func Green(s string) string {
	return "\u001b[32m" + s + END
}

func Magenta(s string) string {
	return "\u001b[35m" + s + END
}
