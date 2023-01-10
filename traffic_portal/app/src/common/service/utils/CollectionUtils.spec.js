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

describe("CollectionUtils tests", function () {
    describe("minimizeArrayDiff tests", function () {

        let minimizeArrayDiff;
        beforeEach(angular.mock.module("trafficPortal.utils"));
        beforeEach(inject(function () {
            const $injector = angular.injector(["trafficPortal.utils"]);
            minimizeArrayDiff = $injector.get("collectionUtils").minimizeArrayDiff;
        }));

        it("pads removed capabilities with undefined", async function () {
            let oldCapabilities = ["cap1"];
            let newCapabilities = ["cap2"];
            let expected = [undefined, "cap2"];
            expect(minimizeArrayDiff(oldCapabilities, newCapabilities)).toEqual(expected);

            oldCapabilities = ["cap1", "cap2", "cap3"];
            newCapabilities = ["cap1", "cap3"];
            expected = ["cap1", undefined, "cap3"];
            expect(minimizeArrayDiff(oldCapabilities, newCapabilities)).toEqual(expected);

            oldCapabilities = ["cap1", "cap2", "cap3"];
            newCapabilities = ["cap1", "cap2"];
            expected = ["cap1", "cap2"];
            expect(minimizeArrayDiff(oldCapabilities, newCapabilities)).toEqual(expected);
        });

        it("appends prepended capabilities", async function () {
            let oldCapabilities = ["cap2"];
            let newCapabilities = ["cap1", "cap2"];
            let expected = ["cap2", "cap1"];
            expect(minimizeArrayDiff(oldCapabilities, newCapabilities)).toEqual(expected);
        });

        it("appends added capabilities", async function () {
            let oldCapabilities = ["cap1"];
            let newCapabilities = ["cap1", "cap2"];
            let expected = ["cap1", "cap2"];
            expect(minimizeArrayDiff(oldCapabilities, newCapabilities)).toEqual(expected);
        });

        it("leaves equal capabilities arrays untouched", async function () {
            let oldCapabilities = ["cap1", "cap2"];
            let newCapabilities = ["cap1", "cap2"];
            expect(minimizeArrayDiff(oldCapabilities, newCapabilities)).toEqual(newCapabilities);
        });
    });
});
