package domain

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

type Image struct {
	ID        int    `gorm:"primary_key;auto_increment" json:"id"`
	ImgUrl    string    `gorm:"size:256;not null;" json:"imgUrl" validate:"required" json:"imgUrl"`
	UserID    int    `gorm:"size:100;not null;" json:"userId" validate:"required"`
}

type BulkImage struct {
	ImageIds []string `gorm:"size:100;not null;" json:"imageIds" validate:"required"`
}

type service struct {
	db *gorm.DB
}

type ImageService interface {
	DeleteImage(imageId int) error
	BulkDeleteImage(imageIds []string, userId int) error
	GetImageByID(id int) (*Image, error)
}

func NewImageService(db *gorm.DB) *service {
	return &service{db}
}

var _ ImageService = &service{}

func (r *service) GetImageByID(id int) (*Image, error) {

	var image = &Image{}
	err := r.db.Debug().Model(Image{}).Where("id = ?", id).Take(&image).Error
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (r *service) DeleteImage(imageId int) error {

	var imgInfo Image
	err := r.db.Debug().Where("id = ?", imageId).Delete(&imgInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *service) BulkDeleteImage(imageIds []string, userId int) error {

	for _, id := range imageIds {

		imgDelete := &Image{}
		img := &Image{}

		err := r.db.Debug().Where("id = ? AND user_id = ?", id, userId).Take(&img).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				return errors.New(fmt.Sprintf("Ensure that Image with the id %s exists and you have permission to delete it", id))
			}
			return err
		}

		if img != nil && img.UserID != userId {
			return errors.New("Cannot delete an images you didnt own.")
		}

		err = r.db.Debug().Where("id = ? AND user_id = ?", id, userId).Delete(&imgDelete).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *service) SeedImages() error {

	images := []Image{
		{
			1,
			"https://res.cloudinary.com/chikodi/image/upload/v1601634127/bart-christiaanse-7QFmFdOpdFs-unsplash.jpg",
			1,
		},
		{
			2,
			"https://res.cloudinary.com/chikodi/image/upload/v1601634106/lina-castaneda-HcmdstM9IFw-unsplash.jpg",
			1,
		},
		{
			3,
			"https://res.cloudinary.com/chikodi/image/upload/v1601634070/jocke-wulcan-NMwgHV1xdHU-unsplash.jpg",
			1,
		},
		{
			4,
			"https://res.cloudinary.com/chikodi/image/upload/v1601634066/bailey-mahon-aK3qEYH_nO0-unsplash.jpg",
			2,
		},
		{
			5,
			"https://res.cloudinary.com/chikodi/image/upload/v1601634057/karolis-puidokas-3Ruy7rRNevY-unsplash.jpg",
			2,
		},
		{
			6,
			"https://res.cloudinary.com/chikodi/image/upload/v1601634049/jake-colling-9O-l0p38gPw-unsplash.jpg",
			2,
		},
		{
			7,
			"https://res.cloudinary.com/chikodi/image/upload/v1601633944/jc-gellidon-5SdGN6k8zpQ-unsplash.jpg",
			2,
		},
		{
			8,
			"https://res.cloudinary.com/chikodi/image/upload/v1601633704/jaber-ahmed-SIRrK_oox2M-unsplash.jpg",
			2,
		},
		{
			9,
			"https://res.cloudinary.com/chikodi/image/upload/v1601690048/alvin-lenin-QA7KTyc3G5Y-unsplash.jpg",
			2,
		},
		{
			10,
			"https://res.cloudinary.com/chikodi/image/upload/v1601690047/andrew-reshetov-JeJIA7SFRRI-unsplash.jpg",
			2,
		},
		{
			11,
			"https://res.cloudinary.com/chikodi/image/upload/v1601690284/jocke-wulcan-3jnMoTB2mD8-unsplash.jpg",
			2,
		},
		{
			12,
			"https://res.cloudinary.com/chikodi/image/upload/v1601690043/billy-cox-EMQlFWKASTg-unsplash.jpg",
			3,
		},
	}

	for i := range images {
		err := r.db.Model(&Image{}).Create(&images[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
