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
var HTML = Package.htmljs.HTML;
var Tracker = Package.tracker.Tracker;
var Deps = Package.tracker.Deps;
var Blaze = Package.blaze.Blaze;
var UI = Package.blaze.UI;
var Handlebars = Package.blaze.Handlebars;
var ObserveSequence = Package['observe-sequence'].ObserveSequence;
var _ = Package.underscore._;

/* Package-scope variables */
var Spacebars;

(function(){

///////////////////////////////////////////////////////////////////////////////////
//                                                                               //
// packages/spacebars/spacebars-runtime.js                                       //
//                                                                               //
///////////////////////////////////////////////////////////////////////////////////
                                                                                 //
Spacebars = {};                                                                  // 1
                                                                                 // 2
var tripleEquals = function (a, b) { return a === b; };                          // 3
                                                                                 // 4
Spacebars.include = function (templateOrFunction, contentFunc, elseFunc) {       // 5
  if (! templateOrFunction)                                                      // 6
    return null;                                                                 // 7
                                                                                 // 8
  if (typeof templateOrFunction !== 'function') {                                // 9
    var template = templateOrFunction;                                           // 10
    if (! Blaze.isTemplate(template))                                            // 11
      throw new Error("Expected template or null, found: " + template);          // 12
    var view = templateOrFunction.constructView(contentFunc, elseFunc);          // 13
    view.__startsNewLexicalScope = true;                                         // 14
    return view;                                                                 // 15
  }                                                                              // 16
                                                                                 // 17
  var templateVar = Blaze.ReactiveVar(null, tripleEquals);                       // 18
  var view = Blaze.View('Spacebars.include', function () {                       // 19
    var template = templateVar.get();                                            // 20
    if (template === null)                                                       // 21
      return null;                                                               // 22
                                                                                 // 23
    if (! Blaze.isTemplate(template))                                            // 24
      throw new Error("Expected template or null, found: " + template);          // 25
                                                                                 // 26
    return template.constructView(contentFunc, elseFunc);                        // 27
  });                                                                            // 28
  view.__templateVar = templateVar;                                              // 29
  view.onViewCreated(function () {                                               // 30
    this.autorun(function () {                                                   // 31
      templateVar.set(templateOrFunction());                                     // 32
    });                                                                          // 33
  });                                                                            // 34
  view.__startsNewLexicalScope = true;                                           // 35
                                                                                 // 36
  return view;                                                                   // 37
};                                                                               // 38
                                                                                 // 39
// Executes `{{foo bar baz}}` when called on `(foo, bar, baz)`.                  // 40
// If `bar` and `baz` are functions, they are called before                      // 41
// `foo` is called on them.                                                      // 42
//                                                                               // 43
// This is the shared part of Spacebars.mustache and                             // 44
// Spacebars.attrMustache, which differ in how they post-process the             // 45
// result.                                                                       // 46
Spacebars.mustacheImpl = function (value/*, args*/) {                            // 47
  var args = arguments;                                                          // 48
  // if we have any arguments (pos or kw), add an options argument               // 49
  // if there isn't one.                                                         // 50
  if (args.length > 1) {                                                         // 51
    var kw = args[args.length - 1];                                              // 52
    if (! (kw instanceof Spacebars.kw)) {                                        // 53
      kw = Spacebars.kw();                                                       // 54
      // clone arguments into an actual array, then push                         // 55
      // the empty kw object.                                                    // 56
      args = Array.prototype.slice.call(arguments);                              // 57
      args.push(kw);                                                             // 58
    } else {                                                                     // 59
      // For each keyword arg, call it if it's a function                        // 60
      var newHash = {};                                                          // 61
      for (var k in kw.hash) {                                                   // 62
        var v = kw.hash[k];                                                      // 63
        newHash[k] = (typeof v === 'function' ? v() : v);                        // 64
      }                                                                          // 65
      args[args.length - 1] = Spacebars.kw(newHash);                             // 66
    }                                                                            // 67
  }                                                                              // 68
                                                                                 // 69
  return Spacebars.call.apply(null, args);                                       // 70
};                                                                               // 71
                                                                                 // 72
Spacebars.mustache = function (value/*, args*/) {                                // 73
  var result = Spacebars.mustacheImpl.apply(null, arguments);                    // 74
                                                                                 // 75
  if (result instanceof Spacebars.SafeString)                                    // 76
    return HTML.Raw(result.toString());                                          // 77
  else                                                                           // 78
    // map `null`, `undefined`, and `false` to null, which is important          // 79
    // so that attributes with nully values are considered absent.               // 80
    // stringify anything else (e.g. strings, booleans, numbers including 0).    // 81
    return (result == null || result === false) ? null : String(result);         // 82
};                                                                               // 83
                                                                                 // 84
Spacebars.attrMustache = function (value/*, args*/) {                            // 85
  var result = Spacebars.mustacheImpl.apply(null, arguments);                    // 86
                                                                                 // 87
  if (result == null || result === '') {                                         // 88
    return null;                                                                 // 89
  } else if (typeof result === 'object') {                                       // 90
    return result;                                                               // 91
  } else if (typeof result === 'string' && HTML.isValidAttributeName(result)) {  // 92
    var obj = {};                                                                // 93
    obj[result] = '';                                                            // 94
    return obj;                                                                  // 95
  } else {                                                                       // 96
    throw new Error("Expected valid attribute name, '', null, or object");       // 97
  }                                                                              // 98
};                                                                               // 99
                                                                                 // 100
Spacebars.dataMustache = function (value/*, args*/) {                            // 101
  var result = Spacebars.mustacheImpl.apply(null, arguments);                    // 102
                                                                                 // 103
  return result;                                                                 // 104
};                                                                               // 105
                                                                                 // 106
// Idempotently wrap in `HTML.Raw`.                                              // 107
//                                                                               // 108
// Called on the return value from `Spacebars.mustache` in case the              // 109
// template uses triple-stache (`{{{foo bar baz}}}`).                            // 110
Spacebars.makeRaw = function (value) {                                           // 111
  if (value == null) // null or undefined                                        // 112
    return null;                                                                 // 113
  else if (value instanceof HTML.Raw)                                            // 114
    return value;                                                                // 115
  else                                                                           // 116
    return HTML.Raw(value);                                                      // 117
};                                                                               // 118
                                                                                 // 119
// If `value` is a function, evaluate its `args` (by calling them, if they       // 120
// are functions), and then call it on them. Otherwise, return `value`.          // 121
//                                                                               // 122
// If `value` is not a function and is not null, then this method will assert    // 123
// that there are no args. We check for null before asserting because a user     // 124
// may write a template like {{user.fullNameWithPrefix 'Mr.'}}, where the        // 125
// function will be null until data is ready.                                    // 126
Spacebars.call = function (value/*, args*/) {                                    // 127
  if (typeof value === 'function') {                                             // 128
    // Evaluate arguments by calling them if they are functions.                 // 129
    var newArgs = [];                                                            // 130
    for (var i = 1; i < arguments.length; i++) {                                 // 131
      var arg = arguments[i];                                                    // 132
      newArgs[i-1] = (typeof arg === 'function' ? arg() : arg);                  // 133
    }                                                                            // 134
                                                                                 // 135
    return value.apply(null, newArgs);                                           // 136
  } else {                                                                       // 137
    if (value != null && arguments.length > 1) {                                 // 138
      throw new Error("Can't call non-function: " + value);                      // 139
    }                                                                            // 140
    return value;                                                                // 141
  }                                                                              // 142
};                                                                               // 143
                                                                                 // 144
// Call this as `Spacebars.kw({ ... })`.  The return value                       // 145
// is `instanceof Spacebars.kw`.                                                 // 146
Spacebars.kw = function (hash) {                                                 // 147
  if (! (this instanceof Spacebars.kw))                                          // 148
    // called without new; call with new                                         // 149
    return new Spacebars.kw(hash);                                               // 150
                                                                                 // 151
  this.hash = hash || {};                                                        // 152
};                                                                               // 153
                                                                                 // 154
// Call this as `Spacebars.SafeString("some HTML")`.  The return value           // 155
// is `instanceof Spacebars.SafeString` (and `instanceof Handlebars.SafeString).
Spacebars.SafeString = function (html) {                                         // 157
  if (! (this instanceof Spacebars.SafeString))                                  // 158
    // called without new; call with new                                         // 159
    return new Spacebars.SafeString(html);                                       // 160
                                                                                 // 161
  return new Handlebars.SafeString(html);                                        // 162
};                                                                               // 163
Spacebars.SafeString.prototype = Handlebars.SafeString.prototype;                // 164
                                                                                 // 165
// `Spacebars.dot(foo, "bar", "baz")` performs a special kind                    // 166
// of `foo.bar.baz` that allows safe indexing of `null` and                      // 167
// indexing of functions (which calls the function).  If the                     // 168
// result is a function, it is always a bound function (e.g.                     // 169
// a wrapped version of `baz` that always uses `foo.bar` as                      // 170
// `this`).                                                                      // 171
//                                                                               // 172
// In `Spacebars.dot(foo, "bar")`, `foo` is assumed to be either                 // 173
// a non-function value or a "fully-bound" function wrapping a value,            // 174
// where fully-bound means it takes no arguments and ignores `this`.             // 175
//                                                                               // 176
// `Spacebars.dot(foo, "bar")` performs the following steps:                     // 177
//                                                                               // 178
// * If `foo` is falsy, return `foo`.                                            // 179
//                                                                               // 180
// * If `foo` is a function, call it (set `foo` to `foo()`).                     // 181
//                                                                               // 182
// * If `foo` is falsy now, return `foo`.                                        // 183
//                                                                               // 184
// * Return `foo.bar`, binding it to `foo` if it's a function.                   // 185
Spacebars.dot = function (value, id1/*, id2, ...*/) {                            // 186
  if (arguments.length > 2) {                                                    // 187
    // Note: doing this recursively is probably less efficient than              // 188
    // doing it in an iterative loop.                                            // 189
    var argsForRecurse = [];                                                     // 190
    argsForRecurse.push(Spacebars.dot(value, id1));                              // 191
    argsForRecurse.push.apply(argsForRecurse,                                    // 192
                              Array.prototype.slice.call(arguments, 2));         // 193
    return Spacebars.dot.apply(null, argsForRecurse);                            // 194
  }                                                                              // 195
                                                                                 // 196
  if (typeof value === 'function')                                               // 197
    value = value();                                                             // 198
                                                                                 // 199
  if (! value)                                                                   // 200
    return value; // falsy, don't index, pass through                            // 201
                                                                                 // 202
  var result = value[id1];                                                       // 203
  if (typeof result !== 'function')                                              // 204
    return result;                                                               // 205
  // `value[id1]` (or `value()[id1]`) is a function.                             // 206
  // Bind it so that when called, `value` will be placed in `this`.              // 207
  return function (/*arguments*/) {                                              // 208
    return result.apply(value, arguments);                                       // 209
  };                                                                             // 210
};                                                                               // 211
                                                                                 // 212
// Spacebars.With implements the conditional logic of rendering                  // 213
// the `{{else}}` block if the argument is falsy.  It combines                   // 214
// a Blaze.If with a Blaze.With (the latter only in the truthy                   // 215
// case, since the else block is evaluated without entering                      // 216
// a new data context).                                                          // 217
Spacebars.With = function (argFunc, contentFunc, elseFunc) {                     // 218
  var argVar = new Blaze.ReactiveVar;                                            // 219
  var view = Blaze.View('Spacebars_with', function () {                          // 220
    return Blaze.If(function () { return argVar.get(); },                        // 221
                    function () { return Blaze.With(function () {                // 222
                      return argVar.get(); }, contentFunc); },                   // 223
                    elseFunc);                                                   // 224
  });                                                                            // 225
  view.onViewCreated(function () {                                               // 226
    this.autorun(function () {                                                   // 227
      argVar.set(argFunc());                                                     // 228
                                                                                 // 229
      // This is a hack so that autoruns inside the body                         // 230
      // of the #with get stopped sooner.  It reaches inside                     // 231
      // our ReactiveVar to access its dep.                                      // 232
                                                                                 // 233
      Tracker.onInvalidate(function () {                                         // 234
        argVar.dep.changed();                                                    // 235
      });                                                                        // 236
                                                                                 // 237
      // Take the case of `{{#with A}}{{B}}{{/with}}`.  The goal                 // 238
      // is to not re-render `B` if `A` changes to become falsy                  // 239
      // and `B` is simultaneously invalidated.                                  // 240
      //                                                                         // 241
      // A series of autoruns are involved:                                      // 242
      //                                                                         // 243
      // 1. This autorun (argument to Spacebars.With)                            // 244
      // 2. Argument to Blaze.If                                                 // 245
      // 3. Blaze.If view re-render                                              // 246
      // 4. Argument to Blaze.With                                               // 247
      // 5. The template tag `{{B}}`                                             // 248
      //                                                                         // 249
      // When (3) is invalidated, it immediately stops (4) and (5)               // 250
      // because of a Tracker.onInvalidate built into materializeView.           // 251
      // (When a View's render method is invalidated, it immediately             // 252
      // tears down all the subviews, via a Tracker.onInvalidate much            // 253
      // like this one.                                                          // 254
      //                                                                         // 255
      // Suppose `A` changes to become falsy, and `B` changes at the             // 256
      // same time (i.e. without an intervening flush).                          // 257
      // Without the code above, this happens:                                   // 258
      //                                                                         // 259
      // - (1) and (5) are invalidated.                                          // 260
      // - (1) runs, invalidating (2) and (4).                                   // 261
      // - (5) runs.                                                             // 262
      // - (2) runs, invalidating (3), stopping (4) and (5).                     // 263
      //                                                                         // 264
      // With the code above:                                                    // 265
      //                                                                         // 266
      // - (1) and (5) are invalidated, invalidating (2) and (4).                // 267
      // - (1) runs.                                                             // 268
      // - (2) runs, invalidating (3), stopping (4) and (5).                     // 269
      //                                                                         // 270
      // If the re-run of (5) is originally enqueued before (1), all             // 271
      // bets are off, but typically that doesn't seem to be the                 // 272
      // case.  Anyway, doing this is always better than not doing it,           // 273
      // because it might save a bunch of DOM from being updated                 // 274
      // needlessly.                                                             // 275
    });                                                                          // 276
  });                                                                            // 277
                                                                                 // 278
  return view;                                                                   // 279
};                                                                               // 280
                                                                                 // 281
// XXX COMPAT WITH 0.9.0                                                         // 282
Spacebars.TemplateWith = Blaze._TemplateWith;                                    // 283
                                                                                 // 284
///////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.spacebars = {}, {
  Spacebars: Spacebars
});

})();
