package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/hatappi/go-kit/storage/option"
)

type S3 struct {
	bucketName string
	prefixPath string

	s3Service s3iface.S3API
}

func NewS3(bucketName string, prefixPath string, region string) (*S3, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &S3{
		s3Service:  s3.New(sess, aws.NewConfig().WithRegion(region)),
		bucketName: bucketName,
		prefixPath: prefixPath,
	}, nil
}

func (s *S3) Save(ctx context.Context, filePath string, data []byte, opts ...option.SaveOptionFunc) (string, error) {
	var saveOpt option.SaveOption
	for _, opt := range opts {
		opt(&saveOpt)
	}

	key := s.objectKey(filePath)

	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}
	if saveOpt.ContentType != nil {
		input.ContentType = saveOpt.ContentType
	}
	if saveOpt.ContentDisposition != nil {
		input.SetContentDisposition(*saveOpt.ContentDisposition)
	}

	if _, err := s.s3Service.PutObjectWithContext(ctx, input); err != nil {
		return "", err
	}
	uri := fmt.Sprintf("s3://%s/%s", s.bucketName, key)

	return uri, nil
}

func (s *S3) Get(ctx context.Context, filePath string) ([]byte, error) {
	key := s.objectKey(filePath)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	o, err := s.s3Service.GetObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, nil
			}
		}

		return nil, err
	}
	defer o.Body.Close()

	resBody, err := ioutil.ReadAll(o.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func (s *S3) Delete(ctx context.Context, filePath string) error {
	key := s.objectKey(filePath)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	if _, err := s.s3Service.DeleteObjectWithContext(ctx, input); err != nil {
		return err
	}

	return nil
}

func (s *S3) Ping(ctx context.Context) error {
	filePath := "ping"

	if _, err := s.Save(ctx, filePath, []byte("test")); err != nil {
		return err
	}

	b, err := s.Get(ctx, filePath)
	if err != nil {
		return err
	}

	if b == nil {
		return errors.New("file does not exist")
	}

	if err := s.Delete(ctx, filePath); err != nil {
		return err
	}

	return nil
}

func (s *S3) objectKey(filePath string) string {
	return path.Join(s.prefixPath, filePath)
}

func (s *S3) exist(ctx context.Context, filePath string) (bool, error) {
	key := s.objectKey(filePath)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	_, err := s.s3Service.GetObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}
