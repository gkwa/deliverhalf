run: deliverhalf
	./deliverhalf

deliverhalf: ./dist/deliverhalf_darwin_amd64_v1/deliverhalf
	cp $< $@

./dist/deliverhalf_darwin_amd64_v1/deliverhalf: main.go cmd/*/*.go
	gofumpt -w $<
	goreleaser build --single-target --snapshot --clean

all:
	goreleaser build --snapshot --clean

clean:
	rm -f deliverhalf
	rm -rf dist
