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

func runList(tasks kaos.TaskList) kaos.TaskList {
	fmt.Println(tasks)
	return tasks
}

func runCreate(tasks kaos.TaskList) kaos.TaskList {
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
	return tasks
}

func runStart(tasks kaos.TaskList, ref string) kaos.TaskList {
	idx, t, err := tasks.FindMatch(ref)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	t.Start()
	tasks.SetTask(idx, t)

	fmt.Printf("Started #%s: %s\n", t.Ref, t.Description)

	return tasks
}

func runFinish(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Finish", ref)
	return tasks
}

// Destructive actions

func runRemove(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Remove", ref)
	return tasks
}

func runUnstart(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Unstart", ref)
	return tasks
}

func runUnfinish(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Unfinish", ref)
	return tasks
}

// Update actions

func runSetDue(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Set due", ref)
	return tasks
}

func runSetProject(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Set project", ref)
	return tasks
}

func runSetSize(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Set size", ref)
	return tasks
}

func runSetDescription(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Set description", ref)
	return tasks
}

func runAddComment(tasks kaos.TaskList, ref string) kaos.TaskList {
	fmt.Println("Add comment", ref)
	return tasks
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
		tasks = runList(tasks)
	case "create":
		tasks = runCreate(tasks)
	case "start":
		tasks = runStart(tasks, parameters[0])
	case "finish":
		tasks = runFinish(tasks, parameters[0])

	case "remove":
		tasks = runRemove(tasks, parameters[0])
	case "Unstart":
		tasks = runUnstart(tasks, parameters[0])
	case "Unfinish":
		tasks = runUnfinish(tasks, parameters[0])

	case "set-due":
		tasks = runSetDue(tasks, parameters[0])
	case "set-project":
		tasks = runSetProject(tasks, parameters[0])
	case "set-size":
		tasks = runSetSize(tasks, parameters[0])
	case "set-description":
		tasks = runSetDescription(tasks, parameters[0])
	case "add-comment":
		tasks = runAddComment(tasks, parameters[0])

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
