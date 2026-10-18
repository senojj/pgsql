package pgmsg

const (
	sizeMessageType   = 1
	sizeMessageLength = 4
	sizeAuth          = 4
)

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
	encHigh               int32 = 1234
	encLow                int32 = 5680
	CodeEncryptionRequest int32 = encLow | encHigh<<16
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
	AuthTypeOk                int32 = 0
	AuthTypeKerberosV5        int32 = 2
	AuthTypeClearTextPassword int32 = 3
	AuthTypeMD5Password       int32 = 5
	AuthTypeGSS               int32 = 7
	AuthTypeGSSContinue       int32 = 8
	AuthTypeSSPI              int32 = 9
	AuthTypeSASL              int32 = 10
	AuthTypeSASLContinue      int32 = 11
	AuthTypeSASLFinal         int32 = 12
)
