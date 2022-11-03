# Golang protobuf FieldMask path generation tool

Despite Google's AIP identifing the use of '*' wildcard in a field mask path, the actual `fieldmaskpb` implementation does not support this. This utility simply expands a field mask path containing a wildcard to the list of paths that the wildcard would specify.

### Simple use case

```go
message := &someProtobufMessage{}

mask, err := fieldmaskpb.New(message, fieldmaskpbexpand.ExpandAsterisk(message)...)
```

### Supporting paths with wildcards
```go
message := &someProtobufMessage{}

mask, err := fieldmaskpb.New(message, fieldmaskpbexpand.ExpandWithWildcards(message, "path.*.with.wildcard")...)
```

### What this doesn't fix

According to https://google.aip.dev/161:

"Field masks may permit the use of the * character on a repeated field or map to indicate the specification of particular sub-fields in the collection"

Unfortunately this utility cannot make fieldmaskpb accept these kinds of paths, since they are required to use the '*' character and cannot be automatically expanded. The utility should return the path back in the expanded correct format if this specification is combined with supported functionality.

EXCEPT IN THE CASE OF:

`message.*.*.xyz`

Where the first asterisk designates a list and the second designates any nested field on that list item. This functionality could be added in the future to correctly expand this to:

```go
[]string{"message.*.one.xyz", "message.*.two.xyz"}
```