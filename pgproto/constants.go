package pgproto

const (
	major3             int32 = 3
	minor2             int32 = 2
	ProtocolVersion3_2 int32 = minor2 | major3<<16
)

func DecodeProtocolVersion(v int32) (int32, int32) {
	return v >> 16, v & 0xFF
}

func EncodeProtocolVersion(major int32, minor int32) int32 {
	return minor | major<<16
}

const (
	ParamUser        string = "user"
	ParamDatabase    string = "database"
	ParamOptions     string = "options"
	ParamReplication string = "replication"
)

const (
	FieldSeverity         byte = 'S'
	FieldSeverityRaw      byte = 'V'
	FieldCode             byte = 'C'
	FieldMessage          byte = 'M'
	FieldDetail           byte = 'D'
	FieldHint             byte = 'H'
	FieldPosition         byte = 'P'
	FieldInternalPosition byte = 'p'
	FieldInternalQuery    byte = 'q'
	FieldWhere            byte = 'W'
	FieldSchema           byte = 's'
	FieldTable            byte = 't'
	FieldColumn           byte = 'c'
	FieldDataType         byte = 'd'
	FieldConstraint       byte = 'n'
	FieldFile             byte = 'F'
	FieldLine             byte = 'L'
	FieldRoutine          byte = 'R'
)

const (
	ObjectStatement byte = 'S'
	ObjectPortal    byte = 'P'
)

const (
	FormatText   byte = 0
	FormatBinary byte = 1
)

const (
	TransactionStatusIdle   byte = 'I'
	TransactionStatusActive byte = 'T'
	TransactionStatusError  byte = 'E'
)
