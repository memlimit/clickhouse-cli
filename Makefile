### protobuf: Generate protobuf files

protobuf:
	@echo "> Generating protobuf"
	@cd $(shell pwd)/pkg/clickhouse/grpc && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative clickhouse.proto
	@echo "> Done"

config:
	@echo "> Creating sample config"
	@cd $(shell pwd)/docs && cp .clickhouse-cli-config.yaml $(HOME)/.clickhouse-cli-config.yaml
	@echo "> Done"
	