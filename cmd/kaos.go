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

	fileChanged := false

	if len(args) == 0 {
		args = append(args, "list")
	}

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
		fmt.Println(tasks.Sorted().Print())
	case "all":
		fmt.Println(tasks.Sorted().PrintAll())
	case "find":
		fmt.Println(tasks.Search(parameters[0]).Sorted().Print())
	case "fold":
		fileChanged = true
		tasks.RescheduleOverdue()
	case "create":
		fileChanged = true
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
		fmt.Println(t.Print())
	case "start":
		fileChanged = true
		target.Start()
		fmt.Println("Started", target.Print())
	case "finish":
		fileChanged = true
		target.Finish()
		fmt.Println("Finished", target.Print())

	case "remove":
		fileChanged = true
		target.Delete()
		fmt.Println("Removed", target.Print())
	case "unstart":
		fileChanged = true
		target.Unstart()
		fmt.Println("Unstarted", target.Print())
	case "unfinish":
		fileChanged = true
		target.Unfinish()
		fmt.Println("Unfinished", target.Print())

	case "due":
		fileChanged = true
		dateStr := Prompt("Due Date?")
		date, err := time.Parse("2006/01/02T15:04:05", dateStr)
		if err != nil {
			date, err = time.Parse("2006/01/02", dateStr)
		}
		if err != nil {
			fmt.Println("Your date was invalid")
		} else {
			target.Due = date
			fmt.Println(target.Print())
		}
	case "project":
		fileChanged = true
		project := Prompt("Project?")
		target.Project = project
		fmt.Println("Updated.")
		fmt.Println(target.Print())
	case "size":
		fileChanged = true
		sizeStr := Prompt("Size?")
		size, _ := strconv.Atoi(sizeStr)
		target.Size = size
		fmt.Println("Updated.")
		fmt.Println(target.Print())
	case "describe":
		fileChanged = true
		description := Prompt("Description?")
		target.Description = description
		fmt.Println("Updated.")
		fmt.Println(target.Print())
	case "comment":
		fileChanged = true
		newComment := Prompt("New Comment?")
		target.Comments = append(target.Comments, newComment)
		fmt.Println("Updated.")
		fmt.Println(target.Print())

	default:
		fmt.Printf("Unknown kaos action '%s'\n", action)
		os.Exit(1)
	}

	if fileChanged {
		writtenBytes, err := kaos.Write(taskFile, tasks)
		err = taskFile.Truncate(int64(writtenBytes))
		if err != nil {
			fmt.Println("Error writing kaosfile:", err)
			os.Exit(1)
		}
	}
}
