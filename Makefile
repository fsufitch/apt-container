SOURCES := go.sum go.mod *.go 

build/apt-container: ${SOURCES}
	go build -o build/apt-container *.go