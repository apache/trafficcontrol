Template.browseCameras.helpers({
     availableCameras: function() {
         //AvailableCameras.insert({cameraName: "livingroom", date: "jan 1st 2016"});
         return AvailableCameras.find();
    }
});

Template.browseCameras.onCreated(function() {
   alert("on created method fired!");

    var login_data = Session.get('login_response');
    var token = null;
    if (login_data.hasOwnProperty('token')) {
        token = login_data.token;
    }

    Meteor.call('getCameras', token, function(err, res) {
        if (err) {
            console.log('getCameras error... ', err.message);
        } else {
            
            console.log("gerCameras response... ", res.content);
        }
        return res;

    });

});