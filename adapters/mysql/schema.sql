CREATE TABLE `bibles` (
  `id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `versification` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'protestant',
  `numeral_system_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `scope` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `script` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `derived` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `copyright` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `priority` tinyint unsigned NOT NULL DEFAULT '0',
  `reviewed` tinyint(1) DEFAULT '0',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bibles_id_unique` (`id`),
  KEY `bibles_language_id_foreign` (`language_id`),
  KEY `bibles_numeral_system_id_foreign` (`numeral_system_id`),
  KEY `bibles_script_foreign` (`script`),
  KEY `priority` (`priority`),
  CONSTRAINT `FK_alphabets_bibles` FOREIGN KEY (`script`) REFERENCES `alphabets` (`script`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bibles` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_numeral_systems_bibles` FOREIGN KEY (`numeral_system_id`) REFERENCES `numeral_systems` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_filesets` (
  `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `set_type_code` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `set_size_code` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `hidden` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `archived` tinyint(1) NOT NULL DEFAULT '0',
  `content_loaded` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`hash_id`),
  UNIQUE KEY `unique_prefix_for_s3` (`id`,`asset_id`,`set_type_code`),
  KEY `bible_filesets_bucket_id_foreign` (`asset_id`),
  KEY `bible_filesets_set_type_code_foreign` (`set_type_code`),
  KEY `bible_filesets_set_size_code_foreign` (`set_size_code`),
  KEY `bible_filesets_id_index` (`id`),
  KEY `bible_filesets_hash_id_index` (`hash_id`),
  CONSTRAINT `FK_assets_bible_filesets` FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_fileset_sizes_bible_filesets` FOREIGN KEY (`set_size_code`) REFERENCES `bible_fileset_sizes` (`set_size_code`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_fileset_types_bible_filesets` FOREIGN KEY (`set_type_code`) REFERENCES `bible_fileset_types` (`set_type_code`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_fileset_connections` (
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`hash_id`,`bible_id`),
  KEY `bible_fileset_connections_hash_id_foreign` (`hash_id`),
  KEY `bible_fileset_connections_bible_id_index` (`bible_id`),
  KEY `index_hash_id` (`hash_id`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_connections` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bibles_bible_fileset_connections` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_fileset_tags` (
  `hash_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `admin_only` tinyint(1) NOT NULL,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'UTC',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'UTC',
  PRIMARY KEY (`hash_id`,`name`,`language_id`),
  KEY `bible_fileset_tags_hash_id_index` (`hash_id`),
  KEY `bible_fileset_tags_iso_index` (`iso`),
  KEY `language_id` (`language_id`),
  KEY `description` (`description`(4)),
  KEY `hashid_name_index` (`hash_id`,`name`),
  KEY `name_index` (`name`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_tags` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bible_fileset_tags` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_languages_bible_fileset_tags_iso` FOREIGN KEY (`iso`) REFERENCES `languages` (`iso`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_files` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `hash_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapter_start` tinyint unsigned DEFAULT NULL,
  `chapter_end` tinyint unsigned DEFAULT NULL,
  `verse_start` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `verse_end` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `file_size` int unsigned DEFAULT NULL,
  `duration` int unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `verse_sequence` tinyint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_bible_file_by_reference` (`hash_id`,`book_id`,`chapter_start`,`verse_start`),
  KEY `bible_files_book_id_foreign` (`book_id`),
  CONSTRAINT `FK_bible_filesets_bible_files` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_books_bible_files` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=3844762 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_verses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapter` tinyint unsigned NOT NULL,
  `verse_start` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `verse_end` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `verse_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `verse_sequence` tinyint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_text_reference` (`hash_id`,`book_id`,`chapter`,`verse_start`),
  KEY `bible_text_book_id_foreign` (`book_id`),
  KEY `bible_text_hash_id_index` (`hash_id`),
  KEY `index_hash_id` (`hash_id`),
  FULLTEXT KEY `verse_text` (`verse_text`),
  CONSTRAINT `bible_text_hash_id_foreign` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=38563148 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_translations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `vernacular` tinyint(1) NOT NULL DEFAULT '0',
  `vernacular_trade` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `background` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`language_id`,`bible_id`,`vernacular`),
  KEY `bible_translations_language_id_foreign` (`language_id`),
  KEY `bible_translations_bible_id_foreign` (`bible_id`),
  FULLTEXT KEY `ft_index_bible_translations_name` (`name`),
  CONSTRAINT `FK_bibles_bible_translations` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bible_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=10334 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_translator` (
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `translator_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `bible_translator_translator_id_foreign` (`translator_id`),
  KEY `bible_translator_bible_id_index` (`bible_id`),
  CONSTRAINT `FK_bibles_bible_translator` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `translators` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `born` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `died` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `countries` (
  `id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `iso_a3` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `fips` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `wfb` tinyint(1) NOT NULL DEFAULT '0',
  `ethnologue` tinyint(1) NOT NULL DEFAULT '0',
  `continent` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `introduction` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `overview` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `countries_iso_a3_unique` (`iso_a3`),
  FULLTEXT KEY `ft_index_countries_name_iso_a3` (`name`,`iso_a3`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_fileset_modes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_fileset_types` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `set_type_code` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mode_id` tinyint unsigned NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bible_fileset_types_set_type_code_unique` (`set_type_code`),
  UNIQUE KEY `bible_fileset_types_name_unique` (`name`),
  KEY `bible_fileset_types_mode_index` (`mode_id`),
  CONSTRAINT `FK_bible_fileset_modes_bible_fileset_types` FOREIGN KEY (`mode_id`) REFERENCES `bible_fileset_modes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_file_stream_bandwidths` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `bible_file_id` int unsigned NOT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bandwidth` int unsigned NOT NULL,
  `resolution_width` int unsigned DEFAULT NULL,
  `resolution_height` int unsigned DEFAULT NULL,
  `codec` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `stream` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `bible_file_video_resolutions_bible_file_id_foreign` (`bible_file_id`),
  CONSTRAINT `FK_bible_files_bible_file_stream_bandwidths` FOREIGN KEY (`bible_file_id`) REFERENCES `bible_files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `bible_file_stream_ts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `stream_bandwidth_id` int unsigned NOT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `runtime` double(8,2) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_name` (`file_name`),
  KEY `bible_file_video_transport_stream_video_resolution_id_foreign` (`stream_bandwidth_id`),
  CONSTRAINT `FK_stream_bandwidths_stream_ts` FOREIGN KEY (`stream_bandwidth_id`) REFERENCES `bible_file_stream_bandwidths` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;