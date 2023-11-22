package meetup

import (
	"fmt"
	"os"
	"strings"
)

func MeetingFromPath(gs GroupStrategy, p string) (Meeting, error) {
	components := strings.Split(strings.Trim(p, string(os.PathSeparator)), string(os.PathSeparator))

	if len(components) < 3 {
		return Meeting{}, fmt.Errorf("path does not have enough components '%s'", p)
	}

	meeting := Meeting{
		Name: components[len(components)-1],
	}

	switch gs {
	case GroupByDomain:
		meeting.Date = components[len(components)-2]
		meeting.Domain = strings.Join(components[:len(components)-2], ".")
	case GroupByDate:
		meeting.Date = components[0]
		meeting.Domain = strings.Join(components[1:len(components)-1], ".")
	}

	return meeting, nil
}
