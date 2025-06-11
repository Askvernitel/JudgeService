FROM debian:latest

WORKDIR /app
COPY judge .

ENV PROBLEMS_PATH=./problems
ENV BIN_OUTPUT_PATH=./uploaded-files-tmp/
ENV SERVER_PORT=localhost:4040
ENV GRPC_SERVER_ADDRESS=localhost:50000
ENV DOCKER_OUT_BIND_PATH=/home/$(USER)/Desktop/backend-project/JudgeService/uploaded-files-tmp/:/uploaded-files-tmp
ENV DOCKER_IMAGE_CMDLIMITER=debian:latest




CMD ["./judge"]

