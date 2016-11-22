package lib

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

const (
	SrcCreations string = "creations"
	SrcPackages  string = "packages"
)

type Storage struct {
	Error    error
	Session  *session.Session
	Source   string
	Username string
	Version  string
}

func NewStorage(src string, username string, v string) *Storage {
	s, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	stor := &Storage{
		nil,
		s,
		src,
		username,
		v,
	}

	return stor
}

func (s *Storage) GetFileContent(title string, filename string) string {
	svc := s3.New(s.Session)

	obj := &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Key: aws.String(getBucketPath(s.Source, s.Username, title, s.Version, filename)),
	}
	out, err := svc.GetObject(obj)
	if err != nil {
		s.Error = err
	}
	bf := new(bytes.Buffer)
	bf.ReadFrom(out.Body)
	return bf.String()
}

func (s *Storage) StoreFile(content string, contentType string, title string, filename string) {
	svc := s3.New(s.Session)

	obj := &s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Body:        bytes.NewReader([]byte(content)),
		Key:         aws.String(getBucketPath(s.Source, s.Username, title, s.Version, filename)),
		ContentType: aws.String(contentType),
	}
	_, s.Error = svc.PutObject(obj)
}

func getBucketPath(src string, username string, title string, version string, filename string) string {
	if version != "" {
		return fmt.Sprintf("%s/%s/%s/%s/%s", src, username, title, version, filename)
	} else {
		return fmt.Sprintf("%s/%s/%s/%s", src, username, title, filename)
	}
}
