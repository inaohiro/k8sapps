SET CHARACTER_SET_CLIENT = utf8mb4;
SET CHARACTER_SET_CONNECTION = utf8mb4;

USE db;

DROP TABLE IF EXISTS images;
CREATE TABLE images
(
  name  VARCHAR(64) NOT NULL COMMENT 'コンテナイメージ名',
  PRIMARY KEY (name)
)
  COMMENT = 'コンテナイメージ';

DROP TABLE IF EXISTS flavors;
CREATE TABLE flavors
(
  name  VARCHAR(64) NOT NULL COMMENT 'フレーバー名',
  PRIMARY KEY (name)
)
  COMMENT = 'フレーバー';
