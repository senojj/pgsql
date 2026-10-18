package pgsql

const (
	major3             int32 = 3
	minor2             int32 = 2
	protocolVersion3_2 int32 = minor2 | major3<<16
)

func decodeProtocolVersion(v int32) (int32, int32) {
	return v >> 16, v & 0xFF
}

func encodeProtocolVersion(major int32, minor int32) int32 {
	return minor | major<<16
}

const (
	paramUser        string = "user"
	paramDatabase    string = "database"
	paramOptions     string = "options"
	paramReplication string = "replication"
)

const (
	fieldSeverity         byte = 'S'
	fieldSeverityRaw      byte = 'V'
	fieldCode             byte = 'C'
	fieldMessage          byte = 'M'
	fieldDetail           byte = 'D'
	fieldHint             byte = 'H'
	fieldPosition         byte = 'P'
	fieldInternalPosition byte = 'p'
	fieldInternalQuery    byte = 'q'
	fieldWhere            byte = 'W'
	fieldSchema           byte = 's'
	fieldTable            byte = 't'
	fieldColumn           byte = 'c'
	fieldDataType         byte = 'd'
	fieldConstraint       byte = 'n'
	fieldFile             byte = 'F'
	fieldLine             byte = 'L'
	fieldRoutine          byte = 'R'
)

const (
	objectStatement byte = 'S'
	objectPortal    byte = 'P'
)

const (
	formatText   byte = 0
	formatBinary byte = 1
)

const (
	transactionStatusIdle   byte = 'I'
	transactionStatusActive byte = 'T'
	transactionStatusError  byte = 'E'
)
