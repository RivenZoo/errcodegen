. PHONY: gen

conf = $(shell find . -name 'error_code_def.*')
modules = example_module
target := $(foreach n,$(modules),$(n)_gen.go)

gen: $(target)

$(target): $(conf)
	errcodegen --err_func=Errorf --err_func_pkg=fmt --pkg=example --err_def=error_code_def.conf

use_image:
	docker run --rm -it -v $(shell pwd):/data/app rivenzoe/errcodegen errcodegen --err_def=error_code_def.conf --err_func=Errorf --err_func_pkg=fmt --pkg=example

clean:
	-rm $(target) 
	-rm *_gen.go
