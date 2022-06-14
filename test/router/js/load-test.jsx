{/*
	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/}

/**
 * Creates an array of arrays of elements that return the same value from the
 * given function.
 *
 * @template T
 * @param {Array<T>} a The array to group.
 * @param {(e: T) => string} f The grouping function; values in `a` that return
 * the same value from this function will be grouped together.
 * @returns {Array<Array<T>>} All of the elements of `a`, grouped by `f`.
 */
function groupBy(a, f) {
    /** @type Record<string, Array<T>> */
    const groups = {};
    for (const element of a) {
        const group = f(element);
        if (group in groups) {
            groups[group].push(element);
        } else {
            groups[group] = [element];
        }
    }
    return Object.keys(groups).map(g=>groups[g]);
}

/**
 * Gets the unique elements of an Array as a new Array. Does not modify Arrays
 * in-place.
 *
 * @template T
 * @param {Array<T>} a The array to be unique'd.
 * @returns {Array<T>} A new array with all of the values in `a`, but without
 * any duplicate values. Duplicates are determined with `===`.
 */
function unique(a) {
    return Array.from(new Set(a));
}

var CdnTestForm = React.createClass({
    getInitialState: function() {
        return {
            cdn: "none",
            numHttp: 10,
            numHttps: 10,
            txCount: 1000,
            connections: 10,
            caFile: './ca.crt'
        }
    },
    handleCdnChange: function(e) {
        this.setState({cdn: e.target.value})
    },
    handleNumHttpChange: function(e) {
      this.setState({numHttp: e.target.value})
    },
    handleNumHttpsChange: function(e) {
      this.setState({numHttps: e.target.value})
    },
    handleTxCountChange: function(e) {
        this.setState({txCount: e.target.value})
    },
    handleConnectionsChange: function(e) {
        this.setState({connections: e.target.value})
    },
    handleCaFileChange: function(e) {
        this.setState({caFile: e.target.value})
    },
    handleSubmit: function(e) {
        e.preventDefault();
        var cdn = this.state.cdn.trim();
        var caFile = this.state.caFile.trim();
        var txCount = this.state.txCount;
        var connections = this.state.connections;
        var numHttp = this.state.numHttp;
        var numHttps = this.state.numHttps;

        var testCdn = this.props.cdns.find(function (item) { return item.name == cdn});

        var blah = testCdn.deliveryServices[0].exampleURLs[0];

        blah = blah.substring(blah.indexOf(testCdn.name) + testCdn.name.length);

        testCdn.name = testCdn.name + blah;

        var httpDsList = testCdn.deliveryServices.filter(function (item) {
            return item.type.toLowerCase().includes("dns") == false && item.protocol == 0;
        });


        var randomIndices = [];

        for (var i = 0; i < numHttp; i++) {
            randomIndices.push(Math.floor((Math.random() * httpDsList.length)))
        }

        var httpDeliveryServices = [];

        for (var i = 0; i < randomIndices.length; i++) {
            var u = httpDsList[randomIndices[i]].exampleURLs[0];
            if (u.indexOf("edge") != -1) {
                console.log("????! " + u);
            }
            var id = u.substring(u.indexOf("ccr") + 4, u.indexOf(testCdn.name) - 1);

            httpDeliveryServices.push(id);
        }

        var httpsDsList = testCdn.deliveryServices.filter(function (item) {
            return item.type.toLowerCase().includes("dns") == false && item.protocol > 0;
        });

        randomIndices = [];

        for (var i = 0; i < numHttps; i++) {
            randomIndices.push(Math.floor((Math.random() * httpsDsList.length)))
        }

        var httpsDeliveryServices = [];
        for (var i = 0; i < randomIndices.length; i++) {
            var u = httpsDsList[randomIndices[i]].exampleURLs[0];
            var id = u.substring(u.indexOf("ccr") + 4, u.indexOf(testCdn.name) - 1);
            httpsDeliveryServices.push(id);
        }

        this.props.onLoadTestSubmit({
            cdn: testCdn.name,
            httpDeliveryServices: httpDeliveryServices,
            httpsDeliveryServices: httpsDeliveryServices,
            txCount: txCount,
            connections: connections,
            caFile: caFile
        });
    },

    render: function() {
        var cdnOptions = this.props.cdns.map(function (cdn) {
            return (
                <option key={cdn.name} value={cdn.name}>{cdn.name}</option>
            )
        });

        return (
            <div className="cdnList">
                <h3>3. Run Load Test</h3>
                <form onSubmit={this.handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="testCdn">CDN</label>
                        <select id="testCdnInput" className="form-control" value={this.state.cdn} onChange={this.handleCdnChange}>
                            {cdnOptions}
                        </select>
                    </div>
                    <div className="form-group">
                        <label htmlFor="numHttpInput"># of Http DS to exercise</label>
                        <input id="numHttpInput" type="text" className="form-control" value={this.state.numHttp} onChange={this.handleNumHttpChange} placeholder="0"/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="numHttpsInput"># of HTTPS DS to exercise</label>
                        <input id="numHttpsInput" type="text" className="form-control" value={this.state.numHttps} onChange={this.handleNumHttpsChange}/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="txCountInput"># of transactions per Delivery Services</label>
                        <input id="txCountInput" type="text" className="form-control" value={this.state.txCount} onChange={this.handleTxCountChange}/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="connectionsInput"># of concurrent requests per Delivery Services</label>
                        <input id="connectionsInput" type="text" className="form-control" value={this.state.connections} onChange={this.handleConnectionsChange} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="caFileInput">CA file path on server</label>
                        <input id="caFileInput" className="form-control" type="text" value={this.state.caFile} onChange={this.handleCaFileChange} placeholder="./ca.crt"/>
                    </div>
                    <button type="submit" className="btn btn-info">Run Test</button>
                </form>
            </div>
        )
    }
});

var CdnList = React.createClass({
    render: function() {

        var cdnRows = this.props.cdns.map(function (cdn) {
            var name = cdn.name;
            var httpDsList = cdn.deliveryServices.filter(function (ds) {
                return ds.type.toLowerCase().includes("dns") == false;
            });

            var httpCount = httpDsList.filter(function (ds) { return ds.protocol == 0; }).length;
            var httpsCount = httpDsList.filter(function (ds) { return ds.protocol > 0; }).length;

            return (
                <div key={cdn.name} className="row strong">
                    <div className="col-sm-4 ">{name}</div>
                    <div className="col-sm-4">{httpCount}</div>
                    <div className="col-sm-4">{httpsCount}</div>
                </div>
            )
        });

        return (
            <div className="well">
            <h4>Current Lab Data</h4>
            <div className="container">
                <div className="row">
                    <div className="col-sm-4">Name</div>
                    <div className="col-sm-4"># of Http Delivery Services</div>
                    <div className="col-sm-4"># of HTTPS delivery services</div>
                </div>
                {cdnRows}
            </div>
            </div>
        )
    }
});

var LabBox = React.createClass({
    getInitialState: function() {
        return {deliveryServices: [], opsHost: ""}
    },
    handleOpsHostChange: function(e) {
        this.setState({opsHost: e.target.value})
    },
    loadResultsFromServer: function() {
        $.ajax({
            url: this.props.url,
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({data: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    componentDidMount: function() {
        this.loadResultsFromServer();
        setInterval(this.loadResultsFromServer, this.props.pollInterval);
    },
    handleOpsAuthSubmit: function(formData) {
        console.log(formData);
        var params = jQuery.param({
            opsHost: formData.opsHost.trim()
        });

        this.opsHost = formData.opsHost.trim();

        $.ajax({
            url: "http://localhost:8888/api/4.0/user/login?" + params,
            dataType: 'json',
            type: 'POST',
            crossDomain: true,
            data: JSON.stringify(formData.credentials)//,
        });
    },
    handleGetDeliveryServicesSubmit: function(data) {
        return $.ajax({
            url: "http://localhost:8888/api/4.0/deliveryservices.json?opsHost=" + this.opsHost,
            dataType: 'json',
            type: 'GET',
            success: function(data) {
                console.log("Got deliveryservices in ajax call");
                this.deliveryServices = data.response;
            }.bind(this)
        });
    },
    handleLoadTestSubmit: function(formData) {
        this.averageLatency = NaN;
        this.minLatency = NaN;
        this.maxLatency = NaN;
        this.medianLatency = NaN;
        this.startTime = moment();
        this.runningTime = 0;

        this.setState({latencies: []});

        formData.txCount = parseInt(formData.txCount);
        formData.connections = parseInt(formData.connections);

        this.totalRequests = formData.txCount * (formData.httpDeliveryServices.length * formData.httpsDeliveryServices.length);

        $.ajax({
            url: this.props.url,
            dataType: 'json',
            type: 'POST',
            data: JSON.stringify(formData),
            success: function(stuff) {
                console.log("posted request to start load test");
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },
    handleSubmit: function(e) {
        e.preventDefault();
        console.log("handle submit");
        var opsHost = this.state.opsHost.trim();

        console.log("http://localhost:8888/api/4.0/deliveryservices.json?opsHost=" + opsHost);

        $.ajax({
            url: "http://localhost:8888/api/4.0/deliveryservices.json?opsHost=" + opsHost,
            dataType: 'json',
            type: 'GET',
            success: function(data) {
                console.log("Got deliveryservices in ajax call");
                this.deliveryServices = data.response;
            }.bind(this)
        });
    },
    render: function() {
        this.cdns = [];


        if (this.deliveryServices != null) {
            var cdnNames = unique(this.deliveryServices.map(function (item) { return item.cdnName}))
                .filter(function(item) { return item.toLowerCase() != 'all'});

            cdnNames.forEach(function (cdnName) {
                var dsList = this.deliveryServices.filter(function (item) {
                    return item.cdnName.toLowerCase() == cdnName;
                });

                var cdn = {
                    name : cdnName,
                    deliveryServices: dsList
                };
                this.cdns.push(cdn);
            }.bind(this));
        }

        if (this.state.data != null && this.state.data.length > 0) {
            var sortedData = this.state.data.sort(function (result1, result2) {
                return result1.latency - result2.latency;
            });

            this.minLatency = (sortedData[0].latency / 1000.000).toFixed(3);
            this.maxLatency = (sortedData[sortedData.length - 1].latency / 1000.000).toFixed(3);

            var s = 0.000;

            for (var i = 0; i < sortedData.length; i++) {
                var l = (sortedData[i].latency / 1000.0).toFixed(3);
                if (this.medianLatency == NaN && (i > this.state.data.length / 2)) {
                    this.medianLatency = l;
                }

                s = +s + +l;
            }

            this.averageLatency = (s / this.state.data.length).toFixed(3);

            if (this.state.data.length < this.totalRequests) {
                this.runningTime = moment.duration(moment(new Date()).diff(this.startTime));
            }

            this.tps = NaN;
            if (this.runningTime != null && this.state.data.length > 0) {
                this.tps = (this.state.data.length / this.runningTime.asSeconds()).toFixed(3);
            }

            this.startTimePretty = "";
            if (this.startTime != null) {
                this.startTimePretty = this.startTime.format();
            }

            this.runningTimePretty = "";
            if (this.runningTime != null) {
                this.runningTimePretty = this.runningTime.asSeconds();
            }
        }

        return (
            <div className="well">
                <LoginBox onOpsAuthSubmit={this.handleOpsAuthSubmit}/>
                <h3>2. Fetch current data from lab</h3>
                <form className="form-inline" onSubmit={this.handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="opsHostInput">Traffic Ops Host</label>
                        <input id="opsHostInput" className="form-control" type="text" value={this.state.opsHost} onChange={this.handleOpsHostChange}/>
                    </div>
                    <button type="submit" className="btn btn-info">Get Lab Data</button>
                </form>
                <br/>
                <CdnList cdns={this.cdns}/>
                <CdnTestForm cdns={this.cdns} onLoadTestSubmit={this.handleLoadTestSubmit}/>
                <h1>Results</h1>
                <dl className="dl-horizontal">
                    <dt>Start Time</dt><dd>{this.startTimePretty}</dd>
                    {/*<dt>Run Time</dt><dd>{this.runningTimePretty}</dd>*/}
                    {/*<dt>TPS</dt><dd>{this.tps}</dd>*/}
                    {/*<dt>Median Latency</dt><dd>{this.medianLatency}</dd>*/}
                    {/*<dt>Average Latency</dt><dd>{this.averageLatency}</dd>*/}
                    {/*<dt>Max Latency</dt><dd>{this.maxLatency}</dd>*/}
                    {/*<dt>Min Latency</dt><dd>{this.minLatency}</dd>*/}
                </dl>
                <div id="asdf"></div>
                <ResultsList data={this.state.data}/>
            </div>
        );
    }
});

var LoginBox = React.createClass({
    getInitialState: function() {
        return {opsUser: "guest", opsPassword: "foo", opsHost: "traffic-ops.example.com", opsCookie: "peanutbutter"}
    },
    handleOpsHostChange: function(e) {
        this.setState({opsHost: e.target.value})
    },
    handleOpsUserChange: function(e) {
        this.setState({opsUser: e.target.value})
    },
    handleOpsPasswordChange: function(e) {
        this.setState({opsPassword: e.target.value})
    },
    handleSubmit: function(e) {
        e.preventDefault();
        var opsUser = this.state.opsUser.trim();
        var opsPassword = this.state.opsPassword.trim();
        var opsHost = this.state.opsHost.trim();

        this.props.onOpsAuthSubmit({opsHost: opsHost, credentials: {u: opsUser, p: opsPassword}});
    },
    render: function() {
        return(
            <div className="loginBox">
                <h3>1. Sign In</h3>
                <form className="form-inline" onSubmit={this.handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="opsHost">Traffic Ops Host</label>
                        <input id="opsHost" className="form-control" type="text" value={this.state.opsHost} onChange={this.handleOpsHostChange}/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="opsUser">User</label>
                        <input id="opsUser" className="form-control" type="text" value={this.state.opsUser} onChange={this.handleOpsUserChange}/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="opsPassword">Password</label>
                        <input id="opsPassword" className="form-control" type="password" value={this.state.opsPassword} onChange={this.handleOpsPasswordChange} />
                    </div>
                    <button type="submit" className="btn btn-info">Authenticate</button>
                </form>
            </div>
        )
    }
});

var SubResult = React.createClass({
    render: function() {
        return (
            <div className="subresult row">
                <div className="col-sm-3">{this.props.subresult.requestTime}</div>
                <div className="col-sm-3">{this.props.subresult.latency}</div>
                <div className="col-sm-3">{this.props.subresult.error}</div>
                <div className="col-sm-3">{this.props.subresult.httpStatus}</div>
            </div>
        )
    }
});

var SubresultsList = React.createClass({
    render: function() {
        for (var i = 0; i < this.props.subresults.length; i++) {
            this.props.subresults[i].id = i;
        }

        var resultNodes = this.props.subresults.map(function (subresult) {
            return (
                <SubResult key={subresult.id} subresult={subresult}/>
            )
        });

        return (
            <div className="SubResultsList">
                {resultNodes}
            </div>
        )
    }
});

var Result = React.createClass({
    render: function () {
        this.props.result.hid = "#" + this.props.result.id;
        return (
            <li className="list-group-item">
                <h4>
                    <button type="button" className="btn btn-info" data-toggle="collapse" data-target={this.props.result.hid}>Details</button>
                    {this.props.result.avgLatency} mSec avg latency for {this.props.result[0].host}
                </h4>
                <div id={this.props.result.id} className="collapse">
                    <div className="row grid-table-header">
                        <div className="col-sm-3">Request Time</div>
                        <div className="col-sm-3">Latency uSec</div>
                        <div className="col-sm-3">Error</div>
                        <div className="col-sm-3">Http Status</div>
                    </div>
                    <SubresultsList subresults={this.props.result}/>
                </div>
            </li>
        )
    }
});

var ResultsList = React.createClass({
    render: function() {
        if (this.props.data == null) {
            return (
                <div>Waiting for results</div>
            )
        }

        histogram(this.props.data.map(function (r, i) { return (r.latency / 1000.0).toFixed(3)}));
        var groupedResults = groupBy(this.props.data, result => JSON.stringify([result.host]));

        for (var i = 0; i < groupedResults.length; i++) {
            groupedResults[i].host = groupedResults[i][0].host;
            groupedResults[i].id = i;
            var sum = 0;
            for (var j = 0; j < groupedResults[i].length; j++) {
                sum += groupedResults[i][j].latency;
            }
            groupedResults[i].avgLatency = ((sum / groupedResults[i].length) / 1000.0).toFixed(3);
        }

        var resultNodes = groupedResults.map(function(result) {
            return (
                <Result key={result.id} result={result}/>
            )
        });

        return (
            <ul className="list-group">
                {resultNodes}
            </ul>
        )
    }
});

ReactDOM.render(
    <LabBox url="http://localhost:8888/loadtest" pollInterval={500} />,
    document.getElementById('loadtest')
);

var histogram = function(latencies) {
    var formatCount = d3.format(",.000f");

    var margin = {top: 10, right: 30, bottom: 30, left: 30},
        width = 960 - margin.left - margin.right,
        height = 500 - margin.top - margin.bottom;

    var x = d3.scaleLog()
        .domain([latencies[0], latencies[latencies.length -1]])
        .rangeRound([0, width]);

    var bins = d3.histogram()
        .domain(x.domain())
        .thresholds(x.ticks(20))
        (latencies);

    var y = d3.scaleLinear()
        .domain([0, d3.max(bins, function(d) { return d.length; })])
        .range([height, 0]);

    d3.select("#asdf > svg").remove();

    var svg = d3.select("#asdf").append("svg:svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");


    var bar = svg.selectAll(".bar")
        .data(bins)
        .enter().append("g")
        .attr("class", "bar")
        .attr("transform", function(d) { return "translate(" + x(d.x0) + "," + y(d.length) + ")"; });

    bar.append("rect")
        .attr("x", 1)
        .attr("width", x(bins[0].x1) - x(bins[0].x0) - 1)
        .attr("height", function(d) { return height - y(d.length); });

    bar.append("text")
        .attr("dy", "-0.75em")
        .attr("y", 6)
        .attr("x", (x(bins[0].x1) - x(bins[0].x0)) / 2)
        .attr("text-anchor", "middle")
        .text(function(d) { return formatCount(d.length); });

    svg.append("g")
        .attr("class", "axis axis--x")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x).tickFormat(d3.format(",d")));

    svg.append("text")
        .attr("text-anchor", "middle")
        .attr("transform", "translate(" + (width/4) + ",0)")
        .text("Latency in msec");

    svg.append("text")
        .attr("text-anchor", "middle")
        .attr("transform", "translate(0," + (height/2) + ")rotate(-90)")
        .text("# of Requests");
};
