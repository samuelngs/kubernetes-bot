package client

import (
	"time"

	"k8s.io/kubernetes/pkg/api"
)

// EventType type
type EventType int8

// Event types
const (
	EventNewPod EventType = iota
	EventDelPod
	EventUpdatePod
	EventNewNode
	EventDelNode
	EventUpdateNode
)

// String returns event type in string format
func (e EventType) String() string {
	switch e {
	case EventNewPod:
		return "new-pod"
	case EventDelPod:
		return "delete-pod"
	case EventUpdatePod:
		return "update-pod"
	case EventNewNode:
		return "new-node"
	case EventDelNode:
		return "del-node"
	case EventUpdateNode:
		return "update-node"
	default:
		return "unknown"
	}
}

// Event interface
type Event interface {
	Type() EventType
	Namespace() string
	Pod() *api.Pod
	Node() *api.Node
	Name() string
	GenerateName() string
	Labels() map[string]string
	Annotations() map[string]string
	NodeName() string
	Time() time.Time
}

// client event
type event struct {
	typ  EventType
	pod  *api.Pod
	node *api.Node
}

// type returns the type of the event
func (v *event) Type() EventType {
	return v.typ
}

// Namespace returns the namespace name
func (v *event) Namespace() string {
	switch {
	case v.pod != nil:
		return v.pod.Namespace
	case v.node != nil:
		return v.node.Namespace
	default:
		return "default"
	}
}

// Pod returns the pod information
func (v *event) Pod() *api.Pod {
	return v.pod
}

func (v *event) Node() *api.Node {
	return v.node
}

// Name returns the name of the pod
func (v *event) Name() string {
	switch {
	case v.pod != nil:
		return v.pod.Name
	case v.node != nil:
		return v.node.Name
	default:
		return "unknown"
	}
}

// GenerateName returns the generated name of the pod
func (v *event) GenerateName() string {
	switch {
	case v.pod != nil:
		return v.pod.GenerateName
	case v.node != nil:
		return v.node.GenerateName
	default:
		return "unknown"
	}
}

// Labels to retrieve all the pod label names
func (v *event) Labels() map[string]string {
	switch {
	case v.pod != nil:
		return v.pod.Labels
	case v.node != nil:
		return v.node.Labels
	default:
		return make(map[string]string, 0)
	}
}

// Annotations to retrieve pod annotations
func (v *event) Annotations() map[string]string {
	switch {
	case v.pod != nil:
		return v.pod.Annotations
	case v.node != nil:
		return v.node.Annotations
	default:
		return make(map[string]string, 0)
	}
}

// NodeName returns node server name
func (v *event) NodeName() string {
	switch {
	case v.pod != nil:
		return v.pod.Spec.NodeName
	case v.node != nil:
		return v.node.Name
	default:
		return "unknown"
	}
}

// Time returns event time
func (v *event) Time() time.Time {
	switch {
	case v.pod != nil:
		return time.Time(v.pod.Status.StartTime.Time)
	default:
		return time.Now().UTC()
	}
}
