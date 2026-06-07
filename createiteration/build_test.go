package createiteration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildIterations(t *testing.T) {
	existing := []Iteration{
		{ID: "00000a0b", Title: "Iteration 1", StartDate: "2026-06-01", Duration: 14},
		{ID: "11111a1b", Title: "Iteration 2", StartDate: "2026-06-15", Duration: 14},
	}
	config := IterationFieldConfiguration{
		Duration:  14,
		StartDate: "2026-06-01",
		Existing:  existing,
	}

	t.Run("keeps existing iterations verbatim and appends new ones", func(t *testing.T) {
		got, err := BuildIterations(config, 2, 14, "", "Iteration ")
		require.NoError(t, err)

		want := []IterationInput{
			// existing, re-sent verbatim (no loss)
			{Title: "Iteration 1", StartDate: "2026-06-01", Duration: 14},
			{Title: "Iteration 2", StartDate: "2026-06-15", Duration: 14},
			// new, appended after the last existing iteration ends
			{Title: "Iteration 3", StartDate: "2026-06-29", Duration: 14},
			{Title: "Iteration 4", StartDate: "2026-07-13", Duration: 14},
		}
		assert.Equal(t, want, got)
	})

	t.Run("existing iterations are never dropped", func(t *testing.T) {
		got, err := BuildIterations(config, 1, 14, "", "Iteration ")
		require.NoError(t, err)

		for _, e := range existing {
			assert.Contains(t, got, IterationInput{
				Title:     e.Title,
				StartDate: e.StartDate,
				Duration:  e.Duration,
			})
		}
		assert.Len(t, got, len(existing)+1)
	})

	t.Run("duration 0 inherits the field configuration", func(t *testing.T) {
		got, err := BuildIterations(config, 1, 0, "", "Iteration ")
		require.NoError(t, err)
		assert.Equal(t, 14, got[len(got)-1].Duration)
	})

	t.Run("explicit start date is honored", func(t *testing.T) {
		got, err := BuildIterations(config, 1, 7, "2026-08-01", "Iteration ")
		require.NoError(t, err)
		last := got[len(got)-1]
		assert.Equal(t, "2026-08-01", last.StartDate)
		assert.Equal(t, 7, last.Duration)
	})

	t.Run("title prefix is configurable", func(t *testing.T) {
		got, err := BuildIterations(config, 1, 14, "", "Sprint ")
		require.NoError(t, err)
		assert.Equal(t, "Sprint 3", got[len(got)-1].Title)
	})

	t.Run("no existing iterations uses the configuration start date", func(t *testing.T) {
		empty := IterationFieldConfiguration{Duration: 7, StartDate: "2026-01-05"}
		got, err := BuildIterations(empty, 2, 0, "", "Iteration ")
		require.NoError(t, err)

		want := []IterationInput{
			{Title: "Iteration 1", StartDate: "2026-01-05", Duration: 7},
			{Title: "Iteration 2", StartDate: "2026-01-12", Duration: 7},
		}
		assert.Equal(t, want, got)
	})

	t.Run("errors", func(t *testing.T) {
		cases := []struct {
			name        string
			config      IterationFieldConfiguration
			count       int
			duration    int
			startDate   string
			titlePrefix string
		}{
			{name: "count below 1", config: config, count: 0, duration: 14},
			{name: "no resolvable duration", config: IterationFieldConfiguration{StartDate: "2026-01-01"}, count: 1, duration: 0},
			{name: "no resolvable start date", config: IterationFieldConfiguration{Duration: 7}, count: 1, duration: 0},
			{name: "invalid start date", config: config, count: 1, duration: 14, startDate: "2026/01/01"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				_, err := BuildIterations(c.config, c.count, c.duration, c.startDate, c.titlePrefix)
				require.Error(t, err)
			})
		}
	})
}
