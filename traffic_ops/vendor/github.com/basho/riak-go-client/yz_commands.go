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

	rpbRiakSCH "github.com/basho/riak-go-client/rpb/riak_search"
	rpbRiakYZ "github.com/basho/riak-go-client/rpb/riak_yokozuna"
	proto "github.com/golang/protobuf/proto"
)

// SearchIndex object representing the solr index that is returned from Riak
type SearchIndex struct {
	Name   string
	Schema string
	NVal   uint32
}

// Schema object representing the solr schema that is returned from Riak
type Schema struct {
	Name    string
	Content string
}

// StoreIndex
// RpbYokozunaIndexPutReq
// RpbPutResp

// StoreIndexCommand is sused to store a new search index on Riak
type StoreIndexCommand struct {
	commandImpl
	timeoutImpl
	retryableCommandImpl
	Response bool
	protobuf *rpbRiakYZ.RpbYokozunaIndexPutReq
}

// Name identifies this command
func (cmd *StoreIndexCommand) Name() string {
	return cmd.getName("StoreIndex")
}

func (cmd *StoreIndexCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *StoreIndexCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = true
	return nil
}

func (cmd *StoreIndexCommand) getRequestCode() byte {
	return rpbCode_RpbYokozunaIndexPutReq
}

func (cmd *StoreIndexCommand) getResponseCode() byte {
	return rpbCode_RpbPutResp
}

func (cmd *StoreIndexCommand) getResponseProtobufMessage() proto.Message {
	return nil
}

// StoreIndexCommandBuilder type is required for creating new instances of StoreIndexCommand
//
//	command, err := NewStoreIndexCommandBuilder().
//		WithIndexName("myIndexName").
//		WithSchemaName("mySchemaName").
//		Build()
type StoreIndexCommandBuilder struct {
	timeout  time.Duration
	protobuf *rpbRiakYZ.RpbYokozunaIndexPutReq
}

// NewStoreIndexCommandBuilder is a factory function for generating the command builder struct
func NewStoreIndexCommandBuilder() *StoreIndexCommandBuilder {
	protobuf := &rpbRiakYZ.RpbYokozunaIndexPutReq{
		Index: &rpbRiakYZ.RpbYokozunaIndex{},
	}
	builder := &StoreIndexCommandBuilder{protobuf: protobuf}
	return builder
}

// WithIndexName sets the index to use for the command
func (builder *StoreIndexCommandBuilder) WithIndexName(indexName string) *StoreIndexCommandBuilder {
	builder.protobuf.Index.Name = []byte(indexName)
	return builder
}

// WithSchemaName sets the schema that the command will use
func (builder *StoreIndexCommandBuilder) WithSchemaName(schemaName string) *StoreIndexCommandBuilder {
	builder.protobuf.Index.Schema = []byte(schemaName)
	return builder
}

// WithNVal sets the number of times this command operation is replicated in the Cluster. If
// ommitted, the ring default is used.
//
// See http://basho.com/posts/technical/riaks-config-behaviors-part-2/
func (builder *StoreIndexCommandBuilder) WithNVal(nval uint32) *StoreIndexCommandBuilder {
	builder.protobuf.Index.NVal = &nval
	return builder
}

// WithTimeout sets a timeout in milliseconds to be used for this command operation
func (builder *StoreIndexCommandBuilder) WithTimeout(timeout time.Duration) *StoreIndexCommandBuilder {
	timeoutMilliseconds := uint32(timeout / time.Millisecond)
	builder.timeout = timeout
	builder.protobuf.Timeout = &timeoutMilliseconds
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *StoreIndexCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	return &StoreIndexCommand{
		timeoutImpl: timeoutImpl{
			timeout: builder.timeout,
		},
		protobuf: builder.protobuf,
	}, nil
}

// FetchIndex
// RpbYokozunaIndexGetReq
// RpbYokozunaIndexGetResp

// FetchIndexCommand is used to fetch a search index from Riak
type FetchIndexCommand struct {
	commandImpl
	retryableCommandImpl
	Response []*SearchIndex
	protobuf *rpbRiakYZ.RpbYokozunaIndexGetReq
}

// Name identifies this command
func (cmd *FetchIndexCommand) Name() string {
	return cmd.getName("FetchIndex")
}

func (cmd *FetchIndexCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchIndexCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbYokozunaIndexGetResp, ok := msg.(*rpbRiakYZ.RpbYokozunaIndexGetResp); ok {
			rpbIndexes := rpbYokozunaIndexGetResp.GetIndex()
			if rpbIndexes != nil {
				cmd.Response = make([]*SearchIndex, len(rpbIndexes))
				for i, rpbIndex := range rpbIndexes {
					index := &SearchIndex{
						Name:   string(rpbIndex.GetName()),
						Schema: string(rpbIndex.GetSchema()),
						NVal:   rpbIndex.GetNVal(),
					}
					cmd.Response[i] = index
				}
			}
		} else {
			return fmt.Errorf("[FetchIndexCommand] could not convert %v to RpbYokozunaIndexGetResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchIndexCommand) getRequestCode() byte {
	return rpbCode_RpbYokozunaIndexGetReq
}

func (cmd *FetchIndexCommand) getResponseCode() byte {
	return rpbCode_RpbYokozunaIndexGetResp
}

func (cmd *FetchIndexCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakYZ.RpbYokozunaIndexGetResp{}
}

// FetchIndexCommandBuilder type is required for creating new instances of FetchIndexCommand
//
//	command, err := NewFetchIndexCommandBuilder().
//		WithIndexName("myIndexName").
//		Build()
type FetchIndexCommandBuilder struct {
	protobuf *rpbRiakYZ.RpbYokozunaIndexGetReq
}

// NewFetchIndexCommandBuilder is a factory function for generating the command builder struct
func NewFetchIndexCommandBuilder() *FetchIndexCommandBuilder {
	builder := &FetchIndexCommandBuilder{protobuf: &rpbRiakYZ.RpbYokozunaIndexGetReq{}}
	return builder
}

// WithIndexName sets the index to use for the command
func (builder *FetchIndexCommandBuilder) WithIndexName(indexName string) *FetchIndexCommandBuilder {
	builder.protobuf.Name = []byte(indexName)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchIndexCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	return &FetchIndexCommand{protobuf: builder.protobuf}, nil
}

// DeleteIndex
// RpbYokozunaIndexDeleteReq
// RpbDelResp

// DeleteIndexCommand is used to delete a search index from Riak
type DeleteIndexCommand struct {
	commandImpl
	retryableCommandImpl
	Response bool
	protobuf *rpbRiakYZ.RpbYokozunaIndexDeleteReq
}

// Name identifies this command
func (cmd *DeleteIndexCommand) Name() string {
	return cmd.getName("DeleteIndex")
}

func (cmd *DeleteIndexCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *DeleteIndexCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = true
	return nil
}

func (cmd *DeleteIndexCommand) getRequestCode() byte {
	return rpbCode_RpbYokozunaIndexDeleteReq
}

func (cmd *DeleteIndexCommand) getResponseCode() byte {
	return rpbCode_RpbDelResp
}

func (cmd *DeleteIndexCommand) getResponseProtobufMessage() proto.Message {
	return nil
}

// DeleteIndexCommandBuilder type is required for creating new instances of DeleteIndexCommand
//
//	command, err := NewDeleteIndexCommandBuilder().
//		WithIndexName("myIndexName").
//		Build()
type DeleteIndexCommandBuilder struct {
	protobuf *rpbRiakYZ.RpbYokozunaIndexDeleteReq
}

// NewDeleteIndexCommandBuilder is a factory function for generating the command builder struct
func NewDeleteIndexCommandBuilder() *DeleteIndexCommandBuilder {
	builder := &DeleteIndexCommandBuilder{protobuf: &rpbRiakYZ.RpbYokozunaIndexDeleteReq{}}
	return builder
}

// WithIndexName sets the index to use for the command
func (builder *DeleteIndexCommandBuilder) WithIndexName(indexName string) *DeleteIndexCommandBuilder {
	builder.protobuf.Name = []byte(indexName)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *DeleteIndexCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	return &DeleteIndexCommand{protobuf: builder.protobuf}, nil
}

// StoreSchema
// RpbYokozunaSchemaPutReq
// RpbPutResp

// StoreSchemaCommand is used to store / update a search schema in Riak
type StoreSchemaCommand struct {
	commandImpl
	retryableCommandImpl
	Response bool
	protobuf *rpbRiakYZ.RpbYokozunaSchemaPutReq
}

// Name identifies this command
func (cmd *StoreSchemaCommand) Name() string {
	return cmd.getName("StoreSchema")
}

func (cmd *StoreSchemaCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *StoreSchemaCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = true
	return nil
}

func (cmd *StoreSchemaCommand) getRequestCode() byte {
	return rpbCode_RpbYokozunaSchemaPutReq
}

func (cmd *StoreSchemaCommand) getResponseCode() byte {
	return rpbCode_RpbPutResp
}

func (cmd *StoreSchemaCommand) getResponseProtobufMessage() proto.Message {
	return nil
}

// StoreSchemaCommandBuilder type is required for creating new instances of StoreSchemaCommand
//
//	command, err := NewStoreSchemaCommandBuilder().
//		WithSchemaName("mySchemaName").
//		WithSchema("mySchemaXML").
//		Build()
type StoreSchemaCommandBuilder struct {
	protobuf *rpbRiakYZ.RpbYokozunaSchemaPutReq
}

// NewStoreSchemaCommandBuilder is a factory function for generating the command builder struct
func NewStoreSchemaCommandBuilder() *StoreSchemaCommandBuilder {
	protobuf := &rpbRiakYZ.RpbYokozunaSchemaPutReq{
		Schema: &rpbRiakYZ.RpbYokozunaSchema{},
	}
	builder := &StoreSchemaCommandBuilder{protobuf: protobuf}
	return builder
}

// WithSchemaName sets the name for the schema to be stored
func (builder *StoreSchemaCommandBuilder) WithSchemaName(schemaName string) *StoreSchemaCommandBuilder {
	builder.protobuf.Schema.Name = []byte(schemaName)
	return builder
}

// WithSchema sets the actual schema that solr will use for indexing and queries
func (builder *StoreSchemaCommandBuilder) WithSchema(schema string) *StoreSchemaCommandBuilder {
	builder.protobuf.Schema.Content = []byte(schema)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *StoreSchemaCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	return &StoreSchemaCommand{protobuf: builder.protobuf}, nil
}

// FetchSchema
// RpbYokozunaSchemaPutReq
// RpbPutResp

// FetchSchemaCommand is used to GET a search schema from Riak
type FetchSchemaCommand struct {
	commandImpl
	retryableCommandImpl
	Response *Schema
	protobuf *rpbRiakYZ.RpbYokozunaSchemaGetReq
}

// Name identifies this command
func (cmd *FetchSchemaCommand) Name() string {
	return cmd.getName("FetchSchema")
}

func (cmd *FetchSchemaCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *FetchSchemaCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	if msg != nil {
		if rpbYokozunaSchemaGetResp, ok := msg.(*rpbRiakYZ.RpbYokozunaSchemaGetResp); ok {
			rpbSchema := rpbYokozunaSchemaGetResp.GetSchema()
			if rpbSchema != nil {
				cmd.Response = &Schema{
					Name:    string(rpbSchema.GetName()),
					Content: string(rpbSchema.GetContent()),
				}
			}
		} else {
			return fmt.Errorf("[FetchSchemaCommand] could not convert %v to RpbYokozunaSchemaGetResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *FetchSchemaCommand) getRequestCode() byte {
	return rpbCode_RpbYokozunaSchemaGetReq
}

func (cmd *FetchSchemaCommand) getResponseCode() byte {
	return rpbCode_RpbYokozunaSchemaGetResp
}

func (cmd *FetchSchemaCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakYZ.RpbYokozunaSchemaGetResp{}
}

// FetchSchemaCommandBuilder type is required for creating new instances of FetchSchemaCommand
//
//	command, err := NewFetchSchemaCommandBuilder().
//		WithSchemaName("mySchemaName").
//		Build()
type FetchSchemaCommandBuilder struct {
	protobuf *rpbRiakYZ.RpbYokozunaSchemaGetReq
}

// NewFetchSchemaCommandBuilder is a factory function for generating the command builder struct
func NewFetchSchemaCommandBuilder() *FetchSchemaCommandBuilder {
	builder := &FetchSchemaCommandBuilder{protobuf: &rpbRiakYZ.RpbYokozunaSchemaGetReq{}}
	return builder
}

// WithSchemaName sets the schema that the command will use
func (builder *FetchSchemaCommandBuilder) WithSchemaName(schemaName string) *FetchSchemaCommandBuilder {
	builder.protobuf.Name = []byte(schemaName)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *FetchSchemaCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	return &FetchSchemaCommand{protobuf: builder.protobuf}, nil
}

// Search
// RpbSearchQueryReq
// RpbSearchQueryResp

// SearchCommand is used to search Riak for values using search indexes & schemas
type SearchCommand struct {
	commandImpl
	Response *SearchResponse
	protobuf *rpbRiakSCH.RpbSearchQueryReq
}

// Name identifies this command
func (cmd *SearchCommand) Name() string {
	return cmd.getName("Search")
}

func (cmd *SearchCommand) constructPbRequest() (proto.Message, error) {
	return cmd.protobuf, nil
}

func (cmd *SearchCommand) onSuccess(msg proto.Message) error {
	cmd.success = true
	cmd.Response = &SearchResponse{}
	if msg != nil {
		if rpbResp, ok := msg.(*rpbRiakSCH.RpbSearchQueryResp); ok {
			resp := &SearchResponse{
				MaxScore: rpbResp.GetMaxScore(),
				NumFound: rpbResp.GetNumFound(),
			}
			rpbDocs := rpbResp.GetDocs()
			if rpbDocs != nil {
				resp.Docs = make([]*SearchDoc, len(rpbDocs))
				for i, rpbDoc := range rpbDocs {
					doc := &SearchDoc{}
					rpbFields := rpbDoc.GetFields()
					if rpbFields != nil {
						doc.Fields = make(map[string][]string)
						for _, rpbDocFld := range rpbFields {
							key := string(rpbDocFld.Key)
							val := string(rpbDocFld.Value)

							if doc.Fields[key] == nil {
								doc.Fields[key] = make([]string, 1)
								doc.Fields[key][0] = val
							} else {
								doc.Fields[key] = append(doc.Fields[key], val)
							}

							switch key {
							case yzBucketTypeFld:
								doc.BucketType = val
							case yzBucketFld:
								doc.Bucket = val
							case yzKeyFld:
								doc.Key = val
							case yzIdFld:
								doc.Id = val
							case yzScoreFld:
								doc.Score = val
							}
						}
					}
					resp.Docs[i] = doc
				}
			}
			cmd.Response = resp
		} else {
			return fmt.Errorf("[SearchCommand] could not convert %v to RpbSearchQueryResp", reflect.TypeOf(msg))
		}
	}
	return nil
}

func (cmd *SearchCommand) getRequestCode() byte {
	return rpbCode_RpbSearchQueryReq
}

func (cmd *SearchCommand) getResponseCode() byte {
	return rpbCode_RpbSearchQueryResp
}

func (cmd *SearchCommand) getResponseProtobufMessage() proto.Message {
	return &rpbRiakSCH.RpbSearchQueryResp{}
}

const yzBucketTypeFld = "_yz_rt"
const yzBucketFld = "_yz_rb"
const yzKeyFld = "_yz_rk"
const yzIdFld = "_yz_id"
const yzScoreFld = "score"

// SearchDoc object representing solr document returned from Riak
type SearchDoc struct {
	BucketType string
	Bucket     string
	Key        string
	Id         string
	Score      string
	Fields     map[string][]string
}

// SearchResponse contains the response data for a SearchCommand
type SearchResponse struct {
	Docs     []*SearchDoc
	MaxScore float32
	NumFound uint32
}

// SearchCommandBuilder type is required for creating new instances of SearchCommand
//
//	command, err := NewSearchCommandBuilder().
//		WithIndexName("myIndexName").
//		WithQuery("mySolrQuery").
//		WithNumRows(100).
//		Build()
type SearchCommandBuilder struct {
	protobuf *rpbRiakSCH.RpbSearchQueryReq
}

// NewSearchCommandBuilder is a factory function for generating the command builder struct
func NewSearchCommandBuilder() *SearchCommandBuilder {
	builder := &SearchCommandBuilder{protobuf: &rpbRiakSCH.RpbSearchQueryReq{}}
	return builder
}

// WithIndexName sets the index to use for the command
func (builder *SearchCommandBuilder) WithIndexName(index string) *SearchCommandBuilder {
	builder.protobuf.Index = []byte(index)
	return builder
}

// WithQuery sets the solr query to be executed on Riak
func (builder *SearchCommandBuilder) WithQuery(query string) *SearchCommandBuilder {
	builder.protobuf.Q = []byte(query)
	return builder
}

// WithNumRows sets the number of documents to be returned by Riak
func (builder *SearchCommandBuilder) WithNumRows(numRows uint32) *SearchCommandBuilder {
	builder.protobuf.Rows = &numRows
	return builder
}

// WithStart sets the document to start the result set with
func (builder *SearchCommandBuilder) WithStart(start uint32) *SearchCommandBuilder {
	builder.protobuf.Start = &start
	return builder
}

// WithSortField defines which field should be used for sorting the result set
func (builder *SearchCommandBuilder) WithSortField(sortField string) *SearchCommandBuilder {
	builder.protobuf.Sort = []byte(sortField)
	return builder
}

// WithFilterQuery sets the solr filter query to be used, the main query runs first, the filter
// query reduces the scope of the result set even further
func (builder *SearchCommandBuilder) WithFilterQuery(filterQuery string) *SearchCommandBuilder {
	builder.protobuf.Filter = []byte(filterQuery)
	return builder
}

// WithDefaultField sets the default field to be used by Riak the search query
//
// See https://wiki.apache.org/solr/SolrQuerySyntax
func (builder *SearchCommandBuilder) WithDefaultField(defaultField string) *SearchCommandBuilder {
	builder.protobuf.Df = []byte(defaultField)
	return builder
}

// WithDefaultOperation sets the default operation to be used by Riak for the search query
//
// See https://wiki.apache.org/solr/SolrQuerySyntax
func (builder *SearchCommandBuilder) WithDefaultOperation(op string) *SearchCommandBuilder {
	builder.protobuf.Op = []byte(op)
	return builder
}

// WithReturnFields sets the fields to be returned within each document
func (builder *SearchCommandBuilder) WithReturnFields(fields ...string) *SearchCommandBuilder {
	builder.protobuf.Fl = make([][]byte, len(fields))
	for i, f := range fields {
		builder.protobuf.Fl[i] = []byte(f)
	}
	return builder
}

// WithPresort allows you to configure Riak to presort the result set by Key or Score
func (builder *SearchCommandBuilder) WithPresort(presort string) *SearchCommandBuilder {
	builder.protobuf.Presort = []byte(presort)
	return builder
}

// Build validates the configuration options provided then builds the command
func (builder *SearchCommandBuilder) Build() (Command, error) {
	if builder.protobuf == nil {
		panic("builder.protobuf must not be nil")
	}
	return &SearchCommand{protobuf: builder.protobuf}, nil
}
