package main

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"strings"
	"time"
)
import "cloud.google.com/go/storage"

var storageClient *storage.Client
var currentBucketName string

func loadGCP() {
	var err error
	ctx := context.Background()
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}

func saveImage(img string, id string) {
	ctx := context.Background()
	if storageClient == nil {
		return
	}
	bkt, err := getBucket()
	if err != nil {
		sentry.CaptureException(err)
		return
	}

	obj := bkt.Object(fmt.Sprintf("%s-%d.png", id, time.Now().Unix()))
	w := obj.NewWriter(ctx)
	if _, err := fmt.Fprintf(w, img); err != nil {
		sentry.CaptureException(err)
	}

	if err := w.Close(); err != nil {
		sentry.CaptureException(err)
	}
}

func getBucket() (*storage.BucketHandle, error){
	ctx := context.Background()
	now := time.Now()
	year := now.Year()
	month := now.Month().String()
	day := now.Day()
	bucketName := strings.ToLower(fmt.Sprintf("%d-%s-%d", year, month, day))

	bkt := storageClient.Bucket(bucketName)
	if bucketName == currentBucketName {
		return bkt, nil // The bucket name is still valid
	}

	_, err := bkt.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		err := bkt.Create(ctx, "leagueai", nil)
		if err != nil {
			sentry.CaptureException(err)
		}
		return bkt, err
	}

	if err != nil {
		return nil, err
	}

	// update the current bucket
	currentBucketName = bucketName
	return bkt, nil
}

