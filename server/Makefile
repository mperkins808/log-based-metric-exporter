.PHONY: deps app server build docker-run

deps:
	# installing air for go hot reloading
	cd ./server && \
	go install github.com/cosmtrek/air@latest

# development
server:
	lsof -ti :10015 | xargs kill && \
	export PORT=10015 && \
	export KUBE_HOST=http://localhost:8001 && \
	export RULE_DIR=./rules && \
	air -c .air.toml

build:
	@echo "building docker image" && \
	docker build -t log-based-metric-exporter . 

docker-run:
	docker run -p 10015:10015 -v ./rules:/rules -e RULE_DIR=/rules log-based-metric-exporter