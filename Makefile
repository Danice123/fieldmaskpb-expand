
proto:
	protoc --go_out=. protos/*.proto
	mv github.com/Danice123/fieldmaskpb-expand/protos/* protos
	rm -r github.com