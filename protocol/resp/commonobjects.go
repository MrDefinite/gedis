package resp

import "github.com/MrDefinite/gedis/basicdata"

type CommonObjectStruct struct {
	Crlf, Ok, Err, EmptyBulk, Zero, One, NegOne, Pong, Space,
	Colon, NullBulk, NullMultiBulk, Queued,
	EmptyMultiBulk, WrongTypeErr, NoKeyErr, SynTaxErr, SameObjectErr,
	OutOfRangeErr, NoScriptErr, LoadingErr, SlowScriptErr, BgSaveErr,
	MasterDownErr, RoSlaveErr, ExecAbortErr, NoAuthErr, NoReplicasErr,
	BusyKeyErr, OomErr, Plus, EmptyScan *basicdata.GedisObject
}

var (
	CommonObjects = CommonObjectStruct{
		Crlf:           basicdata.CreateObject(basicdata.GedisString, "\r\n"),
		Ok:             basicdata.CreateObject(basicdata.GedisString, "+OK\r\n"),
		Err:            basicdata.CreateObject(basicdata.GedisString, "-ERR\r\n"),
		EmptyBulk:      basicdata.CreateObject(basicdata.GedisString, "$0\r\n\r\n"),
		Zero:           basicdata.CreateObject(basicdata.GedisString, ":0\r\n"),
		One:            basicdata.CreateObject(basicdata.GedisString, ":1\r\n"),
		NegOne:         basicdata.CreateObject(basicdata.GedisString, ":-1\r\n"),
		NullBulk:       basicdata.CreateObject(basicdata.GedisString, "$-1\r\n"),
		NullMultiBulk:  basicdata.CreateObject(basicdata.GedisString, "*-1\r\n"),
		EmptyMultiBulk: basicdata.CreateObject(basicdata.GedisString, "*0\r\n"),
		Pong:           basicdata.CreateObject(basicdata.GedisString, "+PONG\r\n"),
		Queued:         basicdata.CreateObject(basicdata.GedisString, "+QUEUED\r\n"),
		EmptyScan:      basicdata.CreateObject(basicdata.GedisString, "*2\r\n$1\r\n0\r\n*0\r\n"),

		WrongTypeErr:  basicdata.CreateObject(basicdata.GedisString, "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n"),
		NoKeyErr:      basicdata.CreateObject(basicdata.GedisString, "-ERR no such key\r\n"),
		SynTaxErr:     basicdata.CreateObject(basicdata.GedisString, "-ERR syntax error\r\n"),
		SameObjectErr: basicdata.CreateObject(basicdata.GedisString, "-ERR source and destination objects are the same\r\n"),
		OutOfRangeErr: basicdata.CreateObject(basicdata.GedisString, "-ERR index out of range\r\n"),
		NoScriptErr:   basicdata.CreateObject(basicdata.GedisString, "-NOSCRIPT No matching script. Please use EVAL.\r\n"),
		LoadingErr:    basicdata.CreateObject(basicdata.GedisString, "-LOADING Redis is loading the dataset in memory\r\n"),
		SlowScriptErr: basicdata.CreateObject(basicdata.GedisString, "-BUSY Redis is busy running a script. You can only call SCRIPT KILL or SHUTDOWN NOSAVE.\r\n"),
		MasterDownErr: basicdata.CreateObject(basicdata.GedisString, "-MASTERDOWN Link with MASTER is down and slave-serve-stale-data is set to 'no'.\r\n"),
		BgSaveErr:     basicdata.CreateObject(basicdata.GedisString, "-MISCONF Redis is configured to save RDB snapshots, but is currently not able to persist on disk. Commands that may modify the data set are disabled. Please check Redis logs for details about the error.\r\n"),
		RoSlaveErr:    basicdata.CreateObject(basicdata.GedisString, "-READONLY You can't write against a read only slave.\r\n"),
		NoAuthErr:     basicdata.CreateObject(basicdata.GedisString, "-NOAUTH Authentication required.\r\n"),
		OomErr:        basicdata.CreateObject(basicdata.GedisString, "-OOM command not allowed when used memory > 'maxmemory'.\r\n"),
		ExecAbortErr:  basicdata.CreateObject(basicdata.GedisString, "-EXECABORT Transaction discarded because of previous errors.\r\n"),
		NoReplicasErr: basicdata.CreateObject(basicdata.GedisString, "-NOREPLICAS Not enough good slaves to write.\r\n"),
		BusyKeyErr:    basicdata.CreateObject(basicdata.GedisString, "-BUSYKEY Target key name already exists.\r\n"),

		Space: basicdata.CreateObject(basicdata.GedisString, " "),
		Colon: basicdata.CreateObject(basicdata.GedisString, ":"),
		Plus:  basicdata.CreateObject(basicdata.GedisString, "+"),
	}
)
