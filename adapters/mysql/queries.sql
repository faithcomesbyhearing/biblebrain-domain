-- queries.sql

-- name: GetStoragePrefixForFileset :one
select CONCAT(REGEXP_SUBSTR(bfs.set_type_code, 'audio|video|text'), "/", bfc.bible_id, "/", bfs.id, "/")
from bible_fileset_connections bfc 
join bible_filesets bfs on bfs.hash_id = bfc.hash_id
where bfs.id =?;

-- name: GetFilesForBible :many
select REGEXP_SUBSTR(bfs.set_type_code, 'audio|video|text') as type,
bfc.bible_id, bfs.id, bf.file_name
from bible_files bf
join bible_fileset_connections bfc on bfc.hash_id = bf.hash_id
join bible_filesets bfs on bfs.hash_id = bf.hash_id
where bfc.bible_id = ?
order by bfs.asset_id, bfc.bible_id, bfs.id, bf.file_name;

-- name: GetFilesForFileset :many
select REGEXP_SUBSTR(bfs.set_type_code, 'audio|video|text') as type,
bfc.bible_id, bfs.id, bf.file_name
from bible_files bf
join bible_fileset_connections bfc on bfc.hash_id = bf.hash_id
join bible_filesets bfs on bfs.hash_id = bf.hash_id
where bfs.id = ?
order by bfs.asset_id, bfc.bible_id, bfs.id, bf.file_name;

-- name: HlsStitchingSearch :many
select distinct b.id as bible
from bibles b
join bible_fileset_connections bfc on bfc.bible_id = b.id
join bible_filesets bfs on bfs.hash_id = bfc.hash_id
join bible_files bf on bf.hash_id = bfs.hash_id
left join bible_file_stream_bandwidths bfsb on bfsb.bible_file_id = bf.id
left join bible_file_stream_ts bfst on bfst.stream_bandwidth_id = bfsb.id
where bfs.archived is false
and set_type_code like 'video%'
and bfst.id is null 
order by b.id;

-- name: HlsStitchingDrillDown :many
select b.id as bible, bfs.id as fileset, bf.id as bible_file_id, bfsb.id as stream_bandwidth_id, bfsb.resolution_height, bfsb.file_name, bfst.id as stream_ts_id, bfsb.updated_at
from bibles b
join bible_fileset_connections bfc on bfc.bible_id = b.id
join bible_filesets bfs on bfs.hash_id = bfc.hash_id
join bible_files bf on bf.hash_id = bfs.hash_id
join bible_file_stream_bandwidths bfsb on bfsb.bible_file_id = bf.id
left join bible_file_stream_ts bfst on bfst.stream_bandwidth_id = bfsb.id
where b.id = ?
and bfst.id is null
and set_type_code like 'video%'
order by bfsb.resolution_height, bf.book_id, bf.chapter_start, bf.verse_start;
