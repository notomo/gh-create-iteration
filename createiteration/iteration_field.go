package createiteration

import (
	"fmt"
	"slices"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// Iteration is an existing iteration read from the project iteration field.
// ID is read but not used in the update input (the input object has no id field);
// existing iterations are preserved by re-sending Title/StartDate/Duration verbatim.
type Iteration struct {
	ID        string
	Title     string
	StartDate string
	Duration  int
}

// IterationFieldConfiguration is the current configuration of an iteration field.
// Existing holds the completed and active iterations combined in start-date order.
// StartDate is only populated when there are no existing iterations to derive a start
// date from; the read-side schema does not expose a configuration-level start date, so
// in practice it stays empty and -start-date is required for an empty field.
type IterationFieldConfiguration struct {
	Duration  int
	StartDate string
	Existing  []Iteration
}

type ProjectV2IterationFieldConfigurationQuery struct {
	Duration            int
	Iterations          []Iteration
	CompletedIterations []Iteration
}

type ProjectV2IterationFieldQuery struct {
	ID            string
	Configuration ProjectV2IterationFieldConfigurationQuery
}

type ProjectV2Query struct {
	Field struct {
		ProjectV2IterationFieldQuery `graphql:"... on ProjectV2IterationField"`
	} `graphql:"field(name: $fieldName)"`
}

type GetUserIterationFieldQuery struct {
	User struct {
		ProjectV2Query `graphql:"projectV2(number: $projectNumber)"`
	} `graphql:"user(login: $owner)"`
}

type GetOrganizationIterationFieldQuery struct {
	Organization struct {
		ProjectV2Query `graphql:"projectV2(number: $projectNumber)"`
	} `graphql:"organization(login: $owner)"`
}

func GetIterationField(
	gql *api.GraphQLClient,
	descriptor ProjectDescriptor,
	iterationFieldName string,
) (string, *IterationFieldConfiguration, error) {
	vars := map[string]any{
		"owner":         graphql.String(descriptor.Owner),
		"projectNumber": graphql.Int(descriptor.Number),
		"fieldName":     graphql.String(iterationFieldName),
	}

	var field ProjectV2IterationFieldQuery
	if descriptor.OwnerIsOrganization {
		var query GetOrganizationIterationFieldQuery
		if err := gql.Query("GetIterationField", &query, vars); err != nil {
			return "", nil, err
		}
		field = query.Organization.Field.ProjectV2IterationFieldQuery
	} else {
		var query GetUserIterationFieldQuery
		if err := gql.Query("GetIterationField", &query, vars); err != nil {
			return "", nil, err
		}
		field = query.User.Field.ProjectV2IterationFieldQuery
	}

	if field.ID == "" {
		return "", nil, fmt.Errorf("no iteration field found: %s", iterationFieldName)
	}

	existing := []Iteration{}
	existing = append(existing, field.Configuration.CompletedIterations...)
	existing = append(existing, field.Configuration.Iterations...)
	slices.SortStableFunc(existing, func(a, b Iteration) int {
		if a.StartDate < b.StartDate {
			return -1
		}
		if a.StartDate > b.StartDate {
			return 1
		}
		return 0
	})

	return field.ID, &IterationFieldConfiguration{
		Duration: field.Configuration.Duration,
		Existing: existing,
	}, nil
}
