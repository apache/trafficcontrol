/*


 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

*/

var PropertiesModel = function() {

	this.properties = {};
	this.loaded = false;

	this.setProperties = function(properties) {
		this.properties = properties;
		this.loaded = true;
	};

	this.setTopology = function(topology) {
		this.topology = topology;
	};

	this.topology = {
		"name": "FooTopology",
		"desc": "a topology for Foo DSes",
		"nodes": [
			{
				"cachegroup": "us-de-newcastle",
				"parents": [
					4,
					5
				]
			},
			{
				"cachegroup": "us-ga-atlanta",
				"parents": [
					3,
					5
				]
			},
			{
				"cachegroup": "us-co-denver",
				"parents": [
					5,
					4
				]
			},
			{
				"cachegroup": "mid-southcentral",
				"parents": []
			},
			{
				"cachegroup": "mid-northeast",
				"parents": []
			},
			{
				"cachegroup": "mid-west",
				"parents": []
			}
		]
	};

};

PropertiesModel.$inject = [];
module.exports = PropertiesModel;