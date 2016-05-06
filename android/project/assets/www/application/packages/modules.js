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
var meteorInstall = Package['modules-runtime'].meteorInstall;

/* Package-scope variables */
var Buffer, process;

var require = meteorInstall({"node_modules":{"meteor":{"modules":{"client.js":["./install-packages.js","./stubs.js","./buffer.js","./process.js","./css",function(require,exports){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/modules/client.js                                             //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
require("./install-packages.js");                                         // 1
require("./stubs.js");                                                    // 2
require("./buffer.js");                                                   // 3
require("./process.js");                                                  // 4
                                                                          // 5
exports.addStyles = require("./css").addStyles;                           // 6
                                                                          // 7
////////////////////////////////////////////////////////////////////////////

}],"buffer.js":["buffer",function(require){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/modules/buffer.js                                             //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
try {                                                                     // 1
  Buffer = global.Buffer || require("buffer").Buffer;                     // 2
} catch (noBuffer) {}                                                     // 3
                                                                          // 4
////////////////////////////////////////////////////////////////////////////

}],"css.js":function(require,exports){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/modules/css.js                                                //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
var doc = document;                                                       // 1
var head = doc.getElementsByTagName("head").item(0);                      // 2
                                                                          // 3
exports.addStyles = function (css) {                                      // 4
  var style = doc.createElement("style");                                 // 5
                                                                          // 6
  style.setAttribute("type", "text/css");                                 // 7
                                                                          // 8
  // https://msdn.microsoft.com/en-us/library/ms535871(v=vs.85).aspx      // 9
  var internetExplorerSheetObject =                                       // 10
    style.sheet || // Edge/IE11.                                          // 11
    style.styleSheet; // Older IEs.                                       // 12
                                                                          // 13
  if (internetExplorerSheetObject) {                                      // 14
    internetExplorerSheetObject.cssText = css;                            // 15
  } else {                                                                // 16
    style.appendChild(doc.createTextNode(css));                           // 17
  }                                                                       // 18
                                                                          // 19
  return head.appendChild(style);                                         // 20
};                                                                        // 21
                                                                          // 22
////////////////////////////////////////////////////////////////////////////

},"install-packages.js":function(require,exports,module){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/modules/install-packages.js                                   //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
function install(name) {                                                  // 1
  var meteorDir = {};                                                     // 2
                                                                          // 3
  // Given a package name <name>, install a stub module in the            // 4
  // /node_modules/meteor directory called <name>.js, so that             // 5
  // require.resolve("meteor/<name>") will always return                  // 6
  // /node_modules/meteor/<name>.js instead of something like             // 7
  // /node_modules/meteor/<name>/index.js, in the rare but possible event
  // that the package contains a file called index.js (#6590).            // 9
  meteorDir[name + ".js"] = function (r, e, module) {                     // 10
    module.exports = Package[name];                                       // 11
  };                                                                      // 12
                                                                          // 13
  meteorInstall({                                                         // 14
    node_modules: {                                                       // 15
      meteor: meteorDir                                                   // 16
    }                                                                     // 17
  });                                                                     // 18
}                                                                         // 19
                                                                          // 20
// This file will be modified during computeJsOutputFilesMap to include   // 21
// install(<name>) calls for every Meteor package.                        // 22
                                                                          // 23
install("underscore");                                                    // 24
install("meteor");                                                        // 25
install("meteor-base");                                                   // 26
install("mobile-experience");                                             // 27
install("babel-compiler");                                                // 28
install("ecmascript");                                                    // 29
install("base64");                                                        // 30
install("ejson");                                                         // 31
install("id-map");                                                        // 32
install("ordered-dict");                                                  // 33
install("tracker");                                                       // 34
install("modules-runtime");                                               // 35
install("modules");                                                       // 36
install("es5-shim");                                                      // 37
install("promise");                                                       // 38
install("ecmascript-runtime");                                            // 39
install("babel-runtime");                                                 // 40
install("random");                                                        // 41
install("mongo-id");                                                      // 42
install("diff-sequence");                                                 // 43
install("geojson-utils");                                                 // 44
install("minimongo");                                                     // 45
install("check");                                                         // 46
install("retry");                                                         // 47
install("ddp-common");                                                    // 48
install("reload");                                                        // 49
install("ddp-client");                                                    // 50
install("ddp");                                                           // 51
install("ddp-server");                                                    // 52
install("allow-deny");                                                    // 53
install("mongo");                                                         // 54
install("blaze-html-templates");                                          // 55
install("reactive-dict");                                                 // 56
install("session");                                                       // 57
install("jquery");                                                        // 58
install("twbs:bootstrap");                                                // 59
install("deps");                                                          // 60
install("htmljs");                                                        // 61
install("observe-sequence");                                              // 62
install("reactive-var");                                                  // 63
install("blaze");                                                         // 64
install("ui");                                                            // 65
install("spacebars");                                                     // 66
install("templating");                                                    // 67
install("iron:core");                                                     // 68
install("iron:dynamic-template");                                         // 69
install("iron:layout");                                                   // 70
install("iron:url");                                                      // 71
install("iron:middleware-stack");                                         // 72
install("iron:location");                                                 // 73
install("iron:controller");                                               // 74
install("iron:router");                                                   // 75
install("sacha:spin");                                                    // 76
install("npm-bcrypt");                                                    // 77
install("ddp-rate-limiter");                                              // 78
install("localstorage");                                                  // 79
install("callback-hook");                                                 // 80
install("accounts-base");                                                 // 81
install("sha");                                                           // 82
install("srp");                                                           // 83
install("accounts-password");                                             // 84
install("stylus");                                                        // 85
install("anti:i18n");                                                     // 86
install("ian:accounts-ui-bootstrap-3");                                   // 87
install("audit-argument-checks");                                         // 88
install("url");                                                           // 89
install("http");                                                          // 90
install("autopublish");                                                   // 91
install("browser-policy");                                                // 92
install("kit:sweetalert");                                                // 93
install("standard-minifier-css");                                         // 94
install("standard-minifier-js");                                          // 95
install("webapp");                                                        // 96
install("livedata");                                                      // 97
install("hot-code-push");                                                 // 98
install("fastclick");                                                     // 99
install("mobile-status-bar");                                             // 100
install("launch-screen");                                                 // 101
install("autoupdate");                                                    // 102
install("service-configuration");                                         // 103
                                                                          // 104
////////////////////////////////////////////////////////////////////////////

},"process.js":["process",function(require,exports,module){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/modules/process.js                                            //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
try {                                                                     // 1
  // The application can run `npm install process` to provide its own     // 2
  // process stub; otherwise this module will provide a partial stub.     // 3
  process = global.process || require("process");                         // 4
} catch (noProcess) {                                                     // 5
  process = {};                                                           // 6
}                                                                         // 7
                                                                          // 8
if (Meteor.isServer) {                                                    // 9
  // Make require("process") work on the server in all versions of Node.  // 10
  meteorInstall({                                                         // 11
    node_modules: {                                                       // 12
      "process.js": function (r, e, module) {                             // 13
        module.exports = process;                                         // 14
      }                                                                   // 15
    }                                                                     // 16
  });                                                                     // 17
} else {                                                                  // 18
  process.platform = "browser";                                           // 19
  process.nextTick = process.nextTick || Meteor._setImmediate;            // 20
}                                                                         // 21
                                                                          // 22
if (typeof process.env !== "object") {                                    // 23
  process.env = {};                                                       // 24
}                                                                         // 25
                                                                          // 26
_.extend(process.env, meteorEnv);                                         // 27
                                                                          // 28
////////////////////////////////////////////////////////////////////////////

}],"stubs.js":["meteor-node-stubs",function(require){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/modules/stubs.js                                              //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
try {                                                                     // 1
  // When meteor-node-stubs is installed in the application's root        // 2
  // node_modules directory, requiring it here installs aliases for stubs
  // for all Node built-in modules, such as fs, util, and http.           // 4
  require("meteor-node-stubs");                                           // 5
} catch (noStubs) {}                                                      // 6
                                                                          // 7
////////////////////////////////////////////////////////////////////////////

}]}}}},{"extensions":[".js",".json"]});
var exports = require("./node_modules/meteor/modules/client.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.modules = exports, {
  meteorInstall: meteorInstall,
  Buffer: Buffer,
  process: process
});

})();
