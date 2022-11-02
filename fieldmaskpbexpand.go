package fieldmaskpbexpand

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Expand(m proto.Message, path string) ([]string, error) {
	root := m.ProtoReflect().Descriptor()
	if path == "*" {
		return allFieldsFromDescriptor(root), nil
	}

	prepend := ""
	for _, part := range strings.Split(path, ".") {
		if part == "*" {
			fields := []string{}
			for _, f := range allFieldsFromDescriptor(root) {
				fields = append(fields, prepend+f)
			}
			return fields, nil
		}

		field := root.Fields().ByTextName(part)
		if field == nil {
			return nil, fmt.Errorf("path part %s is invalid", part)
		}

		switch field.Kind() {
		case protoreflect.GroupKind:
			fallthrough
		case protoreflect.MessageKind:
			root = field.Message()
			prepend += field.TextName() + "."
		default:
			return []string{prepend + field.TextName()}, nil
		}

	}
	return nil, errors.New("something went very wrong")
}

func allFieldsFromDescriptor(message protoreflect.MessageDescriptor) []string {
	fields := []string{}
	for i := 0; i < message.Fields().Len(); i++ {
		field := message.Fields().Get(i)
		switch field.Kind() {
		case protoreflect.MessageKind, protoreflect.GroupKind:
			childFields := allFieldsFromDescriptor(field.Message())
			for _, c := range childFields {
				fields = append(fields, field.TextName()+"."+c)
			}
		default:
			fields = append(fields, field.TextName())
		}
	}
	return fields
}

type Node interface {
	Expand() []Node
	Name() string
	Path() string
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
	for i := 0; i < m.Fields().Len(); i++ {
		nodes[i] = FieldNode{
			field: m.Fields().Get(i),
			path:  ths.path + "." + m.Fields().Get(i).TextName(),
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

func ExpandWithIntermediateWildcards(m proto.Message, path string) ([]string, error) {
	nodes := []Node{BaseNode{message: m.ProtoReflect().Descriptor()}}

	lastPath := ""
	for _, part := range strings.Split(path, ".") {
		lastPath = part
		if part == "*" {
			newNodes := []Node{}
			for _, n := range nodes {
				newNodes = append(newNodes, n.Expand()...)
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

	if lastPath == "*" {
		newNodes := []Node{}
		for _, n := range nodes {
			newNodes = append(newNodes, fullyExpand(n)...)
		}
		nodes = newNodes
	}

	paths := make([]string, len(nodes))
	for i, n := range nodes {
		paths[i] = n.Path()
	}
	return paths, nil
}

func fullyExpand(node Node) []Node {
	nodes := []Node{}
	exp := node.Expand()
	if len(exp) == 0 {
		return []Node{node}
	}

	for {
		newExp := []Node{}
		for _, n := range exp {
			exp2 := n.Expand()
			if len(exp2) == 0 {
				nodes = append(nodes, n)
			} else {
				newExp = append(newExp, exp2...)
			}
		}

		if len(newExp) == 0 {
			break
		}
		exp = newExp
	}
	return nodes
}
