package createiteration

import (
	"fmt"
	"io"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
)

func Run(
	gql *api.GraphQLClient,
	projectUrl string,
	iterationFieldName string,
	count int,
	duration int,
	startDate string,
	titlePrefix string,
	dryRun bool,
	writer io.Writer,
) error {
	projectDescriptor, err := GetProjectDescriptor(projectUrl)
	if err != nil {
		return err
	}

	fieldID, config, err := GetIterationField(gql, *projectDescriptor, iterationFieldName)
	if err != nil {
		return err
	}

	iterations, err := BuildIterations(*config, count, duration, startDate, titlePrefix)
	if err != nil {
		return err
	}

	if err := CreateIterations(gql, fieldID, *config, iterations, dryRun); err != nil {
		return err
	}

	created := iterations[len(config.Existing):]
	var b strings.Builder
	fmt.Fprintf(&b, "\nCreated %d iteration(s) (kept %d existing):\n", len(created), len(config.Existing))
	for _, it := range created {
		fmt.Fprintf(&b, "- %s (start date: %s, duration: %d)\n", it.Title, it.StartDate, it.Duration)
	}
	if _, err := writer.Write([]byte(b.String())); err != nil {
		return err
	}

	return nil
}
