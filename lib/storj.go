package lib

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"

	"storj.io/uplink"
)

var access *uplink.Access
var project *uplink.Project
var ctx context.Context

const storjAccessTokenEnc = "fcb0bb57a74aa4e817dc06481fb40ec5f07c5d04da18c148c0a05ac11babb6d0b5b14ac5908199a5f0d7027ff480e4a4c40947b902e60ddf1531ac1ec30afcc691d3de1c163d7a2e6a0ba92b8c58024104cfe208e567446cb1b8f63501d7c23adfc55c0227b3023252f76092f2d7b936e72d22b9aec38e3c970aded52daae6a413a1fd525a7fc11eb3d26a18e8e6b9c84b008692cec2a87eeb25e1b8bbb82f43284211dc0ce0996b1f38a4db4913287282b68a1aa7dc67be5bc50f28c33d94931ca714db4562c8f7f0b8ae2cd1c88c697f24ca443acbb6794c8f314acc11088703edb9c4bbb15d4e373dcc0e0a5256c4c947155a25c5d8b296abff0be80ab6aef2ba85cc1602c8fd96fb71df6b851cd064e366a68b293bb69593de902ba80a964137bafa6323d5ad5a31140b27775dcf5d87a627bb683e4f163b3d0997dda7cdb8efabcda03f3785498bcd0f65c37ee6581f41048503edbe"
const bucketName = "shift"

func initStorj(context context.Context) {
	ctx = context
	accessGrant := flag.String("access", DecryptStringGCM(storjAccessTokenEnc), "access grant from satellite")
	access, _ = uplink.ParseAccess(*accessGrant)
}

func put(key string, dataToUpload []byte) error {
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}

	upload, err := project.UploadObject(ctx, bucketName, key, nil)
	if err != nil {
		return fmt.Errorf("could not initiate upload: %v", err)
	}

	buf := bytes.NewBuffer(dataToUpload)
	_, err = io.Copy(upload, buf)
	if err != nil {
		_ = upload.Abort()
		return fmt.Errorf("could not upload data: %v", err)
	}

	err = upload.Commit()
	if err != nil {
		return fmt.Errorf("could not commit uploaded object: %v", err)
	}
	return nil
}

func get(key string) ([]byte, error) {
	empty := make([]byte, 0)

	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return empty, fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return empty, fmt.Errorf("could not ensure bucket: %v", err)
	}

	download, err := project.DownloadObject(ctx, bucketName, key, nil)
	if err != nil {
		return empty, fmt.Errorf("could not open object: %v", err)
	}
	defer download.Close()

	receivedContents, err := io.ReadAll(download)
	if err != nil {
		return empty, fmt.Errorf("could not read data: %v", err)
	}
	return receivedContents, nil
}

func delete(key string) error {
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("could not ensure bucket: %v", err)
	}

	_, err = project.DeleteObject(ctx, bucketName, key)
	if err != nil {
		return fmt.Errorf("could not delete object: %v", err)
	}
	return nil
}
