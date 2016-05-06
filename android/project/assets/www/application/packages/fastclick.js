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
var module, FastClick;

(function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/fastclick/pre.js                                                                                           //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// Define an object named module.exports. This will cause fastclick.js to put                                          // 1
// FastClick as a field on it, instead of in the global namespace.                                                     // 2
// See also post.js.                                                                                                   // 3
module = {                                                                                                             // 4
  exports: {}                                                                                                          // 5
};                                                                                                                     // 6
                                                                                                                       // 7
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/fastclick/fastclick.js                                                                                     //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
	/**                                                                                                                   // 1
	 * @preserve FastClick: polyfill to remove click delays on browsers with touch UIs.                                   // 2
	 *                                                                                                                    // 3
	 * @codingstandard ftlabs-jsv2                                                                                        // 4
	 * @copyright The Financial Times Limited [All Rights Reserved]                                                       // 5
	 * @license MIT License (see LICENSE.txt)                                                                             // 6
	 */                                                                                                                   // 7
                                                                                                                       // 8
	/*jslint browser:true, node:true*/                                                                                    // 9
	/*global define, Event, Node*/                                                                                        // 10
                                                                                                                       // 11
                                                                                                                       // 12
	/**                                                                                                                   // 13
	 * Instantiate fast-clicking listeners on the specified layer.                                                        // 14
	 *                                                                                                                    // 15
	 * @constructor                                                                                                       // 16
	 * @param {Element} layer The layer to listen on                                                                      // 17
	 * @param {Object} [options={}] The options to override the defaults                                                  // 18
	 */                                                                                                                   // 19
	function FastClick(layer, options) {                                                                                  // 20
    'use strict';                                                                                                      // 21
		var oldOnClick;                                                                                                      // 22
                                                                                                                       // 23
		options = options || {};                                                                                             // 24
                                                                                                                       // 25
		/**                                                                                                                  // 26
		 * Whether a click is currently being tracked.                                                                       // 27
		 *                                                                                                                   // 28
		 * @type boolean                                                                                                     // 29
		 */                                                                                                                  // 30
		this.trackingClick = false;                                                                                          // 31
                                                                                                                       // 32
                                                                                                                       // 33
		/**                                                                                                                  // 34
		 * Timestamp for when click tracking started.                                                                        // 35
		 *                                                                                                                   // 36
		 * @type number                                                                                                      // 37
		 */                                                                                                                  // 38
		this.trackingClickStart = 0;                                                                                         // 39
                                                                                                                       // 40
                                                                                                                       // 41
		/**                                                                                                                  // 42
		 * The element being tracked for a click.                                                                            // 43
		 *                                                                                                                   // 44
		 * @type EventTarget                                                                                                 // 45
		 */                                                                                                                  // 46
		this.targetElement = null;                                                                                           // 47
                                                                                                                       // 48
                                                                                                                       // 49
		/**                                                                                                                  // 50
		 * X-coordinate of touch start event.                                                                                // 51
		 *                                                                                                                   // 52
		 * @type number                                                                                                      // 53
		 */                                                                                                                  // 54
		this.touchStartX = 0;                                                                                                // 55
                                                                                                                       // 56
                                                                                                                       // 57
		/**                                                                                                                  // 58
		 * Y-coordinate of touch start event.                                                                                // 59
		 *                                                                                                                   // 60
		 * @type number                                                                                                      // 61
		 */                                                                                                                  // 62
		this.touchStartY = 0;                                                                                                // 63
                                                                                                                       // 64
                                                                                                                       // 65
		/**                                                                                                                  // 66
		 * ID of the last touch, retrieved from Touch.identifier.                                                            // 67
		 *                                                                                                                   // 68
		 * @type number                                                                                                      // 69
		 */                                                                                                                  // 70
		this.lastTouchIdentifier = 0;                                                                                        // 71
                                                                                                                       // 72
                                                                                                                       // 73
		/**                                                                                                                  // 74
		 * Touchmove boundary, beyond which a click will be cancelled.                                                       // 75
		 *                                                                                                                   // 76
		 * @type number                                                                                                      // 77
		 */                                                                                                                  // 78
		this.touchBoundary = options.touchBoundary || 10;                                                                    // 79
                                                                                                                       // 80
                                                                                                                       // 81
		/**                                                                                                                  // 82
		 * The FastClick layer.                                                                                              // 83
		 *                                                                                                                   // 84
		 * @type Element                                                                                                     // 85
		 */                                                                                                                  // 86
		this.layer = layer;                                                                                                  // 87
                                                                                                                       // 88
		/**                                                                                                                  // 89
		 * The minimum time between tap(touchstart and touchend) events                                                      // 90
		 *                                                                                                                   // 91
		 * @type number                                                                                                      // 92
		 */                                                                                                                  // 93
		this.tapDelay = options.tapDelay || 200;                                                                             // 94
                                                                                                                       // 95
		/**                                                                                                                  // 96
		 * The maximum time for a tap                                                                                        // 97
		 *                                                                                                                   // 98
		 * @type number                                                                                                      // 99
		 */                                                                                                                  // 100
		this.tapTimeout = options.tapTimeout || 700;                                                                         // 101
                                                                                                                       // 102
		if (FastClick.notNeeded(layer)) {                                                                                    // 103
			return;                                                                                                             // 104
		}                                                                                                                    // 105
                                                                                                                       // 106
		// Some old versions of Android don't have Function.prototype.bind                                                   // 107
		function bind(method, context) {                                                                                     // 108
			return function() { return method.apply(context, arguments); };                                                     // 109
		}                                                                                                                    // 110
                                                                                                                       // 111
                                                                                                                       // 112
		var methods = ['onMouse', 'onClick', 'onTouchStart', 'onTouchMove', 'onTouchEnd', 'onTouchCancel'];                  // 113
		var context = this;                                                                                                  // 114
		for (var i = 0, l = methods.length; i < l; i++) {                                                                    // 115
			context[methods[i]] = bind(context[methods[i]], context);                                                           // 116
		}                                                                                                                    // 117
                                                                                                                       // 118
		// Set up event handlers as required                                                                                 // 119
		if (deviceIsAndroid) {                                                                                               // 120
			layer.addEventListener('mouseover', this.onMouse, true);                                                            // 121
			layer.addEventListener('mousedown', this.onMouse, true);                                                            // 122
			layer.addEventListener('mouseup', this.onMouse, true);                                                              // 123
		}                                                                                                                    // 124
                                                                                                                       // 125
		layer.addEventListener('click', this.onClick, true);                                                                 // 126
		layer.addEventListener('touchstart', this.onTouchStart, false);                                                      // 127
		layer.addEventListener('touchmove', this.onTouchMove, false);                                                        // 128
		layer.addEventListener('touchend', this.onTouchEnd, false);                                                          // 129
		layer.addEventListener('touchcancel', this.onTouchCancel, false);                                                    // 130
                                                                                                                       // 131
		// Hack is required for browsers that don't support Event#stopImmediatePropagation (e.g. Android 2)                  // 132
		// which is how FastClick normally stops click events bubbling to callbacks registered on the FastClick              // 133
		// layer when they are cancelled.                                                                                    // 134
		if (!Event.prototype.stopImmediatePropagation) {                                                                     // 135
			layer.removeEventListener = function(type, callback, capture) {                                                     // 136
				var rmv = Node.prototype.removeEventListener;                                                                      // 137
				if (type === 'click') {                                                                                            // 138
					rmv.call(layer, type, callback.hijacked || callback, capture);                                                    // 139
				} else {                                                                                                           // 140
					rmv.call(layer, type, callback, capture);                                                                         // 141
				}                                                                                                                  // 142
			};                                                                                                                  // 143
                                                                                                                       // 144
			layer.addEventListener = function(type, callback, capture) {                                                        // 145
				var adv = Node.prototype.addEventListener;                                                                         // 146
				if (type === 'click') {                                                                                            // 147
					adv.call(layer, type, callback.hijacked || (callback.hijacked = function(event) {                                 // 148
						if (!event.propagationStopped) {                                                                                 // 149
							callback(event);                                                                                                // 150
						}                                                                                                                // 151
					}), capture);                                                                                                     // 152
				} else {                                                                                                           // 153
					adv.call(layer, type, callback, capture);                                                                         // 154
				}                                                                                                                  // 155
			};                                                                                                                  // 156
		}                                                                                                                    // 157
                                                                                                                       // 158
		// If a handler is already declared in the element's onclick attribute, it will be fired before                      // 159
		// FastClick's onClick handler. Fix this by pulling out the user-defined handler function and                        // 160
		// adding it as listener.                                                                                            // 161
		if (typeof layer.onclick === 'function') {                                                                           // 162
                                                                                                                       // 163
			// Android browser on at least 3.2 requires a new reference to the function in layer.onclick                        // 164
			// - the old one won't work if passed to addEventListener directly.                                                 // 165
			oldOnClick = layer.onclick;                                                                                         // 166
			layer.addEventListener('click', function(event) {                                                                   // 167
				oldOnClick(event);                                                                                                 // 168
			}, false);                                                                                                          // 169
			layer.onclick = null;                                                                                               // 170
		}                                                                                                                    // 171
	}                                                                                                                     // 172
                                                                                                                       // 173
	/**                                                                                                                   // 174
	* Windows Phone 8.1 fakes user agent string to look like Android and iPhone.                                          // 175
	*                                                                                                                     // 176
	* @type boolean                                                                                                       // 177
	*/                                                                                                                    // 178
	var deviceIsWindowsPhone = navigator.userAgent.indexOf("Windows Phone") >= 0;                                         // 179
                                                                                                                       // 180
	/**                                                                                                                   // 181
	 * Android requires exceptions.                                                                                       // 182
	 *                                                                                                                    // 183
	 * @type boolean                                                                                                      // 184
	 */                                                                                                                   // 185
	var deviceIsAndroid = navigator.userAgent.indexOf('Android') > 0 && !deviceIsWindowsPhone;                            // 186
                                                                                                                       // 187
                                                                                                                       // 188
	/**                                                                                                                   // 189
	 * iOS requires exceptions.                                                                                           // 190
	 *                                                                                                                    // 191
	 * @type boolean                                                                                                      // 192
	 */                                                                                                                   // 193
	var deviceIsIOS = /iP(ad|hone|od)/.test(navigator.userAgent) && !deviceIsWindowsPhone;                                // 194
                                                                                                                       // 195
                                                                                                                       // 196
	/**                                                                                                                   // 197
	 * iOS 4 requires an exception for select elements.                                                                   // 198
	 *                                                                                                                    // 199
	 * @type boolean                                                                                                      // 200
	 */                                                                                                                   // 201
	var deviceIsIOS4 = deviceIsIOS && (/OS 4_\d(_\d)?/).test(navigator.userAgent);                                        // 202
                                                                                                                       // 203
                                                                                                                       // 204
	/**                                                                                                                   // 205
	 * iOS 6.0-7.* requires the target element to be manually derived                                                     // 206
	 *                                                                                                                    // 207
	 * @type boolean                                                                                                      // 208
	 */                                                                                                                   // 209
	var deviceIsIOSWithBadTarget = deviceIsIOS && (/OS [6-7]_\d/).test(navigator.userAgent);                              // 210
                                                                                                                       // 211
	/**                                                                                                                   // 212
	 * BlackBerry requires exceptions.                                                                                    // 213
	 *                                                                                                                    // 214
	 * @type boolean                                                                                                      // 215
	 */                                                                                                                   // 216
	var deviceIsBlackBerry10 = navigator.userAgent.indexOf('BB10') > 0;                                                   // 217
                                                                                                                       // 218
	/**                                                                                                                   // 219
	 * Determine whether a given element requires a native click.                                                         // 220
	 *                                                                                                                    // 221
	 * @param {EventTarget|Element} target Target DOM element                                                             // 222
	 * @returns {boolean} Returns true if the element needs a native click                                                // 223
	 */                                                                                                                   // 224
	FastClick.prototype.needsClick = function(target) {                                                                   // 225
		switch (target.nodeName.toLowerCase()) {                                                                             // 226
                                                                                                                       // 227
		// Don't send a synthetic click to disabled inputs (issue #62)                                                       // 228
		case 'button':                                                                                                       // 229
		case 'select':                                                                                                       // 230
		case 'textarea':                                                                                                     // 231
			if (target.disabled) {                                                                                              // 232
				return true;                                                                                                       // 233
			}                                                                                                                   // 234
                                                                                                                       // 235
			break;                                                                                                              // 236
		case 'input':                                                                                                        // 237
                                                                                                                       // 238
			// File inputs need real clicks on iOS 6 due to a browser bug (issue #68)                                           // 239
			if ((deviceIsIOS && target.type === 'file') || target.disabled) {                                                   // 240
				return true;                                                                                                       // 241
			}                                                                                                                   // 242
                                                                                                                       // 243
			break;                                                                                                              // 244
		case 'label':                                                                                                        // 245
		case 'iframe': // iOS8 homescreen apps can prevent events bubbling into frames                                       // 246
		case 'video':                                                                                                        // 247
			return true;                                                                                                        // 248
		}                                                                                                                    // 249
                                                                                                                       // 250
		return (/\bneedsclick\b/).test(target.className);                                                                    // 251
	};                                                                                                                    // 252
                                                                                                                       // 253
                                                                                                                       // 254
	/**                                                                                                                   // 255
	 * Determine whether a given element requires a call to focus to simulate click into element.                         // 256
	 *                                                                                                                    // 257
	 * @param {EventTarget|Element} target Target DOM element                                                             // 258
	 * @returns {boolean} Returns true if the element requires a call to focus to simulate native click.                  // 259
	 */                                                                                                                   // 260
	FastClick.prototype.needsFocus = function(target) {                                                                   // 261
		switch (target.nodeName.toLowerCase()) {                                                                             // 262
		case 'textarea':                                                                                                     // 263
			return true;                                                                                                        // 264
		case 'select':                                                                                                       // 265
			return !deviceIsAndroid;                                                                                            // 266
		case 'input':                                                                                                        // 267
			switch (target.type) {                                                                                              // 268
			case 'button':                                                                                                      // 269
			case 'checkbox':                                                                                                    // 270
			case 'file':                                                                                                        // 271
			case 'image':                                                                                                       // 272
			case 'radio':                                                                                                       // 273
			case 'submit':                                                                                                      // 274
				return false;                                                                                                      // 275
			}                                                                                                                   // 276
                                                                                                                       // 277
			// No point in attempting to focus disabled inputs                                                                  // 278
			return !target.disabled && !target.readOnly;                                                                        // 279
		default:                                                                                                             // 280
			return (/\bneedsfocus\b/).test(target.className);                                                                   // 281
		}                                                                                                                    // 282
	};                                                                                                                    // 283
                                                                                                                       // 284
                                                                                                                       // 285
	/**                                                                                                                   // 286
	 * Send a click event to the specified element.                                                                       // 287
	 *                                                                                                                    // 288
	 * @param {EventTarget|Element} targetElement                                                                         // 289
	 * @param {Event} event                                                                                               // 290
	 */                                                                                                                   // 291
	FastClick.prototype.sendClick = function(targetElement, event) {                                                      // 292
		var clickEvent, touch;                                                                                               // 293
                                                                                                                       // 294
		// On some Android devices activeElement needs to be blurred otherwise the synthetic click will have no effect (#24)
		if (document.activeElement && document.activeElement !== targetElement) {                                            // 296
			document.activeElement.blur();                                                                                      // 297
		}                                                                                                                    // 298
                                                                                                                       // 299
		touch = event.changedTouches[0];                                                                                     // 300
                                                                                                                       // 301
		// Synthesise a click event, with an extra attribute so it can be tracked                                            // 302
		clickEvent = document.createEvent('MouseEvents');                                                                    // 303
		clickEvent.initMouseEvent(this.determineEventType(targetElement), true, true, window, 1, touch.screenX, touch.screenY, touch.clientX, touch.clientY, false, false, false, false, 0, null);
		clickEvent.forwardedTouchEvent = true;                                                                               // 305
		targetElement.dispatchEvent(clickEvent);                                                                             // 306
	};                                                                                                                    // 307
                                                                                                                       // 308
	FastClick.prototype.determineEventType = function(targetElement) {                                                    // 309
                                                                                                                       // 310
		//Issue #159: Android Chrome Select Box does not open with a synthetic click event                                   // 311
		if (deviceIsAndroid && targetElement.tagName.toLowerCase() === 'select') {                                           // 312
			return 'mousedown';                                                                                                 // 313
		}                                                                                                                    // 314
                                                                                                                       // 315
		return 'click';                                                                                                      // 316
	};                                                                                                                    // 317
                                                                                                                       // 318
                                                                                                                       // 319
	/**                                                                                                                   // 320
	 * @param {EventTarget|Element} targetElement                                                                         // 321
	 */                                                                                                                   // 322
	FastClick.prototype.focus = function(targetElement) {                                                                 // 323
		var length;                                                                                                          // 324
                                                                                                                       // 325
		// Issue #160: on iOS 7, some input elements (e.g. date datetime month) throw a vague TypeError on setSelectionRange. These elements don't have an integer value for the selectionStart and selectionEnd properties, but unfortunately that can't be used for detection because accessing the properties also throws a TypeError. Just check the type instead. Filed as Apple bug #15122724.
		if (deviceIsIOS && targetElement.setSelectionRange && targetElement.type.indexOf('date') !== 0 && targetElement.type !== 'time' && targetElement.type !== 'month') {
			length = targetElement.value.length;                                                                                // 328
			targetElement.setSelectionRange(length, length);                                                                    // 329
		} else {                                                                                                             // 330
			targetElement.focus();                                                                                              // 331
		}                                                                                                                    // 332
	};                                                                                                                    // 333
                                                                                                                       // 334
                                                                                                                       // 335
	/**                                                                                                                   // 336
	 * Check whether the given target element is a child of a scrollable layer and if so, set a flag on it.               // 337
	 *                                                                                                                    // 338
	 * @param {EventTarget|Element} targetElement                                                                         // 339
	 */                                                                                                                   // 340
	FastClick.prototype.updateScrollParent = function(targetElement) {                                                    // 341
		var scrollParent, parentElement;                                                                                     // 342
                                                                                                                       // 343
		scrollParent = targetElement.fastClickScrollParent;                                                                  // 344
                                                                                                                       // 345
		// Attempt to discover whether the target element is contained within a scrollable layer. Re-check if the            // 346
		// target element was moved to another parent.                                                                       // 347
		if (!scrollParent || !scrollParent.contains(targetElement)) {                                                        // 348
			parentElement = targetElement;                                                                                      // 349
			do {                                                                                                                // 350
				if (parentElement.scrollHeight > parentElement.offsetHeight) {                                                     // 351
					scrollParent = parentElement;                                                                                     // 352
					targetElement.fastClickScrollParent = parentElement;                                                              // 353
					break;                                                                                                            // 354
				}                                                                                                                  // 355
                                                                                                                       // 356
				parentElement = parentElement.parentElement;                                                                       // 357
			} while (parentElement);                                                                                            // 358
		}                                                                                                                    // 359
                                                                                                                       // 360
		// Always update the scroll top tracker if possible.                                                                 // 361
		if (scrollParent) {                                                                                                  // 362
			scrollParent.fastClickLastScrollTop = scrollParent.scrollTop;                                                       // 363
		}                                                                                                                    // 364
	};                                                                                                                    // 365
                                                                                                                       // 366
                                                                                                                       // 367
	/**                                                                                                                   // 368
	 * @param {EventTarget} targetElement                                                                                 // 369
	 * @returns {Element|EventTarget}                                                                                     // 370
	 */                                                                                                                   // 371
	FastClick.prototype.getTargetElementFromEventTarget = function(eventTarget) {                                         // 372
                                                                                                                       // 373
		// On some older browsers (notably Safari on iOS 4.1 - see issue #56) the event target may be a text node.           // 374
		if (eventTarget.nodeType === Node.TEXT_NODE) {                                                                       // 375
			return eventTarget.parentNode;                                                                                      // 376
		}                                                                                                                    // 377
                                                                                                                       // 378
		return eventTarget;                                                                                                  // 379
	};                                                                                                                    // 380
                                                                                                                       // 381
                                                                                                                       // 382
	/**                                                                                                                   // 383
	 * On touch start, record the position and scroll offset.                                                             // 384
	 *                                                                                                                    // 385
	 * @param {Event} event                                                                                               // 386
	 * @returns {boolean}                                                                                                 // 387
	 */                                                                                                                   // 388
	FastClick.prototype.onTouchStart = function(event) {                                                                  // 389
		var targetElement, touch, selection;                                                                                 // 390
                                                                                                                       // 391
		// Ignore multiple touches, otherwise pinch-to-zoom is prevented if both fingers are on the FastClick element (issue #111).
		if (event.targetTouches.length > 1) {                                                                                // 393
			return true;                                                                                                        // 394
		}                                                                                                                    // 395
                                                                                                                       // 396
		targetElement = this.getTargetElementFromEventTarget(event.target);                                                  // 397
		touch = event.targetTouches[0];                                                                                      // 398
                                                                                                                       // 399
		if (deviceIsIOS) {                                                                                                   // 400
                                                                                                                       // 401
			// Only trusted events will deselect text on iOS (issue #49)                                                        // 402
			selection = window.getSelection();                                                                                  // 403
			if (selection.rangeCount && !selection.isCollapsed) {                                                               // 404
				return true;                                                                                                       // 405
			}                                                                                                                   // 406
                                                                                                                       // 407
			if (!deviceIsIOS4) {                                                                                                // 408
                                                                                                                       // 409
				// Weird things happen on iOS when an alert or confirm dialog is opened from a click event callback (issue #23):   // 410
				// when the user next taps anywhere else on the page, new touchstart and touchend events are dispatched            // 411
				// with the same identifier as the touch event that previously triggered the click that triggered the alert.       // 412
				// Sadly, there is an issue on iOS 4 that causes some normal touch events to have the same identifier as an        // 413
				// immediately preceeding touch event (issue #52), so this fix is unavailable on that platform.                    // 414
				// Issue 120: touch.identifier is 0 when Chrome dev tools 'Emulate touch events' is set with an iOS device UA string,
				// which causes all touch events to be ignored. As this block only applies to iOS, and iOS identifiers are always long,
				// random integers, it's safe to to continue if the identifier is 0 here.                                          // 417
				if (touch.identifier && touch.identifier === this.lastTouchIdentifier) {                                           // 418
					event.preventDefault();                                                                                           // 419
					return false;                                                                                                     // 420
				}                                                                                                                  // 421
                                                                                                                       // 422
				this.lastTouchIdentifier = touch.identifier;                                                                       // 423
                                                                                                                       // 424
				// If the target element is a child of a scrollable layer (using -webkit-overflow-scrolling: touch) and:           // 425
				// 1) the user does a fling scroll on the scrollable layer                                                         // 426
				// 2) the user stops the fling scroll with another tap                                                             // 427
				// then the event.target of the last 'touchend' event will be the element that was under the user's finger         // 428
				// when the fling scroll was started, causing FastClick to send a click event to that layer - unless a check       // 429
				// is made to ensure that a parent layer was not scrolled before sending a synthetic click (issue #42).            // 430
				this.updateScrollParent(targetElement);                                                                            // 431
			}                                                                                                                   // 432
		}                                                                                                                    // 433
                                                                                                                       // 434
		this.trackingClick = true;                                                                                           // 435
		this.trackingClickStart = event.timeStamp;                                                                           // 436
		this.targetElement = targetElement;                                                                                  // 437
                                                                                                                       // 438
		this.touchStartX = touch.pageX;                                                                                      // 439
		this.touchStartY = touch.pageY;                                                                                      // 440
                                                                                                                       // 441
		// Prevent phantom clicks on fast double-tap (issue #36)                                                             // 442
		if ((event.timeStamp - this.lastClickTime) < this.tapDelay) {                                                        // 443
			event.preventDefault();                                                                                             // 444
		}                                                                                                                    // 445
                                                                                                                       // 446
		return true;                                                                                                         // 447
	};                                                                                                                    // 448
                                                                                                                       // 449
                                                                                                                       // 450
	/**                                                                                                                   // 451
	 * Based on a touchmove event object, check whether the touch has moved past a boundary since it started.             // 452
	 *                                                                                                                    // 453
	 * @param {Event} event                                                                                               // 454
	 * @returns {boolean}                                                                                                 // 455
	 */                                                                                                                   // 456
	FastClick.prototype.touchHasMoved = function(event) {                                                                 // 457
		var touch = event.changedTouches[0], boundary = this.touchBoundary;                                                  // 458
                                                                                                                       // 459
		if (Math.abs(touch.pageX - this.touchStartX) > boundary || Math.abs(touch.pageY - this.touchStartY) > boundary) {    // 460
			return true;                                                                                                        // 461
		}                                                                                                                    // 462
                                                                                                                       // 463
		return false;                                                                                                        // 464
	};                                                                                                                    // 465
                                                                                                                       // 466
                                                                                                                       // 467
	/**                                                                                                                   // 468
	 * Update the last position.                                                                                          // 469
	 *                                                                                                                    // 470
	 * @param {Event} event                                                                                               // 471
	 * @returns {boolean}                                                                                                 // 472
	 */                                                                                                                   // 473
	FastClick.prototype.onTouchMove = function(event) {                                                                   // 474
		if (!this.trackingClick) {                                                                                           // 475
			return true;                                                                                                        // 476
		}                                                                                                                    // 477
                                                                                                                       // 478
		// If the touch has moved, cancel the click tracking                                                                 // 479
		if (this.targetElement !== this.getTargetElementFromEventTarget(event.target) || this.touchHasMoved(event)) {        // 480
			this.trackingClick = false;                                                                                         // 481
			this.targetElement = null;                                                                                          // 482
		}                                                                                                                    // 483
                                                                                                                       // 484
		return true;                                                                                                         // 485
	};                                                                                                                    // 486
                                                                                                                       // 487
                                                                                                                       // 488
	/**                                                                                                                   // 489
	 * Attempt to find the labelled control for the given label element.                                                  // 490
	 *                                                                                                                    // 491
	 * @param {EventTarget|HTMLLabelElement} labelElement                                                                 // 492
	 * @returns {Element|null}                                                                                            // 493
	 */                                                                                                                   // 494
	FastClick.prototype.findControl = function(labelElement) {                                                            // 495
                                                                                                                       // 496
		// Fast path for newer browsers supporting the HTML5 control attribute                                               // 497
		if (labelElement.control !== undefined) {                                                                            // 498
			return labelElement.control;                                                                                        // 499
		}                                                                                                                    // 500
                                                                                                                       // 501
		// All browsers under test that support touch events also support the HTML5 htmlFor attribute                        // 502
		if (labelElement.htmlFor) {                                                                                          // 503
			return document.getElementById(labelElement.htmlFor);                                                               // 504
		}                                                                                                                    // 505
                                                                                                                       // 506
		// If no for attribute exists, attempt to retrieve the first labellable descendant element                           // 507
		// the list of which is defined here: http://www.w3.org/TR/html5/forms.html#category-label                           // 508
		return labelElement.querySelector('button, input:not([type=hidden]), keygen, meter, output, progress, select, textarea');
	};                                                                                                                    // 510
                                                                                                                       // 511
                                                                                                                       // 512
	/**                                                                                                                   // 513
	 * On touch end, determine whether to send a click event at once.                                                     // 514
	 *                                                                                                                    // 515
	 * @param {Event} event                                                                                               // 516
	 * @returns {boolean}                                                                                                 // 517
	 */                                                                                                                   // 518
	FastClick.prototype.onTouchEnd = function(event) {                                                                    // 519
		var forElement, trackingClickStart, targetTagName, scrollParent, touch, targetElement = this.targetElement;          // 520
                                                                                                                       // 521
		if (!this.trackingClick) {                                                                                           // 522
			return true;                                                                                                        // 523
		}                                                                                                                    // 524
                                                                                                                       // 525
		// Prevent phantom clicks on fast double-tap (issue #36)                                                             // 526
		if ((event.timeStamp - this.lastClickTime) < this.tapDelay) {                                                        // 527
			this.cancelNextClick = true;                                                                                        // 528
			return true;                                                                                                        // 529
		}                                                                                                                    // 530
                                                                                                                       // 531
		if ((event.timeStamp - this.trackingClickStart) > this.tapTimeout) {                                                 // 532
			return true;                                                                                                        // 533
		}                                                                                                                    // 534
                                                                                                                       // 535
		// Reset to prevent wrong click cancel on input (issue #156).                                                        // 536
		this.cancelNextClick = false;                                                                                        // 537
                                                                                                                       // 538
		this.lastClickTime = event.timeStamp;                                                                                // 539
                                                                                                                       // 540
		trackingClickStart = this.trackingClickStart;                                                                        // 541
		this.trackingClick = false;                                                                                          // 542
		this.trackingClickStart = 0;                                                                                         // 543
                                                                                                                       // 544
		// On some iOS devices, the targetElement supplied with the event is invalid if the layer                            // 545
		// is performing a transition or scroll, and has to be re-detected manually. Note that                               // 546
		// for this to function correctly, it must be called *after* the event target is checked!                            // 547
		// See issue #57; also filed as rdar://13048589 .                                                                    // 548
		if (deviceIsIOSWithBadTarget) {                                                                                      // 549
			touch = event.changedTouches[0];                                                                                    // 550
                                                                                                                       // 551
			// In certain cases arguments of elementFromPoint can be negative, so prevent setting targetElement to null         // 552
			targetElement = document.elementFromPoint(touch.pageX - window.pageXOffset, touch.pageY - window.pageYOffset) || targetElement;
			targetElement.fastClickScrollParent = this.targetElement.fastClickScrollParent;                                     // 554
		}                                                                                                                    // 555
                                                                                                                       // 556
		targetTagName = targetElement.tagName.toLowerCase();                                                                 // 557
		if (targetTagName === 'label') {                                                                                     // 558
			forElement = this.findControl(targetElement);                                                                       // 559
			if (forElement) {                                                                                                   // 560
				this.focus(targetElement);                                                                                         // 561
				if (deviceIsAndroid) {                                                                                             // 562
					return false;                                                                                                     // 563
				}                                                                                                                  // 564
                                                                                                                       // 565
				targetElement = forElement;                                                                                        // 566
			}                                                                                                                   // 567
		} else if (this.needsFocus(targetElement)) {                                                                         // 568
                                                                                                                       // 569
			// Case 1: If the touch started a while ago (best guess is 100ms based on tests for issue #36) then focus will be triggered anyway. Return early and unset the target element reference so that the subsequent click will be allowed through.
			// Case 2: Without this exception for input elements tapped when the document is contained in an iframe, then any inputted text won't be visible even though the value attribute is updated as the user types (issue #37).
			if ((event.timeStamp - trackingClickStart) > 100 || (deviceIsIOS && window.top !== window && targetTagName === 'input')) {
				this.targetElement = null;                                                                                         // 573
				return false;                                                                                                      // 574
			}                                                                                                                   // 575
                                                                                                                       // 576
			this.focus(targetElement);                                                                                          // 577
			this.sendClick(targetElement, event);                                                                               // 578
                                                                                                                       // 579
			// Select elements need the event to go through on iOS 4, otherwise the selector menu won't open.                   // 580
			// Also this breaks opening selects when VoiceOver is active on iOS6, iOS7 (and possibly others)                    // 581
			if (!deviceIsIOS || targetTagName !== 'select') {                                                                   // 582
				this.targetElement = null;                                                                                         // 583
				event.preventDefault();                                                                                            // 584
			}                                                                                                                   // 585
                                                                                                                       // 586
			return false;                                                                                                       // 587
		}                                                                                                                    // 588
                                                                                                                       // 589
		if (deviceIsIOS && !deviceIsIOS4) {                                                                                  // 590
                                                                                                                       // 591
			// Don't send a synthetic click event if the target element is contained within a parent layer that was scrolled    // 592
			// and this tap is being used to stop the scrolling (usually initiated by a fling - issue #42).                     // 593
			scrollParent = targetElement.fastClickScrollParent;                                                                 // 594
			if (scrollParent && scrollParent.fastClickLastScrollTop !== scrollParent.scrollTop) {                               // 595
				return true;                                                                                                       // 596
			}                                                                                                                   // 597
		}                                                                                                                    // 598
                                                                                                                       // 599
		// Prevent the actual click from going though - unless the target node is marked as requiring                        // 600
		// real clicks or if it is in the whitelist in which case only non-programmatic clicks are permitted.                // 601
		if (!this.needsClick(targetElement)) {                                                                               // 602
			event.preventDefault();                                                                                             // 603
			this.sendClick(targetElement, event);                                                                               // 604
		}                                                                                                                    // 605
                                                                                                                       // 606
		return false;                                                                                                        // 607
	};                                                                                                                    // 608
                                                                                                                       // 609
                                                                                                                       // 610
	/**                                                                                                                   // 611
	 * On touch cancel, stop tracking the click.                                                                          // 612
	 *                                                                                                                    // 613
	 * @returns {void}                                                                                                    // 614
	 */                                                                                                                   // 615
	FastClick.prototype.onTouchCancel = function() {                                                                      // 616
		this.trackingClick = false;                                                                                          // 617
		this.targetElement = null;                                                                                           // 618
	};                                                                                                                    // 619
                                                                                                                       // 620
                                                                                                                       // 621
	/**                                                                                                                   // 622
	 * Determine mouse events which should be permitted.                                                                  // 623
	 *                                                                                                                    // 624
	 * @param {Event} event                                                                                               // 625
	 * @returns {boolean}                                                                                                 // 626
	 */                                                                                                                   // 627
	FastClick.prototype.onMouse = function(event) {                                                                       // 628
                                                                                                                       // 629
		// If a target element was never set (because a touch event was never fired) allow the event                         // 630
		if (!this.targetElement) {                                                                                           // 631
			return true;                                                                                                        // 632
		}                                                                                                                    // 633
                                                                                                                       // 634
		if (event.forwardedTouchEvent) {                                                                                     // 635
			return true;                                                                                                        // 636
		}                                                                                                                    // 637
                                                                                                                       // 638
		// Programmatically generated events targeting a specific element should be permitted                                // 639
		if (!event.cancelable) {                                                                                             // 640
			return true;                                                                                                        // 641
		}                                                                                                                    // 642
                                                                                                                       // 643
		// Derive and check the target element to see whether the mouse event needs to be permitted;                         // 644
		// unless explicitly enabled, prevent non-touch click events from triggering actions,                                // 645
		// to prevent ghost/doubleclicks.                                                                                    // 646
		if (!this.needsClick(this.targetElement) || this.cancelNextClick) {                                                  // 647
                                                                                                                       // 648
			// Prevent any user-added listeners declared on FastClick element from being fired.                                 // 649
			if (event.stopImmediatePropagation) {                                                                               // 650
				event.stopImmediatePropagation();                                                                                  // 651
			} else {                                                                                                            // 652
                                                                                                                       // 653
				// Part of the hack for browsers that don't support Event#stopImmediatePropagation (e.g. Android 2)                // 654
				event.propagationStopped = true;                                                                                   // 655
			}                                                                                                                   // 656
                                                                                                                       // 657
			// Cancel the event                                                                                                 // 658
			event.stopPropagation();                                                                                            // 659
			event.preventDefault();                                                                                             // 660
                                                                                                                       // 661
			return false;                                                                                                       // 662
		}                                                                                                                    // 663
                                                                                                                       // 664
		// If the mouse event is permitted, return true for the action to go through.                                        // 665
		return true;                                                                                                         // 666
	};                                                                                                                    // 667
                                                                                                                       // 668
                                                                                                                       // 669
	/**                                                                                                                   // 670
	 * On actual clicks, determine whether this is a touch-generated click, a click action occurring                      // 671
	 * naturally after a delay after a touch (which needs to be cancelled to avoid duplication), or                       // 672
	 * an actual click which should be permitted.                                                                         // 673
	 *                                                                                                                    // 674
	 * @param {Event} event                                                                                               // 675
	 * @returns {boolean}                                                                                                 // 676
	 */                                                                                                                   // 677
	FastClick.prototype.onClick = function(event) {                                                                       // 678
		var permitted;                                                                                                       // 679
                                                                                                                       // 680
		// It's possible for another FastClick-like library delivered with third-party code to fire a click event before FastClick does (issue #44). In that case, set the click-tracking flag back to false and return early. This will cause onTouchEnd to return early.
		if (this.trackingClick) {                                                                                            // 682
			this.targetElement = null;                                                                                          // 683
			this.trackingClick = false;                                                                                         // 684
			return true;                                                                                                        // 685
		}                                                                                                                    // 686
                                                                                                                       // 687
		// Very odd behaviour on iOS (issue #18): if a submit element is present inside a form and the user hits enter in the iOS simulator or clicks the Go button on the pop-up OS keyboard the a kind of 'fake' click event will be triggered with the submit-type input element as the target.
		if (event.target.type === 'submit' && event.detail === 0) {                                                          // 689
			return true;                                                                                                        // 690
		}                                                                                                                    // 691
                                                                                                                       // 692
		permitted = this.onMouse(event);                                                                                     // 693
                                                                                                                       // 694
		// Only unset targetElement if the click is not permitted. This will ensure that the check for !targetElement in onMouse fails and the browser's click doesn't go through.
		if (!permitted) {                                                                                                    // 696
			this.targetElement = null;                                                                                          // 697
		}                                                                                                                    // 698
                                                                                                                       // 699
		// If clicks are permitted, return true for the action to go through.                                                // 700
		return permitted;                                                                                                    // 701
	};                                                                                                                    // 702
                                                                                                                       // 703
                                                                                                                       // 704
	/**                                                                                                                   // 705
	 * Remove all FastClick's event listeners.                                                                            // 706
	 *                                                                                                                    // 707
	 * @returns {void}                                                                                                    // 708
	 */                                                                                                                   // 709
	FastClick.prototype.destroy = function() {                                                                            // 710
		var layer = this.layer;                                                                                              // 711
                                                                                                                       // 712
		if (deviceIsAndroid) {                                                                                               // 713
			layer.removeEventListener('mouseover', this.onMouse, true);                                                         // 714
			layer.removeEventListener('mousedown', this.onMouse, true);                                                         // 715
			layer.removeEventListener('mouseup', this.onMouse, true);                                                           // 716
		}                                                                                                                    // 717
                                                                                                                       // 718
		layer.removeEventListener('click', this.onClick, true);                                                              // 719
		layer.removeEventListener('touchstart', this.onTouchStart, false);                                                   // 720
		layer.removeEventListener('touchmove', this.onTouchMove, false);                                                     // 721
		layer.removeEventListener('touchend', this.onTouchEnd, false);                                                       // 722
		layer.removeEventListener('touchcancel', this.onTouchCancel, false);                                                 // 723
	};                                                                                                                    // 724
                                                                                                                       // 725
                                                                                                                       // 726
	/**                                                                                                                   // 727
	 * Check whether FastClick is needed.                                                                                 // 728
	 *                                                                                                                    // 729
	 * @param {Element} layer The layer to listen on                                                                      // 730
	 */                                                                                                                   // 731
	FastClick.notNeeded = function(layer) {                                                                               // 732
		var metaViewport;                                                                                                    // 733
		var chromeVersion;                                                                                                   // 734
		var blackberryVersion;                                                                                               // 735
		var firefoxVersion;                                                                                                  // 736
                                                                                                                       // 737
		// Devices that don't support touch don't need FastClick                                                             // 738
		if (typeof window.ontouchstart === 'undefined') {                                                                    // 739
			return true;                                                                                                        // 740
		}                                                                                                                    // 741
                                                                                                                       // 742
		// Chrome version - zero for other browsers                                                                          // 743
		chromeVersion = +(/Chrome\/([0-9]+)/.exec(navigator.userAgent) || [,0])[1];                                          // 744
                                                                                                                       // 745
		if (chromeVersion) {                                                                                                 // 746
                                                                                                                       // 747
			if (deviceIsAndroid) {                                                                                              // 748
				metaViewport = document.querySelector('meta[name=viewport]');                                                      // 749
                                                                                                                       // 750
				if (metaViewport) {                                                                                                // 751
					// Chrome on Android with user-scalable="no" doesn't need FastClick (issue #89)                                   // 752
					if (metaViewport.content.indexOf('user-scalable=no') !== -1) {                                                    // 753
						return true;                                                                                                     // 754
					}                                                                                                                 // 755
					// Chrome 32 and above with width=device-width or less don't need FastClick                                       // 756
					if (chromeVersion > 31 && document.documentElement.scrollWidth <= window.outerWidth) {                            // 757
						return true;                                                                                                     // 758
					}                                                                                                                 // 759
				}                                                                                                                  // 760
                                                                                                                       // 761
			// Chrome desktop doesn't need FastClick (issue #15)                                                                // 762
			} else {                                                                                                            // 763
				return true;                                                                                                       // 764
			}                                                                                                                   // 765
		}                                                                                                                    // 766
                                                                                                                       // 767
		if (deviceIsBlackBerry10) {                                                                                          // 768
			blackberryVersion = navigator.userAgent.match(/Version\/([0-9]*)\.([0-9]*)/);                                       // 769
                                                                                                                       // 770
			// BlackBerry 10.3+ does not require Fastclick library.                                                             // 771
			// https://github.com/ftlabs/fastclick/issues/251                                                                   // 772
			if (blackberryVersion[1] >= 10 && blackberryVersion[2] >= 3) {                                                      // 773
				metaViewport = document.querySelector('meta[name=viewport]');                                                      // 774
                                                                                                                       // 775
				if (metaViewport) {                                                                                                // 776
					// user-scalable=no eliminates click delay.                                                                       // 777
					if (metaViewport.content.indexOf('user-scalable=no') !== -1) {                                                    // 778
						return true;                                                                                                     // 779
					}                                                                                                                 // 780
					// width=device-width (or less than device-width) eliminates click delay.                                         // 781
					if (document.documentElement.scrollWidth <= window.outerWidth) {                                                  // 782
						return true;                                                                                                     // 783
					}                                                                                                                 // 784
				}                                                                                                                  // 785
			}                                                                                                                   // 786
		}                                                                                                                    // 787
                                                                                                                       // 788
		// IE10 with -ms-touch-action: none or manipulation, which disables double-tap-to-zoom (issue #97)                   // 789
		if (layer.style.msTouchAction === 'none' || layer.style.touchAction === 'manipulation') {                            // 790
			return true;                                                                                                        // 791
		}                                                                                                                    // 792
                                                                                                                       // 793
		// Firefox version - zero for other browsers                                                                         // 794
		firefoxVersion = +(/Firefox\/([0-9]+)/.exec(navigator.userAgent) || [,0])[1];                                        // 795
                                                                                                                       // 796
		if (firefoxVersion >= 27) {                                                                                          // 797
			// Firefox 27+ does not have tap delay if the content is not zoomable - https://bugzilla.mozilla.org/show_bug.cgi?id=922896
                                                                                                                       // 799
			metaViewport = document.querySelector('meta[name=viewport]');                                                       // 800
			if (metaViewport && (metaViewport.content.indexOf('user-scalable=no') !== -1 || document.documentElement.scrollWidth <= window.outerWidth)) {
				return true;                                                                                                       // 802
			}                                                                                                                   // 803
		}                                                                                                                    // 804
                                                                                                                       // 805
		// IE11: prefixed -ms-touch-action is no longer supported and it's recomended to use non-prefixed version            // 806
		// http://msdn.microsoft.com/en-us/library/windows/apps/Hh767313.aspx                                                // 807
		if (layer.style.touchAction === 'none' || layer.style.touchAction === 'manipulation') {                              // 808
			return true;                                                                                                        // 809
		}                                                                                                                    // 810
                                                                                                                       // 811
		return false;                                                                                                        // 812
	};                                                                                                                    // 813
                                                                                                                       // 814
                                                                                                                       // 815
	/**                                                                                                                   // 816
	 * Factory method for creating a FastClick object                                                                     // 817
	 *                                                                                                                    // 818
	 * @param {Element} layer The layer to listen on                                                                      // 819
	 * @param {Object} [options={}] The options to override the defaults                                                  // 820
	 */                                                                                                                   // 821
	FastClick.attach = function(layer, options) {                                                                         // 822
		return new FastClick(layer, options);                                                                                // 823
	};                                                                                                                    // 824
                                                                                                                       // 825
                                                                                                                       // 826
	if (typeof define === 'function' && typeof define.amd === 'object' && define.amd) {                                   // 827
                                                                                                                       // 828
		// AMD. Register as an anonymous module.                                                                             // 829
		define(function() {                                                                                                  // 830
			return FastClick;                                                                                                   // 831
		});                                                                                                                  // 832
	} else if (typeof module !== 'undefined' && module.exports) {                                                         // 833
		module.exports = FastClick.attach;                                                                                   // 834
		module.exports.FastClick = FastClick;                                                                                // 835
	} else {                                                                                                              // 836
		window.FastClick = FastClick;                                                                                        // 837
	}                                                                                                                     // 838
                                                                                                                       // 839
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// packages/fastclick/post.js                                                                                          //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// This exports object was created in pre.js.  Now copy the 'FastClick' object                                         // 1
// from it into the package-scope variable `FastClick`, which will get exported.                                       // 2
                                                                                                                       // 3
FastClick = module.exports.FastClick;                                                                                  // 4
                                                                                                                       // 5
Meteor.startup(function () {                                                                                           // 6
  FastClick.attach(document.body);                                                                                     // 7
});                                                                                                                    // 8
                                                                                                                       // 9
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.fastclick = {}, {
  FastClick: FastClick
});

})();
