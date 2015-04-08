	function updateAjaxComponents(arg) {
		for ( var key in arg) {
			if (arg.hasOwnProperty(key)) {
				var $e = $("#"+key);
				var o = arg[key];
				for (var key2 in o) {
					if(key2 === "v") {
						$e.text(o["v"]);
						var graphId = $e.attr("data-graph-id");
						if(graphId != null) {
							var index = $e.attr("data-graph-index");
						}
					} else {
						$e.attr(key2, o[key2]);
					}
				}
			}
		}
	}

