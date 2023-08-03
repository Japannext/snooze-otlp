package server

import (
	"fmt"

	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	logv1 "go.opentelemetry.io/proto/otlp/logs/v1"
	resv1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

var ok bool

type SnoozeAlertV1 struct {
	Source     string            `json:"source"`
	Timestamp  string            `json:"timestamp"`
	Host       string            `json:"host"`
	Process    string            `json:"process"`
	Severity   string            `json:"severity"`
	Message    string            `json:"message"`
	Attributes map[string]string `json:"attributes"`
}

// Populate the alert with the kubernetes fields
func populateKubernetes(alert *SnoozeAlertV1, ra map[string]string) {

	if name, ok := ra["k8s.deployment.name"]; ok {
		alert.Attributes["k8s_kind"] = "deployment"
		alert.Attributes["k8s_deployment_name"] = name
		alert.Process = fmt.Sprintf("deploy/%s", name)

	} else if name, ok := ra["k8s.statefulset.name"]; ok {
		alert.Attributes["k8s_kind"] = "statefulset"
		alert.Attributes["k8s_statefulset_name"] = name
		alert.Process = fmt.Sprintf("sts/%s", name)

	} else if name, ok := ra["k8s.daemonset.name"]; ok {
		alert.Attributes["k8s_kind"] = "daemonset"
		alert.Attributes["k8s_daemonset_name"] = name
		alert.Process = fmt.Sprintf("ds/%s", name)

	} else if name, ok := ra["k8s.cronjob.name"]; ok {
		alert.Attributes["k8s_kind"] = "cronjob"
		alert.Attributes["k8s_cronjob_name"] = name
		alert.Process = fmt.Sprintf("cronjob/%s", name)

	} else if name, ok := ra["k8s.job.name"]; ok {
		alert.Attributes["k8s_kind"] = "job"
		alert.Attributes["k8s_job_name"] = name
		alert.Process = fmt.Sprintf("job/%s", name)

	}

	var cluster string
	var ns string

	if cluster, ok = ra["k8s.cluster.name"]; ok {
		alert.Attributes["k8s_cluster_name"] = cluster
	} else {
		cluster = "-"
	}
	if ns, ok = ra["k8s.namespace.name"]; ok {
		alert.Attributes["k8s_namespace_name"] = ns
	} else {
		ns = "-"
	}

	alert.Host = fmt.Sprintf("%s/%s", cluster, ns)

	if name, ok := ra["k8s.replicaset.name"]; ok {
		alert.Attributes["k8s_replicaset_name"] = name
	}
	if name, ok := ra["k8s.pod.name"]; ok {
		alert.Attributes["k8s_pod_name"] = name
	}
	if name, ok := ra["k8s.container.name"]; ok {
		alert.Attributes["k8s_container_name"] = name
	}

}

// Convert an opentelemetry record log (with resource and scope contexts) to a snooze alert
func convertAlert(resource *resv1.Resource, scope *commonv1.InstrumentationScope, lr *logv1.LogRecord) SnoozeAlertV1 {
	var alert SnoozeAlertV1

	alert.Attributes = make(map[string]string)

	// Building attributes
	ra := kvToMap(resource.Attributes)
	//sa := kvToMap(scope.Attributes)
	//la := kvToMap(lr.Attributes)

	if hasPrefixedKey(ra, "k8s.") {
		populateKubernetes(&alert, ra)
	}

	alert.Source = "otel"

	alert.Timestamp = formatTime(lr.TimeUnixNano)
	if lr.ObservedTimeUnixNano != 0 {
		alert.Attributes["observed_time"] = formatTime(lr.ObservedTimeUnixNano)
	}

	if name, ok := ra["service.name"]; ok {
		alert.Process = name
	}

	alert.Attributes["trace_id"] = Hex(lr.TraceId)
	alert.Attributes["span_id"] = Hex(lr.SpanId)
	alert.Severity = lr.SeverityText
	alert.Message = lr.Body.GetStringValue()

	return alert
}
