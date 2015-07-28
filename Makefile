all: build

clean:
	rm -rf pkg bin

test: clean
	gb test -v

build: test
	gb build all

update:
	gb vendor update --all
