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
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;
var Symbol = Package['ecmascript-runtime'].Symbol;
var Map = Package['ecmascript-runtime'].Map;
var Set = Package['ecmascript-runtime'].Set;
var meteorBabelHelpers = Package['babel-runtime'].meteorBabelHelpers;
var Promise = Package.promise.Promise;

/* Package-scope variables */
var Random;

var require = meteorInstall({"node_modules":{"meteor":{"random":{"random.js":function(require){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/random/random.js                                                                                           //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// We use cryptographically strong PRNGs (crypto.getRandomBytes() on the server,                                       //
// window.crypto.getRandomValues() in the browser) when available. If these                                            //
// PRNGs fail, we fall back to the Alea PRNG, which is not cryptographically                                           //
// strong, and we seed it with various sources such as the date, Math.random,                                          //
// and window size on the client.  When using crypto.getRandomValues(), our                                            //
// primitive is hexString(), from which we construct fraction(). When using                                            //
// window.crypto.getRandomValues() or alea, the primitive is fraction and we use                                       //
// that to construct hex string.                                                                                       //
                                                                                                                       //
if (Meteor.isServer) var nodeCrypto = Npm.require('crypto');                                                           // 10
                                                                                                                       //
// see http://baagoe.org/en/wiki/Better_random_numbers_for_javascript                                                  //
// for a full discussion and Alea implementation.                                                                      //
var Alea = function Alea() {                                                                                           // 15
  function Mash() {                                                                                                    // 16
    var n = 0xefc8249d;                                                                                                // 17
                                                                                                                       //
    var mash = function mash(data) {                                                                                   // 19
      data = data.toString();                                                                                          // 20
      for (var i = 0; i < data.length; i++) {                                                                          // 21
        n += data.charCodeAt(i);                                                                                       // 22
        var h = 0.02519603282416938 * n;                                                                               // 23
        n = h >>> 0;                                                                                                   // 24
        h -= n;                                                                                                        // 25
        h *= n;                                                                                                        // 26
        n = h >>> 0;                                                                                                   // 27
        h -= n;                                                                                                        // 28
        n += h * 0x100000000; // 2^32                                                                                  // 29
      }                                                                                                                // 21
      return (n >>> 0) * 2.3283064365386963e-10; // 2^-32                                                              // 31
    };                                                                                                                 // 19
                                                                                                                       //
    mash.version = 'Mash 0.9';                                                                                         // 34
    return mash;                                                                                                       // 35
  }                                                                                                                    //
                                                                                                                       //
  return function (args) {                                                                                             // 38
    var s0 = 0;                                                                                                        // 39
    var s1 = 0;                                                                                                        // 40
    var s2 = 0;                                                                                                        // 41
    var c = 1;                                                                                                         // 42
                                                                                                                       //
    if (args.length == 0) {                                                                                            // 44
      args = [+new Date()];                                                                                            // 45
    }                                                                                                                  //
    var mash = Mash();                                                                                                 // 47
    s0 = mash(' ');                                                                                                    // 48
    s1 = mash(' ');                                                                                                    // 49
    s2 = mash(' ');                                                                                                    // 50
                                                                                                                       //
    for (var i = 0; i < args.length; i++) {                                                                            // 52
      s0 -= mash(args[i]);                                                                                             // 53
      if (s0 < 0) {                                                                                                    // 54
        s0 += 1;                                                                                                       // 55
      }                                                                                                                //
      s1 -= mash(args[i]);                                                                                             // 57
      if (s1 < 0) {                                                                                                    // 58
        s1 += 1;                                                                                                       // 59
      }                                                                                                                //
      s2 -= mash(args[i]);                                                                                             // 61
      if (s2 < 0) {                                                                                                    // 62
        s2 += 1;                                                                                                       // 63
      }                                                                                                                //
    }                                                                                                                  //
    mash = null;                                                                                                       // 66
                                                                                                                       //
    var random = function random() {                                                                                   // 68
      var t = 2091639 * s0 + c * 2.3283064365386963e-10; // 2^-32                                                      // 69
      s0 = s1;                                                                                                         // 68
      s1 = s2;                                                                                                         // 71
      return s2 = t - (c = t | 0);                                                                                     // 72
    };                                                                                                                 //
    random.uint32 = function () {                                                                                      // 74
      return random() * 0x100000000; // 2^32                                                                           // 75
    };                                                                                                                 // 74
    random.fract53 = function () {                                                                                     // 77
      return random() + (random() * 0x200000 | 0) * 1.1102230246251565e-16; // 2^-53                                   // 78
    };                                                                                                                 // 77
    random.version = 'Alea 0.9';                                                                                       // 81
    random.args = args;                                                                                                // 82
    return random;                                                                                                     // 83
  }(Array.prototype.slice.call(arguments));                                                                            //
};                                                                                                                     //
                                                                                                                       //
var UNMISTAKABLE_CHARS = "23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz";                                    // 88
var BASE64_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789-_";                            // 89
                                                                                                                       //
// `type` is one of `RandomGenerator.Type` as defined below.                                                           //
//                                                                                                                     //
// options:                                                                                                            //
// - seeds: (required, only for RandomGenerator.Type.ALEA) an array                                                    //
//   whose items will be `toString`ed and used as the seed to the Alea                                                 //
//   algorithm                                                                                                         //
var RandomGenerator = function RandomGenerator(type, options) {                                                        // 98
  var self = this;                                                                                                     // 99
  self.type = type;                                                                                                    // 100
                                                                                                                       //
  if (!RandomGenerator.Type[type]) {                                                                                   // 102
    throw new Error("Unknown random generator type: " + type);                                                         // 103
  }                                                                                                                    //
                                                                                                                       //
  if (type === RandomGenerator.Type.ALEA) {                                                                            // 106
    if (!options.seeds) {                                                                                              // 107
      throw new Error("No seeds were provided for Alea PRNG");                                                         // 108
    }                                                                                                                  //
    self.alea = Alea.apply(null, options.seeds);                                                                       // 110
  }                                                                                                                    //
};                                                                                                                     //
                                                                                                                       //
// Types of PRNGs supported by the `RandomGenerator` class                                                             //
RandomGenerator.Type = {                                                                                               // 115
  // Use Node's built-in `crypto.getRandomBytes` (cryptographically                                                    //
  // secure but not seedable, runs only on the server). Reverts to                                                     //
  // `crypto.getPseudoRandomBytes` in the extremely uncommon case that                                                 //
  // there isn't enough entropy yet                                                                                    //
  NODE_CRYPTO: "NODE_CRYPTO",                                                                                          // 120
                                                                                                                       //
  // Use non-IE browser's built-in `window.crypto.getRandomValues`                                                     //
  // (cryptographically secure but not seedable, runs only in the                                                      //
  // browser).                                                                                                         //
  BROWSER_CRYPTO: "BROWSER_CRYPTO",                                                                                    // 125
                                                                                                                       //
  // Use the *fast*, seedaable and not cryptographically secure                                                        //
  // Alea algorithm                                                                                                    //
  ALEA: "ALEA"                                                                                                         // 129
};                                                                                                                     //
                                                                                                                       //
RandomGenerator.prototype.fraction = function () {                                                                     // 132
  var self = this;                                                                                                     // 133
  if (self.type === RandomGenerator.Type.ALEA) {                                                                       // 134
    return self.alea();                                                                                                // 135
  } else if (self.type === RandomGenerator.Type.NODE_CRYPTO) {                                                         //
    var numerator = parseInt(self.hexString(8), 16);                                                                   // 137
    return numerator * 2.3283064365386963e-10; // 2^-32                                                                // 138
  } else if (self.type === RandomGenerator.Type.BROWSER_CRYPTO) {                                                      // 136
      var array = new Uint32Array(1);                                                                                  // 140
      window.crypto.getRandomValues(array);                                                                            // 141
      return array[0] * 2.3283064365386963e-10; // 2^-32                                                               // 142
    } else {                                                                                                           // 139
        throw new Error('Unknown random generator type: ' + self.type);                                                // 144
      }                                                                                                                //
};                                                                                                                     //
                                                                                                                       //
RandomGenerator.prototype.hexString = function (digits) {                                                              // 148
  var self = this;                                                                                                     // 149
  if (self.type === RandomGenerator.Type.NODE_CRYPTO) {                                                                // 150
    var numBytes = Math.ceil(digits / 2);                                                                              // 151
    var bytes;                                                                                                         // 152
    // Try to get cryptographically strong randomness. Fall back to                                                    //
    // non-cryptographically strong if not available.                                                                  //
    try {                                                                                                              // 150
      bytes = nodeCrypto.randomBytes(numBytes);                                                                        // 156
    } catch (e) {                                                                                                      //
      // XXX should re-throw any error except insufficient entropy                                                     //
      bytes = nodeCrypto.pseudoRandomBytes(numBytes);                                                                  // 159
    }                                                                                                                  //
    var result = bytes.toString("hex");                                                                                // 161
    // If the number of digits is odd, we'll have generated an extra 4 bits                                            //
    // of randomness, so we need to trim the last digit.                                                               //
    return result.substring(0, digits);                                                                                // 150
  } else {                                                                                                             //
    return this._randomString(digits, "0123456789abcdef");                                                             // 166
  }                                                                                                                    //
};                                                                                                                     //
                                                                                                                       //
RandomGenerator.prototype._randomString = function (charsCount, alphabet) {                                            // 170
  var self = this;                                                                                                     // 172
  var digits = [];                                                                                                     // 173
  for (var i = 0; i < charsCount; i++) {                                                                               // 174
    digits[i] = self.choice(alphabet);                                                                                 // 175
  }                                                                                                                    //
  return digits.join("");                                                                                              // 177
};                                                                                                                     //
                                                                                                                       //
RandomGenerator.prototype.id = function (charsCount) {                                                                 // 180
  var self = this;                                                                                                     // 181
  // 17 characters is around 96 bits of entropy, which is the amount of                                                //
  // state in the Alea PRNG.                                                                                           //
  if (charsCount === undefined) charsCount = 17;                                                                       // 180
                                                                                                                       //
  return self._randomString(charsCount, UNMISTAKABLE_CHARS);                                                           // 187
};                                                                                                                     //
                                                                                                                       //
RandomGenerator.prototype.secret = function (charsCount) {                                                             // 190
  var self = this;                                                                                                     // 191
  // Default to 256 bits of entropy, or 43 characters at 6 bits per                                                    //
  // character.                                                                                                        //
  if (charsCount === undefined) charsCount = 43;                                                                       // 190
  return self._randomString(charsCount, BASE64_CHARS);                                                                 // 196
};                                                                                                                     //
                                                                                                                       //
RandomGenerator.prototype.choice = function (arrayOrString) {                                                          // 199
  var index = Math.floor(this.fraction() * arrayOrString.length);                                                      // 200
  if (typeof arrayOrString === "string") return arrayOrString.substr(index, 1);else return arrayOrString[index];       // 201
};                                                                                                                     //
                                                                                                                       //
// instantiate RNG.  Heuristically collect entropy from various sources when a                                         //
// cryptographic PRNG isn't available.                                                                                 //
                                                                                                                       //
// client sources                                                                                                      //
var height = typeof window !== 'undefined' && window.innerHeight || typeof document !== 'undefined' && document.documentElement && document.documentElement.clientHeight || typeof document !== 'undefined' && document.body && document.body.clientHeight || 1;
                                                                                                                       //
var width = typeof window !== 'undefined' && window.innerWidth || typeof document !== 'undefined' && document.documentElement && document.documentElement.clientWidth || typeof document !== 'undefined' && document.body && document.body.clientWidth || 1;
                                                                                                                       //
var agent = typeof navigator !== 'undefined' && navigator.userAgent || "";                                             // 229
                                                                                                                       //
function createAleaGeneratorWithGeneratedSeed() {                                                                      // 231
  return new RandomGenerator(RandomGenerator.Type.ALEA, { seeds: [new Date(), height, width, agent, Math.random()] });
};                                                                                                                     //
                                                                                                                       //
if (Meteor.isServer) {                                                                                                 // 237
  Random = new RandomGenerator(RandomGenerator.Type.NODE_CRYPTO);                                                      // 238
} else {                                                                                                               //
  if (typeof window !== "undefined" && window.crypto && window.crypto.getRandomValues) {                               // 240
    Random = new RandomGenerator(RandomGenerator.Type.BROWSER_CRYPTO);                                                 // 242
  } else {                                                                                                             //
    // On IE 10 and below, there's no browser crypto API                                                               //
    // available. Fall back to Alea                                                                                    //
    //                                                                                                                 //
    // XXX looks like at the moment, we use Alea in IE 11 as well,                                                     //
    // which has `window.msCrypto` instead of `window.crypto`.                                                         //
    Random = createAleaGeneratorWithGeneratedSeed();                                                                   // 249
  }                                                                                                                    //
}                                                                                                                      //
                                                                                                                       //
// Create a non-cryptographically secure PRNG with a given seed (using                                                 //
// the Alea algorithm)                                                                                                 //
Random.createWithSeeds = function () {                                                                                 // 255
  for (var _len = arguments.length, seeds = Array(_len), _key = 0; _key < _len; _key++) {                              //
    seeds[_key] = arguments[_key];                                                                                     //
  }                                                                                                                    //
                                                                                                                       //
  if (seeds.length === 0) {                                                                                            // 256
    throw new Error("No seeds were provided");                                                                         // 257
  }                                                                                                                    //
  return new RandomGenerator(RandomGenerator.Type.ALEA, { seeds: seeds });                                             // 259
};                                                                                                                     //
                                                                                                                       //
// Used like `Random`, but much faster and not cryptographically                                                       //
// secure                                                                                                              //
Random.insecure = createAleaGeneratorWithGeneratedSeed();                                                              // 264
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"deprecated.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/random/deprecated.js                                                                                       //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// Before this package existed, we used to use this Meteor.uuid()                                                      //
// implementing the RFC 4122 v4 UUID. It is no longer documented                                                       //
// and will go away.                                                                                                   //
// XXX COMPAT WITH 0.5.6                                                                                               //
Meteor.uuid = function () {                                                                                            // 5
  var HEX_DIGITS = "0123456789abcdef";                                                                                 // 6
  var s = [];                                                                                                          // 7
  for (var i = 0; i < 36; i++) {                                                                                       // 8
    s[i] = Random.choice(HEX_DIGITS);                                                                                  // 9
  }                                                                                                                    //
  s[14] = "4";                                                                                                         // 11
  s[19] = HEX_DIGITS.substr(parseInt(s[19], 16) & 0x3 | 0x8, 1);                                                       // 12
  s[8] = s[13] = s[18] = s[23] = "-";                                                                                  // 13
                                                                                                                       //
  var uuid = s.join("");                                                                                               // 15
  return uuid;                                                                                                         // 16
};                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}},{"extensions":[".js",".json"]});
require("./node_modules/meteor/random/random.js");
require("./node_modules/meteor/random/deprecated.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.random = {}, {
  Random: Random
});

})();
