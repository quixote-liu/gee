package binding

import (
	"errors"
	"testing"
)

func TestSliceValidateError(t *testing.T) {
	tests := []struct {
		name string
		err  sliceValidateError
		want string
	}{
		{"has nil elements", sliceValidateError{errors.New("test error"), nil}, "[0]: test error"},
		{"has zero elements", sliceValidateError{}, ""},
		{"has one elements", sliceValidateError{errors.New("one error")}, "[0]: one error"},
		{"has two elements",
			sliceValidateError{
				errors.New("first error"),
				errors.New("second error"),
			},
			"[0]: first error\n[1]: second error",
		},
		{"has many elements",
			sliceValidateError{
				errors.New("first error"),
				errors.New("second error"),
				errors.New("third error"),
				nil,
				nil,
				errors.New("sixth error"),
				nil,
				errors.New("eighth error"),
			},
			"[0]: first error\n[1]: second error\n[2]: third error\n[5]: sixth error\n[7]: eighth error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("sliceValidateError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultValidator(t *testing.T) {
	type exampleStruct struct {
		A string `binding:"max=8"`
		B int    `binding:"gt=0"`
	}
	tests := []struct {
		name    string
		v       *defaultValidator
		obj     interface{}
		wantErr bool
	}{
		{"validate nil obj", &defaultValidator{}, nil, false},
		{"validate int obj", &defaultValidator{}, 3, false},
		{"validate struct failed-1", &defaultValidator{}, exampleStruct{A: "123456789", B: 1}, true},
		{"validate struct failed-2", &defaultValidator{}, exampleStruct{A: "123", B: -1}, true},
		{"validate struct pass", &defaultValidator{}, exampleStruct{A: "12345678", B: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.ValidateStruct(tt.obj); (err != nil) != tt.wantErr {
				t.Errorf("defaultValidator.Validate() error= %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
