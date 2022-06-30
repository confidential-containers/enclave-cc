.PHONY: all clean

DEBUG ?=
ifdef DEBUG
    release :=
else
    release := --release
endif

all:
	$(RUST_FLAGS) cargo build $(release)

clean:
	cargo clean && rm -f Cargo.lock
