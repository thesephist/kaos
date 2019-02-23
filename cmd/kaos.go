package main

import (
	"bufio"
	"fmt"
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
	fmt.Println("Create")

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

	fmt.Println()
	fmt.Println("Created:")
	fmt.Println(t)
}

func runStart(tasks *kaos.TaskList, ref string) {
	fmt.Println("Start", ref)
}

func runFinish(tasks *kaos.TaskList, ref string) {
	fmt.Println("Finish", ref)
}

// Destructive actions

func runRemove(tasks *kaos.TaskList, ref string) {
	fmt.Println("Remove", ref)
}

func runUnstart(task *kaos.TaskList, ref string) {
	fmt.Println("Unstart", ref)
}

func runUnfinish(task *kaos.TaskList, ref string) {
	fmt.Println("Unfinish", ref)
}

// Update actions

func runSetDue(task *kaos.TaskList, ref string) {
	fmt.Println("Set due", ref)
}

func runSetProject(task *kaos.TaskList, ref string) {
	fmt.Println("Set project", ref)
}

func runSetSize(task *kaos.TaskList, ref string) {
	fmt.Println("Set size", ref)
}

func runSetDescription(task *kaos.TaskList, ref string) {
	fmt.Println("Set description", ref)
}

func runAddComment(task *kaos.TaskList, ref string) {
	fmt.Println("Add comment", ref)
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
		runFinish(&tasks, parameters[0])

	case "remove":
		runRemove(&tasks, parameters[0])
	case "Unstart":
		runUnstart(&tasks, parameters[0])
	case "Unfinish":
		runUnfinish(&tasks, parameters[0])

	case "set-due":
		runSetDue(&tasks, parameters[0])
	case "set-project":
		runSetProject(&tasks, parameters[0])
	case "set-size":
		runSetSize(&tasks, parameters[0])
	case "set-description":
		runSetDescription(&tasks, parameters[0])
	case "add-comment":
		runAddComment(&tasks, parameters[0])

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
