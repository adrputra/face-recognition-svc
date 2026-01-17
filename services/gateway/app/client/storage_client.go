package client

import (
	"bytes"
	"context"
	"errors"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type InterfaceStorageClient interface {
	UploadFile(ctx context.Context, req *model.File, bucket string, path string) (string, error)
	UploadFiles(ctx context.Context, req []*model.File, bucket string, path string) error
	StoreFileData(ctx context.Context, tx *gorm.DB, req *model.Dataset) error

	DeleteDatasetDB(ctx context.Context, tx *gorm.DB, username string) error
	DeleteObject(ctx context.Context, bucket string, prefix string) error

	GetDatasetsByUsername(ctx context.Context, bucket string, username string) ([]string, error)
}

type StorageClient struct {
	s3 *s3.S3
	db *gorm.DB
}

func NewStorageClient(s3 *s3.S3, db *gorm.DB) *StorageClient {
	return &StorageClient{
		s3: s3,
		db: db,
	}
}

func (c *StorageClient) UploadFile(ctx context.Context, req *model.File, bucket string, path string) (string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: UploadFile")
	defer span.Finish()

	_, err := c.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("%s.%s", path, req.Extension)),
		Body:   bytes.NewReader(req.BytesObject),
	})

	if err != nil {
		utils.LogEventError(span, err)
		return "", err
	}

	urlStr := fmt.Sprintf("https://minioc.eventarry.com/api/v1/buckets/bpkp/objects/download?preview=true&prefix=%s.%s", path, req.Extension)
	return urlStr, nil
}

func (c *StorageClient) UploadFiles(ctx context.Context, req []*model.File, bucket string, path string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: UploadFiles")
	defer span.Finish()

	for _, file := range req {
		_, err := c.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(fmt.Sprintf("%s/%s", path, file.FileName)),
			Body:   bytes.NewReader(file.BytesObject),
		})

		if err != nil {
			utils.LogEventError(span, err)
			return err
		}
	}

	return nil
}

func (c *StorageClient) DeleteObject(ctx context.Context, bucket string, prefix string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteObject")
	defer span.Finish()

	utils.LogEvent(span, "Request", bucket)

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}

	log.Info().Str("bucket", bucket).Str("prefix", prefix).Msg("Deleting objects under prefix")

	for {
		// Get a batch of objects
		listOutput, err := c.s3.ListObjectsV2(listInput)
		if err != nil {
			utils.LogEventError(span, err)
			return err
		}

		// Create delete requests for each object
		var deleteObjects []*s3.ObjectIdentifier
		for _, object := range listOutput.Contents {
			deleteObjects = append(deleteObjects, &s3.ObjectIdentifier{Key: object.Key})
		}

		if len(deleteObjects) == 0 {
			utils.LogEvent(span, "", "No Objects to delete")
			return model.ThrowError(http.StatusNotFound, errors.New("No Objects to delete"))
		}

		log.Info().Interface("delete_objects", deleteObjects).Msg("Objects to delete")

		// Perform the delete operation
		_, err = c.s3.DeleteObjectsWithContext(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &s3.Delete{
				Objects: deleteObjects},
		})

		if err != nil {
			utils.LogEventError(span, err)
			return err
		}

		fmt.Printf("Deleted objects under prefix %s\n", prefix)

		// Check if there are more objects to delete
		if *listOutput.IsTruncated {
			listInput.ContinuationToken = listOutput.NextContinuationToken
		} else {
			break
		}
	}

	utils.LogEvent(span, "Response", "Success Delete Bucket")

	return nil
}

func (c *StorageClient) StoreFileData(ctx context.Context, tx *gorm.DB, req *model.Dataset) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: StoreFileData")
	defer span.Finish()

	var args []interface{}
	args = append(args, req.Username, req.Bucket, time.Now())

	var result *gorm.DB
	query := "INSERT INTO face_datasets (username, dataset, created_at) VALUES (?, ?, ?)"
	if tx != nil {
		result = tx.Debug().WithContext(ctx).Exec(query, args...)
	} else {
		result = c.db.Debug().WithContext(ctx).Exec(query, args...)
	}

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	return nil
}

func (c *StorageClient) DeleteDatasetDB(ctx context.Context, tx *gorm.DB, username string) error {
	span, ctx := utils.SpanFromContext(ctx, "Client: DeleteDatasetDB")
	defer span.Finish()

	var result *gorm.DB
	query := "DELETE FROM face_datasets WHERE username = ?"

	if tx != nil {
		result = tx.Debug().WithContext(ctx).Exec(query, username)
	} else {
		result = c.db.Debug().WithContext(ctx).Exec(query, username)
	}

	if result.Error != nil {
		utils.LogEventError(span, result.Error)
		return result.Error
	}

	utils.LogEvent(span, "Response", fmt.Sprintf("deleted %d rows", result.RowsAffected))

	return nil
}

func (c *StorageClient) GetDatasetsByUsername(ctx context.Context, bucket string, prefix string) ([]string, error) {
	span, ctx := utils.SpanFromContext(ctx, "Client: GetDatasetByUsername")
	defer span.Finish()

	utils.LogEvent(span, "Request", prefix)

	objectCh, err := c.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		utils.LogEventError(span, err)
		return nil, err
	}

	var res []string

	for _, object := range objectCh.Contents {
		// Generate presigned URL for each object
		req, _ := c.s3.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*object.Key),
		})
		urlStr, err := req.Presign(2 * time.Hour) // URL valid for 15 minutes
		if err != nil {
			log.Error().Str("key", *object.Key).Err(err).Msg("Failed to generate URL")
			continue
		}

		res = append(res, urlStr)
	}

	utils.LogEvent(span, "Response", res)

	return res, nil
}
