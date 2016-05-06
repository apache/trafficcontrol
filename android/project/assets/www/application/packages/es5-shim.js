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
var Date, parseInt, parseFloat, originalStringReplace;

var require = meteorInstall({"node_modules":{"meteor":{"es5-shim":{"client.js":["./import_globals.js","es5-shim/es5-shim.js","es5-shim/es5-sham.js","./console.js","./export_globals.js",function(require){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/es5-shim/client.js                                                                                         //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
require("./import_globals.js");                                                                                        // 1
require("es5-shim/es5-shim.js");                                                                                       // 2
require("es5-shim/es5-sham.js");                                                                                       // 3
require("./console.js");                                                                                               // 4
require("./export_globals.js");                                                                                        // 5
                                                                                                                       // 6
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"console.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/es5-shim/console.js                                                                                        //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
var hasOwn = Object.prototype.hasOwnProperty;                                                                          // 1
                                                                                                                       // 2
function wrap(method) {                                                                                                // 3
  var original = console[method];                                                                                      // 4
  if (original && typeof original === "object") {                                                                      // 5
    // Turn callable console method objects into actual functions.                                                     // 6
    console[method] = function () {                                                                                    // 7
      return Function.prototype.apply.call(                                                                            // 8
        original, console, arguments                                                                                   // 9
      );                                                                                                               // 10
    };                                                                                                                 // 11
  }                                                                                                                    // 12
}                                                                                                                      // 13
                                                                                                                       // 14
if (typeof console === "object" &&                                                                                     // 15
    // In older Internet Explorers, methods like console.log are actually                                              // 16
    // callable objects rather than functions.                                                                         // 17
    typeof console.log === "object") {                                                                                 // 18
  for (var method in console) {                                                                                        // 19
    // In most browsers, this hasOwn check will fail for all console                                                   // 20
    // methods anyway, but fortunately in IE8 the method objects we care                                               // 21
    // about are own properties.                                                                                       // 22
    if (hasOwn.call(console, method)) {                                                                                // 23
      wrap(method);                                                                                                    // 24
    }                                                                                                                  // 25
  }                                                                                                                    // 26
}                                                                                                                      // 27
                                                                                                                       // 28
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"export_globals.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/es5-shim/export_globals.js                                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
if (global.Date !== Date) {                                                                                            // 1
  global.Date = Date;                                                                                                  // 2
}                                                                                                                      // 3
                                                                                                                       // 4
if (global.parseInt !== parseInt) {                                                                                    // 5
  global.parseInt = parseInt;                                                                                          // 6
}                                                                                                                      // 7
                                                                                                                       // 8
if (global.parseFloat !== parseFloat) {                                                                                // 9
  global.parseFloat = parseFloat;                                                                                      // 10
}                                                                                                                      // 11
                                                                                                                       // 12
var Sp = String.prototype;                                                                                             // 13
if (Sp.replace !== originalStringReplace) {                                                                            // 14
  // Restore the original value of String#replace, because the es5-shim                                                // 15
  // reimplementation is buggy. See also import_globals.js.                                                            // 16
  Sp.replace = originalStringReplace;                                                                                  // 17
}                                                                                                                      // 18
                                                                                                                       // 19
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"import_globals.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/es5-shim/import_globals.js                                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// Because the es5-{shim,sham}.js code assigns to Date and parseInt,                                                   // 1
// Meteor treats them as package variables, and so declares them as                                                    // 2
// variables in package scope, which causes some references to Date and                                                // 3
// parseInt in the shim/sham code to refer to those undefined package                                                  // 4
// variables. The simplest solution seems to be to initialize the package                                              // 5
// variables to their appropriate global values.                                                                       // 6
Date = global.Date;                                                                                                    // 7
parseInt = global.parseInt;                                                                                            // 8
parseFloat = global.parseFloat;                                                                                        // 9
                                                                                                                       // 10
// Save the original String#replace method, because es5-shim's                                                         // 11
// reimplementation of it causes problems in markdown/showdown.js.                                                     // 12
// This original method will be restored in export_globals.js.                                                         // 13
originalStringReplace = String.prototype.replace;                                                                      // 14
                                                                                                                       // 15
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"node_modules":{"es5-shim":{"es5-shim.js":function(require,exports,module){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// node_modules/meteor/es5-shim/node_modules/es5-shim/es5-shim.js                                                      //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
/*!                                                                                                                    // 1
 * https://github.com/es-shims/es5-shim                                                                                // 2
 * @license es5-shim Copyright 2009-2015 by contributors, MIT License                                                  // 3
 * see https://github.com/es-shims/es5-shim/blob/master/LICENSE                                                        // 4
 */                                                                                                                    // 5
                                                                                                                       // 6
// vim: ts=4 sts=4 sw=4 expandtab                                                                                      // 7
                                                                                                                       // 8
// Add semicolon to prevent IIFE from being passed as argument to concatenated code.                                   // 9
;                                                                                                                      // 10
                                                                                                                       // 11
// UMD (Universal Module Definition)                                                                                   // 12
// see https://github.com/umdjs/umd/blob/master/templates/returnExports.js                                             // 13
(function (root, factory) {                                                                                            // 14
    'use strict';                                                                                                      // 15
                                                                                                                       // 16
    /* global define, exports, module */                                                                               // 17
    if (typeof define === 'function' && define.amd) {                                                                  // 18
        // AMD. Register as an anonymous module.                                                                       // 19
        define(factory);                                                                                               // 20
    } else if (typeof exports === 'object') {                                                                          // 21
        // Node. Does not work with strict CommonJS, but                                                               // 22
        // only CommonJS-like enviroments that support module.exports,                                                 // 23
        // like Node.                                                                                                  // 24
        module.exports = factory();                                                                                    // 25
    } else {                                                                                                           // 26
        // Browser globals (root is window)                                                                            // 27
        root.returnExports = factory();                                                                                // 28
    }                                                                                                                  // 29
}(this, function () {                                                                                                  // 30
                                                                                                                       // 31
/**                                                                                                                    // 32
 * Brings an environment as close to ECMAScript 5 compliance                                                           // 33
 * as is possible with the facilities of erstwhile engines.                                                            // 34
 *                                                                                                                     // 35
 * Annotated ES5: http://es5.github.com/ (specific links below)                                                        // 36
 * ES5 Spec: http://www.ecma-international.org/publications/files/ECMA-ST/Ecma-262.pdf                                 // 37
 * Required reading: http://javascriptweblog.wordpress.com/2011/12/05/extending-javascript-natives/                    // 38
 */                                                                                                                    // 39
                                                                                                                       // 40
// Shortcut to an often accessed properties, in order to avoid multiple                                                // 41
// dereference that costs universally. This also holds a reference to known-good                                       // 42
// functions.                                                                                                          // 43
var $Array = Array;                                                                                                    // 44
var ArrayPrototype = $Array.prototype;                                                                                 // 45
var $Object = Object;                                                                                                  // 46
var ObjectPrototype = $Object.prototype;                                                                               // 47
var $Function = Function;                                                                                              // 48
var FunctionPrototype = $Function.prototype;                                                                           // 49
var $String = String;                                                                                                  // 50
var StringPrototype = $String.prototype;                                                                               // 51
var $Number = Number;                                                                                                  // 52
var NumberPrototype = $Number.prototype;                                                                               // 53
var array_slice = ArrayPrototype.slice;                                                                                // 54
var array_splice = ArrayPrototype.splice;                                                                              // 55
var array_push = ArrayPrototype.push;                                                                                  // 56
var array_unshift = ArrayPrototype.unshift;                                                                            // 57
var array_concat = ArrayPrototype.concat;                                                                              // 58
var array_join = ArrayPrototype.join;                                                                                  // 59
var call = FunctionPrototype.call;                                                                                     // 60
var apply = FunctionPrototype.apply;                                                                                   // 61
var max = Math.max;                                                                                                    // 62
var min = Math.min;                                                                                                    // 63
                                                                                                                       // 64
// Having a toString local variable name breaks in Opera so use to_string.                                             // 65
var to_string = ObjectPrototype.toString;                                                                              // 66
                                                                                                                       // 67
/* global Symbol */                                                                                                    // 68
/* eslint-disable one-var-declaration-per-line, no-redeclare */                                                        // 69
var hasToStringTag = typeof Symbol === 'function' && typeof Symbol.toStringTag === 'symbol';                           // 70
var isCallable; /* inlined from https://npmjs.com/is-callable */ var fnToStr = Function.prototype.toString, constructorRegex = /^\s*class /, isES6ClassFn = function isES6ClassFn(value) { try { var fnStr = fnToStr.call(value); var singleStripped = fnStr.replace(/\/\/.*\n/g, ''); var multiStripped = singleStripped.replace(/\/\*[.\s\S]*\*\//g, ''); var spaceStripped = multiStripped.replace(/\n/mg, ' ').replace(/ {2}/g, ' '); return constructorRegex.test(spaceStripped); } catch (e) { return false; /* not a function */ } }, tryFunctionObject = function tryFunctionObject(value) { try { if (isES6ClassFn(value)) { return false; } fnToStr.call(value); return true; } catch (e) { return false; } }, fnClass = '[object Function]', genClass = '[object GeneratorFunction]', isCallable = function isCallable(value) { if (!value) { return false; } if (typeof value !== 'function' && typeof value !== 'object') { return false; } if (hasToStringTag) { return tryFunctionObject(value); } if (isES6ClassFn(value)) { return false; } var strClass = to_string.call(value); return strClass === fnClass || strClass === genClass; };
                                                                                                                       // 72
var isRegex; /* inlined from https://npmjs.com/is-regex */ var regexExec = RegExp.prototype.exec, tryRegexExec = function tryRegexExec(value) { try { regexExec.call(value); return true; } catch (e) { return false; } }, regexClass = '[object RegExp]'; isRegex = function isRegex(value) { if (typeof value !== 'object') { return false; } return hasToStringTag ? tryRegexExec(value) : to_string.call(value) === regexClass; };
var isString; /* inlined from https://npmjs.com/is-string */ var strValue = String.prototype.valueOf, tryStringObject = function tryStringObject(value) { try { strValue.call(value); return true; } catch (e) { return false; } }, stringClass = '[object String]'; isString = function isString(value) { if (typeof value === 'string') { return true; } if (typeof value !== 'object') { return false; } return hasToStringTag ? tryStringObject(value) : to_string.call(value) === stringClass; };
/* eslint-enable one-var-declaration-per-line, no-redeclare */                                                         // 75
                                                                                                                       // 76
/* inlined from http://npmjs.com/define-properties */                                                                  // 77
var supportsDescriptors = $Object.defineProperty && (function () {                                                     // 78
    try {                                                                                                              // 79
        var obj = {};                                                                                                  // 80
        $Object.defineProperty(obj, 'x', { enumerable: false, value: obj });                                           // 81
        for (var _ in obj) { return false; }                                                                           // 82
        return obj.x === obj;                                                                                          // 83
    } catch (e) { /* this is ES3 */                                                                                    // 84
        return false;                                                                                                  // 85
    }                                                                                                                  // 86
}());                                                                                                                  // 87
var defineProperties = (function (has) {                                                                               // 88
  // Define configurable, writable, and non-enumerable props                                                           // 89
  // if they don't exist.                                                                                              // 90
  var defineProperty;                                                                                                  // 91
  if (supportsDescriptors) {                                                                                           // 92
      defineProperty = function (object, name, method, forceAssign) {                                                  // 93
          if (!forceAssign && (name in object)) { return; }                                                            // 94
          $Object.defineProperty(object, name, {                                                                       // 95
              configurable: true,                                                                                      // 96
              enumerable: false,                                                                                       // 97
              writable: true,                                                                                          // 98
              value: method                                                                                            // 99
          });                                                                                                          // 100
      };                                                                                                               // 101
  } else {                                                                                                             // 102
      defineProperty = function (object, name, method, forceAssign) {                                                  // 103
          if (!forceAssign && (name in object)) { return; }                                                            // 104
          object[name] = method;                                                                                       // 105
      };                                                                                                               // 106
  }                                                                                                                    // 107
  return function defineProperties(object, map, forceAssign) {                                                         // 108
      for (var name in map) {                                                                                          // 109
          if (has.call(map, name)) {                                                                                   // 110
            defineProperty(object, name, map[name], forceAssign);                                                      // 111
          }                                                                                                            // 112
      }                                                                                                                // 113
  };                                                                                                                   // 114
}(ObjectPrototype.hasOwnProperty));                                                                                    // 115
                                                                                                                       // 116
//                                                                                                                     // 117
// Util                                                                                                                // 118
// ======                                                                                                              // 119
//                                                                                                                     // 120
                                                                                                                       // 121
/* replaceable with https://npmjs.com/package/es-abstract /helpers/isPrimitive */                                      // 122
var isPrimitive = function isPrimitive(input) {                                                                        // 123
    var type = typeof input;                                                                                           // 124
    return input === null || (type !== 'object' && type !== 'function');                                               // 125
};                                                                                                                     // 126
                                                                                                                       // 127
var isActualNaN = $Number.isNaN || function (x) { return x !== x; };                                                   // 128
                                                                                                                       // 129
var ES = {                                                                                                             // 130
    // ES5 9.4                                                                                                         // 131
    // http://es5.github.com/#x9.4                                                                                     // 132
    // http://jsperf.com/to-integer                                                                                    // 133
    /* replaceable with https://npmjs.com/package/es-abstract ES5.ToInteger */                                         // 134
    ToInteger: function ToInteger(num) {                                                                               // 135
        var n = +num;                                                                                                  // 136
        if (isActualNaN(n)) {                                                                                          // 137
            n = 0;                                                                                                     // 138
        } else if (n !== 0 && n !== (1 / 0) && n !== -(1 / 0)) {                                                       // 139
            n = (n > 0 || -1) * Math.floor(Math.abs(n));                                                               // 140
        }                                                                                                              // 141
        return n;                                                                                                      // 142
    },                                                                                                                 // 143
                                                                                                                       // 144
    /* replaceable with https://npmjs.com/package/es-abstract ES5.ToPrimitive */                                       // 145
    ToPrimitive: function ToPrimitive(input) {                                                                         // 146
        var val, valueOf, toStr;                                                                                       // 147
        if (isPrimitive(input)) {                                                                                      // 148
            return input;                                                                                              // 149
        }                                                                                                              // 150
        valueOf = input.valueOf;                                                                                       // 151
        if (isCallable(valueOf)) {                                                                                     // 152
            val = valueOf.call(input);                                                                                 // 153
            if (isPrimitive(val)) {                                                                                    // 154
                return val;                                                                                            // 155
            }                                                                                                          // 156
        }                                                                                                              // 157
        toStr = input.toString;                                                                                        // 158
        if (isCallable(toStr)) {                                                                                       // 159
            val = toStr.call(input);                                                                                   // 160
            if (isPrimitive(val)) {                                                                                    // 161
                return val;                                                                                            // 162
            }                                                                                                          // 163
        }                                                                                                              // 164
        throw new TypeError();                                                                                         // 165
    },                                                                                                                 // 166
                                                                                                                       // 167
    // ES5 9.9                                                                                                         // 168
    // http://es5.github.com/#x9.9                                                                                     // 169
    /* replaceable with https://npmjs.com/package/es-abstract ES5.ToObject */                                          // 170
    ToObject: function (o) {                                                                                           // 171
        if (o == null) { // this matches both null and undefined                                                       // 172
            throw new TypeError("can't convert " + o + ' to object');                                                  // 173
        }                                                                                                              // 174
        return $Object(o);                                                                                             // 175
    },                                                                                                                 // 176
                                                                                                                       // 177
    /* replaceable with https://npmjs.com/package/es-abstract ES5.ToUint32 */                                          // 178
    ToUint32: function ToUint32(x) {                                                                                   // 179
        return x >>> 0;                                                                                                // 180
    }                                                                                                                  // 181
};                                                                                                                     // 182
                                                                                                                       // 183
//                                                                                                                     // 184
// Function                                                                                                            // 185
// ========                                                                                                            // 186
//                                                                                                                     // 187
                                                                                                                       // 188
// ES-5 15.3.4.5                                                                                                       // 189
// http://es5.github.com/#x15.3.4.5                                                                                    // 190
                                                                                                                       // 191
var Empty = function Empty() {};                                                                                       // 192
                                                                                                                       // 193
defineProperties(FunctionPrototype, {                                                                                  // 194
    bind: function bind(that) { // .length is 1                                                                        // 195
        // 1. Let Target be the this value.                                                                            // 196
        var target = this;                                                                                             // 197
        // 2. If IsCallable(Target) is false, throw a TypeError exception.                                             // 198
        if (!isCallable(target)) {                                                                                     // 199
            throw new TypeError('Function.prototype.bind called on incompatible ' + target);                           // 200
        }                                                                                                              // 201
        // 3. Let A be a new (possibly empty) internal list of all of the                                              // 202
        //   argument values provided after thisArg (arg1, arg2 etc), in order.                                        // 203
        // XXX slicedArgs will stand in for "A" if used                                                                // 204
        var args = array_slice.call(arguments, 1); // for normal call                                                  // 205
        // 4. Let F be a new native ECMAScript object.                                                                 // 206
        // 11. Set the [[Prototype]] internal property of F to the standard                                            // 207
        //   built-in Function prototype object as specified in 15.3.3.1.                                              // 208
        // 12. Set the [[Call]] internal property of F as described in                                                 // 209
        //   15.3.4.5.1.                                                                                               // 210
        // 13. Set the [[Construct]] internal property of F as described in                                            // 211
        //   15.3.4.5.2.                                                                                               // 212
        // 14. Set the [[HasInstance]] internal property of F as described in                                          // 213
        //   15.3.4.5.3.                                                                                               // 214
        var bound;                                                                                                     // 215
        var binder = function () {                                                                                     // 216
                                                                                                                       // 217
            if (this instanceof bound) {                                                                               // 218
                // 15.3.4.5.2 [[Construct]]                                                                            // 219
                // When the [[Construct]] internal method of a function object,                                        // 220
                // F that was created using the bind function is called with a                                         // 221
                // list of arguments ExtraArgs, the following steps are taken:                                         // 222
                // 1. Let target be the value of F's [[TargetFunction]]                                                // 223
                //   internal property.                                                                                // 224
                // 2. If target has no [[Construct]] internal method, a                                                // 225
                //   TypeError exception is thrown.                                                                    // 226
                // 3. Let boundArgs be the value of F's [[BoundArgs]] internal                                         // 227
                //   property.                                                                                         // 228
                // 4. Let args be a new list containing the same values as the                                         // 229
                //   list boundArgs in the same order followed by the same                                             // 230
                //   values as the list ExtraArgs in the same order.                                                   // 231
                // 5. Return the result of calling the [[Construct]] internal                                          // 232
                //   method of target providing args as the arguments.                                                 // 233
                                                                                                                       // 234
                var result = apply.call(                                                                               // 235
                    target,                                                                                            // 236
                    this,                                                                                              // 237
                    array_concat.call(args, array_slice.call(arguments))                                               // 238
                );                                                                                                     // 239
                if ($Object(result) === result) {                                                                      // 240
                    return result;                                                                                     // 241
                }                                                                                                      // 242
                return this;                                                                                           // 243
                                                                                                                       // 244
            } else {                                                                                                   // 245
                // 15.3.4.5.1 [[Call]]                                                                                 // 246
                // When the [[Call]] internal method of a function object, F,                                          // 247
                // which was created using the bind function is called with a                                          // 248
                // this value and a list of arguments ExtraArgs, the following                                         // 249
                // steps are taken:                                                                                    // 250
                // 1. Let boundArgs be the value of F's [[BoundArgs]] internal                                         // 251
                //   property.                                                                                         // 252
                // 2. Let boundThis be the value of F's [[BoundThis]] internal                                         // 253
                //   property.                                                                                         // 254
                // 3. Let target be the value of F's [[TargetFunction]] internal                                       // 255
                //   property.                                                                                         // 256
                // 4. Let args be a new list containing the same values as the                                         // 257
                //   list boundArgs in the same order followed by the same                                             // 258
                //   values as the list ExtraArgs in the same order.                                                   // 259
                // 5. Return the result of calling the [[Call]] internal method                                        // 260
                //   of target providing boundThis as the this value and                                               // 261
                //   providing args as the arguments.                                                                  // 262
                                                                                                                       // 263
                // equiv: target.call(this, ...boundArgs, ...args)                                                     // 264
                return apply.call(                                                                                     // 265
                    target,                                                                                            // 266
                    that,                                                                                              // 267
                    array_concat.call(args, array_slice.call(arguments))                                               // 268
                );                                                                                                     // 269
                                                                                                                       // 270
            }                                                                                                          // 271
                                                                                                                       // 272
        };                                                                                                             // 273
                                                                                                                       // 274
        // 15. If the [[Class]] internal property of Target is "Function", then                                        // 275
        //     a. Let L be the length property of Target minus the length of A.                                        // 276
        //     b. Set the length own property of F to either 0 or L, whichever is                                      // 277
        //       larger.                                                                                               // 278
        // 16. Else set the length own property of F to 0.                                                             // 279
                                                                                                                       // 280
        var boundLength = max(0, target.length - args.length);                                                         // 281
                                                                                                                       // 282
        // 17. Set the attributes of the length own property of F to the values                                        // 283
        //   specified in 15.3.5.1.                                                                                    // 284
        var boundArgs = [];                                                                                            // 285
        for (var i = 0; i < boundLength; i++) {                                                                        // 286
            array_push.call(boundArgs, '$' + i);                                                                       // 287
        }                                                                                                              // 288
                                                                                                                       // 289
        // XXX Build a dynamic function with desired amount of arguments is the only                                   // 290
        // way to set the length property of a function.                                                               // 291
        // In environments where Content Security Policies enabled (Chrome extensions,                                 // 292
        // for ex.) all use of eval or Function costructor throws an exception.                                        // 293
        // However in all of these environments Function.prototype.bind exists                                         // 294
        // and so this code will never be executed.                                                                    // 295
        bound = $Function('binder', 'return function (' + array_join.call(boundArgs, ',') + '){ return binder.apply(this, arguments); }')(binder);
                                                                                                                       // 297
        if (target.prototype) {                                                                                        // 298
            Empty.prototype = target.prototype;                                                                        // 299
            bound.prototype = new Empty();                                                                             // 300
            // Clean up dangling references.                                                                           // 301
            Empty.prototype = null;                                                                                    // 302
        }                                                                                                              // 303
                                                                                                                       // 304
        // TODO                                                                                                        // 305
        // 18. Set the [[Extensible]] internal property of F to true.                                                  // 306
                                                                                                                       // 307
        // TODO                                                                                                        // 308
        // 19. Let thrower be the [[ThrowTypeError]] function Object (13.2.3).                                         // 309
        // 20. Call the [[DefineOwnProperty]] internal method of F with                                                // 310
        //   arguments "caller", PropertyDescriptor {[[Get]]: thrower, [[Set]]:                                        // 311
        //   thrower, [[Enumerable]]: false, [[Configurable]]: false}, and                                             // 312
        //   false.                                                                                                    // 313
        // 21. Call the [[DefineOwnProperty]] internal method of F with                                                // 314
        //   arguments "arguments", PropertyDescriptor {[[Get]]: thrower,                                              // 315
        //   [[Set]]: thrower, [[Enumerable]]: false, [[Configurable]]: false},                                        // 316
        //   and false.                                                                                                // 317
                                                                                                                       // 318
        // TODO                                                                                                        // 319
        // NOTE Function objects created using Function.prototype.bind do not                                          // 320
        // have a prototype property or the [[Code]], [[FormalParameters]], and                                        // 321
        // [[Scope]] internal properties.                                                                              // 322
        // XXX can't delete prototype in pure-js.                                                                      // 323
                                                                                                                       // 324
        // 22. Return F.                                                                                               // 325
        return bound;                                                                                                  // 326
    }                                                                                                                  // 327
});                                                                                                                    // 328
                                                                                                                       // 329
// _Please note: Shortcuts are defined after `Function.prototype.bind` as we                                           // 330
// use it in defining shortcuts.                                                                                       // 331
var owns = call.bind(ObjectPrototype.hasOwnProperty);                                                                  // 332
var toStr = call.bind(ObjectPrototype.toString);                                                                       // 333
var arraySlice = call.bind(array_slice);                                                                               // 334
var arraySliceApply = apply.bind(array_slice);                                                                         // 335
var strSlice = call.bind(StringPrototype.slice);                                                                       // 336
var strSplit = call.bind(StringPrototype.split);                                                                       // 337
var strIndexOf = call.bind(StringPrototype.indexOf);                                                                   // 338
var pushCall = call.bind(array_push);                                                                                  // 339
var isEnum = call.bind(ObjectPrototype.propertyIsEnumerable);                                                          // 340
var arraySort = call.bind(ArrayPrototype.sort);                                                                        // 341
                                                                                                                       // 342
//                                                                                                                     // 343
// Array                                                                                                               // 344
// =====                                                                                                               // 345
//                                                                                                                     // 346
                                                                                                                       // 347
var isArray = $Array.isArray || function isArray(obj) {                                                                // 348
    return toStr(obj) === '[object Array]';                                                                            // 349
};                                                                                                                     // 350
                                                                                                                       // 351
// ES5 15.4.4.12                                                                                                       // 352
// http://es5.github.com/#x15.4.4.13                                                                                   // 353
// Return len+argCount.                                                                                                // 354
// [bugfix, ielt8]                                                                                                     // 355
// IE < 8 bug: [].unshift(0) === undefined but should be "1"                                                           // 356
var hasUnshiftReturnValueBug = [].unshift(0) !== 1;                                                                    // 357
defineProperties(ArrayPrototype, {                                                                                     // 358
    unshift: function () {                                                                                             // 359
        array_unshift.apply(this, arguments);                                                                          // 360
        return this.length;                                                                                            // 361
    }                                                                                                                  // 362
}, hasUnshiftReturnValueBug);                                                                                          // 363
                                                                                                                       // 364
// ES5 15.4.3.2                                                                                                        // 365
// http://es5.github.com/#x15.4.3.2                                                                                    // 366
// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Array/isArray                                  // 367
defineProperties($Array, { isArray: isArray });                                                                        // 368
                                                                                                                       // 369
// The IsCallable() check in the Array functions                                                                       // 370
// has been replaced with a strict check on the                                                                        // 371
// internal class of the object to trap cases where                                                                    // 372
// the provided function was actually a regular                                                                        // 373
// expression literal, which in V8 and                                                                                 // 374
// JavaScriptCore is a typeof "function".  Only in                                                                     // 375
// V8 are regular expression literals permitted as                                                                     // 376
// reduce parameters, so it is desirable in the                                                                        // 377
// general case for the shim to match the more                                                                         // 378
// strict and common behavior of rejecting regular                                                                     // 379
// expressions.                                                                                                        // 380
                                                                                                                       // 381
// ES5 15.4.4.18                                                                                                       // 382
// http://es5.github.com/#x15.4.4.18                                                                                   // 383
// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/array/forEach                                  // 384
                                                                                                                       // 385
// Check failure of by-index access of string characters (IE < 9)                                                      // 386
// and failure of `0 in boxedString` (Rhino)                                                                           // 387
var boxedString = $Object('a');                                                                                        // 388
var splitString = boxedString[0] !== 'a' || !(0 in boxedString);                                                       // 389
                                                                                                                       // 390
var properlyBoxesContext = function properlyBoxed(method) {                                                            // 391
    // Check node 0.6.21 bug where third parameter is not boxed                                                        // 392
    var properlyBoxesNonStrict = true;                                                                                 // 393
    var properlyBoxesStrict = true;                                                                                    // 394
    var threwException = false;                                                                                        // 395
    if (method) {                                                                                                      // 396
        try {                                                                                                          // 397
            method.call('foo', function (_, __, context) {                                                             // 398
                if (typeof context !== 'object') { properlyBoxesNonStrict = false; }                                   // 399
            });                                                                                                        // 400
                                                                                                                       // 401
            method.call([1], function () {                                                                             // 402
                'use strict';                                                                                          // 403
                                                                                                                       // 404
                properlyBoxesStrict = typeof this === 'string';                                                        // 405
            }, 'x');                                                                                                   // 406
        } catch (e) {                                                                                                  // 407
            threwException = true;                                                                                     // 408
        }                                                                                                              // 409
    }                                                                                                                  // 410
    return !!method && !threwException && properlyBoxesNonStrict && properlyBoxesStrict;                               // 411
};                                                                                                                     // 412
                                                                                                                       // 413
defineProperties(ArrayPrototype, {                                                                                     // 414
    forEach: function forEach(callbackfn/*, thisArg*/) {                                                               // 415
        var object = ES.ToObject(this);                                                                                // 416
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 417
        var i = -1;                                                                                                    // 418
        var length = ES.ToUint32(self.length);                                                                         // 419
        var T;                                                                                                         // 420
        if (arguments.length > 1) {                                                                                    // 421
          T = arguments[1];                                                                                            // 422
        }                                                                                                              // 423
                                                                                                                       // 424
        // If no callback function or if callback is not a callable function                                           // 425
        if (!isCallable(callbackfn)) {                                                                                 // 426
            throw new TypeError('Array.prototype.forEach callback must be a function');                                // 427
        }                                                                                                              // 428
                                                                                                                       // 429
        while (++i < length) {                                                                                         // 430
            if (i in self) {                                                                                           // 431
                // Invoke the callback function with call, passing arguments:                                          // 432
                // context, property value, property key, thisArg object                                               // 433
                if (typeof T === 'undefined') {                                                                        // 434
                    callbackfn(self[i], i, object);                                                                    // 435
                } else {                                                                                               // 436
                    callbackfn.call(T, self[i], i, object);                                                            // 437
                }                                                                                                      // 438
            }                                                                                                          // 439
        }                                                                                                              // 440
    }                                                                                                                  // 441
}, !properlyBoxesContext(ArrayPrototype.forEach));                                                                     // 442
                                                                                                                       // 443
// ES5 15.4.4.19                                                                                                       // 444
// http://es5.github.com/#x15.4.4.19                                                                                   // 445
// https://developer.mozilla.org/en/Core_JavaScript_1.5_Reference/Objects/Array/map                                    // 446
defineProperties(ArrayPrototype, {                                                                                     // 447
    map: function map(callbackfn/*, thisArg*/) {                                                                       // 448
        var object = ES.ToObject(this);                                                                                // 449
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 450
        var length = ES.ToUint32(self.length);                                                                         // 451
        var result = $Array(length);                                                                                   // 452
        var T;                                                                                                         // 453
        if (arguments.length > 1) {                                                                                    // 454
            T = arguments[1];                                                                                          // 455
        }                                                                                                              // 456
                                                                                                                       // 457
        // If no callback function or if callback is not a callable function                                           // 458
        if (!isCallable(callbackfn)) {                                                                                 // 459
            throw new TypeError('Array.prototype.map callback must be a function');                                    // 460
        }                                                                                                              // 461
                                                                                                                       // 462
        for (var i = 0; i < length; i++) {                                                                             // 463
            if (i in self) {                                                                                           // 464
                if (typeof T === 'undefined') {                                                                        // 465
                    result[i] = callbackfn(self[i], i, object);                                                        // 466
                } else {                                                                                               // 467
                    result[i] = callbackfn.call(T, self[i], i, object);                                                // 468
                }                                                                                                      // 469
            }                                                                                                          // 470
        }                                                                                                              // 471
        return result;                                                                                                 // 472
    }                                                                                                                  // 473
}, !properlyBoxesContext(ArrayPrototype.map));                                                                         // 474
                                                                                                                       // 475
// ES5 15.4.4.20                                                                                                       // 476
// http://es5.github.com/#x15.4.4.20                                                                                   // 477
// https://developer.mozilla.org/en/Core_JavaScript_1.5_Reference/Objects/Array/filter                                 // 478
defineProperties(ArrayPrototype, {                                                                                     // 479
    filter: function filter(callbackfn/*, thisArg*/) {                                                                 // 480
        var object = ES.ToObject(this);                                                                                // 481
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 482
        var length = ES.ToUint32(self.length);                                                                         // 483
        var result = [];                                                                                               // 484
        var value;                                                                                                     // 485
        var T;                                                                                                         // 486
        if (arguments.length > 1) {                                                                                    // 487
            T = arguments[1];                                                                                          // 488
        }                                                                                                              // 489
                                                                                                                       // 490
        // If no callback function or if callback is not a callable function                                           // 491
        if (!isCallable(callbackfn)) {                                                                                 // 492
            throw new TypeError('Array.prototype.filter callback must be a function');                                 // 493
        }                                                                                                              // 494
                                                                                                                       // 495
        for (var i = 0; i < length; i++) {                                                                             // 496
            if (i in self) {                                                                                           // 497
                value = self[i];                                                                                       // 498
                if (typeof T === 'undefined' ? callbackfn(value, i, object) : callbackfn.call(T, value, i, object)) {  // 499
                    pushCall(result, value);                                                                           // 500
                }                                                                                                      // 501
            }                                                                                                          // 502
        }                                                                                                              // 503
        return result;                                                                                                 // 504
    }                                                                                                                  // 505
}, !properlyBoxesContext(ArrayPrototype.filter));                                                                      // 506
                                                                                                                       // 507
// ES5 15.4.4.16                                                                                                       // 508
// http://es5.github.com/#x15.4.4.16                                                                                   // 509
// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Array/every                                    // 510
defineProperties(ArrayPrototype, {                                                                                     // 511
    every: function every(callbackfn/*, thisArg*/) {                                                                   // 512
        var object = ES.ToObject(this);                                                                                // 513
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 514
        var length = ES.ToUint32(self.length);                                                                         // 515
        var T;                                                                                                         // 516
        if (arguments.length > 1) {                                                                                    // 517
            T = arguments[1];                                                                                          // 518
        }                                                                                                              // 519
                                                                                                                       // 520
        // If no callback function or if callback is not a callable function                                           // 521
        if (!isCallable(callbackfn)) {                                                                                 // 522
            throw new TypeError('Array.prototype.every callback must be a function');                                  // 523
        }                                                                                                              // 524
                                                                                                                       // 525
        for (var i = 0; i < length; i++) {                                                                             // 526
            if (i in self && !(typeof T === 'undefined' ? callbackfn(self[i], i, object) : callbackfn.call(T, self[i], i, object))) {
                return false;                                                                                          // 528
            }                                                                                                          // 529
        }                                                                                                              // 530
        return true;                                                                                                   // 531
    }                                                                                                                  // 532
}, !properlyBoxesContext(ArrayPrototype.every));                                                                       // 533
                                                                                                                       // 534
// ES5 15.4.4.17                                                                                                       // 535
// http://es5.github.com/#x15.4.4.17                                                                                   // 536
// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Array/some                                     // 537
defineProperties(ArrayPrototype, {                                                                                     // 538
    some: function some(callbackfn/*, thisArg */) {                                                                    // 539
        var object = ES.ToObject(this);                                                                                // 540
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 541
        var length = ES.ToUint32(self.length);                                                                         // 542
        var T;                                                                                                         // 543
        if (arguments.length > 1) {                                                                                    // 544
            T = arguments[1];                                                                                          // 545
        }                                                                                                              // 546
                                                                                                                       // 547
        // If no callback function or if callback is not a callable function                                           // 548
        if (!isCallable(callbackfn)) {                                                                                 // 549
            throw new TypeError('Array.prototype.some callback must be a function');                                   // 550
        }                                                                                                              // 551
                                                                                                                       // 552
        for (var i = 0; i < length; i++) {                                                                             // 553
            if (i in self && (typeof T === 'undefined' ? callbackfn(self[i], i, object) : callbackfn.call(T, self[i], i, object))) {
                return true;                                                                                           // 555
            }                                                                                                          // 556
        }                                                                                                              // 557
        return false;                                                                                                  // 558
    }                                                                                                                  // 559
}, !properlyBoxesContext(ArrayPrototype.some));                                                                        // 560
                                                                                                                       // 561
// ES5 15.4.4.21                                                                                                       // 562
// http://es5.github.com/#x15.4.4.21                                                                                   // 563
// https://developer.mozilla.org/en/Core_JavaScript_1.5_Reference/Objects/Array/reduce                                 // 564
var reduceCoercesToObject = false;                                                                                     // 565
if (ArrayPrototype.reduce) {                                                                                           // 566
    reduceCoercesToObject = typeof ArrayPrototype.reduce.call('es5', function (_, __, ___, list) { return list; }) === 'object';
}                                                                                                                      // 568
defineProperties(ArrayPrototype, {                                                                                     // 569
    reduce: function reduce(callbackfn/*, initialValue*/) {                                                            // 570
        var object = ES.ToObject(this);                                                                                // 571
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 572
        var length = ES.ToUint32(self.length);                                                                         // 573
                                                                                                                       // 574
        // If no callback function or if callback is not a callable function                                           // 575
        if (!isCallable(callbackfn)) {                                                                                 // 576
            throw new TypeError('Array.prototype.reduce callback must be a function');                                 // 577
        }                                                                                                              // 578
                                                                                                                       // 579
        // no value to return if no initial value and an empty array                                                   // 580
        if (length === 0 && arguments.length === 1) {                                                                  // 581
            throw new TypeError('reduce of empty array with no initial value');                                        // 582
        }                                                                                                              // 583
                                                                                                                       // 584
        var i = 0;                                                                                                     // 585
        var result;                                                                                                    // 586
        if (arguments.length >= 2) {                                                                                   // 587
            result = arguments[1];                                                                                     // 588
        } else {                                                                                                       // 589
            do {                                                                                                       // 590
                if (i in self) {                                                                                       // 591
                    result = self[i++];                                                                                // 592
                    break;                                                                                             // 593
                }                                                                                                      // 594
                                                                                                                       // 595
                // if array contains no values, no initial value to return                                             // 596
                if (++i >= length) {                                                                                   // 597
                    throw new TypeError('reduce of empty array with no initial value');                                // 598
                }                                                                                                      // 599
            } while (true);                                                                                            // 600
        }                                                                                                              // 601
                                                                                                                       // 602
        for (; i < length; i++) {                                                                                      // 603
            if (i in self) {                                                                                           // 604
                result = callbackfn(result, self[i], i, object);                                                       // 605
            }                                                                                                          // 606
        }                                                                                                              // 607
                                                                                                                       // 608
        return result;                                                                                                 // 609
    }                                                                                                                  // 610
}, !reduceCoercesToObject);                                                                                            // 611
                                                                                                                       // 612
// ES5 15.4.4.22                                                                                                       // 613
// http://es5.github.com/#x15.4.4.22                                                                                   // 614
// https://developer.mozilla.org/en/Core_JavaScript_1.5_Reference/Objects/Array/reduceRight                            // 615
var reduceRightCoercesToObject = false;                                                                                // 616
if (ArrayPrototype.reduceRight) {                                                                                      // 617
    reduceRightCoercesToObject = typeof ArrayPrototype.reduceRight.call('es5', function (_, __, ___, list) { return list; }) === 'object';
}                                                                                                                      // 619
defineProperties(ArrayPrototype, {                                                                                     // 620
    reduceRight: function reduceRight(callbackfn/*, initial*/) {                                                       // 621
        var object = ES.ToObject(this);                                                                                // 622
        var self = splitString && isString(this) ? strSplit(this, '') : object;                                        // 623
        var length = ES.ToUint32(self.length);                                                                         // 624
                                                                                                                       // 625
        // If no callback function or if callback is not a callable function                                           // 626
        if (!isCallable(callbackfn)) {                                                                                 // 627
            throw new TypeError('Array.prototype.reduceRight callback must be a function');                            // 628
        }                                                                                                              // 629
                                                                                                                       // 630
        // no value to return if no initial value, empty array                                                         // 631
        if (length === 0 && arguments.length === 1) {                                                                  // 632
            throw new TypeError('reduceRight of empty array with no initial value');                                   // 633
        }                                                                                                              // 634
                                                                                                                       // 635
        var result;                                                                                                    // 636
        var i = length - 1;                                                                                            // 637
        if (arguments.length >= 2) {                                                                                   // 638
            result = arguments[1];                                                                                     // 639
        } else {                                                                                                       // 640
            do {                                                                                                       // 641
                if (i in self) {                                                                                       // 642
                    result = self[i--];                                                                                // 643
                    break;                                                                                             // 644
                }                                                                                                      // 645
                                                                                                                       // 646
                // if array contains no values, no initial value to return                                             // 647
                if (--i < 0) {                                                                                         // 648
                    throw new TypeError('reduceRight of empty array with no initial value');                           // 649
                }                                                                                                      // 650
            } while (true);                                                                                            // 651
        }                                                                                                              // 652
                                                                                                                       // 653
        if (i < 0) {                                                                                                   // 654
            return result;                                                                                             // 655
        }                                                                                                              // 656
                                                                                                                       // 657
        do {                                                                                                           // 658
            if (i in self) {                                                                                           // 659
                result = callbackfn(result, self[i], i, object);                                                       // 660
            }                                                                                                          // 661
        } while (i--);                                                                                                 // 662
                                                                                                                       // 663
        return result;                                                                                                 // 664
    }                                                                                                                  // 665
}, !reduceRightCoercesToObject);                                                                                       // 666
                                                                                                                       // 667
// ES5 15.4.4.14                                                                                                       // 668
// http://es5.github.com/#x15.4.4.14                                                                                   // 669
// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Array/indexOf                                  // 670
var hasFirefox2IndexOfBug = ArrayPrototype.indexOf && [0, 1].indexOf(1, 2) !== -1;                                     // 671
defineProperties(ArrayPrototype, {                                                                                     // 672
    indexOf: function indexOf(searchElement/*, fromIndex */) {                                                         // 673
        var self = splitString && isString(this) ? strSplit(this, '') : ES.ToObject(this);                             // 674
        var length = ES.ToUint32(self.length);                                                                         // 675
                                                                                                                       // 676
        if (length === 0) {                                                                                            // 677
            return -1;                                                                                                 // 678
        }                                                                                                              // 679
                                                                                                                       // 680
        var i = 0;                                                                                                     // 681
        if (arguments.length > 1) {                                                                                    // 682
            i = ES.ToInteger(arguments[1]);                                                                            // 683
        }                                                                                                              // 684
                                                                                                                       // 685
        // handle negative indices                                                                                     // 686
        i = i >= 0 ? i : max(0, length + i);                                                                           // 687
        for (; i < length; i++) {                                                                                      // 688
            if (i in self && self[i] === searchElement) {                                                              // 689
                return i;                                                                                              // 690
            }                                                                                                          // 691
        }                                                                                                              // 692
        return -1;                                                                                                     // 693
    }                                                                                                                  // 694
}, hasFirefox2IndexOfBug);                                                                                             // 695
                                                                                                                       // 696
// ES5 15.4.4.15                                                                                                       // 697
// http://es5.github.com/#x15.4.4.15                                                                                   // 698
// https://developer.mozilla.org/en/JavaScript/Reference/Global_Objects/Array/lastIndexOf                              // 699
var hasFirefox2LastIndexOfBug = ArrayPrototype.lastIndexOf && [0, 1].lastIndexOf(0, -3) !== -1;                        // 700
defineProperties(ArrayPrototype, {                                                                                     // 701
    lastIndexOf: function lastIndexOf(searchElement/*, fromIndex */) {                                                 // 702
        var self = splitString && isString(this) ? strSplit(this, '') : ES.ToObject(this);                             // 703
        var length = ES.ToUint32(self.length);                                                                         // 704
                                                                                                                       // 705
        if (length === 0) {                                                                                            // 706
            return -1;                                                                                                 // 707
        }                                                                                                              // 708
        var i = length - 1;                                                                                            // 709
        if (arguments.length > 1) {                                                                                    // 710
            i = min(i, ES.ToInteger(arguments[1]));                                                                    // 711
        }                                                                                                              // 712
        // handle negative indices                                                                                     // 713
        i = i >= 0 ? i : length - Math.abs(i);                                                                         // 714
        for (; i >= 0; i--) {                                                                                          // 715
            if (i in self && searchElement === self[i]) {                                                              // 716
                return i;                                                                                              // 717
            }                                                                                                          // 718
        }                                                                                                              // 719
        return -1;                                                                                                     // 720
    }                                                                                                                  // 721
}, hasFirefox2LastIndexOfBug);                                                                                         // 722
                                                                                                                       // 723
// ES5 15.4.4.12                                                                                                       // 724
// http://es5.github.com/#x15.4.4.12                                                                                   // 725
var spliceNoopReturnsEmptyArray = (function () {                                                                       // 726
    var a = [1, 2];                                                                                                    // 727
    var result = a.splice();                                                                                           // 728
    return a.length === 2 && isArray(result) && result.length === 0;                                                   // 729
}());                                                                                                                  // 730
defineProperties(ArrayPrototype, {                                                                                     // 731
    // Safari 5.0 bug where .splice() returns undefined                                                                // 732
    splice: function splice(start, deleteCount) {                                                                      // 733
        if (arguments.length === 0) {                                                                                  // 734
            return [];                                                                                                 // 735
        } else {                                                                                                       // 736
            return array_splice.apply(this, arguments);                                                                // 737
        }                                                                                                              // 738
    }                                                                                                                  // 739
}, !spliceNoopReturnsEmptyArray);                                                                                      // 740
                                                                                                                       // 741
var spliceWorksWithEmptyObject = (function () {                                                                        // 742
    var obj = {};                                                                                                      // 743
    ArrayPrototype.splice.call(obj, 0, 0, 1);                                                                          // 744
    return obj.length === 1;                                                                                           // 745
}());                                                                                                                  // 746
defineProperties(ArrayPrototype, {                                                                                     // 747
    splice: function splice(start, deleteCount) {                                                                      // 748
        if (arguments.length === 0) { return []; }                                                                     // 749
        var args = arguments;                                                                                          // 750
        this.length = max(ES.ToInteger(this.length), 0);                                                               // 751
        if (arguments.length > 0 && typeof deleteCount !== 'number') {                                                 // 752
            args = arraySlice(arguments);                                                                              // 753
            if (args.length < 2) {                                                                                     // 754
                pushCall(args, this.length - start);                                                                   // 755
            } else {                                                                                                   // 756
                args[1] = ES.ToInteger(deleteCount);                                                                   // 757
            }                                                                                                          // 758
        }                                                                                                              // 759
        return array_splice.apply(this, args);                                                                         // 760
    }                                                                                                                  // 761
}, !spliceWorksWithEmptyObject);                                                                                       // 762
var spliceWorksWithLargeSparseArrays = (function () {                                                                  // 763
    // Per https://github.com/es-shims/es5-shim/issues/295                                                             // 764
    // Safari 7/8 breaks with sparse arrays of size 1e5 or greater                                                     // 765
    var arr = new $Array(1e5);                                                                                         // 766
    // note: the index MUST be 8 or larger or the test will false pass                                                 // 767
    arr[8] = 'x';                                                                                                      // 768
    arr.splice(1, 1);                                                                                                  // 769
    // note: this test must be defined *after* the indexOf shim                                                        // 770
    // per https://github.com/es-shims/es5-shim/issues/313                                                             // 771
    return arr.indexOf('x') === 7;                                                                                     // 772
}());                                                                                                                  // 773
var spliceWorksWithSmallSparseArrays = (function () {                                                                  // 774
    // Per https://github.com/es-shims/es5-shim/issues/295                                                             // 775
    // Opera 12.15 breaks on this, no idea why.                                                                        // 776
    var n = 256;                                                                                                       // 777
    var arr = [];                                                                                                      // 778
    arr[n] = 'a';                                                                                                      // 779
    arr.splice(n + 1, 0, 'b');                                                                                         // 780
    return arr[n] === 'a';                                                                                             // 781
}());                                                                                                                  // 782
defineProperties(ArrayPrototype, {                                                                                     // 783
    splice: function splice(start, deleteCount) {                                                                      // 784
        var O = ES.ToObject(this);                                                                                     // 785
        var A = [];                                                                                                    // 786
        var len = ES.ToUint32(O.length);                                                                               // 787
        var relativeStart = ES.ToInteger(start);                                                                       // 788
        var actualStart = relativeStart < 0 ? max((len + relativeStart), 0) : min(relativeStart, len);                 // 789
        var actualDeleteCount = min(max(ES.ToInteger(deleteCount), 0), len - actualStart);                             // 790
                                                                                                                       // 791
        var k = 0;                                                                                                     // 792
        var from;                                                                                                      // 793
        while (k < actualDeleteCount) {                                                                                // 794
            from = $String(actualStart + k);                                                                           // 795
            if (owns(O, from)) {                                                                                       // 796
                A[k] = O[from];                                                                                        // 797
            }                                                                                                          // 798
            k += 1;                                                                                                    // 799
        }                                                                                                              // 800
                                                                                                                       // 801
        var items = arraySlice(arguments, 2);                                                                          // 802
        var itemCount = items.length;                                                                                  // 803
        var to;                                                                                                        // 804
        if (itemCount < actualDeleteCount) {                                                                           // 805
            k = actualStart;                                                                                           // 806
            var maxK = len - actualDeleteCount;                                                                        // 807
            while (k < maxK) {                                                                                         // 808
                from = $String(k + actualDeleteCount);                                                                 // 809
                to = $String(k + itemCount);                                                                           // 810
                if (owns(O, from)) {                                                                                   // 811
                    O[to] = O[from];                                                                                   // 812
                } else {                                                                                               // 813
                    delete O[to];                                                                                      // 814
                }                                                                                                      // 815
                k += 1;                                                                                                // 816
            }                                                                                                          // 817
            k = len;                                                                                                   // 818
            var minK = len - actualDeleteCount + itemCount;                                                            // 819
            while (k > minK) {                                                                                         // 820
                delete O[k - 1];                                                                                       // 821
                k -= 1;                                                                                                // 822
            }                                                                                                          // 823
        } else if (itemCount > actualDeleteCount) {                                                                    // 824
            k = len - actualDeleteCount;                                                                               // 825
            while (k > actualStart) {                                                                                  // 826
                from = $String(k + actualDeleteCount - 1);                                                             // 827
                to = $String(k + itemCount - 1);                                                                       // 828
                if (owns(O, from)) {                                                                                   // 829
                    O[to] = O[from];                                                                                   // 830
                } else {                                                                                               // 831
                    delete O[to];                                                                                      // 832
                }                                                                                                      // 833
                k -= 1;                                                                                                // 834
            }                                                                                                          // 835
        }                                                                                                              // 836
        k = actualStart;                                                                                               // 837
        for (var i = 0; i < items.length; ++i) {                                                                       // 838
            O[k] = items[i];                                                                                           // 839
            k += 1;                                                                                                    // 840
        }                                                                                                              // 841
        O.length = len - actualDeleteCount + itemCount;                                                                // 842
                                                                                                                       // 843
        return A;                                                                                                      // 844
    }                                                                                                                  // 845
}, !spliceWorksWithLargeSparseArrays || !spliceWorksWithSmallSparseArrays);                                            // 846
                                                                                                                       // 847
var originalJoin = ArrayPrototype.join;                                                                                // 848
var hasStringJoinBug;                                                                                                  // 849
try {                                                                                                                  // 850
    hasStringJoinBug = Array.prototype.join.call('123', ',') !== '1,2,3';                                              // 851
} catch (e) {                                                                                                          // 852
    hasStringJoinBug = true;                                                                                           // 853
}                                                                                                                      // 854
if (hasStringJoinBug) {                                                                                                // 855
    defineProperties(ArrayPrototype, {                                                                                 // 856
        join: function join(separator) {                                                                               // 857
            var sep = typeof separator === 'undefined' ? ',' : separator;                                              // 858
            return originalJoin.call(isString(this) ? strSplit(this, '') : this, sep);                                 // 859
        }                                                                                                              // 860
    }, hasStringJoinBug);                                                                                              // 861
}                                                                                                                      // 862
                                                                                                                       // 863
var hasJoinUndefinedBug = [1, 2].join(undefined) !== '1,2';                                                            // 864
if (hasJoinUndefinedBug) {                                                                                             // 865
    defineProperties(ArrayPrototype, {                                                                                 // 866
        join: function join(separator) {                                                                               // 867
            var sep = typeof separator === 'undefined' ? ',' : separator;                                              // 868
            return originalJoin.call(this, sep);                                                                       // 869
        }                                                                                                              // 870
    }, hasJoinUndefinedBug);                                                                                           // 871
}                                                                                                                      // 872
                                                                                                                       // 873
var pushShim = function push(item) {                                                                                   // 874
    var O = ES.ToObject(this);                                                                                         // 875
    var n = ES.ToUint32(O.length);                                                                                     // 876
    var i = 0;                                                                                                         // 877
    while (i < arguments.length) {                                                                                     // 878
        O[n + i] = arguments[i];                                                                                       // 879
        i += 1;                                                                                                        // 880
    }                                                                                                                  // 881
    O.length = n + i;                                                                                                  // 882
    return n + i;                                                                                                      // 883
};                                                                                                                     // 884
                                                                                                                       // 885
var pushIsNotGeneric = (function () {                                                                                  // 886
    var obj = {};                                                                                                      // 887
    var result = Array.prototype.push.call(obj, undefined);                                                            // 888
    return result !== 1 || obj.length !== 1 || typeof obj[0] !== 'undefined' || !owns(obj, 0);                         // 889
}());                                                                                                                  // 890
defineProperties(ArrayPrototype, {                                                                                     // 891
    push: function push(item) {                                                                                        // 892
        if (isArray(this)) {                                                                                           // 893
            return array_push.apply(this, arguments);                                                                  // 894
        }                                                                                                              // 895
        return pushShim.apply(this, arguments);                                                                        // 896
    }                                                                                                                  // 897
}, pushIsNotGeneric);                                                                                                  // 898
                                                                                                                       // 899
// This fixes a very weird bug in Opera 10.6 when pushing `undefined                                                   // 900
var pushUndefinedIsWeird = (function () {                                                                              // 901
    var arr = [];                                                                                                      // 902
    var result = arr.push(undefined);                                                                                  // 903
    return result !== 1 || arr.length !== 1 || typeof arr[0] !== 'undefined' || !owns(arr, 0);                         // 904
}());                                                                                                                  // 905
defineProperties(ArrayPrototype, { push: pushShim }, pushUndefinedIsWeird);                                            // 906
                                                                                                                       // 907
// ES5 15.2.3.14                                                                                                       // 908
// http://es5.github.io/#x15.4.4.10                                                                                    // 909
// Fix boxed string bug                                                                                                // 910
defineProperties(ArrayPrototype, {                                                                                     // 911
    slice: function (start, end) {                                                                                     // 912
        var arr = isString(this) ? strSplit(this, '') : this;                                                          // 913
        return arraySliceApply(arr, arguments);                                                                        // 914
    }                                                                                                                  // 915
}, splitString);                                                                                                       // 916
                                                                                                                       // 917
var sortIgnoresNonFunctions = (function () {                                                                           // 918
    try {                                                                                                              // 919
        [1, 2].sort(null);                                                                                             // 920
        [1, 2].sort({});                                                                                               // 921
        return true;                                                                                                   // 922
    } catch (e) { /**/ }                                                                                               // 923
    return false;                                                                                                      // 924
}());                                                                                                                  // 925
var sortThrowsOnRegex = (function () {                                                                                 // 926
    // this is a problem in Firefox 4, in which `typeof /a/ === 'function'`                                            // 927
    try {                                                                                                              // 928
        [1, 2].sort(/a/);                                                                                              // 929
        return false;                                                                                                  // 930
    } catch (e) { /**/ }                                                                                               // 931
    return true;                                                                                                       // 932
}());                                                                                                                  // 933
var sortIgnoresUndefined = (function () {                                                                              // 934
    // applies in IE 8, for one.                                                                                       // 935
    try {                                                                                                              // 936
        [1, 2].sort(undefined);                                                                                        // 937
        return true;                                                                                                   // 938
    } catch (e) { /**/ }                                                                                               // 939
    return false;                                                                                                      // 940
}());                                                                                                                  // 941
defineProperties(ArrayPrototype, {                                                                                     // 942
    sort: function sort(compareFn) {                                                                                   // 943
        if (typeof compareFn === 'undefined') {                                                                        // 944
            return arraySort(this);                                                                                    // 945
        }                                                                                                              // 946
        if (!isCallable(compareFn)) {                                                                                  // 947
            throw new TypeError('Array.prototype.sort callback must be a function');                                   // 948
        }                                                                                                              // 949
        return arraySort(this, compareFn);                                                                             // 950
    }                                                                                                                  // 951
}, sortIgnoresNonFunctions || !sortIgnoresUndefined || !sortThrowsOnRegex);                                            // 952
                                                                                                                       // 953
//                                                                                                                     // 954
// Object                                                                                                              // 955
// ======                                                                                                              // 956
//                                                                                                                     // 957
                                                                                                                       // 958
// ES5 15.2.3.14                                                                                                       // 959
// http://es5.github.com/#x15.2.3.14                                                                                   // 960
                                                                                                                       // 961
// http://whattheheadsaid.com/2010/10/a-safer-object-keys-compatibility-implementation                                 // 962
var hasDontEnumBug = !({ 'toString': null }).propertyIsEnumerable('toString');                                         // 963
var hasProtoEnumBug = function () {}.propertyIsEnumerable('prototype');                                                // 964
var hasStringEnumBug = !owns('x', '0');                                                                                // 965
var equalsConstructorPrototype = function (o) {                                                                        // 966
    var ctor = o.constructor;                                                                                          // 967
    return ctor && ctor.prototype === o;                                                                               // 968
};                                                                                                                     // 969
var blacklistedKeys = {                                                                                                // 970
    $window: true,                                                                                                     // 971
    $console: true,                                                                                                    // 972
    $parent: true,                                                                                                     // 973
    $self: true,                                                                                                       // 974
    $frame: true,                                                                                                      // 975
    $frames: true,                                                                                                     // 976
    $frameElement: true,                                                                                               // 977
    $webkitIndexedDB: true,                                                                                            // 978
    $webkitStorageInfo: true,                                                                                          // 979
    $external: true                                                                                                    // 980
};                                                                                                                     // 981
var hasAutomationEqualityBug = (function () {                                                                          // 982
    /* globals window */                                                                                               // 983
    if (typeof window === 'undefined') { return false; }                                                               // 984
    for (var k in window) {                                                                                            // 985
        try {                                                                                                          // 986
            if (!blacklistedKeys['$' + k] && owns(window, k) && window[k] !== null && typeof window[k] === 'object') {
                equalsConstructorPrototype(window[k]);                                                                 // 988
            }                                                                                                          // 989
        } catch (e) {                                                                                                  // 990
            return true;                                                                                               // 991
        }                                                                                                              // 992
    }                                                                                                                  // 993
    return false;                                                                                                      // 994
}());                                                                                                                  // 995
var equalsConstructorPrototypeIfNotBuggy = function (object) {                                                         // 996
    if (typeof window === 'undefined' || !hasAutomationEqualityBug) { return equalsConstructorPrototype(object); }     // 997
    try {                                                                                                              // 998
        return equalsConstructorPrototype(object);                                                                     // 999
    } catch (e) {                                                                                                      // 1000
        return false;                                                                                                  // 1001
    }                                                                                                                  // 1002
};                                                                                                                     // 1003
var dontEnums = [                                                                                                      // 1004
    'toString',                                                                                                        // 1005
    'toLocaleString',                                                                                                  // 1006
    'valueOf',                                                                                                         // 1007
    'hasOwnProperty',                                                                                                  // 1008
    'isPrototypeOf',                                                                                                   // 1009
    'propertyIsEnumerable',                                                                                            // 1010
    'constructor'                                                                                                      // 1011
];                                                                                                                     // 1012
var dontEnumsLength = dontEnums.length;                                                                                // 1013
                                                                                                                       // 1014
// taken directly from https://github.com/ljharb/is-arguments/blob/master/index.js                                     // 1015
// can be replaced with require('is-arguments') if we ever use a build process instead                                 // 1016
var isStandardArguments = function isArguments(value) {                                                                // 1017
    return toStr(value) === '[object Arguments]';                                                                      // 1018
};                                                                                                                     // 1019
var isLegacyArguments = function isArguments(value) {                                                                  // 1020
    return value !== null &&                                                                                           // 1021
        typeof value === 'object' &&                                                                                   // 1022
        typeof value.length === 'number' &&                                                                            // 1023
        value.length >= 0 &&                                                                                           // 1024
        !isArray(value) &&                                                                                             // 1025
        isCallable(value.callee);                                                                                      // 1026
};                                                                                                                     // 1027
var isArguments = isStandardArguments(arguments) ? isStandardArguments : isLegacyArguments;                            // 1028
                                                                                                                       // 1029
defineProperties($Object, {                                                                                            // 1030
    keys: function keys(object) {                                                                                      // 1031
        var isFn = isCallable(object);                                                                                 // 1032
        var isArgs = isArguments(object);                                                                              // 1033
        var isObject = object !== null && typeof object === 'object';                                                  // 1034
        var isStr = isObject && isString(object);                                                                      // 1035
                                                                                                                       // 1036
        if (!isObject && !isFn && !isArgs) {                                                                           // 1037
            throw new TypeError('Object.keys called on a non-object');                                                 // 1038
        }                                                                                                              // 1039
                                                                                                                       // 1040
        var theKeys = [];                                                                                              // 1041
        var skipProto = hasProtoEnumBug && isFn;                                                                       // 1042
        if ((isStr && hasStringEnumBug) || isArgs) {                                                                   // 1043
            for (var i = 0; i < object.length; ++i) {                                                                  // 1044
                pushCall(theKeys, $String(i));                                                                         // 1045
            }                                                                                                          // 1046
        }                                                                                                              // 1047
                                                                                                                       // 1048
        if (!isArgs) {                                                                                                 // 1049
            for (var name in object) {                                                                                 // 1050
                if (!(skipProto && name === 'prototype') && owns(object, name)) {                                      // 1051
                    pushCall(theKeys, $String(name));                                                                  // 1052
                }                                                                                                      // 1053
            }                                                                                                          // 1054
        }                                                                                                              // 1055
                                                                                                                       // 1056
        if (hasDontEnumBug) {                                                                                          // 1057
            var skipConstructor = equalsConstructorPrototypeIfNotBuggy(object);                                        // 1058
            for (var j = 0; j < dontEnumsLength; j++) {                                                                // 1059
                var dontEnum = dontEnums[j];                                                                           // 1060
                if (!(skipConstructor && dontEnum === 'constructor') && owns(object, dontEnum)) {                      // 1061
                    pushCall(theKeys, dontEnum);                                                                       // 1062
                }                                                                                                      // 1063
            }                                                                                                          // 1064
        }                                                                                                              // 1065
        return theKeys;                                                                                                // 1066
    }                                                                                                                  // 1067
});                                                                                                                    // 1068
                                                                                                                       // 1069
var keysWorksWithArguments = $Object.keys && (function () {                                                            // 1070
    // Safari 5.0 bug                                                                                                  // 1071
    return $Object.keys(arguments).length === 2;                                                                       // 1072
}(1, 2));                                                                                                              // 1073
var keysHasArgumentsLengthBug = $Object.keys && (function () {                                                         // 1074
    var argKeys = $Object.keys(arguments);                                                                             // 1075
    return arguments.length !== 1 || argKeys.length !== 1 || argKeys[0] !== 1;                                         // 1076
}(1));                                                                                                                 // 1077
var originalKeys = $Object.keys;                                                                                       // 1078
defineProperties($Object, {                                                                                            // 1079
    keys: function keys(object) {                                                                                      // 1080
        if (isArguments(object)) {                                                                                     // 1081
            return originalKeys(arraySlice(object));                                                                   // 1082
        } else {                                                                                                       // 1083
            return originalKeys(object);                                                                               // 1084
        }                                                                                                              // 1085
    }                                                                                                                  // 1086
}, !keysWorksWithArguments || keysHasArgumentsLengthBug);                                                              // 1087
                                                                                                                       // 1088
//                                                                                                                     // 1089
// Date                                                                                                                // 1090
// ====                                                                                                                // 1091
//                                                                                                                     // 1092
                                                                                                                       // 1093
var hasNegativeMonthYearBug = new Date(-3509827329600292).getUTCMonth() !== 0;                                         // 1094
var aNegativeTestDate = new Date(-1509842289600292);                                                                   // 1095
var aPositiveTestDate = new Date(1449662400000);                                                                       // 1096
var hasToUTCStringFormatBug = aNegativeTestDate.toUTCString() !== 'Mon, 01 Jan -45875 11:59:59 GMT';                   // 1097
var hasToDateStringFormatBug;                                                                                          // 1098
var hasToStringFormatBug;                                                                                              // 1099
var timeZoneOffset = aNegativeTestDate.getTimezoneOffset();                                                            // 1100
if (timeZoneOffset < -720) {                                                                                           // 1101
    hasToDateStringFormatBug = aNegativeTestDate.toDateString() !== 'Tue Jan 02 -45875';                               // 1102
    hasToStringFormatBug = !(/^Thu Dec 10 2015 \d\d:\d\d:\d\d GMT[-\+]\d\d\d\d(?: |$)/).test(aPositiveTestDate.toString());
} else {                                                                                                               // 1104
    hasToDateStringFormatBug = aNegativeTestDate.toDateString() !== 'Mon Jan 01 -45875';                               // 1105
    hasToStringFormatBug = !(/^Wed Dec 09 2015 \d\d:\d\d:\d\d GMT[-\+]\d\d\d\d(?: |$)/).test(aPositiveTestDate.toString());
}                                                                                                                      // 1107
                                                                                                                       // 1108
var originalGetFullYear = call.bind(Date.prototype.getFullYear);                                                       // 1109
var originalGetMonth = call.bind(Date.prototype.getMonth);                                                             // 1110
var originalGetDate = call.bind(Date.prototype.getDate);                                                               // 1111
var originalGetUTCFullYear = call.bind(Date.prototype.getUTCFullYear);                                                 // 1112
var originalGetUTCMonth = call.bind(Date.prototype.getUTCMonth);                                                       // 1113
var originalGetUTCDate = call.bind(Date.prototype.getUTCDate);                                                         // 1114
var originalGetUTCDay = call.bind(Date.prototype.getUTCDay);                                                           // 1115
var originalGetUTCHours = call.bind(Date.prototype.getUTCHours);                                                       // 1116
var originalGetUTCMinutes = call.bind(Date.prototype.getUTCMinutes);                                                   // 1117
var originalGetUTCSeconds = call.bind(Date.prototype.getUTCSeconds);                                                   // 1118
var originalGetUTCMilliseconds = call.bind(Date.prototype.getUTCMilliseconds);                                         // 1119
var dayName = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];                                                       // 1120
var monthName = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];                  // 1121
var daysInMonth = function daysInMonth(month, year) {                                                                  // 1122
    return originalGetDate(new Date(year, month, 0));                                                                  // 1123
};                                                                                                                     // 1124
                                                                                                                       // 1125
defineProperties(Date.prototype, {                                                                                     // 1126
    getFullYear: function getFullYear() {                                                                              // 1127
        if (!this || !(this instanceof Date)) {                                                                        // 1128
            throw new TypeError('this is not a Date object.');                                                         // 1129
        }                                                                                                              // 1130
        var year = originalGetFullYear(this);                                                                          // 1131
        if (year < 0 && originalGetMonth(this) > 11) {                                                                 // 1132
            return year + 1;                                                                                           // 1133
        }                                                                                                              // 1134
        return year;                                                                                                   // 1135
    },                                                                                                                 // 1136
    getMonth: function getMonth() {                                                                                    // 1137
        if (!this || !(this instanceof Date)) {                                                                        // 1138
            throw new TypeError('this is not a Date object.');                                                         // 1139
        }                                                                                                              // 1140
        var year = originalGetFullYear(this);                                                                          // 1141
        var month = originalGetMonth(this);                                                                            // 1142
        if (year < 0 && month > 11) {                                                                                  // 1143
            return 0;                                                                                                  // 1144
        }                                                                                                              // 1145
        return month;                                                                                                  // 1146
    },                                                                                                                 // 1147
    getDate: function getDate() {                                                                                      // 1148
        if (!this || !(this instanceof Date)) {                                                                        // 1149
            throw new TypeError('this is not a Date object.');                                                         // 1150
        }                                                                                                              // 1151
        var year = originalGetFullYear(this);                                                                          // 1152
        var month = originalGetMonth(this);                                                                            // 1153
        var date = originalGetDate(this);                                                                              // 1154
        if (year < 0 && month > 11) {                                                                                  // 1155
            if (month === 12) {                                                                                        // 1156
                return date;                                                                                           // 1157
            }                                                                                                          // 1158
            var days = daysInMonth(0, year + 1);                                                                       // 1159
            return (days - date) + 1;                                                                                  // 1160
        }                                                                                                              // 1161
        return date;                                                                                                   // 1162
    },                                                                                                                 // 1163
    getUTCFullYear: function getUTCFullYear() {                                                                        // 1164
        if (!this || !(this instanceof Date)) {                                                                        // 1165
            throw new TypeError('this is not a Date object.');                                                         // 1166
        }                                                                                                              // 1167
        var year = originalGetUTCFullYear(this);                                                                       // 1168
        if (year < 0 && originalGetUTCMonth(this) > 11) {                                                              // 1169
            return year + 1;                                                                                           // 1170
        }                                                                                                              // 1171
        return year;                                                                                                   // 1172
    },                                                                                                                 // 1173
    getUTCMonth: function getUTCMonth() {                                                                              // 1174
        if (!this || !(this instanceof Date)) {                                                                        // 1175
            throw new TypeError('this is not a Date object.');                                                         // 1176
        }                                                                                                              // 1177
        var year = originalGetUTCFullYear(this);                                                                       // 1178
        var month = originalGetUTCMonth(this);                                                                         // 1179
        if (year < 0 && month > 11) {                                                                                  // 1180
            return 0;                                                                                                  // 1181
        }                                                                                                              // 1182
        return month;                                                                                                  // 1183
    },                                                                                                                 // 1184
    getUTCDate: function getUTCDate() {                                                                                // 1185
        if (!this || !(this instanceof Date)) {                                                                        // 1186
            throw new TypeError('this is not a Date object.');                                                         // 1187
        }                                                                                                              // 1188
        var year = originalGetUTCFullYear(this);                                                                       // 1189
        var month = originalGetUTCMonth(this);                                                                         // 1190
        var date = originalGetUTCDate(this);                                                                           // 1191
        if (year < 0 && month > 11) {                                                                                  // 1192
            if (month === 12) {                                                                                        // 1193
                return date;                                                                                           // 1194
            }                                                                                                          // 1195
            var days = daysInMonth(0, year + 1);                                                                       // 1196
            return (days - date) + 1;                                                                                  // 1197
        }                                                                                                              // 1198
        return date;                                                                                                   // 1199
    }                                                                                                                  // 1200
}, hasNegativeMonthYearBug);                                                                                           // 1201
                                                                                                                       // 1202
defineProperties(Date.prototype, {                                                                                     // 1203
    toUTCString: function toUTCString() {                                                                              // 1204
        if (!this || !(this instanceof Date)) {                                                                        // 1205
            throw new TypeError('this is not a Date object.');                                                         // 1206
        }                                                                                                              // 1207
        var day = originalGetUTCDay(this);                                                                             // 1208
        var date = originalGetUTCDate(this);                                                                           // 1209
        var month = originalGetUTCMonth(this);                                                                         // 1210
        var year = originalGetUTCFullYear(this);                                                                       // 1211
        var hour = originalGetUTCHours(this);                                                                          // 1212
        var minute = originalGetUTCMinutes(this);                                                                      // 1213
        var second = originalGetUTCSeconds(this);                                                                      // 1214
        return dayName[day] + ', ' +                                                                                   // 1215
            (date < 10 ? '0' + date : date) + ' ' +                                                                    // 1216
            monthName[month] + ' ' +                                                                                   // 1217
            year + ' ' +                                                                                               // 1218
            (hour < 10 ? '0' + hour : hour) + ':' +                                                                    // 1219
            (minute < 10 ? '0' + minute : minute) + ':' +                                                              // 1220
            (second < 10 ? '0' + second : second) + ' GMT';                                                            // 1221
    }                                                                                                                  // 1222
}, hasNegativeMonthYearBug || hasToUTCStringFormatBug);                                                                // 1223
                                                                                                                       // 1224
// Opera 12 has `,`                                                                                                    // 1225
defineProperties(Date.prototype, {                                                                                     // 1226
    toDateString: function toDateString() {                                                                            // 1227
        if (!this || !(this instanceof Date)) {                                                                        // 1228
            throw new TypeError('this is not a Date object.');                                                         // 1229
        }                                                                                                              // 1230
        var day = this.getDay();                                                                                       // 1231
        var date = this.getDate();                                                                                     // 1232
        var month = this.getMonth();                                                                                   // 1233
        var year = this.getFullYear();                                                                                 // 1234
        return dayName[day] + ' ' +                                                                                    // 1235
            monthName[month] + ' ' +                                                                                   // 1236
            (date < 10 ? '0' + date : date) + ' ' +                                                                    // 1237
            year;                                                                                                      // 1238
    }                                                                                                                  // 1239
}, hasNegativeMonthYearBug || hasToDateStringFormatBug);                                                               // 1240
                                                                                                                       // 1241
// can't use defineProperties here because of toString enumeration issue in IE <= 8                                    // 1242
if (hasNegativeMonthYearBug || hasToStringFormatBug) {                                                                 // 1243
    Date.prototype.toString = function toString() {                                                                    // 1244
        if (!this || !(this instanceof Date)) {                                                                        // 1245
            throw new TypeError('this is not a Date object.');                                                         // 1246
        }                                                                                                              // 1247
        var day = this.getDay();                                                                                       // 1248
        var date = this.getDate();                                                                                     // 1249
        var month = this.getMonth();                                                                                   // 1250
        var year = this.getFullYear();                                                                                 // 1251
        var hour = this.getHours();                                                                                    // 1252
        var minute = this.getMinutes();                                                                                // 1253
        var second = this.getSeconds();                                                                                // 1254
        var timezoneOffset = this.getTimezoneOffset();                                                                 // 1255
        var hoursOffset = Math.floor(Math.abs(timezoneOffset) / 60);                                                   // 1256
        var minutesOffset = Math.floor(Math.abs(timezoneOffset) % 60);                                                 // 1257
        return dayName[day] + ' ' +                                                                                    // 1258
            monthName[month] + ' ' +                                                                                   // 1259
            (date < 10 ? '0' + date : date) + ' ' +                                                                    // 1260
            year + ' ' +                                                                                               // 1261
            (hour < 10 ? '0' + hour : hour) + ':' +                                                                    // 1262
            (minute < 10 ? '0' + minute : minute) + ':' +                                                              // 1263
            (second < 10 ? '0' + second : second) + ' GMT' +                                                           // 1264
            (timezoneOffset > 0 ? '-' : '+') +                                                                         // 1265
            (hoursOffset < 10 ? '0' + hoursOffset : hoursOffset) +                                                     // 1266
            (minutesOffset < 10 ? '0' + minutesOffset : minutesOffset);                                                // 1267
    };                                                                                                                 // 1268
    if (supportsDescriptors) {                                                                                         // 1269
        $Object.defineProperty(Date.prototype, 'toString', {                                                           // 1270
            configurable: true,                                                                                        // 1271
            enumerable: false,                                                                                         // 1272
            writable: true                                                                                             // 1273
        });                                                                                                            // 1274
    }                                                                                                                  // 1275
}                                                                                                                      // 1276
                                                                                                                       // 1277
// ES5 15.9.5.43                                                                                                       // 1278
// http://es5.github.com/#x15.9.5.43                                                                                   // 1279
// This function returns a String value represent the instance in time                                                 // 1280
// represented by this Date object. The format of the String is the Date Time                                          // 1281
// string format defined in 15.9.1.15. All fields are present in the String.                                           // 1282
// The time zone is always UTC, denoted by the suffix Z. If the time value of                                          // 1283
// this object is not a finite Number a RangeError exception is thrown.                                                // 1284
var negativeDate = -62198755200000;                                                                                    // 1285
var negativeYearString = '-000001';                                                                                    // 1286
var hasNegativeDateBug = Date.prototype.toISOString && new Date(negativeDate).toISOString().indexOf(negativeYearString) === -1;
var hasSafari51DateBug = Date.prototype.toISOString && new Date(-1).toISOString() !== '1969-12-31T23:59:59.999Z';      // 1288
                                                                                                                       // 1289
var getTime = call.bind(Date.prototype.getTime);                                                                       // 1290
                                                                                                                       // 1291
defineProperties(Date.prototype, {                                                                                     // 1292
    toISOString: function toISOString() {                                                                              // 1293
        if (!isFinite(this) || !isFinite(getTime(this))) {                                                             // 1294
            // Adope Photoshop requires the second check.                                                              // 1295
            throw new RangeError('Date.prototype.toISOString called on non-finite value.');                            // 1296
        }                                                                                                              // 1297
                                                                                                                       // 1298
        var year = originalGetUTCFullYear(this);                                                                       // 1299
                                                                                                                       // 1300
        var month = originalGetUTCMonth(this);                                                                         // 1301
        // see https://github.com/es-shims/es5-shim/issues/111                                                         // 1302
        year += Math.floor(month / 12);                                                                                // 1303
        month = (month % 12 + 12) % 12;                                                                                // 1304
                                                                                                                       // 1305
        // the date time string format is specified in 15.9.1.15.                                                      // 1306
        var result = [month + 1, originalGetUTCDate(this), originalGetUTCHours(this), originalGetUTCMinutes(this), originalGetUTCSeconds(this)];
        year = (                                                                                                       // 1308
            (year < 0 ? '-' : (year > 9999 ? '+' : '')) +                                                              // 1309
            strSlice('00000' + Math.abs(year), (0 <= year && year <= 9999) ? -4 : -6)                                  // 1310
        );                                                                                                             // 1311
                                                                                                                       // 1312
        for (var i = 0; i < result.length; ++i) {                                                                      // 1313
          // pad months, days, hours, minutes, and seconds to have two digits.                                         // 1314
          result[i] = strSlice('00' + result[i], -2);                                                                  // 1315
        }                                                                                                              // 1316
        // pad milliseconds to have three digits.                                                                      // 1317
        return (                                                                                                       // 1318
            year + '-' + arraySlice(result, 0, 2).join('-') +                                                          // 1319
            'T' + arraySlice(result, 2).join(':') + '.' +                                                              // 1320
            strSlice('000' + originalGetUTCMilliseconds(this), -3) + 'Z'                                               // 1321
        );                                                                                                             // 1322
    }                                                                                                                  // 1323
}, hasNegativeDateBug || hasSafari51DateBug);                                                                          // 1324
                                                                                                                       // 1325
// ES5 15.9.5.44                                                                                                       // 1326
// http://es5.github.com/#x15.9.5.44                                                                                   // 1327
// This function provides a String representation of a Date object for use by                                          // 1328
// JSON.stringify (15.12.3).                                                                                           // 1329
var dateToJSONIsSupported = (function () {                                                                             // 1330
    try {                                                                                                              // 1331
        return Date.prototype.toJSON &&                                                                                // 1332
            new Date(NaN).toJSON() === null &&                                                                         // 1333
            new Date(negativeDate).toJSON().indexOf(negativeYearString) !== -1 &&                                      // 1334
            Date.prototype.toJSON.call({ // generic                                                                    // 1335
                toISOString: function () { return true; }                                                              // 1336
            });                                                                                                        // 1337
    } catch (e) {                                                                                                      // 1338
        return false;                                                                                                  // 1339
    }                                                                                                                  // 1340
}());                                                                                                                  // 1341
if (!dateToJSONIsSupported) {                                                                                          // 1342
    Date.prototype.toJSON = function toJSON(key) {                                                                     // 1343
        // When the toJSON method is called with argument key, the following                                           // 1344
        // steps are taken:                                                                                            // 1345
                                                                                                                       // 1346
        // 1.  Let O be the result of calling ToObject, giving it the this                                             // 1347
        // value as its argument.                                                                                      // 1348
        // 2. Let tv be ES.ToPrimitive(O, hint Number).                                                                // 1349
        var O = $Object(this);                                                                                         // 1350
        var tv = ES.ToPrimitive(O);                                                                                    // 1351
        // 3. If tv is a Number and is not finite, return null.                                                        // 1352
        if (typeof tv === 'number' && !isFinite(tv)) {                                                                 // 1353
            return null;                                                                                               // 1354
        }                                                                                                              // 1355
        // 4. Let toISO be the result of calling the [[Get]] internal method of                                        // 1356
        // O with argument "toISOString".                                                                              // 1357
        var toISO = O.toISOString;                                                                                     // 1358
        // 5. If IsCallable(toISO) is false, throw a TypeError exception.                                              // 1359
        if (!isCallable(toISO)) {                                                                                      // 1360
            throw new TypeError('toISOString property is not callable');                                               // 1361
        }                                                                                                              // 1362
        // 6. Return the result of calling the [[Call]] internal method of                                             // 1363
        //  toISO with O as the this value and an empty argument list.                                                 // 1364
        return toISO.call(O);                                                                                          // 1365
                                                                                                                       // 1366
        // NOTE 1 The argument is ignored.                                                                             // 1367
                                                                                                                       // 1368
        // NOTE 2 The toJSON function is intentionally generic; it does not                                            // 1369
        // require that its this value be a Date object. Therefore, it can be                                          // 1370
        // transferred to other kinds of objects for use as a method. However,                                         // 1371
        // it does require that any such object have a toISOString method. An                                          // 1372
        // object is free to use the argument key to filter its                                                        // 1373
        // stringification.                                                                                            // 1374
    };                                                                                                                 // 1375
}                                                                                                                      // 1376
                                                                                                                       // 1377
// ES5 15.9.4.2                                                                                                        // 1378
// http://es5.github.com/#x15.9.4.2                                                                                    // 1379
// based on work shared by Daniel Friesen (dantman)                                                                    // 1380
// http://gist.github.com/303249                                                                                       // 1381
var supportsExtendedYears = Date.parse('+033658-09-27T01:46:40.000Z') === 1e15;                                        // 1382
var acceptsInvalidDates = !isNaN(Date.parse('2012-04-04T24:00:00.500Z')) || !isNaN(Date.parse('2012-11-31T23:59:59.000Z')) || !isNaN(Date.parse('2012-12-31T23:59:60.000Z'));
var doesNotParseY2KNewYear = isNaN(Date.parse('2000-01-01T00:00:00.000Z'));                                            // 1384
if (doesNotParseY2KNewYear || acceptsInvalidDates || !supportsExtendedYears) {                                         // 1385
    // XXX global assignment won't work in embeddings that use                                                         // 1386
    // an alternate object for the context.                                                                            // 1387
    /* global Date: true */                                                                                            // 1388
    /* eslint-disable no-undef */                                                                                      // 1389
    var maxSafeUnsigned32Bit = Math.pow(2, 31) - 1;                                                                    // 1390
    var hasSafariSignedIntBug = isActualNaN(new Date(1970, 0, 1, 0, 0, 0, maxSafeUnsigned32Bit + 1).getTime());        // 1391
    /* eslint-disable no-implicit-globals */                                                                           // 1392
    Date = (function (NativeDate) {                                                                                    // 1393
    /* eslint-enable no-implicit-globals */                                                                            // 1394
    /* eslint-enable no-undef */                                                                                       // 1395
        // Date.length === 7                                                                                           // 1396
        var DateShim = function Date(Y, M, D, h, m, s, ms) {                                                           // 1397
            var length = arguments.length;                                                                             // 1398
            var date;                                                                                                  // 1399
            if (this instanceof NativeDate) {                                                                          // 1400
                var seconds = s;                                                                                       // 1401
                var millis = ms;                                                                                       // 1402
                if (hasSafariSignedIntBug && length >= 7 && ms > maxSafeUnsigned32Bit) {                               // 1403
                    // work around a Safari 8/9 bug where it treats the seconds as signed                              // 1404
                    var msToShift = Math.floor(ms / maxSafeUnsigned32Bit) * maxSafeUnsigned32Bit;                      // 1405
                    var sToShift = Math.floor(msToShift / 1e3);                                                        // 1406
                    seconds += sToShift;                                                                               // 1407
                    millis -= sToShift * 1e3;                                                                          // 1408
                }                                                                                                      // 1409
                date = length === 1 && $String(Y) === Y ? // isString(Y)                                               // 1410
                    // We explicitly pass it through parse:                                                            // 1411
                    new NativeDate(DateShim.parse(Y)) :                                                                // 1412
                    // We have to manually make calls depending on argument                                            // 1413
                    // length here                                                                                     // 1414
                    length >= 7 ? new NativeDate(Y, M, D, h, m, seconds, millis) :                                     // 1415
                    length >= 6 ? new NativeDate(Y, M, D, h, m, seconds) :                                             // 1416
                    length >= 5 ? new NativeDate(Y, M, D, h, m) :                                                      // 1417
                    length >= 4 ? new NativeDate(Y, M, D, h) :                                                         // 1418
                    length >= 3 ? new NativeDate(Y, M, D) :                                                            // 1419
                    length >= 2 ? new NativeDate(Y, M) :                                                               // 1420
                    length >= 1 ? new NativeDate(Y instanceof NativeDate ? +Y : Y) :                                   // 1421
                                  new NativeDate();                                                                    // 1422
            } else {                                                                                                   // 1423
                date = NativeDate.apply(this, arguments);                                                              // 1424
            }                                                                                                          // 1425
            if (!isPrimitive(date)) {                                                                                  // 1426
              // Prevent mixups with unfixed Date object                                                               // 1427
              defineProperties(date, { constructor: DateShim }, true);                                                 // 1428
            }                                                                                                          // 1429
            return date;                                                                                               // 1430
        };                                                                                                             // 1431
                                                                                                                       // 1432
        // 15.9.1.15 Date Time String Format.                                                                          // 1433
        var isoDateExpression = new RegExp('^' +                                                                       // 1434
            '(\\d{4}|[+-]\\d{6})' + // four-digit year capture or sign +                                               // 1435
                                      // 6-digit extended year                                                         // 1436
            '(?:-(\\d{2})' + // optional month capture                                                                 // 1437
            '(?:-(\\d{2})' + // optional day capture                                                                   // 1438
            '(?:' + // capture hours:minutes:seconds.milliseconds                                                      // 1439
                'T(\\d{2})' + // hours capture                                                                         // 1440
                ':(\\d{2})' + // minutes capture                                                                       // 1441
                '(?:' + // optional :seconds.milliseconds                                                              // 1442
                    ':(\\d{2})' + // seconds capture                                                                   // 1443
                    '(?:(\\.\\d{1,}))?' + // milliseconds capture                                                      // 1444
                ')?' +                                                                                                 // 1445
            '(' + // capture UTC offset component                                                                      // 1446
                'Z|' + // UTC capture                                                                                  // 1447
                '(?:' + // offset specifier +/-hours:minutes                                                           // 1448
                    '([-+])' + // sign capture                                                                         // 1449
                    '(\\d{2})' + // hours offset capture                                                               // 1450
                    ':(\\d{2})' + // minutes offset capture                                                            // 1451
                ')' +                                                                                                  // 1452
            ')?)?)?)?' +                                                                                               // 1453
        '$');                                                                                                          // 1454
                                                                                                                       // 1455
        var months = [0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365];                                     // 1456
                                                                                                                       // 1457
        var dayFromMonth = function dayFromMonth(year, month) {                                                        // 1458
            var t = month > 1 ? 1 : 0;                                                                                 // 1459
            return (                                                                                                   // 1460
                months[month] +                                                                                        // 1461
                Math.floor((year - 1969 + t) / 4) -                                                                    // 1462
                Math.floor((year - 1901 + t) / 100) +                                                                  // 1463
                Math.floor((year - 1601 + t) / 400) +                                                                  // 1464
                365 * (year - 1970)                                                                                    // 1465
            );                                                                                                         // 1466
        };                                                                                                             // 1467
                                                                                                                       // 1468
        var toUTC = function toUTC(t) {                                                                                // 1469
            var s = 0;                                                                                                 // 1470
            var ms = t;                                                                                                // 1471
            if (hasSafariSignedIntBug && ms > maxSafeUnsigned32Bit) {                                                  // 1472
                // work around a Safari 8/9 bug where it treats the seconds as signed                                  // 1473
                var msToShift = Math.floor(ms / maxSafeUnsigned32Bit) * maxSafeUnsigned32Bit;                          // 1474
                var sToShift = Math.floor(msToShift / 1e3);                                                            // 1475
                s += sToShift;                                                                                         // 1476
                ms -= sToShift * 1e3;                                                                                  // 1477
            }                                                                                                          // 1478
            return $Number(new NativeDate(1970, 0, 1, 0, 0, s, ms));                                                   // 1479
        };                                                                                                             // 1480
                                                                                                                       // 1481
        // Copy any custom methods a 3rd party library may have added                                                  // 1482
        for (var key in NativeDate) {                                                                                  // 1483
            if (owns(NativeDate, key)) {                                                                               // 1484
                DateShim[key] = NativeDate[key];                                                                       // 1485
            }                                                                                                          // 1486
        }                                                                                                              // 1487
                                                                                                                       // 1488
        // Copy "native" methods explicitly; they may be non-enumerable                                                // 1489
        defineProperties(DateShim, {                                                                                   // 1490
            now: NativeDate.now,                                                                                       // 1491
            UTC: NativeDate.UTC                                                                                        // 1492
        }, true);                                                                                                      // 1493
        DateShim.prototype = NativeDate.prototype;                                                                     // 1494
        defineProperties(DateShim.prototype, {                                                                         // 1495
            constructor: DateShim                                                                                      // 1496
        }, true);                                                                                                      // 1497
                                                                                                                       // 1498
        // Upgrade Date.parse to handle simplified ISO 8601 strings                                                    // 1499
        var parseShim = function parse(string) {                                                                       // 1500
            var match = isoDateExpression.exec(string);                                                                // 1501
            if (match) {                                                                                               // 1502
                // parse months, days, hours, minutes, seconds, and milliseconds                                       // 1503
                // provide default values if necessary                                                                 // 1504
                // parse the UTC offset component                                                                      // 1505
                var year = $Number(match[1]),                                                                          // 1506
                    month = $Number(match[2] || 1) - 1,                                                                // 1507
                    day = $Number(match[3] || 1) - 1,                                                                  // 1508
                    hour = $Number(match[4] || 0),                                                                     // 1509
                    minute = $Number(match[5] || 0),                                                                   // 1510
                    second = $Number(match[6] || 0),                                                                   // 1511
                    millisecond = Math.floor($Number(match[7] || 0) * 1000),                                           // 1512
                    // When time zone is missed, local offset should be used                                           // 1513
                    // (ES 5.1 bug)                                                                                    // 1514
                    // see https://bugs.ecmascript.org/show_bug.cgi?id=112                                             // 1515
                    isLocalTime = Boolean(match[4] && !match[8]),                                                      // 1516
                    signOffset = match[9] === '-' ? 1 : -1,                                                            // 1517
                    hourOffset = $Number(match[10] || 0),                                                              // 1518
                    minuteOffset = $Number(match[11] || 0),                                                            // 1519
                    result;                                                                                            // 1520
                var hasMinutesOrSecondsOrMilliseconds = minute > 0 || second > 0 || millisecond > 0;                   // 1521
                if (                                                                                                   // 1522
                    hour < (hasMinutesOrSecondsOrMilliseconds ? 24 : 25) &&                                            // 1523
                    minute < 60 && second < 60 && millisecond < 1000 &&                                                // 1524
                    month > -1 && month < 12 && hourOffset < 24 &&                                                     // 1525
                    minuteOffset < 60 && // detect invalid offsets                                                     // 1526
                    day > -1 &&                                                                                        // 1527
                    day < (dayFromMonth(year, month + 1) - dayFromMonth(year, month))                                  // 1528
                ) {                                                                                                    // 1529
                    result = (                                                                                         // 1530
                        (dayFromMonth(year, month) + day) * 24 +                                                       // 1531
                        hour +                                                                                         // 1532
                        hourOffset * signOffset                                                                        // 1533
                    ) * 60;                                                                                            // 1534
                    result = (                                                                                         // 1535
                        (result + minute + minuteOffset * signOffset) * 60 +                                           // 1536
                        second                                                                                         // 1537
                    ) * 1000 + millisecond;                                                                            // 1538
                    if (isLocalTime) {                                                                                 // 1539
                        result = toUTC(result);                                                                        // 1540
                    }                                                                                                  // 1541
                    if (-8.64e15 <= result && result <= 8.64e15) {                                                     // 1542
                        return result;                                                                                 // 1543
                    }                                                                                                  // 1544
                }                                                                                                      // 1545
                return NaN;                                                                                            // 1546
            }                                                                                                          // 1547
            return NativeDate.parse.apply(this, arguments);                                                            // 1548
        };                                                                                                             // 1549
        defineProperties(DateShim, { parse: parseShim });                                                              // 1550
                                                                                                                       // 1551
        return DateShim;                                                                                               // 1552
    }(Date));                                                                                                          // 1553
    /* global Date: false */                                                                                           // 1554
}                                                                                                                      // 1555
                                                                                                                       // 1556
// ES5 15.9.4.4                                                                                                        // 1557
// http://es5.github.com/#x15.9.4.4                                                                                    // 1558
if (!Date.now) {                                                                                                       // 1559
    Date.now = function now() {                                                                                        // 1560
        return new Date().getTime();                                                                                   // 1561
    };                                                                                                                 // 1562
}                                                                                                                      // 1563
                                                                                                                       // 1564
//                                                                                                                     // 1565
// Number                                                                                                              // 1566
// ======                                                                                                              // 1567
//                                                                                                                     // 1568
                                                                                                                       // 1569
// ES5.1 15.7.4.5                                                                                                      // 1570
// http://es5.github.com/#x15.7.4.5                                                                                    // 1571
var hasToFixedBugs = NumberPrototype.toFixed && (                                                                      // 1572
  (0.00008).toFixed(3) !== '0.000' ||                                                                                  // 1573
  (0.9).toFixed(0) !== '1' ||                                                                                          // 1574
  (1.255).toFixed(2) !== '1.25' ||                                                                                     // 1575
  (1000000000000000128).toFixed(0) !== '1000000000000000128'                                                           // 1576
);                                                                                                                     // 1577
                                                                                                                       // 1578
var toFixedHelpers = {                                                                                                 // 1579
  base: 1e7,                                                                                                           // 1580
  size: 6,                                                                                                             // 1581
  data: [0, 0, 0, 0, 0, 0],                                                                                            // 1582
  multiply: function multiply(n, c) {                                                                                  // 1583
      var i = -1;                                                                                                      // 1584
      var c2 = c;                                                                                                      // 1585
      while (++i < toFixedHelpers.size) {                                                                              // 1586
          c2 += n * toFixedHelpers.data[i];                                                                            // 1587
          toFixedHelpers.data[i] = c2 % toFixedHelpers.base;                                                           // 1588
          c2 = Math.floor(c2 / toFixedHelpers.base);                                                                   // 1589
      }                                                                                                                // 1590
  },                                                                                                                   // 1591
  divide: function divide(n) {                                                                                         // 1592
      var i = toFixedHelpers.size;                                                                                     // 1593
      var c = 0;                                                                                                       // 1594
      while (--i >= 0) {                                                                                               // 1595
          c += toFixedHelpers.data[i];                                                                                 // 1596
          toFixedHelpers.data[i] = Math.floor(c / n);                                                                  // 1597
          c = (c % n) * toFixedHelpers.base;                                                                           // 1598
      }                                                                                                                // 1599
  },                                                                                                                   // 1600
  numToString: function numToString() {                                                                                // 1601
      var i = toFixedHelpers.size;                                                                                     // 1602
      var s = '';                                                                                                      // 1603
      while (--i >= 0) {                                                                                               // 1604
          if (s !== '' || i === 0 || toFixedHelpers.data[i] !== 0) {                                                   // 1605
              var t = $String(toFixedHelpers.data[i]);                                                                 // 1606
              if (s === '') {                                                                                          // 1607
                  s = t;                                                                                               // 1608
              } else {                                                                                                 // 1609
                  s += strSlice('0000000', 0, 7 - t.length) + t;                                                       // 1610
              }                                                                                                        // 1611
          }                                                                                                            // 1612
      }                                                                                                                // 1613
      return s;                                                                                                        // 1614
  },                                                                                                                   // 1615
  pow: function pow(x, n, acc) {                                                                                       // 1616
      return (n === 0 ? acc : (n % 2 === 1 ? pow(x, n - 1, acc * x) : pow(x * x, n / 2, acc)));                        // 1617
  },                                                                                                                   // 1618
  log: function log(x) {                                                                                               // 1619
      var n = 0;                                                                                                       // 1620
      var x2 = x;                                                                                                      // 1621
      while (x2 >= 4096) {                                                                                             // 1622
          n += 12;                                                                                                     // 1623
          x2 /= 4096;                                                                                                  // 1624
      }                                                                                                                // 1625
      while (x2 >= 2) {                                                                                                // 1626
          n += 1;                                                                                                      // 1627
          x2 /= 2;                                                                                                     // 1628
      }                                                                                                                // 1629
      return n;                                                                                                        // 1630
  }                                                                                                                    // 1631
};                                                                                                                     // 1632
                                                                                                                       // 1633
var toFixedShim = function toFixed(fractionDigits) {                                                                   // 1634
    var f, x, s, m, e, z, j, k;                                                                                        // 1635
                                                                                                                       // 1636
    // Test for NaN and round fractionDigits down                                                                      // 1637
    f = $Number(fractionDigits);                                                                                       // 1638
    f = isActualNaN(f) ? 0 : Math.floor(f);                                                                            // 1639
                                                                                                                       // 1640
    if (f < 0 || f > 20) {                                                                                             // 1641
        throw new RangeError('Number.toFixed called with invalid number of decimals');                                 // 1642
    }                                                                                                                  // 1643
                                                                                                                       // 1644
    x = $Number(this);                                                                                                 // 1645
                                                                                                                       // 1646
    if (isActualNaN(x)) {                                                                                              // 1647
        return 'NaN';                                                                                                  // 1648
    }                                                                                                                  // 1649
                                                                                                                       // 1650
    // If it is too big or small, return the string value of the number                                                // 1651
    if (x <= -1e21 || x >= 1e21) {                                                                                     // 1652
        return $String(x);                                                                                             // 1653
    }                                                                                                                  // 1654
                                                                                                                       // 1655
    s = '';                                                                                                            // 1656
                                                                                                                       // 1657
    if (x < 0) {                                                                                                       // 1658
        s = '-';                                                                                                       // 1659
        x = -x;                                                                                                        // 1660
    }                                                                                                                  // 1661
                                                                                                                       // 1662
    m = '0';                                                                                                           // 1663
                                                                                                                       // 1664
    if (x > 1e-21) {                                                                                                   // 1665
        // 1e-21 < x < 1e21                                                                                            // 1666
        // -70 < log2(x) < 70                                                                                          // 1667
        e = toFixedHelpers.log(x * toFixedHelpers.pow(2, 69, 1)) - 69;                                                 // 1668
        z = (e < 0 ? x * toFixedHelpers.pow(2, -e, 1) : x / toFixedHelpers.pow(2, e, 1));                              // 1669
        z *= 0x10000000000000; // Math.pow(2, 52);                                                                     // 1670
        e = 52 - e;                                                                                                    // 1671
                                                                                                                       // 1672
        // -18 < e < 122                                                                                               // 1673
        // x = z / 2 ^ e                                                                                               // 1674
        if (e > 0) {                                                                                                   // 1675
            toFixedHelpers.multiply(0, z);                                                                             // 1676
            j = f;                                                                                                     // 1677
                                                                                                                       // 1678
            while (j >= 7) {                                                                                           // 1679
                toFixedHelpers.multiply(1e7, 0);                                                                       // 1680
                j -= 7;                                                                                                // 1681
            }                                                                                                          // 1682
                                                                                                                       // 1683
            toFixedHelpers.multiply(toFixedHelpers.pow(10, j, 1), 0);                                                  // 1684
            j = e - 1;                                                                                                 // 1685
                                                                                                                       // 1686
            while (j >= 23) {                                                                                          // 1687
                toFixedHelpers.divide(1 << 23);                                                                        // 1688
                j -= 23;                                                                                               // 1689
            }                                                                                                          // 1690
                                                                                                                       // 1691
            toFixedHelpers.divide(1 << j);                                                                             // 1692
            toFixedHelpers.multiply(1, 1);                                                                             // 1693
            toFixedHelpers.divide(2);                                                                                  // 1694
            m = toFixedHelpers.numToString();                                                                          // 1695
        } else {                                                                                                       // 1696
            toFixedHelpers.multiply(0, z);                                                                             // 1697
            toFixedHelpers.multiply(1 << (-e), 0);                                                                     // 1698
            m = toFixedHelpers.numToString() + strSlice('0.00000000000000000000', 2, 2 + f);                           // 1699
        }                                                                                                              // 1700
    }                                                                                                                  // 1701
                                                                                                                       // 1702
    if (f > 0) {                                                                                                       // 1703
        k = m.length;                                                                                                  // 1704
                                                                                                                       // 1705
        if (k <= f) {                                                                                                  // 1706
            m = s + strSlice('0.0000000000000000000', 0, f - k + 2) + m;                                               // 1707
        } else {                                                                                                       // 1708
            m = s + strSlice(m, 0, k - f) + '.' + strSlice(m, k - f);                                                  // 1709
        }                                                                                                              // 1710
    } else {                                                                                                           // 1711
        m = s + m;                                                                                                     // 1712
    }                                                                                                                  // 1713
                                                                                                                       // 1714
    return m;                                                                                                          // 1715
};                                                                                                                     // 1716
defineProperties(NumberPrototype, { toFixed: toFixedShim }, hasToFixedBugs);                                           // 1717
                                                                                                                       // 1718
var hasToPrecisionUndefinedBug = (function () {                                                                        // 1719
    try {                                                                                                              // 1720
        return 1.0.toPrecision(undefined) === '1';                                                                     // 1721
    } catch (e) {                                                                                                      // 1722
        return true;                                                                                                   // 1723
    }                                                                                                                  // 1724
}());                                                                                                                  // 1725
var originalToPrecision = NumberPrototype.toPrecision;                                                                 // 1726
defineProperties(NumberPrototype, {                                                                                    // 1727
    toPrecision: function toPrecision(precision) {                                                                     // 1728
        return typeof precision === 'undefined' ? originalToPrecision.call(this) : originalToPrecision.call(this, precision);
    }                                                                                                                  // 1730
}, hasToPrecisionUndefinedBug);                                                                                        // 1731
                                                                                                                       // 1732
//                                                                                                                     // 1733
// String                                                                                                              // 1734
// ======                                                                                                              // 1735
//                                                                                                                     // 1736
                                                                                                                       // 1737
// ES5 15.5.4.14                                                                                                       // 1738
// http://es5.github.com/#x15.5.4.14                                                                                   // 1739
                                                                                                                       // 1740
// [bugfix, IE lt 9, firefox 4, Konqueror, Opera, obscure browsers]                                                    // 1741
// Many browsers do not split properly with regular expressions or they                                                // 1742
// do not perform the split correctly under obscure conditions.                                                        // 1743
// See http://blog.stevenlevithan.com/archives/cross-browser-split                                                     // 1744
// I've tested in many browsers and this seems to cover the deviant ones:                                              // 1745
//    'ab'.split(/(?:ab)*/) should be ["", ""], not [""]                                                               // 1746
//    '.'.split(/(.?)(.?)/) should be ["", ".", "", ""], not ["", ""]                                                  // 1747
//    'tesst'.split(/(s)*/) should be ["t", undefined, "e", "s", "t"], not                                             // 1748
//       [undefined, "t", undefined, "e", ...]                                                                         // 1749
//    ''.split(/.?/) should be [], not [""]                                                                            // 1750
//    '.'.split(/()()/) should be ["."], not ["", "", "."]                                                             // 1751
                                                                                                                       // 1752
if (                                                                                                                   // 1753
    'ab'.split(/(?:ab)*/).length !== 2 ||                                                                              // 1754
    '.'.split(/(.?)(.?)/).length !== 4 ||                                                                              // 1755
    'tesst'.split(/(s)*/)[1] === 't' ||                                                                                // 1756
    'test'.split(/(?:)/, -1).length !== 4 ||                                                                           // 1757
    ''.split(/.?/).length ||                                                                                           // 1758
    '.'.split(/()()/).length > 1                                                                                       // 1759
) {                                                                                                                    // 1760
    (function () {                                                                                                     // 1761
        var compliantExecNpcg = typeof (/()??/).exec('')[1] === 'undefined'; // NPCG: nonparticipating capturing group
        var maxSafe32BitInt = Math.pow(2, 32) - 1;                                                                     // 1763
                                                                                                                       // 1764
        StringPrototype.split = function (separator, limit) {                                                          // 1765
            var string = String(this);                                                                                 // 1766
            if (typeof separator === 'undefined' && limit === 0) {                                                     // 1767
                return [];                                                                                             // 1768
            }                                                                                                          // 1769
                                                                                                                       // 1770
            // If `separator` is not a regex, use native split                                                         // 1771
            if (!isRegex(separator)) {                                                                                 // 1772
                return strSplit(this, separator, limit);                                                               // 1773
            }                                                                                                          // 1774
                                                                                                                       // 1775
            var output = [];                                                                                           // 1776
            var flags = (separator.ignoreCase ? 'i' : '') +                                                            // 1777
                        (separator.multiline ? 'm' : '') +                                                             // 1778
                        (separator.unicode ? 'u' : '') + // in ES6                                                     // 1779
                        (separator.sticky ? 'y' : ''), // Firefox 3+ and ES6                                           // 1780
                lastLastIndex = 0,                                                                                     // 1781
                // Make `global` and avoid `lastIndex` issues by working with a copy                                   // 1782
                separator2, match, lastIndex, lastLength;                                                              // 1783
            var separatorCopy = new RegExp(separator.source, flags + 'g');                                             // 1784
            if (!compliantExecNpcg) {                                                                                  // 1785
                // Doesn't need flags gy, but they don't hurt                                                          // 1786
                separator2 = new RegExp('^' + separatorCopy.source + '$(?!\\s)', flags);                               // 1787
            }                                                                                                          // 1788
            /* Values for `limit`, per the spec:                                                                       // 1789
             * If undefined: 4294967295 // maxSafe32BitInt                                                             // 1790
             * If 0, Infinity, or NaN: 0                                                                               // 1791
             * If positive number: limit = Math.floor(limit); if (limit > 4294967295) limit -= 4294967296;             // 1792
             * If negative number: 4294967296 - Math.floor(Math.abs(limit))                                            // 1793
             * If other: Type-convert, then use the above rules                                                        // 1794
             */                                                                                                        // 1795
            var splitLimit = typeof limit === 'undefined' ? maxSafe32BitInt : ES.ToUint32(limit);                      // 1796
            match = separatorCopy.exec(string);                                                                        // 1797
            while (match) {                                                                                            // 1798
                // `separatorCopy.lastIndex` is not reliable cross-browser                                             // 1799
                lastIndex = match.index + match[0].length;                                                             // 1800
                if (lastIndex > lastLastIndex) {                                                                       // 1801
                    pushCall(output, strSlice(string, lastLastIndex, match.index));                                    // 1802
                    // Fix browsers whose `exec` methods don't consistently return `undefined` for                     // 1803
                    // nonparticipating capturing groups                                                               // 1804
                    if (!compliantExecNpcg && match.length > 1) {                                                      // 1805
                        /* eslint-disable no-loop-func */                                                              // 1806
                        match[0].replace(separator2, function () {                                                     // 1807
                            for (var i = 1; i < arguments.length - 2; i++) {                                           // 1808
                                if (typeof arguments[i] === 'undefined') {                                             // 1809
                                    match[i] = void 0;                                                                 // 1810
                                }                                                                                      // 1811
                            }                                                                                          // 1812
                        });                                                                                            // 1813
                        /* eslint-enable no-loop-func */                                                               // 1814
                    }                                                                                                  // 1815
                    if (match.length > 1 && match.index < string.length) {                                             // 1816
                        array_push.apply(output, arraySlice(match, 1));                                                // 1817
                    }                                                                                                  // 1818
                    lastLength = match[0].length;                                                                      // 1819
                    lastLastIndex = lastIndex;                                                                         // 1820
                    if (output.length >= splitLimit) {                                                                 // 1821
                        break;                                                                                         // 1822
                    }                                                                                                  // 1823
                }                                                                                                      // 1824
                if (separatorCopy.lastIndex === match.index) {                                                         // 1825
                    separatorCopy.lastIndex++; // Avoid an infinite loop                                               // 1826
                }                                                                                                      // 1827
                match = separatorCopy.exec(string);                                                                    // 1828
            }                                                                                                          // 1829
            if (lastLastIndex === string.length) {                                                                     // 1830
                if (lastLength || !separatorCopy.test('')) {                                                           // 1831
                    pushCall(output, '');                                                                              // 1832
                }                                                                                                      // 1833
            } else {                                                                                                   // 1834
                pushCall(output, strSlice(string, lastLastIndex));                                                     // 1835
            }                                                                                                          // 1836
            return output.length > splitLimit ? arraySlice(output, 0, splitLimit) : output;                            // 1837
        };                                                                                                             // 1838
    }());                                                                                                              // 1839
                                                                                                                       // 1840
// [bugfix, chrome]                                                                                                    // 1841
// If separator is undefined, then the result array contains just one String,                                          // 1842
// which is the this value (converted to a String). If limit is not undefined,                                         // 1843
// then the output array is truncated so that it contains no more than limit                                           // 1844
// elements.                                                                                                           // 1845
// "0".split(undefined, 0) -> []                                                                                       // 1846
} else if ('0'.split(void 0, 0).length) {                                                                              // 1847
    StringPrototype.split = function split(separator, limit) {                                                         // 1848
        if (typeof separator === 'undefined' && limit === 0) { return []; }                                            // 1849
        return strSplit(this, separator, limit);                                                                       // 1850
    };                                                                                                                 // 1851
}                                                                                                                      // 1852
                                                                                                                       // 1853
var str_replace = StringPrototype.replace;                                                                             // 1854
var replaceReportsGroupsCorrectly = (function () {                                                                     // 1855
    var groups = [];                                                                                                   // 1856
    'x'.replace(/x(.)?/g, function (match, group) {                                                                    // 1857
        pushCall(groups, group);                                                                                       // 1858
    });                                                                                                                // 1859
    return groups.length === 1 && typeof groups[0] === 'undefined';                                                    // 1860
}());                                                                                                                  // 1861
                                                                                                                       // 1862
if (!replaceReportsGroupsCorrectly) {                                                                                  // 1863
    StringPrototype.replace = function replace(searchValue, replaceValue) {                                            // 1864
        var isFn = isCallable(replaceValue);                                                                           // 1865
        var hasCapturingGroups = isRegex(searchValue) && (/\)[*?]/).test(searchValue.source);                          // 1866
        if (!isFn || !hasCapturingGroups) {                                                                            // 1867
            return str_replace.call(this, searchValue, replaceValue);                                                  // 1868
        } else {                                                                                                       // 1869
            var wrappedReplaceValue = function (match) {                                                               // 1870
                var length = arguments.length;                                                                         // 1871
                var originalLastIndex = searchValue.lastIndex;                                                         // 1872
                searchValue.lastIndex = 0;                                                                             // 1873
                var args = searchValue.exec(match) || [];                                                              // 1874
                searchValue.lastIndex = originalLastIndex;                                                             // 1875
                pushCall(args, arguments[length - 2], arguments[length - 1]);                                          // 1876
                return replaceValue.apply(this, args);                                                                 // 1877
            };                                                                                                         // 1878
            return str_replace.call(this, searchValue, wrappedReplaceValue);                                           // 1879
        }                                                                                                              // 1880
    };                                                                                                                 // 1881
}                                                                                                                      // 1882
                                                                                                                       // 1883
// ECMA-262, 3rd B.2.3                                                                                                 // 1884
// Not an ECMAScript standard, although ECMAScript 3rd Edition has a                                                   // 1885
// non-normative section suggesting uniform semantics and it should be                                                 // 1886
// normalized across all browsers                                                                                      // 1887
// [bugfix, IE lt 9] IE < 9 substr() with negative value not working in IE                                             // 1888
var string_substr = StringPrototype.substr;                                                                            // 1889
var hasNegativeSubstrBug = ''.substr && '0b'.substr(-1) !== 'b';                                                       // 1890
defineProperties(StringPrototype, {                                                                                    // 1891
    substr: function substr(start, length) {                                                                           // 1892
        var normalizedStart = start;                                                                                   // 1893
        if (start < 0) {                                                                                               // 1894
            normalizedStart = max(this.length + start, 0);                                                             // 1895
        }                                                                                                              // 1896
        return string_substr.call(this, normalizedStart, length);                                                      // 1897
    }                                                                                                                  // 1898
}, hasNegativeSubstrBug);                                                                                              // 1899
                                                                                                                       // 1900
// ES5 15.5.4.20                                                                                                       // 1901
// whitespace from: http://es5.github.io/#x15.5.4.20                                                                   // 1902
var ws = '\x09\x0A\x0B\x0C\x0D\x20\xA0\u1680\u180E\u2000\u2001\u2002\u2003' +                                          // 1903
    '\u2004\u2005\u2006\u2007\u2008\u2009\u200A\u202F\u205F\u3000\u2028' +                                             // 1904
    '\u2029\uFEFF';                                                                                                    // 1905
var zeroWidth = '\u200b';                                                                                              // 1906
var wsRegexChars = '[' + ws + ']';                                                                                     // 1907
var trimBeginRegexp = new RegExp('^' + wsRegexChars + wsRegexChars + '*');                                             // 1908
var trimEndRegexp = new RegExp(wsRegexChars + wsRegexChars + '*$');                                                    // 1909
var hasTrimWhitespaceBug = StringPrototype.trim && (ws.trim() || !zeroWidth.trim());                                   // 1910
defineProperties(StringPrototype, {                                                                                    // 1911
    // http://blog.stevenlevithan.com/archives/faster-trim-javascript                                                  // 1912
    // http://perfectionkills.com/whitespace-deviations/                                                               // 1913
    trim: function trim() {                                                                                            // 1914
        if (typeof this === 'undefined' || this === null) {                                                            // 1915
            throw new TypeError("can't convert " + this + ' to object');                                               // 1916
        }                                                                                                              // 1917
        return $String(this).replace(trimBeginRegexp, '').replace(trimEndRegexp, '');                                  // 1918
    }                                                                                                                  // 1919
}, hasTrimWhitespaceBug);                                                                                              // 1920
var trim = call.bind(String.prototype.trim);                                                                           // 1921
                                                                                                                       // 1922
var hasLastIndexBug = StringPrototype.lastIndexOf && 'abcあい'.lastIndexOf('あい', 2) !== -1;                              // 1923
defineProperties(StringPrototype, {                                                                                    // 1924
    lastIndexOf: function lastIndexOf(searchString) {                                                                  // 1925
        if (typeof this === 'undefined' || this === null) {                                                            // 1926
            throw new TypeError("can't convert " + this + ' to object');                                               // 1927
        }                                                                                                              // 1928
        var S = $String(this);                                                                                         // 1929
        var searchStr = $String(searchString);                                                                         // 1930
        var numPos = arguments.length > 1 ? $Number(arguments[1]) : NaN;                                               // 1931
        var pos = isActualNaN(numPos) ? Infinity : ES.ToInteger(numPos);                                               // 1932
        var start = min(max(pos, 0), S.length);                                                                        // 1933
        var searchLen = searchStr.length;                                                                              // 1934
        var k = start + searchLen;                                                                                     // 1935
        while (k > 0) {                                                                                                // 1936
            k = max(0, k - searchLen);                                                                                 // 1937
            var index = strIndexOf(strSlice(S, k, start + searchLen), searchStr);                                      // 1938
            if (index !== -1) {                                                                                        // 1939
                return k + index;                                                                                      // 1940
            }                                                                                                          // 1941
        }                                                                                                              // 1942
        return -1;                                                                                                     // 1943
    }                                                                                                                  // 1944
}, hasLastIndexBug);                                                                                                   // 1945
                                                                                                                       // 1946
var originalLastIndexOf = StringPrototype.lastIndexOf;                                                                 // 1947
defineProperties(StringPrototype, {                                                                                    // 1948
    lastIndexOf: function lastIndexOf(searchString) {                                                                  // 1949
        return originalLastIndexOf.apply(this, arguments);                                                             // 1950
    }                                                                                                                  // 1951
}, StringPrototype.lastIndexOf.length !== 1);                                                                          // 1952
                                                                                                                       // 1953
// ES-5 15.1.2.2                                                                                                       // 1954
/* eslint-disable radix */                                                                                             // 1955
if (parseInt(ws + '08') !== 8 || parseInt(ws + '0x16') !== 22) {                                                       // 1956
/* eslint-enable radix */                                                                                              // 1957
    /* global parseInt: true */                                                                                        // 1958
    parseInt = (function (origParseInt) {                                                                              // 1959
        var hexRegex = /^[\-+]?0[xX]/;                                                                                 // 1960
        return function parseInt(str, radix) {                                                                         // 1961
            var string = trim(str);                                                                                    // 1962
            var defaultedRadix = $Number(radix) || (hexRegex.test(string) ? 16 : 10);                                  // 1963
            return origParseInt(string, defaultedRadix);                                                               // 1964
        };                                                                                                             // 1965
    }(parseInt));                                                                                                      // 1966
}                                                                                                                      // 1967
                                                                                                                       // 1968
// https://es5.github.io/#x15.1.2.3                                                                                    // 1969
if (1 / parseFloat('-0') !== -Infinity) {                                                                              // 1970
    /* global parseFloat: true */                                                                                      // 1971
    parseFloat = (function (origParseFloat) {                                                                          // 1972
        return function parseFloat(string) {                                                                           // 1973
            var inputString = trim(string);                                                                            // 1974
            var result = origParseFloat(inputString);                                                                  // 1975
            return result === 0 && strSlice(inputString, 0, 1) === '-' ? -0 : result;                                  // 1976
        };                                                                                                             // 1977
    }(parseFloat));                                                                                                    // 1978
}                                                                                                                      // 1979
                                                                                                                       // 1980
if (String(new RangeError('test')) !== 'RangeError: test') {                                                           // 1981
    var errorToStringShim = function toString() {                                                                      // 1982
        if (typeof this === 'undefined' || this === null) {                                                            // 1983
            throw new TypeError("can't convert " + this + ' to object');                                               // 1984
        }                                                                                                              // 1985
        var name = this.name;                                                                                          // 1986
        if (typeof name === 'undefined') {                                                                             // 1987
            name = 'Error';                                                                                            // 1988
        } else if (typeof name !== 'string') {                                                                         // 1989
            name = $String(name);                                                                                      // 1990
        }                                                                                                              // 1991
        var msg = this.message;                                                                                        // 1992
        if (typeof msg === 'undefined') {                                                                              // 1993
            msg = '';                                                                                                  // 1994
        } else if (typeof msg !== 'string') {                                                                          // 1995
            msg = $String(msg);                                                                                        // 1996
        }                                                                                                              // 1997
        if (!name) {                                                                                                   // 1998
            return msg;                                                                                                // 1999
        }                                                                                                              // 2000
        if (!msg) {                                                                                                    // 2001
            return name;                                                                                               // 2002
        }                                                                                                              // 2003
        return name + ': ' + msg;                                                                                      // 2004
    };                                                                                                                 // 2005
    // can't use defineProperties here because of toString enumeration issue in IE <= 8                                // 2006
    Error.prototype.toString = errorToStringShim;                                                                      // 2007
}                                                                                                                      // 2008
                                                                                                                       // 2009
if (supportsDescriptors) {                                                                                             // 2010
    var ensureNonEnumerable = function (obj, prop) {                                                                   // 2011
        if (isEnum(obj, prop)) {                                                                                       // 2012
            var desc = Object.getOwnPropertyDescriptor(obj, prop);                                                     // 2013
            desc.enumerable = false;                                                                                   // 2014
            Object.defineProperty(obj, prop, desc);                                                                    // 2015
        }                                                                                                              // 2016
    };                                                                                                                 // 2017
    ensureNonEnumerable(Error.prototype, 'message');                                                                   // 2018
    if (Error.prototype.message !== '') {                                                                              // 2019
      Error.prototype.message = '';                                                                                    // 2020
    }                                                                                                                  // 2021
    ensureNonEnumerable(Error.prototype, 'name');                                                                      // 2022
}                                                                                                                      // 2023
                                                                                                                       // 2024
if (String(/a/mig) !== '/a/gim') {                                                                                     // 2025
    var regexToString = function toString() {                                                                          // 2026
        var str = '/' + this.source + '/';                                                                             // 2027
        if (this.global) {                                                                                             // 2028
            str += 'g';                                                                                                // 2029
        }                                                                                                              // 2030
        if (this.ignoreCase) {                                                                                         // 2031
            str += 'i';                                                                                                // 2032
        }                                                                                                              // 2033
        if (this.multiline) {                                                                                          // 2034
            str += 'm';                                                                                                // 2035
        }                                                                                                              // 2036
        return str;                                                                                                    // 2037
    };                                                                                                                 // 2038
    // can't use defineProperties here because of toString enumeration issue in IE <= 8                                // 2039
    RegExp.prototype.toString = regexToString;                                                                         // 2040
}                                                                                                                      // 2041
                                                                                                                       // 2042
}));                                                                                                                   // 2043
                                                                                                                       // 2044
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"es5-sham.js":function(require,exports,module){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// node_modules/meteor/es5-shim/node_modules/es5-shim/es5-sham.js                                                      //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
/*!                                                                                                                    // 1
 * https://github.com/es-shims/es5-shim                                                                                // 2
 * @license es5-shim Copyright 2009-2015 by contributors, MIT License                                                  // 3
 * see https://github.com/es-shims/es5-shim/blob/master/LICENSE                                                        // 4
 */                                                                                                                    // 5
                                                                                                                       // 6
// vim: ts=4 sts=4 sw=4 expandtab                                                                                      // 7
                                                                                                                       // 8
// Add semicolon to prevent IIFE from being passed as argument to concatenated code.                                   // 9
;                                                                                                                      // 10
                                                                                                                       // 11
// UMD (Universal Module Definition)                                                                                   // 12
// see https://github.com/umdjs/umd/blob/master/templates/returnExports.js                                             // 13
(function (root, factory) {                                                                                            // 14
    'use strict';                                                                                                      // 15
                                                                                                                       // 16
    /* global define, exports, module */                                                                               // 17
    if (typeof define === 'function' && define.amd) {                                                                  // 18
        // AMD. Register as an anonymous module.                                                                       // 19
        define(factory);                                                                                               // 20
    } else if (typeof exports === 'object') {                                                                          // 21
        // Node. Does not work with strict CommonJS, but                                                               // 22
        // only CommonJS-like enviroments that support module.exports,                                                 // 23
        // like Node.                                                                                                  // 24
        module.exports = factory();                                                                                    // 25
    } else {                                                                                                           // 26
        // Browser globals (root is window)                                                                            // 27
        root.returnExports = factory();                                                                                // 28
  }                                                                                                                    // 29
}(this, function () {                                                                                                  // 30
                                                                                                                       // 31
var call = Function.call;                                                                                              // 32
var prototypeOfObject = Object.prototype;                                                                              // 33
var owns = call.bind(prototypeOfObject.hasOwnProperty);                                                                // 34
var isEnumerable = call.bind(prototypeOfObject.propertyIsEnumerable);                                                  // 35
var toStr = call.bind(prototypeOfObject.toString);                                                                     // 36
                                                                                                                       // 37
// If JS engine supports accessors creating shortcuts.                                                                 // 38
var defineGetter;                                                                                                      // 39
var defineSetter;                                                                                                      // 40
var lookupGetter;                                                                                                      // 41
var lookupSetter;                                                                                                      // 42
var supportsAccessors = owns(prototypeOfObject, '__defineGetter__');                                                   // 43
if (supportsAccessors) {                                                                                               // 44
    /* eslint-disable no-underscore-dangle */                                                                          // 45
    defineGetter = call.bind(prototypeOfObject.__defineGetter__);                                                      // 46
    defineSetter = call.bind(prototypeOfObject.__defineSetter__);                                                      // 47
    lookupGetter = call.bind(prototypeOfObject.__lookupGetter__);                                                      // 48
    lookupSetter = call.bind(prototypeOfObject.__lookupSetter__);                                                      // 49
    /* eslint-enable no-underscore-dangle */                                                                           // 50
}                                                                                                                      // 51
                                                                                                                       // 52
// ES5 15.2.3.2                                                                                                        // 53
// http://es5.github.com/#x15.2.3.2                                                                                    // 54
if (!Object.getPrototypeOf) {                                                                                          // 55
    // https://github.com/es-shims/es5-shim/issues#issue/2                                                             // 56
    // http://ejohn.org/blog/objectgetprototypeof/                                                                     // 57
    // recommended by fschaefer on github                                                                              // 58
    //                                                                                                                 // 59
    // sure, and webreflection says ^_^                                                                                // 60
    // ... this will nerever possibly return null                                                                      // 61
    // ... Opera Mini breaks here with infinite loops                                                                  // 62
    Object.getPrototypeOf = function getPrototypeOf(object) {                                                          // 63
        /* eslint-disable no-proto */                                                                                  // 64
        var proto = object.__proto__;                                                                                  // 65
        /* eslint-enable no-proto */                                                                                   // 66
        if (proto || proto === null) {                                                                                 // 67
            return proto;                                                                                              // 68
        } else if (toStr(object.constructor) === '[object Function]') {                                                // 69
            return object.constructor.prototype;                                                                       // 70
        } else if (object instanceof Object) {                                                                         // 71
          return prototypeOfObject;                                                                                    // 72
        } else {                                                                                                       // 73
          // Correctly return null for Objects created with `Object.create(null)`                                      // 74
          // (shammed or native) or `{ __proto__: null}`.  Also returns null for                                       // 75
          // cross-realm objects on browsers that lack `__proto__` support (like                                       // 76
          // IE <11), but that's the best we can do.                                                                   // 77
          return null;                                                                                                 // 78
        }                                                                                                              // 79
    };                                                                                                                 // 80
}                                                                                                                      // 81
                                                                                                                       // 82
// ES5 15.2.3.3                                                                                                        // 83
// http://es5.github.com/#x15.2.3.3                                                                                    // 84
                                                                                                                       // 85
var doesGetOwnPropertyDescriptorWork = function doesGetOwnPropertyDescriptorWork(object) {                             // 86
    try {                                                                                                              // 87
        object.sentinel = 0;                                                                                           // 88
        return Object.getOwnPropertyDescriptor(object, 'sentinel').value === 0;                                        // 89
    } catch (exception) {                                                                                              // 90
        return false;                                                                                                  // 91
    }                                                                                                                  // 92
};                                                                                                                     // 93
                                                                                                                       // 94
// check whether getOwnPropertyDescriptor works if it's given. Otherwise, shim partially.                              // 95
if (Object.defineProperty) {                                                                                           // 96
    var getOwnPropertyDescriptorWorksOnObject = doesGetOwnPropertyDescriptorWork({});                                  // 97
    var getOwnPropertyDescriptorWorksOnDom = typeof document === 'undefined' ||                                        // 98
    doesGetOwnPropertyDescriptorWork(document.createElement('div'));                                                   // 99
    if (!getOwnPropertyDescriptorWorksOnDom || !getOwnPropertyDescriptorWorksOnObject) {                               // 100
        var getOwnPropertyDescriptorFallback = Object.getOwnPropertyDescriptor;                                        // 101
    }                                                                                                                  // 102
}                                                                                                                      // 103
                                                                                                                       // 104
if (!Object.getOwnPropertyDescriptor || getOwnPropertyDescriptorFallback) {                                            // 105
    var ERR_NON_OBJECT = 'Object.getOwnPropertyDescriptor called on a non-object: ';                                   // 106
                                                                                                                       // 107
    /* eslint-disable no-proto */                                                                                      // 108
    Object.getOwnPropertyDescriptor = function getOwnPropertyDescriptor(object, property) {                            // 109
        if ((typeof object !== 'object' && typeof object !== 'function') || object === null) {                         // 110
            throw new TypeError(ERR_NON_OBJECT + object);                                                              // 111
        }                                                                                                              // 112
                                                                                                                       // 113
        // make a valiant attempt to use the real getOwnPropertyDescriptor                                             // 114
        // for I8's DOM elements.                                                                                      // 115
        if (getOwnPropertyDescriptorFallback) {                                                                        // 116
            try {                                                                                                      // 117
                return getOwnPropertyDescriptorFallback.call(Object, object, property);                                // 118
            } catch (exception) {                                                                                      // 119
                // try the shim if the real one doesn't work                                                           // 120
            }                                                                                                          // 121
        }                                                                                                              // 122
                                                                                                                       // 123
        var descriptor;                                                                                                // 124
                                                                                                                       // 125
        // If object does not owns property return undefined immediately.                                              // 126
        if (!owns(object, property)) {                                                                                 // 127
            return descriptor;                                                                                         // 128
        }                                                                                                              // 129
                                                                                                                       // 130
        // If object has a property then it's for sure `configurable`, and                                             // 131
        // probably `enumerable`. Detect enumerability though.                                                         // 132
        descriptor = {                                                                                                 // 133
            enumerable: isEnumerable(object, property),                                                                // 134
            configurable: true                                                                                         // 135
        };                                                                                                             // 136
                                                                                                                       // 137
        // If JS engine supports accessor properties then property may be a                                            // 138
        // getter or setter.                                                                                           // 139
        if (supportsAccessors) {                                                                                       // 140
            // Unfortunately `__lookupGetter__` will return a getter even                                              // 141
            // if object has own non getter property along with a same named                                           // 142
            // inherited getter. To avoid misbehavior we temporary remove                                              // 143
            // `__proto__` so that `__lookupGetter__` will return getter only                                          // 144
            // if it's owned by an object.                                                                             // 145
            var prototype = object.__proto__;                                                                          // 146
            var notPrototypeOfObject = object !== prototypeOfObject;                                                   // 147
            // avoid recursion problem, breaking in Opera Mini when                                                    // 148
            // Object.getOwnPropertyDescriptor(Object.prototype, 'toString')                                           // 149
            // or any other Object.prototype accessor                                                                  // 150
            if (notPrototypeOfObject) {                                                                                // 151
                object.__proto__ = prototypeOfObject;                                                                  // 152
            }                                                                                                          // 153
                                                                                                                       // 154
            var getter = lookupGetter(object, property);                                                               // 155
            var setter = lookupSetter(object, property);                                                               // 156
                                                                                                                       // 157
            if (notPrototypeOfObject) {                                                                                // 158
                // Once we have getter and setter we can put values back.                                              // 159
                object.__proto__ = prototype;                                                                          // 160
            }                                                                                                          // 161
                                                                                                                       // 162
            if (getter || setter) {                                                                                    // 163
                if (getter) {                                                                                          // 164
                    descriptor.get = getter;                                                                           // 165
                }                                                                                                      // 166
                if (setter) {                                                                                          // 167
                    descriptor.set = setter;                                                                           // 168
                }                                                                                                      // 169
                // If it was accessor property we're done and return here                                              // 170
                // in order to avoid adding `value` to the descriptor.                                                 // 171
                return descriptor;                                                                                     // 172
            }                                                                                                          // 173
        }                                                                                                              // 174
                                                                                                                       // 175
        // If we got this far we know that object has an own property that is                                          // 176
        // not an accessor so we set it as a value and return descriptor.                                              // 177
        descriptor.value = object[property];                                                                           // 178
        descriptor.writable = true;                                                                                    // 179
        return descriptor;                                                                                             // 180
    };                                                                                                                 // 181
    /* eslint-enable no-proto */                                                                                       // 182
}                                                                                                                      // 183
                                                                                                                       // 184
// ES5 15.2.3.4                                                                                                        // 185
// http://es5.github.com/#x15.2.3.4                                                                                    // 186
if (!Object.getOwnPropertyNames) {                                                                                     // 187
    Object.getOwnPropertyNames = function getOwnPropertyNames(object) {                                                // 188
        return Object.keys(object);                                                                                    // 189
    };                                                                                                                 // 190
}                                                                                                                      // 191
                                                                                                                       // 192
// ES5 15.2.3.5                                                                                                        // 193
// http://es5.github.com/#x15.2.3.5                                                                                    // 194
if (!Object.create) {                                                                                                  // 195
                                                                                                                       // 196
    // Contributed by Brandon Benvie, October, 2012                                                                    // 197
    var createEmpty;                                                                                                   // 198
    var supportsProto = !({ __proto__: null } instanceof Object);                                                      // 199
                        // the following produces false positives                                                      // 200
                        // in Opera Mini => not a reliable check                                                       // 201
                        // Object.prototype.__proto__ === null                                                         // 202
                                                                                                                       // 203
    // Check for document.domain and active x support                                                                  // 204
    // No need to use active x approach when document.domain is not set                                                // 205
    // see https://github.com/es-shims/es5-shim/issues/150                                                             // 206
    // variation of https://github.com/kitcambridge/es5-shim/commit/4f738ac066346                                      // 207
    /* global ActiveXObject */                                                                                         // 208
    var shouldUseActiveX = function shouldUseActiveX() {                                                               // 209
        // return early if document.domain not set                                                                     // 210
        if (!document.domain) {                                                                                        // 211
            return false;                                                                                              // 212
        }                                                                                                              // 213
                                                                                                                       // 214
        try {                                                                                                          // 215
            return !!new ActiveXObject('htmlfile');                                                                    // 216
        } catch (exception) {                                                                                          // 217
            return false;                                                                                              // 218
        }                                                                                                              // 219
    };                                                                                                                 // 220
                                                                                                                       // 221
    // This supports IE8 when document.domain is used                                                                  // 222
    // see https://github.com/es-shims/es5-shim/issues/150                                                             // 223
    // variation of https://github.com/kitcambridge/es5-shim/commit/4f738ac066346                                      // 224
    var getEmptyViaActiveX = function getEmptyViaActiveX() {                                                           // 225
        var empty;                                                                                                     // 226
        var xDoc;                                                                                                      // 227
                                                                                                                       // 228
        xDoc = new ActiveXObject('htmlfile');                                                                          // 229
                                                                                                                       // 230
        xDoc.write('<script><\/script>');                                                                              // 231
        xDoc.close();                                                                                                  // 232
                                                                                                                       // 233
        empty = xDoc.parentWindow.Object.prototype;                                                                    // 234
        xDoc = null;                                                                                                   // 235
                                                                                                                       // 236
        return empty;                                                                                                  // 237
    };                                                                                                                 // 238
                                                                                                                       // 239
    // The original implementation using an iframe                                                                     // 240
    // before the activex approach was added                                                                           // 241
    // see https://github.com/es-shims/es5-shim/issues/150                                                             // 242
    var getEmptyViaIFrame = function getEmptyViaIFrame() {                                                             // 243
        var iframe = document.createElement('iframe');                                                                 // 244
        var parent = document.body || document.documentElement;                                                        // 245
        var empty;                                                                                                     // 246
                                                                                                                       // 247
        iframe.style.display = 'none';                                                                                 // 248
        parent.appendChild(iframe);                                                                                    // 249
        /* eslint-disable no-script-url */                                                                             // 250
        iframe.src = 'javascript:';                                                                                    // 251
        /* eslint-enable no-script-url */                                                                              // 252
                                                                                                                       // 253
        empty = iframe.contentWindow.Object.prototype;                                                                 // 254
        parent.removeChild(iframe);                                                                                    // 255
        iframe = null;                                                                                                 // 256
                                                                                                                       // 257
        return empty;                                                                                                  // 258
    };                                                                                                                 // 259
                                                                                                                       // 260
    /* global document */                                                                                              // 261
    if (supportsProto || typeof document === 'undefined') {                                                            // 262
        createEmpty = function () {                                                                                    // 263
            return { __proto__: null };                                                                                // 264
        };                                                                                                             // 265
    } else {                                                                                                           // 266
        // In old IE __proto__ can't be used to manually set `null`, nor does                                          // 267
        // any other method exist to make an object that inherits from nothing,                                        // 268
        // aside from Object.prototype itself. Instead, create a new global                                            // 269
        // object and *steal* its Object.prototype and strip it bare. This is                                          // 270
        // used as the prototype to create nullary objects.                                                            // 271
        createEmpty = function () {                                                                                    // 272
            // Determine which approach to use                                                                         // 273
            // see https://github.com/es-shims/es5-shim/issues/150                                                     // 274
            var empty = shouldUseActiveX() ? getEmptyViaActiveX() : getEmptyViaIFrame();                               // 275
                                                                                                                       // 276
            delete empty.constructor;                                                                                  // 277
            delete empty.hasOwnProperty;                                                                               // 278
            delete empty.propertyIsEnumerable;                                                                         // 279
            delete empty.isPrototypeOf;                                                                                // 280
            delete empty.toLocaleString;                                                                               // 281
            delete empty.toString;                                                                                     // 282
            delete empty.valueOf;                                                                                      // 283
                                                                                                                       // 284
            var Empty = function Empty() {};                                                                           // 285
            Empty.prototype = empty;                                                                                   // 286
            // short-circuit future calls                                                                              // 287
            createEmpty = function () {                                                                                // 288
                return new Empty();                                                                                    // 289
            };                                                                                                         // 290
            return new Empty();                                                                                        // 291
        };                                                                                                             // 292
    }                                                                                                                  // 293
                                                                                                                       // 294
    Object.create = function create(prototype, properties) {                                                           // 295
                                                                                                                       // 296
        var object;                                                                                                    // 297
        var Type = function Type() {}; // An empty constructor.                                                        // 298
                                                                                                                       // 299
        if (prototype === null) {                                                                                      // 300
            object = createEmpty();                                                                                    // 301
        } else {                                                                                                       // 302
            if (typeof prototype !== 'object' && typeof prototype !== 'function') {                                    // 303
                // In the native implementation `parent` can be `null`                                                 // 304
                // OR *any* `instanceof Object`  (Object|Function|Array|RegExp|etc)                                    // 305
                // Use `typeof` tho, b/c in old IE, DOM elements are not `instanceof Object`                           // 306
                // like they are in modern browsers. Using `Object.create` on DOM elements                             // 307
                // is...err...probably inappropriate, but the native version allows for it.                            // 308
                throw new TypeError('Object prototype may only be an Object or null'); // same msg as Chrome           // 309
            }                                                                                                          // 310
            Type.prototype = prototype;                                                                                // 311
            object = new Type();                                                                                       // 312
            // IE has no built-in implementation of `Object.getPrototypeOf`                                            // 313
            // neither `__proto__`, but this manually setting `__proto__` will                                         // 314
            // guarantee that `Object.getPrototypeOf` will work as expected with                                       // 315
            // objects created using `Object.create`                                                                   // 316
            /* eslint-disable no-proto */                                                                              // 317
            object.__proto__ = prototype;                                                                              // 318
            /* eslint-enable no-proto */                                                                               // 319
        }                                                                                                              // 320
                                                                                                                       // 321
        if (properties !== void 0) {                                                                                   // 322
            Object.defineProperties(object, properties);                                                               // 323
        }                                                                                                              // 324
                                                                                                                       // 325
        return object;                                                                                                 // 326
    };                                                                                                                 // 327
}                                                                                                                      // 328
                                                                                                                       // 329
// ES5 15.2.3.6                                                                                                        // 330
// http://es5.github.com/#x15.2.3.6                                                                                    // 331
                                                                                                                       // 332
// Patch for WebKit and IE8 standard mode                                                                              // 333
// Designed by hax <hax.github.com>                                                                                    // 334
// related issue: https://github.com/es-shims/es5-shim/issues#issue/5                                                  // 335
// IE8 Reference:                                                                                                      // 336
//     http://msdn.microsoft.com/en-us/library/dd282900.aspx                                                           // 337
//     http://msdn.microsoft.com/en-us/library/dd229916.aspx                                                           // 338
// WebKit Bugs:                                                                                                        // 339
//     https://bugs.webkit.org/show_bug.cgi?id=36423                                                                   // 340
                                                                                                                       // 341
var doesDefinePropertyWork = function doesDefinePropertyWork(object) {                                                 // 342
    try {                                                                                                              // 343
        Object.defineProperty(object, 'sentinel', {});                                                                 // 344
        return 'sentinel' in object;                                                                                   // 345
    } catch (exception) {                                                                                              // 346
        return false;                                                                                                  // 347
    }                                                                                                                  // 348
};                                                                                                                     // 349
                                                                                                                       // 350
// check whether defineProperty works if it's given. Otherwise,                                                        // 351
// shim partially.                                                                                                     // 352
if (Object.defineProperty) {                                                                                           // 353
    var definePropertyWorksOnObject = doesDefinePropertyWork({});                                                      // 354
    var definePropertyWorksOnDom = typeof document === 'undefined' ||                                                  // 355
        doesDefinePropertyWork(document.createElement('div'));                                                         // 356
    if (!definePropertyWorksOnObject || !definePropertyWorksOnDom) {                                                   // 357
        var definePropertyFallback = Object.defineProperty,                                                            // 358
            definePropertiesFallback = Object.defineProperties;                                                        // 359
    }                                                                                                                  // 360
}                                                                                                                      // 361
                                                                                                                       // 362
if (!Object.defineProperty || definePropertyFallback) {                                                                // 363
    var ERR_NON_OBJECT_DESCRIPTOR = 'Property description must be an object: ';                                        // 364
    var ERR_NON_OBJECT_TARGET = 'Object.defineProperty called on non-object: ';                                        // 365
    var ERR_ACCESSORS_NOT_SUPPORTED = 'getters & setters can not be defined on this javascript engine';                // 366
                                                                                                                       // 367
    Object.defineProperty = function defineProperty(object, property, descriptor) {                                    // 368
        if ((typeof object !== 'object' && typeof object !== 'function') || object === null) {                         // 369
            throw new TypeError(ERR_NON_OBJECT_TARGET + object);                                                       // 370
        }                                                                                                              // 371
        if ((typeof descriptor !== 'object' && typeof descriptor !== 'function') || descriptor === null) {             // 372
            throw new TypeError(ERR_NON_OBJECT_DESCRIPTOR + descriptor);                                               // 373
        }                                                                                                              // 374
        // make a valiant attempt to use the real defineProperty                                                       // 375
        // for I8's DOM elements.                                                                                      // 376
        if (definePropertyFallback) {                                                                                  // 377
            try {                                                                                                      // 378
                return definePropertyFallback.call(Object, object, property, descriptor);                              // 379
            } catch (exception) {                                                                                      // 380
                // try the shim if the real one doesn't work                                                           // 381
            }                                                                                                          // 382
        }                                                                                                              // 383
                                                                                                                       // 384
        // If it's a data property.                                                                                    // 385
        if ('value' in descriptor) {                                                                                   // 386
            // fail silently if 'writable', 'enumerable', or 'configurable'                                            // 387
            // are requested but not supported                                                                         // 388
            /*                                                                                                         // 389
            // alternate approach:                                                                                     // 390
            if ( // can't implement these features; allow false but not true                                           // 391
                ('writable' in descriptor && !descriptor.writable) ||                                                  // 392
                ('enumerable' in descriptor && !descriptor.enumerable) ||                                              // 393
                ('configurable' in descriptor && !descriptor.configurable)                                             // 394
            ))                                                                                                         // 395
                throw new RangeError(                                                                                  // 396
                    'This implementation of Object.defineProperty does not support configurable, enumerable, or writable.'
                );                                                                                                     // 398
            */                                                                                                         // 399
                                                                                                                       // 400
            if (supportsAccessors && (lookupGetter(object, property) || lookupSetter(object, property))) {             // 401
                // As accessors are supported only on engines implementing                                             // 402
                // `__proto__` we can safely override `__proto__` while defining                                       // 403
                // a property to make sure that we don't hit an inherited                                              // 404
                // accessor.                                                                                           // 405
                /* eslint-disable no-proto */                                                                          // 406
                var prototype = object.__proto__;                                                                      // 407
                object.__proto__ = prototypeOfObject;                                                                  // 408
                // Deleting a property anyway since getter / setter may be                                             // 409
                // defined on object itself.                                                                           // 410
                delete object[property];                                                                               // 411
                object[property] = descriptor.value;                                                                   // 412
                // Setting original `__proto__` back now.                                                              // 413
                object.__proto__ = prototype;                                                                          // 414
                /* eslint-enable no-proto */                                                                           // 415
            } else {                                                                                                   // 416
                object[property] = descriptor.value;                                                                   // 417
            }                                                                                                          // 418
        } else {                                                                                                       // 419
            if (!supportsAccessors && (('get' in descriptor) || ('set' in descriptor))) {                              // 420
                throw new TypeError(ERR_ACCESSORS_NOT_SUPPORTED);                                                      // 421
            }                                                                                                          // 422
            // If we got that far then getters and setters can be defined !!                                           // 423
            if ('get' in descriptor) {                                                                                 // 424
                defineGetter(object, property, descriptor.get);                                                        // 425
            }                                                                                                          // 426
            if ('set' in descriptor) {                                                                                 // 427
                defineSetter(object, property, descriptor.set);                                                        // 428
            }                                                                                                          // 429
        }                                                                                                              // 430
        return object;                                                                                                 // 431
    };                                                                                                                 // 432
}                                                                                                                      // 433
                                                                                                                       // 434
// ES5 15.2.3.7                                                                                                        // 435
// http://es5.github.com/#x15.2.3.7                                                                                    // 436
if (!Object.defineProperties || definePropertiesFallback) {                                                            // 437
    Object.defineProperties = function defineProperties(object, properties) {                                          // 438
        // make a valiant attempt to use the real defineProperties                                                     // 439
        if (definePropertiesFallback) {                                                                                // 440
            try {                                                                                                      // 441
                return definePropertiesFallback.call(Object, object, properties);                                      // 442
            } catch (exception) {                                                                                      // 443
                // try the shim if the real one doesn't work                                                           // 444
            }                                                                                                          // 445
        }                                                                                                              // 446
                                                                                                                       // 447
        Object.keys(properties).forEach(function (property) {                                                          // 448
            if (property !== '__proto__') {                                                                            // 449
                Object.defineProperty(object, property, properties[property]);                                         // 450
            }                                                                                                          // 451
        });                                                                                                            // 452
        return object;                                                                                                 // 453
    };                                                                                                                 // 454
}                                                                                                                      // 455
                                                                                                                       // 456
// ES5 15.2.3.8                                                                                                        // 457
// http://es5.github.com/#x15.2.3.8                                                                                    // 458
if (!Object.seal) {                                                                                                    // 459
    Object.seal = function seal(object) {                                                                              // 460
        if (Object(object) !== object) {                                                                               // 461
            throw new TypeError('Object.seal can only be called on Objects.');                                         // 462
        }                                                                                                              // 463
        // this is misleading and breaks feature-detection, but                                                        // 464
        // allows "securable" code to "gracefully" degrade to working                                                  // 465
        // but insecure code.                                                                                          // 466
        return object;                                                                                                 // 467
    };                                                                                                                 // 468
}                                                                                                                      // 469
                                                                                                                       // 470
// ES5 15.2.3.9                                                                                                        // 471
// http://es5.github.com/#x15.2.3.9                                                                                    // 472
if (!Object.freeze) {                                                                                                  // 473
    Object.freeze = function freeze(object) {                                                                          // 474
        if (Object(object) !== object) {                                                                               // 475
            throw new TypeError('Object.freeze can only be called on Objects.');                                       // 476
        }                                                                                                              // 477
        // this is misleading and breaks feature-detection, but                                                        // 478
        // allows "securable" code to "gracefully" degrade to working                                                  // 479
        // but insecure code.                                                                                          // 480
        return object;                                                                                                 // 481
    };                                                                                                                 // 482
}                                                                                                                      // 483
                                                                                                                       // 484
// detect a Rhino bug and patch it                                                                                     // 485
try {                                                                                                                  // 486
    Object.freeze(function () {});                                                                                     // 487
} catch (exception) {                                                                                                  // 488
    Object.freeze = (function (freezeObject) {                                                                         // 489
        return function freeze(object) {                                                                               // 490
            if (typeof object === 'function') {                                                                        // 491
                return object;                                                                                         // 492
            } else {                                                                                                   // 493
                return freezeObject(object);                                                                           // 494
            }                                                                                                          // 495
        };                                                                                                             // 496
    }(Object.freeze));                                                                                                 // 497
}                                                                                                                      // 498
                                                                                                                       // 499
// ES5 15.2.3.10                                                                                                       // 500
// http://es5.github.com/#x15.2.3.10                                                                                   // 501
if (!Object.preventExtensions) {                                                                                       // 502
    Object.preventExtensions = function preventExtensions(object) {                                                    // 503
        if (Object(object) !== object) {                                                                               // 504
            throw new TypeError('Object.preventExtensions can only be called on Objects.');                            // 505
        }                                                                                                              // 506
        // this is misleading and breaks feature-detection, but                                                        // 507
        // allows "securable" code to "gracefully" degrade to working                                                  // 508
        // but insecure code.                                                                                          // 509
        return object;                                                                                                 // 510
    };                                                                                                                 // 511
}                                                                                                                      // 512
                                                                                                                       // 513
// ES5 15.2.3.11                                                                                                       // 514
// http://es5.github.com/#x15.2.3.11                                                                                   // 515
if (!Object.isSealed) {                                                                                                // 516
    Object.isSealed = function isSealed(object) {                                                                      // 517
        if (Object(object) !== object) {                                                                               // 518
            throw new TypeError('Object.isSealed can only be called on Objects.');                                     // 519
        }                                                                                                              // 520
        return false;                                                                                                  // 521
    };                                                                                                                 // 522
}                                                                                                                      // 523
                                                                                                                       // 524
// ES5 15.2.3.12                                                                                                       // 525
// http://es5.github.com/#x15.2.3.12                                                                                   // 526
if (!Object.isFrozen) {                                                                                                // 527
    Object.isFrozen = function isFrozen(object) {                                                                      // 528
        if (Object(object) !== object) {                                                                               // 529
            throw new TypeError('Object.isFrozen can only be called on Objects.');                                     // 530
        }                                                                                                              // 531
        return false;                                                                                                  // 532
    };                                                                                                                 // 533
}                                                                                                                      // 534
                                                                                                                       // 535
// ES5 15.2.3.13                                                                                                       // 536
// http://es5.github.com/#x15.2.3.13                                                                                   // 537
if (!Object.isExtensible) {                                                                                            // 538
    Object.isExtensible = function isExtensible(object) {                                                              // 539
        // 1. If Type(O) is not Object throw a TypeError exception.                                                    // 540
        if (Object(object) !== object) {                                                                               // 541
            throw new TypeError('Object.isExtensible can only be called on Objects.');                                 // 542
        }                                                                                                              // 543
        // 2. Return the Boolean value of the [[Extensible]] internal property of O.                                   // 544
        var name = '';                                                                                                 // 545
        while (owns(object, name)) {                                                                                   // 546
            name += '?';                                                                                               // 547
        }                                                                                                              // 548
        object[name] = true;                                                                                           // 549
        var returnValue = owns(object, name);                                                                          // 550
        delete object[name];                                                                                           // 551
        return returnValue;                                                                                            // 552
    };                                                                                                                 // 553
}                                                                                                                      // 554
                                                                                                                       // 555
}));                                                                                                                   // 556
                                                                                                                       // 557
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}}}},{"extensions":[".js",".json"]});
var exports = require("./node_modules/meteor/es5-shim/client.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
Package['es5-shim'] = exports;

})();
