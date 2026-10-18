package pgmsg

const (
	sizeMessageType   = 1
	sizeMessageLength = 4
	sizeAuth          = 4
)

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
	encHigh               int32 = 1234
	encLow                int32 = 5680
	codeEncryptionRequest int32 = encLow | encHigh<<16
)

// Backend messages
const (
	typeAuthentication           byte = 'R'
	typeBackendKeyData           byte = 'K'
	typeBindComplete             byte = '2'
	typeCloseComplete            byte = '3'
	typeCommandComplete          byte = 'C'
	typeCopyInResponse           byte = 'G'
	typeCopyOutResponse          byte = 'H'
	typeCopyBothResponse         byte = 'W'
	typeDataRow                  byte = 'D'
	typeEmptyQueryResponse       byte = 'I'
	typeErrorResponse            byte = 'E'
	typeFunctionCallResponse     byte = 'V'
	typeNegotiateProtocolVersion byte = 'v'
	typeNoData                   byte = 'n'
	typeNoticeResponse           byte = 'N'
	typeNotificationResponse     byte = 'A'
	typeParameterDescription     byte = 't'
	typeParameterStatus          byte = 'S'
	typeParseComplete            byte = '1'
	typePortalSuspend            byte = 's'
	typeReadyForQuery            byte = 'Z'
	typeRowDescription           byte = 'T'
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
	authTypeOk                int32 = 0
	authTypeKerberosV5        int32 = 2
	authTypeClearTextPassword int32 = 3
	authTypeMD5Password       int32 = 5
	authTypeGSS               int32 = 7
	authTypeGSSContinue       int32 = 8
	authTypeSSPI              int32 = 9
	authTypeSASL              int32 = 10
	authTypeSASLContinue      int32 = 11
	authTypeSASLFinal         int32 = 12
)
