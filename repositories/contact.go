package repositories

import (
	"contact/models"

	"gorm.io/gorm"
)

type ContactRepository interface {
	FindContacts(offset, limit int) ([]models.Contact, int64, error)
	GetContact(ID string) (models.Contact, error)
	CreateContact(contact models.Contact) (models.Contact, error)
	UpdateContact(contact models.Contact) (models.Contact, error)
	DeleteContact(contact models.Contact, ID string) (models.Contact, error)
}

type repository struct {
	db *gorm.DB
}

func RepositoryContact(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindContacts(offset, limit int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	// Menghitung total jumlah kontak
	err := r.db.Model(&models.Contact{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Mengambil daftar kontak dengan paginasi
	err = r.db.Offset(offset).Limit(limit).Find(&contacts).Error
	if err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (r *repository) GetContact(ID string) (models.Contact, error) {
	var contact models.Contact
	err := r.db.First(&contact, ID).Error

	return contact, err
}

func (r *repository) CreateContact(contact models.Contact) (models.Contact, error) {
	err := r.db.Create(&contact).Error
	return contact, err
}

// Write this code
func (r *repository) UpdateContact(contact models.Contact) (models.Contact, error) {
	err := r.db.Save(&contact).Error

	return contact, err
}

// Write this code
func (r *repository) DeleteContact(contact models.Contact, ID string) (models.Contact, error) {
	err := r.db.Delete(&contact).Error

	return contact, err
}
