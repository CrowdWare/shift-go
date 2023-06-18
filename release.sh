#!/bin/bash
version="1.0.0"
go run . vars
mv ./lib/crypto_vars.go.temp ./lib/crypto_vars.go

go run . secret
mv ./lib/crypto_vars.go.temp ./lib/crypto_vars.go

export GOOS=android
export CGO_ENABLED=1
export ANDROID_NDK_HOME=/Volumes/LaCie/Android/sdk/ndk/25.2.9519653
# Build the AAR file for all architectures
gomobile bind -androidapi 28 -target=android -o libshift-release-$version.aar github.com/crowdware/shift-go/lib
cp libshift-release-$version.aar /Users/art/SourceCode/Shift/Android/app/libs