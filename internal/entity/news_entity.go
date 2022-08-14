package entity

import (
	"fmt"
	"strconv"
	"time"
)

// Error is an error type for domain game.
type Error string

// Error implements error interface.
func (e Error) Error() string {
	return string(e)
}

// NewsModel represents a news model
type NewsModel struct {
	ID              int64     `json:"id" gorm:"column:id"`
	Title           string    `json:"title" gorm:"column:title;index"`
	RssID           int64     `json:"rss_id" gorm:"column:rss_id;index"`
	Cover           string    `json:"cover" gorm:"column:image_url"`
	Link            string    `json:"link" gorm:"column:link"`
	Author          string    `json:"author" gorm:"column:author"`
	CategorySource  string    `json:"category_source" gorm:"column:category_source"`
	Category        string    `json:"subcategory_name" gorm:"column:category"`
	CategoryID      int64     `json:"subcategory_id" gorm:"column:category_id;index"`
	PubDate         string    `json:"pubDate" gorm:"column:pubdate" sql:"DEFAULT:current_timestamp"`
	Permalink       string    `json:"permalink" gorm:"column:permalink"`
	Content         string    `json:"content" gorm:"column:content"`
	Source          string    `json:"source" gorm:"column:source"`
	Tags            string    `json:"tags" gorm:"column:tags"`
	Description     string    `json:"description" gorm:"column:description"`
	Publish         int64     `json:"publish" gorm:"column:publish"`
	IsHeadline      int64     `json:"is_headline" gorm:"column:is_headline"`
	CreateTs        time.Time `json:"-" gorm:"column:create_ts" sql:"DEFAULT:current_timestamp"`
	Count           int64     `json:"count"`
	ShareLink       string    `json:"share_link"`
	GaPartnerId     string    `json:"ga_partner_id"`
	Exclusive       string    `json:"exclusive"`
	GoogleIndex     string    `json:"google_index"`
	CountryId       int64     `json:"country_id" gorm:"column:country_id;index"`
	PublisherId     int64     `json:"publisher_id" gorm:"column:publisher_id;index"`
	CountryName     string    `json:"country_name"`
	TotalLike       int64     `json:"total_like"`
	Image           string    `json:"image"`
	Pinned          string    `json:"pinned"`
	Sorting         int64     `json:"sorting" gorm:"column:sorting"`
	TotalViews      int64     `json:"total_views"`
	MetaTitle       string    `json:"meta_title"`
	MetaDescription string    `json:"meta_description"`
	MetaKeyword     string    `json:"meta_keyword"`
	CreatedAt       string    `json:"created_at" gorm:"column:created_at"`
	PublishDate     time.Time `json:"publish_date" gorm:"column:publish_date"`
	Deeplink        string    `json:"deeplink"`
}

type TotalNews struct {
	Total string `json:"total" gorm:"column:total"`
}

// ErrXXX are known errors for domain news.
const (
	ErrNotFound     = Error("News not found")
	ErrAlreadyExist = Error("News already exist")
)

// Timestamp is a marshal & unmarshall field pubDate.
type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(*t).Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(int64(ts), 0))

	return nil
}
