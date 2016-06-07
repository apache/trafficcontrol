module.exports = {
    timeout: '120s',
    useSSL: false, // set to true if using ssl
    port: 8080,
    sslPort: 8443,
    proxyPort: 8009,
    ssl: {
        key:    './ssl/tls/private/ssl.key',
        cert:   './ssl/tls/certs/ssl.crt',
        ca:     [
            './ssl/tls/certs/ssl-bundle.crt'
        ]
    },
    api: {
        base_url: 'http://localhost:3000',
        key: ''
    },
    files: {
        static: './app/dist/public/'
    },
    log: {
        stream: './server/log/access.log'
    },
    reject_unauthorized: false // false if using self-signed certs, true if trusted certs
};

