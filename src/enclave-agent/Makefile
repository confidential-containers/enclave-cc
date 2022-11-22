.PHONY: all clean

DEBUG ?=
ifdef DEBUG
    release :=
else
    release := --release
endif

rats-tls:
	RUSTFLAGS="-C link-args=-Wl,-rpath=/usr/local/lib/rats-tls:/usr/local/lib" cargo build $(release)

all:
	$(RUST_FLAGS) cargo build $(release)

clean:
	cargo clean && rm -f Cargo.lock
