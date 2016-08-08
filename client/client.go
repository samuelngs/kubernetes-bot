package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/util/wait"
)

// Client interface
type Client interface {
	Run() error
	Stop() error
	Options() *Options
	Watch() <-chan Event
}

// resource types
type resource int8

// list of resources types
const (
	nodes resource = iota
	pods
)

func (v resource) name() string {
	switch v {
	case nodes:
		return "nodes"
	case pods:
		return "pods"
	default:
		return "unknown"
	}
}

// resource kind
type kind struct {
	resource  resource
	namespace string
}

// list of resources kind
var kinds = []kind{
	kind{
		resource:  nodes,
		namespace: api.NamespaceAll,
	},
	kind{
		resource:  pods,
		namespace: api.NamespaceAll,
	},
}

type client struct {
	option *Options
	client *unversioned.Client
	events chan Event
}

// New client
func New(opts ...Option) (Client, error) {
	c := &client{
		option: newOptions(opts...),
		events: make(chan Event),
	}
	switch client, err := unversioned.New(c.option.RestOpt()); {
	case err != nil:
		return nil, err
	default:
		c.client = client
	}
	return c, nil
}

// Watch events
func (v *client) Watch() <-chan Event {
	return v.events
}

// Options return client options
func (v *client) Options() *Options {
	return v.option
}

// Run to start client
func (v *client) Run() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	go v.monitor()
	log.Printf("Received signal %s", <-ch)
	return v.Stop()
}

// Stop to close client
func (v *client) Stop() error {
	defer close(v.events)
	return nil
}

// monitor handler
func (v *client) monitor() {
	for _, kind := range kinds {
		switch kind.resource {
		case nodes:
			handler := framework.ResourceEventHandlerFuncs{
				AddFunc:    v.addNode,
				DeleteFunc: v.delNode,
				UpdateFunc: v.updateNode,
			}
			_, controller := framework.NewInformer(
				cache.NewListWatchFromClient(v.client, kind.resource.name(), kind.namespace, fields.Everything()),
				new(api.Node),
				v.option.Interval,
				handler,
			)
			go controller.Run(wait.NeverStop)
		case pods:
			handler := framework.ResourceEventHandlerFuncs{
				AddFunc:    v.addPod,
				DeleteFunc: v.delPod,
				UpdateFunc: v.updatePod,
			}
			_, controller := framework.NewInformer(
				cache.NewListWatchFromClient(v.client, kind.resource.name(), kind.namespace, fields.Everything()),
				new(api.Pod),
				v.option.Interval,
				handler,
			)
			go controller.Run(wait.NeverStop)
		}
	}
}

// addNode handler
func (v *client) addNode(obj interface{}) {
	v.events <- &event{
		typ:  EventNewNode,
		node: obj.(*api.Node),
	}
}

// delNode handler
func (v *client) delNode(obj interface{}) {
	v.events <- &event{
		typ:  EventDelNode,
		node: obj.(*api.Node),
	}
}

// updateNode handler
func (v *client) updateNode(old, new interface{}) {
	v.events <- &event{
		typ:  EventUpdateNode,
		node: new.(*api.Node),
	}
}

// addPod handler
func (v *client) addPod(obj interface{}) {
	v.events <- &event{
		typ: EventNewPod,
		pod: obj.(*api.Pod),
	}
}

// delPod handler
func (v *client) delPod(obj interface{}) {
	v.events <- &event{
		typ: EventDelPod,
		pod: obj.(*api.Pod),
	}
}

// updatePod handler
func (v *client) updatePod(old, new interface{}) {
	v.events <- &event{
		typ: EventUpdatePod,
		pod: new.(*api.Pod),
	}
}
