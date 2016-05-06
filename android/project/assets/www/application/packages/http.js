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
var URL = Package.url.URL;
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;
var Symbol = Package['ecmascript-runtime'].Symbol;
var Map = Package['ecmascript-runtime'].Map;
var Set = Package['ecmascript-runtime'].Set;
var meteorBabelHelpers = Package['babel-runtime'].meteorBabelHelpers;
var Promise = Package.promise.Promise;

/* Package-scope variables */
var makeErrorByStatus, populateData, HTTP;

var require = meteorInstall({"node_modules":{"meteor":{"http":{"httpcall_common.js":function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/http/httpcall_common.js                                                                                   //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
var MAX_LENGTH = 500; // if you change this, also change the appropriate test                                         // 1
                                                                                                                      //
makeErrorByStatus = function makeErrorByStatus(statusCode, content) {                                                 // 3
  var message = 'failed [' + statusCode + ']';                                                                        // 4
                                                                                                                      //
  if (content) {                                                                                                      // 6
    var stringContent = typeof content == "string" ? content : content.toString();                                    // 7
                                                                                                                      //
    message += ' ' + truncate(stringContent.replace(/\n/g, ' '), MAX_LENGTH);                                         // 10
  }                                                                                                                   //
                                                                                                                      //
  return new Error(message);                                                                                          // 13
};                                                                                                                    //
                                                                                                                      //
function truncate(str, length) {                                                                                      // 16
  return str.length > length ? str.slice(0, length) + '...' : str;                                                    // 17
}                                                                                                                     //
                                                                                                                      //
// Fill in `response.data` if the content-type is JSON.                                                               //
populateData = function populateData(response) {                                                                      // 21
  // Read Content-Type header, up to a ';' if there is one.                                                           //
  // A typical header might be "application/json; charset=utf-8"                                                      //
  // or just "application/json".                                                                                      //
  var contentType = (response.headers['content-type'] || ';').split(';')[0];                                          // 25
                                                                                                                      //
  // Only try to parse data as JSON if server sets correct content type.                                              //
  if (_.include(['application/json', 'text/javascript', 'application/javascript', 'application/x-javascript'], contentType)) {
    try {                                                                                                             // 30
      response.data = JSON.parse(response.content);                                                                   // 31
    } catch (err) {                                                                                                   //
      response.data = null;                                                                                           // 33
    }                                                                                                                 //
  } else {                                                                                                            //
    response.data = null;                                                                                             // 36
  }                                                                                                                   //
};                                                                                                                    //
                                                                                                                      //
HTTP = {};                                                                                                            // 40
                                                                                                                      //
/**                                                                                                                   //
 * @summary Send an HTTP `GET` request. Equivalent to calling [`HTTP.call`](#http_call) with "GET" as the first argument.
 * @param {String} url The URL to which the request should be sent.                                                   //
 * @param {Object} [callOptions] Options passed on to [`HTTP.call`](#http_call).                                      //
 * @param {Function} [asyncCallback] Callback that is called when the request is completed. Required on the client.   //
 * @locus Anywhere                                                                                                    //
 */                                                                                                                   //
HTTP.get = function () /* varargs */{                                                                                 // 49
  return HTTP.call.apply(this, ["GET"].concat(_.toArray(arguments)));                                                 // 50
};                                                                                                                    //
                                                                                                                      //
/**                                                                                                                   //
 * @summary Send an HTTP `POST` request. Equivalent to calling [`HTTP.call`](#http_call) with "POST" as the first argument.
 * @param {String} url The URL to which the request should be sent.                                                   //
 * @param {Object} [callOptions] Options passed on to [`HTTP.call`](#http_call).                                      //
 * @param {Function} [asyncCallback] Callback that is called when the request is completed. Required on the client.   //
 * @locus Anywhere                                                                                                    //
 */                                                                                                                   //
HTTP.post = function () /* varargs */{                                                                                // 60
  return HTTP.call.apply(this, ["POST"].concat(_.toArray(arguments)));                                                // 61
};                                                                                                                    //
                                                                                                                      //
/**                                                                                                                   //
 * @summary Send an HTTP `PUT` request. Equivalent to calling [`HTTP.call`](#http_call) with "PUT" as the first argument.
 * @param {String} url The URL to which the request should be sent.                                                   //
 * @param {Object} [callOptions] Options passed on to [`HTTP.call`](#http_call).                                      //
 * @param {Function} [asyncCallback] Callback that is called when the request is completed. Required on the client.   //
 * @locus Anywhere                                                                                                    //
 */                                                                                                                   //
HTTP.put = function () /* varargs */{                                                                                 // 71
  return HTTP.call.apply(this, ["PUT"].concat(_.toArray(arguments)));                                                 // 72
};                                                                                                                    //
                                                                                                                      //
/**                                                                                                                   //
 * @summary Send an HTTP `DELETE` request. Equivalent to calling [`HTTP.call`](#http_call) with "DELETE" as the first argument. (Named `del` to avoid conflict with the Javascript keyword `delete`)
 * @param {String} url The URL to which the request should be sent.                                                   //
 * @param {Object} [callOptions] Options passed on to [`HTTP.call`](#http_call).                                      //
 * @param {Function} [asyncCallback] Callback that is called when the request is completed. Required on the client.   //
 * @locus Anywhere                                                                                                    //
 */                                                                                                                   //
HTTP.del = function () /* varargs */{                                                                                 // 82
  return HTTP.call.apply(this, ["DELETE"].concat(_.toArray(arguments)));                                              // 83
};                                                                                                                    //
                                                                                                                      //
/**                                                                                                                   //
 * @summary Send an HTTP `PATCH` request. Equivalent to calling [`HTTP.call`](#http_call) with "PATCH" as the first argument.
 * @param {String} url The URL to which the request should be sent.                                                   //
 * @param {Object} [callOptions] Options passed on to [`HTTP.call`](#http_call).                                      //
 * @param {Function} [asyncCallback] Callback that is called when the request is completed. Required on the client.   //
 * @locus Anywhere                                                                                                    //
 */                                                                                                                   //
HTTP.patch = function () /* varargs */{                                                                               // 93
  return HTTP.call.apply(this, ["PATCH"].concat(_.toArray(arguments)));                                               // 94
};                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"httpcall_client.js":function(require,exports,module){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/http/httpcall_client.js                                                                                   //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
/**                                                                                                                   //
 * @summary Perform an outbound HTTP request.                                                                         //
 * @locus Anywhere                                                                                                    //
 * @param {String} method The [HTTP method](http://en.wikipedia.org/wiki/HTTP_method) to use, such as "`GET`", "`POST`", or "`HEAD`".
 * @param {String} url The URL to retrieve.                                                                           //
 * @param {Object} [options]                                                                                          //
 * @param {String} options.content String to use as the HTTP request body.                                            //
 * @param {Object} options.data JSON-able object to stringify and use as the HTTP request body. Overwrites `content`.
 * @param {String} options.query Query string to go in the URL. Overwrites any query string in `url`.                 //
 * @param {Object} options.params Dictionary of request parameters to be encoded and placed in the URL (for GETs) or request body (for POSTs).  If `content` or `data` is specified, `params` will always be placed in the URL.
 * @param {String} options.auth HTTP basic authentication string of the form `"username:password"`                    //
 * @param {Object} options.headers Dictionary of strings, headers to add to the HTTP request.                         //
 * @param {Number} options.timeout Maximum time in milliseconds to wait for the request before failing.  There is no timeout by default.
 * @param {Boolean} options.followRedirects If `true`, transparently follow HTTP redirects. Cannot be set to `false` on the client. Default `true`.
 * @param {Object} options.npmRequestOptions On the server, `HTTP.call` is implemented by using the [npm `request` module](https://www.npmjs.com/package/request). Any options in this object will be passed directly to the `request` invocation.
 * @param {Function} options.beforeSend On the client, this will be called before the request is sent to allow for more direct manipulation of the underlying XMLHttpRequest object, which will be passed as the first argument. If the callback returns `false`, the request will be not be send.
 * @param {Function} [asyncCallback] Optional callback.  If passed, the method runs asynchronously, instead of synchronously, and calls asyncCallback.  On the client, this callback is required.
 */                                                                                                                   //
HTTP.call = function (method, url, options, callback) {                                                               // 19
                                                                                                                      //
  ////////// Process arguments //////////                                                                             //
                                                                                                                      //
  if (!callback && typeof options === "function") {                                                                   // 23
    // support (method, url, callback) argument list                                                                  //
    callback = options;                                                                                               // 25
    options = null;                                                                                                   // 26
  }                                                                                                                   //
                                                                                                                      //
  options = options || {};                                                                                            // 29
                                                                                                                      //
  if (typeof callback !== "function") throw new Error("Can't make a blocking HTTP call from the client; callback required.");
                                                                                                                      //
  method = (method || "").toUpperCase();                                                                              // 35
                                                                                                                      //
  var headers = {};                                                                                                   // 37
                                                                                                                      //
  var content = options.content;                                                                                      // 39
  if (options.data) {                                                                                                 // 40
    content = JSON.stringify(options.data);                                                                           // 41
    headers['Content-Type'] = 'application/json';                                                                     // 42
  }                                                                                                                   //
                                                                                                                      //
  var params_for_url, params_for_body;                                                                                // 45
  if (content || method === "GET" || method === "HEAD") params_for_url = options.params;else params_for_body = options.params;
                                                                                                                      //
  url = URL._constructUrl(url, options.query, params_for_url);                                                        // 51
                                                                                                                      //
  if (options.followRedirects === false) throw new Error("Option followRedirects:false not supported on client.");    // 53
                                                                                                                      //
  if (_.has(options, 'npmRequestOptions')) {                                                                          // 56
    throw new Error("Option npmRequestOptions not supported on client.");                                             // 57
  }                                                                                                                   //
                                                                                                                      //
  var username, password;                                                                                             // 60
  if (options.auth) {                                                                                                 // 61
    var colonLoc = options.auth.indexOf(':');                                                                         // 62
    if (colonLoc < 0) throw new Error('auth option should be of the form "username:password"');                       // 63
    username = options.auth.substring(0, colonLoc);                                                                   // 65
    password = options.auth.substring(colonLoc + 1);                                                                  // 66
  }                                                                                                                   //
                                                                                                                      //
  if (params_for_body) {                                                                                              // 69
    content = URL._encodeParams(params_for_body);                                                                     // 70
  }                                                                                                                   //
                                                                                                                      //
  _.extend(headers, options.headers || {});                                                                           // 73
                                                                                                                      //
  ////////// Callback wrapping //////////                                                                             //
                                                                                                                      //
  // wrap callback to add a 'response' property on an error, in case                                                  //
  // we have both (http 4xx/5xx error, which has a response payload)                                                  //
  callback = function (callback) {                                                                                    // 19
    return function (error, response) {                                                                               // 80
      if (error && response) error.response = response;                                                               // 81
      callback(error, response);                                                                                      // 83
    };                                                                                                                //
  }(callback);                                                                                                        //
                                                                                                                      //
  // safety belt: only call the callback once.                                                                        //
  callback = _.once(callback);                                                                                        // 19
                                                                                                                      //
  ////////// Kickoff! //////////                                                                                      //
                                                                                                                      //
  // from this point on, errors are because of something remote, not                                                  //
  // something we should check in advance. Turn exceptions into error                                                 //
  // results.                                                                                                         //
  try {                                                                                                               // 19
    // setup XHR object                                                                                               //
    var xhr;                                                                                                          // 98
    if (typeof XMLHttpRequest !== "undefined") xhr = new XMLHttpRequest();else if (typeof ActiveXObject !== "undefined") xhr = new ActiveXObject("Microsoft.XMLHttp"); // IE6
    else throw new Error("Can't create XMLHttpRequest"); // ???                                                       // 101
                                                                                                                      //
    xhr.open(method, url, true, username, password);                                                                  // 96
                                                                                                                      //
    for (var k in meteorBabelHelpers.sanitizeForInObject(headers)) {                                                  // 108
      xhr.setRequestHeader(k, headers[k]);                                                                            // 109
    } // setup timeout                                                                                                //
    var timed_out = false;                                                                                            // 96
    var timer;                                                                                                        // 114
    if (options.timeout) {                                                                                            // 115
      timer = Meteor.setTimeout(function () {                                                                         // 116
        timed_out = true;                                                                                             // 117
        xhr.abort();                                                                                                  // 118
      }, options.timeout);                                                                                            //
    };                                                                                                                //
                                                                                                                      //
    // callback on complete                                                                                           //
    xhr.onreadystatechange = function (evt) {                                                                         // 96
      if (xhr.readyState === 4) {                                                                                     // 124
        // COMPLETE                                                                                                   //
        if (timer) Meteor.clearTimeout(timer);                                                                        // 125
                                                                                                                      //
        if (timed_out) {                                                                                              // 128
          callback(new Error("timeout"));                                                                             // 129
        } else if (!xhr.status) {                                                                                     //
          // no HTTP response                                                                                         //
          callback(new Error("network"));                                                                             // 132
        } else {                                                                                                      //
                                                                                                                      //
          var response = {};                                                                                          // 135
          response.statusCode = xhr.status;                                                                           // 136
          response.content = xhr.responseText;                                                                        // 137
                                                                                                                      //
          response.headers = {};                                                                                      // 139
          var header_str = xhr.getAllResponseHeaders();                                                               // 140
                                                                                                                      //
          // https://github.com/meteor/meteor/issues/553                                                              //
          //                                                                                                          //
          // In Firefox there is a weird issue, sometimes                                                             //
          // getAllResponseHeaders returns the empty string, but                                                      //
          // getResponseHeader returns correct results. Possibly this                                                 //
          // issue:                                                                                                   //
          // https://bugzilla.mozilla.org/show_bug.cgi?id=608735                                                      //
          //                                                                                                          //
          // If this happens we can't get a full list of headers, but                                                 //
          // at least get content-type so our JSON decoding happens                                                   //
          // correctly. In theory, we could try and rescue more header                                                //
          // values with a list of common headers, but content-type is                                                //
          // the only vital one for now.                                                                              //
          if ("" === header_str && xhr.getResponseHeader("content-type")) header_str = "content-type: " + xhr.getResponseHeader("content-type");
                                                                                                                      //
          var headers_raw = header_str.split(/\r?\n/);                                                                // 159
          _.each(headers_raw, function (h) {                                                                          // 160
            var m = /^(.*?):(?:\s+)(.*)$/.exec(h);                                                                    // 161
            if (m && m.length === 3) response.headers[m[1].toLowerCase()] = m[2];                                     // 162
          });                                                                                                         //
                                                                                                                      //
          populateData(response);                                                                                     // 166
                                                                                                                      //
          var error = null;                                                                                           // 168
          if (response.statusCode >= 400) error = makeErrorByStatus(response.statusCode, response.content);           // 169
                                                                                                                      //
          callback(error, response);                                                                                  // 172
        }                                                                                                             //
      }                                                                                                               //
    };                                                                                                                //
                                                                                                                      //
    // Allow custom control over XHR and abort early.                                                                 //
    if (options.beforeSend) {                                                                                         // 96
      // Sanity                                                                                                       //
      var beforeSend = _.once(options.beforeSend);                                                                    // 180
                                                                                                                      //
      // Call the callback and check to see if the request was aborted                                                //
      if (false === beforeSend.call(null, xhr, options)) {                                                            // 178
        return xhr.abort();                                                                                           // 184
      }                                                                                                               //
    }                                                                                                                 //
                                                                                                                      //
    // send it on its way                                                                                             //
    xhr.send(content);                                                                                                // 96
  } catch (err) {                                                                                                     //
    callback(err);                                                                                                    // 192
  }                                                                                                                   //
};                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"deprecated.js":function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                    //
// packages/http/deprecated.js                                                                                        //
//                                                                                                                    //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                      //
// The HTTP object used to be called Meteor.http.                                                                     //
// XXX COMPAT WITH 0.6.4                                                                                              //
Meteor.http = HTTP;                                                                                                   // 3
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}},{"extensions":[".js",".json"]});
require("./node_modules/meteor/http/httpcall_common.js");
require("./node_modules/meteor/http/httpcall_client.js");
require("./node_modules/meteor/http/deprecated.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.http = {}, {
  HTTP: HTTP
});

})();
