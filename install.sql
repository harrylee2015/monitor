DROP TABLE IF EXISTS 'GroupInfo'
;

DROP TABLE IF EXISTS 'HostInfo'
;

DROP TABLE IF EXISTS 'PaymentAddress'
;

DROP TABLE IF EXISTS 'ResourceInfo'
;

DROP TABLE IF EXISTS 'Monitor'
;

DROP TABLE IF EXISTS 'Balance'
;
DROP TABLE IF EXISTS 'Warning'
;
CREATE TABLE 'GroupInfo'
(
  'groupId' INTEGER PRIMARY KEY AUTOINCREMENT,
	'groupName' TEXT NOT NULL,
	'describe' TEXT NOT NULL,
	'title' TEXT NOT NULL
)
;

CREATE TABLE 'PaymentAddress'
(
  'id' INTEGER PRIMARY KEY AUTOINCREMENT,
  'groupId' INTEGER NOT NULL,
	'address' TEXT NOT NULL
)
;

CREATE TABLE 'HostInfo'
(
  'hostId' INTEGER PRIMARY KEY AUTOINCREMENT,
  'groupId' INTEGER NOT NULL,
  'hostIp' TEXT NOT NULL,
  'sshPort' INTEGER ,
  'userName' TEXT,
  'passWd' TEXT,
  'isCheckResource' NUMERIC NOT NULL,
	'processName' TEXT NOT NULL,
	'serverPort' INTEGER NOT NULL,
	'createTime' INTEGER NOT NULL,
	'updateTime' INTEGER NOT NULL
)
;

CREATE TABLE 'ResourceInfo'
(
  'id' INTEGER PRIMARY KEY AUTOINCREMENT,
  'groupId' INTEGER NOT NULL,
  'hostId' INTEGER NOT NULL,
  'mem' INTEGER NOT NULL,
  'cpu' INTEGER NOT NULL,
  'disk' INTEGER NOT NULL,
  'createTime' INTEGER NOT NULL
)
;

CREATE TABLE 'Monitor'
(
	'id' INTEGER PRIMARY KEY AUTOINCREMENT,
	'groupId' INTEGER NOT NULL,
  'hostId' INTEGER NOT NULL,
  'hostIp' TEXT NOT NULL,
  'serverPort' INTEGER NOT NULL,
  'lastBlockHeight' INTEGER NOT NULL,
  'isSync' NUMERIC NOT NULL ,
  'lastBlockHash' TEXT NOT NULL,
  'updateTime' INTEGER NOT NULL
)
;

CREATE TABLE 'Balance'
(
	'id' INTEGER PRIMARY KEY AUTOINCREMENT,
	'groupId' INTEGER NOT NULL,
	'address' TEXT NOT NULL,
	'balance' INTEGER NOT NULL,
	'createTime' INTEGER NOT NULL
)
;

CREATE TABLE 'Warning'
(
	'id' INTEGER PRIMARY KEY AUTOINCREMENT,
	'hostId' INTEGER NOT NULL,
	'groupId' INTEGER NOT NULL,
	'type' NUMERIC NOT NULL,
	'warning' TEXT NOT NULL,
	'blockHeight' INTEGER NOT NULL,
	'createTime' INTEGER NOT NULL,
	'isClosed' NUMERIC NOT NULL,
	'updateTime' INTEGER NOT NULL
)
;

CREATE INDEX 'IDX_HostInfo'
 ON 'HostInfo' ('groupId' ASC)
;

CREATE INDEX 'IDX_ResourceInfo'
  ON 'ResourceInfo' ('hostId' ASC)
;

CREATE INDEX 'IDX_Monitor'
 ON 'Monitor' ('groupId' ASC)
;

CREATE INDEX 'IDX_Warning'
  ON 'Warning' ('groupId' ASC)
;
