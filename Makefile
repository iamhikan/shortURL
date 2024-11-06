CURR_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
UID = $(shell id -u)
GID = $(shell id -g)
-include ${CURR_DIR}/.env

# Gen

.PHONY: gen-mocks
gen-mocks:
	docker run --rm -u ${UID}:${GID} -v "${CURR_DIR}":/src -w /src vektra/mockery:v2.44.2