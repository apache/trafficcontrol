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

{
    "deliveryServices": {
        "omg-08": {
            "health.threshold.total.kbps": 100000,
            "health.threshold.total.tps_total": 20,
            "status": "REPORTED"
        }
    },
    "profiles": {
        "CCR": {
            "CCR1": null
        },
        "EDGE": {
            "EDGE1": {
                "health.connection.timeout": "2000",
                "health.polling.url": "http://${hostname}/_astats?application=&inf.name=${interface_name}",
                "health.threshold.availableBandwidthInKbps": ">800000",
                "health.threshold.availableBandwidthInMbps": ">800000",
                "health.threshold.loadavg": "25.0",
                "health.threshold.myNewParam": ">444",
                "health.threshold.myOtherNewParam": "<0",
                "health.threshold.queryTime": "1000",
                "health.timepad": "400",
                "history.count": "30"
            },
        },
        "MID": {
            "MID1": {
                "health.connection.timeout": "2000",
                "health.polling.url": "http://${hostname}/_astats?application=&inf.name=${interface_name}",
                "health.threshold.availableBandwidthInKbps": ">50000",
                "health.threshold.availableBandwidthInMbps": ">800000",
                "health.threshold.loadavg": "25.0",
                "health.threshold.myNewParam": ">444",
                "health.threshold.myOtherNewParam": "<0",
                "health.threshold.queryTime": "1000",
                "health.timepad": "400",
                "history.count": "30"
            }
        }
    },
    "rascal-config": {
        "CDN_name": "jenkins",
        "hack.ttl": "30",
        "health.event-count": "200",
        "health.polling.interval": "8000",
        "health.threadPool": "4",
        "health.timepad": "100",
        "tm.dataServer.polling.url": "https://${tmHostname}/dataserver/orderby/id",
        "tm.healthParams.polling.url": "https://${tmHostname}/health/${cdnName}",
        "tm.polling.interval": "5000"
    }
}
