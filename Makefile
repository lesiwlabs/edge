unexport GOFLAGS

.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen@v0.18.0
	tinygo build \
	    -o ./build/app.wasm \
	    -target wasm \
	    -no-debug \
	    -panic=trap \
	    -gc=leaking \
	    -opt=2 \
	    ./...

.PHONY: deploy
deploy:
	wrangler deploy
