package domain

import (
	"path/filepath"
)

// changed from AudioInputBucket
type AudioContentBucketBase struct {
	Name   string
	Prefix string
}

// changed from AudioInputBucket
type AudioMediaBucket struct {
	Base  AudioContentBucketBase
	Media Media
}

func (a AudioMediaBucket) PrefixBase() string {
	return filepath.Join(a.Base.Prefix, a.Media.GetBibleId(), a.Media.Id)
}

func PrefixOptions64kMp3(bibleId, mediaId string) []string {
	var paths []string
	// basePrefix will be: audio/ENGESV/ENGESVN2DA
	paths = append(paths, filepath.Join("audio", bibleId, mediaId))
	//audio/ENGESV/ENGESVN2DA/B01___21_Matthew_______N2ESVESV.mp3 (old, most common location)

	// if there are no mp3 files in that location, then append -mp3-64 and try again
	paths = append(paths, filepath.Join("audio", bibleId, mediaId+"-mp3-64"))
	//audio/ENGESV/ENGESVN2DA-mp3-64/B23___04_1John_________N2ESVESV.mp3 (newer location)
	return paths
}

func CalculateOutputPrefix(bibleId, mediaId string) string {
	return filepath.Join("audio", mediaId)
}
