package blogcore

import "github.com/jinzhu/gorm"

type BlogCore struct {
	db *gorm.DB
}

func NewBlogCore(blogDB *gorm.DB) *BlogCore {
	return &BlogCore{
		db: blogDB,
	}
}
func (blogCore *BlogCore) DB() *gorm.DB {
	return blogCore.db
}
