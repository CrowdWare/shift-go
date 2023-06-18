#!/bin/bash
version="1.0.0"
export GOOS=android
export CGO_ENABLED=1
export ANDROID_NDK_HOME=/Volumes/LaCie/Android/sdk/ndk/25.2.9519653

# Build the AAR file for all architectures
gomobile bind -androidapi 28 -target=android -tags debug -o libshift-debug-$version.aar github.com/crowdware/shift-go/lib
cp libshift-debug-$version.aar /Users/art/SourceCode/Shift/Android/app/libs
