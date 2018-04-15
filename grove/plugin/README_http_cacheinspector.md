# Cache Insepctor Plugin

The cache inspector plugin allows you to view the contents of the cache.  Example output:

```
Jump to:  disk  my-disk-cache-two  


*** Cache "" ***

  * Size of in use cache:      9.5M 
  * Cache capacity:            9.5M 
  * Number of elements in LRU: 169
  * Objects in cache sorted by Least Recently Used on top, showing only first 100 and last 100:

     #		Code	Size	        Age			            Key
     00000	200	    34K	            16h19m9.135951506s  	GET:http://localhost/34k.bin?hah
     00001	200	    35K	            16h19m9.103181146s  	GET:http://localhost/35k.bin?hah
     00002	200	    36K	            16h19m9.079490592s  	GET:http://localhost/36k.bin?hah
     00003	200	    37K	            16h19m9.050852429s  	GET:http://localhost/37k.bin?hah
     00004	200	    38K	            16h19m9.020264008s  	GET:http://localhost/38k.bin?hah
     00005	200	    39K	            16h19m8.989722029s  	GET:http://localhost/39k.bin?hah
     00006	200 	40K	            16h19m8.967823063s  	GET:http://localhost/40k.bin?hah
     00007	200	    41K	            16h19m8.939422783s  	GET:http://localhost/41k.bin?hah
     00008	200	    42K	            16h19m8.913485123s  	GET:http://localhost/42k.bin?h
```


Any of the keys can be clicked to peek at the details of this object in cache, and this will not update the LRU list:

```
Key: GET:http://localhost/34k.bin?hah cache: ""

  > User-Agent: curl/7.54.0
  > Accept: */*
  > Host: localhost

  < Last-Modified: Sun, 01 Apr 2018 19:42:43 GMT
  < Etag: "8800-568cead43b2c0"
  < Accept-Ranges: bytes
  < Content-Length: 34816
  < Content-Type: application/octet-stream
  < Date: Sat, 14 Apr 2018 22:18:35 GMT
  < Server: Apache/2.4.29 (Unix)

  Code:                         200
  OriginCode:                   200
  ProxyURL:                     
  ReqTime:                      2018-04-14 16:18:35.098419 -0600 MDT m=+33.256292902
  ReqRespTime:                  2018-04-14 16:18:35.098843 -0600 MDT m=+33.256716928
  RespRespTime:                 2018-04-14 22:18:35 +0000 GMT
  LastModified:                 2018-04-01 19:42:43 +0000 GMT
```

The following querystrings can be used: 

- `cache=<cachename>`
Only display the contents of the cache `<cachename>`. Use `cache=` (empty cachename) for the default memory cachename `""`.
- `head=<number>`
Number of items to list from the top of the top of the LRU. Default is 100.
- `key=<keystring>`
Show details page of the given key.
- `search=<searchstring>`
List only items that have `<searchstring>` in the key. This overrules `<head>` and `<tail>`, when search is used these are ignored. 
- `tail=<number>`
Number of items to list from the bottom of the top of the LRU. Default is 100.