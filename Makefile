.PHONY: proto run-backend run-ai docker-up

proto:
	mkdir -p backend/gen/ai_service
	protoc --go_out=backend/gen/ai_service --go_opt=paths=source_relative \
		--go-grpc_out=backend/gen/ai_service --go-grpc_opt=paths=source_relative \
		--proto_path=protos protos/ai_service.proto

	python3 -m grpc_tools.protoc -Iprotos --python_out=ai_service/gen --grpc_python_out=ai_service/gen protos/ai_service.proto

run-backend:
	cd backend && go run main.go

run-ai:
	cd ai_service && python3 server.py

docker-up:
	docker-compose up --build
