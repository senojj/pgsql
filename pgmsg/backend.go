package pgmsg

import (
	"math"
	
	"github.com/senojj/pgsql/pgio"
)

type BackendKeyData struct {
	ProcessID int32
	SecretKey []byte
}

func (x *BackendKeyData) AppendBinary(b []byte) ([]byte, error) {
	const sizeProcessID = 4
	sizeSecretKey := len(x.SecretKey)

	if sizeSecretKey < 4 {
		return b, invalidFormat(pgio.ErrValueUnderflow)
	}

	if sizeSecretKey > 256 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	length := sizeMessageLength + sizeProcessID + sizeSecretKey

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeBackendKeyData)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.ProcessID)
	buf.AppendByte(x.SecretKey...)
	return buf.Bytes, nil
}

func (x *BackendKeyData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeBackendKeyData, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	processID, err := buf.ShiftInt32()
	if err != nil {
		return invalidFormat(err)
	}

	if buf.Len() < 4 {
		return invalidFormat(pgio.ErrValueUnderflow)
	}

	if buf.Len() > 256 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	secretKey := make([]byte, buf.Len())
	copy(secretKey, buf.Bytes)

	x.ProcessID = processID
	x.SecretKey = secretKey
	return nil
}

type BindComplete struct{}

func (x *BindComplete) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeBindComplete)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *BindComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeBindComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type CloseComplete struct{}

func (x *CloseComplete) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCloseComplete)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *CloseComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCloseComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type CommandComplete struct {
	Tag string
}

func (x *CommandComplete) AppendBinary(b []byte) ([]byte, error) {
	sizeTag := len(x.Tag) + 1 // null terminated string
	length := sizeMessageLength + sizeTag

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCommandComplete)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Tag)
	return buf.Bytes, nil
}

func (x *CommandComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCommandComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

	tag, err := buf.ShiftString()
	if err != nil {
		return invalidFormat(err)
	}
	x.Tag = tag
	return nil
}

type CopyInResponse struct {
	Format  int8
	Columns []int16
}

func (x *CopyInResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeFormat = 1
	const sizeColCount = 2
	const sizeCol = 2

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeCols := countCols * sizeCol

	length := sizeMessageLength +
		sizeFormat +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCopyInResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt8(x.Format)
	buf.AppendInt16(int16(countCols))
	buf.AppendInt16(x.Columns...)
	return buf.Bytes, nil
}

func (x *CopyInResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCopyInResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type CopyOutResponse struct {
	Format  int8
	Columns []int16
}

func (x *CopyOutResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeFormat = 1
	const sizeColCount = 2
	const sizeCol = 2

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeCols := countCols * sizeCol

	length := sizeMessageLength +
		sizeFormat +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCopyOutResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt8(x.Format)
	buf.AppendInt16(int16(countCols))
	buf.AppendInt16(x.Columns...)
	return buf.Bytes, nil
}

func (x *CopyOutResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCopyOutResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type CopyBothResponse struct {
	Format  int8
	Columns []int16
}

func (x *CopyBothResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeFormat = 1
	const sizeColCount = 2
	const sizeCol = 2

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeCols := countCols * sizeCol

	length := sizeMessageLength +
		sizeFormat +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeCopyBothResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt8(x.Format)
	buf.AppendInt16(int16(countCols))
	buf.AppendInt16(x.Columns...)
	return buf.Bytes, nil
}

func (x *CopyBothResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeCopyBothResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type DataRow struct {
	// First dimention is columns (can have 0 elements).
	// Second dimension is value data (can have 0 elements
	// or be nil).
	Columns [][]byte
}

func (x *DataRow) AppendBinary(b []byte) ([]byte, error) {
	const sizeColCount = 2
	const sizeColLength = 4

	countCols := len(x.Columns)

	if countCols > math.MaxInt16 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	sizeCols := 0

	for i := range countCols {
		lengthCol := len(x.Columns[i])

		if lengthCol > math.MaxInt32 {
			return b, invalidFormat(pgio.ErrValueOverflow)
		}
		sizeCols += sizeColLength + lengthCol
	}

	length := sizeMessageLength +
		sizeColCount +
		sizeCols

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeDataRow)
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
	return buf.Bytes, nil
}

func (x *DataRow) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeDataRow, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type EmptyQueryResponse struct{}

func (x *EmptyQueryResponse) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeEmptyQueryResponse)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *EmptyQueryResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeEmptyQueryResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type ErrorResponse struct {
	Fields []byte
	Values []string
}

func (x *ErrorResponse) AppendBinary(b []byte) ([]byte, error) {
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
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeErrorResponse)
	buf.AppendInt32(int32(length))

	for i := range countFields {
		buf.AppendByte(x.Fields[i])
		buf.AppendString(x.Values[i])
	}
	buf.AppendByte(0)
	return buf.Bytes, nil
}

func (x *ErrorResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeErrorResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type FunctionCallResponse struct {
	// Can be zero length or nil.
	Result []byte
}

func (x *FunctionCallResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeResultLength = 4
	sizeResult := len(x.Result)

	if sizeResult > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	length := sizeMessageLength +
		sizeResultLength +
		sizeResult

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeFunctionCallResponse)
	buf.AppendInt32(int32(length))

	if x.Result == nil {
		buf.AppendInt32(int32(-1))
	} else {
		buf.AppendInt32(int32(sizeResult))
	}
	buf.AppendByte(x.Result...)
	return buf.Bytes, nil
}

func (x *FunctionCallResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeFunctionCallResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	length, b, err := pgio.ShiftInt32(b)
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

type NegotiateProtocolVersion struct {
	MinorVersionSupported int32
	UnrecognizedOptions   []string
}

func (x *NegotiateProtocolVersion) AppendBinary(b []byte) ([]byte, error) {
	const sizeMinorVersion = 4
	const sizeUnrecognizedOptionCount = 4

	length := sizeMessageLength +
		sizeMinorVersion +
		sizeUnrecognizedOptionCount

	countUnrecognizedOptions := len(x.UnrecognizedOptions)

	if countUnrecognizedOptions > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	for _, option := range x.UnrecognizedOptions {
		length += len(option) + 1 // +1 null byte for string termination
	}

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeNegotiateProtocolVersion)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.MinorVersionSupported)
	buf.AppendInt32(int32(countUnrecognizedOptions))
	buf.AppendString(x.UnrecognizedOptions...)
	return buf.Bytes, nil
}

func (x *NegotiateProtocolVersion) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeNegotiateProtocolVersion, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type NoData struct{}

func (x *NoData) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeNoData)
	buf.AppendInt32(length)
	return buf.Bytes, nil
}

func (x *NoData) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeNoData, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type NoticeResponse struct {
	Fields []byte
	Values []string
}

func (x *NoticeResponse) AppendBinary(b []byte) ([]byte, error) {
	countFields := len(x.Fields)

	if countFields != len(x.Values) {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	length := sizeMessageLength + countFields

	for _, value := range x.Values {
		length += len(value) + 1 // null terminated strings
	}
	length += 1 // null terminated list

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeNoticeResponse)
	buf.AppendInt32(int32(length))

	for i := range countFields {
		buf.AppendByte(x.Fields[i])
		buf.AppendString(x.Values[i])
	}
	buf.AppendByte(0)
	return buf.Bytes, nil
}

func (x *NoticeResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeNoticeResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type NotificationResponse struct {
	ProcessID int32
	Channel   string
	Payload   string
}

func (x *NotificationResponse) AppendBinary(b []byte) ([]byte, error) {
	const sizeProcessID = 4

	sizeChannel := len(x.Channel) + 1 // null terminated string
	sizePayload := len(x.Payload) + 1 // null terminated string

	length := sizeMessageLength +
		sizeProcessID +
		sizeChannel +
		sizePayload

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeNotificationResponse)
	buf.AppendInt32(int32(length))
	buf.AppendInt32(x.ProcessID)
	buf.AppendString(x.Channel)
	buf.AppendString(x.Payload)
	return buf.Bytes, nil
}

func (x *NotificationResponse) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeNotificationResponse, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type ParameterDescription struct {
	Parameters []int32
}

func (x *ParameterDescription) AppendBinary(b []byte) ([]byte, error) {
	const sizeParameterCount = 2
	const sizeParameter = 4

	countParameter := len(x.Parameters)

	length := sizeMessageLength +
		sizeParameterCount +
		sizeParameter*countParameter

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}
	buf.Grow(size)
	buf.AppendByte(TypeParameterDescription)
	buf.AppendInt32(int32(length))
	buf.AppendInt16(int16(countParameter))

	for i := range countParameter {
		buf.AppendInt32(x.Parameters[i])
	}
	return buf.Bytes, nil
}

func (x *ParameterDescription) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeParameterDescription, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type ParameterStatus struct {
	Name  string
	Value string
}

func (x *ParameterStatus) AppendBinary(b []byte) ([]byte, error) {
	sizeName := len(x.Name) + 1   // null terminated string
	sizeValue := len(x.Value) + 1 // null terminated string

	length := sizeMessageLength +
		sizeName +
		sizeValue

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeParameterStatus)
	buf.AppendInt32(int32(length))
	buf.AppendString(x.Name)
	buf.AppendString(x.Value)
	return buf.Bytes, nil
}

func (x *ParameterStatus) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeParameterStatus, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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

type ParseComplete struct{}

func (x *ParseComplete) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeParseComplete)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *ParseComplete) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeParseComplete, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type PortalSuspended struct{}

func (x *PortalSuspended) AppendBinary(b []byte) ([]byte, error) {
	const length = sizeMessageLength
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypePortalSuspend)
	buf.AppendInt32(int32(length))
	return buf.Bytes, nil
}

func (x *PortalSuspended) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypePortalSuspend, b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}
	return nil
}

type ReadyForQuery struct {
	TxStatus byte
}

func (x *ReadyForQuery) AppendBinary(b []byte) ([]byte, error) {
	const sizeTxStatus = 1
	const length = sizeMessageLength + sizeTxStatus
	const size = sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeReadyForQuery)
	buf.AppendInt32(int32(length))
	buf.AppendByte(x.TxStatus)
	return buf.Bytes, nil
}

func (x *ReadyForQuery) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeReadyForQuery, b)
	if err != nil {
		return invalidFormat(err)
	}

	status, b, err := pgio.ShiftByte(b)
	if err != nil {
		return invalidFormat(err)
	}

	if len(b) > 0 {
		return invalidFormat(pgio.ErrValueOverflow)
	}

	x.TxStatus = status
	return nil
}

type RowDescription struct {
	Names     []string
	Tables    []int32
	Columns   []int16
	DataTypes []int32
	Sizes     []int16
	Modifiers []int32
	Formats   []int16
}

func (x *RowDescription) AppendBinary(b []byte) ([]byte, error) {
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
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	var sizeRows int

	for i := range countFields {
		sizeName := len(x.Names[i]) + 1 // null terminated string
		sizeRows += (sizeName + sizeRowInvariant)
	}

	length := sizeMessageLength + sizeFieldCount + sizeRows

	if length > math.MaxInt32 {
		return b, invalidFormat(pgio.ErrValueOverflow)
	}

	size := sizeMessageType + length

	buf := pgio.Buffer{Bytes: b}

	buf.Grow(size)
	buf.AppendByte(TypeRowDescription)
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
	return buf.Bytes, nil
}

func (x *RowDescription) UnmarshalBinary(b []byte) error {
	b, err := shiftHeader(TypeRowDescription, b)
	if err != nil {
		return invalidFormat(err)
	}

	buf := pgio.Buffer{Bytes: b}

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
