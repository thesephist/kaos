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

	var target *kaos.Task
	if len(parameters) == 1 {
		target, err = tasks.FindMatch(parameters[0])
	}

	switch action {
	case "list":
		fmt.Println(tasks)
	case "create":
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

		fmt.Println("Created:")
		fmt.Println(t)
	case "start":
		target.Start()
		fmt.Printf("Started #%s: %s\n", target.Ref, target.Description)
	case "finish":
		target.Finish()
		fmt.Printf("Finished#%s: %s\n", target.Ref, target.Description)

	case "remove":
		target.Delete()
		fmt.Printf("Removed #%s: %s\n", target.Ref, target.Description)
	case "Unstart":
		target.Unstart()
		fmt.Printf("Unstarted #%s: %s\n", target.Ref, target.Description)
	case "Unfinish":
		target.Unfinish()
		fmt.Printf("Unfinished#%s: %s\n", target.Ref, target.Description)

	case "set-due":
	case "set-project":
		project := Prompt("Project?")
		target.Project = project
		fmt.Println("Updated.")
		fmt.Println(target)
	case "set-size":
		sizeStr := Prompt("Size?")
		size, _ := strconv.Atoi(sizeStr)
		target.Size = size
		fmt.Println("Updated.")
		fmt.Println(target)
	case "set-description":
		description := Prompt("Description?")
		target.Description = description
		fmt.Println("Updated.")
		fmt.Println(target)
	case "add-comment":

	default:
		fmt.Printf("Unknown kaos action '%s'\n", action)
		os.Exit(1)
	}

	writtenBytes, err := kaos.Write(taskFile, tasks)
	err = taskFile.Truncate(int64(writtenBytes))
	if err != nil {
		fmt.Println("Error writing kaosfile:", err)
		os.Exit(1)
	}
}
