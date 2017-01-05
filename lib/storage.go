package lib

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

// AWS S3 storage locations
const (
	SrcCreations string = "creations"
	SrcPackages  string = "packages"
)

// Storage is Wooble cloud storage interface
type Storage struct {
	Error   error
	Session *session.Session
	Source  string
	Version string
}

// NewStorage initialized a Storage session
func NewStorage(src string, v string) *Storage {
	s, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	stor := &Storage{
		nil,
		s,
		src,
		v,
	}

	return stor
}

// GetFileContent returns requested file from the cloud
func (s *Storage) GetFileContent(username string, title string, filename string, key string) string {
	svc := s3.New(s.Session)

	obj := &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Key: aws.String(s.getBucketPath(makeID(username, key, title), filename)),
	}
	out, err := svc.GetObject(obj)
	if err != nil {
		s.Error = err
		return ""
	}
	bf := new(bytes.Buffer)
	bf.ReadFrom(out.Body)
	return bf.String()
}

// StoreFile stores the file in the cloud
func (s *Storage) StoreFile(content string, contentType string, username string, title string, filename string, key string) string {
	svc := s3.New(s.Session)

	path := s.getBucketPath(makeID(username, key, title), filename)

	obj := &s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Body:        bytes.NewReader([]byte(content)),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	}
	_, s.Error = svc.PutObject(obj)
	return path
}

func (s *Storage) getBucketPath(key []byte, filename string) string {
	var path string
	switch s.Source {
	case SrcCreations:
		path = fmt.Sprintf("%s/%x/%s/%s", s.Source, key, s.Version, filename)
	case SrcPackages:
		path = fmt.Sprintf("%s/%x/%s", s.Source, key, filename)
	}
	return path
}

func makeID(username string, key string, title string) []byte {
	h := sha1.New()
	h.Write([]byte(username))
	if key != "" {
		h.Write([]byte(key))
	}
	h.Write([]byte(title))

	return h.Sum(nil)
}
