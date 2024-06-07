package storage

import (
	"context"

	"go.uber.org/fx"
	"k8s.io/utils/env"

	"github.com/browserkube/browserkube/storage"
)

var Module = fx.Options(
	fx.Provide(
		provideSessionRecordStorage,
		provideSessionArchiveStorage,
	),
)

func provideSessionRecordStorage() (storage.BlobSessionStorage, error) {
	blobURL := env.GetString("BLOB_URL", "")
	return storage.New(context.Background(), blobURL)
}

func provideSessionArchiveStorage() (storage.BlobSessionArchiveStorage, error) {
	blobURL := env.GetString("BLOB_URL_ARCHIVE", "")
	return storage.New(context.Background(), blobURL)
}
