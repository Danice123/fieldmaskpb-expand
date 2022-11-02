# Golang protobuf FieldMask path generation tool

Despite Google's AIP identifing the use of '*' wildcard in a field mask path, the actual `fieldmaskpb` implementation does not support this. This tool simply expands a field mask path containing a wildcard to the list of paths that the wildcard would specify.

### Simple use case

```go
message := &someProtobufMessage{}

mask, err := fieldmaskpb.New(message, fieldmaskpbexpand.Expand(message, "path.with.wildcard.*")...)
```

### Slower expansion but supports intermediate wildcards
```go
message := &someProtobufMessage{}

mask, err := fieldmaskpb.New(message, fieldmaskpbexpand.ExpandWithIntermediateWildcards(message, "path.*.with.wildcard")...)
```