unexport GOFLAGS

.DEFAULT_GOAL := run

.PHONY: build
build: ## builds the executable
	go build

.PHONY: run
run: build download ## runs the conversion
	./ungegn-slo-exonyms

.PHONY: download
download: ## downloads the source
	wget --timestamping http://ungegn.zrc-sazu.si/Portals/7/VELIKA%20PREGLEDNICA_slo.xlsx

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
