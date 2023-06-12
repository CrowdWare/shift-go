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

const storjAccessTokenEnc = "64c0c12683560b96fa83646868d263724781303c8a7c0f9acc7c6a0a4c928994f317d2bb44118755a8c3ea44b01d578a3621de615862c662cdf42378b81eaa3dc7d048eff7cf572baedb6901b4daa6841e4c6c7fa559cb9d08d788acabfc4131ae430b0ecdbcd8e46324f53a1753474096c12361fe9ef26fc8331e0875615b57b0fd23b325828c8107a8a8dd2c0df14d28165980134cac7d064fbf4e74c749849676644d32f99c1c4e6a3232ab8fd90f8cc9cc1f70d986be04aaecd9a0f2da4f3b6791b7612753725668f495065670d2c7591aca3e4d70d74300f95a4f80c697e757c459f204cfd6027714ee6ffeb3edfc021346e2591ec79f21acdb673f799b933f5a68ac0e4bf3c220c5322be745c535aa4d258b9fadaf0edd79f89409a8ba893b8159d47c4ce3509ab8cba6cebe47e52e34cc2d590ddb8616767fe4afed2bae87b6db028472ada8b2feae5325385c30f1af7f4c82b1c5"
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
