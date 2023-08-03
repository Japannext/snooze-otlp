package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resv1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

var alerttests = []struct {
	name      string
	resource  *resv1.Resource
	scope     *commonv1.InstrumentationScope
	logRecord *logv1.LogRecord
	expected  SnoozeAlertV1
}{
	{
		name: "kubernetes deploy log",
		resource: &resv1.Resource{
			Attributes: Kv(map[string]string{
				"k8s.cluster.name":    "dev",
				"k8s.namespace.name":  "myns",
				"k8s.deployment.name": "myapp",
				"k8s.replicaset.name": "myapp-1234",
				"k8s.pod.name":        "myapp-1234-5678",
				"k8s.container.name":  "myapp",
			}),
		},
		scope: &commonv1.InstrumentationScope{},
		logRecord: &logv1.LogRecord{
			TimeUnixNano:   1692320543952328400,
			SeverityNumber: 17,
			SeverityText:   "error",
			Body:           AnyString("Error loading /etc/config.yaml: No such file or directory"),
			TraceId:        DeHex("99197824582792068c2aa8880f8c5300"),
			SpanId:         DeHex("ec7e83dfdb711f0e"),
		},
		expected: SnoozeAlertV1{
			Source:    "otel",
			Timestamp: "2023-08-18T10:02:23",
			Host:      "dev/myns",
			Process:   "deploy/myapp",
			Severity:  "error",
			Message:   "Error loading /etc/config.yaml: No such file or directory",
			Attributes: map[string]string{
				"k8s_cluster_name":    "dev",
				"k8s_namespace_name":  "myns",
				"k8s_kind":            "deployment",
				"k8s_deployment_name": "myapp",
				"k8s_replicaset_name": "myapp-1234",
				"k8s_pod_name":        "myapp-1234-5678",
				"k8s_container_name":  "myapp",
				"trace_id":            "99197824582792068c2aa8880f8c5300",
				"span_id":             "ec7e83dfdb711f0e",
			},
		},
	},
	{
		name: "kubernetes statefulset log",
		resource: &resv1.Resource{
			Attributes: Kv(map[string]string{
				"k8s.cluster.name":     "dev",
				"k8s.namespace.name":   "myns",
				"k8s.statefulset.name": "mydb",
				"k8s.pod.name":         "mydb-0",
				"k8s.container.name":   "db",
			}),
		},
		scope: &commonv1.InstrumentationScope{},
		logRecord: &logv1.LogRecord{
			TimeUnixNano:   1692320543952328400,
			SeverityNumber: 17,
			SeverityText:   "error",
			Body:           AnyString("Client connection timed out"),
			TraceId:        DeHex("31b25a9c5c9657e7809208dc4ffc2cf3"),
			SpanId:         DeHex("15a18a9501602d55"),
		},
		expected: SnoozeAlertV1{
			Source:    "otel",
			Timestamp: "2023-08-18T10:02:23",
			Host:      "dev/myns",
			Process:   "sts/mydb",
			Severity:  "error",
			Message:   "Client connection timed out",
			Attributes: map[string]string{
				"k8s_cluster_name":     "dev",
				"k8s_namespace_name":   "myns",
				"k8s_kind":             "statefulset",
				"k8s_statefulset_name": "mydb",
				"k8s_pod_name":         "mydb-0",
				"k8s_container_name":   "db",
				"trace_id":             "31b25a9c5c9657e7809208dc4ffc2cf3",
				"span_id":              "15a18a9501602d55",
			},
		},
	},
}

func TestConvertAlert(t *testing.T) {

	for _, tt := range alerttests {
		t.Run(tt.name, func(t *testing.T) {
			alert := convertAlert(tt.resource, tt.scope, tt.logRecord)
			assert.Equal(t, tt.expected, alert)
		})
	}

}
