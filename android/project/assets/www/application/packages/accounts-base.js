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
var Tracker = Package.tracker.Tracker;
var Deps = Package.tracker.Deps;
var Random = Package.random.Random;
var Hook = Package['callback-hook'].Hook;
var DDP = Package['ddp-client'].DDP;
var Mongo = Package.mongo.Mongo;
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;
var Symbol = Package['ecmascript-runtime'].Symbol;
var Map = Package['ecmascript-runtime'].Map;
var Set = Package['ecmascript-runtime'].Set;
var meteorBabelHelpers = Package['babel-runtime'].meteorBabelHelpers;
var Promise = Package.promise.Promise;

/* Package-scope variables */
var Accounts, EXPIRE_TOKENS_INTERVAL_MS, CONNECTION_CLOSE_DELAY_MS;

var require = meteorInstall({"node_modules":{"meteor":{"accounts-base":{"client_main.js":["./accounts_client.js","./url_client.js","./localstorage_token.js",function(require,exports,module){

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                   //
// packages/accounts-base/client_main.js                                                                             //
//                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                     //
exports.__esModule = true;                                                                                           //
exports.AccountsTest = exports.AccountsClient = undefined;                                                           //
                                                                                                                     //
var _accounts_client = require("./accounts_client.js");                                                              // 1
                                                                                                                     //
var _url_client = require("./url_client.js");                                                                        // 2
                                                                                                                     //
require("./localstorage_token.js");                                                                                  // 3
                                                                                                                     //
/**                                                                                                                  //
 * @namespace Accounts                                                                                               //
 * @summary The namespace for all client-side accounts-related methods.                                              //
 */                                                                                                                  //
Accounts = new _accounts_client.AccountsClient();                                                                    // 9
                                                                                                                     //
/**                                                                                                                  //
 * @summary A [Mongo.Collection](#collections) containing user documents.                                            //
 * @locus Anywhere                                                                                                   //
 * @type {Mongo.Collection}                                                                                          //
 * @importFromPackage meteor                                                                                         //
 */                                                                                                                  //
Meteor.users = Accounts.users;                                                                                       // 17
                                                                                                                     //
exports.                                                                                                             //
// Since this file is the main module for the client version of the                                                  //
// accounts-base package, properties of non-entry-point modules need to                                              //
// be re-exported in order to be accessible to modules that import the                                               //
// accounts-base package.                                                                                            //
AccountsClient = _accounts_client.AccountsClient;                                                                    // 24
exports.AccountsTest = _url_client.AccountsTest;                                                                     //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"accounts_client.js":["babel-runtime/helpers/classCallCheck","babel-runtime/helpers/possibleConstructorReturn","babel-runtime/helpers/inherits","./accounts_common.js",function(require,exports){

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                   //
// packages/accounts-base/accounts_client.js                                                                         //
//                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                     //
exports.__esModule = true;                                                                                           //
exports.AccountsClient = undefined;                                                                                  //
                                                                                                                     //
var _classCallCheck2 = require("babel-runtime/helpers/classCallCheck");                                              //
                                                                                                                     //
var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);                                                     //
                                                                                                                     //
var _possibleConstructorReturn2 = require("babel-runtime/helpers/possibleConstructorReturn");                        //
                                                                                                                     //
var _possibleConstructorReturn3 = _interopRequireDefault(_possibleConstructorReturn2);                               //
                                                                                                                     //
var _inherits2 = require("babel-runtime/helpers/inherits");                                                          //
                                                                                                                     //
var _inherits3 = _interopRequireDefault(_inherits2);                                                                 //
                                                                                                                     //
var _accounts_common = require("./accounts_common.js");                                                              // 1
                                                                                                                     //
function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }                    //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Constructor for the `Accounts` object on the client.                                                     //
 * @locus Client                                                                                                     //
 * @class AccountsClient                                                                                             //
 * @extends AccountsCommon                                                                                           //
 * @instancename accountsClient                                                                                      //
 * @param {Object} options an object with fields:                                                                    //
 * @param {Object} options.connection Optional DDP connection to reuse.                                              //
 * @param {String} options.ddpUrl Optional URL for creating a new DDP connection.                                    //
 */                                                                                                                  //
                                                                                                                     //
var AccountsClient = exports.AccountsClient = function (_AccountsCommon) {                                           //
  (0, _inherits3["default"])(AccountsClient, _AccountsCommon);                                                       //
                                                                                                                     //
  function AccountsClient(options) {                                                                                 // 14
    (0, _classCallCheck3["default"])(this, AccountsClient);                                                          //
                                                                                                                     //
    var _this = (0, _possibleConstructorReturn3["default"])(this, _AccountsCommon.call(this, options));              //
                                                                                                                     //
    _this._loggingIn = false;                                                                                        // 17
    _this._loggingInDeps = new Tracker.Dependency();                                                                 // 18
                                                                                                                     //
    _this._loginServicesHandle = _this.connection.subscribe("meteor.loginServiceConfiguration");                     // 20
                                                                                                                     //
    _this._pageLoadLoginCallbacks = [];                                                                              // 23
    _this._pageLoadLoginAttemptInfo = null;                                                                          // 24
                                                                                                                     //
    // Defined in url_client.js.                                                                                     //
    _this._initUrlMatching();                                                                                        // 14
                                                                                                                     //
    // Defined in localstorage_token.js.                                                                             //
    _this._initLocalStorage();                                                                                       // 14
    return _this;                                                                                                    //
  }                                                                                                                  //
                                                                                                                     //
  ///                                                                                                                //
  /// CURRENT USER                                                                                                   //
  ///                                                                                                                //
                                                                                                                     //
  // @override                                                                                                       //
                                                                                                                     //
                                                                                                                     //
  AccountsClient.prototype.userId = function () {                                                                    // 13
    function userId() {                                                                                              //
      return this.connection.userId();                                                                               // 39
    }                                                                                                                //
                                                                                                                     //
    return userId;                                                                                                   //
  }();                                                                                                               //
                                                                                                                     //
  // This is mostly just called within this file, but Meteor.loginWithPassword                                       //
  // also uses it to make loggingIn() be true during the beginPasswordExchange                                       //
  // method call too.                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsClient.prototype._setLoggingIn = function () {                                                             // 13
    function _setLoggingIn(x) {                                                                                      //
      if (this._loggingIn !== x) {                                                                                   // 46
        this._loggingIn = x;                                                                                         // 47
        this._loggingInDeps.changed();                                                                               // 48
      }                                                                                                              //
    }                                                                                                                //
                                                                                                                     //
    return _setLoggingIn;                                                                                            //
  }();                                                                                                               //
                                                                                                                     //
  /**                                                                                                                //
   * @summary True if a login method (such as `Meteor.loginWithPassword`, `Meteor.loginWithFacebook`, or `Accounts.createUser`) is currently in progress. A reactive data source.
   * @locus Client                                                                                                   //
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsClient.prototype.loggingIn = function () {                                                                 // 13
    function loggingIn() {                                                                                           //
      this._loggingInDeps.depend();                                                                                  // 57
      return this._loggingIn;                                                                                        // 58
    }                                                                                                                //
                                                                                                                     //
    return loggingIn;                                                                                                //
  }();                                                                                                               //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Log the user out.                                                                                      //
   * @locus Client                                                                                                   //
   * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsClient.prototype.logout = function () {                                                                    // 13
    function logout(callback) {                                                                                      //
      var self = this;                                                                                               // 67
      self.connection.apply('logout', [], {                                                                          // 68
        wait: true                                                                                                   // 69
      }, function (error, result) {                                                                                  //
        if (error) {                                                                                                 // 71
          callback && callback(error);                                                                               // 72
        } else {                                                                                                     //
          self.makeClientLoggedOut();                                                                                // 74
          callback && callback();                                                                                    // 75
        }                                                                                                            //
      });                                                                                                            //
    }                                                                                                                //
                                                                                                                     //
    return logout;                                                                                                   //
  }();                                                                                                               //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Log out other clients logged in as the current user, but does not log out the client that calls this function.
   * @locus Client                                                                                                   //
   * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsClient.prototype.logoutOtherClients = function () {                                                        // 13
    function logoutOtherClients(callback) {                                                                          //
      var self = this;                                                                                               // 86
                                                                                                                     //
      // We need to make two method calls: one to replace our current token,                                         //
      // and another to remove all tokens except the current one. We want to                                         //
      // call these two methods one after the other, without any other                                               //
      // methods running between them. For example, we don't want `logout`                                           //
      // to be called in between our two method calls (otherwise the second                                          //
      // method call would return an error). Another example: we don't want                                          //
      // logout to be called before the callback for `getNewToken`;                                                  //
      // otherwise we would momentarily log the user out and then write a                                            //
      // new token to localStorage.                                                                                  //
      //                                                                                                             //
      // To accomplish this, we make both calls as wait methods, and queue                                           //
      // them one after the other, without spinning off the event loop in                                            //
      // between. Even though we queue `removeOtherTokens` before                                                    //
      // `getNewToken`, we won't actually send the `removeOtherTokens` call                                          //
      // until the `getNewToken` callback has finished running, because they                                         //
      // are both wait methods.                                                                                      //
      self.connection.apply('getNewToken', [], { wait: true }, function (err, result) {                              // 85
        if (!err) {                                                                                                  // 109
          self._storeLoginToken(self.userId(), result.token, result.tokenExpires);                                   // 110
        }                                                                                                            //
      });                                                                                                            //
                                                                                                                     //
      self.connection.apply('removeOtherTokens', [], { wait: true }, function (err) {                                // 119
        callback && callback(err);                                                                                   // 124
      });                                                                                                            //
    }                                                                                                                //
                                                                                                                     //
    return logoutOtherClients;                                                                                       //
  }();                                                                                                               //
                                                                                                                     //
  return AccountsClient;                                                                                             //
}(_accounts_common.AccountsCommon);                                                                                  //
                                                                                                                     //
;                                                                                                                    // 128
                                                                                                                     //
var Ap = AccountsClient.prototype;                                                                                   // 130
                                                                                                                     //
/**                                                                                                                  //
 * @summary True if a login method (such as `Meteor.loginWithPassword`, `Meteor.loginWithFacebook`, or `Accounts.createUser`) is currently in progress. A reactive data source.
 * @locus Client                                                                                                     //
 * @importFromPackage meteor                                                                                         //
 */                                                                                                                  //
Meteor.loggingIn = function () {                                                                                     // 137
  return Accounts.loggingIn();                                                                                       // 138
};                                                                                                                   //
                                                                                                                     //
///                                                                                                                  //
/// LOGIN METHODS                                                                                                    //
///                                                                                                                  //
                                                                                                                     //
// Call a login method on the server.                                                                                //
//                                                                                                                   //
// A login method is a method which on success calls `this.setUserId(id)` and                                        //
// `Accounts._setLoginToken` on the server and returns an object with fields                                         //
// 'id' (containing the user id), 'token' (containing a resume token), and                                           //
// optionally `tokenExpires`.                                                                                        //
//                                                                                                                   //
// This function takes care of:                                                                                      //
//   - Updating the Meteor.loggingIn() reactive data source                                                          //
//   - Calling the method in 'wait' mode                                                                             //
//   - On success, saving the resume token to localStorage                                                           //
//   - On success, calling Accounts.connection.setUserId()                                                           //
//   - Setting up an onReconnect handler which logs in with                                                          //
//     the resume token                                                                                              //
//                                                                                                                   //
// Options:                                                                                                          //
// - methodName: The method to call (default 'login')                                                                //
// - methodArguments: The arguments for the method                                                                   //
// - validateResult: If provided, will be called with the result of the                                              //
//                 method. If it throws, the client will not be logged in (and                                       //
//                 its error will be passed to the callback).                                                        //
// - userCallback: Will be called with no arguments once the user is fully                                           //
//                 logged in, or with the error on error.                                                            //
//                                                                                                                   //
Ap.callLoginMethod = function (options) {                                                                            // 169
  var self = this;                                                                                                   // 170
                                                                                                                     //
  options = _.extend({                                                                                               // 172
    methodName: 'login',                                                                                             // 173
    methodArguments: [{}],                                                                                           // 174
    _suppressLoggingIn: false                                                                                        // 175
  }, options);                                                                                                       //
                                                                                                                     //
  // Set defaults for callback arguments to no-op functions; make sure we                                            //
  // override falsey values too.                                                                                     //
  _.each(['validateResult', 'userCallback'], function (f) {                                                          // 169
    if (!options[f]) options[f] = function () {};                                                                    // 181
  });                                                                                                                //
                                                                                                                     //
  // Prepare callbacks: user provided and onLogin/onLoginFailure hooks.                                              //
  var loginCallbacks = _.once(function (error) {                                                                     // 169
    if (!error) {                                                                                                    // 187
      self._onLoginHook.each(function (callback) {                                                                   // 188
        callback();                                                                                                  // 189
        return true;                                                                                                 // 190
      });                                                                                                            //
    } else {                                                                                                         //
      self._onLoginFailureHook.each(function (callback) {                                                            // 193
        callback();                                                                                                  // 194
        return true;                                                                                                 // 195
      });                                                                                                            //
    }                                                                                                                //
    options.userCallback.apply(this, arguments);                                                                     // 198
  });                                                                                                                //
                                                                                                                     //
  var reconnected = false;                                                                                           // 201
                                                                                                                     //
  // We want to set up onReconnect as soon as we get a result token back from                                        //
  // the server, without having to wait for subscriptions to rerun. This is                                          //
  // because if we disconnect and reconnect between getting the result and                                           //
  // getting the results of subscription rerun, we WILL NOT re-send this                                             //
  // method (because we never re-send methods whose results we've received)                                          //
  // but we WILL call loggedInAndDataReadyCallback at "reconnect quiesce"                                            //
  // time. This will lead to makeClientLoggedIn(result.id) even though we                                            //
  // haven't actually sent a login method!                                                                           //
  //                                                                                                                 //
  // But by making sure that we send this "resume" login in that case (and                                           //
  // calling makeClientLoggedOut if it fails), we'll end up with an accurate                                         //
  // client-side userId. (It's important that livedata_connection guarantees                                         //
  // that the "reconnect quiesce"-time call to loggedInAndDataReadyCallback                                          //
  // will occur before the callback from the resume login call.)                                                     //
  var onResultReceived = function onResultReceived(err, result) {                                                    // 169
    if (err || !result || !result.token) {                                                                           // 218
      // Leave onReconnect alone if there was an error, so that if the user was                                      //
      // already logged in they will still get logged in on reconnect.                                               //
      // See issue #4970.                                                                                            //
    } else {                                                                                                         //
        self.connection.onReconnect = function () {                                                                  // 223
          reconnected = true;                                                                                        // 224
          // If our token was updated in storage, use the latest one.                                                //
          var storedToken = self._storedLoginToken();                                                                // 223
          if (storedToken) {                                                                                         // 227
            result = {                                                                                               // 228
              token: storedToken,                                                                                    // 229
              tokenExpires: self._storedLoginTokenExpires()                                                          // 230
            };                                                                                                       //
          }                                                                                                          //
          if (!result.tokenExpires) result.tokenExpires = self._tokenExpiration(new Date());                         // 233
          if (self._tokenExpiresSoon(result.tokenExpires)) {                                                         // 235
            self.makeClientLoggedOut();                                                                              // 236
          } else {                                                                                                   //
            self.callLoginMethod({                                                                                   // 238
              methodArguments: [{ resume: result.token }],                                                           // 239
              // Reconnect quiescence ensures that the user doesn't see an                                           //
              // intermediate state before the login method finishes. So we don't                                    //
              // need to show a logging-in animation.                                                                //
              _suppressLoggingIn: true,                                                                              // 243
              userCallback: function () {                                                                            // 244
                function userCallback(error) {                                                                       // 244
                  var storedTokenNow = self._storedLoginToken();                                                     // 245
                  if (error) {                                                                                       // 246
                    // If we had a login error AND the current stored token is the                                   //
                    // one that we tried to log in with, then declare ourselves                                      //
                    // logged out. If there's a token in storage but it's not the                                    //
                    // token that we tried to log in with, we don't know anything                                    //
                    // about whether that token is valid or not, so do nothing. The                                  //
                    // periodic localStorage poll will decide if we are logged in or                                 //
                    // out with this token, if it hasn't already. Of course, even                                    //
                    // with this check, another tab could insert a new valid token                                   //
                    // immediately before we clear localStorage here, which would                                    //
                    // lead to both tabs being logged out, but by checking the token                                 //
                    // in storage right now we hope to make that unlikely to happen.                                 //
                    //                                                                                               //
                    // If there is no token in storage right now, we don't have to                                   //
                    // do anything; whatever code removed the token from storage was                                 //
                    // responsible for calling `makeClientLoggedOut()`, or the                                       //
                    // periodic localStorage poll will call `makeClientLoggedOut`                                    //
                    // eventually if another tab wiped the token from storage.                                       //
                    if (storedTokenNow && storedTokenNow === result.token) {                                         // 264
                      self.makeClientLoggedOut();                                                                    // 265
                    }                                                                                                //
                  }                                                                                                  //
                  // Possibly a weird callback to call, but better than nothing if                                   //
                  // there is a reconnect between "login result received" and "data                                  //
                  // ready".                                                                                         //
                  loginCallbacks(error);                                                                             // 244
                }                                                                                                    //
                                                                                                                     //
                return userCallback;                                                                                 //
              }() });                                                                                                //
          }                                                                                                          //
        };                                                                                                           //
      }                                                                                                              //
  };                                                                                                                 //
                                                                                                                     //
  // This callback is called once the local cache of the current-user                                                //
  // subscription (and all subscriptions, in fact) are guaranteed to be up to                                        //
  // date.                                                                                                           //
  var loggedInAndDataReadyCallback = function loggedInAndDataReadyCallback(error, result) {                          // 169
    // If the login method returns its result but the connection is lost                                             //
    // before the data is in the local cache, it'll set an onReconnect (see                                          //
    // above). The onReconnect will try to log in using the token, and *it*                                          //
    // will call userCallback via its own version of this                                                            //
    // loggedInAndDataReadyCallback. So we don't have to do anything here.                                           //
    if (reconnected) return;                                                                                         // 287
                                                                                                                     //
    // Note that we need to call this even if _suppressLoggingIn is true,                                            //
    // because it could be matching a _setLoggingIn(true) from a                                                     //
    // half-completed pre-reconnect login method.                                                                    //
    self._setLoggingIn(false);                                                                                       // 281
    if (error || !result) {                                                                                          // 294
      error = error || new Error("No result from call to " + options.methodName);                                    // 295
      loginCallbacks(error);                                                                                         // 297
      return;                                                                                                        // 298
    }                                                                                                                //
    try {                                                                                                            // 300
      options.validateResult(result);                                                                                // 301
    } catch (e) {                                                                                                    //
      loginCallbacks(e);                                                                                             // 303
      return;                                                                                                        // 304
    }                                                                                                                //
                                                                                                                     //
    // Make the client logged in. (The user data should already be loaded!)                                          //
    self.makeClientLoggedIn(result.id, result.token, result.tokenExpires);                                           // 281
    loginCallbacks();                                                                                                // 309
  };                                                                                                                 //
                                                                                                                     //
  if (!options._suppressLoggingIn) self._setLoggingIn(true);                                                         // 312
  self.connection.apply(options.methodName, options.methodArguments, { wait: true, onResultReceived: onResultReceived }, loggedInAndDataReadyCallback);
};                                                                                                                   //
                                                                                                                     //
Ap.makeClientLoggedOut = function () {                                                                               // 321
  this._unstoreLoginToken();                                                                                         // 322
  this.connection.setUserId(null);                                                                                   // 323
  this.connection.onReconnect = null;                                                                                // 324
};                                                                                                                   //
                                                                                                                     //
Ap.makeClientLoggedIn = function (userId, token, tokenExpires) {                                                     // 327
  this._storeLoginToken(userId, token, tokenExpires);                                                                // 328
  this.connection.setUserId(userId);                                                                                 // 329
};                                                                                                                   //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Log the user out.                                                                                        //
 * @locus Client                                                                                                     //
 * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage meteor                                                                                         //
 */                                                                                                                  //
Meteor.logout = function (callback) {                                                                                // 338
  return Accounts.logout(callback);                                                                                  // 339
};                                                                                                                   //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Log out other clients logged in as the current user, but does not log out the client that calls this function.
 * @locus Client                                                                                                     //
 * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage meteor                                                                                         //
 */                                                                                                                  //
Meteor.logoutOtherClients = function (callback) {                                                                    // 348
  return Accounts.logoutOtherClients(callback);                                                                      // 349
};                                                                                                                   //
                                                                                                                     //
///                                                                                                                  //
/// LOGIN SERVICES                                                                                                   //
///                                                                                                                  //
                                                                                                                     //
// A reactive function returning whether the loginServiceConfiguration                                               //
// subscription is ready. Used by accounts-ui to hide the login button                                               //
// until we have all the configuration loaded                                                                        //
//                                                                                                                   //
Ap.loginServicesConfigured = function () {                                                                           // 361
  return this._loginServicesHandle.ready();                                                                          // 362
};                                                                                                                   //
                                                                                                                     //
// Some login services such as the redirect login flow or the resume                                                 //
// login handler can log the user in at page load time.  The                                                         //
// Meteor.loginWithX functions have a callback argument, but the                                                     //
// callback function instance won't be in memory any longer if the                                                   //
// page was reloaded.  The `onPageLoadLogin` function allows a                                                       //
// callback to be registered for the case where the login was                                                        //
// initiated in a previous VM, and we now have the result of the login                                               //
// attempt in a new VM.                                                                                              //
                                                                                                                     //
// Register a callback to be called if we have information about a                                                   //
// login attempt at page load time.  Call the callback immediately if                                                //
// we already have the page load login attempt info, otherwise stash                                                 //
// the callback to be called if and when we do get the attempt info.                                                 //
//                                                                                                                   //
Ap.onPageLoadLogin = function (f) {                                                                                  // 380
  if (this._pageLoadLoginAttemptInfo) {                                                                              // 381
    f(this._pageLoadLoginAttemptInfo);                                                                               // 382
  } else {                                                                                                           //
    this._pageLoadLoginCallbacks.push(f);                                                                            // 384
  }                                                                                                                  //
};                                                                                                                   //
                                                                                                                     //
// Receive the information about the login attempt at page load time.                                                //
// Call registered callbacks, and also record the info in case                                                       //
// someone's callback hasn't been registered yet.                                                                    //
//                                                                                                                   //
Ap._pageLoadLogin = function (attemptInfo) {                                                                         // 393
  if (this._pageLoadLoginAttemptInfo) {                                                                              // 394
    Meteor._debug("Ignoring unexpected duplicate page load login attempt info");                                     // 395
    return;                                                                                                          // 396
  }                                                                                                                  //
                                                                                                                     //
  _.each(this._pageLoadLoginCallbacks, function (callback) {                                                         // 399
    callback(attemptInfo);                                                                                           // 400
  });                                                                                                                //
                                                                                                                     //
  this._pageLoadLoginCallbacks = [];                                                                                 // 403
  this._pageLoadLoginAttemptInfo = attemptInfo;                                                                      // 404
};                                                                                                                   //
                                                                                                                     //
///                                                                                                                  //
/// HANDLEBARS HELPERS                                                                                               //
///                                                                                                                  //
                                                                                                                     //
// If our app has a Blaze, register the {{currentUser}} and {{loggingIn}}                                            //
// global helpers.                                                                                                   //
if (Package.blaze) {                                                                                                 // 414
  /**                                                                                                                //
   * @global                                                                                                         //
   * @name  currentUser                                                                                              //
   * @isHelper true                                                                                                  //
   * @summary Calls [Meteor.user()](#meteor_user). Use `{{#if currentUser}}` to check whether the user is logged in.
   */                                                                                                                //
  Package.blaze.Blaze.Template.registerHelper('currentUser', function () {                                           // 421
    return Meteor.user();                                                                                            // 422
  });                                                                                                                //
                                                                                                                     //
  /**                                                                                                                //
   * @global                                                                                                         //
   * @name  loggingIn                                                                                                //
   * @isHelper true                                                                                                  //
   * @summary Calls [Meteor.loggingIn()](#meteor_loggingin).                                                         //
   */                                                                                                                //
  Package.blaze.Blaze.Template.registerHelper('loggingIn', function () {                                             // 414
    return Meteor.loggingIn();                                                                                       // 432
  });                                                                                                                //
}                                                                                                                    //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"accounts_common.js":["babel-runtime/helpers/classCallCheck",function(require,exports){

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                   //
// packages/accounts-base/accounts_common.js                                                                         //
//                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                     //
exports.__esModule = true;                                                                                           //
exports.AccountsCommon = undefined;                                                                                  //
                                                                                                                     //
var _classCallCheck2 = require("babel-runtime/helpers/classCallCheck");                                              //
                                                                                                                     //
var _classCallCheck3 = _interopRequireDefault(_classCallCheck2);                                                     //
                                                                                                                     //
function _interopRequireDefault(obj) { return obj && obj.__esModule ? obj : { "default": obj }; }                    //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Super-constructor for AccountsClient and AccountsServer.                                                 //
 * @locus Anywhere                                                                                                   //
 * @class AccountsCommon                                                                                             //
 * @instancename accountsClientOrServer                                                                              //
 * @param options {Object} an object with fields:                                                                    //
 * - connection {Object} Optional DDP connection to reuse.                                                           //
 * - ddpUrl {String} Optional URL for creating a new DDP connection.                                                 //
 */                                                                                                                  //
                                                                                                                     //
var AccountsCommon = exports.AccountsCommon = function () {                                                          //
  function AccountsCommon(options) {                                                                                 // 11
    (0, _classCallCheck3["default"])(this, AccountsCommon);                                                          //
                                                                                                                     //
    // Currently this is read directly by packages like accounts-password                                            //
    // and accounts-ui-unstyled.                                                                                     //
    this._options = {};                                                                                              // 14
                                                                                                                     //
    // Note that setting this.connection = null causes this.users to be a                                            //
    // LocalCollection, which is not what we want.                                                                   //
    this.connection = undefined;                                                                                     // 11
    this._initConnection(options || {});                                                                             // 19
                                                                                                                     //
    // There is an allow call in accounts_server.js that restricts writes to                                         //
    // this collection.                                                                                              //
    this.users = new Mongo.Collection("users", {                                                                     // 11
      _preventAutopublish: true,                                                                                     // 24
      connection: this.connection                                                                                    // 25
    });                                                                                                              //
                                                                                                                     //
    // Callback exceptions are printed with Meteor._debug and ignored.                                               //
    this._onLoginHook = new Hook({                                                                                   // 11
      bindEnvironment: false,                                                                                        // 30
      debugPrintExceptions: "onLogin callback"                                                                       // 31
    });                                                                                                              //
                                                                                                                     //
    this._onLoginFailureHook = new Hook({                                                                            // 34
      bindEnvironment: false,                                                                                        // 35
      debugPrintExceptions: "onLoginFailure callback"                                                                // 36
    });                                                                                                              //
  }                                                                                                                  //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Get the current user id, or `null` if no user is logged in. A reactive data source.                    //
   * @locus Anywhere but publish functions                                                                           //
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsCommon.prototype.userId = function () {                                                                    // 10
    function userId() {                                                                                              //
      throw new Error("userId method not implemented");                                                              // 45
    }                                                                                                                //
                                                                                                                     //
    return userId;                                                                                                   //
  }();                                                                                                               //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Get the current user record, or `null` if no user is logged in. A reactive data source.                //
   * @locus Anywhere but publish functions                                                                           //
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsCommon.prototype.user = function () {                                                                      // 10
    function user() {                                                                                                //
      var userId = this.userId();                                                                                    // 53
      return userId ? this.users.findOne(userId) : null;                                                             // 54
    }                                                                                                                //
                                                                                                                     //
    return user;                                                                                                     //
  }();                                                                                                               //
                                                                                                                     //
  // Set up config for the accounts system. Call this on both the client                                             //
  // and the server.                                                                                                 //
  //                                                                                                                 //
  // Note that this method gets overridden on AccountsServer.prototype, but                                          //
  // the overriding method calls the overridden method.                                                              //
  //                                                                                                                 //
  // XXX we should add some enforcement that this is called on both the                                              //
  // client and the server. Otherwise, a user can                                                                    //
  // 'forbidClientAccountCreation' only on the client and while it looks                                             //
  // like their app is secure, the server will still accept createUser                                               //
  // calls. https://github.com/meteor/meteor/issues/828                                                              //
  //                                                                                                                 //
  // @param options {Object} an object with fields:                                                                  //
  // - sendVerificationEmail {Boolean}                                                                               //
  //     Send email address verification emails to new users created from                                            //
  //     client signups.                                                                                             //
  // - forbidClientAccountCreation {Boolean}                                                                         //
  //     Do not allow clients to create accounts directly.                                                           //
  // - restrictCreationByEmailDomain {Function or String}                                                            //
  //     Require created users to have an email matching the function or                                             //
  //     having the string as domain.                                                                                //
  // - loginExpirationInDays {Number}                                                                                //
  //     Number of days since login until a user is logged out (login token                                          //
  //     expires).                                                                                                   //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Set global accounts options.                                                                           //
   * @locus Anywhere                                                                                                 //
   * @param {Object} options                                                                                         //
   * @param {Boolean} options.sendVerificationEmail New users with an email address will receive an address verification email.
   * @param {Boolean} options.forbidClientAccountCreation Calls to [`createUser`](#accounts_createuser) from the client will be rejected. In addition, if you are using [accounts-ui](#accountsui), the "Create account" link will not be available.
   * @param {String | Function} options.restrictCreationByEmailDomain If set to a string, only allows new users if the domain part of their email address matches the string. If set to a function, only allows new users if the function returns true.  The function is passed the full email address of the proposed new user.  Works with password-based sign-in and external services that expose email addresses (Google, Facebook, GitHub). All existing users still can log in after enabling this option. Example: `Accounts.config({ restrictCreationByEmailDomain: 'school.edu' })`.
   * @param {Number} options.loginExpirationInDays The number of days from when a user logs in until their token expires and they are logged out. Defaults to 90. Set to `null` to disable login expiration.
   * @param {String} options.oauthSecretKey When using the `oauth-encryption` package, the 16 byte key using to encrypt sensitive account credentials in the database, encoded in base64.  This option may only be specifed on the server.  See packages/oauth-encryption/README.md for details.
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsCommon.prototype.config = function () {                                                                    // 10
    function config(options) {                                                                                       //
      var self = this;                                                                                               // 93
                                                                                                                     //
      // We don't want users to accidentally only call Accounts.config on the                                        //
      // client, where some of the options will have partial effects (eg removing                                    //
      // the "create account" button from accounts-ui if forbidClientAccountCreation                                 //
      // is set, or redirecting Google login to a specific-domain page) without                                      //
      // having their full effects.                                                                                  //
      if (Meteor.isServer) {                                                                                         // 92
        __meteor_runtime_config__.accountsConfigCalled = true;                                                       // 101
      } else if (!__meteor_runtime_config__.accountsConfigCalled) {                                                  //
        // XXX would be nice to "crash" the client and replace the UI with an error                                  //
        // message, but there's no trivial way to do this.                                                           //
        Meteor._debug("Accounts.config was called on the client but not on the " + "server; some configuration options may not take effect.");
      }                                                                                                              //
                                                                                                                     //
      // We need to validate the oauthSecretKey option at the time                                                   //
      // Accounts.config is called. We also deliberately don't store the                                             //
      // oauthSecretKey in Accounts._options.                                                                        //
      if (_.has(options, "oauthSecretKey")) {                                                                        // 92
        if (Meteor.isClient) throw new Error("The oauthSecretKey option may only be specified on the server");       // 113
        if (!Package["oauth-encryption"]) throw new Error("The oauth-encryption package must be loaded to set oauthSecretKey");
        Package["oauth-encryption"].OAuthEncryption.loadKey(options.oauthSecretKey);                                 // 117
        options = _.omit(options, "oauthSecretKey");                                                                 // 118
      }                                                                                                              //
                                                                                                                     //
      // validate option keys                                                                                        //
      var VALID_KEYS = ["sendVerificationEmail", "forbidClientAccountCreation", "restrictCreationByEmailDomain", "loginExpirationInDays"];
      _.each(_.keys(options), function (key) {                                                                       // 124
        if (!_.contains(VALID_KEYS, key)) {                                                                          // 125
          throw new Error("Accounts.config: Invalid key: " + key);                                                   // 126
        }                                                                                                            //
      });                                                                                                            //
                                                                                                                     //
      // set values in Accounts._options                                                                             //
      _.each(VALID_KEYS, function (key) {                                                                            // 92
        if (key in options) {                                                                                        // 132
          if (key in self._options) {                                                                                // 133
            throw new Error("Can't set `" + key + "` more than once");                                               // 134
          }                                                                                                          //
          self._options[key] = options[key];                                                                         // 136
        }                                                                                                            //
      });                                                                                                            //
    }                                                                                                                //
                                                                                                                     //
    return config;                                                                                                   //
  }();                                                                                                               //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Register a callback to be called after a login attempt succeeds.                                       //
   * @locus Anywhere                                                                                                 //
   * @param {Function} func The callback to be called when login is successful.                                      //
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsCommon.prototype.onLogin = function () {                                                                   // 10
    function onLogin(func) {                                                                                         //
      return this._onLoginHook.register(func);                                                                       // 147
    }                                                                                                                //
                                                                                                                     //
    return onLogin;                                                                                                  //
  }();                                                                                                               //
                                                                                                                     //
  /**                                                                                                                //
   * @summary Register a callback to be called after a login attempt fails.                                          //
   * @locus Anywhere                                                                                                 //
   * @param {Function} func The callback to be called after the login has failed.                                    //
   */                                                                                                                //
                                                                                                                     //
                                                                                                                     //
  AccountsCommon.prototype.onLoginFailure = function () {                                                            // 10
    function onLoginFailure(func) {                                                                                  //
      return this._onLoginFailureHook.register(func);                                                                // 156
    }                                                                                                                //
                                                                                                                     //
    return onLoginFailure;                                                                                           //
  }();                                                                                                               //
                                                                                                                     //
  AccountsCommon.prototype._initConnection = function () {                                                           // 10
    function _initConnection(options) {                                                                              //
      if (!Meteor.isClient) {                                                                                        // 160
        return;                                                                                                      // 161
      }                                                                                                              //
                                                                                                                     //
      // The connection used by the Accounts system. This is the connection                                          //
      // that will get logged in by Meteor.login(), and this is the                                                  //
      // connection whose login state will be reflected by Meteor.userId().                                          //
      //                                                                                                             //
      // It would be much preferable for this to be in accounts_client.js,                                           //
      // but it has to be here because it's needed to create the                                                     //
      // Meteor.users collection.                                                                                    //
                                                                                                                     //
      if (options.connection) {                                                                                      // 159
        this.connection = options.connection;                                                                        // 173
      } else if (options.ddpUrl) {                                                                                   //
        this.connection = DDP.connect(options.ddpUrl);                                                               // 175
      } else if (typeof __meteor_runtime_config__ !== "undefined" && __meteor_runtime_config__.ACCOUNTS_CONNECTION_URL) {
        // Temporary, internal hook to allow the server to point the client                                          //
        // to a different authentication server. This is for a very                                                  //
        // particular use case that comes up when implementing a oauth                                               //
        // server. Unsupported and may go away at any point in time.                                                 //
        //                                                                                                           //
        // We will eventually provide a general way to use account-base                                              //
        // against any DDP connection, not just one special one.                                                     //
        this.connection = DDP.connect(__meteor_runtime_config__.ACCOUNTS_CONNECTION_URL);                            // 185
      } else {                                                                                                       //
        this.connection = Meteor.connection;                                                                         // 188
      }                                                                                                              //
    }                                                                                                                //
                                                                                                                     //
    return _initConnection;                                                                                          //
  }();                                                                                                               //
                                                                                                                     //
  AccountsCommon.prototype._getTokenLifetimeMs = function () {                                                       // 10
    function _getTokenLifetimeMs() {                                                                                 //
      return (this._options.loginExpirationInDays || DEFAULT_LOGIN_EXPIRATION_DAYS) * 24 * 60 * 60 * 1000;           // 193
    }                                                                                                                //
                                                                                                                     //
    return _getTokenLifetimeMs;                                                                                      //
  }();                                                                                                               //
                                                                                                                     //
  AccountsCommon.prototype._tokenExpiration = function () {                                                          // 10
    function _tokenExpiration(when) {                                                                                //
      // We pass when through the Date constructor for backwards compatibility;                                      //
      // `when` used to be a number.                                                                                 //
      return new Date(new Date(when).getTime() + this._getTokenLifetimeMs());                                        // 200
    }                                                                                                                //
                                                                                                                     //
    return _tokenExpiration;                                                                                         //
  }();                                                                                                               //
                                                                                                                     //
  AccountsCommon.prototype._tokenExpiresSoon = function () {                                                         // 10
    function _tokenExpiresSoon(when) {                                                                               //
      var minLifetimeMs = .1 * this._getTokenLifetimeMs();                                                           // 204
      var minLifetimeCapMs = MIN_TOKEN_LIFETIME_CAP_SECS * 1000;                                                     // 205
      if (minLifetimeMs > minLifetimeCapMs) minLifetimeMs = minLifetimeCapMs;                                        // 206
      return new Date() > new Date(when) - minLifetimeMs;                                                            // 208
    }                                                                                                                //
                                                                                                                     //
    return _tokenExpiresSoon;                                                                                        //
  }();                                                                                                               //
                                                                                                                     //
  return AccountsCommon;                                                                                             //
}();                                                                                                                 //
                                                                                                                     //
var Ap = AccountsCommon.prototype;                                                                                   // 212
                                                                                                                     //
// Note that Accounts is defined separately in accounts_client.js and                                                //
// accounts_server.js.                                                                                               //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Get the current user id, or `null` if no user is logged in. A reactive data source.                      //
 * @locus Anywhere but publish functions                                                                             //
 * @importFromPackage meteor                                                                                         //
 */                                                                                                                  //
Meteor.userId = function () {                                                                                        // 222
  return Accounts.userId();                                                                                          // 223
};                                                                                                                   //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Get the current user record, or `null` if no user is logged in. A reactive data source.                  //
 * @locus Anywhere but publish functions                                                                             //
 * @importFromPackage meteor                                                                                         //
 */                                                                                                                  //
Meteor.user = function () {                                                                                          // 231
  return Accounts.user();                                                                                            // 232
};                                                                                                                   //
                                                                                                                     //
// how long (in days) until a login token expires                                                                    //
var DEFAULT_LOGIN_EXPIRATION_DAYS = 90;                                                                              // 236
// Clients don't try to auto-login with a token that is going to expire within                                       //
// .1 * DEFAULT_LOGIN_EXPIRATION_DAYS, capped at MIN_TOKEN_LIFETIME_CAP_SECS.                                        //
// Tries to avoid abrupt disconnects from expiring tokens.                                                           //
var MIN_TOKEN_LIFETIME_CAP_SECS = 3600; // one hour                                                                  // 240
// how often (in milliseconds) we check for expired tokens                                                           //
EXPIRE_TOKENS_INTERVAL_MS = 600 * 1000; // 10 minutes                                                                // 242
// how long we wait before logging out clients when Meteor.logoutOtherClients is                                     //
// called                                                                                                            //
CONNECTION_CLOSE_DELAY_MS = 10 * 1000;                                                                               // 245
                                                                                                                     //
// loginServiceConfiguration and ConfigError are maintained for backwards compatibility                              //
Meteor.startup(function () {                                                                                         // 248
  var ServiceConfiguration = Package['service-configuration'].ServiceConfiguration;                                  // 249
  Ap.loginServiceConfiguration = ServiceConfiguration.configurations;                                                // 251
  Ap.ConfigError = ServiceConfiguration.ConfigError;                                                                 // 252
});                                                                                                                  //
                                                                                                                     //
// Thrown when the user cancels the login process (eg, closes an oauth                                               //
// popup, declines retina scan, etc)                                                                                 //
var lceName = 'Accounts.LoginCancelledError';                                                                        // 257
Ap.LoginCancelledError = Meteor.makeErrorType(lceName, function (description) {                                      // 258
  this.message = description;                                                                                        // 261
});                                                                                                                  //
Ap.LoginCancelledError.prototype.name = lceName;                                                                     // 264
                                                                                                                     //
// This is used to transmit specific subclass errors over the wire. We should                                        //
// come up with a more generic way to do this (eg, with some sort of symbolic                                        //
// error code rather than a number).                                                                                 //
Ap.LoginCancelledError.numericError = 0x8acdc2f;                                                                     // 269
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"localstorage_token.js":["./accounts_client.js",function(require){

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                   //
// packages/accounts-base/localstorage_token.js                                                                      //
//                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                     //
var _accounts_client = require("./accounts_client.js");                                                              // 1
                                                                                                                     //
var Ap = _accounts_client.AccountsClient.prototype;                                                                  // 2
                                                                                                                     //
// This file deals with storing a login token and user id in the                                                     //
// browser's localStorage facility. It polls local storage every few                                                 //
// seconds to synchronize login state between multiple tabs in the same                                              //
// browser.                                                                                                          //
                                                                                                                     //
// Login with a Meteor access token. This is the only public function                                                //
// here.                                                                                                             //
Meteor.loginWithToken = function (token, callback) {                                                                 // 11
  return Accounts.loginWithToken(token, callback);                                                                   // 12
};                                                                                                                   //
                                                                                                                     //
Ap.loginWithToken = function (token, callback) {                                                                     // 15
  this.callLoginMethod({                                                                                             // 16
    methodArguments: [{                                                                                              // 17
      resume: token                                                                                                  // 18
    }],                                                                                                              //
    userCallback: callback                                                                                           // 20
  });                                                                                                                //
};                                                                                                                   //
                                                                                                                     //
// Semi-internal API. Call this function to re-enable auto login after                                               //
// if it was disabled at startup.                                                                                    //
Ap._enableAutoLogin = function () {                                                                                  // 26
  this._autoLoginEnabled = true;                                                                                     // 27
  this._pollStoredLoginToken();                                                                                      // 28
};                                                                                                                   //
                                                                                                                     //
///                                                                                                                  //
/// STORING                                                                                                          //
///                                                                                                                  //
                                                                                                                     //
// Call this from the top level of the test file for any test that does                                              //
// logging in and out, to protect multiple tabs running the same tests                                               //
// simultaneously from interfering with each others' localStorage.                                                   //
Ap._isolateLoginTokenForTest = function () {                                                                         // 39
  this.LOGIN_TOKEN_KEY = this.LOGIN_TOKEN_KEY + Random.id();                                                         // 40
  this.USER_ID_KEY = this.USER_ID_KEY + Random.id();                                                                 // 41
};                                                                                                                   //
                                                                                                                     //
Ap._storeLoginToken = function (userId, token, tokenExpires) {                                                       // 44
  Meteor._localStorage.setItem(this.USER_ID_KEY, userId);                                                            // 45
  Meteor._localStorage.setItem(this.LOGIN_TOKEN_KEY, token);                                                         // 46
  if (!tokenExpires) tokenExpires = this._tokenExpiration(new Date());                                               // 47
  Meteor._localStorage.setItem(this.LOGIN_TOKEN_EXPIRES_KEY, tokenExpires);                                          // 49
                                                                                                                     //
  // to ensure that the localstorage poller doesn't end up trying to                                                 //
  // connect a second time                                                                                           //
  this._lastLoginTokenWhenPolled = token;                                                                            // 44
};                                                                                                                   //
                                                                                                                     //
Ap._unstoreLoginToken = function () {                                                                                // 56
  Meteor._localStorage.removeItem(this.USER_ID_KEY);                                                                 // 57
  Meteor._localStorage.removeItem(this.LOGIN_TOKEN_KEY);                                                             // 58
  Meteor._localStorage.removeItem(this.LOGIN_TOKEN_EXPIRES_KEY);                                                     // 59
                                                                                                                     //
  // to ensure that the localstorage poller doesn't end up trying to                                                 //
  // connect a second time                                                                                           //
  this._lastLoginTokenWhenPolled = null;                                                                             // 56
};                                                                                                                   //
                                                                                                                     //
// This is private, but it is exported for now because it is used by a                                               //
// test in accounts-password.                                                                                        //
//                                                                                                                   //
Ap._storedLoginToken = function () {                                                                                 // 69
  return Meteor._localStorage.getItem(this.LOGIN_TOKEN_KEY);                                                         // 70
};                                                                                                                   //
                                                                                                                     //
Ap._storedLoginTokenExpires = function () {                                                                          // 73
  return Meteor._localStorage.getItem(this.LOGIN_TOKEN_EXPIRES_KEY);                                                 // 74
};                                                                                                                   //
                                                                                                                     //
Ap._storedUserId = function () {                                                                                     // 77
  return Meteor._localStorage.getItem(this.USER_ID_KEY);                                                             // 78
};                                                                                                                   //
                                                                                                                     //
Ap._unstoreLoginTokenIfExpiresSoon = function () {                                                                   // 81
  var tokenExpires = this._storedLoginTokenExpires();                                                                // 82
  if (tokenExpires && this._tokenExpiresSoon(new Date(tokenExpires))) {                                              // 83
    this._unstoreLoginToken();                                                                                       // 84
  }                                                                                                                  //
};                                                                                                                   //
                                                                                                                     //
///                                                                                                                  //
/// AUTO-LOGIN                                                                                                       //
///                                                                                                                  //
                                                                                                                     //
Ap._initLocalStorage = function () {                                                                                 // 92
  var self = this;                                                                                                   // 93
                                                                                                                     //
  // Key names to use in localStorage                                                                                //
  self.LOGIN_TOKEN_KEY = "Meteor.loginToken";                                                                        // 92
  self.LOGIN_TOKEN_EXPIRES_KEY = "Meteor.loginTokenExpires";                                                         // 97
  self.USER_ID_KEY = "Meteor.userId";                                                                                // 98
                                                                                                                     //
  var rootUrlPathPrefix = __meteor_runtime_config__.ROOT_URL_PATH_PREFIX;                                            // 100
  if (rootUrlPathPrefix || this.connection !== Meteor.connection) {                                                  // 101
    // We want to keep using the same keys for existing apps that do not                                             //
    // set a custom ROOT_URL_PATH_PREFIX, so that most users will not have                                           //
    // to log in again after an app updates to a version of Meteor that                                              //
    // contains this code, but it's generally preferable to namespace the                                            //
    // keys so that connections from distinct apps to distinct DDP URLs                                              //
    // will be distinct in Meteor._localStorage.                                                                     //
    var namespace = ":" + this.connection._stream.rawUrl;                                                            // 108
    if (rootUrlPathPrefix) {                                                                                         // 109
      namespace += ":" + rootUrlPathPrefix;                                                                          // 110
    }                                                                                                                //
    self.LOGIN_TOKEN_KEY += namespace;                                                                               // 112
    self.LOGIN_TOKEN_EXPIRES_KEY += namespace;                                                                       // 113
    self.USER_ID_KEY += namespace;                                                                                   // 114
  }                                                                                                                  //
                                                                                                                     //
  if (self._autoLoginEnabled) {                                                                                      // 117
    // Immediately try to log in via local storage, so that any DDP                                                  //
    // messages are sent after we have established our user account                                                  //
    self._unstoreLoginTokenIfExpiresSoon();                                                                          // 120
    var token = self._storedLoginToken();                                                                            // 121
    if (token) {                                                                                                     // 122
      // On startup, optimistically present us as logged in while the                                                //
      // request is in flight. This reduces page flicker on startup.                                                 //
      var userId = self._storedUserId();                                                                             // 125
      userId && self.connection.setUserId(userId);                                                                   // 126
      self.loginWithToken(token, function (err) {                                                                    // 127
        if (err) {                                                                                                   // 128
          Meteor._debug("Error logging in with token: " + err);                                                      // 129
          self.makeClientLoggedOut();                                                                                // 130
        }                                                                                                            //
                                                                                                                     //
        self._pageLoadLogin({                                                                                        // 133
          type: "resume",                                                                                            // 134
          allowed: !err,                                                                                             // 135
          error: err,                                                                                                // 136
          methodName: "login",                                                                                       // 137
          // XXX This is duplicate code with loginWithToken, but                                                     //
          // loginWithToken can also be called at other times besides                                                //
          // page load.                                                                                              //
          methodArguments: [{ resume: token }]                                                                       // 141
        });                                                                                                          //
      });                                                                                                            //
    }                                                                                                                //
  }                                                                                                                  //
                                                                                                                     //
  // Poll local storage every 3 seconds to login if someone logged in in                                             //
  // another tab                                                                                                     //
  self._lastLoginTokenWhenPolled = token;                                                                            // 92
                                                                                                                     //
  if (self._pollIntervalTimer) {                                                                                     // 151
    // Unlikely that _initLocalStorage will be called more than once for                                             //
    // the same AccountsClient instance, but just in case...                                                         //
    clearInterval(self._pollIntervalTimer);                                                                          // 154
  }                                                                                                                  //
                                                                                                                     //
  self._pollIntervalTimer = setInterval(function () {                                                                // 157
    self._pollStoredLoginToken();                                                                                    // 158
  }, 3000);                                                                                                          //
};                                                                                                                   //
                                                                                                                     //
Ap._pollStoredLoginToken = function () {                                                                             // 162
  var self = this;                                                                                                   // 163
                                                                                                                     //
  if (!self._autoLoginEnabled) {                                                                                     // 165
    return;                                                                                                          // 166
  }                                                                                                                  //
                                                                                                                     //
  var currentLoginToken = self._storedLoginToken();                                                                  // 169
                                                                                                                     //
  // != instead of !== just to make sure undefined and null are treated the same                                     //
  if (self._lastLoginTokenWhenPolled != currentLoginToken) {                                                         // 162
    if (currentLoginToken) {                                                                                         // 173
      self.loginWithToken(currentLoginToken, function (err) {                                                        // 174
        if (err) {                                                                                                   // 175
          self.makeClientLoggedOut();                                                                                // 176
        }                                                                                                            //
      });                                                                                                            //
    } else {                                                                                                         //
      self.logout();                                                                                                 // 180
    }                                                                                                                //
  }                                                                                                                  //
                                                                                                                     //
  self._lastLoginTokenWhenPolled = currentLoginToken;                                                                // 184
};                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}],"url_client.js":["./accounts_client.js",function(require,exports){

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                   //
// packages/accounts-base/url_client.js                                                                              //
//                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                     //
exports.__esModule = true;                                                                                           //
exports.AccountsTest = undefined;                                                                                    //
                                                                                                                     //
var _accounts_client = require("./accounts_client.js");                                                              // 1
                                                                                                                     //
var Ap = _accounts_client.AccountsClient.prototype;                                                                  // 3
                                                                                                                     //
// All of the special hash URLs we support for accounts interactions                                                 //
var accountsPaths = ["reset-password", "verify-email", "enroll-account"];                                            // 6
                                                                                                                     //
var savedHash = window.location.hash;                                                                                // 8
                                                                                                                     //
Ap._initUrlMatching = function () {                                                                                  // 10
  // By default, allow the autologin process to happen.                                                              //
  this._autoLoginEnabled = true;                                                                                     // 12
                                                                                                                     //
  // We only support one callback per URL.                                                                           //
  this._accountsCallbacks = {};                                                                                      // 10
                                                                                                                     //
  // Try to match the saved value of window.location.hash.                                                           //
  this._attemptToMatchHash();                                                                                        // 10
};                                                                                                                   //
                                                                                                                     //
// Separate out this functionality for testing                                                                       //
                                                                                                                     //
Ap._attemptToMatchHash = function () {                                                                               // 23
  _attemptToMatchHash(this, savedHash, defaultSuccessHandler);                                                       // 24
};                                                                                                                   //
                                                                                                                     //
// Note that both arguments are optional and are currently only passed by                                            //
// accounts_url_tests.js.                                                                                            //
function _attemptToMatchHash(accounts, hash, success) {                                                              // 29
  _.each(accountsPaths, function (urlPart) {                                                                         // 30
    var token;                                                                                                       // 31
                                                                                                                     //
    var tokenRegex = new RegExp("^\\#\\/" + urlPart + "\\/(.*)$");                                                   // 33
    var match = hash.match(tokenRegex);                                                                              // 34
                                                                                                                     //
    if (match) {                                                                                                     // 36
      token = match[1];                                                                                              // 37
                                                                                                                     //
      // XXX COMPAT WITH 0.9.3                                                                                       //
      if (urlPart === "reset-password") {                                                                            // 36
        accounts._resetPasswordToken = token;                                                                        // 41
      } else if (urlPart === "verify-email") {                                                                       //
        accounts._verifyEmailToken = token;                                                                          // 43
      } else if (urlPart === "enroll-account") {                                                                     //
        accounts._enrollAccountToken = token;                                                                        // 45
      }                                                                                                              //
    } else {                                                                                                         //
      return;                                                                                                        // 48
    }                                                                                                                //
                                                                                                                     //
    // If no handlers match the hash, then maybe it's meant to be consumed                                           //
    // by some entirely different code, so we only clear it the first time                                           //
    // a handler successfully matches. Note that later handlers reuse the                                            //
    // savedHash, so clearing window.location.hash here will not interfere                                           //
    // with their needs.                                                                                             //
    window.location.hash = "";                                                                                       // 30
                                                                                                                     //
    // Do some stuff with the token we matched                                                                       //
    success.call(accounts, token, urlPart);                                                                          // 30
  });                                                                                                                //
}                                                                                                                    //
                                                                                                                     //
function defaultSuccessHandler(token, urlPart) {                                                                     // 63
  var self = this;                                                                                                   // 64
                                                                                                                     //
  // put login in a suspended state to wait for the interaction to finish                                            //
  self._autoLoginEnabled = false;                                                                                    // 63
                                                                                                                     //
  // wait for other packages to register callbacks                                                                   //
  Meteor.startup(function () {                                                                                       // 63
    // if a callback has been registered for this kind of token, call it                                             //
    if (self._accountsCallbacks[urlPart]) {                                                                          // 72
      self._accountsCallbacks[urlPart](token, function () {                                                          // 73
        self._enableAutoLogin();                                                                                     // 74
      });                                                                                                            //
    }                                                                                                                //
  });                                                                                                                //
}                                                                                                                    //
                                                                                                                     //
// Export for testing                                                                                                //
var AccountsTest = exports.AccountsTest = {                                                                          // 81
  attemptToMatchHash: function () {                                                                                  // 82
    function attemptToMatchHash(hash, success) {                                                                     // 82
      return _attemptToMatchHash(Accounts, hash, success);                                                           // 83
    }                                                                                                                //
                                                                                                                     //
    return attemptToMatchHash;                                                                                       //
  }()                                                                                                                //
};                                                                                                                   //
                                                                                                                     //
// XXX these should be moved to accounts-password eventually. Right now                                              //
// this is prevented by the need to set autoLoginEnabled=false, but in                                               //
// some bright future we won't need to do that anymore.                                                              //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Register a function to call when a reset password link is clicked                                        //
 * in an email sent by                                                                                               //
 * [`Accounts.sendResetPasswordEmail`](#accounts_sendresetpasswordemail).                                            //
 * This function should be called in top-level code, not inside                                                      //
 * `Meteor.startup()`.                                                                                               //
 * @memberof! Accounts                                                                                               //
 * @name onResetPasswordLink                                                                                         //
 * @param  {Function} callback The function to call. It is given two arguments:                                      //
 *                                                                                                                   //
 * 1. `token`: A password reset token that can be passed to                                                          //
 * [`Accounts.resetPassword`](#accounts_resetpassword).                                                              //
 * 2. `done`: A function to call when the password reset UI flow is complete. The normal                             //
 * login process is suspended until this function is called, so that the                                             //
 * password for user A can be reset even if user B was logged in.                                                    //
 * @locus Client                                                                                                     //
 */                                                                                                                  //
Ap.onResetPasswordLink = function (callback) {                                                                       // 108
  if (this._accountsCallbacks["reset-password"]) {                                                                   // 109
    Meteor._debug("Accounts.onResetPasswordLink was called more than once. " + "Only one callback added will be executed.");
  }                                                                                                                  //
                                                                                                                     //
  this._accountsCallbacks["reset-password"] = callback;                                                              // 114
};                                                                                                                   //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Register a function to call when an email verification link is                                           //
 * clicked in an email sent by                                                                                       //
 * [`Accounts.sendVerificationEmail`](#accounts_sendverificationemail).                                              //
 * This function should be called in top-level code, not inside                                                      //
 * `Meteor.startup()`.                                                                                               //
 * @memberof! Accounts                                                                                               //
 * @name onEmailVerificationLink                                                                                     //
 * @param  {Function} callback The function to call. It is given two arguments:                                      //
 *                                                                                                                   //
 * 1. `token`: An email verification token that can be passed to                                                     //
 * [`Accounts.verifyEmail`](#accounts_verifyemail).                                                                  //
 * 2. `done`: A function to call when the email verification UI flow is complete.                                    //
 * The normal login process is suspended until this function is called, so                                           //
 * that the user can be notified that they are verifying their email before                                          //
 * being logged in.                                                                                                  //
 * @locus Client                                                                                                     //
 */                                                                                                                  //
Ap.onEmailVerificationLink = function (callback) {                                                                   // 135
  if (this._accountsCallbacks["verify-email"]) {                                                                     // 136
    Meteor._debug("Accounts.onEmailVerificationLink was called more than once. " + "Only one callback added will be executed.");
  }                                                                                                                  //
                                                                                                                     //
  this._accountsCallbacks["verify-email"] = callback;                                                                // 141
};                                                                                                                   //
                                                                                                                     //
/**                                                                                                                  //
 * @summary Register a function to call when an account enrollment link is                                           //
 * clicked in an email sent by                                                                                       //
 * [`Accounts.sendEnrollmentEmail`](#accounts_sendenrollmentemail).                                                  //
 * This function should be called in top-level code, not inside                                                      //
 * `Meteor.startup()`.                                                                                               //
 * @memberof! Accounts                                                                                               //
 * @name onEnrollmentLink                                                                                            //
 * @param  {Function} callback The function to call. It is given two arguments:                                      //
 *                                                                                                                   //
 * 1. `token`: A password reset token that can be passed to                                                          //
 * [`Accounts.resetPassword`](#accounts_resetpassword) to give the newly                                             //
 * enrolled account a password.                                                                                      //
 * 2. `done`: A function to call when the enrollment UI flow is complete.                                            //
 * The normal login process is suspended until this function is called, so that                                      //
 * user A can be enrolled even if user B was logged in.                                                              //
 * @locus Client                                                                                                     //
 */                                                                                                                  //
Ap.onEnrollmentLink = function (callback) {                                                                          // 162
  if (this._accountsCallbacks["enroll-account"]) {                                                                   // 163
    Meteor._debug("Accounts.onEnrollmentLink was called more than once. " + "Only one callback added will be executed.");
  }                                                                                                                  //
                                                                                                                     //
  this._accountsCallbacks["enroll-account"] = callback;                                                              // 168
};                                                                                                                   //
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}]}}}},{"extensions":[".js",".json"]});
var exports = require("./node_modules/meteor/accounts-base/client_main.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['accounts-base'] = exports, {
  Accounts: Accounts
});

})();
