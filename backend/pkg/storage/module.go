package storage

import (
	"context"

	"go.uber.org/fx"
	"k8s.io/utils/env"
)

var Module = fx.Options(
	fx.Provide(
		provideSessionRecordStorage,
		provideSessionArchiveStorage,
	),
)

func provideSessionRecordStorage() (BlobSessionStorage, error) {
	blobURL := env.GetString("BLOB_URL", "")
	return New(context.Background(), blobURL)
}

func provideSessionArchiveStorage() (BlobSessionArchiveStorage, error) {
	blobURL := env.GetString("BLOB_URL_ARCHIVE", "")
	return New(context.Background(), blobURL)
}
