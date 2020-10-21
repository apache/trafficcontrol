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
	"time"

	rpbRiak "github.com/basho/riak-go-client/rpb/riak"
	rpbRiakKV "github.com/basho/riak-go-client/rpb/riak_kv"
)

// Link is used to represent a Riak KV object link, which is a one way link to another object within
// Riak
type Link struct {
	Bucket string
	Key    string
	Tag    string
}

// Pair is used to store user defined meta data with a key and value
type Pair struct {
	Key   string
	Value string
}

// Object structure used for representing a KV Riak object
type Object struct {
	BucketType      string
	Bucket          string
	Key             string
	IsTombstone     bool
	Value           []byte
	ContentType     string
	Charset         string
	ContentEncoding string
	VTag            string
	LastModified    time.Time
	UserMeta        []*Pair
	Indexes         map[string][]string
	Links           []*Link
	VClock          []byte
}

// HasIndexes is a bool check to determine if the object contains any secondary indexes for searching
func (o *Object) HasIndexes() bool {
	return len(o.Indexes) > 0
}

// HasUserMeta is a bool check to determine if the object contains any user defined meta data
func (o *Object) HasUserMeta() bool {
	return len(o.UserMeta) > 0
}

// HasLinks is a bool check to determine if the object contains any links
func (o *Object) HasLinks() bool {
	return len(o.Links) > 0
}

// AddToIntIndex adds the object to the specified secondary index with the integer value to be used
// for index searches
func (o *Object) AddToIntIndex(indexName string, indexValue int) {
	o.AddToIndex(indexName, fmt.Sprintf("%v", indexValue))
}

// AddToIndex adds the object to the specified secondary index with the string value to be used for
// index searches
func (o *Object) AddToIndex(indexName string, indexValue string) {
	if o.Indexes == nil {
		o.Indexes = make(map[string][]string)
	}
	if o.Indexes[indexName] == nil {
		o.Indexes[indexName] = make([]string, 1)
		o.Indexes[indexName][0] = indexValue
	} else {
		o.Indexes[indexName] = append(o.Indexes[indexName], indexValue)
	}
}

func fromRpbContent(rpbContent *rpbRiakKV.RpbContent) (ro *Object, err error) {
	// NB: ro = "Riak Object"
	ro = &Object{
		IsTombstone: rpbContent.GetDeleted(),
	}

	if ro.IsTombstone {
		ro.Value = nil
	} else {
		ro.Value = rpbContent.GetValue()
	}

	ro.ContentType = string(rpbContent.GetContentType())
	ro.Charset = string(rpbContent.GetCharset())
	ro.ContentEncoding = string(rpbContent.GetContentEncoding())
	ro.VTag = string(rpbContent.GetVtag())
	ro.LastModified = time.Unix(int64(rpbContent.GetLastMod()), int64(rpbContent.GetLastModUsecs()))

	rpbUserMeta := rpbContent.GetUsermeta()
	if len(rpbUserMeta) > 0 {
		ro.UserMeta = make([]*Pair, len(rpbUserMeta))
		for i, userMeta := range rpbUserMeta {
			ro.UserMeta[i] = &Pair{
				Key:   string(userMeta.Key),
				Value: string(userMeta.Value),
			}
		}
	}

	rpbIndexes := rpbContent.GetIndexes()
	if len(rpbIndexes) > 0 {
		ro.Indexes = make(map[string][]string)
		for _, index := range rpbIndexes {
			indexName := string(index.Key)
			indexValue := string(index.Value)
			if ro.Indexes[indexName] == nil {
				ro.Indexes[indexName] = make([]string, 1)
				ro.Indexes[indexName][0] = indexValue
			} else {
				ro.Indexes[indexName] = append(ro.Indexes[indexName], indexValue)
			}
		}
	}

	rpbLinks := rpbContent.GetLinks()
	if len(rpbLinks) > 0 {
		ro.Links = make([]*Link, len(rpbLinks))
		for i, link := range rpbLinks {
			ro.Links[i] = &Link{
				Bucket: string(link.Bucket),
				Key:    string(link.Key),
				Tag:    string(link.Tag),
			}
		}
	}

	return
}

func toRpbContent(ro *Object) (*rpbRiakKV.RpbContent, error) {
	rpbContent := &rpbRiakKV.RpbContent{
		Value:           ro.Value,
		ContentType:     rpbBytes(ro.ContentType),
		Charset:         rpbBytes(ro.Charset),
		ContentEncoding: rpbBytes(ro.ContentEncoding),
	}

	if ro.HasIndexes() {
		count := 0
		for _, idxValues := range ro.Indexes {
			count += len(idxValues)
		}
		idx := 0
		rpbIndexes := make([]*rpbRiak.RpbPair, count)
		for idxName, idxValues := range ro.Indexes {
			idxNameBytes := []byte(idxName)
			for _, idxVal := range idxValues {
				pair := &rpbRiak.RpbPair{
					Key:   idxNameBytes,
					Value: []byte(idxVal),
				}
				rpbIndexes[idx] = pair
				idx++
			}
		}
		rpbContent.Indexes = rpbIndexes
	}

	if ro.HasUserMeta() {
		rpbUserMeta := make([]*rpbRiak.RpbPair, len(ro.UserMeta))
		for i, userMeta := range ro.UserMeta {
			rpbUserMeta[i] = &rpbRiak.RpbPair{
				Key:   []byte(userMeta.Key),
				Value: []byte(userMeta.Value),
			}
		}
		rpbContent.Usermeta = rpbUserMeta
	}

	if ro.HasLinks() {
		rpbLinks := make([]*rpbRiakKV.RpbLink, len(ro.Links))
		for i, link := range ro.Links {
			rpbLinks[i] = &rpbRiakKV.RpbLink{
				Bucket: []byte(link.Bucket),
				Key:    []byte(link.Key),
				Tag:    []byte(link.Tag),
			}
		}
		rpbContent.Links = rpbLinks
	}

	return rpbContent, nil
}
