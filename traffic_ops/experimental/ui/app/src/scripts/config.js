"use strict";

 angular.module("config", [])

.constant("ENV", {
  "apiEndpoint": {
    "login": "/api/2.0/login",
    "get_current_user": "/api/2.0/tm_user/current",
    "update_current_user": "/api/2.0/tm_user/current",
    "get_users": "/api/2.0/tm_user"
  }
})

;