package kaos

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strings"
	"time"
)

type TaskList []Task

func (tasks TaskList) String() string {
	var s []string
	for _, t := range tasks {
		if !t.deleted {
			s = append(s, t.String())
		}
	}
	return strings.Join(s, "\n")
}

func (tasks TaskList) Print() string {
	var s []string
	for _, t := range tasks {
		// don't show finished tasks by default
		if !t.deleted && t.Finished.IsZero() {
			s = append(s, t.Print())
		}
	}
	return strings.Join(s, "\n")
}

func (tasks TaskList) PrintAll() string {
	var s []string
	for _, t := range tasks {
		if !t.deleted {
			s = append(s, t.Print())
		}
	}
	return strings.Join(s, "\n")
}

func (tasks TaskList) Sorted() TaskList {
	sorted := TaskList(tasks[:])
	sort.Slice(sorted, func(i, j int) bool {
		iT, jT := sorted[i].Due, sorted[j].Due
		if iT.IsZero() {
			return false
		} else {
			if jT.IsZero() {
				return true
			} else {
				return iT.Before(jT)
			}
		}
	})
	return TaskList(sorted)
}

func iContains(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}

func (tasks *TaskList) Search(sub string) (matches TaskList) {
	for _, t := range *tasks {
		if iContains(t.Description, sub) || iContains(t.Project, sub) {
			matches = append(matches, t)
		}
	}
	return matches
}

func (tasks *TaskList) RescheduleOverdue() {
	now := time.Now()
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)
	for i, t := range *tasks {
		if t.IsOverdue() {
			(*tasks)[i].Due = today
		}
	}
}

func (tasks *TaskList) FindMatch(ref string) (match *Task, err error) {
	matchIdx, count := -1, 0
	for idx, t := range *tasks {
		if t.Matches(ref) {
			matchIdx = idx
			count++
		}
		if count > 1 {
			break
		}
	}

	switch count {
	case 0:
		err = errors.New("No match found")
		return
	case 1:
		match = &(*tasks)[matchIdx]
		return
	default:
		err = errors.New("More than one matches found")
		return
	}
}

func (tasks *TaskList) AddTask(t Task) {
	*tasks = append(*tasks, t)
}

func Write(writer io.Writer, tasks TaskList) (n int, err error) {
	buf := new(bytes.Buffer)
	buf.WriteString(tasks.String())
	n, err = writer.Write(buf.Bytes())
	return n, err
}

func Parse(reader io.Reader) (tasks TaskList, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)

	// We split by newline + #, so we pad the front
	str := "\n" + strings.TrimSpace(buf.String())

	if err != nil {
		return
	}

	taskParts := strings.Split(str, "\n#")[1:]

	for _, taskPart := range taskParts {
		var t Task
		t, err = ParseTask(taskPart)
		tasks.AddTask(t)
	}

	return
}
