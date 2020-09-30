.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures && ./start.sh && docker-compose  -f docker-compose.yaml -f docker-compose-ca.yaml up  --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures && docker-compose down --volumes --remove-orphans
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./education

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/education-* education
	@cd fixtures && ./stop.sh
	@echo "Clean up done"

