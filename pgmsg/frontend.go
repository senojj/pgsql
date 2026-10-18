package pgmsg

import (
	"math"

	"github.com/senojj/pgsql/pgio"
)

type Bind struct {
	DestinationName      string
	SourceName           string
	ParameterFormatCodes []byte
	ParameterData        [][]byte
	ColumnFormatCodes    []byte
}

func (x *Bind) AppendBinary(b []byte) ([]byte, error) {
	const sizeParamFmtCodeCount = 2
	const sizeParamFmtCode = 2
	const sizeParamDataCount = 2
	const sizeParamDatum = 4
	const sizeColFmtCodeCount = 2
	const sizeColFmtCode = 2

	sizeDestinationName := len(x.DestinationName) + 1 // null terminated string
	sizeSourceName := len(x.SourceName) + 1           // null terminated string

	paramFmtCodeCount := len(x.ParameterFormatCodes)
	if paramFmtCodeCount > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeParamFmtCodes := paramFmtCodeCount * sizeParamFmtCode

	paramDataCount := len(x.ParameterData)
	if paramDataCount > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeParamData := 0

	for i := range paramDataCount {
		length := len(x.ParameterData[i])
		if length > math.MaxInt32 {
			return b, invalidFormat(pgio.ErrValueOverflow)
		}
		sizeParamData += sizeParamDatum + length // size prefixed
	}

	colFmtCodeCount := len(x.ColumnFormatCodes)
	if colFmtCodeCount > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeColFmtCodes := colFmtCodeCount * sizeColFmtCode

	length := sizeMessageLength +
		sizeDestinationName +
		sizeSourceName +
		sizeParamFmtCodeCount +
		sizeParamFmtCodes +
		sizeParamDataCount +
		sizeParamData +
		sizeColFmtCodeCount +
		sizeColFmtCodes

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeBind)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.DestinationName)
	buf.AppendString(x.SourceName)
	buf.AppendInt16(int16(paramFmtCodeCount))

	for _, format := range x.ParameterFormatCodes {
		buf.AppendInt16(int16(format))
	}
	buf.AppendInt16(int16(paramDataCount))

	for i := range paramDataCount {
		sizeData := len(x.ParameterData[i])
		buf.AppendInt32(int32(sizeData))
		buf.AppendByte(x.ParameterData[i]...)
	}

	buf.AppendInt16(int16(colFmtCodeCount))

	for _, format := range x.ColumnFormatCodes {
		buf.AppendInt16(int16(format))
	}
	return buf.Bytes, nil
}

func (x *Bind) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeBind, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	destination, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	source, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	paramFmtCodeCount, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}
	parameterFormatCodes := make([]byte, paramFmtCodeCount)

	for i := range paramFmtCodeCount {
		code, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		parameterFormatCodes[i] = byte(code)
	}

	paramDataCount, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}
	parameterData := make([][]byte, paramDataCount)

	for i := range paramDataCount {
		dataLen, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}

		data, err := buf.ShiftBytes(int(dataLen))
		if err != nil {
			return invalidFormat(err)
		}
		parameterData[i] = data
	}

	colFmtCodeCount, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}
	columnFormatCodes := make([]byte, colFmtCodeCount)

	for i := range colFmtCodeCount {
		data, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		columnFormatCodes[i] = byte(data)
	}

	if buf.Len() > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}

	x.DestinationName = destination
	x.SourceName = source
	x.ParameterFormatCodes = parameterFormatCodes
	x.ParameterData = parameterData
	x.ColumnFormatCodes = columnFormatCodes
	return nil
}

type CancelRequest struct {
	ProcessID int32
	SecretKey []byte
}

func (x *CancelRequest) AppendBinary(b []byte) ([]byte, error) {
	const sizeCode = 4
	const sizeProcessID = 4

	sizeSecretKey := len(x.SecretKey)

	if sizeSecretKey > 256 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	if sizeSecretKey < 4 {
		return b, invalidFormat(pgio.ErrValueUnderflow)
	}

	length := sizeMessageLength +
		sizeCode +
		sizeProcessID +
		sizeSecretKey

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(length)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(CodeCancelRequest)
	buf.AppendInt32(x.ProcessID)
	buf.AppendByte(x.SecretKey...)
	return buf.Bytes, nil
}

func (x *CancelRequest) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	code, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if code != CodeCancelRequest {
		return invalidFormat(ErrUnknownCode)
	}

	processID, err := buf.ShiftInt32()
	if err != nil {
		return err
	}

	x.ProcessID = processID
	x.SecretKey = make([]byte, buf.Len())
	copy(x.SecretKey, buf.Bytes)
	return nil
}

type Close struct {
	ObjectKind byte
	ObjectName string
}

func (x *Close) AppendBinary(b []byte) ([]byte, error) {
	const sizeType = 1

	sizeName := len(x.ObjectName) + 1 // null terminated string

	length := sizeMessageLength + sizeType + sizeName

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeClose)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.ObjectKind)
	buf.AppendString(x.ObjectName)
	return buf.Bytes, nil
}

func (x *Close) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeClose, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	kind, err := buf.ShiftByte()
	if err != nil {
		return invalidFormat(err)
	}

	name, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}
	x.ObjectKind = kind
	x.ObjectName = name
	return nil
}

type CopyFail struct {
	Message string
}

func (x *CopyFail) AppendBinary(b []byte) ([]byte, error) {
	sizeMessage := len(x.Message) + 1 // null terminated string
	length := sizeMessageLength + sizeMessage

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCopyFail)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Message)
	return buf.Bytes, nil
}

func (x *CopyFail) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCopyFail, b)
	if err != nil {
		return invalidFormat(err)
	}

	m, b, err := pgio.ShiftString(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	x.Message = m
	return nil
}

type Describe struct {
	ObjectKind byte
	ObjectName string
}

func (x *Describe) AppendBinary(b []byte) ([]byte, error) {
	const sizeType = 1

	length := sizeMessageLength +
		sizeType +
		len(x.ObjectName) + 1 // null terminated string

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeDescribe)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.ObjectKind)
	buf.AppendString(x.ObjectName)
	return buf.Bytes, nil
}

func (x *Describe) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeDescribe, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	kind, err := buf.ShiftByte()
	if err != nil {
		return invalidFormat(err)
	}

	name, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	x.ObjectKind = kind
	x.ObjectName = name
	return nil
}

type Execute struct {
	PortalName string
	RowLimit   int32
}

func (x *Execute) AppendBinary(b []byte) ([]byte, error) {
	const sizeRowLimit = 4

	sizePortal := len(x.PortalName) + 1 // null terminated string

	length := sizeMessageLength + sizePortal + sizeRowLimit

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeExecute)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.PortalName)
	buf.AppendInt32(x.RowLimit)
	return buf.Bytes, nil
}

func (x *Execute) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeExecute, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	portal, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	limit, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	x.PortalName = portal
	x.RowLimit = limit
	return nil
}

type Flush struct{}

func (x *Flush) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeFlush)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *Flush) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeFlush, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type FunctionCall struct {
	ObjectID int32

	// ArgumentFormats may have zero elements, indicating that there are no
	// arguments, or that all arguments use the default format (text); or one,
	// in which case the specified format code is applied to all arguments; or
	// its element count may equal the total number of arguments.
	ArgumentFormats []byte

	ArgumentValues [][]byte
	ResultFormat   byte
}

func (x *FunctionCall) AppendBinary(b []byte) ([]byte, error) {
	const sizeObjectID = 4
	const sizeCountFormats = 2
	const sizeFormat = 2
	const sizeCountArguments = 2
	const sizeArgumentLength = 4
	const sizeResultFormat = 2

	countFormats := len(x.ArgumentFormats)

	if countFormats > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	countArguments := len(x.ArgumentValues)

	if countArguments > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeFormats := countFormats * sizeFormat
	sizeArguments := 0

	for i := range countArguments {
		sizeArgumentValue := len(x.ArgumentValues[i])

		if sizeArgumentValue > math.MaxInt32 {
			return b, invalidFormat(pgio.ErrValueOverflow)
		}
		sizeArguments += sizeArgumentLength + sizeArgumentValue
	}

	length := sizeMessageLength +
		sizeObjectID +
		sizeCountFormats +
		sizeFormats +
		sizeCountArguments +
		sizeArguments +
		sizeResultFormat

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeFunctionCall)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.ObjectID)
	buf.AppendInt16(int16(countFormats))

	for _, format := range x.ArgumentFormats {
		buf.AppendInt16(int16(format))
	}
	buf.AppendInt16(int16(countArguments))

	for i := range countArguments {
		value := x.ArgumentValues[i]
		length := len(value)
		buf.AppendInt32(int32(length))
		buf.AppendByte(value...)
	}
	buf.AppendInt16(int16(x.ResultFormat))
	return buf.Bytes, nil
}

func (x *FunctionCall) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeFunctionCall, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	objectID, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	countFormats, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	formats := make([]byte, 0, countFormats)
	for range countFormats {
		format, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		formats = append(formats, byte(format))
	}

	countArguments, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	arguments := make([][]byte, 0, countArguments)
	for range countArguments {
		length, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}
		value, err := buf.ShiftBytes(int(length))
		if err != nil {
			return invalidFormat(err)
		}
		arguments = append(arguments, value)
	}
	resultFormat, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}
	x.ObjectID = objectID
	x.ArgumentFormats = formats
	x.ArgumentValues = arguments
	x.ResultFormat = byte(resultFormat)
	return nil
}

type GSSENCRequest struct{}

func (x *GSSENCRequest) AppendBinary(b []byte) ([]byte, error) {
	const sizeCode = 4
	const length = sizeMessageLength + sizeCode

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(length)
	buf.AppendInt32(length)
	buf.AppendInt32(CodeEncryptionRequest)
	return buf.Bytes, nil
}

func (x *GSSENCRequest) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	code, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if code != CodeEncryptionRequest {
		return invalidFormat(ErrUnknownCode)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type GSSResponse struct {
	Data []byte
}

func (x *GSSResponse) AppendBinary(b []byte) ([]byte, error) {
	length := sizeMessageLength + len(x.Data)

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}
	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeGSSResponse)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *GSSResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeGSSResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type Parse struct {
	ParameterDataTypes       []int32
	DestinationStatementName string
	Query                    string
}

func (x *Parse) AppendBinary(b []byte) ([]byte, error) {
	const sizeParameterDataType = 4
	const sizeParameterDataTypeCount = 2

	sizeDestinationStatementName := len(x.DestinationStatementName) + 1 // null terminated string
	sizeQuery := len(x.Query) + 1                                       // null terminated string

	countParameterDataTypes := len(x.ParameterDataTypes)

	if countParameterDataTypes > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	length := sizeMessageLength +
		sizeDestinationStatementName +
		sizeQuery +
		sizeParameterDataTypeCount +
		countParameterDataTypes*sizeParameterDataType

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(byte(TypeParse))
	buf.AppendInt32(int32(length))
	buf.AppendString(x.DestinationStatementName)
	buf.AppendString(x.Query)
	buf.AppendInt16(int16(countParameterDataTypes))

	for i := range countParameterDataTypes {
		buf.AppendInt32(x.ParameterDataTypes[i])
	}
	return buf.Bytes, nil
}

func (x *Parse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeParse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	destinationStatementName, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	query, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	countParameterDataTypes, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	parameterDataTypes := make([]int32, 0, countParameterDataTypes)

	for range countParameterDataTypes {
		parameterDataType, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}
		parameterDataTypes = append(parameterDataTypes, parameterDataType)
	}

	x.DestinationStatementName = destinationStatementName
	x.Query = query
	x.ParameterDataTypes = parameterDataTypes
	return nil
}

type PasswordMessage struct {
	Password string
}

func (x *PasswordMessage) AppendBinary(b []byte) ([]byte, error) {
	sizePassword := len(x.Password) + 1 // null terminated string

	length := sizeMessageLength + sizePassword

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypePasswordMessage)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Password)
	return buf.Bytes, nil
}

func (x *PasswordMessage) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypePasswordMessage, b)
	if err != nil {
		return invalidFormat(err)
	}

	password, b, err := pgio.ShiftString(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}

	x.Password = password
	return nil
}

type Query struct {
	Value string
}

func (x *Query) AppendBinary(b []byte) ([]byte, error) {
	sizeQuery := len(x.Value)

	length := sizeMessageLength + sizeQuery

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeQuery)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Value)
	return buf.Bytes, nil
}

func (x *Query) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeQuery, b)
	if err != nil {
		return invalidFormat(err)
	}

	query, b, err := pgio.ShiftString(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}

	x.Value = query
	return nil
}

type SASLInitialResponse struct {
	Name     string
	Response []byte // will be nil when there is no initial response.
}

func (x *SASLInitialResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeMechanism = 4

	sizeName := len(x.Name) + 1 // null terminated string
	sizeResponse := len(x.Response)

	if sizeResponse > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	length := sizeMessageLength +
		sizeName +
		sizeMechanism +
		sizeResponse

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeSASLInitialResponse)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Name)
	buf.AppendInt32(int32(sizeResponse))
	buf.AppendByte(x.Response...)
	return buf.Bytes, nil
}

func (x *SASLInitialResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeSASLInitialResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	name, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	sizeResponse, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	response, err := buf.ShiftBytes(int(sizeResponse))
	if err != nil {
		return invalidFormat(err)
	}

	if buf.Len() > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}

	x.Name = name
	x.Response = response
	return nil
}

type SASLResponse struct {
	Data []byte
}

func (x *SASLResponse) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeSASLResponse)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes, nil
}

func (x *SASLResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeSASLResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type SSLRequest struct{}

func (x *SSLRequest) AppendBinary(b []byte) ([]byte, error) {
	const sizeCode = 4

	length := sizeMessageLength + sizeCode

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(length)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(CodeSSLRequest)
	return buf.Bytes, nil
}

func (x *SSLRequest) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	code, b, err := pgio.ShiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if code != CodeSSLRequest {
		return invalidFormat(ErrUnknownCode)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type StartupMessage struct {
	ProtocolVersion int32
	Parameters      map[string]string
}

func (x *StartupMessage) AppendBinary(b []byte) ([]byte, error) {
	const sizeProtocolVersion = 4

	var sizeParameters int

	for key, value := range x.Parameters {
		sizeParameters += (len(key) + 1)   // null terminated string
		sizeParameters += (len(value) + 1) // null terminated string
	}
	sizeParameters += 1 // null terminated list

	length := sizeMessageLength +
		sizeProtocolVersion +
		sizeParameters

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(length)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(int32(x.ProtocolVersion))

	for key, value := range x.Parameters {
		buf.AppendString(key)
		buf.AppendString(value)
	}
	buf.AppendByte(0)
	return buf.Bytes, nil
}

func (x *StartupMessage) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	protocolVersion, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	parameters := make(map[string]string)

	key, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	for len(key) > 0 {
		value, err := buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}
		parameters[key] = value

		key, err = buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}
	}

	x.ProtocolVersion = protocolVersion
	x.Parameters = parameters
	return nil
}
