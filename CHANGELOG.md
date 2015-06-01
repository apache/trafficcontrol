## 1.1.3

### Features
- [156](https://github.com/Comcast/traffic_control/issues/156): Update SELECT * to return tag values.
	- Updated the record stored in Riak to contian an array of keys for KSK and ZSK.
	- Added a status field to the keys stored in Riak ('new' for new keys and 'expired' for expired keys)
		- Use this key to determine which key to check for expiration ('new' key).
		- If key is expired set it to 'expired' and add back to array with new key


### Bugfixes
