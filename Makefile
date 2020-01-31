OUTPUT_IMG ?= bifrost-output:latest
BIFROST_IPAM_IMG ?= bifrost-ipam:latest

all: bifrost-output bifrost-ipam

bifrost-output:
	docker build . -t ${OUTPUT_IMG}

bifrost-ipam: bifrost-output
	docker build build/package/bifrost-ipam -t ${BIFROST_IPAM_IMG}

install-ipam:
	kubectl apply -f build/package/bifrost-ipam/bifrost-ipam.yaml

test: bifrost-output


.PHONY: all bifrost_output bifrost-ipam install-ipam test