.PHONY: clean build deploy

build:
	GOOS=linux go build -ldflags="-s -w" -tags lambda -o bin/http lambda/lambda.go

deploy: clean build
	sls deploy --verbose

serverfull:
	env GOOS=darwin go build -o bin/serverfull

run: serverfull
	bin/serverfull

docker-build:
	docker build -t covid-decoder:latest .

docker-run: docker-build
	docker run -d -p 8081:8080 covid-decoder:latest

clean:
	rm -rf ./bin
