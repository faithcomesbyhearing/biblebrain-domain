package bible

import (
	"time"
)

type BibleFileTimestamps struct {
	ID            uint `gorm:"primary_key"`
	BibleFileId   uint
	VerseStart    int
	VerseEnd      string
	VerseSequence uint
	Timestamp     float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (BibleFileTimestamps) TableName() string {
	return "bible_file_timestamps"
}

// BibleFileWithGaps is a temp entity which holds information about
// Bible files with gaps in their timestamp sequences.
type BibleFileWithGaps struct {
	BibleFileId   int64 `gorm:"primary_key"`
	BookID        string
	ChapterStart  int
	BibleFileName string
	SumVerses     int
	CountVerses   int
	// FormulaSumOfN int64
}

func (BibleFileWithGaps) TableName() string {
	return "bible_files_with_timestamp_gaps"
}
