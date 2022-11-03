package fieldmaskpbexpand

import (
	"reflect"
	"testing"

	"github.com/Danice123/fieldmaskpb-expand/protos"
	"google.golang.org/protobuf/proto"
)

func TestExpandAsterisk(t *testing.T) {
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
			want: []string{"nested", "two"},
		},
		{
			name: "One of",
			arg:  &protos.MessageOneOf{},
			want: []string{"s", "i", "m"},
		},
		{
			name: "Repeated",
			arg:  &protos.ListMessage{},
			want: []string{"list"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandAsterisk(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpandAsterisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandWithWildcards(t *testing.T) {
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
			want: []string{"nested", "two"},
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
				path: "*.nested",
			}, want: []string{"double_one.nested", "double_two.nested"},
		},
		{
			name: "Repeated wildcard nested",
			args: args{
				m:    &protos.ListMessage{},
				path: "list.*.one",
			}, want: []string{"list.*.one"},
		},
		{
			name: "Repeated double wildcard nesting",
			args: args{
				m:    &protos.ComplexListMessage{},
				path: "list.*.double_one.*.two",
			}, want: []string{"list.*.double_one.nested.two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandWithWildcards(tt.args.m, tt.args.path)
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
