package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

// Client interface
type Client interface {
	Run() error
	Stop() error
	Options() *Options
	Watch() <-chan Event
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
	observe := v.client.Events(api.NamespaceAll)
	options := api.ListOptions{
		LabelSelector: labels.Everything(),
	}
	w, err := observe.Watch(options)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case res, _ := <-w.ResultChan():
			e, ok := res.Object.(*api.Event)
			if !ok {
				continue
			}
			var color string
			switch {
			case e.Count > 1:
				continue
			case e.Reason == "BackOff" && e.Count == 3:
				color = "danger"
			case res.Type == watch.Added:
				color = "good"
			case res.Type == watch.Deleted:
				color = "warning"
			case e.Reason == "SuccessfulCreate":
				color = "good"
			case e.Reason == "NodeReady":
				color = "good"
			case e.Reason == "NodeNotReady":
				color = "danger"
			case e.Reason == "NodeOutOfDisk":
				color = "danger"
			}
			switch e.Source.Component {
			case "kubelet", "controllermanager", "default-scheduler":
				continue
			}
			v.events <- &event{
				level:        color,
				name:         e.GetObjectMeta().GetName(),
				namespace:    e.GetObjectMeta().GetNamespace(),
				generatename: e.GetObjectMeta().GetGenerateName(),
				reason:       e.Reason,
				message:      e.Message,
				kind:         e.InvolvedObject.Kind,
				component:    e.Source.Component,
			}
		}
	}
}
