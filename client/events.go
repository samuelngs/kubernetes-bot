package client

// Event interface
type Event interface {
	Level() string
	Name() string
	Namespace() string
	GenerateName() string
	Reason() string
	Message() string
	Kind() string
	Component() string
}

// client event
type event struct {
	level        string
	name         string
	namespace    string
	generatename string
	reason       string
	message      string
	kind         string
	component    string
}

func (v *event) Level() string {
	return v.level
}

func (v *event) Name() string {
	return v.name
}

func (v *event) Namespace() string {
	return v.namespace
}

func (v *event) GenerateName() string {
	return v.generatename
}

func (v *event) Reason() string {
	return v.reason
}

func (v *event) Message() string {
	return v.message
}

func (v *event) Kind() string {
	return v.kind
}

func (v *event) Component() string {
	return v.component
}
