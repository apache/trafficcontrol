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
var Tracker = Package.tracker.Tracker;
var Deps = Package.tracker.Deps;
var Retry = Package.retry.Retry;
var DDP = Package['ddp-client'].DDP;
var Mongo = Package.mongo.Mongo;
var _ = Package.underscore._;
var HTTP = Package.http.HTTP;
var Random = Package.random.Random;

/* Package-scope variables */
var ClientVersions, Autoupdate;

(function(){

//////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                              //
// packages/autoupdate/autoupdate_cordova.js                                                    //
//                                                                                              //
//////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                //
var autoupdateVersionCordova = __meteor_runtime_config__.autoupdateVersionCordova || "unknown";
                                                                                                // 2
// The collection of acceptable client versions.                                                // 3
ClientVersions = new Mongo.Collection("meteor_autoupdate_clientVersions");                      // 4
                                                                                                // 5
Autoupdate = {};                                                                                // 6
                                                                                                // 7
Autoupdate.newClientAvailable = function() {                                                    // 8
  return !! ClientVersions.findOne({                                                            // 9
    _id: 'version-cordova',                                                                     // 10
    version: {$ne: autoupdateVersionCordova}                                                    // 11
  });                                                                                           // 12
};                                                                                              // 13
                                                                                                // 14
var retry = new Retry({                                                                         // 15
  // Unlike the stream reconnect use of Retry, which we want to be instant                      // 16
  // in normal operation, this is a wacky failure. We don't want to retry                       // 17
  // right away, we can start slowly.                                                           // 18
  //                                                                                            // 19
  // A better way than timeconstants here might be to use the knowledge                         // 20
  // of when we reconnect to help trigger these retries. Typically, the                         // 21
  // server fixing code will result in a restart and reconnect, but                             // 22
  // potentially the subscription could have a transient error.                                 // 23
  minCount: 0, // don't do any immediate retries                                                // 24
  baseTimeout: 30*1000 // start with 30s                                                        // 25
});                                                                                             // 26
var failures = 0;                                                                               // 27
                                                                                                // 28
Autoupdate._retrySubscription = function() {                                                    // 29
  var appId = __meteor_runtime_config__.appId;                                                  // 30
  Meteor.subscribe("meteor_autoupdate_clientVersions", appId, {                                 // 31
    onError: function(error) {                                                                  // 32
      console.log("autoupdate subscription failed:", error);                                    // 33
      failures++;                                                                               // 34
      retry.retryLater(failures, function() {                                                   // 35
        // Just retry making the subscription, don't reload the whole                           // 36
        // page. While reloading would catch more cases (for example,                           // 37
        // the server went back a version and is now doing old-style hot                        // 38
        // code push), it would also be more prone to reload loops,                             // 39
        // which look really bad to the user. Just retrying the                                 // 40
        // subscription over DDP means it is at least possible to fix by                        // 41
        // updating the server.                                                                 // 42
        Autoupdate._retrySubscription();                                                        // 43
      });                                                                                       // 44
    },                                                                                          // 45
    onReady: function() {                                                                       // 46
      if (Package.reload) {                                                                     // 47
        var checkNewVersionDocument = function(doc) {                                           // 48
          var self = this;                                                                      // 49
          if (doc.version !== autoupdateVersionCordova) {                                       // 50
            newVersionAvailable();                                                              // 51
          }                                                                                     // 52
        };                                                                                      // 53
                                                                                                // 54
        var handle = ClientVersions.find({_id: 'version-cordova'}).observe({                    // 55
          added: checkNewVersionDocument,                                                       // 56
          changed: checkNewVersionDocument                                                      // 57
        });                                                                                     // 58
      }                                                                                         // 59
    }                                                                                           // 60
  });                                                                                           // 61
};                                                                                              // 62
                                                                                                // 63
Meteor.startup(function() {                                                                     // 64
  WebAppLocalServer.onNewVersionReady(function() {                                              // 65
    if (Package.reload) {                                                                       // 66
      Package.reload.Reload._reload();                                                          // 67
    }                                                                                           // 68
  });                                                                                           // 69
                                                                                                // 70
  Autoupdate._retrySubscription();                                                              // 71
});                                                                                             // 72
                                                                                                // 73
var newVersionAvailable = function() {                                                          // 74
  WebAppLocalServer.checkForUpdates();                                                          // 75
}                                                                                               // 76
                                                                                                // 77
//////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.autoupdate = {}, {
  Autoupdate: Autoupdate
});

})();
