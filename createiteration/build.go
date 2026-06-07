package createiteration

import (
	"fmt"
	"time"
)

// IterationInput is one element of the iteration configuration sent to GitHub.
// It matches the GraphQL input object `ProjectV2Iteration` (no id field).
type IterationInput struct {
	Title     string `json:"title"`
	StartDate string `json:"startDate"`
	Duration  int    `json:"duration"`
}

func shiftDate(date string, offsetDays int) (string, error) {
	at, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return "", err
	}
	return at.AddDate(0, 0, offsetDays).Format(time.DateOnly), nil
}

// BuildIterations builds the iterations array for updateProjectV2Field.
// It re-sends every existing iteration verbatim (so existing iterations and their
// item associations are preserved) and then appends `count` new iterations.
func BuildIterations(
	config IterationFieldConfiguration,
	count int,
	duration int,
	startDate string,
	titlePrefix string,
) ([]IterationInput, error) {
	if count < 1 {
		return nil, fmt.Errorf("count must be >= 1, but actual %d", count)
	}

	resolvedDuration := duration
	if resolvedDuration == 0 {
		resolvedDuration = config.Duration
	}
	if resolvedDuration < 1 {
		return nil, fmt.Errorf("duration must be >= 1 (or set via field configuration), but actual %d", resolvedDuration)
	}

	firstStart, err := firstStartDate(config, startDate)
	if err != nil {
		return nil, err
	}

	inputs := make([]IterationInput, 0, len(config.Existing)+count)
	for _, e := range config.Existing {
		inputs = append(inputs, IterationInput{
			Title:     e.Title,
			StartDate: e.StartDate,
			Duration:  e.Duration,
		})
	}

	for i := range count {
		start, err := shiftDate(firstStart, i*resolvedDuration)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, IterationInput{
			Title:     fmt.Sprintf("%s%d", titlePrefix, len(config.Existing)+i+1),
			StartDate: start,
			Duration:  resolvedDuration,
		})
	}

	return inputs, nil
}

func firstStartDate(config IterationFieldConfiguration, startDate string) (string, error) {
	if startDate != "" {
		if _, err := time.Parse(time.DateOnly, startDate); err != nil {
			return "", fmt.Errorf("invalid start-date (want yyyy-mm-dd): %w", err)
		}
		return startDate, nil
	}

	if len(config.Existing) > 0 {
		last := config.Existing[0]
		for _, e := range config.Existing[1:] {
			if e.StartDate > last.StartDate {
				last = e
			}
		}
		return shiftDate(last.StartDate, last.Duration)
	}

	if config.StartDate == "" {
		return "", fmt.Errorf("cannot determine start date: specify -start-date")
	}
	return config.StartDate, nil
}
