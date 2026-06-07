package createiteration

import (
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// ProjectV2IterationFieldConfigurationInput mirrors the GraphQL input object of the
// same name. The Go type name is used by shurcooL-graphql to derive the GraphQL
// variable type, so it must match the schema exactly.
type ProjectV2IterationFieldConfigurationInput struct {
	StartDate  string           `json:"startDate"`
	Duration   int              `json:"duration"`
	Iterations []IterationInput `json:"iterations"`
}

type CreateIterationsMutation struct {
	UpdateProjectV2Field struct {
		ClientMutationId string
	} `graphql:"updateProjectV2Field(input: {fieldId: $fieldId, iterationConfiguration: $iterationConfiguration})"`
}

func CreateIterations(
	gql *api.GraphQLClient,
	fieldID string,
	config IterationFieldConfiguration,
	iterations []IterationInput,
	dryRun bool,
) error {
	if dryRun {
		return nil
	}

	// The read-side schema does not expose a configuration-level start date, so derive
	// the required top-level startDate from the earliest iteration we are sending. Using
	// an existing iteration's start date keeps the field's cadence baseline intact.
	startDate := iterations[0].StartDate
	for _, it := range iterations[1:] {
		if it.StartDate < startDate {
			startDate = it.StartDate
		}
	}
	duration := config.Duration
	if duration == 0 {
		duration = iterations[len(iterations)-1].Duration
	}

	var mutation CreateIterationsMutation
	vars := map[string]any{
		"fieldId": graphql.ID(fieldID),
		"iterationConfiguration": ProjectV2IterationFieldConfigurationInput{
			StartDate:  startDate,
			Duration:   duration,
			Iterations: iterations,
		},
	}
	return gql.Mutate("CreateIterations", &mutation, vars)
}
