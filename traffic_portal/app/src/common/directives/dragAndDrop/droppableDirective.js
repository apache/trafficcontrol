/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

var DroppableDirective = function () {
    return {
        scope: {
            drop: '&',
            dragsmart: '&',
            bin: '='
        },
        link: function(scope, element) {
            var el = element[0];

            el.addEventListener(
                'dragover',
                function(e) {
                    e.dataTransfer.dropEffect = 'move';
                    // allows us to drop
                    if (e.preventDefault) e.preventDefault();
                    this.classList.add('over');
                    var rect = e.target.getBoundingClientRect();
                    var y = e.offsetY;
                    if (y <= (rect.height / 2)) {
                        if (e.target.classList.contains('drop-child')) {
                            e.target.parentElement.style.borderTop = '2px solid red';
                            e.target.parentElement.style.borderBottom = '';
                        } else {
                            e.target.style.borderTop = '2px solid red';
                            e.target.style.borderBottom = '';
                        }
                        scope.$parent.$parent.moveAbove = true;
                    } else {
                        if (e.target.classList.contains('drop-child')) {
                            e.target.parentElement.style.borderBottom = '2px solid red';
                            e.target.parentElement.style.borderTop = '';
                        } else {
                            e.target.style.borderBottom = '2px solid red';
                            e.target.style.borderTop = '';
                        }
                        scope.$parent.$parent.moveAbove = false;
                    }
                    return false;
                },
                false
            );

            el.addEventListener(
                'dragenter',
                function(e) {
                    this.classList.add('over');
                    if (e.preventDefault) e.preventDefault();
                    return false;
                },
                false
            );

            el.addEventListener(
                'dragleave',
                function(e) {
                    this.classList.remove('over');
                    if (e.target.classList.contains('drop-child')) {
                        e.target.parentElement.style.borderTop = '';
                        e.target.parentElement.style.borderBottom = '';
                    } else {
                        e.target.style.borderTop = '';
                        e.target.style.borderBottom = '';
                    }
                    return false;
                },
                false
            );

            el.addEventListener(
                'drop',
                function(e) {
                    if (e.stopPropagation) e.stopPropagation();
                    if (e.preventDefault) e.preventDefault(); // *needed* for firefox
                    this.classList.remove('over');
                    if (e.target.classList.contains('drop-child')) {
                        e.target.parentElement.style.borderTop = '';
                        e.target.parentElement.style.borderBottom = '';
                    } else {
                        e.target.style.borderTop = '';
                        e.target.style.borderBottom = '';
                    }
                    var from = e.dataTransfer.getData("Text");
                    if (from != this.textContent) { // ignore drop on self
                        scope.$apply(function(scope) {
                            var fn = scope.drop();
                            if ('undefined' !== typeof fn) {
                                fn();
                            }
                        });
                    }
                    return false;
                },
                false
            );

            el.draggable = true;

            el.addEventListener(
                'dragstart',
                function (e) {
                    if (e.stopPropagation) e.stopPropagation();
                    e.dataTransfer.effectAllowed = 'move';
                    e.dataTransfer.setData('Text',this.textContent);  // some value *required* for firefox; "Text" - not "text/plain" for IE
                    this.classList.add('drag');
                    scope.$apply(function (scope) {
                        var fn = scope.dragsmart();
                        if ('undefined' !== typeof fn) {
                            fn();
                        }
                    });
                    return false;
                },
                false
            );

            el.addEventListener(
                'dragend',
                function (e) {
                    this.classList.remove('drag');
                    return false;
                },
                false
            );
        }
    };
};

module.exports = DroppableDirective;