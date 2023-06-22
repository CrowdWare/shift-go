package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"storj.io/uplink"
)

func put(key string, dataToUpload []byte, bucketName string, ctx context.Context, access *uplink.Access) error {
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

func get(key string, bucketName string, ctx context.Context, access *uplink.Access) ([]byte, error) {
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

func delete(key string, bucketName string, ctx context.Context, access *uplink.Access) error {
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

func exists(key string, bucketName string, ctx context.Context, access *uplink.Access) (bool, error) {
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

func listObjects(bucketName string, path string, ctx context.Context, access *uplink.Access) ([]string, error) {
	project, err := uplink.OpenProject(ctx, access)
	if err != nil {
		return nil, err
	}

	defer project.Close()

	// Ensure the bucket exists
	_, err = project.EnsureBucket(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	// List objects in the specified path
	objects := project.ListObjects(ctx, bucketName, &uplink.ListObjectsOptions{
		Prefix:    path,
		Recursive: false,
	})

	var itemKeys []string
	for objects.Next() {
		item := objects.Item()
		if !item.IsPrefix {
			itemKeys = append(itemKeys, item.Key)
		}
	}

	if err := objects.Err(); err != nil {
		return nil, err
	}

	return itemKeys, nil
}
