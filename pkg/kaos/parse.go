package kaos

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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
		var err error

		lines := strings.Split(taskPart, "\n")

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
			result, err = time.Parse("2006/01/02T15:04:05", dateStr)
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

		tasks.AddTask(Task{
			Ref:         refPart,
			Created:     dates[0],
			Started:     dates[1],
			Finished:    dates[2],
			Due:         dates[3],
			Project:     projectPart,
			Size:        sizePart,
			Description: descriptionPart,
			Comments:    commentParts,
		})
	}

	return
}
