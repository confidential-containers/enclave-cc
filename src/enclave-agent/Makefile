.PHONY: eaakbc clean cckbc

DEBUG ?=
ifdef DEBUG
    release :=
else
    release := --release
endif

all cc-kbc:
	RUSTFLAGS="-C link-args=-W" cargo build $(release)

eaa-kbc:
	RUSTFLAGS="-C link-args=-Wl,-rpath=/usr/local/lib/rats-tls:/usr/local/lib" cargo build $(release) --no-default-features --features eaa-kbc-rustls-tls

sample-kbc:
	RUSTFLAGS="-C link-args=-W" cargo build $(release) --no-default-features --features simple

clean:
	cargo clean
