// Copyright 2015-present Basho Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package riak

import (
	"fmt"
	"reflect"
	"time"

	"github.com/basho/riak-go-client/rpb/riak_ts"

	"github.com/golang/protobuf/proto"
)

// TsCell represents a cell value within a time series row
type TsCell struct {
	columnType riak_ts.TsColumnType
	cell       *riak_ts.TsCell
}

// GetDataType returns the data type of the value stored within the cell
func (c *TsCell) GetDataType() string {
	var dType string
	switch {
	case c.cell.VarcharValue != nil:
		switch c.columnType {
		case riak_ts.TsColumnType_VARCHAR:
			dType = riak_ts.TsColumnType_VARCHAR.String()
		case riak_ts.TsColumnType_BLOB:
			dType = riak_ts.TsColumnType_BLOB.String()
		default:
			dType = riak_ts.TsColumnType_VARCHAR.String()
		}
	case c.cell.Sint64Value != nil:
		dType = riak_ts.TsColumnType_SINT64.String()
	case c.cell.TimestampValue != nil:
		dType = riak_ts.TsColumnType_TIMESTAMP.String()
	case c.cell.BooleanValue != nil:
		dType = riak_ts.TsColumnType_BOOLEAN.String()
	case c.cell.DoubleValue != nil:
		dType = riak_ts.TsColumnType_DOUBLE.String()
	}

	return dType
}

// GetStringValue returns the string value stored within the cell
func (c *TsCell) GetStringValue() string {
	return string(c.cell.GetVarcharValue())
}

// GetBlobValue returns the blob value stored within the cell
func (c *TsCell) GetBlobValue() []byte {
	return c.cell.GetVarcharValue()
}

// GetBooleanValue returns the boolean value stored within the cell
func (c *TsCell) GetBooleanValue() bool {
	return c.cell.GetBooleanValue()
}

// GetDoubleValue returns the double value stored within the cell
func (c *TsCell) GetDoubleValue() float64 {
	return c.cell.GetDoubleValue()
}

// GetSint64Value returns the sint64 value stored within the cell
func (c *TsCell) GetSint64Value() int64 {
	return c.cell.GetSint64Value()
}

// GetTimeValue returns the timestamp value stored within the cell as a time.Time
func (c *TsCell) GetTimeValue() time.Time {
	ts := c.cell.GetTimestampValue()
	s := ts / int64(1000)
	ms := time.Duration(ts%int64(1000)) * time.Millisecond
	return time.Unix(s, ms.Nanoseconds())
}

// GetTimestampValue returns the timestamp value stored within the cell
func (c *TsCell) GetTimestampValue() int64 {
	return c.cell.GetTimestampValue()
}

func (c *TsCell) setCell(tsc *riak_ts.TsCell, tsct riak_ts.TsColumnType) {
	c.cell = tsc
	c.columnType = tsct
}

// NewStringTsCell creates a TsCell from a string
func NewStringTsCell(v string) TsCell {
	tsc := riak_ts.TsCell{VarcharValue: []byte(v)}
	tsct := riak_ts.TsColumnType_VARCHAR
	return TsCell{columnType: tsct, cell: &tsc}
}

// NewBlobTsCell creates a TsCell from a []byte
func NewBlobTsCell(v []byte) TsCell {
	tsc := riak_ts.TsCell{VarcharValue: v}
	tsct := riak_ts.TsColumnType_BLOB
	return TsCell{columnType: tsct, cell: &tsc}
}

// NewBooleanTsCell creates a TsCell from a boolean
func NewBooleanTsCell(v bool) TsCell {
	tsc := riak_ts.TsCell{BooleanValue: &v}
	tsct := riak_ts.TsColumnType_BOOLEAN
	return TsCell{columnType: tsct, cell: &tsc}
}

// NewDoubleTsCell creates a TsCell from an floating point number
func NewDoubleTsCell(v float64) TsCell {
	tsc := riak_ts.TsCell{DoubleValue: &v}
	tsct := riak_ts.TsColumnType_DOUBLE
	return TsCell{columnType: tsct, cell: &tsc}
}

// NewSint64TsCell creates a TsCell from an integer
func NewSint64TsCell(v int64) TsCell {
	tsc := riak_ts.TsCell{Sint64Value: &v}
	tsct := riak_ts.TsColumnType_SINT64
	return TsCell{columnType: tsct, cell: &tsc}
}

// NewTimestampTsCell creates a TsCell from a time.Time struct
func NewTimestampTsCell(t time.Time) TsCell {
	v := ToUnixMillis(t)
	tsc := riak_ts.TsCell{TimestampValue: &v}
	tsct := riak_ts.TsColumnType_TIMESTAMP
	return TsCell{columnType: tsct, cell: &tsc}
}

// NewTimestampTsCellFromInt64 creates a TsCell from an int64 value
// that represents *milliseconds* since UTC epoch
func NewTimestampTsCellFromInt64(v int64) TsCell {
	tsc := riak_ts.TsCell{TimestampValue: &v}
	tsct := riak_ts.TsColumnType_TIMESTAMP
	return TsCell{columnType: tsct, cell: &tsc}
}

// ToUnixMillis converts a time.Time to Unix milliseconds since UTC epoch
func ToUnixMillis(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// TsColumnDescription describes a Time Series column
type TsColumnDescription struct {
	column *riak_ts.TsColumnDescription
}

// GetName returns the name for the column
func (c *TsColumnDescription) GetName() string {
	return string(c.column.GetName())
}

// GetType returns the data type for the column
func (c *TsColumnDescription) GetType() string {
	return riak_ts.TsColumnType_name[int32(c.column.GetType())]
}

func (c *TsColumnDescription) setColumn(tsCol *riak_ts.TsColumnDescription) {
	c.column = tsCol
}

// TsStoreRows
// TsPutReq
// TsPutResp

// TsStoreRowsCommand is sused to store a new row/s in Riak TS
type TsStoreRowsCommand struct {
	commandImpl
	retryableCommandImpl
	Response bool
	protobuf *riak_ts.TsPutReq
}

// Name identifies this command
func (cmd *TsStoreRowsCommand) Name() string {
	return cmd.getName("TsStoreRows")
}

func (cmd *TsStoreRowsCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *TsStoreRowsCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = true
	return nil
}

func (cmd *TsStoreRowsCommand) getRequestCode() byte {
	return rpbCode_TsPutReq
}

func (cmd *TsStoreRowsCommand) getResponseCode() byte {
	return rpbCode_TsPutResp
}

func (cmd *TsStoreRowsCommand) getResponseProtobufMessage() proto.Message {
	return nil
}

// TsStoreRowsCommandBuilder type is required for creating new instances of StoreIndexCommand
//
//	cmd, err := NewTsStoreRowsCommandBuilder().
//		WithTable("myTableName").
//		WithRows(rows).
//		Build()
type TsStoreRowsCommandBuilder struct {
	protobuf *riak_ts.TsPutReq
}

// NewTsStoreRowsCommandBuilder is a factory function for generating the command builder struct
func NewTsStoreRowsCommandBuilder() *TsStoreRowsCommandBuilder {
	return &TsStoreRowsCommandBuilder{protobuf: &riak_ts.TsPutReq{}}
}

// WithTable sets the table to use for the command
func (builder *TsStoreRowsCommandBuilder) WithTable(table string) *TsStoreRowsCommandBuilder {
	builder.protobuf.Table = []byte(table)
	return builder
}

// WithRows sets the rows to be stored by the command
func (builder *TsStoreRowsCommandBuilder) WithRows(rows [][]TsCell) *TsStoreRowsCommandBuilder {
	builder.protobuf.Rows = convertFromTsRows(rows)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *TsStoreRowsCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}

	if len(builder.protobuf.GetTable()) == 0 {
		return nil, ErrTableRequired
	}

	return &TsStoreRowsCommand{
		protobuf: builder.protobuf,
	}, nil
}

// TsFetchRowCommand is used to fetch / get a value from Riak KV
type TsFetchRowCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response *TsFetchRowResponse
	protobuf *riak_ts.TsGetReq
}

// Name identifies this command
func (cmd *TsFetchRowCommand) Name() string {
	return cmd.getName("TsFetchRow")
}

func (cmd *TsFetchRowCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *TsFetchRowCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	var col TsColumnDescription

	if msg == nil {
		cmd.Response = &TsFetchRowResponse{
			IsNotFound: true,
		}
	} else {
		if tsTsFetchRowResp, ok := msg.(*riak_ts.TsGetResp); ok {
			tsCols := tsTsFetchRowResp.GetColumns()
			tsRows := tsTsFetchRowResp.GetRows()
			if tsCols != nil && tsRows != nil {
				cmd.Response = &TsFetchRowResponse{
					IsNotFound: false,
					Columns:    make([]TsColumnDescription, 0),
					Row:        make([]TsCell, 0),
				}

				for _, tsCol := range tsCols {
					col.setColumn(tsCol)
					cmd.Response.Columns = append(cmd.Response.Columns, col)
				}

				// grab only the first row if any
				rows := convertFromPbTsRows(tsRows, tsCols)
				if len(rows) > 0 {
					cmd.Response.Row = rows[0]
				} else {
					return fmt.Errorf("[TsFetchRowCommand] could not retrieve row from %v to TsGetResp", reflect.TypeOf(msg))
				}
			}
		} else {
			return fmt.Errorf("[TsFetchRowCommand] could not convert %v to TsGetResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *TsFetchRowCommand) getRequestCode() byte {
	return rpbCode_TsGetReq
}

func (cmd *TsFetchRowCommand) getResponseCode() byte {
	return rpbCode_TsGetResp
}

func (cmd *TsFetchRowCommand) getResponseProtobufMessage() proto.Message {
	return &riak_ts.TsGetResp{}
}

// TsFetchRowResponse contains the response data for a TsFetchRowCommand
type TsFetchRowResponse struct {
	IsNotFound bool
	Columns    []TsColumnDescription
	Row        []TsCell
}

// TsFetchRowCommandBuilder type is required for creating new instances of TsFetchRowCommand
//
//	key := make([]riak.TsCell, 3)
//	key[0] = NewStringTsCell("South Atlantic")
//	key[1] = NewStringTsCell("South Carolina")
//	key[2] = NewTimestampTsCell(1420113600)
//
//	cmd, err := NewTsFetchRowCommandBuilder().
//		WithBucketType("myBucketType").
//		WithTable("myTable").
//		WithKey(key).
//		Build()
type TsFetchRowCommandBuilder struct {
	timeout  time.Duration
	protobuf *riak_ts.TsGetReq
}

// NewTsFetchRowCommandBuilder is a factory function for generating the command builder struct
func NewTsFetchRowCommandBuilder() *TsFetchRowCommandBuilder {
	builder := &TsFetchRowCommandBuilder{protobuf: &riak_ts.TsGetReq{}}
	return builder
}

// WithTable sets the table to be used by the command
func (builder *TsFetchRowCommandBuilder) WithTable(table string) *TsFetchRowCommandBuilder {
	builder.protobuf.Table = []byte(table)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *TsFetchRowCommandBuilder) WithKey(key []TsCell) *TsFetchRowCommandBuilder {
	tsKey := make([]*riak_ts.TsCell, len(key))

	for i, v := range key {
		tsKey[i] = v.cell
	}

	builder.protobuf.Key = tsKey

	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *TsFetchRowCommandBuilder) WithTimeout(timeout time.Duration) *TsFetchRowCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *TsFetchRowCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}

	if len(builder.protobuf.GetTable()) == 0 {
		return nil, ErrTableRequired
	}

	if len(builder.protobuf.GetKey()) == 0 {
		return nil, ErrKeyRequired
	}

	return &TsFetchRowCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// TsDeleteRowCommand is used to delete a value from Riak TS
type TsDeleteRowCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response bool
	protobuf *riak_ts.TsDelReq
}

// Name identifies this command
func (cmd *TsDeleteRowCommand) Name() string {
	return cmd.getName("TsDeleteRow")
}

func (cmd *TsDeleteRowCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *TsDeleteRowCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = true
	return nil
}

func (cmd *TsDeleteRowCommand) getRequestCode() byte {
	return rpbCode_TsDelReq
}

func (cmd *TsDeleteRowCommand) getResponseCode() byte {
	return rpbCode_TsDelResp
}

func (cmd *TsDeleteRowCommand) getResponseProtobufMessage() proto.Message {
	return &riak_ts.TsDelResp{}
}

// TsDeleteRowCommandBuilder type is required for creating new instances of TsDeleteRowCommand
//
//	cmd, err := NewTsDeleteRowCommandBuilder().
//		WithTable("myTable").
//		WithKey(key).
//		Build()
type TsDeleteRowCommandBuilder struct {
	timeout  time.Duration
	protobuf *riak_ts.TsDelReq
}

// NewTsDeleteRowCommandBuilder is a factory function for generating the command builder struct
func NewTsDeleteRowCommandBuilder() *TsDeleteRowCommandBuilder {
	builder := &TsDeleteRowCommandBuilder{protobuf: &riak_ts.TsDelReq{}}
	return builder
}

// WithTable sets the table to be used by the command
func (builder *TsDeleteRowCommandBuilder) WithTable(table string) *TsDeleteRowCommandBuilder {
	builder.protobuf.Table = []byte(table)
	return builder
}

// WithKey sets the key to be used by the command to read / write values
func (builder *TsDeleteRowCommandBuilder) WithKey(key []TsCell) *TsDeleteRowCommandBuilder {
	tsKey := make([]*riak_ts.TsCell, len(key))

	for i, v := range key {
		tsKey[i] = v.cell
	}

	builder.protobuf.Key = tsKey

	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *TsDeleteRowCommandBuilder) WithTimeout(timeout time.Duration) *TsDeleteRowCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *TsDeleteRowCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}

	if len(builder.protobuf.GetTable()) == 0 {
		return nil, ErrTableRequired
	}

	if len(builder.protobuf.GetKey()) == 0 {
		return nil, ErrKeyRequired
	}

	return &TsDeleteRowCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// TsQueryCommand is used to fetch / get a value from Riak TS
type TsQueryCommand struct {
	commandImpl
	Response *TsQueryResponse
	protobuf *riak_ts.TsQueryReq
	callback func([][]TsCell) error
	done     bool
}

// Name identifies this command
func (cmd *TsQueryCommand) Name() string {
	return cmd.getName("TsQuery")
}

func (cmd *TsQueryCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *TsQueryCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	var col TsColumnDescription

	if msg == nil {
		cmd.done = true
		cmd.Response = &TsQueryResponse{}
	} else {
		if queryResp, ok := msg.(*riak_ts.TsQueryResp); ok {
			if cmd.Response == nil {
				cmd.Response = &TsQueryResponse{}
			}
			response := cmd.Response

			cmd.done = queryResp.GetDone()

			tsCols := queryResp.GetColumns()
			tsRows := queryResp.GetRows()
			if tsCols != nil && tsRows != nil {
				if tsCols != nil && response.Columns == nil {
					response.Columns = make([]TsColumnDescription, 0)
				}
				if tsRows != nil && response.Rows == nil {
					response.Rows = make([][]TsCell, 0)
				}

				for _, tsCol := range tsCols {
					col.setColumn(tsCol)
					response.Columns = append(response.Columns, col)
				}

				rows := convertFromPbTsRows(tsRows, tsCols)

				if cmd.protobuf.GetStream() {
					if cmd.callback == nil {
						panic("[TsQueryCommand] requires a callback when streaming.")
					} else {
						if err := cmd.callback(rows); err != nil {
							cmd.Response = nil
							return err
						}
					}
				} else {
					response.Rows = append(response.Rows, rows...)
				}
			}
		} else {
			cmd.done = true
			return fmt.Errorf("[TsQueryCommand] could not convert %v to TsQueryResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *TsQueryCommand) getRequestCode() byte {
	return rpbCode_TsQueryReq
}

func (cmd *TsQueryCommand) getResponseCode() byte {
	return rpbCode_TsQueryResp
}

func (cmd *TsQueryCommand) getResponseProtobufMessage() proto.Message {
	return &riak_ts.TsQueryResp{}
}

// TsQueryResponse contains the response data for a TsQueryCommand
type TsQueryResponse struct {
	Columns []TsColumnDescription
	Rows    [][]TsCell
}

// TsQueryCommandBuilder type is required for creating new instances of TsQueryCommand
//
//	cmd, err := NewTsQueryCommandBuilder().
//		WithQuery("select * from GeoCheckin where time > 1234560 and time < 1234569 and region = 'South Atlantic'").
//		WithStreaming(true).
//		WithCallback(cb).
//		Build()
type TsQueryCommandBuilder struct {
	protobuf *riak_ts.TsQueryReq
	callback func(rows [][]TsCell) error
}

// NewTsQueryCommandBuilder is a factory function for generating the command builder struct
func NewTsQueryCommandBuilder() *TsQueryCommandBuilder {
	builder := &TsQueryCommandBuilder{protobuf: &riak_ts.TsQueryReq{}}
	return builder
}

// WithQuery sets the query to be used by the command
func (builder *TsQueryCommandBuilder) WithQuery(query string) *TsQueryCommandBuilder {
	builder.protobuf.Query = &riak_ts.TsInterpolation{Base: []byte(query)}
	return builder
}

// WithStreaming sets the command to provide a streamed response
//
// If true, a callback must be provided via WithCallback()
func (builder *TsQueryCommandBuilder) WithStreaming(streaming bool) *TsQueryCommandBuilder {
	builder.protobuf.Stream = &streaming
	return builder
}

// WithCallback sets the callback to be used when handling a streaming response
//
// Requires WithStreaming(true)
func (builder *TsQueryCommandBuilder) WithCallback(callback func([][]TsCell) error) *TsQueryCommandBuilder {
	builder.callback = callback
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *TsQueryCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}

	if len(builder.protobuf.GetQuery().GetBase()) == 0 {
		return nil, ErrQueryRequired
	}

	if builder.protobuf.GetStream() && builder.callback == nil {
		return nil, newClientError("TsQueryCommand requires a callback when streaming.", nil)
	}

	return &TsQueryCommand{
		protobuf: builder.protobuf,
		callback: builder.callback,
	}, nil
}

// TsListKeys
// TsListKeysReq
// TsListKeysResp

// TsListKeysCommand is used to fetch values from a table in Riak TS
type TsListKeysCommand struct {
	commandImpl
	timeoutImpl
	listingImpl
	Response  *TsListKeysResponse
	protobuf  *riak_ts.TsListKeysReq
	streaming bool
	callback  func(keys [][]TsCell) error
	done      bool
}

// Name identifies this command
func (cmd *TsListKeysCommand) Name() string {
	return cmd.getName("TsListKeys")
}

func (cmd *TsListKeysCommand) isDone() bool {
	// NB: TsListKeysReq is *always* streaming so no need to take
	// cmd.streaming into account here, unlike RpbListBucketsReq
	return cmd.done
}

func (cmd *TsListKeysCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *TsListKeysCommand) onSuccess(msg proto.Message) error {
	cmd.success = true

	if msg == nil {
		cmd.done = true
		cmd.Response = &TsListKeysResponse{}
	} else {
		if keysResp, ok := msg.(*riak_ts.TsListKeysResp); ok {
			if cmd.Response == nil {
				cmd.Response = &TsListKeysResponse{}
			}

			cmd.done = keysResp.GetDone()
			response := cmd.Response

			if keysResp.GetKeys() != nil && len(keysResp.GetKeys()) > 0 {
				rows := convertFromPbTsRows(keysResp.GetKeys(), nil)
				if cmd.streaming {
					if cmd.callback == nil {
						panic("[TsListKeysCommand] requires a callback when streaming.")
					} else {
						if err := cmd.callback(rows); err != nil {
							cmd.Response = nil
							return err
						}
					}
				} else {
					// append slice to slice
					response.Keys = append(response.Keys, rows...)
				}
			}
		} else {
			cmd.done = true
			return fmt.Errorf("[TsListKeysCommand] could not convert %v to TsListKeysResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *TsListKeysCommand) getRequestCode() byte {
	return rpbCode_TsListKeysReq
}

func (cmd *TsListKeysCommand) getResponseCode() byte {
	return rpbCode_TsListKeysResp
}

func (cmd *TsListKeysCommand) getResponseProtobufMessage() proto.Message {
	return &riak_ts.TsListKeysResp{}
}

// TsListKeysResponse contains the response data for a TsListKeysCommand
type TsListKeysResponse struct {
	Keys [][]TsCell
}

// TsListKeysCommandBuilder type is required for creating new instances of TsListKeysCommand
//
//	cmd, err := NewTsListKeysCommandBuilder().
//		WithTable("myTable").
//		WithStreaming(true).
//		WithCallback(cb).
//		Build()
type TsListKeysCommandBuilder struct {
	allowListing bool
	timeout      time.Duration
	protobuf     *riak_ts.TsListKeysReq
	streaming    bool
	callback     func(keys [][]TsCell) error
}

// NewTsListKeysCommandBuilder is a factory function for generating the command builder struct
func NewTsListKeysCommandBuilder() *TsListKeysCommandBuilder {
	builder := &TsListKeysCommandBuilder{protobuf: &riak_ts.TsListKeysReq{}}
	return builder
}

// WithAllowListing will allow this command to be built and execute
func (builder *TsListKeysCommandBuilder) WithAllowListing() *TsListKeysCommandBuilder {
	builder.allowListing = true
	return builder
}

// WithTable sets the table to be used by the command
func (builder *TsListKeysCommandBuilder) WithTable(table string) *TsListKeysCommandBuilder {
	builder.protobuf.Table = []byte(table)
	return builder
}

// WithStreaming sets the command to provide a streamed response
//
// If true, a callback must be provided via WithCallback()
func (builder *TsListKeysCommandBuilder) WithStreaming(streaming bool) *TsListKeysCommandBuilder {
	builder.streaming = streaming
	return builder
}

// WithCallback sets the callback to be used when handling a streaming response
//
// Requires WithStreaming(true)
func (builder *TsListKeysCommandBuilder) WithCallback(callback func([][]TsCell) error) *TsListKeysCommandBuilder {
	builder.callback = callback
	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *TsListKeysCommandBuilder) WithTimeout(timeout time.Duration) *TsListKeysCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *TsListKeysCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	if len(builder.protobuf.GetTable()) == 0 {
		return nil, ErrTableRequired
	}
	if builder.streaming && builder.callback == nil {
		return nil, newClientError("ListKeysCommand requires a callback when streaming.", nil)
	}
	if !builder.allowListing {
		return nil, ErrListingDisabled
	}
	return &TsListKeysCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf:  builder.protobuf,
		streaming: builder.streaming,
		callback:  builder.callback,
	}, nil
}

// Converts a slice of riak_ts.TsRow to a slice of .TsRows
func convertFromPbTsRows(tsRows []*riak_ts.TsRow, tsCols []*riak_ts.TsColumnDescription) [][]TsCell {
	var rows [][]TsCell
	var row []TsCell
	var cell TsCell

	for _, tsRow := range tsRows {
		row = make([]TsCell, 0)

		for i, tsCell := range tsRow.Cells {
			tsColumnType := riak_ts.TsColumnType_VARCHAR
			if tsCols != nil {
				tsColumnType = tsCols[i].GetType()
			}
			cell.setCell(tsCell, tsColumnType)
			row = append(row, cell)
		}

		if len(rows) < 1 {
			rows = make([][]TsCell, 0)
		}

		rows = append(rows, row)
	}

	return rows
}

// Converts a slice of .TsRows to a slice of riak_ts.TsRow
func convertFromTsRows(tsRows [][]TsCell) []*riak_ts.TsRow {
	var rows []*riak_ts.TsRow
	var cells []*riak_ts.TsCell
	for _, tsRow := range tsRows {
		cells = make([]*riak_ts.TsCell, 0)

		for _, tsCell := range tsRow {
			cells = append(cells, tsCell.cell)
		}

		if len(rows) < 1 {
			rows = make([]*riak_ts.TsRow, 0)
		}

		rows = append(rows, &riak_ts.TsRow{Cells: cells})
	}

	return rows
}
