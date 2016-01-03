
package main
//This file is generated automatically. Do not try to edit it manually.

var resourceListingJson = `{
    "apiVersion": "2.0 alpha",
    "swaggerVersion": "1.2",
    "basePath": "{{.}}",
    "apis": [
        {
            "path": "/api/2.0",
            "description": "retrieves Asn information"
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
var apiDescriptionsJson = map[string]string{"api/2.0":`{
    "apiVersion": "2.0 alpha",
    "swaggerVersion": "1.2",
    "basePath": "{{.}}",
    "resourcePath": "/api/2.0",
    "apis": [
        {
            "path": "/api/2.0/asn/{id}",
            "description": "retrieves Asn information",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getAsn",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Asn"
                    },
                    "summary": "retrieves Asn information",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "Asn id",
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
            "path": "/api/2.0/cachegroup/{id}",
            "description": "retrieves cachegroup information",
            "operations": [
                {
                    "httpMethod": "GET",
                    "nickname": "getCachegroup",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                    },
                    "summary": "retrieves cachegroup information",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "cachegroup id",
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
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                    },
                    "summary": "modify an existing cachegroup",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "cachegroup id",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "shortName",
                            "description": "cachegroup short name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "name",
                            "description": "cachegroup name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "longitude",
                            "description": "Location longitude",
                            "dataType": "gopkg.in.guregu.null.v3.Float",
                            "type": "gopkg.in.guregu.null.v3.Float",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "latitide",
                            "description": "Location latitiude",
                            "dataType": "gopkg.in.guregu.null.v3.Float",
                            "type": "gopkg.in.guregu.null.v3.Float",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "type",
                            "description": "CG type",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "parentCachegroupId",
                            "description": "Parent cachegroup id",
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
                    "httpMethod": "DELETE",
                    "nickname": "delCachegroup",
                    "type": "array",
                    "items": {
                        "$ref": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                    },
                    "summary": "deletes a cachegroup",
                    "parameters": [
                        {
                            "paramType": "path",
                            "name": "id",
                            "description": "cachegroup id",
                            "dataType": "int",
                            "type": "int",
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
                            "responseType": "array",
                            "responseModel": "github.com.Comcast.traffic_control.traffic_ops.goto2.api.Cachegroup"
                        }
                    ]
                }
            ]
        },
        {
            "path": "/api/2.0/cachegroup",
            "description": "enter a new cachegroup",
            "operations": [
                {
                    "httpMethod": "POST",
                    "nickname": "postCachegroup",
                    "type": "github.com.Comcast.traffic_control.traffic_ops.goto2.output_format.ApiWrapper",
                    "items": {},
                    "summary": "enter a new cachegroup",
                    "parameters": [
                        {
                            "paramType": "json",
                            "name": "shortName",
                            "description": "cachegroup short name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "name",
                            "description": "cachegroup name",
                            "dataType": "string",
                            "type": "string",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "longitude",
                            "description": "Location longitude",
                            "dataType": "gopkg.in.guregu.null.v3.Float",
                            "type": "gopkg.in.guregu.null.v3.Float",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "latitide",
                            "description": "Location latitiude",
                            "dataType": "gopkg.in.guregu.null.v3.Float",
                            "type": "gopkg.in.guregu.null.v3.Float",
                            "format": "",
                            "allowMultiple": false,
                            "required": false,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "type",
                            "description": "CG type",
                            "dataType": "int",
                            "type": "int",
                            "format": "",
                            "allowMultiple": false,
                            "required": true,
                            "minimum": 0,
                            "maximum": 0
                        },
                        {
                            "paramType": "json",
                            "name": "parentCachegroupId",
                            "description": "Parent cachegroup id",
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
        }
    }
}`,}
