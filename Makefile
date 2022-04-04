.PHONY: build

build:
	docker build -t ciokan/ipecho:latest .

push: build
	docker push ciokan/ipecho:latest