#
#
export PROBLEMS_PATH=./problems
export BIN_OUTPUT_PATH=./uploaded-files-tmp/
export DOCKER_OUT_BIND_PATH=/home/$(USER)/Desktop/backend-project/JudgeService/uploaded-files-tmp/:/uploaded-files-tmp
export DOCKER_IMAGE_CMDLIMITER=debian:latest

build:
	@echo "building binary..."
	@go build -o bin/judge
run:build 
	@echo "Ready!"
	@./bin/judge
