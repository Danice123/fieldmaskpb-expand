package fieldmaskpbexpand

import (
	"reflect"
	"testing"

	"github.com/Danice123/fieldmaskpb-expand/protos"
	"google.golang.org/protobuf/proto"
)

func TestExpand(t *testing.T) {
	type args struct {
		m    proto.Message
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Simple Message",
			args: args{
				m:    &protos.SimpleMessage{},
				path: "*",
			},
			want: []string{"one", "two", "three", "four_nice"},
		},
		{
			name: "Nested Message",
			args: args{
				m:    &protos.NestedMessage{},
				path: "*",
			},
			want: []string{"nested.one", "nested.two", "nested.three", "nested.four_nice", "two"},
		},
		{
			name: "Half of double nested",
			args: args{
				m:    &protos.DoubleNestedMessage{},
				path: "double_one.*",
			},
			want: []string{"double_one.nested.one", "double_one.nested.two", "double_one.nested.three", "double_one.nested.four_nice", "double_one.two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Expand(tt.args.m, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_allFieldsFromDescriptor(t *testing.T) {
	tests := []struct {
		name string
		arg  proto.Message
		want []string
	}{
		{
			name: "Simple Message",
			arg:  &protos.SimpleMessage{},
			want: []string{"one", "two", "three", "four_nice"},
		},
		{
			name: "Nested Message",
			arg:  &protos.NestedMessage{},
			want: []string{"nested.one", "nested.two", "nested.three", "nested.four_nice", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allFieldsFromDescriptor(tt.arg.ProtoReflect().Descriptor()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allFieldsFromDescriptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandWithIntermediateWildcards(t *testing.T) {
	type args struct {
		m    proto.Message
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Simple Message",
			args: args{
				m:    &protos.SimpleMessage{},
				path: "*",
			},
			want: []string{"one", "two", "three", "four_nice"},
		},
		{
			name: "Nested Message",
			args: args{
				m:    &protos.NestedMessage{},
				path: "*",
			},
			want: []string{"nested.one", "nested.two", "nested.three", "nested.four_nice", "two"},
		},
		{
			name: "Half of double nested",
			args: args{
				m:    &protos.DoubleNestedMessage{},
				path: "double_one.*",
			},
			want: []string{"double_one.nested.one", "double_one.nested.two", "double_one.nested.three", "double_one.nested.four_nice", "double_one.two"},
		},
		{
			name: "Middle wildcard",
			args: args{
				m:    &protos.DoubleNestedMessage{},
				path: "*.two",
			}, want: []string{"double_one.two", "double_two.two"},
		},
		{
			name: "Middle wildcard nested",
			args: args{
				m:    &protos.DoubleNestedMessage{},
				path: "*.nested.*",
			}, want: []string{"double_one.nested.one", "double_one.nested.two", "double_one.nested.three", "double_one.nested.four_nice",
				"double_two.nested.one", "double_two.nested.two", "double_two.nested.three", "double_two.nested.four_nice"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandWithIntermediateWildcards(tt.args.m, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expand2() = %v, want %v", got, tt.want)
			}
		})
	}
}
