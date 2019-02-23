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

/*
	"#%s [%s|%s|%s|%s]
	%s (%d): %s",
	t.Ref,
	formatTime(t.Created),
	formatTime(t.Started),
	formatTime(t.Finished),
	formatTime(t.Due),
	t.Project,
	t.Size,
	t.Description,
*/

func Parse(reader io.Reader) (tasks TaskList, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	str := "\n" + strings.TrimSpace(buf.String())

	if err != nil {
		return
	}

	taskParts := strings.Split(str, "\n#")[1:]
	fmt.Println(taskParts)

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

		fmt.Println("ref:", refPart)
		fmt.Println("Dates:", dates)

		secondLineParts := strings.Split(secondLine, " ")
		projectPart := secondLineParts[0]
		sizePartStr := secondLineParts[1][1 : len(secondLineParts[1])-2]
		sizePart, err := strconv.Atoi(sizePartStr)
		descriptionPart := strings.Join(secondLineParts[2:], " ")

		fmt.Println("Project:", projectPart)
		fmt.Println("Size:", sizePart)
		fmt.Println("Description:", descriptionPart)
		fmt.Println("Comments:", commentParts)

		if err != nil {
			fmt.Println("There was an error reading tasks from disk", err)
			os.Exit(1)
		}
	}

	return
}
