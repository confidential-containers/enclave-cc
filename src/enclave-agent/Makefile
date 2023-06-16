.PHONY: eaakbc clean cckbc

DEBUG ?=
ifdef DEBUG
    release :=
else
    release := --release
endif

all cc-kbc:
	RUSTFLAGS="-C link-args=-W" cargo build $(release)

sample-kbc:
	RUSTFLAGS="-C link-args=-W" cargo build $(release) --no-default-features --features simple

clean:
	cargo clean
