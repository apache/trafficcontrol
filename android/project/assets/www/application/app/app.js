var require = meteorInstall({"client":{"templates":{"application":{"template.layout.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/application/template.layout.js                                                                     //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("layout");                                                                                        // 2
Template["layout"] = new Template("Template.layout", (function() {                                                     // 3
  var view = this;                                                                                                     // 4
  return HTML.DIV({                                                                                                    // 5
    "class": "container"                                                                                               // 6
  }, "\n    ", Spacebars.include(view.lookupTemplate("header")), "\n    ", Spacebars.include(view.lookupTemplate("errors")), "\n    ", HTML.DIV({
    id: "main"                                                                                                         // 8
  }, "\n      ", Spacebars.include(view.lookupTemplate("yield")), "\n    "), "\n  ");                                  // 9
}));                                                                                                                   // 10
                                                                                                                       // 11
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.not_found.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/application/template.not_found.js                                                                  //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("notFound");                                                                                      // 2
Template["notFound"] = new Template("Template.notFound", (function() {                                                 // 3
  var view = this;                                                                                                     // 4
  return HTML.Raw('<div class="not-found page jumbotron">\n    <h2>404</h2>\n    <p>Sorry, we couldn\'t find a page at this address.</p>\n  </div>');
}));                                                                                                                   // 6
                                                                                                                       // 7
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"layout.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/application/layout.js                                                                              //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Template.layout.onRendered(function () {                                                                               // 1
  this.find('#main')._uihooks = {                                                                                      // 2
    insertElement: function () {                                                                                       // 3
      function insertElement(node, next) {                                                                             // 3
        $(node).hide().insertBefore(next).fadeIn();                                                                    // 4
      }                                                                                                                //
                                                                                                                       //
      return insertElement;                                                                                            //
    }(),                                                                                                               //
    removeElement: function () {                                                                                       // 9
      function removeElement(node) {                                                                                   // 9
        $(node).fadeOut(function () {                                                                                  // 10
          $(this).remove();                                                                                            // 11
        });                                                                                                            //
      }                                                                                                                //
                                                                                                                       //
      return removeElement;                                                                                            //
    }()                                                                                                                //
  };                                                                                                                   //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}},"includes":{"template.access_denied.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/includes/template.access_denied.js                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("accessDenied");                                                                                  // 2
Template["accessDenied"] = new Template("Template.accessDenied", (function() {                                         // 3
  var view = this;                                                                                                     // 4
  return HTML.Raw('<div class="access-denied page jumbotron">\n    <h2>Access Denied</h2>\n    <p>You can\'t get here! Please log in.</p>\n  </div>');
}));                                                                                                                   // 6
                                                                                                                       // 7
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.errors.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/includes/template.errors.js                                                                        //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("errors");                                                                                        // 2
Template["errors"] = new Template("Template.errors", (function() {                                                     // 3
  var view = this;                                                                                                     // 4
  return HTML.DIV({                                                                                                    // 5
    "class": "errors"                                                                                                  // 6
  }, "\n    ", Blaze.Each(function() {                                                                                 // 7
    return Spacebars.call(view.lookup("errors"));                                                                      // 8
  }, function() {                                                                                                      // 9
    return [ "\n      ", Spacebars.include(view.lookupTemplate("error")), "\n    " ];                                  // 10
  }), "\n  ");                                                                                                         // 11
}));                                                                                                                   // 12
                                                                                                                       // 13
Template.__checkName("error");                                                                                         // 14
Template["error"] = new Template("Template.error", (function() {                                                       // 15
  var view = this;                                                                                                     // 16
  return HTML.DIV({                                                                                                    // 17
    "class": "alert alert-danger",                                                                                     // 18
    role: "alert"                                                                                                      // 19
  }, HTML.Raw('\n    <button type="button" class="close" data-dismiss="alert">&times;</button>\n    '), Blaze.View("lookup:message", function() {
    return Spacebars.mustache(view.lookup("message"));                                                                 // 21
  }), "\n  ");                                                                                                         // 22
}));                                                                                                                   // 23
                                                                                                                       // 24
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.header.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/includes/template.header.js                                                                        //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("header");                                                                                        // 2
Template["header"] = new Template("Template.header", (function() {                                                     // 3
  var view = this;                                                                                                     // 4
  return HTML.NAV({                                                                                                    // 5
    "class": "navbar navbar-default",                                                                                  // 6
    role: "navigation"                                                                                                 // 7
  }, "\n    ", HTML.DIV({                                                                                              // 8
    "class": "navbar-header"                                                                                           // 9
  }, "\n      ", HTML.Raw('<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navigation">\n        <span class="sr-only">Toggle navigation</span>\n        <span class="icon-bar"></span>\n        <span class="icon-bar"></span>\n        <span class="icon-bar"></span>\n      </button>'), "\n      ", HTML.A({
    "class": "navbar-brand",                                                                                           // 11
    href: function() {                                                                                                 // 12
      return Spacebars.mustache(view.lookup("pathFor"), "homePage");                                                   // 13
    }                                                                                                                  // 14
  }, "Project 5799"), "\n    "), "\n    ", HTML.DIV({                                                                  // 15
    "class": "collapse navbar-collapse",                                                                               // 16
    id: "navigation"                                                                                                   // 17
  }, "\n      ", HTML.UL({                                                                                             // 18
    "class": "nav navbar-nav"                                                                                          // 19
  }, "\n        ", Spacebars.With(function() {                                                                         // 20
    return Spacebars.call(view.lookup("login_response"));                                                              // 21
  }, function() {                                                                                                      // 22
    return [ "\n            ", Blaze.If(function() {                                                                   // 23
      return Spacebars.call(view.lookup("token"));                                                                     // 24
    }, function() {                                                                                                    // 25
      return [ "\n                ", HTML.LI({                                                                         // 26
        "class": function() {                                                                                          // 27
          return Spacebars.mustache(view.lookup("activeRouteClass"), "browseCameras");                                 // 28
        }                                                                                                              // 29
      }, "\n                    ", HTML.A({                                                                            // 30
        href: function() {                                                                                             // 31
          return Spacebars.mustache(view.lookup("pathFor"), "browseCameras");                                          // 32
        }                                                                                                              // 33
      }, "Cameras"), "\n                "), "\n                ", HTML.LI({                                            // 34
        "class": function() {                                                                                          // 35
          return Spacebars.mustache(view.lookup("activeRouteClass"), "browseVideos");                                  // 36
        }                                                                                                              // 37
      }, "\n                    ", HTML.A({                                                                            // 38
        href: function() {                                                                                             // 39
          return Spacebars.mustache(view.lookup("pathFor"), "browseVideos");                                           // 40
        }                                                                                                              // 41
      }, "Videos"), "\n                "), "\n                ", HTML.LI({                                             // 42
        "class": function() {                                                                                          // 43
          return Spacebars.mustache(view.lookup("activeRouteClass"), "addCamera");                                     // 44
        }                                                                                                              // 45
      }, "\n                    ", HTML.A({                                                                            // 46
        href: function() {                                                                                             // 47
          return Spacebars.mustache(view.lookup("pathFor"), "addCamera");                                              // 48
        }                                                                                                              // 49
      }, "Add Camera"), "\n                "), "\n                ", HTML.LI({                                         // 50
        "class": function() {                                                                                          // 51
          return Spacebars.mustache(view.lookup("activeRouteClass"), "editUser");                                      // 52
        }                                                                                                              // 53
      }, "\n                    ", HTML.A({                                                                            // 54
        href: function() {                                                                                             // 55
          return Spacebars.mustache(view.lookup("pathFor"), "editUser");                                               // 56
        }                                                                                                              // 57
      }, "Profile"), "\n                "), "\n                ", HTML.LI("\n                    ", HTML.A({           // 58
        href: "#",                                                                                                     // 59
        id: "logout_button",                                                                                           // 60
        style: "color:red;"                                                                                            // 61
      }, "Logout"), "\n                "), "\n            " ];                                                         // 62
    }), "\n        " ];                                                                                                // 63
  }), "\n      "), "\n    "), "\n  ");                                                                                 // 64
}));                                                                                                                   // 65
                                                                                                                       // 66
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.loading.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/includes/template.loading.js                                                                       //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("loading");                                                                                       // 2
Template["loading"] = new Template("Template.loading", (function() {                                                   // 3
  var view = this;                                                                                                     // 4
  return Spacebars.include(view.lookupTemplate("spinner"));                                                            // 5
}));                                                                                                                   // 6
                                                                                                                       // 7
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"errors.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/includes/errors.js                                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Template.errors.helpers({                                                                                              // 1
  errors: function () {                                                                                                // 2
    function errors() {                                                                                                // 2
      return Errors.find();                                                                                            // 3
    }                                                                                                                  //
                                                                                                                       //
    return errors;                                                                                                     //
  }()                                                                                                                  //
});                                                                                                                    //
                                                                                                                       //
Template.error.onRendered(function () {                                                                                // 7
  var error = this.data;                                                                                               // 8
  Meteor.setTimeout(function () {                                                                                      // 9
    Errors.remove(error._id);                                                                                          // 10
  }, 3000);                                                                                                            //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"header.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/includes/header.js                                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Template.header.helpers({                                                                                              // 1
                                                                                                                       //
  activeRouteClass: function () {                                                                                      // 3
    function activeRouteClass() /* route names */{                                                                     // 3
      var args = Array.prototype.slice.call(arguments, 0);                                                             // 4
      args.pop();                                                                                                      // 5
                                                                                                                       //
      var active = _.any(args, function (name) {                                                                       // 7
        return Router.current() && Router.current().route.getName() === name;                                          // 8
      });                                                                                                              //
                                                                                                                       //
      return active && 'active';                                                                                       // 11
    }                                                                                                                  //
                                                                                                                       //
    return activeRouteClass;                                                                                           //
  }(),                                                                                                                 //
                                                                                                                       //
  login_response: function () {                                                                                        // 14
    function login_response() {                                                                                        // 14
      return Session.get('login_response');                                                                            // 15
    }                                                                                                                  //
                                                                                                                       //
    return login_response;                                                                                             //
  }()                                                                                                                  //
});                                                                                                                    //
                                                                                                                       //
Template.header.events({                                                                                               // 19
  'click #logout_button': function () {                                                                                // 20
    function clickLogout_button(e) {                                                                                   // 20
      e.preventDefault();                                                                                              // 21
      localStorage.removeItem('login_response');                                                                       // 22
      Session.set('login_response', null);                                                                             // 23
                                                                                                                       //
      // remove all of the client collections on logout                                                                //
      var globalObject = Meteor.isClient ? window : global;                                                            // 20
      for (var property in meteorBabelHelpers.sanitizeForInObject(globalObject)) {                                     // 27
        var object = globalObject[property];                                                                           // 28
        if (object instanceof Meteor.Collection) {                                                                     // 29
          object.remove({});                                                                                           // 30
        }                                                                                                              //
      }                                                                                                                //
                                                                                                                       //
      Router.go('homePage');                                                                                           // 34
    }                                                                                                                  //
                                                                                                                       //
    return clickLogout_button;                                                                                         //
  }()                                                                                                                  //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}},"views":{"template.add_camera.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.add_camera.js                                                                       //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("addCamera");                                                                                     // 2
Template["addCamera"] = new Template("Template.addCamera", (function() {                                               // 3
  var view = this;                                                                                                     // 4
  return HTML.Raw('<!-- Basic Form -->\n    <div class="panel panel-blue margin-bottom-40">\n        <div class="panel-heading" style="background-color:#428bca">\n            <h3 class="panel-title" style="color:white;"><i class="fa fa-tasks"></i>Register Wifi Camera</h3>\n        </div>\n        <div class="panel-body">\n            <form id="form-register-camera" class="margin-bottom-40" role="form">\n                <div class="form-group">\n                    <label for="register-cameraname">Name</label>\n                    <input type="text" class="form-control" id="register-cameraname" placeholder="Enter camera name">\n                </div>\n                <div class="form-group">\n                    <label for="register-cameralocation">Location</label>\n                    <input type="text" class="form-control" id="register-cameralocation" placeholder="Enter camera location">\n                </div>\n                <div class="form-group">\n                    <label for="register-cameraurl">URL</label>\n                    <input type="text" class="form-control" id="register-cameraurl" placeholder="Enter camera url">\n                </div>\n                <div class="form-group">\n                    <label for="register-camerausername">Camera username</label>\n                    <input type="text" class="form-control" id="register-camerausername" placeholder="Enter camera username">\n                </div>\n                <div class="form-group">\n                    <label for="register-camerapassword">Camera password</label>\n                    <input type="text" class="form-control" id="register-camerapassword" placeholder="Enter camera password">\n                </div>\n                <a class="btn btn-success" href="#" id="btn-register-camera">Submit</a>\n                <!--<button type="submit" class="btn btn-lg btn-primary">Submit</button>-->\n            </form>\n        </div>\n    </div>\n    <!-- End Basic Form -->');
}));                                                                                                                   // 6
                                                                                                                       // 7
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.browse_cameras.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.browse_cameras.js                                                                   //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("browseCameras");                                                                                 // 2
Template["browseCameras"] = new Template("Template.browseCameras", (function() {                                       // 3
  var view = this;                                                                                                     // 4
  return HTML.DIV({                                                                                                    // 5
    "class": "page list-group"                                                                                         // 6
  }, HTML.Raw('\n        <span class="list-group-item active">\n            Registered Cameras\n        </span>\n        '), Blaze.Each(function() {
    return Spacebars.call(view.lookup("availableCameras"));                                                            // 8
  }, function() {                                                                                                      // 9
    return [ "\n            ", HTML.A({                                                                                // 10
      href: function() {                                                                                               // 11
        return Spacebars.mustache(view.lookup("pathFor"), "cameraDetail");                                             // 12
      },                                                                                                               // 13
      "class": "list-group-item"                                                                                       // 14
    }, Blaze.View("lookup:name", function() {                                                                          // 15
      return Spacebars.mustache(view.lookup("name"));                                                                  // 16
    })), "\n        " ];                                                                                               // 17
  }), HTML.Raw('\n        <a href="/cameraDetail" class="list-group-item">Work camera</a>\n    '));                    // 18
}));                                                                                                                   // 19
                                                                                                                       // 20
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.browse_videos.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.browse_videos.js                                                                    //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("browseVideos");                                                                                  // 2
Template["browseVideos"] = new Template("Template.browseVideos", (function() {                                         // 3
  var view = this;                                                                                                     // 4
  return HTML.Raw('<div class="page list-group">\n        <span class="list-group-item active">\n            Browse Videos\n        </span>\n        <a href="#" class="list-group-item">Jan 1st 2016</a>\n        <a href="#" class="list-group-item">Feb 23rd 2016</a>\n        <a href="#" class="list-group-item">Feb 28th 2016</a>\n        <a href="#" class="list-group-item">March 2nd 2016</a>\n    </div>');
}));                                                                                                                   // 6
                                                                                                                       // 7
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.camera_detail.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.camera_detail.js                                                                    //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("cameraDetail");                                                                                  // 2
Template["cameraDetail"] = new Template("Template.cameraDetail", (function() {                                         // 3
  var view = this;                                                                                                     // 4
  return HTML.DIV({                                                                                                    // 5
    "class": "page panel panel-default"                                                                                // 6
  }, "\n        ", HTML.DIV({                                                                                          // 7
    "class": "panel-heading clearfix"                                                                                  // 8
  }, "Camera: ", Blaze.View("lookup:name", function() {                                                                // 9
    return Spacebars.mustache(view.lookup("name"));                                                                    // 10
  }), HTML.A({                                                                                                         // 11
    href: function() {                                                                                                 // 12
      return Spacebars.mustache(view.lookup("pathFor"), "editCamera");                                                 // 13
    },                                                                                                                 // 14
    "class": "btn btn-warning",                                                                                        // 15
    style: "float:right;"                                                                                              // 16
  }, "Edit")), HTML.Raw('\n        <div class="panel-body">\n                <div class="col-md-6">\n                    <video width="100%" src="http://v2v.cc/~j/theora_testsuite/320x240.ogg" controls="">\n                        Your browser does not support the <code>video</code> element.\n                    </video>\n                    <div class="span">\n                        <p><button class="btn btn-success btn-block">Record</button></p>\n                    </div>\n                </div>\n                <div class="col-md-6">\n                    <!-- Main component for a primary marketing message or call to action -->\n                    <div class="jumbotron" style="padding-top: 1em;">\n                        <h3>Controls</h3>\n                        <div class="span">\n                            <p><button class="btn btn-primary btn-block">Up</button></p>\n                            <p><button class="btn btn-primary btn-block">Down</button></p>\n                            <p><button class="btn btn-primary btn-block">Left</button></p>\n                            <p><button class="btn btn-primary btn-block">Right</button></p>\n                        </div>\n                    </div>\n                </div>\n        </div>\n    '));
}));                                                                                                                   // 18
                                                                                                                       // 19
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.edit_camera.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.edit_camera.js                                                                      //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("editCamera");                                                                                    // 2
Template["editCamera"] = new Template("Template.editCamera", (function() {                                             // 3
  var view = this;                                                                                                     // 4
  return Spacebars.With(function() {                                                                                   // 5
    return Spacebars.call(view.lookup("cameraToEdit"));                                                                // 6
  }, function() {                                                                                                      // 7
    return [ "\n    ", HTML.Comment(" Basic Form "), "\n    ", HTML.DIV({                                              // 8
      "class": "panel panel-blue margin-bottom-40"                                                                     // 9
    }, "\n        ", HTML.DIV({                                                                                        // 10
      "class": "panel-heading",                                                                                        // 11
      style: "background-color:#428bca"                                                                                // 12
    }, "\n            ", HTML.H3({                                                                                     // 13
      "class": "panel-title",                                                                                          // 14
      style: "color:white;"                                                                                            // 15
    }, HTML.I({                                                                                                        // 16
      "class": "fa fa-tasks"                                                                                           // 17
    }), "Edit Camera ", Blaze.View("lookup:name", function() {                                                         // 18
      return Spacebars.mustache(view.lookup("name"));                                                                  // 19
    }), " Information"), "\n        "), "\n        ", HTML.DIV({                                                       // 20
      "class": "panel-body"                                                                                            // 21
    }, "\n            ", HTML.FORM({                                                                                   // 22
      id: "form-edit-camera",                                                                                          // 23
      "class": "margin-bottom-40",                                                                                     // 24
      role: "form"                                                                                                     // 25
    }, "\n                ", HTML.INPUT({                                                                              // 26
      type: "hidden",                                                                                                  // 27
      id: "edit-cameraname-current",                                                                                   // 28
      value: function() {                                                                                              // 29
        return Spacebars.mustache(view.lookup("name"));                                                                // 30
      }                                                                                                                // 31
    }), "\n                ", HTML.DIV({                                                                               // 32
      "class": "form-group"                                                                                            // 33
    }, "\n                    ", HTML.LABEL({                                                                          // 34
      "for": "edit-cameraname"                                                                                         // 35
    }, "Name"), "\n                    ", HTML.INPUT({                                                                 // 36
      type: "text",                                                                                                    // 37
      "class": "form-control",                                                                                         // 38
      id: "edit-cameraname",                                                                                           // 39
      value: function() {                                                                                              // 40
        return Spacebars.mustache(view.lookup("name"));                                                                // 41
      }                                                                                                                // 42
    }), "\n                "), "\n                ", HTML.DIV({                                                        // 43
      "class": "form-group"                                                                                            // 44
    }, "\n                    ", HTML.LABEL({                                                                          // 45
      "for": "edit-cameralocation"                                                                                     // 46
    }, "Location"), "\n                    ", HTML.INPUT({                                                             // 47
      type: "text",                                                                                                    // 48
      "class": "form-control",                                                                                         // 49
      id: "edit-cameralocation",                                                                                       // 50
      value: function() {                                                                                              // 51
        return Spacebars.mustache(view.lookup("location"));                                                            // 52
      }                                                                                                                // 53
    }), "\n                "), "\n                ", HTML.DIV({                                                        // 54
      "class": "form-group"                                                                                            // 55
    }, "\n                    ", HTML.LABEL({                                                                          // 56
      "for": "edit-cameraurl"                                                                                          // 57
    }, "URL"), "\n                    ", HTML.INPUT({                                                                  // 58
      type: "text",                                                                                                    // 59
      "class": "form-control",                                                                                         // 60
      id: "edit-cameraurl",                                                                                            // 61
      value: function() {                                                                                              // 62
        return Spacebars.mustache(view.lookup("url"));                                                                 // 63
      }                                                                                                                // 64
    }), "\n                "), "\n                ", HTML.DIV({                                                        // 65
      "class": "form-group"                                                                                            // 66
    }, "\n                    ", HTML.LABEL({                                                                          // 67
      "for": "edit-camerausername"                                                                                     // 68
    }, "Camera username"), "\n                    ", HTML.INPUT({                                                      // 69
      type: "text",                                                                                                    // 70
      "class": "form-control",                                                                                         // 71
      id: "edit-camerausername",                                                                                       // 72
      value: function() {                                                                                              // 73
        return Spacebars.mustache(view.lookup("username"));                                                            // 74
      }                                                                                                                // 75
    }), "\n                "), "\n                ", HTML.DIV({                                                        // 76
      "class": "form-group"                                                                                            // 77
    }, "\n                    ", HTML.LABEL({                                                                          // 78
      "for": "edit-camerapassword"                                                                                     // 79
    }, "Camera password"), "\n                    ", HTML.INPUT({                                                      // 80
      type: "password",                                                                                                // 81
      "class": "form-control",                                                                                         // 82
      id: "edit-camerapassword",                                                                                       // 83
      value: function() {                                                                                              // 84
        return Spacebars.mustache(view.lookup("password"));                                                            // 85
      }                                                                                                                // 86
    }), "\n                "), "\n                ", HTML.A({                                                          // 87
      "class": "btn btn-success",                                                                                      // 88
      href: "#",                                                                                                       // 89
      id: "btn-edit-camera"                                                                                            // 90
    }, "Update"), "\n                ", HTML.A({                                                                       // 91
      "class": "btn btn-danger",                                                                                       // 92
      href: "#",                                                                                                       // 93
      id: "btn-edit-camera-delete"                                                                                     // 94
    }, "Delete"), "\n            "), "\n        "), "\n    "), "\n    ", HTML.Comment(" End Basic Form "), "\n    " ];
  });                                                                                                                  // 96
}));                                                                                                                   // 97
                                                                                                                       // 98
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.edit_user.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.edit_user.js                                                                        //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("editUser");                                                                                      // 2
Template["editUser"] = new Template("Template.editUser", (function() {                                                 // 3
  var view = this;                                                                                                     // 4
  return Spacebars.With(function() {                                                                                   // 5
    return Spacebars.call(view.lookup("userData"));                                                                    // 6
  }, function() {                                                                                                      // 7
    return [ "\n        ", HTML.Comment(" Basic Form "), "\n        ", HTML.DIV({                                      // 8
      "class": "panel panel-blue margin-bottom-40"                                                                     // 9
    }, "\n            ", HTML.DIV({                                                                                    // 10
      "class": "panel-heading",                                                                                        // 11
      style: "background-color:#428bca"                                                                                // 12
    }, "\n                ", HTML.H3({                                                                                 // 13
      "class": "panel-title",                                                                                          // 14
      style: "color:white;"                                                                                            // 15
    }, HTML.I({                                                                                                        // 16
      "class": "fa fa-tasks"                                                                                           // 17
    }), "Edit Profile for ", Blaze.View("lookup:username", function() {                                                // 18
      return Spacebars.mustache(view.lookup("username"));                                                              // 19
    })), "\n            "), "\n            ", HTML.DIV({                                                               // 20
      "class": "panel-body"                                                                                            // 21
    }, "\n                ", HTML.FORM({                                                                               // 22
      id: "form-edit-user",                                                                                            // 23
      "class": "margin-bottom-40",                                                                                     // 24
      role: "form"                                                                                                     // 25
    }, "\n                    ", HTML.INPUT({                                                                          // 26
      type: "hidden",                                                                                                  // 27
      id: "edit-username-current",                                                                                     // 28
      value: function() {                                                                                              // 29
        return Spacebars.mustache(view.lookup("username"));                                                            // 30
      }                                                                                                                // 31
    }), "\n                    ", HTML.DIV({                                                                           // 32
      "class": "form-group"                                                                                            // 33
    }, "\n                        ", HTML.LABEL({                                                                      // 34
      "for": "edit-firstName"                                                                                          // 35
    }, "First name"), "\n                        ", HTML.INPUT({                                                       // 36
      type: "text",                                                                                                    // 37
      "class": "form-control",                                                                                         // 38
      id: "edit-firstName",                                                                                            // 39
      value: function() {                                                                                              // 40
        return Spacebars.mustache(view.lookup("firstName"));                                                           // 41
      }                                                                                                                // 42
    }), "\n                    "), "\n                    ", HTML.DIV({                                                // 43
      "class": "form-group"                                                                                            // 44
    }, "\n                        ", HTML.LABEL({                                                                      // 45
      "for": "edit-lastName"                                                                                           // 46
    }, "Last name"), "\n                        ", HTML.INPUT({                                                        // 47
      type: "text",                                                                                                    // 48
      "class": "form-control",                                                                                         // 49
      id: "edit-lastName",                                                                                             // 50
      value: function() {                                                                                              // 51
        return Spacebars.mustache(view.lookup("lastName"));                                                            // 52
      }                                                                                                                // 53
    }), "\n                    "), "\n                    ", HTML.DIV({                                                // 54
      "class": "form-group"                                                                                            // 55
    }, "\n                        ", HTML.LABEL({                                                                      // 56
      "for": "edit-password"                                                                                           // 57
    }, "Password"), "\n                        ", HTML.INPUT({                                                         // 58
      type: "password",                                                                                                // 59
      "class": "form-control",                                                                                         // 60
      id: "edit-password",                                                                                             // 61
      value: function() {                                                                                              // 62
        return Spacebars.mustache(view.lookup("password"));                                                            // 63
      }                                                                                                                // 64
    }), "\n                    "), "\n                    ", HTML.A({                                                  // 65
      "class": "btn btn-success",                                                                                      // 66
      href: "#",                                                                                                       // 67
      id: "btn-edit-user"                                                                                              // 68
    }, "Update"), "\n                    ", HTML.A({                                                                   // 69
      "class": "btn btn-danger",                                                                                       // 70
      href: "#",                                                                                                       // 71
      id: "btn-delete-user"                                                                                            // 72
    }, "Delete Account"), "\n                "), "\n            "), "\n        "), "\n        ", HTML.Comment(" End Basic Form "), "\n    " ];
  });                                                                                                                  // 74
}));                                                                                                                   // 75
                                                                                                                       // 76
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"template.home_page.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/template.home_page.js                                                                        //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       // 1
Template.__checkName("homePage");                                                                                      // 2
Template["homePage"] = new Template("Template.homePage", (function() {                                                 // 3
  var view = this;                                                                                                     // 4
  return [ HTML.Raw("<!-- set the data context -->\n    "), Blaze.If(function() {                                      // 5
    return Spacebars.call(view.lookup("login_response"));                                                              // 6
  }, function() {                                                                                                      // 7
    return [ "\n        ", HTML.DIV({                                                                                  // 8
      "class": "page jumbotron"                                                                                        // 9
    }, "\n            ", HTML.H2("You're logged in"), "\n            ", HTML.P("You can now register new cameras, browse archive videos or view a live feed"), "\n        "), "\n        " ];
  }, function() {                                                                                                      // 11
    return [ "\n            ", HTML.DIV({                                                                              // 12
      "class": "page container"                                                                                        // 13
    }, "\n                ", HTML.DIV({                                                                                // 14
      id: "loginbox",                                                                                                  // 15
      style: "margin-top:50px;",                                                                                       // 16
      "class": "mainbox col-md-6 col-md-offset-3 col-sm-8 col-sm-offset-2"                                             // 17
    }, "\n                    ", HTML.DIV({                                                                            // 18
      "class": "panel panel-info"                                                                                      // 19
    }, "\n                        ", HTML.DIV({                                                                        // 20
      "class": "panel-heading"                                                                                         // 21
    }, "\n                            ", HTML.DIV({                                                                    // 22
      "class": "panel-title"                                                                                           // 23
    }, "Sign In"), "\n                        "), "\n                        ", HTML.DIV({                             // 24
      style: "padding-top:30px",                                                                                       // 25
      "class": "panel-body"                                                                                            // 26
    }, "\n                            ", HTML.DIV({                                                                    // 27
      style: "display:none",                                                                                           // 28
      id: "login-alert",                                                                                               // 29
      "class": "alert alert-danger col-sm-12"                                                                          // 30
    }), "\n                            ", HTML.FORM({                                                                  // 31
      id: "loginform",                                                                                                 // 32
      "class": "form-horizontal",                                                                                      // 33
      role: "form"                                                                                                     // 34
    }, "\n                                ", HTML.DIV({                                                                // 35
      style: "margin-bottom: 25px",                                                                                    // 36
      "class": "input-group"                                                                                           // 37
    }, "\n                                    ", HTML.SPAN({                                                           // 38
      "class": "input-group-addon"                                                                                     // 39
    }, HTML.I({                                                                                                        // 40
      "class": "glyphicon glyphicon-user"                                                                              // 41
    })), "\n                                    ", HTML.INPUT({                                                        // 42
      id: "login-username",                                                                                            // 43
      type: "text",                                                                                                    // 44
      "class": "form-control",                                                                                         // 45
      name: "username",                                                                                                // 46
      value: "",                                                                                                       // 47
      placeholder: "username or email"                                                                                 // 48
    }), "\n                                "), "\n                                ", HTML.DIV({                        // 49
      style: "margin-bottom: 25px",                                                                                    // 50
      "class": "input-group"                                                                                           // 51
    }, "\n                                    ", HTML.SPAN({                                                           // 52
      "class": "input-group-addon"                                                                                     // 53
    }, HTML.I({                                                                                                        // 54
      "class": "glyphicon glyphicon-lock"                                                                              // 55
    })), "\n                                    ", HTML.INPUT({                                                        // 56
      id: "login-password",                                                                                            // 57
      type: "password",                                                                                                // 58
      "class": "form-control",                                                                                         // 59
      name: "password",                                                                                                // 60
      placeholder: "password"                                                                                          // 61
    }), "\n                                "), "\n                                ", HTML.DIV({                        // 62
      style: "margin-top:10px",                                                                                        // 63
      "class": "form-group"                                                                                            // 64
    }, "\n                                    ", HTML.Comment(" Button "), "\n                                    ", HTML.DIV({
      "class": "col-sm-12 controls"                                                                                    // 66
    }, "\n                                        ", HTML.A({                                                          // 67
      id: "btn-login",                                                                                                 // 68
      href: "#",                                                                                                       // 69
      "class": "btn btn-success"                                                                                       // 70
    }, "  Login  "), "\n                                    "), "\n                                "), "\n                                ", HTML.DIV({
      "class": "form-group"                                                                                            // 72
    }, "\n                                    ", HTML.DIV({                                                            // 73
      "class": "col-md-12 control"                                                                                     // 74
    }, "\n                                        ", HTML.DIV({                                                        // 75
      style: "border-top: 1px solid#888; padding-top:15px; font-size:85%"                                              // 76
    }, "\n                                            Don't have an account!\n                                            ", HTML.A({
      href: "#",                                                                                                       // 78
      onclick: "$('#loginbox').hide(); $('#signupbox').show()"                                                         // 79
    }, "\n                                                Sign Up Here\n                                            "), "\n                                        "), "\n                                    "), "\n                                "), "\n                            "), "\n                        "), "\n                    "), "\n                "), "\n                ", HTML.DIV({
      id: "signupbox",                                                                                                 // 81
      style: "display:none; margin-top:50px",                                                                          // 82
      "class": "mainbox col-md-6 col-md-offset-3 col-sm-8 col-sm-offset-2"                                             // 83
    }, "\n                    ", HTML.DIV({                                                                            // 84
      "class": "panel panel-info"                                                                                      // 85
    }, "\n                        ", HTML.DIV({                                                                        // 86
      "class": "panel-heading"                                                                                         // 87
    }, "\n                            ", HTML.DIV({                                                                    // 88
      "class": "panel-title"                                                                                           // 89
    }, "Sign Up"), "\n                            ", HTML.DIV({                                                        // 90
      style: "float:right; font-size: 85%; position: relative; top:-10px"                                              // 91
    }, HTML.A({                                                                                                        // 92
      id: "signinlink",                                                                                                // 93
      href: "#",                                                                                                       // 94
      onclick: "$('#signupbox').hide(); $('#loginbox').show()"                                                         // 95
    }, "Sign In")), "\n                        "), "\n                        ", HTML.DIV({                            // 96
      "class": "panel-body"                                                                                            // 97
    }, "\n                            ", HTML.FORM({                                                                   // 98
      id: "signupform",                                                                                                // 99
      "class": "form-horizontal",                                                                                      // 100
      role: "form"                                                                                                     // 101
    }, "\n                                ", HTML.DIV({                                                                // 102
      id: "signupalert",                                                                                               // 103
      style: "display:none",                                                                                           // 104
      "class": "alert alert-danger"                                                                                    // 105
    }, "\n                                    ", HTML.P("Error:"), "\n                                    ", HTML.SPAN(), "\n                                "), "\n                                ", HTML.DIV({
      "class": "form-group"                                                                                            // 107
    }, "\n                                    ", HTML.LABEL({                                                          // 108
      "for": "register-username",                                                                                      // 109
      "class": "col-md-3 control-label"                                                                                // 110
    }, "Username"), "\n                                    ", HTML.DIV({                                               // 111
      "class": "col-md-9"                                                                                              // 112
    }, "\n                                        ", HTML.INPUT({                                                      // 113
      id: "register-username",                                                                                         // 114
      type: "text",                                                                                                    // 115
      "class": "form-control",                                                                                         // 116
      name: "email",                                                                                                   // 117
      placeholder: "Username"                                                                                          // 118
    }), "\n                                    "), "\n                                "), "\n                                ", HTML.DIV({
      "class": "form-group"                                                                                            // 120
    }, "\n                                    ", HTML.LABEL({                                                          // 121
      "for": "firstname",                                                                                              // 122
      "class": "col-md-3 control-label"                                                                                // 123
    }, "First Name"), "\n                                    ", HTML.DIV({                                             // 124
      "class": "col-md-9"                                                                                              // 125
    }, "\n                                        ", HTML.INPUT({                                                      // 126
      id: "register-firstname",                                                                                        // 127
      type: "text",                                                                                                    // 128
      "class": "form-control",                                                                                         // 129
      name: "firstname",                                                                                               // 130
      placeholder: "First Name"                                                                                        // 131
    }), "\n                                    "), "\n                                "), "\n                                ", HTML.DIV({
      "class": "form-group"                                                                                            // 133
    }, "\n                                    ", HTML.LABEL({                                                          // 134
      "for": "lastname",                                                                                               // 135
      "class": "col-md-3 control-label"                                                                                // 136
    }, "Last Name"), "\n                                    ", HTML.DIV({                                              // 137
      "class": "col-md-9"                                                                                              // 138
    }, "\n                                        ", HTML.INPUT({                                                      // 139
      id: "register-lastname",                                                                                         // 140
      type: "text",                                                                                                    // 141
      "class": "form-control",                                                                                         // 142
      name: "lastname",                                                                                                // 143
      placeholder: "Last Name"                                                                                         // 144
    }), "\n                                    "), "\n                                "), "\n                                ", HTML.DIV({
      "class": "form-group"                                                                                            // 146
    }, "\n                                    ", HTML.LABEL({                                                          // 147
      "for": "password",                                                                                               // 148
      "class": "col-md-3 control-label"                                                                                // 149
    }, "Password"), "\n                                    ", HTML.DIV({                                               // 150
      "class": "col-md-9"                                                                                              // 151
    }, "\n                                        ", HTML.INPUT({                                                      // 152
      id: "register-password",                                                                                         // 153
      type: "password",                                                                                                // 154
      "class": "form-control",                                                                                         // 155
      name: "passwd",                                                                                                  // 156
      placeholder: "Password"                                                                                          // 157
    }), "\n                                    "), "\n                                "), "\n                                ", HTML.DIV({
      "class": "form-group"                                                                                            // 159
    }, "\n                                    ", HTML.Comment(" Button "), "\n                                    ", HTML.DIV({
      "class": "col-md-offset-3 col-md-9"                                                                              // 161
    }, "\n                                        ", HTML.BUTTON({                                                     // 162
      id: "btn-signup",                                                                                                // 163
      type: "button",                                                                                                  // 164
      "class": "btn btn-info"                                                                                          // 165
    }, HTML.I({                                                                                                        // 166
      "class": "icon-hand-right"                                                                                       // 167
    }), "Sign Up"), "\n                                    "), "\n                                "), "\n                            "), "\n                        "), "\n                    "), "\n                "), "\n            "), "\n        " ];
  }) ];                                                                                                                // 169
}));                                                                                                                   // 170
                                                                                                                       // 171
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"add_camera.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/add_camera.js                                                                                //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
var addCameraCalls = {                                                                                                 // 1
                                                                                                                       //
    registerCamera: function () {                                                                                      // 3
        function registerCamera(registerObj) {                                                                         // 3
                                                                                                                       //
            var titleAlert = 'Error';                                                                                  // 5
            var messageAlert = null;                                                                                   // 6
            var typeAlert = 'error';                                                                                   // 7
                                                                                                                       //
            Meteor.call('registerCamera', Utilities.getUserToken(), Utilities.getUsername(), registerObj, function (err, res) {
                if (err) {                                                                                             // 10
                    messageAlert = JSON.stringify(err);                                                                // 11
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 13
                    if (res.statusCode == 200) {                                                                       // 14
                        if (res.hasOwnProperty('content')) {                                                           // 15
                            res = JSON.parse(res.content);                                                             // 16
                            if (res.hasOwnProperty('Message')) {                                                       // 17
                                Utilities.clearForm('form-register-camera');                                           // 18
                                titleAlert = 'Success';                                                                // 19
                                messageAlert = res.Message;                                                            // 20
                                typeAlert = 'success';                                                                 // 21
                                Router.go('browseCameras');                                                            // 22
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
                swal(titleAlert, messageAlert, typeAlert);                                                             // 27
                return res;                                                                                            // 28
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return registerCamera;                                                                                         //
    }()                                                                                                                //
                                                                                                                       //
};                                                                                                                     //
                                                                                                                       //
Template.addCamera.events({                                                                                            // 34
                                                                                                                       //
    'click #btn-register-camera': function () {                                                                        // 36
        function clickBtnRegisterCamera(evt, tpl) {                                                                    // 36
                                                                                                                       //
            var cameraName = tpl.find('input#register-cameraname').value;                                              // 38
            var cameraLocation = tpl.find('input#register-cameralocation').value;                                      // 39
            var cameraURL = tpl.find('input#register-cameraurl').value;                                                // 40
            var cameraUsername = tpl.find('input#register-camerausername').value;                                      // 41
            var cameraPassword = tpl.find('input#register-camerapassword').value;                                      // 42
                                                                                                                       //
            if (cameraName && cameraLocation && cameraURL && cameraUsername && cameraPassword) {                       // 44
                var dataObj = {                                                                                        // 45
                    name: cameraName,                                                                                  // 46
                    location: cameraLocation,                                                                          // 47
                    url: cameraURL,                                                                                    // 48
                    username: cameraUsername,                                                                          // 49
                    password: cameraPassword                                                                           // 50
                };                                                                                                     //
                addCameraCalls.registerCamera(dataObj);                                                                // 52
            } else {                                                                                                   //
                swal('All fields required', 'Please fill all form fields', 'info');                                    // 55
            }                                                                                                          //
        }                                                                                                              //
                                                                                                                       //
        return clickBtnRegisterCamera;                                                                                 //
    }()                                                                                                                //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"browse_cameras.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/browse_cameras.js                                                                            //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Template.browseCameras.helpers({                                                                                       // 1
    availableCameras: function () {                                                                                    // 2
        function availableCameras() {                                                                                  // 2
            return AvailableCameras.find();                                                                            // 3
        }                                                                                                              //
                                                                                                                       //
        return availableCameras;                                                                                       //
    }()                                                                                                                //
});                                                                                                                    //
                                                                                                                       //
Template.browseCameras.onCreated(function () {                                                                         // 7
                                                                                                                       //
    // TODO: If there are no cameras, server response with a 500 error                                                 //
                                                                                                                       //
    var username = Utilities.getUsername();                                                                            // 11
    var token = Utilities.getUserToken();                                                                              // 12
    var titleAlert = 'error';                                                                                          // 13
    var messageAlert = null;                                                                                           // 14
    var typeAlert = 'error';                                                                                           // 15
    var showAlert = true;                                                                                              // 16
                                                                                                                       //
    Meteor.call('getCameras', token, username, function (err, res) {                                                   // 18
        if (err) {                                                                                                     // 19
            messageAlert = JSON.stringify(err);                                                                        // 20
        } else {                                                                                                       //
            messageAlert = JSON.stringify(res);                                                                        // 22
            if (res.statusCode == 200) {                                                                               // 23
                if (res.hasOwnProperty('content')) {                                                                   // 24
                    res = JSON.parse(res.content);                                                                     // 25
                    if (res.hasOwnProperty('CameraData')) {                                                            // 26
                        showAlert = false;                                                                             // 27
                        AvailableCameras.remove({});                                                                   // 28
                        for (var i = 0; i < res.CameraData.length; i++) {                                              // 29
                            AvailableCameras.insert(res.CameraData[i]);                                                // 30
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
            }                                                                                                          //
        }                                                                                                              //
        if (showAlert) {                                                                                               // 36
            swal(titleAlert, messageAlert, typeAlert);                                                                 // 36
        }                                                                                                              //
        return res;                                                                                                    // 37
    });                                                                                                                //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"camera_detail.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/camera_detail.js                                                                             //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
/**                                                                                                                    //
 * Created by arenivar on 4/16/16.                                                                                     //
 */                                                                                                                    //
Template.cameraDetail.helpers({                                                                                        // 4
    cameraId: function () {                                                                                            // 5
        function cameraId() {                                                                                          // 5
            return this.cameraId;                                                                                      // 6
        }                                                                                                              //
                                                                                                                       //
        return cameraId;                                                                                               //
    }(),                                                                                                               //
    cameraName: function () {                                                                                          // 8
        function cameraName() {                                                                                        // 8
            return this.cameraName;                                                                                    // 9
        }                                                                                                              //
                                                                                                                       //
        return cameraName;                                                                                             //
    }()                                                                                                                //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"edit_camera.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/edit_camera.js                                                                               //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
var editCameraCalls = {                                                                                                // 1
                                                                                                                       //
    editCamera: function () {                                                                                          // 3
        function editCamera(editCameraObj, currentCameraName) {                                                        // 3
                                                                                                                       //
            var userName = Utilities.getUsername();                                                                    // 5
            var token = Utilities.getUserToken();                                                                      // 6
                                                                                                                       //
            var titleAlert = 'Error';                                                                                  // 8
            var messageAlert = null;                                                                                   // 9
            var typeAlert = 'error';                                                                                   // 10
                                                                                                                       //
            Meteor.call('editCameraInformation', token, userName, currentCameraName, editCameraObj, function (err, res) {
                if (err) {                                                                                             // 13
                    messageAlert = JSON.stringify(err);                                                                // 14
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 17
                    if (res.statusCode == 200) {                                                                       // 18
                        if (res.hasOwnProperty('content')) {                                                           // 19
                            res = JSON.parse(res.content);                                                             // 20
                            if (res.hasOwnProperty('Status') && res.Status == "Success") {                             // 21
                                if (res.hasOwnProperty('Message')) {                                                   // 22
                                    titleAlert = 'Success';                                                            // 23
                                    messageAlert = res.Message;                                                        // 24
                                    typeAlert = 'success';                                                             // 25
                                    Router.go('browseCameras');                                                        // 26
                                }                                                                                      //
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
                swal(titleAlert, messageAlert, typeAlert);                                                             // 32
                return res;                                                                                            // 33
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return editCamera;                                                                                             //
    }(),                                                                                                               //
                                                                                                                       //
    deleteCamera: function () {                                                                                        // 37
        function deleteCamera(currentCameraName) {                                                                     // 37
            var userName = Utilities.getUsername();                                                                    // 38
            var token = Utilities.getUserToken();                                                                      // 39
                                                                                                                       //
            //alert text                                                                                               //
            var messageAlert = null;                                                                                   // 37
            var titleAlert = 'Error';                                                                                  // 43
            var typeAlert = 'error';                                                                                   // 44
                                                                                                                       //
            Meteor.call('deleteCamera', token, userName, currentCameraName, function (err, res) {                      // 46
                if (err) {                                                                                             // 47
                    messageAlert = JSON.stringify(err);                                                                // 48
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 51
                    if (res.statusCode == 200) {                                                                       // 52
                        typeAlert = 'success';                                                                         // 53
                        titleAlert = 'Deleted';                                                                        // 54
                        if (res.hasOwnProperty('content')) {                                                           // 55
                            res = JSON.parse(res.content);                                                             // 56
                            if (res.hasOwnProperty('Status') && res.Status == "Success") {                             // 57
                                if (res.hasOwnProperty('Message')) {                                                   // 58
                                    typeAlert = 'success';                                                             // 59
                                    titleAlert = 'Success';                                                            // 60
                                    messageAlert = res.Message;                                                        // 61
                                    Router.go('browseCameras');                                                        // 62
                                }                                                                                      //
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
                swal(titleAlert, messageAlert, typeAlert);                                                             // 68
                return res;                                                                                            // 69
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return deleteCamera;                                                                                           //
    }()                                                                                                                //
};                                                                                                                     //
                                                                                                                       //
Template.editCamera.events({                                                                                           // 74
                                                                                                                       //
    'click #btn-edit-camera': function () {                                                                            // 76
        function clickBtnEditCamera(evt, tpl) {                                                                        // 76
            var name = tpl.find('input#edit-cameraname').value;                                                        // 77
            var location = tpl.find('input#edit-cameralocation').value;                                                // 78
            var url = tpl.find('input#edit-cameraurl').value;                                                          // 79
            var cameraUsername = tpl.find('input#edit-camerausername').value;                                          // 80
            var cameraPassword = tpl.find('input#edit-camerapassword').value;                                          // 81
            var currentCameraName = tpl.find('input#edit-cameraname-current').value;                                   // 82
                                                                                                                       //
            if (name && location && url && cameraUsername && cameraPassword && currentCameraName) {                    // 84
                var cameraObj = {                                                                                      // 85
                    name: name,                                                                                        // 86
                    location: location,                                                                                // 87
                    url: url,                                                                                          // 88
                    username: cameraUsername,                                                                          // 89
                    password: cameraPassword                                                                           // 90
                };                                                                                                     //
                editCameraCalls.editCamera(cameraObj, currentCameraName);                                              // 92
            } else {                                                                                                   //
                swal('All fields required', 'Please fill all form fields', 'info');                                    // 95
            }                                                                                                          //
        }                                                                                                              //
                                                                                                                       //
        return clickBtnEditCamera;                                                                                     //
    }(),                                                                                                               //
                                                                                                                       //
    'click #btn-edit-camera-delete': function () {                                                                     // 99
        function clickBtnEditCameraDelete(evt, tpl) {                                                                  // 99
                                                                                                                       //
            var currentCameraName = tpl.find('input#edit-cameraname-current').value;                                   // 101
                                                                                                                       //
            swal({                                                                                                     // 103
                title: "Are you sure?",                                                                                // 104
                text: "You are about to delete the " + currentCameraName + " camera",                                  // 105
                type: "warning",                                                                                       // 106
                showCancelButton: true,                                                                                // 107
                confirmButtonColor: "#DD6B55",                                                                         // 108
                confirmButtonText: "Yes, delete it!",                                                                  // 109
                closeOnConfirm: false,                                                                                 // 110
                html: false                                                                                            // 111
            }, function () {                                                                                           //
                editCameraCalls.deleteCamera(currentCameraName);                                                       // 113
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return clickBtnEditCameraDelete;                                                                               //
    }()                                                                                                                //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"edit_user.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/edit_user.js                                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
var editUserPageCalls = {                                                                                              // 1
                                                                                                                       //
    editUser: function () {                                                                                            // 3
        function editUser(userObj) {                                                                                   // 3
                                                                                                                       //
            var username = Utilities.getUsername();                                                                    // 5
            var token = Utilities.getUserToken();                                                                      // 6
                                                                                                                       //
            var titleAlert = 'Error';                                                                                  // 8
            var messageAlert = null;                                                                                   // 9
            var typeAlert = 'error';                                                                                   // 10
                                                                                                                       //
            Meteor.call('editUserInfo', username, token, userObj, function (err, res) {                                // 12
                if (err) {                                                                                             // 13
                    messageAlert = JSON.stringify(err);                                                                // 14
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 16
                    if (res.statusCode == 200) {                                                                       // 17
                        if (res.hasOwnProperty('content')) {                                                           // 18
                            res = JSON.parse(res.content);                                                             // 19
                            if (res.hasOwnProperty('Status')) {                                                        // 20
                                titleAlert = 'Success';                                                                // 21
                                messageAlert = res.Message;                                                            // 22
                                typeAlert = 'success';                                                                 // 23
                                UserData.remove({});                                                                   // 24
                                userObj.username = username;                                                           // 25
                                UserData.insert(userObj);                                                              // 26
                                Router.go('browseCameras');                                                            // 27
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
                swal(titleAlert, messageAlert, typeAlert);                                                             // 32
                return res;                                                                                            // 33
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return editUser;                                                                                               //
    }(),                                                                                                               //
                                                                                                                       //
    deleteUser: function () {                                                                                          // 37
        function deleteUser() {                                                                                        // 37
                                                                                                                       //
            var username = Utilities.getUsername();                                                                    // 39
            var token = Utilities.getUserToken();                                                                      // 40
                                                                                                                       //
            var titleAlert = 'Error';                                                                                  // 42
            var messageAlert = null;                                                                                   // 43
            var typeAlert = 'error';                                                                                   // 44
                                                                                                                       //
            Meteor.call('deleteUser', username, token, function (err, res) {                                           // 46
                if (err) {                                                                                             // 47
                    messageAlert = JSON.stringify(err);                                                                // 48
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 50
                    if (res.statusCode == 200) {                                                                       // 51
                        if (res.hasOwnProperty('content')) {                                                           // 52
                            res = JSON.parse(res.content);                                                             // 53
                            if (res.hasOwnProperty('Status')) {                                                        // 54
                                titleAlert = 'Success';                                                                // 55
                                messageAlert = res.Message;                                                            // 56
                                typeAlert = 'success';                                                                 // 57
                                localStorage.removeItem('login_response');                                             // 58
                                Session.set('login_response', null);                                                   // 59
                                // remove all of the client collections on logout                                      //
                                var globalObject = Meteor.isClient ? window : global;                                  // 54
                                for (var property in meteorBabelHelpers.sanitizeForInObject(globalObject)) {           // 62
                                    var object = globalObject[property];                                               // 63
                                    if (object instanceof Meteor.Collection) {                                         // 64
                                        object.remove({});                                                             // 65
                                    }                                                                                  //
                                }                                                                                      //
                                Router.go('homePage');                                                                 // 68
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                    swal(titleAlert, messageAlert, typeAlert);                                                         // 72
                }                                                                                                      //
                return res;                                                                                            // 74
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return deleteUser;                                                                                             //
    }()                                                                                                                //
};                                                                                                                     //
                                                                                                                       //
Template.editUser.events({                                                                                             // 79
                                                                                                                       //
    'click #btn-edit-user': function () {                                                                              // 81
        function clickBtnEditUser(evt, tpl) {                                                                          // 81
                                                                                                                       //
            var username = tpl.find('input#edit-username-current').value;                                              // 83
            var firstName = tpl.find('input#edit-firstName').value;                                                    // 84
            var lastName = tpl.find('input#edit-lastName').value;                                                      // 85
            var password = tpl.find('input#edit-password').value;                                                      // 86
                                                                                                                       //
            if (username && firstName && lastName && password) {                                                       // 88
                var userObj = {                                                                                        // 89
                    firstName: firstName,                                                                              // 90
                    lastName: lastName,                                                                                // 91
                    password: password                                                                                 // 92
                };                                                                                                     //
                editUserPageCalls.editUser(userObj);                                                                   // 94
            } else {                                                                                                   //
                swal('All fields are required!', 'Please fill all form fields', 'warning');                            // 96
            }                                                                                                          //
        }                                                                                                              //
                                                                                                                       //
        return clickBtnEditUser;                                                                                       //
    }(),                                                                                                               //
                                                                                                                       //
    'click #btn-delete-user': function () {                                                                            // 100
        function clickBtnDeleteUser(evt, tpl) {                                                                        // 100
            swal({                                                                                                     // 101
                title: "Are you sure?",                                                                                // 102
                text: "You are about to delete your account",                                                          // 103
                type: "warning",                                                                                       // 104
                showCancelButton: true,                                                                                // 105
                confirmButtonColor: "#DD6B55",                                                                         // 106
                confirmButtonText: "Yes, delete it!",                                                                  // 107
                closeOnConfirm: false,                                                                                 // 108
                html: false                                                                                            // 109
            }, function () {                                                                                           //
                editUserPageCalls.deleteUser();                                                                        // 111
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return clickBtnDeleteUser;                                                                                     //
    }()                                                                                                                //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"home_page.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/templates/views/home_page.js                                                                                 //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
var homePageCalls = {                                                                                                  // 1
                                                                                                                       //
    login: function () {                                                                                               // 3
        function login(username, password) {                                                                           // 3
                                                                                                                       //
            var titleAlert = 'Error';                                                                                  // 5
            var messageAlert = null;                                                                                   // 6
            var typeAlert = 'error';                                                                                   // 7
            var showAlert = true;                                                                                      // 8
                                                                                                                       //
            Meteor.call('loginCall', username, password, function (err, res) {                                         // 10
                if (err) {                                                                                             // 11
                    messageAlert = JSON.stringify(err);                                                                // 12
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 14
                    if (res.statusCode == 200) {                                                                       // 15
                        if (res.hasOwnProperty('data')) {                                                              // 16
                            var theData = res.data;                                                                    // 17
                            if (theData.hasOwnProperty('Token')) {                                                     // 18
                                showAlert = false;                                                                     // 19
                                localStorage.setItem('login_response', JSON.stringify({ token: theData.Token, username: username }));
                                Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));     // 21
                                homePageCalls.userData();                                                              // 22
                                Router.go('browseCameras');                                                            // 23
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
                if (showAlert) {                                                                                       // 28
                    swal(titleAlert, messageAlert, typeAlert);                                                         // 28
                }                                                                                                      //
                return res;                                                                                            // 29
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return login;                                                                                                  //
    }(),                                                                                                               //
                                                                                                                       //
    register: function () {                                                                                            // 33
        function register(username, firstName, lastName, password) {                                                   // 33
                                                                                                                       //
            var titleAlert = 'Error';                                                                                  // 35
            var messageAlert = null;                                                                                   // 36
            var typeAlert = 'error';                                                                                   // 37
                                                                                                                       //
            Meteor.call('registerUser', username, firstName, lastName, password, function (err, res) {                 // 39
                if (err) {                                                                                             // 40
                    messageAlert = JSON.stringify(err);                                                                // 41
                } else {                                                                                               //
                    messageAlert = JSON.stringify(res);                                                                // 43
                    if (res.statusCode == 200) {                                                                       // 44
                        if (res.hasOwnProperty('content')) {                                                           // 45
                            res = JSON.parse(res.content);                                                             // 46
                            if (res.hasOwnProperty('Message')) {                                                       // 47
                                titleAlert = 'Success';                                                                // 48
                                messageAlert = res.Message;                                                            // 49
                                typeAlert = 'success';                                                                 // 50
                                homePageCalls.login(username, password);                                               // 51
                            }                                                                                          //
                        }                                                                                              //
                    }                                                                                                  //
                }                                                                                                      //
                swal(titleAlert, messageAlert, typeAlert);                                                             // 56
                return res;                                                                                            // 57
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return register;                                                                                               //
    }(),                                                                                                               //
                                                                                                                       //
    userData: function () {                                                                                            // 61
        function userData() {                                                                                          // 61
            if (Utilities.getUsername() && Utilities.getUserToken()) {                                                 // 62
                Meteor.call('userInfo', Utilities.getUsername(), Utilities.getUserToken(), function (err, res) {       // 63
                    if (err) {                                                                                         // 64
                        swal('Error', JSON.stringify(err), 'warning');                                                 // 65
                    } else {                                                                                           //
                        if (res.statusCode == 200) {                                                                   // 67
                            if (res.hasOwnProperty('content')) {                                                       // 68
                                res = JSON.parse(res.content);                                                         // 69
                                if (res.hasOwnProperty('UserData')) {                                                  // 70
                                    res = res.UserData[0];                                                             // 71
                                    UserData.insert(res);                                                              // 72
                                }                                                                                      //
                            }                                                                                          //
                        } else {                                                                                       //
                            swal('Error', JSON.stringify(res.content), 'warning');                                     // 78
                        }                                                                                              //
                    }                                                                                                  //
                    return res;                                                                                        // 81
                });                                                                                                    //
            } else {                                                                                                   //
                swal('Error', 'Error trying to get user info, please logout and login again', 'warning');              // 84
            }                                                                                                          //
        }                                                                                                              //
                                                                                                                       //
        return userData;                                                                                               //
    }()                                                                                                                //
                                                                                                                       //
};                                                                                                                     //
                                                                                                                       //
Meteor.startup(function () {                                                                                           // 90
    Session.set('login_response', JSON.parse(localStorage.getItem('login_response')));                                 // 91
});                                                                                                                    //
                                                                                                                       //
Template.homePage.helpers({                                                                                            // 94
    login_response: function () {                                                                                      // 95
        function login_response() {                                                                                    // 95
            return Session.get('login_response');                                                                      // 96
        }                                                                                                              //
                                                                                                                       //
        return login_response;                                                                                         //
    }()                                                                                                                //
});                                                                                                                    //
                                                                                                                       //
Template.homePage.events({                                                                                             // 102
                                                                                                                       //
    'click #btn-login': function () {                                                                                  // 104
        function clickBtnLogin(evt, tpl) {                                                                             // 104
                                                                                                                       //
            var username = tpl.find('input#login-username').value;                                                     // 106
            var password = tpl.find('input#login-password').value;                                                     // 107
            homePageCalls.login(username, password);                                                                   // 108
        }                                                                                                              //
                                                                                                                       //
        return clickBtnLogin;                                                                                          //
    }(),                                                                                                               //
                                                                                                                       //
    'click #btn-signup': function () {                                                                                 // 111
        function clickBtnSignup(evt, tpl) {                                                                            // 111
                                                                                                                       //
            var username = tpl.find('input#register-username').value;                                                  // 113
            var firstName = tpl.find('input#register-firstname').value;                                                // 114
            var lastName = tpl.find('input#register-lastname').value;                                                  // 115
            var password = tpl.find('input#register-password').value;                                                  // 116
                                                                                                                       //
            if (username && firstName && lastName && password) {                                                       // 118
                homePageCalls.register(username, firstName, lastName, password);                                       // 119
            } else {                                                                                                   //
                swal('All fields required', 'Please fill all form fields', 'info');                                    // 122
            }                                                                                                          //
        }                                                                                                              //
                                                                                                                       //
        return clickBtnSignup;                                                                                         //
    }()                                                                                                                //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}},"helpers":{"config.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/helpers/config.js                                                                                            //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Accounts.ui.config({                                                                                                   // 1
  passwordSignupFields: 'USERNAME_ONLY'                                                                                // 2
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"errors.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/helpers/errors.js                                                                                            //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// Local (client-only) collection                                                                                      //
Errors = new Mongo.Collection(null);                                                                                   // 2
                                                                                                                       //
throwError = function throwError(message) {                                                                            // 4
  Errors.insert({ message: message });                                                                                 // 5
};                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"utilities.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/helpers/utilities.js                                                                                         //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Utilities = {                                                                                                          // 1
                                                                                                                       //
    getUsername: function () {                                                                                         // 3
        function getUsername() {                                                                                       // 3
            if (Session.get('login_response')) {                                                                       // 4
                if (Session.get('login_response')['username']) {                                                       // 5
                    return Session.get('login_response')['username'];                                                  // 6
                }                                                                                                      //
            }                                                                                                          //
            return null;                                                                                               // 9
        }                                                                                                              //
                                                                                                                       //
        return getUsername;                                                                                            //
    }(),                                                                                                               //
                                                                                                                       //
    getUserToken: function () {                                                                                        // 12
        function getUserToken() {                                                                                      // 12
            if (Session.get('login_response')) {                                                                       // 13
                if (Session.get('login_response')['token']) {                                                          // 14
                    return Session.get('login_response')['token'];                                                     // 15
                }                                                                                                      //
            }                                                                                                          //
            return null;                                                                                               // 18
        }                                                                                                              //
                                                                                                                       //
        return getUserToken;                                                                                           //
    }(),                                                                                                               //
                                                                                                                       //
    clearForm: function () {                                                                                           // 21
        function clearForm(formId) {                                                                                   // 21
            $('#' + formId).find(':input').each(function () {                                                          // 22
                switch (this.type) {                                                                                   // 23
                    case 'password':                                                                                   // 24
                    case 'select-multiple':                                                                            // 25
                    case 'select-one':                                                                                 // 26
                    case 'text':                                                                                       // 27
                    case 'textarea':                                                                                   // 28
                        $(this).val('');                                                                               // 29
                        break;                                                                                         // 30
                    case 'checkbox':                                                                                   // 23
                    case 'radio':                                                                                      // 32
                        this.checked = false;                                                                          // 33
                }                                                                                                      // 23
            });                                                                                                        //
        }                                                                                                              //
                                                                                                                       //
        return clearForm;                                                                                              //
    }()                                                                                                                //
};                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}},"data_collections.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/data_collections.js                                                                                          //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
// Local collections only                                                                                              //
AvailableCameras = new Mongo.Collection(null);                                                                         // 2
UserData = new Mongo.Collection(null);                                                                                 // 3
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

},"main.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// client/main.js                                                                                                      //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
                                                                                                                       //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}},"lib":{"router.js":function(){

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                                                                                     //
// lib/router.js                                                                                                       //
//                                                                                                                     //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
                                                                                                                       //
Router.configure({                                                                                                     // 1
  layoutTemplate: 'layout',                                                                                            // 2
  loadingTemplate: 'loading',                                                                                          // 3
  notFoundTemplate: 'notFound'                                                                                         // 4
});                                                                                                                    //
                                                                                                                       //
Router.route('/', {                                                                                                    // 7
  name: 'homePage'                                                                                                     // 8
});                                                                                                                    //
                                                                                                                       //
Router.route('/addCamera/', {                                                                                          // 11
  name: 'addCamera'                                                                                                    // 12
});                                                                                                                    //
                                                                                                                       //
Router.route('/browseVideos/', {                                                                                       // 15
  name: 'browseVideos'                                                                                                 // 16
});                                                                                                                    //
                                                                                                                       //
Router.route('/browseCameras/', {                                                                                      // 19
  name: 'browseCameras'                                                                                                // 20
});                                                                                                                    //
                                                                                                                       //
Router.route('/editCamera/:name', {                                                                                    // 23
  name: 'editCamera',                                                                                                  // 24
  data: function () {                                                                                                  // 25
    function data() {                                                                                                  // 25
      return {                                                                                                         // 26
        cameraToEdit: AvailableCameras.findOne({ name: this.params.name })                                             // 27
      };                                                                                                               //
    }                                                                                                                  //
                                                                                                                       //
    return data;                                                                                                       //
  }()                                                                                                                  //
});                                                                                                                    //
                                                                                                                       //
Router.route('/cameraDetail/:name', {                                                                                  // 32
  name: 'cameraDetail',                                                                                                // 33
  data: function () {                                                                                                  // 34
    function data() {                                                                                                  // 34
      return {                                                                                                         // 35
        name: this.params.name                                                                                         // 36
      };                                                                                                               //
    }                                                                                                                  //
                                                                                                                       //
    return data;                                                                                                       //
  }()                                                                                                                  //
});                                                                                                                    //
                                                                                                                       //
Router.route('/editUser/:username', {                                                                                  // 41
  name: 'editUser',                                                                                                    // 42
  data: function () {                                                                                                  // 43
    function data() {                                                                                                  // 43
      return {                                                                                                         // 44
        userData: UserData.findOne()                                                                                   // 45
      };                                                                                                               //
    }                                                                                                                  //
                                                                                                                       //
    return data;                                                                                                       //
  }()                                                                                                                  //
});                                                                                                                    //
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

}}},{"extensions":[".js",".json",".html",".css"]});
require("./client/templates/application/template.layout.js");
require("./client/templates/application/template.not_found.js");
require("./client/templates/includes/template.access_denied.js");
require("./client/templates/includes/template.errors.js");
require("./client/templates/includes/template.header.js");
require("./client/templates/includes/template.loading.js");
require("./client/templates/views/template.add_camera.js");
require("./client/templates/views/template.browse_cameras.js");
require("./client/templates/views/template.browse_videos.js");
require("./client/templates/views/template.camera_detail.js");
require("./client/templates/views/template.edit_camera.js");
require("./client/templates/views/template.edit_user.js");
require("./client/templates/views/template.home_page.js");
require("./lib/router.js");
require("./client/templates/application/layout.js");
require("./client/templates/includes/errors.js");
require("./client/templates/includes/header.js");
require("./client/templates/views/add_camera.js");
require("./client/templates/views/browse_cameras.js");
require("./client/templates/views/camera_detail.js");
require("./client/templates/views/edit_camera.js");
require("./client/templates/views/edit_user.js");
require("./client/templates/views/home_page.js");
require("./client/helpers/config.js");
require("./client/helpers/errors.js");
require("./client/helpers/utilities.js");
require("./client/data_collections.js");
require("./client/main.js");