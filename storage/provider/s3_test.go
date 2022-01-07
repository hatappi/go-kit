package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type mockS3Client struct {
	s3iface.S3API

	mockPutObjectWithContext    func(aws.Context, *s3.PutObjectInput, ...request.Option) (*s3.PutObjectOutput, error)
	mockGetObjectWithContext    func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error)
	mockDeleteObjectWithContext func(aws.Context, *s3.DeleteObjectInput, ...request.Option) (*s3.DeleteObjectOutput, error)
}

func (m *mockS3Client) PutObjectWithContext(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
	return m.mockPutObjectWithContext(ctx, input, opts...)
}

func (m *mockS3Client) GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
	return m.mockGetObjectWithContext(ctx, input, opts...)
}

func (m *mockS3Client) DeleteObjectWithContext(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
	return m.mockDeleteObjectWithContext(ctx, input, opts...)
}

func TestS3Save(t *testing.T) {
	type args struct {
		filepath string
		data     []byte
	}

	testCases := []struct {
		name                     string
		args                     args
		mockPutObjectWithContext func(aws.Context, *s3.PutObjectInput, ...request.Option) (*s3.PutObjectOutput, error)
		wantSavedPath            string
		wantErr                  bool
	}{
		{
			name: "success",
			args: args{
				filepath: "foo",
				data:     []byte("test"),
			},
			mockPutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/foo"),
				}

				opt := cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body")
				if d := cmp.Diff(*expected, *input, opt); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			wantErr:       false,
			wantSavedPath: "test_prefix/foo",
		},
		{
			name: "fail",
			args: args{
				filepath: "foo",
				data:     []byte("test"),
			},
			mockPutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, fmt.Errorf("error")
			},
			wantErr:       true,
			wantSavedPath: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			s3Provider := &S3{
				bucketName: "test_bucket",
				prefixPath: "test_prefix",
				s3Service: &mockS3Client{
					mockPutObjectWithContext: tc.mockPutObjectWithContext,
				},
			}

			ctx := context.Background()
			savedPath, err := s3Provider.Save(ctx, tc.args.filepath, tc.args.data)
			if (err != nil) != tc.wantErr {
				t.Errorf("err: %v", err)
			}

			if savedPath != tc.wantSavedPath {
				t.Errorf("savedPath was a mismatch. expected: %s, actual: %s", tc.wantSavedPath, savedPath)
			}
		})
	}
}

func TestS3Get(t *testing.T) {
	type args struct {
		filepath string
	}

	testCases := []struct {
		name                     string
		args                     args
		mockGetObjectWithContext func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error)
		wantBody                 []byte
		wantErr                  bool
	}{
		{
			name: "success",
			args: args{
				filepath: "foo",
			},
			mockGetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				expected := &s3.GetObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/foo"),
				}

				if d := cmp.Diff(*expected, *input); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test"))),
				}, nil
			},
			wantErr:  false,
			wantBody: []byte("test"),
		},
		{
			name: "fail",
			args: args{
				filepath: "foo",
			},
			mockGetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				return nil, fmt.Errorf("error")
			},
			wantErr:  true,
			wantBody: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			s3Provider := &S3{
				bucketName: "test_bucket",
				prefixPath: "test_prefix",
				s3Service: &mockS3Client{
					mockGetObjectWithContext: tc.mockGetObjectWithContext,
				},
			}

			ctx := context.Background()
			body, err := s3Provider.Get(ctx, tc.args.filepath)
			if (err != nil) != tc.wantErr {
				t.Errorf("err: %v", err)
			}

			if string(body) != string(tc.wantBody) {
				t.Errorf("body was a mismatch. expected: %s, actual: %s", tc.wantBody, body)
			}
		})
	}
}

func TestS3Delete(t *testing.T) {
	type args struct {
		filepath string
	}

	testCases := []struct {
		name                        string
		args                        args
		mockDeleteObjectWithContext func(aws.Context, *s3.DeleteObjectInput, ...request.Option) (*s3.DeleteObjectOutput, error)
		wantErr                     bool
	}{
		{
			name: "success",
			args: args{
				filepath: "foo",
			},
			mockDeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				expected := &s3.DeleteObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/foo"),
				}

				if d := cmp.Diff(*expected, *input); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.DeleteObjectOutput{}, nil
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				filepath: "foo",
			},
			mockDeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				return nil, fmt.Errorf("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			s3Provider := &S3{
				bucketName: "test_bucket",
				prefixPath: "test_prefix",
				s3Service: &mockS3Client{
					mockDeleteObjectWithContext: tc.mockDeleteObjectWithContext,
				},
			}

			ctx := context.Background()
			err := s3Provider.Delete(ctx, tc.args.filepath)
			if (err != nil) != tc.wantErr {
				t.Errorf("err: %v", err)
			}
		})
	}
}

func TestS3Ping(t *testing.T) {
	testCases := []struct {
		name                        string
		mockPutObjectWithContext    func(aws.Context, *s3.PutObjectInput, ...request.Option) (*s3.PutObjectOutput, error)
		mockGetObjectWithContext    func(aws.Context, *s3.GetObjectInput, ...request.Option) (*s3.GetObjectOutput, error)
		mockDeleteObjectWithContext func(aws.Context, *s3.DeleteObjectInput, ...request.Option) (*s3.DeleteObjectOutput, error)
		wantErr                     bool
	}{
		{
			name: "success",
			mockPutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				opt := cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body")
				if d := cmp.Diff(*expected, *input, opt); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			mockGetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				expected := &s3.GetObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test"))),
				}, nil
			},
			mockDeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				expected := &s3.DeleteObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.DeleteObjectOutput{}, nil
			},
			wantErr: false,
		},
		{
			name: "PutObject returns error",
			mockPutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "GetObject returns error",
			mockPutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				opt := cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body")
				if d := cmp.Diff(*expected, *input, opt); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			mockGetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				return nil, fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "DeleteObject returns error",
			mockPutObjectWithContext: func(ctx aws.Context, input *s3.PutObjectInput, opts ...request.Option) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				opt := cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body")
				if d := cmp.Diff(*expected, *input, opt); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			mockGetObjectWithContext: func(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
				expected := &s3.GetObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test"))),
				}, nil
			},
			mockDeleteObjectWithContext: func(ctx aws.Context, input *s3.DeleteObjectInput, opts ...request.Option) (*s3.DeleteObjectOutput, error) {
				return nil, fmt.Errorf("error")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			s3Provider := &S3{
				bucketName: "test_bucket",
				prefixPath: "test_prefix",
				s3Service: &mockS3Client{
					mockGetObjectWithContext:    tc.mockGetObjectWithContext,
					mockPutObjectWithContext:    tc.mockPutObjectWithContext,
					mockDeleteObjectWithContext: tc.mockDeleteObjectWithContext,
				},
			}

			ctx := context.Background()
			err := s3Provider.Ping(ctx)
			if (err != nil) != tc.wantErr {
				t.Errorf("err: %v", err)
			}
		})
	}
}
