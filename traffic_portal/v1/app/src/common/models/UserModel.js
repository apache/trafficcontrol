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

var UserModel = function($rootScope, messageModel) {

    var user = {};
    user.loaded = false;
    this.user = user;

    this.setUser = function(userData) {
        user.loaded = true;
        user = angular.extend(user, userData);
        if (user.newUser) {
            user.username = ''; // new users were given a temp username that needs to be ditched
        }
        if (!user.localUser) {
            messageModel.setMessages([ { level: 'success', text: 'Logged in as read-only user.' } ], false);
        }
        $rootScope.$broadcast('userModel::userUpdated', user);
    };

    this.resetUser = function() {
        user = {};
        user.loaded = false;
        this.user = user;
        $rootScope.$broadcast('userModel::userUpdated', user);
    };

};

UserModel.$inject = ['$rootScope', 'messageModel'];
module.exports = UserModel;