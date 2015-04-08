/************************************************************************
*************************************************************************
@Name    :      jMenu - jQuery Plugin
@Revison :      2.0
@Date    :      08/2013
@Author  :      ALPIXEL - (www.myjqueryplugins.com - www.alpixel.fr)
@Support :      FF, IE7, IE8, MAC Firefox, MAC Safari
@License :      Open Source - MIT License : http://www.opensource.org/licenses/mit-license.php

**************************************************************************
*************************************************************************/

/** jMenu Plugin **/
(function($) {
    $.jMenu = {
        /**************/
        /** OPTIONS **/
        /**************/
        defaults: {
            ulWidth:           'auto',
            absoluteTop:       33,
            absoluteLeft:      0,
            TimeBeforeOpening: 100,
            TimeBeforeClosing: 100,
            animatedText:      true,
            paddingLeft:       7,
            openClick:         false,
            effects: {
                effectSpeedOpen:  150,
                effectSpeedClose: 150,
                effectTypeOpen:   'slide',
                effectTypeClose:  'slide',
                effectOpen:       'swing',
                effectClose:      'swing'
            }
        },

        /*****************/
        /** Init Method **/
        /*****************/
        init: function(options) {
            /* vars **/
            opts = $.extend({}, $.jMenu.defaults, options);

            // Set global width of the sub-menus links
            if(opts.ulWidth == 'auto')
                $width = $('.fNiv').outerWidth(false);
            else
                $width = opts.ulWidth;



            $(".jMenu li").each(function() {
                var
                    $thisChild = $(this).find('a:first'),
                    $allUl = $(this).find('ul');

                if($.jMenu._IsParent($thisChild))
                {
                    $thisChild.addClass('isParent');

                    var
                        $ul = $thisChild.next(),
                        $position = $thisChild.position();

                    if($(this).hasClass('jmenu-level-0'))
                        $ul.css({
                            top:   $position.top + opts.absoluteTop,
                            left:  $position.left + opts.absoluteLeft,
                            width : $width
                        });
                    else
                        $ul.css({
                            top:   $position.top,
                            left:  $position.left + $width,
                            width : $width
                        });


                    if(!opts.openClick)
                        $(this).bind({
                            mouseenter:function() {
                                if($(this).hasClass('jmenu-level-0')) {
									$position = $(this).position();
									$ul.css({
										left: $position.left + opts.absoluteLeft,
										top: $position.top + opts.absoluteTop
									});
								}
								$.jMenu._show($ul);
                            },
                            mouseleave:function(){
                                $.jMenu._closeList($ul);
                            }
                        });
                    else
                        $(this).bind({
                            click:function(e) {
                                e.preventDefault();
                                $.jMenu._show($ul);
                            },
                            mouseleave:function(){
                                $.jMenu._closeList($ul);
                            }
                        });
                }
            });
        },


        /****************************
        *****************************
        ** jMenu Methods Below     **
        *****************************
        ****************************/
        _show: function(el) {
            switch(opts.effects.effectTypeOpen) {
                case 'slide':
                    el.stop(true, true).delay(opts.TimeBeforeOpening).slideDown(opts.effects.effectSpeedOpen, opts.effects.effectOpen);
                    break;
                case 'fade':
                    el.stop(true, true).delay(opts.TimeBeforeOpening).fadeIn(opts.effects.effectSpeedOpen, opts.effects.effectOpen);
                    break;
                default:
                    el.stop(true, true).delay(opts.TimeBeforeOpening).show();
            }
        },

        _closeList: function(el) {
            switch(opts.effects.effectTypeClose) {
                case 'slide':
                    el.stop(true,true).delay(opts.TimeBeforeClosing).slideUp(opts.effects.effectSpeedClose, opts.effects.effectClose);
                    break;
                case 'fade':
                    el.stop(true,true).delay(opts.TimeBeforeClosing).fadeOut(opts.effects.effectSpeedClose, opts.effects.effectClose);
                    break;
                default:
                    el.delay(opts.TimeBeforeClosing).hide();
            }
        },

        _animateText: function(el) {
            var paddingInit = parseInt(el.css('padding-left'));
            el.hover(
                function() {
                    $(this).stop(true,false).animate({paddingLeft: paddingInit + opts.paddingLeft}, 100);
                },
                function() {
                    $(this).stop(true,false).animate({paddingLeft:paddingInit}, 100);
                }
            );
        },

        _IsParent: function(el) {
            if (el.next().is('ul'))
                return true;
            else
                return false;
        },

        _isReadable: function() {
            if ($(".jmenu-level-0").length > 0)
                return true;
            else
                return false;
        },

        _error: function() {
            alert('jMenu plugin can\'t be initialized. Please, check you have the \'.jmenu-level-0\' class on your first level <li> elements.');
        }
    };

    jQuery.fn.jMenu = function(options){
        $(this).addClass('jMenu');
        $(this)
            .children('li').addClass('jmenu-level-0')
            .children('a').addClass('fNiv');

        if($.jMenu._isReadable()) {
            $.jMenu.init(options);
        } else {
            $.jMenu._error();
        }
    };
})(jQuery);
