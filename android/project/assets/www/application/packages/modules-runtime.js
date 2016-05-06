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

/* Package-scope variables */
var makeInstaller, meteorInstall;

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
// packages/modules-runtime/.npm/package/node_modules/install/install.js   //
// This file is in bare mode and is not in its own closure.                //
//                                                                         //
/////////////////////////////////////////////////////////////////////////////
                                                                           //
makeInstaller = function (options) {                                       // 1
  options = options || {};                                                 // 2
                                                                           // 3
  // These file extensions will be appended to required module identifiers
  // if they do not exactly match an installed module.                     // 5
  var defaultExtensions = options.extensions || [".js", ".json"];          // 6
                                                                           // 7
  // This constructor will be used to instantiate the module objects       // 8
  // passed to module factory functions (i.e. the third argument after     // 9
  // require and exports).                                                 // 10
  var Module = options.Module || function Module(id) {                     // 11
    this.id = id;                                                          // 12
    this.children = [];                                                    // 13
  };                                                                       // 14
                                                                           // 15
  // If defined, the options.onInstall function will be called any time    // 16
  // new modules are installed.                                            // 17
  var onInstall = options.onInstall;                                       // 18
                                                                           // 19
  // If defined, the options.override function will be called before       // 20
  // looking up any top-level package identifiers in node_modules          // 21
  // directories. It can either return a string to provide an alternate    // 22
  // package identifier, or a non-string value to prevent the lookup from  // 23
  // proceeding.                                                           // 24
  var override = options.override;                                         // 25
                                                                           // 26
  // If defined, the options.fallback function will be called when no      // 27
  // installed module is found for a required module identifier. Often     // 28
  // options.fallback will be implemented in terms of the native Node      // 29
  // require function, which has the ability to load binary modules.       // 30
  var fallback = options.fallback;                                         // 31
                                                                           // 32
  // Nothing special about MISSING.hasOwnProperty, except that it's fewer  // 33
  // characters than Object.prototype.hasOwnProperty after minification.   // 34
  var hasOwn = {}.hasOwnProperty;                                          // 35
                                                                           // 36
  // The file object representing the root directory of the installed      // 37
  // module tree.                                                          // 38
  var root = new File("/", new File("/.."));                               // 39
  var rootRequire = makeRequire(root);                                     // 40
                                                                           // 41
  // Merges the given tree of directories and module factory functions     // 42
  // into the tree of installed modules and returns a require function     // 43
  // that behaves as if called from a module in the root directory.        // 44
  function install(tree, options) {                                        // 45
    if (isObject(tree)) {                                                  // 46
      fileMergeContents(root, tree, options);                              // 47
      if (isFunction(onInstall)) {                                         // 48
        onInstall(rootRequire);                                            // 49
      }                                                                    // 50
    }                                                                      // 51
    return rootRequire;                                                    // 52
  }                                                                        // 53
                                                                           // 54
  function getOwn(obj, key) {                                              // 55
    return hasOwn.call(obj, key) && obj[key];                              // 56
  }                                                                        // 57
                                                                           // 58
  function isObject(value) {                                               // 59
    return value && typeof value === "object";                             // 60
  }                                                                        // 61
                                                                           // 62
  function isFunction(value) {                                             // 63
    return typeof value === "function";                                    // 64
  }                                                                        // 65
                                                                           // 66
  function isString(value) {                                               // 67
    return typeof value === "string";                                      // 68
  }                                                                        // 69
                                                                           // 70
  function makeRequire(file) {                                             // 71
    function require(id) {                                                 // 72
      var result = fileResolve(file, id);                                  // 73
      if (result) {                                                        // 74
        return fileEvaluate(result, file.m);                               // 75
      }                                                                    // 76
                                                                           // 77
      var error = new Error("Cannot find module '" + id + "'");            // 78
                                                                           // 79
      if (isFunction(fallback)) {                                          // 80
        return fallback(                                                   // 81
          id, // The missing module identifier.                            // 82
          file.m.id, // The path of the requiring file.                    // 83
          error // The error we would have thrown.                         // 84
        );                                                                 // 85
      }                                                                    // 86
                                                                           // 87
      throw error;                                                         // 88
    }                                                                      // 89
                                                                           // 90
    require.resolve = function (id) {                                      // 91
      var f = fileResolve(file, id);                                       // 92
      if (f) return f.m.id;                                                // 93
      throw new Error("Cannot find module '" + id + "'");                  // 94
    };                                                                     // 95
                                                                           // 96
    return require;                                                        // 97
  }                                                                        // 98
                                                                           // 99
  // File objects represent either directories or modules that have been   // 100
  // installed. When a `File` respresents a directory, its `.c` (contents)
  // property is an object containing the names of the files (or           // 102
  // directories) that it contains. When a `File` represents a module, its
  // `.c` property is a function that can be invoked with the appropriate  // 104
  // `(require, exports, module)` arguments to evaluate the module. If the
  // `.c` property is a string, that string will be resolved as a module   // 106
  // identifier, and the exports of the resulting module will provide the  // 107
  // exports of the original file. The `.p` (parent) property of a File is
  // either a directory `File` or `null`. Note that a child may claim      // 109
  // another `File` as its parent even if the parent does not have an      // 110
  // entry for that child in its `.c` object.  This is important for       // 111
  // implementing anonymous files, and preventing child modules from using
  // `../relative/identifier` syntax to examine unrelated modules.         // 113
  function File(name, parent) {                                            // 114
    var file = this;                                                       // 115
                                                                           // 116
    // Link to the parent file.                                            // 117
    file.p = parent = parent || null;                                      // 118
                                                                           // 119
    // The module object for this File, which will eventually boast an     // 120
    // .exports property when/if the file is evaluated.                    // 121
    file.m = new Module(name);                                             // 122
  }                                                                        // 123
                                                                           // 124
  function fileEvaluate(file, parentModule) {                              // 125
    var contents = file && file.c;                                         // 126
    var module = file.m;                                                   // 127
    if (! hasOwn.call(module, "exports")) {                                // 128
      if (parentModule) {                                                  // 129
        module.parent = parentModule;                                      // 130
        var children = parentModule.children;                              // 131
        if (Array.isArray(children)) {                                     // 132
          children.push(module);                                           // 133
        }                                                                  // 134
      }                                                                    // 135
                                                                           // 136
      // If a Module.prototype.useNode method is defined, give it a chance
      // to define module.exports based on module.id using Node.           // 138
      if (! isFunction(module.useNode) ||                                  // 139
          ! module.useNode()) {                                            // 140
        contents(                                                          // 141
          file.r = file.r || makeRequire(file),                            // 142
          module.exports = {},                                             // 143
          module,                                                          // 144
          file.m.id,                                                       // 145
          file.p.m.id                                                      // 146
        );                                                                 // 147
      }                                                                    // 148
    }                                                                      // 149
    return module.exports;                                                 // 150
  }                                                                        // 151
                                                                           // 152
  function fileIsDirectory(file) {                                         // 153
    return file && isObject(file.c);                                       // 154
  }                                                                        // 155
                                                                           // 156
  function fileMergeContents(file, contents, options) {                    // 157
    // If contents is an array of strings and functions, return the last   // 158
    // function with a `.d` property containing all the strings.           // 159
    if (Array.isArray(contents)) {                                         // 160
      var deps = [];                                                       // 161
                                                                           // 162
      contents.forEach(function (item) {                                   // 163
        if (isString(item)) {                                              // 164
          deps.push(item);                                                 // 165
        } else if (isFunction(item)) {                                     // 166
          contents = item;                                                 // 167
        }                                                                  // 168
      });                                                                  // 169
                                                                           // 170
      if (isFunction(contents)) {                                          // 171
        contents.d = deps;                                                 // 172
      } else {                                                             // 173
        // If the array did not contain a function, merge nothing.         // 174
        contents = null;                                                   // 175
      }                                                                    // 176
                                                                           // 177
    } else if (isFunction(contents)) {                                     // 178
      // If contents is already a function, make sure it has `.d`.         // 179
      contents.d = contents.d || [];                                       // 180
                                                                           // 181
    } else if (! isString(contents) &&                                     // 182
               ! isObject(contents)) {                                     // 183
      // If contents is neither an array nor a function nor a string nor   // 184
      // an object, just give up and merge nothing.                        // 185
      contents = null;                                                     // 186
    }                                                                      // 187
                                                                           // 188
    if (contents) {                                                        // 189
      file.c = file.c || (isObject(contents) ? {} : contents);             // 190
      if (isObject(contents) && fileIsDirectory(file)) {                   // 191
        Object.keys(contents).forEach(function (key) {                     // 192
          if (key === "..") {                                              // 193
            child = file.p;                                                // 194
                                                                           // 195
          } else {                                                         // 196
            var child = getOwn(file.c, key);                               // 197
            if (! child) {                                                 // 198
              child = file.c[key] = new File(                              // 199
                file.m.id.replace(/\/*$/, "/") + key,                      // 200
                file                                                       // 201
              );                                                           // 202
                                                                           // 203
              child.o = options;                                           // 204
            }                                                              // 205
          }                                                                // 206
                                                                           // 207
          fileMergeContents(child, contents[key], options);                // 208
        });                                                                // 209
      }                                                                    // 210
    }                                                                      // 211
  }                                                                        // 212
                                                                           // 213
  function fileGetExtensions(file) {                                       // 214
    return file.o && file.o.extensions || defaultExtensions;               // 215
  }                                                                        // 216
                                                                           // 217
  function fileAppendIdPart(file, part, extensions) {                      // 218
    // Always append relative to a directory.                              // 219
    while (file && ! fileIsDirectory(file)) {                              // 220
      file = file.p;                                                       // 221
    }                                                                      // 222
                                                                           // 223
    if (! file || ! part || part === ".") {                                // 224
      return file;                                                         // 225
    }                                                                      // 226
                                                                           // 227
    if (part === "..") {                                                   // 228
      return file.p;                                                       // 229
    }                                                                      // 230
                                                                           // 231
    var exactChild = getOwn(file.c, part);                                 // 232
                                                                           // 233
    // Only consider multiple file extensions if this part is the last     // 234
    // part of a module identifier and not equal to `.` or `..`, and there
    // was no exact match or the exact match was a directory.              // 236
    if (extensions && (! exactChild || fileIsDirectory(exactChild))) {     // 237
      for (var e = 0; e < extensions.length; ++e) {                        // 238
        var child = getOwn(file.c, part + extensions[e]);                  // 239
        if (child) {                                                       // 240
          return child;                                                    // 241
        }                                                                  // 242
      }                                                                    // 243
    }                                                                      // 244
                                                                           // 245
    return exactChild;                                                     // 246
  }                                                                        // 247
                                                                           // 248
  function fileAppendId(file, id, extensions) {                            // 249
    var parts = id.split("/");                                             // 250
                                                                           // 251
    // Use `Array.prototype.every` to terminate iteration early if         // 252
    // `fileAppendIdPart` returns a falsy value.                           // 253
    parts.every(function (part, i) {                                       // 254
      return file = i < parts.length - 1                                   // 255
        ? fileAppendIdPart(file, part)                                     // 256
        : fileAppendIdPart(file, part, extensions);                        // 257
    });                                                                    // 258
                                                                           // 259
    return file;                                                           // 260
  }                                                                        // 261
                                                                           // 262
  function fileResolve(file, id, seenDirFiles) {                           // 263
    var extensions = fileGetExtensions(file);                              // 264
                                                                           // 265
    file =                                                                 // 266
      // Absolute module identifiers (i.e. those that begin with a `/`     // 267
      // character) are interpreted relative to the root directory, which  // 268
      // is a slight deviation from Node, which has access to the entire   // 269
      // file system.                                                      // 270
      id.charAt(0) === "/" ? fileAppendId(root, id, extensions) :          // 271
      // Relative module identifiers are interpreted relative to the       // 272
      // current file, naturally.                                          // 273
      id.charAt(0) === "." ? fileAppendId(file, id, extensions) :          // 274
      // Top-level module identifiers are interpreted as referring to      // 275
      // packages in `node_modules` directories.                           // 276
      nodeModulesLookup(file, id, extensions);                             // 277
                                                                           // 278
    // If the identifier resolves to a directory, we use the same logic as
    // Node to find an `index.js` or `package.json` file to evaluate.      // 280
    while (fileIsDirectory(file)) {                                        // 281
      seenDirFiles = seenDirFiles || [];                                   // 282
                                                                           // 283
      // If the "main" field of a `package.json` file resolves to a        // 284
      // directory we've already considered, then we should not attempt to
      // read the same `package.json` file again. Using an array as a set  // 286
      // is acceptable here because the number of directories to consider  // 287
      // is rarely greater than 1 or 2. Also, using indexOf allows us to   // 288
      // store File objects instead of strings.                            // 289
      if (seenDirFiles.indexOf(file) < 0) {                                // 290
        seenDirFiles.push(file);                                           // 291
                                                                           // 292
        var pkgJsonFile = fileAppendIdPart(file, "package.json");          // 293
        var main = pkgJsonFile && fileEvaluate(pkgJsonFile).main;          // 294
        if (isString(main)) {                                              // 295
          // The "main" field of package.json does not have to begin with  // 296
          // ./ to be considered relative, so first we try simply          // 297
          // appending it to the directory path before falling back to a   // 298
          // full fileResolve, which might return a package from a         // 299
          // node_modules directory.                                       // 300
          file = fileAppendId(file, main, extensions) ||                   // 301
            fileResolve(file, main, seenDirFiles);                         // 302
                                                                           // 303
          if (file) {                                                      // 304
            // The fileAppendId call above may have returned a directory,  // 305
            // so continue the loop to make sure we resolve it to a        // 306
            // non-directory file.                                         // 307
            continue;                                                      // 308
          }                                                                // 309
        }                                                                  // 310
      }                                                                    // 311
                                                                           // 312
      // If we didn't find a `package.json` file, or it didn't have a      // 313
      // resolvable `.main` property, the only possibility left to         // 314
      // consider is that this directory contains an `index.js` module.    // 315
      // This assignment almost always terminates the while loop, because  // 316
      // there's very little chance `fileIsDirectory(file)` will be true   // 317
      // for the result of `fileAppendIdPart(file, "index.js")`. However,  // 318
      // in principle it is remotely possible that a file called           // 319
      // `index.js` could be a directory instead of a file.                // 320
      file = fileAppendIdPart(file, "index.js");                           // 321
    }                                                                      // 322
                                                                           // 323
    if (file && isString(file.c)) {                                        // 324
      file = fileResolve(file, file.c, seenDirFiles);                      // 325
    }                                                                      // 326
                                                                           // 327
    return file;                                                           // 328
  };                                                                       // 329
                                                                           // 330
  function nodeModulesLookup(file, id, extensions) {                       // 331
    if (isFunction(override)) {                                            // 332
      id = override(id, file.m.id);                                        // 333
    }                                                                      // 334
                                                                           // 335
    if (isString(id)) {                                                    // 336
      for (var resolved; file && ! resolved; file = file.p) {              // 337
        resolved = fileIsDirectory(file) &&                                // 338
          fileAppendId(file, "node_modules/" + id, extensions);            // 339
      }                                                                    // 340
                                                                           // 341
      return resolved;                                                     // 342
    }                                                                      // 343
  }                                                                        // 344
                                                                           // 345
  return install;                                                          // 346
};                                                                         // 347
                                                                           // 348
if (typeof exports === "object") {                                         // 349
  exports.makeInstaller = makeInstaller;                                   // 350
}                                                                          // 351
                                                                           // 352
/////////////////////////////////////////////////////////////////////////////







(function(){

/////////////////////////////////////////////////////////////////////////////
//                                                                         //
// packages/modules-runtime/modules-runtime.js                             //
//                                                                         //
/////////////////////////////////////////////////////////////////////////////
                                                                           //
var options = {};                                                          // 1
var hasOwn = options.hasOwnProperty;                                       // 2
                                                                           // 3
// RegExp matching strings that don't start with a `.` or a `/`.           // 4
var topLevelIdPattern = /^[^./]/;                                          // 5
                                                                           // 6
// This function will be called whenever a module identifier that hasn't   // 7
// been installed is required. For backwards compatibility, and so that we
// can require binary dependencies on the server, we implement the         // 9
// fallback in terms of Npm.require.                                       // 10
options.fallback = function (id, dir, error) {                             // 11
  // For simplicity, we honor only top-level module identifiers here.      // 12
  // We could try to honor relative and absolute module identifiers by     // 13
  // somehow combining `id` with `dir`, but we'd have to be really careful
  // that the resulting modules were located in a known directory (not     // 15
  // some arbitrary location on the file system), and we only really need  // 16
  // the fallback for dependencies installed in node_modules directories.  // 17
  if (topLevelIdPattern.test(id)) {                                        // 18
    if (typeof Npm === "object" &&                                         // 19
        typeof Npm.require === "function") {                               // 20
      return Npm.require(id);                                              // 21
    }                                                                      // 22
  }                                                                        // 23
                                                                           // 24
  throw error;                                                             // 25
};                                                                         // 26
                                                                           // 27
if (Meteor.isServer) {                                                     // 28
  // Defining Module.prototype.useNode allows the module system to         // 29
  // delegate evaluation to Node, unless useNode returns false.            // 30
  (options.Module = function Module(id) {                                  // 31
    // Same as the default Module constructor implementation.              // 32
    this.id = id;                                                          // 33
    this.children = [];                                                    // 34
  }).prototype.useNode = function () {                                     // 35
    if (typeof npmRequire !== "function") {                                // 36
      // Can't use Node if npmRequire is not defined.                      // 37
      return false;                                                        // 38
    }                                                                      // 39
                                                                           // 40
    var parts = this.id.split("/");                                        // 41
    var start = 0;                                                         // 42
    if (parts[start] === "") ++start;                                      // 43
    if (parts[start] === "node_modules" &&                                 // 44
        parts[start + 1] === "meteor") {                                   // 45
      start += 2;                                                          // 46
    }                                                                      // 47
                                                                           // 48
    if (parts.indexOf("node_modules", start) < 0) {                        // 49
      // Don't try to use Node for modules that aren't in node_modules     // 50
      // directories.                                                      // 51
      return false;                                                        // 52
    }                                                                      // 53
                                                                           // 54
    try {                                                                  // 55
      npmRequire.resolve(this.id);                                         // 56
    } catch (e) {                                                          // 57
      return false;                                                        // 58
    }                                                                      // 59
                                                                           // 60
    this.exports = npmRequire(this.id);                                    // 61
                                                                           // 62
    return true;                                                           // 63
  };                                                                       // 64
}                                                                          // 65
                                                                           // 66
meteorInstall = makeInstaller(options);                                    // 67
                                                                           // 68
/////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['modules-runtime'] = {}, {
  meteorInstall: meteorInstall
});

})();
