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

// From http://www.webappers.com/progressBar/lib/progress.js
// use:
// <script>display('<%= $short %>_percent_indicator',0,1);</script><
// update:
// setProgress(cdn + "_percent_indicator", Math.round(lastVal/maxGbps[cdn] * 100));

/* WebAppers Progress Bar, version 0.2
* (c) 2007 Ray Cheung
*
* WebAppers Progress Bar is freely distributable under the terms of an Creative Commons license.
* For details, see the WebAppers web site: http://wwww.Webappers.com/
*
/*--------------------------------------------------------------------------*/
var initial = -120;
var imageWidth=240;
var eachPercent = (imageWidth/2)/100;

function percentLabel( id, percentage )
{ 
    document.write('<div id="'+id+'Text" class="percent_label">'+percentage+'%</div>'); 
}

function progressBar( id, percentage )
{ 
	document.write('<div id="'+id+'" class="progress_bar"></div>');
}


function setText (id, percent)
{
    $('#' + id +'Text')[0].innerHTML = percent+"%";
}

function setProgress(id, percentage)
{
    if (isNaN(percentage) || !isFinite(percentage)){
	   percentage = 0;
	}
    var percentageWidth = eachPercent * percentage;
    var pb = $('#' + id);
	/* Sample Colors
	var green = '#A4E142';
	var yellow = '#D5E043';
	var red = '#E05043';
	*/

	/* Mark T's colors*/
	var green = '#76AA5E';
	var yellow = '#C07818';
	var red = '#6D1E05';
	var color = green;
	var backgroundColor = '#ccc';
	if (percentage = 0) {
	  color = backgroundColor;
	} else if ((percentage > 50) && (percentage < 75)) {
	  color = yellow;
	} else if ((percentage > 75) && (percentage <=100))  {
	  color = red;
	} else {
	  color = green;
	}

    pb.progressbar({
	     value: percentageWidth,
    });
    pb.css({
	  "background": backgroundColor
    });
	var pbValue = pb.find( ".ui-progressbar-value" );
	pbValue.css({
	  "background": color
	});
	//#A4E142 - Green
	//#D5E043 - Yellow
	//#E05043 - Red
    var newProgress = eval(initial)+eval(percentageWidth)+'px';
    $('#' + id)[0].style.backgroundPosition=newProgress+' 0';
    setText(id,percentage);
}

/*--------------------------------------------------------------------------*/
