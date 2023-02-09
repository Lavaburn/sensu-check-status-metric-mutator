package main

import (
	"fmt"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu-community/sensu-plugin-sdk/templates"
	"github.com/sensu/sensu-go/types"
)

// Config represents the mutator plugin config.
type Config struct {
	sensu.PluginConfig
	
	MetricNameTemplate string
	MetricNameTemplateOccurrences string
	
	TagNameEntity string
	TagNameCheck string
	TagNameState string
	TagNameOccurrences string
	
	// Not using bool because ShowOccurrencesAsTag is default true for backwards compatibility
	ShowOccurrencesAsMetric string
	ShowOccurrencesAsTag string
	EnableWatermark string
}

var (
	mutatorConfig = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-check-status-metric-mutator",
			Short:    "Sensu Check Status Metric Mutator",
			Keyspace: "sensu.io/plugins/sensu-check-status-metric-mutator/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "metric-name-template",
			Env:       "METRIC_NAME_TEMPLATE",
			Argument:  "metric-name-template",
			Shorthand: "t",
			Default:   "{{.Check.Name}}.status",
			Usage:     "Template for naming the metric point for the check status",
			Value:     &mutatorConfig.MetricNameTemplate,
		},
		{
			Path:      "metric-name-template-occurrences",
			Env:       "METRIC_NAME_TEMPLATE_OCCURRENCES",
			Argument:  "metric-name-template-occurrences",
			Shorthand: "u",
			Default:   "{{.Check.Name}}.occurrences",
			Usage:     "Template for naming the metric point for the check status",
			Value:     &mutatorConfig.MetricNameTemplateOccurrences,
		},
		{
			Path:      "tag-name-entity",
			Env:       "TAG_NAME_ENTITY",
			Argument:  "tag-name-entity",
			Shorthand: "e",
			Default:   "entity",
			Usage:     "The tag name that contains the entity name",
			Value:     &mutatorConfig.TagNameEntity,
		},	
		{
			Path:      "tag-name-check",
			Env:       "TAG_NAME_CHECK",
			Argument:  "tag-name-check",
			Shorthand: "c",
			Default:   "check",
			Usage:     "The tag name that contains the check name",
			Value:     &mutatorConfig.TagNameCheck,
		},		
		{
			Path:      "tag-name-state",
			Env:       "TAG_NAME_STATE",
			Argument:  "tag-name-state",
			Shorthand: "s",
			Default:   "state",
			Usage:     "The tag name that contains the state name",
			Value:     &mutatorConfig.TagNameState,
		},		
		{
			Path:      "tag-name-occurrences",
			Env:       "TAG_NAME_OCCURRENCES",
			Argument:  "tag-name-occurrences",
			Shorthand: "o",
			Default:   "occurrences",
			Usage:     "The tag name that contains the occurrences name",
			Value:     &mutatorConfig.TagNameOccurrences,
		},
		{
			Path:      "show-occurrences-as-metric",
			Env:       "SHOW_OCCURRENCES_AS_METRIC",
			Argument:  "show-occurrences-as-metric",
			Shorthand: "m",
			Default:   "false",
			Usage:     "Whether to add a metric for Occurrences",
			Value:     &mutatorConfig.ShowOccurrencesAsMetric,
		},		
		{
			Path:      "show-occurrences-as-tag",
			Env:       "SHOW_OCCURRENCES_AS_TAG",
			Argument:  "show-occurrences-as-tag",
			Shorthand: "r",
			Default:   "true",
			Usage:     "Whether to add a tag for Occurrences",
			Value:     &mutatorConfig.ShowOccurrencesAsTag,
		},
		{
			Path:      "enable-watermark",
			Env:       "ENABLE_WATERMARK",
			Argument:  "enable-watermark",
			Shorthand: "w",
			Default:   "true",
			Usage:     "The tag name that contains the occurrences name",
			Value:     &mutatorConfig.EnableWatermark,
		},
	}
)

func main() {
	mutator := sensu.NewGoMutator(&mutatorConfig.PluginConfig, options, checkArgs, executeMutator)
	mutator.Execute()
}

func checkArgs(_ *types.Event) error {
	if len(mutatorConfig.MetricNameTemplate) == 0 {
		return fmt.Errorf("--MetricNameTemplate or METRIC_NAME_TEMPLATE environment variable is required")
	}
	if len(mutatorConfig.MetricNameTemplateOccurrences) == 0 {
		return fmt.Errorf("--MetricNameTemplateOccurrences or METRIC_NAME_TEMPLATE_OCCURRENCES environment variable is required")
	}
		
	if len(mutatorConfig.TagNameEntity) == 0 {
		return fmt.Errorf("--TagNameEntity or TAG_NAME_ENTITY environment variable is required")
	}
	if len(mutatorConfig.TagNameCheck) == 0 {
		return fmt.Errorf("--TagNameCheck or TAG_NAME_CHECK environment variable is required")
	}	
	if len(mutatorConfig.TagNameState) == 0 {
		return fmt.Errorf("--TagNameState or TAG_NAME_STATE environment variable is required")
	}	
	if len(mutatorConfig.TagNameOccurrences) == 0 {
		return fmt.Errorf("--TagNameOccurrences or TAG_NAME_OCCURRENCES environment variable is required")
	}
	
	if mutatorConfig.ShowOccurrencesAsMetric != "true" && mutatorConfig.ShowOccurrencesAsMetric != "false" {
		return fmt.Errorf("--ShowOccurrencesAsMetric or SHOW_OCCURRENCES_AS_METRIC environment variable is required and should be -true- or -false-")
	}
	if mutatorConfig.ShowOccurrencesAsTag != "true" && mutatorConfig.ShowOccurrencesAsTag != "false" {
		return fmt.Errorf("--ShowOccurrencesAsTag or SHOW_OCCURRENCES_AS_TAG environment variable is required and should be -true- or -false-")
	}
	if mutatorConfig.EnableWatermark != "true" && mutatorConfig.EnableWatermark != "false" {
		return fmt.Errorf("--EnableWatermark or ENABLE_WATERMARK environment variable is required and should be -true- or -false-")
	}
	
	return nil
}

func executeMutator(event *types.Event) (*types.Event, error) {
	// Sanity Check
	if !event.HasCheck() {
		return &types.Event{}, fmt.Errorf("Event does not have a check defined.")
	}
	
	// Replace templates
	metricName, err := templates.EvalTemplate("metricName", mutatorConfig.MetricNameTemplate, event)
	if err != nil {
		return &types.Event{}, fmt.Errorf("Failed to evalutate template: %v", err)
	}
	metricNameOccurrences, err := templates.EvalTemplate("metricNameOccurrences", mutatorConfig.MetricNameTemplateOccurrences, event)
	if err != nil {
		return &types.Event{}, fmt.Errorf("Failed to evalutate template: %v", err)
	}

	// Possible TODO:  replace any spaces, periods, dashes from the templated metricName

	// Initialize (if no metrics yet)
	if !event.HasMetrics() {
		event.Metrics = new(types.Metrics)
	}

	// Provide some extra information in the tags
	mt := make([]*types.MetricTag, 0)
	
	mt = append(mt, &types.MetricTag{Name: mutatorConfig.TagNameEntity, Value: event.Entity.Name})
	mt = append(mt, &types.MetricTag{Name: mutatorConfig.TagNameCheck, Value: event.Check.Name})
	mt = append(mt, &types.MetricTag{Name: mutatorConfig.TagNameState, Value: event.Check.State})
	
	if mutatorConfig.ShowOccurrencesAsTag == "true" {
		mt = append(mt, &types.MetricTag{Name: mutatorConfig.TagNameOccurrences, Value: fmt.Sprintf("%d", event.Check.Occurrences)})
		if mutatorConfig.EnableWatermark == "true" {
			mt = append(mt, &types.MetricTag{Name: mutatorConfig.TagNameOccurrences+"_watermark", Value: fmt.Sprintf("%d", event.Check.OccurrencesWatermark)})
		}
	}

	mp := &types.MetricPoint{
		Name:      metricName,
		Value:     float64(event.Check.Status),
		Timestamp: event.Timestamp,
		Tags:      mt,
	}
	event.Metrics.Points = append(event.Metrics.Points, mp)
	
	if mutatorConfig.ShowOccurrencesAsMetric == "true" {
		mp2 := &types.MetricPoint{
			Name:      metricNameOccurrences,
			Value:     float64(event.Check.Occurrences),
			Timestamp: event.Timestamp,
			Tags:      mt,
		}
		event.Metrics.Points = append(event.Metrics.Points, mp2)
		if mutatorConfig.EnableWatermark == "true" {
			mp3 := &types.MetricPoint{
				Name:      metricNameOccurrences+"_watermark",
				Value:     float64(event.Check.OccurrencesWatermark),
				Timestamp: event.Timestamp,
				Tags:      mt,
			}
			event.Metrics.Points = append(event.Metrics.Points, mp3)
		}
	}
	return event, nil
}
