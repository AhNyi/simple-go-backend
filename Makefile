.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf .aws-sam

build:
	sam build 

deploy: build
	sam deploy \
		--resolve-s3 \
		--template-file sam.yml \
		--stack-name simple-go-backend \
		--capabilities CAPABILITY_IAM \
		--profile nk-test
.PHONY: deploy

test:
	go test ./... -coverprofile cover.out -count=1 --cover -tags=test ./ -v -p=1
.PHONY: test 

format:
	go fmt ./...
.PHONY: format