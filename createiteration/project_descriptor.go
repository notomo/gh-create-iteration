package createiteration

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

type ProjectDescriptor struct {
	OwnerIsOrganization bool
	Owner               string
	Number              int
}

var projectUrlRegex = regexp.MustCompile(`^/(users|orgs)/([^/]+)/projects/(\d+)`)

func GetProjectDescriptor(
	projectUrl string,
) (*ProjectDescriptor, error) {
	u, err := url.Parse(projectUrl)
	if err != nil {
		return nil, err
	}

	matches := projectUrlRegex.FindStringSubmatch(u.Path)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid project url path: %s", u.Path)
	}

	number, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, err
	}

	return &ProjectDescriptor{
		OwnerIsOrganization: matches[1] == "orgs",
		Owner:               matches[2],
		Number:              number,
	}, nil
}
