# GraphQL Testing

## Getting started
1. Get docker and docker compose working
2. `docker compose up -d`
3. Connect to PGAdmin at `localhost:80`
4. Create a new database `traffic_ops` owned by a new role `traffic_ops`
5. Restore a TO database dump
6. Stop the local env `docker compose down`
7. `docker compose up`
8. Connect to the interactive query tester `localhost:5000/graphiql` and start experimenting

### PGModler (optional, mac)
1. `brew cask install xquartz`
2. reboot
3. `brew install socat`
4. `source launch.pgmodler.rc`
5. `pgmodeler.indockernet`

## Gotchas (Out-of-box generation)
1. Since we used autoinc primary keys, conditionals are going to make less sense.
2. Since we jammed raw json into the database, you can't filter DS Requests.
3. Since we reimplement security in the application rather than the database, this should never be publicly exposed (https://blog.2ndquadrant.com/application-users-vs-row-level-security/)

## References
* https://philsturgeon.uk/api/2017/01/24/graphql-vs-rest-overview/
* https://github.com/graphile/postgraphile
* https://www.graphile.org/postgraphile/
* https://graphql.org/users/

## Sample Query

### overly verbbose, filter to taste
```
{
  allServers(condition: {hostName: "odol-atsec-den-01"}) {
    edges {
      node {
        hostName
        domainName
        cdnByCdnId {
          name
        }
        profileByProfile {
          name
          profileParametersByProfile {
            nodes {
              parameterByParameter {
                name
                value
                configFile
              }
            }
          }
        }
        deliveryserviceServersByServer {
          nodes {
            deliveryserviceByDeliveryservice {
              xmlId
              active
              dscp
              signingAlgorithm
              qstringIgnore
              geoLimit
              httpBypassFqdn
              dnsBypassIp
              dnsBypassIp6
              typeByType {
                name
              }
              profileByProfile {
                name
                profileParametersByProfile {
                  nodes {
                    parameterByParameter {
                      name
                      value
                      configFile
                    }
                  }
                }
              }
              ccrDnsTtl
              globalMaxTps
              globalMaxMbps
              longDesc
              longDesc1
              longDesc2
              maxDnsAnswers
              infoUrl
              missLat
              missLong
              checkPath
              protocol
              sslKeyVersion
              ipv6RoutingEnabled
              rangeRequestHandling
              edgeHeaderRewrite
              originShield
              midHeaderRewrite
              regexRemap
              cacheurl
              remapText
              multiSiteOrigin
              displayName
              trResponseHeaders
              initialDispersion
              dnsBypassCname
              trRequestHeaders
              regionalGeoBlocking
              geoProvider
              geoLimitCountries
              logsEnabled
              multiSiteOriginAlgorithm
              geolimitRedirectUrl
              tenantByTenantId {
                name
              }
              routingName
              deepCachingType
              fqPacingRate
              anonymousBlockingEnabled
              deliveryserviceRegexesByDeliveryservice {
                nodes {
                  regexByRegex {
                    pattern
                    typeByType {
                      name
                    }
                  }
                }
              }
              federationDeliveryservicesByDeliveryservice {
                nodes {
                  federationByFederation {
                    cname
                    description
                    ttl
                    federationFederationResolversByFederation {
                      nodes {
                        federationResolverByFederationResolver {
                          ipAddress
                          typeByType {
                            name
                          }
                        }
                      }
                    }
                  }
                }
              }
              jobsByJobDeliveryservice {
                nodes {
                  jobAgentByAgent {
                    name
                    description
                    active
                  }
                  objectType
                  objectName
                  keyword
                  parameters
                  assetUrl
                  assetType
                  jobStatusByStatus {
                    name
                    description
                  }
                  startTime
                  enteredTime
                  tmUserByJobUser {
                    username
                    fullName
                  }
                }
              }
              originsByDeliveryservice {
                nodes {
                  name
                  fqdn
                  protocol
                  isPrimary
                  port
                  ipAddress
                  ip6Address
                  coordinateByCoordinate {
                    name
                    latitude
                    longitude
                  }
                  profileByProfile {
                    name
                    profileParametersByProfile {
                      nodes {
                        parameterByParameter {
                          name
                          value
                          configFile
                        }
                      }
                    }
                  }
                }
              }
              staticdnsentriesByDeliveryservice {
                nodes {
                  host
                  address
                  typeByType {
                    name
                  }
                  ttl
                }
              }
              steeringTargetsByDeliveryservice {
                nodes {
                  value
                  typeByType {
                    name
                  }
                  deliveryserviceByTarget {
                    xmlId
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

### Pagination
```
{
  allServers(first:3 after:"WyJwcmltYXJ5X2tleV9hc2MiLFsyMiwxMSwxLDMsODExXV0=") {
    totalCount
    pageInfo{
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    edges {
      node {
        hostName
      }
    }
  }
}
```

## Testing
`curl -H "Content-Type: application/json" http://localhost:5000/graphql --data '{ "query": "{ allServers { edges { node { hostName } } } }"}' | jq`
