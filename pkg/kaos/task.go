package kaos

import (
	"../wordid"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Task struct {
	Ref string

	Created  time.Time
	Started  time.Time
	Finished time.Time
	Due      time.Time

	// can include /s
	Project string

	Size        int
	Description string
	Comments    []string

	deleted bool
}

type TaskList struct {
	list []Task
}

func NewRef() string {
	return wordid.Generate()
}

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

func (t Task) Print() string {
	taskStr := fmt.Sprintf(
		Bold("#%s")+" "+Grey("[%s|%s|%s]")+Red(" @%s")+"\n"+Blue("%s")+Yellow("\t(%d)")+": %s",
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

func (t *Task) Start() {
	t.Started = time.Now()
}

func (t *Task) Finish() {
	t.Finished = time.Now()
}

func (t *Task) Unstart() {
	t.Started = time.Time{}
	t.Finished = time.Time{}
}

func (t *Task) Unfinish() {
	t.Finished = time.Time{}
}

func (t *Task) Delete() {
	t.deleted = true
}

func formatTime(t time.Time) string {
	zero := time.Time{}
	if t == zero {
		return "-"
	} else {
		return t.Format("2006/01/02T15:04:05")
	}
}

func (t *Task) Matches(ref string) bool {
	return strings.Contains(t.Ref, ref)
}

func (tasks TaskList) String() string {
	var s []string
	for _, t := range tasks.list {
		if !t.deleted {
			s = append(s, t.String())
		}
	}
	return strings.Join(s, "\n")
}

func (tasks TaskList) Print() string {
	var s []string
	for _, t := range tasks.list {
		if !t.deleted {
			s = append(s, t.Print())
		}
	}
	return strings.Join(s, "\n")
}

func (tasks TaskList) Sorted() TaskList {
	zeroTime := time.Time{}
	sorted := tasks.list[:]
	sort.Slice(sorted, func(i, j int) bool {
		iT := sorted[i].Due
		jT := sorted[j].Due
		if iT == zeroTime {
			return false
		} else {
			if jT == zeroTime {
				return true
			} else {
				return iT.Before(jT)
			}
		}
	})
	return TaskList{
		list: sorted,
	}
}

func (tasks *TaskList) FindMatch(ref string) (match *Task, err error) {
	matchIdx, count := -1, 0
	for idx, t := range tasks.list {
		if t.Matches(ref) {
			matchIdx = idx
			count++
		}
	}
	match = &tasks.list[matchIdx]

	switch count {
	case 0:
		err = errors.New("No match found")
		return
	case 1:
		return
	default:
		err = errors.New("More than one matches found")
		return
	}
}

func (tasks *TaskList) AddTask(t Task) {
	tasks.list = append(tasks.list, t)
}
