Template.layout.onRendered(function() {
  this.find('#main')._uihooks = {
    insertElement: function(node, next) {
      console.log("insert element");
      $(node)
          .hide()
          .insertBefore(next)
          .fadeIn();
      //setTimeout(function(){
      //  $(node)
      //      .hide()
      //      .insertBefore(next)
      //      .fadeIn();
      //}, 500);

    },
    removeElement: function(node) {
      console.log("remove element");
      $(node).fadeOut(function() {
        $(this).remove();
      });
    }
  }
});