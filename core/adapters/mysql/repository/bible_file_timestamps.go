package repository

import (
	domain "biblebrain-domain/core/domain/bible"
	"biblebrain-domain/core/ports/repository"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository to handle the records of BibleFileWithGaps temp entity
type BibleFileWithGapsRepository struct {
	db *gorm.DB
}

func NewBibleFileWithGapsRepository(db *gorm.DB) repository.BibleFileWithGapsRepository {
	return &BibleFileWithGapsRepository{db: db}
}

// This method creates a temporary table called bible_files_with_timestamp_gaps which holds information about
// Bible files with gaps in their timestamp sequences.
func (r BibleFileWithGapsRepository) CreateTempTable() error {
	sql := `CREATE TEMPORARY TABLE bible_files_with_timestamp_gaps SELECT bfc.bible_id,
		bf2.id as bible_file_id,
		bf2.book_id,
		bf2.chapter_start,
		bf2.file_name as bible_file_name,
		SUM(bft.verse_sequence) AS sum_verses,
		COUNT(bft.verse_sequence) AS count_verses,
		(COUNT(bft.verse_sequence)*(COUNT(bft.verse_sequence)+1))/2 AS formula_sum_of_n
		FROM bible_fileset_connections bfc 
		JOIN bible_filesets bf ON bfc.hash_id = bf.hash_id
		JOIN bible_files bf2 ON bf2.hash_id = bf.hash_id 
		JOIN bible_file_timestamps bft ON bft.bible_file_id = bf2.id
		WHERE bft.verse_sequence <> 0
		GROUP BY bfc.bible_id, bf2.id, bf2.book_id, bf2.file_name,bf2.chapter_start
		HAVING SUM(bft.verse_sequence) <> (COUNT(bft.verse_sequence)*(COUNT(bft.verse_sequence)+1))/2;`

	return r.db.Exec(sql).Error
}

// This method retrieves a list of BibleFileWithGaps records from the database, applying the provided limit and offset to the query.
func (r BibleFileWithGapsRepository) List(limit int, offset int) ([]domain.BibleFileWithGaps, error) {
	var bibleFiles []domain.BibleFileWithGaps
	result := r.db.Limit(limit).Offset(offset).Find(&bibleFiles)
	if result.Error != nil {
		fmt.Println("Error fetching BibleFileWithGaps:", result.Error)
		return nil, result.Error
	}

	return bibleFiles, nil
}

// This method returns the total count of records in the BibleFileWithGaps table.
func (r BibleFileWithGapsRepository) Total() (int64, error) {
	var count int64
	result := r.db.Model(&domain.BibleFileWithGaps{}).Count(&count)
	if result.Error != nil {
		fmt.Println("Error counting BibleFileWithGaps:", result.Error)
		return 0, result.Error
	}

	return count, nil
}

// Repository to handle the records of the BibleFileTimestamps entity
type BibleFileTimestampsRepository struct {
	db *gorm.DB
}

func NewBibleFileTimestampsRepository(db *gorm.DB) repository.BibleFileTimestampsRepository {
	return &BibleFileTimestampsRepository{db: db}
}

// This method updates a single BibleFileTimestamps record in the database, using the provided bibleFileTimestamp object.
func (r *BibleFileTimestampsRepository) Update(bibleFileTimestamp *domain.BibleFileTimestamps) error {
	return r.db.Save(bibleFileTimestamp).Error
}

// This method retrieves a list of BibleFileTimestamps records associated with the provided bibleFileId.
// It only returns records where the verse_start field is not equal to 0.
func (r *BibleFileTimestampsRepository) GetByBibleFileId(bibleFileId int64) ([]domain.BibleFileTimestamps, error) {
	var bibleFileTimestamps []domain.BibleFileTimestamps
	err := r.db.Where(
		"bible_file_id = ?",
		bibleFileId,
	).Where(
		"verse_start <> ?",
		0,
	).Find(&bibleFileTimestamps).Error

	if err != nil {
		return nil, err
	}
	return bibleFileTimestamps, nil
}

// This method updates a batch of BibleFileTimestamps records in the database, specifically updating the verse_end field
// for each record in the provided listToUpdate. This method uses the clause.OnConflict option to handle conflicts during
// batch updates.
func (r BibleFileTimestampsRepository) UpdateBatchForVerseEnd(listToUpdate []domain.BibleFileTimestamps) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Table: domain.BibleFileTimestamps{}.TableName(), Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"verse_end"}),
	}).Create(&listToUpdate).Error
}
