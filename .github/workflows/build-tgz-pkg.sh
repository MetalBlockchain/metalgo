PKG_ROOT=/tmp
VERSION=$TAG
METAL_ROOT=$PKG_ROOT/metalgo-$VERSION

mkdir -p $METAL_ROOT

OK=`cp ./build/metalgo $METAL_ROOT`
if [[ $OK -ne 0 ]]; then
  exit $OK;
fi


echo "Build tgz package..."
cd $PKG_ROOT
echo "Version: $VERSION"
tar -czvf "metalgo-linux-$ARCH-$VERSION.tar.gz" metalgo-$VERSION
aws s3 cp metalgo-linux-$ARCH-$VERSION.tar.gz s3://$BUCKET/linux/binaries/ubuntu/$RELEASE/$ARCH/
rm -rf $PKG_ROOT/metalgo*
