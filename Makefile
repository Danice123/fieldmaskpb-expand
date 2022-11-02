
proto:
	go mod vendor
	protoc --go_out=vendor protos/*.proto

clean:
	rm -r vendor