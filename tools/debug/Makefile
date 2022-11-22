.PHONY: build
build:
	cargo build --example async-client

.PHONY: deps
deps:
	rustup update stable
	rustup default stable
	rustup component add rustfmt clippy
