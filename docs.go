package main

//This file is generated automatically. Do not try to edit it manually.

var resourceListingJson = `{
    "apiVersion": "2.0 alpha",
    "swaggerVersion": "1.2",
    "basePath": "{{.}}",
    "apis": [
        {
            "path": "/api/2.0",
            "description": "retrieves the profile_parameter information for a certain id"
        },
        {
            "path": "/api/2",
            "description": "Version 2.0 API"
        }
    ],
    "info": {
        "title": "Traffic Operations",
        "description": "Traffic Ops API",
        "contact": "https://traffic-control-cdn.net",
        "license": "Apache 2.0",
        "licenseUrl": "http://www.apache.org/licenses/LICENSE-2.0"
    }
}`
var apiDescriptionsJson = map[string]string{"api/2.0": `{
    "apiVersion": "2.0 alpha",
    "swaggerVersion": "1.2",
    "basePath": "{{.}}",
    "resourcePath": "/api/2.0",
    "apis": [
        {
            "path": "/api/2.0/profile_parameter/{id}",
            "description": "retrieves the profile_parameter information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getProfileParameterById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter"
                    },
                    "summary": "retrieves the profile_parameter information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putProfileParameter",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing profile_parameterentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "ProfileParameter object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delProfileParameterById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter"
                    },
                    "summary": "deletes profile_parameter information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/profile_parameter",
            "description": "retrieves the profile_parameter information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getProfileParameters",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter"
                    },
                    "summary": "retrieves the profile_parameter information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postProfileParameter",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new profile_parameter",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "ProfileParameter object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice_regex/{id}",
            "description": "retrieves the deliveryservice_regex information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceRegexById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex"
                    },
                    "summary": "retrieves the deliveryservice_regex information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putDeliveryserviceRegex",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing deliveryservice_regexentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "DeliveryserviceRegex object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delDeliveryserviceRegexById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex"
                    },
                    "summary": "deletes deliveryservice_regex information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice_regex",
            "description": "retrieves the deliveryservice_regex information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceRegexs",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex"
                    },
                    "summary": "retrieves the deliveryservice_regex information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postDeliveryserviceRegex",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new deliveryservice_regex",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "DeliveryserviceRegex object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/servercheck/{id}",
            "description": "retrieves the servercheck information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getServercheckById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck"
                    },
                    "summary": "retrieves the servercheck information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putServercheck",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing servercheckentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Servercheck object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delServercheckById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck"
                    },
                    "summary": "deletes servercheck information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/servercheck",
            "description": "retrieves the servercheck information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getServerchecks",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck"
                    },
                    "summary": "retrieves the servercheck information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postServercheck",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new servercheck",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Servercheck object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/staticdnsentry/{id}",
            "description": "retrieves the staticdnsentry information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getStaticdnsentryById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry"
                    },
                    "summary": "retrieves the staticdnsentry information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putStaticdnsentry",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing staticdnsentryentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Staticdnsentry object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delStaticdnsentryById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry"
                    },
                    "summary": "deletes staticdnsentry information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/staticdnsentry",
            "description": "retrieves the staticdnsentry information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getStaticdnsentrys",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry"
                    },
                    "summary": "retrieves the staticdnsentry information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postStaticdnsentry",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new staticdnsentry",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Staticdnsentry object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_resolver/{id}",
            "description": "retrieves the federation_resolver information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationResolverById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver"
                    },
                    "summary": "retrieves the federation_resolver information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putFederationResolver",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing federation_resolverentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationResolver object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delFederationResolverById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver"
                    },
                    "summary": "deletes federation_resolver information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_resolver",
            "description": "retrieves the federation_resolver information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationResolvers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver"
                    },
                    "summary": "retrieves the federation_resolver information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postFederationResolver",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new federation_resolver",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationResolver object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/role/{id}",
            "description": "retrieves the role information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getRoleById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role"
                    },
                    "summary": "retrieves the role information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putRole",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing roleentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Role object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delRoleById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role"
                    },
                    "summary": "deletes role information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/role",
            "description": "retrieves the role information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getRoles",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role"
                    },
                    "summary": "retrieves the role information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postRole",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new role",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Role object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cachegroup_parameter/{id}",
            "description": "retrieves the cachegroup_parameter information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCachegroupParameterById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter"
                    },
                    "summary": "retrieves the cachegroup_parameter information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putCachegroupParameter",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing cachegroup_parameterentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "CachegroupParameter object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delCachegroupParameterById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter"
                    },
                    "summary": "deletes cachegroup_parameter information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cachegroup_parameter",
            "description": "retrieves the cachegroup_parameter information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCachegroupParameters",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter"
                    },
                    "summary": "retrieves the cachegroup_parameter information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postCachegroupParameter",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new cachegroup_parameter",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "CachegroupParameter object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/division/{id}",
            "description": "retrieves the division information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDivisionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division"
                    },
                    "summary": "retrieves the division information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putDivision",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing divisionentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Division object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delDivisionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division"
                    },
                    "summary": "deletes division information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/division",
            "description": "retrieves the division information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDivisions",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division"
                    },
                    "summary": "retrieves the division information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postDivision",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new division",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Division object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/profile/{id}",
            "description": "retrieves the profile information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getProfileById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile"
                    },
                    "summary": "retrieves the profile information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putProfile",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing profileentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Profile object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delProfileById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile"
                    },
                    "summary": "deletes profile information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/profile",
            "description": "retrieves the profile information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getProfiles",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile"
                    },
                    "summary": "retrieves the profile information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postProfile",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new profile",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Profile object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/status/{id}",
            "description": "retrieves the status information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getStatusById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status"
                    },
                    "summary": "retrieves the status information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putStatus",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing statusentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Status object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delStatusById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status"
                    },
                    "summary": "deletes status information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/status",
            "description": "retrieves the status information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getStatuss",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status"
                    },
                    "summary": "retrieves the status information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postStatus",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new status",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Status object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/to_extension/{id}",
            "description": "retrieves the to_extension information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getToExtensionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension"
                    },
                    "summary": "retrieves the to_extension information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putToExtension",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing to_extensionentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "ToExtension object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delToExtensionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension"
                    },
                    "summary": "deletes to_extension information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/to_extension",
            "description": "retrieves the to_extension information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getToExtensions",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension"
                    },
                    "summary": "retrieves the to_extension information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postToExtension",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new to_extension",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "ToExtension object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cdn/{id}",
            "description": "retrieves the cdn information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCdnById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn"
                    },
                    "summary": "retrieves the cdn information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putCdn",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing cdnentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Cdn object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delCdnById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn"
                    },
                    "summary": "deletes cdn information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cdn",
            "description": "retrieves the cdn information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCdns",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn"
                    },
                    "summary": "retrieves the cdn information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postCdn",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new cdn",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Cdn object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/server/{id}",
            "description": "retrieves the server information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getServerById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server"
                    },
                    "summary": "retrieves the server information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putServer",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing serverentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Server object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delServerById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server"
                    },
                    "summary": "deletes server information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/server",
            "description": "retrieves the server information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getServers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server"
                    },
                    "summary": "retrieves the server information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postServer",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new server",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Server object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/hwinfo/{id}",
            "description": "retrieves the hwinfo information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getHwinfoById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo"
                    },
                    "summary": "retrieves the hwinfo information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putHwinfo",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing hwinfoentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Hwinfo object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delHwinfoById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo"
                    },
                    "summary": "deletes hwinfo information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/hwinfo",
            "description": "retrieves the hwinfo information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getHwinfos",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo"
                    },
                    "summary": "retrieves the hwinfo information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postHwinfo",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new hwinfo",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Hwinfo object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job_agent/{id}",
            "description": "retrieves the job_agent information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobAgentById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent"
                    },
                    "summary": "retrieves the job_agent information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putJobAgent",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing job_agententry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "JobAgent object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delJobAgentById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent"
                    },
                    "summary": "deletes job_agent information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job_agent",
            "description": "retrieves the job_agent information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobAgents",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent"
                    },
                    "summary": "retrieves the job_agent information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postJobAgent",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new job_agent",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "JobAgent object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cachegroup/{id}",
            "description": "retrieves the cachegroup information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCachegroupById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                    },
                    "summary": "retrieves the cachegroup information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putCachegroup",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing cachegroupentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Cachegroup object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delCachegroupById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                    },
                    "summary": "deletes cachegroup information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cachegroup",
            "description": "retrieves the cachegroup information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCachegroups",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                    },
                    "summary": "retrieves the cachegroup information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postCachegroup",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new cachegroup",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Cachegroup object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/phys_location/{id}",
            "description": "retrieves the phys_location information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getPhysLocationById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation"
                    },
                    "summary": "retrieves the phys_location information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putPhysLocation",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing phys_locationentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "PhysLocation object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delPhysLocationById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation"
                    },
                    "summary": "deletes phys_location information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/phys_location",
            "description": "retrieves the phys_location information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getPhysLocations",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation"
                    },
                    "summary": "retrieves the phys_location information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postPhysLocation",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new phys_location",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "PhysLocation object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job_status/{id}",
            "description": "retrieves the job_status information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobStatusById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus"
                    },
                    "summary": "retrieves the job_status information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putJobStatus",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing job_statusentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "JobStatus object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delJobStatusById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus"
                    },
                    "summary": "deletes job_status information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job_status",
            "description": "retrieves the job_status information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobStatuss",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus"
                    },
                    "summary": "retrieves the job_status information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postJobStatus",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new job_status",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "JobStatus object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation/{id}",
            "description": "retrieves the federation information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation"
                    },
                    "summary": "retrieves the federation information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putFederation",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing federationentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Federation object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delFederationById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation"
                    },
                    "summary": "deletes federation information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation",
            "description": "retrieves the federation information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederations",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation"
                    },
                    "summary": "retrieves the federation information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postFederation",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new federation",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Federation object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/log/{id}",
            "description": "retrieves the log information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getLogById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log"
                    },
                    "summary": "retrieves the log information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putLog",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing logentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Log object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delLogById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log"
                    },
                    "summary": "deletes log information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/log",
            "description": "retrieves the log information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getLogs",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log"
                    },
                    "summary": "retrieves the log information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postLog",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new log",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Log object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job_result/{id}",
            "description": "retrieves the job_result information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobResultById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult"
                    },
                    "summary": "retrieves the job_result information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putJobResult",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing job_resultentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "JobResult object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delJobResultById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult"
                    },
                    "summary": "deletes job_result information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job_result",
            "description": "retrieves the job_result information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobResults",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult"
                    },
                    "summary": "retrieves the job_result information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postJobResult",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new job_result",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "JobResult object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/type/{id}",
            "description": "retrieves the type information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getTypeById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type"
                    },
                    "summary": "retrieves the type information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putType",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing typeentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Type object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delTypeById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type"
                    },
                    "summary": "deletes type information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/type",
            "description": "retrieves the type information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getTypes",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type"
                    },
                    "summary": "retrieves the type information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postType",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new type",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Type object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/parameter/{id}",
            "description": "retrieves the parameter information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getParameterById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter"
                    },
                    "summary": "retrieves the parameter information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putParameter",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing parameterentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Parameter object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delParameterById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter"
                    },
                    "summary": "deletes parameter information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/parameter",
            "description": "retrieves the parameter information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getParameters",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter"
                    },
                    "summary": "retrieves the parameter information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postParameter",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new parameter",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Parameter object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/regex/{id}",
            "description": "retrieves the regex information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getRegexById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex"
                    },
                    "summary": "retrieves the regex information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putRegex",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing regexentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Regex object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delRegexById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex"
                    },
                    "summary": "deletes regex information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/regex",
            "description": "retrieves the regex information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getRegexs",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex"
                    },
                    "summary": "retrieves the regex information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postRegex",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new regex",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Regex object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/stats_summary/{id}",
            "description": "retrieves the stats_summary information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getStatsSummaryById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary"
                    },
                    "summary": "retrieves the stats_summary information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putStatsSummary",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing stats_summaryentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "StatsSummary object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delStatsSummaryById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary"
                    },
                    "summary": "deletes stats_summary information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/stats_summary",
            "description": "retrieves the stats_summary information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getStatsSummarys",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary"
                    },
                    "summary": "retrieves the stats_summary information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postStatsSummary",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new stats_summary",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "StatsSummary object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job/{id}",
            "description": "retrieves the job information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job"
                    },
                    "summary": "retrieves the job information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putJob",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing jobentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Job object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delJobById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job"
                    },
                    "summary": "deletes job information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/job",
            "description": "retrieves the job information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getJobs",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job"
                    },
                    "summary": "retrieves the job information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postJob",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new job",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Job object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/region/{id}",
            "description": "retrieves the region information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getRegionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region"
                    },
                    "summary": "retrieves the region information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putRegion",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing regionentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Region object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delRegionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region"
                    },
                    "summary": "deletes region information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/region",
            "description": "retrieves the region information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getRegions",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region"
                    },
                    "summary": "retrieves the region information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postRegion",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new region",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Region object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_federation_resolver/{id}",
            "description": "retrieves the federation_federation_resolver information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationFederationResolverById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver"
                    },
                    "summary": "retrieves the federation_federation_resolver information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putFederationFederationResolver",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing federation_federation_resolverentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationFederationResolver object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delFederationFederationResolverById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver"
                    },
                    "summary": "deletes federation_federation_resolver information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_federation_resolver",
            "description": "retrieves the federation_federation_resolver information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationFederationResolvers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver"
                    },
                    "summary": "retrieves the federation_federation_resolver information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postFederationFederationResolver",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new federation_federation_resolver",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationFederationResolver object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice_server/{id}",
            "description": "retrieves the deliveryservice_server information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceServerById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer"
                    },
                    "summary": "retrieves the deliveryservice_server information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putDeliveryserviceServer",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing deliveryservice_serverentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "DeliveryserviceServer object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delDeliveryserviceServerById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer"
                    },
                    "summary": "deletes deliveryservice_server information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice_server",
            "description": "retrieves the deliveryservice_server information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceServers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer"
                    },
                    "summary": "retrieves the deliveryservice_server information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postDeliveryserviceServer",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new deliveryservice_server",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "DeliveryserviceServer object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_tmuser/{id}",
            "description": "retrieves the federation_tmuser information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationTmuserById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser"
                    },
                    "summary": "retrieves the federation_tmuser information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putFederationTmuser",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing federation_tmuserentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationTmuser object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delFederationTmuserById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser"
                    },
                    "summary": "deletes federation_tmuser information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_tmuser",
            "description": "retrieves the federation_tmuser information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationTmusers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser"
                    },
                    "summary": "retrieves the federation_tmuser information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postFederationTmuser",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new federation_tmuser",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationTmuser object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice_tmuser/{id}",
            "description": "retrieves the deliveryservice_tmuser information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceTmuserById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser"
                    },
                    "summary": "retrieves the deliveryservice_tmuser information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putDeliveryserviceTmuser",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing deliveryservice_tmuserentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "DeliveryserviceTmuser object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delDeliveryserviceTmuserById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser"
                    },
                    "summary": "deletes deliveryservice_tmuser information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice_tmuser",
            "description": "retrieves the deliveryservice_tmuser information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceTmusers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser"
                    },
                    "summary": "retrieves the deliveryservice_tmuser information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postDeliveryserviceTmuser",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new deliveryservice_tmuser",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "DeliveryserviceTmuser object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/tm_user/{id}",
            "description": "retrieves the tm_user information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getTmUserById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser"
                    },
                    "summary": "retrieves the tm_user information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putTmUser",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing tm_userentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "TmUser object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delTmUserById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser"
                    },
                    "summary": "deletes tm_user information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/tm_user",
            "description": "retrieves the tm_user information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getTmUsers",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser"
                    },
                    "summary": "retrieves the tm_user information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postTmUser",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new tm_user",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "TmUser object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/asn/{id}",
            "description": "retrieves the asn information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getAsnById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                    },
                    "summary": "retrieves the asn information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putAsn",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing asnentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Asn object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delAsnById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                    },
                    "summary": "deletes asn information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/asn",
            "description": "retrieves the asn information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getAsns",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                    },
                    "summary": "retrieves the asn information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postAsn",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new asn",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Asn object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/goose_db_version/{id}",
            "description": "retrieves the goose_db_version information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getGooseDbVersionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion"
                    },
                    "summary": "retrieves the goose_db_version information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putGooseDbVersion",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing goose_db_versionentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "GooseDbVersion object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delGooseDbVersionById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion"
                    },
                    "summary": "deletes goose_db_version information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/goose_db_version",
            "description": "retrieves the goose_db_version information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getGooseDbVersions",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion"
                    },
                    "summary": "retrieves the goose_db_version information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postGooseDbVersion",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new goose_db_version",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "GooseDbVersion object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_deliveryservice/{id}",
            "description": "retrieves the federation_deliveryservice information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationDeliveryserviceById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice"
                    },
                    "summary": "retrieves the federation_deliveryservice information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putFederationDeliveryservice",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing federation_deliveryserviceentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationDeliveryservice object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delFederationDeliveryserviceById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice"
                    },
                    "summary": "deletes federation_deliveryservice information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/federation_deliveryservice",
            "description": "retrieves the federation_deliveryservice information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getFederationDeliveryservices",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice"
                    },
                    "summary": "retrieves the federation_deliveryservice information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postFederationDeliveryservice",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new federation_deliveryservice",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "FederationDeliveryservice object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice/{id}",
            "description": "retrieves the deliveryservice information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryserviceById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice"
                    },
                    "summary": "retrieves the deliveryservice information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice"
                        }
                    ]
                },
                {
                    "httpMethod": "PUT",
                    "nickname": "putDeliveryservice",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "modify an existing deliveryserviceentry",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Deliveryservice object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                },
                {
                    "httpMethod": "DELETE",
                    "nickname": "delDeliveryserviceById",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice"
                    },
                    "summary": "deletes deliveryservice information for a certain id",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "The row id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/deliveryservice",
            "description": "retrieves the deliveryservice information for a certain id",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getDeliveryservices",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice"
                    },
                    "summary": "retrieves the deliveryservice information for a certain id",
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice"
                        }
                    ]
                },
                {
                    "httpMethod": "POST",
                    "nickname": "postDeliveryservice",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new deliveryservice",
                    "parameters": [
                        {
                            "paramType": "body",
                            "name": "Body",
                            "description": "Deliveryservice object that should be added to the table",
                            "dataType": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice",
                            "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        }
                    ],
                    "responseMessages": [
                        {
                            "code": 200,
                            "message": "",
                            "responseType": "object",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper"
                        }
                    ]
                }
            ]
        }
    ],
    "models": {
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn",
            "properties": {
                "asn": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "cachegroup": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup",
            "properties": {
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "latitude": {
                    "type": "gopkg.in.guregu.null.v3.Float",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "longitude": {
                    "type": "gopkg.in.guregu.null.v3.Float",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "parentCachegroupId": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "shortName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.CachegroupParameter",
            "properties": {
                "cachegroup": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "parameter": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cdn",
            "properties": {
                "dnssecEnabled": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Deliveryservice",
            "properties": {
                "active": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "cacheurl": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ccrDnsTtl": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "cdnId": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "checkPath": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "displayName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "dnsBypassCname": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "dnsBypassIp": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "dnsBypassIp6": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "dnsBypassTtl": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "dscp": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "edgeHeaderRewrite": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "geoLimit": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "globalMaxMbps": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "globalMaxTps": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "httpBypassFqdn": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "infoUrl": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "initialDispersion": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ipv6RoutingEnabled": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "longDesc": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "longDesc1": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "longDesc2": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "maxDnsAnswers": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "midHeaderRewrite": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "missLat": {
                    "type": "gopkg.in.guregu.null.v3.Float",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "missLong": {
                    "type": "gopkg.in.guregu.null.v3.Float",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "multiSiteOrigin": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "orgServerFqdn": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "originShield": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "profile": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "protocol": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "qstringIgnore": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "rangeRequestHandling": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "regexRemap": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "remapText": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "signed": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "sslKeyVersion": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "trRequestHeaders": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "trResponseHeaders": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "xmlId": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceRegex",
            "properties": {
                "deliveryservice": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "regex": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "setNumber": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceServer",
            "properties": {
                "deliveryservice": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "server": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.DeliveryserviceTmuser",
            "properties": {
                "deliveryservice": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "tmUserId": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Division",
            "properties": {
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Federation",
            "properties": {
                "cname": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ttl": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationDeliveryservice",
            "properties": {
                "deliveryservice": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "federation": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationFederationResolver",
            "properties": {
                "federation": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "federationResolver": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationResolver",
            "properties": {
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ipAddress": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.FederationTmuser",
            "properties": {
                "federation": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "role": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "tmUser": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.GooseDbVersion",
            "properties": {
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "isApplied": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "tstamp": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "versionId": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Hwinfo",
            "properties": {
                "description": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "serverid": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "val": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Job",
            "properties": {
                "agent": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "assetType": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "assetUrl": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "enteredTime": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "jobDeliveryservice": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "jobUser": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "keyword": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "objectName": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "objectType": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "parameters": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "startTime": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "status": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobAgent",
            "properties": {
                "active": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobResult",
            "properties": {
                "agent": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "job": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "result": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.JobStatus",
            "properties": {
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Log",
            "properties": {
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "level": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "message": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ticketnum": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "tmUser": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Parameter",
            "properties": {
                "configFile": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "value": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.PhysLocation",
            "properties": {
                "address": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "city": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "comments": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "email": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "phone": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "poc": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "region": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "shortName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "state": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "zip": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Profile",
            "properties": {
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ProfileParameter",
            "properties": {
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "parameter": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "profile": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Regex",
            "properties": {
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "pattern": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Region",
            "properties": {
                "division": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Role",
            "properties": {
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "privLevel": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Server",
            "properties": {
                "cachegroup": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "cdnId": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "domainName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "hostName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "iloIpAddress": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "iloIpGateway": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "iloIpNetmask": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "iloPassword": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "iloUsername": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "interfaceMtu": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "interfaceName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ip6Address": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ip6Gateway": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ipAddress": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ipGateway": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ipNetmask": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "mgmtIpAddress": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "mgmtIpGateway": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "mgmtIpNetmask": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "physLocation": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "profile": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "rack": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "routerHostName": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "routerPortName": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "status": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "tcpPort": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "updPending": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "xmppId": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "xmppPasswd": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Servercheck",
            "properties": {
                "aa": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ab": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ac": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ad": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ae": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "af": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ag": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ah": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ai": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "aj": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ak": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "al": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "am": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "an": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ao": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ap": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "aq": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ar": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "as": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "at": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "au": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "av": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "aw": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ax": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ay": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "az": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ba": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "bb": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "bc": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "bd": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "be": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "server": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Staticdnsentry",
            "properties": {
                "address": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "cachegroup": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "deliveryservice": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "host": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "ttl": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.StatsSummary",
            "properties": {
                "cdnName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "deliveryserviceName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "statDate": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "statName": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "statValue": {
                    "type": "float64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "summaryTime": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Status",
            "properties": {
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.TmUser",
            "properties": {
                "addressLine1": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "addressLine2": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "city": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "company": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "confirmLocalPasswd": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "country": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "email": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "fullName": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "gid": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "localPasswd": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "localUser": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "newUser": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "phoneNumber": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "postalCode": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "registrationSent": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "role": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "stateOrProvince": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "token": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "uid": {
                    "type": "gopkg.in.guregu.null.v3.Int",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "username": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.ToExtension",
            "properties": {
                "additionalConfigJson": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "infoUrl": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "isactive": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "scriptFile": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "servercheckColumnName": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "servercheckShortName": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "type": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "version": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Type",
            "properties": {
                "description": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "id": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "lastUpdated": {
                    "type": "Time",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "name": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "useInTable": {
                    "type": "gopkg.in.guregu.null.v3.String",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.Alert": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.Alert",
            "properties": {
                "level": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "text": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper": {
            "id": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
            "properties": {
                "alerts": {
                    "type": "array",
                    "description": "",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.Alert"
                    },
                    "format": ""
                },
                "error": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "response": {
                    "type": "interface",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "version": {
                    "type": "float64",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "gopkg.in.guregu.null.v3.Float": {
            "id": "gopkg.in.guregu.null.v3.Float",
            "properties": {
                "Float64": {
                    "type": "float64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "Valid": {
                    "type": "bool",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "gopkg.in.guregu.null.v3.Int": {
            "id": "gopkg.in.guregu.null.v3.Int",
            "properties": {
                "Int64": {
                    "type": "int64",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "Valid": {
                    "type": "bool",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        },
        "gopkg.in.guregu.null.v3.String": {
            "id": "gopkg.in.guregu.null.v3.String",
            "properties": {
                "String": {
                    "type": "string",
                    "description": "",
                    "items": {},
                    "format": ""
                },
                "Valid": {
                    "type": "bool",
                    "description": "",
                    "items": {},
                    "format": ""
                }
            }
        }
    }
}`}
