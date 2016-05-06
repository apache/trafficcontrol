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
var LaunchScreen;

(function(){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/launch-screen/mobile-launch-screen.js                         //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
// XXX This currently implements loading screens for mobile apps only,    // 1
// but in the future can be expanded to all apps.                         // 2
                                                                          // 3
var holdCount = 0;                                                        // 4
var alreadyHidden = false;                                                // 5
                                                                          // 6
LaunchScreen = {                                                          // 7
  hold: function () {                                                     // 8
    if (! Meteor.isCordova) {                                             // 9
      return {                                                            // 10
        release: function () { /* noop */ }                               // 11
      };                                                                  // 12
    }                                                                     // 13
                                                                          // 14
    if (alreadyHidden) {                                                  // 15
      throw new Error("Can't show launch screen once it's hidden");       // 16
    }                                                                     // 17
                                                                          // 18
    holdCount++;                                                          // 19
                                                                          // 20
    var released = false;                                                 // 21
    var release = function () {                                           // 22
      if (! Meteor.isCordova)                                             // 23
        return;                                                           // 24
                                                                          // 25
      if (! released) {                                                   // 26
        released = true;                                                  // 27
        holdCount--;                                                      // 28
        if (holdCount === 0 &&                                            // 29
            typeof navigator !== 'undefined' && navigator.splashscreen) {
          alreadyHidden = true;                                           // 31
          navigator.splashscreen.hide();                                  // 32
        }                                                                 // 33
      }                                                                   // 34
    };                                                                    // 35
                                                                          // 36
    // Returns a launch screen handle with a release method               // 37
    return {                                                              // 38
      release: release                                                    // 39
    };                                                                    // 40
  }                                                                       // 41
};                                                                        // 42
                                                                          // 43
////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////
//                                                                        //
// packages/launch-screen/default-behavior.js                             //
//                                                                        //
////////////////////////////////////////////////////////////////////////////
                                                                          //
// Hold launch screen on app load. This reflects the fact that Meteor     // 1
// mobile apps that use this package always start with a launch screen    // 2
// visible. (see XXX comment at the top of package.js for more            // 3
// details)                                                               // 4
var handle = LaunchScreen.hold();                                         // 5
                                                                          // 6
var Template = Package.templating && Package.templating.Template;         // 7
                                                                          // 8
Meteor.startup(function () {                                              // 9
  if (! Template) {                                                       // 10
    handle.release();                                                     // 11
  } else if (Package['iron:router']) {                                    // 12
    // XXX Instead of doing this here, this code should be in             // 13
    // iron:router directly. Note that since we're in a                   // 14
    // `Meteor.startup` block it's ok that we don't have a                // 15
    // weak dependency on iron:router in package.js.                      // 16
    Package['iron:router'].Router.onAfterAction(function () {             // 17
      handle.release();                                                   // 18
    });                                                                   // 19
  } else {                                                                // 20
    Template.body.onRendered(function () {                                // 21
      handle.release();                                                   // 22
    });                                                                   // 23
                                                                          // 24
    // In case `Template.body` never gets rendered (due to some bug),     // 25
    // hide the launch screen after 6 seconds. This matches the           // 26
    // observed timeout that Cordova apps on Android (but not iOS)        // 27
    // have on hiding the launch screen (even if you don't call           // 28
    // `navigator.splashscreen.hide()`)                                   // 29
    setTimeout(function () {                                              // 30
      handle.release();                                                   // 31
    }, 6000);                                                             // 32
  }                                                                       // 33
});                                                                       // 34
                                                                          // 35
////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['launch-screen'] = {}, {
  LaunchScreen: LaunchScreen
});

})();
