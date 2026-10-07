package pgwire

const (
	cancelHigh        int32 = 1234
	cancelLow         int32 = 5678
	CodeCancelRequest int32 = cancelLow | cancelHigh<<16
)

const (
	sslHigh        int32 = 1234
	sslLow         int32 = 5679
	CodeSSLRequest int32 = sslLow | sslHigh<<16
)

const (
	major3             int32 = 3
	minor2             int32 = 2
	ProtocolVersion3_2 int32 = minor2 | major3<<16
)

const (
	ParamUser        string = "user"
	ParamDatabase    string = "database"
	ParamOptions     string = "options"
	ParamReplication string = "replication"
)

// Backend messages
const (
	TypeAuthentication           byte = 'R'
	TypeBackendKeyData           byte = 'K'
	TypeBindComplete             byte = '2'
	TypeCloseComplete            byte = '3'
	TypeCommandComplete          byte = 'C'
	TypeCopyInResponse           byte = 'G'
	TypeCopyOutResponse          byte = 'H'
	TypeCopyBothResponse         byte = 'W'
	TypeDataRow                  byte = 'D'
	TypeEmptyQueryResponse       byte = 'I'
	TypeErrorResponse            byte = 'E'
	TypeFunctionCallResponse     byte = 'V'
	TypeNegotiateProtocolVersion byte = 'v'
	TypeNoData                   byte = 'n'
	TypeNoticeResponse           byte = 'N'
	TypeNotificationResponse     byte = 'A'
	TypeParameterDescription     byte = 't'
	TypeParameterStatus          byte = 'S'
	TypeParseComplete            byte = '1'
	TypePortalSuspend            byte = 's'
	TypeReadyForQuery            byte = 'Z'
	TypeRowDescription           byte = 'T'
)

// Frontend messages
const (
	TypeBind                byte = 'B'
	TypeClose               byte = 'C'
	TypeCopyFail            byte = 'f'
	TypeDescribe            byte = 'D'
	TypeExecute             byte = 'E'
	TypeFlush               byte = 'H'
	TypeFunctionCall        byte = 'F'
	TypeGSSResponse         byte = 'p'
	TypeParse               byte = 'P'
	TypePasswordMessage     byte = 'p'
	TypeQuery               byte = 'Q'
	TypeSASLInitialResponse byte = 'p'
	TypeSASLResponse        byte = 'p'
	TypeSync                byte = 'S'
	TypeTerminate           byte = 'X'
)

// Backend + Frontend messages
const (
	TypeCopyData byte = 'd'
	TypeCopyDone byte = 'c'
)

const (
	AuthOk                int32 = 0
	AuthKerberosV5        int32 = 2
	AuthClearTextPassword int32 = 3
	AuthMD5Password       int32 = 5
	AuthGSS               int32 = 7
	AuthGSSContinue       int32 = 8
	AuthSSPI              int32 = 9
	AuthSASL              int32 = 10
	AuthSASLContinue      int32 = 11
	AuthSASLFinal         int32 = 12
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

const (
	sizeMessageType   = 1
	sizeMessageLength = 4
	sizeAuth      	  = 4
)
