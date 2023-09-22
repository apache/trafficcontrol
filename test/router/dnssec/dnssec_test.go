package dnssec_test

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
	"github.com/apache/trafficcontrol/v8/test/router/dnssec"
	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dnssec", func() {
	Context("The Interwebs", func() {
		It("Makes Label Hierarchy", func() {
			Expect(dnssec.MakeLabelHierarchy("example.com.")).To(Equal([]string{".", "com.", "example.com."}))
		})

		It("Uses Parent Zone Key to validate DS", func() {
			signedDSSets := d.DelegationSignerData(nameserver, deliveryService)

			Expect(len(signedDSSets)).ToNot(Equal(0))
			Expect(len(signedDSSets[0].RRSet)).ToNot(Equal(0))

			verifiedCount := 0
			for _, signedDSSet := range signedDSSets {

				signedKeys := d.SigningData(nameserver, signedDSSet.RRSIG.SignerName)

				Expect(len(signedKeys.SignedKsks)).ToNot(Equal(0))
				Expect(len(signedKeys.SignedZsks)).ToNot(Equal(0))

				for _, sk := range signedKeys.SignedZsks {
					for _, k := range sk.RRSet {
						switch kk := k.(type) {
						case *dns.DNSKEY:
							if kk.KeyTag() == signedDSSet.RRSIG.KeyTag {
								Expect(signedDSSet.RRSIG.Verify(kk, signedDSSet.RRSet)).To(BeNil())
								verifiedCount++
							}
						}
					}
				}

				for _, sk := range signedKeys.SignedKsks {
					for _, k := range sk.RRSet {
						switch kk := k.(type) {
						case *dns.DNSKEY:
							if kk.KeyTag() == signedDSSet.RRSIG.KeyTag {
								Expect(signedDSSet.RRSIG.Verify(kk, signedDSSet.RRSet)).To(BeNil())
								verifiedCount++
							}
						}
					}
				}
			}

			Expect(verifiedCount).ToNot(Equal(0))
		})

		It("Uses DS to validate Public Key", func() {
			signedKeys := d.SigningData(nameserver, deliveryService)
			signedDSSets := d.DelegationSignerData(nameserver, deliveryService)

			Expect(len(signedDSSets)).ToNot(Equal(0))

			count := 0
			for _, signedZsk := range signedKeys.SignedZsks {
				for _, zsk := range signedZsk.RRSet {
					switch z := zsk.(type) {
					case *dns.DNSKEY:
						for _, signedDs := range signedDSSets {
							for _, ds := range signedDs.RRSet {
								switch d := ds.(type) {
								case *dns.DS:
									if d.KeyTag == z.KeyTag() {
										computedDS := z.ToDS(d.DigestType)
										Expect(d.Digest).To(Equal(computedDS.Digest))
										count++
									}
								}
							}
						}
					}
				}
			}

			Expect(count).ToNot(Equal(0))
		})

		It("Uses KSK public key to verify ZSK RRSig", func() {
			signedKeys := d.SigningData(nameserver, deliveryService)

			count := 0
			for _, signedZsk := range signedKeys.SignedZsks {
				for _, signedKsk := range signedKeys.SignedKsks {
					for _, ksk := range signedKsk.RRSet {
						switch k := ksk.(type) {
						case *dns.DNSKEY:
							if k.KeyTag() == signedZsk.RRSIG.KeyTag {
								Expect(signedZsk.RRSIG.Verify(k, signedZsk.RRSet)).To(BeNil())
								count++
							}
						}
					}
				}
			}
			Expect(count).ToNot(Equal(0))
		})
	})
})
