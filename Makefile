
run-config:
	go run $$(ls -1 *.go | grep -v _test.go) --config=sample.config.yml
run-inline:
	go run $$(ls -1 *.go | grep -v _test.go) --periodic=false --path=$(FILEPATH)
run-inline-mediafaker:
	go run $$(ls -1 *.go | grep -v _test.go) --periodic=false --path=$(FILEPATH) --algorithm_name="MediafakerTreeWalk" --ignore=".git" --ignore=".circleci" --show_tree_diff=false --gzip=true
run-inline-mediafaker-moni: 
	moni --periodic=false --path=$(FILEPATH) --algorithm_name="MediafakerTreeWalk" --ignore=".git" --ignore=".circleci"