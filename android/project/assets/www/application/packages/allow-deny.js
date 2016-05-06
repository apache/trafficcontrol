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
var LocalCollection = Package.minimongo.LocalCollection;
var Minimongo = Package.minimongo.Minimongo;
var check = Package.check.check;
var Match = Package.check.Match;
var EJSON = Package.ejson.EJSON;
var DDP = Package['ddp-client'].DDP;
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;
var Symbol = Package['ecmascript-runtime'].Symbol;
var Map = Package['ecmascript-runtime'].Map;
var Set = Package['ecmascript-runtime'].Set;
var meteorBabelHelpers = Package['babel-runtime'].meteorBabelHelpers;
var Promise = Package.promise.Promise;

/* Package-scope variables */
var AllowDeny;

var require = meteorInstall({"node_modules":{"meteor":{"allow-deny":{"allow-deny.js":function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/allow-deny/allow-deny.js                                                                              //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
///                                                                                                               //
/// Remote methods and access control.                                                                            //
///                                                                                                               //
                                                                                                                  //
// Restrict default mutators on collection. allow() and deny() take the                                           //
// same options:                                                                                                  //
//                                                                                                                //
// options.insert {Function(userId, doc)}                                                                         //
//   return true to allow/deny adding this document                                                               //
//                                                                                                                //
// options.update {Function(userId, docs, fields, modifier)}                                                      //
//   return true to allow/deny updating these documents.                                                          //
//   `fields` is passed as an array of fields that are to be modified                                             //
//                                                                                                                //
// options.remove {Function(userId, docs)}                                                                        //
//   return true to allow/deny removing these documents                                                           //
//                                                                                                                //
// options.fetch {Array}                                                                                          //
//   Fields to fetch for these validators. If any call to allow or deny                                           //
//   does not have this option then all fields are loaded.                                                        //
//                                                                                                                //
// allow and deny can be called multiple times. The validators are                                                //
// evaluated as follows:                                                                                          //
// - If neither deny() nor allow() has been called on the collection,                                             //
//   then the request is allowed if and only if the "insecure" smart                                              //
//   package is in use.                                                                                           //
// - Otherwise, if any deny() function returns true, the request is denied.                                       //
// - Otherwise, if any allow() function returns true, the request is allowed.                                     //
// - Otherwise, the request is denied.                                                                            //
//                                                                                                                //
// Meteor may call your deny() and allow() functions in any order, and may not                                    //
// call all of them if it is able to make a decision without calling them all                                     //
// (so don't include side effects).                                                                               //
                                                                                                                  //
AllowDeny = {                                                                                                     // 35
  CollectionPrototype: {}                                                                                         // 36
};                                                                                                                //
                                                                                                                  //
// In the `mongo` package, we will extend Mongo.Collection.prototype with these                                   //
// methods                                                                                                        //
var CollectionPrototype = AllowDeny.CollectionPrototype;                                                          // 41
                                                                                                                  //
/**                                                                                                               //
 * @summary Allow users to write directly to this collection from client code, subject to limitations you define.
 * @locus Server                                                                                                  //
 * @method allow                                                                                                  //
 * @memberOf Mongo.Collection                                                                                     //
 * @instance                                                                                                      //
 * @param {Object} options                                                                                        //
 * @param {Function} options.insert,update,remove Functions that look at a proposed modification to the database and return true if it should be allowed.
 * @param {String[]} options.fetch Optional performance enhancement. Limits the fields that will be fetched from the database for inspection by your `update` and `remove` functions.
 * @param {Function} options.transform Overrides `transform` on the  [`Collection`](#collections).  Pass `null` to disable transformation.
 */                                                                                                               //
CollectionPrototype.allow = function (options) {                                                                  // 54
  addValidator(this, 'allow', options);                                                                           // 55
};                                                                                                                //
                                                                                                                  //
/**                                                                                                               //
 * @summary Override `allow` rules.                                                                               //
 * @locus Server                                                                                                  //
 * @method deny                                                                                                   //
 * @memberOf Mongo.Collection                                                                                     //
 * @instance                                                                                                      //
 * @param {Object} options                                                                                        //
 * @param {Function} options.insert,update,remove Functions that look at a proposed modification to the database and return true if it should be denied, even if an [allow](#allow) rule says otherwise.
 * @param {String[]} options.fetch Optional performance enhancement. Limits the fields that will be fetched from the database for inspection by your `update` and `remove` functions.
 * @param {Function} options.transform Overrides `transform` on the  [`Collection`](#collections).  Pass `null` to disable transformation.
 */                                                                                                               //
CollectionPrototype.deny = function (options) {                                                                   // 69
  addValidator(this, 'deny', options);                                                                            // 70
};                                                                                                                //
                                                                                                                  //
CollectionPrototype._defineMutationMethods = function () {                                                        // 73
  var self = this;                                                                                                // 74
                                                                                                                  //
  // set to true once we call any allow or deny methods. If true, use                                             //
  // allow/deny semantics. If false, use insecure mode semantics.                                                 //
  self._restricted = false;                                                                                       // 73
                                                                                                                  //
  // Insecure mode (default to allowing writes). Defaults to 'undefined' which                                    //
  // means insecure iff the insecure package is loaded. This property can be                                      //
  // overriden by tests or packages wishing to change insecure mode behavior of                                   //
  // their collections.                                                                                           //
  self._insecure = undefined;                                                                                     // 73
                                                                                                                  //
  self._validators = {                                                                                            // 86
    insert: { allow: [], deny: [] },                                                                              // 87
    update: { allow: [], deny: [] },                                                                              // 88
    remove: { allow: [], deny: [] },                                                                              // 89
    upsert: { allow: [], deny: [] }, // dummy arrays; can't set these!                                            // 90
    fetch: [],                                                                                                    // 91
    fetchAllFields: false                                                                                         // 92
  };                                                                                                              //
                                                                                                                  //
  if (!self._name) return; // anonymous collection                                                                // 95
                                                                                                                  //
  // XXX Think about method namespacing. Maybe methods should be                                                  //
  // "Meteor:Mongo:insert/NAME"?                                                                                  //
  self._prefix = '/' + self._name + '/';                                                                          // 73
                                                                                                                  //
  // mutation methods                                                                                             //
  if (self._connection) {                                                                                         // 73
    (function () {                                                                                                //
      var m = {};                                                                                                 // 104
                                                                                                                  //
      _.each(['insert', 'update', 'remove'], function (method) {                                                  // 106
        m[self._prefix + method] = function () /* ... */{                                                         // 107
          // All the methods do their own validation, instead of using check().                                   //
          check(arguments, [Match.Any]);                                                                          // 109
          var args = _.toArray(arguments);                                                                        // 110
          try {                                                                                                   // 111
            // For an insert, if the client didn't specify an _id, generate one                                   //
            // now; because this uses DDP.randomStream, it will be consistent with                                //
            // what the client generated. We generate it now rather than later so                                 //
            // that if (eg) an allow/deny rule does an insert to the same                                         //
            // collection (not that it really should), the generated _id will                                     //
            // still be the first use of the stream and will be consistent.                                       //
            //                                                                                                    //
            // However, we don't actually stick the _id onto the document yet,                                    //
            // because we want allow/deny rules to be able to differentiate                                       //
            // between arbitrary client-specified _id fields and merely                                           //
            // client-controlled-via-randomSeed fields.                                                           //
            var generatedId = null;                                                                               // 123
            if (method === "insert" && !_.has(args[0], '_id')) {                                                  // 124
              generatedId = self._makeNewID();                                                                    // 125
            }                                                                                                     //
                                                                                                                  //
            if (this.isSimulation) {                                                                              // 128
              // In a client simulation, you can do any mutation (even with a                                     //
              // complex selector).                                                                               //
              if (generatedId !== null) args[0]._id = generatedId;                                                // 131
              return self._collection[method].apply(self._collection, args);                                      // 133
            }                                                                                                     //
                                                                                                                  //
            // This is the server receiving a method call from the client.                                        //
                                                                                                                  //
            // We don't allow arbitrary selectors in mutations from the client: only                              //
            // single-ID selectors.                                                                               //
            if (method !== 'insert') throwIfSelectorIsNotId(args[0], method);                                     // 111
                                                                                                                  //
            if (self._restricted) {                                                                               // 144
              // short circuit if there is no way it will pass.                                                   //
              if (self._validators[method].allow.length === 0) {                                                  // 146
                throw new Meteor.Error(403, "Access denied. No allow validators set on restricted " + "collection for method '" + method + "'.");
              }                                                                                                   //
                                                                                                                  //
              var validatedMethodName = '_validated' + method.charAt(0).toUpperCase() + method.slice(1);          // 152
              args.unshift(this.userId);                                                                          // 154
              method === 'insert' && args.push(generatedId);                                                      // 155
              return self[validatedMethodName].apply(self, args);                                                 // 156
            } else if (self._isInsecure()) {                                                                      //
              if (generatedId !== null) args[0]._id = generatedId;                                                // 158
              // In insecure mode, allow any mutation (with a simple selector).                                   //
              // XXX This is kind of bogus.  Instead of blindly passing whatever                                  //
              //     we get from the network to this function, we should actually                                 //
              //     know the correct arguments for the function and pass just                                    //
              //     them.  For example, if you have an extraneous extra null                                     //
              //     argument and this is Mongo on the server, the .wrapAsync'd                                   //
              //     functions like update will get confused and pass the                                         //
              //     "fut.resolver()" in the wrong slot, where _update will never                                 //
              //     invoke it. Bam, broken DDP connection.  Probably should just                                 //
              //     take this whole method and write it three times, invoking                                    //
              //     helpers for the common code.                                                                 //
              return self._collection[method].apply(self._collection, args);                                      // 157
            } else {                                                                                              //
              // In secure mode, if we haven't called allow or deny, then nothing                                 //
              // is permitted.                                                                                    //
              throw new Meteor.Error(403, "Access denied");                                                       // 175
            }                                                                                                     //
          } catch (e) {                                                                                           //
            if (e.name === 'MongoError' || e.name === 'MinimongoError') {                                         // 178
              throw new Meteor.Error(409, e.toString());                                                          // 179
            } else {                                                                                              //
              throw e;                                                                                            // 181
            }                                                                                                     //
          }                                                                                                       //
        };                                                                                                        //
      });                                                                                                         //
      // Minimongo on the server gets no stubs; instead, by default                                               //
      // it wait()s until its result is ready, yielding.                                                          //
      // This matches the behavior of macromongo on the server better.                                            //
      // XXX see #MeteorServerNull                                                                                //
      if (Meteor.isClient || self._connection === Meteor.server) self._connection.methods(m);                     // 190
    })();                                                                                                         //
  }                                                                                                               //
};                                                                                                                //
                                                                                                                  //
CollectionPrototype._updateFetch = function (fields) {                                                            // 195
  var self = this;                                                                                                // 196
                                                                                                                  //
  if (!self._validators.fetchAllFields) {                                                                         // 198
    if (fields) {                                                                                                 // 199
      self._validators.fetch = _.union(self._validators.fetch, fields);                                           // 200
    } else {                                                                                                      //
      self._validators.fetchAllFields = true;                                                                     // 202
      // clear fetch just to make sure we don't accidentally read it                                              //
      self._validators.fetch = null;                                                                              // 201
    }                                                                                                             //
  }                                                                                                               //
};                                                                                                                //
                                                                                                                  //
CollectionPrototype._isInsecure = function () {                                                                   // 209
  var self = this;                                                                                                // 210
  if (self._insecure === undefined) return !!Package.insecure;                                                    // 211
  return self._insecure;                                                                                          // 213
};                                                                                                                //
                                                                                                                  //
CollectionPrototype._validatedInsert = function (userId, doc, generatedId) {                                      // 216
  var self = this;                                                                                                // 218
                                                                                                                  //
  // call user validators.                                                                                        //
  // Any deny returns true means denied.                                                                          //
  if (_.any(self._validators.insert.deny, function (validator) {                                                  // 217
    return validator(userId, docToValidate(validator, doc, generatedId));                                         // 223
  })) {                                                                                                           //
    throw new Meteor.Error(403, "Access denied");                                                                 // 225
  }                                                                                                               //
  // Any allow returns true means proceed. Throw error if they all fail.                                          //
  if (_.all(self._validators.insert.allow, function (validator) {                                                 // 217
    return !validator(userId, docToValidate(validator, doc, generatedId));                                        // 229
  })) {                                                                                                           //
    throw new Meteor.Error(403, "Access denied");                                                                 // 231
  }                                                                                                               //
                                                                                                                  //
  // If we generated an ID above, insert it now: after the validation, but                                        //
  // before actually inserting.                                                                                   //
  if (generatedId !== null) doc._id = generatedId;                                                                // 217
                                                                                                                  //
  self._collection.insert.call(self._collection, doc);                                                            // 239
};                                                                                                                //
                                                                                                                  //
// Simulate a mongo `update` operation while validating that the access                                           //
// control rules set by calls to `allow/deny` are satisfied. If all                                               //
// pass, rewrite the mongo operation to use $in to set the list of                                                //
// document ids to change ##ValidatedChange                                                                       //
CollectionPrototype._validatedUpdate = function (userId, selector, mutator, options) {                            // 246
  var self = this;                                                                                                // 248
                                                                                                                  //
  check(mutator, Object);                                                                                         // 250
                                                                                                                  //
  options = _.clone(options) || {};                                                                               // 252
                                                                                                                  //
  if (!LocalCollection._selectorIsIdPerhapsAsObject(selector)) throw new Error("validated update should be of a single ID");
                                                                                                                  //
  // We don't support upserts because they don't fit nicely into allow/deny                                       //
  // rules.                                                                                                       //
  if (options.upsert) throw new Meteor.Error(403, "Access denied. Upserts not " + "allowed in a restricted collection.");
                                                                                                                  //
  var noReplaceError = "Access denied. In a restricted collection you can only" + " update documents, not replace them. Use a Mongo update operator, such " + "as '$set'.";
                                                                                                                  //
  // compute modified fields                                                                                      //
  var fields = [];                                                                                                // 247
  if (_.isEmpty(mutator)) {                                                                                       // 269
    throw new Meteor.Error(403, noReplaceError);                                                                  // 270
  }                                                                                                               //
  _.each(mutator, function (params, op) {                                                                         // 272
    if (op.charAt(0) !== '$') {                                                                                   // 273
      throw new Meteor.Error(403, noReplaceError);                                                                // 274
    } else if (!_.has(ALLOWED_UPDATE_OPERATIONS, op)) {                                                           //
      throw new Meteor.Error(403, "Access denied. Operator " + op + " not allowed in a restricted collection.");  // 276
    } else {                                                                                                      //
      _.each(_.keys(params), function (field) {                                                                   // 279
        // treat dotted fields as if they are replacing their                                                     //
        // top-level part                                                                                         //
        if (field.indexOf('.') !== -1) field = field.substring(0, field.indexOf('.'));                            // 282
                                                                                                                  //
        // record the field we are trying to change                                                               //
        if (!_.contains(fields, field)) fields.push(field);                                                       // 279
      });                                                                                                         //
    }                                                                                                             //
  });                                                                                                             //
                                                                                                                  //
  var findOptions = { transform: null };                                                                          // 292
  if (!self._validators.fetchAllFields) {                                                                         // 293
    findOptions.fields = {};                                                                                      // 294
    _.each(self._validators.fetch, function (fieldName) {                                                         // 295
      findOptions.fields[fieldName] = 1;                                                                          // 296
    });                                                                                                           //
  }                                                                                                               //
                                                                                                                  //
  var doc = self._collection.findOne(selector, findOptions);                                                      // 300
  if (!doc) // none satisfied!                                                                                    // 301
    return 0;                                                                                                     // 302
                                                                                                                  //
  // call user validators.                                                                                        //
  // Any deny returns true means denied.                                                                          //
  if (_.any(self._validators.update.deny, function (validator) {                                                  // 247
    var factoriedDoc = transformDoc(validator, doc);                                                              // 307
    return validator(userId, factoriedDoc, fields, mutator);                                                      // 308
  })) {                                                                                                           //
    throw new Meteor.Error(403, "Access denied");                                                                 // 313
  }                                                                                                               //
  // Any allow returns true means proceed. Throw error if they all fail.                                          //
  if (_.all(self._validators.update.allow, function (validator) {                                                 // 247
    var factoriedDoc = transformDoc(validator, doc);                                                              // 317
    return !validator(userId, factoriedDoc, fields, mutator);                                                     // 318
  })) {                                                                                                           //
    throw new Meteor.Error(403, "Access denied");                                                                 // 323
  }                                                                                                               //
                                                                                                                  //
  options._forbidReplace = true;                                                                                  // 326
                                                                                                                  //
  // Back when we supported arbitrary client-provided selectors, we actually                                      //
  // rewrote the selector to include an _id clause before passing to Mongo to                                     //
  // avoid races, but since selector is guaranteed to already just be an ID, we                                   //
  // don't have to any more.                                                                                      //
                                                                                                                  //
  return self._collection.update.call(self._collection, selector, mutator, options);                              // 247
};                                                                                                                //
                                                                                                                  //
// Only allow these operations in validated updates. Specifically                                                 //
// whitelist operations, rather than blacklist, so new complex                                                    //
// operations that are added aren't automatically allowed. A complex                                              //
// operation is one that does more than just modify its target                                                    //
// field. For now this contains all update operations except '$rename'.                                           //
// http://docs.mongodb.org/manual/reference/operators/#update                                                     //
var ALLOWED_UPDATE_OPERATIONS = {                                                                                 // 343
  $inc: 1, $set: 1, $unset: 1, $addToSet: 1, $pop: 1, $pullAll: 1, $pull: 1,                                      // 344
  $pushAll: 1, $push: 1, $bit: 1                                                                                  // 345
};                                                                                                                //
                                                                                                                  //
// Simulate a mongo `remove` operation while validating access control                                            //
// rules. See #ValidatedChange                                                                                    //
CollectionPrototype._validatedRemove = function (userId, selector) {                                              // 350
  var self = this;                                                                                                // 351
                                                                                                                  //
  var findOptions = { transform: null };                                                                          // 353
  if (!self._validators.fetchAllFields) {                                                                         // 354
    findOptions.fields = {};                                                                                      // 355
    _.each(self._validators.fetch, function (fieldName) {                                                         // 356
      findOptions.fields[fieldName] = 1;                                                                          // 357
    });                                                                                                           //
  }                                                                                                               //
                                                                                                                  //
  var doc = self._collection.findOne(selector, findOptions);                                                      // 361
  if (!doc) return 0;                                                                                             // 362
                                                                                                                  //
  // call user validators.                                                                                        //
  // Any deny returns true means denied.                                                                          //
  if (_.any(self._validators.remove.deny, function (validator) {                                                  // 350
    return validator(userId, transformDoc(validator, doc));                                                       // 368
  })) {                                                                                                           //
    throw new Meteor.Error(403, "Access denied");                                                                 // 370
  }                                                                                                               //
  // Any allow returns true means proceed. Throw error if they all fail.                                          //
  if (_.all(self._validators.remove.allow, function (validator) {                                                 // 350
    return !validator(userId, transformDoc(validator, doc));                                                      // 374
  })) {                                                                                                           //
    throw new Meteor.Error(403, "Access denied");                                                                 // 376
  }                                                                                                               //
                                                                                                                  //
  // Back when we supported arbitrary client-provided selectors, we actually                                      //
  // rewrote the selector to {_id: {$in: [ids that we found]}} before passing to                                  //
  // Mongo to avoid races, but since selector is guaranteed to already just be                                    //
  // an ID, we don't have to any more.                                                                            //
                                                                                                                  //
  return self._collection.remove.call(self._collection, selector);                                                // 350
};                                                                                                                //
                                                                                                                  //
CollectionPrototype._callMutatorMethod = function () {                                                            // 387
  function _callMutatorMethod(name, args, callback) {                                                             // 387
    if (Meteor.isClient && !callback && !alreadyInSimulation()) {                                                 // 388
      // Client can't block, so it can't report errors by exception,                                              //
      // only by callback. If they forget the callback, give them a                                               //
      // default one that logs the error, so they aren't totally                                                  //
      // baffled if their writes don't work because their database is                                             //
      // down.                                                                                                    //
      // Don't give a default callback in simulation, because inside stubs we                                     //
      // want to return the results from the local collection immediately and                                     //
      // not force a callback.                                                                                    //
      callback = function () {                                                                                    // 397
        function callback(err) {                                                                                  // 397
          if (err) Meteor._debug(name + " failed: " + (err.reason || err.stack));                                 // 398
        }                                                                                                         //
                                                                                                                  //
        return callback;                                                                                          //
      }();                                                                                                        //
    }                                                                                                             //
                                                                                                                  //
    // For two out of three mutator methods, the first argument is a selector                                     //
    var firstArgIsSelector = name === "update" || name === "remove";                                              // 387
    if (firstArgIsSelector && !alreadyInSimulation()) {                                                           // 405
      // If we're about to actually send an RPC, we should throw an error if                                      //
      // this is a non-ID selector, because the mutation methods only allow                                       //
      // single-ID selectors. (If we don't throw here, we'll see flicker.)                                        //
      throwIfSelectorIsNotId(args[0], name);                                                                      // 409
    }                                                                                                             //
                                                                                                                  //
    var mutatorMethodName = this._prefix + name;                                                                  // 412
    return this._connection.apply(mutatorMethodName, args, { returnStubValue: true }, callback);                  // 413
  }                                                                                                               //
                                                                                                                  //
  return _callMutatorMethod;                                                                                      //
}();                                                                                                              //
                                                                                                                  //
function transformDoc(validator, doc) {                                                                           // 417
  if (validator.transform) return validator.transform(doc);                                                       // 418
  return doc;                                                                                                     // 420
}                                                                                                                 //
                                                                                                                  //
function docToValidate(validator, doc, generatedId) {                                                             // 423
  var ret = doc;                                                                                                  // 424
  if (validator.transform) {                                                                                      // 425
    ret = EJSON.clone(doc);                                                                                       // 426
    // If you set a server-side transform on your collection, then you don't get                                  //
    // to tell the difference between "client specified the ID" and "server                                       //
    // generated the ID", because transforms expect to get _id.  If you want to                                   //
    // do that check, you can do it with a specific                                                               //
    // `C.allow({insert: f, transform: null})` validator.                                                         //
    if (generatedId !== null) {                                                                                   // 425
      ret._id = generatedId;                                                                                      // 433
    }                                                                                                             //
    ret = validator.transform(ret);                                                                               // 435
  }                                                                                                               //
  return ret;                                                                                                     // 437
}                                                                                                                 //
                                                                                                                  //
function addValidator(collection, allowOrDeny, options) {                                                         // 440
  // validate keys                                                                                                //
  var VALID_KEYS = ['insert', 'update', 'remove', 'fetch', 'transform'];                                          // 442
  _.each(_.keys(options), function (key) {                                                                        // 443
    if (!_.contains(VALID_KEYS, key)) throw new Error(allowOrDeny + ": Invalid key: " + key);                     // 444
  });                                                                                                             //
                                                                                                                  //
  collection._restricted = true;                                                                                  // 448
                                                                                                                  //
  _.each(['insert', 'update', 'remove'], function (name) {                                                        // 450
    if (options.hasOwnProperty(name)) {                                                                           // 451
      if (!(options[name] instanceof Function)) {                                                                 // 452
        throw new Error(allowOrDeny + ": Value for `" + name + "` must be a function");                           // 453
      }                                                                                                           //
                                                                                                                  //
      // If the transform is specified at all (including as 'null') in this                                       //
      // call, then take that; otherwise, take the transform from the                                             //
      // collection.                                                                                              //
      if (options.transform === undefined) {                                                                      // 451
        options[name].transform = collection._transform; // already wrapped                                       // 460
      } else {                                                                                                    // 459
          options[name].transform = LocalCollection.wrapTransform(options.transform);                             // 462
        }                                                                                                         //
                                                                                                                  //
      collection._validators[name][allowOrDeny].push(options[name]);                                              // 466
    }                                                                                                             //
  });                                                                                                             //
                                                                                                                  //
  // Only update the fetch fields if we're passed things that affect                                              //
  // fetching. This way allow({}) and allow({insert: f}) don't result in                                          //
  // setting fetchAllFields                                                                                       //
  if (options.update || options.remove || options.fetch) {                                                        // 440
    if (options.fetch && !(options.fetch instanceof Array)) {                                                     // 474
      throw new Error(allowOrDeny + ": Value for `fetch` must be an array");                                      // 475
    }                                                                                                             //
    collection._updateFetch(options.fetch);                                                                       // 477
  }                                                                                                               //
}                                                                                                                 //
                                                                                                                  //
function throwIfSelectorIsNotId(selector, methodName) {                                                           // 481
  if (!LocalCollection._selectorIsIdPerhapsAsObject(selector)) {                                                  // 482
    throw new Meteor.Error(403, "Not permitted. Untrusted code may only " + methodName + " documents by ID.");    // 483
  }                                                                                                               //
};                                                                                                                //
                                                                                                                  //
// Determine if we are in a DDP method simulation                                                                 //
function alreadyInSimulation() {                                                                                  // 490
  var enclosing = DDP._CurrentInvocation.get();                                                                   // 491
  return enclosing && enclosing.isSimulation;                                                                     // 492
}                                                                                                                 //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}},{"extensions":[".js",".json"]});
require("./node_modules/meteor/allow-deny/allow-deny.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['allow-deny'] = {}, {
  AllowDeny: AllowDeny
});

})();
