#!/bin/bash
version="1.2.0"
export GOOS=android
export CGO_ENABLED=1
export ANDROID_NDK_HOME=/Volumes/LaCie/Android/sdk/ndk/25.2.9519653

# Build the AAR file for all architectures
gomobile bind -androidapi 21 -target=android -tags debug -o libshift-debug-$version.aar github.com/crowdware/shift-go/lib
cp libshift-debug-$version.aar /Users/art/SourceCode/shift/Android/app/libs
cp libshift-debug-$version.aar /Users/art/SourceCode/shift/plugins/PluginSample/app/libs
unzip -jo "libshift-debug-$version.aar" "classes.jar" -d "/Users/art/SourceCode/shift/android/shiftapi/libs"
