. PHONY: build

src := ${shell find . -name '*.go'}
target := ./build/errcodegen

build: ${target}

${target}:	${src}
	go build -o $@ .

test:
	go test -v ./...

release:
	gox -os="linux darwin windows" -arch="amd64" ./...

clean:
	-rm -f ${target}
	-rm errcodegen_*_amd64*
