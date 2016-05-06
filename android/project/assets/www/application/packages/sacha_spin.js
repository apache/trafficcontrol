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
var Template = Package.templating.Template;
var _ = Package.underscore._;
var Blaze = Package.blaze.Blaze;
var UI = Package.blaze.UI;
var Handlebars = Package.blaze.Handlebars;
var Spacebars = Package.spacebars.Spacebars;
var HTML = Package.htmljs.HTML;

(function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/sacha_spin/packages/sacha_spin.js                                                                          //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
(function () {                                                                                                         // 1
                                                                                                                       // 2
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// packages/sacha:spin/.npm/package/node_modules/spin.js/spin.js                                                    //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
/**                                                                                                                 // 1
 * Copyright (c) 2011-2014 Felix Gnass                                                                              // 2
 * Licensed under the MIT license                                                                                   // 3
 * http://spin.js.org/                                                                                              // 4
 *                                                                                                                  // 5
 * Example:                                                                                                         // 6
    var opts = {                                                                                                    // 7
      lines: 12             // The number of lines to draw                                                          // 8
    , length: 7             // The length of each line                                                              // 9
    , width: 5              // The line thickness                                                                   // 10
    , radius: 10            // The radius of the inner circle                                                       // 11
    , scale: 1.0            // Scales overall size of the spinner                                                   // 12
    , corners: 1            // Roundness (0..1)                                                                     // 13
    , color: '#000'         // #rgb or #rrggbb                                                                      // 14
    , opacity: 1/4          // Opacity of the lines                                                                 // 15
    , rotate: 0             // Rotation offset                                                                      // 16
    , direction: 1          // 1: clockwise, -1: counterclockwise                                                   // 17
    , speed: 1              // Rounds per second                                                                    // 18
    , trail: 100            // Afterglow percentage                                                                 // 19
    , fps: 20               // Frames per second when using setTimeout()                                            // 20
    , zIndex: 2e9           // Use a high z-index by default                                                        // 21
    , className: 'spinner'  // CSS class to assign to the element                                                   // 22
    , top: '50%'            // center vertically                                                                    // 23
    , left: '50%'           // center horizontally                                                                  // 24
    , shadow: false         // Whether to render a shadow                                                           // 25
    , hwaccel: false        // Whether to use hardware acceleration (might be buggy)                                // 26
    , position: 'absolute'  // Element positioning                                                                  // 27
    }                                                                                                               // 28
    var target = document.getElementById('foo')                                                                     // 29
    var spinner = new Spinner(opts).spin(target)                                                                    // 30
 */                                                                                                                 // 31
;(function (root, factory) {                                                                                        // 32
                                                                                                                    // 33
  /* CommonJS */                                                                                                    // 34
  if (typeof exports == 'object') module.exports = factory()                                                        // 35
                                                                                                                    // 36
  /* AMD module */                                                                                                  // 37
  else if (typeof define == 'function' && define.amd) define(factory)                                               // 38
                                                                                                                    // 39
  /* Browser global */                                                                                              // 40
  else root.Spinner = factory()                                                                                     // 41
}(this, function () {                                                                                               // 42
  "use strict"                                                                                                      // 43
                                                                                                                    // 44
  var prefixes = ['webkit', 'Moz', 'ms', 'O'] /* Vendor prefixes */                                                 // 45
    , animations = {} /* Animation rules keyed by their name */                                                     // 46
    , useCssAnimations /* Whether to use CSS animations or setTimeout */                                            // 47
    , sheet /* A stylesheet to hold the @keyframe or VML rules. */                                                  // 48
                                                                                                                    // 49
  /**                                                                                                               // 50
   * Utility function to create elements. If no tag name is given,                                                  // 51
   * a DIV is created. Optionally properties can be passed.                                                         // 52
   */                                                                                                               // 53
  function createEl (tag, prop) {                                                                                   // 54
    var el = document.createElement(tag || 'div')                                                                   // 55
      , n                                                                                                           // 56
                                                                                                                    // 57
    for (n in prop) el[n] = prop[n]                                                                                 // 58
    return el                                                                                                       // 59
  }                                                                                                                 // 60
                                                                                                                    // 61
  /**                                                                                                               // 62
   * Appends children and returns the parent.                                                                       // 63
   */                                                                                                               // 64
  function ins (parent /* child1, child2, ...*/) {                                                                  // 65
    for (var i = 1, n = arguments.length; i < n; i++) {                                                             // 66
      parent.appendChild(arguments[i])                                                                              // 67
    }                                                                                                               // 68
                                                                                                                    // 69
    return parent                                                                                                   // 70
  }                                                                                                                 // 71
                                                                                                                    // 72
  /**                                                                                                               // 73
   * Creates an opacity keyframe animation rule and returns its name.                                               // 74
   * Since most mobile Webkits have timing issues with animation-delay,                                             // 75
   * we create separate rules for each line/segment.                                                                // 76
   */                                                                                                               // 77
  function addAnimation (alpha, trail, i, lines) {                                                                  // 78
    var name = ['opacity', trail, ~~(alpha * 100), i, lines].join('-')                                              // 79
      , start = 0.01 + i/lines * 100                                                                                // 80
      , z = Math.max(1 - (1-alpha) / trail * (100-start), alpha)                                                    // 81
      , prefix = useCssAnimations.substring(0, useCssAnimations.indexOf('Animation')).toLowerCase()                 // 82
      , pre = prefix && '-' + prefix + '-' || ''                                                                    // 83
                                                                                                                    // 84
    if (!animations[name]) {                                                                                        // 85
      sheet.insertRule(                                                                                             // 86
        '@' + pre + 'keyframes ' + name + '{' +                                                                     // 87
        '0%{opacity:' + z + '}' +                                                                                   // 88
        start + '%{opacity:' + alpha + '}' +                                                                        // 89
        (start+0.01) + '%{opacity:1}' +                                                                             // 90
        (start+trail) % 100 + '%{opacity:' + alpha + '}' +                                                          // 91
        '100%{opacity:' + z + '}' +                                                                                 // 92
        '}', sheet.cssRules.length)                                                                                 // 93
                                                                                                                    // 94
      animations[name] = 1                                                                                          // 95
    }                                                                                                               // 96
                                                                                                                    // 97
    return name                                                                                                     // 98
  }                                                                                                                 // 99
                                                                                                                    // 100
  /**                                                                                                               // 101
   * Tries various vendor prefixes and returns the first supported property.                                        // 102
   */                                                                                                               // 103
  function vendor (el, prop) {                                                                                      // 104
    var s = el.style                                                                                                // 105
      , pp                                                                                                          // 106
      , i                                                                                                           // 107
                                                                                                                    // 108
    prop = prop.charAt(0).toUpperCase() + prop.slice(1)                                                             // 109
    if (s[prop] !== undefined) return prop                                                                          // 110
    for (i = 0; i < prefixes.length; i++) {                                                                         // 111
      pp = prefixes[i]+prop                                                                                         // 112
      if (s[pp] !== undefined) return pp                                                                            // 113
    }                                                                                                               // 114
  }                                                                                                                 // 115
                                                                                                                    // 116
  /**                                                                                                               // 117
   * Sets multiple style properties at once.                                                                        // 118
   */                                                                                                               // 119
  function css (el, prop) {                                                                                         // 120
    for (var n in prop) {                                                                                           // 121
      el.style[vendor(el, n) || n] = prop[n]                                                                        // 122
    }                                                                                                               // 123
                                                                                                                    // 124
    return el                                                                                                       // 125
  }                                                                                                                 // 126
                                                                                                                    // 127
  /**                                                                                                               // 128
   * Fills in default values.                                                                                       // 129
   */                                                                                                               // 130
  function merge (obj) {                                                                                            // 131
    for (var i = 1; i < arguments.length; i++) {                                                                    // 132
      var def = arguments[i]                                                                                        // 133
      for (var n in def) {                                                                                          // 134
        if (obj[n] === undefined) obj[n] = def[n]                                                                   // 135
      }                                                                                                             // 136
    }                                                                                                               // 137
    return obj                                                                                                      // 138
  }                                                                                                                 // 139
                                                                                                                    // 140
  /**                                                                                                               // 141
   * Returns the line color from the given string or array.                                                         // 142
   */                                                                                                               // 143
  function getColor (color, idx) {                                                                                  // 144
    return typeof color == 'string' ? color : color[idx % color.length]                                             // 145
  }                                                                                                                 // 146
                                                                                                                    // 147
  // Built-in defaults                                                                                              // 148
                                                                                                                    // 149
  var defaults = {                                                                                                  // 150
    lines: 12             // The number of lines to draw                                                            // 151
  , length: 7             // The length of each line                                                                // 152
  , width: 5              // The line thickness                                                                     // 153
  , radius: 10            // The radius of the inner circle                                                         // 154
  , scale: 1.0            // Scales overall size of the spinner                                                     // 155
  , corners: 1            // Roundness (0..1)                                                                       // 156
  , color: '#000'         // #rgb or #rrggbb                                                                        // 157
  , opacity: 1/4          // Opacity of the lines                                                                   // 158
  , rotate: 0             // Rotation offset                                                                        // 159
  , direction: 1          // 1: clockwise, -1: counterclockwise                                                     // 160
  , speed: 1              // Rounds per second                                                                      // 161
  , trail: 100            // Afterglow percentage                                                                   // 162
  , fps: 20               // Frames per second when using setTimeout()                                              // 163
  , zIndex: 2e9           // Use a high z-index by default                                                          // 164
  , className: 'spinner'  // CSS class to assign to the element                                                     // 165
  , top: '50%'            // center vertically                                                                      // 166
  , left: '50%'           // center horizontally                                                                    // 167
  , shadow: false         // Whether to render a shadow                                                             // 168
  , hwaccel: false        // Whether to use hardware acceleration (might be buggy)                                  // 169
  , position: 'absolute'  // Element positioning                                                                    // 170
  }                                                                                                                 // 171
                                                                                                                    // 172
  /** The constructor */                                                                                            // 173
  function Spinner (o) {                                                                                            // 174
    this.opts = merge(o || {}, Spinner.defaults, defaults)                                                          // 175
  }                                                                                                                 // 176
                                                                                                                    // 177
  // Global defaults that override the built-ins:                                                                   // 178
  Spinner.defaults = {}                                                                                             // 179
                                                                                                                    // 180
  merge(Spinner.prototype, {                                                                                        // 181
    /**                                                                                                             // 182
     * Adds the spinner to the given target element. If this instance is already                                    // 183
     * spinning, it is automatically removed from its previous target b calling                                     // 184
     * stop() internally.                                                                                           // 185
     */                                                                                                             // 186
    spin: function (target) {                                                                                       // 187
      this.stop()                                                                                                   // 188
                                                                                                                    // 189
      var self = this                                                                                               // 190
        , o = self.opts                                                                                             // 191
        , el = self.el = createEl(null, {className: o.className})                                                   // 192
                                                                                                                    // 193
      css(el, {                                                                                                     // 194
        position: o.position                                                                                        // 195
      , width: 0                                                                                                    // 196
      , zIndex: o.zIndex                                                                                            // 197
      , left: o.left                                                                                                // 198
      , top: o.top                                                                                                  // 199
      })                                                                                                            // 200
                                                                                                                    // 201
      if (target) {                                                                                                 // 202
        target.insertBefore(el, target.firstChild || null)                                                          // 203
      }                                                                                                             // 204
                                                                                                                    // 205
      el.setAttribute('role', 'progressbar')                                                                        // 206
      self.lines(el, self.opts)                                                                                     // 207
                                                                                                                    // 208
      if (!useCssAnimations) {                                                                                      // 209
        // No CSS animation support, use setTimeout() instead                                                       // 210
        var i = 0                                                                                                   // 211
          , start = (o.lines - 1) * (1 - o.direction) / 2                                                           // 212
          , alpha                                                                                                   // 213
          , fps = o.fps                                                                                             // 214
          , f = fps / o.speed                                                                                       // 215
          , ostep = (1 - o.opacity) / (f * o.trail / 100)                                                           // 216
          , astep = f / o.lines                                                                                     // 217
                                                                                                                    // 218
        ;(function anim () {                                                                                        // 219
          i++                                                                                                       // 220
          for (var j = 0; j < o.lines; j++) {                                                                       // 221
            alpha = Math.max(1 - (i + (o.lines - j) * astep) % f * ostep, o.opacity)                                // 222
                                                                                                                    // 223
            self.opacity(el, j * o.direction + start, alpha, o)                                                     // 224
          }                                                                                                         // 225
          self.timeout = self.el && setTimeout(anim, ~~(1000 / fps))                                                // 226
        })()                                                                                                        // 227
      }                                                                                                             // 228
      return self                                                                                                   // 229
    }                                                                                                               // 230
                                                                                                                    // 231
    /**                                                                                                             // 232
     * Stops and removes the Spinner.                                                                               // 233
     */                                                                                                             // 234
  , stop: function () {                                                                                             // 235
      var el = this.el                                                                                              // 236
      if (el) {                                                                                                     // 237
        clearTimeout(this.timeout)                                                                                  // 238
        if (el.parentNode) el.parentNode.removeChild(el)                                                            // 239
        this.el = undefined                                                                                         // 240
      }                                                                                                             // 241
      return this                                                                                                   // 242
    }                                                                                                               // 243
                                                                                                                    // 244
    /**                                                                                                             // 245
     * Internal method that draws the individual lines. Will be overwritten                                         // 246
     * in VML fallback mode below.                                                                                  // 247
     */                                                                                                             // 248
  , lines: function (el, o) {                                                                                       // 249
      var i = 0                                                                                                     // 250
        , start = (o.lines - 1) * (1 - o.direction) / 2                                                             // 251
        , seg                                                                                                       // 252
                                                                                                                    // 253
      function fill (color, shadow) {                                                                               // 254
        return css(createEl(), {                                                                                    // 255
          position: 'absolute'                                                                                      // 256
        , width: o.scale * (o.length + o.width) + 'px'                                                              // 257
        , height: o.scale * o.width + 'px'                                                                          // 258
        , background: color                                                                                         // 259
        , boxShadow: shadow                                                                                         // 260
        , transformOrigin: 'left'                                                                                   // 261
        , transform: 'rotate(' + ~~(360/o.lines*i + o.rotate) + 'deg) translate(' + o.scale*o.radius + 'px' + ',0)' // 262
        , borderRadius: (o.corners * o.scale * o.width >> 1) + 'px'                                                 // 263
        })                                                                                                          // 264
      }                                                                                                             // 265
                                                                                                                    // 266
      for (; i < o.lines; i++) {                                                                                    // 267
        seg = css(createEl(), {                                                                                     // 268
          position: 'absolute'                                                                                      // 269
        , top: 1 + ~(o.scale * o.width / 2) + 'px'                                                                  // 270
        , transform: o.hwaccel ? 'translate3d(0,0,0)' : ''                                                          // 271
        , opacity: o.opacity                                                                                        // 272
        , animation: useCssAnimations && addAnimation(o.opacity, o.trail, start + i * o.direction, o.lines) + ' ' + 1 / o.speed + 's linear infinite'
        })                                                                                                          // 274
                                                                                                                    // 275
        if (o.shadow) ins(seg, css(fill('#000', '0 0 4px #000'), {top: '2px'}))                                     // 276
        ins(el, ins(seg, fill(getColor(o.color, i), '0 0 1px rgba(0,0,0,.1)')))                                     // 277
      }                                                                                                             // 278
      return el                                                                                                     // 279
    }                                                                                                               // 280
                                                                                                                    // 281
    /**                                                                                                             // 282
     * Internal method that adjusts the opacity of a single line.                                                   // 283
     * Will be overwritten in VML fallback mode below.                                                              // 284
     */                                                                                                             // 285
  , opacity: function (el, i, val) {                                                                                // 286
      if (i < el.childNodes.length) el.childNodes[i].style.opacity = val                                            // 287
    }                                                                                                               // 288
                                                                                                                    // 289
  })                                                                                                                // 290
                                                                                                                    // 291
                                                                                                                    // 292
  function initVML () {                                                                                             // 293
                                                                                                                    // 294
    /* Utility function to create a VML tag */                                                                      // 295
    function vml (tag, attr) {                                                                                      // 296
      return createEl('<' + tag + ' xmlns="urn:schemas-microsoft.com:vml" class="spin-vml">', attr)                 // 297
    }                                                                                                               // 298
                                                                                                                    // 299
    // No CSS transforms but VML support, add a CSS rule for VML elements:                                          // 300
    sheet.addRule('.spin-vml', 'behavior:url(#default#VML)')                                                        // 301
                                                                                                                    // 302
    Spinner.prototype.lines = function (el, o) {                                                                    // 303
      var r = o.scale * (o.length + o.width)                                                                        // 304
        , s = o.scale * 2 * r                                                                                       // 305
                                                                                                                    // 306
      function grp () {                                                                                             // 307
        return css(                                                                                                 // 308
          vml('group', {                                                                                            // 309
            coordsize: s + ' ' + s                                                                                  // 310
          , coordorigin: -r + ' ' + -r                                                                              // 311
          })                                                                                                        // 312
        , { width: s, height: s }                                                                                   // 313
        )                                                                                                           // 314
      }                                                                                                             // 315
                                                                                                                    // 316
      var margin = -(o.width + o.length) * o.scale * 2 + 'px'                                                       // 317
        , g = css(grp(), {position: 'absolute', top: margin, left: margin})                                         // 318
        , i                                                                                                         // 319
                                                                                                                    // 320
      function seg (i, dx, filter) {                                                                                // 321
        ins(                                                                                                        // 322
          g                                                                                                         // 323
        , ins(                                                                                                      // 324
            css(grp(), {rotation: 360 / o.lines * i + 'deg', left: ~~dx})                                           // 325
          , ins(                                                                                                    // 326
              css(                                                                                                  // 327
                vml('roundrect', {arcsize: o.corners})                                                              // 328
              , { width: r                                                                                          // 329
                , height: o.scale * o.width                                                                         // 330
                , left: o.scale * o.radius                                                                          // 331
                , top: -o.scale * o.width >> 1                                                                      // 332
                , filter: filter                                                                                    // 333
                }                                                                                                   // 334
              )                                                                                                     // 335
            , vml('fill', {color: getColor(o.color, i), opacity: o.opacity})                                        // 336
            , vml('stroke', {opacity: 0}) // transparent stroke to fix color bleeding upon opacity change           // 337
            )                                                                                                       // 338
          )                                                                                                         // 339
        )                                                                                                           // 340
      }                                                                                                             // 341
                                                                                                                    // 342
      if (o.shadow)                                                                                                 // 343
        for (i = 1; i <= o.lines; i++) {                                                                            // 344
          seg(i, -2, 'progid:DXImageTransform.Microsoft.Blur(pixelradius=2,makeshadow=1,shadowopacity=.3)')         // 345
        }                                                                                                           // 346
                                                                                                                    // 347
      for (i = 1; i <= o.lines; i++) seg(i)                                                                         // 348
      return ins(el, g)                                                                                             // 349
    }                                                                                                               // 350
                                                                                                                    // 351
    Spinner.prototype.opacity = function (el, i, val, o) {                                                          // 352
      var c = el.firstChild                                                                                         // 353
      o = o.shadow && o.lines || 0                                                                                  // 354
      if (c && i + o < c.childNodes.length) {                                                                       // 355
        c = c.childNodes[i + o]; c = c && c.firstChild; c = c && c.firstChild                                       // 356
        if (c) c.opacity = val                                                                                      // 357
      }                                                                                                             // 358
    }                                                                                                               // 359
  }                                                                                                                 // 360
                                                                                                                    // 361
  if (typeof document !== 'undefined') {                                                                            // 362
    sheet = (function () {                                                                                          // 363
      var el = createEl('style', {type : 'text/css'})                                                               // 364
      ins(document.getElementsByTagName('head')[0], el)                                                             // 365
      return el.sheet || el.styleSheet                                                                              // 366
    }())                                                                                                            // 367
                                                                                                                    // 368
    var probe = css(createEl('group'), {behavior: 'url(#default#VML)'})                                             // 369
                                                                                                                    // 370
    if (!vendor(probe, 'transform') && probe.adj) initVML()                                                         // 371
    else useCssAnimations = vendor(probe, 'animation')                                                              // 372
  }                                                                                                                 // 373
                                                                                                                    // 374
  return Spinner                                                                                                    // 375
                                                                                                                    // 376
}));                                                                                                                // 377
                                                                                                                    // 378
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       // 388
}).call(this);                                                                                                         // 389
                                                                                                                       // 390
                                                                                                                       // 391
                                                                                                                       // 392
                                                                                                                       // 393
                                                                                                                       // 394
                                                                                                                       // 395
(function () {                                                                                                         // 396
                                                                                                                       // 397
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// packages/sacha:spin/lib/template.spinner.js                                                                      //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
                                                                                                                    // 1
Template.__checkName("spinner");                                                                                    // 2
Template["spinner"] = new Template("Template.spinner", (function() {                                                // 3
  var view = this;                                                                                                  // 4
  return HTML.Raw('<div class="spinner-container"></div>');                                                         // 5
}));                                                                                                                // 6
                                                                                                                    // 7
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       // 412
}).call(this);                                                                                                         // 413
                                                                                                                       // 414
                                                                                                                       // 415
                                                                                                                       // 416
                                                                                                                       // 417
                                                                                                                       // 418
                                                                                                                       // 419
(function () {                                                                                                         // 420
                                                                                                                       // 421
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                  //
// packages/sacha:spin/lib/spinner.js                                                                               //
//                                                                                                                  //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                    //
Template.spinner.onRendered(function(){                                                                             // 1
  var options = _.extend({}, Meteor.Spinner.options, this.data);                                                    // 2
                                                                                                                    // 3
  this.spinner = new Spinner(options);                                                                              // 4
  this.spinner.spin(this.firstNode);                                                                                // 5
});                                                                                                                 // 6
                                                                                                                    // 7
                                                                                                                    // 8
Template.spinner.onDestroyed(function(){                                                                            // 9
  this.spinner && this.spinner.stop();                                                                              // 10
});                                                                                                                 // 11
                                                                                                                    // 12
                                                                                                                    // 13
Meteor.Spinner = {                                                                                                  // 14
  options: {                                                                                                        // 15
    lines: 13,  // The number of lines to draw                                                                      // 16
    length: 8,  // The length of each line                                                                          // 17
    width: 3,  // The line thickness                                                                                // 18
    radius: 12,  // The radius of the inner circle                                                                  // 19
    corners: 1,  // Corner roundness (0..1)                                                                         // 20
    rotate: 0,  // The rotation offset                                                                              // 21
    direction: 1,  // 1: clockwise, -1: counterclockwise                                                            // 22
    color: '#000',  // #rgb or #rrggbb                                                                              // 23
    speed: 1.2,  // Rounds per second                                                                               // 24
    trail: 60,  // Afterglow percentage                                                                             // 25
    shadow: false,  // Whether to render a shadow                                                                   // 26
    hwaccel: false,  // Whether to use hardware acceleration                                                        // 27
    className: 'spinner', // The CSS class to assign to the spinner                                                 // 28
    zIndex: 2e9,  // The z-index (defaults to 2000000000)                                                           // 29
    top: '50%',  // Top position relative to parent in px                                                           // 30
    left: '50%'  // Left position relative to parent in px                                                          // 31
  }                                                                                                                 // 32
};                                                                                                                  // 33
                                                                                                                    // 34
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       // 463
}).call(this);                                                                                                         // 464
                                                                                                                       // 465
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
Package['sacha:spin'] = {};

})();
