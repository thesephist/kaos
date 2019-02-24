package kaos

import (
	"bytes"
	"io"
	"strings"
)

func Write(writer io.Writer, tasks TaskList) (n int, err error) {
	buf := new(bytes.Buffer)
	buf.WriteString(tasks.String())
	n, err = writer.Write(buf.Bytes())
	return n, err
}

func Parse(reader io.Reader) (tasks TaskList, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
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
