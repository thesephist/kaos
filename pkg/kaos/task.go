package kaos

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"../ansi"
	"../wordid"
)

// Task represents a completable unit of work tracked in kaos
type Task struct {
	// string reference code used as unique ID
	Ref string

	// timestamps
	Created  time.Time
	Started  time.Time
	Finished time.Time
	Due      time.Time

	// project slug, all-lowercase, separated by '/'
	Project string

	// how much effort does this task require?
	Size int
	// content of the task
	Description string
	// each comment is a single-line string related to the task
	Comments []string

	// used in live structs to remove a task from the list
	// 	without modifying the slice of tasks directly
	deleted bool
}

// NewRef returns a random lowercase-alphabetical string ID
//	The choice to restrict it to lowercase alphabets is for
//	the purpose of ease of use, when referencing tasks in a CLI.
func NewRef() string {
	return wordid.Generate()
}

// String returns a serialized string representation of a Task,
//	which can be serialized back with Task.FromString
func (t Task) String() string {
	taskStr := fmt.Sprintf(
		"#%s [%s|%s|%s|%s]\n%s (%d): %s",
		t.Ref,
		formatTime(t.Created),
		formatTime(t.Started),
		formatTime(t.Finished),
		formatTime(t.Due),
		t.Project,
		t.Size,
		t.Description,
	)
	for _, commentStr := range t.Comments {
		taskStr += "\n\t" + commentStr
	}
	return taskStr
}

// Parse parses a given string and deserializes it into a Task
func ParseTask(taskStr string) (Task, error) {
	var err error

	lines := strings.Split(taskStr, "\n")

	firstLine := lines[0]
	secondLine := lines[1]
	commentParts := lines[2:]
	for idx, comm := range commentParts {
		commentParts[idx] = strings.TrimSpace(comm)
	}

	firstLineParts := strings.Split(firstLine, " ")
	refPart := firstLineParts[0]
	datePart := firstLineParts[1]
	dateParts := strings.Split(datePart[1:len(datePart)-1], "|")

	var dates []time.Time
	for _, dateStr := range dateParts {
		var result time.Time
		result, err = parseTime(dateStr)
		dates = append(dates, result)
	}

	secondLineParts := strings.Split(secondLine, " ")
	projectPart := secondLineParts[0]
	sizePartStr := secondLineParts[1][1 : len(secondLineParts[1])-2]
	sizePart, err := strconv.Atoi(sizePartStr)
	descriptionPart := strings.Join(secondLineParts[2:], " ")

	if err != nil {
		fmt.Println("There was an error reading tasks from disk", err)
		os.Exit(1)
	}

	t := Task{
		Ref:         refPart,
		Created:     dates[0],
		Started:     dates[1],
		Finished:    dates[2],
		Due:         dates[3],
		Project:     projectPart,
		Size:        sizePart,
		Description: descriptionPart,
		Comments:    commentParts,
	}
	return t, err
}

// Print returns a serialized string representation of a task, colorized and formatted
// 	for an ANSI interactive terminal
func (t Task) Print() string {
	coloredDueTime := "@" + formatTime(t.Due)
	now := time.Now()
	switch {
	case t.Due.IsZero():
		coloredDueTime = ansi.Grey(coloredDueTime)
	case t.Due.Before(now):
		coloredDueTime = ansi.Red(coloredDueTime)
	case time.Until(t.Due).Hours() < 24:
		coloredDueTime = ansi.Yellow(coloredDueTime)
	default:
		coloredDueTime = ansi.Green(coloredDueTime)
	}

	taskStr := fmt.Sprintf(
		ansi.Bold("#%s ")+ansi.Grey("[%s|%s|%s] ")+"%s\n"+ansi.Blue("%s")+ansi.Magenta("\t(%d)")+": %s",
		t.Ref,
		formatTime(t.Created),
		formatTime(t.Started),
		formatTime(t.Finished),
		coloredDueTime,
		t.Project,
		t.Size,
		t.Description,
	)
	for _, commentStr := range t.Comments {
		taskStr += "\n\t" + commentStr
	}
	return taskStr
}

// Reports whether a task's due date is in the past
func (t *Task) IsOverdue() bool {
	return !t.Due.IsZero() &&
		t.Finished.IsZero() &&
		t.Due.Before(time.Now())
}

// Start marks the task as started now
func (t *Task) Start() {
	t.Started = time.Now()
	t.Finished = time.Time{}
}

// Finish marks the task as finished now
func (t *Task) Finish() {
	now := time.Now()
	if t.Started.IsZero() {
		t.Started = now
	}
	t.Finished = now
}

// Unstart undoes Task.Start()
func (t *Task) Unstart() {
	t.Started = time.Time{}
	t.Finished = time.Time{}
}

// Unfinish undoes Task.Finish()
func (t *Task) Unfinish() {
	t.Finished = time.Time{}
}

// Delete deletes a task from a task list, so it is not
//	serialized back out in String() or Print()
func (t *Task) Delete() {
	t.deleted = true
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	} else {
		return t.Format("2006/01/02T15:04:05")
	}
}

func parseTime(timeStr string) (time.Time, error) {
	result, err := time.Parse("2006/01/02T15:04:05", timeStr)
	return result, err
}

// Matches reports whether a given task's Ref includes the given string
//	Think of it like the reverse of git rev-parse
func (t *Task) Matches(ref string) bool {
	return strings.Contains(t.Ref, ref)
}
