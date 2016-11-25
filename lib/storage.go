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

const (
	SrcCreations string = "creations"
	SrcPackages  string = "packages"
)

type Storage struct {
	Error   error
	Session *session.Session
	Source  string
	Version string
}

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

func (s *Storage) GetFileContent(username string, title string, filename string, key string) string {
	svc := s3.New(s.Session)

	obj := &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Key: aws.String(s.getBucketPath(makeId(username, key, title), filename)),
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

func (s *Storage) StoreFile(content string, contentType string, username string, title string, filename string, key string) {
	svc := s3.New(s.Session)

	obj := &s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Body:        bytes.NewReader([]byte(content)),
		Key:         aws.String(s.getBucketPath(makeId(username, key, title), filename)),
		ContentType: aws.String(contentType),
	}
	_, s.Error = svc.PutObject(obj)
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

func makeId(username string, key string, title string) []byte {
	h := sha1.New()
	h.Write([]byte(username))
	if key != "" {
		h.Write([]byte(key))
	}
	h.Write([]byte(title))

	return h.Sum(nil)
}
