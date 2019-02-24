package kaos

import (
	"errors"
	"sort"
	"strings"
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

func (tasks *TaskList) FindMatch(ref string) (match *Task, err error) {
	matchIdx, count := -1, 0
	for idx, t := range *tasks {
		if t.Matches(ref) {
			matchIdx = idx
			count++
		}
	}
	match = &(*tasks)[matchIdx]

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
	*tasks = append(*tasks, t)
}
