package pgsql

const (
	cancelHigh        int32 = 1234
	cancelLow         int32 = 5678
	codeCancelRequest int32 = cancelLow | cancelHigh<<16
)

const (
	sslHigh        int32 = 1234
	sslLow         int32 = 5679
	codeSSLRequest int32 = sslLow | sslHigh<<16
)

const (
	major3             int32 = 3
	minor2             int32 = 2
	protocolVersion3_2 int32 = minor2 | major3<<16
)

const (
	paramUser        string = "user"
	paramDatabase    string = "database"
	paramOptions     string = "options"
	paramReplication string = "replication"
)

// Frontend messages
const (
	typeBind                byte = 'B'
	typeClose               byte = 'C'
	typeCopyFail            byte = 'f'
	typeDescribe            byte = 'D'
	typeExecute             byte = 'E'
	typeFlush               byte = 'H'
	typeFunctionCall        byte = 'F'
	typeGSSResponse         byte = 'p'
	typeParse               byte = 'P'
	typePasswordMessage     byte = 'p'
	typeQuery               byte = 'Q'
	typeSASLInitialResponse byte = 'p'
	typeSASLResponse        byte = 'p'
	typeSync                byte = 'S'
	typeTerminate           byte = 'X'
)

// Backend + Frontend messages
const (
	typeCopyData byte = 'd'
	typeCopyDone byte = 'c'
)

const (
	authOk                int32 = 0
	authKerberosV5        int32 = 2
	authClearTextPassword int32 = 3
	authMD5Password       int32 = 5
	authGSS               int32 = 7
	authGSSContinue       int32 = 8
	authSSPI              int32 = 9
	authSASL              int32 = 10
	authSASLContinue      int32 = 11
	authSASLFinal         int32 = 12
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

const (
	sizeMessageType   = 1
	sizeMessageLength = 4
	sizeAuth      	  = 4
)
