// Package feature implements feature functionality.
package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// Every feature gate should add method here following this template:
	//
	// // owner: @username
	// // alpha: v1.4
	// MyFeature featuregate.Feature = "MyFeature"
	//
	// Feature gates should be listed in alphabetical, case-sensitive
	// (upper before any lower case character) order. This reduces the risk
	// of code conflicts because changes are more likely to be scattered
	// across the file.

	// owner: @colin404
	// alpha: v1.26
	//
	// MachinePool is a feature gate for MachinePool functionality.
	MachinePool featuregate.Feature = "MachinePool"
)

func init() {
	// runtime.Must(utilfeature.DefaultMutableFeatureGate.Add(defaultstreamingFeatureGates))
	runtime.Must(DefaultMutableFeatureGate.Add(defaultstreamingFeatureGates))
}

// defaultstreamingFeatureGates consists of all known streaming-specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultstreamingFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	// Every feature should be initiated here:
	MachinePool: {Default: false, PreRelease: featuregate.Alpha},
	// ClusterResourceSet:             {Default: true, PreRelease: featuregate.Beta},
}
