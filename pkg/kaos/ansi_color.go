package kaos

const END string = "\u001b[0m"

func Bold(s string) string {
	return "\u001b[1m" + s + END
}

func Grey(s string) string {
	return "\u001b[30;1m" + s + END
}

func Blue(s string) string {
	return "\u001b[34m" + s + END
}

func Yellow(s string) string {
	return "\u001b[33;1m" + s + END
}
