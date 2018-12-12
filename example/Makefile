. PHONY: gen

conf = $(shell find . -name 'error_code_def.*')
modules = example_module
target := $(foreach n,$(modules),$(n)_gen.go)

gen: target

target: $(conf)
	go generate error_code_def.go
