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