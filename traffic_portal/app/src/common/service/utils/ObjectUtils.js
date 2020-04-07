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

var ObjectUtils = function() {

	let merge = function(objects) {
		let out = {};

		for (let i = 0; i < objects.length; i++) {
			for (let p in objects[i]) {
				out[p] = objects[i][p];
			}
		}

		return out;
	};

	this.flatten = function(obj, name, branch) {
		let out = {},
			newBranch = (typeof branch !== 'undefined' && branch !== '') ? branch + '.' + name : name;

		if (typeof obj !== 'object') {
			out[newBranch] = obj;
			return out;
		}

		for (let p in obj) {
			let prop = this.flatten(obj[p], p, newBranch);
			out = merge([out, prop]);
		}

		return out;
	};

	this.removeKeysWithout = function (obj, without) {
		let out = _.clone(obj);

		Object.keys(out).forEach(function(key) {
			if(key.indexOf(without) == -1) delete out[key];
		});

		return out;
	};

	this.traverse = function(obj) {
		_.each(obj, function (val, key, obj) {
			if (_.isArray(val)) {
				val.forEach(function(el) {
					traverse(el);
				});
			} else if (_.isObject(val)) {
				traverse(val);
			} else {
				console.log('i am a leaf');
				console.log(val);
			}
		});
	};

};

ObjectUtils.$inject = [];
module.exports = ObjectUtils;
