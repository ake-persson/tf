all: build

clean:
	rm -rf pkg bin

build: clean
	gb build all

update:
	gb vendor update --all
