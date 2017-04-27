package lib

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AWS S3 storage locations
const (
	SrcCreations string = "creations"
	SrcPackages  string = "packages"
	SrcPreview   string = "previews"
	SrcProfile   string = "profiles"
	SrcCreaThumb string = "crea_thumb"
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

// PushBulkFile prepares multiple files to be processed in the cloud
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

	bucket := s.getBucket()

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
		Bucket: aws.String(s.getBucket()),

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
func (s *Storage) StoreFile(content interface{}, contentType string, userID string, objID string, version string, filename string) string {
	svc := s3.New(s.Session)

	path := s.getFilePath(makeID(userID, objID), version, filename)

	var contentByte []byte
	switch content.(type) {
	case string:
		contentByte = []byte(content.(string))
	case io.Reader:
		var buff bytes.Buffer
		buff.ReadFrom(content.(io.Reader))
		contentByte = buff.Bytes()
	}

	obj := &s3.PutObjectInput{
		Bucket: aws.String(s.getBucket()),

		Body:        bytes.NewReader(contentByte),
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
		Bucket: aws.String(s.getBucket()),
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
		Bucket: aws.String(s.getBucket()),

		Key: aws.String(path),
	}

	if _, err := svc.DeleteObject(obj); err != nil {
		s.errs = append(s.errs, err)
	}
}

// SetSource set storage source
func (s *Storage) SetSource(src string) {
	s.Source = src
}

// GetPathFor returns object path
func (s *Storage) GetPathFor(userID string, objID string, version string, filename string) string {
	return s.getFilePath(makeID(userID, objID), version, filename)
}

func (s *Storage) getBucket() string {
	var bucket string
	if s.Source == SrcPackages {
		bucket = GetPkgRepo()
	} else {
		bucket = GetCloudRepo()
	}
	return bucket
}

func (s *Storage) getFilePath(id []byte, version string, filename string) string {
	var path string
	switch s.Source {
	case SrcCreations:
		path = fmt.Sprintf("%s/%x/%s/%s", s.Source, id, version, filename)
	case SrcPackages:
		path = fmt.Sprintf("%x-%s/%s", id, version, filename)
	case SrcPreview:
		path = fmt.Sprintf("public/%s/%x/%s/%s", s.Source, id, version, filename)
	case SrcProfile:
		// the filename is the id
		filenameSplit := strings.Split(filename, ".")
		ext := filenameSplit[len(filenameSplit)-1]
		path = fmt.Sprintf("public/%s/%x.%s", s.Source, id, ext)
	case SrcCreaThumb:
		filenameSplit := strings.Split(filename, ".")
		ext := filenameSplit[len(filenameSplit)-1]
		path = fmt.Sprintf("public/thumbnails/%x.%s", id, ext)
	}
	return path
}

func makeID(userID string, objID string) []byte {
	h := sha1.New()
	h.Write([]byte(userID))
	h.Write([]byte(objID))
	return h.Sum(nil)
}
