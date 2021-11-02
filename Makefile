.PHONY: clean build deploy

build:
	GOOS=linux go build -ldflags="-s -w" -tags lambda -o bin/http lambda/lambda.go

deploy: clean build
	sls deploy --verbose

serverfull:
	env GOOS=darwin go build -o bin/serverfull

run: serverfull
	bin/serverfull

clean:
	rm -rf ./bin
