dev:
	air server -c .air.toml

protoc-gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative src/proto/*.proto