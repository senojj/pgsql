package pgsql

import (
	"math"
)

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

type msgBackendKeyData struct {
	ProcessID int32
	SecretKey []byte
}

func (x *msgBackendKeyData) message() {}

func (x *msgBackendKeyData) backend() {}

func (x *msgBackendKeyData) AppendBinary(b []byte) ([]byte, error) {
	const sizeProcessID = 4
	sizeSecretKey := len(x.SecretKey)

	if sizeSecretKey < 4 {
		return b, invalidFormat(ErrValueUnderflow)
	}

	if sizeSecretKey > 256 {
		return b, invalidFormat(ErrValueOverflow)
	}

	length := sizeMessageLength + sizeProcessID + sizeSecretKey

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeBackendKeyData)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.ProcessID)
	buf.AppendByte(x.SecretKey...)
	return buf.Bytes(), nil
}

func (x *msgBackendKeyData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeBackendKeyData, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	processID, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if buf.Len() < 4 {
		return invalidFormat(ErrValueUnderflow)
	}

	if buf.Len() > 256 {
		return invalidFormat(ErrValueOverflow)
	}
	secretKey := make([]byte, buf.Len())
	copy(secretKey, buf.Bytes())

	x.ProcessID = processID
	x.SecretKey = secretKey
	return nil
}

type msgBindComplete struct{}

func (x *msgBindComplete) message() {}

func (x *msgBindComplete) backend() {}

func (x *msgBindComplete) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeBindComplete)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgBindComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeBindComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgCloseComplete struct{}

func (x *msgCloseComplete) message() {}

func (x *msgCloseComplete) backend() {}

func (x *msgCloseComplete) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCloseComplete)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgCloseComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCloseComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgCommandComplete struct {
	Tag string
}

func (x *msgCommandComplete) message() {}

func (x *msgCommandComplete) backend() {}

func (x *msgCommandComplete) AppendBinary(b []byte) ([]byte, error) {
	sizeTag := len(x.Tag) + 1 // null terminated string
	length := sizeMessageLength + sizeTag

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCommandComplete)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Tag)
	return buf.Bytes(), nil
}

func (x *msgCommandComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCommandComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	tag, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}
	x.Tag = tag
	return nil
}

type msgCopyInResponse struct {
	Format  int8
	Columns []int16
}

func (x *msgCopyInResponse) message() {}

func (x *msgCopyInResponse) backend() {}

func (x *msgCopyInResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeFormat = 1
	const sizeColCount = 2
	const sizeCol = 2

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeCols := countCols * sizeCol

	length := sizeMessageLength +
		sizeFormat +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCopyInResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt8(x.Format)
	buf.AppendInt16(int16(countCols))
	buf.AppendInt16(x.Columns...)
	return buf.Bytes(), nil
}

func (x *msgCopyInResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyInResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	format, err := buf.ShiftInt8()
	if err != nil {
		return invalidFormat(err)
	}

	length, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	columns := make([]int16, 0, length)

	for range length {
		value, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		columns = append(columns, value)
	}
	x.Format = format
	x.Columns = columns
	return nil
}

type msgCopyOutResponse struct {
	Format  int8
	Columns []int16
}

func (x *msgCopyOutResponse) message() {}

func (x *msgCopyOutResponse) backend() {}

func (x *msgCopyOutResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeFormat = 1
	const sizeColCount = 2
	const sizeCol = 2

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeCols := countCols * sizeCol

	length := sizeMessageLength +
		sizeFormat +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCopyOutResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt8(x.Format)
	buf.AppendInt16(int16(countCols))
	buf.AppendInt16(x.Columns...)
	return buf.Bytes(), nil
}

func (x *msgCopyOutResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyOutResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	format, err := buf.ShiftInt8()
	if err != nil {
		return invalidFormat(err)
	}

	length, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	columns := make([]int16, 0, length)

	for range length {
		value, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		columns = append(columns, value)
	}
	x.Format = format
	x.Columns = columns
	return nil
}

type msgCopyBothResponse struct {
	Format  int8
	Columns []int16
}

func (x *msgCopyBothResponse) message() {}

func (x *msgCopyBothResponse) backend() {}

func (x *msgCopyBothResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeFormat = 1
	const sizeColCount = 2
	const sizeCol = 2

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeCols := countCols * sizeCol

	length := sizeMessageLength +
		sizeFormat +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeCopyBothResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt8(x.Format)
	buf.AppendInt16(int16(countCols))
	buf.AppendInt16(x.Columns...)
	return buf.Bytes(), nil
}

func (x *msgCopyBothResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeCopyBothResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	format, err := buf.ShiftInt8()
	if err != nil {
		return invalidFormat(err)
	}

	length, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	columns := make([]int16, 0, length)

	for range length {
		value, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		columns = append(columns, value)
	}
	x.Format = format
	x.Columns = columns
	return nil
}

type msgDataRow struct {
	// First dimention is columns (can have 0 elements).
	// Second dimension is value data (can have 0 elements
	// or be nil).
	Columns [][]byte
}

func (x *msgDataRow) message() {}

func (x *msgDataRow) backend() {}

func (x *msgDataRow) AppendBinary(b []byte) ([]byte, error) {
	const sizeColCount = 2
	const sizeColLength = 4

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	sizeCols := 0

	for i := range countCols {
		lengthCol := len(x.Columns[i])

		if lengthCol > math.MaxInt32 {
			return b, invalidFormat(ErrValueOverflow)
		}
		sizeCols += sizeColLength + lengthCol
	}

	length := sizeMessageLength +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeDataRow)
	buf.AppendInt32(int32(length))
	buf.AppendInt16(int16(countCols))

	for i := range countCols {
		column := x.Columns[i]
		if column == nil {
			buf.AppendInt32(-1)
			continue
		}
		lengthCol := len(column)
		buf.AppendInt32(int32(lengthCol))
		buf.AppendByte(column...)
	}
	return buf.Bytes(), nil
}

func (x *msgDataRow) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeDataRow, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	countCols, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}
	columns := make([][]byte, 0, countCols)

	for range countCols {
		length, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}

		if length == -1 {
			columns = append(columns, nil)
			continue
		}

		data, err := buf.ShiftBytes(int(length))
		if err != nil {
			return invalidFormat(err)
		}
		columns = append(columns, data)
	}

	x.Columns = columns
	return nil
}

type msgEmptyQueryResponse struct{}

func (x *msgEmptyQueryResponse) message() {}

func (x *msgEmptyQueryResponse) backend() {}

func (x *msgEmptyQueryResponse) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeEmptyQueryResponse)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgEmptyQueryResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeEmptyQueryResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgErrorResponse struct {
	Fields []byte
	Values []string
}

func (x *msgErrorResponse) message() {}

func (x *msgErrorResponse) backend() {}

func (x *msgErrorResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeField = 1

	countFields := len(x.Fields)
	countValues := len(x.Values)

	sizeFields := countFields * sizeField

	length := sizeMessageLength + sizeFields

	for i := range countValues {
		value := x.Values[i]
		length += len(value) + 1 // null terminated string
	}

	length += 1 // null terminated list

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeErrorResponse)
	buf.AppendInt32(int32(length))

	for i := range countFields {
		buf.AppendByte(x.Fields[i])
		buf.AppendString(x.Values[i])
	}
	buf.AppendByte(0)
	return buf.Bytes(), nil
}

func (x *msgErrorResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeErrorResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	var fields []byte
	var values []string

	field, err := buf.ShiftByte()
	if err != nil {
		return invalidFormat(err)
	}

	for field != 0 {
		fields = append(fields, field)

		value, err := buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}
		values = append(values, value)

		field, err = buf.ShiftByte()
		if err != nil {
			return invalidFormat(err)
		}
	}

	x.Fields = fields
	x.Values = values
	return nil
}

type msgFunctionCallResponse struct {
	// Can be zero length or nil.
	Result []byte
}

func (x *msgFunctionCallResponse) message() {}

func (x *msgFunctionCallResponse) backend() {}

func (x *msgFunctionCallResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeResultLength = 4
	sizeResult := len(x.Result)

	if sizeResult > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	length := sizeMessageLength +
		sizeResultLength +
		sizeResult

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeFunctionCallResponse)
	buf.AppendInt32(int32(length))

	if x.Result == nil {
		buf.AppendInt32(int32(-1))
	} else {
		buf.AppendInt32(int32(sizeResult))
	}
	buf.AppendByte(x.Result...)
	return buf.Bytes(), nil
}

func (x *msgFunctionCallResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeFunctionCallResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	length, b, err := shiftInt32(b)
	if err != nil {
		return invalidFormat(err)
	}

	if length >= 0 {
		x.Result = make([]byte, length)
		copy(x.Result, b)
	}
	// Result remains nil when length < 0
	return nil
}

type msgNegotiateProtocolVersion struct {
	MinorVersionSupported int32
	UnrecognizedOptions   []string
}

func (x *msgNegotiateProtocolVersion) message() {}

func (x *msgNegotiateProtocolVersion) backend() {}

func (x *msgNegotiateProtocolVersion) AppendBinary(b []byte) ([]byte, error) {
	const sizeMinorVersion = 4
	const sizeUnrecognizedOptionCount = 4

	length := sizeMessageLength +
		sizeMinorVersion +
		sizeUnrecognizedOptionCount

	countUnrecognizedOptions := len(x.UnrecognizedOptions)

	if countUnrecognizedOptions > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	for _, option := range x.UnrecognizedOptions {
		length += len(option) + 1 // +1 null byte for string termination
	}

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeNegotiateProtocolVersion)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.MinorVersionSupported)
	buf.AppendInt32(int32(countUnrecognizedOptions))
	buf.AppendString(x.UnrecognizedOptions...)
	return buf.Bytes(), nil
}

func (x *msgNegotiateProtocolVersion) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeNegotiateProtocolVersion, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	minorVersion, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	countUnsupportedOptions, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	options := make([]string, 0, countUnsupportedOptions)

	for range countUnsupportedOptions {
		option, err := buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}
		options = append(options, option)
	}
	x.MinorVersionSupported = minorVersion
	x.UnrecognizedOptions = options
	return nil
}

type msgNoData struct{}

func (x *msgNoData) message() {}

func (x *msgNoData) backend() {}

func (x *msgNoData) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeNoData)
	buf.AppendInt32(length)
	return buf.Bytes(), nil
}

func (x *msgNoData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeNoData, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgNoticeResponse struct {
	Fields []byte
	Values []string
}

func (x *msgNoticeResponse) message() {}

func (x *msgNoticeResponse) backend() {}

func (x *msgNoticeResponse) AppendBinary(b []byte) ([]byte, error) {
	countFields := len(x.Fields)

	if countFields != len(x.Values) {
		return b, invalidFormat(ErrValueOverflow)
	}

	length := sizeMessageLength + countFields

	for _, value := range x.Values {
		length += len(value) + 1 // null terminated strings
	}
	length += 1 // null terminated list

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeNoticeResponse)
	buf.AppendInt32(int32(length))

	for i := range countFields {
		buf.AppendByte(x.Fields[i])
		buf.AppendString(x.Values[i])
	}
	buf.AppendByte(0)
	return buf.Bytes(), nil
}

func (x *msgNoticeResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeNoticeResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	field, err := buf.ShiftByte()
	if err != nil {
		return invalidFormat(err)
	}

	var fields []byte
	var values []string

	for field != 0 {
		fields = append(fields, field)

		value, err := buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}
		values = append(values, value)

		field, err = buf.ShiftByte()
		if err != nil {
			return invalidFormat(err)
		}
	}
	x.Fields = fields
	x.Values = values
	return nil
}

type msgNotificationResponse struct {
	ProcessID int32
	Channel   string
	Payload   string
}

func (x *msgNotificationResponse) message() {}

func (x *msgNotificationResponse) backend() {}

func (x *msgNotificationResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeProcessID = 4

	sizeChannel := len(x.Channel) + 1 // null terminated string
	sizePayload := len(x.Payload) + 1 // null terminated string

	length := sizeMessageLength +
		sizeProcessID +
		sizeChannel +
		sizePayload

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeNotificationResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.ProcessID)
	buf.AppendString(x.Channel)
	buf.AppendString(x.Payload)
	return buf.Bytes(), nil
}

func (x *msgNotificationResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeNotificationResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	processID, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	channel, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	payload, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	x.ProcessID = processID
	x.Channel = channel
	x.Payload = payload
	return nil
}

type msgParameterDescription struct {
	Parameters []int32
}

func (x *msgParameterDescription) message() {}

func (x *msgParameterDescription) backend() {}

func (x *msgParameterDescription) AppendBinary(b []byte) ([]byte, error) {
	const sizeParameterCount = 2
	const sizeParameter = 4

	countParameter := len(x.Parameters)

	length := sizeMessageLength +
		sizeParameterCount +
		sizeParameter*countParameter

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}
	buf.Grow(size)
	buf.AppendByte(typeParameterDescription)
	buf.AppendInt32(int32(length))
	buf.AppendInt16(int16(countParameter))

	for i := range countParameter {
		buf.AppendInt32(x.Parameters[i])
	}
	return buf.Bytes(), nil
}

func (x *msgParameterDescription) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeParameterDescription, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	countParameters, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	parameters := make([]int32, 0, countParameters)
	for range countParameters {
		parameter, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}
		parameters = append(parameters, parameter)
	}
	x.Parameters = parameters
	return nil
}

type msgParameterStatus struct {
	Name  string
	Value string
}

func (x *msgParameterStatus) message() {}

func (x *msgParameterStatus) backend() {}

func (x *msgParameterStatus) AppendBinary(b []byte) ([]byte, error) {
	sizeName := len(x.Name) + 1   // null terminated string
	sizeValue := len(x.Value) + 1 // null terminated string

	length := sizeMessageLength +
		sizeName +
		sizeValue

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeParameterStatus)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Name)
	buf.AppendString(x.Value)
	return buf.Bytes(), nil
}

func (x *msgParameterStatus) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeParameterStatus, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	name, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	value, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}

	x.Name = name
	x.Value = value
	return nil
}

type msgParseComplete struct{}

func (x *msgParseComplete) message() {}

func (x *msgParseComplete) backend() {}

func (x *msgParseComplete) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeParseComplete)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgParseComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeParseComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgPortalSuspended struct{}

func (x *msgPortalSuspended) message() {}

func (x *msgPortalSuspended) backend() {}

func (x *msgPortalSuspended) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typePortalSuspend)
	buf.AppendInt32(int32(length))
	return buf.Bytes(), nil
}

func (x *msgPortalSuspended) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typePortalSuspend, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}
	return nil
}

type msgReadyForQuery struct {
	TxStatus byte
}

func (x *msgReadyForQuery) message() {}

func (x *msgReadyForQuery) backend() {}

func (x *msgReadyForQuery) AppendBinary(b []byte) ([]byte, error) {
	const sizeTxStatus = 1
	const length = sizeMessageLength + sizeTxStatus
	const size = sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeReadyForQuery)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.TxStatus)
	return buf.Bytes(), nil
}

func (x *msgReadyForQuery) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeReadyForQuery, b)
	if err != nil {
		return invalidFormat(err)
	}

	status, b, err := shiftByte(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(ErrValueOverflow)
	}

	x.TxStatus = status
	return nil
}

type msgRowDescription struct {
	Names     []string
	Tables    []int32
	Columns   []int16
	DataTypes []int32
	Sizes     []int16
	Modifiers []int32
	Formats   []int16
}

func (x *msgRowDescription) message() {}

func (x *msgRowDescription) backend() {}

func (x *msgRowDescription) AppendBinary(b []byte) ([]byte, error) {
	const (
		sizeFieldCount = 2
		sizeTable      = 4
		sizeColumn     = 2
		sizeDataType   = 4
		sizeSize       = 2
		sizeModifier   = 4
		sizeFormat     = 2

		sizeRowInvariant = sizeTable +
			sizeColumn +
			sizeDataType +
			sizeSize +
			sizeModifier +
			sizeFormat
	)

	countFields := len(x.Names)

	if countFields > math.MaxInt16 {
		return b, invalidFormat(ErrValueOverflow)
	}

	var sizeRows int

	for i := range countFields {
		sizeName := len(x.Names[i]) + 1 // null terminated string
		sizeRows += (sizeName + sizeRowInvariant)
	}

	length := sizeMessageLength + sizeFieldCount + sizeRows

	if length > math.MaxInt32 {
		return b, invalidFormat(ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := buffer{b}

	buf.Grow(size)
	buf.AppendByte(typeRowDescription)
	buf.AppendInt32(int32(length))
	buf.AppendInt16(int16(countFields))

	for i := range countFields {
		buf.AppendString(x.Names[i])
		buf.AppendInt32(x.Tables[i])
		buf.AppendInt16(x.Columns[i])
		buf.AppendInt32(x.DataTypes[i])
		buf.AppendInt16(x.Sizes[i])
		buf.AppendInt32(x.Modifiers[i])
		buf.AppendInt16(x.Formats[i])
	}
	return buf.Bytes(), nil
}

func (x *msgRowDescription) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(typeRowDescription, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := buffer{b}

	countFields, err := buf.ShiftInt16()
	if err != nil {
		return invalidFormat(err)
	}

	names := make([]string, 0, countFields)
	tables := make([]int32, 0, countFields)
	columns := make([]int16, 0, countFields)
	dataTypes := make([]int32, 0, countFields)
	sizes := make([]int16, 0, countFields)
	modifiers := make([]int32, 0, countFields)
	formats := make([]int16, 0, countFields)

	for range countFields {
		name, err := buf.ShiftString()
		if err != nil {
			return invalidFormat(err)
		}
		names = append(names, name)

		table, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}
		tables = append(tables, table)

		column, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		columns = append(columns, column)

		dataType, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}
		dataTypes = append(dataTypes, dataType)

		sz, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		sizes = append(sizes, sz)

		modifier, err := buf.ShiftInt32()
		if err != nil {
			return invalidFormat(err)
		}
		modifiers = append(modifiers, modifier)

		format, err := buf.ShiftInt16()
		if err != nil {
			return invalidFormat(err)
		}
		formats = append(formats, format)
	}

	x.Names = names
	x.Tables = tables
	x.Columns = columns
	x.DataTypes = dataTypes
	x.Sizes = sizes
	x.Modifiers = modifiers
	x.Formats = formats
	return nil
}
