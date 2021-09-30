#!/bin/sh

TUN="hijack0"

main() {
    ip tuntap add mode tun dev $TUN
    ip addr add 198.18.0.1/15 dev $TUN  # optional
    ip link set dev $TUN up
    ip route replace default dev $TUN

    echo "tunhijack $*"
    exec tunhijack $*
}

main || exit 1