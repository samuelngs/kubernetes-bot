package bots

// Message interface
type Message interface {
	Body() interface{}
}

// Metadata interface
type Metadata interface {
	Key() string
	Val() string
}

// field struct
type field struct {
	key, val string
}

// text struct
type text struct {
	s string
}

// Field func
func Field(key, val string) Metadata {
	return &field{key, val}
}

// Text creates text message
func Text(s string) Message {
	return &text{s}
}

// Key to return field key
func (v *field) Key() string {
	return v.key
}

// Val to return field value
func (v *field) Val() string {
	return v.val
}

func (v *text) Body() interface{} {
	return v.s
}
