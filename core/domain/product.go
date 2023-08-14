package domain

import (
	"database/sql"
)

type Product struct {
	MediaRecords []MediaRecord `json:"media_records"`
}

type MediaRecord struct {
	Stocknumber string  `json:"stocknumber"`
	Language    string  `json:"language"`
	Version     string  `json:"version"`
	Bibleid     string  `json:"bible_id"`
	Mediaid     string  `json:"media_id"`
	Description string  `json:"description"`
	Codec       *string `json:"codec"`
	Bitrate     *string `json:"bitrate"`
}

func GetMediaRecordsByProductCodeQuery(db *sql.DB, productCode string) (*sql.Rows, error) {
	return db.Query(
		`SELECT DISTINCT replace(bft.description, '/', '') as stocknumber,
				l.name AS language,
				bt.name AS version,
				b.id AS bibleid,
				bf.id AS mediaid,
				bft.description,
				bft_codec.description as codec,
				bft_bitrate.description as bitrate
		FROM bible_filesets bf
		JOIN bible_fileset_connections bfc on bfc.hash_id = bf.hash_id
		JOIN bible_fileset_tags bft on bft.hash_id = bf.hash_id
		JOIN bibles b on b.id = bfc.bible_id
		JOIN bible_translations bt on bt.bible_id = b.id
		JOIN languages l on l.id = b.language_id
		LEFT JOIN bible_fileset_tags bft_codec on bft_codec.hash_id = bf.hash_id AND bft_codec.name = 'codec' 
		LEFT JOIN bible_fileset_tags bft_bitrate on bft_bitrate.hash_id = bf.hash_id AND bft_bitrate.name = 'bitrate' 
		WHERE bt.language_id = 6414
		AND bft.name = 'stock_no'
		AND replace(bft.description, '/', '') = ?`,
		productCode,
	)
}
