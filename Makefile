SOURCES := $(shell find . -name '*.go')
TARGET := ./dist/deliverhalf_darwin_amd64_v1/deliverhalf

deliverhalf: $(TARGET)
	cp $< $@


run: deliverhalf
	./deliverhalf

all:
	goreleaser build --snapshot --clean

$(TARGET): $(SOURCES)
	gofumpt -w $<
	goreleaser build --single-target --snapshot --clean

.PHONY: clean
clean:
	rm -f deliverhalf
	rm -f $(TARGET)
	rm -rf dist
