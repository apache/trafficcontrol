// needed for self signed certificate
// link... https://github.com/meteor/meteor/issues/2866
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
Future = Npm.require('fibers/future');

Meteor.methods({

    // The method expects a valid IPv4 address
    'loginCall': function (username, password) {
        // Construct the API URL
        var myFuture = new Future();
        check(username, String);
        check(password, String);
        var apiUrl = 'https://ec2-52-37-126-44.us-west-2.compute.amazonaws.com:9000/login';
        var response = null;
        // query the API
        //var response = HTTP.get(apiUrl).data;
        HTTP.call("POST", apiUrl,
            {data: {"username": username, "password": password}},
            function (error, result) {
                if (!error) {
                    /*var tokenValue = null;
                    if (result.hasOwnProperty('content')) {
                        tokenValue = result.content;
                        tokenValue = JSON.parse(tokenValue);
                        if (tokenValue.hasOwnProperty('Token')) {
                            tokenValue = tokenValue.Token;
                            myFuture.return(tokenValue);
                        }
                    }*/
                    myFuture.return(result);
                    console.log("bif pipma ==> ", result.content);
                } else {
                    console.log("error2 ===> ", error.toString());
                    myFuture.return(error.response);
                    //myFuture.throw(error);
            }
        });
        return myFuture.wait();
    },

    'registerUser': function (username, firstName, lastName, password) {
        // Construct the API URL
        var myFuture = new Future();
        check(username, String);
        check(password, String);
        check(firstName, String);
        check(lastName, String);
        var apiUrl = 'https://ec2-52-37-126-44.us-west-2.compute.amazonaws.com:9000/users';

        HTTP.call("POST", apiUrl,
            {data: {"username": username, "password": password, "lastName": lastName, "firstName": firstName}},
            function (error, result) {
                if (!error) {
                    myFuture.return(result);
                } else {
                    myFuture.return(error.response);
                    //myFuture.throw(error);
                }
            });
        return myFuture.wait();
    },

    'getCameras': function (token) {
        var myFuture = new Future();
        check(token, String);
        token = 'Bearer ' + token;
        var apiURL = "https://ec2-52-37-126-44.us-west-2.compute.amazonaws.com:9000/8001/r";

        HTTP.call("POST", apiURL,
            { headers: { 'Authorization': token} },
            function (error, result) {
                //result.statuscode
                if (!error) {
                    console.log("getcameras server, response ", result.content);
                    myFuture.return(result.content);
                } else {
                    console.log("error ===> ", error.toString());
                    myFuture.throw(error);
                }
            });
        return myFuture.wait();
    },

    'getVideos': function (token) {

    },

    'getLiveFeed': function (token) {

    },

    'saveFeed': function (token) {

    },

    'controlFeed': function (token) {

    },

    'registerCamera': function (token) {

    }
});