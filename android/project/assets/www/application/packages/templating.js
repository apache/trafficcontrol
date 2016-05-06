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
var Spacebars = Package.spacebars.Spacebars;
var HTML = Package.htmljs.HTML;

/* Package-scope variables */
var Template;

(function(){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                              //
// packages/templating/templating.js                                                                            //
//                                                                                                              //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                //
                                                                                                                // 1
// Packages and apps add templates on to this object.                                                           // 2
                                                                                                                // 3
/**                                                                                                             // 4
 * @summary The class for defining templates                                                                    // 5
 * @class                                                                                                       // 6
 * @instanceName Template.myTemplate                                                                            // 7
 */                                                                                                             // 8
Template = Blaze.Template;                                                                                      // 9
                                                                                                                // 10
var RESERVED_TEMPLATE_NAMES = "__proto__ name".split(" ");                                                      // 11
                                                                                                                // 12
// Check for duplicate template names and illegal names that won't work.                                        // 13
Template.__checkName = function (name) {                                                                        // 14
  // Some names can't be used for Templates. These include:                                                     // 15
  //  - Properties Blaze sets on the Template object.                                                           // 16
  //  - Properties that some browsers don't let the code to set.                                                // 17
  //    These are specified in RESERVED_TEMPLATE_NAMES.                                                         // 18
  if (name in Template || _.contains(RESERVED_TEMPLATE_NAMES, name)) {                                          // 19
    if ((Template[name] instanceof Template) && name !== "body")                                                // 20
      throw new Error("There are multiple templates named '" + name + "'. Each template needs a unique name.");
    throw new Error("This template name is reserved: " + name);                                                 // 22
  }                                                                                                             // 23
};                                                                                                              // 24
                                                                                                                // 25
// XXX COMPAT WITH 0.8.3                                                                                        // 26
Template.__define__ = function (name, renderFunc) {                                                             // 27
  Template.__checkName(name);                                                                                   // 28
  Template[name] = new Template("Template." + name, renderFunc);                                                // 29
  // Exempt packages built pre-0.9.0 from warnings about using old                                              // 30
  // helper syntax, because we can.  It's not very useful to get a                                              // 31
  // warning about someone else's code (like a package on Atmosphere),                                          // 32
  // and this should at least put a bit of a dent in number of warnings                                         // 33
  // that come from packages that haven't been updated lately.                                                  // 34
  Template[name]._NOWARN_OLDSTYLE_HELPERS = true;                                                               // 35
};                                                                                                              // 36
                                                                                                                // 37
// Define a template `Template.body` that renders its                                                           // 38
// `contentRenderFuncs`.  `<body>` tags (of which there may be                                                  // 39
// multiple) will have their contents added to it.                                                              // 40
                                                                                                                // 41
/**                                                                                                             // 42
 * @summary The [template object](#templates_api) representing your `<body>`                                    // 43
 * tag.                                                                                                         // 44
 * @locus Client                                                                                                // 45
 */                                                                                                             // 46
Template.body = new Template('body', function () {                                                              // 47
  var view = this;                                                                                              // 48
  return _.map(Template.body.contentRenderFuncs, function (func) {                                              // 49
    return func.apply(view);                                                                                    // 50
  });                                                                                                           // 51
});                                                                                                             // 52
Template.body.contentRenderFuncs = []; // array of Blaze.Views                                                  // 53
Template.body.view = null;                                                                                      // 54
                                                                                                                // 55
Template.body.addContent = function (renderFunc) {                                                              // 56
  Template.body.contentRenderFuncs.push(renderFunc);                                                            // 57
};                                                                                                              // 58
                                                                                                                // 59
// This function does not use `this` and so it may be called                                                    // 60
// as `Meteor.startup(Template.body.renderIntoDocument)`.                                                       // 61
Template.body.renderToDocument = function () {                                                                  // 62
  // Only do it once.                                                                                           // 63
  if (Template.body.view)                                                                                       // 64
    return;                                                                                                     // 65
                                                                                                                // 66
  var view = Blaze.render(Template.body, document.body);                                                        // 67
  Template.body.view = view;                                                                                    // 68
};                                                                                                              // 69
                                                                                                                // 70
// XXX COMPAT WITH 0.9.0                                                                                        // 71
UI.body = Template.body;                                                                                        // 72
                                                                                                                // 73
// XXX COMPAT WITH 0.9.0                                                                                        // 74
// (<body> tags in packages built with 0.9.0)                                                                   // 75
Template.__body__ = Template.body;                                                                              // 76
Template.__body__.__contentParts = Template.body.contentViews;                                                  // 77
Template.__body__.__instantiate = Template.body.renderToDocument;                                               // 78
                                                                                                                // 79
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                              //
// packages/templating/template.dynamic.js                                                                      //
//                                                                                                              //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                //
                                                                                                                // 1
Template.__checkName("__dynamic");                                                                              // 2
Template["__dynamic"] = new Template("Template.__dynamic", (function() {                                        // 3
  var view = this;                                                                                              // 4
  return [ Blaze.View("lookup:checkContext", function() {                                                       // 5
    return Spacebars.mustache(view.lookup("checkContext"));                                                     // 6
  }), "\n  ", Blaze.If(function() {                                                                             // 7
    return Spacebars.call(view.lookup("dataContextPresent"));                                                   // 8
  }, function() {                                                                                               // 9
    return [ "\n    ", Spacebars.include(view.lookupTemplate("__dynamicWithDataContext"), function() {          // 10
      return Blaze._InOuterTemplateScope(view, function() {                                                     // 11
        return Spacebars.include(function() {                                                                   // 12
          return Spacebars.call(view.templateContentBlock);                                                     // 13
        });                                                                                                     // 14
      });                                                                                                       // 15
    }), "\n  " ];                                                                                               // 16
  }, function() {                                                                                               // 17
    return [ "\n    \n    ", Blaze._TemplateWith(function() {                                                   // 18
      return {                                                                                                  // 19
        template: Spacebars.call(view.lookup("template")),                                                      // 20
        data: Spacebars.call(view.lookup(".."))                                                                 // 21
      };                                                                                                        // 22
    }, function() {                                                                                             // 23
      return Spacebars.include(view.lookupTemplate("__dynamicWithDataContext"), function() {                    // 24
        return Blaze._InOuterTemplateScope(view, function() {                                                   // 25
          return Spacebars.include(function() {                                                                 // 26
            return Spacebars.call(view.templateContentBlock);                                                   // 27
          });                                                                                                   // 28
        });                                                                                                     // 29
      });                                                                                                       // 30
    }), "\n  " ];                                                                                               // 31
  }) ];                                                                                                         // 32
}));                                                                                                            // 33
                                                                                                                // 34
Template.__checkName("__dynamicWithDataContext");                                                               // 35
Template["__dynamicWithDataContext"] = new Template("Template.__dynamicWithDataContext", (function() {          // 36
  var view = this;                                                                                              // 37
  return Spacebars.With(function() {                                                                            // 38
    return Spacebars.dataMustache(view.lookup("chooseTemplate"), view.lookup("template"));                      // 39
  }, function() {                                                                                               // 40
    return [ "\n    \n    ", Blaze._TemplateWith(function() {                                                   // 41
      return Spacebars.call(Spacebars.dot(view.lookup(".."), "data"));                                          // 42
    }, function() {                                                                                             // 43
      return Spacebars.include(view.lookupTemplate(".."), function() {                                          // 44
        return Blaze._InOuterTemplateScope(view, function() {                                                   // 45
          return Spacebars.include(function() {                                                                 // 46
            return Spacebars.call(view.templateContentBlock);                                                   // 47
          });                                                                                                   // 48
        });                                                                                                     // 49
      });                                                                                                       // 50
    }), "\n  " ];                                                                                               // 51
  });                                                                                                           // 52
}));                                                                                                            // 53
                                                                                                                // 54
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);






(function(){

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                              //
// packages/templating/dynamic.js                                                                               //
//                                                                                                              //
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                //
var Template = Blaze.Template;                                                                                  // 1
                                                                                                                // 2
/**                                                                                                             // 3
 * @isTemplate true                                                                                             // 4
 * @memberOf Template                                                                                           // 5
 * @function dynamic                                                                                            // 6
 * @summary Choose a template to include dynamically, by name.                                                  // 7
 * @locus Templates                                                                                             // 8
 * @param {String} template The name of the template to include.                                                // 9
 * @param {Object} [data] Optional. The data context in which to include the                                    // 10
 * template.                                                                                                    // 11
 */                                                                                                             // 12
                                                                                                                // 13
Template.__dynamicWithDataContext.helpers({                                                                     // 14
  chooseTemplate: function (name) {                                                                             // 15
    return Blaze._getTemplate(name, function () {                                                               // 16
      return Template.instance();                                                                               // 17
    });                                                                                                         // 18
  }                                                                                                             // 19
});                                                                                                             // 20
                                                                                                                // 21
Template.__dynamic.helpers({                                                                                    // 22
  dataContextPresent: function () {                                                                             // 23
    return _.has(this, "data");                                                                                 // 24
  },                                                                                                            // 25
  checkContext: function () {                                                                                   // 26
    if (! _.has(this, "template")) {                                                                            // 27
      throw new Error("Must specify name in the 'template' argument " +                                         // 28
                      "to {{> Template.dynamic}}.");                                                            // 29
    }                                                                                                           // 30
                                                                                                                // 31
    _.each(this, function (v, k) {                                                                              // 32
      if (k !== "template" && k !== "data") {                                                                   // 33
        throw new Error("Invalid argument to {{> Template.dynamic}}: " +                                        // 34
                        k);                                                                                     // 35
      }                                                                                                         // 36
    });                                                                                                         // 37
  }                                                                                                             // 38
});                                                                                                             // 39
                                                                                                                // 40
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}).call(this);


/* Exports */
if (typeof Package === 'undefined') Package = {};
(function (pkg, symbols) {
  for (var s in symbols)
    (s in pkg) || (pkg[s] = symbols[s]);
})(Package.templating = {}, {
  Template: Template
});

})();
