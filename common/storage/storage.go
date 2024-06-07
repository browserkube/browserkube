package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob" // gocloud.dev api's imports
	_ "gocloud.dev/blob/gcsblob"  // gocloud.dev api's imports
	_ "gocloud.dev/blob/s3blob"   // gocloud.dev api's imports
)

type Storage interface {
	GetFile(ctx context.Context, sessionID, filename string) (*BlobFile, error)
	ListFileNames(ctx context.Context, sessionID, prefix string) ([]string, error)
	ListPage(ctx context.Context, sessionID, prefix, pageToken string, pageSize int) ([]string, []byte, error)
	SaveFile(ctx context.Context, sessionID, prefix string, sr *BlobFile) error
	DeleteFile(ctx context.Context, sessionID, filename string) error
	Exists(ctx context.Context, sessionID, filename string) (bool, error)

	SizeUsed() (int64, error)

	Close() error
}

type BlobSessionStorage Storage
type BlobSessionArchiveStorage Storage

type blobStorage struct {
	bucket *blob.Bucket
	log    *zap.SugaredLogger
}

// New creates Storage driver for the given url.
// For local filesystem storage use url like: "file:///absolute/path/" (note the trailing slash is required)
// For AWS S3 storage first set the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env vars and
// then pass url: "s3://your_bucket_name/?awssdk=2"
// Also you may explicitly set region: "s3://your_bucket_name/?region=us-west-1&awssdk=2"
func New(ctx context.Context, blobURL string) (Storage, error) {
	bucket, err := blob.OpenBucket(ctx, blobURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bs := &blobStorage{
		bucket: bucket,
		log:    zap.S(),
	}
	accessible, err := bucket.IsAccessible(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !accessible {
		return nil, errors.New("bucket isn't accessible")
	}

	return bs, nil
}

type BlobFile struct {
	FileName    string
	ContentType string
	Content     io.Reader
}

// GetFile returns session record file for given sessionID and filename
func (s *blobStorage) GetFile(ctx context.Context, sessionID, filename string) (*BlobFile, error) {
	r, err := s.bucket.NewReader(ctx, buildPath(sessionID, filename), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer r.Close()

	contentType := r.ContentType()
	if contentType == "" {
		ext := filepath.Ext(filename)
		switch ext {
		case ".log":
			contentType = "text/plain"
		case ".json":
			contentType = "application/json"
		case ".png":
			contentType = "image/png"
		default:
			return nil, fmt.Errorf("unable to recognize content type for extension: %s", ext)
		}
	}

	payload := &bytes.Buffer{}
	if _, err := io.Copy(payload, r); err != nil {
		return nil, errors.WithStack(err)
	}

	return &BlobFile{
		FileName:    filename,
		ContentType: contentType,
		Content:     payload,
	}, nil
}

// Exists checks if there is a file in the storage
func (s *blobStorage) Exists(ctx context.Context, sessionID, filename string) (bool, error) {
	return s.bucket.Exists(ctx, buildPath(sessionID, filename))
}

// ListFileNames returns file names that exist in storage for the given sessionID
// Prefix is package name. In case if package sessionID consider package records
func (s *blobStorage) ListFileNames(ctx context.Context, sessionID, prefix string) ([]string, error) {
	iter := s.bucket.List(&blob.ListOptions{Prefix: filepath.Join(sessionID, prefix)})

	result := make([]string, 0)

	for {
		obj, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		result = append(result, keyToFileName(obj.Key, sessionID))
	}

	return result, nil
}

// ListPage returns a page of ListObject results for blobs in a bucket.
// To fetch the first page, pass "first page" as the pageToken.
// For subsequent pages, pass the pageToken returned from a previous call to ListPage.
// It is not possible to "skip ahead" pages.
func (s *blobStorage) ListPage(ctx context.Context, sessionID, prefix, pageToken string, pageSize int) ([]string, []byte, error) {
	page, nextPageToken, err := s.bucket.ListPage(ctx, []byte(pageToken), pageSize, &blob.ListOptions{Prefix: filepath.Join(sessionID, prefix)})
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	result := make([]string, len(page))

	for i, obj := range page {
		result[i] = keyToFileName(obj.Key, sessionID)
	}

	return result, nextPageToken, nil
}

// SaveFile saves a session record file into the storage
// Prefix is package name in case if you want to create new package in package sessionID
func (s *blobStorage) SaveFile(ctx context.Context, sessionID, prefix string, sr *BlobFile) error {
	path := filepath.Join(sessionID, prefix, sr.FileName)
	writer, err := s.bucket.NewWriter(ctx, path, &blob.WriterOptions{
		ContentType: sr.ContentType,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	defer writer.Close()

	if _, err := io.Copy(writer, sr.Content); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DeleteFile deletes session record file from storage.
func (s *blobStorage) DeleteFile(ctx context.Context, sessionID, filename string) error {
	return errors.WithStack(s.bucket.Delete(ctx, buildPath(sessionID, filename)))
}

func (s *blobStorage) SizeUsed() (int64, error) {
	sizeTotal, err := s.readSizeOfAllFiles(context.Background())
	if err != nil {
		s.log.Errorf("failed to read size of all files in bucket: %s", err)
	}
	return sizeTotal, err
}

// Close releases any used resources.
func (s *blobStorage) Close() error {
	return errors.WithStack(s.bucket.Close())
}

func buildPath(sessionID, filename string) string {
	return sessionID + "/" + filename
}

func keyToFileName(key, sessionID string) string {
	return strings.TrimPrefix(key, sessionID+"/")
}

func (s *blobStorage) readSizeOfAllFiles(ctx context.Context) (int64, error) {
	iter := s.bucket.List(nil)
	var result int64
	for {
		obj, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return 0, errors.WithStack(err)
		}

		result += obj.Size
	}

	return result, nil
}
