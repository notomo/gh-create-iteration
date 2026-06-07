package createiteration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/notomo/gh-create-iteration/createiteration/gqltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type capturedMutation struct {
	Variables struct {
		FieldID                string `json:"fieldId"`
		IterationConfiguration struct {
			StartDate  string           `json:"startDate"`
			Duration   int              `json:"duration"`
			Iterations []IterationInput `json:"iterations"`
		} `json:"iterationConfiguration"`
	} `json:"variables"`
}

func TestRun(t *testing.T) {
	var captured capturedMutation

	gql, err := gqltest.New(
		t,

		gqltest.WithQueryOK("GetIterationField", `
{
  "data": {
    "user": {
      "projectV2": {
        "field": {
          "id": "PVTIF_22222222222222222222",
          "configuration": {
            "duration": 14,
            "completedIterations": [
              {
                "id": "00000a0b",
                "title": "Iteration 1",
                "startDate": "2026-06-01",
                "duration": 14
              }
            ],
            "iterations": [
              {
                "id": "11111a1b",
                "title": "Iteration 2",
                "startDate": "2026-06-15",
                "duration": 14
              }
            ]
          }
        }
      }
    }
  }
}
`),

		gqltest.WithMutate("CreateIterations", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, &captured))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "data": { "updateProjectV2Field": { "clientMutationId": null } } }`))
		}),
	)
	require.NoError(t, err)

	output := &bytes.Buffer{}
	require.NoError(t, Run(
		gql,
		"https://github.com/users/notomo/projects/1",
		"Iteration",
		2,  // count
		14, // duration
		"", // start-date (default: after last existing)
		"Iteration ",
		false, // dry-run
		output,
	))

	// The mutation must preserve every existing iteration verbatim and append the new ones.
	want := []IterationInput{
		{Title: "Iteration 1", StartDate: "2026-06-01", Duration: 14},
		{Title: "Iteration 2", StartDate: "2026-06-15", Duration: 14},
		{Title: "Iteration 3", StartDate: "2026-06-29", Duration: 14},
		{Title: "Iteration 4", StartDate: "2026-07-13", Duration: 14},
	}
	assert.Equal(t, want, captured.Variables.IterationConfiguration.Iterations)
	assert.Equal(t, "PVTIF_22222222222222222222", captured.Variables.FieldID)
	// Field-level defaults are preserved.
	assert.Equal(t, "2026-06-01", captured.Variables.IterationConfiguration.StartDate)
	assert.Equal(t, 14, captured.Variables.IterationConfiguration.Duration)

	wantOutput := `
Created 2 iteration(s) (kept 2 existing):
- Iteration 3 (start date: 2026-06-29, duration: 14)
- Iteration 4 (start date: 2026-07-13, duration: 14)
`
	assert.Equal(t, wantOutput, output.String())
}

func TestRunDryRun(t *testing.T) {
	mutateCalled := false
	gql, err := gqltest.New(
		t,
		gqltest.WithQueryOK("GetIterationField", `
{
  "data": {
    "user": {
      "projectV2": {
        "field": {
          "id": "PVTIF_22222222222222222222",
          "configuration": {
            "duration": 14,
            "completedIterations": [],
            "iterations": []
          }
        }
      }
    }
  }
}
`),
		gqltest.WithMutate("CreateIterations", func(w http.ResponseWriter, r *http.Request) {
			mutateCalled = true
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "data": { "updateProjectV2Field": { "clientMutationId": null } } }`))
		}),
	)
	require.NoError(t, err)

	output := &bytes.Buffer{}
	require.NoError(t, Run(
		gql,
		"https://github.com/users/notomo/projects/1",
		"Iteration",
		1,
		7,
		"2026-06-01",
		"Iteration ",
		true, // dry-run
		output,
	))

	assert.False(t, mutateCalled, "dry-run must not send the mutation")
}
