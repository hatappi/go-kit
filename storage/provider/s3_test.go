package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type mockS3Client struct {
	mockPutObject    func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	mockGetObject    func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	mockDeleteObject func(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func (m *mockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m.mockPutObject(ctx, input, opts...)
}

func (m *mockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.mockGetObject(ctx, input, opts...)
}

func (m *mockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.mockDeleteObject(ctx, input, opts...)
}

func TestS3Save(t *testing.T) {
	type args struct {
		filepath string
		data     []byte
	}

	testCases := []struct {
		name          string
		args          args
		mockPutObject func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
		wantSavedPath string
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				filepath: "foo",
				data:     []byte("test"),
			},
			mockPutObject: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/foo"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body"), cmpopts.IgnoreUnexported(s3.PutObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			wantErr:       false,
			wantSavedPath: "s3://test_bucket/test_prefix/foo",
		},
		{
			name: "fail",
			args: args{
				filepath: "foo",
				data:     []byte("test"),
			},
			mockPutObject: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
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
				s3Client: &mockS3Client{
					mockPutObject: tc.mockPutObject,
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
		name          string
		args          args
		mockGetObject func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		wantBody      []byte
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				filepath: "foo",
			},
			mockGetObject: func(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				expected := &s3.GetObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/foo"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreUnexported(s3.GetObjectInput{})); d != "" {
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
			mockGetObject: func(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
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
				s3Client: &mockS3Client{
					mockGetObject: tc.mockGetObject,
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
		name             string
		args             args
		mockDeleteObject func(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
		wantErr          bool
	}{
		{
			name: "success",
			args: args{
				filepath: "foo",
			},
			mockDeleteObject: func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				expected := &s3.DeleteObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/foo"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreUnexported(s3.DeleteObjectInput{})); d != "" {
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
			mockDeleteObject: func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
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
				s3Client: &mockS3Client{
					mockDeleteObject: tc.mockDeleteObject,
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
		name             string
		mockPutObject    func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
		mockGetObject    func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
		mockDeleteObject func(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
		wantErr          bool
	}{
		{
			name: "success",
			mockPutObject: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body"), cmpopts.IgnoreUnexported(s3.PutObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			mockGetObject: func(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				expected := &s3.GetObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreUnexported(s3.GetObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test"))),
				}, nil
			},
			mockDeleteObject: func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
				expected := &s3.DeleteObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreUnexported(s3.DeleteObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.DeleteObjectOutput{}, nil
			},
			wantErr: false,
		},
		{
			name: "PutObject returns error",
			mockPutObject: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "GetObject returns error",
			mockPutObject: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body"), cmpopts.IgnoreUnexported(s3.PutObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			mockGetObject: func(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return nil, fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "DeleteObject returns error",
			mockPutObject: func(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
				expected := &s3.PutObjectInput{
					Body:   bytes.NewReader([]byte("test")),
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreFields(s3.PutObjectInput{}, "Body"), cmpopts.IgnoreUnexported(s3.PutObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.PutObjectOutput{}, nil
			},
			mockGetObject: func(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				expected := &s3.GetObjectInput{
					Bucket: aws.String("test_bucket"),
					Key:    aws.String("test_prefix/ping"),
				}

				if d := cmp.Diff(*expected, *input, cmpopts.IgnoreUnexported(s3.GetObjectInput{})); d != "" {
					t.Fatalf("unexpected input. %s", d)
				}

				return &s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("test"))),
				}, nil
			},
			mockDeleteObject: func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
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
				s3Client: &mockS3Client{
					mockGetObject:    tc.mockGetObject,
					mockPutObject:    tc.mockPutObject,
					mockDeleteObject: tc.mockDeleteObject,
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
