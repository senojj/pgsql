package pgsql

import (
	"math"
)

type protocolVersion int32

func (x protocolVersion) Major() int32 {
	return int32(x) >> 16
}

func (x protocolVersion) Minor() int32 {
	return int32(x) & 0xFF
}

type msgBind struct {
	DestinationName      string
	SourceName           string
	ParameterFormatCodes []byte
	ParameterData        [][]byte
	ColumnFormatCodes    []byte
}

func (x *msgBind) message() {}

func (x *msgBind) frontend() {}

func (x *msgBind) AppendBinary(b []byte) ([]byte, error) {
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
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeParamFmtCodes := paramFmtCodeCount * sizeParamFmtCode

	paramDataCount := len(x.ParameterData)
	if paramDataCount > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeParamData := 0

	for i := range paramDataCount {
		length := len(x.ParameterData[i])
		if length > math.MaxInt32 {
			return b, invalidFormat(ErrValueOverflow)
		}
		sizeParamData += sizeParamDatum + length // size prefixed
	}

	colFmtCodeCount := len(x.ColumnFormatCodes)
	if colFmtCodeCount > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
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
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeBind)
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
	return buf.Bytes(), nil
}

func (x *msgBind) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeBind, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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
		return invalidFormat(ErrValueOverflow)
	}

	x.DestinationName = destination
	x.SourceName = source
	x.ParameterFormatCodes = parameterFormatCodes
	x.ParameterData = parameterData
	x.ColumnFormatCodes = columnFormatCodes
	return nil
}

type msgCancelRequest struct {
	ProcessID int32
	SecretKey []byte
}

func (x *msgCancelRequest) message() {}

func (x *msgCancelRequest) frontend() {}

func (x *msgCancelRequest) AppendBinary(b []byte) ([]byte, error) {
	const sizeCode = 4
	const sizeProcessID = 4

	sizeSecretKey := len(x.SecretKey)

	if sizeSecretKey > 256 {
		return b, invalidFormat(ErrValueOverflow)
	}

	if sizeSecretKey < 4 {
		return b, invalidFormat(ErrValueUnderflow)
	}

	length := sizeMessageLength +
		sizeCode +
		sizeProcessID +
		sizeSecretKey

	buf := buffer{b}
	buf.Grow(length)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(codeCancelRequest)
	buf.AppendInt32(x.ProcessID)
	buf.AppendByte(x.SecretKey...)
	return buf.Bytes(), nil
}

func (x *msgCancelRequest) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	code, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if code != codeCancelRequest {
		return invalidFormat(ErrUnknownCode)
	}

	processID, err := buf.ShiftInt32()
	if err != nil {
		return err
	}

	x.ProcessID = processID
	x.SecretKey = make([]byte, buf.Len())
	copy(x.SecretKey, buf.Bytes())
	return nil
}

type msgClose struct {
	ObjectKind byte
	ObjectName string
}

func (x *msgClose) message() {}

func (x *msgClose) frontend() {}

func (x *msgClose) AppendBinary(b []byte) ([]byte, error) {
	const sizeType = 1

	sizeName := len(x.ObjectName) + 1 // null terminated string

	length := sizeMessageLength + sizeType + sizeName

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeClose)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.ObjectKind)
	buf.AppendString(x.ObjectName)
	return buf.Bytes(), nil
}

func (x *msgClose) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeClose, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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

type msgCopyFail struct {
	Message string
}

func (x *msgCopyFail) message() {}

func (x *msgCopyFail) frontend() {}

func (x *msgCopyFail) AppendBinary(b []byte) ([]byte, error) {
	sizeMessage := len(x.Message) + 1 // null terminated string
	length := sizeMessageLength + sizeMessage

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCopyFail)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Message)
	return buf.Bytes(), nil
}

func (x *msgCopyFail) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyFail, b)
	if err != nil {
		return invalidFormat(err)
	}

	m, b, err := shiftString(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	x.Message = m
	return nil
}

type msgDescribe struct {
	ObjectKind byte
	ObjectName string
}

func (x *msgDescribe) message() {}

func (x *msgDescribe) frontend() {}

func (x *msgDescribe) AppendBinary(b []byte) ([]byte, error) {
	const sizeType = 1

	length := sizeMessageLength +
		sizeType +
		len(x.ObjectName) + 1 // null terminated string

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeDescribe)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.ObjectKind)
	buf.AppendString(x.ObjectName)
	return buf.Bytes(), nil
}

func (x *msgDescribe) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeDescribe, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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

type msgExecute struct {
	PortalName string
	RowLimit   int32
}

func (x *msgExecute) message() {}

func (x *msgExecute) frontend() {}

func (x *msgExecute) AppendBinary(b []byte) ([]byte, error) {
	const sizeRowLimit = 4

	sizePortal := len(x.PortalName) + 1 // null terminated string

	length := sizeMessageLength + sizePortal + sizeRowLimit

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeExecute)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.PortalName)
	buf.AppendInt32(x.RowLimit)
	return buf.Bytes(), nil
}

func (x *msgExecute) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeExecute, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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

type msgFlush struct{}

func (x *msgFlush) message() {}

func (x *msgFlush) frontend() {}

func (x *msgFlush) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeFlush)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgFlush) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeFlush, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgFunctionCall struct {
	ObjectID int32

	// ArgumentFormats may have zero elements, indicating that there are no
	// arguments, or that all arguments use the default format (text); or one,
	// in which case the specified format code is applied to all arguments; or
	// its element count may equal the total number of arguments.
	ArgumentFormats []byte

	ArgumentValues [][]byte
	ResultFormat   byte
}

func (x *msgFunctionCall) message() {}

func (x *msgFunctionCall) frontend() {}

func (x *msgFunctionCall) AppendBinary(b []byte) ([]byte, error) {
	const sizeObjectID = 4
	const sizeCountFormats = 2
	const sizeFormat = 2
	const sizeCountArguments = 2
	const sizeArgumentLength = 4
	const sizeResultFormat = 2

	countFormats := len(x.ArgumentFormats)

	if countFormats > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	countArguments := len(x.ArgumentValues)

	if countArguments > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeFormats := countFormats * sizeFormat
	sizeArguments := 0

	for i := range countArguments {
		sizeArgumentValue := len(x.ArgumentValues[i])

		if sizeArgumentValue > math.MaxInt32 {
			return b, invalidFormat(ErrValueOverflow)
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
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeFunctionCall)
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
	return buf.Bytes(), nil
}

func (x *msgFunctionCall) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeFunctionCall, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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

type msgGSSENCRequest struct{}

func (x *msgGSSENCRequest) message() {}

func (x *msgGSSENCRequest) frontend() {}

func (x *msgGSSENCRequest) AppendBinary(b []byte) ([]byte, error) {
	const sizeCode = 4
	const length = sizeMessageLength + sizeCode

	buf := buffer{b}
	buf.Grow(length)
	buf.AppendInt32(length)
	buf.AppendInt32(codeEncryptionRequest)
	return buf.Bytes(), nil
}

func (x *msgGSSENCRequest) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	code, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if code != codeEncryptionRequest {
		return invalidFormat(ErrUnknownCode)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgGSSResponse struct {
	Data []byte
}

func (x *msgGSSResponse) message() {}

func (x *msgGSSResponse) frontend() {}

func (x *msgGSSResponse) AppendBinary(b []byte) ([]byte, error) {
	length := sizeMessageLength + len(x.Data)

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}
	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeGSSResponse)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *msgGSSResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeGSSResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type msgParse struct {
	ParameterDataTypes       []int32
	DestinationStatementName string
	Query                    string
}

func (x *msgParse) message() {}

func (x *msgParse) frontend() {}

func (x *msgParse) AppendBinary(b []byte) ([]byte, error) {
	const sizeParameterDataType = 4
	const sizeParameterDataTypeCount = 2

	sizeDestinationStatementName := len(x.DestinationStatementName) + 1 // null terminated string
	sizeQuery := len(x.Query) + 1                                       // null terminated string

	countParameterDataTypes := len(x.ParameterDataTypes)

	if countParameterDataTypes > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	length := sizeMessageLength +
		sizeDestinationStatementName +
		sizeQuery +
		sizeParameterDataTypeCount +
		countParameterDataTypes*sizeParameterDataType

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(byte(typeParse))
	buf.AppendInt32(int32(length))
	buf.AppendString(x.DestinationStatementName)
	buf.AppendString(x.Query)
	buf.AppendInt16(int16(countParameterDataTypes))

	for i := range countParameterDataTypes {
		buf.AppendInt32(x.ParameterDataTypes[i])
	}
	return buf.Bytes(), nil
}

func (x *msgParse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeParse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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

type msgPasswordMessage struct {
	Password string
}

func (x *msgPasswordMessage) message() {}

func (x *msgPasswordMessage) frontend() {}

func (x *msgPasswordMessage) AppendBinary(b []byte) ([]byte, error) {
	sizePassword := len(x.Password) + 1 // null terminated string

	length := sizeMessageLength + sizePassword

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typePasswordMessage)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Password)
	return buf.Bytes(), nil
}

func (x *msgPasswordMessage) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typePasswordMessage, b)
	if err != nil {
		return invalidFormat(err)
	}

	password, b, err := shiftString(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}

	x.Password = password
	return nil
}

type msgQuery struct {
	Value string
}

func (x *msgQuery) message() {}

func (x *msgQuery) frontend() {}

func (x *msgQuery) AppendBinary(b []byte) ([]byte, error) {
	sizeQuery := len(x.Value)

	length := sizeMessageLength + sizeQuery

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeQuery)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Value)
	return buf.Bytes(), nil
}

func (x *msgQuery) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeQuery, b)
	if err != nil {
		return invalidFormat(err)
	}

	query, b, err := shiftString(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}

	x.Value = query
	return nil
}

type msgSASLInitialResponse struct {
	Name     string
	Response []byte // will be nil when there is no initial response.
}

func (x *msgSASLInitialResponse) message() {}

func (x *msgSASLInitialResponse) frontend() {}

func (x *msgSASLInitialResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeMechanism = 4

	sizeName := len(x.Name) + 1 // null terminated string
	sizeResponse := len(x.Response)

	if sizeResponse > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	length := sizeMessageLength +
		sizeName +
		sizeMechanism +
		sizeResponse

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeSASLInitialResponse)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Name)
	buf.AppendInt32(int32(sizeResponse))
	buf.AppendByte(x.Response...)
	return buf.Bytes(), nil
}

func (x *msgSASLInitialResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeSASLInitialResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

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
		return invalidFormat(ErrValueOverflow)
	}

	x.Name = name
	x.Response = response
	return nil
}

type msgSASLResponse struct {
	Data []byte
}

func (x *msgSASLResponse) message() {}

func (x *msgSASLResponse) frontend() {}

func (x *msgSASLResponse) AppendBinary(b []byte) ([]byte, error) {
	sizeData := len(x.Data)

	length := sizeMessageLength + sizeData

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeSASLResponse)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.Data...)
	return buf.Bytes(), nil
}

func (x *msgSASLResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeSASLResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	x.Data = make([]byte, len(b))
	copy(x.Data, b)
	return nil
}

type msgSSLRequest struct{}

func (x *msgSSLRequest) message() {}

func (x *msgSSLRequest) frontend() {}

func (x *msgSSLRequest) AppendBinary(b []byte) ([]byte, error) {
	const sizeCode = 4

	length := sizeMessageLength + sizeCode

	buf := buffer{b}

	buf.Grow(length)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(codeSSLRequest)
	return buf.Bytes(), nil
}

func (x *msgSSLRequest) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	code, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if code != codeSSLRequest {
		return invalidFormat(ErrUnknownCode)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgStartupMessage struct {
	ProtocolVersion protocolVersion
	Parameters      map[string]string
}

func (x *msgStartupMessage) message() {}

func (x *msgStartupMessage) frontend() {}

func (x *msgStartupMessage) AppendBinary(b []byte) ([]byte, error) {
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
		return b, invalidFormat(ErrValueOverflow)
	}

	buf := buffer{b}

	buf.Grow(length)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(int32(x.ProtocolVersion))

	for key, value := range x.Parameters {
		buf.AppendString(key)
		buf.AppendString(value)
	}
	buf.AppendByte(0)
	return buf.Bytes(), nil
}

func (x *msgStartupMessage) UnmarshalBinary(b []byte) error {
	b, err := shiftLength(b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	pv, err := buf.ShiftInt32()
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

	x.ProtocolVersion = protocolVersion(pv)
	x.Parameters = parameters
	return nil
}
