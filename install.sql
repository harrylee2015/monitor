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
	'title' TEXT NOT NULL,
	'email' TEXT NOT NULL
)
;

CREATE TABLE 'PaymentAddress'
(
  'id' INTEGER PRIMARY KEY AUTOINCREMENT,
  'groupId' INTEGER NOT NULL,
  'groupName' TEXT NOT NULL,
	'address' TEXT NOT NULL
)
;

CREATE TABLE 'HostInfo'
(
  'hostId' INTEGER PRIMARY KEY AUTOINCREMENT,
  'hostName'  TEXT NOT NULL,
  'groupId' INTEGER NOT NULL,
  'groupName' TEXT NOT NULL,
  'hostIp' TEXT NOT NULL,
  'sshPort' INTEGER ,
  'userName' TEXT,
  'passWd' TEXT,
  'isCheckResource' NUMERIC NOT NULL,
	'processName' TEXT NOT NULL,
	'serverPort' INTEGER NOT NULL,
	'mainNet' TEXT NOT NULL,
    'netPort' INTEGER ,
	'createTime' INTEGER NOT NULL,
	'updateTime' INTEGER NOT NULL
)
;

CREATE TABLE 'ResourceInfo'
(
  'id' INTEGER PRIMARY KEY AUTOINCREMENT,
  'groupId' INTEGER NOT NULL,
  'hostId' INTEGER NOT NULL,
  'memTotal' INTEGER NOT NULL,
  'memUsedPercent' REAL NOT NULL,
  'cpuTotal' INTEGER NOT NULL,
  'cpuUsedPercent' REAL NOT NULL,
  'diskTotal' INTEGER NOT NULL,
  'diskUsedPercent' REAL NOT NULL,
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
  'serverStatus' NUMERIC NOT NULL,
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

CREATE INDEX 'IDX_Balance'
  ON 'Balance' ('createTime' ASC)
;

CREATE INDEX 'IDX_Balance_Time'
  ON 'Balance' ('groupId' ASC)
;

CREATE INDEX 'IDX_ResourceInfo_Time'
  ON 'ResourceInfo' ('createTime' ASC)
;

CREATE INDEX 'IDX_Monitor'
 ON 'Monitor' ('groupId','hostId' ASC)
;

CREATE INDEX 'IDX_Warning'
  ON 'Warning' ('groupId','hostId' ASC)
;
