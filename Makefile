docker_dev:
	docker build -t go_nvm_test .
	docker run -it go_nvm_test:latest