// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type      = component.MustNewType("k8sattributes")
	ScopeName = "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"
)

const (
	LogsStability    = component.StabilityLevelBeta
	MetricsStability = component.StabilityLevelBeta
	TracesStability  = component.StabilityLevelBeta
)
