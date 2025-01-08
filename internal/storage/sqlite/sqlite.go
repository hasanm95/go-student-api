package sqlite

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Storage struct {
    db *gorm.DB
}

type Student struct {
    gorm.Model
    Name    string `gorm:"not null"`
    Email   string `gorm:"uniqueIndex;not null"`
    Age     int    `gorm:"not null"`
}

func New(dbPath string) (*Storage, error) {
    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Auto migrate the schema
    err = db.AutoMigrate(&Student{})
    if err != nil {
        return nil, err
    }

    return &Storage{db: db}, nil
}

func (s *Storage) CreateStudent(student *Student) error {
    return s.db.Create(student).Error
}

func (s *Storage) GetStudent(id uint64) (*Student, error) {
    var student Student
    err := s.db.First(&student, id).Error
    return &student, err
}

func (s *Storage) GetStudents() ([]Student, error){
    var students []Student
    err := s.db.Find(&students).Error
    return students, err
}