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
var check = Package.check.check;
var Match = Package.check.Match;
var Random = Package.random.Random;
var EJSON = Package.ejson.EJSON;
var _ = Package.underscore._;
var Tracker = Package.tracker.Tracker;
var Deps = Package.tracker.Deps;
var Retry = Package.retry.Retry;
var IdMap = Package['id-map'].IdMap;
var DDPCommon = Package['ddp-common'].DDPCommon;
var DiffSequence = Package['diff-sequence'].DiffSequence;
var MongoID = Package['mongo-id'].MongoID;

/* Package-scope variables */
var DDP, LivedataTest, MongoIDMap, SockJS, toSockjsUrl, toWebsocketUrl, allConnections;

(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/namespace.js                                                                               //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
/**                                                                                                               // 1
 * @namespace DDP                                                                                                 // 2
 * @summary Namespace for DDP-related methods/classes.                                                            // 3
 */                                                                                                               // 4
DDP          = {};                                                                                                // 5
LivedataTest = {};                                                                                                // 6
                                                                                                                  // 7
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/id_map.js                                                                                  //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
MongoIDMap = function () {                                                                                        // 1
  var self = this;                                                                                                // 2
  IdMap.call(self, MongoID.idStringify, MongoID.idParse);                                                         // 3
};                                                                                                                // 4
                                                                                                                  // 5
Meteor._inherits(MongoIDMap, IdMap);                                                                              // 6
                                                                                                                  // 7
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/sockjs-0.3.4.js                                                                            //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
// XXX METEOR changes in <METEOR>                                                                                 // 1
                                                                                                                  // 2
/* SockJS client, version 0.3.4, http://sockjs.org, MIT License                                                   // 3
                                                                                                                  // 4
Copyright (c) 2011-2012 VMware, Inc.                                                                              // 5
                                                                                                                  // 6
Permission is hereby granted, free of charge, to any person obtaining a copy                                      // 7
of this software and associated documentation files (the "Software"), to deal                                     // 8
in the Software without restriction, including without limitation the rights                                      // 9
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell                                         // 10
copies of the Software, and to permit persons to whom the Software is                                             // 11
furnished to do so, subject to the following conditions:                                                          // 12
                                                                                                                  // 13
The above copyright notice and this permission notice shall be included in                                        // 14
all copies or substantial portions of the Software.                                                               // 15
                                                                                                                  // 16
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR                                        // 17
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,                                          // 18
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE                                       // 19
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER                                            // 20
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,                                     // 21
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN                                         // 22
THE SOFTWARE.                                                                                                     // 23
*/                                                                                                                // 24
                                                                                                                  // 25
// <METEOR> Commented out JSO implementation (use json package instead).                                          // 26
// JSON2 by Douglas Crockford (minified).                                                                         // 27
// var JSON;JSON||(JSON={}),function(){function str(a,b){var c,d,e,f,g=gap,h,i=b[a];i&&typeof i=="object"&&typeof i.toJSON=="function"&&(i=i.toJSON(a)),typeof rep=="function"&&(i=rep.call(b,a,i));switch(typeof i){case"string":return quote(i);case"number":return isFinite(i)?String(i):"null";case"boolean":case"null":return String(i);case"object":if(!i)return"null";gap+=indent,h=[];if(Object.prototype.toString.apply(i)==="[object Array]"){f=i.length;for(c=0;c<f;c+=1)h[c]=str(c,i)||"null";e=h.length===0?"[]":gap?"[\n"+gap+h.join(",\n"+gap)+"\n"+g+"]":"["+h.join(",")+"]",gap=g;return e}if(rep&&typeof rep=="object"){f=rep.length;for(c=0;c<f;c+=1)typeof rep[c]=="string"&&(d=rep[c],e=str(d,i),e&&h.push(quote(d)+(gap?": ":":")+e))}else for(d in i)Object.prototype.hasOwnProperty.call(i,d)&&(e=str(d,i),e&&h.push(quote(d)+(gap?": ":":")+e));e=h.length===0?"{}":gap?"{\n"+gap+h.join(",\n"+gap)+"\n"+g+"}":"{"+h.join(",")+"}",gap=g;return e}}function quote(a){escapable.lastIndex=0;return escapable.test(a)?'"'+a.replace(escapable,function(a){var b=meta[a];return typeof b=="string"?b:"\\u"+("0000"+a.charCodeAt(0).toString(16)).slice(-4)})+'"':'"'+a+'"'}function f(a){return a<10?"0"+a:a}"use strict",typeof Date.prototype.toJSON!="function"&&(Date.prototype.toJSON=function(a){return isFinite(this.valueOf())?this.getUTCFullYear()+"-"+f(this.getUTCMonth()+1)+"-"+f(this.getUTCDate())+"T"+f(this.getUTCHours())+":"+f(this.getUTCMinutes())+":"+f(this.getUTCSeconds())+"Z":null},String.prototype.toJSON=Number.prototype.toJSON=Boolean.prototype.toJSON=function(a){return this.valueOf()});var cx=/[\u0000\u00ad\u0600-\u0604\u070f\u17b4\u17b5\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufeff\ufff0-\uffff]/g,escapable=/[\\\"\x00-\x1f\x7f-\x9f\u00ad\u0600-\u0604\u070f\u17b4\u17b5\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufeff\ufff0-\uffff]/g,gap,indent,meta={"\b":"\\b","\t":"\\t","\n":"\\n","\f":"\\f","\r":"\\r",'"':'\\"',"\\":"\\\\"},rep;typeof JSON.stringify!="function"&&(JSON.stringify=function(a,b,c){var d;gap="",indent="";if(typeof c=="number")for(d=0;d<c;d+=1)indent+=" ";else typeof c=="string"&&(indent=c);rep=b;if(!b||typeof b=="function"||typeof b=="object"&&typeof b.length=="number")return str("",{"":a});throw new Error("JSON.stringify")}),typeof JSON.parse!="function"&&(JSON.parse=function(text,reviver){function walk(a,b){var c,d,e=a[b];if(e&&typeof e=="object")for(c in e)Object.prototype.hasOwnProperty.call(e,c)&&(d=walk(e,c),d!==undefined?e[c]=d:delete e[c]);return reviver.call(a,b,e)}var j;text=String(text),cx.lastIndex=0,cx.test(text)&&(text=text.replace(cx,function(a){return"\\u"+("0000"+a.charCodeAt(0).toString(16)).slice(-4)}));if(/^[\],:{}\s]*$/.test(text.replace(/\\(?:["\\\/bfnrt]|u[0-9a-fA-F]{4})/g,"@").replace(/"[^"\\\n\r]*"|true|false|null|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?/g,"]").replace(/(?:^|:|,)(?:\s*\[)+/g,""))){j=eval("("+text+")");return typeof reviver=="function"?walk({"":j},""):j}throw new SyntaxError("JSON.parse")})}()
// </METEOR>                                                                                                      // 29
                                                                                                                  // 30
//     [*] Including lib/index.js                                                                                 // 31
// Public object                                                                                                  // 32
SockJS = (function(){                                                                                             // 33
              var _document = document;                                                                           // 34
              var _window = window;                                                                               // 35
              var utils = {};                                                                                     // 36
                                                                                                                  // 37
                                                                                                                  // 38
//         [*] Including lib/reventtarget.js                                                                      // 39
/*                                                                                                                // 40
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 41
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 42
 *                                                                                                                // 43
 * For the license see COPYING.                                                                                   // 44
 * ***** END LICENSE BLOCK *****                                                                                  // 45
 */                                                                                                               // 46
                                                                                                                  // 47
/* Simplified implementation of DOM2 EventTarget.                                                                 // 48
 *   http://www.w3.org/TR/DOM-Level-2-Events/events.html#Events-EventTarget                                       // 49
 */                                                                                                               // 50
var REventTarget = function() {};                                                                                 // 51
REventTarget.prototype.addEventListener = function (eventType, listener) {                                        // 52
    if(!this._listeners) {                                                                                        // 53
         this._listeners = {};                                                                                    // 54
    }                                                                                                             // 55
    if(!(eventType in this._listeners)) {                                                                         // 56
        this._listeners[eventType] = [];                                                                          // 57
    }                                                                                                             // 58
    var arr = this._listeners[eventType];                                                                         // 59
    if(utils.arrIndexOf(arr, listener) === -1) {                                                                  // 60
        arr.push(listener);                                                                                       // 61
    }                                                                                                             // 62
    return;                                                                                                       // 63
};                                                                                                                // 64
                                                                                                                  // 65
REventTarget.prototype.removeEventListener = function (eventType, listener) {                                     // 66
    if(!(this._listeners && (eventType in this._listeners))) {                                                    // 67
        return;                                                                                                   // 68
    }                                                                                                             // 69
    var arr = this._listeners[eventType];                                                                         // 70
    var idx = utils.arrIndexOf(arr, listener);                                                                    // 71
    if (idx !== -1) {                                                                                             // 72
        if(arr.length > 1) {                                                                                      // 73
            this._listeners[eventType] = arr.slice(0, idx).concat( arr.slice(idx+1) );                            // 74
        } else {                                                                                                  // 75
            delete this._listeners[eventType];                                                                    // 76
        }                                                                                                         // 77
        return;                                                                                                   // 78
    }                                                                                                             // 79
    return;                                                                                                       // 80
};                                                                                                                // 81
                                                                                                                  // 82
REventTarget.prototype.dispatchEvent = function (event) {                                                         // 83
    var t = event.type;                                                                                           // 84
    var args = Array.prototype.slice.call(arguments, 0);                                                          // 85
    if (this['on'+t]) {                                                                                           // 86
        this['on'+t].apply(this, args);                                                                           // 87
    }                                                                                                             // 88
    if (this._listeners && t in this._listeners) {                                                                // 89
        for(var i=0; i < this._listeners[t].length; i++) {                                                        // 90
            this._listeners[t][i].apply(this, args);                                                              // 91
        }                                                                                                         // 92
    }                                                                                                             // 93
};                                                                                                                // 94
//         [*] End of lib/reventtarget.js                                                                         // 95
                                                                                                                  // 96
                                                                                                                  // 97
//         [*] Including lib/simpleevent.js                                                                       // 98
/*                                                                                                                // 99
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 100
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 101
 *                                                                                                                // 102
 * For the license see COPYING.                                                                                   // 103
 * ***** END LICENSE BLOCK *****                                                                                  // 104
 */                                                                                                               // 105
                                                                                                                  // 106
var SimpleEvent = function(type, obj) {                                                                           // 107
    this.type = type;                                                                                             // 108
    if (typeof obj !== 'undefined') {                                                                             // 109
        for(var k in obj) {                                                                                       // 110
            if (!obj.hasOwnProperty(k)) continue;                                                                 // 111
            this[k] = obj[k];                                                                                     // 112
        }                                                                                                         // 113
    }                                                                                                             // 114
};                                                                                                                // 115
                                                                                                                  // 116
SimpleEvent.prototype.toString = function() {                                                                     // 117
    var r = [];                                                                                                   // 118
    for(var k in this) {                                                                                          // 119
        if (!this.hasOwnProperty(k)) continue;                                                                    // 120
        var v = this[k];                                                                                          // 121
        if (typeof v === 'function') v = '[function]';                                                            // 122
        r.push(k + '=' + v);                                                                                      // 123
    }                                                                                                             // 124
    return 'SimpleEvent(' + r.join(', ') + ')';                                                                   // 125
};                                                                                                                // 126
//         [*] End of lib/simpleevent.js                                                                          // 127
                                                                                                                  // 128
                                                                                                                  // 129
//         [*] Including lib/eventemitter.js                                                                      // 130
/*                                                                                                                // 131
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 132
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 133
 *                                                                                                                // 134
 * For the license see COPYING.                                                                                   // 135
 * ***** END LICENSE BLOCK *****                                                                                  // 136
 */                                                                                                               // 137
                                                                                                                  // 138
var EventEmitter = function(events) {                                                                             // 139
    var that = this;                                                                                              // 140
    that._events = events || [];                                                                                  // 141
    that._listeners = {};                                                                                         // 142
};                                                                                                                // 143
EventEmitter.prototype.emit = function(type) {                                                                    // 144
    var that = this;                                                                                              // 145
    that._verifyType(type);                                                                                       // 146
    if (that._nuked) return;                                                                                      // 147
                                                                                                                  // 148
    var args = Array.prototype.slice.call(arguments, 1);                                                          // 149
    if (that['on'+type]) {                                                                                        // 150
        that['on'+type].apply(that, args);                                                                        // 151
    }                                                                                                             // 152
    if (type in that._listeners) {                                                                                // 153
        for(var i = 0; i < that._listeners[type].length; i++) {                                                   // 154
            that._listeners[type][i].apply(that, args);                                                           // 155
        }                                                                                                         // 156
    }                                                                                                             // 157
};                                                                                                                // 158
                                                                                                                  // 159
EventEmitter.prototype.on = function(type, callback) {                                                            // 160
    var that = this;                                                                                              // 161
    that._verifyType(type);                                                                                       // 162
    if (that._nuked) return;                                                                                      // 163
                                                                                                                  // 164
    if (!(type in that._listeners)) {                                                                             // 165
        that._listeners[type] = [];                                                                               // 166
    }                                                                                                             // 167
    that._listeners[type].push(callback);                                                                         // 168
};                                                                                                                // 169
                                                                                                                  // 170
EventEmitter.prototype._verifyType = function(type) {                                                             // 171
    var that = this;                                                                                              // 172
    if (utils.arrIndexOf(that._events, type) === -1) {                                                            // 173
        utils.log('Event ' + JSON.stringify(type) +                                                               // 174
                  ' not listed ' + JSON.stringify(that._events) +                                                 // 175
                  ' in ' + that);                                                                                 // 176
    }                                                                                                             // 177
};                                                                                                                // 178
                                                                                                                  // 179
EventEmitter.prototype.nuke = function() {                                                                        // 180
    var that = this;                                                                                              // 181
    that._nuked = true;                                                                                           // 182
    for(var i=0; i<that._events.length; i++) {                                                                    // 183
        delete that[that._events[i]];                                                                             // 184
    }                                                                                                             // 185
    that._listeners = {};                                                                                         // 186
};                                                                                                                // 187
//         [*] End of lib/eventemitter.js                                                                         // 188
                                                                                                                  // 189
                                                                                                                  // 190
//         [*] Including lib/utils.js                                                                             // 191
/*                                                                                                                // 192
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 193
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 194
 *                                                                                                                // 195
 * For the license see COPYING.                                                                                   // 196
 * ***** END LICENSE BLOCK *****                                                                                  // 197
 */                                                                                                               // 198
                                                                                                                  // 199
var random_string_chars = 'abcdefghijklmnopqrstuvwxyz0123456789_';                                                // 200
utils.random_string = function(length, max) {                                                                     // 201
    max = max || random_string_chars.length;                                                                      // 202
    var i, ret = [];                                                                                              // 203
    for(i=0; i < length; i++) {                                                                                   // 204
        ret.push( random_string_chars.substr(Math.floor(Math.random() * max),1) );                                // 205
    }                                                                                                             // 206
    return ret.join('');                                                                                          // 207
};                                                                                                                // 208
utils.random_number = function(max) {                                                                             // 209
    return Math.floor(Math.random() * max);                                                                       // 210
};                                                                                                                // 211
utils.random_number_string = function(max) {                                                                      // 212
    var t = (''+(max - 1)).length;                                                                                // 213
    var p = Array(t+1).join('0');                                                                                 // 214
    return (p + utils.random_number(max)).slice(-t);                                                              // 215
};                                                                                                                // 216
                                                                                                                  // 217
// Assuming that url looks like: http://asdasd:111/asd                                                            // 218
utils.getOrigin = function(url) {                                                                                 // 219
    url += '/';                                                                                                   // 220
    var parts = url.split('/').slice(0, 3);                                                                       // 221
    return parts.join('/');                                                                                       // 222
};                                                                                                                // 223
                                                                                                                  // 224
utils.isSameOriginUrl = function(url_a, url_b) {                                                                  // 225
    // location.origin would do, but it's not always available.                                                   // 226
    if (!url_b) url_b = _window.location.href;                                                                    // 227
                                                                                                                  // 228
    return (url_a.split('/').slice(0,3).join('/')                                                                 // 229
                ===                                                                                               // 230
            url_b.split('/').slice(0,3).join('/'));                                                               // 231
};                                                                                                                // 232
                                                                                                                  // 233
// <METEOR>                                                                                                       // 234
// https://github.com/sockjs/sockjs-client/issues/79                                                              // 235
utils.isSameOriginScheme = function(url_a, url_b) {                                                               // 236
    if (!url_b) url_b = _window.location.href;                                                                    // 237
                                                                                                                  // 238
    return (url_a.split(':')[0]                                                                                   // 239
                ===                                                                                               // 240
            url_b.split(':')[0]);                                                                                 // 241
};                                                                                                                // 242
// </METEOR>                                                                                                      // 243
                                                                                                                  // 244
                                                                                                                  // 245
utils.getParentDomain = function(url) {                                                                           // 246
    // ipv4 ip address                                                                                            // 247
    if (/^[0-9.]*$/.test(url)) return url;                                                                        // 248
    // ipv6 ip address                                                                                            // 249
    if (/^\[/.test(url)) return url;                                                                              // 250
    // no dots                                                                                                    // 251
    if (!(/[.]/.test(url))) return url;                                                                           // 252
                                                                                                                  // 253
    var parts = url.split('.').slice(1);                                                                          // 254
    return parts.join('.');                                                                                       // 255
};                                                                                                                // 256
                                                                                                                  // 257
utils.objectExtend = function(dst, src) {                                                                         // 258
    for(var k in src) {                                                                                           // 259
        if (src.hasOwnProperty(k)) {                                                                              // 260
            dst[k] = src[k];                                                                                      // 261
        }                                                                                                         // 262
    }                                                                                                             // 263
    return dst;                                                                                                   // 264
};                                                                                                                // 265
                                                                                                                  // 266
var WPrefix = '_jp';                                                                                              // 267
                                                                                                                  // 268
utils.polluteGlobalNamespace = function() {                                                                       // 269
    if (!(WPrefix in _window)) {                                                                                  // 270
        _window[WPrefix] = {};                                                                                    // 271
    }                                                                                                             // 272
};                                                                                                                // 273
                                                                                                                  // 274
utils.closeFrame = function (code, reason) {                                                                      // 275
    return 'c'+JSON.stringify([code, reason]);                                                                    // 276
};                                                                                                                // 277
                                                                                                                  // 278
utils.userSetCode = function (code) {                                                                             // 279
    return code === 1000 || (code >= 3000 && code <= 4999);                                                       // 280
};                                                                                                                // 281
                                                                                                                  // 282
// See: http://www.erg.abdn.ac.uk/~gerrit/dccp/notes/ccid2/rto_estimator/                                         // 283
// and RFC 2988.                                                                                                  // 284
utils.countRTO = function (rtt) {                                                                                 // 285
    var rto;                                                                                                      // 286
    if (rtt > 100) {                                                                                              // 287
        rto = 3 * rtt; // rto > 300msec                                                                           // 288
    } else {                                                                                                      // 289
        rto = rtt + 200; // 200msec < rto <= 300msec                                                              // 290
    }                                                                                                             // 291
    return rto;                                                                                                   // 292
}                                                                                                                 // 293
                                                                                                                  // 294
utils.log = function() {                                                                                          // 295
    if (_window.console && console.log && console.log.apply) {                                                    // 296
        console.log.apply(console, arguments);                                                                    // 297
    }                                                                                                             // 298
};                                                                                                                // 299
                                                                                                                  // 300
utils.bind = function(fun, that) {                                                                                // 301
    if (fun.bind) {                                                                                               // 302
        return fun.bind(that);                                                                                    // 303
    } else {                                                                                                      // 304
        return function() {                                                                                       // 305
            return fun.apply(that, arguments);                                                                    // 306
        };                                                                                                        // 307
    }                                                                                                             // 308
};                                                                                                                // 309
                                                                                                                  // 310
utils.flatUrl = function(url) {                                                                                   // 311
    return url.indexOf('?') === -1 && url.indexOf('#') === -1;                                                    // 312
};                                                                                                                // 313
                                                                                                                  // 314
// `relativeTo` is an optional absolute URL. If provided, `url` will be                                           // 315
// interpreted relative to `relativeTo`. Defaults to `document.location`.                                         // 316
// <METEOR>                                                                                                       // 317
utils.amendUrl = function(url, relativeTo) {                                                                      // 318
    var baseUrl;                                                                                                  // 319
    if (relativeTo === undefined) {                                                                               // 320
      baseUrl = _document.location;                                                                               // 321
    } else {                                                                                                      // 322
      var protocolMatch = /^([a-z0-9.+-]+:)/i.exec(relativeTo);                                                   // 323
      if (protocolMatch) {                                                                                        // 324
        var protocol = protocolMatch[0].toLowerCase();                                                            // 325
        var rest = relativeTo.substring(protocol.length);                                                         // 326
        var hostMatch = /[a-z0-9\.-]+(:[0-9]+)?/.exec(rest);                                                      // 327
        if (hostMatch)                                                                                            // 328
          var host = hostMatch[0];                                                                                // 329
      }                                                                                                           // 330
      if (! protocol || ! host)                                                                                   // 331
        throw new Error("relativeTo must be an absolute url");                                                    // 332
      baseUrl = {                                                                                                 // 333
        protocol: protocol,                                                                                       // 334
        host: host                                                                                                // 335
      };                                                                                                          // 336
    }                                                                                                             // 337
    if (!url) {                                                                                                   // 338
        throw new Error('Wrong url for SockJS');                                                                  // 339
    }                                                                                                             // 340
    if (!utils.flatUrl(url)) {                                                                                    // 341
        throw new Error('Only basic urls are supported in SockJS');                                               // 342
    }                                                                                                             // 343
                                                                                                                  // 344
    //  '//abc' --> 'http://abc'                                                                                  // 345
    if (url.indexOf('//') === 0) {                                                                                // 346
        url = baseUrl.protocol + url;                                                                             // 347
    }                                                                                                             // 348
    // '/abc' --> 'http://localhost:1234/abc'                                                                     // 349
    if (url.indexOf('/') === 0) {                                                                                 // 350
        url = baseUrl.protocol + '//' + baseUrl.host + url;                                                       // 351
    }                                                                                                             // 352
    // </METEOR>                                                                                                  // 353
    // strip trailing slashes                                                                                     // 354
    url = url.replace(/[/]+$/,'');                                                                                // 355
                                                                                                                  // 356
    // We have a full url here, with proto and host. For some browsers                                            // 357
    // http://localhost:80/ is not in the same origin as http://localhost/                                        // 358
	// Remove explicit :80 or :443 in such cases. See #74                                                            // 359
    var parts = url.split("/");                                                                                   // 360
    if ((parts[0] === "http:" && /:80$/.test(parts[2])) ||                                                        // 361
	    (parts[0] === "https:" && /:443$/.test(parts[2]))) {                                                         // 362
		parts[2] = parts[2].replace(/:(80|443)$/, "");                                                                  // 363
	}                                                                                                                // 364
    url = parts.join("/");                                                                                        // 365
    return url;                                                                                                   // 366
};                                                                                                                // 367
                                                                                                                  // 368
// IE doesn't support [].indexOf.                                                                                 // 369
utils.arrIndexOf = function(arr, obj){                                                                            // 370
    for(var i=0; i < arr.length; i++){                                                                            // 371
        if(arr[i] === obj){                                                                                       // 372
            return i;                                                                                             // 373
        }                                                                                                         // 374
    }                                                                                                             // 375
    return -1;                                                                                                    // 376
};                                                                                                                // 377
                                                                                                                  // 378
utils.arrSkip = function(arr, obj) {                                                                              // 379
    var idx = utils.arrIndexOf(arr, obj);                                                                         // 380
    if (idx === -1) {                                                                                             // 381
        return arr.slice();                                                                                       // 382
    } else {                                                                                                      // 383
        var dst = arr.slice(0, idx);                                                                              // 384
        return dst.concat(arr.slice(idx+1));                                                                      // 385
    }                                                                                                             // 386
};                                                                                                                // 387
                                                                                                                  // 388
// Via: https://gist.github.com/1133122/2121c601c5549155483f50be3da5305e83b8c5df                                  // 389
utils.isArray = Array.isArray || function(value) {                                                                // 390
    return {}.toString.call(value).indexOf('Array') >= 0                                                          // 391
};                                                                                                                // 392
                                                                                                                  // 393
utils.delay = function(t, fun) {                                                                                  // 394
    if(typeof t === 'function') {                                                                                 // 395
        fun = t;                                                                                                  // 396
        t = 0;                                                                                                    // 397
    }                                                                                                             // 398
    return setTimeout(fun, t);                                                                                    // 399
};                                                                                                                // 400
                                                                                                                  // 401
                                                                                                                  // 402
// Chars worth escaping, as defined by Douglas Crockford:                                                         // 403
//   https://github.com/douglascrockford/JSON-js/blob/47a9882cddeb1e8529e07af9736218075372b8ac/json2.js#L196      // 404
var json_escapable = /[\\\"\x00-\x1f\x7f-\x9f\u00ad\u0600-\u0604\u070f\u17b4\u17b5\u200c-\u200f\u2028-\u202f\u2060-\u206f\ufeff\ufff0-\uffff]/g,
    json_lookup = {                                                                                               // 406
"\u0000":"\\u0000","\u0001":"\\u0001","\u0002":"\\u0002","\u0003":"\\u0003",                                      // 407
"\u0004":"\\u0004","\u0005":"\\u0005","\u0006":"\\u0006","\u0007":"\\u0007",                                      // 408
"\b":"\\b","\t":"\\t","\n":"\\n","\u000b":"\\u000b","\f":"\\f","\r":"\\r",                                        // 409
"\u000e":"\\u000e","\u000f":"\\u000f","\u0010":"\\u0010","\u0011":"\\u0011",                                      // 410
"\u0012":"\\u0012","\u0013":"\\u0013","\u0014":"\\u0014","\u0015":"\\u0015",                                      // 411
"\u0016":"\\u0016","\u0017":"\\u0017","\u0018":"\\u0018","\u0019":"\\u0019",                                      // 412
"\u001a":"\\u001a","\u001b":"\\u001b","\u001c":"\\u001c","\u001d":"\\u001d",                                      // 413
"\u001e":"\\u001e","\u001f":"\\u001f","\"":"\\\"","\\":"\\\\",                                                    // 414
"\u007f":"\\u007f","\u0080":"\\u0080","\u0081":"\\u0081","\u0082":"\\u0082",                                      // 415
"\u0083":"\\u0083","\u0084":"\\u0084","\u0085":"\\u0085","\u0086":"\\u0086",                                      // 416
"\u0087":"\\u0087","\u0088":"\\u0088","\u0089":"\\u0089","\u008a":"\\u008a",                                      // 417
"\u008b":"\\u008b","\u008c":"\\u008c","\u008d":"\\u008d","\u008e":"\\u008e",                                      // 418
"\u008f":"\\u008f","\u0090":"\\u0090","\u0091":"\\u0091","\u0092":"\\u0092",                                      // 419
"\u0093":"\\u0093","\u0094":"\\u0094","\u0095":"\\u0095","\u0096":"\\u0096",                                      // 420
"\u0097":"\\u0097","\u0098":"\\u0098","\u0099":"\\u0099","\u009a":"\\u009a",                                      // 421
"\u009b":"\\u009b","\u009c":"\\u009c","\u009d":"\\u009d","\u009e":"\\u009e",                                      // 422
"\u009f":"\\u009f","\u00ad":"\\u00ad","\u0600":"\\u0600","\u0601":"\\u0601",                                      // 423
"\u0602":"\\u0602","\u0603":"\\u0603","\u0604":"\\u0604","\u070f":"\\u070f",                                      // 424
"\u17b4":"\\u17b4","\u17b5":"\\u17b5","\u200c":"\\u200c","\u200d":"\\u200d",                                      // 425
"\u200e":"\\u200e","\u200f":"\\u200f","\u2028":"\\u2028","\u2029":"\\u2029",                                      // 426
"\u202a":"\\u202a","\u202b":"\\u202b","\u202c":"\\u202c","\u202d":"\\u202d",                                      // 427
"\u202e":"\\u202e","\u202f":"\\u202f","\u2060":"\\u2060","\u2061":"\\u2061",                                      // 428
"\u2062":"\\u2062","\u2063":"\\u2063","\u2064":"\\u2064","\u2065":"\\u2065",                                      // 429
"\u2066":"\\u2066","\u2067":"\\u2067","\u2068":"\\u2068","\u2069":"\\u2069",                                      // 430
"\u206a":"\\u206a","\u206b":"\\u206b","\u206c":"\\u206c","\u206d":"\\u206d",                                      // 431
"\u206e":"\\u206e","\u206f":"\\u206f","\ufeff":"\\ufeff","\ufff0":"\\ufff0",                                      // 432
"\ufff1":"\\ufff1","\ufff2":"\\ufff2","\ufff3":"\\ufff3","\ufff4":"\\ufff4",                                      // 433
"\ufff5":"\\ufff5","\ufff6":"\\ufff6","\ufff7":"\\ufff7","\ufff8":"\\ufff8",                                      // 434
"\ufff9":"\\ufff9","\ufffa":"\\ufffa","\ufffb":"\\ufffb","\ufffc":"\\ufffc",                                      // 435
"\ufffd":"\\ufffd","\ufffe":"\\ufffe","\uffff":"\\uffff"};                                                        // 436
                                                                                                                  // 437
// Some extra characters that Chrome gets wrong, and substitutes with                                             // 438
// something else on the wire.                                                                                    // 439
var extra_escapable = /[\x00-\x1f\ud800-\udfff\ufffe\uffff\u0300-\u0333\u033d-\u0346\u034a-\u034c\u0350-\u0352\u0357-\u0358\u035c-\u0362\u0374\u037e\u0387\u0591-\u05af\u05c4\u0610-\u0617\u0653-\u0654\u0657-\u065b\u065d-\u065e\u06df-\u06e2\u06eb-\u06ec\u0730\u0732-\u0733\u0735-\u0736\u073a\u073d\u073f-\u0741\u0743\u0745\u0747\u07eb-\u07f1\u0951\u0958-\u095f\u09dc-\u09dd\u09df\u0a33\u0a36\u0a59-\u0a5b\u0a5e\u0b5c-\u0b5d\u0e38-\u0e39\u0f43\u0f4d\u0f52\u0f57\u0f5c\u0f69\u0f72-\u0f76\u0f78\u0f80-\u0f83\u0f93\u0f9d\u0fa2\u0fa7\u0fac\u0fb9\u1939-\u193a\u1a17\u1b6b\u1cda-\u1cdb\u1dc0-\u1dcf\u1dfc\u1dfe\u1f71\u1f73\u1f75\u1f77\u1f79\u1f7b\u1f7d\u1fbb\u1fbe\u1fc9\u1fcb\u1fd3\u1fdb\u1fe3\u1feb\u1fee-\u1fef\u1ff9\u1ffb\u1ffd\u2000-\u2001\u20d0-\u20d1\u20d4-\u20d7\u20e7-\u20e9\u2126\u212a-\u212b\u2329-\u232a\u2adc\u302b-\u302c\uaab2-\uaab3\uf900-\ufa0d\ufa10\ufa12\ufa15-\ufa1e\ufa20\ufa22\ufa25-\ufa26\ufa2a-\ufa2d\ufa30-\ufa6d\ufa70-\ufad9\ufb1d\ufb1f\ufb2a-\ufb36\ufb38-\ufb3c\ufb3e\ufb40-\ufb41\ufb43-\ufb44\ufb46-\ufb4e\ufff0-\uffff]/g,
    extra_lookup;                                                                                                 // 441
                                                                                                                  // 442
// JSON Quote string. Use native implementation when possible.                                                    // 443
var JSONQuote = (JSON && JSON.stringify) || function(string) {                                                    // 444
    json_escapable.lastIndex = 0;                                                                                 // 445
    if (json_escapable.test(string)) {                                                                            // 446
        string = string.replace(json_escapable, function(a) {                                                     // 447
            return json_lookup[a];                                                                                // 448
        });                                                                                                       // 449
    }                                                                                                             // 450
    return '"' + string + '"';                                                                                    // 451
};                                                                                                                // 452
                                                                                                                  // 453
// This may be quite slow, so let's delay until user actually uses bad                                            // 454
// characters.                                                                                                    // 455
var unroll_lookup = function(escapable) {                                                                         // 456
    var i;                                                                                                        // 457
    var unrolled = {}                                                                                             // 458
    var c = []                                                                                                    // 459
    for(i=0; i<65536; i++) {                                                                                      // 460
        c.push( String.fromCharCode(i) );                                                                         // 461
    }                                                                                                             // 462
    escapable.lastIndex = 0;                                                                                      // 463
    c.join('').replace(escapable, function (a) {                                                                  // 464
        unrolled[ a ] = '\\u' + ('0000' + a.charCodeAt(0).toString(16)).slice(-4);                                // 465
        return '';                                                                                                // 466
    });                                                                                                           // 467
    escapable.lastIndex = 0;                                                                                      // 468
    return unrolled;                                                                                              // 469
};                                                                                                                // 470
                                                                                                                  // 471
// Quote string, also taking care of unicode characters that browsers                                             // 472
// often break. Especially, take care of unicode surrogates:                                                      // 473
//    http://en.wikipedia.org/wiki/Mapping_of_Unicode_characters#Surrogates                                       // 474
utils.quote = function(string) {                                                                                  // 475
    var quoted = JSONQuote(string);                                                                               // 476
                                                                                                                  // 477
    // In most cases this should be very fast and good enough.                                                    // 478
    extra_escapable.lastIndex = 0;                                                                                // 479
    if(!extra_escapable.test(quoted)) {                                                                           // 480
        return quoted;                                                                                            // 481
    }                                                                                                             // 482
                                                                                                                  // 483
    if(!extra_lookup) extra_lookup = unroll_lookup(extra_escapable);                                              // 484
                                                                                                                  // 485
    return quoted.replace(extra_escapable, function(a) {                                                          // 486
        return extra_lookup[a];                                                                                   // 487
    });                                                                                                           // 488
}                                                                                                                 // 489
                                                                                                                  // 490
var _all_protocols = ['websocket',                                                                                // 491
                      'xdr-streaming',                                                                            // 492
                      'xhr-streaming',                                                                            // 493
                      'iframe-eventsource',                                                                       // 494
                      'iframe-htmlfile',                                                                          // 495
                      'xdr-polling',                                                                              // 496
                      'xhr-polling',                                                                              // 497
                      'iframe-xhr-polling',                                                                       // 498
                      'jsonp-polling'];                                                                           // 499
                                                                                                                  // 500
utils.probeProtocols = function() {                                                                               // 501
    var probed = {};                                                                                              // 502
    for(var i=0; i<_all_protocols.length; i++) {                                                                  // 503
        var protocol = _all_protocols[i];                                                                         // 504
        // User can have a typo in protocol name.                                                                 // 505
        probed[protocol] = SockJS[protocol] &&                                                                    // 506
                           SockJS[protocol].enabled();                                                            // 507
    }                                                                                                             // 508
    return probed;                                                                                                // 509
};                                                                                                                // 510
                                                                                                                  // 511
utils.detectProtocols = function(probed, protocols_whitelist, info) {                                             // 512
    var pe = {},                                                                                                  // 513
        protocols = [];                                                                                           // 514
    if (!protocols_whitelist) protocols_whitelist = _all_protocols;                                               // 515
    for(var i=0; i<protocols_whitelist.length; i++) {                                                             // 516
        var protocol = protocols_whitelist[i];                                                                    // 517
        pe[protocol] = probed[protocol];                                                                          // 518
    }                                                                                                             // 519
    var maybe_push = function(protos) {                                                                           // 520
        var proto = protos.shift();                                                                               // 521
        if (pe[proto]) {                                                                                          // 522
            protocols.push(proto);                                                                                // 523
        } else {                                                                                                  // 524
            if (protos.length > 0) {                                                                              // 525
                maybe_push(protos);                                                                               // 526
            }                                                                                                     // 527
        }                                                                                                         // 528
    }                                                                                                             // 529
                                                                                                                  // 530
    // 1. Websocket                                                                                               // 531
    if (info.websocket !== false) {                                                                               // 532
        maybe_push(['websocket']);                                                                                // 533
    }                                                                                                             // 534
                                                                                                                  // 535
    // 2. Streaming                                                                                               // 536
    if (pe['xhr-streaming'] && !info.null_origin) {                                                               // 537
        protocols.push('xhr-streaming');                                                                          // 538
    } else {                                                                                                      // 539
        if (pe['xdr-streaming'] && !info.cookie_needed && !info.null_origin) {                                    // 540
            protocols.push('xdr-streaming');                                                                      // 541
        } else {                                                                                                  // 542
            maybe_push(['iframe-eventsource',                                                                     // 543
                        'iframe-htmlfile']);                                                                      // 544
        }                                                                                                         // 545
    }                                                                                                             // 546
                                                                                                                  // 547
    // 3. Polling                                                                                                 // 548
    if (pe['xhr-polling'] && !info.null_origin) {                                                                 // 549
        protocols.push('xhr-polling');                                                                            // 550
    } else {                                                                                                      // 551
        if (pe['xdr-polling'] && !info.cookie_needed && !info.null_origin) {                                      // 552
            protocols.push('xdr-polling');                                                                        // 553
        } else {                                                                                                  // 554
            maybe_push(['iframe-xhr-polling',                                                                     // 555
                        'jsonp-polling']);                                                                        // 556
        }                                                                                                         // 557
    }                                                                                                             // 558
    return protocols;                                                                                             // 559
}                                                                                                                 // 560
//         [*] End of lib/utils.js                                                                                // 561
                                                                                                                  // 562
                                                                                                                  // 563
//         [*] Including lib/dom.js                                                                               // 564
/*                                                                                                                // 565
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 566
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 567
 *                                                                                                                // 568
 * For the license see COPYING.                                                                                   // 569
 * ***** END LICENSE BLOCK *****                                                                                  // 570
 */                                                                                                               // 571
                                                                                                                  // 572
// May be used by htmlfile jsonp and transports.                                                                  // 573
var MPrefix = '_sockjs_global';                                                                                   // 574
utils.createHook = function() {                                                                                   // 575
    var window_id = 'a' + utils.random_string(8);                                                                 // 576
    if (!(MPrefix in _window)) {                                                                                  // 577
        var map = {};                                                                                             // 578
        _window[MPrefix] = function(window_id) {                                                                  // 579
            if (!(window_id in map)) {                                                                            // 580
                map[window_id] = {                                                                                // 581
                    id: window_id,                                                                                // 582
                    del: function() {delete map[window_id];}                                                      // 583
                };                                                                                                // 584
            }                                                                                                     // 585
            return map[window_id];                                                                                // 586
        }                                                                                                         // 587
    }                                                                                                             // 588
    return _window[MPrefix](window_id);                                                                           // 589
};                                                                                                                // 590
                                                                                                                  // 591
                                                                                                                  // 592
                                                                                                                  // 593
utils.attachMessage = function(listener) {                                                                        // 594
    utils.attachEvent('message', listener);                                                                       // 595
};                                                                                                                // 596
utils.attachEvent = function(event, listener) {                                                                   // 597
    if (typeof _window.addEventListener !== 'undefined') {                                                        // 598
        _window.addEventListener(event, listener, false);                                                         // 599
    } else {                                                                                                      // 600
        // IE quirks.                                                                                             // 601
        // According to: http://stevesouders.com/misc/test-postmessage.php                                        // 602
        // the message gets delivered only to 'document', not 'window'.                                           // 603
        _document.attachEvent("on" + event, listener);                                                            // 604
        // I get 'window' for ie8.                                                                                // 605
        _window.attachEvent("on" + event, listener);                                                              // 606
    }                                                                                                             // 607
};                                                                                                                // 608
                                                                                                                  // 609
utils.detachMessage = function(listener) {                                                                        // 610
    utils.detachEvent('message', listener);                                                                       // 611
};                                                                                                                // 612
utils.detachEvent = function(event, listener) {                                                                   // 613
    if (typeof _window.addEventListener !== 'undefined') {                                                        // 614
        _window.removeEventListener(event, listener, false);                                                      // 615
    } else {                                                                                                      // 616
        _document.detachEvent("on" + event, listener);                                                            // 617
        _window.detachEvent("on" + event, listener);                                                              // 618
    }                                                                                                             // 619
};                                                                                                                // 620
                                                                                                                  // 621
                                                                                                                  // 622
var on_unload = {};                                                                                               // 623
// Things registered after beforeunload are to be called immediately.                                             // 624
var after_unload = false;                                                                                         // 625
                                                                                                                  // 626
var trigger_unload_callbacks = function() {                                                                       // 627
    for(var ref in on_unload) {                                                                                   // 628
        on_unload[ref]();                                                                                         // 629
        delete on_unload[ref];                                                                                    // 630
    };                                                                                                            // 631
};                                                                                                                // 632
                                                                                                                  // 633
var unload_triggered = function() {                                                                               // 634
    if(after_unload) return;                                                                                      // 635
    after_unload = true;                                                                                          // 636
    trigger_unload_callbacks();                                                                                   // 637
};                                                                                                                // 638
                                                                                                                  // 639
// 'unload' alone is not reliable in opera within an iframe, but we                                               // 640
// can't use `beforeunload` as IE fires it on javascript: links.                                                  // 641
utils.attachEvent('unload', unload_triggered);                                                                    // 642
                                                                                                                  // 643
utils.unload_add = function(listener) {                                                                           // 644
    var ref = utils.random_string(8);                                                                             // 645
    on_unload[ref] = listener;                                                                                    // 646
    if (after_unload) {                                                                                           // 647
        utils.delay(trigger_unload_callbacks);                                                                    // 648
    }                                                                                                             // 649
    return ref;                                                                                                   // 650
};                                                                                                                // 651
utils.unload_del = function(ref) {                                                                                // 652
    if (ref in on_unload)                                                                                         // 653
        delete on_unload[ref];                                                                                    // 654
};                                                                                                                // 655
                                                                                                                  // 656
                                                                                                                  // 657
utils.createIframe = function (iframe_url, error_callback) {                                                      // 658
    var iframe = _document.createElement('iframe');                                                               // 659
    var tref, unload_ref;                                                                                         // 660
    var unattach = function() {                                                                                   // 661
        clearTimeout(tref);                                                                                       // 662
        // Explorer had problems with that.                                                                       // 663
        try {iframe.onload = null;} catch (x) {}                                                                  // 664
        iframe.onerror = null;                                                                                    // 665
    };                                                                                                            // 666
    var cleanup = function() {                                                                                    // 667
        if (iframe) {                                                                                             // 668
            unattach();                                                                                           // 669
            // This timeout makes chrome fire onbeforeunload event                                                // 670
            // within iframe. Without the timeout it goes straight to                                             // 671
            // onunload.                                                                                          // 672
            setTimeout(function() {                                                                               // 673
                if(iframe) {                                                                                      // 674
                    iframe.parentNode.removeChild(iframe);                                                        // 675
                }                                                                                                 // 676
                iframe = null;                                                                                    // 677
            }, 0);                                                                                                // 678
            utils.unload_del(unload_ref);                                                                         // 679
        }                                                                                                         // 680
    };                                                                                                            // 681
    var onerror = function(r) {                                                                                   // 682
        if (iframe) {                                                                                             // 683
            cleanup();                                                                                            // 684
            error_callback(r);                                                                                    // 685
        }                                                                                                         // 686
    };                                                                                                            // 687
    var post = function(msg, origin) {                                                                            // 688
        try {                                                                                                     // 689
            // When the iframe is not loaded, IE raises an exception                                              // 690
            // on 'contentWindow'.                                                                                // 691
            if (iframe && iframe.contentWindow) {                                                                 // 692
                iframe.contentWindow.postMessage(msg, origin);                                                    // 693
            }                                                                                                     // 694
        } catch (x) {};                                                                                           // 695
    };                                                                                                            // 696
                                                                                                                  // 697
    iframe.src = iframe_url;                                                                                      // 698
    iframe.style.display = 'none';                                                                                // 699
    iframe.style.position = 'absolute';                                                                           // 700
    iframe.onerror = function(){onerror('onerror');};                                                             // 701
    iframe.onload = function() {                                                                                  // 702
        // `onload` is triggered before scripts on the iframe are                                                 // 703
        // executed. Give it few seconds to actually load stuff.                                                  // 704
        clearTimeout(tref);                                                                                       // 705
        tref = setTimeout(function(){onerror('onload timeout');}, 2000);                                          // 706
    };                                                                                                            // 707
    _document.body.appendChild(iframe);                                                                           // 708
    tref = setTimeout(function(){onerror('timeout');}, 15000);                                                    // 709
    unload_ref = utils.unload_add(cleanup);                                                                       // 710
    return {                                                                                                      // 711
        post: post,                                                                                               // 712
        cleanup: cleanup,                                                                                         // 713
        loaded: unattach                                                                                          // 714
    };                                                                                                            // 715
};                                                                                                                // 716
                                                                                                                  // 717
utils.createHtmlfile = function (iframe_url, error_callback) {                                                    // 718
    var doc = new ActiveXObject('htmlfile');                                                                      // 719
    var tref, unload_ref;                                                                                         // 720
    var iframe;                                                                                                   // 721
    var unattach = function() {                                                                                   // 722
        clearTimeout(tref);                                                                                       // 723
    };                                                                                                            // 724
    var cleanup = function() {                                                                                    // 725
        if (doc) {                                                                                                // 726
            unattach();                                                                                           // 727
            utils.unload_del(unload_ref);                                                                         // 728
            iframe.parentNode.removeChild(iframe);                                                                // 729
            iframe = doc = null;                                                                                  // 730
            CollectGarbage();                                                                                     // 731
        }                                                                                                         // 732
    };                                                                                                            // 733
    var onerror = function(r)  {                                                                                  // 734
        if (doc) {                                                                                                // 735
            cleanup();                                                                                            // 736
            error_callback(r);                                                                                    // 737
        }                                                                                                         // 738
    };                                                                                                            // 739
    var post = function(msg, origin) {                                                                            // 740
        try {                                                                                                     // 741
            // When the iframe is not loaded, IE raises an exception                                              // 742
            // on 'contentWindow'.                                                                                // 743
            if (iframe && iframe.contentWindow) {                                                                 // 744
                iframe.contentWindow.postMessage(msg, origin);                                                    // 745
            }                                                                                                     // 746
        } catch (x) {};                                                                                           // 747
    };                                                                                                            // 748
                                                                                                                  // 749
    doc.open();                                                                                                   // 750
    doc.write('<html><s' + 'cript>' +                                                                             // 751
              'document.domain="' + document.domain + '";' +                                                      // 752
              '</s' + 'cript></html>');                                                                           // 753
    doc.close();                                                                                                  // 754
    doc.parentWindow[WPrefix] = _window[WPrefix];                                                                 // 755
    var c = doc.createElement('div');                                                                             // 756
    doc.body.appendChild(c);                                                                                      // 757
    iframe = doc.createElement('iframe');                                                                         // 758
    c.appendChild(iframe);                                                                                        // 759
    iframe.src = iframe_url;                                                                                      // 760
    tref = setTimeout(function(){onerror('timeout');}, 15000);                                                    // 761
    unload_ref = utils.unload_add(cleanup);                                                                       // 762
    return {                                                                                                      // 763
        post: post,                                                                                               // 764
        cleanup: cleanup,                                                                                         // 765
        loaded: unattach                                                                                          // 766
    };                                                                                                            // 767
};                                                                                                                // 768
//         [*] End of lib/dom.js                                                                                  // 769
                                                                                                                  // 770
                                                                                                                  // 771
//         [*] Including lib/dom2.js                                                                              // 772
/*                                                                                                                // 773
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 774
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 775
 *                                                                                                                // 776
 * For the license see COPYING.                                                                                   // 777
 * ***** END LICENSE BLOCK *****                                                                                  // 778
 */                                                                                                               // 779
                                                                                                                  // 780
var AbstractXHRObject = function(){};                                                                             // 781
AbstractXHRObject.prototype = new EventEmitter(['chunk', 'finish']);                                              // 782
                                                                                                                  // 783
AbstractXHRObject.prototype._start = function(method, url, payload, opts) {                                       // 784
    var that = this;                                                                                              // 785
                                                                                                                  // 786
    try {                                                                                                         // 787
        that.xhr = new XMLHttpRequest();                                                                          // 788
    } catch(x) {};                                                                                                // 789
                                                                                                                  // 790
    if (!that.xhr) {                                                                                              // 791
        try {                                                                                                     // 792
            that.xhr = new _window.ActiveXObject('Microsoft.XMLHTTP');                                            // 793
        } catch(x) {};                                                                                            // 794
    }                                                                                                             // 795
    if (_window.ActiveXObject || _window.XDomainRequest) {                                                        // 796
        // IE8 caches even POSTs                                                                                  // 797
        url += ((url.indexOf('?') === -1) ? '?' : '&') + 't='+(+new Date);                                        // 798
    }                                                                                                             // 799
                                                                                                                  // 800
    // Explorer tends to keep connection open, even after the                                                     // 801
    // tab gets closed: http://bugs.jquery.com/ticket/5280                                                        // 802
    that.unload_ref = utils.unload_add(function(){that._cleanup(true);});                                         // 803
    try {                                                                                                         // 804
        that.xhr.open(method, url, true);                                                                         // 805
    } catch(e) {                                                                                                  // 806
        // IE raises an exception on wrong port.                                                                  // 807
        that.emit('finish', 0, '');                                                                               // 808
        that._cleanup();                                                                                          // 809
        return;                                                                                                   // 810
    };                                                                                                            // 811
                                                                                                                  // 812
    if (!opts || !opts.no_credentials) {                                                                          // 813
        // Mozilla docs says https://developer.mozilla.org/en/XMLHttpRequest :                                    // 814
        // "This never affects same-site requests."                                                               // 815
        that.xhr.withCredentials = 'true';                                                                        // 816
    }                                                                                                             // 817
    if (opts && opts.headers) {                                                                                   // 818
        for(var key in opts.headers) {                                                                            // 819
            that.xhr.setRequestHeader(key, opts.headers[key]);                                                    // 820
        }                                                                                                         // 821
    }                                                                                                             // 822
                                                                                                                  // 823
    that.xhr.onreadystatechange = function() {                                                                    // 824
        if (that.xhr) {                                                                                           // 825
            var x = that.xhr;                                                                                     // 826
            switch (x.readyState) {                                                                               // 827
            case 3:                                                                                               // 828
                // IE doesn't like peeking into responseText or status                                            // 829
                // on Microsoft.XMLHTTP and readystate=3                                                          // 830
                try {                                                                                             // 831
                    var status = x.status;                                                                        // 832
                    var text = x.responseText;                                                                    // 833
                } catch (x) {};                                                                                   // 834
                // IE returns 1223 for 204: http://bugs.jquery.com/ticket/1450                                    // 835
                if (status === 1223) status = 204;                                                                // 836
                                                                                                                  // 837
                // IE does return readystate == 3 for 404 answers.                                                // 838
                if (text && text.length > 0) {                                                                    // 839
                    that.emit('chunk', status, text);                                                             // 840
                }                                                                                                 // 841
                break;                                                                                            // 842
            case 4:                                                                                               // 843
                var status = x.status;                                                                            // 844
                // IE returns 1223 for 204: http://bugs.jquery.com/ticket/1450                                    // 845
                if (status === 1223) status = 204;                                                                // 846
                                                                                                                  // 847
                that.emit('finish', status, x.responseText);                                                      // 848
                that._cleanup(false);                                                                             // 849
                break;                                                                                            // 850
            }                                                                                                     // 851
        }                                                                                                         // 852
    };                                                                                                            // 853
    that.xhr.send(payload);                                                                                       // 854
};                                                                                                                // 855
                                                                                                                  // 856
AbstractXHRObject.prototype._cleanup = function(abort) {                                                          // 857
    var that = this;                                                                                              // 858
    if (!that.xhr) return;                                                                                        // 859
    utils.unload_del(that.unload_ref);                                                                            // 860
                                                                                                                  // 861
    // IE needs this field to be a function                                                                       // 862
    that.xhr.onreadystatechange = function(){};                                                                   // 863
                                                                                                                  // 864
    if (abort) {                                                                                                  // 865
        try {                                                                                                     // 866
            that.xhr.abort();                                                                                     // 867
        } catch(x) {};                                                                                            // 868
    }                                                                                                             // 869
    that.unload_ref = that.xhr = null;                                                                            // 870
};                                                                                                                // 871
                                                                                                                  // 872
AbstractXHRObject.prototype.close = function() {                                                                  // 873
    var that = this;                                                                                              // 874
    that.nuke();                                                                                                  // 875
    that._cleanup(true);                                                                                          // 876
};                                                                                                                // 877
                                                                                                                  // 878
var XHRCorsObject = utils.XHRCorsObject = function() {                                                            // 879
    var that = this, args = arguments;                                                                            // 880
    utils.delay(function(){that._start.apply(that, args);});                                                      // 881
};                                                                                                                // 882
XHRCorsObject.prototype = new AbstractXHRObject();                                                                // 883
                                                                                                                  // 884
var XHRLocalObject = utils.XHRLocalObject = function(method, url, payload) {                                      // 885
    var that = this;                                                                                              // 886
    utils.delay(function(){                                                                                       // 887
        that._start(method, url, payload, {                                                                       // 888
            no_credentials: true                                                                                  // 889
        });                                                                                                       // 890
    });                                                                                                           // 891
};                                                                                                                // 892
XHRLocalObject.prototype = new AbstractXHRObject();                                                               // 893
                                                                                                                  // 894
                                                                                                                  // 895
                                                                                                                  // 896
// References:                                                                                                    // 897
//   http://ajaxian.com/archives/100-line-ajax-wrapper                                                            // 898
//   http://msdn.microsoft.com/en-us/library/cc288060(v=VS.85).aspx                                               // 899
var XDRObject = utils.XDRObject = function(method, url, payload) {                                                // 900
    var that = this;                                                                                              // 901
    utils.delay(function(){that._start(method, url, payload);});                                                  // 902
};                                                                                                                // 903
XDRObject.prototype = new EventEmitter(['chunk', 'finish']);                                                      // 904
XDRObject.prototype._start = function(method, url, payload) {                                                     // 905
    var that = this;                                                                                              // 906
    var xdr = new XDomainRequest();                                                                               // 907
    // IE caches even POSTs                                                                                       // 908
    url += ((url.indexOf('?') === -1) ? '?' : '&') + 't='+(+new Date);                                            // 909
                                                                                                                  // 910
    var onerror = xdr.ontimeout = xdr.onerror = function() {                                                      // 911
        that.emit('finish', 0, '');                                                                               // 912
        that._cleanup(false);                                                                                     // 913
    };                                                                                                            // 914
    xdr.onprogress = function() {                                                                                 // 915
        that.emit('chunk', 200, xdr.responseText);                                                                // 916
    };                                                                                                            // 917
    xdr.onload = function() {                                                                                     // 918
        that.emit('finish', 200, xdr.responseText);                                                               // 919
        that._cleanup(false);                                                                                     // 920
    };                                                                                                            // 921
    that.xdr = xdr;                                                                                               // 922
    that.unload_ref = utils.unload_add(function(){that._cleanup(true);});                                         // 923
    try {                                                                                                         // 924
        // Fails with AccessDenied if port number is bogus                                                        // 925
        that.xdr.open(method, url);                                                                               // 926
        that.xdr.send(payload);                                                                                   // 927
    } catch(x) {                                                                                                  // 928
        onerror();                                                                                                // 929
    }                                                                                                             // 930
};                                                                                                                // 931
                                                                                                                  // 932
XDRObject.prototype._cleanup = function(abort) {                                                                  // 933
    var that = this;                                                                                              // 934
    if (!that.xdr) return;                                                                                        // 935
    utils.unload_del(that.unload_ref);                                                                            // 936
                                                                                                                  // 937
    that.xdr.ontimeout = that.xdr.onerror = that.xdr.onprogress =                                                 // 938
        that.xdr.onload = null;                                                                                   // 939
    if (abort) {                                                                                                  // 940
        try {                                                                                                     // 941
            that.xdr.abort();                                                                                     // 942
        } catch(x) {};                                                                                            // 943
    }                                                                                                             // 944
    that.unload_ref = that.xdr = null;                                                                            // 945
};                                                                                                                // 946
                                                                                                                  // 947
XDRObject.prototype.close = function() {                                                                          // 948
    var that = this;                                                                                              // 949
    that.nuke();                                                                                                  // 950
    that._cleanup(true);                                                                                          // 951
};                                                                                                                // 952
                                                                                                                  // 953
// 1. Is natively via XHR                                                                                         // 954
// 2. Is natively via XDR                                                                                         // 955
// 3. Nope, but postMessage is there so it should work via the Iframe.                                            // 956
// 4. Nope, sorry.                                                                                                // 957
utils.isXHRCorsCapable = function() {                                                                             // 958
    if (_window.XMLHttpRequest && 'withCredentials' in new XMLHttpRequest()) {                                    // 959
        return 1;                                                                                                 // 960
    }                                                                                                             // 961
    // XDomainRequest doesn't work if page is served from file://                                                 // 962
    if (_window.XDomainRequest && _document.domain) {                                                             // 963
        return 2;                                                                                                 // 964
    }                                                                                                             // 965
    if (IframeTransport.enabled()) {                                                                              // 966
        return 3;                                                                                                 // 967
    }                                                                                                             // 968
    return 4;                                                                                                     // 969
};                                                                                                                // 970
//         [*] End of lib/dom2.js                                                                                 // 971
                                                                                                                  // 972
                                                                                                                  // 973
//         [*] Including lib/sockjs.js                                                                            // 974
/*                                                                                                                // 975
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 976
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 977
 *                                                                                                                // 978
 * For the license see COPYING.                                                                                   // 979
 * ***** END LICENSE BLOCK *****                                                                                  // 980
 */                                                                                                               // 981
                                                                                                                  // 982
var SockJS = function(url, dep_protocols_whitelist, options) {                                                    // 983
    if (!(this instanceof SockJS)) {                                                                              // 984
        // makes `new` optional                                                                                   // 985
        return new SockJS(url, dep_protocols_whitelist, options);                                                 // 986
    }                                                                                                             // 987
                                                                                                                  // 988
    var that = this, protocols_whitelist;                                                                         // 989
    that._options = {devel: false, debug: false, protocols_whitelist: [],                                         // 990
                     info: undefined, rtt: undefined};                                                            // 991
    if (options) {                                                                                                // 992
        utils.objectExtend(that._options, options);                                                               // 993
    }                                                                                                             // 994
    that._base_url = utils.amendUrl(url);                                                                         // 995
    that._server = that._options.server || utils.random_number_string(1000);                                      // 996
    if (that._options.protocols_whitelist &&                                                                      // 997
        that._options.protocols_whitelist.length) {                                                               // 998
        protocols_whitelist = that._options.protocols_whitelist;                                                  // 999
    } else {                                                                                                      // 1000
        // Deprecated API                                                                                         // 1001
        if (typeof dep_protocols_whitelist === 'string' &&                                                        // 1002
            dep_protocols_whitelist.length > 0) {                                                                 // 1003
            protocols_whitelist = [dep_protocols_whitelist];                                                      // 1004
        } else if (utils.isArray(dep_protocols_whitelist)) {                                                      // 1005
            protocols_whitelist = dep_protocols_whitelist                                                         // 1006
        } else {                                                                                                  // 1007
            protocols_whitelist = null;                                                                           // 1008
        }                                                                                                         // 1009
        if (protocols_whitelist) {                                                                                // 1010
            that._debug('Deprecated API: Use "protocols_whitelist" option ' +                                     // 1011
                        'instead of supplying protocol list as a second ' +                                       // 1012
                        'parameter to SockJS constructor.');                                                      // 1013
        }                                                                                                         // 1014
    }                                                                                                             // 1015
    that._protocols = [];                                                                                         // 1016
    that.protocol = null;                                                                                         // 1017
    that.readyState = SockJS.CONNECTING;                                                                          // 1018
    that._ir = createInfoReceiver(that._base_url);                                                                // 1019
    that._ir.onfinish = function(info, rtt) {                                                                     // 1020
        that._ir = null;                                                                                          // 1021
        if (info) {                                                                                               // 1022
            if (that._options.info) {                                                                             // 1023
                // Override if user supplies the option                                                           // 1024
                info = utils.objectExtend(info, that._options.info);                                              // 1025
            }                                                                                                     // 1026
            if (that._options.rtt) {                                                                              // 1027
                rtt = that._options.rtt;                                                                          // 1028
            }                                                                                                     // 1029
            that._applyInfo(info, rtt, protocols_whitelist);                                                      // 1030
            that._didClose();                                                                                     // 1031
        } else {                                                                                                  // 1032
            that._didClose(1002, 'Can\'t connect to server', true);                                               // 1033
        }                                                                                                         // 1034
    };                                                                                                            // 1035
};                                                                                                                // 1036
// Inheritance                                                                                                    // 1037
SockJS.prototype = new REventTarget();                                                                            // 1038
                                                                                                                  // 1039
SockJS.version = "0.3.4";                                                                                         // 1040
                                                                                                                  // 1041
SockJS.CONNECTING = 0;                                                                                            // 1042
SockJS.OPEN = 1;                                                                                                  // 1043
SockJS.CLOSING = 2;                                                                                               // 1044
SockJS.CLOSED = 3;                                                                                                // 1045
                                                                                                                  // 1046
SockJS.prototype._debug = function() {                                                                            // 1047
    if (this._options.debug)                                                                                      // 1048
        utils.log.apply(utils, arguments);                                                                        // 1049
};                                                                                                                // 1050
                                                                                                                  // 1051
SockJS.prototype._dispatchOpen = function() {                                                                     // 1052
    var that = this;                                                                                              // 1053
    if (that.readyState === SockJS.CONNECTING) {                                                                  // 1054
        if (that._transport_tref) {                                                                               // 1055
            clearTimeout(that._transport_tref);                                                                   // 1056
            that._transport_tref = null;                                                                          // 1057
        }                                                                                                         // 1058
        that.readyState = SockJS.OPEN;                                                                            // 1059
        that.dispatchEvent(new SimpleEvent("open"));                                                              // 1060
    } else {                                                                                                      // 1061
        // The server might have been restarted, and lost track of our                                            // 1062
        // connection.                                                                                            // 1063
        that._didClose(1006, "Server lost session");                                                              // 1064
    }                                                                                                             // 1065
};                                                                                                                // 1066
                                                                                                                  // 1067
SockJS.prototype._dispatchMessage = function(data) {                                                              // 1068
    var that = this;                                                                                              // 1069
    if (that.readyState !== SockJS.OPEN)                                                                          // 1070
            return;                                                                                               // 1071
    that.dispatchEvent(new SimpleEvent("message", {data: data}));                                                 // 1072
};                                                                                                                // 1073
                                                                                                                  // 1074
SockJS.prototype._dispatchHeartbeat = function(data) {                                                            // 1075
    var that = this;                                                                                              // 1076
    if (that.readyState !== SockJS.OPEN)                                                                          // 1077
        return;                                                                                                   // 1078
    that.dispatchEvent(new SimpleEvent('heartbeat', {}));                                                         // 1079
};                                                                                                                // 1080
                                                                                                                  // 1081
SockJS.prototype._didClose = function(code, reason, force) {                                                      // 1082
    var that = this;                                                                                              // 1083
    if (that.readyState !== SockJS.CONNECTING &&                                                                  // 1084
        that.readyState !== SockJS.OPEN &&                                                                        // 1085
        that.readyState !== SockJS.CLOSING)                                                                       // 1086
            throw new Error('INVALID_STATE_ERR');                                                                 // 1087
    if (that._ir) {                                                                                               // 1088
        that._ir.nuke();                                                                                          // 1089
        that._ir = null;                                                                                          // 1090
    }                                                                                                             // 1091
                                                                                                                  // 1092
    if (that._transport) {                                                                                        // 1093
        that._transport.doCleanup();                                                                              // 1094
        that._transport = null;                                                                                   // 1095
    }                                                                                                             // 1096
                                                                                                                  // 1097
    var close_event = new SimpleEvent("close", {                                                                  // 1098
        code: code,                                                                                               // 1099
        reason: reason,                                                                                           // 1100
        wasClean: utils.userSetCode(code)});                                                                      // 1101
                                                                                                                  // 1102
    if (!utils.userSetCode(code) &&                                                                               // 1103
        that.readyState === SockJS.CONNECTING && !force) {                                                        // 1104
        if (that._try_next_protocol(close_event)) {                                                               // 1105
            return;                                                                                               // 1106
        }                                                                                                         // 1107
        close_event = new SimpleEvent("close", {code: 2000,                                                       // 1108
                                                reason: "All transports failed",                                  // 1109
                                                wasClean: false,                                                  // 1110
                                                last_event: close_event});                                        // 1111
    }                                                                                                             // 1112
    that.readyState = SockJS.CLOSED;                                                                              // 1113
                                                                                                                  // 1114
    utils.delay(function() {                                                                                      // 1115
                   that.dispatchEvent(close_event);                                                               // 1116
                });                                                                                               // 1117
};                                                                                                                // 1118
                                                                                                                  // 1119
SockJS.prototype._didMessage = function(data) {                                                                   // 1120
    var that = this;                                                                                              // 1121
    var type = data.slice(0, 1);                                                                                  // 1122
    switch(type) {                                                                                                // 1123
    case 'o':                                                                                                     // 1124
        that._dispatchOpen();                                                                                     // 1125
        break;                                                                                                    // 1126
    case 'a':                                                                                                     // 1127
        var payload = JSON.parse(data.slice(1) || '[]');                                                          // 1128
        for(var i=0; i < payload.length; i++){                                                                    // 1129
            that._dispatchMessage(payload[i]);                                                                    // 1130
        }                                                                                                         // 1131
        break;                                                                                                    // 1132
    case 'm':                                                                                                     // 1133
        var payload = JSON.parse(data.slice(1) || 'null');                                                        // 1134
        that._dispatchMessage(payload);                                                                           // 1135
        break;                                                                                                    // 1136
    case 'c':                                                                                                     // 1137
        var payload = JSON.parse(data.slice(1) || '[]');                                                          // 1138
        that._didClose(payload[0], payload[1]);                                                                   // 1139
        break;                                                                                                    // 1140
    case 'h':                                                                                                     // 1141
        that._dispatchHeartbeat();                                                                                // 1142
        break;                                                                                                    // 1143
    }                                                                                                             // 1144
};                                                                                                                // 1145
                                                                                                                  // 1146
SockJS.prototype._try_next_protocol = function(close_event) {                                                     // 1147
    var that = this;                                                                                              // 1148
    if (that.protocol) {                                                                                          // 1149
        that._debug('Closed transport:', that.protocol, ''+close_event);                                          // 1150
        that.protocol = null;                                                                                     // 1151
    }                                                                                                             // 1152
    if (that._transport_tref) {                                                                                   // 1153
        clearTimeout(that._transport_tref);                                                                       // 1154
        that._transport_tref = null;                                                                              // 1155
    }                                                                                                             // 1156
                                                                                                                  // 1157
    while(1) {                                                                                                    // 1158
        var protocol = that.protocol = that._protocols.shift();                                                   // 1159
        if (!protocol) {                                                                                          // 1160
            return false;                                                                                         // 1161
        }                                                                                                         // 1162
        // Some protocols require access to `body`, what if were in                                               // 1163
        // the `head`?                                                                                            // 1164
        if (SockJS[protocol] &&                                                                                   // 1165
            SockJS[protocol].need_body === true &&                                                                // 1166
            (!_document.body ||                                                                                   // 1167
             (typeof _document.readyState !== 'undefined'                                                         // 1168
              && _document.readyState !== 'complete'))) {                                                         // 1169
            that._protocols.unshift(protocol);                                                                    // 1170
            that.protocol = 'waiting-for-load';                                                                   // 1171
            utils.attachEvent('load', function(){                                                                 // 1172
                that._try_next_protocol();                                                                        // 1173
            });                                                                                                   // 1174
            return true;                                                                                          // 1175
        }                                                                                                         // 1176
                                                                                                                  // 1177
        if (!SockJS[protocol] ||                                                                                  // 1178
              !SockJS[protocol].enabled(that._options)) {                                                         // 1179
            that._debug('Skipping transport:', protocol);                                                         // 1180
        } else {                                                                                                  // 1181
            var roundTrips = SockJS[protocol].roundTrips || 1;                                                    // 1182
            var to = ((that._options.rto || 0) * roundTrips) || 5000;                                             // 1183
            that._transport_tref = utils.delay(to, function() {                                                   // 1184
                if (that.readyState === SockJS.CONNECTING) {                                                      // 1185
                    // I can't understand how it is possible to run                                               // 1186
                    // this timer, when the state is CLOSED, but                                                  // 1187
                    // apparently in IE everythin is possible.                                                    // 1188
                    that._didClose(2007, "Transport timeouted");                                                  // 1189
                }                                                                                                 // 1190
            });                                                                                                   // 1191
                                                                                                                  // 1192
            var connid = utils.random_string(8);                                                                  // 1193
            var trans_url = that._base_url + '/' + that._server + '/' + connid;                                   // 1194
            that._debug('Opening transport:', protocol, ' url:'+trans_url,                                        // 1195
                        ' RTO:'+that._options.rto);                                                               // 1196
            that._transport = new SockJS[protocol](that, trans_url,                                               // 1197
                                                   that._base_url);                                               // 1198
            return true;                                                                                          // 1199
        }                                                                                                         // 1200
    }                                                                                                             // 1201
};                                                                                                                // 1202
                                                                                                                  // 1203
SockJS.prototype.close = function(code, reason) {                                                                 // 1204
    var that = this;                                                                                              // 1205
    if (code && !utils.userSetCode(code))                                                                         // 1206
        throw new Error("INVALID_ACCESS_ERR");                                                                    // 1207
    if(that.readyState !== SockJS.CONNECTING &&                                                                   // 1208
       that.readyState !== SockJS.OPEN) {                                                                         // 1209
        return false;                                                                                             // 1210
    }                                                                                                             // 1211
    that.readyState = SockJS.CLOSING;                                                                             // 1212
    that._didClose(code || 1000, reason || "Normal closure");                                                     // 1213
    return true;                                                                                                  // 1214
};                                                                                                                // 1215
                                                                                                                  // 1216
SockJS.prototype.send = function(data) {                                                                          // 1217
    var that = this;                                                                                              // 1218
    if (that.readyState === SockJS.CONNECTING)                                                                    // 1219
        throw new Error('INVALID_STATE_ERR');                                                                     // 1220
    if (that.readyState === SockJS.OPEN) {                                                                        // 1221
        that._transport.doSend(utils.quote('' + data));                                                           // 1222
    }                                                                                                             // 1223
    return true;                                                                                                  // 1224
};                                                                                                                // 1225
                                                                                                                  // 1226
SockJS.prototype._applyInfo = function(info, rtt, protocols_whitelist) {                                          // 1227
    var that = this;                                                                                              // 1228
    that._options.info = info;                                                                                    // 1229
    that._options.rtt = rtt;                                                                                      // 1230
    that._options.rto = utils.countRTO(rtt);                                                                      // 1231
    that._options.info.null_origin = !_document.domain;                                                           // 1232
    // Servers can override base_url, eg to provide a randomized domain name and                                  // 1233
    // avoid browser per-domain connection limits.                                                                // 1234
    if (info.base_url)                                                                                            // 1235
      // <METEOR>                                                                                                 // 1236
      that._base_url = utils.amendUrl(info.base_url, that._base_url);                                             // 1237
      // </METEOR>                                                                                                // 1238
    var probed = utils.probeProtocols();                                                                          // 1239
    that._protocols = utils.detectProtocols(probed, protocols_whitelist, info);                                   // 1240
// <METEOR>                                                                                                       // 1241
// https://github.com/sockjs/sockjs-client/issues/79                                                              // 1242
    // Hack to avoid XDR when using different protocols                                                           // 1243
    // We're on IE trying to do cross-protocol. jsonp only.                                                       // 1244
    if (!utils.isSameOriginScheme(that._base_url) &&                                                              // 1245
        2 === utils.isXHRCorsCapable()) {                                                                         // 1246
        that._protocols = ['jsonp-polling'];                                                                      // 1247
    }                                                                                                             // 1248
// </METEOR>                                                                                                      // 1249
};                                                                                                                // 1250
//         [*] End of lib/sockjs.js                                                                               // 1251
                                                                                                                  // 1252
                                                                                                                  // 1253
//         [*] Including lib/trans-websocket.js                                                                   // 1254
/*                                                                                                                // 1255
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1256
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1257
 *                                                                                                                // 1258
 * For the license see COPYING.                                                                                   // 1259
 * ***** END LICENSE BLOCK *****                                                                                  // 1260
 */                                                                                                               // 1261
                                                                                                                  // 1262
var WebSocketTransport = SockJS.websocket = function(ri, trans_url) {                                             // 1263
    var that = this;                                                                                              // 1264
    var url = trans_url + '/websocket';                                                                           // 1265
    if (url.slice(0, 5) === 'https') {                                                                            // 1266
        url = 'wss' + url.slice(5);                                                                               // 1267
    } else {                                                                                                      // 1268
        url = 'ws' + url.slice(4);                                                                                // 1269
    }                                                                                                             // 1270
    that.ri = ri;                                                                                                 // 1271
    that.url = url;                                                                                               // 1272
    var Constructor = _window.WebSocket || _window.MozWebSocket;                                                  // 1273
                                                                                                                  // 1274
    that.ws = new Constructor(that.url);                                                                          // 1275
    that.ws.onmessage = function(e) {                                                                             // 1276
        that.ri._didMessage(e.data);                                                                              // 1277
    };                                                                                                            // 1278
    // Firefox has an interesting bug. If a websocket connection is                                               // 1279
    // created after onunload, it stays alive even when user                                                      // 1280
    // navigates away from the page. In such situation let's lie -                                                // 1281
    // let's not open the ws connection at all. See:                                                              // 1282
    // https://github.com/sockjs/sockjs-client/issues/28                                                          // 1283
    // https://bugzilla.mozilla.org/show_bug.cgi?id=696085                                                        // 1284
    that.unload_ref = utils.unload_add(function(){that.ws.close()});                                              // 1285
    that.ws.onclose = function() {                                                                                // 1286
        that.ri._didMessage(utils.closeFrame(1006, "WebSocket connection broken"));                               // 1287
    };                                                                                                            // 1288
};                                                                                                                // 1289
                                                                                                                  // 1290
WebSocketTransport.prototype.doSend = function(data) {                                                            // 1291
    this.ws.send('[' + data + ']');                                                                               // 1292
};                                                                                                                // 1293
                                                                                                                  // 1294
WebSocketTransport.prototype.doCleanup = function() {                                                             // 1295
    var that = this;                                                                                              // 1296
    var ws = that.ws;                                                                                             // 1297
    if (ws) {                                                                                                     // 1298
        ws.onmessage = ws.onclose = null;                                                                         // 1299
        ws.close();                                                                                               // 1300
        utils.unload_del(that.unload_ref);                                                                        // 1301
        that.unload_ref = that.ri = that.ws = null;                                                               // 1302
    }                                                                                                             // 1303
};                                                                                                                // 1304
                                                                                                                  // 1305
WebSocketTransport.enabled = function() {                                                                         // 1306
    return !!(_window.WebSocket || _window.MozWebSocket);                                                         // 1307
};                                                                                                                // 1308
                                                                                                                  // 1309
// In theory, ws should require 1 round trip. But in chrome, this is                                              // 1310
// not very stable over SSL. Most likely a ws connection requires a                                               // 1311
// separate SSL connection, in which case 2 round trips are an                                                    // 1312
// absolute minumum.                                                                                              // 1313
WebSocketTransport.roundTrips = 2;                                                                                // 1314
//         [*] End of lib/trans-websocket.js                                                                      // 1315
                                                                                                                  // 1316
                                                                                                                  // 1317
//         [*] Including lib/trans-sender.js                                                                      // 1318
/*                                                                                                                // 1319
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1320
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1321
 *                                                                                                                // 1322
 * For the license see COPYING.                                                                                   // 1323
 * ***** END LICENSE BLOCK *****                                                                                  // 1324
 */                                                                                                               // 1325
                                                                                                                  // 1326
var BufferedSender = function() {};                                                                               // 1327
BufferedSender.prototype.send_constructor = function(sender) {                                                    // 1328
    var that = this;                                                                                              // 1329
    that.send_buffer = [];                                                                                        // 1330
    that.sender = sender;                                                                                         // 1331
};                                                                                                                // 1332
BufferedSender.prototype.doSend = function(message) {                                                             // 1333
    var that = this;                                                                                              // 1334
    that.send_buffer.push(message);                                                                               // 1335
    if (!that.send_stop) {                                                                                        // 1336
        that.send_schedule();                                                                                     // 1337
    }                                                                                                             // 1338
};                                                                                                                // 1339
                                                                                                                  // 1340
// For polling transports in a situation when in the message callback,                                            // 1341
// new message is being send. If the sending connection was started                                               // 1342
// before receiving one, it is possible to saturate the network and                                               // 1343
// timeout due to the lack of receiving socket. To avoid that we delay                                            // 1344
// sending messages by some small time, in order to let receiving                                                 // 1345
// connection be started beforehand. This is only a halfmeasure and                                               // 1346
// does not fix the big problem, but it does make the tests go more                                               // 1347
// stable on slow networks.                                                                                       // 1348
BufferedSender.prototype.send_schedule_wait = function() {                                                        // 1349
    var that = this;                                                                                              // 1350
    var tref;                                                                                                     // 1351
    that.send_stop = function() {                                                                                 // 1352
        that.send_stop = null;                                                                                    // 1353
        clearTimeout(tref);                                                                                       // 1354
    };                                                                                                            // 1355
    tref = utils.delay(25, function() {                                                                           // 1356
        that.send_stop = null;                                                                                    // 1357
        that.send_schedule();                                                                                     // 1358
    });                                                                                                           // 1359
};                                                                                                                // 1360
                                                                                                                  // 1361
BufferedSender.prototype.send_schedule = function() {                                                             // 1362
    var that = this;                                                                                              // 1363
    if (that.send_buffer.length > 0) {                                                                            // 1364
        var payload = '[' + that.send_buffer.join(',') + ']';                                                     // 1365
        that.send_stop = that.sender(that.trans_url, payload, function(success, abort_reason) {                   // 1366
            that.send_stop = null;                                                                                // 1367
            if (success === false) {                                                                              // 1368
                that.ri._didClose(1006, 'Sending error ' + abort_reason);                                         // 1369
            } else {                                                                                              // 1370
                that.send_schedule_wait();                                                                        // 1371
            }                                                                                                     // 1372
        });                                                                                                       // 1373
        that.send_buffer = [];                                                                                    // 1374
    }                                                                                                             // 1375
};                                                                                                                // 1376
                                                                                                                  // 1377
BufferedSender.prototype.send_destructor = function() {                                                           // 1378
    var that = this;                                                                                              // 1379
    if (that._send_stop) {                                                                                        // 1380
        that._send_stop();                                                                                        // 1381
    }                                                                                                             // 1382
    that._send_stop = null;                                                                                       // 1383
};                                                                                                                // 1384
                                                                                                                  // 1385
var jsonPGenericSender = function(url, payload, callback) {                                                       // 1386
    var that = this;                                                                                              // 1387
                                                                                                                  // 1388
    if (!('_send_form' in that)) {                                                                                // 1389
        var form = that._send_form = _document.createElement('form');                                             // 1390
        var area = that._send_area = _document.createElement('textarea');                                         // 1391
        area.name = 'd';                                                                                          // 1392
        form.style.display = 'none';                                                                              // 1393
        form.style.position = 'absolute';                                                                         // 1394
        form.method = 'POST';                                                                                     // 1395
        form.enctype = 'application/x-www-form-urlencoded';                                                       // 1396
        form.acceptCharset = "UTF-8";                                                                             // 1397
        form.appendChild(area);                                                                                   // 1398
        _document.body.appendChild(form);                                                                         // 1399
    }                                                                                                             // 1400
    var form = that._send_form;                                                                                   // 1401
    var area = that._send_area;                                                                                   // 1402
    var id = 'a' + utils.random_string(8);                                                                        // 1403
    form.target = id;                                                                                             // 1404
    form.action = url + '/jsonp_send?i=' + id;                                                                    // 1405
                                                                                                                  // 1406
    var iframe;                                                                                                   // 1407
    try {                                                                                                         // 1408
        // ie6 dynamic iframes with target="" support (thanks Chris Lambacher)                                    // 1409
        iframe = _document.createElement('<iframe name="'+ id +'">');                                             // 1410
    } catch(x) {                                                                                                  // 1411
        iframe = _document.createElement('iframe');                                                               // 1412
        iframe.name = id;                                                                                         // 1413
    }                                                                                                             // 1414
    iframe.id = id;                                                                                               // 1415
    form.appendChild(iframe);                                                                                     // 1416
    iframe.style.display = 'none';                                                                                // 1417
                                                                                                                  // 1418
    try {                                                                                                         // 1419
        area.value = payload;                                                                                     // 1420
    } catch(e) {                                                                                                  // 1421
        utils.log('Your browser is seriously broken. Go home! ' + e.message);                                     // 1422
    }                                                                                                             // 1423
    form.submit();                                                                                                // 1424
                                                                                                                  // 1425
    var completed = function(e) {                                                                                 // 1426
        if (!iframe.onerror) return;                                                                              // 1427
        iframe.onreadystatechange = iframe.onerror = iframe.onload = null;                                        // 1428
        // Opera mini doesn't like if we GC iframe                                                                // 1429
        // immediately, thus this timeout.                                                                        // 1430
        utils.delay(500, function() {                                                                             // 1431
                       iframe.parentNode.removeChild(iframe);                                                     // 1432
                       iframe = null;                                                                             // 1433
                   });                                                                                            // 1434
        area.value = '';                                                                                          // 1435
        // It is not possible to detect if the iframe succeeded or                                                // 1436
        // failed to submit our form.                                                                             // 1437
        callback(true);                                                                                           // 1438
    };                                                                                                            // 1439
    iframe.onerror = iframe.onload = completed;                                                                   // 1440
    iframe.onreadystatechange = function(e) {                                                                     // 1441
        if (iframe.readyState == 'complete') completed();                                                         // 1442
    };                                                                                                            // 1443
    return completed;                                                                                             // 1444
};                                                                                                                // 1445
                                                                                                                  // 1446
var createAjaxSender = function(AjaxObject) {                                                                     // 1447
    return function(url, payload, callback) {                                                                     // 1448
        var xo = new AjaxObject('POST', url + '/xhr_send', payload);                                              // 1449
        xo.onfinish = function(status, text) {                                                                    // 1450
            callback(status === 200 || status === 204,                                                            // 1451
                     'http status ' + status);                                                                    // 1452
        };                                                                                                        // 1453
        return function(abort_reason) {                                                                           // 1454
            callback(false, abort_reason);                                                                        // 1455
        };                                                                                                        // 1456
    };                                                                                                            // 1457
};                                                                                                                // 1458
//         [*] End of lib/trans-sender.js                                                                         // 1459
                                                                                                                  // 1460
                                                                                                                  // 1461
//         [*] Including lib/trans-jsonp-receiver.js                                                              // 1462
/*                                                                                                                // 1463
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1464
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1465
 *                                                                                                                // 1466
 * For the license see COPYING.                                                                                   // 1467
 * ***** END LICENSE BLOCK *****                                                                                  // 1468
 */                                                                                                               // 1469
                                                                                                                  // 1470
// Parts derived from Socket.io:                                                                                  // 1471
//    https://github.com/LearnBoost/socket.io/blob/0.6.17/lib/socket.io/transports/jsonp-polling.js               // 1472
// and jQuery-JSONP:                                                                                              // 1473
//    https://code.google.com/p/jquery-jsonp/source/browse/trunk/core/jquery.jsonp.js                             // 1474
var jsonPGenericReceiver = function(url, callback) {                                                              // 1475
    var tref;                                                                                                     // 1476
    var script = _document.createElement('script');                                                               // 1477
    var script2;  // Opera synchronous load trick.                                                                // 1478
    var close_script = function(frame) {                                                                          // 1479
        if (script2) {                                                                                            // 1480
            script2.parentNode.removeChild(script2);                                                              // 1481
            script2 = null;                                                                                       // 1482
        }                                                                                                         // 1483
        if (script) {                                                                                             // 1484
            clearTimeout(tref);                                                                                   // 1485
            // Unfortunately, you can't really abort script loading of                                            // 1486
            // the script.                                                                                        // 1487
            script.parentNode.removeChild(script);                                                                // 1488
            script.onreadystatechange = script.onerror =                                                          // 1489
                script.onload = script.onclick = null;                                                            // 1490
            script = null;                                                                                        // 1491
            callback(frame);                                                                                      // 1492
            callback = null;                                                                                      // 1493
        }                                                                                                         // 1494
    };                                                                                                            // 1495
                                                                                                                  // 1496
    // IE9 fires 'error' event after orsc or before, in random order.                                             // 1497
    var loaded_okay = false;                                                                                      // 1498
    var error_timer = null;                                                                                       // 1499
                                                                                                                  // 1500
    script.id = 'a' + utils.random_string(8);                                                                     // 1501
    script.src = url;                                                                                             // 1502
    script.type = 'text/javascript';                                                                              // 1503
    script.charset = 'UTF-8';                                                                                     // 1504
    script.onerror = function(e) {                                                                                // 1505
        if (!error_timer) {                                                                                       // 1506
            // Delay firing close_script.                                                                         // 1507
            error_timer = setTimeout(function() {                                                                 // 1508
                if (!loaded_okay) {                                                                               // 1509
                    close_script(utils.closeFrame(                                                                // 1510
                        1006,                                                                                     // 1511
                        "JSONP script loaded abnormally (onerror)"));                                             // 1512
                }                                                                                                 // 1513
            }, 1000);                                                                                             // 1514
        }                                                                                                         // 1515
    };                                                                                                            // 1516
    script.onload = function(e) {                                                                                 // 1517
        close_script(utils.closeFrame(1006, "JSONP script loaded abnormally (onload)"));                          // 1518
    };                                                                                                            // 1519
                                                                                                                  // 1520
    script.onreadystatechange = function(e) {                                                                     // 1521
        if (/loaded|closed/.test(script.readyState)) {                                                            // 1522
            if (script && script.htmlFor && script.onclick) {                                                     // 1523
                loaded_okay = true;                                                                               // 1524
                try {                                                                                             // 1525
                    // In IE, actually execute the script.                                                        // 1526
                    script.onclick();                                                                             // 1527
                } catch (x) {}                                                                                    // 1528
            }                                                                                                     // 1529
            if (script) {                                                                                         // 1530
                close_script(utils.closeFrame(1006, "JSONP script loaded abnormally (onreadystatechange)"));      // 1531
            }                                                                                                     // 1532
        }                                                                                                         // 1533
    };                                                                                                            // 1534
    // IE: event/htmlFor/onclick trick.                                                                           // 1535
    // One can't rely on proper order for onreadystatechange. In order to                                         // 1536
    // make sure, set a 'htmlFor' and 'event' properties, so that                                                 // 1537
    // script code will be installed as 'onclick' handler for the                                                 // 1538
    // script object. Later, onreadystatechange, manually execute this                                            // 1539
    // code. FF and Chrome doesn't work with 'event' and 'htmlFor'                                                // 1540
    // set. For reference see:                                                                                    // 1541
    //   http://jaubourg.net/2010/07/loading-script-as-onclick-handler-of.html                                    // 1542
    // Also, read on that about script ordering:                                                                  // 1543
    //   http://wiki.whatwg.org/wiki/Dynamic_Script_Execution_Order                                               // 1544
    if (typeof script.async === 'undefined' && _document.attachEvent) {                                           // 1545
        // According to mozilla docs, in recent browsers script.async defaults                                    // 1546
        // to 'true', so we may use it to detect a good browser:                                                  // 1547
        // https://developer.mozilla.org/en/HTML/Element/script                                                   // 1548
        if (!/opera/i.test(navigator.userAgent)) {                                                                // 1549
            // Naively assume we're in IE                                                                         // 1550
            try {                                                                                                 // 1551
                script.htmlFor = script.id;                                                                       // 1552
                script.event = "onclick";                                                                         // 1553
            } catch (x) {}                                                                                        // 1554
            script.async = true;                                                                                  // 1555
        } else {                                                                                                  // 1556
            // Opera, second sync script hack                                                                     // 1557
            script2 = _document.createElement('script');                                                          // 1558
            script2.text = "try{var a = document.getElementById('"+script.id+"'); if(a)a.onerror();}catch(x){};";
            script.async = script2.async = false;                                                                 // 1560
        }                                                                                                         // 1561
    }                                                                                                             // 1562
    if (typeof script.async !== 'undefined') {                                                                    // 1563
        script.async = true;                                                                                      // 1564
    }                                                                                                             // 1565
                                                                                                                  // 1566
    // Fallback mostly for Konqueror - stupid timer, 35 seconds shall be plenty.                                  // 1567
    tref = setTimeout(function() {                                                                                // 1568
                          close_script(utils.closeFrame(1006, "JSONP script loaded abnormally (timeout)"));       // 1569
                      }, 35000);                                                                                  // 1570
                                                                                                                  // 1571
    var head = _document.getElementsByTagName('head')[0];                                                         // 1572
    head.insertBefore(script, head.firstChild);                                                                   // 1573
    if (script2) {                                                                                                // 1574
        head.insertBefore(script2, head.firstChild);                                                              // 1575
    }                                                                                                             // 1576
    return close_script;                                                                                          // 1577
};                                                                                                                // 1578
//         [*] End of lib/trans-jsonp-receiver.js                                                                 // 1579
                                                                                                                  // 1580
                                                                                                                  // 1581
//         [*] Including lib/trans-jsonp-polling.js                                                               // 1582
/*                                                                                                                // 1583
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1584
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1585
 *                                                                                                                // 1586
 * For the license see COPYING.                                                                                   // 1587
 * ***** END LICENSE BLOCK *****                                                                                  // 1588
 */                                                                                                               // 1589
                                                                                                                  // 1590
// The simplest and most robust transport, using the well-know cross                                              // 1591
// domain hack - JSONP. This transport is quite inefficient - one                                                 // 1592
// mssage could use up to one http request. But at least it works almost                                          // 1593
// everywhere.                                                                                                    // 1594
// Known limitations:                                                                                             // 1595
//   o you will get a spinning cursor                                                                             // 1596
//   o for Konqueror a dumb timer is needed to detect errors                                                      // 1597
                                                                                                                  // 1598
                                                                                                                  // 1599
var JsonPTransport = SockJS['jsonp-polling'] = function(ri, trans_url) {                                          // 1600
    utils.polluteGlobalNamespace();                                                                               // 1601
    var that = this;                                                                                              // 1602
    that.ri = ri;                                                                                                 // 1603
    that.trans_url = trans_url;                                                                                   // 1604
    that.send_constructor(jsonPGenericSender);                                                                    // 1605
    that._schedule_recv();                                                                                        // 1606
};                                                                                                                // 1607
                                                                                                                  // 1608
// Inheritnace                                                                                                    // 1609
JsonPTransport.prototype = new BufferedSender();                                                                  // 1610
                                                                                                                  // 1611
JsonPTransport.prototype._schedule_recv = function() {                                                            // 1612
    var that = this;                                                                                              // 1613
    var callback = function(data) {                                                                               // 1614
        that._recv_stop = null;                                                                                   // 1615
        if (data) {                                                                                               // 1616
            // no data - heartbeat;                                                                               // 1617
            if (!that._is_closing) {                                                                              // 1618
                that.ri._didMessage(data);                                                                        // 1619
            }                                                                                                     // 1620
        }                                                                                                         // 1621
        // The message can be a close message, and change is_closing state.                                       // 1622
        if (!that._is_closing) {                                                                                  // 1623
            that._schedule_recv();                                                                                // 1624
        }                                                                                                         // 1625
    };                                                                                                            // 1626
    that._recv_stop = jsonPReceiverWrapper(that.trans_url + '/jsonp',                                             // 1627
                                           jsonPGenericReceiver, callback);                                       // 1628
};                                                                                                                // 1629
                                                                                                                  // 1630
JsonPTransport.enabled = function() {                                                                             // 1631
    return true;                                                                                                  // 1632
};                                                                                                                // 1633
                                                                                                                  // 1634
JsonPTransport.need_body = true;                                                                                  // 1635
                                                                                                                  // 1636
                                                                                                                  // 1637
JsonPTransport.prototype.doCleanup = function() {                                                                 // 1638
    var that = this;                                                                                              // 1639
    that._is_closing = true;                                                                                      // 1640
    if (that._recv_stop) {                                                                                        // 1641
        that._recv_stop();                                                                                        // 1642
    }                                                                                                             // 1643
    that.ri = that._recv_stop = null;                                                                             // 1644
    that.send_destructor();                                                                                       // 1645
};                                                                                                                // 1646
                                                                                                                  // 1647
                                                                                                                  // 1648
// Abstract away code that handles global namespace pollution.                                                    // 1649
var jsonPReceiverWrapper = function(url, constructReceiver, user_callback) {                                      // 1650
    var id = 'a' + utils.random_string(6);                                                                        // 1651
    var url_id = url + '?c=' + escape(WPrefix + '.' + id);                                                        // 1652
                                                                                                                  // 1653
    // Unfortunately it is not possible to abort loading of the                                                   // 1654
    // script. We need to keep track of frake close frames.                                                       // 1655
    var aborting = 0;                                                                                             // 1656
                                                                                                                  // 1657
    // Callback will be called exactly once.                                                                      // 1658
    var callback = function(frame) {                                                                              // 1659
        switch(aborting) {                                                                                        // 1660
        case 0:                                                                                                   // 1661
            // Normal behaviour - delete hook _and_ emit message.                                                 // 1662
            delete _window[WPrefix][id];                                                                          // 1663
            user_callback(frame);                                                                                 // 1664
            break;                                                                                                // 1665
        case 1:                                                                                                   // 1666
            // Fake close frame - emit but don't delete hook.                                                     // 1667
            user_callback(frame);                                                                                 // 1668
            aborting = 2;                                                                                         // 1669
            break;                                                                                                // 1670
        case 2:                                                                                                   // 1671
            // Got frame after connection was closed, delete hook, don't emit.                                    // 1672
            delete _window[WPrefix][id];                                                                          // 1673
            break;                                                                                                // 1674
        }                                                                                                         // 1675
    };                                                                                                            // 1676
                                                                                                                  // 1677
    var close_script = constructReceiver(url_id, callback);                                                       // 1678
    _window[WPrefix][id] = close_script;                                                                          // 1679
    var stop = function() {                                                                                       // 1680
        if (_window[WPrefix][id]) {                                                                               // 1681
            aborting = 1;                                                                                         // 1682
            _window[WPrefix][id](utils.closeFrame(1000, "JSONP user aborted read"));                              // 1683
        }                                                                                                         // 1684
    };                                                                                                            // 1685
    return stop;                                                                                                  // 1686
};                                                                                                                // 1687
//         [*] End of lib/trans-jsonp-polling.js                                                                  // 1688
                                                                                                                  // 1689
                                                                                                                  // 1690
//         [*] Including lib/trans-xhr.js                                                                         // 1691
/*                                                                                                                // 1692
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1693
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1694
 *                                                                                                                // 1695
 * For the license see COPYING.                                                                                   // 1696
 * ***** END LICENSE BLOCK *****                                                                                  // 1697
 */                                                                                                               // 1698
                                                                                                                  // 1699
var AjaxBasedTransport = function() {};                                                                           // 1700
AjaxBasedTransport.prototype = new BufferedSender();                                                              // 1701
                                                                                                                  // 1702
AjaxBasedTransport.prototype.run = function(ri, trans_url,                                                        // 1703
                                            url_suffix, Receiver, AjaxObject) {                                   // 1704
    var that = this;                                                                                              // 1705
    that.ri = ri;                                                                                                 // 1706
    that.trans_url = trans_url;                                                                                   // 1707
    that.send_constructor(createAjaxSender(AjaxObject));                                                          // 1708
    that.poll = new Polling(ri, Receiver,                                                                         // 1709
                            trans_url + url_suffix, AjaxObject);                                                  // 1710
};                                                                                                                // 1711
                                                                                                                  // 1712
AjaxBasedTransport.prototype.doCleanup = function() {                                                             // 1713
    var that = this;                                                                                              // 1714
    if (that.poll) {                                                                                              // 1715
        that.poll.abort();                                                                                        // 1716
        that.poll = null;                                                                                         // 1717
    }                                                                                                             // 1718
};                                                                                                                // 1719
                                                                                                                  // 1720
// xhr-streaming                                                                                                  // 1721
var XhrStreamingTransport = SockJS['xhr-streaming'] = function(ri, trans_url) {                                   // 1722
    this.run(ri, trans_url, '/xhr_streaming', XhrReceiver, utils.XHRCorsObject);                                  // 1723
};                                                                                                                // 1724
                                                                                                                  // 1725
XhrStreamingTransport.prototype = new AjaxBasedTransport();                                                       // 1726
                                                                                                                  // 1727
XhrStreamingTransport.enabled = function() {                                                                      // 1728
    // Support for CORS Ajax aka Ajax2? Opera 12 claims CORS but                                                  // 1729
    // doesn't do streaming.                                                                                      // 1730
    return (_window.XMLHttpRequest &&                                                                             // 1731
            'withCredentials' in new XMLHttpRequest() &&                                                          // 1732
            (!/opera/i.test(navigator.userAgent)));                                                               // 1733
};                                                                                                                // 1734
XhrStreamingTransport.roundTrips = 2; // preflight, ajax                                                          // 1735
                                                                                                                  // 1736
// Safari gets confused when a streaming ajax request is started                                                  // 1737
// before onload. This causes the load indicator to spin indefinetely.                                            // 1738
XhrStreamingTransport.need_body = true;                                                                           // 1739
                                                                                                                  // 1740
                                                                                                                  // 1741
// According to:                                                                                                  // 1742
//   http://stackoverflow.com/questions/1641507/detect-browser-support-for-cross-domain-xmlhttprequests           // 1743
//   http://hacks.mozilla.org/2009/07/cross-site-xmlhttprequest-with-cors/                                        // 1744
                                                                                                                  // 1745
                                                                                                                  // 1746
// xdr-streaming                                                                                                  // 1747
var XdrStreamingTransport = SockJS['xdr-streaming'] = function(ri, trans_url) {                                   // 1748
    this.run(ri, trans_url, '/xhr_streaming', XhrReceiver, utils.XDRObject);                                      // 1749
};                                                                                                                // 1750
                                                                                                                  // 1751
XdrStreamingTransport.prototype = new AjaxBasedTransport();                                                       // 1752
                                                                                                                  // 1753
XdrStreamingTransport.enabled = function() {                                                                      // 1754
    return !!_window.XDomainRequest;                                                                              // 1755
};                                                                                                                // 1756
XdrStreamingTransport.roundTrips = 2; // preflight, ajax                                                          // 1757
                                                                                                                  // 1758
                                                                                                                  // 1759
                                                                                                                  // 1760
// xhr-polling                                                                                                    // 1761
var XhrPollingTransport = SockJS['xhr-polling'] = function(ri, trans_url) {                                       // 1762
    this.run(ri, trans_url, '/xhr', XhrReceiver, utils.XHRCorsObject);                                            // 1763
};                                                                                                                // 1764
                                                                                                                  // 1765
XhrPollingTransport.prototype = new AjaxBasedTransport();                                                         // 1766
                                                                                                                  // 1767
XhrPollingTransport.enabled = XhrStreamingTransport.enabled;                                                      // 1768
XhrPollingTransport.roundTrips = 2; // preflight, ajax                                                            // 1769
                                                                                                                  // 1770
                                                                                                                  // 1771
// xdr-polling                                                                                                    // 1772
var XdrPollingTransport = SockJS['xdr-polling'] = function(ri, trans_url) {                                       // 1773
    this.run(ri, trans_url, '/xhr', XhrReceiver, utils.XDRObject);                                                // 1774
};                                                                                                                // 1775
                                                                                                                  // 1776
XdrPollingTransport.prototype = new AjaxBasedTransport();                                                         // 1777
                                                                                                                  // 1778
XdrPollingTransport.enabled = XdrStreamingTransport.enabled;                                                      // 1779
XdrPollingTransport.roundTrips = 2; // preflight, ajax                                                            // 1780
//         [*] End of lib/trans-xhr.js                                                                            // 1781
                                                                                                                  // 1782
                                                                                                                  // 1783
//         [*] Including lib/trans-iframe.js                                                                      // 1784
/*                                                                                                                // 1785
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1786
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1787
 *                                                                                                                // 1788
 * For the license see COPYING.                                                                                   // 1789
 * ***** END LICENSE BLOCK *****                                                                                  // 1790
 */                                                                                                               // 1791
                                                                                                                  // 1792
// Few cool transports do work only for same-origin. In order to make                                             // 1793
// them working cross-domain we shall use iframe, served form the                                                 // 1794
// remote domain. New browsers, have capabilities to communicate with                                             // 1795
// cross domain iframe, using postMessage(). In IE it was implemented                                             // 1796
// from IE 8+, but of course, IE got some details wrong:                                                          // 1797
//    http://msdn.microsoft.com/en-us/library/cc197015(v=VS.85).aspx                                              // 1798
//    http://stevesouders.com/misc/test-postmessage.php                                                           // 1799
                                                                                                                  // 1800
var IframeTransport = function() {};                                                                              // 1801
                                                                                                                  // 1802
IframeTransport.prototype.i_constructor = function(ri, trans_url, base_url) {                                     // 1803
    var that = this;                                                                                              // 1804
    that.ri = ri;                                                                                                 // 1805
    that.origin = utils.getOrigin(base_url);                                                                      // 1806
    that.base_url = base_url;                                                                                     // 1807
    that.trans_url = trans_url;                                                                                   // 1808
                                                                                                                  // 1809
    var iframe_url = base_url + '/iframe.html';                                                                   // 1810
    if (that.ri._options.devel) {                                                                                 // 1811
        iframe_url += '?t=' + (+new Date);                                                                        // 1812
    }                                                                                                             // 1813
    that.window_id = utils.random_string(8);                                                                      // 1814
    iframe_url += '#' + that.window_id;                                                                           // 1815
                                                                                                                  // 1816
    that.iframeObj = utils.createIframe(iframe_url, function(r) {                                                 // 1817
                                            that.ri._didClose(1006, "Unable to load an iframe (" + r + ")");      // 1818
                                        });                                                                       // 1819
                                                                                                                  // 1820
    that.onmessage_cb = utils.bind(that.onmessage, that);                                                         // 1821
    utils.attachMessage(that.onmessage_cb);                                                                       // 1822
};                                                                                                                // 1823
                                                                                                                  // 1824
IframeTransport.prototype.doCleanup = function() {                                                                // 1825
    var that = this;                                                                                              // 1826
    if (that.iframeObj) {                                                                                         // 1827
        utils.detachMessage(that.onmessage_cb);                                                                   // 1828
        try {                                                                                                     // 1829
            // When the iframe is not loaded, IE raises an exception                                              // 1830
            // on 'contentWindow'.                                                                                // 1831
            if (that.iframeObj.iframe.contentWindow) {                                                            // 1832
                that.postMessage('c');                                                                            // 1833
            }                                                                                                     // 1834
        } catch (x) {}                                                                                            // 1835
        that.iframeObj.cleanup();                                                                                 // 1836
        that.iframeObj = null;                                                                                    // 1837
        that.onmessage_cb = that.iframeObj = null;                                                                // 1838
    }                                                                                                             // 1839
};                                                                                                                // 1840
                                                                                                                  // 1841
IframeTransport.prototype.onmessage = function(e) {                                                               // 1842
    var that = this;                                                                                              // 1843
    if (e.origin !== that.origin) return;                                                                         // 1844
    var window_id = e.data.slice(0, 8);                                                                           // 1845
    var type = e.data.slice(8, 9);                                                                                // 1846
    var data = e.data.slice(9);                                                                                   // 1847
                                                                                                                  // 1848
    if (window_id !== that.window_id) return;                                                                     // 1849
                                                                                                                  // 1850
    switch(type) {                                                                                                // 1851
    case 's':                                                                                                     // 1852
        that.iframeObj.loaded();                                                                                  // 1853
        that.postMessage('s', JSON.stringify([SockJS.version, that.protocol, that.trans_url, that.base_url]));    // 1854
        break;                                                                                                    // 1855
    case 't':                                                                                                     // 1856
        that.ri._didMessage(data);                                                                                // 1857
        break;                                                                                                    // 1858
    }                                                                                                             // 1859
};                                                                                                                // 1860
                                                                                                                  // 1861
IframeTransport.prototype.postMessage = function(type, data) {                                                    // 1862
    var that = this;                                                                                              // 1863
    that.iframeObj.post(that.window_id + type + (data || ''), that.origin);                                       // 1864
};                                                                                                                // 1865
                                                                                                                  // 1866
IframeTransport.prototype.doSend = function (message) {                                                           // 1867
    this.postMessage('m', message);                                                                               // 1868
};                                                                                                                // 1869
                                                                                                                  // 1870
IframeTransport.enabled = function() {                                                                            // 1871
    // postMessage misbehaves in konqueror 4.6.5 - the messages are delivered with                                // 1872
    // huge delay, or not at all.                                                                                 // 1873
    var konqueror = navigator && navigator.userAgent && navigator.userAgent.indexOf('Konqueror') !== -1;          // 1874
    return ((typeof _window.postMessage === 'function' ||                                                         // 1875
            typeof _window.postMessage === 'object') && (!konqueror));                                            // 1876
};                                                                                                                // 1877
//         [*] End of lib/trans-iframe.js                                                                         // 1878
                                                                                                                  // 1879
                                                                                                                  // 1880
//         [*] Including lib/trans-iframe-within.js                                                               // 1881
/*                                                                                                                // 1882
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1883
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1884
 *                                                                                                                // 1885
 * For the license see COPYING.                                                                                   // 1886
 * ***** END LICENSE BLOCK *****                                                                                  // 1887
 */                                                                                                               // 1888
                                                                                                                  // 1889
var curr_window_id;                                                                                               // 1890
                                                                                                                  // 1891
var postMessage = function (type, data) {                                                                         // 1892
    if(parent !== _window) {                                                                                      // 1893
        parent.postMessage(curr_window_id + type + (data || ''), '*');                                            // 1894
    } else {                                                                                                      // 1895
        utils.log("Can't postMessage, no parent window.", type, data);                                            // 1896
    }                                                                                                             // 1897
};                                                                                                                // 1898
                                                                                                                  // 1899
var FacadeJS = function() {};                                                                                     // 1900
FacadeJS.prototype._didClose = function (code, reason) {                                                          // 1901
    postMessage('t', utils.closeFrame(code, reason));                                                             // 1902
};                                                                                                                // 1903
FacadeJS.prototype._didMessage = function (frame) {                                                               // 1904
    postMessage('t', frame);                                                                                      // 1905
};                                                                                                                // 1906
FacadeJS.prototype._doSend = function (data) {                                                                    // 1907
    this._transport.doSend(data);                                                                                 // 1908
};                                                                                                                // 1909
FacadeJS.prototype._doCleanup = function () {                                                                     // 1910
    this._transport.doCleanup();                                                                                  // 1911
};                                                                                                                // 1912
                                                                                                                  // 1913
utils.parent_origin = undefined;                                                                                  // 1914
                                                                                                                  // 1915
SockJS.bootstrap_iframe = function() {                                                                            // 1916
    var facade;                                                                                                   // 1917
    curr_window_id = _document.location.hash.slice(1);                                                            // 1918
    var onMessage = function(e) {                                                                                 // 1919
        if(e.source !== parent) return;                                                                           // 1920
        if(typeof utils.parent_origin === 'undefined')                                                            // 1921
            utils.parent_origin = e.origin;                                                                       // 1922
        if (e.origin !== utils.parent_origin) return;                                                             // 1923
                                                                                                                  // 1924
        var window_id = e.data.slice(0, 8);                                                                       // 1925
        var type = e.data.slice(8, 9);                                                                            // 1926
        var data = e.data.slice(9);                                                                               // 1927
        if (window_id !== curr_window_id) return;                                                                 // 1928
        switch(type) {                                                                                            // 1929
        case 's':                                                                                                 // 1930
            var p = JSON.parse(data);                                                                             // 1931
            var version = p[0];                                                                                   // 1932
            var protocol = p[1];                                                                                  // 1933
            var trans_url = p[2];                                                                                 // 1934
            var base_url = p[3];                                                                                  // 1935
            if (version !== SockJS.version) {                                                                     // 1936
                utils.log("Incompatibile SockJS! Main site uses:" +                                               // 1937
                          " \"" + version + "\", the iframe:" +                                                   // 1938
                          " \"" + SockJS.version + "\".");                                                        // 1939
            }                                                                                                     // 1940
            if (!utils.flatUrl(trans_url) || !utils.flatUrl(base_url)) {                                          // 1941
                utils.log("Only basic urls are supported in SockJS");                                             // 1942
                return;                                                                                           // 1943
            }                                                                                                     // 1944
                                                                                                                  // 1945
            if (!utils.isSameOriginUrl(trans_url) ||                                                              // 1946
                !utils.isSameOriginUrl(base_url)) {                                                               // 1947
                utils.log("Can't connect to different domain from within an " +                                   // 1948
                          "iframe. (" + JSON.stringify([_window.location.href, trans_url, base_url]) +            // 1949
                          ")");                                                                                   // 1950
                return;                                                                                           // 1951
            }                                                                                                     // 1952
            facade = new FacadeJS();                                                                              // 1953
            facade._transport = new FacadeJS[protocol](facade, trans_url, base_url);                              // 1954
            break;                                                                                                // 1955
        case 'm':                                                                                                 // 1956
            facade._doSend(data);                                                                                 // 1957
            break;                                                                                                // 1958
        case 'c':                                                                                                 // 1959
            if (facade)                                                                                           // 1960
                facade._doCleanup();                                                                              // 1961
            facade = null;                                                                                        // 1962
            break;                                                                                                // 1963
        }                                                                                                         // 1964
    };                                                                                                            // 1965
                                                                                                                  // 1966
    // alert('test ticker');                                                                                      // 1967
    // facade = new FacadeJS();                                                                                   // 1968
    // facade._transport = new FacadeJS['w-iframe-xhr-polling'](facade, 'http://host.com:9999/ticker/12/basd');   // 1969
                                                                                                                  // 1970
    utils.attachMessage(onMessage);                                                                               // 1971
                                                                                                                  // 1972
    // Start                                                                                                      // 1973
    postMessage('s');                                                                                             // 1974
};                                                                                                                // 1975
//         [*] End of lib/trans-iframe-within.js                                                                  // 1976
                                                                                                                  // 1977
                                                                                                                  // 1978
//         [*] Including lib/info.js                                                                              // 1979
/*                                                                                                                // 1980
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 1981
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 1982
 *                                                                                                                // 1983
 * For the license see COPYING.                                                                                   // 1984
 * ***** END LICENSE BLOCK *****                                                                                  // 1985
 */                                                                                                               // 1986
                                                                                                                  // 1987
var InfoReceiver = function(base_url, AjaxObject) {                                                               // 1988
    var that = this;                                                                                              // 1989
    utils.delay(function(){that.doXhr(base_url, AjaxObject);});                                                   // 1990
};                                                                                                                // 1991
                                                                                                                  // 1992
InfoReceiver.prototype = new EventEmitter(['finish']);                                                            // 1993
                                                                                                                  // 1994
InfoReceiver.prototype.doXhr = function(base_url, AjaxObject) {                                                   // 1995
    var that = this;                                                                                              // 1996
    var t0 = (new Date()).getTime();                                                                              // 1997
                                                                                                                  // 1998
// <METEOR>                                                                                                       // 1999
  // https://github.com/sockjs/sockjs-client/pull/129                                                             // 2000
  // var xo = new AjaxObject('GET', base_url + '/info');                                                          // 2001
                                                                                                                  // 2002
    var xo = new AjaxObject(                                                                                      // 2003
      // add cachebusting parameter to url to work around a chrome bug:                                           // 2004
      // https://code.google.com/p/chromium/issues/detail?id=263981                                               // 2005
      // or misbehaving proxies.                                                                                  // 2006
      'GET', base_url + '/info?cb=' + utils.random_string(10))                                                    // 2007
// </METEOR>                                                                                                      // 2008
                                                                                                                  // 2009
    var tref = utils.delay(8000,                                                                                  // 2010
                           function(){xo.ontimeout();});                                                          // 2011
                                                                                                                  // 2012
    xo.onfinish = function(status, text) {                                                                        // 2013
        clearTimeout(tref);                                                                                       // 2014
        tref = null;                                                                                              // 2015
        if (status === 200) {                                                                                     // 2016
            var rtt = (new Date()).getTime() - t0;                                                                // 2017
            var info = JSON.parse(text);                                                                          // 2018
            if (typeof info !== 'object') info = {};                                                              // 2019
            that.emit('finish', info, rtt);                                                                       // 2020
        } else {                                                                                                  // 2021
            that.emit('finish');                                                                                  // 2022
        }                                                                                                         // 2023
    };                                                                                                            // 2024
    xo.ontimeout = function() {                                                                                   // 2025
        xo.close();                                                                                               // 2026
        that.emit('finish');                                                                                      // 2027
    };                                                                                                            // 2028
};                                                                                                                // 2029
                                                                                                                  // 2030
var InfoReceiverIframe = function(base_url) {                                                                     // 2031
    var that = this;                                                                                              // 2032
    var go = function() {                                                                                         // 2033
        var ifr = new IframeTransport();                                                                          // 2034
        ifr.protocol = 'w-iframe-info-receiver';                                                                  // 2035
        var fun = function(r) {                                                                                   // 2036
            if (typeof r === 'string' && r.substr(0,1) === 'm') {                                                 // 2037
                var d = JSON.parse(r.substr(1));                                                                  // 2038
                var info = d[0], rtt = d[1];                                                                      // 2039
                that.emit('finish', info, rtt);                                                                   // 2040
            } else {                                                                                              // 2041
                that.emit('finish');                                                                              // 2042
            }                                                                                                     // 2043
            ifr.doCleanup();                                                                                      // 2044
            ifr = null;                                                                                           // 2045
        };                                                                                                        // 2046
        var mock_ri = {                                                                                           // 2047
            _options: {},                                                                                         // 2048
            _didClose: fun,                                                                                       // 2049
            _didMessage: fun                                                                                      // 2050
        };                                                                                                        // 2051
        ifr.i_constructor(mock_ri, base_url, base_url);                                                           // 2052
    }                                                                                                             // 2053
    if(!_document.body) {                                                                                         // 2054
        utils.attachEvent('load', go);                                                                            // 2055
    } else {                                                                                                      // 2056
        go();                                                                                                     // 2057
    }                                                                                                             // 2058
};                                                                                                                // 2059
InfoReceiverIframe.prototype = new EventEmitter(['finish']);                                                      // 2060
                                                                                                                  // 2061
                                                                                                                  // 2062
var InfoReceiverFake = function() {                                                                               // 2063
    // It may not be possible to do cross domain AJAX to get the info                                             // 2064
    // data, for example for IE7. But we want to run JSONP, so let's                                              // 2065
    // fake the response, with rtt=2s (rto=6s).                                                                   // 2066
    var that = this;                                                                                              // 2067
    utils.delay(function() {                                                                                      // 2068
        that.emit('finish', {}, 2000);                                                                            // 2069
    });                                                                                                           // 2070
};                                                                                                                // 2071
InfoReceiverFake.prototype = new EventEmitter(['finish']);                                                        // 2072
                                                                                                                  // 2073
var createInfoReceiver = function(base_url) {                                                                     // 2074
    if (utils.isSameOriginUrl(base_url)) {                                                                        // 2075
        // If, for some reason, we have SockJS locally - there's no                                               // 2076
        // need to start up the complex machinery. Just use ajax.                                                 // 2077
        return new InfoReceiver(base_url, utils.XHRLocalObject);                                                  // 2078
    }                                                                                                             // 2079
    switch (utils.isXHRCorsCapable()) {                                                                           // 2080
    case 1:                                                                                                       // 2081
        // XHRLocalObject -> no_credentials=true                                                                  // 2082
        return new InfoReceiver(base_url, utils.XHRLocalObject);                                                  // 2083
    case 2:                                                                                                       // 2084
// <METEOR>                                                                                                       // 2085
// https://github.com/sockjs/sockjs-client/issues/79                                                              // 2086
        // XDR doesn't work across different schemes                                                              // 2087
        // http://blogs.msdn.com/b/ieinternals/archive/2010/05/13/xdomainrequest-restrictions-limitations-and-workarounds.aspx
        if (utils.isSameOriginScheme(base_url))                                                                   // 2089
            return new InfoReceiver(base_url, utils.XDRObject);                                                   // 2090
        else                                                                                                      // 2091
            return new InfoReceiverFake();                                                                        // 2092
// </METEOR>                                                                                                      // 2093
    case 3:                                                                                                       // 2094
        // Opera                                                                                                  // 2095
        return new InfoReceiverIframe(base_url);                                                                  // 2096
    default:                                                                                                      // 2097
        // IE 7                                                                                                   // 2098
        return new InfoReceiverFake();                                                                            // 2099
    };                                                                                                            // 2100
};                                                                                                                // 2101
                                                                                                                  // 2102
                                                                                                                  // 2103
var WInfoReceiverIframe = FacadeJS['w-iframe-info-receiver'] = function(ri, _trans_url, base_url) {               // 2104
    var ir = new InfoReceiver(base_url, utils.XHRLocalObject);                                                    // 2105
    ir.onfinish = function(info, rtt) {                                                                           // 2106
        ri._didMessage('m'+JSON.stringify([info, rtt]));                                                          // 2107
        ri._didClose();                                                                                           // 2108
    }                                                                                                             // 2109
};                                                                                                                // 2110
WInfoReceiverIframe.prototype.doCleanup = function() {};                                                          // 2111
//         [*] End of lib/info.js                                                                                 // 2112
                                                                                                                  // 2113
                                                                                                                  // 2114
//         [*] Including lib/trans-iframe-eventsource.js                                                          // 2115
/*                                                                                                                // 2116
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2117
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2118
 *                                                                                                                // 2119
 * For the license see COPYING.                                                                                   // 2120
 * ***** END LICENSE BLOCK *****                                                                                  // 2121
 */                                                                                                               // 2122
                                                                                                                  // 2123
var EventSourceIframeTransport = SockJS['iframe-eventsource'] = function () {                                     // 2124
    var that = this;                                                                                              // 2125
    that.protocol = 'w-iframe-eventsource';                                                                       // 2126
    that.i_constructor.apply(that, arguments);                                                                    // 2127
};                                                                                                                // 2128
                                                                                                                  // 2129
EventSourceIframeTransport.prototype = new IframeTransport();                                                     // 2130
                                                                                                                  // 2131
EventSourceIframeTransport.enabled = function () {                                                                // 2132
    return ('EventSource' in _window) && IframeTransport.enabled();                                               // 2133
};                                                                                                                // 2134
                                                                                                                  // 2135
EventSourceIframeTransport.need_body = true;                                                                      // 2136
EventSourceIframeTransport.roundTrips = 3; // html, javascript, eventsource                                       // 2137
                                                                                                                  // 2138
                                                                                                                  // 2139
// w-iframe-eventsource                                                                                           // 2140
var EventSourceTransport = FacadeJS['w-iframe-eventsource'] = function(ri, trans_url) {                           // 2141
    this.run(ri, trans_url, '/eventsource', EventSourceReceiver, utils.XHRLocalObject);                           // 2142
}                                                                                                                 // 2143
EventSourceTransport.prototype = new AjaxBasedTransport();                                                        // 2144
//         [*] End of lib/trans-iframe-eventsource.js                                                             // 2145
                                                                                                                  // 2146
                                                                                                                  // 2147
//         [*] Including lib/trans-iframe-xhr-polling.js                                                          // 2148
/*                                                                                                                // 2149
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2150
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2151
 *                                                                                                                // 2152
 * For the license see COPYING.                                                                                   // 2153
 * ***** END LICENSE BLOCK *****                                                                                  // 2154
 */                                                                                                               // 2155
                                                                                                                  // 2156
var XhrPollingIframeTransport = SockJS['iframe-xhr-polling'] = function () {                                      // 2157
    var that = this;                                                                                              // 2158
    that.protocol = 'w-iframe-xhr-polling';                                                                       // 2159
    that.i_constructor.apply(that, arguments);                                                                    // 2160
};                                                                                                                // 2161
                                                                                                                  // 2162
XhrPollingIframeTransport.prototype = new IframeTransport();                                                      // 2163
                                                                                                                  // 2164
XhrPollingIframeTransport.enabled = function () {                                                                 // 2165
    return _window.XMLHttpRequest && IframeTransport.enabled();                                                   // 2166
};                                                                                                                // 2167
                                                                                                                  // 2168
XhrPollingIframeTransport.need_body = true;                                                                       // 2169
XhrPollingIframeTransport.roundTrips = 3; // html, javascript, xhr                                                // 2170
                                                                                                                  // 2171
                                                                                                                  // 2172
// w-iframe-xhr-polling                                                                                           // 2173
var XhrPollingITransport = FacadeJS['w-iframe-xhr-polling'] = function(ri, trans_url) {                           // 2174
    this.run(ri, trans_url, '/xhr', XhrReceiver, utils.XHRLocalObject);                                           // 2175
};                                                                                                                // 2176
                                                                                                                  // 2177
XhrPollingITransport.prototype = new AjaxBasedTransport();                                                        // 2178
//         [*] End of lib/trans-iframe-xhr-polling.js                                                             // 2179
                                                                                                                  // 2180
                                                                                                                  // 2181
//         [*] Including lib/trans-iframe-htmlfile.js                                                             // 2182
/*                                                                                                                // 2183
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2184
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2185
 *                                                                                                                // 2186
 * For the license see COPYING.                                                                                   // 2187
 * ***** END LICENSE BLOCK *****                                                                                  // 2188
 */                                                                                                               // 2189
                                                                                                                  // 2190
// This transport generally works in any browser, but will cause a                                                // 2191
// spinning cursor to appear in any browser other than IE.                                                        // 2192
// We may test this transport in all browsers - why not, but in                                                   // 2193
// production it should be only run in IE.                                                                        // 2194
                                                                                                                  // 2195
var HtmlFileIframeTransport = SockJS['iframe-htmlfile'] = function () {                                           // 2196
    var that = this;                                                                                              // 2197
    that.protocol = 'w-iframe-htmlfile';                                                                          // 2198
    that.i_constructor.apply(that, arguments);                                                                    // 2199
};                                                                                                                // 2200
                                                                                                                  // 2201
// Inheritance.                                                                                                   // 2202
HtmlFileIframeTransport.prototype = new IframeTransport();                                                        // 2203
                                                                                                                  // 2204
HtmlFileIframeTransport.enabled = function() {                                                                    // 2205
    return IframeTransport.enabled();                                                                             // 2206
};                                                                                                                // 2207
                                                                                                                  // 2208
HtmlFileIframeTransport.need_body = true;                                                                         // 2209
HtmlFileIframeTransport.roundTrips = 3; // html, javascript, htmlfile                                             // 2210
                                                                                                                  // 2211
                                                                                                                  // 2212
// w-iframe-htmlfile                                                                                              // 2213
var HtmlFileTransport = FacadeJS['w-iframe-htmlfile'] = function(ri, trans_url) {                                 // 2214
    this.run(ri, trans_url, '/htmlfile', HtmlfileReceiver, utils.XHRLocalObject);                                 // 2215
};                                                                                                                // 2216
HtmlFileTransport.prototype = new AjaxBasedTransport();                                                           // 2217
//         [*] End of lib/trans-iframe-htmlfile.js                                                                // 2218
                                                                                                                  // 2219
                                                                                                                  // 2220
//         [*] Including lib/trans-polling.js                                                                     // 2221
/*                                                                                                                // 2222
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2223
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2224
 *                                                                                                                // 2225
 * For the license see COPYING.                                                                                   // 2226
 * ***** END LICENSE BLOCK *****                                                                                  // 2227
 */                                                                                                               // 2228
                                                                                                                  // 2229
var Polling = function(ri, Receiver, recv_url, AjaxObject) {                                                      // 2230
    var that = this;                                                                                              // 2231
    that.ri = ri;                                                                                                 // 2232
    that.Receiver = Receiver;                                                                                     // 2233
    that.recv_url = recv_url;                                                                                     // 2234
    that.AjaxObject = AjaxObject;                                                                                 // 2235
    that._scheduleRecv();                                                                                         // 2236
};                                                                                                                // 2237
                                                                                                                  // 2238
Polling.prototype._scheduleRecv = function() {                                                                    // 2239
    var that = this;                                                                                              // 2240
    var poll = that.poll = new that.Receiver(that.recv_url, that.AjaxObject);                                     // 2241
    var msg_counter = 0;                                                                                          // 2242
    poll.onmessage = function(e) {                                                                                // 2243
        msg_counter += 1;                                                                                         // 2244
        that.ri._didMessage(e.data);                                                                              // 2245
    };                                                                                                            // 2246
    poll.onclose = function(e) {                                                                                  // 2247
        that.poll = poll = poll.onmessage = poll.onclose = null;                                                  // 2248
        if (!that.poll_is_closing) {                                                                              // 2249
            if (e.reason === 'permanent') {                                                                       // 2250
                that.ri._didClose(1006, 'Polling error (' + e.reason + ')');                                      // 2251
            } else {                                                                                              // 2252
                that._scheduleRecv();                                                                             // 2253
            }                                                                                                     // 2254
        }                                                                                                         // 2255
    };                                                                                                            // 2256
};                                                                                                                // 2257
                                                                                                                  // 2258
Polling.prototype.abort = function() {                                                                            // 2259
    var that = this;                                                                                              // 2260
    that.poll_is_closing = true;                                                                                  // 2261
    if (that.poll) {                                                                                              // 2262
        that.poll.abort();                                                                                        // 2263
    }                                                                                                             // 2264
};                                                                                                                // 2265
//         [*] End of lib/trans-polling.js                                                                        // 2266
                                                                                                                  // 2267
                                                                                                                  // 2268
//         [*] Including lib/trans-receiver-eventsource.js                                                        // 2269
/*                                                                                                                // 2270
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2271
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2272
 *                                                                                                                // 2273
 * For the license see COPYING.                                                                                   // 2274
 * ***** END LICENSE BLOCK *****                                                                                  // 2275
 */                                                                                                               // 2276
                                                                                                                  // 2277
var EventSourceReceiver = function(url) {                                                                         // 2278
    var that = this;                                                                                              // 2279
    var es = new EventSource(url);                                                                                // 2280
    es.onmessage = function(e) {                                                                                  // 2281
        that.dispatchEvent(new SimpleEvent('message',                                                             // 2282
                                           {'data': unescape(e.data)}));                                          // 2283
    };                                                                                                            // 2284
    that.es_close = es.onerror = function(e, abort_reason) {                                                      // 2285
        // ES on reconnection has readyState = 0 or 1.                                                            // 2286
        // on network error it's CLOSED = 2                                                                       // 2287
        var reason = abort_reason ? 'user' :                                                                      // 2288
            (es.readyState !== 2 ? 'network' : 'permanent');                                                      // 2289
        that.es_close = es.onmessage = es.onerror = null;                                                         // 2290
        // EventSource reconnects automatically.                                                                  // 2291
        es.close();                                                                                               // 2292
        es = null;                                                                                                // 2293
        // Safari and chrome < 15 crash if we close window before                                                 // 2294
        // waiting for ES cleanup. See:                                                                           // 2295
        //   https://code.google.com/p/chromium/issues/detail?id=89155                                            // 2296
        utils.delay(200, function() {                                                                             // 2297
                        that.dispatchEvent(new SimpleEvent('close', {reason: reason}));                           // 2298
                    });                                                                                           // 2299
    };                                                                                                            // 2300
};                                                                                                                // 2301
                                                                                                                  // 2302
EventSourceReceiver.prototype = new REventTarget();                                                               // 2303
                                                                                                                  // 2304
EventSourceReceiver.prototype.abort = function() {                                                                // 2305
    var that = this;                                                                                              // 2306
    if (that.es_close) {                                                                                          // 2307
        that.es_close({}, true);                                                                                  // 2308
    }                                                                                                             // 2309
};                                                                                                                // 2310
//         [*] End of lib/trans-receiver-eventsource.js                                                           // 2311
                                                                                                                  // 2312
                                                                                                                  // 2313
//         [*] Including lib/trans-receiver-htmlfile.js                                                           // 2314
/*                                                                                                                // 2315
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2316
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2317
 *                                                                                                                // 2318
 * For the license see COPYING.                                                                                   // 2319
 * ***** END LICENSE BLOCK *****                                                                                  // 2320
 */                                                                                                               // 2321
                                                                                                                  // 2322
var _is_ie_htmlfile_capable;                                                                                      // 2323
var isIeHtmlfileCapable = function() {                                                                            // 2324
    if (_is_ie_htmlfile_capable === undefined) {                                                                  // 2325
        if ('ActiveXObject' in _window) {                                                                         // 2326
            try {                                                                                                 // 2327
                _is_ie_htmlfile_capable = !!new ActiveXObject('htmlfile');                                        // 2328
            } catch (x) {}                                                                                        // 2329
        } else {                                                                                                  // 2330
            _is_ie_htmlfile_capable = false;                                                                      // 2331
        }                                                                                                         // 2332
    }                                                                                                             // 2333
    return _is_ie_htmlfile_capable;                                                                               // 2334
};                                                                                                                // 2335
                                                                                                                  // 2336
                                                                                                                  // 2337
var HtmlfileReceiver = function(url) {                                                                            // 2338
    var that = this;                                                                                              // 2339
    utils.polluteGlobalNamespace();                                                                               // 2340
                                                                                                                  // 2341
    that.id = 'a' + utils.random_string(6, 26);                                                                   // 2342
    url += ((url.indexOf('?') === -1) ? '?' : '&') +                                                              // 2343
        'c=' + escape(WPrefix + '.' + that.id);                                                                   // 2344
                                                                                                                  // 2345
    var constructor = isIeHtmlfileCapable() ?                                                                     // 2346
        utils.createHtmlfile : utils.createIframe;                                                                // 2347
                                                                                                                  // 2348
    var iframeObj;                                                                                                // 2349
    _window[WPrefix][that.id] = {                                                                                 // 2350
        start: function () {                                                                                      // 2351
            iframeObj.loaded();                                                                                   // 2352
        },                                                                                                        // 2353
        message: function (data) {                                                                                // 2354
            that.dispatchEvent(new SimpleEvent('message', {'data': data}));                                       // 2355
        },                                                                                                        // 2356
        stop: function () {                                                                                       // 2357
            that.iframe_close({}, 'network');                                                                     // 2358
        }                                                                                                         // 2359
    };                                                                                                            // 2360
    that.iframe_close = function(e, abort_reason) {                                                               // 2361
        iframeObj.cleanup();                                                                                      // 2362
        that.iframe_close = iframeObj = null;                                                                     // 2363
        delete _window[WPrefix][that.id];                                                                         // 2364
        that.dispatchEvent(new SimpleEvent('close', {reason: abort_reason}));                                     // 2365
    };                                                                                                            // 2366
    iframeObj = constructor(url, function(e) {                                                                    // 2367
                                that.iframe_close({}, 'permanent');                                               // 2368
                            });                                                                                   // 2369
};                                                                                                                // 2370
                                                                                                                  // 2371
HtmlfileReceiver.prototype = new REventTarget();                                                                  // 2372
                                                                                                                  // 2373
HtmlfileReceiver.prototype.abort = function() {                                                                   // 2374
    var that = this;                                                                                              // 2375
    if (that.iframe_close) {                                                                                      // 2376
        that.iframe_close({}, 'user');                                                                            // 2377
    }                                                                                                             // 2378
};                                                                                                                // 2379
//         [*] End of lib/trans-receiver-htmlfile.js                                                              // 2380
                                                                                                                  // 2381
                                                                                                                  // 2382
//         [*] Including lib/trans-receiver-xhr.js                                                                // 2383
/*                                                                                                                // 2384
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2385
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2386
 *                                                                                                                // 2387
 * For the license see COPYING.                                                                                   // 2388
 * ***** END LICENSE BLOCK *****                                                                                  // 2389
 */                                                                                                               // 2390
                                                                                                                  // 2391
var XhrReceiver = function(url, AjaxObject) {                                                                     // 2392
    var that = this;                                                                                              // 2393
    var buf_pos = 0;                                                                                              // 2394
                                                                                                                  // 2395
    that.xo = new AjaxObject('POST', url, null);                                                                  // 2396
    that.xo.onchunk = function(status, text) {                                                                    // 2397
        if (status !== 200) return;                                                                               // 2398
        while (1) {                                                                                               // 2399
            var buf = text.slice(buf_pos);                                                                        // 2400
            var p = buf.indexOf('\n');                                                                            // 2401
            if (p === -1) break;                                                                                  // 2402
            buf_pos += p+1;                                                                                       // 2403
            var msg = buf.slice(0, p);                                                                            // 2404
            that.dispatchEvent(new SimpleEvent('message', {data: msg}));                                          // 2405
        }                                                                                                         // 2406
    };                                                                                                            // 2407
    that.xo.onfinish = function(status, text) {                                                                   // 2408
        that.xo.onchunk(status, text);                                                                            // 2409
        that.xo = null;                                                                                           // 2410
        var reason = status === 200 ? 'network' : 'permanent';                                                    // 2411
        that.dispatchEvent(new SimpleEvent('close', {reason: reason}));                                           // 2412
    }                                                                                                             // 2413
};                                                                                                                // 2414
                                                                                                                  // 2415
XhrReceiver.prototype = new REventTarget();                                                                       // 2416
                                                                                                                  // 2417
XhrReceiver.prototype.abort = function() {                                                                        // 2418
    var that = this;                                                                                              // 2419
    if (that.xo) {                                                                                                // 2420
        that.xo.close();                                                                                          // 2421
        that.dispatchEvent(new SimpleEvent('close', {reason: 'user'}));                                           // 2422
        that.xo = null;                                                                                           // 2423
    }                                                                                                             // 2424
};                                                                                                                // 2425
//         [*] End of lib/trans-receiver-xhr.js                                                                   // 2426
                                                                                                                  // 2427
                                                                                                                  // 2428
//         [*] Including lib/test-hooks.js                                                                        // 2429
/*                                                                                                                // 2430
 * ***** BEGIN LICENSE BLOCK *****                                                                                // 2431
 * Copyright (c) 2011-2012 VMware, Inc.                                                                           // 2432
 *                                                                                                                // 2433
 * For the license see COPYING.                                                                                   // 2434
 * ***** END LICENSE BLOCK *****                                                                                  // 2435
 */                                                                                                               // 2436
                                                                                                                  // 2437
// For testing                                                                                                    // 2438
SockJS.getUtils = function(){                                                                                     // 2439
    return utils;                                                                                                 // 2440
};                                                                                                                // 2441
                                                                                                                  // 2442
SockJS.getIframeTransport = function(){                                                                           // 2443
    return IframeTransport;                                                                                       // 2444
};                                                                                                                // 2445
//         [*] End of lib/test-hooks.js                                                                           // 2446
                                                                                                                  // 2447
                  return SockJS;                                                                                  // 2448
          })();                                                                                                   // 2449
if ('_sockjs_onload' in window) setTimeout(_sockjs_onload, 1);                                                    // 2450
                                                                                                                  // 2451
// AMD compliance                                                                                                 // 2452
if (typeof define === 'function' && define.amd) {                                                                 // 2453
    define('sockjs', [], function(){return SockJS;});                                                             // 2454
}                                                                                                                 // 2455
//     [*] End of lib/index.js                                                                                    // 2456
                                                                                                                  // 2457
// [*] End of lib/all.js                                                                                          // 2458
                                                                                                                  // 2459
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/stream_client_sockjs.js                                                                    //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
// @param url {String} URL to Meteor app                                                                          // 1
//   "http://subdomain.meteor.com/" or "/" or                                                                     // 2
//   "ddp+sockjs://foo-**.meteor.com/sockjs"                                                                      // 3
LivedataTest.ClientStream = function (url, options) {                                                             // 4
  var self = this;                                                                                                // 5
  self.options = _.extend({                                                                                       // 6
    retry: true                                                                                                   // 7
  }, options);                                                                                                    // 8
  self._initCommon(self.options);                                                                                 // 9
                                                                                                                  // 10
  //// Constants                                                                                                  // 11
                                                                                                                  // 12
                                                                                                                  // 13
  // how long between hearing heartbeat from the server until we declare                                          // 14
  // the connection dead. heartbeats come every 45s (stream_server.js)                                            // 15
  //                                                                                                              // 16
  // NOTE: this is a older timeout mechanism. We now send heartbeats at                                           // 17
  // the DDP level (https://github.com/meteor/meteor/pull/1865), and                                              // 18
  // expect those timeouts to kill a non-responsive connection before                                             // 19
  // this timeout fires. This is kept around for compatibility (when                                              // 20
  // talking to a server that doesn't support DDP heartbeats) and can be                                          // 21
  // removed later.                                                                                               // 22
  self.HEARTBEAT_TIMEOUT = 100*1000;                                                                              // 23
                                                                                                                  // 24
  self.rawUrl = url;                                                                                              // 25
  self.socket = null;                                                                                             // 26
                                                                                                                  // 27
  self.heartbeatTimer = null;                                                                                     // 28
                                                                                                                  // 29
  // Listen to global 'online' event if we are running in a browser.                                              // 30
  // (IE8 does not support addEventListener)                                                                      // 31
  if (typeof window !== 'undefined' && window.addEventListener)                                                   // 32
    window.addEventListener("online", _.bind(self._online, self),                                                 // 33
                            false /* useCapture. make FF3.6 happy. */);                                           // 34
                                                                                                                  // 35
  //// Kickoff!                                                                                                   // 36
  self._launchConnection();                                                                                       // 37
};                                                                                                                // 38
                                                                                                                  // 39
_.extend(LivedataTest.ClientStream.prototype, {                                                                   // 40
                                                                                                                  // 41
  // data is a utf8 string. Data sent while not connected is dropped on                                           // 42
  // the floor, and it is up the user of this API to retransmit lost                                              // 43
  // messages on 'reset'                                                                                          // 44
  send: function (data) {                                                                                         // 45
    var self = this;                                                                                              // 46
    if (self.currentStatus.connected) {                                                                           // 47
      self.socket.send(data);                                                                                     // 48
    }                                                                                                             // 49
  },                                                                                                              // 50
                                                                                                                  // 51
  // Changes where this connection points                                                                         // 52
  _changeUrl: function (url) {                                                                                    // 53
    var self = this;                                                                                              // 54
    self.rawUrl = url;                                                                                            // 55
  },                                                                                                              // 56
                                                                                                                  // 57
  _connected: function () {                                                                                       // 58
    var self = this;                                                                                              // 59
                                                                                                                  // 60
    if (self.connectionTimer) {                                                                                   // 61
      clearTimeout(self.connectionTimer);                                                                         // 62
      self.connectionTimer = null;                                                                                // 63
    }                                                                                                             // 64
                                                                                                                  // 65
    if (self.currentStatus.connected) {                                                                           // 66
      // already connected. do nothing. this probably shouldn't happen.                                           // 67
      return;                                                                                                     // 68
    }                                                                                                             // 69
                                                                                                                  // 70
    // update status                                                                                              // 71
    self.currentStatus.status = "connected";                                                                      // 72
    self.currentStatus.connected = true;                                                                          // 73
    self.currentStatus.retryCount = 0;                                                                            // 74
    self.statusChanged();                                                                                         // 75
                                                                                                                  // 76
    // fire resets. This must come after status change so that clients                                            // 77
    // can call send from within a reset callback.                                                                // 78
    _.each(self.eventCallbacks.reset, function (callback) { callback(); });                                       // 79
                                                                                                                  // 80
  },                                                                                                              // 81
                                                                                                                  // 82
  _cleanup: function (maybeError) {                                                                               // 83
    var self = this;                                                                                              // 84
                                                                                                                  // 85
    self._clearConnectionAndHeartbeatTimers();                                                                    // 86
    if (self.socket) {                                                                                            // 87
      self.socket.onmessage = self.socket.onclose                                                                 // 88
        = self.socket.onerror = self.socket.onheartbeat = function () {};                                         // 89
      self.socket.close();                                                                                        // 90
      self.socket = null;                                                                                         // 91
    }                                                                                                             // 92
                                                                                                                  // 93
    _.each(self.eventCallbacks.disconnect, function (callback) {                                                  // 94
      callback(maybeError);                                                                                       // 95
    });                                                                                                           // 96
  },                                                                                                              // 97
                                                                                                                  // 98
  _clearConnectionAndHeartbeatTimers: function () {                                                               // 99
    var self = this;                                                                                              // 100
    if (self.connectionTimer) {                                                                                   // 101
      clearTimeout(self.connectionTimer);                                                                         // 102
      self.connectionTimer = null;                                                                                // 103
    }                                                                                                             // 104
    if (self.heartbeatTimer) {                                                                                    // 105
      clearTimeout(self.heartbeatTimer);                                                                          // 106
      self.heartbeatTimer = null;                                                                                 // 107
    }                                                                                                             // 108
  },                                                                                                              // 109
                                                                                                                  // 110
  _heartbeat_timeout: function () {                                                                               // 111
    var self = this;                                                                                              // 112
    Meteor._debug("Connection timeout. No sockjs heartbeat received.");                                           // 113
    self._lostConnection(new DDP.ConnectionError("Heartbeat timed out"));                                         // 114
  },                                                                                                              // 115
                                                                                                                  // 116
  _heartbeat_received: function () {                                                                              // 117
    var self = this;                                                                                              // 118
    // If we've already permanently shut down this stream, the timeout is                                         // 119
    // already cleared, and we don't need to set it again.                                                        // 120
    if (self._forcedToDisconnect)                                                                                 // 121
      return;                                                                                                     // 122
    if (self.heartbeatTimer)                                                                                      // 123
      clearTimeout(self.heartbeatTimer);                                                                          // 124
    self.heartbeatTimer = setTimeout(                                                                             // 125
      _.bind(self._heartbeat_timeout, self),                                                                      // 126
      self.HEARTBEAT_TIMEOUT);                                                                                    // 127
  },                                                                                                              // 128
                                                                                                                  // 129
  _sockjsProtocolsWhitelist: function () {                                                                        // 130
    // only allow polling protocols. no streaming.  streaming                                                     // 131
    // makes safari spin.                                                                                         // 132
    var protocolsWhitelist = [                                                                                    // 133
      'xdr-polling', 'xhr-polling', 'iframe-xhr-polling', 'jsonp-polling'];                                       // 134
                                                                                                                  // 135
    // iOS 4 and 5 and below crash when using websockets over certain                                             // 136
    // proxies. this seems to be resolved with iOS 6. eg                                                          // 137
    // https://github.com/LearnBoost/socket.io/issues/193#issuecomment-7308865.                                   // 138
    //                                                                                                            // 139
    // iOS <4 doesn't support websockets at all so sockjs will just                                               // 140
    // immediately fall back to http                                                                              // 141
    var noWebsockets = navigator &&                                                                               // 142
          /iPhone|iPad|iPod/.test(navigator.userAgent) &&                                                         // 143
          /OS 4_|OS 5_/.test(navigator.userAgent);                                                                // 144
                                                                                                                  // 145
    if (!noWebsockets)                                                                                            // 146
      protocolsWhitelist = ['websocket'].concat(protocolsWhitelist);                                              // 147
                                                                                                                  // 148
    return protocolsWhitelist;                                                                                    // 149
  },                                                                                                              // 150
                                                                                                                  // 151
  _launchConnection: function () {                                                                                // 152
    var self = this;                                                                                              // 153
    self._cleanup(); // cleanup the old socket, if there was one.                                                 // 154
                                                                                                                  // 155
    var options = _.extend({                                                                                      // 156
      protocols_whitelist:self._sockjsProtocolsWhitelist()                                                        // 157
    }, self.options._sockjsOptions);                                                                              // 158
                                                                                                                  // 159
    // Convert raw URL to SockJS URL each time we open a connection, so that we                                   // 160
    // can connect to random hostnames and get around browser per-host                                            // 161
    // connection limits.                                                                                         // 162
    self.socket = new SockJS(toSockjsUrl(self.rawUrl), undefined, options);                                       // 163
    self.socket.onopen = function (data) {                                                                        // 164
      self._connected();                                                                                          // 165
    };                                                                                                            // 166
    self.socket.onmessage = function (data) {                                                                     // 167
      self._heartbeat_received();                                                                                 // 168
                                                                                                                  // 169
      if (self.currentStatus.connected)                                                                           // 170
        _.each(self.eventCallbacks.message, function (callback) {                                                 // 171
          callback(data.data);                                                                                    // 172
        });                                                                                                       // 173
    };                                                                                                            // 174
    self.socket.onclose = function () {                                                                           // 175
      self._lostConnection();                                                                                     // 176
    };                                                                                                            // 177
    self.socket.onerror = function () {                                                                           // 178
      // XXX is this ever called?                                                                                 // 179
      Meteor._debug("stream error", _.toArray(arguments), (new Date()).toDateString());                           // 180
    };                                                                                                            // 181
                                                                                                                  // 182
    self.socket.onheartbeat =  function () {                                                                      // 183
      self._heartbeat_received();                                                                                 // 184
    };                                                                                                            // 185
                                                                                                                  // 186
    if (self.connectionTimer)                                                                                     // 187
      clearTimeout(self.connectionTimer);                                                                         // 188
    self.connectionTimer = setTimeout(function () {                                                               // 189
      self._lostConnection(                                                                                       // 190
        new DDP.ConnectionError("DDP connection timed out"));                                                     // 191
    }, self.CONNECT_TIMEOUT);                                                                                     // 192
  }                                                                                                               // 193
});                                                                                                               // 194
                                                                                                                  // 195
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/stream_client_common.js                                                                    //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
// XXX from Underscore.String (http://epeli.github.com/underscore.string/)                                        // 1
var startsWith = function(str, starts) {                                                                          // 2
  return str.length >= starts.length &&                                                                           // 3
    str.substring(0, starts.length) === starts;                                                                   // 4
};                                                                                                                // 5
var endsWith = function(str, ends) {                                                                              // 6
  return str.length >= ends.length &&                                                                             // 7
    str.substring(str.length - ends.length) === ends;                                                             // 8
};                                                                                                                // 9
                                                                                                                  // 10
// @param url {String} URL to Meteor app, eg:                                                                     // 11
//   "/" or "madewith.meteor.com" or "https://foo.meteor.com"                                                     // 12
//   or "ddp+sockjs://ddp--****-foo.meteor.com/sockjs"                                                            // 13
// @returns {String} URL to the endpoint with the specific scheme and subPath, e.g.                               // 14
// for scheme "http" and subPath "sockjs"                                                                         // 15
//   "http://subdomain.meteor.com/sockjs" or "/sockjs"                                                            // 16
//   or "https://ddp--1234-foo.meteor.com/sockjs"                                                                 // 17
var translateUrl =  function(url, newSchemeBase, subPath) {                                                       // 18
  if (! newSchemeBase) {                                                                                          // 19
    newSchemeBase = "http";                                                                                       // 20
  }                                                                                                               // 21
                                                                                                                  // 22
  var ddpUrlMatch = url.match(/^ddp(i?)\+sockjs:\/\//);                                                           // 23
  var httpUrlMatch = url.match(/^http(s?):\/\//);                                                                 // 24
  var newScheme;                                                                                                  // 25
  if (ddpUrlMatch) {                                                                                              // 26
    // Remove scheme and split off the host.                                                                      // 27
    var urlAfterDDP = url.substr(ddpUrlMatch[0].length);                                                          // 28
    newScheme = ddpUrlMatch[1] === "i" ? newSchemeBase : newSchemeBase + "s";                                     // 29
    var slashPos = urlAfterDDP.indexOf('/');                                                                      // 30
    var host =                                                                                                    // 31
          slashPos === -1 ? urlAfterDDP : urlAfterDDP.substr(0, slashPos);                                        // 32
    var rest = slashPos === -1 ? '' : urlAfterDDP.substr(slashPos);                                               // 33
                                                                                                                  // 34
    // In the host (ONLY!), change '*' characters into random digits. This                                        // 35
    // allows different stream connections to connect to different hostnames                                      // 36
    // and avoid browser per-hostname connection limits.                                                          // 37
    host = host.replace(/\*/g, function () {                                                                      // 38
      return Math.floor(Random.fraction()*10);                                                                    // 39
    });                                                                                                           // 40
                                                                                                                  // 41
    return newScheme + '://' + host + rest;                                                                       // 42
  } else if (httpUrlMatch) {                                                                                      // 43
    newScheme = !httpUrlMatch[1] ? newSchemeBase : newSchemeBase + "s";                                           // 44
    var urlAfterHttp = url.substr(httpUrlMatch[0].length);                                                        // 45
    url = newScheme + "://" + urlAfterHttp;                                                                       // 46
  }                                                                                                               // 47
                                                                                                                  // 48
  // Prefix FQDNs but not relative URLs                                                                           // 49
  if (url.indexOf("://") === -1 && !startsWith(url, "/")) {                                                       // 50
    url = newSchemeBase + "://" + url;                                                                            // 51
  }                                                                                                               // 52
                                                                                                                  // 53
  // XXX This is not what we should be doing: if I have a site                                                    // 54
  // deployed at "/foo", then DDP.connect("/") should actually connect                                            // 55
  // to "/", not to "/foo". "/" is an absolute path. (Contrast: if                                                // 56
  // deployed at "/foo", it would be reasonable for DDP.connect("bar")                                            // 57
  // to connect to "/foo/bar").                                                                                   // 58
  //                                                                                                              // 59
  // We should make this properly honor absolute paths rather than                                                // 60
  // forcing the path to be relative to the site root. Simultaneously,                                            // 61
  // we should set DDP_DEFAULT_CONNECTION_URL to include the site                                                 // 62
  // root. See also client_convenience.js #RationalizingRelativeDDPURLs                                           // 63
  url = Meteor._relativeToSiteRootUrl(url);                                                                       // 64
                                                                                                                  // 65
  if (endsWith(url, "/"))                                                                                         // 66
    return url + subPath;                                                                                         // 67
  else                                                                                                            // 68
    return url + "/" + subPath;                                                                                   // 69
};                                                                                                                // 70
                                                                                                                  // 71
toSockjsUrl = function (url) {                                                                                    // 72
  return translateUrl(url, "http", "sockjs");                                                                     // 73
};                                                                                                                // 74
                                                                                                                  // 75
toWebsocketUrl = function (url) {                                                                                 // 76
  var ret = translateUrl(url, "ws", "websocket");                                                                 // 77
  return ret;                                                                                                     // 78
};                                                                                                                // 79
                                                                                                                  // 80
LivedataTest.toSockjsUrl = toSockjsUrl;                                                                           // 81
                                                                                                                  // 82
                                                                                                                  // 83
_.extend(LivedataTest.ClientStream.prototype, {                                                                   // 84
                                                                                                                  // 85
  // Register for callbacks.                                                                                      // 86
  on: function (name, callback) {                                                                                 // 87
    var self = this;                                                                                              // 88
                                                                                                                  // 89
    if (name !== 'message' && name !== 'reset' && name !== 'disconnect')                                          // 90
      throw new Error("unknown event type: " + name);                                                             // 91
                                                                                                                  // 92
    if (!self.eventCallbacks[name])                                                                               // 93
      self.eventCallbacks[name] = [];                                                                             // 94
    self.eventCallbacks[name].push(callback);                                                                     // 95
  },                                                                                                              // 96
                                                                                                                  // 97
                                                                                                                  // 98
  _initCommon: function (options) {                                                                               // 99
    var self = this;                                                                                              // 100
    options = options || {};                                                                                      // 101
                                                                                                                  // 102
    //// Constants                                                                                                // 103
                                                                                                                  // 104
    // how long to wait until we declare the connection attempt                                                   // 105
    // failed.                                                                                                    // 106
    self.CONNECT_TIMEOUT = options.connectTimeoutMs || 10000;                                                     // 107
                                                                                                                  // 108
    self.eventCallbacks = {}; // name -> [callback]                                                               // 109
                                                                                                                  // 110
    self._forcedToDisconnect = false;                                                                             // 111
                                                                                                                  // 112
    //// Reactive status                                                                                          // 113
    self.currentStatus = {                                                                                        // 114
      status: "connecting",                                                                                       // 115
      connected: false,                                                                                           // 116
      retryCount: 0                                                                                               // 117
    };                                                                                                            // 118
                                                                                                                  // 119
                                                                                                                  // 120
    self.statusListeners = typeof Tracker !== 'undefined' && new Tracker.Dependency;                              // 121
    self.statusChanged = function () {                                                                            // 122
      if (self.statusListeners)                                                                                   // 123
        self.statusListeners.changed();                                                                           // 124
    };                                                                                                            // 125
                                                                                                                  // 126
    //// Retry logic                                                                                              // 127
    self._retry = new Retry;                                                                                      // 128
    self.connectionTimer = null;                                                                                  // 129
                                                                                                                  // 130
  },                                                                                                              // 131
                                                                                                                  // 132
  // Trigger a reconnect.                                                                                         // 133
  reconnect: function (options) {                                                                                 // 134
    var self = this;                                                                                              // 135
    options = options || {};                                                                                      // 136
                                                                                                                  // 137
    if (options.url) {                                                                                            // 138
      self._changeUrl(options.url);                                                                               // 139
    }                                                                                                             // 140
                                                                                                                  // 141
    if (options._sockjsOptions) {                                                                                 // 142
      self.options._sockjsOptions = options._sockjsOptions;                                                       // 143
    }                                                                                                             // 144
                                                                                                                  // 145
    if (self.currentStatus.connected) {                                                                           // 146
      if (options._force || options.url) {                                                                        // 147
        // force reconnect.                                                                                       // 148
        self._lostConnection(new DDP.ForcedReconnectError);                                                       // 149
      } // else, noop.                                                                                            // 150
      return;                                                                                                     // 151
    }                                                                                                             // 152
                                                                                                                  // 153
    // if we're mid-connection, stop it.                                                                          // 154
    if (self.currentStatus.status === "connecting") {                                                             // 155
      // Pretend it's a clean close.                                                                              // 156
      self._lostConnection();                                                                                     // 157
    }                                                                                                             // 158
                                                                                                                  // 159
    self._retry.clear();                                                                                          // 160
    self.currentStatus.retryCount -= 1; // don't count manual retries                                             // 161
    self._retryNow();                                                                                             // 162
  },                                                                                                              // 163
                                                                                                                  // 164
  disconnect: function (options) {                                                                                // 165
    var self = this;                                                                                              // 166
    options = options || {};                                                                                      // 167
                                                                                                                  // 168
    // Failed is permanent. If we're failed, don't let people go back                                             // 169
    // online by calling 'disconnect' then 'reconnect'.                                                           // 170
    if (self._forcedToDisconnect)                                                                                 // 171
      return;                                                                                                     // 172
                                                                                                                  // 173
    // If _permanent is set, permanently disconnect a stream. Once a stream                                       // 174
    // is forced to disconnect, it can never reconnect. This is for                                               // 175
    // error cases such as ddp version mismatch, where trying again                                               // 176
    // won't fix the problem.                                                                                     // 177
    if (options._permanent) {                                                                                     // 178
      self._forcedToDisconnect = true;                                                                            // 179
    }                                                                                                             // 180
                                                                                                                  // 181
    self._cleanup();                                                                                              // 182
    self._retry.clear();                                                                                          // 183
                                                                                                                  // 184
    self.currentStatus = {                                                                                        // 185
      status: (options._permanent ? "failed" : "offline"),                                                        // 186
      connected: false,                                                                                           // 187
      retryCount: 0                                                                                               // 188
    };                                                                                                            // 189
                                                                                                                  // 190
    if (options._permanent && options._error)                                                                     // 191
      self.currentStatus.reason = options._error;                                                                 // 192
                                                                                                                  // 193
    self.statusChanged();                                                                                         // 194
  },                                                                                                              // 195
                                                                                                                  // 196
  // maybeError is set unless it's a clean protocol-level close.                                                  // 197
  _lostConnection: function (maybeError) {                                                                        // 198
    var self = this;                                                                                              // 199
                                                                                                                  // 200
    self._cleanup(maybeError);                                                                                    // 201
    self._retryLater(maybeError); // sets status. no need to do it here.                                          // 202
  },                                                                                                              // 203
                                                                                                                  // 204
  // fired when we detect that we've gone online. try to reconnect                                                // 205
  // immediately.                                                                                                 // 206
  _online: function () {                                                                                          // 207
    // if we've requested to be offline by disconnecting, don't reconnect.                                        // 208
    if (this.currentStatus.status != "offline")                                                                   // 209
      this.reconnect();                                                                                           // 210
  },                                                                                                              // 211
                                                                                                                  // 212
  _retryLater: function (maybeError) {                                                                            // 213
    var self = this;                                                                                              // 214
                                                                                                                  // 215
    var timeout = 0;                                                                                              // 216
    if (self.options.retry ||                                                                                     // 217
        (maybeError && maybeError.errorType === "DDP.ForcedReconnectError")) {                                    // 218
      timeout = self._retry.retryLater(                                                                           // 219
        self.currentStatus.retryCount,                                                                            // 220
        _.bind(self._retryNow, self)                                                                              // 221
      );                                                                                                          // 222
      self.currentStatus.status = "waiting";                                                                      // 223
      self.currentStatus.retryTime = (new Date()).getTime() + timeout;                                            // 224
    } else {                                                                                                      // 225
      self.currentStatus.status = "failed";                                                                       // 226
      delete self.currentStatus.retryTime;                                                                        // 227
    }                                                                                                             // 228
                                                                                                                  // 229
    self.currentStatus.connected = false;                                                                         // 230
    self.statusChanged();                                                                                         // 231
  },                                                                                                              // 232
                                                                                                                  // 233
  _retryNow: function () {                                                                                        // 234
    var self = this;                                                                                              // 235
                                                                                                                  // 236
    if (self._forcedToDisconnect)                                                                                 // 237
      return;                                                                                                     // 238
                                                                                                                  // 239
    self.currentStatus.retryCount += 1;                                                                           // 240
    self.currentStatus.status = "connecting";                                                                     // 241
    self.currentStatus.connected = false;                                                                         // 242
    delete self.currentStatus.retryTime;                                                                          // 243
    self.statusChanged();                                                                                         // 244
                                                                                                                  // 245
    self._launchConnection();                                                                                     // 246
  },                                                                                                              // 247
                                                                                                                  // 248
                                                                                                                  // 249
  // Get current status. Reactive.                                                                                // 250
  status: function () {                                                                                           // 251
    var self = this;                                                                                              // 252
    if (self.statusListeners)                                                                                     // 253
      self.statusListeners.depend();                                                                              // 254
    return self.currentStatus;                                                                                    // 255
  }                                                                                                               // 256
});                                                                                                               // 257
                                                                                                                  // 258
DDP.ConnectionError = Meteor.makeErrorType(                                                                       // 259
  "DDP.ConnectionError", function (message) {                                                                     // 260
    var self = this;                                                                                              // 261
    self.message = message;                                                                                       // 262
});                                                                                                               // 263
                                                                                                                  // 264
DDP.ForcedReconnectError = Meteor.makeErrorType(                                                                  // 265
  "DDP.ForcedReconnectError", function () {});                                                                    // 266
                                                                                                                  // 267
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/livedata_common.js                                                                         //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
LivedataTest.SUPPORTED_DDP_VERSIONS = DDPCommon.SUPPORTED_DDP_VERSIONS;                                           // 1
                                                                                                                  // 2
// This is private but it's used in a few places. accounts-base uses                                              // 3
// it to get the current user. Meteor.setTimeout and friends clear                                                // 4
// it. We can probably find a better way to factor this.                                                          // 5
DDP._CurrentInvocation = new Meteor.EnvironmentVariable;                                                          // 6
                                                                                                                  // 7
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/random_stream.js                                                                           //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
// Returns the named sequence of pseudo-random values.                                                            // 1
// The scope will be DDP._CurrentInvocation.get(), so the stream will produce                                     // 2
// consistent values for method calls on the client and server.                                                   // 3
DDP.randomStream = function (name) {                                                                              // 4
  var scope = DDP._CurrentInvocation.get();                                                                       // 5
  return DDPCommon.RandomStream.get(scope, name);                                                                 // 6
};                                                                                                                // 7
                                                                                                                  // 8
                                                                                                                  // 9
                                                                                                                  // 10
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/livedata_connection.js                                                                     //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
if (Meteor.isServer) {                                                                                            // 1
  var path = Npm.require('path');                                                                                 // 2
  var Fiber = Npm.require('fibers');                                                                              // 3
  var Future = Npm.require(path.join('fibers', 'future'));                                                        // 4
}                                                                                                                 // 5
                                                                                                                  // 6
// @param url {String|Object} URL to Meteor app,                                                                  // 7
//   or an object as a test hook (see code)                                                                       // 8
// Options:                                                                                                       // 9
//   reloadWithOutstanding: is it OK to reload if there are outstanding methods?                                  // 10
//   headers: extra headers to send on the websockets connection, for                                             // 11
//     server-to-server DDP only                                                                                  // 12
//   _sockjsOptions: Specifies options to pass through to the sockjs client                                       // 13
//   onDDPNegotiationVersionFailure: callback when version negotiation fails.                                     // 14
//                                                                                                                // 15
// XXX There should be a way to destroy a DDP connection, causing all                                             // 16
// outstanding method calls to fail.                                                                              // 17
//                                                                                                                // 18
// XXX Our current way of handling failure and reconnection is great                                              // 19
// for an app (where we want to tolerate being disconnected as an                                                 // 20
// expect state, and keep trying forever to reconnect) but cumbersome                                             // 21
// for something like a command line tool that wants to make a                                                    // 22
// connection, call a method, and print an error if connection                                                    // 23
// fails. We should have better usability in the latter case (while                                               // 24
// still transparently reconnecting if it's just a transient failure                                              // 25
// or the server migrating us).                                                                                   // 26
var Connection = function (url, options) {                                                                        // 27
  var self = this;                                                                                                // 28
  options = _.extend({                                                                                            // 29
    onConnected: function () {},                                                                                  // 30
    onDDPVersionNegotiationFailure: function (description) {                                                      // 31
      Meteor._debug(description);                                                                                 // 32
    },                                                                                                            // 33
    heartbeatInterval: 17500,                                                                                     // 34
    heartbeatTimeout: 15000,                                                                                      // 35
    npmFayeOptions: {},                                                                                           // 36
    // These options are only for testing.                                                                        // 37
    reloadWithOutstanding: false,                                                                                 // 38
    supportedDDPVersions: DDPCommon.SUPPORTED_DDP_VERSIONS,                                                       // 39
    retry: true,                                                                                                  // 40
    respondToPings: true                                                                                          // 41
  }, options);                                                                                                    // 42
                                                                                                                  // 43
  // If set, called when we reconnect, queuing method calls _before_ the                                          // 44
  // existing outstanding ones. This is the only data member that is part of the                                  // 45
  // public API!                                                                                                  // 46
  self.onReconnect = null;                                                                                        // 47
                                                                                                                  // 48
  // as a test hook, allow passing a stream instead of a url.                                                     // 49
  if (typeof url === "object") {                                                                                  // 50
    self._stream = url;                                                                                           // 51
  } else {                                                                                                        // 52
    self._stream = new LivedataTest.ClientStream(url, {                                                           // 53
      retry: options.retry,                                                                                       // 54
      headers: options.headers,                                                                                   // 55
      _sockjsOptions: options._sockjsOptions,                                                                     // 56
      // Used to keep some tests quiet, or for other cases in which                                               // 57
      // the right thing to do with connection errors is to silently                                              // 58
      // fail (e.g. sending package usage stats). At some point we                                                // 59
      // should have a real API for handling client-stream-level                                                  // 60
      // errors.                                                                                                  // 61
      _dontPrintErrors: options._dontPrintErrors,                                                                 // 62
      connectTimeoutMs: options.connectTimeoutMs,                                                                 // 63
      npmFayeOptions: options.npmFayeOptions                                                                      // 64
    });                                                                                                           // 65
  }                                                                                                               // 66
                                                                                                                  // 67
  self._lastSessionId = null;                                                                                     // 68
  self._versionSuggestion = null;  // The last proposed DDP version.                                              // 69
  self._version = null;   // The DDP version agreed on by client and server.                                      // 70
  self._stores = {}; // name -> object with methods                                                               // 71
  self._methodHandlers = {}; // name -> func                                                                      // 72
  self._nextMethodId = 1;                                                                                         // 73
  self._supportedDDPVersions = options.supportedDDPVersions;                                                      // 74
                                                                                                                  // 75
  self._heartbeatInterval = options.heartbeatInterval;                                                            // 76
  self._heartbeatTimeout = options.heartbeatTimeout;                                                              // 77
                                                                                                                  // 78
  // Tracks methods which the user has tried to call but which have not yet                                       // 79
  // called their user callback (ie, they are waiting on their result or for all                                  // 80
  // of their writes to be written to the local cache). Map from method ID to                                     // 81
  // MethodInvoker object.                                                                                        // 82
  self._methodInvokers = {};                                                                                      // 83
                                                                                                                  // 84
  // Tracks methods which the user has called but whose result messages have not                                  // 85
  // arrived yet.                                                                                                 // 86
  //                                                                                                              // 87
  // _outstandingMethodBlocks is an array of blocks of methods. Each block                                        // 88
  // represents a set of methods that can run at the same time. The first block                                   // 89
  // represents the methods which are currently in flight; subsequent blocks                                      // 90
  // must wait for previous blocks to be fully finished before they can be sent                                   // 91
  // to the server.                                                                                               // 92
  //                                                                                                              // 93
  // Each block is an object with the following fields:                                                           // 94
  // - methods: a list of MethodInvoker objects                                                                   // 95
  // - wait: a boolean; if true, this block had a single method invoked with                                      // 96
  //         the "wait" option                                                                                    // 97
  //                                                                                                              // 98
  // There will never be adjacent blocks with wait=false, because the only thing                                  // 99
  // that makes methods need to be serialized is a wait method.                                                   // 100
  //                                                                                                              // 101
  // Methods are removed from the first block when their "result" is                                              // 102
  // received. The entire first block is only removed when all of the in-flight                                   // 103
  // methods have received their results (so the "methods" list is empty) *AND*                                   // 104
  // all of the data written by those methods are visible in the local cache. So                                  // 105
  // it is possible for the first block's methods list to be empty, if we are                                     // 106
  // still waiting for some objects to quiesce.                                                                   // 107
  //                                                                                                              // 108
  // Example:                                                                                                     // 109
  //  _outstandingMethodBlocks = [                                                                                // 110
  //    {wait: false, methods: []},                                                                               // 111
  //    {wait: true, methods: [<MethodInvoker for 'login'>]},                                                     // 112
  //    {wait: false, methods: [<MethodInvoker for 'foo'>,                                                        // 113
  //                            <MethodInvoker for 'bar'>]}]                                                      // 114
  // This means that there were some methods which were sent to the server and                                    // 115
  // which have returned their results, but some of the data written by                                           // 116
  // the methods may not be visible in the local cache. Once all that data is                                     // 117
  // visible, we will send a 'login' method. Once the login method has returned                                   // 118
  // and all the data is visible (including re-running subs if userId changes),                                   // 119
  // we will send the 'foo' and 'bar' methods in parallel.                                                        // 120
  self._outstandingMethodBlocks = [];                                                                             // 121
                                                                                                                  // 122
  // method ID -> array of objects with keys 'collection' and 'id', listing                                       // 123
  // documents written by a given method's stub. keys are associated with                                         // 124
  // methods whose stub wrote at least one document, and whose data-done message                                  // 125
  // has not yet been received.                                                                                   // 126
  self._documentsWrittenByStub = {};                                                                              // 127
  // collection -> IdMap of "server document" object. A "server document" has:                                    // 128
  // - "document": the version of the document according the                                                      // 129
  //   server (ie, the snapshot before a stub wrote it, amended by any changes                                    // 130
  //   received from the server)                                                                                  // 131
  //   It is undefined if we think the document does not exist                                                    // 132
  // - "writtenByStubs": a set of method IDs whose stubs wrote to the document                                    // 133
  //   whose "data done" messages have not yet been processed                                                     // 134
  self._serverDocuments = {};                                                                                     // 135
                                                                                                                  // 136
  // Array of callbacks to be called after the next update of the local                                           // 137
  // cache. Used for:                                                                                             // 138
  //  - Calling methodInvoker.dataVisible and sub ready callbacks after                                           // 139
  //    the relevant data is flushed.                                                                             // 140
  //  - Invoking the callbacks of "half-finished" methods after reconnect                                         // 141
  //    quiescence. Specifically, methods whose result was received over the old                                  // 142
  //    connection (so we don't re-send it) but whose data had not been made                                      // 143
  //    visible.                                                                                                  // 144
  self._afterUpdateCallbacks = [];                                                                                // 145
                                                                                                                  // 146
  // In two contexts, we buffer all incoming data messages and then process them                                  // 147
  // all at once in a single update:                                                                              // 148
  //   - During reconnect, we buffer all data messages until all subs that had                                    // 149
  //     been ready before reconnect are ready again, and all methods that are                                    // 150
  //     active have returned their "data done message"; then                                                     // 151
  //   - During the execution of a "wait" method, we buffer all data messages                                     // 152
  //     until the wait method gets its "data done" message. (If the wait method                                  // 153
  //     occurs during reconnect, it doesn't get any special handling.)                                           // 154
  // all data messages are processed in one update.                                                               // 155
  //                                                                                                              // 156
  // The following fields are used for this "quiescence" process.                                                 // 157
                                                                                                                  // 158
  // This buffers the messages that aren't being processed yet.                                                   // 159
  self._messagesBufferedUntilQuiescence = [];                                                                     // 160
  // Map from method ID -> true. Methods are removed from this when their                                         // 161
  // "data done" message is received, and we will not quiesce until it is                                         // 162
  // empty.                                                                                                       // 163
  self._methodsBlockingQuiescence = {};                                                                           // 164
  // map from sub ID -> true for subs that were ready (ie, called the sub                                         // 165
  // ready callback) before reconnect but haven't become ready again yet                                          // 166
  self._subsBeingRevived = {}; // map from sub._id -> true                                                        // 167
  // if true, the next data update should reset all stores. (set during                                           // 168
  // reconnect.)                                                                                                  // 169
  self._resetStores = false;                                                                                      // 170
                                                                                                                  // 171
  // name -> array of updates for (yet to be created) collections                                                 // 172
  self._updatesForUnknownStores = {};                                                                             // 173
  // if we're blocking a migration, the retry func                                                                // 174
  self._retryMigrate = null;                                                                                      // 175
                                                                                                                  // 176
  // metadata for subscriptions.  Map from sub ID to object with keys:                                            // 177
  //   - id                                                                                                       // 178
  //   - name                                                                                                     // 179
  //   - params                                                                                                   // 180
  //   - inactive (if true, will be cleaned up if not reused in re-run)                                           // 181
  //   - ready (has the 'ready' message been received?)                                                           // 182
  //   - readyCallback (an optional callback to call when ready)                                                  // 183
  //   - errorCallback (an optional callback to call if the sub terminates with                                   // 184
  //                    an error, XXX COMPAT WITH 1.0.3.1)                                                        // 185
  //   - stopCallback (an optional callback to call when the sub terminates                                       // 186
  //     for any reason, with an error argument if an error triggered the stop)                                   // 187
  self._subscriptions = {};                                                                                       // 188
                                                                                                                  // 189
  // Reactive userId.                                                                                             // 190
  self._userId = null;                                                                                            // 191
  self._userIdDeps = new Tracker.Dependency;                                                                      // 192
                                                                                                                  // 193
  // Block auto-reload while we're waiting for method responses.                                                  // 194
  if (Meteor.isClient && Package.reload && !options.reloadWithOutstanding) {                                      // 195
    Package.reload.Reload._onMigrate(function (retry) {                                                           // 196
      if (!self._readyToMigrate()) {                                                                              // 197
        if (self._retryMigrate)                                                                                   // 198
          throw new Error("Two migrations in progress?");                                                         // 199
        self._retryMigrate = retry;                                                                               // 200
        return false;                                                                                             // 201
      } else {                                                                                                    // 202
        return [true];                                                                                            // 203
      }                                                                                                           // 204
    });                                                                                                           // 205
  }                                                                                                               // 206
                                                                                                                  // 207
  var onMessage = function (raw_msg) {                                                                            // 208
    try {                                                                                                         // 209
      var msg = DDPCommon.parseDDP(raw_msg);                                                                      // 210
    } catch (e) {                                                                                                 // 211
      Meteor._debug("Exception while parsing DDP", e);                                                            // 212
      return;                                                                                                     // 213
    }                                                                                                             // 214
                                                                                                                  // 215
    // Any message counts as receiving a pong, as it demonstrates that                                            // 216
    // the server is still alive.                                                                                 // 217
    if (self._heartbeat) {                                                                                        // 218
      self._heartbeat.messageReceived();                                                                          // 219
    }                                                                                                             // 220
                                                                                                                  // 221
    if (msg === null || !msg.msg) {                                                                               // 222
      // XXX COMPAT WITH 0.6.6. ignore the old welcome message for back                                           // 223
      // compat.  Remove this 'if' once the server stops sending welcome                                          // 224
      // messages (stream_server.js).                                                                             // 225
      if (! (msg && msg.server_id))                                                                               // 226
        Meteor._debug("discarding invalid livedata message", msg);                                                // 227
      return;                                                                                                     // 228
    }                                                                                                             // 229
                                                                                                                  // 230
    if (msg.msg === 'connected') {                                                                                // 231
      self._version = self._versionSuggestion;                                                                    // 232
      self._livedata_connected(msg);                                                                              // 233
      options.onConnected();                                                                                      // 234
    }                                                                                                             // 235
    else if (msg.msg === 'failed') {                                                                              // 236
      if (_.contains(self._supportedDDPVersions, msg.version)) {                                                  // 237
        self._versionSuggestion = msg.version;                                                                    // 238
        self._stream.reconnect({_force: true});                                                                   // 239
      } else {                                                                                                    // 240
        var description =                                                                                         // 241
              "DDP version negotiation failed; server requested version " + msg.version;                          // 242
        self._stream.disconnect({_permanent: true, _error: description});                                         // 243
        options.onDDPVersionNegotiationFailure(description);                                                      // 244
      }                                                                                                           // 245
    }                                                                                                             // 246
    else if (msg.msg === 'ping' && options.respondToPings) {                                                      // 247
      self._send({msg: "pong", id: msg.id});                                                                      // 248
    }                                                                                                             // 249
    else if (msg.msg === 'pong') {                                                                                // 250
      // noop, as we assume everything's a pong                                                                   // 251
    }                                                                                                             // 252
    else if (_.include(['added', 'changed', 'removed', 'ready', 'updated'], msg.msg))                             // 253
      self._livedata_data(msg);                                                                                   // 254
    else if (msg.msg === 'nosub')                                                                                 // 255
      self._livedata_nosub(msg);                                                                                  // 256
    else if (msg.msg === 'result')                                                                                // 257
      self._livedata_result(msg);                                                                                 // 258
    else if (msg.msg === 'error')                                                                                 // 259
      self._livedata_error(msg);                                                                                  // 260
    else                                                                                                          // 261
      Meteor._debug("discarding unknown livedata message type", msg);                                             // 262
  };                                                                                                              // 263
                                                                                                                  // 264
  var onReset = function () {                                                                                     // 265
    // Send a connect message at the beginning of the stream.                                                     // 266
    // NOTE: reset is called even on the first connection, so this is                                             // 267
    // the only place we send this message.                                                                       // 268
    var msg = {msg: 'connect'};                                                                                   // 269
    if (self._lastSessionId)                                                                                      // 270
      msg.session = self._lastSessionId;                                                                          // 271
    msg.version = self._versionSuggestion || self._supportedDDPVersions[0];                                       // 272
    self._versionSuggestion = msg.version;                                                                        // 273
    msg.support = self._supportedDDPVersions;                                                                     // 274
    self._send(msg);                                                                                              // 275
                                                                                                                  // 276
    // Mark non-retry calls as failed. This has to be done early as getting these methods out of the              // 277
    // current block is pretty important to making sure that quiescence is properly calculated, as                // 278
    // well as possibly moving on to another useful block.                                                        // 279
                                                                                                                  // 280
    // Only bother testing if there is an outstandingMethodBlock (there might not be, especially if               // 281
    // we are connecting for the first time.                                                                      // 282
    if (self._outstandingMethodBlocks.length > 0) {                                                               // 283
      // If there is an outstanding method block, we only care about the first one as that is the                 // 284
      // one that could have already sent messages with no response, that are not allowed to retry.               // 285
      _.each(self._outstandingMethodBlocks[0].methods, function(methodInvoker) {                                  // 286
        // If the message wasn't sent or it's allowed to retry, do nothing.                                       // 287
        if (methodInvoker.sentMessage && methodInvoker.noRetry) {                                                 // 288
          // The next loop serves to get the index in the current method block of this method.                    // 289
          var currentMethodBlock = self._outstandingMethodBlocks[0].methods;                                      // 290
          var loopMethod;                                                                                         // 291
          for (var i = 0; i < currentMethodBlock.length; i++) {                                                   // 292
            loopMethod = currentMethodBlock[i];                                                                   // 293
            if (loopMethod.methodId === methodInvoker.methodId) {                                                 // 294
              break;                                                                                              // 295
            }                                                                                                     // 296
          }                                                                                                       // 297
                                                                                                                  // 298
          // Remove from current method block. This may leave the block empty, but we                             // 299
          // don't move on to the next block until the callback has been delivered, in                            // 300
          // _outstandingMethodFinished.                                                                          // 301
          currentMethodBlock.splice(i, 1);                                                                        // 302
                                                                                                                  // 303
          // Make sure that the method is told that it failed.                                                    // 304
          methodInvoker.receiveResult(Meteor.Error('invocation-failed',                                           // 305
            'Method invocation might have failed due to dropped connection. ' +                                   // 306
            'Failing because `noRetry` option was passed to Meteor.apply.'));                                     // 307
        }                                                                                                         // 308
      });                                                                                                         // 309
    }                                                                                                             // 310
                                                                                                                  // 311
    // Now, to minimize setup latency, go ahead and blast out all of                                              // 312
    // our pending methods ands subscriptions before we've even taken                                             // 313
    // the necessary RTT to know if we successfully reconnected. (1)                                              // 314
    // They're supposed to be idempotent, and where they are not,                                                 // 315
    // they can block retry in apply; (2) even if we did reconnect,                                               // 316
    // we're not sure what messages might have gotten lost                                                        // 317
    // (in either direction) since we were disconnected (TCP being                                                // 318
    // sloppy about that.)                                                                                        // 319
                                                                                                                  // 320
    // If the current block of methods all got their results (but didn't all get                                  // 321
    // their data visible), discard the empty block now.                                                          // 322
    if (! _.isEmpty(self._outstandingMethodBlocks) &&                                                             // 323
        _.isEmpty(self._outstandingMethodBlocks[0].methods)) {                                                    // 324
      self._outstandingMethodBlocks.shift();                                                                      // 325
    }                                                                                                             // 326
                                                                                                                  // 327
    // Mark all messages as unsent, they have not yet been sent on this                                           // 328
    // connection.                                                                                                // 329
    _.each(self._methodInvokers, function (m) {                                                                   // 330
      m.sentMessage = false;                                                                                      // 331
    });                                                                                                           // 332
                                                                                                                  // 333
    // If an `onReconnect` handler is set, call it first. Go through                                              // 334
    // some hoops to ensure that methods that are called from within                                              // 335
    // `onReconnect` get executed _before_ ones that were originally                                              // 336
    // outstanding (since `onReconnect` is used to re-establish auth                                              // 337
    // certificates)                                                                                              // 338
    if (self.onReconnect)                                                                                         // 339
      self._callOnReconnectAndSendAppropriateOutstandingMethods();                                                // 340
    else                                                                                                          // 341
      self._sendOutstandingMethods();                                                                             // 342
                                                                                                                  // 343
    // add new subscriptions at the end. this way they take effect after                                          // 344
    // the handlers and we don't see flicker.                                                                     // 345
    _.each(self._subscriptions, function (sub, id) {                                                              // 346
      self._send({                                                                                                // 347
        msg: 'sub',                                                                                               // 348
        id: id,                                                                                                   // 349
        name: sub.name,                                                                                           // 350
        params: sub.params                                                                                        // 351
      });                                                                                                         // 352
    });                                                                                                           // 353
  };                                                                                                              // 354
                                                                                                                  // 355
  var onDisconnect = function () {                                                                                // 356
    if (self._heartbeat) {                                                                                        // 357
      self._heartbeat.stop();                                                                                     // 358
      self._heartbeat = null;                                                                                     // 359
    }                                                                                                             // 360
  };                                                                                                              // 361
                                                                                                                  // 362
  if (Meteor.isServer) {                                                                                          // 363
    self._stream.on('message', Meteor.bindEnvironment(onMessage, "handling DDP message"));                        // 364
    self._stream.on('reset', Meteor.bindEnvironment(onReset, "handling DDP reset"));                              // 365
    self._stream.on('disconnect', Meteor.bindEnvironment(onDisconnect, "handling DDP disconnect"));               // 366
  } else {                                                                                                        // 367
    self._stream.on('message', onMessage);                                                                        // 368
    self._stream.on('reset', onReset);                                                                            // 369
    self._stream.on('disconnect', onDisconnect);                                                                  // 370
  }                                                                                                               // 371
};                                                                                                                // 372
                                                                                                                  // 373
// A MethodInvoker manages sending a method to the server and calling the user's                                  // 374
// callbacks. On construction, it registers itself in the connection's                                            // 375
// _methodInvokers map; it removes itself once the method is fully finished and                                   // 376
// the callback is invoked. This occurs when it has both received a result,                                       // 377
// and the data written by it is fully visible.                                                                   // 378
var MethodInvoker = function (options) {                                                                          // 379
  var self = this;                                                                                                // 380
                                                                                                                  // 381
  // Public (within this file) fields.                                                                            // 382
  self.methodId = options.methodId;                                                                               // 383
  self.sentMessage = false;                                                                                       // 384
                                                                                                                  // 385
  self._callback = options.callback;                                                                              // 386
  self._connection = options.connection;                                                                          // 387
  self._message = options.message;                                                                                // 388
  self._onResultReceived = options.onResultReceived || function () {};                                            // 389
  self._wait = options.wait;                                                                                      // 390
  self.noRetry = options.noRetry;                                                                                 // 391
  self._methodResult = null;                                                                                      // 392
  self._dataVisible = false;                                                                                      // 393
                                                                                                                  // 394
  // Register with the connection.                                                                                // 395
  self._connection._methodInvokers[self.methodId] = self;                                                         // 396
};                                                                                                                // 397
_.extend(MethodInvoker.prototype, {                                                                               // 398
  // Sends the method message to the server. May be called additional times if                                    // 399
  // we lose the connection and reconnect before receiving a result.                                              // 400
  sendMessage: function () {                                                                                      // 401
    var self = this;                                                                                              // 402
    // This function is called before sending a method (including resending on                                    // 403
    // reconnect). We should only (re)send methods where we don't already have a                                  // 404
    // result!                                                                                                    // 405
    if (self.gotResult())                                                                                         // 406
      throw new Error("sendingMethod is called on method with result");                                           // 407
                                                                                                                  // 408
                                                                                                                  // 409
    // If we're re-sending it, it doesn't matter if data was written the first                                    // 410
    // time.                                                                                                      // 411
    self._dataVisible = false;                                                                                    // 412
    self.sentMessage = true;                                                                                      // 413
                                                                                                                  // 414
    // If this is a wait method, make all data messages be buffered until it is                                   // 415
    // done.                                                                                                      // 416
    if (self._wait)                                                                                               // 417
      self._connection._methodsBlockingQuiescence[self.methodId] = true;                                          // 418
                                                                                                                  // 419
    // Actually send the message.                                                                                 // 420
    self._connection._send(self._message);                                                                        // 421
  },                                                                                                              // 422
  // Invoke the callback, if we have both a result and know that all data has                                     // 423
  // been written to the local cache.                                                                             // 424
  _maybeInvokeCallback: function () {                                                                             // 425
    var self = this;                                                                                              // 426
    if (self._methodResult && self._dataVisible) {                                                                // 427
      // Call the callback. (This won't throw: the callback was wrapped with                                      // 428
      // bindEnvironment.)                                                                                        // 429
      self._callback(self._methodResult[0], self._methodResult[1]);                                               // 430
                                                                                                                  // 431
      // Forget about this method.                                                                                // 432
      delete self._connection._methodInvokers[self.methodId];                                                     // 433
                                                                                                                  // 434
      // Let the connection know that this method is finished, so it can try to                                   // 435
      // move on to the next block of methods.                                                                    // 436
      self._connection._outstandingMethodFinished();                                                              // 437
    }                                                                                                             // 438
  },                                                                                                              // 439
  // Call with the result of the method from the server. Only may be called                                       // 440
  // once; once it is called, you should not call sendMessage again.                                              // 441
  // If the user provided an onResultReceived callback, call it immediately.                                      // 442
  // Then invoke the main callback if data is also visible.                                                       // 443
  receiveResult: function (err, result) {                                                                         // 444
    var self = this;                                                                                              // 445
    if (self.gotResult())                                                                                         // 446
      throw new Error("Methods should only receive results once");                                                // 447
    self._methodResult = [err, result];                                                                           // 448
    self._onResultReceived(err, result);                                                                          // 449
    self._maybeInvokeCallback();                                                                                  // 450
  },                                                                                                              // 451
  // Call this when all data written by the method is visible. This means that                                    // 452
  // the method has returns its "data is done" message *AND* all server                                           // 453
  // documents that are buffered at that time have been written to the local                                      // 454
  // cache. Invokes the main callback if the result has been received.                                            // 455
  dataVisible: function () {                                                                                      // 456
    var self = this;                                                                                              // 457
    self._dataVisible = true;                                                                                     // 458
    self._maybeInvokeCallback();                                                                                  // 459
  },                                                                                                              // 460
  // True if receiveResult has been called.                                                                       // 461
  gotResult: function () {                                                                                        // 462
    var self = this;                                                                                              // 463
    return !!self._methodResult;                                                                                  // 464
  }                                                                                                               // 465
});                                                                                                               // 466
                                                                                                                  // 467
_.extend(Connection.prototype, {                                                                                  // 468
  // 'name' is the name of the data on the wire that should go in the                                             // 469
  // store. 'wrappedStore' should be an object with methods beginUpdate, update,                                  // 470
  // endUpdate, saveOriginals, retrieveOriginals. see Collection for an example.                                  // 471
  registerStore: function (name, wrappedStore) {                                                                  // 472
    var self = this;                                                                                              // 473
                                                                                                                  // 474
    if (name in self._stores)                                                                                     // 475
      return false;                                                                                               // 476
                                                                                                                  // 477
    // Wrap the input object in an object which makes any store method not                                        // 478
    // implemented by 'store' into a no-op.                                                                       // 479
    var store = {};                                                                                               // 480
    _.each(['update', 'beginUpdate', 'endUpdate', 'saveOriginals',                                                // 481
            'retrieveOriginals', 'getDoc',                                                                        // 482
			'_getCollection'], function (method) {                                                                         // 483
              store[method] = function () {                                                                       // 484
                return (wrappedStore[method]                                                                      // 485
                        ? wrappedStore[method].apply(wrappedStore, arguments)                                     // 486
                        : undefined);                                                                             // 487
              };                                                                                                  // 488
            });                                                                                                   // 489
                                                                                                                  // 490
    self._stores[name] = store;                                                                                   // 491
                                                                                                                  // 492
    var queued = self._updatesForUnknownStores[name];                                                             // 493
    if (queued) {                                                                                                 // 494
      store.beginUpdate(queued.length, false);                                                                    // 495
      _.each(queued, function (msg) {                                                                             // 496
        store.update(msg);                                                                                        // 497
      });                                                                                                         // 498
      store.endUpdate();                                                                                          // 499
      delete self._updatesForUnknownStores[name];                                                                 // 500
    }                                                                                                             // 501
                                                                                                                  // 502
    return true;                                                                                                  // 503
  },                                                                                                              // 504
                                                                                                                  // 505
  /**                                                                                                             // 506
   * @memberOf Meteor                                                                                             // 507
   * @importFromPackage meteor                                                                                    // 508
   * @summary Subscribe to a record set.  Returns a handle that provides                                          // 509
   * `stop()` and `ready()` methods.                                                                              // 510
   * @locus Client                                                                                                // 511
   * @param {String} name Name of the subscription.  Matches the name of the                                      // 512
   * server's `publish()` call.                                                                                   // 513
   * @param {EJSONable} [arg1,arg2...] Optional arguments passed to publisher                                     // 514
   * function on server.                                                                                          // 515
   * @param {Function|Object} [callbacks] Optional. May include `onStop`                                          // 516
   * and `onReady` callbacks. If there is an error, it is passed as an                                            // 517
   * argument to `onStop`. If a function is passed instead of an object, it                                       // 518
   * is interpreted as an `onReady` callback.                                                                     // 519
   */                                                                                                             // 520
  subscribe: function (name /* .. [arguments] .. (callback|callbacks) */) {                                       // 521
    var self = this;                                                                                              // 522
                                                                                                                  // 523
    var params = Array.prototype.slice.call(arguments, 1);                                                        // 524
    var callbacks = {};                                                                                           // 525
    if (params.length) {                                                                                          // 526
      var lastParam = params[params.length - 1];                                                                  // 527
      if (_.isFunction(lastParam)) {                                                                              // 528
        callbacks.onReady = params.pop();                                                                         // 529
      } else if (lastParam &&                                                                                     // 530
        // XXX COMPAT WITH 1.0.3.1 onError used to exist, but now we use                                          // 531
        // onStop with an error callback instead.                                                                 // 532
        _.any([lastParam.onReady, lastParam.onError, lastParam.onStop],                                           // 533
          _.isFunction)) {                                                                                        // 534
        callbacks = params.pop();                                                                                 // 535
      }                                                                                                           // 536
    }                                                                                                             // 537
                                                                                                                  // 538
    // Is there an existing sub with the same name and param, run in an                                           // 539
    // invalidated Computation? This will happen if we are rerunning an                                           // 540
    // existing computation.                                                                                      // 541
    //                                                                                                            // 542
    // For example, consider a rerun of:                                                                          // 543
    //                                                                                                            // 544
    //     Tracker.autorun(function () {                                                                          // 545
    //       Meteor.subscribe("foo", Session.get("foo"));                                                         // 546
    //       Meteor.subscribe("bar", Session.get("bar"));                                                         // 547
    //     });                                                                                                    // 548
    //                                                                                                            // 549
    // If "foo" has changed but "bar" has not, we will match the "bar"                                            // 550
    // subcribe to an existing inactive subscription in order to not                                              // 551
    // unsub and resub the subscription unnecessarily.                                                            // 552
    //                                                                                                            // 553
    // We only look for one such sub; if there are N apparently-identical subs                                    // 554
    // being invalidated, we will require N matching subscribe calls to keep                                      // 555
    // them all active.                                                                                           // 556
    var existing = _.find(self._subscriptions, function (sub) {                                                   // 557
      return sub.inactive && sub.name === name &&                                                                 // 558
        EJSON.equals(sub.params, params);                                                                         // 559
    });                                                                                                           // 560
                                                                                                                  // 561
    var id;                                                                                                       // 562
    if (existing) {                                                                                               // 563
      id = existing.id;                                                                                           // 564
      existing.inactive = false; // reactivate                                                                    // 565
                                                                                                                  // 566
      if (callbacks.onReady) {                                                                                    // 567
        // If the sub is not already ready, replace any ready callback with the                                   // 568
        // one provided now. (It's not really clear what users would expect for                                   // 569
        // an onReady callback inside an autorun; the semantics we provide is                                     // 570
        // that at the time the sub first becomes ready, we call the last                                         // 571
        // onReady callback provided, if any.)                                                                    // 572
        if (!existing.ready)                                                                                      // 573
          existing.readyCallback = callbacks.onReady;                                                             // 574
      }                                                                                                           // 575
                                                                                                                  // 576
      // XXX COMPAT WITH 1.0.3.1 we used to have onError but now we call                                          // 577
      // onStop with an optional error argument                                                                   // 578
      if (callbacks.onError) {                                                                                    // 579
        // Replace existing callback if any, so that errors aren't                                                // 580
        // double-reported.                                                                                       // 581
        existing.errorCallback = callbacks.onError;                                                               // 582
      }                                                                                                           // 583
                                                                                                                  // 584
      if (callbacks.onStop) {                                                                                     // 585
        existing.stopCallback = callbacks.onStop;                                                                 // 586
      }                                                                                                           // 587
    } else {                                                                                                      // 588
      // New sub! Generate an id, save it locally, and send message.                                              // 589
      id = Random.id();                                                                                           // 590
      self._subscriptions[id] = {                                                                                 // 591
        id: id,                                                                                                   // 592
        name: name,                                                                                               // 593
        params: EJSON.clone(params),                                                                              // 594
        inactive: false,                                                                                          // 595
        ready: false,                                                                                             // 596
        readyDeps: new Tracker.Dependency,                                                                        // 597
        readyCallback: callbacks.onReady,                                                                         // 598
        // XXX COMPAT WITH 1.0.3.1 #errorCallback                                                                 // 599
        errorCallback: callbacks.onError,                                                                         // 600
        stopCallback: callbacks.onStop,                                                                           // 601
        connection: self,                                                                                         // 602
        remove: function() {                                                                                      // 603
          delete this.connection._subscriptions[this.id];                                                         // 604
          this.ready && this.readyDeps.changed();                                                                 // 605
        },                                                                                                        // 606
        stop: function() {                                                                                        // 607
          this.connection._send({msg: 'unsub', id: id});                                                          // 608
          this.remove();                                                                                          // 609
                                                                                                                  // 610
          if (callbacks.onStop) {                                                                                 // 611
            callbacks.onStop();                                                                                   // 612
          }                                                                                                       // 613
        }                                                                                                         // 614
      };                                                                                                          // 615
      self._send({msg: 'sub', id: id, name: name, params: params});                                               // 616
    }                                                                                                             // 617
                                                                                                                  // 618
    // return a handle to the application.                                                                        // 619
    var handle = {                                                                                                // 620
      stop: function () {                                                                                         // 621
        if (!_.has(self._subscriptions, id))                                                                      // 622
          return;                                                                                                 // 623
                                                                                                                  // 624
        self._subscriptions[id].stop();                                                                           // 625
      },                                                                                                          // 626
      ready: function () {                                                                                        // 627
        // return false if we've unsubscribed.                                                                    // 628
        if (!_.has(self._subscriptions, id))                                                                      // 629
          return false;                                                                                           // 630
        var record = self._subscriptions[id];                                                                     // 631
        record.readyDeps.depend();                                                                                // 632
        return record.ready;                                                                                      // 633
      },                                                                                                          // 634
      subscriptionId: id                                                                                          // 635
    };                                                                                                            // 636
                                                                                                                  // 637
    if (Tracker.active) {                                                                                         // 638
      // We're in a reactive computation, so we'd like to unsubscribe when the                                    // 639
      // computation is invalidated... but not if the rerun just re-subscribes                                    // 640
      // to the same subscription!  When a rerun happens, we use onInvalidate                                     // 641
      // as a change to mark the subscription "inactive" so that it can                                           // 642
      // be reused from the rerun.  If it isn't reused, it's killed from                                          // 643
      // an afterFlush.                                                                                           // 644
      Tracker.onInvalidate(function (c) {                                                                         // 645
        if (_.has(self._subscriptions, id))                                                                       // 646
          self._subscriptions[id].inactive = true;                                                                // 647
                                                                                                                  // 648
        Tracker.afterFlush(function () {                                                                          // 649
          if (_.has(self._subscriptions, id) &&                                                                   // 650
              self._subscriptions[id].inactive)                                                                   // 651
            handle.stop();                                                                                        // 652
        });                                                                                                       // 653
      });                                                                                                         // 654
    }                                                                                                             // 655
                                                                                                                  // 656
    return handle;                                                                                                // 657
  },                                                                                                              // 658
                                                                                                                  // 659
  // options:                                                                                                     // 660
  // - onLateError {Function(error)} called if an error was received after the ready event.                       // 661
  //     (errors received before ready cause an error to be thrown)                                               // 662
  _subscribeAndWait: function (name, args, options) {                                                             // 663
    var self = this;                                                                                              // 664
    var f = new Future();                                                                                         // 665
    var ready = false;                                                                                            // 666
    var handle;                                                                                                   // 667
    args = args || [];                                                                                            // 668
    args.push({                                                                                                   // 669
      onReady: function () {                                                                                      // 670
        ready = true;                                                                                             // 671
        f['return']();                                                                                            // 672
      },                                                                                                          // 673
      onError: function (e) {                                                                                     // 674
        if (!ready)                                                                                               // 675
          f['throw'](e);                                                                                          // 676
        else                                                                                                      // 677
          options && options.onLateError && options.onLateError(e);                                               // 678
      }                                                                                                           // 679
    });                                                                                                           // 680
                                                                                                                  // 681
    handle = self.subscribe.apply(self, [name].concat(args));                                                     // 682
    f.wait();                                                                                                     // 683
    return handle;                                                                                                // 684
  },                                                                                                              // 685
                                                                                                                  // 686
  methods: function (methods) {                                                                                   // 687
    var self = this;                                                                                              // 688
    _.each(methods, function (func, name) {                                                                       // 689
      if (typeof func !== 'function')                                                                             // 690
        throw new Error("Method '" + name + "' must be a function");                                              // 691
      if (self._methodHandlers[name])                                                                             // 692
        throw new Error("A method named '" + name + "' is already defined");                                      // 693
      self._methodHandlers[name] = func;                                                                          // 694
    });                                                                                                           // 695
  },                                                                                                              // 696
                                                                                                                  // 697
  /**                                                                                                             // 698
   * @memberOf Meteor                                                                                             // 699
   * @importFromPackage meteor                                                                                    // 700
   * @summary Invokes a method passing any number of arguments.                                                   // 701
   * @locus Anywhere                                                                                              // 702
   * @param {String} name Name of method to invoke                                                                // 703
   * @param {EJSONable} [arg1,arg2...] Optional method arguments                                                  // 704
   * @param {Function} [asyncCallback] Optional callback, which is called asynchronously with the error or result after the method is complete. If not provided, the method runs synchronously if possible (see below).
   */                                                                                                             // 706
  call: function (name /* .. [arguments] .. callback */) {                                                        // 707
    // if it's a function, the last argument is the result callback,                                              // 708
    // not a parameter to the remote method.                                                                      // 709
    var args = Array.prototype.slice.call(arguments, 1);                                                          // 710
    if (args.length && typeof args[args.length - 1] === "function")                                               // 711
      var callback = args.pop();                                                                                  // 712
    return this.apply(name, args, callback);                                                                      // 713
  },                                                                                                              // 714
                                                                                                                  // 715
  // @param options {Optional Object}                                                                             // 716
  //   wait: Boolean - Should we wait to call this until all current methods                                      // 717
  //                   are fully finished, and block subsequent method calls                                      // 718
  //                   until this method is fully finished?                                                       // 719
  //                   (does not affect methods called from within this method)                                   // 720
  //   onResultReceived: Function - a callback to call as soon as the method                                      // 721
  //                                result is received. the data written by                                       // 722
  //                                the method may not yet be in the cache!                                       // 723
  //   returnStubValue: Boolean - If true then in cases where we would have                                       // 724
  //                              otherwise discarded the stub's return value                                     // 725
  //                              and returned undefined, instead we go ahead                                     // 726
  //                              and return it.  Specifically, this is any                                       // 727
  //                              time other than when (a) we are already                                         // 728
  //                              inside a stub or (b) we are in Node and no                                      // 729
  //                              callback was provided.  Currently we require                                    // 730
  //                              this flag to be explicitly passed to reduce                                     // 731
  //                              the likelihood that stub return values will                                     // 732
  //                              be confused with server return values; we                                       // 733
  //                              may improve this in future.                                                     // 734
  // @param callback {Optional Function}                                                                          // 735
                                                                                                                  // 736
  /**                                                                                                             // 737
   * @memberOf Meteor                                                                                             // 738
   * @importFromPackage meteor                                                                                    // 739
   * @summary Invoke a method passing an array of arguments.                                                      // 740
   * @locus Anywhere                                                                                              // 741
   * @param {String} name Name of method to invoke                                                                // 742
   * @param {EJSONable[]} args Method arguments                                                                   // 743
   * @param {Object} [options]                                                                                    // 744
   * @param {Boolean} options.wait (Client only) If true, don't send this method until all previous method calls have completed, and don't send any subsequent method calls until this one is completed.
   * @param {Function} options.onResultReceived (Client only) This callback is invoked with the error or result of the method (just like `asyncCallback`) as soon as the error or result is available. The local cache may not yet reflect the writes performed by the method.
   * @param (Boolean) options.noRetry (Client only) if true, don't send this method again on reload, simply call the callback an error with the error code 'invocation-failed'.
   * @param {Function} [asyncCallback] Optional callback; same semantics as in [`Meteor.call`](#meteor_call).     // 748
   */                                                                                                             // 749
  apply: function (name, args, options, callback) {                                                               // 750
    var self = this;                                                                                              // 751
                                                                                                                  // 752
    // We were passed 3 arguments. They may be either (name, args, options)                                       // 753
    // or (name, args, callback)                                                                                  // 754
    if (!callback && typeof options === 'function') {                                                             // 755
      callback = options;                                                                                         // 756
      options = {};                                                                                               // 757
    }                                                                                                             // 758
    options = options || {};                                                                                      // 759
                                                                                                                  // 760
    if (callback) {                                                                                               // 761
      // XXX would it be better form to do the binding in stream.on,                                              // 762
      // or caller, instead of here?                                                                              // 763
      // XXX improve error message (and how we report it)                                                         // 764
      callback = Meteor.bindEnvironment(                                                                          // 765
        callback,                                                                                                 // 766
        "delivering result of invoking '" + name + "'"                                                            // 767
      );                                                                                                          // 768
    }                                                                                                             // 769
                                                                                                                  // 770
    // Keep our args safe from mutation (eg if we don't send the message for a                                    // 771
    // while because of a wait method).                                                                           // 772
    args = EJSON.clone(args);                                                                                     // 773
                                                                                                                  // 774
    // Lazily allocate method ID once we know that it'll be needed.                                               // 775
    var methodId = (function () {                                                                                 // 776
      var id;                                                                                                     // 777
      return function () {                                                                                        // 778
        if (id === undefined)                                                                                     // 779
          id = '' + (self._nextMethodId++);                                                                       // 780
        return id;                                                                                                // 781
      };                                                                                                          // 782
    })();                                                                                                         // 783
                                                                                                                  // 784
    var enclosing = DDP._CurrentInvocation.get();                                                                 // 785
    var alreadyInSimulation = enclosing && enclosing.isSimulation;                                                // 786
                                                                                                                  // 787
    // Lazily generate a randomSeed, only if it is requested by the stub.                                         // 788
    // The random streams only have utility if they're used on both the client                                    // 789
    // and the server; if the client doesn't generate any 'random' values                                         // 790
    // then we don't expect the server to generate any either.                                                    // 791
    // Less commonly, the server may perform different actions from the client,                                   // 792
    // and may in fact generate values where the client did not, but we don't                                     // 793
    // have any client-side values to match, so even here we may as well just                                     // 794
    // use a random seed on the server.  In that case, we don't pass the                                          // 795
    // randomSeed to save bandwidth, and we don't even generate it to save a                                      // 796
    // bit of CPU and to avoid consuming entropy.                                                                 // 797
    var randomSeed = null;                                                                                        // 798
    var randomSeedGenerator = function () {                                                                       // 799
      if (randomSeed === null) {                                                                                  // 800
        randomSeed = DDPCommon.makeRpcSeed(enclosing, name);                                                      // 801
      }                                                                                                           // 802
      return randomSeed;                                                                                          // 803
    };                                                                                                            // 804
                                                                                                                  // 805
    // Run the stub, if we have one. The stub is supposed to make some                                            // 806
    // temporary writes to the database to give the user a smooth experience                                      // 807
    // until the actual result of executing the method comes back from the                                        // 808
    // server (whereupon the temporary writes to the database will be reversed                                    // 809
    // during the beginUpdate/endUpdate process.)                                                                 // 810
    //                                                                                                            // 811
    // Normally, we ignore the return value of the stub (even if it is an                                         // 812
    // exception), in favor of the real return value from the server. The                                         // 813
    // exception is if the *caller* is a stub. In that case, we're not going                                      // 814
    // to do a RPC, so we use the return value of the stub as our return                                          // 815
    // value.                                                                                                     // 816
                                                                                                                  // 817
    var stub = self._methodHandlers[name];                                                                        // 818
    if (stub) {                                                                                                   // 819
      var setUserId = function(userId) {                                                                          // 820
        self.setUserId(userId);                                                                                   // 821
      };                                                                                                          // 822
                                                                                                                  // 823
      var invocation = new DDPCommon.MethodInvocation({                                                           // 824
        isSimulation: true,                                                                                       // 825
        userId: self.userId(),                                                                                    // 826
        setUserId: setUserId,                                                                                     // 827
        randomSeed: function () { return randomSeedGenerator(); }                                                 // 828
      });                                                                                                         // 829
                                                                                                                  // 830
      if (!alreadyInSimulation)                                                                                   // 831
        self._saveOriginals();                                                                                    // 832
                                                                                                                  // 833
      try {                                                                                                       // 834
        // Note that unlike in the corresponding server code, we never audit                                      // 835
        // that stubs check() their arguments.                                                                    // 836
        var stubReturnValue = DDP._CurrentInvocation.withValue(invocation, function () {                          // 837
          if (Meteor.isServer) {                                                                                  // 838
            // Because saveOriginals and retrieveOriginals aren't reentrant,                                      // 839
            // don't allow stubs to yield.                                                                        // 840
            return Meteor._noYieldsAllowed(function () {                                                          // 841
              // re-clone, so that the stub can't affect our caller's values                                      // 842
              return stub.apply(invocation, EJSON.clone(args));                                                   // 843
            });                                                                                                   // 844
          } else {                                                                                                // 845
            return stub.apply(invocation, EJSON.clone(args));                                                     // 846
          }                                                                                                       // 847
        });                                                                                                       // 848
      }                                                                                                           // 849
      catch (e) {                                                                                                 // 850
        var exception = e;                                                                                        // 851
      }                                                                                                           // 852
                                                                                                                  // 853
      if (!alreadyInSimulation)                                                                                   // 854
        self._retrieveAndStoreOriginals(methodId());                                                              // 855
    }                                                                                                             // 856
                                                                                                                  // 857
    // If we're in a simulation, stop and return the result we have,                                              // 858
    // rather than going on to do an RPC. If there was no stub,                                                   // 859
    // we'll end up returning undefined.                                                                          // 860
    if (alreadyInSimulation) {                                                                                    // 861
      if (callback) {                                                                                             // 862
        callback(exception, stubReturnValue);                                                                     // 863
        return undefined;                                                                                         // 864
      }                                                                                                           // 865
      if (exception)                                                                                              // 866
        throw exception;                                                                                          // 867
      return stubReturnValue;                                                                                     // 868
    }                                                                                                             // 869
                                                                                                                  // 870
    // If an exception occurred in a stub, and we're ignoring it                                                  // 871
    // because we're doing an RPC and want to use what the server                                                 // 872
    // returns instead, log it so the developer knows                                                             // 873
    // (unless they explicitly ask to see the error).                                                             // 874
    //                                                                                                            // 875
    // Tests can set the 'expected' flag on an exception so it won't                                              // 876
    // go to log.                                                                                                 // 877
    if (exception) {                                                                                              // 878
      if (options.throwStubExceptions) {                                                                          // 879
        throw exception;                                                                                          // 880
      } else if (!exception.expected) {                                                                           // 881
        Meteor._debug("Exception while simulating the effect of invoking '" +                                     // 882
          name + "'", exception, exception.stack);                                                                // 883
      }                                                                                                           // 884
    }                                                                                                             // 885
                                                                                                                  // 886
                                                                                                                  // 887
    // At this point we're definitely doing an RPC, and we're going to                                            // 888
    // return the value of the RPC to the caller.                                                                 // 889
                                                                                                                  // 890
    // If the caller didn't give a callback, decide what to do.                                                   // 891
    if (!callback) {                                                                                              // 892
      if (Meteor.isClient) {                                                                                      // 893
        // On the client, we don't have fibers, so we can't block. The                                            // 894
        // only thing we can do is to return undefined and discard the                                            // 895
        // result of the RPC. If an error occurred then print the error                                           // 896
        // to the console.                                                                                        // 897
        callback = function (err) {                                                                               // 898
          err && Meteor._debug("Error invoking Method '" + name + "':",                                           // 899
                               err.message);                                                                      // 900
        };                                                                                                        // 901
      } else {                                                                                                    // 902
        // On the server, make the function synchronous. Throw on                                                 // 903
        // errors, return on success.                                                                             // 904
        var future = new Future;                                                                                  // 905
        callback = future.resolver();                                                                             // 906
      }                                                                                                           // 907
    }                                                                                                             // 908
    // Send the RPC. Note that on the client, it is important that the                                            // 909
    // stub have finished before we send the RPC, so that we know we have                                         // 910
    // a complete list of which local documents the stub wrote.                                                   // 911
    var message = {                                                                                               // 912
      msg: 'method',                                                                                              // 913
      method: name,                                                                                               // 914
      params: args,                                                                                               // 915
      id: methodId()                                                                                              // 916
    };                                                                                                            // 917
                                                                                                                  // 918
    // Send the randomSeed only if we used it                                                                     // 919
    if (randomSeed !== null) {                                                                                    // 920
      message.randomSeed = randomSeed;                                                                            // 921
    }                                                                                                             // 922
                                                                                                                  // 923
    var methodInvoker = new MethodInvoker({                                                                       // 924
      methodId: methodId(),                                                                                       // 925
      callback: callback,                                                                                         // 926
      connection: self,                                                                                           // 927
      onResultReceived: options.onResultReceived,                                                                 // 928
      wait: !!options.wait,                                                                                       // 929
      message: message,                                                                                           // 930
      noRetry: !!options.noRetry                                                                                  // 931
    });                                                                                                           // 932
                                                                                                                  // 933
    if (options.wait) {                                                                                           // 934
      // It's a wait method! Wait methods go in their own block.                                                  // 935
      self._outstandingMethodBlocks.push(                                                                         // 936
        {wait: true, methods: [methodInvoker]});                                                                  // 937
    } else {                                                                                                      // 938
      // Not a wait method. Start a new block if the previous block was a wait                                    // 939
      // block, and add it to the last block of methods.                                                          // 940
      if (_.isEmpty(self._outstandingMethodBlocks) ||                                                             // 941
          _.last(self._outstandingMethodBlocks).wait)                                                             // 942
        self._outstandingMethodBlocks.push({wait: false, methods: []});                                           // 943
      _.last(self._outstandingMethodBlocks).methods.push(methodInvoker);                                          // 944
    }                                                                                                             // 945
                                                                                                                  // 946
    // If we added it to the first block, send it out now.                                                        // 947
    if (self._outstandingMethodBlocks.length === 1)                                                               // 948
      methodInvoker.sendMessage();                                                                                // 949
                                                                                                                  // 950
    // If we're using the default callback on the server,                                                         // 951
    // block waiting for the result.                                                                              // 952
    if (future) {                                                                                                 // 953
      return future.wait();                                                                                       // 954
    }                                                                                                             // 955
    return options.returnStubValue ? stubReturnValue : undefined;                                                 // 956
  },                                                                                                              // 957
                                                                                                                  // 958
  // Before calling a method stub, prepare all stores to track changes and allow                                  // 959
  // _retrieveAndStoreOriginals to get the original versions of changed                                           // 960
  // documents.                                                                                                   // 961
  _saveOriginals: function () {                                                                                   // 962
    var self = this;                                                                                              // 963
    _.each(self._stores, function (s) {                                                                           // 964
      s.saveOriginals();                                                                                          // 965
    });                                                                                                           // 966
  },                                                                                                              // 967
  // Retrieves the original versions of all documents modified by the stub for                                    // 968
  // method 'methodId' from all stores and saves them to _serverDocuments (keyed                                  // 969
  // by document) and _documentsWrittenByStub (keyed by method ID).                                               // 970
  _retrieveAndStoreOriginals: function (methodId) {                                                               // 971
    var self = this;                                                                                              // 972
    if (self._documentsWrittenByStub[methodId])                                                                   // 973
      throw new Error("Duplicate methodId in _retrieveAndStoreOriginals");                                        // 974
                                                                                                                  // 975
    var docsWritten = [];                                                                                         // 976
    _.each(self._stores, function (s, collection) {                                                               // 977
      var originals = s.retrieveOriginals();                                                                      // 978
      // not all stores define retrieveOriginals                                                                  // 979
      if (!originals)                                                                                             // 980
        return;                                                                                                   // 981
      originals.forEach(function (doc, id) {                                                                      // 982
        docsWritten.push({collection: collection, id: id});                                                       // 983
        if (!_.has(self._serverDocuments, collection))                                                            // 984
          self._serverDocuments[collection] = new MongoIDMap;                                                     // 985
        var serverDoc = self._serverDocuments[collection].setDefault(id, {});                                     // 986
        if (serverDoc.writtenByStubs) {                                                                           // 987
          // We're not the first stub to write this doc. Just add our method ID                                   // 988
          // to the record.                                                                                       // 989
          serverDoc.writtenByStubs[methodId] = true;                                                              // 990
        } else {                                                                                                  // 991
          // First stub! Save the original value and our method ID.                                               // 992
          serverDoc.document = doc;                                                                               // 993
          serverDoc.flushCallbacks = [];                                                                          // 994
          serverDoc.writtenByStubs = {};                                                                          // 995
          serverDoc.writtenByStubs[methodId] = true;                                                              // 996
        }                                                                                                         // 997
      });                                                                                                         // 998
    });                                                                                                           // 999
    if (!_.isEmpty(docsWritten)) {                                                                                // 1000
      self._documentsWrittenByStub[methodId] = docsWritten;                                                       // 1001
    }                                                                                                             // 1002
  },                                                                                                              // 1003
                                                                                                                  // 1004
  // This is very much a private function we use to make the tests                                                // 1005
  // take up fewer server resources after they complete.                                                          // 1006
  _unsubscribeAll: function () {                                                                                  // 1007
    var self = this;                                                                                              // 1008
    _.each(_.clone(self._subscriptions), function (sub, id) {                                                     // 1009
      // Avoid killing the autoupdate subscription so that developers                                             // 1010
      // still get hot code pushes when writing tests.                                                            // 1011
      //                                                                                                          // 1012
      // XXX it's a hack to encode knowledge about autoupdate here,                                               // 1013
      // but it doesn't seem worth it yet to have a special API for                                               // 1014
      // subscriptions to preserve after unit tests.                                                              // 1015
      if (sub.name !== 'meteor_autoupdate_clientVersions') {                                                      // 1016
        self._subscriptions[id].stop();                                                                           // 1017
      }                                                                                                           // 1018
    });                                                                                                           // 1019
  },                                                                                                              // 1020
                                                                                                                  // 1021
  // Sends the DDP stringification of the given message object                                                    // 1022
  _send: function (obj) {                                                                                         // 1023
    var self = this;                                                                                              // 1024
    self._stream.send(DDPCommon.stringifyDDP(obj));                                                               // 1025
  },                                                                                                              // 1026
                                                                                                                  // 1027
  // We detected via DDP-level heartbeats that we've lost the                                                     // 1028
  // connection.  Unlike `disconnect` or `close`, a lost connection                                               // 1029
  // will be automatically retried.                                                                               // 1030
  _lostConnection: function (error) {                                                                             // 1031
    var self = this;                                                                                              // 1032
    self._stream._lostConnection(error);                                                                          // 1033
  },                                                                                                              // 1034
                                                                                                                  // 1035
  /**                                                                                                             // 1036
   * @summary Get the current connection status. A reactive data source.                                          // 1037
   * @locus Client                                                                                                // 1038
   * @memberOf Meteor                                                                                             // 1039
   * @importFromPackage meteor                                                                                    // 1040
   */                                                                                                             // 1041
  status: function (/*passthrough args*/) {                                                                       // 1042
    var self = this;                                                                                              // 1043
    return self._stream.status.apply(self._stream, arguments);                                                    // 1044
  },                                                                                                              // 1045
                                                                                                                  // 1046
  /**                                                                                                             // 1047
   * @summary Force an immediate reconnection attempt if the client is not connected to the server.               // 1048
                                                                                                                  // 1049
  This method does nothing if the client is already connected.                                                    // 1050
   * @locus Client                                                                                                // 1051
   * @memberOf Meteor                                                                                             // 1052
   * @importFromPackage meteor                                                                                    // 1053
   */                                                                                                             // 1054
  reconnect: function (/*passthrough args*/) {                                                                    // 1055
    var self = this;                                                                                              // 1056
    return self._stream.reconnect.apply(self._stream, arguments);                                                 // 1057
  },                                                                                                              // 1058
                                                                                                                  // 1059
  /**                                                                                                             // 1060
   * @summary Disconnect the client from the server.                                                              // 1061
   * @locus Client                                                                                                // 1062
   * @memberOf Meteor                                                                                             // 1063
   * @importFromPackage meteor                                                                                    // 1064
   */                                                                                                             // 1065
  disconnect: function (/*passthrough args*/) {                                                                   // 1066
    var self = this;                                                                                              // 1067
    return self._stream.disconnect.apply(self._stream, arguments);                                                // 1068
  },                                                                                                              // 1069
                                                                                                                  // 1070
  close: function () {                                                                                            // 1071
    var self = this;                                                                                              // 1072
    return self._stream.disconnect({_permanent: true});                                                           // 1073
  },                                                                                                              // 1074
                                                                                                                  // 1075
  ///                                                                                                             // 1076
  /// Reactive user system                                                                                        // 1077
  ///                                                                                                             // 1078
  userId: function () {                                                                                           // 1079
    var self = this;                                                                                              // 1080
    if (self._userIdDeps)                                                                                         // 1081
      self._userIdDeps.depend();                                                                                  // 1082
    return self._userId;                                                                                          // 1083
  },                                                                                                              // 1084
                                                                                                                  // 1085
  setUserId: function (userId) {                                                                                  // 1086
    var self = this;                                                                                              // 1087
    // Avoid invalidating dependents if setUserId is called with current value.                                   // 1088
    if (self._userId === userId)                                                                                  // 1089
      return;                                                                                                     // 1090
    self._userId = userId;                                                                                        // 1091
    if (self._userIdDeps)                                                                                         // 1092
      self._userIdDeps.changed();                                                                                 // 1093
  },                                                                                                              // 1094
                                                                                                                  // 1095
  // Returns true if we are in a state after reconnect of waiting for subs to be                                  // 1096
  // revived or early methods to finish their data, or we are waiting for a                                       // 1097
  // "wait" method to finish.                                                                                     // 1098
  _waitingForQuiescence: function () {                                                                            // 1099
    var self = this;                                                                                              // 1100
    return (! _.isEmpty(self._subsBeingRevived) ||                                                                // 1101
            ! _.isEmpty(self._methodsBlockingQuiescence));                                                        // 1102
  },                                                                                                              // 1103
                                                                                                                  // 1104
  // Returns true if any method whose message has been sent to the server has                                     // 1105
  // not yet invoked its user callback.                                                                           // 1106
  _anyMethodsAreOutstanding: function () {                                                                        // 1107
    var self = this;                                                                                              // 1108
    return _.any(_.pluck(self._methodInvokers, 'sentMessage'));                                                   // 1109
  },                                                                                                              // 1110
                                                                                                                  // 1111
  _livedata_connected: function (msg) {                                                                           // 1112
    var self = this;                                                                                              // 1113
                                                                                                                  // 1114
    if (self._version !== 'pre1' && self._heartbeatInterval !== 0) {                                              // 1115
      self._heartbeat = new DDPCommon.Heartbeat({                                                                 // 1116
        heartbeatInterval: self._heartbeatInterval,                                                               // 1117
        heartbeatTimeout: self._heartbeatTimeout,                                                                 // 1118
        onTimeout: function () {                                                                                  // 1119
          self._lostConnection(                                                                                   // 1120
            new DDP.ConnectionError("DDP heartbeat timed out"));                                                  // 1121
        },                                                                                                        // 1122
        sendPing: function () {                                                                                   // 1123
          self._send({msg: 'ping'});                                                                              // 1124
        }                                                                                                         // 1125
      });                                                                                                         // 1126
      self._heartbeat.start();                                                                                    // 1127
    }                                                                                                             // 1128
                                                                                                                  // 1129
    // If this is a reconnect, we'll have to reset all stores.                                                    // 1130
    if (self._lastSessionId)                                                                                      // 1131
      self._resetStores = true;                                                                                   // 1132
                                                                                                                  // 1133
    if (typeof (msg.session) === "string") {                                                                      // 1134
      var reconnectedToPreviousSession = (self._lastSessionId === msg.session);                                   // 1135
      self._lastSessionId = msg.session;                                                                          // 1136
    }                                                                                                             // 1137
                                                                                                                  // 1138
    if (reconnectedToPreviousSession) {                                                                           // 1139
      // Successful reconnection -- pick up where we left off.  Note that right                                   // 1140
      // now, this never happens: the server never connects us to a previous                                      // 1141
      // session, because DDP doesn't provide enough data for the server to know                                  // 1142
      // what messages the client has processed. We need to improve DDP to make                                   // 1143
      // this possible, at which point we'll probably need more code here.                                        // 1144
      return;                                                                                                     // 1145
    }                                                                                                             // 1146
                                                                                                                  // 1147
    // Server doesn't have our data any more. Re-sync a new session.                                              // 1148
                                                                                                                  // 1149
    // Forget about messages we were buffering for unknown collections. They'll                                   // 1150
    // be resent if still relevant.                                                                               // 1151
    self._updatesForUnknownStores = {};                                                                           // 1152
                                                                                                                  // 1153
    if (self._resetStores) {                                                                                      // 1154
      // Forget about the effects of stubs. We'll be resetting all collections                                    // 1155
      // anyway.                                                                                                  // 1156
      self._documentsWrittenByStub = {};                                                                          // 1157
      self._serverDocuments = {};                                                                                 // 1158
    }                                                                                                             // 1159
                                                                                                                  // 1160
    // Clear _afterUpdateCallbacks.                                                                               // 1161
    self._afterUpdateCallbacks = [];                                                                              // 1162
                                                                                                                  // 1163
    // Mark all named subscriptions which are ready (ie, we already called the                                    // 1164
    // ready callback) as needing to be revived.                                                                  // 1165
    // XXX We should also block reconnect quiescence until unnamed subscriptions                                  // 1166
    //     (eg, autopublish) are done re-publishing to avoid flicker!                                             // 1167
    self._subsBeingRevived = {};                                                                                  // 1168
    _.each(self._subscriptions, function (sub, id) {                                                              // 1169
      if (sub.ready)                                                                                              // 1170
        self._subsBeingRevived[id] = true;                                                                        // 1171
    });                                                                                                           // 1172
                                                                                                                  // 1173
    // Arrange for "half-finished" methods to have their callbacks run, and                                       // 1174
    // track methods that were sent on this connection so that we don't                                           // 1175
    // quiesce until they are all done.                                                                           // 1176
    //                                                                                                            // 1177
    // Start by clearing _methodsBlockingQuiescence: methods sent before                                          // 1178
    // reconnect don't matter, and any "wait" methods sent on the new connection                                  // 1179
    // that we drop here will be restored by the loop below.                                                      // 1180
    self._methodsBlockingQuiescence = {};                                                                         // 1181
    if (self._resetStores) {                                                                                      // 1182
      _.each(self._methodInvokers, function (invoker) {                                                           // 1183
        if (invoker.gotResult()) {                                                                                // 1184
          // This method already got its result, but it didn't call its callback                                  // 1185
          // because its data didn't become visible. We did not resend the                                        // 1186
          // method RPC. We'll call its callback when we get a full quiesce,                                      // 1187
          // since that's as close as we'll get to "data must be visible".                                        // 1188
          self._afterUpdateCallbacks.push(_.bind(invoker.dataVisible, invoker));                                  // 1189
        } else if (invoker.sentMessage) {                                                                         // 1190
          // This method has been sent on this connection (maybe as a resend                                      // 1191
          // from the last connection, maybe from onReconnect, maybe just very                                    // 1192
          // quickly before processing the connected message).                                                    // 1193
          //                                                                                                      // 1194
          // We don't need to do anything special to ensure its callbacks get                                     // 1195
          // called, but we'll count it as a method which is preventing                                           // 1196
          // reconnect quiescence. (eg, it might be a login method that was run                                   // 1197
          // from onReconnect, and we don't want to see flicker by seeing a                                       // 1198
          // logged-out state.)                                                                                   // 1199
          self._methodsBlockingQuiescence[invoker.methodId] = true;                                               // 1200
        }                                                                                                         // 1201
      });                                                                                                         // 1202
    }                                                                                                             // 1203
                                                                                                                  // 1204
    self._messagesBufferedUntilQuiescence = [];                                                                   // 1205
                                                                                                                  // 1206
    // If we're not waiting on any methods or subs, we can reset the stores and                                   // 1207
    // call the callbacks immediately.                                                                            // 1208
    if (!self._waitingForQuiescence()) {                                                                          // 1209
      if (self._resetStores) {                                                                                    // 1210
        _.each(self._stores, function (s) {                                                                       // 1211
          s.beginUpdate(0, true);                                                                                 // 1212
          s.endUpdate();                                                                                          // 1213
        });                                                                                                       // 1214
        self._resetStores = false;                                                                                // 1215
      }                                                                                                           // 1216
      self._runAfterUpdateCallbacks();                                                                            // 1217
    }                                                                                                             // 1218
  },                                                                                                              // 1219
                                                                                                                  // 1220
                                                                                                                  // 1221
  _processOneDataMessage: function (msg, updates) {                                                               // 1222
    var self = this;                                                                                              // 1223
    // Using underscore here so as not to need to capitalize.                                                     // 1224
    self['_process_' + msg.msg](msg, updates);                                                                    // 1225
  },                                                                                                              // 1226
                                                                                                                  // 1227
                                                                                                                  // 1228
  _livedata_data: function (msg) {                                                                                // 1229
    var self = this;                                                                                              // 1230
                                                                                                                  // 1231
    // collection name -> array of messages                                                                       // 1232
    var updates = {};                                                                                             // 1233
                                                                                                                  // 1234
    if (self._waitingForQuiescence()) {                                                                           // 1235
      self._messagesBufferedUntilQuiescence.push(msg);                                                            // 1236
                                                                                                                  // 1237
      if (msg.msg === "nosub")                                                                                    // 1238
        delete self._subsBeingRevived[msg.id];                                                                    // 1239
                                                                                                                  // 1240
      _.each(msg.subs || [], function (subId) {                                                                   // 1241
        delete self._subsBeingRevived[subId];                                                                     // 1242
      });                                                                                                         // 1243
      _.each(msg.methods || [], function (methodId) {                                                             // 1244
        delete self._methodsBlockingQuiescence[methodId];                                                         // 1245
      });                                                                                                         // 1246
                                                                                                                  // 1247
      if (self._waitingForQuiescence())                                                                           // 1248
        return;                                                                                                   // 1249
                                                                                                                  // 1250
      // No methods or subs are blocking quiescence!                                                              // 1251
      // We'll now process and all of our buffered messages, reset all stores,                                    // 1252
      // and apply them all at once.                                                                              // 1253
      _.each(self._messagesBufferedUntilQuiescence, function (bufferedMsg) {                                      // 1254
        self._processOneDataMessage(bufferedMsg, updates);                                                        // 1255
      });                                                                                                         // 1256
      self._messagesBufferedUntilQuiescence = [];                                                                 // 1257
    } else {                                                                                                      // 1258
      self._processOneDataMessage(msg, updates);                                                                  // 1259
    }                                                                                                             // 1260
                                                                                                                  // 1261
    if (self._resetStores || !_.isEmpty(updates)) {                                                               // 1262
      // Begin a transactional update of each store.                                                              // 1263
      _.each(self._stores, function (s, storeName) {                                                              // 1264
        s.beginUpdate(_.has(updates, storeName) ? updates[storeName].length : 0,                                  // 1265
                      self._resetStores);                                                                         // 1266
      });                                                                                                         // 1267
      self._resetStores = false;                                                                                  // 1268
                                                                                                                  // 1269
      _.each(updates, function (updateMessages, storeName) {                                                      // 1270
        var store = self._stores[storeName];                                                                      // 1271
        if (store) {                                                                                              // 1272
          _.each(updateMessages, function (updateMessage) {                                                       // 1273
            store.update(updateMessage);                                                                          // 1274
          });                                                                                                     // 1275
        } else {                                                                                                  // 1276
          // Nobody's listening for this data. Queue it up until                                                  // 1277
          // someone wants it.                                                                                    // 1278
          // XXX memory use will grow without bound if you forget to                                              // 1279
          // create a collection or just don't care about it... going                                             // 1280
          // to have to do something about that.                                                                  // 1281
          if (!_.has(self._updatesForUnknownStores, storeName))                                                   // 1282
            self._updatesForUnknownStores[storeName] = [];                                                        // 1283
          Array.prototype.push.apply(self._updatesForUnknownStores[storeName],                                    // 1284
                                     updateMessages);                                                             // 1285
        }                                                                                                         // 1286
      });                                                                                                         // 1287
                                                                                                                  // 1288
      // End update transaction.                                                                                  // 1289
      _.each(self._stores, function (s) { s.endUpdate(); });                                                      // 1290
    }                                                                                                             // 1291
                                                                                                                  // 1292
    self._runAfterUpdateCallbacks();                                                                              // 1293
  },                                                                                                              // 1294
                                                                                                                  // 1295
  // Call any callbacks deferred with _runWhenAllServerDocsAreFlushed whose                                       // 1296
  // relevant docs have been flushed, as well as dataVisible callbacks at                                         // 1297
  // reconnect-quiescence time.                                                                                   // 1298
  _runAfterUpdateCallbacks: function () {                                                                         // 1299
    var self = this;                                                                                              // 1300
    var callbacks = self._afterUpdateCallbacks;                                                                   // 1301
    self._afterUpdateCallbacks = [];                                                                              // 1302
    _.each(callbacks, function (c) {                                                                              // 1303
      c();                                                                                                        // 1304
    });                                                                                                           // 1305
  },                                                                                                              // 1306
                                                                                                                  // 1307
  _pushUpdate: function (updates, collection, msg) {                                                              // 1308
    var self = this;                                                                                              // 1309
    if (!_.has(updates, collection)) {                                                                            // 1310
      updates[collection] = [];                                                                                   // 1311
    }                                                                                                             // 1312
    updates[collection].push(msg);                                                                                // 1313
  },                                                                                                              // 1314
                                                                                                                  // 1315
  _getServerDoc: function (collection, id) {                                                                      // 1316
    var self = this;                                                                                              // 1317
    if (!_.has(self._serverDocuments, collection))                                                                // 1318
      return null;                                                                                                // 1319
    var serverDocsForCollection = self._serverDocuments[collection];                                              // 1320
    return serverDocsForCollection.get(id) || null;                                                               // 1321
  },                                                                                                              // 1322
                                                                                                                  // 1323
  _process_added: function (msg, updates) {                                                                       // 1324
    var self = this;                                                                                              // 1325
    var id = MongoID.idParse(msg.id);                                                                             // 1326
    var serverDoc = self._getServerDoc(msg.collection, id);                                                       // 1327
    if (serverDoc) {                                                                                              // 1328
      // Some outstanding stub wrote here.                                                                        // 1329
      var isExisting = (serverDoc.document !== undefined);                                                        // 1330
                                                                                                                  // 1331
      serverDoc.document = msg.fields || {};                                                                      // 1332
      serverDoc.document._id = id;                                                                                // 1333
                                                                                                                  // 1334
      if (self._resetStores) {                                                                                    // 1335
        // During reconnect the server is sending adds for existing ids.                                          // 1336
        // Always push an update so that document stays in the store after                                        // 1337
        // reset. Use current version of the document for this update, so                                         // 1338
        // that stub-written values are preserved.                                                                // 1339
        var currentDoc = self._stores[msg.collection].getDoc(msg.id);                                             // 1340
        if (currentDoc !== undefined)                                                                             // 1341
          msg.fields = currentDoc;                                                                                // 1342
                                                                                                                  // 1343
        self._pushUpdate(updates, msg.collection, msg);                                                           // 1344
      } else if (isExisting) {                                                                                    // 1345
        throw new Error("Server sent add for existing id: " + msg.id);                                            // 1346
      }                                                                                                           // 1347
    } else {                                                                                                      // 1348
      self._pushUpdate(updates, msg.collection, msg);                                                             // 1349
    }                                                                                                             // 1350
  },                                                                                                              // 1351
                                                                                                                  // 1352
  _process_changed: function (msg, updates) {                                                                     // 1353
    var self = this;                                                                                              // 1354
    var serverDoc = self._getServerDoc(                                                                           // 1355
      msg.collection, MongoID.idParse(msg.id));                                                                   // 1356
    if (serverDoc) {                                                                                              // 1357
      if (serverDoc.document === undefined)                                                                       // 1358
        throw new Error("Server sent changed for nonexisting id: " + msg.id);                                     // 1359
      DiffSequence.applyChanges(serverDoc.document, msg.fields);                                                  // 1360
    } else {                                                                                                      // 1361
      self._pushUpdate(updates, msg.collection, msg);                                                             // 1362
    }                                                                                                             // 1363
  },                                                                                                              // 1364
                                                                                                                  // 1365
  _process_removed: function (msg, updates) {                                                                     // 1366
    var self = this;                                                                                              // 1367
    var serverDoc = self._getServerDoc(                                                                           // 1368
      msg.collection, MongoID.idParse(msg.id));                                                                   // 1369
    if (serverDoc) {                                                                                              // 1370
      // Some outstanding stub wrote here.                                                                        // 1371
      if (serverDoc.document === undefined)                                                                       // 1372
        throw new Error("Server sent removed for nonexisting id:" + msg.id);                                      // 1373
      serverDoc.document = undefined;                                                                             // 1374
    } else {                                                                                                      // 1375
      self._pushUpdate(updates, msg.collection, {                                                                 // 1376
        msg: 'removed',                                                                                           // 1377
        collection: msg.collection,                                                                               // 1378
        id: msg.id                                                                                                // 1379
      });                                                                                                         // 1380
    }                                                                                                             // 1381
  },                                                                                                              // 1382
                                                                                                                  // 1383
  _process_updated: function (msg, updates) {                                                                     // 1384
    var self = this;                                                                                              // 1385
    // Process "method done" messages.                                                                            // 1386
    _.each(msg.methods, function (methodId) {                                                                     // 1387
      _.each(self._documentsWrittenByStub[methodId], function (written) {                                         // 1388
        var serverDoc = self._getServerDoc(written.collection, written.id);                                       // 1389
        if (!serverDoc)                                                                                           // 1390
          throw new Error("Lost serverDoc for " + JSON.stringify(written));                                       // 1391
        if (!serverDoc.writtenByStubs[methodId])                                                                  // 1392
          throw new Error("Doc " + JSON.stringify(written) +                                                      // 1393
                          " not written by  method " + methodId);                                                 // 1394
        delete serverDoc.writtenByStubs[methodId];                                                                // 1395
        if (_.isEmpty(serverDoc.writtenByStubs)) {                                                                // 1396
          // All methods whose stubs wrote this method have completed! We can                                     // 1397
          // now copy the saved document to the database (reverting the stub's                                    // 1398
          // change if the server did not write to this object, or applying the                                   // 1399
          // server's writes if it did).                                                                          // 1400
                                                                                                                  // 1401
          // This is a fake ddp 'replace' message.  It's just for talking                                         // 1402
          // between livedata connections and minimongo.  (We have to stringify                                   // 1403
          // the ID because it's supposed to look like a wire message.)                                           // 1404
          self._pushUpdate(updates, written.collection, {                                                         // 1405
            msg: 'replace',                                                                                       // 1406
            id: MongoID.idStringify(written.id),                                                                  // 1407
            replace: serverDoc.document                                                                           // 1408
          });                                                                                                     // 1409
          // Call all flush callbacks.                                                                            // 1410
          _.each(serverDoc.flushCallbacks, function (c) {                                                         // 1411
            c();                                                                                                  // 1412
          });                                                                                                     // 1413
                                                                                                                  // 1414
          // Delete this completed serverDocument. Don't bother to GC empty                                       // 1415
          // IdMaps inside self._serverDocuments, since there probably aren't                                     // 1416
          // many collections and they'll be written repeatedly.                                                  // 1417
          self._serverDocuments[written.collection].remove(written.id);                                           // 1418
        }                                                                                                         // 1419
      });                                                                                                         // 1420
      delete self._documentsWrittenByStub[methodId];                                                              // 1421
                                                                                                                  // 1422
      // We want to call the data-written callback, but we can't do so until all                                  // 1423
      // currently buffered messages are flushed.                                                                 // 1424
      var callbackInvoker = self._methodInvokers[methodId];                                                       // 1425
      if (!callbackInvoker)                                                                                       // 1426
        throw new Error("No callback invoker for method " + methodId);                                            // 1427
      self._runWhenAllServerDocsAreFlushed(                                                                       // 1428
        _.bind(callbackInvoker.dataVisible, callbackInvoker));                                                    // 1429
    });                                                                                                           // 1430
  },                                                                                                              // 1431
                                                                                                                  // 1432
  _process_ready: function (msg, updates) {                                                                       // 1433
    var self = this;                                                                                              // 1434
    // Process "sub ready" messages. "sub ready" messages don't take effect                                       // 1435
    // until all current server documents have been flushed to the local                                          // 1436
    // database. We can use a write fence to implement this.                                                      // 1437
    _.each(msg.subs, function (subId) {                                                                           // 1438
      self._runWhenAllServerDocsAreFlushed(function () {                                                          // 1439
        var subRecord = self._subscriptions[subId];                                                               // 1440
        // Did we already unsubscribe?                                                                            // 1441
        if (!subRecord)                                                                                           // 1442
          return;                                                                                                 // 1443
        // Did we already receive a ready message? (Oops!)                                                        // 1444
        if (subRecord.ready)                                                                                      // 1445
          return;                                                                                                 // 1446
        subRecord.ready = true;                                                                                   // 1447
        subRecord.readyCallback && subRecord.readyCallback();                                                     // 1448
        subRecord.readyDeps.changed();                                                                            // 1449
      });                                                                                                         // 1450
    });                                                                                                           // 1451
  },                                                                                                              // 1452
                                                                                                                  // 1453
  // Ensures that "f" will be called after all documents currently in                                             // 1454
  // _serverDocuments have been written to the local cache. f will not be called                                  // 1455
  // if the connection is lost before then!                                                                       // 1456
  _runWhenAllServerDocsAreFlushed: function (f) {                                                                 // 1457
    var self = this;                                                                                              // 1458
    var runFAfterUpdates = function () {                                                                          // 1459
      self._afterUpdateCallbacks.push(f);                                                                         // 1460
    };                                                                                                            // 1461
    var unflushedServerDocCount = 0;                                                                              // 1462
    var onServerDocFlush = function () {                                                                          // 1463
      --unflushedServerDocCount;                                                                                  // 1464
      if (unflushedServerDocCount === 0) {                                                                        // 1465
        // This was the last doc to flush! Arrange to run f after the updates                                     // 1466
        // have been applied.                                                                                     // 1467
        runFAfterUpdates();                                                                                       // 1468
      }                                                                                                           // 1469
    };                                                                                                            // 1470
    _.each(self._serverDocuments, function (collectionDocs) {                                                     // 1471
      collectionDocs.forEach(function (serverDoc) {                                                               // 1472
        var writtenByStubForAMethodWithSentMessage = _.any(                                                       // 1473
          serverDoc.writtenByStubs, function (dummy, methodId) {                                                  // 1474
            var invoker = self._methodInvokers[methodId];                                                         // 1475
            return invoker && invoker.sentMessage;                                                                // 1476
          });                                                                                                     // 1477
        if (writtenByStubForAMethodWithSentMessage) {                                                             // 1478
          ++unflushedServerDocCount;                                                                              // 1479
          serverDoc.flushCallbacks.push(onServerDocFlush);                                                        // 1480
        }                                                                                                         // 1481
      });                                                                                                         // 1482
    });                                                                                                           // 1483
    if (unflushedServerDocCount === 0) {                                                                          // 1484
      // There aren't any buffered docs --- we can call f as soon as the current                                  // 1485
      // round of updates is applied!                                                                             // 1486
      runFAfterUpdates();                                                                                         // 1487
    }                                                                                                             // 1488
  },                                                                                                              // 1489
                                                                                                                  // 1490
  _livedata_nosub: function (msg) {                                                                               // 1491
    var self = this;                                                                                              // 1492
                                                                                                                  // 1493
    // First pass it through _livedata_data, which only uses it to help get                                       // 1494
    // towards quiescence.                                                                                        // 1495
    self._livedata_data(msg);                                                                                     // 1496
                                                                                                                  // 1497
    // Do the rest of our processing immediately, with no                                                         // 1498
    // buffering-until-quiescence.                                                                                // 1499
                                                                                                                  // 1500
    // we weren't subbed anyway, or we initiated the unsub.                                                       // 1501
    if (!_.has(self._subscriptions, msg.id))                                                                      // 1502
      return;                                                                                                     // 1503
                                                                                                                  // 1504
    // XXX COMPAT WITH 1.0.3.1 #errorCallback                                                                     // 1505
    var errorCallback = self._subscriptions[msg.id].errorCallback;                                                // 1506
    var stopCallback = self._subscriptions[msg.id].stopCallback;                                                  // 1507
                                                                                                                  // 1508
    self._subscriptions[msg.id].remove();                                                                         // 1509
                                                                                                                  // 1510
    var meteorErrorFromMsg = function (msgArg) {                                                                  // 1511
      return msgArg && msgArg.error && new Meteor.Error(                                                          // 1512
        msgArg.error.error, msgArg.error.reason, msgArg.error.details);                                           // 1513
    }                                                                                                             // 1514
                                                                                                                  // 1515
    // XXX COMPAT WITH 1.0.3.1 #errorCallback                                                                     // 1516
    if (errorCallback && msg.error) {                                                                             // 1517
      errorCallback(meteorErrorFromMsg(msg));                                                                     // 1518
    }                                                                                                             // 1519
                                                                                                                  // 1520
    if (stopCallback) {                                                                                           // 1521
      stopCallback(meteorErrorFromMsg(msg));                                                                      // 1522
    }                                                                                                             // 1523
  },                                                                                                              // 1524
                                                                                                                  // 1525
  _process_nosub: function () {                                                                                   // 1526
    // This is called as part of the "buffer until quiescence" process, but                                       // 1527
    // nosub's effect is always immediate. It only goes in the buffer at all                                      // 1528
    // because it's possible for a nosub to be the thing that triggers                                            // 1529
    // quiescence, if we were waiting for a sub to be revived and it dies                                         // 1530
    // instead.                                                                                                   // 1531
  },                                                                                                              // 1532
                                                                                                                  // 1533
  _livedata_result: function (msg) {                                                                              // 1534
    // id, result or error. error has error (code), reason, details                                               // 1535
                                                                                                                  // 1536
    var self = this;                                                                                              // 1537
                                                                                                                  // 1538
    // find the outstanding request                                                                               // 1539
    // should be O(1) in nearly all realistic use cases                                                           // 1540
    if (_.isEmpty(self._outstandingMethodBlocks)) {                                                               // 1541
      Meteor._debug("Received method result but no methods outstanding");                                         // 1542
      return;                                                                                                     // 1543
    }                                                                                                             // 1544
    var currentMethodBlock = self._outstandingMethodBlocks[0].methods;                                            // 1545
    var m;                                                                                                        // 1546
    for (var i = 0; i < currentMethodBlock.length; i++) {                                                         // 1547
      m = currentMethodBlock[i];                                                                                  // 1548
      if (m.methodId === msg.id)                                                                                  // 1549
        break;                                                                                                    // 1550
    }                                                                                                             // 1551
                                                                                                                  // 1552
    if (!m) {                                                                                                     // 1553
      Meteor._debug("Can't match method response to original method call", msg);                                  // 1554
      return;                                                                                                     // 1555
    }                                                                                                             // 1556
                                                                                                                  // 1557
    // Remove from current method block. This may leave the block empty, but we                                   // 1558
    // don't move on to the next block until the callback has been delivered, in                                  // 1559
    // _outstandingMethodFinished.                                                                                // 1560
    currentMethodBlock.splice(i, 1);                                                                              // 1561
                                                                                                                  // 1562
    if (_.has(msg, 'error')) {                                                                                    // 1563
      m.receiveResult(new Meteor.Error(                                                                           // 1564
        msg.error.error, msg.error.reason,                                                                        // 1565
        msg.error.details));                                                                                      // 1566
    } else {                                                                                                      // 1567
      // msg.result may be undefined if the method didn't return a                                                // 1568
      // value                                                                                                    // 1569
      m.receiveResult(undefined, msg.result);                                                                     // 1570
    }                                                                                                             // 1571
  },                                                                                                              // 1572
                                                                                                                  // 1573
  // Called by MethodInvoker after a method's callback is invoked.  If this was                                   // 1574
  // the last outstanding method in the current block, runs the next block. If                                    // 1575
  // there are no more methods, consider accepting a hot code push.                                               // 1576
  _outstandingMethodFinished: function () {                                                                       // 1577
    var self = this;                                                                                              // 1578
    if (self._anyMethodsAreOutstanding())                                                                         // 1579
      return;                                                                                                     // 1580
                                                                                                                  // 1581
    // No methods are outstanding. This should mean that the first block of                                       // 1582
    // methods is empty. (Or it might not exist, if this was a method that                                        // 1583
    // half-finished before disconnect/reconnect.)                                                                // 1584
    if (! _.isEmpty(self._outstandingMethodBlocks)) {                                                             // 1585
      var firstBlock = self._outstandingMethodBlocks.shift();                                                     // 1586
      if (! _.isEmpty(firstBlock.methods))                                                                        // 1587
        throw new Error("No methods outstanding but nonempty block: " +                                           // 1588
                        JSON.stringify(firstBlock));                                                              // 1589
                                                                                                                  // 1590
      // Send the outstanding methods now in the first block.                                                     // 1591
      if (!_.isEmpty(self._outstandingMethodBlocks))                                                              // 1592
        self._sendOutstandingMethods();                                                                           // 1593
    }                                                                                                             // 1594
                                                                                                                  // 1595
    // Maybe accept a hot code push.                                                                              // 1596
    self._maybeMigrate();                                                                                         // 1597
  },                                                                                                              // 1598
                                                                                                                  // 1599
  // Sends messages for all the methods in the first block in                                                     // 1600
  // _outstandingMethodBlocks.                                                                                    // 1601
  _sendOutstandingMethods: function() {                                                                           // 1602
    var self = this;                                                                                              // 1603
    if (_.isEmpty(self._outstandingMethodBlocks))                                                                 // 1604
      return;                                                                                                     // 1605
    _.each(self._outstandingMethodBlocks[0].methods, function (m) {                                               // 1606
      m.sendMessage();                                                                                            // 1607
    });                                                                                                           // 1608
  },                                                                                                              // 1609
                                                                                                                  // 1610
  _livedata_error: function (msg) {                                                                               // 1611
    Meteor._debug("Received error from server: ", msg.reason);                                                    // 1612
    if (msg.offendingMessage)                                                                                     // 1613
      Meteor._debug("For: ", msg.offendingMessage);                                                               // 1614
  },                                                                                                              // 1615
                                                                                                                  // 1616
  _callOnReconnectAndSendAppropriateOutstandingMethods: function() {                                              // 1617
    var self = this;                                                                                              // 1618
    var oldOutstandingMethodBlocks = self._outstandingMethodBlocks;                                               // 1619
    self._outstandingMethodBlocks = [];                                                                           // 1620
                                                                                                                  // 1621
    self.onReconnect();                                                                                           // 1622
                                                                                                                  // 1623
    if (_.isEmpty(oldOutstandingMethodBlocks))                                                                    // 1624
      return;                                                                                                     // 1625
                                                                                                                  // 1626
    // We have at least one block worth of old outstanding methods to try                                         // 1627
    // again. First: did onReconnect actually send anything? If not, we just                                      // 1628
    // restore all outstanding methods and run the first block.                                                   // 1629
    if (_.isEmpty(self._outstandingMethodBlocks)) {                                                               // 1630
      self._outstandingMethodBlocks = oldOutstandingMethodBlocks;                                                 // 1631
      self._sendOutstandingMethods();                                                                             // 1632
      return;                                                                                                     // 1633
    }                                                                                                             // 1634
                                                                                                                  // 1635
    // OK, there are blocks on both sides. Special case: merge the last block of                                  // 1636
    // the reconnect methods with the first block of the original methods, if                                     // 1637
    // neither of them are "wait" blocks.                                                                         // 1638
    if (!_.last(self._outstandingMethodBlocks).wait &&                                                            // 1639
        !oldOutstandingMethodBlocks[0].wait) {                                                                    // 1640
      _.each(oldOutstandingMethodBlocks[0].methods, function (m) {                                                // 1641
        _.last(self._outstandingMethodBlocks).methods.push(m);                                                    // 1642
                                                                                                                  // 1643
        // If this "last block" is also the first block, send the message.                                        // 1644
        if (self._outstandingMethodBlocks.length === 1)                                                           // 1645
          m.sendMessage();                                                                                        // 1646
      });                                                                                                         // 1647
                                                                                                                  // 1648
      oldOutstandingMethodBlocks.shift();                                                                         // 1649
    }                                                                                                             // 1650
                                                                                                                  // 1651
    // Now add the rest of the original blocks on.                                                                // 1652
    _.each(oldOutstandingMethodBlocks, function (block) {                                                         // 1653
      self._outstandingMethodBlocks.push(block);                                                                  // 1654
    });                                                                                                           // 1655
  },                                                                                                              // 1656
                                                                                                                  // 1657
  // We can accept a hot code push if there are no methods in flight.                                             // 1658
  _readyToMigrate: function() {                                                                                   // 1659
    var self = this;                                                                                              // 1660
    return _.isEmpty(self._methodInvokers);                                                                       // 1661
  },                                                                                                              // 1662
                                                                                                                  // 1663
  // If we were blocking a migration, see if it's now possible to continue.                                       // 1664
  // Call whenever the set of outstanding/blocked methods shrinks.                                                // 1665
  _maybeMigrate: function () {                                                                                    // 1666
    var self = this;                                                                                              // 1667
    if (self._retryMigrate && self._readyToMigrate()) {                                                           // 1668
      self._retryMigrate();                                                                                       // 1669
      self._retryMigrate = null;                                                                                  // 1670
    }                                                                                                             // 1671
  }                                                                                                               // 1672
});                                                                                                               // 1673
                                                                                                                  // 1674
LivedataTest.Connection = Connection;                                                                             // 1675
                                                                                                                  // 1676
// @param url {String} URL to Meteor app,                                                                         // 1677
//     e.g.:                                                                                                      // 1678
//     "subdomain.meteor.com",                                                                                    // 1679
//     "http://subdomain.meteor.com",                                                                             // 1680
//     "/",                                                                                                       // 1681
//     "ddp+sockjs://ddp--****-foo.meteor.com/sockjs"                                                             // 1682
                                                                                                                  // 1683
/**                                                                                                               // 1684
 * @summary Connect to the server of a different Meteor application to subscribe to its document sets and invoke its remote methods.
 * @locus Anywhere                                                                                                // 1686
 * @param {String} url The URL of another Meteor application.                                                     // 1687
 */                                                                                                               // 1688
DDP.connect = function (url, options) {                                                                           // 1689
  var ret = new Connection(url, options);                                                                         // 1690
  allConnections.push(ret); // hack. see below.                                                                   // 1691
  return ret;                                                                                                     // 1692
};                                                                                                                // 1693
                                                                                                                  // 1694
// Hack for `spiderable` package: a way to see if the page is done                                                // 1695
// loading all the data it needs.                                                                                 // 1696
//                                                                                                                // 1697
allConnections = [];                                                                                              // 1698
DDP._allSubscriptionsReady = function () {                                                                        // 1699
  return _.all(allConnections, function (conn) {                                                                  // 1700
    return _.all(conn._subscriptions, function (sub) {                                                            // 1701
      return sub.ready;                                                                                           // 1702
    });                                                                                                           // 1703
  });                                                                                                             // 1704
};                                                                                                                // 1705
                                                                                                                  // 1706
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                //
// packages/ddp-client/client_convenience.js                                                                      //
//                                                                                                                //
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                  //
// Meteor.refresh can be called on the client (if you're in common code) but it                                   // 1
// only has an effect on the server.                                                                              // 2
Meteor.refresh = function (notification) {                                                                        // 3
};                                                                                                                // 4
                                                                                                                  // 5
if (Meteor.isClient) {                                                                                            // 6
  // By default, try to connect back to the same endpoint as the page                                             // 7
  // was served from.                                                                                             // 8
  //                                                                                                              // 9
  // XXX We should be doing this a different way. Right now we don't                                              // 10
  // include ROOT_URL_PATH_PREFIX when computing ddpUrl. (We don't                                                // 11
  // include it on the server when computing                                                                      // 12
  // DDP_DEFAULT_CONNECTION_URL, and we don't include it in our                                                   // 13
  // default, '/'.) We get by with this because DDP.connect then                                                  // 14
  // forces the URL passed to it to be interpreted relative to the                                                // 15
  // app's deploy path, even if it is absolute. Instead, we should                                                // 16
  // make DDP_DEFAULT_CONNECTION_URL, if set, include the path prefix;                                            // 17
  // make the default ddpUrl be '' rather that '/'; and make                                                      // 18
  // _translateUrl in stream_client_common.js not force absolute paths                                            // 19
  // to be treated like relative paths. See also                                                                  // 20
  // stream_client_common.js #RationalizingRelativeDDPURLs                                                        // 21
  var ddpUrl = '/';                                                                                               // 22
  if (typeof __meteor_runtime_config__ !== "undefined") {                                                         // 23
    if (__meteor_runtime_config__.DDP_DEFAULT_CONNECTION_URL)                                                     // 24
      ddpUrl = __meteor_runtime_config__.DDP_DEFAULT_CONNECTION_URL;                                              // 25
  }                                                                                                               // 26
                                                                                                                  // 27
  var retry = new Retry();                                                                                        // 28
                                                                                                                  // 29
  var onDDPVersionNegotiationFailure = function (description) {                                                   // 30
    Meteor._debug(description);                                                                                   // 31
    if (Package.reload) {                                                                                         // 32
      var migrationData = Package.reload.Reload._migrationData('livedata') || {};                                 // 33
      var failures = migrationData.DDPVersionNegotiationFailures || 0;                                            // 34
      ++failures;                                                                                                 // 35
      Package.reload.Reload._onMigrate('livedata', function () {                                                  // 36
        return [true, {DDPVersionNegotiationFailures: failures}];                                                 // 37
      });                                                                                                         // 38
      retry.retryLater(failures, function () {                                                                    // 39
        Package.reload.Reload._reload();                                                                          // 40
      });                                                                                                         // 41
    }                                                                                                             // 42
  };                                                                                                              // 43
                                                                                                                  // 44
  Meteor.connection =                                                                                             // 45
    DDP.connect(ddpUrl, {                                                                                         // 46
      onDDPVersionNegotiationFailure: onDDPVersionNegotiationFailure                                              // 47
    });                                                                                                           // 48
                                                                                                                  // 49
  // Proxy the public methods of Meteor.connection so they can                                                    // 50
  // be called directly on Meteor.                                                                                // 51
  _.each(['subscribe', 'methods', 'call', 'apply', 'status', 'reconnect',                                         // 52
          'disconnect'],                                                                                          // 53
         function (name) {                                                                                        // 54
           Meteor[name] = _.bind(Meteor.connection[name], Meteor.connection);                                     // 55
         });                                                                                                      // 56
} else {                                                                                                          // 57
  // Never set up a default connection on the server. Don't even map                                              // 58
  // subscribe/call/etc onto Meteor.                                                                              // 59
  Meteor.connection = null;                                                                                       // 60
}                                                                                                                 // 61
                                                                                                                  // 62
// Meteor.connection used to be called                                                                            // 63
// Meteor.default_connection. Provide backcompat as a courtesy even                                               // 64
// though it was never documented.                                                                                // 65
// XXX COMPAT WITH 0.6.4                                                                                          // 66
Meteor.default_connection = Meteor.connection;                                                                    // 67
                                                                                                                  // 68
// We should transition from Meteor.connect to DDP.connect.                                                       // 69
// XXX COMPAT WITH 0.6.4                                                                                          // 70
Meteor.connect = DDP.connect;                                                                                     // 71
                                                                                                                  // 72
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['ddp-client'] = {}, {
  DDP: DDP,
  LivedataTest: LivedataTest
});

})();
