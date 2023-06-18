#!/bin/bash
export GOOS=android
export CGO_ENABLED=1
export ANDROID_NDK_HOME=/Volumes/LaCie/Android/sdk/ndk/25.2.9519653

go run . vars
mv ./lib/crypto_vars.go.temp ./lib/crypto_vars.go

go run . secret
mv ./lib/crypto_vars.go.temp ./lib/crypto_vars.go

# Build the AAR file for all architectures
gomobile bind -androidapi 28 -target=android -o libshift.aar github.com/crowdware/shift-go/lib
cp libshift.aar /Users/art/SourceCode/Shift/Android/app/libs