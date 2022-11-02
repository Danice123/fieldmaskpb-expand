package fieldmaskpbexpand

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Expand(m proto.Message, path string) ([]string, error) {
	root := m.ProtoReflect().Descriptor()
	if path == "*" {
		return allFieldsFromDescriptor(root), nil
	}

	// for _, part := range strings.Split(path, ".") {
	// 	if part == "*" {
	// 		return allFieldsFromDescriptor(root), nil
	// 	} else {
	// 		field := root.Fields().ByJSONName(part)
	// 		root = field.Message()
	// 	}
	// }

	return nil, nil
}

func allFieldsFromDescriptor(message protoreflect.MessageDescriptor) []string {
	fields := make([]string, message.Fields().Len())
	for i := 0; i < message.Fields().Len(); i++ {
		fields[i] = message.Fields().Get(i).JSONName()
	}
	return fields
}
