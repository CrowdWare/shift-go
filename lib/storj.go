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

const storjAccessToken = "1GW7L5Hab3vR4twJARK4mMuatA2D319NyYboQXnRQU9JcLDj2BEwwtiZ5whRtwDV4KRPvsfV4HcSjq9DutvF2NLr6yMgij6N6debnCzeLEfPZJds2uLtj4PcQHPXUyzqStdxwTAZrMDJX4RQcvdpqAtbRUVxtbrkg7hRCrjgwTFNCAoATvfeeoXacMkUBMSxpNXLfp3NYWk9KjGgbRC9SkFHDurkrHg8aVs1mMs2vRqW2Y1mcHbpzYthWJxfJB1sQP1shfRyCUZxTY4okb5gnZH3tSSyCPSsSkbLh6KSYnVrb2bqRAr1AgvfQVaB"
const bucketName = "shift"

func initStorj(context context.Context) {
	ctx = context
	accessGrant := flag.String("access", storjAccessToken, "access grant from satellite")
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
