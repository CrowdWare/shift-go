package lib

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"storj.io/uplink"
)

var access *uplink.Access
var project *uplink.Project
var ctx context.Context

const bucketName = "shift"

func initStorj(context context.Context) error {
	ctx = context
	plain, err := decryptStringGCM(storj_access_token_enc, false)
	if err != nil {
		if debug {
			log.Println("Error decrypt storj access token: " + err.Error())
		}
		return err
	}
	accessGrant := flag.String("access", plain, "access grant from satellite")
	access, err = uplink.ParseAccess(*accessGrant)
	return err
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

func exists(key string) (bool, error) {
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return false, fmt.Errorf("could not open project: %v", err)
	}
	defer project.Close()

	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("could not ensure bucket: %v", err)
	}
	info, err := project.StatObject(ctx, bucketName, key)
	if err != nil {
		return false, fmt.Errorf("could not get info: %v", err)
	}
	log.Println(info)
	return info.System.ContentLength > 0, nil
}
