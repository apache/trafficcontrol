/*

     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.
 */

$('.fancybox').fancybox({
    'height': '100%',
//    'width': '100%',
'autoDimensions': false,
'autoScale' : false,
'type' : 'iframe',
'showclosebutton' : true,
'shownavarrows' : false,
'padding': 0,
'beforeClose': function() { 
    if (parent.activeTable != undefined) {
        parent.activeTable.fnReloadAjax();
    } else {
       location.reload(); 
    }
   }
});

if ($.fn.dataTableExt != undefined) {
    var calcDataTableHeight = function() {
        var $window = $(window);
        return Math.round($window.height() ) -196 ; // this gets rid off the scrollbar in Chrome _and FF, looks a bit too 
        // short in FF I think.... Need to test in Safari
    };

    jQuery.extend( jQuery.fn.dataTableExt.oSort, {
        "alt-string-pre": function ( a ) {
            return a.match(/alt="(.*?)"/)[1].toLowerCase();
        },

        "alt-string-asc": function( a, b ) {
            return ((a < b) ? -1 : ((a > b) ? 1 : 0));
        },
    
        "alt-string-desc": function(a,b) {
            return ((a < b) ? 1 : ((a > b) ? -1 : 0));
        },
       "alt-string-pre": function ( a ) {
            return a.match(/alt="(.*?)"/)[1].toLowerCase();
        },
         
        "alt-string-asc": function( a, b ) {
            return ((a < b) ? -1 : ((a > b) ? 1 : 0));
        },
     
        "alt-string-desc": function(a,b) {
            return ((a < b) ? 1 : ((a > b) ? -1 : 0));
        },

        "alt-number-pre": function ( a ) {
            var x = a.match(/alt="(.*?)"/)[1];
            return parseFloat(x);
        },
         
        "alt-number-asc": function( a, b ) {
            return ((a < b) ? -1 : ((a > b) ? 1 : 0));
        },
     
        "alt-number-desc": function(a,b) {
            return ((a < b) ? 1 : ((a > b) ? -1 : 0));
        }, 

        "day-hour-pre": function(a) {
            var days = 0;
            var hrs = 0;
            if (a.match(/(.*?)d/)) {
             days = a.match(/(.*?)d/)[1];
             hrs = a.match(/d(.*?)h/)[1]
            } else {
                hrs = a.match(/(.*?)h/)[1]
            }
            return(parseInt(days) * 24 + parseInt(hrs));
        },

        "day-hour-asc": function(a,b) {
            return ((a < b) ? -1 : ((a > b) ? 1 : 0));
        },

        "day-hour-desc": function(a,b) {
            return ((a < b) ? 1 : ((a > b) ? -1 : 0));
        }
    } );

    // Just reload on resize... Will that suffice?
    $(window).resize(function () { 
        if ($("#gbps_flot").length > 0)
            return; // don't redraw the graph status page
        var uri_string = window.location.pathname;
        if (uri_string.indexOf('custom_charts') == -1) // do not reload on the customcharts page
        location.reload();
    });

    // For the table refreshing, see http://datatables.net/plug-ins/api
    jQuery.fn.dataTableExt.oApi.fnReloadAjax = function ( oSettings, sNewSource, fnCallback, bStandingRedraw )
    {
        // DataTables 1.10 compatibility - if 1.10 then `versionCheck` exists.
        // 1.10's API has ajax reloading built in, so we use those abilities
        // directly.
        if ( jQuery.fn.dataTable.versionCheck ) {
            var api = new jQuery.fn.dataTable.Api( oSettings );
    
            if ( sNewSource ) {
                api.ajax.url( sNewSource ).load( fnCallback, !bStandingRedraw );
            }
            else {
                api.ajax.reload( fnCallback, !bStandingRedraw );
            }
            return;
        }
    
        if ( sNewSource !== undefined && sNewSource !== null ) {
            oSettings.sAjaxSource = sNewSource;
        }
    
        // Server-side processing should just call fnDraw
        if ( oSettings.oFeatures.bServerSide ) {
            this.fnDraw();
            return;
        }
    
        this.oApi._fnProcessingDisplay( oSettings, true );
        var that = this;
        var iStart = oSettings._iDisplayStart;
        var aData = [];
    
        this.oApi._fnServerParams( oSettings, aData );
    
        oSettings.fnServerData.call( oSettings.oInstance, oSettings.sAjaxSource, aData, function(json) {
            /* Clear the old information from the table */
            that.oApi._fnClearTable( oSettings );
    
            /* Got the data - add it to the table */
            var aData =  (oSettings.sAjaxDataProp !== "") ?
                that.oApi._fnGetObjectDataFn( oSettings.sAjaxDataProp )( json ) : json;
    
            for ( var i=0 ; i<aData.length ; i++ )
            {
                that.oApi._fnAddData( oSettings, aData[i] );
            }
    
            oSettings.aiDisplay = oSettings.aiDisplayMaster.slice();
    
            that.fnDraw();
    
            if ( bStandingRedraw === true )
            {
                oSettings._iDisplayStart = iStart;
                that.oApi._fnCalculateEnd( oSettings );
                that.fnDraw( false );
            }
    
            that.oApi._fnProcessingDisplay( oSettings, false );
    
            /* Callback user function - for event handlers etc */
            if ( typeof fnCallback == 'function' && fnCallback !== null )
            {
                fnCallback( oSettings );
            }
        }, oSettings );
    };
}
