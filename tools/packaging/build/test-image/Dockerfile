FROM ubuntu:20.04 as builder

RUN apt-get update && \
    env DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    build-essential

COPY hello.c .

RUN gcc hello.c -fPIE -pie -o hello_world

FROM scratch

LABEL enc_key="LieOhvkqFcGMzZrVzt6vPWlj/F/bgYMNe45vhQpdxAA="
LABEL key_id="kbs:///default/image-kek/11032d96-dccd-46a3-9244-93644d76745f"

COPY --from=builder /hello_world /bin/
