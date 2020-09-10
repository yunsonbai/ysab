#!/bin/bash
# 该脚本打包使用

platform=$1

verssion=`grep -rnw 'VERSION' ./conf/version.go | awk -F '"' '{print $2}'`
app=ysab-$platform-$verssion
package=$app.tgz

echo ""
echo "you will build $app"

sleep 2


case $platform in
    linux)
        env GOOS=linux go build -o $app
        ;;
    mac)
        go build -o $app
        ;;
    *)
      echo "请使用:"
      echo "    ./build.sh linux"
      echo "    or"
      echo "    ./build.sh mac"
      echo "  "
      exit 1
      ;;
esac

echo "you will get $package"

tar -zcvf $package $app
rm $app

echo "success"
echo ""
exit 0