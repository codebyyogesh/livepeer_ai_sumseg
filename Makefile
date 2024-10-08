CC=go
GO_BUILD_FLAGS=-trimpath
PROJECT_NAME=lp_ai_sumseg
build:
	$(CC) build -o bin/$(PROJECT_NAME) $(GO_BUILD_FLAGS) .

run: build
	 ./bin/$(PROJECT_NAME)

clean:
	rm -rf bin/$(PROJECT_NAME)