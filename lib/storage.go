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
	Session *session.Session
	Source  string

	errs        []error
	bulkObjects []*s3.ObjectIdentifier
}

func (s *Storage) Error() error {
	if len(s.errs) > 0 {
		return s.errs[0]
	}
	return nil
}

// NewStorage initialized a Storage session
func NewStorage(src string) *Storage {
	s, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	stor := &Storage{
		s,
		src,
		make([]error, 0),
		make([]*s3.ObjectIdentifier, 0),
	}

	return stor
}

// PushBulkObject push an object to process in the cloud
func (s *Storage) PushBulkFile(userID string, objID string, version string, filename string) {
	path := s.getFilePath(makeID(userID, objID), version, filename)
	obj := &s3.ObjectIdentifier{
		Key: aws.String(path),
	}
	s.bulkObjects = append(s.bulkObjects, obj)
}

// CopyAndStoreFile copy and store cloud object
func (s *Storage) CopyAndStoreFile(userID string, objID string, prevVersion string, version string, filename string) {
	svc := s3.New(s.Session)

	path := s.getFilePath(makeID(userID, objID), prevVersion, filename)
	newPath := s.getFilePath(makeID(userID, objID), version, filename)

	bucket := viper.GetString("cloud_repo")

	obj := &s3.CopyObjectInput{
		Bucket: aws.String(bucket),

		Key:        aws.String(newPath),
		CopySource: aws.String(bucket + "/" + path),
	}

	if _, err := svc.CopyObject(obj); err != nil {
		s.errs = append(s.errs, err)
	}
}

// GetFileContent returns requested file from the cloud
func (s *Storage) GetFileContent(userID string, objID string, version string, filename string) string {
	svc := s3.New(s.Session)

	path := s.getFilePath(makeID(userID, objID), version, filename)

	obj := &s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Key: aws.String(path),
	}
	out, err := svc.GetObject(obj)
	if err != nil {
		s.errs = append(s.errs, err)
		return ""
	}
	bf := new(bytes.Buffer)
	bf.ReadFrom(out.Body)
	return bf.String()
}

// StoreFile stores the file in the cloud
func (s *Storage) StoreFile(content string, contentType string, userID string, objID string, version string, filename string) string {
	svc := s3.New(s.Session)

	path := s.getFilePath(makeID(userID, objID), version, filename)

	obj := &s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Body:        bytes.NewReader([]byte(content)),
		Key:         aws.String(path),
		ContentType: aws.String(contentType),
	}
	if _, err := svc.PutObject(obj); err != nil {
		s.errs = append(s.errs, err)
	}
	return path
}

// BulkDeleteFiles delete pushed objects
func (s *Storage) BulkDeleteFiles() {
	svc := s3.New(s.Session)

	params := &s3.DeleteObjectsInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),
		Delete: &s3.Delete{
			Objects: s.bulkObjects,
		},
	}

	_, err := svc.DeleteObjects(params)
	s.errs = append(s.errs, err)
}

// DeleteFile delete a file from the cloud
func (s *Storage) DeleteFile(userID string, objID string, version string, filename string) {
	svc := s3.New(s.Session)

	path := s.getFilePath(makeID(userID, objID), version, filename)

	obj := &s3.DeleteObjectInput{
		Bucket: aws.String(viper.GetString("cloud_repo")),

		Key: aws.String(path),
	}

	if _, err := svc.DeleteObject(obj); err != nil {
		s.errs = append(s.errs, err)
	}
}

func (s *Storage) getFilePath(id []byte, version string, filename string) string {
	var path string
	switch s.Source {
	case SrcCreations:
		path = fmt.Sprintf("%s/%x/%s/%s", s.Source, id, version, filename)
	case SrcPackages:
		path = fmt.Sprintf("%s/%x/%s", s.Source, id, filename)
	}
	return path
}

func makeID(userID string, objID string) []byte {
	h := sha1.New()
	h.Write([]byte(userID))
	h.Write([]byte(objID))
	return h.Sum(nil)
}
