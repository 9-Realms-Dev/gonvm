docker_dev:
	docker build -t go_nvm_dev .
	docker run -it --rm go_nvm_dev:latest

docker_test:
	docker build -f docker/test.dockerfile -t go_nvm_test .
	docker run --rm go_nvm_test:latest