build:
	@echo "building binary..."
	@go build -o bin/judge
run:build 
	@echo "Ready!"
	@./bin/judge
