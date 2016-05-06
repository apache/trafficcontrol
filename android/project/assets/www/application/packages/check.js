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
var _ = Package.underscore._;
var EJSON = Package.ejson.EJSON;

/* Package-scope variables */
var check, Match;

var require = meteorInstall({"node_modules":{"meteor":{"check":{"match.js":["./isPlainObject.js",function(require,exports){

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                       //
// packages/check/match.js                                                                               //
//                                                                                                       //
///////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                         //
// XXX docs                                                                                              // 1
                                                                                                         // 2
// Things we explicitly do NOT support:                                                                  // 3
//    - heterogenous arrays                                                                              // 4
                                                                                                         // 5
var currentArgumentChecker = new Meteor.EnvironmentVariable;                                             // 6
var isPlainObject = require("./isPlainObject.js").isPlainObject;                                         // 7
                                                                                                         // 8
/**                                                                                                      // 9
 * @summary Check that a value matches a [pattern](#matchpatterns).                                      // 10
 * If the value does not match the pattern, throw a `Match.Error`.                                       // 11
 *                                                                                                       // 12
 * Particularly useful to assert that arguments to a function have the right                             // 13
 * types and structure.                                                                                  // 14
 * @locus Anywhere                                                                                       // 15
 * @param {Any} value The value to check                                                                 // 16
 * @param {MatchPattern} pattern The pattern to match                                                    // 17
 * `value` against                                                                                       // 18
 */                                                                                                      // 19
var check = exports.check = function (value, pattern) {                                                  // 20
  // Record that check got called, if somebody cared.                                                    // 21
  //                                                                                                     // 22
  // We use getOrNullIfOutsideFiber so that it's OK to call check()                                      // 23
  // from non-Fiber server contexts; the downside is that if you forget to                               // 24
  // bindEnvironment on some random callback in your method/publisher,                                   // 25
  // it might not find the argumentChecker and you'll get an error about                                 // 26
  // not checking an argument that it looks like you're checking (instead                                // 27
  // of just getting a "Node code must run in a Fiber" error).                                           // 28
  var argChecker = currentArgumentChecker.getOrNullIfOutsideFiber();                                     // 29
  if (argChecker)                                                                                        // 30
    argChecker.checking(value);                                                                          // 31
  var result = testSubtree(value, pattern);                                                              // 32
  if (result) {                                                                                          // 33
    var err = new Match.Error(result.message);                                                           // 34
    if (result.path) {                                                                                   // 35
      err.message += " in field " + result.path;                                                         // 36
      err.path = result.path;                                                                            // 37
    }                                                                                                    // 38
    throw err;                                                                                           // 39
  }                                                                                                      // 40
};                                                                                                       // 41
                                                                                                         // 42
/**                                                                                                      // 43
 * @namespace Match                                                                                      // 44
 * @summary The namespace for all Match types and methods.                                               // 45
 */                                                                                                      // 46
var Match = exports.Match = {                                                                            // 47
  Optional: function (pattern) {                                                                         // 48
    return new Optional(pattern);                                                                        // 49
  },                                                                                                     // 50
  Maybe: function (pattern) {                                                                            // 51
    return new Maybe(pattern);                                                                           // 52
  },                                                                                                     // 53
  OneOf: function (/*arguments*/) {                                                                      // 54
    return new OneOf(_.toArray(arguments));                                                              // 55
  },                                                                                                     // 56
  Any: ['__any__'],                                                                                      // 57
  Where: function (condition) {                                                                          // 58
    return new Where(condition);                                                                         // 59
  },                                                                                                     // 60
  ObjectIncluding: function (pattern) {                                                                  // 61
    return new ObjectIncluding(pattern);                                                                 // 62
  },                                                                                                     // 63
  ObjectWithValues: function (pattern) {                                                                 // 64
    return new ObjectWithValues(pattern);                                                                // 65
  },                                                                                                     // 66
  // Matches only signed 32-bit integers                                                                 // 67
  Integer: ['__integer__'],                                                                              // 68
                                                                                                         // 69
  // XXX matchers should know how to describe themselves for errors                                      // 70
  Error: Meteor.makeErrorType("Match.Error", function (msg) {                                            // 71
    this.message = "Match error: " + msg;                                                                // 72
    // The path of the value that failed to match. Initially empty, this gets                            // 73
    // populated by catching and rethrowing the exception as it goes back up the                         // 74
    // stack.                                                                                            // 75
    // E.g.: "vals[3].entity.created"                                                                    // 76
    this.path = "";                                                                                      // 77
    // If this gets sent over DDP, don't give full internal details but at least                         // 78
    // provide something better than 500 Internal server error.                                          // 79
    this.sanitizedError = new Meteor.Error(400, "Match failed");                                         // 80
  }),                                                                                                    // 81
                                                                                                         // 82
  // Tests to see if value matches pattern. Unlike check, it merely returns true                         // 83
  // or false (unless an error other than Match.Error was thrown). It does not                           // 84
  // interact with _failIfArgumentsAreNotAllChecked.                                                     // 85
  // XXX maybe also implement a Match.match which returns more information about                         // 86
  //     failures but without using exception handling or doing what check()                             // 87
  //     does with _failIfArgumentsAreNotAllChecked and Meteor.Error conversion                          // 88
                                                                                                         // 89
  /**                                                                                                    // 90
   * @summary Returns true if the value matches the pattern.                                             // 91
   * @locus Anywhere                                                                                     // 92
   * @param {Any} value The value to check                                                               // 93
   * @param {MatchPattern} pattern The pattern to match `value` against                                  // 94
   */                                                                                                    // 95
  test: function (value, pattern) {                                                                      // 96
    return !testSubtree(value, pattern);                                                                 // 97
  },                                                                                                     // 98
                                                                                                         // 99
  // Runs `f.apply(context, args)`. If check() is not called on every element of                         // 100
  // `args` (either directly or in the first level of an array), throws an error                         // 101
  // (using `description` in the message).                                                               // 102
  //                                                                                                     // 103
  _failIfArgumentsAreNotAllChecked: function (f, context, args, description) {                           // 104
    var argChecker = new ArgumentChecker(args, description);                                             // 105
    var result = currentArgumentChecker.withValue(argChecker, function () {                              // 106
      return f.apply(context, args);                                                                     // 107
    });                                                                                                  // 108
    // If f didn't itself throw, make sure it checked all of its arguments.                              // 109
    argChecker.throwUnlessAllArgumentsHaveBeenChecked();                                                 // 110
    return result;                                                                                       // 111
  }                                                                                                      // 112
};                                                                                                       // 113
                                                                                                         // 114
var Optional = function (pattern) {                                                                      // 115
  this.pattern = pattern;                                                                                // 116
};                                                                                                       // 117
                                                                                                         // 118
var Maybe = function (pattern) {                                                                         // 119
  this.pattern = pattern;                                                                                // 120
};                                                                                                       // 121
                                                                                                         // 122
var OneOf = function (choices) {                                                                         // 123
  if (_.isEmpty(choices))                                                                                // 124
    throw new Error("Must provide at least one choice to Match.OneOf");                                  // 125
  this.choices = choices;                                                                                // 126
};                                                                                                       // 127
                                                                                                         // 128
var Where = function (condition) {                                                                       // 129
  this.condition = condition;                                                                            // 130
};                                                                                                       // 131
                                                                                                         // 132
var ObjectIncluding = function (pattern) {                                                               // 133
  this.pattern = pattern;                                                                                // 134
};                                                                                                       // 135
                                                                                                         // 136
var ObjectWithValues = function (pattern) {                                                              // 137
  this.pattern = pattern;                                                                                // 138
};                                                                                                       // 139
                                                                                                         // 140
var typeofChecks = [                                                                                     // 141
  [String, "string"],                                                                                    // 142
  [Number, "number"],                                                                                    // 143
  [Boolean, "boolean"],                                                                                  // 144
  // While we don't allow undefined in EJSON, this is good for optional                                  // 145
  // arguments with OneOf.                                                                               // 146
  [undefined, "undefined"]                                                                               // 147
];                                                                                                       // 148
                                                                                                         // 149
// Return `false` if it matches. Otherwise, return an object with a `message` and a `path` field.        // 150
var testSubtree = function (value, pattern) {                                                            // 151
  // Match anything!                                                                                     // 152
  if (pattern === Match.Any)                                                                             // 153
    return false;                                                                                        // 154
                                                                                                         // 155
  // Basic atomic types.                                                                                 // 156
  // Do not match boxed objects (e.g. String, Boolean)                                                   // 157
  for (var i = 0; i < typeofChecks.length; ++i) {                                                        // 158
    if (pattern === typeofChecks[i][0]) {                                                                // 159
      if (typeof value === typeofChecks[i][1])                                                           // 160
        return false;                                                                                    // 161
      return {                                                                                           // 162
        message: "Expected " + typeofChecks[i][1] + ", got " + (value === null ? "null" : typeof value),
        path: ""                                                                                         // 164
      };                                                                                                 // 165
    }                                                                                                    // 166
  }                                                                                                      // 167
  if (pattern === null) {                                                                                // 168
    if (value === null)                                                                                  // 169
      return false;                                                                                      // 170
    return {                                                                                             // 171
      message: "Expected null, got " + EJSON.stringify(value),                                           // 172
      path: ""                                                                                           // 173
    };                                                                                                   // 174
  }                                                                                                      // 175
                                                                                                         // 176
  // Strings, numbers, and booleans match literally. Goes well with Match.OneOf.                         // 177
  if (typeof pattern === "string" || typeof pattern === "number" || typeof pattern === "boolean") {      // 178
    if (value === pattern)                                                                               // 179
      return false;                                                                                      // 180
    return {                                                                                             // 181
      message: "Expected " + pattern + ", got " + EJSON.stringify(value),                                // 182
      path: ""                                                                                           // 183
    };                                                                                                   // 184
  }                                                                                                      // 185
                                                                                                         // 186
  // Match.Integer is special type encoded with array                                                    // 187
  if (pattern === Match.Integer) {                                                                       // 188
    // There is no consistent and reliable way to check if variable is a 64-bit                          // 189
    // integer. One of the popular solutions is to get reminder of division by 1                         // 190
    // but this method fails on really large floats with big precision.                                  // 191
    // E.g.: 1.348192308491824e+23 % 1 === 0 in V8                                                       // 192
    // Bitwise operators work consistantly but always cast variable to 32-bit                            // 193
    // signed integer according to JavaScript specs.                                                     // 194
    if (typeof value === "number" && (value | 0) === value)                                              // 195
      return false;                                                                                      // 196
    return {                                                                                             // 197
      message: "Expected Integer, got " + (value instanceof Object ? EJSON.stringify(value) : value),    // 198
      path: ""                                                                                           // 199
    };                                                                                                   // 200
  }                                                                                                      // 201
                                                                                                         // 202
  // "Object" is shorthand for Match.ObjectIncluding({});                                                // 203
  if (pattern === Object)                                                                                // 204
    pattern = Match.ObjectIncluding({});                                                                 // 205
                                                                                                         // 206
  // Array (checked AFTER Any, which is implemented as an Array).                                        // 207
  if (pattern instanceof Array) {                                                                        // 208
    if (pattern.length !== 1) {                                                                          // 209
      return {                                                                                           // 210
        message: "Bad pattern: arrays must have one type element" + EJSON.stringify(pattern),            // 211
        path: ""                                                                                         // 212
      };                                                                                                 // 213
    }                                                                                                    // 214
    if (!_.isArray(value) && !_.isArguments(value)) {                                                    // 215
      return {                                                                                           // 216
        message: "Expected array, got " + EJSON.stringify(value),                                        // 217
        path: ""                                                                                         // 218
      };                                                                                                 // 219
    }                                                                                                    // 220
                                                                                                         // 221
    for (var i = 0, length = value.length; i < length; i++) {                                            // 222
      var result = testSubtree(value[i], pattern[0]);                                                    // 223
      if (result) {                                                                                      // 224
        result.path = _prependPath(i, result.path);                                                      // 225
        return result;                                                                                   // 226
      }                                                                                                  // 227
    }                                                                                                    // 228
    return false;                                                                                        // 229
  }                                                                                                      // 230
                                                                                                         // 231
  // Arbitrary validation checks. The condition can return false or throw a                              // 232
  // Match.Error (ie, it can internally use check()) to fail.                                            // 233
  if (pattern instanceof Where) {                                                                        // 234
    var result;                                                                                          // 235
    try {                                                                                                // 236
      result = pattern.condition(value);                                                                 // 237
    } catch (err) {                                                                                      // 238
      if (!(err instanceof Match.Error))                                                                 // 239
        throw err;                                                                                       // 240
      return {                                                                                           // 241
        message: err.message,                                                                            // 242
        path: err.path                                                                                   // 243
      };                                                                                                 // 244
    }                                                                                                    // 245
    if (result)                                                                                          // 246
      return false;                                                                                      // 247
    // XXX this error is terrible                                                                        // 248
    return {                                                                                             // 249
      message: "Failed Match.Where validation",                                                          // 250
      path: ""                                                                                           // 251
    };                                                                                                   // 252
  }                                                                                                      // 253
                                                                                                         // 254
                                                                                                         // 255
  if (pattern instanceof Maybe) {                                                                        // 256
    pattern = Match.OneOf(undefined, null, pattern.pattern);                                             // 257
  }                                                                                                      // 258
  else if (pattern instanceof Optional) {                                                                // 259
    pattern = Match.OneOf(undefined, pattern.pattern);                                                   // 260
  }                                                                                                      // 261
                                                                                                         // 262
  if (pattern instanceof OneOf) {                                                                        // 263
    for (var i = 0; i < pattern.choices.length; ++i) {                                                   // 264
      var result = testSubtree(value, pattern.choices[i]);                                               // 265
      if (!result) {                                                                                     // 266
        // No error? Yay, return.                                                                        // 267
        return false;                                                                                    // 268
      }                                                                                                  // 269
      // Match errors just mean try another choice.                                                      // 270
    }                                                                                                    // 271
    // XXX this error is terrible                                                                        // 272
    return {                                                                                             // 273
      message: "Failed Match.OneOf, Match.Maybe or Match.Optional validation",                           // 274
      path: ""                                                                                           // 275
    };                                                                                                   // 276
  }                                                                                                      // 277
                                                                                                         // 278
  // A function that isn't something we special-case is assumed to be a                                  // 279
  // constructor.                                                                                        // 280
  if (pattern instanceof Function) {                                                                     // 281
    if (value instanceof pattern)                                                                        // 282
      return false;                                                                                      // 283
    return {                                                                                             // 284
      message: "Expected " + (pattern.name ||"particular constructor"),                                  // 285
      path: ""                                                                                           // 286
    };                                                                                                   // 287
  }                                                                                                      // 288
                                                                                                         // 289
  var unknownKeysAllowed = false;                                                                        // 290
  var unknownKeyPattern;                                                                                 // 291
  if (pattern instanceof ObjectIncluding) {                                                              // 292
    unknownKeysAllowed = true;                                                                           // 293
    pattern = pattern.pattern;                                                                           // 294
  }                                                                                                      // 295
  if (pattern instanceof ObjectWithValues) {                                                             // 296
    unknownKeysAllowed = true;                                                                           // 297
    unknownKeyPattern = [pattern.pattern];                                                               // 298
    pattern = {};  // no required keys                                                                   // 299
  }                                                                                                      // 300
                                                                                                         // 301
  if (typeof pattern !== "object") {                                                                     // 302
    return {                                                                                             // 303
      message: "Bad pattern: unknown pattern type",                                                      // 304
      path: ""                                                                                           // 305
    };                                                                                                   // 306
  }                                                                                                      // 307
                                                                                                         // 308
  // An object, with required and optional keys. Note that this does NOT do                              // 309
  // structural matches against objects of special types that happen to match                            // 310
  // the pattern: this really needs to be a plain old {Object}!                                          // 311
  if (typeof value !== 'object') {                                                                       // 312
    return {                                                                                             // 313
      message: "Expected object, got " + typeof value,                                                   // 314
      path: ""                                                                                           // 315
    };                                                                                                   // 316
  }                                                                                                      // 317
  if (value === null) {                                                                                  // 318
    return {                                                                                             // 319
      message: "Expected object, got null",                                                              // 320
      path: ""                                                                                           // 321
    };                                                                                                   // 322
  }                                                                                                      // 323
  if (! isPlainObject(value)) {                                                                          // 324
    return {                                                                                             // 325
      message: "Expected plain object",                                                                  // 326
      path: ""                                                                                           // 327
    };                                                                                                   // 328
  }                                                                                                      // 329
                                                                                                         // 330
  var requiredPatterns = {};                                                                             // 331
  var optionalPatterns = {};                                                                             // 332
  _.each(pattern, function (subPattern, key) {                                                           // 333
    if (subPattern instanceof Optional || subPattern instanceof Maybe)                                   // 334
      optionalPatterns[key] = subPattern.pattern;                                                        // 335
    else                                                                                                 // 336
      requiredPatterns[key] = subPattern;                                                                // 337
  });                                                                                                    // 338
                                                                                                         // 339
  //XXX: replace with underscore's _.allKeys if Meteor updates underscore to 1.8+ (or lodash)            // 340
  var allKeys = function(obj){                                                                           // 341
    var keys = [];                                                                                       // 342
    if (_.isObject(obj)){                                                                                // 343
      for (var key in obj) keys.push(key);                                                               // 344
    }                                                                                                    // 345
    return keys;                                                                                         // 346
  }                                                                                                      // 347
                                                                                                         // 348
  for (var keys = allKeys(value), i = 0, length = keys.length; i < length; i++) {                        // 349
    var key = keys[i];                                                                                   // 350
    var subValue = value[key];                                                                           // 351
    if (_.has(requiredPatterns, key)) {                                                                  // 352
      var result = testSubtree(subValue, requiredPatterns[key]);                                         // 353
      if (result) {                                                                                      // 354
        result.path = _prependPath(key, result.path);                                                    // 355
        return result;                                                                                   // 356
      }                                                                                                  // 357
      delete requiredPatterns[key];                                                                      // 358
    } else if (_.has(optionalPatterns, key)) {                                                           // 359
      var result = testSubtree(subValue, optionalPatterns[key]);                                         // 360
      if (result) {                                                                                      // 361
        result.path = _prependPath(key, result.path);                                                    // 362
        return result;                                                                                   // 363
      }                                                                                                  // 364
    } else {                                                                                             // 365
      if (!unknownKeysAllowed) {                                                                         // 366
        return {                                                                                         // 367
          message: "Unknown key",                                                                        // 368
          path: key                                                                                      // 369
        };                                                                                               // 370
      }                                                                                                  // 371
      if (unknownKeyPattern) {                                                                           // 372
        var result = testSubtree(subValue, unknownKeyPattern[0]);                                        // 373
        if (result) {                                                                                    // 374
          result.path = _prependPath(key, result.path);                                                  // 375
          return result;                                                                                 // 376
        }                                                                                                // 377
      }                                                                                                  // 378
    }                                                                                                    // 379
  }                                                                                                      // 380
                                                                                                         // 381
  var keys = _.keys(requiredPatterns);                                                                   // 382
  if (keys.length) {                                                                                     // 383
    return {                                                                                             // 384
      message: "Missing key '" + keys[0] + "'",                                                          // 385
      path: ""                                                                                           // 386
    };                                                                                                   // 387
  }                                                                                                      // 388
};                                                                                                       // 389
                                                                                                         // 390
var ArgumentChecker = function (args, description) {                                                     // 391
  var self = this;                                                                                       // 392
  // Make a SHALLOW copy of the arguments. (We'll be doing identity checks                               // 393
  // against its contents.)                                                                              // 394
  self.args = _.clone(args);                                                                             // 395
  // Since the common case will be to check arguments in order, and we splice                            // 396
  // out arguments when we check them, make it so we splice out from the end                             // 397
  // rather than the beginning.                                                                          // 398
  self.args.reverse();                                                                                   // 399
  self.description = description;                                                                        // 400
};                                                                                                       // 401
                                                                                                         // 402
_.extend(ArgumentChecker.prototype, {                                                                    // 403
  checking: function (value) {                                                                           // 404
    var self = this;                                                                                     // 405
    if (self._checkingOneValue(value))                                                                   // 406
      return;                                                                                            // 407
    // Allow check(arguments, [String]) or check(arguments.slice(1), [String])                           // 408
    // or check([foo, bar], [String]) to count... but only if value wasn't                               // 409
    // itself an argument.                                                                               // 410
    if (_.isArray(value) || _.isArguments(value)) {                                                      // 411
      _.each(value, _.bind(self._checkingOneValue, self));                                               // 412
    }                                                                                                    // 413
  },                                                                                                     // 414
  _checkingOneValue: function (value) {                                                                  // 415
    var self = this;                                                                                     // 416
    for (var i = 0; i < self.args.length; ++i) {                                                         // 417
      // Is this value one of the arguments? (This can have a false positive if                          // 418
      // the argument is an interned primitive, but it's still a good enough                             // 419
      // check.)                                                                                         // 420
      // (NaN is not === to itself, so we have to check specially.)                                      // 421
      if (value === self.args[i] || (_.isNaN(value) && _.isNaN(self.args[i]))) {                         // 422
        self.args.splice(i, 1);                                                                          // 423
        return true;                                                                                     // 424
      }                                                                                                  // 425
    }                                                                                                    // 426
    return false;                                                                                        // 427
  },                                                                                                     // 428
  throwUnlessAllArgumentsHaveBeenChecked: function () {                                                  // 429
    var self = this;                                                                                     // 430
    if (!_.isEmpty(self.args))                                                                           // 431
      throw new Error("Did not check() all arguments during " +                                          // 432
                      self.description);                                                                 // 433
  }                                                                                                      // 434
});                                                                                                      // 435
                                                                                                         // 436
var _jsKeywords = ["do", "if", "in", "for", "let", "new", "try", "var", "case",                          // 437
  "else", "enum", "eval", "false", "null", "this", "true", "void", "with",                               // 438
  "break", "catch", "class", "const", "super", "throw", "while", "yield",                                // 439
  "delete", "export", "import", "public", "return", "static", "switch",                                  // 440
  "typeof", "default", "extends", "finally", "package", "private", "continue",                           // 441
  "debugger", "function", "arguments", "interface", "protected", "implements",                           // 442
  "instanceof"];                                                                                         // 443
                                                                                                         // 444
// Assumes the base of path is already escaped properly                                                  // 445
// returns key + base                                                                                    // 446
var _prependPath = function (key, base) {                                                                // 447
  if ((typeof key) === "number" || key.match(/^[0-9]+$/))                                                // 448
    key = "[" + key + "]";                                                                               // 449
  else if (!key.match(/^[a-z_$][0-9a-z_$]*$/i) || _.contains(_jsKeywords, key))                          // 450
    key = JSON.stringify([key]);                                                                         // 451
                                                                                                         // 452
  if (base && base[0] !== "[")                                                                           // 453
    return key + '.' + base;                                                                             // 454
  return key + base;                                                                                     // 455
};                                                                                                       // 456
                                                                                                         // 457
                                                                                                         // 458
///////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"isPlainObject.js":function(require,exports){

///////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                       //
// packages/check/isPlainObject.js                                                                       //
//                                                                                                       //
///////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                         //
// Copy of jQuery.isPlainObject for the server side from jQuery v1.11.2.                                 // 1
                                                                                                         // 2
var class2type = {};                                                                                     // 3
                                                                                                         // 4
var toString = class2type.toString;                                                                      // 5
                                                                                                         // 6
var hasOwn = class2type.hasOwnProperty;                                                                  // 7
                                                                                                         // 8
var support = {};                                                                                        // 9
                                                                                                         // 10
// Populate the class2type map                                                                           // 11
_.each("Boolean Number String Function Array Date RegExp Object Error".split(" "), function(name, i) {   // 12
  class2type[ "[object " + name + "]" ] = name.toLowerCase();                                            // 13
});                                                                                                      // 14
                                                                                                         // 15
function type( obj ) {                                                                                   // 16
  if ( obj == null ) {                                                                                   // 17
    return obj + "";                                                                                     // 18
  }                                                                                                      // 19
  return typeof obj === "object" || typeof obj === "function" ?                                          // 20
    class2type[ toString.call(obj) ] || "object" :                                                       // 21
    typeof obj;                                                                                          // 22
}                                                                                                        // 23
                                                                                                         // 24
function isWindow( obj ) {                                                                               // 25
  /* jshint eqeqeq: false */                                                                             // 26
  return obj != null && obj == obj.window;                                                               // 27
}                                                                                                        // 28
                                                                                                         // 29
exports.isPlainObject = function( obj ) {                                                                // 30
  var key;                                                                                               // 31
                                                                                                         // 32
  // Must be an Object.                                                                                  // 33
  // Because of IE, we also have to check the presence of the constructor property.                      // 34
  // Make sure that DOM nodes and window objects don't pass through, as well                             // 35
  if ( !obj || type(obj) !== "object" || obj.nodeType || isWindow( obj ) ) {                             // 36
    return false;                                                                                        // 37
  }                                                                                                      // 38
                                                                                                         // 39
  try {                                                                                                  // 40
    // Not own constructor property must be Object                                                       // 41
    if ( obj.constructor &&                                                                              // 42
         !hasOwn.call(obj, "constructor") &&                                                             // 43
         !hasOwn.call(obj.constructor.prototype, "isPrototypeOf") ) {                                    // 44
      return false;                                                                                      // 45
    }                                                                                                    // 46
  } catch ( e ) {                                                                                        // 47
    // IE8,9 Will throw exceptions on certain host objects #9897                                         // 48
    return false;                                                                                        // 49
  }                                                                                                      // 50
                                                                                                         // 51
  // Support: IE<9                                                                                       // 52
  // Handle iteration over inherited properties before own properties.                                   // 53
  if ( support.ownLast ) {                                                                               // 54
    for ( key in obj ) {                                                                                 // 55
      return hasOwn.call( obj, key );                                                                    // 56
    }                                                                                                    // 57
  }                                                                                                      // 58
                                                                                                         // 59
  // Own properties are enumerated firstly, so to speed up,                                              // 60
  // if last one is own, then all properties are own.                                                    // 61
  for ( key in obj ) {}                                                                                  // 62
                                                                                                         // 63
  return key === undefined || hasOwn.call( obj, key );                                                   // 64
};                                                                                                       // 65
                                                                                                         // 66
///////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}},{"extensions":[".js",".json"]});
var exports = require("./node_modules/meteor/check/match.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.check = exports, {
  check: check,
  Match: Match
});

})();
