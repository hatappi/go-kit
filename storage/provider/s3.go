package provider

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
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

func (s *S3) Save(filePath string, data []byte) (string, error) {
	key := s.objectKey(filePath)
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	ctx := context.Background()
	if _, err := s.s3Service.PutObjectWithContext(ctx, input); err != nil {
		return "", err
	}

	return key, nil
}

func (s *S3) Get(filePath string) ([]byte, error) {
	key := s.objectKey(filePath)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	ctx := context.Background()
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

func (s *S3) Ping() error {
	filepath := "ping"

	if _, err := s.Save(filepath, []byte("test")); err != nil {
		return err
	}

	ok, err := s.exist(filepath)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("file does not exist")
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(s.objectKey(filepath)),
	}

	if _, err := s.s3Service.DeleteObject(input); err != nil {
		return err
	}

	return nil
}

func (s *S3) objectKey(filePath string) string {
	return path.Join(s.prefixPath, filePath)
}

func (s *S3) exist(filePath string) (bool, error) {
	key := s.objectKey(filePath)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	ctx := context.Background()
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
