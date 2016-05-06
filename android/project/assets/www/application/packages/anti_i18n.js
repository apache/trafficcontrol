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
var Blaze = Package.blaze.Blaze;
var UI = Package.blaze.UI;
var Handlebars = Package.blaze.Handlebars;
var Tracker = Package.tracker.Tracker;
var Deps = Package.tracker.Deps;
var HTML = Package.htmljs.HTML;

/* Package-scope variables */
var i18n;

(function(){

//////////////////////////////////////////////////////////////////////////////
//                                                                          //
// packages/anti_i18n/packages/anti_i18n.js                                 //
//                                                                          //
//////////////////////////////////////////////////////////////////////////////
                                                                            //
(function () {                                                              // 1
                                                                            // 2
///////////////////////////////////////////////////////////////////////     // 3
//                                                                   //     // 4
// packages/anti:i18n/i18n.js                                        //     // 5
//                                                                   //     // 6
///////////////////////////////////////////////////////////////////////     // 7
                                                                     //     // 8
/*                                                                   // 1   // 9
  just-i18n package for Meteor.js                                    // 2   // 10
  author: Hubert OG <hubert@orlikarnia.com>                          // 3   // 11
*/                                                                   // 4   // 12
                                                                     // 5   // 13
                                                                     // 6   // 14
var maps            = {};                                            // 7   // 15
var language        = '';                                            // 8   // 16
var defaultLanguage = 'en';                                          // 9   // 17
var missingTemplate = '';                                            // 10  // 18
var showMissing     = false;                                         // 11  // 19
var dep             = new Deps.Dependency();                         // 12  // 20
                                                                     // 13  // 21
                                                                     // 14  // 22
/*                                                                   // 15  // 23
  Convert key to internationalized version                           // 16  // 24
*/                                                                   // 17  // 25
i18n = function() {                                                  // 18  // 26
  dep.depend();                                                      // 19  // 27
                                                                     // 20  // 28
  var label;                                                         // 21  // 29
  var args = _.toArray(arguments);                                   // 22  // 30
                                                                     // 23  // 31
  /* remove extra parameter added by blaze */                        // 24  // 32
  if(typeof args[args.length-1] === 'object') {                      // 25  // 33
    args.pop();                                                      // 26  // 34
  }                                                                  // 27  // 35
                                                                     // 28  // 36
  var label = args[0];                                               // 29  // 37
  args.shift();                                                      // 30  // 38
                                                                     // 31  // 39
                                                                     // 32  // 40
  if(typeof label !== 'string') return '';                           // 33  // 41
  var str = (maps[language] && maps[language][label]) ||             // 34  // 42
         (maps[defaultLanguage] && maps[defaultLanguage][label]) ||  // 35  // 43
         (showMissing && _.template(missingTemplate, {language: language, defaultLanguage: defaultLanguage, label: label})) ||
         '';                                                         // 37  // 45
  str = replaceWithParams(str, args)                                 // 38  // 46
  return str;                                                        // 39  // 47
};                                                                   // 40  // 48
                                                                     // 41  // 49
/*                                                                   // 42  // 50
  Register handlebars helper                                         // 43  // 51
*/                                                                   // 44  // 52
if(Meteor.isClient) {                                                // 45  // 53
  if(UI) {                                                           // 46  // 54
    UI.registerHelper('i18n', function () {                          // 47  // 55
      return i18n.apply(this, arguments);                            // 48  // 56
    });                                                              // 49  // 57
  } else if(Handlebars) {                                            // 50  // 58
    Handlebars.registerHelper('i18n', function () {                  // 51  // 59
      return i18n.apply(this, arguments);                            // 52  // 60
    });                                                              // 53  // 61
  }                                                                  // 54  // 62
}                                                                    // 55  // 63
                                                                     // 56  // 64
function replaceWithParams(string, params) {                         // 57  // 65
  var formatted = string;                                            // 58  // 66
  params.forEach(function(param , index){                            // 59  // 67
    var pos = index + 1;                                             // 60  // 68
    formatted = formatted.replace("{$" + pos + "}", param);          // 61  // 69
  });                                                                // 62  // 70
                                                                     // 63  // 71
  return formatted;                                                  // 64  // 72
};                                                                   // 65  // 73
                                                                     // 66  // 74
/*                                                                   // 67  // 75
  Settings                                                           // 68  // 76
*/                                                                   // 69  // 77
i18n.setLanguage = function(lng) {                                   // 70  // 78
  language = lng;                                                    // 71  // 79
  dep.changed();                                                     // 72  // 80
};                                                                   // 73  // 81
                                                                     // 74  // 82
i18n.setDefaultLanguage = function(lng) {                            // 75  // 83
  defaultLanguage = lng;                                             // 76  // 84
  dep.changed();                                                     // 77  // 85
};                                                                   // 78  // 86
                                                                     // 79  // 87
i18n.getLanguage = function() {                                      // 80  // 88
  dep.depend();                                                      // 81  // 89
  return language;                                                   // 82  // 90
};                                                                   // 83  // 91
                                                                     // 84  // 92
i18n.showMissing = function(template) {                              // 85  // 93
  if(template) {                                                     // 86  // 94
    if(typeof template === 'string') {                               // 87  // 95
      missingTemplate = template;                                    // 88  // 96
    } else {                                                         // 89  // 97
      missingTemplate = '[<%= label %>]';                            // 90  // 98
    }                                                                // 91  // 99
    showMissing = true;                                              // 92  // 100
  } else {                                                           // 93  // 101
    missingTemplate = '';                                            // 94  // 102
    showMissing = false;                                             // 95  // 103
  }                                                                  // 96  // 104
};                                                                   // 97  // 105
                                                                     // 98  // 106
/*                                                                   // 99  // 107
  Register map                                                       // 100
*/                                                                   // 101
i18n.map = function(language, map) {                                 // 102
  if(!maps[language]) maps[language] = {};                           // 103
  registerMap(language, '', false, map);                             // 104
  dep.changed();                                                     // 105
};                                                                   // 106
                                                                     // 107
var registerMap = function(language, prefix, dot, map) {             // 108
  if(typeof map === 'string') {                                      // 109
    maps[language][prefix] = map;                                    // 110
  } else if(typeof map === 'object') {                               // 111
    if(dot) prefix = prefix + '.';                                   // 112
    _.each(map, function(value, key) {                               // 113
      registerMap(language, prefix + key, true, value);              // 114
    });                                                              // 115
  }                                                                  // 116
};                                                                   // 117
                                                                     // 118
                                                                     // 119
///////////////////////////////////////////////////////////////////////     // 128
                                                                            // 129
}).call(this);                                                              // 130
                                                                            // 131
//////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package['anti:i18n'] = {}, {
  i18n: i18n
});

})();
