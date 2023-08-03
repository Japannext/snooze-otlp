package server

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
)

// Convert a string to the AnyValue type used
// in opentelemetry
func AnyString(s string) *commonv1.AnyValue {
	return &commonv1.AnyValue{
		Value: &commonv1.AnyValue_StringValue{
			StringValue: s,
		},
	}
}

// Create a KeyValue array from a map[string]string
func Kv(d map[string]string) []*commonv1.KeyValue {
	var kvs []*commonv1.KeyValue
	for key, value := range d {
		kvs = append(kvs, &commonv1.KeyValue{
			Key:   key,
			Value: AnyString(value),
		})
	}
	return kvs
}

func Hex(b []byte) string {
	return hex.EncodeToString(b)
}

// Decode an hex string into bytes
func DeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func Base64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Decode a Base64 string into bytes
func DeBase64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func kvToMap(kvs []*commonv1.KeyValue) map[string]string {
	var mapping map[string]string = make(map[string]string)

	for _, kv := range kvs {
		if kv != nil {
			mapping[kv.Key] = kv.Value.GetStringValue()
		}
	}
	return mapping
}

// Check if a map has at least one key prefixed in a certain way
func hasPrefixedKey(m map[string]string, prefix string) bool {
	for k, _ := range m {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

// Convert a Unix time in nanoseconds to a format snooze can understand.
func formatTime(nano uint64) string {
	t := time.Unix(0, int64(nano))
	return t.Format("2006-01-02T15:04:05")
}
