package model

import "k8s.io/apimachinery/pkg/runtime"

type MicroFrontendAttribute struct {
	Name  string               `json:"name"`
	Value runtime.RawExtension `json:"value"`
}
