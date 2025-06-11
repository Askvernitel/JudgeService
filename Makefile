#
#
export PROBLEMS_PATH=/home/$(USER)/Desktop/backend-project/JudgeService/problems
export BIN_OUTPUT_PATH=./uploaded-files-tmp/
export DOCKER_OUT_BIND_PATH=/home/$(USER)/Desktop/backend-project/JudgeService/uploaded-files-tmp/:/uploaded-files-tmp
export DOCKER_IMAGE_CMDLIMITER=debian:latest
export SERVER_PORT=localhost:4040
export GRPC_SERVER_ADDRESS=localhost:50000
build:
	@echo "building binary..."
	@go build -o bin/judge
run:build 
	@echo "Ready!"
	@./bin/judge
