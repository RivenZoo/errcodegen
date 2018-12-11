. PHONY: build

src := ${shell find . -name '*.go'}
target := ./build/errcodegen

build: ${target}

${target}:	${src}
	go build -o $@ .

test:
	go test -v ./...

clean:
	-rm -f ${target}
