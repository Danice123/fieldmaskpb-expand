package fieldmaskpbexpand

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func Test_allFieldsFromDescriptor(t *testing.T) {
	m := testprotos.SimpleMessage{}

	type args struct {
		message protoreflect.MessageDescriptor
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allFieldsFromDescriptor(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allFieldsFromDescriptor() = %v, want %v", got, tt.want)
			}
		})
	}
}
