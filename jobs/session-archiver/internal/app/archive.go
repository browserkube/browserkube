package app

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	clientv1 "github.com/browserkube/browserkube/operator/pkg/client/v1"
	"github.com/browserkube/browserkube/storage"

	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SessionResultArchiver struct {
	SessionResults            clientv1.SessionResultsInterface
	BlobSessionStorage        storage.Storage
	BlobSessionArchiveStorage storage.Storage
	ctx                       context.Context
}

func (a *SessionResultArchiver) Archive() error {
	res, err := a.SessionResults.List(a.ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(res.Items) == 0 {
		return nil
	}

	filesToDelete := make(map[string][]string, len(res.Items))

	writer := new(bytes.Buffer)
	zipW := zip.NewWriter(writer)

	for i := range res.Items {
		archive, err := zipW.Create(fmt.Sprintf("%s/", res.Items[i].Name) + "sessionresult.json")
		if err != nil {
			return err
		}

		err = json.NewEncoder(archive).Encode(res.Items[i])
		if err != nil {
			return err
		}

		files, err := a.BlobSessionStorage.ListFileNames(a.ctx, res.Items[i].Name, "")
		if err != nil {
			return err
		}

		for _, name := range files {
			file, err := a.BlobSessionStorage.GetFile(a.ctx, res.Items[i].Name, name)
			if err != nil {
				return err
			}

			archive, err := zipW.Create(fmt.Sprintf("%s/data/", res.Items[i].Name) + file.FileName)
			if err != nil {
				return err
			}

			_, err = io.Copy(archive, file.Content)
			if err != nil {
				return err
			}
		}

		filesToDelete[res.Items[i].Name] = files
	}

	err = zipW.Close()
	if err != nil {
		return err
	}

	currentTime := time.Now()
	err = a.BlobSessionArchiveStorage.SaveFile(a.ctx, "", "", &storage.BlobFile{
		FileName:    filepath.Join("archive-" + currentTime.Format("2006-01-02")),
		ContentType: "application/zip",
		Content:     writer,
	})
	if err != nil {
		return err
	}

	for sessionName, files := range filesToDelete {
		for _, fileName := range files {
			err = a.BlobSessionStorage.DeleteFile(a.ctx, sessionName, fileName)
			if err != nil {
				return err
			}
		}

		err = a.SessionResults.Delete(a.ctx, sessionName, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
