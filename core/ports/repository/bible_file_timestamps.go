package repository

import (
	domain "github.com/faithcomesbyhearing/biblebrain-domain/core/domain/bible"
)

type BibleFileWithGapsRepository interface {
	CreateTempTable() error
	List(limit int, offset int) ([]domain.BibleFileWithGaps, error)
	Total() (int64, error)
}

type BibleFileTimestampsRepository interface {
	UpdateBatchForVerseEnd(listToUpdate []domain.BibleFileTimestamps) error
	Update(bibleFileTimestamp *domain.BibleFileTimestamps) error
	GetByBibleFileId(bibleFileId int64) ([]domain.BibleFileTimestamps, error)
}
