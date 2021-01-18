package domain_test

import (
	"gallery/server/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteImage_Success(t *testing.T) {
	conn, err := domain.DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := domain.NewImageService(conn)
	delErr := repo.DeleteImage(1) // seeded image id

	assert.Nil(t, delErr)
}

func TestBulkDeleteImage_Success(t *testing.T) {
	conn, err := domain.DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := domain.NewImageService(conn)
	delErr := repo.BulkDeleteImage([]string{"1","2","3"}, 1) // seeded image id

	assert.Nil(t, delErr)
}
