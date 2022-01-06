package binding

import "net/http"

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON = jsonBinding{}
	YAML = yamlBinding{}
	XML  = xmlBinding{}
)

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

// BindingBody adds BindBody method to Binding. BindBody is similar with Bind,
// but it reads the body from supplied bytes instead of req.Body.
type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

// BindingUri adds BindUri method to Binding. BindUri is similar with Bind,
// but it read the params.
type BindingUri interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

type StructValidator interface {
	ValidateStruct(interface{}) error
	Engine() interface{}
}

var Validator StructValidator = &defaultValidator{}

func validate(obj interface{}) error {
	if Validator != nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
