.PHONY: all clean

DEBUG ?=
ifdef DEBUG
    release :=
else
    release := --release
endif

all:
	RUSTFLAGS="-C link-args=-Wl,-rpath=/usr/local/lib/rats-tls:/usr/local/lib" cargo build $(release)

clean:
	cargo clean
