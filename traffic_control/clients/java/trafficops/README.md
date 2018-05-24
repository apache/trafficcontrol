<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# TrafficOps Api Client

Simple Java API client for communicating with the TrafficOps API

## Example Usage

### Create Traffic Ops session

**Using a provided URI**

```java
//Construct TrafficOps Session
final URI trafficOpsUri = new URI("https://trafficops.mycdn.com:443");
final TOSession.Builder toSessionBuilder = TOSession.builder()
	.fromURI(trafficOpsUri);
final RestApiSession.Builder restSession = toSessionBuilder.restClientBuilder();
final TOSession toSession = toSessionBuilder.build();
```
**Explicitly set properties**

```java
//Construct TrafficOps Session
final URI trafficOpsUri = new URI("http://trafficops.mycdn.com:443");
final TOSession.Builder toSessionBuilder = TOSession.builder()
	.setHost("trafficops.mycdn.com")
	.setPort(443)
	.setSsl(true);
final RestApiSession.Builder restSession = toSessionBuilder.restClientBuilder();
final TOSession toSession = toSessionBuilder.build();
```

### Logging In

```java
//Login
final CompletableFuture<Boolean> loginFuture = toSession
	.login("MyUsername", "MyPassword");
try {
	//Timeout if login takes longer then 1sec
	loginFuture.get(1000, TimeUnit.MILLISECONDS);
} catch(TimeoutException e) {
	loginFuture.cancel(true);
	LOG.error("Timeout while logging in");
	System.exit(1);
}
```

### Getting a list of All Servers

**Synchronously**

```java
final CollectionResponse response = toSession.getServers().get();
```

**Asynchronously**

```java
toSession
	.getServers()
	.whenCompleteAsync((servers, exception)->{
		if(exception != null){
			//Handle Exception
		} else {
			//Do something with your server list
		}
	});

```