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
var _ = Package.underscore._;
var Tracker = Package.tracker.Tracker;
var Deps = Package.tracker.Deps;
var EJSON = Package.ejson.EJSON;
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;
var Symbol = Package['ecmascript-runtime'].Symbol;
var Map = Package['ecmascript-runtime'].Map;
var Set = Package['ecmascript-runtime'].Set;
var meteorBabelHelpers = Package['babel-runtime'].meteorBabelHelpers;
var Promise = Package.promise.Promise;

/* Package-scope variables */
var ReactiveDict;

var require = meteorInstall({"node_modules":{"meteor":{"reactive-dict":{"reactive-dict.js":["babel-runtime/helpers/typeof",function(require){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// packages/reactive-dict/reactive-dict.js                                                                          //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
var _typeof2 = require('babel-runtime/helpers/typeof');                                                             //
                                                                                                                    //
var _typeof3 = _interopRequireDefault(_typeof2);                                                                    //
                                                                                                                    //
function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { 'default': obj }; }                   //
                                                                                                                    //
// XXX come up with a serialization method which canonicalizes object key                                           //
// order, which would allow us to use objects as values for equals.                                                 //
var stringify = function stringify(value) {                                                                         // 3
  if (value === undefined) return 'undefined';                                                                      // 4
  return EJSON.stringify(value);                                                                                    // 6
};                                                                                                                  //
var parse = function parse(serialized) {                                                                            // 8
  if (serialized === undefined || serialized === 'undefined') return undefined;                                     // 9
  return EJSON.parse(serialized);                                                                                   // 11
};                                                                                                                  //
                                                                                                                    //
var changed = function changed(v) {                                                                                 // 14
  v && v.changed();                                                                                                 // 15
};                                                                                                                  //
                                                                                                                    //
// XXX COMPAT WITH 0.9.1 : accept migrationData instead of dictName                                                 //
ReactiveDict = function (_ReactiveDict) {                                                                           // 19
  function ReactiveDict(_x) {                                                                                       //
    return _ReactiveDict.apply(this, arguments);                                                                    //
  }                                                                                                                 //
                                                                                                                    //
  ReactiveDict.toString = function () {                                                                             //
    return _ReactiveDict.toString();                                                                                //
  };                                                                                                                //
                                                                                                                    //
  return ReactiveDict;                                                                                              //
}(function (dictName) {                                                                                             //
  // this.keys: key -> value                                                                                        //
  if (dictName) {                                                                                                   // 21
    if (typeof dictName === 'string') {                                                                             // 22
      // the normal case, argument is a string name.                                                                //
      // _registerDictForMigrate will throw an error on duplicate name.                                             //
      ReactiveDict._registerDictForMigrate(dictName, this);                                                         // 25
      this.keys = ReactiveDict._loadMigratedDict(dictName) || {};                                                   // 26
      this.name = dictName;                                                                                         // 27
    } else if ((typeof dictName === 'undefined' ? 'undefined' : (0, _typeof3['default'])(dictName)) === 'object') {
      // back-compat case: dictName is actually migrationData                                                       //
      this.keys = dictName;                                                                                         // 30
    } else {                                                                                                        //
      throw new Error("Invalid ReactiveDict argument: " + dictName);                                                // 32
    }                                                                                                               //
  } else {                                                                                                          //
    // no name given; no migration will be performed                                                                //
    this.keys = {};                                                                                                 // 36
  }                                                                                                                 //
                                                                                                                    //
  this.allDeps = new Tracker.Dependency();                                                                          // 39
  this.keyDeps = {}; // key -> Dependency                                                                           // 40
  this.keyValueDeps = {}; // key -> Dependency                                                                      // 19
});                                                                                                                 // 19
                                                                                                                    //
_.extend(ReactiveDict.prototype, {                                                                                  // 44
  // set() began as a key/value method, but we are now overloading it                                               //
  // to take an object of key/value pairs, similar to backbone                                                      //
  // http://backbonejs.org/#Model-set                                                                               //
                                                                                                                    //
  set: function () {                                                                                                // 49
    function set(keyOrObject, value) {                                                                              // 49
      var self = this;                                                                                              // 50
                                                                                                                    //
      if ((typeof keyOrObject === 'undefined' ? 'undefined' : (0, _typeof3['default'])(keyOrObject)) === 'object' && value === undefined) {
        // Called as `dict.set({...})`                                                                              //
        self._setObject(keyOrObject);                                                                               // 54
        return;                                                                                                     // 55
      }                                                                                                             //
      // the input isn't an object, so it must be a key                                                             //
      // and we resume with the rest of the function                                                                //
      var key = keyOrObject;                                                                                        // 49
                                                                                                                    //
      value = stringify(value);                                                                                     // 61
                                                                                                                    //
      var keyExisted = _.has(self.keys, key);                                                                       // 63
      var oldSerializedValue = keyExisted ? self.keys[key] : 'undefined';                                           // 64
      var isNewValue = value !== oldSerializedValue;                                                                // 65
                                                                                                                    //
      self.keys[key] = value;                                                                                       // 67
                                                                                                                    //
      if (isNewValue || !keyExisted) {                                                                              // 69
        self.allDeps.changed();                                                                                     // 70
      }                                                                                                             //
                                                                                                                    //
      if (isNewValue) {                                                                                             // 73
        changed(self.keyDeps[key]);                                                                                 // 74
        if (self.keyValueDeps[key]) {                                                                               // 75
          changed(self.keyValueDeps[key][oldSerializedValue]);                                                      // 76
          changed(self.keyValueDeps[key][value]);                                                                   // 77
        }                                                                                                           //
      }                                                                                                             //
    }                                                                                                               //
                                                                                                                    //
    return set;                                                                                                     //
  }(),                                                                                                              //
                                                                                                                    //
  setDefault: function () {                                                                                         // 82
    function setDefault(key, value) {                                                                               // 82
      var self = this;                                                                                              // 83
      if (!_.has(self.keys, key)) {                                                                                 // 84
        self.set(key, value);                                                                                       // 85
      }                                                                                                             //
    }                                                                                                               //
                                                                                                                    //
    return setDefault;                                                                                              //
  }(),                                                                                                              //
                                                                                                                    //
  get: function () {                                                                                                // 89
    function get(key) {                                                                                             // 89
      var self = this;                                                                                              // 90
      self._ensureKey(key);                                                                                         // 91
      self.keyDeps[key].depend();                                                                                   // 92
      return parse(self.keys[key]);                                                                                 // 93
    }                                                                                                               //
                                                                                                                    //
    return get;                                                                                                     //
  }(),                                                                                                              //
                                                                                                                    //
  equals: function () {                                                                                             // 96
    function equals(key, value) {                                                                                   // 96
      var self = this;                                                                                              // 97
                                                                                                                    //
      // Mongo.ObjectID is in the 'mongo' package                                                                   //
      var ObjectID = null;                                                                                          // 96
      if (Package.mongo) {                                                                                          // 101
        ObjectID = Package.mongo.Mongo.ObjectID;                                                                    // 102
      }                                                                                                             //
                                                                                                                    //
      // We don't allow objects (or arrays that might include objects) for                                          //
      // .equals, because JSON.stringify doesn't canonicalize object key                                            //
      // order. (We can make equals have the right return value by parsing the                                      //
      // current value and using EJSON.equals, but we won't have a canonical                                        //
      // element of keyValueDeps[key] to store the dependency.) You can still use                                   //
      // "EJSON.equals(reactiveDict.get(key), value)".                                                              //
      //                                                                                                            //
      // XXX we could allow arrays as long as we recursively check that there                                       //
      // are no objects                                                                                             //
      if (typeof value !== 'string' && typeof value !== 'number' && typeof value !== 'boolean' && typeof value !== 'undefined' && !(value instanceof Date) && !(ObjectID && value instanceof ObjectID) && value !== null) {
        throw new Error("ReactiveDict.equals: value must be scalar");                                               // 121
      }                                                                                                             //
      var serializedValue = stringify(value);                                                                       // 123
                                                                                                                    //
      if (Tracker.active) {                                                                                         // 125
        self._ensureKey(key);                                                                                       // 126
                                                                                                                    //
        if (!_.has(self.keyValueDeps[key], serializedValue)) self.keyValueDeps[key][serializedValue] = new Tracker.Dependency();
                                                                                                                    //
        var isNew = self.keyValueDeps[key][serializedValue].depend();                                               // 131
        if (isNew) {                                                                                                // 132
          Tracker.onInvalidate(function () {                                                                        // 133
            // clean up [key][serializedValue] if it's now empty, so we don't                                       //
            // use O(n) memory for n = values seen ever                                                             //
            if (!self.keyValueDeps[key][serializedValue].hasDependents()) delete self.keyValueDeps[key][serializedValue];
          });                                                                                                       //
        }                                                                                                           //
      }                                                                                                             //
                                                                                                                    //
      var oldValue = undefined;                                                                                     // 142
      if (_.has(self.keys, key)) oldValue = parse(self.keys[key]);                                                  // 143
      return EJSON.equals(oldValue, value);                                                                         // 144
    }                                                                                                               //
                                                                                                                    //
    return equals;                                                                                                  //
  }(),                                                                                                              //
                                                                                                                    //
  all: function () {                                                                                                // 147
    function all() {                                                                                                // 147
      this.allDeps.depend();                                                                                        // 148
      var ret = {};                                                                                                 // 149
      _.each(this.keys, function (value, key) {                                                                     // 150
        ret[key] = parse(value);                                                                                    // 151
      });                                                                                                           //
      return ret;                                                                                                   // 153
    }                                                                                                               //
                                                                                                                    //
    return all;                                                                                                     //
  }(),                                                                                                              //
                                                                                                                    //
  clear: function () {                                                                                              // 156
    function clear() {                                                                                              // 156
      var self = this;                                                                                              // 157
                                                                                                                    //
      var oldKeys = self.keys;                                                                                      // 159
      self.keys = {};                                                                                               // 160
                                                                                                                    //
      self.allDeps.changed();                                                                                       // 162
                                                                                                                    //
      _.each(oldKeys, function (value, key) {                                                                       // 164
        changed(self.keyDeps[key]);                                                                                 // 165
        if (self.keyValueDeps[key]) {                                                                               // 166
          changed(self.keyValueDeps[key][value]);                                                                   // 167
          changed(self.keyValueDeps[key]['undefined']);                                                             // 168
        }                                                                                                           //
      });                                                                                                           //
    }                                                                                                               //
                                                                                                                    //
    return clear;                                                                                                   //
  }(),                                                                                                              //
                                                                                                                    //
  'delete': function () {                                                                                           // 174
    function _delete(key) {                                                                                         // 174
      var self = this;                                                                                              // 175
      var didRemove = false;                                                                                        // 176
                                                                                                                    //
      if (_.has(self.keys, key)) {                                                                                  // 178
        var oldValue = self.keys[key];                                                                              // 179
        delete self.keys[key];                                                                                      // 180
        changed(self.keyDeps[key]);                                                                                 // 181
        if (self.keyValueDeps[key]) {                                                                               // 182
          changed(self.keyValueDeps[key][oldValue]);                                                                // 183
          changed(self.keyValueDeps[key]['undefined']);                                                             // 184
        }                                                                                                           //
        self.allDeps.changed();                                                                                     // 186
        didRemove = true;                                                                                           // 187
      }                                                                                                             //
                                                                                                                    //
      return didRemove;                                                                                             // 190
    }                                                                                                               //
                                                                                                                    //
    return _delete;                                                                                                 //
  }(),                                                                                                              //
                                                                                                                    //
  _setObject: function () {                                                                                         // 193
    function _setObject(object) {                                                                                   // 193
      var self = this;                                                                                              // 194
                                                                                                                    //
      _.each(object, function (value, key) {                                                                        // 196
        self.set(key, value);                                                                                       // 197
      });                                                                                                           //
    }                                                                                                               //
                                                                                                                    //
    return _setObject;                                                                                              //
  }(),                                                                                                              //
                                                                                                                    //
  _ensureKey: function () {                                                                                         // 201
    function _ensureKey(key) {                                                                                      // 201
      var self = this;                                                                                              // 202
      if (!(key in self.keyDeps)) {                                                                                 // 203
        self.keyDeps[key] = new Tracker.Dependency();                                                               // 204
        self.keyValueDeps[key] = {};                                                                                // 205
      }                                                                                                             //
    }                                                                                                               //
                                                                                                                    //
    return _ensureKey;                                                                                              //
  }(),                                                                                                              //
                                                                                                                    //
  // Get a JSON value that can be passed to the constructor to                                                      //
  // create a new ReactiveDict with the same contents as this one                                                   //
  _getMigrationData: function () {                                                                                  // 211
    function _getMigrationData() {                                                                                  // 211
      // XXX sanitize and make sure it's JSONible?                                                                  //
      return this.keys;                                                                                             // 213
    }                                                                                                               //
                                                                                                                    //
    return _getMigrationData;                                                                                       //
  }()                                                                                                               //
});                                                                                                                 //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"migration.js":function(){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// packages/reactive-dict/migration.js                                                                              //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
ReactiveDict._migratedDictData = {}; // name -> data                                                                // 1
ReactiveDict._dictsToMigrate = {}; // name -> ReactiveDict                                                          // 2
                                                                                                                    //
ReactiveDict._loadMigratedDict = function (dictName) {                                                              // 4
  if (_.has(ReactiveDict._migratedDictData, dictName)) return ReactiveDict._migratedDictData[dictName];             // 5
                                                                                                                    //
  return null;                                                                                                      // 8
};                                                                                                                  //
                                                                                                                    //
ReactiveDict._registerDictForMigrate = function (dictName, dict) {                                                  // 11
  if (_.has(ReactiveDict._dictsToMigrate, dictName)) throw new Error("Duplicate ReactiveDict name: " + dictName);   // 12
                                                                                                                    //
  ReactiveDict._dictsToMigrate[dictName] = dict;                                                                    // 15
};                                                                                                                  //
                                                                                                                    //
if (Meteor.isClient && Package.reload) {                                                                            // 18
  // Put old migrated data into ReactiveDict._migratedDictData,                                                     //
  // where it can be accessed by ReactiveDict._loadMigratedDict.                                                    //
  var migrationData = Package.reload.Reload._migrationData('reactive-dict');                                        // 21
  if (migrationData && migrationData.dicts) ReactiveDict._migratedDictData = migrationData.dicts;                   // 22
                                                                                                                    //
  // On migration, assemble the data from all the dicts that have been                                              //
  // registered.                                                                                                    //
  Package.reload.Reload._onMigrate('reactive-dict', function () {                                                   // 18
    var dictsToMigrate = ReactiveDict._dictsToMigrate;                                                              // 28
    var dataToMigrate = {};                                                                                         // 29
                                                                                                                    //
    for (var dictName in meteorBabelHelpers.sanitizeForInObject(dictsToMigrate)) {                                  // 31
      dataToMigrate[dictName] = dictsToMigrate[dictName]._getMigrationData();                                       // 32
    }return [true, { dicts: dataToMigrate }];                                                                       //
  });                                                                                                               //
}                                                                                                                   //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}},{"extensions":[".js",".json"]});
require("./node_modules/meteor/reactive-dict/reactive-dict.js");
require("./node_modules/meteor/reactive-dict/migration.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['reactive-dict'] = {}, {
  ReactiveDict: ReactiveDict
});

})();
