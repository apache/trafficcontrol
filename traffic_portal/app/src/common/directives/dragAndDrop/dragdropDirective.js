var DndEnable = function($compile) {
    return {
        restrict: 'A',  // use as an Attribute only
        replace: false,
        terminal: true, //this setting is crucial since I'm adding attributes that need to be compiled
        priority: 1050, // "me first", ng-repeat priority is 1000
        link: function(scope, element, attrs) {
            if (attrs.dndEnable) { // if dnd-enable evaluates to true
                if (attrs.pageid === 'cacheGroupFallback') {
                    element.attr('dragsmart', 'handleDrag(fb)');
                    element.attr('droppable', 'true');
                    element.attr('drop', 'handleDrop(fb)');
                }
            }
            element.removeAttr('dnd-enable'); // prevent infinite loop on compile
            $compile(element)(scope);
        }
    };
};

module.exports = DndEnable;