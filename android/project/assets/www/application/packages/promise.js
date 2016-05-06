//////////////////////////////////////////////////////////////////////////
//                                                                      //
// This is a generated file. You can view the original                  //
// source in your browser if your browser supports source maps.         //
// Source maps are supported by all recent versions of Chrome, Safari,  //
// and Firefox, and by Internet Explorer 11.                            //
//                                                                      //
//////////////////////////////////////////////////////////////////////////


(function () {

/* Imports */
var Meteor = Package.meteor.Meteor;
var global = Package.meteor.global;
var meteorEnv = Package.meteor.meteorEnv;
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;

/* Package-scope variables */
var exports, Promise;

var require = meteorInstall({"node_modules":{"meteor":{"promise":{"client.js":["meteor-promise",function(require,exports){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// packages/promise/client.js                                                                                       //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
exports.Promise = require("meteor-promise");                                                                        // 1
                                                                                                                    // 2
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"node_modules":{"meteor-promise":{"package.json":function(require,exports){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// ../npm/node_modules/meteor-promise/package.json                                                                  //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
exports.name = "meteor-promise";                                                                                    // 1
exports.version = "0.6.3";                                                                                          // 2
exports.main = "promise_client.js";                                                                                 // 3
                                                                                                                    // 4
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"promise_client.js":["./promise.js",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/promise_client.js                                        //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
var MeteorPromise = require("./promise.js");                                                                        // 1
                                                                                                                    // 2
var es6PromiseThen = MeteorPromise.prototype.then;                                                                  // 3
MeteorPromise.prototype.then = function (onResolved, onRejected) {                                                  // 4
  if (typeof Meteor === "object" &&                                                                                 // 5
      typeof Meteor.bindEnvironment === "function") {                                                               // 6
    return es6PromiseThen.call(                                                                                     // 7
      this,                                                                                                         // 8
      onResolved && Meteor.bindEnvironment(onResolved, raise),                                                      // 9
      onRejected && Meteor.bindEnvironment(onRejected, raise)                                                       // 10
    );                                                                                                              // 11
  }                                                                                                                 // 12
  return es6PromiseThen.call(this, onResolved, onRejected);                                                         // 13
};                                                                                                                  // 14
                                                                                                                    // 15
function raise(exception) {                                                                                         // 16
  throw exception;                                                                                                  // 17
}                                                                                                                   // 18
                                                                                                                    // 19
module.exports = exports = MeteorPromise;                                                                           // 20
                                                                                                                    // 21
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"promise.js":["promise",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/promise.js                                               //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
var hasOwn = Object.prototype.hasOwnProperty;                                                                       // 1
                                                                                                                    // 2
var g =                                                                                                             // 3
  typeof global === "object" ? global :                                                                             // 4
  typeof window === "object" ? window :                                                                             // 5
  typeof self === "object" ? self : this;                                                                           // 6
                                                                                                                    // 7
var GlobalPromise = g.Promise;                                                                                      // 8
var NpmPromise = require("promise");                                                                                // 9
                                                                                                                    // 10
function copyMethods(target, source) {                                                                              // 11
  Object.keys(source).forEach(function (key) {                                                                      // 12
    var value = source[key];                                                                                        // 13
    if (typeof value === "function" &&                                                                              // 14
        ! hasOwn.call(target, key)) {                                                                               // 15
      target[key] = value;                                                                                          // 16
    }                                                                                                               // 17
  });                                                                                                               // 18
}                                                                                                                   // 19
                                                                                                                    // 20
if (typeof GlobalPromise === "function") {                                                                          // 21
  copyMethods(GlobalPromise, NpmPromise);                                                                           // 22
  copyMethods(GlobalPromise.prototype, NpmPromise.prototype);                                                       // 23
  module.exports = GlobalPromise;                                                                                   // 24
} else {                                                                                                            // 25
  module.exports = NpmPromise;                                                                                      // 26
}                                                                                                                   // 27
                                                                                                                    // 28
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"node_modules":{"promise":{"package.json":function(require,exports){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// ../npm/node_modules/meteor-promise/node_modules/promise/package.json                                             //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
exports.name = "promise";                                                                                           // 1
exports.version = "7.1.1";                                                                                          // 2
exports.main = "index.js";                                                                                          // 3
                                                                                                                    // 4
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"index.js":["./lib",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/index.js                            //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
module.exports = require('./lib')                                                                                   // 3
                                                                                                                    // 4
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"lib":{"index.js":["./core.js","./done.js","./finally.js","./es6-extensions.js","./node-extensions.js","./synchronous.js",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/index.js                        //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
module.exports = require('./core.js');                                                                              // 3
require('./done.js');                                                                                               // 4
require('./finally.js');                                                                                            // 5
require('./es6-extensions.js');                                                                                     // 6
require('./node-extensions.js');                                                                                    // 7
require('./synchronous.js');                                                                                        // 8
                                                                                                                    // 9
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"core.js":["asap/raw",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/core.js                         //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
var asap = require('asap/raw');                                                                                     // 3
                                                                                                                    // 4
function noop() {}                                                                                                  // 5
                                                                                                                    // 6
// States:                                                                                                          // 7
//                                                                                                                  // 8
// 0 - pending                                                                                                      // 9
// 1 - fulfilled with _value                                                                                        // 10
// 2 - rejected with _value                                                                                         // 11
// 3 - adopted the state of another promise, _value                                                                 // 12
//                                                                                                                  // 13
// once the state is no longer pending (0) it is immutable                                                          // 14
                                                                                                                    // 15
// All `_` prefixed properties will be reduced to `_{random number}`                                                // 16
// at build time to obfuscate them and discourage their use.                                                        // 17
// We don't use symbols or Object.defineProperty to fully hide them                                                 // 18
// because the performance isn't good enough.                                                                       // 19
                                                                                                                    // 20
                                                                                                                    // 21
// to avoid using try/catch inside critical functions, we                                                           // 22
// extract them to here.                                                                                            // 23
var LAST_ERROR = null;                                                                                              // 24
var IS_ERROR = {};                                                                                                  // 25
function getThen(obj) {                                                                                             // 26
  try {                                                                                                             // 27
    return obj.then;                                                                                                // 28
  } catch (ex) {                                                                                                    // 29
    LAST_ERROR = ex;                                                                                                // 30
    return IS_ERROR;                                                                                                // 31
  }                                                                                                                 // 32
}                                                                                                                   // 33
                                                                                                                    // 34
function tryCallOne(fn, a) {                                                                                        // 35
  try {                                                                                                             // 36
    return fn(a);                                                                                                   // 37
  } catch (ex) {                                                                                                    // 38
    LAST_ERROR = ex;                                                                                                // 39
    return IS_ERROR;                                                                                                // 40
  }                                                                                                                 // 41
}                                                                                                                   // 42
function tryCallTwo(fn, a, b) {                                                                                     // 43
  try {                                                                                                             // 44
    fn(a, b);                                                                                                       // 45
  } catch (ex) {                                                                                                    // 46
    LAST_ERROR = ex;                                                                                                // 47
    return IS_ERROR;                                                                                                // 48
  }                                                                                                                 // 49
}                                                                                                                   // 50
                                                                                                                    // 51
module.exports = Promise;                                                                                           // 52
                                                                                                                    // 53
function Promise(fn) {                                                                                              // 54
  if (typeof this !== 'object') {                                                                                   // 55
    throw new TypeError('Promises must be constructed via new');                                                    // 56
  }                                                                                                                 // 57
  if (typeof fn !== 'function') {                                                                                   // 58
    throw new TypeError('not a function');                                                                          // 59
  }                                                                                                                 // 60
  this._45 = 0;                                                                                                     // 61
  this._81 = 0;                                                                                                     // 62
  this._65 = null;                                                                                                  // 63
  this._54 = null;                                                                                                  // 64
  if (fn === noop) return;                                                                                          // 65
  doResolve(fn, this);                                                                                              // 66
}                                                                                                                   // 67
Promise._10 = null;                                                                                                 // 68
Promise._97 = null;                                                                                                 // 69
Promise._61 = noop;                                                                                                 // 70
                                                                                                                    // 71
Promise.prototype.then = function(onFulfilled, onRejected) {                                                        // 72
  if (this.constructor !== Promise) {                                                                               // 73
    return safeThen(this, onFulfilled, onRejected);                                                                 // 74
  }                                                                                                                 // 75
  var res = new Promise(noop);                                                                                      // 76
  handle(this, new Handler(onFulfilled, onRejected, res));                                                          // 77
  return res;                                                                                                       // 78
};                                                                                                                  // 79
                                                                                                                    // 80
function safeThen(self, onFulfilled, onRejected) {                                                                  // 81
  return new self.constructor(function (resolve, reject) {                                                          // 82
    var res = new Promise(noop);                                                                                    // 83
    res.then(resolve, reject);                                                                                      // 84
    handle(self, new Handler(onFulfilled, onRejected, res));                                                        // 85
  });                                                                                                               // 86
};                                                                                                                  // 87
function handle(self, deferred) {                                                                                   // 88
  while (self._81 === 3) {                                                                                          // 89
    self = self._65;                                                                                                // 90
  }                                                                                                                 // 91
  if (Promise._10) {                                                                                                // 92
    Promise._10(self);                                                                                              // 93
  }                                                                                                                 // 94
  if (self._81 === 0) {                                                                                             // 95
    if (self._45 === 0) {                                                                                           // 96
      self._45 = 1;                                                                                                 // 97
      self._54 = deferred;                                                                                          // 98
      return;                                                                                                       // 99
    }                                                                                                               // 100
    if (self._45 === 1) {                                                                                           // 101
      self._45 = 2;                                                                                                 // 102
      self._54 = [self._54, deferred];                                                                              // 103
      return;                                                                                                       // 104
    }                                                                                                               // 105
    self._54.push(deferred);                                                                                        // 106
    return;                                                                                                         // 107
  }                                                                                                                 // 108
  handleResolved(self, deferred);                                                                                   // 109
}                                                                                                                   // 110
                                                                                                                    // 111
function handleResolved(self, deferred) {                                                                           // 112
  asap(function() {                                                                                                 // 113
    var cb = self._81 === 1 ? deferred.onFulfilled : deferred.onRejected;                                           // 114
    if (cb === null) {                                                                                              // 115
      if (self._81 === 1) {                                                                                         // 116
        resolve(deferred.promise, self._65);                                                                        // 117
      } else {                                                                                                      // 118
        reject(deferred.promise, self._65);                                                                         // 119
      }                                                                                                             // 120
      return;                                                                                                       // 121
    }                                                                                                               // 122
    var ret = tryCallOne(cb, self._65);                                                                             // 123
    if (ret === IS_ERROR) {                                                                                         // 124
      reject(deferred.promise, LAST_ERROR);                                                                         // 125
    } else {                                                                                                        // 126
      resolve(deferred.promise, ret);                                                                               // 127
    }                                                                                                               // 128
  });                                                                                                               // 129
}                                                                                                                   // 130
function resolve(self, newValue) {                                                                                  // 131
  // Promise Resolution Procedure: https://github.com/promises-aplus/promises-spec#the-promise-resolution-procedure
  if (newValue === self) {                                                                                          // 133
    return reject(                                                                                                  // 134
      self,                                                                                                         // 135
      new TypeError('A promise cannot be resolved with itself.')                                                    // 136
    );                                                                                                              // 137
  }                                                                                                                 // 138
  if (                                                                                                              // 139
    newValue &&                                                                                                     // 140
    (typeof newValue === 'object' || typeof newValue === 'function')                                                // 141
  ) {                                                                                                               // 142
    var then = getThen(newValue);                                                                                   // 143
    if (then === IS_ERROR) {                                                                                        // 144
      return reject(self, LAST_ERROR);                                                                              // 145
    }                                                                                                               // 146
    if (                                                                                                            // 147
      then === self.then &&                                                                                         // 148
      newValue instanceof Promise                                                                                   // 149
    ) {                                                                                                             // 150
      self._81 = 3;                                                                                                 // 151
      self._65 = newValue;                                                                                          // 152
      finale(self);                                                                                                 // 153
      return;                                                                                                       // 154
    } else if (typeof then === 'function') {                                                                        // 155
      doResolve(then.bind(newValue), self);                                                                         // 156
      return;                                                                                                       // 157
    }                                                                                                               // 158
  }                                                                                                                 // 159
  self._81 = 1;                                                                                                     // 160
  self._65 = newValue;                                                                                              // 161
  finale(self);                                                                                                     // 162
}                                                                                                                   // 163
                                                                                                                    // 164
function reject(self, newValue) {                                                                                   // 165
  self._81 = 2;                                                                                                     // 166
  self._65 = newValue;                                                                                              // 167
  if (Promise._97) {                                                                                                // 168
    Promise._97(self, newValue);                                                                                    // 169
  }                                                                                                                 // 170
  finale(self);                                                                                                     // 171
}                                                                                                                   // 172
function finale(self) {                                                                                             // 173
  if (self._45 === 1) {                                                                                             // 174
    handle(self, self._54);                                                                                         // 175
    self._54 = null;                                                                                                // 176
  }                                                                                                                 // 177
  if (self._45 === 2) {                                                                                             // 178
    for (var i = 0; i < self._54.length; i++) {                                                                     // 179
      handle(self, self._54[i]);                                                                                    // 180
    }                                                                                                               // 181
    self._54 = null;                                                                                                // 182
  }                                                                                                                 // 183
}                                                                                                                   // 184
                                                                                                                    // 185
function Handler(onFulfilled, onRejected, promise){                                                                 // 186
  this.onFulfilled = typeof onFulfilled === 'function' ? onFulfilled : null;                                        // 187
  this.onRejected = typeof onRejected === 'function' ? onRejected : null;                                           // 188
  this.promise = promise;                                                                                           // 189
}                                                                                                                   // 190
                                                                                                                    // 191
/**                                                                                                                 // 192
 * Take a potentially misbehaving resolver function and make sure                                                   // 193
 * onFulfilled and onRejected are only called once.                                                                 // 194
 *                                                                                                                  // 195
 * Makes no guarantees about asynchrony.                                                                            // 196
 */                                                                                                                 // 197
function doResolve(fn, promise) {                                                                                   // 198
  var done = false;                                                                                                 // 199
  var res = tryCallTwo(fn, function (value) {                                                                       // 200
    if (done) return;                                                                                               // 201
    done = true;                                                                                                    // 202
    resolve(promise, value);                                                                                        // 203
  }, function (reason) {                                                                                            // 204
    if (done) return;                                                                                               // 205
    done = true;                                                                                                    // 206
    reject(promise, reason);                                                                                        // 207
  })                                                                                                                // 208
  if (!done && res === IS_ERROR) {                                                                                  // 209
    done = true;                                                                                                    // 210
    reject(promise, LAST_ERROR);                                                                                    // 211
  }                                                                                                                 // 212
}                                                                                                                   // 213
                                                                                                                    // 214
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"done.js":["./core.js",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/done.js                         //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
var Promise = require('./core.js');                                                                                 // 3
                                                                                                                    // 4
module.exports = Promise;                                                                                           // 5
Promise.prototype.done = function (onFulfilled, onRejected) {                                                       // 6
  var self = arguments.length ? this.then.apply(this, arguments) : this;                                            // 7
  self.then(null, function (err) {                                                                                  // 8
    setTimeout(function () {                                                                                        // 9
      throw err;                                                                                                    // 10
    }, 0);                                                                                                          // 11
  });                                                                                                               // 12
};                                                                                                                  // 13
                                                                                                                    // 14
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"finally.js":["./core.js",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/finally.js                      //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
var Promise = require('./core.js');                                                                                 // 3
                                                                                                                    // 4
module.exports = Promise;                                                                                           // 5
Promise.prototype['finally'] = function (f) {                                                                       // 6
  return this.then(function (value) {                                                                               // 7
    return Promise.resolve(f()).then(function () {                                                                  // 8
      return value;                                                                                                 // 9
    });                                                                                                             // 10
  }, function (err) {                                                                                               // 11
    return Promise.resolve(f()).then(function () {                                                                  // 12
      throw err;                                                                                                    // 13
    });                                                                                                             // 14
  });                                                                                                               // 15
};                                                                                                                  // 16
                                                                                                                    // 17
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"es6-extensions.js":["./core.js",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/es6-extensions.js               //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
//This file contains the ES6 extensions to the core Promises/A+ API                                                 // 3
                                                                                                                    // 4
var Promise = require('./core.js');                                                                                 // 5
                                                                                                                    // 6
module.exports = Promise;                                                                                           // 7
                                                                                                                    // 8
/* Static Functions */                                                                                              // 9
                                                                                                                    // 10
var TRUE = valuePromise(true);                                                                                      // 11
var FALSE = valuePromise(false);                                                                                    // 12
var NULL = valuePromise(null);                                                                                      // 13
var UNDEFINED = valuePromise(undefined);                                                                            // 14
var ZERO = valuePromise(0);                                                                                         // 15
var EMPTYSTRING = valuePromise('');                                                                                 // 16
                                                                                                                    // 17
function valuePromise(value) {                                                                                      // 18
  var p = new Promise(Promise._61);                                                                                 // 19
  p._81 = 1;                                                                                                        // 20
  p._65 = value;                                                                                                    // 21
  return p;                                                                                                         // 22
}                                                                                                                   // 23
Promise.resolve = function (value) {                                                                                // 24
  if (value instanceof Promise) return value;                                                                       // 25
                                                                                                                    // 26
  if (value === null) return NULL;                                                                                  // 27
  if (value === undefined) return UNDEFINED;                                                                        // 28
  if (value === true) return TRUE;                                                                                  // 29
  if (value === false) return FALSE;                                                                                // 30
  if (value === 0) return ZERO;                                                                                     // 31
  if (value === '') return EMPTYSTRING;                                                                             // 32
                                                                                                                    // 33
  if (typeof value === 'object' || typeof value === 'function') {                                                   // 34
    try {                                                                                                           // 35
      var then = value.then;                                                                                        // 36
      if (typeof then === 'function') {                                                                             // 37
        return new Promise(then.bind(value));                                                                       // 38
      }                                                                                                             // 39
    } catch (ex) {                                                                                                  // 40
      return new Promise(function (resolve, reject) {                                                               // 41
        reject(ex);                                                                                                 // 42
      });                                                                                                           // 43
    }                                                                                                               // 44
  }                                                                                                                 // 45
  return valuePromise(value);                                                                                       // 46
};                                                                                                                  // 47
                                                                                                                    // 48
Promise.all = function (arr) {                                                                                      // 49
  var args = Array.prototype.slice.call(arr);                                                                       // 50
                                                                                                                    // 51
  return new Promise(function (resolve, reject) {                                                                   // 52
    if (args.length === 0) return resolve([]);                                                                      // 53
    var remaining = args.length;                                                                                    // 54
    function res(i, val) {                                                                                          // 55
      if (val && (typeof val === 'object' || typeof val === 'function')) {                                          // 56
        if (val instanceof Promise && val.then === Promise.prototype.then) {                                        // 57
          while (val._81 === 3) {                                                                                   // 58
            val = val._65;                                                                                          // 59
          }                                                                                                         // 60
          if (val._81 === 1) return res(i, val._65);                                                                // 61
          if (val._81 === 2) reject(val._65);                                                                       // 62
          val.then(function (val) {                                                                                 // 63
            res(i, val);                                                                                            // 64
          }, reject);                                                                                               // 65
          return;                                                                                                   // 66
        } else {                                                                                                    // 67
          var then = val.then;                                                                                      // 68
          if (typeof then === 'function') {                                                                         // 69
            var p = new Promise(then.bind(val));                                                                    // 70
            p.then(function (val) {                                                                                 // 71
              res(i, val);                                                                                          // 72
            }, reject);                                                                                             // 73
            return;                                                                                                 // 74
          }                                                                                                         // 75
        }                                                                                                           // 76
      }                                                                                                             // 77
      args[i] = val;                                                                                                // 78
      if (--remaining === 0) {                                                                                      // 79
        resolve(args);                                                                                              // 80
      }                                                                                                             // 81
    }                                                                                                               // 82
    for (var i = 0; i < args.length; i++) {                                                                         // 83
      res(i, args[i]);                                                                                              // 84
    }                                                                                                               // 85
  });                                                                                                               // 86
};                                                                                                                  // 87
                                                                                                                    // 88
Promise.reject = function (value) {                                                                                 // 89
  return new Promise(function (resolve, reject) {                                                                   // 90
    reject(value);                                                                                                  // 91
  });                                                                                                               // 92
};                                                                                                                  // 93
                                                                                                                    // 94
Promise.race = function (values) {                                                                                  // 95
  return new Promise(function (resolve, reject) {                                                                   // 96
    values.forEach(function(value){                                                                                 // 97
      Promise.resolve(value).then(resolve, reject);                                                                 // 98
    });                                                                                                             // 99
  });                                                                                                               // 100
};                                                                                                                  // 101
                                                                                                                    // 102
/* Prototype Methods */                                                                                             // 103
                                                                                                                    // 104
Promise.prototype['catch'] = function (onRejected) {                                                                // 105
  return this.then(null, onRejected);                                                                               // 106
};                                                                                                                  // 107
                                                                                                                    // 108
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"node-extensions.js":["./core.js","asap",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/node-extensions.js              //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
// This file contains then/promise specific extensions that are only useful                                         // 3
// for node.js interop                                                                                              // 4
                                                                                                                    // 5
var Promise = require('./core.js');                                                                                 // 6
var asap = require('asap');                                                                                         // 7
                                                                                                                    // 8
module.exports = Promise;                                                                                           // 9
                                                                                                                    // 10
/* Static Functions */                                                                                              // 11
                                                                                                                    // 12
Promise.denodeify = function (fn, argumentCount) {                                                                  // 13
  if (                                                                                                              // 14
    typeof argumentCount === 'number' && argumentCount !== Infinity                                                 // 15
  ) {                                                                                                               // 16
    return denodeifyWithCount(fn, argumentCount);                                                                   // 17
  } else {                                                                                                          // 18
    return denodeifyWithoutCount(fn);                                                                               // 19
  }                                                                                                                 // 20
}                                                                                                                   // 21
                                                                                                                    // 22
var callbackFn = (                                                                                                  // 23
  'function (err, res) {' +                                                                                         // 24
  'if (err) { rj(err); } else { rs(res); }' +                                                                       // 25
  '}'                                                                                                               // 26
);                                                                                                                  // 27
function denodeifyWithCount(fn, argumentCount) {                                                                    // 28
  var args = [];                                                                                                    // 29
  for (var i = 0; i < argumentCount; i++) {                                                                         // 30
    args.push('a' + i);                                                                                             // 31
  }                                                                                                                 // 32
  var body = [                                                                                                      // 33
    'return function (' + args.join(',') + ') {',                                                                   // 34
    'var self = this;',                                                                                             // 35
    'return new Promise(function (rs, rj) {',                                                                       // 36
    'var res = fn.call(',                                                                                           // 37
    ['self'].concat(args).concat([callbackFn]).join(','),                                                           // 38
    ');',                                                                                                           // 39
    'if (res &&',                                                                                                   // 40
    '(typeof res === "object" || typeof res === "function") &&',                                                    // 41
    'typeof res.then === "function"',                                                                               // 42
    ') {rs(res);}',                                                                                                 // 43
    '});',                                                                                                          // 44
    '};'                                                                                                            // 45
  ].join('');                                                                                                       // 46
  return Function(['Promise', 'fn'], body)(Promise, fn);                                                            // 47
}                                                                                                                   // 48
function denodeifyWithoutCount(fn) {                                                                                // 49
  var fnLength = Math.max(fn.length - 1, 3);                                                                        // 50
  var args = [];                                                                                                    // 51
  for (var i = 0; i < fnLength; i++) {                                                                              // 52
    args.push('a' + i);                                                                                             // 53
  }                                                                                                                 // 54
  var body = [                                                                                                      // 55
    'return function (' + args.join(',') + ') {',                                                                   // 56
    'var self = this;',                                                                                             // 57
    'var args;',                                                                                                    // 58
    'var argLength = arguments.length;',                                                                            // 59
    'if (arguments.length > ' + fnLength + ') {',                                                                   // 60
    'args = new Array(arguments.length + 1);',                                                                      // 61
    'for (var i = 0; i < arguments.length; i++) {',                                                                 // 62
    'args[i] = arguments[i];',                                                                                      // 63
    '}',                                                                                                            // 64
    '}',                                                                                                            // 65
    'return new Promise(function (rs, rj) {',                                                                       // 66
    'var cb = ' + callbackFn + ';',                                                                                 // 67
    'var res;',                                                                                                     // 68
    'switch (argLength) {',                                                                                         // 69
    args.concat(['extra']).map(function (_, index) {                                                                // 70
      return (                                                                                                      // 71
        'case ' + (index) + ':' +                                                                                   // 72
        'res = fn.call(' + ['self'].concat(args.slice(0, index)).concat('cb').join(',') + ');' +                    // 73
        'break;'                                                                                                    // 74
      );                                                                                                            // 75
    }).join(''),                                                                                                    // 76
    'default:',                                                                                                     // 77
    'args[argLength] = cb;',                                                                                        // 78
    'res = fn.apply(self, args);',                                                                                  // 79
    '}',                                                                                                            // 80
                                                                                                                    // 81
    'if (res &&',                                                                                                   // 82
    '(typeof res === "object" || typeof res === "function") &&',                                                    // 83
    'typeof res.then === "function"',                                                                               // 84
    ') {rs(res);}',                                                                                                 // 85
    '});',                                                                                                          // 86
    '};'                                                                                                            // 87
  ].join('');                                                                                                       // 88
                                                                                                                    // 89
  return Function(                                                                                                  // 90
    ['Promise', 'fn'],                                                                                              // 91
    body                                                                                                            // 92
  )(Promise, fn);                                                                                                   // 93
}                                                                                                                   // 94
                                                                                                                    // 95
Promise.nodeify = function (fn) {                                                                                   // 96
  return function () {                                                                                              // 97
    var args = Array.prototype.slice.call(arguments);                                                               // 98
    var callback =                                                                                                  // 99
      typeof args[args.length - 1] === 'function' ? args.pop() : null;                                              // 100
    var ctx = this;                                                                                                 // 101
    try {                                                                                                           // 102
      return fn.apply(this, arguments).nodeify(callback, ctx);                                                      // 103
    } catch (ex) {                                                                                                  // 104
      if (callback === null || typeof callback == 'undefined') {                                                    // 105
        return new Promise(function (resolve, reject) {                                                             // 106
          reject(ex);                                                                                               // 107
        });                                                                                                         // 108
      } else {                                                                                                      // 109
        asap(function () {                                                                                          // 110
          callback.call(ctx, ex);                                                                                   // 111
        })                                                                                                          // 112
      }                                                                                                             // 113
    }                                                                                                               // 114
  }                                                                                                                 // 115
}                                                                                                                   // 116
                                                                                                                    // 117
Promise.prototype.nodeify = function (callback, ctx) {                                                              // 118
  if (typeof callback != 'function') return this;                                                                   // 119
                                                                                                                    // 120
  this.then(function (value) {                                                                                      // 121
    asap(function () {                                                                                              // 122
      callback.call(ctx, null, value);                                                                              // 123
    });                                                                                                             // 124
  }, function (err) {                                                                                               // 125
    asap(function () {                                                                                              // 126
      callback.call(ctx, err);                                                                                      // 127
    });                                                                                                             // 128
  });                                                                                                               // 129
}                                                                                                                   // 130
                                                                                                                    // 131
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"synchronous.js":["./core.js",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/lib/synchronous.js                  //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
'use strict';                                                                                                       // 1
                                                                                                                    // 2
var Promise = require('./core.js');                                                                                 // 3
                                                                                                                    // 4
module.exports = Promise;                                                                                           // 5
Promise.enableSynchronous = function () {                                                                           // 6
  Promise.prototype.isPending = function() {                                                                        // 7
    return this.getState() == 0;                                                                                    // 8
  };                                                                                                                // 9
                                                                                                                    // 10
  Promise.prototype.isFulfilled = function() {                                                                      // 11
    return this.getState() == 1;                                                                                    // 12
  };                                                                                                                // 13
                                                                                                                    // 14
  Promise.prototype.isRejected = function() {                                                                       // 15
    return this.getState() == 2;                                                                                    // 16
  };                                                                                                                // 17
                                                                                                                    // 18
  Promise.prototype.getValue = function () {                                                                        // 19
    if (this._81 === 3) {                                                                                           // 20
      return this._65.getValue();                                                                                   // 21
    }                                                                                                               // 22
                                                                                                                    // 23
    if (!this.isFulfilled()) {                                                                                      // 24
      throw new Error('Cannot get a value of an unfulfilled promise.');                                             // 25
    }                                                                                                               // 26
                                                                                                                    // 27
    return this._65;                                                                                                // 28
  };                                                                                                                // 29
                                                                                                                    // 30
  Promise.prototype.getReason = function () {                                                                       // 31
    if (this._81 === 3) {                                                                                           // 32
      return this._65.getReason();                                                                                  // 33
    }                                                                                                               // 34
                                                                                                                    // 35
    if (!this.isRejected()) {                                                                                       // 36
      throw new Error('Cannot get a rejection reason of a non-rejected promise.');                                  // 37
    }                                                                                                               // 38
                                                                                                                    // 39
    return this._65;                                                                                                // 40
  };                                                                                                                // 41
                                                                                                                    // 42
  Promise.prototype.getState = function () {                                                                        // 43
    if (this._81 === 3) {                                                                                           // 44
      return this._65.getState();                                                                                   // 45
    }                                                                                                               // 46
    if (this._81 === -1 || this._81 === -2) {                                                                       // 47
      return 0;                                                                                                     // 48
    }                                                                                                               // 49
                                                                                                                    // 50
    return this._81;                                                                                                // 51
  };                                                                                                                // 52
};                                                                                                                  // 53
                                                                                                                    // 54
Promise.disableSynchronous = function() {                                                                           // 55
  Promise.prototype.isPending = undefined;                                                                          // 56
  Promise.prototype.isFulfilled = undefined;                                                                        // 57
  Promise.prototype.isRejected = undefined;                                                                         // 58
  Promise.prototype.getValue = undefined;                                                                           // 59
  Promise.prototype.getReason = undefined;                                                                          // 60
  Promise.prototype.getState = undefined;                                                                           // 61
};                                                                                                                  // 62
                                                                                                                    // 63
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}]},"node_modules":{"asap":{"raw.js":["domain",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/node_modules/asap/raw.js            //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
"use strict";                                                                                                       // 1
                                                                                                                    // 2
var domain; // The domain module is executed on demand                                                              // 3
var hasSetImmediate = typeof setImmediate === "function";                                                           // 4
                                                                                                                    // 5
// Use the fastest means possible to execute a task in its own turn, with                                           // 6
// priority over other events including network IO events in Node.js.                                               // 7
//                                                                                                                  // 8
// An exception thrown by a task will permanently interrupt the processing of                                       // 9
// subsequent tasks. The higher level `asap` function ensures that if an                                            // 10
// exception is thrown by a task, that the task queue will continue flushing as                                     // 11
// soon as possible, but if you use `rawAsap` directly, you are responsible to                                      // 12
// either ensure that no exceptions are thrown from your task, or to manually                                       // 13
// call `rawAsap.requestFlush` if an exception is thrown.                                                           // 14
module.exports = rawAsap;                                                                                           // 15
function rawAsap(task) {                                                                                            // 16
    if (!queue.length) {                                                                                            // 17
        requestFlush();                                                                                             // 18
        flushing = true;                                                                                            // 19
    }                                                                                                               // 20
    // Avoids a function call                                                                                       // 21
    queue[queue.length] = task;                                                                                     // 22
}                                                                                                                   // 23
                                                                                                                    // 24
var queue = [];                                                                                                     // 25
// Once a flush has been requested, no further calls to `requestFlush` are                                          // 26
// necessary until the next `flush` completes.                                                                      // 27
var flushing = false;                                                                                               // 28
// The position of the next task to execute in the task queue. This is                                              // 29
// preserved between calls to `flush` so that it can be resumed if                                                  // 30
// a task throws an exception.                                                                                      // 31
var index = 0;                                                                                                      // 32
// If a task schedules additional tasks recursively, the task queue can grow                                        // 33
// unbounded. To prevent memory excaustion, the task queue will periodically                                        // 34
// truncate already-completed tasks.                                                                                // 35
var capacity = 1024;                                                                                                // 36
                                                                                                                    // 37
// The flush function processes all tasks that have been scheduled with                                             // 38
// `rawAsap` unless and until one of those tasks throws an exception.                                               // 39
// If a task throws an exception, `flush` ensures that its state will remain                                        // 40
// consistent and will resume where it left off when called again.                                                  // 41
// However, `flush` does not make any arrangements to be called again if an                                         // 42
// exception is thrown.                                                                                             // 43
function flush() {                                                                                                  // 44
    while (index < queue.length) {                                                                                  // 45
        var currentIndex = index;                                                                                   // 46
        // Advance the index before calling the task. This ensures that we will                                     // 47
        // begin flushing on the next task the task throws an error.                                                // 48
        index = index + 1;                                                                                          // 49
        queue[currentIndex].call();                                                                                 // 50
        // Prevent leaking memory for long chains of recursive calls to `asap`.                                     // 51
        // If we call `asap` within tasks scheduled by `asap`, the queue will                                       // 52
        // grow, but to avoid an O(n) walk for every task we execute, we don't                                      // 53
        // shift tasks off the queue after they have been executed.                                                 // 54
        // Instead, we periodically shift 1024 tasks off the queue.                                                 // 55
        if (index > capacity) {                                                                                     // 56
            // Manually shift all values starting at the index back to the                                          // 57
            // beginning of the queue.                                                                              // 58
            for (var scan = 0, newLength = queue.length - index; scan < newLength; scan++) {                        // 59
                queue[scan] = queue[scan + index];                                                                  // 60
            }                                                                                                       // 61
            queue.length -= index;                                                                                  // 62
            index = 0;                                                                                              // 63
        }                                                                                                           // 64
    }                                                                                                               // 65
    queue.length = 0;                                                                                               // 66
    index = 0;                                                                                                      // 67
    flushing = false;                                                                                               // 68
}                                                                                                                   // 69
                                                                                                                    // 70
rawAsap.requestFlush = requestFlush;                                                                                // 71
function requestFlush() {                                                                                           // 72
    // Ensure flushing is not bound to any domain.                                                                  // 73
    // It is not sufficient to exit the domain, because domains exist on a stack.                                   // 74
    // To execute code outside of any domain, the following dance is necessary.                                     // 75
    var parentDomain = process.domain;                                                                              // 76
    if (parentDomain) {                                                                                             // 77
        if (!domain) {                                                                                              // 78
            // Lazy execute the domain module.                                                                      // 79
            // Only employed if the user elects to use domains.                                                     // 80
            domain = require("domain");                                                                             // 81
        }                                                                                                           // 82
        domain.active = process.domain = null;                                                                      // 83
    }                                                                                                               // 84
                                                                                                                    // 85
    // `setImmediate` is slower that `process.nextTick`, but `process.nextTick`                                     // 86
    // cannot handle recursion.                                                                                     // 87
    // `requestFlush` will only be called recursively from `asap.js`, to resume                                     // 88
    // flushing after an error is thrown into a domain.                                                             // 89
    // Conveniently, `setImmediate` was introduced in the same version                                              // 90
    // `process.nextTick` started throwing recursion errors.                                                        // 91
    if (flushing && hasSetImmediate) {                                                                              // 92
        setImmediate(flush);                                                                                        // 93
    } else {                                                                                                        // 94
        process.nextTick(flush);                                                                                    // 95
    }                                                                                                               // 96
                                                                                                                    // 97
    if (parentDomain) {                                                                                             // 98
        domain.active = process.domain = parentDomain;                                                              // 99
    }                                                                                                               // 100
}                                                                                                                   // 101
                                                                                                                    // 102
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"package.json":function(require,exports){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// ../npm/node_modules/meteor-promise/node_modules/promise/node_modules/asap/package.json                           //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
exports.name = "asap";                                                                                              // 1
exports.version = "2.0.3";                                                                                          // 2
exports.main = "./asap.js";                                                                                         // 3
                                                                                                                    // 4
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"asap.js":["./raw",function(require,exports,module){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// node_modules/meteor/promise/node_modules/meteor-promise/node_modules/promise/node_modules/asap/asap.js           //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
"use strict";                                                                                                       // 1
                                                                                                                    // 2
var rawAsap = require("./raw");                                                                                     // 3
var freeTasks = [];                                                                                                 // 4
                                                                                                                    // 5
/**                                                                                                                 // 6
 * Calls a task as soon as possible after returning, in its own event, with                                         // 7
 * priority over IO events. An exception thrown in a task can be handled by                                         // 8
 * `process.on("uncaughtException") or `domain.on("error")`, but will otherwise                                     // 9
 * crash the process. If the error is handled, all subsequent tasks will                                            // 10
 * resume.                                                                                                          // 11
 *                                                                                                                  // 12
 * @param {{call}} task A callable object, typically a function that takes no                                       // 13
 * arguments.                                                                                                       // 14
 */                                                                                                                 // 15
module.exports = asap;                                                                                              // 16
function asap(task) {                                                                                               // 17
    var rawTask;                                                                                                    // 18
    if (freeTasks.length) {                                                                                         // 19
        rawTask = freeTasks.pop();                                                                                  // 20
    } else {                                                                                                        // 21
        rawTask = new RawTask();                                                                                    // 22
    }                                                                                                               // 23
    rawTask.task = task;                                                                                            // 24
    rawTask.domain = process.domain;                                                                                // 25
    rawAsap(rawTask);                                                                                               // 26
}                                                                                                                   // 27
                                                                                                                    // 28
function RawTask() {                                                                                                // 29
    this.task = null;                                                                                               // 30
    this.domain = null;                                                                                             // 31
}                                                                                                                   // 32
                                                                                                                    // 33
RawTask.prototype.call = function () {                                                                              // 34
    if (this.domain) {                                                                                              // 35
        this.domain.enter();                                                                                        // 36
    }                                                                                                               // 37
    var threw = true;                                                                                               // 38
    try {                                                                                                           // 39
        this.task.call();                                                                                           // 40
        threw = false;                                                                                              // 41
        // If the task throws an exception (presumably) Node.js restores the                                        // 42
        // domain stack for the next event.                                                                         // 43
        if (this.domain) {                                                                                          // 44
            this.domain.exit();                                                                                     // 45
        }                                                                                                           // 46
    } finally {                                                                                                     // 47
        // We use try/finally and a threw flag to avoid messing up stack traces                                     // 48
        // when we catch and release errors.                                                                        // 49
        if (threw) {                                                                                                // 50
            // In Node.js, uncaught exceptions are considered fatal errors.                                         // 51
            // Re-throw them to interrupt flushing!                                                                 // 52
            // Ensure that flushing continues if an uncaught exception is                                           // 53
            // suppressed listening process.on("uncaughtException") or                                              // 54
            // domain.on("error").                                                                                  // 55
            rawAsap.requestFlush();                                                                                 // 56
        }                                                                                                           // 57
        // If the task threw an error, we do not want to exit the domain here.                                      // 58
        // Exiting the domain would prevent the domain from catching the error.                                     // 59
        this.task = null;                                                                                           // 60
        this.domain = null;                                                                                         // 61
        freeTasks.push(this);                                                                                       // 62
    }                                                                                                               // 63
};                                                                                                                  // 64
                                                                                                                    // 65
                                                                                                                    // 66
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}]}}}}}}}}}},{"extensions":[".js",".json"]});
var exports = require("./node_modules/meteor/promise/client.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.promise = exports, {
  Promise: Promise
});

})();
