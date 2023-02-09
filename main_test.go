package main

import (
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func configureDefaultArgs() {
	mutatorConfig.MetricNameTemplate = "check_name"
	mutatorConfig.MetricNameTemplateOccurrences = "check_name_occurrences"
	mutatorConfig.TagNameEntity = "entity"
	mutatorConfig.TagNameCheck = "check"
	mutatorConfig.TagNameState = "state"
	mutatorConfig.TagNameOccurrences = "occurrences"
	mutatorConfig.ShowOccurrencesAsMetric = "false"
	mutatorConfig.ShowOccurrencesAsTag = "true"
	mutatorConfig.EnableWatermark = "true"
}

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := types.FixtureEvent("entity1", "check1")
	assert.Error(checkArgs(event))
	configureDefaultArgs()
	assert.NoError(checkArgs(event))
}

func TestExecuteMutator(t *testing.T) {
	assert := assert.New(t)

	// Event with no metrics
	event := types.FixtureEvent("entity1", "check1")
	configureDefaultArgs()
	ev, err := executeMutator(event)
	assert.NoError(err)
	assert.Equal(len(ev.Metrics.Points), 1)
	var mps []string
	for _, v := range ev.Metrics.Points {
		mps = append(mps, v.Name)
	}
	assert.Contains(mps, "check_name")

	// Event with existing metrics
	event.Metrics = types.FixtureMetrics()
	ev, err = executeMutator(event)
	assert.NoError(err)
	assert.Equal(len(ev.Metrics.Points), 2)
	mps = nil
	for _, v := range ev.Metrics.Points {
		mps = append(mps, v.Name)
	}
	assert.Contains(mps, "answer")
	assert.Contains(mps, "check_name")

	// Event without check
	event.Check = nil
	_, err = executeMutator(event)
	assert.Error(err)
}
