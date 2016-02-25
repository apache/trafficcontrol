"use strict";

 angular.module("config", [])

.constant("ENV", {
  "apiEndpoint": {
    "1.1": "/api/1.1/",
    "1.2": "/api/1.2/",
    "2.0": "/api/2.0/"
  }
})

;