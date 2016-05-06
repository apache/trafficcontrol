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
var $ = Package.jquery.$;
var jQuery = Package.jquery.jQuery;

(function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/twbs_bootstrap/dist/js/bootstrap.js                                                                        //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
/*!                                                                                                                    // 1
 * Bootstrap v3.3.6 (http://getbootstrap.com)                                                                          // 2
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 3
 * Licensed under the MIT license                                                                                      // 4
 */                                                                                                                    // 5
                                                                                                                       // 6
if (typeof jQuery === 'undefined') {                                                                                   // 7
  throw new Error('Bootstrap\'s JavaScript requires jQuery')                                                           // 8
}                                                                                                                      // 9
                                                                                                                       // 10
+function ($) {                                                                                                        // 11
  'use strict';                                                                                                        // 12
  var version = $.fn.jquery.split(' ')[0].split('.')                                                                   // 13
  if ((version[0] < 2 && version[1] < 9) || (version[0] == 1 && version[1] == 9 && version[2] < 1) || (version[0] > 2)) {
    throw new Error('Bootstrap\'s JavaScript requires jQuery version 1.9.1 or higher, but lower than version 3')       // 15
  }                                                                                                                    // 16
}(jQuery);                                                                                                             // 17
                                                                                                                       // 18
/* ========================================================================                                            // 19
 * Bootstrap: transition.js v3.3.6                                                                                     // 20
 * http://getbootstrap.com/javascript/#transitions                                                                     // 21
 * ========================================================================                                            // 22
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 23
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 24
 * ======================================================================== */                                         // 25
                                                                                                                       // 26
                                                                                                                       // 27
+function ($) {                                                                                                        // 28
  'use strict';                                                                                                        // 29
                                                                                                                       // 30
  // CSS TRANSITION SUPPORT (Shoutout: http://www.modernizr.com/)                                                      // 31
  // ============================================================                                                      // 32
                                                                                                                       // 33
  function transitionEnd() {                                                                                           // 34
    var el = document.createElement('bootstrap')                                                                       // 35
                                                                                                                       // 36
    var transEndEventNames = {                                                                                         // 37
      WebkitTransition : 'webkitTransitionEnd',                                                                        // 38
      MozTransition    : 'transitionend',                                                                              // 39
      OTransition      : 'oTransitionEnd otransitionend',                                                              // 40
      transition       : 'transitionend'                                                                               // 41
    }                                                                                                                  // 42
                                                                                                                       // 43
    for (var name in transEndEventNames) {                                                                             // 44
      if (el.style[name] !== undefined) {                                                                              // 45
        return { end: transEndEventNames[name] }                                                                       // 46
      }                                                                                                                // 47
    }                                                                                                                  // 48
                                                                                                                       // 49
    return false // explicit for ie8 (  ._.)                                                                           // 50
  }                                                                                                                    // 51
                                                                                                                       // 52
  // http://blog.alexmaccaw.com/css-transitions                                                                        // 53
  $.fn.emulateTransitionEnd = function (duration) {                                                                    // 54
    var called = false                                                                                                 // 55
    var $el = this                                                                                                     // 56
    $(this).one('bsTransitionEnd', function () { called = true })                                                      // 57
    var callback = function () { if (!called) $($el).trigger($.support.transition.end) }                               // 58
    setTimeout(callback, duration)                                                                                     // 59
    return this                                                                                                        // 60
  }                                                                                                                    // 61
                                                                                                                       // 62
  $(function () {                                                                                                      // 63
    $.support.transition = transitionEnd()                                                                             // 64
                                                                                                                       // 65
    if (!$.support.transition) return                                                                                  // 66
                                                                                                                       // 67
    $.event.special.bsTransitionEnd = {                                                                                // 68
      bindType: $.support.transition.end,                                                                              // 69
      delegateType: $.support.transition.end,                                                                          // 70
      handle: function (e) {                                                                                           // 71
        if ($(e.target).is(this)) return e.handleObj.handler.apply(this, arguments)                                    // 72
      }                                                                                                                // 73
    }                                                                                                                  // 74
  })                                                                                                                   // 75
                                                                                                                       // 76
}(jQuery);                                                                                                             // 77
                                                                                                                       // 78
/* ========================================================================                                            // 79
 * Bootstrap: alert.js v3.3.6                                                                                          // 80
 * http://getbootstrap.com/javascript/#alerts                                                                          // 81
 * ========================================================================                                            // 82
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 83
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 84
 * ======================================================================== */                                         // 85
                                                                                                                       // 86
                                                                                                                       // 87
+function ($) {                                                                                                        // 88
  'use strict';                                                                                                        // 89
                                                                                                                       // 90
  // ALERT CLASS DEFINITION                                                                                            // 91
  // ======================                                                                                            // 92
                                                                                                                       // 93
  var dismiss = '[data-dismiss="alert"]'                                                                               // 94
  var Alert   = function (el) {                                                                                        // 95
    $(el).on('click', dismiss, this.close)                                                                             // 96
  }                                                                                                                    // 97
                                                                                                                       // 98
  Alert.VERSION = '3.3.6'                                                                                              // 99
                                                                                                                       // 100
  Alert.TRANSITION_DURATION = 150                                                                                      // 101
                                                                                                                       // 102
  Alert.prototype.close = function (e) {                                                                               // 103
    var $this    = $(this)                                                                                             // 104
    var selector = $this.attr('data-target')                                                                           // 105
                                                                                                                       // 106
    if (!selector) {                                                                                                   // 107
      selector = $this.attr('href')                                                                                    // 108
      selector = selector && selector.replace(/.*(?=#[^\s]*$)/, '') // strip for ie7                                   // 109
    }                                                                                                                  // 110
                                                                                                                       // 111
    var $parent = $(selector)                                                                                          // 112
                                                                                                                       // 113
    if (e) e.preventDefault()                                                                                          // 114
                                                                                                                       // 115
    if (!$parent.length) {                                                                                             // 116
      $parent = $this.closest('.alert')                                                                                // 117
    }                                                                                                                  // 118
                                                                                                                       // 119
    $parent.trigger(e = $.Event('close.bs.alert'))                                                                     // 120
                                                                                                                       // 121
    if (e.isDefaultPrevented()) return                                                                                 // 122
                                                                                                                       // 123
    $parent.removeClass('in')                                                                                          // 124
                                                                                                                       // 125
    function removeElement() {                                                                                         // 126
      // detach from parent, fire event then clean up data                                                             // 127
      $parent.detach().trigger('closed.bs.alert').remove()                                                             // 128
    }                                                                                                                  // 129
                                                                                                                       // 130
    $.support.transition && $parent.hasClass('fade') ?                                                                 // 131
      $parent                                                                                                          // 132
        .one('bsTransitionEnd', removeElement)                                                                         // 133
        .emulateTransitionEnd(Alert.TRANSITION_DURATION) :                                                             // 134
      removeElement()                                                                                                  // 135
  }                                                                                                                    // 136
                                                                                                                       // 137
                                                                                                                       // 138
  // ALERT PLUGIN DEFINITION                                                                                           // 139
  // =======================                                                                                           // 140
                                                                                                                       // 141
  function Plugin(option) {                                                                                            // 142
    return this.each(function () {                                                                                     // 143
      var $this = $(this)                                                                                              // 144
      var data  = $this.data('bs.alert')                                                                               // 145
                                                                                                                       // 146
      if (!data) $this.data('bs.alert', (data = new Alert(this)))                                                      // 147
      if (typeof option == 'string') data[option].call($this)                                                          // 148
    })                                                                                                                 // 149
  }                                                                                                                    // 150
                                                                                                                       // 151
  var old = $.fn.alert                                                                                                 // 152
                                                                                                                       // 153
  $.fn.alert             = Plugin                                                                                      // 154
  $.fn.alert.Constructor = Alert                                                                                       // 155
                                                                                                                       // 156
                                                                                                                       // 157
  // ALERT NO CONFLICT                                                                                                 // 158
  // =================                                                                                                 // 159
                                                                                                                       // 160
  $.fn.alert.noConflict = function () {                                                                                // 161
    $.fn.alert = old                                                                                                   // 162
    return this                                                                                                        // 163
  }                                                                                                                    // 164
                                                                                                                       // 165
                                                                                                                       // 166
  // ALERT DATA-API                                                                                                    // 167
  // ==============                                                                                                    // 168
                                                                                                                       // 169
  $(document).on('click.bs.alert.data-api', dismiss, Alert.prototype.close)                                            // 170
                                                                                                                       // 171
}(jQuery);                                                                                                             // 172
                                                                                                                       // 173
/* ========================================================================                                            // 174
 * Bootstrap: button.js v3.3.6                                                                                         // 175
 * http://getbootstrap.com/javascript/#buttons                                                                         // 176
 * ========================================================================                                            // 177
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 178
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 179
 * ======================================================================== */                                         // 180
                                                                                                                       // 181
                                                                                                                       // 182
+function ($) {                                                                                                        // 183
  'use strict';                                                                                                        // 184
                                                                                                                       // 185
  // BUTTON PUBLIC CLASS DEFINITION                                                                                    // 186
  // ==============================                                                                                    // 187
                                                                                                                       // 188
  var Button = function (element, options) {                                                                           // 189
    this.$element  = $(element)                                                                                        // 190
    this.options   = $.extend({}, Button.DEFAULTS, options)                                                            // 191
    this.isLoading = false                                                                                             // 192
  }                                                                                                                    // 193
                                                                                                                       // 194
  Button.VERSION  = '3.3.6'                                                                                            // 195
                                                                                                                       // 196
  Button.DEFAULTS = {                                                                                                  // 197
    loadingText: 'loading...'                                                                                          // 198
  }                                                                                                                    // 199
                                                                                                                       // 200
  Button.prototype.setState = function (state) {                                                                       // 201
    var d    = 'disabled'                                                                                              // 202
    var $el  = this.$element                                                                                           // 203
    var val  = $el.is('input') ? 'val' : 'html'                                                                        // 204
    var data = $el.data()                                                                                              // 205
                                                                                                                       // 206
    state += 'Text'                                                                                                    // 207
                                                                                                                       // 208
    if (data.resetText == null) $el.data('resetText', $el[val]())                                                      // 209
                                                                                                                       // 210
    // push to event loop to allow forms to submit                                                                     // 211
    setTimeout($.proxy(function () {                                                                                   // 212
      $el[val](data[state] == null ? this.options[state] : data[state])                                                // 213
                                                                                                                       // 214
      if (state == 'loadingText') {                                                                                    // 215
        this.isLoading = true                                                                                          // 216
        $el.addClass(d).attr(d, d)                                                                                     // 217
      } else if (this.isLoading) {                                                                                     // 218
        this.isLoading = false                                                                                         // 219
        $el.removeClass(d).removeAttr(d)                                                                               // 220
      }                                                                                                                // 221
    }, this), 0)                                                                                                       // 222
  }                                                                                                                    // 223
                                                                                                                       // 224
  Button.prototype.toggle = function () {                                                                              // 225
    var changed = true                                                                                                 // 226
    var $parent = this.$element.closest('[data-toggle="buttons"]')                                                     // 227
                                                                                                                       // 228
    if ($parent.length) {                                                                                              // 229
      var $input = this.$element.find('input')                                                                         // 230
      if ($input.prop('type') == 'radio') {                                                                            // 231
        if ($input.prop('checked')) changed = false                                                                    // 232
        $parent.find('.active').removeClass('active')                                                                  // 233
        this.$element.addClass('active')                                                                               // 234
      } else if ($input.prop('type') == 'checkbox') {                                                                  // 235
        if (($input.prop('checked')) !== this.$element.hasClass('active')) changed = false                             // 236
        this.$element.toggleClass('active')                                                                            // 237
      }                                                                                                                // 238
      $input.prop('checked', this.$element.hasClass('active'))                                                         // 239
      if (changed) $input.trigger('change')                                                                            // 240
    } else {                                                                                                           // 241
      this.$element.attr('aria-pressed', !this.$element.hasClass('active'))                                            // 242
      this.$element.toggleClass('active')                                                                              // 243
    }                                                                                                                  // 244
  }                                                                                                                    // 245
                                                                                                                       // 246
                                                                                                                       // 247
  // BUTTON PLUGIN DEFINITION                                                                                          // 248
  // ========================                                                                                          // 249
                                                                                                                       // 250
  function Plugin(option) {                                                                                            // 251
    return this.each(function () {                                                                                     // 252
      var $this   = $(this)                                                                                            // 253
      var data    = $this.data('bs.button')                                                                            // 254
      var options = typeof option == 'object' && option                                                                // 255
                                                                                                                       // 256
      if (!data) $this.data('bs.button', (data = new Button(this, options)))                                           // 257
                                                                                                                       // 258
      if (option == 'toggle') data.toggle()                                                                            // 259
      else if (option) data.setState(option)                                                                           // 260
    })                                                                                                                 // 261
  }                                                                                                                    // 262
                                                                                                                       // 263
  var old = $.fn.button                                                                                                // 264
                                                                                                                       // 265
  $.fn.button             = Plugin                                                                                     // 266
  $.fn.button.Constructor = Button                                                                                     // 267
                                                                                                                       // 268
                                                                                                                       // 269
  // BUTTON NO CONFLICT                                                                                                // 270
  // ==================                                                                                                // 271
                                                                                                                       // 272
  $.fn.button.noConflict = function () {                                                                               // 273
    $.fn.button = old                                                                                                  // 274
    return this                                                                                                        // 275
  }                                                                                                                    // 276
                                                                                                                       // 277
                                                                                                                       // 278
  // BUTTON DATA-API                                                                                                   // 279
  // ===============                                                                                                   // 280
                                                                                                                       // 281
  $(document)                                                                                                          // 282
    .on('click.bs.button.data-api', '[data-toggle^="button"]', function (e) {                                          // 283
      var $btn = $(e.target)                                                                                           // 284
      if (!$btn.hasClass('btn')) $btn = $btn.closest('.btn')                                                           // 285
      Plugin.call($btn, 'toggle')                                                                                      // 286
      if (!($(e.target).is('input[type="radio"]') || $(e.target).is('input[type="checkbox"]'))) e.preventDefault()     // 287
    })                                                                                                                 // 288
    .on('focus.bs.button.data-api blur.bs.button.data-api', '[data-toggle^="button"]', function (e) {                  // 289
      $(e.target).closest('.btn').toggleClass('focus', /^focus(in)?$/.test(e.type))                                    // 290
    })                                                                                                                 // 291
                                                                                                                       // 292
}(jQuery);                                                                                                             // 293
                                                                                                                       // 294
/* ========================================================================                                            // 295
 * Bootstrap: carousel.js v3.3.6                                                                                       // 296
 * http://getbootstrap.com/javascript/#carousel                                                                        // 297
 * ========================================================================                                            // 298
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 299
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 300
 * ======================================================================== */                                         // 301
                                                                                                                       // 302
                                                                                                                       // 303
+function ($) {                                                                                                        // 304
  'use strict';                                                                                                        // 305
                                                                                                                       // 306
  // CAROUSEL CLASS DEFINITION                                                                                         // 307
  // =========================                                                                                         // 308
                                                                                                                       // 309
  var Carousel = function (element, options) {                                                                         // 310
    this.$element    = $(element)                                                                                      // 311
    this.$indicators = this.$element.find('.carousel-indicators')                                                      // 312
    this.options     = options                                                                                         // 313
    this.paused      = null                                                                                            // 314
    this.sliding     = null                                                                                            // 315
    this.interval    = null                                                                                            // 316
    this.$active     = null                                                                                            // 317
    this.$items      = null                                                                                            // 318
                                                                                                                       // 319
    this.options.keyboard && this.$element.on('keydown.bs.carousel', $.proxy(this.keydown, this))                      // 320
                                                                                                                       // 321
    this.options.pause == 'hover' && !('ontouchstart' in document.documentElement) && this.$element                    // 322
      .on('mouseenter.bs.carousel', $.proxy(this.pause, this))                                                         // 323
      .on('mouseleave.bs.carousel', $.proxy(this.cycle, this))                                                         // 324
  }                                                                                                                    // 325
                                                                                                                       // 326
  Carousel.VERSION  = '3.3.6'                                                                                          // 327
                                                                                                                       // 328
  Carousel.TRANSITION_DURATION = 600                                                                                   // 329
                                                                                                                       // 330
  Carousel.DEFAULTS = {                                                                                                // 331
    interval: 5000,                                                                                                    // 332
    pause: 'hover',                                                                                                    // 333
    wrap: true,                                                                                                        // 334
    keyboard: true                                                                                                     // 335
  }                                                                                                                    // 336
                                                                                                                       // 337
  Carousel.prototype.keydown = function (e) {                                                                          // 338
    if (/input|textarea/i.test(e.target.tagName)) return                                                               // 339
    switch (e.which) {                                                                                                 // 340
      case 37: this.prev(); break                                                                                      // 341
      case 39: this.next(); break                                                                                      // 342
      default: return                                                                                                  // 343
    }                                                                                                                  // 344
                                                                                                                       // 345
    e.preventDefault()                                                                                                 // 346
  }                                                                                                                    // 347
                                                                                                                       // 348
  Carousel.prototype.cycle = function (e) {                                                                            // 349
    e || (this.paused = false)                                                                                         // 350
                                                                                                                       // 351
    this.interval && clearInterval(this.interval)                                                                      // 352
                                                                                                                       // 353
    this.options.interval                                                                                              // 354
      && !this.paused                                                                                                  // 355
      && (this.interval = setInterval($.proxy(this.next, this), this.options.interval))                                // 356
                                                                                                                       // 357
    return this                                                                                                        // 358
  }                                                                                                                    // 359
                                                                                                                       // 360
  Carousel.prototype.getItemIndex = function (item) {                                                                  // 361
    this.$items = item.parent().children('.item')                                                                      // 362
    return this.$items.index(item || this.$active)                                                                     // 363
  }                                                                                                                    // 364
                                                                                                                       // 365
  Carousel.prototype.getItemForDirection = function (direction, active) {                                              // 366
    var activeIndex = this.getItemIndex(active)                                                                        // 367
    var willWrap = (direction == 'prev' && activeIndex === 0)                                                          // 368
                || (direction == 'next' && activeIndex == (this.$items.length - 1))                                    // 369
    if (willWrap && !this.options.wrap) return active                                                                  // 370
    var delta = direction == 'prev' ? -1 : 1                                                                           // 371
    var itemIndex = (activeIndex + delta) % this.$items.length                                                         // 372
    return this.$items.eq(itemIndex)                                                                                   // 373
  }                                                                                                                    // 374
                                                                                                                       // 375
  Carousel.prototype.to = function (pos) {                                                                             // 376
    var that        = this                                                                                             // 377
    var activeIndex = this.getItemIndex(this.$active = this.$element.find('.item.active'))                             // 378
                                                                                                                       // 379
    if (pos > (this.$items.length - 1) || pos < 0) return                                                              // 380
                                                                                                                       // 381
    if (this.sliding)       return this.$element.one('slid.bs.carousel', function () { that.to(pos) }) // yes, "slid"  // 382
    if (activeIndex == pos) return this.pause().cycle()                                                                // 383
                                                                                                                       // 384
    return this.slide(pos > activeIndex ? 'next' : 'prev', this.$items.eq(pos))                                        // 385
  }                                                                                                                    // 386
                                                                                                                       // 387
  Carousel.prototype.pause = function (e) {                                                                            // 388
    e || (this.paused = true)                                                                                          // 389
                                                                                                                       // 390
    if (this.$element.find('.next, .prev').length && $.support.transition) {                                           // 391
      this.$element.trigger($.support.transition.end)                                                                  // 392
      this.cycle(true)                                                                                                 // 393
    }                                                                                                                  // 394
                                                                                                                       // 395
    this.interval = clearInterval(this.interval)                                                                       // 396
                                                                                                                       // 397
    return this                                                                                                        // 398
  }                                                                                                                    // 399
                                                                                                                       // 400
  Carousel.prototype.next = function () {                                                                              // 401
    if (this.sliding) return                                                                                           // 402
    return this.slide('next')                                                                                          // 403
  }                                                                                                                    // 404
                                                                                                                       // 405
  Carousel.prototype.prev = function () {                                                                              // 406
    if (this.sliding) return                                                                                           // 407
    return this.slide('prev')                                                                                          // 408
  }                                                                                                                    // 409
                                                                                                                       // 410
  Carousel.prototype.slide = function (type, next) {                                                                   // 411
    var $active   = this.$element.find('.item.active')                                                                 // 412
    var $next     = next || this.getItemForDirection(type, $active)                                                    // 413
    var isCycling = this.interval                                                                                      // 414
    var direction = type == 'next' ? 'left' : 'right'                                                                  // 415
    var that      = this                                                                                               // 416
                                                                                                                       // 417
    if ($next.hasClass('active')) return (this.sliding = false)                                                        // 418
                                                                                                                       // 419
    var relatedTarget = $next[0]                                                                                       // 420
    var slideEvent = $.Event('slide.bs.carousel', {                                                                    // 421
      relatedTarget: relatedTarget,                                                                                    // 422
      direction: direction                                                                                             // 423
    })                                                                                                                 // 424
    this.$element.trigger(slideEvent)                                                                                  // 425
    if (slideEvent.isDefaultPrevented()) return                                                                        // 426
                                                                                                                       // 427
    this.sliding = true                                                                                                // 428
                                                                                                                       // 429
    isCycling && this.pause()                                                                                          // 430
                                                                                                                       // 431
    if (this.$indicators.length) {                                                                                     // 432
      this.$indicators.find('.active').removeClass('active')                                                           // 433
      var $nextIndicator = $(this.$indicators.children()[this.getItemIndex($next)])                                    // 434
      $nextIndicator && $nextIndicator.addClass('active')                                                              // 435
    }                                                                                                                  // 436
                                                                                                                       // 437
    var slidEvent = $.Event('slid.bs.carousel', { relatedTarget: relatedTarget, direction: direction }) // yes, "slid"
    if ($.support.transition && this.$element.hasClass('slide')) {                                                     // 439
      $next.addClass(type)                                                                                             // 440
      $next[0].offsetWidth // force reflow                                                                             // 441
      $active.addClass(direction)                                                                                      // 442
      $next.addClass(direction)                                                                                        // 443
      $active                                                                                                          // 444
        .one('bsTransitionEnd', function () {                                                                          // 445
          $next.removeClass([type, direction].join(' ')).addClass('active')                                            // 446
          $active.removeClass(['active', direction].join(' '))                                                         // 447
          that.sliding = false                                                                                         // 448
          setTimeout(function () {                                                                                     // 449
            that.$element.trigger(slidEvent)                                                                           // 450
          }, 0)                                                                                                        // 451
        })                                                                                                             // 452
        .emulateTransitionEnd(Carousel.TRANSITION_DURATION)                                                            // 453
    } else {                                                                                                           // 454
      $active.removeClass('active')                                                                                    // 455
      $next.addClass('active')                                                                                         // 456
      this.sliding = false                                                                                             // 457
      this.$element.trigger(slidEvent)                                                                                 // 458
    }                                                                                                                  // 459
                                                                                                                       // 460
    isCycling && this.cycle()                                                                                          // 461
                                                                                                                       // 462
    return this                                                                                                        // 463
  }                                                                                                                    // 464
                                                                                                                       // 465
                                                                                                                       // 466
  // CAROUSEL PLUGIN DEFINITION                                                                                        // 467
  // ==========================                                                                                        // 468
                                                                                                                       // 469
  function Plugin(option) {                                                                                            // 470
    return this.each(function () {                                                                                     // 471
      var $this   = $(this)                                                                                            // 472
      var data    = $this.data('bs.carousel')                                                                          // 473
      var options = $.extend({}, Carousel.DEFAULTS, $this.data(), typeof option == 'object' && option)                 // 474
      var action  = typeof option == 'string' ? option : options.slide                                                 // 475
                                                                                                                       // 476
      if (!data) $this.data('bs.carousel', (data = new Carousel(this, options)))                                       // 477
      if (typeof option == 'number') data.to(option)                                                                   // 478
      else if (action) data[action]()                                                                                  // 479
      else if (options.interval) data.pause().cycle()                                                                  // 480
    })                                                                                                                 // 481
  }                                                                                                                    // 482
                                                                                                                       // 483
  var old = $.fn.carousel                                                                                              // 484
                                                                                                                       // 485
  $.fn.carousel             = Plugin                                                                                   // 486
  $.fn.carousel.Constructor = Carousel                                                                                 // 487
                                                                                                                       // 488
                                                                                                                       // 489
  // CAROUSEL NO CONFLICT                                                                                              // 490
  // ====================                                                                                              // 491
                                                                                                                       // 492
  $.fn.carousel.noConflict = function () {                                                                             // 493
    $.fn.carousel = old                                                                                                // 494
    return this                                                                                                        // 495
  }                                                                                                                    // 496
                                                                                                                       // 497
                                                                                                                       // 498
  // CAROUSEL DATA-API                                                                                                 // 499
  // =================                                                                                                 // 500
                                                                                                                       // 501
  var clickHandler = function (e) {                                                                                    // 502
    var href                                                                                                           // 503
    var $this   = $(this)                                                                                              // 504
    var $target = $($this.attr('data-target') || (href = $this.attr('href')) && href.replace(/.*(?=#[^\s]+$)/, '')) // strip for ie7
    if (!$target.hasClass('carousel')) return                                                                          // 506
    var options = $.extend({}, $target.data(), $this.data())                                                           // 507
    var slideIndex = $this.attr('data-slide-to')                                                                       // 508
    if (slideIndex) options.interval = false                                                                           // 509
                                                                                                                       // 510
    Plugin.call($target, options)                                                                                      // 511
                                                                                                                       // 512
    if (slideIndex) {                                                                                                  // 513
      $target.data('bs.carousel').to(slideIndex)                                                                       // 514
    }                                                                                                                  // 515
                                                                                                                       // 516
    e.preventDefault()                                                                                                 // 517
  }                                                                                                                    // 518
                                                                                                                       // 519
  $(document)                                                                                                          // 520
    .on('click.bs.carousel.data-api', '[data-slide]', clickHandler)                                                    // 521
    .on('click.bs.carousel.data-api', '[data-slide-to]', clickHandler)                                                 // 522
                                                                                                                       // 523
  $(window).on('load', function () {                                                                                   // 524
    $('[data-ride="carousel"]').each(function () {                                                                     // 525
      var $carousel = $(this)                                                                                          // 526
      Plugin.call($carousel, $carousel.data())                                                                         // 527
    })                                                                                                                 // 528
  })                                                                                                                   // 529
                                                                                                                       // 530
}(jQuery);                                                                                                             // 531
                                                                                                                       // 532
/* ========================================================================                                            // 533
 * Bootstrap: collapse.js v3.3.6                                                                                       // 534
 * http://getbootstrap.com/javascript/#collapse                                                                        // 535
 * ========================================================================                                            // 536
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 537
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 538
 * ======================================================================== */                                         // 539
                                                                                                                       // 540
                                                                                                                       // 541
+function ($) {                                                                                                        // 542
  'use strict';                                                                                                        // 543
                                                                                                                       // 544
  // COLLAPSE PUBLIC CLASS DEFINITION                                                                                  // 545
  // ================================                                                                                  // 546
                                                                                                                       // 547
  var Collapse = function (element, options) {                                                                         // 548
    this.$element      = $(element)                                                                                    // 549
    this.options       = $.extend({}, Collapse.DEFAULTS, options)                                                      // 550
    this.$trigger      = $('[data-toggle="collapse"][href="#' + element.id + '"],' +                                   // 551
                           '[data-toggle="collapse"][data-target="#' + element.id + '"]')                              // 552
    this.transitioning = null                                                                                          // 553
                                                                                                                       // 554
    if (this.options.parent) {                                                                                         // 555
      this.$parent = this.getParent()                                                                                  // 556
    } else {                                                                                                           // 557
      this.addAriaAndCollapsedClass(this.$element, this.$trigger)                                                      // 558
    }                                                                                                                  // 559
                                                                                                                       // 560
    if (this.options.toggle) this.toggle()                                                                             // 561
  }                                                                                                                    // 562
                                                                                                                       // 563
  Collapse.VERSION  = '3.3.6'                                                                                          // 564
                                                                                                                       // 565
  Collapse.TRANSITION_DURATION = 350                                                                                   // 566
                                                                                                                       // 567
  Collapse.DEFAULTS = {                                                                                                // 568
    toggle: true                                                                                                       // 569
  }                                                                                                                    // 570
                                                                                                                       // 571
  Collapse.prototype.dimension = function () {                                                                         // 572
    var hasWidth = this.$element.hasClass('width')                                                                     // 573
    return hasWidth ? 'width' : 'height'                                                                               // 574
  }                                                                                                                    // 575
                                                                                                                       // 576
  Collapse.prototype.show = function () {                                                                              // 577
    if (this.transitioning || this.$element.hasClass('in')) return                                                     // 578
                                                                                                                       // 579
    var activesData                                                                                                    // 580
    var actives = this.$parent && this.$parent.children('.panel').children('.in, .collapsing')                         // 581
                                                                                                                       // 582
    if (actives && actives.length) {                                                                                   // 583
      activesData = actives.data('bs.collapse')                                                                        // 584
      if (activesData && activesData.transitioning) return                                                             // 585
    }                                                                                                                  // 586
                                                                                                                       // 587
    var startEvent = $.Event('show.bs.collapse')                                                                       // 588
    this.$element.trigger(startEvent)                                                                                  // 589
    if (startEvent.isDefaultPrevented()) return                                                                        // 590
                                                                                                                       // 591
    if (actives && actives.length) {                                                                                   // 592
      Plugin.call(actives, 'hide')                                                                                     // 593
      activesData || actives.data('bs.collapse', null)                                                                 // 594
    }                                                                                                                  // 595
                                                                                                                       // 596
    var dimension = this.dimension()                                                                                   // 597
                                                                                                                       // 598
    this.$element                                                                                                      // 599
      .removeClass('collapse')                                                                                         // 600
      .addClass('collapsing')[dimension](0)                                                                            // 601
      .attr('aria-expanded', true)                                                                                     // 602
                                                                                                                       // 603
    this.$trigger                                                                                                      // 604
      .removeClass('collapsed')                                                                                        // 605
      .attr('aria-expanded', true)                                                                                     // 606
                                                                                                                       // 607
    this.transitioning = 1                                                                                             // 608
                                                                                                                       // 609
    var complete = function () {                                                                                       // 610
      this.$element                                                                                                    // 611
        .removeClass('collapsing')                                                                                     // 612
        .addClass('collapse in')[dimension]('')                                                                        // 613
      this.transitioning = 0                                                                                           // 614
      this.$element                                                                                                    // 615
        .trigger('shown.bs.collapse')                                                                                  // 616
    }                                                                                                                  // 617
                                                                                                                       // 618
    if (!$.support.transition) return complete.call(this)                                                              // 619
                                                                                                                       // 620
    var scrollSize = $.camelCase(['scroll', dimension].join('-'))                                                      // 621
                                                                                                                       // 622
    this.$element                                                                                                      // 623
      .one('bsTransitionEnd', $.proxy(complete, this))                                                                 // 624
      .emulateTransitionEnd(Collapse.TRANSITION_DURATION)[dimension](this.$element[0][scrollSize])                     // 625
  }                                                                                                                    // 626
                                                                                                                       // 627
  Collapse.prototype.hide = function () {                                                                              // 628
    if (this.transitioning || !this.$element.hasClass('in')) return                                                    // 629
                                                                                                                       // 630
    var startEvent = $.Event('hide.bs.collapse')                                                                       // 631
    this.$element.trigger(startEvent)                                                                                  // 632
    if (startEvent.isDefaultPrevented()) return                                                                        // 633
                                                                                                                       // 634
    var dimension = this.dimension()                                                                                   // 635
                                                                                                                       // 636
    this.$element[dimension](this.$element[dimension]())[0].offsetHeight                                               // 637
                                                                                                                       // 638
    this.$element                                                                                                      // 639
      .addClass('collapsing')                                                                                          // 640
      .removeClass('collapse in')                                                                                      // 641
      .attr('aria-expanded', false)                                                                                    // 642
                                                                                                                       // 643
    this.$trigger                                                                                                      // 644
      .addClass('collapsed')                                                                                           // 645
      .attr('aria-expanded', false)                                                                                    // 646
                                                                                                                       // 647
    this.transitioning = 1                                                                                             // 648
                                                                                                                       // 649
    var complete = function () {                                                                                       // 650
      this.transitioning = 0                                                                                           // 651
      this.$element                                                                                                    // 652
        .removeClass('collapsing')                                                                                     // 653
        .addClass('collapse')                                                                                          // 654
        .trigger('hidden.bs.collapse')                                                                                 // 655
    }                                                                                                                  // 656
                                                                                                                       // 657
    if (!$.support.transition) return complete.call(this)                                                              // 658
                                                                                                                       // 659
    this.$element                                                                                                      // 660
      [dimension](0)                                                                                                   // 661
      .one('bsTransitionEnd', $.proxy(complete, this))                                                                 // 662
      .emulateTransitionEnd(Collapse.TRANSITION_DURATION)                                                              // 663
  }                                                                                                                    // 664
                                                                                                                       // 665
  Collapse.prototype.toggle = function () {                                                                            // 666
    this[this.$element.hasClass('in') ? 'hide' : 'show']()                                                             // 667
  }                                                                                                                    // 668
                                                                                                                       // 669
  Collapse.prototype.getParent = function () {                                                                         // 670
    return $(this.options.parent)                                                                                      // 671
      .find('[data-toggle="collapse"][data-parent="' + this.options.parent + '"]')                                     // 672
      .each($.proxy(function (i, element) {                                                                            // 673
        var $element = $(element)                                                                                      // 674
        this.addAriaAndCollapsedClass(getTargetFromTrigger($element), $element)                                        // 675
      }, this))                                                                                                        // 676
      .end()                                                                                                           // 677
  }                                                                                                                    // 678
                                                                                                                       // 679
  Collapse.prototype.addAriaAndCollapsedClass = function ($element, $trigger) {                                        // 680
    var isOpen = $element.hasClass('in')                                                                               // 681
                                                                                                                       // 682
    $element.attr('aria-expanded', isOpen)                                                                             // 683
    $trigger                                                                                                           // 684
      .toggleClass('collapsed', !isOpen)                                                                               // 685
      .attr('aria-expanded', isOpen)                                                                                   // 686
  }                                                                                                                    // 687
                                                                                                                       // 688
  function getTargetFromTrigger($trigger) {                                                                            // 689
    var href                                                                                                           // 690
    var target = $trigger.attr('data-target')                                                                          // 691
      || (href = $trigger.attr('href')) && href.replace(/.*(?=#[^\s]+$)/, '') // strip for ie7                         // 692
                                                                                                                       // 693
    return $(target)                                                                                                   // 694
  }                                                                                                                    // 695
                                                                                                                       // 696
                                                                                                                       // 697
  // COLLAPSE PLUGIN DEFINITION                                                                                        // 698
  // ==========================                                                                                        // 699
                                                                                                                       // 700
  function Plugin(option) {                                                                                            // 701
    return this.each(function () {                                                                                     // 702
      var $this   = $(this)                                                                                            // 703
      var data    = $this.data('bs.collapse')                                                                          // 704
      var options = $.extend({}, Collapse.DEFAULTS, $this.data(), typeof option == 'object' && option)                 // 705
                                                                                                                       // 706
      if (!data && options.toggle && /show|hide/.test(option)) options.toggle = false                                  // 707
      if (!data) $this.data('bs.collapse', (data = new Collapse(this, options)))                                       // 708
      if (typeof option == 'string') data[option]()                                                                    // 709
    })                                                                                                                 // 710
  }                                                                                                                    // 711
                                                                                                                       // 712
  var old = $.fn.collapse                                                                                              // 713
                                                                                                                       // 714
  $.fn.collapse             = Plugin                                                                                   // 715
  $.fn.collapse.Constructor = Collapse                                                                                 // 716
                                                                                                                       // 717
                                                                                                                       // 718
  // COLLAPSE NO CONFLICT                                                                                              // 719
  // ====================                                                                                              // 720
                                                                                                                       // 721
  $.fn.collapse.noConflict = function () {                                                                             // 722
    $.fn.collapse = old                                                                                                // 723
    return this                                                                                                        // 724
  }                                                                                                                    // 725
                                                                                                                       // 726
                                                                                                                       // 727
  // COLLAPSE DATA-API                                                                                                 // 728
  // =================                                                                                                 // 729
                                                                                                                       // 730
  $(document).on('click.bs.collapse.data-api', '[data-toggle="collapse"]', function (e) {                              // 731
    var $this   = $(this)                                                                                              // 732
                                                                                                                       // 733
    if (!$this.attr('data-target')) e.preventDefault()                                                                 // 734
                                                                                                                       // 735
    var $target = getTargetFromTrigger($this)                                                                          // 736
    var data    = $target.data('bs.collapse')                                                                          // 737
    var option  = data ? 'toggle' : $this.data()                                                                       // 738
                                                                                                                       // 739
    Plugin.call($target, option)                                                                                       // 740
  })                                                                                                                   // 741
                                                                                                                       // 742
}(jQuery);                                                                                                             // 743
                                                                                                                       // 744
/* ========================================================================                                            // 745
 * Bootstrap: dropdown.js v3.3.6                                                                                       // 746
 * http://getbootstrap.com/javascript/#dropdowns                                                                       // 747
 * ========================================================================                                            // 748
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 749
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 750
 * ======================================================================== */                                         // 751
                                                                                                                       // 752
                                                                                                                       // 753
+function ($) {                                                                                                        // 754
  'use strict';                                                                                                        // 755
                                                                                                                       // 756
  // DROPDOWN CLASS DEFINITION                                                                                         // 757
  // =========================                                                                                         // 758
                                                                                                                       // 759
  var backdrop = '.dropdown-backdrop'                                                                                  // 760
  var toggle   = '[data-toggle="dropdown"]'                                                                            // 761
  var Dropdown = function (element) {                                                                                  // 762
    $(element).on('click.bs.dropdown', this.toggle)                                                                    // 763
  }                                                                                                                    // 764
                                                                                                                       // 765
  Dropdown.VERSION = '3.3.6'                                                                                           // 766
                                                                                                                       // 767
  function getParent($this) {                                                                                          // 768
    var selector = $this.attr('data-target')                                                                           // 769
                                                                                                                       // 770
    if (!selector) {                                                                                                   // 771
      selector = $this.attr('href')                                                                                    // 772
      selector = selector && /#[A-Za-z]/.test(selector) && selector.replace(/.*(?=#[^\s]*$)/, '') // strip for ie7     // 773
    }                                                                                                                  // 774
                                                                                                                       // 775
    var $parent = selector && $(selector)                                                                              // 776
                                                                                                                       // 777
    return $parent && $parent.length ? $parent : $this.parent()                                                        // 778
  }                                                                                                                    // 779
                                                                                                                       // 780
  function clearMenus(e) {                                                                                             // 781
    if (e && e.which === 3) return                                                                                     // 782
    $(backdrop).remove()                                                                                               // 783
    $(toggle).each(function () {                                                                                       // 784
      var $this         = $(this)                                                                                      // 785
      var $parent       = getParent($this)                                                                             // 786
      var relatedTarget = { relatedTarget: this }                                                                      // 787
                                                                                                                       // 788
      if (!$parent.hasClass('open')) return                                                                            // 789
                                                                                                                       // 790
      if (e && e.type == 'click' && /input|textarea/i.test(e.target.tagName) && $.contains($parent[0], e.target)) return
                                                                                                                       // 792
      $parent.trigger(e = $.Event('hide.bs.dropdown', relatedTarget))                                                  // 793
                                                                                                                       // 794
      if (e.isDefaultPrevented()) return                                                                               // 795
                                                                                                                       // 796
      $this.attr('aria-expanded', 'false')                                                                             // 797
      $parent.removeClass('open').trigger($.Event('hidden.bs.dropdown', relatedTarget))                                // 798
    })                                                                                                                 // 799
  }                                                                                                                    // 800
                                                                                                                       // 801
  Dropdown.prototype.toggle = function (e) {                                                                           // 802
    var $this = $(this)                                                                                                // 803
                                                                                                                       // 804
    if ($this.is('.disabled, :disabled')) return                                                                       // 805
                                                                                                                       // 806
    var $parent  = getParent($this)                                                                                    // 807
    var isActive = $parent.hasClass('open')                                                                            // 808
                                                                                                                       // 809
    clearMenus()                                                                                                       // 810
                                                                                                                       // 811
    if (!isActive) {                                                                                                   // 812
      if ('ontouchstart' in document.documentElement && !$parent.closest('.navbar-nav').length) {                      // 813
        // if mobile we use a backdrop because click events don't delegate                                             // 814
        $(document.createElement('div'))                                                                               // 815
          .addClass('dropdown-backdrop')                                                                               // 816
          .insertAfter($(this))                                                                                        // 817
          .on('click', clearMenus)                                                                                     // 818
      }                                                                                                                // 819
                                                                                                                       // 820
      var relatedTarget = { relatedTarget: this }                                                                      // 821
      $parent.trigger(e = $.Event('show.bs.dropdown', relatedTarget))                                                  // 822
                                                                                                                       // 823
      if (e.isDefaultPrevented()) return                                                                               // 824
                                                                                                                       // 825
      $this                                                                                                            // 826
        .trigger('focus')                                                                                              // 827
        .attr('aria-expanded', 'true')                                                                                 // 828
                                                                                                                       // 829
      $parent                                                                                                          // 830
        .toggleClass('open')                                                                                           // 831
        .trigger($.Event('shown.bs.dropdown', relatedTarget))                                                          // 832
    }                                                                                                                  // 833
                                                                                                                       // 834
    return false                                                                                                       // 835
  }                                                                                                                    // 836
                                                                                                                       // 837
  Dropdown.prototype.keydown = function (e) {                                                                          // 838
    if (!/(38|40|27|32)/.test(e.which) || /input|textarea/i.test(e.target.tagName)) return                             // 839
                                                                                                                       // 840
    var $this = $(this)                                                                                                // 841
                                                                                                                       // 842
    e.preventDefault()                                                                                                 // 843
    e.stopPropagation()                                                                                                // 844
                                                                                                                       // 845
    if ($this.is('.disabled, :disabled')) return                                                                       // 846
                                                                                                                       // 847
    var $parent  = getParent($this)                                                                                    // 848
    var isActive = $parent.hasClass('open')                                                                            // 849
                                                                                                                       // 850
    if (!isActive && e.which != 27 || isActive && e.which == 27) {                                                     // 851
      if (e.which == 27) $parent.find(toggle).trigger('focus')                                                         // 852
      return $this.trigger('click')                                                                                    // 853
    }                                                                                                                  // 854
                                                                                                                       // 855
    var desc = ' li:not(.disabled):visible a'                                                                          // 856
    var $items = $parent.find('.dropdown-menu' + desc)                                                                 // 857
                                                                                                                       // 858
    if (!$items.length) return                                                                                         // 859
                                                                                                                       // 860
    var index = $items.index(e.target)                                                                                 // 861
                                                                                                                       // 862
    if (e.which == 38 && index > 0)                 index--         // up                                              // 863
    if (e.which == 40 && index < $items.length - 1) index++         // down                                            // 864
    if (!~index)                                    index = 0                                                          // 865
                                                                                                                       // 866
    $items.eq(index).trigger('focus')                                                                                  // 867
  }                                                                                                                    // 868
                                                                                                                       // 869
                                                                                                                       // 870
  // DROPDOWN PLUGIN DEFINITION                                                                                        // 871
  // ==========================                                                                                        // 872
                                                                                                                       // 873
  function Plugin(option) {                                                                                            // 874
    return this.each(function () {                                                                                     // 875
      var $this = $(this)                                                                                              // 876
      var data  = $this.data('bs.dropdown')                                                                            // 877
                                                                                                                       // 878
      if (!data) $this.data('bs.dropdown', (data = new Dropdown(this)))                                                // 879
      if (typeof option == 'string') data[option].call($this)                                                          // 880
    })                                                                                                                 // 881
  }                                                                                                                    // 882
                                                                                                                       // 883
  var old = $.fn.dropdown                                                                                              // 884
                                                                                                                       // 885
  $.fn.dropdown             = Plugin                                                                                   // 886
  $.fn.dropdown.Constructor = Dropdown                                                                                 // 887
                                                                                                                       // 888
                                                                                                                       // 889
  // DROPDOWN NO CONFLICT                                                                                              // 890
  // ====================                                                                                              // 891
                                                                                                                       // 892
  $.fn.dropdown.noConflict = function () {                                                                             // 893
    $.fn.dropdown = old                                                                                                // 894
    return this                                                                                                        // 895
  }                                                                                                                    // 896
                                                                                                                       // 897
                                                                                                                       // 898
  // APPLY TO STANDARD DROPDOWN ELEMENTS                                                                               // 899
  // ===================================                                                                               // 900
                                                                                                                       // 901
  $(document)                                                                                                          // 902
    .on('click.bs.dropdown.data-api', clearMenus)                                                                      // 903
    .on('click.bs.dropdown.data-api', '.dropdown form', function (e) { e.stopPropagation() })                          // 904
    .on('click.bs.dropdown.data-api', toggle, Dropdown.prototype.toggle)                                               // 905
    .on('keydown.bs.dropdown.data-api', toggle, Dropdown.prototype.keydown)                                            // 906
    .on('keydown.bs.dropdown.data-api', '.dropdown-menu', Dropdown.prototype.keydown)                                  // 907
                                                                                                                       // 908
}(jQuery);                                                                                                             // 909
                                                                                                                       // 910
/* ========================================================================                                            // 911
 * Bootstrap: modal.js v3.3.6                                                                                          // 912
 * http://getbootstrap.com/javascript/#modals                                                                          // 913
 * ========================================================================                                            // 914
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 915
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 916
 * ======================================================================== */                                         // 917
                                                                                                                       // 918
                                                                                                                       // 919
+function ($) {                                                                                                        // 920
  'use strict';                                                                                                        // 921
                                                                                                                       // 922
  // MODAL CLASS DEFINITION                                                                                            // 923
  // ======================                                                                                            // 924
                                                                                                                       // 925
  var Modal = function (element, options) {                                                                            // 926
    this.options             = options                                                                                 // 927
    this.$body               = $(document.body)                                                                        // 928
    this.$element            = $(element)                                                                              // 929
    this.$dialog             = this.$element.find('.modal-dialog')                                                     // 930
    this.$backdrop           = null                                                                                    // 931
    this.isShown             = null                                                                                    // 932
    this.originalBodyPad     = null                                                                                    // 933
    this.scrollbarWidth      = 0                                                                                       // 934
    this.ignoreBackdropClick = false                                                                                   // 935
                                                                                                                       // 936
    if (this.options.remote) {                                                                                         // 937
      this.$element                                                                                                    // 938
        .find('.modal-content')                                                                                        // 939
        .load(this.options.remote, $.proxy(function () {                                                               // 940
          this.$element.trigger('loaded.bs.modal')                                                                     // 941
        }, this))                                                                                                      // 942
    }                                                                                                                  // 943
  }                                                                                                                    // 944
                                                                                                                       // 945
  Modal.VERSION  = '3.3.6'                                                                                             // 946
                                                                                                                       // 947
  Modal.TRANSITION_DURATION = 300                                                                                      // 948
  Modal.BACKDROP_TRANSITION_DURATION = 150                                                                             // 949
                                                                                                                       // 950
  Modal.DEFAULTS = {                                                                                                   // 951
    backdrop: true,                                                                                                    // 952
    keyboard: true,                                                                                                    // 953
    show: true                                                                                                         // 954
  }                                                                                                                    // 955
                                                                                                                       // 956
  Modal.prototype.toggle = function (_relatedTarget) {                                                                 // 957
    return this.isShown ? this.hide() : this.show(_relatedTarget)                                                      // 958
  }                                                                                                                    // 959
                                                                                                                       // 960
  Modal.prototype.show = function (_relatedTarget) {                                                                   // 961
    var that = this                                                                                                    // 962
    var e    = $.Event('show.bs.modal', { relatedTarget: _relatedTarget })                                             // 963
                                                                                                                       // 964
    this.$element.trigger(e)                                                                                           // 965
                                                                                                                       // 966
    if (this.isShown || e.isDefaultPrevented()) return                                                                 // 967
                                                                                                                       // 968
    this.isShown = true                                                                                                // 969
                                                                                                                       // 970
    this.checkScrollbar()                                                                                              // 971
    this.setScrollbar()                                                                                                // 972
    this.$body.addClass('modal-open')                                                                                  // 973
                                                                                                                       // 974
    this.escape()                                                                                                      // 975
    this.resize()                                                                                                      // 976
                                                                                                                       // 977
    this.$element.on('click.dismiss.bs.modal', '[data-dismiss="modal"]', $.proxy(this.hide, this))                     // 978
                                                                                                                       // 979
    this.$dialog.on('mousedown.dismiss.bs.modal', function () {                                                        // 980
      that.$element.one('mouseup.dismiss.bs.modal', function (e) {                                                     // 981
        if ($(e.target).is(that.$element)) that.ignoreBackdropClick = true                                             // 982
      })                                                                                                               // 983
    })                                                                                                                 // 984
                                                                                                                       // 985
    this.backdrop(function () {                                                                                        // 986
      var transition = $.support.transition && that.$element.hasClass('fade')                                          // 987
                                                                                                                       // 988
      if (!that.$element.parent().length) {                                                                            // 989
        that.$element.appendTo(that.$body) // don't move modals dom position                                           // 990
      }                                                                                                                // 991
                                                                                                                       // 992
      that.$element                                                                                                    // 993
        .show()                                                                                                        // 994
        .scrollTop(0)                                                                                                  // 995
                                                                                                                       // 996
      that.adjustDialog()                                                                                              // 997
                                                                                                                       // 998
      if (transition) {                                                                                                // 999
        that.$element[0].offsetWidth // force reflow                                                                   // 1000
      }                                                                                                                // 1001
                                                                                                                       // 1002
      that.$element.addClass('in')                                                                                     // 1003
                                                                                                                       // 1004
      that.enforceFocus()                                                                                              // 1005
                                                                                                                       // 1006
      var e = $.Event('shown.bs.modal', { relatedTarget: _relatedTarget })                                             // 1007
                                                                                                                       // 1008
      transition ?                                                                                                     // 1009
        that.$dialog // wait for modal to slide in                                                                     // 1010
          .one('bsTransitionEnd', function () {                                                                        // 1011
            that.$element.trigger('focus').trigger(e)                                                                  // 1012
          })                                                                                                           // 1013
          .emulateTransitionEnd(Modal.TRANSITION_DURATION) :                                                           // 1014
        that.$element.trigger('focus').trigger(e)                                                                      // 1015
    })                                                                                                                 // 1016
  }                                                                                                                    // 1017
                                                                                                                       // 1018
  Modal.prototype.hide = function (e) {                                                                                // 1019
    if (e) e.preventDefault()                                                                                          // 1020
                                                                                                                       // 1021
    e = $.Event('hide.bs.modal')                                                                                       // 1022
                                                                                                                       // 1023
    this.$element.trigger(e)                                                                                           // 1024
                                                                                                                       // 1025
    if (!this.isShown || e.isDefaultPrevented()) return                                                                // 1026
                                                                                                                       // 1027
    this.isShown = false                                                                                               // 1028
                                                                                                                       // 1029
    this.escape()                                                                                                      // 1030
    this.resize()                                                                                                      // 1031
                                                                                                                       // 1032
    $(document).off('focusin.bs.modal')                                                                                // 1033
                                                                                                                       // 1034
    this.$element                                                                                                      // 1035
      .removeClass('in')                                                                                               // 1036
      .off('click.dismiss.bs.modal')                                                                                   // 1037
      .off('mouseup.dismiss.bs.modal')                                                                                 // 1038
                                                                                                                       // 1039
    this.$dialog.off('mousedown.dismiss.bs.modal')                                                                     // 1040
                                                                                                                       // 1041
    $.support.transition && this.$element.hasClass('fade') ?                                                           // 1042
      this.$element                                                                                                    // 1043
        .one('bsTransitionEnd', $.proxy(this.hideModal, this))                                                         // 1044
        .emulateTransitionEnd(Modal.TRANSITION_DURATION) :                                                             // 1045
      this.hideModal()                                                                                                 // 1046
  }                                                                                                                    // 1047
                                                                                                                       // 1048
  Modal.prototype.enforceFocus = function () {                                                                         // 1049
    $(document)                                                                                                        // 1050
      .off('focusin.bs.modal') // guard against infinite focus loop                                                    // 1051
      .on('focusin.bs.modal', $.proxy(function (e) {                                                                   // 1052
        if (this.$element[0] !== e.target && !this.$element.has(e.target).length) {                                    // 1053
          this.$element.trigger('focus')                                                                               // 1054
        }                                                                                                              // 1055
      }, this))                                                                                                        // 1056
  }                                                                                                                    // 1057
                                                                                                                       // 1058
  Modal.prototype.escape = function () {                                                                               // 1059
    if (this.isShown && this.options.keyboard) {                                                                       // 1060
      this.$element.on('keydown.dismiss.bs.modal', $.proxy(function (e) {                                              // 1061
        e.which == 27 && this.hide()                                                                                   // 1062
      }, this))                                                                                                        // 1063
    } else if (!this.isShown) {                                                                                        // 1064
      this.$element.off('keydown.dismiss.bs.modal')                                                                    // 1065
    }                                                                                                                  // 1066
  }                                                                                                                    // 1067
                                                                                                                       // 1068
  Modal.prototype.resize = function () {                                                                               // 1069
    if (this.isShown) {                                                                                                // 1070
      $(window).on('resize.bs.modal', $.proxy(this.handleUpdate, this))                                                // 1071
    } else {                                                                                                           // 1072
      $(window).off('resize.bs.modal')                                                                                 // 1073
    }                                                                                                                  // 1074
  }                                                                                                                    // 1075
                                                                                                                       // 1076
  Modal.prototype.hideModal = function () {                                                                            // 1077
    var that = this                                                                                                    // 1078
    this.$element.hide()                                                                                               // 1079
    this.backdrop(function () {                                                                                        // 1080
      that.$body.removeClass('modal-open')                                                                             // 1081
      that.resetAdjustments()                                                                                          // 1082
      that.resetScrollbar()                                                                                            // 1083
      that.$element.trigger('hidden.bs.modal')                                                                         // 1084
    })                                                                                                                 // 1085
  }                                                                                                                    // 1086
                                                                                                                       // 1087
  Modal.prototype.removeBackdrop = function () {                                                                       // 1088
    this.$backdrop && this.$backdrop.remove()                                                                          // 1089
    this.$backdrop = null                                                                                              // 1090
  }                                                                                                                    // 1091
                                                                                                                       // 1092
  Modal.prototype.backdrop = function (callback) {                                                                     // 1093
    var that = this                                                                                                    // 1094
    var animate = this.$element.hasClass('fade') ? 'fade' : ''                                                         // 1095
                                                                                                                       // 1096
    if (this.isShown && this.options.backdrop) {                                                                       // 1097
      var doAnimate = $.support.transition && animate                                                                  // 1098
                                                                                                                       // 1099
      this.$backdrop = $(document.createElement('div'))                                                                // 1100
        .addClass('modal-backdrop ' + animate)                                                                         // 1101
        .appendTo(this.$body)                                                                                          // 1102
                                                                                                                       // 1103
      this.$element.on('click.dismiss.bs.modal', $.proxy(function (e) {                                                // 1104
        if (this.ignoreBackdropClick) {                                                                                // 1105
          this.ignoreBackdropClick = false                                                                             // 1106
          return                                                                                                       // 1107
        }                                                                                                              // 1108
        if (e.target !== e.currentTarget) return                                                                       // 1109
        this.options.backdrop == 'static'                                                                              // 1110
          ? this.$element[0].focus()                                                                                   // 1111
          : this.hide()                                                                                                // 1112
      }, this))                                                                                                        // 1113
                                                                                                                       // 1114
      if (doAnimate) this.$backdrop[0].offsetWidth // force reflow                                                     // 1115
                                                                                                                       // 1116
      this.$backdrop.addClass('in')                                                                                    // 1117
                                                                                                                       // 1118
      if (!callback) return                                                                                            // 1119
                                                                                                                       // 1120
      doAnimate ?                                                                                                      // 1121
        this.$backdrop                                                                                                 // 1122
          .one('bsTransitionEnd', callback)                                                                            // 1123
          .emulateTransitionEnd(Modal.BACKDROP_TRANSITION_DURATION) :                                                  // 1124
        callback()                                                                                                     // 1125
                                                                                                                       // 1126
    } else if (!this.isShown && this.$backdrop) {                                                                      // 1127
      this.$backdrop.removeClass('in')                                                                                 // 1128
                                                                                                                       // 1129
      var callbackRemove = function () {                                                                               // 1130
        that.removeBackdrop()                                                                                          // 1131
        callback && callback()                                                                                         // 1132
      }                                                                                                                // 1133
      $.support.transition && this.$element.hasClass('fade') ?                                                         // 1134
        this.$backdrop                                                                                                 // 1135
          .one('bsTransitionEnd', callbackRemove)                                                                      // 1136
          .emulateTransitionEnd(Modal.BACKDROP_TRANSITION_DURATION) :                                                  // 1137
        callbackRemove()                                                                                               // 1138
                                                                                                                       // 1139
    } else if (callback) {                                                                                             // 1140
      callback()                                                                                                       // 1141
    }                                                                                                                  // 1142
  }                                                                                                                    // 1143
                                                                                                                       // 1144
  // these following methods are used to handle overflowing modals                                                     // 1145
                                                                                                                       // 1146
  Modal.prototype.handleUpdate = function () {                                                                         // 1147
    this.adjustDialog()                                                                                                // 1148
  }                                                                                                                    // 1149
                                                                                                                       // 1150
  Modal.prototype.adjustDialog = function () {                                                                         // 1151
    var modalIsOverflowing = this.$element[0].scrollHeight > document.documentElement.clientHeight                     // 1152
                                                                                                                       // 1153
    this.$element.css({                                                                                                // 1154
      paddingLeft:  !this.bodyIsOverflowing && modalIsOverflowing ? this.scrollbarWidth : '',                          // 1155
      paddingRight: this.bodyIsOverflowing && !modalIsOverflowing ? this.scrollbarWidth : ''                           // 1156
    })                                                                                                                 // 1157
  }                                                                                                                    // 1158
                                                                                                                       // 1159
  Modal.prototype.resetAdjustments = function () {                                                                     // 1160
    this.$element.css({                                                                                                // 1161
      paddingLeft: '',                                                                                                 // 1162
      paddingRight: ''                                                                                                 // 1163
    })                                                                                                                 // 1164
  }                                                                                                                    // 1165
                                                                                                                       // 1166
  Modal.prototype.checkScrollbar = function () {                                                                       // 1167
    var fullWindowWidth = window.innerWidth                                                                            // 1168
    if (!fullWindowWidth) { // workaround for missing window.innerWidth in IE8                                         // 1169
      var documentElementRect = document.documentElement.getBoundingClientRect()                                       // 1170
      fullWindowWidth = documentElementRect.right - Math.abs(documentElementRect.left)                                 // 1171
    }                                                                                                                  // 1172
    this.bodyIsOverflowing = document.body.clientWidth < fullWindowWidth                                               // 1173
    this.scrollbarWidth = this.measureScrollbar()                                                                      // 1174
  }                                                                                                                    // 1175
                                                                                                                       // 1176
  Modal.prototype.setScrollbar = function () {                                                                         // 1177
    var bodyPad = parseInt((this.$body.css('padding-right') || 0), 10)                                                 // 1178
    this.originalBodyPad = document.body.style.paddingRight || ''                                                      // 1179
    if (this.bodyIsOverflowing) this.$body.css('padding-right', bodyPad + this.scrollbarWidth)                         // 1180
  }                                                                                                                    // 1181
                                                                                                                       // 1182
  Modal.prototype.resetScrollbar = function () {                                                                       // 1183
    this.$body.css('padding-right', this.originalBodyPad)                                                              // 1184
  }                                                                                                                    // 1185
                                                                                                                       // 1186
  Modal.prototype.measureScrollbar = function () { // thx walsh                                                        // 1187
    var scrollDiv = document.createElement('div')                                                                      // 1188
    scrollDiv.className = 'modal-scrollbar-measure'                                                                    // 1189
    this.$body.append(scrollDiv)                                                                                       // 1190
    var scrollbarWidth = scrollDiv.offsetWidth - scrollDiv.clientWidth                                                 // 1191
    this.$body[0].removeChild(scrollDiv)                                                                               // 1192
    return scrollbarWidth                                                                                              // 1193
  }                                                                                                                    // 1194
                                                                                                                       // 1195
                                                                                                                       // 1196
  // MODAL PLUGIN DEFINITION                                                                                           // 1197
  // =======================                                                                                           // 1198
                                                                                                                       // 1199
  function Plugin(option, _relatedTarget) {                                                                            // 1200
    return this.each(function () {                                                                                     // 1201
      var $this   = $(this)                                                                                            // 1202
      var data    = $this.data('bs.modal')                                                                             // 1203
      var options = $.extend({}, Modal.DEFAULTS, $this.data(), typeof option == 'object' && option)                    // 1204
                                                                                                                       // 1205
      if (!data) $this.data('bs.modal', (data = new Modal(this, options)))                                             // 1206
      if (typeof option == 'string') data[option](_relatedTarget)                                                      // 1207
      else if (options.show) data.show(_relatedTarget)                                                                 // 1208
    })                                                                                                                 // 1209
  }                                                                                                                    // 1210
                                                                                                                       // 1211
  var old = $.fn.modal                                                                                                 // 1212
                                                                                                                       // 1213
  $.fn.modal             = Plugin                                                                                      // 1214
  $.fn.modal.Constructor = Modal                                                                                       // 1215
                                                                                                                       // 1216
                                                                                                                       // 1217
  // MODAL NO CONFLICT                                                                                                 // 1218
  // =================                                                                                                 // 1219
                                                                                                                       // 1220
  $.fn.modal.noConflict = function () {                                                                                // 1221
    $.fn.modal = old                                                                                                   // 1222
    return this                                                                                                        // 1223
  }                                                                                                                    // 1224
                                                                                                                       // 1225
                                                                                                                       // 1226
  // MODAL DATA-API                                                                                                    // 1227
  // ==============                                                                                                    // 1228
                                                                                                                       // 1229
  $(document).on('click.bs.modal.data-api', '[data-toggle="modal"]', function (e) {                                    // 1230
    var $this   = $(this)                                                                                              // 1231
    var href    = $this.attr('href')                                                                                   // 1232
    var $target = $($this.attr('data-target') || (href && href.replace(/.*(?=#[^\s]+$)/, ''))) // strip for ie7        // 1233
    var option  = $target.data('bs.modal') ? 'toggle' : $.extend({ remote: !/#/.test(href) && href }, $target.data(), $this.data())
                                                                                                                       // 1235
    if ($this.is('a')) e.preventDefault()                                                                              // 1236
                                                                                                                       // 1237
    $target.one('show.bs.modal', function (showEvent) {                                                                // 1238
      if (showEvent.isDefaultPrevented()) return // only register focus restorer if modal will actually get shown      // 1239
      $target.one('hidden.bs.modal', function () {                                                                     // 1240
        $this.is(':visible') && $this.trigger('focus')                                                                 // 1241
      })                                                                                                               // 1242
    })                                                                                                                 // 1243
    Plugin.call($target, option, this)                                                                                 // 1244
  })                                                                                                                   // 1245
                                                                                                                       // 1246
}(jQuery);                                                                                                             // 1247
                                                                                                                       // 1248
/* ========================================================================                                            // 1249
 * Bootstrap: tooltip.js v3.3.6                                                                                        // 1250
 * http://getbootstrap.com/javascript/#tooltip                                                                         // 1251
 * Inspired by the original jQuery.tipsy by Jason Frame                                                                // 1252
 * ========================================================================                                            // 1253
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 1254
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 1255
 * ======================================================================== */                                         // 1256
                                                                                                                       // 1257
                                                                                                                       // 1258
+function ($) {                                                                                                        // 1259
  'use strict';                                                                                                        // 1260
                                                                                                                       // 1261
  // TOOLTIP PUBLIC CLASS DEFINITION                                                                                   // 1262
  // ===============================                                                                                   // 1263
                                                                                                                       // 1264
  var Tooltip = function (element, options) {                                                                          // 1265
    this.type       = null                                                                                             // 1266
    this.options    = null                                                                                             // 1267
    this.enabled    = null                                                                                             // 1268
    this.timeout    = null                                                                                             // 1269
    this.hoverState = null                                                                                             // 1270
    this.$element   = null                                                                                             // 1271
    this.inState    = null                                                                                             // 1272
                                                                                                                       // 1273
    this.init('tooltip', element, options)                                                                             // 1274
  }                                                                                                                    // 1275
                                                                                                                       // 1276
  Tooltip.VERSION  = '3.3.6'                                                                                           // 1277
                                                                                                                       // 1278
  Tooltip.TRANSITION_DURATION = 150                                                                                    // 1279
                                                                                                                       // 1280
  Tooltip.DEFAULTS = {                                                                                                 // 1281
    animation: true,                                                                                                   // 1282
    placement: 'top',                                                                                                  // 1283
    selector: false,                                                                                                   // 1284
    template: '<div class="tooltip" role="tooltip"><div class="tooltip-arrow"></div><div class="tooltip-inner"></div></div>',
    trigger: 'hover focus',                                                                                            // 1286
    title: '',                                                                                                         // 1287
    delay: 0,                                                                                                          // 1288
    html: false,                                                                                                       // 1289
    container: false,                                                                                                  // 1290
    viewport: {                                                                                                        // 1291
      selector: 'body',                                                                                                // 1292
      padding: 0                                                                                                       // 1293
    }                                                                                                                  // 1294
  }                                                                                                                    // 1295
                                                                                                                       // 1296
  Tooltip.prototype.init = function (type, element, options) {                                                         // 1297
    this.enabled   = true                                                                                              // 1298
    this.type      = type                                                                                              // 1299
    this.$element  = $(element)                                                                                        // 1300
    this.options   = this.getOptions(options)                                                                          // 1301
    this.$viewport = this.options.viewport && $($.isFunction(this.options.viewport) ? this.options.viewport.call(this, this.$element) : (this.options.viewport.selector || this.options.viewport))
    this.inState   = { click: false, hover: false, focus: false }                                                      // 1303
                                                                                                                       // 1304
    if (this.$element[0] instanceof document.constructor && !this.options.selector) {                                  // 1305
      throw new Error('`selector` option must be specified when initializing ' + this.type + ' on the window.document object!')
    }                                                                                                                  // 1307
                                                                                                                       // 1308
    var triggers = this.options.trigger.split(' ')                                                                     // 1309
                                                                                                                       // 1310
    for (var i = triggers.length; i--;) {                                                                              // 1311
      var trigger = triggers[i]                                                                                        // 1312
                                                                                                                       // 1313
      if (trigger == 'click') {                                                                                        // 1314
        this.$element.on('click.' + this.type, this.options.selector, $.proxy(this.toggle, this))                      // 1315
      } else if (trigger != 'manual') {                                                                                // 1316
        var eventIn  = trigger == 'hover' ? 'mouseenter' : 'focusin'                                                   // 1317
        var eventOut = trigger == 'hover' ? 'mouseleave' : 'focusout'                                                  // 1318
                                                                                                                       // 1319
        this.$element.on(eventIn  + '.' + this.type, this.options.selector, $.proxy(this.enter, this))                 // 1320
        this.$element.on(eventOut + '.' + this.type, this.options.selector, $.proxy(this.leave, this))                 // 1321
      }                                                                                                                // 1322
    }                                                                                                                  // 1323
                                                                                                                       // 1324
    this.options.selector ?                                                                                            // 1325
      (this._options = $.extend({}, this.options, { trigger: 'manual', selector: '' })) :                              // 1326
      this.fixTitle()                                                                                                  // 1327
  }                                                                                                                    // 1328
                                                                                                                       // 1329
  Tooltip.prototype.getDefaults = function () {                                                                        // 1330
    return Tooltip.DEFAULTS                                                                                            // 1331
  }                                                                                                                    // 1332
                                                                                                                       // 1333
  Tooltip.prototype.getOptions = function (options) {                                                                  // 1334
    options = $.extend({}, this.getDefaults(), this.$element.data(), options)                                          // 1335
                                                                                                                       // 1336
    if (options.delay && typeof options.delay == 'number') {                                                           // 1337
      options.delay = {                                                                                                // 1338
        show: options.delay,                                                                                           // 1339
        hide: options.delay                                                                                            // 1340
      }                                                                                                                // 1341
    }                                                                                                                  // 1342
                                                                                                                       // 1343
    return options                                                                                                     // 1344
  }                                                                                                                    // 1345
                                                                                                                       // 1346
  Tooltip.prototype.getDelegateOptions = function () {                                                                 // 1347
    var options  = {}                                                                                                  // 1348
    var defaults = this.getDefaults()                                                                                  // 1349
                                                                                                                       // 1350
    this._options && $.each(this._options, function (key, value) {                                                     // 1351
      if (defaults[key] != value) options[key] = value                                                                 // 1352
    })                                                                                                                 // 1353
                                                                                                                       // 1354
    return options                                                                                                     // 1355
  }                                                                                                                    // 1356
                                                                                                                       // 1357
  Tooltip.prototype.enter = function (obj) {                                                                           // 1358
    var self = obj instanceof this.constructor ?                                                                       // 1359
      obj : $(obj.currentTarget).data('bs.' + this.type)                                                               // 1360
                                                                                                                       // 1361
    if (!self) {                                                                                                       // 1362
      self = new this.constructor(obj.currentTarget, this.getDelegateOptions())                                        // 1363
      $(obj.currentTarget).data('bs.' + this.type, self)                                                               // 1364
    }                                                                                                                  // 1365
                                                                                                                       // 1366
    if (obj instanceof $.Event) {                                                                                      // 1367
      self.inState[obj.type == 'focusin' ? 'focus' : 'hover'] = true                                                   // 1368
    }                                                                                                                  // 1369
                                                                                                                       // 1370
    if (self.tip().hasClass('in') || self.hoverState == 'in') {                                                        // 1371
      self.hoverState = 'in'                                                                                           // 1372
      return                                                                                                           // 1373
    }                                                                                                                  // 1374
                                                                                                                       // 1375
    clearTimeout(self.timeout)                                                                                         // 1376
                                                                                                                       // 1377
    self.hoverState = 'in'                                                                                             // 1378
                                                                                                                       // 1379
    if (!self.options.delay || !self.options.delay.show) return self.show()                                            // 1380
                                                                                                                       // 1381
    self.timeout = setTimeout(function () {                                                                            // 1382
      if (self.hoverState == 'in') self.show()                                                                         // 1383
    }, self.options.delay.show)                                                                                        // 1384
  }                                                                                                                    // 1385
                                                                                                                       // 1386
  Tooltip.prototype.isInStateTrue = function () {                                                                      // 1387
    for (var key in this.inState) {                                                                                    // 1388
      if (this.inState[key]) return true                                                                               // 1389
    }                                                                                                                  // 1390
                                                                                                                       // 1391
    return false                                                                                                       // 1392
  }                                                                                                                    // 1393
                                                                                                                       // 1394
  Tooltip.prototype.leave = function (obj) {                                                                           // 1395
    var self = obj instanceof this.constructor ?                                                                       // 1396
      obj : $(obj.currentTarget).data('bs.' + this.type)                                                               // 1397
                                                                                                                       // 1398
    if (!self) {                                                                                                       // 1399
      self = new this.constructor(obj.currentTarget, this.getDelegateOptions())                                        // 1400
      $(obj.currentTarget).data('bs.' + this.type, self)                                                               // 1401
    }                                                                                                                  // 1402
                                                                                                                       // 1403
    if (obj instanceof $.Event) {                                                                                      // 1404
      self.inState[obj.type == 'focusout' ? 'focus' : 'hover'] = false                                                 // 1405
    }                                                                                                                  // 1406
                                                                                                                       // 1407
    if (self.isInStateTrue()) return                                                                                   // 1408
                                                                                                                       // 1409
    clearTimeout(self.timeout)                                                                                         // 1410
                                                                                                                       // 1411
    self.hoverState = 'out'                                                                                            // 1412
                                                                                                                       // 1413
    if (!self.options.delay || !self.options.delay.hide) return self.hide()                                            // 1414
                                                                                                                       // 1415
    self.timeout = setTimeout(function () {                                                                            // 1416
      if (self.hoverState == 'out') self.hide()                                                                        // 1417
    }, self.options.delay.hide)                                                                                        // 1418
  }                                                                                                                    // 1419
                                                                                                                       // 1420
  Tooltip.prototype.show = function () {                                                                               // 1421
    var e = $.Event('show.bs.' + this.type)                                                                            // 1422
                                                                                                                       // 1423
    if (this.hasContent() && this.enabled) {                                                                           // 1424
      this.$element.trigger(e)                                                                                         // 1425
                                                                                                                       // 1426
      var inDom = $.contains(this.$element[0].ownerDocument.documentElement, this.$element[0])                         // 1427
      if (e.isDefaultPrevented() || !inDom) return                                                                     // 1428
      var that = this                                                                                                  // 1429
                                                                                                                       // 1430
      var $tip = this.tip()                                                                                            // 1431
                                                                                                                       // 1432
      var tipId = this.getUID(this.type)                                                                               // 1433
                                                                                                                       // 1434
      this.setContent()                                                                                                // 1435
      $tip.attr('id', tipId)                                                                                           // 1436
      this.$element.attr('aria-describedby', tipId)                                                                    // 1437
                                                                                                                       // 1438
      if (this.options.animation) $tip.addClass('fade')                                                                // 1439
                                                                                                                       // 1440
      var placement = typeof this.options.placement == 'function' ?                                                    // 1441
        this.options.placement.call(this, $tip[0], this.$element[0]) :                                                 // 1442
        this.options.placement                                                                                         // 1443
                                                                                                                       // 1444
      var autoToken = /\s?auto?\s?/i                                                                                   // 1445
      var autoPlace = autoToken.test(placement)                                                                        // 1446
      if (autoPlace) placement = placement.replace(autoToken, '') || 'top'                                             // 1447
                                                                                                                       // 1448
      $tip                                                                                                             // 1449
        .detach()                                                                                                      // 1450
        .css({ top: 0, left: 0, display: 'block' })                                                                    // 1451
        .addClass(placement)                                                                                           // 1452
        .data('bs.' + this.type, this)                                                                                 // 1453
                                                                                                                       // 1454
      this.options.container ? $tip.appendTo(this.options.container) : $tip.insertAfter(this.$element)                 // 1455
      this.$element.trigger('inserted.bs.' + this.type)                                                                // 1456
                                                                                                                       // 1457
      var pos          = this.getPosition()                                                                            // 1458
      var actualWidth  = $tip[0].offsetWidth                                                                           // 1459
      var actualHeight = $tip[0].offsetHeight                                                                          // 1460
                                                                                                                       // 1461
      if (autoPlace) {                                                                                                 // 1462
        var orgPlacement = placement                                                                                   // 1463
        var viewportDim = this.getPosition(this.$viewport)                                                             // 1464
                                                                                                                       // 1465
        placement = placement == 'bottom' && pos.bottom + actualHeight > viewportDim.bottom ? 'top'    :               // 1466
                    placement == 'top'    && pos.top    - actualHeight < viewportDim.top    ? 'bottom' :               // 1467
                    placement == 'right'  && pos.right  + actualWidth  > viewportDim.width  ? 'left'   :               // 1468
                    placement == 'left'   && pos.left   - actualWidth  < viewportDim.left   ? 'right'  :               // 1469
                    placement                                                                                          // 1470
                                                                                                                       // 1471
        $tip                                                                                                           // 1472
          .removeClass(orgPlacement)                                                                                   // 1473
          .addClass(placement)                                                                                         // 1474
      }                                                                                                                // 1475
                                                                                                                       // 1476
      var calculatedOffset = this.getCalculatedOffset(placement, pos, actualWidth, actualHeight)                       // 1477
                                                                                                                       // 1478
      this.applyPlacement(calculatedOffset, placement)                                                                 // 1479
                                                                                                                       // 1480
      var complete = function () {                                                                                     // 1481
        var prevHoverState = that.hoverState                                                                           // 1482
        that.$element.trigger('shown.bs.' + that.type)                                                                 // 1483
        that.hoverState = null                                                                                         // 1484
                                                                                                                       // 1485
        if (prevHoverState == 'out') that.leave(that)                                                                  // 1486
      }                                                                                                                // 1487
                                                                                                                       // 1488
      $.support.transition && this.$tip.hasClass('fade') ?                                                             // 1489
        $tip                                                                                                           // 1490
          .one('bsTransitionEnd', complete)                                                                            // 1491
          .emulateTransitionEnd(Tooltip.TRANSITION_DURATION) :                                                         // 1492
        complete()                                                                                                     // 1493
    }                                                                                                                  // 1494
  }                                                                                                                    // 1495
                                                                                                                       // 1496
  Tooltip.prototype.applyPlacement = function (offset, placement) {                                                    // 1497
    var $tip   = this.tip()                                                                                            // 1498
    var width  = $tip[0].offsetWidth                                                                                   // 1499
    var height = $tip[0].offsetHeight                                                                                  // 1500
                                                                                                                       // 1501
    // manually read margins because getBoundingClientRect includes difference                                         // 1502
    var marginTop = parseInt($tip.css('margin-top'), 10)                                                               // 1503
    var marginLeft = parseInt($tip.css('margin-left'), 10)                                                             // 1504
                                                                                                                       // 1505
    // we must check for NaN for ie 8/9                                                                                // 1506
    if (isNaN(marginTop))  marginTop  = 0                                                                              // 1507
    if (isNaN(marginLeft)) marginLeft = 0                                                                              // 1508
                                                                                                                       // 1509
    offset.top  += marginTop                                                                                           // 1510
    offset.left += marginLeft                                                                                          // 1511
                                                                                                                       // 1512
    // $.fn.offset doesn't round pixel values                                                                          // 1513
    // so we use setOffset directly with our own function B-0                                                          // 1514
    $.offset.setOffset($tip[0], $.extend({                                                                             // 1515
      using: function (props) {                                                                                        // 1516
        $tip.css({                                                                                                     // 1517
          top: Math.round(props.top),                                                                                  // 1518
          left: Math.round(props.left)                                                                                 // 1519
        })                                                                                                             // 1520
      }                                                                                                                // 1521
    }, offset), 0)                                                                                                     // 1522
                                                                                                                       // 1523
    $tip.addClass('in')                                                                                                // 1524
                                                                                                                       // 1525
    // check to see if placing tip in new offset caused the tip to resize itself                                       // 1526
    var actualWidth  = $tip[0].offsetWidth                                                                             // 1527
    var actualHeight = $tip[0].offsetHeight                                                                            // 1528
                                                                                                                       // 1529
    if (placement == 'top' && actualHeight != height) {                                                                // 1530
      offset.top = offset.top + height - actualHeight                                                                  // 1531
    }                                                                                                                  // 1532
                                                                                                                       // 1533
    var delta = this.getViewportAdjustedDelta(placement, offset, actualWidth, actualHeight)                            // 1534
                                                                                                                       // 1535
    if (delta.left) offset.left += delta.left                                                                          // 1536
    else offset.top += delta.top                                                                                       // 1537
                                                                                                                       // 1538
    var isVertical          = /top|bottom/.test(placement)                                                             // 1539
    var arrowDelta          = isVertical ? delta.left * 2 - width + actualWidth : delta.top * 2 - height + actualHeight
    var arrowOffsetPosition = isVertical ? 'offsetWidth' : 'offsetHeight'                                              // 1541
                                                                                                                       // 1542
    $tip.offset(offset)                                                                                                // 1543
    this.replaceArrow(arrowDelta, $tip[0][arrowOffsetPosition], isVertical)                                            // 1544
  }                                                                                                                    // 1545
                                                                                                                       // 1546
  Tooltip.prototype.replaceArrow = function (delta, dimension, isVertical) {                                           // 1547
    this.arrow()                                                                                                       // 1548
      .css(isVertical ? 'left' : 'top', 50 * (1 - delta / dimension) + '%')                                            // 1549
      .css(isVertical ? 'top' : 'left', '')                                                                            // 1550
  }                                                                                                                    // 1551
                                                                                                                       // 1552
  Tooltip.prototype.setContent = function () {                                                                         // 1553
    var $tip  = this.tip()                                                                                             // 1554
    var title = this.getTitle()                                                                                        // 1555
                                                                                                                       // 1556
    $tip.find('.tooltip-inner')[this.options.html ? 'html' : 'text'](title)                                            // 1557
    $tip.removeClass('fade in top bottom left right')                                                                  // 1558
  }                                                                                                                    // 1559
                                                                                                                       // 1560
  Tooltip.prototype.hide = function (callback) {                                                                       // 1561
    var that = this                                                                                                    // 1562
    var $tip = $(this.$tip)                                                                                            // 1563
    var e    = $.Event('hide.bs.' + this.type)                                                                         // 1564
                                                                                                                       // 1565
    function complete() {                                                                                              // 1566
      if (that.hoverState != 'in') $tip.detach()                                                                       // 1567
      that.$element                                                                                                    // 1568
        .removeAttr('aria-describedby')                                                                                // 1569
        .trigger('hidden.bs.' + that.type)                                                                             // 1570
      callback && callback()                                                                                           // 1571
    }                                                                                                                  // 1572
                                                                                                                       // 1573
    this.$element.trigger(e)                                                                                           // 1574
                                                                                                                       // 1575
    if (e.isDefaultPrevented()) return                                                                                 // 1576
                                                                                                                       // 1577
    $tip.removeClass('in')                                                                                             // 1578
                                                                                                                       // 1579
    $.support.transition && $tip.hasClass('fade') ?                                                                    // 1580
      $tip                                                                                                             // 1581
        .one('bsTransitionEnd', complete)                                                                              // 1582
        .emulateTransitionEnd(Tooltip.TRANSITION_DURATION) :                                                           // 1583
      complete()                                                                                                       // 1584
                                                                                                                       // 1585
    this.hoverState = null                                                                                             // 1586
                                                                                                                       // 1587
    return this                                                                                                        // 1588
  }                                                                                                                    // 1589
                                                                                                                       // 1590
  Tooltip.prototype.fixTitle = function () {                                                                           // 1591
    var $e = this.$element                                                                                             // 1592
    if ($e.attr('title') || typeof $e.attr('data-original-title') != 'string') {                                       // 1593
      $e.attr('data-original-title', $e.attr('title') || '').attr('title', '')                                         // 1594
    }                                                                                                                  // 1595
  }                                                                                                                    // 1596
                                                                                                                       // 1597
  Tooltip.prototype.hasContent = function () {                                                                         // 1598
    return this.getTitle()                                                                                             // 1599
  }                                                                                                                    // 1600
                                                                                                                       // 1601
  Tooltip.prototype.getPosition = function ($element) {                                                                // 1602
    $element   = $element || this.$element                                                                             // 1603
                                                                                                                       // 1604
    var el     = $element[0]                                                                                           // 1605
    var isBody = el.tagName == 'BODY'                                                                                  // 1606
                                                                                                                       // 1607
    var elRect    = el.getBoundingClientRect()                                                                         // 1608
    if (elRect.width == null) {                                                                                        // 1609
      // width and height are missing in IE8, so compute them manually; see https://github.com/twbs/bootstrap/issues/14093
      elRect = $.extend({}, elRect, { width: elRect.right - elRect.left, height: elRect.bottom - elRect.top })         // 1611
    }                                                                                                                  // 1612
    var elOffset  = isBody ? { top: 0, left: 0 } : $element.offset()                                                   // 1613
    var scroll    = { scroll: isBody ? document.documentElement.scrollTop || document.body.scrollTop : $element.scrollTop() }
    var outerDims = isBody ? { width: $(window).width(), height: $(window).height() } : null                           // 1615
                                                                                                                       // 1616
    return $.extend({}, elRect, scroll, outerDims, elOffset)                                                           // 1617
  }                                                                                                                    // 1618
                                                                                                                       // 1619
  Tooltip.prototype.getCalculatedOffset = function (placement, pos, actualWidth, actualHeight) {                       // 1620
    return placement == 'bottom' ? { top: pos.top + pos.height,   left: pos.left + pos.width / 2 - actualWidth / 2 } :
           placement == 'top'    ? { top: pos.top - actualHeight, left: pos.left + pos.width / 2 - actualWidth / 2 } :
           placement == 'left'   ? { top: pos.top + pos.height / 2 - actualHeight / 2, left: pos.left - actualWidth } :
        /* placement == 'right' */ { top: pos.top + pos.height / 2 - actualHeight / 2, left: pos.left + pos.width }    // 1624
                                                                                                                       // 1625
  }                                                                                                                    // 1626
                                                                                                                       // 1627
  Tooltip.prototype.getViewportAdjustedDelta = function (placement, pos, actualWidth, actualHeight) {                  // 1628
    var delta = { top: 0, left: 0 }                                                                                    // 1629
    if (!this.$viewport) return delta                                                                                  // 1630
                                                                                                                       // 1631
    var viewportPadding = this.options.viewport && this.options.viewport.padding || 0                                  // 1632
    var viewportDimensions = this.getPosition(this.$viewport)                                                          // 1633
                                                                                                                       // 1634
    if (/right|left/.test(placement)) {                                                                                // 1635
      var topEdgeOffset    = pos.top - viewportPadding - viewportDimensions.scroll                                     // 1636
      var bottomEdgeOffset = pos.top + viewportPadding - viewportDimensions.scroll + actualHeight                      // 1637
      if (topEdgeOffset < viewportDimensions.top) { // top overflow                                                    // 1638
        delta.top = viewportDimensions.top - topEdgeOffset                                                             // 1639
      } else if (bottomEdgeOffset > viewportDimensions.top + viewportDimensions.height) { // bottom overflow           // 1640
        delta.top = viewportDimensions.top + viewportDimensions.height - bottomEdgeOffset                              // 1641
      }                                                                                                                // 1642
    } else {                                                                                                           // 1643
      var leftEdgeOffset  = pos.left - viewportPadding                                                                 // 1644
      var rightEdgeOffset = pos.left + viewportPadding + actualWidth                                                   // 1645
      if (leftEdgeOffset < viewportDimensions.left) { // left overflow                                                 // 1646
        delta.left = viewportDimensions.left - leftEdgeOffset                                                          // 1647
      } else if (rightEdgeOffset > viewportDimensions.right) { // right overflow                                       // 1648
        delta.left = viewportDimensions.left + viewportDimensions.width - rightEdgeOffset                              // 1649
      }                                                                                                                // 1650
    }                                                                                                                  // 1651
                                                                                                                       // 1652
    return delta                                                                                                       // 1653
  }                                                                                                                    // 1654
                                                                                                                       // 1655
  Tooltip.prototype.getTitle = function () {                                                                           // 1656
    var title                                                                                                          // 1657
    var $e = this.$element                                                                                             // 1658
    var o  = this.options                                                                                              // 1659
                                                                                                                       // 1660
    title = $e.attr('data-original-title')                                                                             // 1661
      || (typeof o.title == 'function' ? o.title.call($e[0]) :  o.title)                                               // 1662
                                                                                                                       // 1663
    return title                                                                                                       // 1664
  }                                                                                                                    // 1665
                                                                                                                       // 1666
  Tooltip.prototype.getUID = function (prefix) {                                                                       // 1667
    do prefix += ~~(Math.random() * 1000000)                                                                           // 1668
    while (document.getElementById(prefix))                                                                            // 1669
    return prefix                                                                                                      // 1670
  }                                                                                                                    // 1671
                                                                                                                       // 1672
  Tooltip.prototype.tip = function () {                                                                                // 1673
    if (!this.$tip) {                                                                                                  // 1674
      this.$tip = $(this.options.template)                                                                             // 1675
      if (this.$tip.length != 1) {                                                                                     // 1676
        throw new Error(this.type + ' `template` option must consist of exactly 1 top-level element!')                 // 1677
      }                                                                                                                // 1678
    }                                                                                                                  // 1679
    return this.$tip                                                                                                   // 1680
  }                                                                                                                    // 1681
                                                                                                                       // 1682
  Tooltip.prototype.arrow = function () {                                                                              // 1683
    return (this.$arrow = this.$arrow || this.tip().find('.tooltip-arrow'))                                            // 1684
  }                                                                                                                    // 1685
                                                                                                                       // 1686
  Tooltip.prototype.enable = function () {                                                                             // 1687
    this.enabled = true                                                                                                // 1688
  }                                                                                                                    // 1689
                                                                                                                       // 1690
  Tooltip.prototype.disable = function () {                                                                            // 1691
    this.enabled = false                                                                                               // 1692
  }                                                                                                                    // 1693
                                                                                                                       // 1694
  Tooltip.prototype.toggleEnabled = function () {                                                                      // 1695
    this.enabled = !this.enabled                                                                                       // 1696
  }                                                                                                                    // 1697
                                                                                                                       // 1698
  Tooltip.prototype.toggle = function (e) {                                                                            // 1699
    var self = this                                                                                                    // 1700
    if (e) {                                                                                                           // 1701
      self = $(e.currentTarget).data('bs.' + this.type)                                                                // 1702
      if (!self) {                                                                                                     // 1703
        self = new this.constructor(e.currentTarget, this.getDelegateOptions())                                        // 1704
        $(e.currentTarget).data('bs.' + this.type, self)                                                               // 1705
      }                                                                                                                // 1706
    }                                                                                                                  // 1707
                                                                                                                       // 1708
    if (e) {                                                                                                           // 1709
      self.inState.click = !self.inState.click                                                                         // 1710
      if (self.isInStateTrue()) self.enter(self)                                                                       // 1711
      else self.leave(self)                                                                                            // 1712
    } else {                                                                                                           // 1713
      self.tip().hasClass('in') ? self.leave(self) : self.enter(self)                                                  // 1714
    }                                                                                                                  // 1715
  }                                                                                                                    // 1716
                                                                                                                       // 1717
  Tooltip.prototype.destroy = function () {                                                                            // 1718
    var that = this                                                                                                    // 1719
    clearTimeout(this.timeout)                                                                                         // 1720
    this.hide(function () {                                                                                            // 1721
      that.$element.off('.' + that.type).removeData('bs.' + that.type)                                                 // 1722
      if (that.$tip) {                                                                                                 // 1723
        that.$tip.detach()                                                                                             // 1724
      }                                                                                                                // 1725
      that.$tip = null                                                                                                 // 1726
      that.$arrow = null                                                                                               // 1727
      that.$viewport = null                                                                                            // 1728
    })                                                                                                                 // 1729
  }                                                                                                                    // 1730
                                                                                                                       // 1731
                                                                                                                       // 1732
  // TOOLTIP PLUGIN DEFINITION                                                                                         // 1733
  // =========================                                                                                         // 1734
                                                                                                                       // 1735
  function Plugin(option) {                                                                                            // 1736
    return this.each(function () {                                                                                     // 1737
      var $this   = $(this)                                                                                            // 1738
      var data    = $this.data('bs.tooltip')                                                                           // 1739
      var options = typeof option == 'object' && option                                                                // 1740
                                                                                                                       // 1741
      if (!data && /destroy|hide/.test(option)) return                                                                 // 1742
      if (!data) $this.data('bs.tooltip', (data = new Tooltip(this, options)))                                         // 1743
      if (typeof option == 'string') data[option]()                                                                    // 1744
    })                                                                                                                 // 1745
  }                                                                                                                    // 1746
                                                                                                                       // 1747
  var old = $.fn.tooltip                                                                                               // 1748
                                                                                                                       // 1749
  $.fn.tooltip             = Plugin                                                                                    // 1750
  $.fn.tooltip.Constructor = Tooltip                                                                                   // 1751
                                                                                                                       // 1752
                                                                                                                       // 1753
  // TOOLTIP NO CONFLICT                                                                                               // 1754
  // ===================                                                                                               // 1755
                                                                                                                       // 1756
  $.fn.tooltip.noConflict = function () {                                                                              // 1757
    $.fn.tooltip = old                                                                                                 // 1758
    return this                                                                                                        // 1759
  }                                                                                                                    // 1760
                                                                                                                       // 1761
}(jQuery);                                                                                                             // 1762
                                                                                                                       // 1763
/* ========================================================================                                            // 1764
 * Bootstrap: popover.js v3.3.6                                                                                        // 1765
 * http://getbootstrap.com/javascript/#popovers                                                                        // 1766
 * ========================================================================                                            // 1767
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 1768
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 1769
 * ======================================================================== */                                         // 1770
                                                                                                                       // 1771
                                                                                                                       // 1772
+function ($) {                                                                                                        // 1773
  'use strict';                                                                                                        // 1774
                                                                                                                       // 1775
  // POPOVER PUBLIC CLASS DEFINITION                                                                                   // 1776
  // ===============================                                                                                   // 1777
                                                                                                                       // 1778
  var Popover = function (element, options) {                                                                          // 1779
    this.init('popover', element, options)                                                                             // 1780
  }                                                                                                                    // 1781
                                                                                                                       // 1782
  if (!$.fn.tooltip) throw new Error('Popover requires tooltip.js')                                                    // 1783
                                                                                                                       // 1784
  Popover.VERSION  = '3.3.6'                                                                                           // 1785
                                                                                                                       // 1786
  Popover.DEFAULTS = $.extend({}, $.fn.tooltip.Constructor.DEFAULTS, {                                                 // 1787
    placement: 'right',                                                                                                // 1788
    trigger: 'click',                                                                                                  // 1789
    content: '',                                                                                                       // 1790
    template: '<div class="popover" role="tooltip"><div class="arrow"></div><h3 class="popover-title"></h3><div class="popover-content"></div></div>'
  })                                                                                                                   // 1792
                                                                                                                       // 1793
                                                                                                                       // 1794
  // NOTE: POPOVER EXTENDS tooltip.js                                                                                  // 1795
  // ================================                                                                                  // 1796
                                                                                                                       // 1797
  Popover.prototype = $.extend({}, $.fn.tooltip.Constructor.prototype)                                                 // 1798
                                                                                                                       // 1799
  Popover.prototype.constructor = Popover                                                                              // 1800
                                                                                                                       // 1801
  Popover.prototype.getDefaults = function () {                                                                        // 1802
    return Popover.DEFAULTS                                                                                            // 1803
  }                                                                                                                    // 1804
                                                                                                                       // 1805
  Popover.prototype.setContent = function () {                                                                         // 1806
    var $tip    = this.tip()                                                                                           // 1807
    var title   = this.getTitle()                                                                                      // 1808
    var content = this.getContent()                                                                                    // 1809
                                                                                                                       // 1810
    $tip.find('.popover-title')[this.options.html ? 'html' : 'text'](title)                                            // 1811
    $tip.find('.popover-content').children().detach().end()[ // we use append for html objects to maintain js events   // 1812
      this.options.html ? (typeof content == 'string' ? 'html' : 'append') : 'text'                                    // 1813
    ](content)                                                                                                         // 1814
                                                                                                                       // 1815
    $tip.removeClass('fade top bottom left right in')                                                                  // 1816
                                                                                                                       // 1817
    // IE8 doesn't accept hiding via the `:empty` pseudo selector, we have to do                                       // 1818
    // this manually by checking the contents.                                                                         // 1819
    if (!$tip.find('.popover-title').html()) $tip.find('.popover-title').hide()                                        // 1820
  }                                                                                                                    // 1821
                                                                                                                       // 1822
  Popover.prototype.hasContent = function () {                                                                         // 1823
    return this.getTitle() || this.getContent()                                                                        // 1824
  }                                                                                                                    // 1825
                                                                                                                       // 1826
  Popover.prototype.getContent = function () {                                                                         // 1827
    var $e = this.$element                                                                                             // 1828
    var o  = this.options                                                                                              // 1829
                                                                                                                       // 1830
    return $e.attr('data-content')                                                                                     // 1831
      || (typeof o.content == 'function' ?                                                                             // 1832
            o.content.call($e[0]) :                                                                                    // 1833
            o.content)                                                                                                 // 1834
  }                                                                                                                    // 1835
                                                                                                                       // 1836
  Popover.prototype.arrow = function () {                                                                              // 1837
    return (this.$arrow = this.$arrow || this.tip().find('.arrow'))                                                    // 1838
  }                                                                                                                    // 1839
                                                                                                                       // 1840
                                                                                                                       // 1841
  // POPOVER PLUGIN DEFINITION                                                                                         // 1842
  // =========================                                                                                         // 1843
                                                                                                                       // 1844
  function Plugin(option) {                                                                                            // 1845
    return this.each(function () {                                                                                     // 1846
      var $this   = $(this)                                                                                            // 1847
      var data    = $this.data('bs.popover')                                                                           // 1848
      var options = typeof option == 'object' && option                                                                // 1849
                                                                                                                       // 1850
      if (!data && /destroy|hide/.test(option)) return                                                                 // 1851
      if (!data) $this.data('bs.popover', (data = new Popover(this, options)))                                         // 1852
      if (typeof option == 'string') data[option]()                                                                    // 1853
    })                                                                                                                 // 1854
  }                                                                                                                    // 1855
                                                                                                                       // 1856
  var old = $.fn.popover                                                                                               // 1857
                                                                                                                       // 1858
  $.fn.popover             = Plugin                                                                                    // 1859
  $.fn.popover.Constructor = Popover                                                                                   // 1860
                                                                                                                       // 1861
                                                                                                                       // 1862
  // POPOVER NO CONFLICT                                                                                               // 1863
  // ===================                                                                                               // 1864
                                                                                                                       // 1865
  $.fn.popover.noConflict = function () {                                                                              // 1866
    $.fn.popover = old                                                                                                 // 1867
    return this                                                                                                        // 1868
  }                                                                                                                    // 1869
                                                                                                                       // 1870
}(jQuery);                                                                                                             // 1871
                                                                                                                       // 1872
/* ========================================================================                                            // 1873
 * Bootstrap: scrollspy.js v3.3.6                                                                                      // 1874
 * http://getbootstrap.com/javascript/#scrollspy                                                                       // 1875
 * ========================================================================                                            // 1876
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 1877
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 1878
 * ======================================================================== */                                         // 1879
                                                                                                                       // 1880
                                                                                                                       // 1881
+function ($) {                                                                                                        // 1882
  'use strict';                                                                                                        // 1883
                                                                                                                       // 1884
  // SCROLLSPY CLASS DEFINITION                                                                                        // 1885
  // ==========================                                                                                        // 1886
                                                                                                                       // 1887
  function ScrollSpy(element, options) {                                                                               // 1888
    this.$body          = $(document.body)                                                                             // 1889
    this.$scrollElement = $(element).is(document.body) ? $(window) : $(element)                                        // 1890
    this.options        = $.extend({}, ScrollSpy.DEFAULTS, options)                                                    // 1891
    this.selector       = (this.options.target || '') + ' .nav li > a'                                                 // 1892
    this.offsets        = []                                                                                           // 1893
    this.targets        = []                                                                                           // 1894
    this.activeTarget   = null                                                                                         // 1895
    this.scrollHeight   = 0                                                                                            // 1896
                                                                                                                       // 1897
    this.$scrollElement.on('scroll.bs.scrollspy', $.proxy(this.process, this))                                         // 1898
    this.refresh()                                                                                                     // 1899
    this.process()                                                                                                     // 1900
  }                                                                                                                    // 1901
                                                                                                                       // 1902
  ScrollSpy.VERSION  = '3.3.6'                                                                                         // 1903
                                                                                                                       // 1904
  ScrollSpy.DEFAULTS = {                                                                                               // 1905
    offset: 10                                                                                                         // 1906
  }                                                                                                                    // 1907
                                                                                                                       // 1908
  ScrollSpy.prototype.getScrollHeight = function () {                                                                  // 1909
    return this.$scrollElement[0].scrollHeight || Math.max(this.$body[0].scrollHeight, document.documentElement.scrollHeight)
  }                                                                                                                    // 1911
                                                                                                                       // 1912
  ScrollSpy.prototype.refresh = function () {                                                                          // 1913
    var that          = this                                                                                           // 1914
    var offsetMethod  = 'offset'                                                                                       // 1915
    var offsetBase    = 0                                                                                              // 1916
                                                                                                                       // 1917
    this.offsets      = []                                                                                             // 1918
    this.targets      = []                                                                                             // 1919
    this.scrollHeight = this.getScrollHeight()                                                                         // 1920
                                                                                                                       // 1921
    if (!$.isWindow(this.$scrollElement[0])) {                                                                         // 1922
      offsetMethod = 'position'                                                                                        // 1923
      offsetBase   = this.$scrollElement.scrollTop()                                                                   // 1924
    }                                                                                                                  // 1925
                                                                                                                       // 1926
    this.$body                                                                                                         // 1927
      .find(this.selector)                                                                                             // 1928
      .map(function () {                                                                                               // 1929
        var $el   = $(this)                                                                                            // 1930
        var href  = $el.data('target') || $el.attr('href')                                                             // 1931
        var $href = /^#./.test(href) && $(href)                                                                        // 1932
                                                                                                                       // 1933
        return ($href                                                                                                  // 1934
          && $href.length                                                                                              // 1935
          && $href.is(':visible')                                                                                      // 1936
          && [[$href[offsetMethod]().top + offsetBase, href]]) || null                                                 // 1937
      })                                                                                                               // 1938
      .sort(function (a, b) { return a[0] - b[0] })                                                                    // 1939
      .each(function () {                                                                                              // 1940
        that.offsets.push(this[0])                                                                                     // 1941
        that.targets.push(this[1])                                                                                     // 1942
      })                                                                                                               // 1943
  }                                                                                                                    // 1944
                                                                                                                       // 1945
  ScrollSpy.prototype.process = function () {                                                                          // 1946
    var scrollTop    = this.$scrollElement.scrollTop() + this.options.offset                                           // 1947
    var scrollHeight = this.getScrollHeight()                                                                          // 1948
    var maxScroll    = this.options.offset + scrollHeight - this.$scrollElement.height()                               // 1949
    var offsets      = this.offsets                                                                                    // 1950
    var targets      = this.targets                                                                                    // 1951
    var activeTarget = this.activeTarget                                                                               // 1952
    var i                                                                                                              // 1953
                                                                                                                       // 1954
    if (this.scrollHeight != scrollHeight) {                                                                           // 1955
      this.refresh()                                                                                                   // 1956
    }                                                                                                                  // 1957
                                                                                                                       // 1958
    if (scrollTop >= maxScroll) {                                                                                      // 1959
      return activeTarget != (i = targets[targets.length - 1]) && this.activate(i)                                     // 1960
    }                                                                                                                  // 1961
                                                                                                                       // 1962
    if (activeTarget && scrollTop < offsets[0]) {                                                                      // 1963
      this.activeTarget = null                                                                                         // 1964
      return this.clear()                                                                                              // 1965
    }                                                                                                                  // 1966
                                                                                                                       // 1967
    for (i = offsets.length; i--;) {                                                                                   // 1968
      activeTarget != targets[i]                                                                                       // 1969
        && scrollTop >= offsets[i]                                                                                     // 1970
        && (offsets[i + 1] === undefined || scrollTop < offsets[i + 1])                                                // 1971
        && this.activate(targets[i])                                                                                   // 1972
    }                                                                                                                  // 1973
  }                                                                                                                    // 1974
                                                                                                                       // 1975
  ScrollSpy.prototype.activate = function (target) {                                                                   // 1976
    this.activeTarget = target                                                                                         // 1977
                                                                                                                       // 1978
    this.clear()                                                                                                       // 1979
                                                                                                                       // 1980
    var selector = this.selector +                                                                                     // 1981
      '[data-target="' + target + '"],' +                                                                              // 1982
      this.selector + '[href="' + target + '"]'                                                                        // 1983
                                                                                                                       // 1984
    var active = $(selector)                                                                                           // 1985
      .parents('li')                                                                                                   // 1986
      .addClass('active')                                                                                              // 1987
                                                                                                                       // 1988
    if (active.parent('.dropdown-menu').length) {                                                                      // 1989
      active = active                                                                                                  // 1990
        .closest('li.dropdown')                                                                                        // 1991
        .addClass('active')                                                                                            // 1992
    }                                                                                                                  // 1993
                                                                                                                       // 1994
    active.trigger('activate.bs.scrollspy')                                                                            // 1995
  }                                                                                                                    // 1996
                                                                                                                       // 1997
  ScrollSpy.prototype.clear = function () {                                                                            // 1998
    $(this.selector)                                                                                                   // 1999
      .parentsUntil(this.options.target, '.active')                                                                    // 2000
      .removeClass('active')                                                                                           // 2001
  }                                                                                                                    // 2002
                                                                                                                       // 2003
                                                                                                                       // 2004
  // SCROLLSPY PLUGIN DEFINITION                                                                                       // 2005
  // ===========================                                                                                       // 2006
                                                                                                                       // 2007
  function Plugin(option) {                                                                                            // 2008
    return this.each(function () {                                                                                     // 2009
      var $this   = $(this)                                                                                            // 2010
      var data    = $this.data('bs.scrollspy')                                                                         // 2011
      var options = typeof option == 'object' && option                                                                // 2012
                                                                                                                       // 2013
      if (!data) $this.data('bs.scrollspy', (data = new ScrollSpy(this, options)))                                     // 2014
      if (typeof option == 'string') data[option]()                                                                    // 2015
    })                                                                                                                 // 2016
  }                                                                                                                    // 2017
                                                                                                                       // 2018
  var old = $.fn.scrollspy                                                                                             // 2019
                                                                                                                       // 2020
  $.fn.scrollspy             = Plugin                                                                                  // 2021
  $.fn.scrollspy.Constructor = ScrollSpy                                                                               // 2022
                                                                                                                       // 2023
                                                                                                                       // 2024
  // SCROLLSPY NO CONFLICT                                                                                             // 2025
  // =====================                                                                                             // 2026
                                                                                                                       // 2027
  $.fn.scrollspy.noConflict = function () {                                                                            // 2028
    $.fn.scrollspy = old                                                                                               // 2029
    return this                                                                                                        // 2030
  }                                                                                                                    // 2031
                                                                                                                       // 2032
                                                                                                                       // 2033
  // SCROLLSPY DATA-API                                                                                                // 2034
  // ==================                                                                                                // 2035
                                                                                                                       // 2036
  $(window).on('load.bs.scrollspy.data-api', function () {                                                             // 2037
    $('[data-spy="scroll"]').each(function () {                                                                        // 2038
      var $spy = $(this)                                                                                               // 2039
      Plugin.call($spy, $spy.data())                                                                                   // 2040
    })                                                                                                                 // 2041
  })                                                                                                                   // 2042
                                                                                                                       // 2043
}(jQuery);                                                                                                             // 2044
                                                                                                                       // 2045
/* ========================================================================                                            // 2046
 * Bootstrap: tab.js v3.3.6                                                                                            // 2047
 * http://getbootstrap.com/javascript/#tabs                                                                            // 2048
 * ========================================================================                                            // 2049
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 2050
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 2051
 * ======================================================================== */                                         // 2052
                                                                                                                       // 2053
                                                                                                                       // 2054
+function ($) {                                                                                                        // 2055
  'use strict';                                                                                                        // 2056
                                                                                                                       // 2057
  // TAB CLASS DEFINITION                                                                                              // 2058
  // ====================                                                                                              // 2059
                                                                                                                       // 2060
  var Tab = function (element) {                                                                                       // 2061
    // jscs:disable requireDollarBeforejQueryAssignment                                                                // 2062
    this.element = $(element)                                                                                          // 2063
    // jscs:enable requireDollarBeforejQueryAssignment                                                                 // 2064
  }                                                                                                                    // 2065
                                                                                                                       // 2066
  Tab.VERSION = '3.3.6'                                                                                                // 2067
                                                                                                                       // 2068
  Tab.TRANSITION_DURATION = 150                                                                                        // 2069
                                                                                                                       // 2070
  Tab.prototype.show = function () {                                                                                   // 2071
    var $this    = this.element                                                                                        // 2072
    var $ul      = $this.closest('ul:not(.dropdown-menu)')                                                             // 2073
    var selector = $this.data('target')                                                                                // 2074
                                                                                                                       // 2075
    if (!selector) {                                                                                                   // 2076
      selector = $this.attr('href')                                                                                    // 2077
      selector = selector && selector.replace(/.*(?=#[^\s]*$)/, '') // strip for ie7                                   // 2078
    }                                                                                                                  // 2079
                                                                                                                       // 2080
    if ($this.parent('li').hasClass('active')) return                                                                  // 2081
                                                                                                                       // 2082
    var $previous = $ul.find('.active:last a')                                                                         // 2083
    var hideEvent = $.Event('hide.bs.tab', {                                                                           // 2084
      relatedTarget: $this[0]                                                                                          // 2085
    })                                                                                                                 // 2086
    var showEvent = $.Event('show.bs.tab', {                                                                           // 2087
      relatedTarget: $previous[0]                                                                                      // 2088
    })                                                                                                                 // 2089
                                                                                                                       // 2090
    $previous.trigger(hideEvent)                                                                                       // 2091
    $this.trigger(showEvent)                                                                                           // 2092
                                                                                                                       // 2093
    if (showEvent.isDefaultPrevented() || hideEvent.isDefaultPrevented()) return                                       // 2094
                                                                                                                       // 2095
    var $target = $(selector)                                                                                          // 2096
                                                                                                                       // 2097
    this.activate($this.closest('li'), $ul)                                                                            // 2098
    this.activate($target, $target.parent(), function () {                                                             // 2099
      $previous.trigger({                                                                                              // 2100
        type: 'hidden.bs.tab',                                                                                         // 2101
        relatedTarget: $this[0]                                                                                        // 2102
      })                                                                                                               // 2103
      $this.trigger({                                                                                                  // 2104
        type: 'shown.bs.tab',                                                                                          // 2105
        relatedTarget: $previous[0]                                                                                    // 2106
      })                                                                                                               // 2107
    })                                                                                                                 // 2108
  }                                                                                                                    // 2109
                                                                                                                       // 2110
  Tab.prototype.activate = function (element, container, callback) {                                                   // 2111
    var $active    = container.find('> .active')                                                                       // 2112
    var transition = callback                                                                                          // 2113
      && $.support.transition                                                                                          // 2114
      && ($active.length && $active.hasClass('fade') || !!container.find('> .fade').length)                            // 2115
                                                                                                                       // 2116
    function next() {                                                                                                  // 2117
      $active                                                                                                          // 2118
        .removeClass('active')                                                                                         // 2119
        .find('> .dropdown-menu > .active')                                                                            // 2120
          .removeClass('active')                                                                                       // 2121
        .end()                                                                                                         // 2122
        .find('[data-toggle="tab"]')                                                                                   // 2123
          .attr('aria-expanded', false)                                                                                // 2124
                                                                                                                       // 2125
      element                                                                                                          // 2126
        .addClass('active')                                                                                            // 2127
        .find('[data-toggle="tab"]')                                                                                   // 2128
          .attr('aria-expanded', true)                                                                                 // 2129
                                                                                                                       // 2130
      if (transition) {                                                                                                // 2131
        element[0].offsetWidth // reflow for transition                                                                // 2132
        element.addClass('in')                                                                                         // 2133
      } else {                                                                                                         // 2134
        element.removeClass('fade')                                                                                    // 2135
      }                                                                                                                // 2136
                                                                                                                       // 2137
      if (element.parent('.dropdown-menu').length) {                                                                   // 2138
        element                                                                                                        // 2139
          .closest('li.dropdown')                                                                                      // 2140
            .addClass('active')                                                                                        // 2141
          .end()                                                                                                       // 2142
          .find('[data-toggle="tab"]')                                                                                 // 2143
            .attr('aria-expanded', true)                                                                               // 2144
      }                                                                                                                // 2145
                                                                                                                       // 2146
      callback && callback()                                                                                           // 2147
    }                                                                                                                  // 2148
                                                                                                                       // 2149
    $active.length && transition ?                                                                                     // 2150
      $active                                                                                                          // 2151
        .one('bsTransitionEnd', next)                                                                                  // 2152
        .emulateTransitionEnd(Tab.TRANSITION_DURATION) :                                                               // 2153
      next()                                                                                                           // 2154
                                                                                                                       // 2155
    $active.removeClass('in')                                                                                          // 2156
  }                                                                                                                    // 2157
                                                                                                                       // 2158
                                                                                                                       // 2159
  // TAB PLUGIN DEFINITION                                                                                             // 2160
  // =====================                                                                                             // 2161
                                                                                                                       // 2162
  function Plugin(option) {                                                                                            // 2163
    return this.each(function () {                                                                                     // 2164
      var $this = $(this)                                                                                              // 2165
      var data  = $this.data('bs.tab')                                                                                 // 2166
                                                                                                                       // 2167
      if (!data) $this.data('bs.tab', (data = new Tab(this)))                                                          // 2168
      if (typeof option == 'string') data[option]()                                                                    // 2169
    })                                                                                                                 // 2170
  }                                                                                                                    // 2171
                                                                                                                       // 2172
  var old = $.fn.tab                                                                                                   // 2173
                                                                                                                       // 2174
  $.fn.tab             = Plugin                                                                                        // 2175
  $.fn.tab.Constructor = Tab                                                                                           // 2176
                                                                                                                       // 2177
                                                                                                                       // 2178
  // TAB NO CONFLICT                                                                                                   // 2179
  // ===============                                                                                                   // 2180
                                                                                                                       // 2181
  $.fn.tab.noConflict = function () {                                                                                  // 2182
    $.fn.tab = old                                                                                                     // 2183
    return this                                                                                                        // 2184
  }                                                                                                                    // 2185
                                                                                                                       // 2186
                                                                                                                       // 2187
  // TAB DATA-API                                                                                                      // 2188
  // ============                                                                                                      // 2189
                                                                                                                       // 2190
  var clickHandler = function (e) {                                                                                    // 2191
    e.preventDefault()                                                                                                 // 2192
    Plugin.call($(this), 'show')                                                                                       // 2193
  }                                                                                                                    // 2194
                                                                                                                       // 2195
  $(document)                                                                                                          // 2196
    .on('click.bs.tab.data-api', '[data-toggle="tab"]', clickHandler)                                                  // 2197
    .on('click.bs.tab.data-api', '[data-toggle="pill"]', clickHandler)                                                 // 2198
                                                                                                                       // 2199
}(jQuery);                                                                                                             // 2200
                                                                                                                       // 2201
/* ========================================================================                                            // 2202
 * Bootstrap: affix.js v3.3.6                                                                                          // 2203
 * http://getbootstrap.com/javascript/#affix                                                                           // 2204
 * ========================================================================                                            // 2205
 * Copyright 2011-2015 Twitter, Inc.                                                                                   // 2206
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)                                          // 2207
 * ======================================================================== */                                         // 2208
                                                                                                                       // 2209
                                                                                                                       // 2210
+function ($) {                                                                                                        // 2211
  'use strict';                                                                                                        // 2212
                                                                                                                       // 2213
  // AFFIX CLASS DEFINITION                                                                                            // 2214
  // ======================                                                                                            // 2215
                                                                                                                       // 2216
  var Affix = function (element, options) {                                                                            // 2217
    this.options = $.extend({}, Affix.DEFAULTS, options)                                                               // 2218
                                                                                                                       // 2219
    this.$target = $(this.options.target)                                                                              // 2220
      .on('scroll.bs.affix.data-api', $.proxy(this.checkPosition, this))                                               // 2221
      .on('click.bs.affix.data-api',  $.proxy(this.checkPositionWithEventLoop, this))                                  // 2222
                                                                                                                       // 2223
    this.$element     = $(element)                                                                                     // 2224
    this.affixed      = null                                                                                           // 2225
    this.unpin        = null                                                                                           // 2226
    this.pinnedOffset = null                                                                                           // 2227
                                                                                                                       // 2228
    this.checkPosition()                                                                                               // 2229
  }                                                                                                                    // 2230
                                                                                                                       // 2231
  Affix.VERSION  = '3.3.6'                                                                                             // 2232
                                                                                                                       // 2233
  Affix.RESET    = 'affix affix-top affix-bottom'                                                                      // 2234
                                                                                                                       // 2235
  Affix.DEFAULTS = {                                                                                                   // 2236
    offset: 0,                                                                                                         // 2237
    target: window                                                                                                     // 2238
  }                                                                                                                    // 2239
                                                                                                                       // 2240
  Affix.prototype.getState = function (scrollHeight, height, offsetTop, offsetBottom) {                                // 2241
    var scrollTop    = this.$target.scrollTop()                                                                        // 2242
    var position     = this.$element.offset()                                                                          // 2243
    var targetHeight = this.$target.height()                                                                           // 2244
                                                                                                                       // 2245
    if (offsetTop != null && this.affixed == 'top') return scrollTop < offsetTop ? 'top' : false                       // 2246
                                                                                                                       // 2247
    if (this.affixed == 'bottom') {                                                                                    // 2248
      if (offsetTop != null) return (scrollTop + this.unpin <= position.top) ? false : 'bottom'                        // 2249
      return (scrollTop + targetHeight <= scrollHeight - offsetBottom) ? false : 'bottom'                              // 2250
    }                                                                                                                  // 2251
                                                                                                                       // 2252
    var initializing   = this.affixed == null                                                                          // 2253
    var colliderTop    = initializing ? scrollTop : position.top                                                       // 2254
    var colliderHeight = initializing ? targetHeight : height                                                          // 2255
                                                                                                                       // 2256
    if (offsetTop != null && scrollTop <= offsetTop) return 'top'                                                      // 2257
    if (offsetBottom != null && (colliderTop + colliderHeight >= scrollHeight - offsetBottom)) return 'bottom'         // 2258
                                                                                                                       // 2259
    return false                                                                                                       // 2260
  }                                                                                                                    // 2261
                                                                                                                       // 2262
  Affix.prototype.getPinnedOffset = function () {                                                                      // 2263
    if (this.pinnedOffset) return this.pinnedOffset                                                                    // 2264
    this.$element.removeClass(Affix.RESET).addClass('affix')                                                           // 2265
    var scrollTop = this.$target.scrollTop()                                                                           // 2266
    var position  = this.$element.offset()                                                                             // 2267
    return (this.pinnedOffset = position.top - scrollTop)                                                              // 2268
  }                                                                                                                    // 2269
                                                                                                                       // 2270
  Affix.prototype.checkPositionWithEventLoop = function () {                                                           // 2271
    setTimeout($.proxy(this.checkPosition, this), 1)                                                                   // 2272
  }                                                                                                                    // 2273
                                                                                                                       // 2274
  Affix.prototype.checkPosition = function () {                                                                        // 2275
    if (!this.$element.is(':visible')) return                                                                          // 2276
                                                                                                                       // 2277
    var height       = this.$element.height()                                                                          // 2278
    var offset       = this.options.offset                                                                             // 2279
    var offsetTop    = offset.top                                                                                      // 2280
    var offsetBottom = offset.bottom                                                                                   // 2281
    var scrollHeight = Math.max($(document).height(), $(document.body).height())                                       // 2282
                                                                                                                       // 2283
    if (typeof offset != 'object')         offsetBottom = offsetTop = offset                                           // 2284
    if (typeof offsetTop == 'function')    offsetTop    = offset.top(this.$element)                                    // 2285
    if (typeof offsetBottom == 'function') offsetBottom = offset.bottom(this.$element)                                 // 2286
                                                                                                                       // 2287
    var affix = this.getState(scrollHeight, height, offsetTop, offsetBottom)                                           // 2288
                                                                                                                       // 2289
    if (this.affixed != affix) {                                                                                       // 2290
      if (this.unpin != null) this.$element.css('top', '')                                                             // 2291
                                                                                                                       // 2292
      var affixType = 'affix' + (affix ? '-' + affix : '')                                                             // 2293
      var e         = $.Event(affixType + '.bs.affix')                                                                 // 2294
                                                                                                                       // 2295
      this.$element.trigger(e)                                                                                         // 2296
                                                                                                                       // 2297
      if (e.isDefaultPrevented()) return                                                                               // 2298
                                                                                                                       // 2299
      this.affixed = affix                                                                                             // 2300
      this.unpin = affix == 'bottom' ? this.getPinnedOffset() : null                                                   // 2301
                                                                                                                       // 2302
      this.$element                                                                                                    // 2303
        .removeClass(Affix.RESET)                                                                                      // 2304
        .addClass(affixType)                                                                                           // 2305
        .trigger(affixType.replace('affix', 'affixed') + '.bs.affix')                                                  // 2306
    }                                                                                                                  // 2307
                                                                                                                       // 2308
    if (affix == 'bottom') {                                                                                           // 2309
      this.$element.offset({                                                                                           // 2310
        top: scrollHeight - height - offsetBottom                                                                      // 2311
      })                                                                                                               // 2312
    }                                                                                                                  // 2313
  }                                                                                                                    // 2314
                                                                                                                       // 2315
                                                                                                                       // 2316
  // AFFIX PLUGIN DEFINITION                                                                                           // 2317
  // =======================                                                                                           // 2318
                                                                                                                       // 2319
  function Plugin(option) {                                                                                            // 2320
    return this.each(function () {                                                                                     // 2321
      var $this   = $(this)                                                                                            // 2322
      var data    = $this.data('bs.affix')                                                                             // 2323
      var options = typeof option == 'object' && option                                                                // 2324
                                                                                                                       // 2325
      if (!data) $this.data('bs.affix', (data = new Affix(this, options)))                                             // 2326
      if (typeof option == 'string') data[option]()                                                                    // 2327
    })                                                                                                                 // 2328
  }                                                                                                                    // 2329
                                                                                                                       // 2330
  var old = $.fn.affix                                                                                                 // 2331
                                                                                                                       // 2332
  $.fn.affix             = Plugin                                                                                      // 2333
  $.fn.affix.Constructor = Affix                                                                                       // 2334
                                                                                                                       // 2335
                                                                                                                       // 2336
  // AFFIX NO CONFLICT                                                                                                 // 2337
  // =================                                                                                                 // 2338
                                                                                                                       // 2339
  $.fn.affix.noConflict = function () {                                                                                // 2340
    $.fn.affix = old                                                                                                   // 2341
    return this                                                                                                        // 2342
  }                                                                                                                    // 2343
                                                                                                                       // 2344
                                                                                                                       // 2345
  // AFFIX DATA-API                                                                                                    // 2346
  // ==============                                                                                                    // 2347
                                                                                                                       // 2348
  $(window).on('load', function () {                                                                                   // 2349
    $('[data-spy="affix"]').each(function () {                                                                         // 2350
      var $spy = $(this)                                                                                               // 2351
      var data = $spy.data()                                                                                           // 2352
                                                                                                                       // 2353
      data.offset = data.offset || {}                                                                                  // 2354
                                                                                                                       // 2355
      if (data.offsetBottom != null) data.offset.bottom = data.offsetBottom                                            // 2356
      if (data.offsetTop    != null) data.offset.top    = data.offsetTop                                               // 2357
                                                                                                                       // 2358
      Plugin.call($spy, data)                                                                                          // 2359
    })                                                                                                                 // 2360
  })                                                                                                                   // 2361
                                                                                                                       // 2362
}(jQuery);                                                                                                             // 2363
                                                                                                                       // 2364
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
Package['twbs:bootstrap'] = {};

})();
