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

/* Package-scope variables */
var URL, buildUrl;

(function(){

/////////////////////////////////////////////////////////////////////////////////////
//                                                                                 //
// packages/url/url_common.js                                                      //
//                                                                                 //
/////////////////////////////////////////////////////////////////////////////////////
                                                                                   //
URL = {};                                                                          // 1
                                                                                   // 2
var encodeString = function(str) {                                                 // 3
  return encodeURIComponent(str).replace(/[!'()]/g, escape).replace(/\*/g, "%2A");
};                                                                                 // 5
                                                                                   // 6
                                                                                   // 7
URL._encodeParams = function(params) {                                             // 8
  var buf = [];                                                                    // 9
  _.each(params, function(value, key) {                                            // 10
    if (buf.length)                                                                // 11
      buf.push('&');                                                               // 12
    buf.push(encodeString(key), '=', encodeString(value));                         // 13
  });                                                                              // 14
  return buf.join('').replace(/%20/g, '+');                                        // 15
};                                                                                 // 16
                                                                                   // 17
                                                                                   // 18
buildUrl = function(before_qmark, from_qmark, opt_query, opt_params) {             // 19
  var url_without_query = before_qmark;                                            // 20
  var query = from_qmark ? from_qmark.slice(1) : null;                             // 21
                                                                                   // 22
  if (typeof opt_query === "string")                                               // 23
    query = String(opt_query);                                                     // 24
                                                                                   // 25
  if (opt_params) {                                                                // 26
    query = query || "";                                                           // 27
    var prms = URL._encodeParams(opt_params);                                      // 28
    if (query && prms)                                                             // 29
      query += '&';                                                                // 30
    query += prms;                                                                 // 31
  }                                                                                // 32
                                                                                   // 33
  var url = url_without_query;                                                     // 34
  if (query !== null)                                                              // 35
    url += ("?"+query);                                                            // 36
                                                                                   // 37
  return url;                                                                      // 38
};                                                                                 // 39
                                                                                   // 40
/////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

/////////////////////////////////////////////////////////////////////////////////////
//                                                                                 //
// packages/url/url_client.js                                                      //
//                                                                                 //
/////////////////////////////////////////////////////////////////////////////////////
                                                                                   //
URL._constructUrl = function (url, query, params) {                                // 1
  var query_match = /^(.*?)(\?.*)?$/.exec(url);                                    // 2
  return buildUrl(query_match[1], query_match[2], query, params);                  // 3
};                                                                                 // 4
                                                                                   // 5
/////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.url = {}, {
  URL: URL
});

})();
