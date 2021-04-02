package dnssec

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	. "github.com/miekg/dns"
	. "github.com/onsi/gomega"
	"log"
)

type DnssecClient struct {
	*Client
}

type SignedRRSet struct {
	RRSIG RRSIG
	RRSet []RR
}

type SignedKeys struct {
	SignedZsks []SignedRRSet
	SignedKsks []SignedRRSet
}

func MakeLabelHierarchy(label string) []string {
	labels := []string{}
	done := false
	i := 0
	for !done {
		label = label[i:]
		labels = append([]string{label}, labels...)
		i, done = NextLabel(label, i)
	}

	return append([]string{"."}, labels...)
}

func (d *DnssecClient) GetRecords(nameserver string, name string, t uint16) *Msg {
	m := new(Msg)
	m.Id = Id()
	m.RecursionDesired = true
	m.SetEdns0(4096, true)
	m.Question = []Question{{Name: name, Qtype: t, Qclass: ClassINET}}
	r, _, err := d.Exchange(m, nameserver)

	Expect(err).Should(BeNil())
	Expect(len(r.Answer)).ToNot(Equal(0), "Received no answers from %v for query of records type %d for zone %v", nameserver, t, name)
	return r
}

func sigCovers(s RRSIG, rr RR) bool {
	return s.TypeCovered == rr.Header().Rrtype &&
		s.Hdr.Class == rr.Header().Class &&
		s.Hdr.Ttl == rr.Header().Ttl
}

func (d *DnssecClient) GetSignedRRSets(nameserver string, name string, t uint16) []SignedRRSet {
	records := []RR{}
	rrsigs := []RR{}

	answers := d.GetRecords(nameserver, name, t).Answer
	for _, ans := range answers {
		if ans.Header().Rrtype == TypeRRSIG {
			rrsigs = append(rrsigs, ans)
		} else {
			records = append(records, ans)
		}
	}

	rrsets := []SignedRRSet{}
	for _, sig := range rrsigs {
		switch s := sig.(type) {
		case *RRSIG:
			rs := RRSIG{
				Hdr:         s.Hdr,
				Signature:   s.Signature,
				Algorithm:   s.Algorithm,
				Expiration:  s.Expiration,
				Inception:   s.Inception,
				KeyTag:      s.KeyTag,
				Labels:      s.Labels,
				OrigTtl:     s.OrigTtl,
				SignerName:  s.SignerName,
				TypeCovered: s.TypeCovered,
			}

			signedSet := SignedRRSet{
				RRSIG: rs,
			}

			for _, rr := range records {
				if sigCovers(*s, rr) {
					signedSet.RRSet = append(signedSet.RRSet, rr)
				} else {
					log.Println("rrsig does not cover record")
					log.Println(s.Header(), s.TypeCovered)
					log.Println(rr.Header(), rr.Header().Rrtype)
				}
			}

			rrsets = append(rrsets, signedSet)
		}

	}

	return rrsets
}

func (d *DnssecClient) DelegationSignerData(nameserver string, name string) []SignedRRSet {
	return d.GetSignedRRSets(nameserver, name, TypeDS)
}

func (d *DnssecClient) SigningData(nameserver string, name string) SignedKeys {
	var signedKeys = SignedKeys{
		SignedZsks: []SignedRRSet{},
		SignedKsks: []SignedRRSet{},
	}

	signedRrsets := d.GetSignedRRSets(nameserver, name, TypeDNSKEY)

	for _, signedRRset := range signedRrsets {
		if len(signedRRset.RRSet) < 1 {
			log.Println("****** no rrset")
			continue
		}

		for _, rr := range signedRRset.RRSet {
			switch k := rr.(type) {
			case *DNSKEY:
				if k.Flags&1 == 1 {
					signedKeys.SignedKsks = append(signedKeys.SignedKsks, signedRRset)
				} else {
					signedKeys.SignedZsks = append(signedKeys.SignedZsks, signedRRset)
				}
			}
		}
	}

	return signedKeys
}
