package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"../pkg/kaos"
)

func Prompt(str string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(str)

	input, _ := reader.ReadString('\n')
	return input
}

func main() {
	args := os.Args[1:]

	action := args[0]
	parameters := args[1:]

	switch action {
	case "create":
		fmt.Println("Create", parameters)
	case "start":
		fmt.Println("Start", parameters)
	case "finish":
		fmt.Println("Finish", parameters)

	case "remove":
		fmt.Println("Remove", parameters)
	case "Unstart":
		fmt.Println("Unstart", parameters)
	case "Unfinish":
		fmt.Println("Unfinish", parameters)

	case "set-due":
		fmt.Println("Set due", parameters)
	case "set-project":
		fmt.Println("Set project", parameters)
	case "set-size":
		fmt.Println("Set size", parameters)
	case "set-description":
		fmt.Println("Set description", parameters)
	case "add-comment":
		fmt.Println("Add comment", parameters)

	default:
		fmt.Printf("Unknown action '%s'\n", action)
		os.Exit(1)
	}

	t := kaos.Task{
		Ref:         kaos.NewRef(),
		Created:     time.Now(),
		Project:     "inbox",
		Size:        5,
		Description: "Test description",
		Comments:    []string{"First comment", "Second"},
	}
	fmt.Println(t)
	t.Start()
	fmt.Println(t)

	input := Prompt("What do you want?")
	fmt.Println("Input was", input)
}
