CREATE TABLE IF NOT EXISTS `hello` (
  `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` VARCHAR(64) NULL
);

INSERT INTO hello (name) VALUES('world')