#!/usr/bin/env bash

set -euo pipefail

METAL_ROOT=$PKG_ROOT/metalgo-$TAG

mkdir -p "$METAL_ROOT"

OK=$(cp ./build/metalgo "$METAL_ROOT")
if [[ $OK -ne 0 ]]; then
  exit "$OK";
fi


echo "Build tgz package..."
cd "$PKG_ROOT"
echo "Tag: $TAG"
tar -czvf "metalgo-linux-$ARCH-$TAG.tar.gz" "metalgo-$TAG"
aws s3 cp "metalgo-linux-$ARCH-$TAG.tar.gz" "s3://$BUCKET/linux/binaries/ubuntu/$RELEASE/$ARCH/"
