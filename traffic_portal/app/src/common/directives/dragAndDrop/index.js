//application directives
module.exports = angular.module('trafficPortal.directives.dragAndDrop',[])
    .directive('dndEnable', require('./dragdropDirective'))
    .directive('droppable', require('./droppableDirective'))
;