#!/bin/bash

KERNEL="kernel/arch/x86_64/boot/bzImage"
INITRD="build/output/initrd.img"

if [ ! -f "$KERNEL" ]; then
    echo "Error: Kernel image not found at $KERNEL"
    exit 1
fi

if [ ! -f "$INITRD" ]; then
    echo "Error: Initrd not found at $INITRD"
    exit 1
fi

qemu-system-x86_64 \
    -kernel "$KERNEL" \
    -initrd "$INITRD" \
    -append "console=ttyS0 root=/dev/ram0 init=/init" \
    -nographic \
    -m 2G \
    -smp 2
