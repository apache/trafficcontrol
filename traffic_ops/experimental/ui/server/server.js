var constants = require('constants'),
    express = require('express'),
    http = require('http'),
    https = require('https'),
    path = require('path'),
    fs = require('fs'),
    morgan = require('morgan'),
    errorhandler = require('errorhandler'),
    modRewrite = require('connect-modrewrite'),
    timeout = require('connect-timeout');

var config;

try {
    // this should exist in prod environment. no need to create this file in dev as it will use the fallback (see catch)
    config = require('/etc/traffic_ops_v2/conf/config');
}
catch(e) {
    // this is used for dev environment
    config = require('../conf/config');
}

var logStream = fs.createWriteStream(config.log.stream, { flags: 'a' });
var useSSL = config.useSSL;

// Disable for self-signed certs in dev/test
process.env.NODE_TLS_REJECT_UNAUTHORIZED = config.reject_unauthorized;

var app = express();
// Add a handler to inspect the req.secure flag (see
// http://expressjs.com/api#req.secure). This allows us
// to know whether the request was via http or https.
app.all ("/*", function (req, res, next) {
    if (useSSL && !req.secure) {
        var headersHost = req.headers.host.split(':');
        var httpsUrl = 'https://' + headersHost[0] + ':' +  config.sslPort + req.url;
        // request was via http, so redirect to https
        res.redirect(httpsUrl);
    } else {
        next();
    }
});

app.use(modRewrite([
        '^/api/(.*?)\\?(.*)$ ' + config.api.base_url + '/api/$1?$2&api_key=' + config.api.key + ' [P]', // match /api/{version}/foos?active=true and replace with api.base_url/api/{version}/foos?active=true&api_key=api.key
        '^/api/(.*)$ ' + config.api.base_url + '/api/$1?api_key=' + config.api.key + ' [P]' // match /api/{version}/foos and replace with api.base_url/api/{version}/foos?api_key=api.key
]));
app.use(express.static(config.files.static));
app.use(morgan('combined', {
    stream: logStream,
    skip: function (req, res) { return res.statusCode < 400 }
}));
app.use(errorhandler());
app.use(timeout(config.timeout));

if (app.get('env') === 'dev') {
    app.use(require('connect-livereload')({
        port: 35728,
        excludeList: ['.woff', '.flv']
    }));
} else {
    app.set('env', "production");
}

// Enable reverse proxy support in Express. This causes the
// the "X-Forwarded-Proto" header field to be trusted so its
// value can be used to determine the protocol. See
// http://expressjs.com/api#app-settings for more details.
app.enable('trust proxy');

// Startup HTTP Server
var httpServer = http.createServer(app);
httpServer.listen(config.port);

if (useSSL) {
    //
    // Supply `SSL_OP_NO_SSLv3` constant as secureOption to disable SSLv3
    // from the list of supported protocols that SSLv23_method supports.
    //
    var sslOptions = {};
    sslOptions['secureOptions'] = constants.SSL_OP_NO_SSLv3;

    sslOptions['key'] = fs.readFileSync(config.ssl.key);
    sslOptions['cert'] = fs.readFileSync(config.ssl.cert);
    sslOptions['ca'] = config.ssl.ca.map(function(cert){
        return fs.readFileSync(cert);
    });

    // Startup HTTPS Server
    var httpsServer = https.createServer(sslOptions, app);
    httpsServer.listen(config.sslPort);

    sslOptions.agent = new https.Agent(sslOptions);
}

console.log("Traffic Ops Port         : %s", config.port);
console.log("Traffic Ops Proxy Port   : %s", config.proxyPort);
console.log("Traffic Ops SSL Port     : %s", config.sslPort);
