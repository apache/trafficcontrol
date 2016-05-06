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
var Accounts = Package['accounts-base'].Accounts;
var SRP = Package.srp.SRP;
var SHA256 = Package.sha.SHA256;
var EJSON = Package.ejson.EJSON;
var DDP = Package['ddp-client'].DDP;
var check = Package.check.check;
var Match = Package.check.Match;
var _ = Package.underscore._;
var meteorInstall = Package.modules.meteorInstall;
var Buffer = Package.modules.Buffer;
var process = Package.modules.process;
var Symbol = Package['ecmascript-runtime'].Symbol;
var Map = Package['ecmascript-runtime'].Map;
var Set = Package['ecmascript-runtime'].Set;
var meteorBabelHelpers = Package['babel-runtime'].meteorBabelHelpers;
var Promise = Package.promise.Promise;

var require = meteorInstall({"node_modules":{"meteor":{"accounts-password":{"password_client.js":function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                            //
// packages/accounts-password/password_client.js                                                              //
//                                                                                                            //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                              //
// Attempt to log in with a password.                                                                         //
//                                                                                                            //
// @param selector {String|Object} One of the following:                                                      //
//   - {username: (username)}                                                                                 //
//   - {email: (email)}                                                                                       //
//   - a string which may be a username or email, depending on whether                                        //
//     it contains "@".                                                                                       //
// @param password {String}                                                                                   //
// @param callback {Function(error|undefined)}                                                                //
                                                                                                              //
/**                                                                                                           //
 * @summary Log the user in with a password.                                                                  //
 * @locus Client                                                                                              //
 * @param {Object | String} user                                                                              //
 *   Either a string interpreted as a username or an email; or an object with a                               //
 *   single key: `email`, `username` or `id`. Username or email match in a case                               //
 *   insensitive manner.                                                                                      //
 * @param {String} password The user's password.                                                              //
 * @param {Function} [callback] Optional callback.                                                            //
 *   Called with no arguments on success, or with a single `Error` argument                                   //
 *   on failure.                                                                                              //
 * @importFromPackage meteor                                                                                  //
 */                                                                                                           //
Meteor.loginWithPassword = function (selector, password, callback) {                                          // 24
  if (typeof selector === 'string') if (selector.indexOf('@') === -1) selector = { username: selector };else selector = { email: selector };
                                                                                                              //
  Accounts.callLoginMethod({                                                                                  // 31
    methodArguments: [{                                                                                       // 32
      user: selector,                                                                                         // 33
      password: Accounts._hashPassword(password)                                                              // 34
    }],                                                                                                       //
    userCallback: function () {                                                                               // 36
      function userCallback(error, result) {                                                                  // 36
        if (error && error.error === 400 && error.reason === 'old password format') {                         // 37
          // The "reason" string should match the error thrown in the                                         //
          // password login handler in password_server.js.                                                    //
                                                                                                              //
          // XXX COMPAT WITH 0.8.1.3                                                                          //
          // If this user's last login was with a previous version of                                         //
          // Meteor that used SRP, then the server throws this error to                                       //
          // indicate that we should try again. The error includes the                                        //
          // user's SRP identity. We provide a value derived from the                                         //
          // identity and the password to prove to the server that we know                                    //
          // the password without requiring a full SRP flow, as well as                                       //
          // SHA256(password), which the server bcrypts and stores in                                         //
          // place of the old SRP information for this user.                                                  //
          srpUpgradePath({                                                                                    // 51
            upgradeError: error,                                                                              // 52
            userSelector: selector,                                                                           // 53
            plaintextPassword: password                                                                       // 54
          }, callback);                                                                                       //
        } else if (error) {                                                                                   //
          callback && callback(error);                                                                        // 58
        } else {                                                                                              //
          callback && callback();                                                                             // 60
        }                                                                                                     //
      }                                                                                                       //
                                                                                                              //
      return userCallback;                                                                                    //
    }()                                                                                                       //
  });                                                                                                         //
};                                                                                                            //
                                                                                                              //
Accounts._hashPassword = function (password) {                                                                // 66
  return {                                                                                                    // 67
    digest: SHA256(password),                                                                                 // 68
    algorithm: "sha-256"                                                                                      // 69
  };                                                                                                          //
};                                                                                                            //
                                                                                                              //
// XXX COMPAT WITH 0.8.1.3                                                                                    //
// The server requested an upgrade from the old SRP password format,                                          //
// so supply the needed SRP identity to login. Options:                                                       //
//   - upgradeError: the error object that the server returned to tell                                        //
//     us to upgrade from SRP to bcrypt.                                                                      //
//   - userSelector: selector to retrieve the user object                                                     //
//   - plaintextPassword: the password as a string                                                            //
var srpUpgradePath = function srpUpgradePath(options, callback) {                                             // 80
  var details;                                                                                                // 81
  try {                                                                                                       // 82
    details = EJSON.parse(options.upgradeError.details);                                                      // 83
  } catch (e) {}                                                                                              //
  if (!(details && details.format === 'srp')) {                                                               // 85
    callback && callback(new Meteor.Error(400, "Password is old. Please reset your " + "password."));         // 86
  } else {                                                                                                    //
    Accounts.callLoginMethod({                                                                                // 90
      methodArguments: [{                                                                                     // 91
        user: options.userSelector,                                                                           // 92
        srp: SHA256(details.identity + ":" + options.plaintextPassword),                                      // 93
        password: Accounts._hashPassword(options.plaintextPassword)                                           // 94
      }],                                                                                                     //
      userCallback: callback                                                                                  // 96
    });                                                                                                       //
  }                                                                                                           //
};                                                                                                            //
                                                                                                              //
// Attempt to log in as a new user.                                                                           //
                                                                                                              //
/**                                                                                                           //
 * @summary Create a new user.                                                                                //
 * @locus Anywhere                                                                                            //
 * @param {Object} options                                                                                    //
 * @param {String} options.username A unique name for this user.                                              //
 * @param {String} options.email The user's email address.                                                    //
 * @param {String} options.password The user's password. This is __not__ sent in plain text over the wire.    //
 * @param {Object} options.profile The user's profile, typically including the `name` field.                  //
 * @param {Function} [callback] Client only, optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage accounts-base                                                                           //
 */                                                                                                           //
Accounts.createUser = function (options, callback) {                                                          // 115
  options = _.clone(options); // we'll be modifying options                                                   // 116
                                                                                                              //
  if (typeof options.password !== 'string') throw new Error("options.password must be a string");             // 115
  if (!options.password) {                                                                                    // 120
    callback(new Meteor.Error(400, "Password may not be empty"));                                             // 121
    return;                                                                                                   // 122
  }                                                                                                           //
                                                                                                              //
  // Replace password with the hashed password.                                                               //
  options.password = Accounts._hashPassword(options.password);                                                // 115
                                                                                                              //
  Accounts.callLoginMethod({                                                                                  // 128
    methodName: 'createUser',                                                                                 // 129
    methodArguments: [options],                                                                               // 130
    userCallback: callback                                                                                    // 131
  });                                                                                                         //
};                                                                                                            //
                                                                                                              //
// Change password. Must be logged in.                                                                        //
//                                                                                                            //
// @param oldPassword {String|null} By default servers no longer allow                                        //
//   changing password without the old password, but they could so we                                         //
//   support passing no password to the server and letting it decide.                                         //
// @param newPassword {String}                                                                                //
// @param callback {Function(error|undefined)}                                                                //
                                                                                                              //
/**                                                                                                           //
 * @summary Change the current user's password. Must be logged in.                                            //
 * @locus Client                                                                                              //
 * @param {String} oldPassword The user's current password. This is __not__ sent in plain text over the wire.
 * @param {String} newPassword A new password for the user. This is __not__ sent in plain text over the wire.
 * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage accounts-base                                                                           //
 */                                                                                                           //
Accounts.changePassword = function (oldPassword, newPassword, callback) {                                     // 151
  if (!Meteor.user()) {                                                                                       // 152
    callback && callback(new Error("Must be logged in to change password."));                                 // 153
    return;                                                                                                   // 154
  }                                                                                                           //
                                                                                                              //
  check(newPassword, String);                                                                                 // 157
  if (!newPassword) {                                                                                         // 158
    callback(new Meteor.Error(400, "Password may not be empty"));                                             // 159
    return;                                                                                                   // 160
  }                                                                                                           //
                                                                                                              //
  Accounts.connection.apply('changePassword', [oldPassword ? Accounts._hashPassword(oldPassword) : null, Accounts._hashPassword(newPassword)], function (error, result) {
    if (error || !result) {                                                                                   // 168
      if (error && error.error === 400 && error.reason === 'old password format') {                           // 169
        // XXX COMPAT WITH 0.8.1.3                                                                            //
        // The server is telling us to upgrade from SRP to bcrypt, as                                         //
        // in Meteor.loginWithPassword.                                                                       //
        srpUpgradePath({                                                                                      // 174
          upgradeError: error,                                                                                // 175
          userSelector: { id: Meteor.userId() },                                                              // 176
          plaintextPassword: oldPassword                                                                      // 177
        }, function (err) {                                                                                   //
          if (err) {                                                                                          // 179
            callback && callback(err);                                                                        // 180
          } else {                                                                                            //
            // Now that we've successfully migrated from srp to                                               //
            // bcrypt, try changing the password again.                                                       //
            Accounts.changePassword(oldPassword, newPassword, callback);                                      // 184
          }                                                                                                   //
        });                                                                                                   //
      } else {                                                                                                //
        // A normal error, not an error telling us to upgrade to bcrypt                                       //
        callback && callback(error || new Error("No result from changePassword."));                           // 189
      }                                                                                                       //
    } else {                                                                                                  //
      callback && callback();                                                                                 // 193
    }                                                                                                         //
  });                                                                                                         //
};                                                                                                            //
                                                                                                              //
// Sends an email to a user with a link that can be used to reset                                             //
// their password                                                                                             //
//                                                                                                            //
// @param options {Object}                                                                                    //
//   - email: (email)                                                                                         //
// @param callback (optional) {Function(error|undefined)}                                                     //
                                                                                                              //
/**                                                                                                           //
 * @summary Request a forgot password email.                                                                  //
 * @locus Client                                                                                              //
 * @param {Object} options                                                                                    //
 * @param {String} options.email The email address to send a password reset link.                             //
 * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage accounts-base                                                                           //
 */                                                                                                           //
Accounts.forgotPassword = function (options, callback) {                                                      // 214
  if (!options.email) throw new Error("Must pass options.email");                                             // 215
  Accounts.connection.call("forgotPassword", options, callback);                                              // 217
};                                                                                                            //
                                                                                                              //
// Resets a password based on a token originally created by                                                   //
// Accounts.forgotPassword, and then logs in the matching user.                                               //
//                                                                                                            //
// @param token {String}                                                                                      //
// @param newPassword {String}                                                                                //
// @param callback (optional) {Function(error|undefined)}                                                     //
                                                                                                              //
/**                                                                                                           //
 * @summary Reset the password for a user using a token received in email. Logs the user in afterwards.       //
 * @locus Client                                                                                              //
 * @param {String} token The token retrieved from the reset password URL.                                     //
 * @param {String} newPassword A new password for the user. This is __not__ sent in plain text over the wire.
 * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage accounts-base                                                                           //
 */                                                                                                           //
Accounts.resetPassword = function (token, newPassword, callback) {                                            // 235
  check(token, String);                                                                                       // 236
  check(newPassword, String);                                                                                 // 237
                                                                                                              //
  if (!newPassword) {                                                                                         // 239
    callback(new Meteor.Error(400, "Password may not be empty"));                                             // 240
    return;                                                                                                   // 241
  }                                                                                                           //
                                                                                                              //
  Accounts.callLoginMethod({                                                                                  // 244
    methodName: 'resetPassword',                                                                              // 245
    methodArguments: [token, Accounts._hashPassword(newPassword)],                                            // 246
    userCallback: callback });                                                                                // 247
};                                                                                                            //
                                                                                                              //
// Verifies a user's email address based on a token originally                                                //
// created by Accounts.sendVerificationEmail                                                                  //
//                                                                                                            //
// @param token {String}                                                                                      //
// @param callback (optional) {Function(error|undefined)}                                                     //
                                                                                                              //
/**                                                                                                           //
 * @summary Marks the user's email address as verified. Logs the user in afterwards.                          //
 * @locus Client                                                                                              //
 * @param {String} token The token retrieved from the verification URL.                                       //
 * @param {Function} [callback] Optional callback. Called with no arguments on success, or with a single `Error` argument on failure.
 * @importFromPackage accounts-base                                                                           //
 */                                                                                                           //
Accounts.verifyEmail = function (token, callback) {                                                           // 263
  if (!token) throw new Error("Need to pass token");                                                          // 264
                                                                                                              //
  Accounts.callLoginMethod({                                                                                  // 267
    methodName: 'verifyEmail',                                                                                // 268
    methodArguments: [token],                                                                                 // 269
    userCallback: callback });                                                                                // 270
};                                                                                                            //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}}}},{"extensions":[".js",".json"]});
require("./node_modules/meteor/accounts-password/password_client.js");

/* Exports */
if (typeof Package === 'undefined') Package = {};
Package['accounts-password'] = {};

})();
