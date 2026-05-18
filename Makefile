spot_proto:
	protoc \
		-I=./Protobuf/proto \
		--go-grpc_out=./Protobuf/gen \
		--go-grpc_opt=paths=source_relative \
		--go_out=./Protobuf/gen \
		--go_opt=paths=source_relative \
		--validate_out="lang=go:Protobuf/gen" \
		--validate_opt=paths=source_relative \
		spot_service/spot_service.proto

order_proto:
	protoc \
		-I=./Protobuf/proto \
		--go-grpc_out=./Protobuf/gen \
		--go-grpc_opt=paths=source_relative \
		--go_out=./Protobuf/gen \
		--go_opt=paths=source_relative \
		--validate_out="lang=go:Protobuf/gen" \
		--validate_opt=paths=source_relative \
		order_service/order_service.proto

user_proto:
	protoc \
		-I=./Protobuf/proto \
		--go-grpc_out=./Protobuf/gen \
		--go-grpc_opt=paths=source_relative \
		--go_out=./Protobuf/gen \
		--go_opt=paths=source_relative \
		--validate_out="lang=go:Protobuf/gen" \
		--validate_opt=paths=source_relative \
		user_service/user_service.proto

common_proto:
	protoc \
		-I=./Protobuf/proto \
		--go_out=./Protobuf/gen \
		--go_opt=paths=source_relative \
		common/common.proto

money_proto:
	protoc \
		-I=./Protobuf/proto \
		--go_out=./Protobuf/gen \
		--go_opt=paths=source_relative \
	 money/money.proto

order_migration_up:
	goose -dir ./OrderService/migrations postgres "postgres://postgres:12345@localhost:5432/order_service?sslmode=disable" up

order_migration_down:
	goose -dir ./OrderService/migrations postgres "postgres://postgres:12345@localhost:5432/order_service?sslmode=disable" down

user_migration_up:
	goose -dir ./UserService/migrations postgres "postgres://postgres:12345@localhost:5432/user_service?sslmode=disable" up

user_migration_down:
	goose -dir ./UserService/migrations postgres "postgres://postgres:12345@localhost:5432/user_service?sslmode=disable" down