.PHONY: clean
clean:
	@rm -rf generated crds
	@find ./apis -maxdepth 3 -mindepth 3 -name 'zz_generated.*.go' -exec rm -f {} \;

.PHONY: update
update: clean
	hack/update-all.sh
