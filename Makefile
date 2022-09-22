
build-macos-app:
	mkdir -p snake.app/Contents/MacOS
	go build -o snake.app/Contents/MacOS/snake cmd/snake.go
	
