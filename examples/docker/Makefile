IMAGE=centos
VERSION=centos7.1.1503
NAME=${IMAGE}-${VERSION}-ntpd
RELEASE:=$(shell git rev-parse --verify --short HEAD)
USER=myuser
PASS=mypass
MAIL=myuser@example.com
REGISTRY=myregistry
YAML={Registry: '${REGISTRY}', Image: '${IMAGE}', ImageVersion: '${VERSION}'}

all: build

clean:
	docker rmi ${NAME} &>/dev/null || true
	for file in $$(find . -name '*.tf'); do \
		rm -f $${file%%.tf} ; \
	done

template: clean
	for file in $$(find . -name '*.tf'); do \
		tf -f input.yaml -i "${YAML}" -t $${file} -o $${file%%.tf} ; \
	done

build: template
	docker build --pull=true --no-cache -t ${REGISTRY}/${NAME}:${RELEASE} .
	docker tag -f ${REGISTRY}/${NAME}:${RELEASE} ${REGISTRY}/${NAME}:latest

push: build
	docker login -u ${USER} -p '${PASS}' -e ${MAIL} ${REGISTRY}
	docker push ${REGISTRY}/${NAME}:${RELEASE}
	docker push ${REGISTRY}/${NAME}:latest

verify: template
	for f in $$(find . -type f -name '*.tf'); do printf "=== $${f%%.tf} ===\n"; cat $${f%%.tf}; printf "\n"; done | more
