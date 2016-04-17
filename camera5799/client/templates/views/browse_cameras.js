Template.browseCameras.helpers({
     availableCameras: function() {
         //AvailableCameras.insert({cameraName: "livingroom", date: "jan 1st 2016"});
         //return AvailableCameras.find();
         var a = [ {cameraName: "living room", cameraId: 123},
                   {cameraName: "family room", cameraId: 456},
                   {cameraName: "back yard", cameraId: 789},
                   {cameraName: "front yard", cameraId: 974}
         ];
         return a;
     }
});

Template.browseCameras.onCreated(function() {
    //alert("on created method fired!");

    var login_data = Session.get('login_response');
    var token = null;
    if (login_data.hasOwnProperty('token')) {
        token = login_data.token;
    }

    Meteor.call('getCameras', token, function(err, res) {
        if (err) {
            console.log('getCameras error... ', err.message);
        } else {
            AvailableCameras.insert({cameraName: "livingroom", date: "jan 1st 2016"});
            console.log("getCameras response from client... ", res.content);
        }
        return res;

    });

});