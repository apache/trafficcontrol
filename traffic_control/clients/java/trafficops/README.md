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