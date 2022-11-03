package fieldmaskpbexpand

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ExpandAsterisk(m proto.Message) []string {
	messageFields := m.ProtoReflect().Descriptor().Fields()
	fields := make([]string, messageFields.Len())
	for i := 0; i < messageFields.Len(); i++ {
		fields[i] = messageFields.Get(i).TextName()
	}
	return fields
}

type Node interface {
	Expand() []Node
	Name() string
	Path() string
	IsList() bool
}

type BaseNode struct {
	message protoreflect.MessageDescriptor
}

func (ths BaseNode) Expand() []Node {
	nodes := make([]Node, ths.message.Fields().Len())
	for i := 0; i < ths.message.Fields().Len(); i++ {
		nodes[i] = FieldNode{
			field: ths.message.Fields().Get(i),
			path:  ths.message.Fields().Get(i).TextName(),
		}
	}
	return nodes
}

func (ths BaseNode) Name() string {
	return ""
}

func (ths BaseNode) Path() string {
	return ""
}

func (ths BaseNode) IsList() bool {
	return false
}

type FieldNode struct {
	field protoreflect.FieldDescriptor
	path  string
}

func (ths FieldNode) Expand() []Node {
	m := ths.field.Message()
	if m == nil {
		return []Node{}
	}

	nodes := make([]Node, m.Fields().Len())
	path := ths.path
	if ths.field.IsList() {
		path = path + ".*"
	}
	for i := 0; i < m.Fields().Len(); i++ {
		nodes[i] = FieldNode{
			field: m.Fields().Get(i),
			path:  path + "." + m.Fields().Get(i).TextName(),
		}
	}

	return nodes
}

func (ths FieldNode) Name() string {
	return ths.field.TextName()
}

func (ths FieldNode) Path() string {
	return ths.path
}

func (ths FieldNode) IsList() bool {
	return ths.field.IsList()
}

func ExpandWithWildcards(m proto.Message, path string) ([]string, error) {
	nodes := []Node{BaseNode{message: m.ProtoReflect().Descriptor()}}
	for _, part := range strings.Split(path, ".") {
		if part == "*" {
			newNodes := []Node{}
			for _, n := range nodes {
				if n.IsList() {
					newNodes = append(newNodes, n)
				} else {
					newNodes = append(newNodes, n.Expand()...)
				}

			}
			nodes = newNodes
		} else {
			newNodes := []Node{}
			for _, n := range nodes {
				for _, sn := range n.Expand() {
					if sn.Name() == part {
						newNodes = append(newNodes, sn)
					}
				}
			}
			nodes = newNodes
		}
		if len(nodes) == 0 {
			return nil, fmt.Errorf("path part %s is invalid", part)
		}
	}

	paths := make([]string, len(nodes))
	for i, n := range nodes {
		paths[i] = n.Path()
	}
	return paths, nil
}
