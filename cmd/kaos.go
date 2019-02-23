package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"../pkg/kaos"
)

func Prompt(str string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(">", str)

	input, _ := reader.ReadString('\n')
	return input[:len(input)-1]
}

// Additive actions

func runList(tasks *kaos.TaskList) {
	fmt.Println(tasks)
}

func runCreate(tasks *kaos.TaskList) {
	project := Prompt("Project?")
	sizeStr := Prompt("Size?")
	size, err := strconv.Atoi(sizeStr)
	description := Prompt("Description?")

	if err != nil {
		fmt.Println("Your input was invalid!", err)
		os.Exit(1)
	}

	t := kaos.Task{
		Ref:         kaos.NewRef(),
		Created:     time.Now(),
		Project:     project,
		Size:        size,
		Description: description,
	}
	tasks.AddTask(t)

	fmt.Println("Created")
	fmt.Println(t)
}

func runStart(tasks *kaos.TaskList, ref string) {
	t, err := tasks.FindMatch(ref)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	t.Start()

	fmt.Printf("Started #%s: %s\n", t.Ref, t.Description)
}

func main() {
	args := os.Args[1:]

	action := args[0]
	parameters := args[1:]

	taskFile, err := os.OpenFile("./kaosfile", os.O_RDWR, 0755)
	defer taskFile.Close()
	if err != nil {
		fmt.Println("Error encountered while reading kaos tasks file:", err)
		os.Exit(1)
	}
	tasks, err := kaos.Parse(taskFile)
	if err != nil {
		fmt.Println("Error encountered while reading kaos tasks file:", err)
		os.Exit(1)
	}
	taskFile.Seek(0, 0)

	switch action {
	case "list":
		runList(&tasks)
	case "create":
		runCreate(&tasks)
	case "start":
		runStart(&tasks, parameters[0])
	case "finish":

	case "remove":
	case "Unstart":
	case "Unfinish":

	case "set-due":
	case "set-project":
	case "set-size":
	case "set-description":
	case "add-comment":

	default:
		fmt.Printf("Unknown kaos action '%s'\n", action)
		os.Exit(1)
	}

	err = kaos.Write(taskFile, tasks)
	if err != nil {
		fmt.Println("Error writing kaosfile:", err)
		os.Exit(1)
	}
}
