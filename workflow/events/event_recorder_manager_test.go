package events

import (
	"testing"

	apiv1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/require"
)

const aggregationWithAnnotationsEnvKey = "EVENT_AGGREGATION_WITH_ANNOTATIONS"

func TestCustomEventAggregatorFuncWithAnnotations(t *testing.T) {
	event := apiv1.Event{}
	key, msg := customEventAggregatorFuncWithAnnotations(&event)
	require.Equal(t, "", key)
	require.Equal(t, "", msg)

	event.Source = apiv1.EventSource{Component: "component1", Host: "host1"}
	event.InvolvedObject.Name = "name1"
	event.Message = "message1"

	key, msg = customEventAggregatorFuncWithAnnotations(&event)
	require.Equal(t, "component1host1name1", key)
	require.Equal(t, "message1", msg)

	// Test default behavior where annotations are not used for aggregation
	event.ObjectMeta.Annotations = map[string]string{"key1": "val1", "key2": "val2"}
	key, msg = customEventAggregatorFuncWithAnnotations(&event)
	require.Equal(t, "component1host1name1", key)
	require.Equal(t, "message1", msg)

	t.Setenv(aggregationWithAnnotationsEnvKey, "true")
	key, msg = customEventAggregatorFuncWithAnnotations(&event)
	require.Equal(t, "component1host1name1val1val2", key)
	require.Equal(t, "message1", msg)

	// Test annotations with values in different order
	t.Setenv(aggregationWithAnnotationsEnvKey, "true")
	event.ObjectMeta.Annotations = map[string]string{"key2": "val2", "key1": "val1"}
	key, msg = customEventAggregatorFuncWithAnnotations(&event)
	require.Equal(t, "component1host1name1val1val2", key)
	require.Equal(t, "message1", msg)
}
