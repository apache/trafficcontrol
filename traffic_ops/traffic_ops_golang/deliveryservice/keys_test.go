package deliveryservice

import (
	"encoding/pem"
	"testing"
)

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
 
const (
	BadCertData = "This is bad cert data and it is not pem encoded"
	BadKeyData  = "This is bad private key data and it not pem encoeded"

	PrivateKeyPKCS1RSA2048 = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEApIye2g6c0irFgoMAztPeEuV4ZX8NdORRsV2Ddi4kV9gKYiDk
xQ9eOal3Nea9kDEaUpRvk2nuVl5jbFG3sOore9Vs6YUsCvICDKCUPl5OrxURJ6BS
iC2l0jVE6FVhXf9pLX+OR8lZmMZ9GKbw1nlGGKwdF3T76gOZFSHQfT57vNFsI0Q+
8NhWbDclNo/AnPpQuBs6ocah83BAVQk3gTRMHO3+rMLuvXVppBBqD+P5OU+LRROr
joDulub2QUPiR0Q5DczqMCjDIQB96LReryN4dx/p4x4P62SLQB6dd2vHvJNmpfl/
CBzAdSLBRna9mSV6xw429tu5b9Z48Da5gp9igQIDAQABAoIBABoepDyS4zvNREri
RqeOJAs117WsxFMQxxLzeCGzU1uKVKOc+xN4zAk1KFIrDV4tHTOMkmWBBC87jmas
Vg9ELKDckQxEcmhOYBrnBoEb8TuDiZSTs2YgcNj8UbLbkrgcCfMJ82jbwlgo8cSP
A13YJFNYRsnpbO+JoKwlEPZAi92901MrwH5r6kf/D+q6SwLpABVvVux9UzlhKh/G
hGUcB1GhLIq1axwAtlKahkutxyUWabiTu2hGlaY8Q7JGid0t4GZHCvZ0cqQbS/bE
M2x8zskgWQrOPLneuxCYVLtpMaPdrJFSghcD39/Qw5kDiSw83m3VjrlluZG868+X
aooOOBECgYEA0YmSG09MppC2fQlmwZzbiZreygiKkrSwDvYwOwbn2z0DYwSqnJ3V
bHI0MIX4479s+1Lt3Rsr308GRIDmnx/u0gDlrAIZqh5Xoy3p2azQzCEUghr4csIp
sGEPE1FKzFL7UdQvJ3K/QeM2Xux11aF5B+jZh+2LrN3GStpZLyhyUdUCgYEAyQlK
5BzfUC6+c1Kuv6EO4YhlDy23AnTBoeGNEMOe8SQ3h4Fcz2TaOpWs2zXOLGoSoIHs
OTN5VYq3R8AQL4w75OS1Au9+nr4m0o7ix1YI3cBMwseMqQtR0OvrzOSZSm2fQZtO
SBeLPnLNg+XsR6Adxy7r2qSz9aYbXQA0WvHw9/0CgYEAxDbpJL27b3awDKKTINb8
Ff16hwI8kWi2PSx4ua2bzIdz9nNWONbsFmNTT+UEznBhY2+i4pwhFznvCpMSYwwK
HYlNiSdmVRGYy2uhQn87/wszIyqSYRRE6a/Z6CMFwhQq19O0XGJtiwtzzKvtJCHT
Ln7zxP/C/hunJk0Vmr1rYAkCgYEAkPvCpwCrjIgpkcHvhQQCV2SmfWvasErD2ptv
wMdTuVUFNxR0ep2hRN7s6qrDJgTZqigI1LfqqWaBB53cDm50Q38tjBBsoM9B8Fhb
9KZ3fnVQ5qhDKSaguqtqQzoZ0zN7xzTaH+Pa6A6jaJxI6t7umtecAPMHVgGVelzL
ZUtXHYECgYEAyOZQ9HevD5jaWE//wXggjdASE3CFu5QkKJ1MEhEz3O9DconflKh0
e2TVphH64nagLC8G3ta9aSHch56r4frnZ7wc6xX2sORP/PCg3iAmvEtzaevQHQue
T5BkD3zmC6RQIDOYN5DzpokgXZLuTCUl95evqAsIZ6Cd3tqNCqrJIfg=
-----END RSA PRIVATE KEY-----
`
	PrivateKeyPKCS8RSA2048 = `
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCkjJ7aDpzSKsWC
gwDO094S5Xhlfw105FGxXYN2LiRX2ApiIOTFD145qXc15r2QMRpSlG+Tae5WXmNs
Ubew6it71WzphSwK8gIMoJQ+Xk6vFREnoFKILaXSNUToVWFd/2ktf45HyVmYxn0Y
pvDWeUYYrB0XdPvqA5kVIdB9Pnu80WwjRD7w2FZsNyU2j8Cc+lC4GzqhxqHzcEBV
CTeBNEwc7f6swu69dWmkEGoP4/k5T4tFE6uOgO6W5vZBQ+JHRDkNzOowKMMhAH3o
tF6vI3h3H+njHg/rZItAHp13a8e8k2al+X8IHMB1IsFGdr2ZJXrHDjb227lv1njw
NrmCn2KBAgMBAAECggEAGh6kPJLjO81ESuJGp44kCzXXtazEUxDHEvN4IbNTW4pU
o5z7E3jMCTUoUisNXi0dM4ySZYEELzuOZqxWD0QsoNyRDERyaE5gGucGgRvxO4OJ
lJOzZiBw2PxRstuSuBwJ8wnzaNvCWCjxxI8DXdgkU1hGyels74mgrCUQ9kCL3b3T
UyvAfmvqR/8P6rpLAukAFW9W7H1TOWEqH8aEZRwHUaEsirVrHAC2UpqGS63HJRZp
uJO7aEaVpjxDskaJ3S3gZkcK9nRypBtL9sQzbHzOySBZCs48ud67EJhUu2kxo92s
kVKCFwPf39DDmQOJLDzebdWOuWW5kbzrz5dqig44EQKBgQDRiZIbT0ymkLZ9CWbB
nNuJmt7KCIqStLAO9jA7BufbPQNjBKqcndVscjQwhfjjv2z7Uu3dGyvfTwZEgOaf
H+7SAOWsAhmqHlejLenZrNDMIRSCGvhywimwYQ8TUUrMUvtR1C8ncr9B4zZe7HXV
oXkH6NmH7Yus3cZK2lkvKHJR1QKBgQDJCUrkHN9QLr5zUq6/oQ7hiGUPLbcCdMGh
4Y0Qw57xJDeHgVzPZNo6lazbNc4sahKggew5M3lVirdHwBAvjDvk5LUC736evibS
juLHVgjdwEzCx4ypC1HQ6+vM5JlKbZ9Bm05IF4s+cs2D5exHoB3HLuvapLP1phtd
ADRa8fD3/QKBgQDENukkvbtvdrAMopMg1vwV/XqHAjyRaLY9LHi5rZvMh3P2c1Y4
1uwWY1NP5QTOcGFjb6LinCEXOe8KkxJjDAodiU2JJ2ZVEZjLa6FCfzv/CzMjKpJh
FETpr9noIwXCFCrX07RcYm2LC3PMq+0kIdMufvPE/8L+G6cmTRWavWtgCQKBgQCQ
+8KnAKuMiCmRwe+FBAJXZKZ9a9qwSsPam2/Ax1O5VQU3FHR6naFE3uzqqsMmBNmq
KAjUt+qpZoEHndwObnRDfy2MEGygz0HwWFv0pnd+dVDmqEMpJqC6q2pDOhnTM3vH
NNof49roDqNonEjq3u6a15wA8wdWAZV6XMtlS1cdgQKBgQDI5lD0d68PmNpYT//B
eCCN0BITcIW7lCQonUwSETPc70Nyid+UqHR7ZNWmEfridqAsLwbe1r1pIdyHnqvh
+udnvBzrFfaw5E/88KDeICa8S3Np69AdC55PkGQPfOYLpFAgM5g3kPOmiSBdku5M
JSX3l6+oCwhnoJ3e2o0Kqskh+A==
-----END PRIVATE KEY-----
`
	PrivateKeyEncryptedRSA2048 = `
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,A6E695C7D7038524

YMB2rcNobMEfCRzLcGV5uSF5g1hMIuHlJjZHG0QYNpeFdP2pONtgGLH+NIVWbb7B
2dYHUbofywcJTGPigSIm0RUEHquz/+seauRXs0jvN0cEKgEyWVaQZB9yem+V5sJx
guicvKmbWfiV9u004n9G4ue9IxpWwoZp1KixOrJ7GOmRhJsNL8E9MUpJ8kD1eOIY
n8RBuo/sQMyDKinL7npk5IYE3N1tQuVAxpEn2KVXjDIIrl0SeCO5pfJdXA40Sw8Q
6bEF8GogDuEaibw2XzOfqzJsmQRks1STtCqKVnN7b7HjpY3pPgYqMqc9BVKWL3PF
a8G5BSGjjosu3SPuzIIlADlQm65O8Qi/6eZ9thByfGBtXcanSdYZtoprreEikhgC
MYnID949McXvnJTjB1snWni2FYhrWuH6iK4/7/UevFg74She5qKfwuw/YXe+p2pP
Z0ON4olUgG68/aReRoICqDurDW88zklbq8c18ECI8ppBy3tDB85Sey+ybwTliSd7
FLBhwk/vLuhBlX6kVba8CkBUHqu8T9zECyLK2f7geDUjsw6EVXYy7whJE9LL1Sq0
cSUpyMb5I6RXOe90zj59sBJddEsllSn2KTMhrJa2LPj+j6Df4PuwmsKdPz+DWo5r
TTAvc3iN9sZbXGjwbjiJulzapNPsfOYXrSlN9KrlkHxvmJQq0Y83ZZqmGE+Z54Z4
TMfA/rKULJ60BSaBIwVb/NyQC+6iKUpiLDMdWhVEqKXtes7fpHnrHwHYPtTMRc6S
EG+txT5gtUKKuddBTx01UZA67L40RyC6gwbg3zipOY5XFChLQ8O4hjSaDZuQEgik
ewZhJ4ExH17+f3GaNAgPE8OtafSF86jcqt2lowrXbPi1CV9BHSriD/2G1WQycjma
Z8tFSb40fxvWHTiuDs9AyymN3OKpVX/IBm8gFxeitryxSYZ4ZOtsI7fTH+hH74C9
ceKwu/3iPednzj1NvTpBdQvuEHRzgN/YfABXdJ71WiYRwVtE+hsdPJkKAdNYRXVs
YWJU5Ry5CAyaKZ//XsqqXM6PeACbRfWt6lqHSxnuaJuuh6dn2btJj4hX4Do81bh/
qwKUwzIVmJCtyEaO1VsslQk7CnLZ2dErIWGnRVkeLY0wMW4qspIZUx7ikgJZZ0sz
XFMwByTkVOtzz24nX3DdhxXfClqk1wrUI50erG6XbCjbuN/XEWbKZNLKUrfxkNys
1GjCXr1qagz6s8igxAHNmK2I9N8lyZSLrqKZf6m5CjyLMHSqj95jEeXL8pw29FLs
0BQnpfuqqTQfonYDSxdHgMPZfT2y+iX5LyIaozTJuQjsDAKXhFbytrSZ9as/kmtB
ne8L04gDLiHcX3K0anLSZN/0N06LiIa9O3qygfBHdtB3iIlBnE+yuGMZGPo2KXSe
4HnU/9E6Sayh3hHEqYtDnfVtNSQhJEwGr+HgX2+wvfQYJCUX/x+2gSN40/aYQ1jI
h9gynY7NK8WybupP8JsjJ8t0UOaFwfXC2har1kq/uChOGnsBI+E+Lkx9mOPkpQ0M
+NkN+HYuZC6dqUJAUZmHdzGPgh5MPZiIwusaW7frswmKko32y9VDfg==
-----END RSA PRIVATE KEY-----
`

	SelfSignedRSACertificate = `
-----BEGIN CERTIFICATE-----
MIIEIjCCAwqgAwIBAgIJAKfsBagyQeuuMA0GCSqGSIb3DQEBCwUAMIGFMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCQ08xHzAdBgNVBAcMFkFwYWNoZSBUcmFmZmljIENv
bnRyb2wxFDASBgNVBAoMC1RyYWZmaWMgT3BzMRUwEwYDVQQLDAxVbml0IFRlc3Rp
bmcxGzAZBgNVBAMMEiouaW52YWxpZDIuaW52YWxpZDAeFw0xOTAzMDYxNzA4NDRa
Fw0zOTAzMDExNzA4NDRaMIGFMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ08xHzAd
BgNVBAcMFkFwYWNoZSBUcmFmZmljIENvbnRyb2wxFDASBgNVBAoMC1RyYWZmaWMg
T3BzMRUwEwYDVQQLDAxVbml0IFRlc3RpbmcxGzAZBgNVBAMMEiouaW52YWxpZDIu
aW52YWxpZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMOXU/IY3uIA
Kb0QjlWC0eEEGQAVBTOWqXclb+Ov2fZbJ+DGfrRmmBRk4Sz3H5p/l79mVDd5+Npm
UchWSut9o/lYm9zbZ6HDAyE7O9vtO7uNGq+Odg3lUzh3QSZzh+vaMeYRMzodGK+7
6vkZ6jxGJRMazZO2LghuLH7g3ft0Dbd/QFGOQHkP/y2KeNt7jeuE4OT6jjA6OdSG
+q/gXAtFADQKkx814pEAyp1T4K59rteN2KoCbRiNrgxuewgkeu3VOGTagkwbyav1
o3uxE+wfc3r0TOdW+2xLUO9KlaG3pioT8Vh1YcwZvghe+mkueyod1VldKl4vgH0N
sUEMpi5zB68CAwEAAaOBkjCBjzAMBgNVHRMBAf8EAjAAMAsGA1UdDwQEAwIEMDAT
BgNVHSUEDDAKBggrBgEFBQcDATAdBgNVHQ4EFgQUJguFG8ibiS0+orhtIBPDUPAT
jL4wHwYDVR0jBBgwFoAUJguFG8ibiS0+orhtIBPDUPATjL4wHQYDVR0RBBYwFIIS
Ki5pbnZhbGlkMi5pbnZhbGlkMA0GCSqGSIb3DQEBCwUAA4IBAQB1th3ZUxczr5ZS
5Yem4KGLAAQck0o71cvmIMLXs4u0ioehM6UDxNrOwy5XaxLweZebIuUgZ44+KZU/
UerqSinWE3zsNpHbV+7j60aFjuCKAj+9WtJVVa/A2KLPVMYaoGITgdohHeYLKqWO
6io1VvLiahXa6G1nElPhRwdBU5qnJ651I821xf9TKeIUZM+bmYyqDE2FKHbZ+hBB
AGrT265gN3/aZKi5C5s2YFZ3Oc7dlZ1ZUhZqiLeTr8cmc9bZAdNY9ppDn3xs8ThI
GQlQzUI4o2FE3OdwA28PO78/DvnZ8oVeF/KD19oExJ9tj6dzLZIrkxaRDu7CD70S
41PQLy8Q
-----END CERTIFICATE-----
`
	SelfSignedRSAPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDDl1PyGN7iACm9
EI5VgtHhBBkAFQUzlql3JW/jr9n2Wyfgxn60ZpgUZOEs9x+af5e/ZlQ3efjaZlHI
VkrrfaP5WJvc22ehwwMhOzvb7Tu7jRqvjnYN5VM4d0Emc4fr2jHmETM6HRivu+r5
Geo8RiUTGs2Tti4Ibix+4N37dA23f0BRjkB5D/8tinjbe43rhODk+o4wOjnUhvqv
4FwLRQA0CpMfNeKRAMqdU+Cufa7XjdiqAm0Yja4MbnsIJHrt1Thk2oJMG8mr9aN7
sRPsH3N69EznVvtsS1DvSpWht6YqE/FYdWHMGb4IXvppLnsqHdVZXSpeL4B9DbFB
DKYucwevAgMBAAECggEBAIwTGVx9sUmbokizzaux58tFmv3zD+mVUcJxfkNK0kdb
myCgJ2fdPbcFVDpWtTx5elzp1RBx+uW2d4WJP1iNf1x4uA8g1oQD3H71I/ZqXOgB
swXdefCTttjulysJfGNNvYSt9sj8w4w/gZVqmNUXyz92Z5oM08TX2mf3dSK7R4OM
jk3hgGYXbgLg9ZMy3DuiZlkzIUY5sk/574j3pYaSaA+ufeIINzovod7zfqIT+5HY
nmoHKJqEjd7gZ2U0EVn3n39NXV9+eBsjuqC5aFB/oGYvH6rkXa+10vILImNoJEl2
PRMIWggvKyG+rFn714Y5QnPwLP/So71Hq00XSMTOCUECgYEA8MVfyRuA8/d0Jyas
NZv0lbtsN1XyZ5MRsNND4DkZR4Gz3Z5JT6mXwlzTLFMl6RKxPwkxLowYGeJQnEnK
5navmcvjr9EbMrAiebmCam1PWfUKxvLVWgw4sLrscD2YlB/ylm3CugLB6p2eS5HV
7bqGhrS56PtjibZbU0dBHIc7Ek8CgYEAz/Zjw8GBV4bh9dqpKB0Kv9VF05d8maLH
X6ybRSacJ24u/icYijq1vIlz5P7fM1enVkIKF3K8V2aToMjFShnAkdWxk1Ss6UrA
DsuRq1WAf9vSIsSDao8Cl4We7M9mbnGj3pEmxdN+Fh5l9Gs09cb8v75UacDrQ3AA
x/izdR9ePKECgYA5uCdZT9WfJuBajmPUSjndN1we++Srvc8M4+iutSGBSe+znkGW
7mIqPxSE2L4K3OdO5EY9EWqpsd5/SRVyIFvCc+V2fZuWkDEDm1xPZTsTprHZtdc9
yhZBD96kna+ZOvtt4ow/CXAxbW3IUgLGBeRAIM08tB9NoUDNexR54VVg8wKBgAhU
ZSCVIpof4/MWBAqAR0rxvu8/tRTuSTAS2NFHP5/wsN8rL7dxrI/VrZexgb0ruJ8i
3AaeaN7TU/xvrj9OksEEny04igh6HwE4tCf5r8DvTBZqap0dB0yMZJY1pHOuB5NF
mxj3ZKh8JTdKCAAELWF8vSLTQFkeJlncI0wAqqShAoGAR/iDslDHU0M/CjvqN+le
m91MxJsaFi3FZ2Pv60ghqG27X8uPbgvGXkTmmcqgAaR/Fh8b1UH+yqOnPUiremxO
oIuOOmvkBAyy0cHT/BSunqxiFOWcu9yb+2MfSWaHHxNj3kteif8Oai1lL84JOucC
/FZTue8KxnrjBcbGYYtwPCU=
-----END PRIVATE KEY-----
`

	SelfSignedNOSKIAKIRSACertificate = `
-----BEGIN CERTIFICATE-----
MIID4DCCAsigAwIBAgIJAI84W5tsR3z4MA0GCSqGSIb3DQEBCwUAMIGFMQswCQYD
VQQGEwJVUzELMAkGA1UECAwCQ08xHzAdBgNVBAcMFkFwYWNoZSBUcmFmZmljIENv
bnRyb2wxFDASBgNVBAoMC1RyYWZmaWMgT3BzMRUwEwYDVQQLDAxVbml0IFRlc3Rp
bmcxGzAZBgNVBAMMEiouaW52YWxpZDIuaW52YWxpZDAeFw0xOTAzMDYxNzExMTZa
Fw0zOTAzMDExNzExMTZaMIGFMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ08xHzAd
BgNVBAcMFkFwYWNoZSBUcmFmZmljIENvbnRyb2wxFDASBgNVBAoMC1RyYWZmaWMg
T3BzMRUwEwYDVQQLDAxVbml0IFRlc3RpbmcxGzAZBgNVBAMMEiouaW52YWxpZDIu
aW52YWxpZDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL2yyphDlkc0
AmpxPAp0GZkPiGM5PgkahtafnIMCWXC5JkWz1Qk7dSgKUNl6hXg3I1dBJ/Kej461
XwRSPNVFp1ehUsl7+D6QJy7+WEuk4Ea3s6/btMByH6+7S5EvJJGivavXlaLfqXnD
2vRKW+6EuabNpeNVplYrZQYQ2aHYIsY+HYInzqdpb/qx6BSbuvuV74N6KcI7dbpC
551y87hq2U/Mftka1TmfT/ll7bKdZ4q192ME14/USGqmyy8HHAiFvsgRjmqUUe3N
6+O0VOjjTDByL1naMw5MlSIhcIOM4O+rklAiKVzL9ONc1rz1U38JrNeUjD6pEleD
cPPBR1ge8DcCAwEAAaNRME8wDAYDVR0TAQH/BAIwADALBgNVHQ8EBAMCBDAwEwYD
VR0lBAwwCgYIKwYBBQUHAwEwHQYDVR0RBBYwFIISKi5pbnZhbGlkMi5pbnZhbGlk
MA0GCSqGSIb3DQEBCwUAA4IBAQBtkA0rtyLblXFDOntBY68gpcgaa1KjUo9YNfku
VywCVan7ZoxiJkCsirIonugJDan1Q+meMsA2tO9UXoxShXSyh6pEcqEa4uSelj1M
IysFFdtCwcUOzL3OwAY1TecC3rRadJD/hWM1F3iHPpG3BGB1LDkEZVCQdi+Gtq/C
UCHb6g1lluLt/pxHuicS4etEd1zTb1oh2CbcijMLFSGimPqkqmSBOkgO6/283xA3
MOSm/w7yrLylq/N7/g5c6kLZSWkYpoxA5PgImJEsOwoFF9VrMzC4JPLdR3zjLq/T
wBmyIzcVFJ2g4V0iS5Ufi6uHwc7A03PVA50g7p7JUVRx5q/B
-----END CERTIFICATE-----
`
	SelfSignedNOSKIAKIRSAPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC9ssqYQ5ZHNAJq
cTwKdBmZD4hjOT4JGobWn5yDAllwuSZFs9UJO3UoClDZeoV4NyNXQSfyno+OtV8E
UjzVRadXoVLJe/g+kCcu/lhLpOBGt7Ov27TAch+vu0uRLySRor2r15Wi36l5w9r0
SlvuhLmmzaXjVaZWK2UGENmh2CLGPh2CJ86naW/6segUm7r7le+DeinCO3W6Qued
cvO4atlPzH7ZGtU5n0/5Ze2ynWeKtfdjBNeP1EhqpssvBxwIhb7IEY5qlFHtzevj
tFTo40wwci9Z2jMOTJUiIXCDjODvq5JQIilcy/TjXNa89VN/CazXlIw+qRJXg3Dz
wUdYHvA3AgMBAAECggEBALZkgyEV0xdRLYV0rJsMeFRPt/5XWotcQwt3WgApMSAO
FXttZkdDMOk3yfbhNBWbRlKt5iAETtmTD/HStIUHPNgn17a8iLp21gX8LZ9FvxKf
rJhI8ikbUdYgio7kug+BX0cruMdqr8PRCeRa3rueR/bWwkqr3ov5m1/Ssb5IV18t
UFj7c2nV9JMB0Wq30Awa2RcAFwC2zBkNcyZ72bR6x2DxV+QHG+xDzLO5y/PTuLEF
GBc3a9ExByC8RtUlC1n2YCwe6HSZBPALoaCIloo2isEpLpAy2ihxxaJ+tJh6iScB
t9RAvffPNN6AgRIM5qfWaoGR8E60yFBPrkHPTgCiICkCgYEA30yqVmU3zxbTvNQH
gmRCJrUW++ReLsa7xnaZmN4Pn72tGzdpmhajUWgw7LXdZ6TtHQ3UzZzt4pU5NoXa
BgE+4XKewW90GoCW66yQpg80OD0UUl5JNT/GM2cVYqROEioFg+7/rEzhF1tv2JAe
+I2uxMJoMlA82giwLJ2qDqctDZsCgYEA2XpyRHVk3GZnQ/q9tK1VyVfUkFJh0FnQ
q7uzWysi3UCOzuJYauM300LhkP5gWsB0J3osBlreNMg9MTOIy1xi7tNXdOq07Hhr
4KsP8TkDhCmQDepx1V++qVBDTlAofCO5V530Ut1+wWDBlcycAz5MIy6j5TeAWrjm
s0sisOrk35UCgYEAlQyWcn6zht6kzNj2fjmv0ih0RATGPRDYS/vkQJ6Q7T5taspN
CdZsgy054vbt53216/vMfMZwHxseCl/EDNgOAexBPbrIU9xbYMpZ7w4c/CEBDI30
7b847By1sJcdqZA1CECiln7mjHGzMWnZ0my4KIvfgx390EeWWOGQnqFGOFsCgYBs
OEZMUq6SDlMsvMVR9z9NJeuctaH+7/KqwoiJwXlj6BAoWvHsnozVD973K93+yu4C
BwWJVAZm9Y2dwis8JwkEFx7aC0FkurfT4MvaGajqR1Rr2FI0/6P81PfpLHI48/3y
36MI6Td+OwuZ42tvIbz5dOgR1ACHJKOIbMciioDB0QKBgCtqIx4aMteQEOvXlx8E
mmD5qNBhcWUgxexqdqrvyqLScGHuhuUHny4ytTzNMulG3S8DEC2ULrq64OZuIdyT
XkhQHb96BsPSpUioe5MLzr5EnDLLV8Hptf4mmJctkd0eczgvtlX9QyFp9LpeqAMQ
FA21NyLKatm942vrWMsBGrcS
-----END PRIVATE KEY-----
`

	SelfSignedECCCertificate = `
-----BEGIN CERTIFICATE-----
MIICiTCCAi+gAwIBAgIJAKVi77tebfBPMAoGCCqGSM49BAMCMIGFMQswCQYDVQQG
EwJVUzELMAkGA1UECAwCQ08xHzAdBgNVBAcMFkFwYWNoZSBUcmFmZmljIENvbnRy
b2wxFDASBgNVBAoMC1RyYWZmaWMgT3BzMRUwEwYDVQQLDAxVbml0IFRlc3Rpbmcx
GzAZBgNVBAMMEiouaW52YWxpZDIuaW52YWxpZDAeFw0xOTAzMDYxNzM3MjNaFw0z
OTAzMDExNzM3MjNaMIGFMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ08xHzAdBgNV
BAcMFkFwYWNoZSBUcmFmZmljIENvbnRyb2wxFDASBgNVBAoMC1RyYWZmaWMgT3Bz
MRUwEwYDVQQLDAxVbml0IFRlc3RpbmcxGzAZBgNVBAMMEiouaW52YWxpZDIuaW52
YWxpZDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABIOpnfYOa2+Fxhxl996A5cvU
OUrpY1/oy6uUdBmgo0Rfal9lv5NiMRhgg7s++jr2yf8ls2zLIosNPShIVLJl9Sej
gYUwgYIwDAYDVR0TAQH/BAIwADATBgNVHSUEDDAKBggrBgEFBQcDATAdBgNVHQ4E
FgQUcjdTM122v7VK7gWXexBLg+1fCu0wHwYDVR0jBBgwFoAUcjdTM122v7VK7gWX
exBLg+1fCu0wHQYDVR0RBBYwFIISKi5pbnZhbGlkMi5pbnZhbGlkMAoGCCqGSM49
BAMCA0gAMEUCIDybXtWSTIeyETgLBGitveWAYh8dmBK/Wt3+Jdj/5GF6AiEA95I8
96/cXDmduClIg4CnYLGbkDPsC+3nENEIr92WqgE=
-----END CERTIFICATE-----
`
	SelfSignedECCPrivateKey = `
-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPrjdSmSp6D/M6KBOwwz7u/NzO70nBT0U74QSCBWmwAOoAoGCCqGSM49
AwEHoUQDQgAEg6md9g5rb4XGHGX33oDly9Q5SuljX+jLq5R0GaCjRF9qX2W/k2Ix
GGCDuz76OvbJ/yWzbMsiiw09KEhUsmX1Jw==
-----END EC PRIVATE KEY-----
`

	CASignedRSACertificateChain = `
-----BEGIN CERTIFICATE-----
MIIE+zCCA+OgAwIBAgIFFVGJEGcwDQYJKoZIhvcNAQELBQAwga0xCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMGRGVudmVyMRQwEgYDVQQK
EwtUcmFmZmljIE9wczEUMBIGA1UECxMLVHJhZmZpYyBPcHMxJDAiBgNVBAMTG1Ry
YWZmaWMgT3BzIEludGVybWVkaWF0ZSBDQTEoMCYGCSqGSIb3DQEJARYZbm8tcmVw
bHlAaW52YWxpZDIuaW52YWxpZDAeFw0xOTAzMDYxNjUxMTNaFw0zOTAzMDExNjUx
MTNaMGkxCzAJBgNVBAYTAlVTMREwDwYDVQQIDAhDb2xvcmFkbzEPMA0GA1UEBwwG
RGVudmVyMRQwEgYDVQQKDAtUcmFmZmljIE9wczEgMB4GA1UEAwwXKi50ZXN0Lmlu
dmFsaWQyLmludmFsaWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDE
oCEYaEl1mrEOJy6jO8+QXXmlb8mYAWQL+lSOPYRvzuPWOvnfvW/BkDYpvzosfDZg
HNS106ROtWFe+kQZTTrm+xPTWMl/JIyn5ka1ZnJMAey0Oxyn2PgSmly++6esOiKH
CZm9Zhx2p7/d/fkYhu6G0X08SKVwZdtkerfKHqPEuEEcWm1H+IbMVweKDjiAqrEe
pcWKr/+i8NmOERCYvPFL6/Kf8/U4MrDM/3vh70DFKWZ3x5xWG7TWl3vCS4I+y/m7
xEEQgF2vofcQO0CD09s4iRKCqbVM2vofWZVhYj63BgfbIXHYCFPyIpaTwc3u4fRW
pnWIUp9BFUy8PqXjh9vfAgMBAAGjggFjMIIBXzAJBgNVHRMEAjAAMBEGCWCGSAGG
+EIBAQQEAwIGQDAdBgNVHQ4EFgQU8RrF6qwRsTAtsdHUDUe//WnjNskwgdYGA1Ud
IwSBzjCBy4AUwVR9U/f9JDLi179Jc1qudzjYjhehgaukgagwgaUxCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMGRGVudmVyMRQwEgYDVQQK
EwtUcmFmZmljIE9wczEUMBIGA1UECxMLVHJhZmZpYyBPcHMxHDAaBgNVBAMTE1Ry
YWZmaWMgT3BzIFJvb3QgQ0ExKDAmBgkqhkiG9w0BCQEWGW5vLXJlcGx5QGludmFs
aWQyLmludmFsaWSCBRVRiRBmMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggr
BgEFBQcDATAiBgNVHREEGzAZghcqLnRlc3QuaW52YWxpZDIuaW52YWxpZDANBgkq
hkiG9w0BAQsFAAOCAQEAbcOigHR11zFrPXHS7ZGG6Dro708QVrECofZJkVAKOUbZ
kWP2h4T+xSS6WEqk/1hKKo2RhiL7UlGwdtqNZ6lk9cwQORmeuIXEgy2BsbR5H4vy
L/2LhZzCSDYfeSbByhFdvLCNJMSGK+o49NUESD1zEJFPDDtMj/NLAgkLCSGi/8iF
VVJpJT3zP1I/C9Zmt8eFf6sjEdoIg6PPReJzu1hcEJLz/jqLuzoGasxea/k2k06k
Rj3GV8f3BZryN3+TKa7IiNsT7SGiZMiPWUvML/dp8aTgCNA/GlhGqzQPnDXBvgTN
JfDXiQDnR1svJZNHEfQDPIEeY6ECBAEaU0rgH9/iZg==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEOTCCAyGgAwIBAgIFFVGJEGYwDQYJKoZIhvcNAQELBQAwgaUxCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMGRGVudmVyMRQwEgYDVQQK
EwtUcmFmZmljIE9wczEUMBIGA1UECxMLVHJhZmZpYyBPcHMxHDAaBgNVBAMTE1Ry
YWZmaWMgT3BzIFJvb3QgQ0ExKDAmBgkqhkiG9w0BCQEWGW5vLXJlcGx5QGludmFs
aWQyLmludmFsaWQwHhcNMTkwMzA2MTY1MTEwWhcNMzkwMzAxMTY1MTEwWjCBrTEL
MAkGA1UEBhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIx
FDASBgNVBAoTC1RyYWZmaWMgT3BzMRQwEgYDVQQLEwtUcmFmZmljIE9wczEkMCIG
A1UEAxMbVHJhZmZpYyBPcHMgSW50ZXJtZWRpYXRlIENBMSgwJgYJKoZIhvcNAQkB
Fhluby1yZXBseUBpbnZhbGlkMi5pbnZhbGlkMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAsu3xVwphWUz/josArTjldncVgAZHVa+NVKgJpohoyJra7IgX
IpNY9dlMd4wydPc2/Y2yc6lcZ/XlfzWQS5hYsOqMoWsylC38ZeX9yminQ+Q0TFsz
EMK3v9aFV9kiFIKPh15wF1Iau1HMypN6e61rE8ThEixLU3gGLHhMvDuD+JIoIJzj
b9DvF0GU8kJr/ur4oJQDSnGiW4cg+IGsySH4y9e9NMZmwpnNprXVkP6Rr94UKAfB
KJDb683VeuuyyknwC6ZNmUYuSLIUxXi5UPbKQi20zwJkHbXui8aBMPvOIUKXSiCd
y5apfzZw22npUUbkarg6d26BU2hVjOd+6VyrqQIDAQABo2YwZDAdBgNVHQ4EFgQU
wVR9U/f9JDLi179Jc1qudzjYjhcwHwYDVR0jBBgwFoAUwYqnYixb6XZqbD7dpFMx
280Z+5YwEgYDVR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZI
hvcNAQELBQADggEBADiiT8vhjan43mJDv0kFS7xipePJeXbLGgAu3PocoR9qUbh+
7Q2Zj5p1rCDaqUoUB8C5p1T3JjPxeGLQ3A9qjWP1FUhbsHTfJzcU989Q73Iomx3k
cmmQhWt+xvCpxGFmvjBpEyyIhCkgoGBUfHQFjcyrLWbaLRXOR1eF9B0YORbJmf10
1lxg8PQ3LclCs9m4XSI2THiURDN4dFANRYqcWvg0ncej/SOXJQ+OZ5K1AJ9z/acB
GT/EpQ6YfSyxOsKVPAu3mrIxfOAqahKFJa4bgV+a5SmiENYEKYaWObAgMpO+Ud5F
kG5IVBDJicJm/kazyAs2f4NmW+8p8CKEdrB/6gI=
-----END CERTIFICATE-----
`
	CASignedRSACertificateChainPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDEoCEYaEl1mrEO
Jy6jO8+QXXmlb8mYAWQL+lSOPYRvzuPWOvnfvW/BkDYpvzosfDZgHNS106ROtWFe
+kQZTTrm+xPTWMl/JIyn5ka1ZnJMAey0Oxyn2PgSmly++6esOiKHCZm9Zhx2p7/d
/fkYhu6G0X08SKVwZdtkerfKHqPEuEEcWm1H+IbMVweKDjiAqrEepcWKr/+i8NmO
ERCYvPFL6/Kf8/U4MrDM/3vh70DFKWZ3x5xWG7TWl3vCS4I+y/m7xEEQgF2vofcQ
O0CD09s4iRKCqbVM2vofWZVhYj63BgfbIXHYCFPyIpaTwc3u4fRWpnWIUp9BFUy8
PqXjh9vfAgMBAAECggEARD43flsjs9eewATFYQ4vOjHXOJ4V39YLvUSC+GNNhejO
ltodQ5RiJ1JAGjkunaX20WDwSrNMAa1eQDKoVAfD+8sE1IOqW6B52QRJYkhOPycj
2mHxzie14e8FZZu+VD5RIYEphNzd4CjUpN2zCNo8CzrGNpgYI2yWuscE5ve/a1TT
/diQVMpnDQkOFbqAmGS70hAjcSJ1yW5nUSSXmjjlEkxOmWt4wqLUYTmrqcW5DKxv
jmxoRb1XdWzT+ryXOVX6ilvjlQVMi4k+ohvV6svFjkRiTdii16/YK9CbRaZvq8eU
+F53LF68Cv4YlX/5C4pcwyEf+I/UolDfyNE4NrY4kQKBgQDsvYpenMs80H14wd3r
O7QpaIvZe6u++nIa9Pa0n4gwsieBEHAhrvJvqln6h2mOuWL8I6cSERBVqqUD3F5w
HdNa8mStpIF+XlBmQW5UGZTqUk6KugJzIB6ZBXKOfW6g/+xaQmC0tFTzuvsV0pZ0
X8EmsW9xbBe/+Y/zOgLXpsfgtQKBgQDUnyNz7Gaq2SldGEVMrI8eicfV8KcH0rK7
ocUnESxL9zvNFkz9wz1t5/D4L/D3XuLAjXtKV1lUvbzVxHh46GQxvz0/5swYjL71
rhEDqQF2t7PISxXFq9yMAg5iOxyLwBFGDcNsL6G5SdXpM/YWGo7wSwsWdYEzZULs
UZ/ZOQIqwwKBgQCfFh/NxH+utkwawexnDw/aY67Wzwxyocnb45GFf079qjpxuKIh
gHbaIxekCyscBehGl47FzUG0z59kIMo1fVVyYEDXjxyV1rsgfAev7CDt9bFh9+19
f7AQFGEO76tP9arWXJSv2h7cSmJAH+uK+G3LmqDRD1pGX2YkhG80i5b1oQKBgQCo
GMypqJuetSObo0WekcpwxVNFVAZqC+0cpI+/DDeuM1+HC/uAoKvfSYFcZmKm39B7
lR+FLbvFYGB7zOHGDUyxe9VLwQdY3WVXzO9MqoAqwJ+VWa9z4STzV+jRRpSR9B5z
+Quoa5v7ZmGFBnynCwY4+cthTTMBVCxtszaiQQzyiwKBgQDnJLZONi3IcsgGOuka
4qYVQl7sSIuBJ1gyZr44WjoUlykxLij59MEnAOolR86PWif6VgNfsJL4MqTkaX++
PgH8k1z8wkIxUP94BXLfSHJyUriHMPX45T3Rr6uBrvtVujRNIHndGKDoa7P16Hzz
otCrxsW0tXSVJ1KfW1Cwb7jtHA==
-----END PRIVATE KEY-----
`
	CASignedRSARootCA = `
-----BEGIN CERTIFICATE-----
MIIENDCCAxygAwIBAgIJAKtfkOE64RynMA0GCSqGSIb3DQEBCwUAMIGlMQswCQYD
VQQGEwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZlcjEUMBIG
A1UEChMLVHJhZmZpYyBPcHMxFDASBgNVBAsTC1RyYWZmaWMgT3BzMRwwGgYDVQQD
ExNUcmFmZmljIE9wcyBSb290IENBMSgwJgYJKoZIhvcNAQkBFhluby1yZXBseUBp
bnZhbGlkMi5pbnZhbGlkMCAXDTE5MDMwNjE2NTEwN1oYDzIwNTkwMjI0MTY1MTA3
WjCBpTELMAkGA1UEBhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZE
ZW52ZXIxFDASBgNVBAoTC1RyYWZmaWMgT3BzMRQwEgYDVQQLEwtUcmFmZmljIE9w
czEcMBoGA1UEAxMTVHJhZmZpYyBPcHMgUm9vdCBDQTEoMCYGCSqGSIb3DQEJARYZ
bm8tcmVwbHlAaW52YWxpZDIuaW52YWxpZDCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAMi3d1btekZha72UWkdJvsuZenp8zNiQW7OaLLyZvLbvCncm6sNG
nZTVTkq5LZnrq2sOc04JIf7gFJc8vnJdBPEZEP5GRbWXmel4Q4qnoOcCr4/6UZ5I
1gG1JkGMsMWq83Zv/jqtjBSpIO9fQFEiv6h9jUAkzln7lfvpM892qCKZCFLh89TZ
5flSmwmNVOW5IFvUu9uFdKDBzOYaXpBXXhua3ewh9YqUtCfwBH/zz8bb3Shm9p55
RhT0Rp4yQpFB4FwG0MxNNqjZ8BLuF9pleTcRlTrVwTLQ2sOq7HndvyyahxtSP3BF
plFrnnzmeOZjtr7RRPvG/22fHxb6lGxnXEkCAwEAAaNjMGEwHQYDVR0OBBYEFMGK
p2IsW+l2amw+3aRTMdvNGfuWMB8GA1UdIwQYMBaAFMGKp2IsW+l2amw+3aRTMdvN
GfuWMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgKEMA0GCSqGSIb3DQEB
CwUAA4IBAQCRtSa8H5mwNKT1SgSjEatLMLjzYk1jpfIns6hR2X5f+D/kZkpgPX/d
4dAke1DnMIyxE43t1RJJzoi1BhFg9YvEyNU5k5L7ECuYmgU5tbmeNx2cDu4Huy5k
6ZmQ4Amh01X3JY3MFARJyRi/+CvmT+SbwdKYs8fwhPbGcLRhlUC0tvFe6JfvceH/
18TRlcDPC51+f+WidVkWC8Er7Mydv+bmNYqEgDEJ0UQWtI/j/hxYj1wuqvUaoN2A
GcDuf0HznG/JMOtoGmMELyxJlIHCBPG1h+myLON9/XuPHlhMQcIhmqNEBBGwSdwm
1Ybc3W8sUkBAEha1Nm4SGJOz2Lfk/rAZ
-----END CERTIFICATE-----
`

	CASignedNoSkiAkiRSACertificateChain = `
-----BEGIN CERTIFICATE-----
MIID/zCCAuegAwIBAgIFFVGJhGIwDQYJKoZIhvcNAQELBQAwga0xCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMGRGVudmVyMRQwEgYDVQQK
EwtUcmFmZmljIE9wczEUMBIGA1UECxMLVHJhZmZpYyBPcHMxJDAiBgNVBAMTG1Ry
YWZmaWMgT3BzIEludGVybWVkaWF0ZSBDQTEoMCYGCSqGSIb3DQEJARYZbm8tcmVw
bHlAaW52YWxpZDIuaW52YWxpZDAeFw0xOTAzMDYxODU0MjVaFw0zOTAzMDExODU0
MjVaMGkxCzAJBgNVBAYTAlVTMREwDwYDVQQIDAhDb2xvcmFkbzEPMA0GA1UEBwwG
RGVudmVyMRQwEgYDVQQKDAtUcmFmZmljIE9wczEgMB4GA1UEAwwXKi50ZXN0Lmlu
dmFsaWQyLmludmFsaWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDS
Mski1L4DVQfYdwXBxBHAWD3ZdpoS+MZa3lJXPMecpNOQfZmFzI16l3oKaePK0t2a
2BmxPmueiQ4UD4n2x3Z4bV/sjz6W13ATZs+N3EZopkr5kjr5T4NxiHaWA2JCFBFF
aMhJwoX4g94pUuHKHm5iMr7mr6nJLKwXW1RrZ3hl8q3vBoo0DzyZ1lIo+s1MqMzG
G86Wdw49lNwgi4eCMBPnX+UOENij0Kl80KfoCsp9FOxbAtZIpP5cgmb8EQnmuaoV
wbetaz2sZ+S/qasSr1nlD9wuZ74dCTxnPzt+hguyhiSLra0sd48ie9QdCqSHyYAc
rmNtELNY3ytbEN2KG0xnAgMBAAGjaTBnMAkGA1UdEwQCMAAwEQYJYIZIAYb4QgEB
BAQDAgZAMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAiBgNV
HREEGzAZghcqLnRlc3QuaW52YWxpZDIuaW52YWxpZDANBgkqhkiG9w0BAQsFAAOC
AQEAHkjRjqQSsLCJUQE6QNBcx5I+hle3InFiyvr5u2DAuL95GzA4LQbk1voYFj2h
fCgAE9ON2aViPjlqQcPYl+2Koyk0vkQwkSb8oJsEucAcAGk65BmlBVDbM9IgVC96
ILYuk0kT+Y0CgvEWBftodYs8xKEAkNuopUTKZH1iPLIFghI3cxuU56G6A4XUFEbZ
Af5qM2UBxFnKPPm7HyKs7WPoud9/V4d7UaaspXwFdvA7f8bGJvgjHGYE8XBZONgQ
YJWLBgCfIKF0LqXAFwVrbbzJEshRT+WXCNvCCq5+AJK4gBaovzk7QkAREkwTe7Du
++NBiHeG0TrPN6tX7Ke7EUYcHg==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIID+TCCAuGgAwIBAgIFFVGJhGAwDQYJKoZIhvcNAQELBQAwgaUxCzAJBgNVBAYT
AlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMGRGVudmVyMRQwEgYDVQQK
EwtUcmFmZmljIE9wczEUMBIGA1UECxMLVHJhZmZpYyBPcHMxHDAaBgNVBAMTE1Ry
YWZmaWMgT3BzIFJvb3QgQ0ExKDAmBgkqhkiG9w0BCQEWGW5vLXJlcGx5QGludmFs
aWQyLmludmFsaWQwHhcNMTkwMzA2MTg1NDIyWhcNMzkwMzAxMTg1NDIyWjCBrTEL
MAkGA1UEBhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIx
FDASBgNVBAoTC1RyYWZmaWMgT3BzMRQwEgYDVQQLEwtUcmFmZmljIE9wczEkMCIG
A1UEAxMbVHJhZmZpYyBPcHMgSW50ZXJtZWRpYXRlIENBMSgwJgYJKoZIhvcNAQkB
Fhluby1yZXBseUBpbnZhbGlkMi5pbnZhbGlkMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAtkDh9gusfqgvuTN/EYfiXS/JbnX34R9Kjau3JRqA+CttbY7i
5pYbb+xdVBiUSEO2hikIF9LbuCrEQUdwPbBdX5gratDcOquqwYj47uh9ii1UEsH+
7Sa4aV2Tr1zOFLrdkuKdJIUk1Gwp6yfoLAsXDfmScIY61EMfEcau9zH/UlCNMbfT
UMYIkJJK8iPyFlS/EehhFuPd+0Mjd/qXPCBBD0bCcJLVka/mAjKyXZPJ6Vyh4/UD
5Ne+bVp7VatiWaiAf7wCF0pOr81yTdEYV7Sp53fXMUEuwACClV88YlET3eGwfSXd
j7ugKBnyDX3TlgOZix88qRhrJro+Zr9i4RMl8wIDAQABoyYwJDASBgNVHRMBAf8E
CDAGAQH/AgEAMA4GA1UdDwEB/wQEAwIChDANBgkqhkiG9w0BAQsFAAOCAQEAWIM+
R69UBhumUvzLsw2xXTPmHElzuM0w5BQMfJcvZ3gsdo0wvCXkKQPSu2yPBH0BuqRt
l29Ng3+ciuvRqk+wnDp8BH651U6U7oqSw2V7alOOy4ugPdQm7NNbJuOSYCCAngEh
ohEloW/4CGLyyY4ffNfoPySgl0xXL8B/QVuTuw/0IPK448ziUJ346sZf1oBQcM2y
+gHfbNp+jQM9lQeQhjv2TKFZKrfy6uS5hHQAlpo17e0oPZA+Yw6+hQ1n91CGZ/Xx
rWm+oBTAyZUuL0NgHBCwGzZMpqKk2hyevvOC3yWQsxhbh3Q92ucxx3ifQv+WsGu5
e8Rwmr3/doh6voOnGQ==
-----END CERTIFICATE-----
`
	CASignedNoSkiAkiRSAPrivateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDSMski1L4DVQfY
dwXBxBHAWD3ZdpoS+MZa3lJXPMecpNOQfZmFzI16l3oKaePK0t2a2BmxPmueiQ4U
D4n2x3Z4bV/sjz6W13ATZs+N3EZopkr5kjr5T4NxiHaWA2JCFBFFaMhJwoX4g94p
UuHKHm5iMr7mr6nJLKwXW1RrZ3hl8q3vBoo0DzyZ1lIo+s1MqMzGG86Wdw49lNwg
i4eCMBPnX+UOENij0Kl80KfoCsp9FOxbAtZIpP5cgmb8EQnmuaoVwbetaz2sZ+S/
qasSr1nlD9wuZ74dCTxnPzt+hguyhiSLra0sd48ie9QdCqSHyYAcrmNtELNY3ytb
EN2KG0xnAgMBAAECggEBAKyW3IXH7nin8bgwCj8OMZEgIzCSbHHVaHCmCS/uDOw2
fiwupMayrRwSkjdIuKwJtcF1XKsm2JCkcjXQiHRjVIgPLmr7NuX94N1dVmBhlEJL
AFapVdjtC71F0jDceGpPNdsq7QF7QitKgzilABXIJNRmXE7nv14aWvcWm1tQ6w+w
1RkAEqmPcmEzGeOxWWMcaUFamEPI/SDGIDSYyeVW3qH2afqgm2B3WMU0A9hM2iDV
rTXYQpQCaTulSSCNYXCou+cOlvN+Jvf21OjxoZ0/UjGqeM42aCdyUU9K0Jc56Ps6
RdW8rObTAYFQJj1wp9NB41QRh02HIBp7yhyyazJoXIkCgYEA8VbRqwNN7g3QmCp9
gDxm/CQ+1tDFfbSteRnvcPwpQoLY8L+XOCws8VPP3IfKSwVU21CS+ji+0uk+x8nj
83iyFHK8K/ToAeijllftACcBgJDQOfRq9LrXwHBCIvTq4pwNoPiSVwtUotaZ5AQX
DMTOB9T6hnMu+YpaOtq84O2xHyMCgYEA3veuizIYnpi1QnTDB5zn08EWLWQYMjrq
K4SUqgL3Y81dls34bPHHAPqcK62pMRnK7qbgdJ7wjJiODC3ZZ7bjZHJFNFO7RzEA
ktFNCZ344CuX7vNR9NywathFA46mu/eCKWJLI/NrF1IclGG08ju+y8lrFm0UdxOe
9G228h30s+0CgYBHOCOvn84Djjgcb42RpkGN7vRMWFevfP4kWq76XK+gXRTAFwn9
Haw1m1If9kKQWQZtoh19kfleLE7GjqGiW9/RgPpezmsZBRohZ9kczmX3FsUcFTDq
/6hjtb0Oq9AVB5BODIzC+ykC1OmdDEfxELLsRMGZo6wdH+L4s0xB5GL8mQKBgGMs
S6iCKc0xIz5h7PWP5tWbBqA96z08Uzf0CqPsGdl8WOpgxuS+TcOztI8A+UZrsIWi
GCgHIfuHR3dHVXH6OP5OjVWPALfTpeunyNpEN5SOD1ArTgLZvmZnt5qzcpocpvp9
S+q7tKB0111wcClmRaEi/8zDy9yDD6qsujjK9jKpAoGAWQCY3cGmJD5tbt2FtHbb
hcAmRHPHGNps2rfW6JzQsLcb47Tp5iqt8aZYTtehnQyb7wyBdP68MGloenI6x3dQ
UjVsFHOBQq/cUVgoIlKcm9kmByjmGKcgWtPHi0lFT8u8kLvxf2NUiZyKDis+P2xq
0ojmGWKT22C4HUSsMTfq2Ew=
-----END PRIVATE KEY-----
`
	CASignedNoSkiAkiRSARootCA = `
-----BEGIN CERTIFICATE-----
MIID9DCCAtygAwIBAgIJAJSR9zIym8jSMA0GCSqGSIb3DQEBCwUAMIGlMQswCQYD
VQQGEwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZlcjEUMBIG
A1UEChMLVHJhZmZpYyBPcHMxFDASBgNVBAsTC1RyYWZmaWMgT3BzMRwwGgYDVQQD
ExNUcmFmZmljIE9wcyBSb290IENBMSgwJgYJKoZIhvcNAQkBFhluby1yZXBseUBp
bnZhbGlkMi5pbnZhbGlkMCAXDTE5MDMwNjE4NTQyMloYDzIwNTkwMjI0MTg1NDIy
WjCBpTELMAkGA1UEBhMCVVMxETAPBgNVBAgTCENvbG9yYWRvMQ8wDQYDVQQHEwZE
ZW52ZXIxFDASBgNVBAoTC1RyYWZmaWMgT3BzMRQwEgYDVQQLEwtUcmFmZmljIE9w
czEcMBoGA1UEAxMTVHJhZmZpYyBPcHMgUm9vdCBDQTEoMCYGCSqGSIb3DQEJARYZ
bm8tcmVwbHlAaW52YWxpZDIuaW52YWxpZDCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAKSMntoOnNIqxYKDAM7T3hLleGV/DXTkUbFdg3YuJFfYCmIg5MUP
XjmpdzXmvZAxGlKUb5Np7lZeY2xRt7DqK3vVbOmFLAryAgyglD5eTq8VESegUogt
pdI1ROhVYV3/aS1/jkfJWZjGfRim8NZ5RhisHRd0++oDmRUh0H0+e7zRbCNEPvDY
Vmw3JTaPwJz6ULgbOqHGofNwQFUJN4E0TBzt/qzC7r11aaQQag/j+TlPi0UTq46A
7pbm9kFD4kdEOQ3M6jAowyEAfei0Xq8jeHcf6eMeD+tki0AenXdrx7yTZqX5fwgc
wHUiwUZ2vZklescONvbbuW/WePA2uYKfYoECAwEAAaMjMCEwDwYDVR0TAQH/BAUw
AwEB/zAOBgNVHQ8BAf8EBAMCAoQwDQYJKoZIhvcNAQELBQADggEBAAO2gylryjq2
Ou+mm4E8On/Ecchci1FUlc8ZLAzxxTdtJDZQd4b6uEEOFnXyIDy1DU3CND3ieIa0
EL12/dsGzzhI0mJEf+oes4aaa4l0znMA7WQ58CUdBWuGoGvUzLkE3TAFqLwKC79h
DARQQp0VaM01pZQd/Sh/3+V3ijEg3818VtTvKHm5t26ggIRhSSGMhEeqCoG+9EKL
igzzaBprKzbSuGpj9QKLO4P4vIvgt2RutfqYBNR6sdzKmQYyPCYshfsoQVWo9mmD
/LLpNoSwwXGgBhI7cprp0CsN5f5/mL2dDSQqEAOl8YitVa4nL44iFrL8d4d8OM7z
jI///X8MUJs=
-----END CERTIFICATE-----
`

	CASignedECCCertificateChain = ``
	CASignedECCCertificateChainPrivateKey = ``
	CASignedECCRootCA = ``

	CASignedNoSkiAkiECCCertificateChain = ``
	CASignedNoSkiAkiECCCertificateChainPrivateKey = ``
	CASignedNoSkiAkiECCRootCA = ``
)

func TestDecodePrivateKeyPKCS8RSA2048(t *testing.T) {

	privateKey, cleanPemPrivateKey, err := decodeRSAPrivateKey(PrivateKeyPKCS8RSA2048)

	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	pBlock, remain := pem.Decode([]byte(cleanPemPrivateKey))

	if pBlock == nil {
		t.Error("can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if privateKey == nil {
		t.Error("RSA private key is nil. expect: not nil")
		t.FailNow()
	}
}

func TestDecodePrivateKeyPKCS1RSA2048(t *testing.T) {

	privateKey, cleanPemPrivateKey, err := decodeRSAPrivateKey(PrivateKeyPKCS1RSA2048)

	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	pBlock, remain := pem.Decode([]byte(cleanPemPrivateKey))

	if pBlock == nil {
		t.Error("can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if privateKey == nil {
		t.Error("RSA private key is nil")
		t.FailNow()
	}
}

func TestDecodePrivateKeyBadData(t *testing.T) {

	// Expected to fail.
	privateKey, _, err := decodeRSAPrivateKey(BadKeyData)
	if err == nil && privateKey != nil {
		t.Error("unexpected result: decoding of bad private key data should have returned an error")
		t.FailNow()
	}
}

func TestDecodePrivateKeyRSAEncrypted(t *testing.T) {

	// Expected to fail on decode of encrypted pem private key
	privateKey, _, err := decodeRSAPrivateKey(PrivateKeyEncryptedRSA2048)
	if err == nil && privateKey != nil {
		t.Error("unexpected result: decoding of encrypted private key should have returned an error")
		t.FailNow()
	}
}

func TestVerifyAndEncodeCertificateBadData(t *testing.T) {
	// should fail bad base64 data
	_, _, _, _, err := verifyCertKeyPair(BadCertData, BadKeyData, "")
	if err == nil {
		t.Errorf("Unexpected result: there should have been a base64 decoding failure")
	}
}

func TestVerifyAndEncodeCertificateSelfSignedNoSkiAkiCertKeyPair(t *testing.T) {

	certChain, certPrivateKey, unknownAuth, _, err := verifyCertKeyPair(SelfSignedNOSKIAKIRSACertificate, SelfSignedNOSKIAKIRSAPrivateKey, "")

	if err != nil {
		t.Errorf("Unexpected result: a certificate verification error should have occured")
	}

	if !unknownAuth {
		t.Errorf("Unexpected result: certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Error("unexpected result: can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("unexpected result: remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}


	if len(certChain) == 0 {
		t.Error("unexpected: certChain should not be empty")
	}
}

func TestVerifyAndEncodeCertificateSelfSignedCertKeyPair(t *testing.T) {

	certChain, certPrivateKey, unknownAuth, _, err := verifyCertKeyPair(SelfSignedRSACertificate, SelfSignedRSAPrivateKey, "")

	if err != nil {
		t.Errorf("Unexpected result, a certificate verification error should have occured")
	}

	if !unknownAuth {
		t.Errorf("Unexpected result, certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Error("unexpected result: can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("unexpected result: remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if len(certChain) == 0 {
		t.Error("unexpected: certchain should not empty")
	}

}

func TestVerifyAndEncodeCertificateSelfSignedCertKeyPairMisMatchedPrivateKey(t *testing.T) {

	// Should fail on cert/private-key mismatch
	_, _, _, _, err := verifyCertKeyPair(SelfSignedRSACertificate, PrivateKeyPKCS1RSA2048, "")

	if err == nil {
		t.Errorf("Unexpected result, a certificate/key modulus mismatch error should have occurred")
	}
}

func TestVerifyAndEncodeCertificateCASignedCertKeyPair(t *testing.T) {

	// Should succeed, but with unknown authority warning.
	certChain, certPrivateKey, unknownAuth, _, err := verifyCertKeyPair(CASignedRSACertificateChain, CASignedRSACertificateChainPrivateKey, "")

	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	if !unknownAuth {
		t.Errorf("Unexpected result, certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Error("unexpected result: can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("unexpected result: remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if len(certChain) == 0 {
		t.Error("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateCASignedCertKeyPairWithRootCA(t *testing.T) {

	// should succeed and be fully validated.
	certChain, certPrivateKey, unknownAuth, _, err := verifyCertKeyPair(CASignedRSACertificateChain, CASignedRSACertificateChainPrivateKey, CASignedRSARootCA)

	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	if unknownAuth {
		t.Errorf("Unexpected result: warning for unknown authority even though rootCA is in certChain")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Error("unexpected result: can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("unexpected result: remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if len(certChain) == 0 {
		t.Error("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateCASignedNoSkiAkiCertKeyPair(t *testing.T) {

	// Should succeed, but with unknown authority warning.
	certChain, certPrivateKey, unknownAuth, _, err := verifyCertKeyPair(CASignedNoSkiAkiRSACertificateChain, CASignedNoSkiAkiRSAPrivateKey, "")

	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	if !unknownAuth {
		t.Errorf("Unexpected result, certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Error("unexpected result: can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("unexpected result: remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if len(certChain) == 0 {
		t.Error("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateCASignedNoSkiAkiCertKeyPairWithRootCA(t *testing.T) {

	// should succeed and be fully validated despite not having subject/authority key identifier(s).
	certChain, certPrivateKey, unknownAuth, _, err := verifyCertKeyPair(CASignedNoSkiAkiRSACertificateChain, CASignedNoSkiAkiRSAPrivateKey, CASignedNoSkiAkiRSARootCA)

	if err != nil {
		t.Errorf("Unexpected result: " + err.Error())
	}

	if unknownAuth {
		t.Errorf("Unexpected result: warning for unknown authority even though rootCA is in certChain")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Error("unexpected result: can't decode cleaned private key pem block")
		t.FailNow()
	} else if len(remain) > 0 {
		t.Error("unexpected result: remaining bytes after decode > 0. expected: 0")
		t.FailNow()
	}

	if len(certChain) == 0 {
		t.Error("unexpected: certchain should not empty")
	}
}


func TestVerifyAndEncodeCertificateECCCertificateKeyPair(t *testing.T) {
	// Should fail due to unsupported private/public key algorithm
	_, _, _, _, err := verifyCertKeyPair(SelfSignedECCCertificate, SelfSignedECCPrivateKey, "")

	if err == nil {
		t.Errorf("Unexpected result, cert/key PKI algorithm for ECC is unsupported.")
	}
}