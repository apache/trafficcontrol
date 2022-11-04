package deliveryservice

import (
	"encoding/pem"
	"strings"
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
	BadKeyData  = "This is bad private key data and it is not pem encoded"

	PrivateKeyPKCS1RSA2048 = `
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:a4:8c:9e:da:0e:9c:d2:2a:c5:82:83:00:ce:d3:
    de:12:e5:78:65:7f:0d:74:e4:51:b1:5d:83:76:2e:
    24:57:d8:0a:62:20:e4:c5:0f:5e:39:a9:77:35:e6:
    bd:90:31:1a:52:94:6f:93:69:ee:56:5e:63:6c:51:
    b7:b0:ea:2b:7b:d5:6c:e9:85:2c:0a:f2:02:0c:a0:
    94:3e:5e:4e:af:15:11:27:a0:52:88:2d:a5:d2:35:
    44:e8:55:61:5d:ff:69:2d:7f:8e:47:c9:59:98:c6:
    7d:18:a6:f0:d6:79:46:18:ac:1d:17:74:fb:ea:03:
    99:15:21:d0:7d:3e:7b:bc:d1:6c:23:44:3e:f0:d8:
    56:6c:37:25:36:8f:c0:9c:fa:50:b8:1b:3a:a1:c6:
    a1:f3:70:40:55:09:37:81:34:4c:1c:ed:fe:ac:c2:
    ee:bd:75:69:a4:10:6a:0f:e3:f9:39:4f:8b:45:13:
    ab:8e:80:ee:96:e6:f6:41:43:e2:47:44:39:0d:cc:
    ea:30:28:c3:21:00:7d:e8:b4:5e:af:23:78:77:1f:
    e9:e3:1e:0f:eb:64:8b:40:1e:9d:77:6b:c7:bc:93:
    66:a5:f9:7f:08:1c:c0:75:22:c1:46:76:bd:99:25:
    7a:c7:0e:36:f6:db:b9:6f:d6:78:f0:36:b9:82:9f:
    62:81
publicExponent: 65537 (0x10001)
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
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:a4:8c:9e:da:0e:9c:d2:2a:c5:82:83:00:ce:d3:
    de:12:e5:78:65:7f:0d:74:e4:51:b1:5d:83:76:2e:
    24:57:d8:0a:62:20:e4:c5:0f:5e:39:a9:77:35:e6:
    bd:90:31:1a:52:94:6f:93:69:ee:56:5e:63:6c:51:
    b7:b0:ea:2b:7b:d5:6c:e9:85:2c:0a:f2:02:0c:a0:
    94:3e:5e:4e:af:15:11:27:a0:52:88:2d:a5:d2:35:
    44:e8:55:61:5d:ff:69:2d:7f:8e:47:c9:59:98:c6:
    7d:18:a6:f0:d6:79:46:18:ac:1d:17:74:fb:ea:03:
    99:15:21:d0:7d:3e:7b:bc:d1:6c:23:44:3e:f0:d8:
    56:6c:37:25:36:8f:c0:9c:fa:50:b8:1b:3a:a1:c6:
    a1:f3:70:40:55:09:37:81:34:4c:1c:ed:fe:ac:c2:
    ee:bd:75:69:a4:10:6a:0f:e3:f9:39:4f:8b:45:13:
    ab:8e:80:ee:96:e6:f6:41:43:e2:47:44:39:0d:cc:
    ea:30:28:c3:21:00:7d:e8:b4:5e:af:23:78:77:1f:
    e9:e3:1e:0f:eb:64:8b:40:1e:9d:77:6b:c7:bc:93:
    66:a5:f9:7f:08:1c:c0:75:22:c1:46:76:bd:99:25:
    7a:c7:0e:36:f6:db:b9:6f:d6:78:f0:36:b9:82:9f:
    62:81
publicExponent: 65537 (0x10001)
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
	PrivateKeyECDSANISTPrime256V1 = `
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
Private-Key: (256 bit)
priv:
    cd:6e:e5:7f:9f:9a:e9:da:d2:32:b9:36:80:4e:f0:
    24:56:26:24:f3:a2:98:3d:c7:3a:a3:b3:7e:1d:07:
    7a:2c
pub:
    04:ab:48:07:ec:1d:6a:d9:c5:8c:0e:f6:13:98:5a:
    fe:c4:d9:9b:ac:90:4c:2b:d1:11:99:99:12:7d:17:
    da:0e:6c:f1:aa:f1:39:e8:5a:c1:a7:76:ae:7b:fc:
    43:89:a1:69:f7:29:e2:b8:cf:25:a9:b4:d6:1f:0c:
    a8:a6:b6:8c:94
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIM1u5X+fmuna0jK5NoBO8CRWJiTzopg9xzqjs34dB3osoAoGCCqGSM49
AwEHoUQDQgAEq0gH7B1q2cWMDvYTmFr+xNmbrJBMK9ERmZkSfRfaDmzxqvE56FrB
p3aue/xDiaFp9yniuM8lqbTWHwyopraMlA==
-----END EC PRIVATE KEY-----
`
	PrivateKeyECDSANISTPrime256V1WithoutParams = `
Private-Key: (256 bit)
priv:
    cd:6e:e5:7f:9f:9a:e9:da:d2:32:b9:36:80:4e:f0:
    24:56:26:24:f3:a2:98:3d:c7:3a:a3:b3:7e:1d:07:
    7a:2c
pub:
    04:ab:48:07:ec:1d:6a:d9:c5:8c:0e:f6:13:98:5a:
    fe:c4:d9:9b:ac:90:4c:2b:d1:11:99:99:12:7d:17:
    da:0e:6c:f1:aa:f1:39:e8:5a:c1:a7:76:ae:7b:fc:
    43:89:a1:69:f7:29:e2:b8:cf:25:a9:b4:d6:1f:0c:
    a8:a6:b6:8c:94
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIM1u5X+fmuna0jK5NoBO8CRWJiTzopg9xzqjs34dB3osoAoGCCqGSM49
AwEHoUQDQgAEq0gH7B1q2cWMDvYTmFr+xNmbrJBMK9ERmZkSfRfaDmzxqvE56FrB
p3aue/xDiaFp9yniuM8lqbTWHwyopraMlA==
-----END EC PRIVATE KEY-----
`
	PrivateKeyECDSANISTSecP384R1 = `
ASN1 OID: secp384r1
NIST CURVE: P-384
-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
Private-Key: (384 bit)
priv:
    79:eb:77:84:04:07:4b:70:b0:69:81:48:b7:58:eb:
    33:25:2e:89:38:79:00:70:07:8f:01:65:ff:7e:f8:
    2e:28:fe:7a:6e:b6:c8:b1:99:e5:fc:89:6e:8f:23:
    3a:44:00
pub:
    04:ad:5b:40:5f:67:62:82:92:60:50:bf:d0:9c:f0:
    64:fd:86:81:83:46:e2:6f:7d:b8:f4:72:68:08:47:
    3d:72:86:a1:f2:d4:fb:30:df:fc:fc:16:18:41:92:
    19:b5:63:42:27:27:bc:71:b7:40:ac:39:8b:36:fd:
    f6:80:3a:63:5c:2a:dd:e6:c3:0b:61:2d:c1:32:f5:
    75:59:af:06:e8:83:9b:3b:e9:d6:ad:97:3f:fb:89:
    57:78:26:8e:ef:4b:cf
ASN1 OID: secp384r1
NIST CURVE: P-384
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDB563eEBAdLcLBpgUi3WOszJS6JOHkAcAePAWX/fvguKP56brbIsZnl
/IlujyM6RACgBwYFK4EEACKhZANiAAStW0BfZ2KCkmBQv9Cc8GT9hoGDRuJvfbj0
cmgIRz1yhqHy1Psw3/z8FhhBkhm1Y0InJ7xxt0CsOYs2/faAOmNcKt3mwwthLcEy
9XVZrwbog5s76datlz/7iVd4Jo7vS88=
-----END EC PRIVATE KEY-----
`
	PrivateKeyECDSANISTSecP384R1WithoutParams = `
Private-Key: (384 bit)
priv:
    79:eb:77:84:04:07:4b:70:b0:69:81:48:b7:58:eb:
    33:25:2e:89:38:79:00:70:07:8f:01:65:ff:7e:f8:
    2e:28:fe:7a:6e:b6:c8:b1:99:e5:fc:89:6e:8f:23:
    3a:44:00
pub:
    04:ad:5b:40:5f:67:62:82:92:60:50:bf:d0:9c:f0:
    64:fd:86:81:83:46:e2:6f:7d:b8:f4:72:68:08:47:
    3d:72:86:a1:f2:d4:fb:30:df:fc:fc:16:18:41:92:
    19:b5:63:42:27:27:bc:71:b7:40:ac:39:8b:36:fd:
    f6:80:3a:63:5c:2a:dd:e6:c3:0b:61:2d:c1:32:f5:
    75:59:af:06:e8:83:9b:3b:e9:d6:ad:97:3f:fb:89:
    57:78:26:8e:ef:4b:cf
ASN1 OID: secp384r1
NIST CURVE: P-384
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDB563eEBAdLcLBpgUi3WOszJS6JOHkAcAePAWX/fvguKP56brbIsZnl
/IlujyM6RACgBwYFK4EEACKhZANiAAStW0BfZ2KCkmBQv9Cc8GT9hoGDRuJvfbj0
cmgIRz1yhqHy1Psw3/z8FhhBkhm1Y0InJ7xxt0CsOYs2/faAOmNcKt3mwwthLcEy
9XVZrwbog5s76datlz/7iVd4Jo7vS88=
-----END EC PRIVATE KEY-----
`
	PrivateKeyECDSANISTSecP384R1Encrypted = `
-----BEGIN EC PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: AES-256-CBC,0B0F937007A25C35FBB3DEC9F09C0343

KIRWqbmNfP3xnnDE8f9Ndx/aaNvQBPCucQrVHc6ZYpImPnVmIzH/eOZMyio7HQkZ
tH2ggwwI+zSg3cJWTehJaR9j9qiFtPH+UDEA03co2QyIyERk1wI5ev4hv822tmtl
/TrYpdjqNkfDZUcZscuf1VHkjSrAwn+3K0NV5hUGfdhWryZ7B16iKyCJrSrbde4x
E34vrABCPJZtg/O7SbXQL8cURtVoEdbT+AveW3qoh5g=
-----END EC PRIVATE KEY-----
`
	SelfSignedRSACertificate = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            a7:ec:05:a8:32:41:eb:ae
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = CO, L = Apache Traffic Control, O = Traffic Ops, OU = Unit Testing, CN = *.invalid2.invalid
        Validity
            Not Before: Mar  6 17:08:44 2019 GMT
            Not After : Mar  1 17:08:44 2039 GMT
        Subject: C = US, ST = CO, L = Apache Traffic Control, O = Traffic Ops, OU = Unit Testing, CN = *.invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:c3:97:53:f2:18:de:e2:00:29:bd:10:8e:55:82:
                    d1:e1:04:19:00:15:05:33:96:a9:77:25:6f:e3:af:
                    d9:f6:5b:27:e0:c6:7e:b4:66:98:14:64:e1:2c:f7:
                    1f:9a:7f:97:bf:66:54:37:79:f8:da:66:51:c8:56:
                    4a:eb:7d:a3:f9:58:9b:dc:db:67:a1:c3:03:21:3b:
                    3b:db:ed:3b:bb:8d:1a:af:8e:76:0d:e5:53:38:77:
                    41:26:73:87:eb:da:31:e6:11:33:3a:1d:18:af:bb:
                    ea:f9:19:ea:3c:46:25:13:1a:cd:93:b6:2e:08:6e:
                    2c:7e:e0:dd:fb:74:0d:b7:7f:40:51:8e:40:79:0f:
                    ff:2d:8a:78:db:7b:8d:eb:84:e0:e4:fa:8e:30:3a:
                    39:d4:86:fa:af:e0:5c:0b:45:00:34:0a:93:1f:35:
                    e2:91:00:ca:9d:53:e0:ae:7d:ae:d7:8d:d8:aa:02:
                    6d:18:8d:ae:0c:6e:7b:08:24:7a:ed:d5:38:64:da:
                    82:4c:1b:c9:ab:f5:a3:7b:b1:13:ec:1f:73:7a:f4:
                    4c:e7:56:fb:6c:4b:50:ef:4a:95:a1:b7:a6:2a:13:
                    f1:58:75:61:cc:19:be:08:5e:fa:69:2e:7b:2a:1d:
                    d5:59:5d:2a:5e:2f:80:7d:0d:b1:41:0c:a6:2e:73:
                    07:af
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Key Usage: 
                Key Encipherment, Data Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
            X509v3 Subject Key Identifier: 
                26:0B:85:1B:C8:9B:89:2D:3E:A2:B8:6D:20:13:C3:50:F0:13:8C:BE
            X509v3 Authority Key Identifier: 
                keyid:26:0B:85:1B:C8:9B:89:2D:3E:A2:B8:6D:20:13:C3:50:F0:13:8C:BE

            X509v3 Subject Alternative Name: 
                DNS:*.invalid2.invalid
    Signature Algorithm: sha256WithRSAEncryption
         75:b6:1d:d9:53:17:33:af:96:52:e5:87:a6:e0:a1:8b:00:04:
         1c:93:4a:3b:d5:cb:e6:20:c2:d7:b3:8b:b4:8a:87:a1:33:a5:
         03:c4:da:ce:c3:2e:57:6b:12:f0:79:97:9b:22:e5:20:67:8e:
         3e:29:95:3f:51:ea:ea:4a:29:d6:13:7c:ec:36:91:db:57:ee:
         e3:eb:46:85:8e:e0:8a:02:3f:bd:5a:d2:55:55:af:c0:d8:a2:
         cf:54:c6:1a:a0:62:13:81:da:21:1d:e6:0b:2a:a5:8e:ea:2a:
         35:56:f2:e2:6a:15:da:e8:6d:67:12:53:e1:47:07:41:53:9a:
         a7:27:ae:75:23:cd:b5:c5:ff:53:29:e2:14:64:cf:9b:99:8c:
         aa:0c:4d:85:28:76:d9:fa:10:41:00:6a:d3:db:ae:60:37:7f:
         da:64:a8:b9:0b:9b:36:60:56:77:39:ce:dd:95:9d:59:52:16:
         6a:88:b7:93:af:c7:26:73:d6:d9:01:d3:58:f6:9a:43:9f:7c:
         6c:f1:38:48:19:09:50:cd:42:38:a3:61:44:dc:e7:70:03:6f:
         0f:3b:bf:3f:0e:f9:d9:f2:85:5e:17:f2:83:d7:da:04:c4:9f:
         6d:8f:a7:73:2d:92:2b:93:16:91:0e:ee:c2:0f:bd:12:e3:53:
         d0:2f:2f:10
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
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:c3:97:53:f2:18:de:e2:00:29:bd:10:8e:55:82:
    d1:e1:04:19:00:15:05:33:96:a9:77:25:6f:e3:af:
    d9:f6:5b:27:e0:c6:7e:b4:66:98:14:64:e1:2c:f7:
    1f:9a:7f:97:bf:66:54:37:79:f8:da:66:51:c8:56:
    4a:eb:7d:a3:f9:58:9b:dc:db:67:a1:c3:03:21:3b:
    3b:db:ed:3b:bb:8d:1a:af:8e:76:0d:e5:53:38:77:
    41:26:73:87:eb:da:31:e6:11:33:3a:1d:18:af:bb:
    ea:f9:19:ea:3c:46:25:13:1a:cd:93:b6:2e:08:6e:
    2c:7e:e0:dd:fb:74:0d:b7:7f:40:51:8e:40:79:0f:
    ff:2d:8a:78:db:7b:8d:eb:84:e0:e4:fa:8e:30:3a:
    39:d4:86:fa:af:e0:5c:0b:45:00:34:0a:93:1f:35:
    e2:91:00:ca:9d:53:e0:ae:7d:ae:d7:8d:d8:aa:02:
    6d:18:8d:ae:0c:6e:7b:08:24:7a:ed:d5:38:64:da:
    82:4c:1b:c9:ab:f5:a3:7b:b1:13:ec:1f:73:7a:f4:
    4c:e7:56:fb:6c:4b:50:ef:4a:95:a1:b7:a6:2a:13:
    f1:58:75:61:cc:19:be:08:5e:fa:69:2e:7b:2a:1d:
    d5:59:5d:2a:5e:2f:80:7d:0d:b1:41:0c:a6:2e:73:
    07:af
publicExponent: 65537 (0x10001)
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
	SelfSignedRSACertificateNoServerAuthExtKeyUsage = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            12:d1:c8:08:54:57:af:f0:73:fb:0f:26:14:02:79:93:2b:b9:95:a1
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.noserverauth.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Validity
            Not Before: May  3 17:12:01 2019 GMT
            Not After : Apr 28 17:12:01 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.noserverauth.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:a7:dc:3b:4d:05:5d:61:27:5b:d7:16:61:2b:f0:
                    fd:2c:e6:14:1d:36:ab:c2:0b:17:27:db:16:69:31:
                    6c:30:c8:67:20:a6:89:92:28:73:9c:50:76:03:94:
                    f5:97:30:c1:2e:17:73:a0:1e:99:56:ef:2e:7f:8b:
                    fe:7c:39:d3:fc:3b:f9:c5:fb:05:20:30:31:b7:00:
                    34:aa:e0:af:43:92:a8:7a:e6:bd:74:99:ef:1a:67:
                    ec:35:9b:de:7a:db:a6:19:c8:cd:13:1f:0a:24:2c:
                    ef:a7:32:81:75:76:de:ef:0d:3b:84:1e:77:f9:a5:
                    6a:95:18:bb:1f:92:03:8b:7c:c1:5a:f7:64:df:20:
                    72:86:a5:e1:fe:24:bf:67:59:fe:76:40:5e:33:d8:
                    10:0b:a0:67:13:fe:86:c2:85:e8:88:65:00:33:85:
                    0e:cb:0e:05:45:97:07:e4:c4:e7:8c:fd:c1:d5:15:
                    c2:8b:ce:48:e4:28:3c:ea:68:56:c7:3c:ad:c5:46:
                    1e:41:a6:a6:35:e7:77:90:9c:80:bb:fd:33:59:3b:
                    81:2e:e7:10:49:54:c0:d5:86:01:8e:1e:d1:12:6a:
                    0f:2d:f6:7e:df:4f:94:66:93:54:b3:71:00:ac:d2:
                    0b:8f:46:40:15:ae:84:91:99:ec:09:a0:c9:12:2a:
                    7c:1f
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:TRUE
            Netscape Cert Type: 
                SSL Server
            X509v3 Key Usage: critical
                Key Encipherment
            X509v3 Extended Key Usage: critical
                Code Signing
            X509v3 Subject Key Identifier: 
                61:8A:5E:90:59:68:6F:07:B9:49:E3:A0:8E:9C:A9:97:A7:25:DC:15
            X509v3 Subject Alternative Name: 
                DNS:*.noserverauth.mycdn.ciab.test
    Signature Algorithm: sha256WithRSAEncryption
         6c:98:03:c8:77:50:e3:f8:b5:9d:08:d8:5d:20:90:34:96:eb:
         60:2d:7e:4b:4b:e1:86:74:f9:6d:c7:18:41:7a:c3:ba:3b:b9:
         d8:e9:b1:28:e9:5c:16:2c:96:b3:35:a0:05:c9:f1:90:4a:d3:
         16:2c:ad:8c:8b:c5:16:fa:d5:26:5a:f1:a0:d0:60:ee:97:f7:
         1e:da:6f:a9:ba:89:a5:80:e0:dd:c6:2a:40:31:9f:04:b5:34:
         30:19:1f:20:f0:6a:9f:a2:34:07:2b:cb:e0:cc:eb:33:87:06:
         a9:64:15:00:04:f6:86:4d:b9:4c:16:8b:d4:03:93:a9:a2:ae:
         27:75:69:d4:ab:15:4e:66:a5:97:9b:ef:c2:ad:3b:da:de:54:
         29:a9:cc:90:21:cc:3f:67:e3:7b:71:ee:6f:4b:eb:cb:8d:1e:
         9f:dd:a0:bf:98:87:48:f0:33:81:a0:2d:08:d5:3b:a0:94:49:
         25:2d:67:04:b2:f0:0c:95:a5:15:3e:7a:fc:d8:1b:fb:14:0e:
         c0:88:2d:7d:56:7a:41:01:fe:29:00:4a:f8:23:af:48:22:e1:
         b9:33:19:49:d7:94:f0:b6:90:41:d2:07:70:ac:3d:3a:ea:28:
         14:0f:ae:b4:e1:42:1b:8f:27:9a:fc:eb:d3:9e:9a:4d:0d:23:
         5a:c8:72:6c
-----BEGIN CERTIFICATE-----
MIIEiTCCA3GgAwIBAgIUEtHICFRXr/Bz+w8mFAJ5kyu5laEwDQYJKoZIhvcNAQEL
BQAwgbExCzAJBgNVBAYTAlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMG
RGVudmVyMRUwEwYDVQQKEwxDRE4taW4tYS1Cb3gxFTATBgNVBAsTDENETi1pbi1h
LUJveDEnMCUGA1UEAxQeKi5ub3NlcnZlcmF1dGgubXljZG4uY2lhYi50ZXN0MScw
JQYJKoZIhvcNAQkBFhhuby1yZXBseUBpbmZyYS5jaWFiLnRlc3QwHhcNMTkwNTAz
MTcxMjAxWhcNMzkwNDI4MTcxMjAxWjCBsTELMAkGA1UEBhMCVVMxETAPBgNVBAgT
CENvbG9yYWRvMQ8wDQYDVQQHEwZEZW52ZXIxFTATBgNVBAoTDENETi1pbi1hLUJv
eDEVMBMGA1UECxMMQ0ROLWluLWEtQm94MScwJQYDVQQDFB4qLm5vc2VydmVyYXV0
aC5teWNkbi5jaWFiLnRlc3QxJzAlBgkqhkiG9w0BCQEWGG5vLXJlcGx5QGluZnJh
LmNpYWIudGVzdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKfcO00F
XWEnW9cWYSvw/SzmFB02q8ILFyfbFmkxbDDIZyCmiZIoc5xQdgOU9ZcwwS4Xc6Ae
mVbvLn+L/nw50/w7+cX7BSAwMbcANKrgr0OSqHrmvXSZ7xpn7DWb3nrbphnIzRMf
CiQs76cygXV23u8NO4Qed/mlapUYux+SA4t8wVr3ZN8gcoal4f4kv2dZ/nZAXjPY
EAugZxP+hsKF6IhlADOFDssOBUWXB+TE54z9wdUVwovOSOQoPOpoVsc8rcVGHkGm
pjXnd5CcgLv9M1k7gS7nEElUwNWGAY4e0RJqDy32ft9PlGaTVLNxAKzSC49GQBWu
hJGZ7AmgyRIqfB8CAwEAAaOBljCBkzAMBgNVHRMEBTADAQH/MBEGCWCGSAGG+EIB
AQQEAwIGQDAOBgNVHQ8BAf8EBAMCBSAwFgYDVR0lAQH/BAwwCgYIKwYBBQUHAwMw
HQYDVR0OBBYEFGGKXpBZaG8HuUnjoI6cqZenJdwVMCkGA1UdEQQiMCCCHioubm9z
ZXJ2ZXJhdXRoLm15Y2RuLmNpYWIudGVzdDANBgkqhkiG9w0BAQsFAAOCAQEAbJgD
yHdQ4/i1nQjYXSCQNJbrYC1+S0vhhnT5bccYQXrDuju52OmxKOlcFiyWszWgBcnx
kErTFiytjIvFFvrVJlrxoNBg7pf3HtpvqbqJpYDg3cYqQDGfBLU0MBkfIPBqn6I0
ByvL4MzrM4cGqWQVAAT2hk25TBaL1AOTqaKuJ3Vp1KsVTmall5vvwq072t5UKanM
kCHMP2fje3Hub0vry40en92gv5iHSPAzgaAtCNU7oJRJJS1nBLLwDJWlFT56/Ngb
+xQOwIgtfVZ6QQH+KQBK+COvSCLhuTMZSdeU8LaQQdIHcKw9OuooFA+utOFCG48n
mvzr056aTQ0jWshybA==
-----END CERTIFICATE-----
`

	SelfSignedRSAPrivateKeyNoServerAuthExtKeyUsage = `
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:a7:dc:3b:4d:05:5d:61:27:5b:d7:16:61:2b:f0:
    fd:2c:e6:14:1d:36:ab:c2:0b:17:27:db:16:69:31:
    6c:30:c8:67:20:a6:89:92:28:73:9c:50:76:03:94:
    f5:97:30:c1:2e:17:73:a0:1e:99:56:ef:2e:7f:8b:
    fe:7c:39:d3:fc:3b:f9:c5:fb:05:20:30:31:b7:00:
    34:aa:e0:af:43:92:a8:7a:e6:bd:74:99:ef:1a:67:
    ec:35:9b:de:7a:db:a6:19:c8:cd:13:1f:0a:24:2c:
    ef:a7:32:81:75:76:de:ef:0d:3b:84:1e:77:f9:a5:
    6a:95:18:bb:1f:92:03:8b:7c:c1:5a:f7:64:df:20:
    72:86:a5:e1:fe:24:bf:67:59:fe:76:40:5e:33:d8:
    10:0b:a0:67:13:fe:86:c2:85:e8:88:65:00:33:85:
    0e:cb:0e:05:45:97:07:e4:c4:e7:8c:fd:c1:d5:15:
    c2:8b:ce:48:e4:28:3c:ea:68:56:c7:3c:ad:c5:46:
    1e:41:a6:a6:35:e7:77:90:9c:80:bb:fd:33:59:3b:
    81:2e:e7:10:49:54:c0:d5:86:01:8e:1e:d1:12:6a:
    0f:2d:f6:7e:df:4f:94:66:93:54:b3:71:00:ac:d2:
    0b:8f:46:40:15:ae:84:91:99:ec:09:a0:c9:12:2a:
    7c:1f
publicExponent: 65537 (0x10001)
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAp9w7TQVdYSdb1xZhK/D9LOYUHTarwgsXJ9sWaTFsMMhnIKaJ
kihznFB2A5T1lzDBLhdzoB6ZVu8uf4v+fDnT/Dv5xfsFIDAxtwA0quCvQ5Koeua9
dJnvGmfsNZveetumGcjNEx8KJCzvpzKBdXbe7w07hB53+aVqlRi7H5IDi3zBWvdk
3yByhqXh/iS/Z1n+dkBeM9gQC6BnE/6GwoXoiGUAM4UOyw4FRZcH5MTnjP3B1RXC
i85I5Cg86mhWxzytxUYeQaamNed3kJyAu/0zWTuBLucQSVTA1YYBjh7REmoPLfZ+
30+UZpNUs3EArNILj0ZAFa6EkZnsCaDJEip8HwIDAQABAoIBAGWau9ZaGfS1szSV
GkpTu5uSxLgOIJb6yZBZX85amQdKNoof5AOxMpF6boSqhKF4ZGY20cko3F4vtrCD
l42wHy19TCnXUHn0UhNYL4kDKXM4cXy68BCFIKKWJvcoGtm43GidD+y0DBprjMBi
pNPqGPUPyGenXa2hv8rxxkpMwpKJ/JHdONq0m+eTu3lCRdxhVVHbOn1xQmJTgP2B
CyrtmBD7rJ5AhDNYQ7GfL9WwVf1UeFNaOuAOEK9T7Weqm3Ak5Pd1dJUwukXk/vND
Rcu+cWhklNTVIUw9XwXawPpFm03iZiTcFurAcrShXrtW7aY38S1e7w1fRSSxhrfQ
GroShXECgYEA3rjHANz5dx/Z0EUQ3rK6MgOuKpsS4rAmpydqY7PFfNkSgQzsUwot
o4VRIiKd16voodvT9CZausgCut6ZqCpVgwZY0Mg9FBEBXfxPIOg5oRNmje7RHAfD
Q9I2yiEgBZ6exi3m4lZKgWSNvKkocY/HLIwsDz4kEq4ux3Gr2+zkIacCgYEAwPD5
N0R25aO06gI/eeYOV0P0Y0h3OKCtF0b/PwU4tKSuqQiojeIvmoJspbCZJbg+D6Da
0kXgEhybchVBSzgeOggZb+sZl0Odg1eqJea/kbz7JyO0NcZu3t7wg5AVynssF6s0
JUcDFxELfW9GYLtdzR8LeWvgaezSwZBeYmV1cMkCgYBtri1CPZAUm/jV2c1O/lE3
ZByXGrsYK4s9cemwo80ziGrWZpjS5AZJqtOjrcxxc1UisHEWoPS5WtoNUKX27LIj
zjJazuFVSnKT6DbHi9Ulf7pXVy5fUWtVsOYOcHWmjtC948j52Wjjg7NRHzStiBKb
24OvFfkJwgGDcnUh3u0RrQKBgG7ZxAVp03nKXYXY9sk9UN34T+++0ah6QBhQlROL
F3JJ74N0Uwr5eeompu9nEAYo3ZczDqWiucMOJo0cAyCJRGyI/LxdcZ2Dnnq4oiwW
b9f2oMFy9PW0ZTytD7g2zx4/OCz9Ev+b1f2psFVH2kJ3Q8Q24uvG++8/vjKxlFip
/BhpAoGAJmo5n5fR9+q1Zwb4Gs9i/gquxcJsWo9eaCgQpt+IK0oqqgi/8xx61PCW
Xeg+Gf0gAv5J/WfbYpBCLD4K8o1COm7ZqZGqG3kFTWYjh0jRm8hLQF+8kkn0Fby5
ohq7Fw9P96opYCnxqau5nfiHExfyfLNQZlj7dPHZZPCEhUFYSSo=
-----END RSA PRIVATE KEY-----
`
	SelfSignedRSACertificateNoKeyEnciphermentKeyUsage = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            1b:09:1a:1a:8a:6a:68:10:cc:72:35:ab:d4:9f:f2:fb:ac:28:79:d1
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.nokeyencipherment.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Validity
            Not Before: May  3 17:14:21 2019 GMT
            Not After : Apr 28 17:14:21 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.nokeyencipherment.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:a2:39:f1:a1:4f:c1:20:a4:3a:0c:2e:ad:3c:f1:
                    1e:47:47:43:17:7e:0c:1c:f5:7a:05:66:60:79:a3:
                    26:6f:4c:37:10:90:62:47:e3:c6:76:b4:51:6d:77:
                    7a:a8:6b:9b:1b:ca:2b:b9:fc:8f:23:fc:9e:3e:35:
                    7b:d8:04:b2:da:4c:d3:13:cf:87:c5:92:71:66:96:
                    cb:32:47:70:97:52:3b:e1:df:09:73:59:52:74:4a:
                    40:3e:18:b7:e7:65:d5:f5:cc:9f:b7:e4:6a:df:fd:
                    b2:0b:79:5d:24:97:f5:8a:a0:51:32:23:d0:99:5b:
                    b5:5b:4d:a4:65:ba:80:b1:ac:83:65:1a:d2:f3:8b:
                    e2:be:90:74:5c:97:1b:f8:11:9e:93:72:02:59:d2:
                    54:f8:ec:e4:b1:51:34:d3:0c:a1:f8:2f:6a:c5:e9:
                    d3:a5:b6:61:a2:19:db:0c:ca:16:6c:22:cc:0a:6c:
                    cf:9e:12:58:84:60:0e:90:dc:24:f9:b3:fd:c5:c2:
                    28:52:db:b8:16:26:af:8e:44:02:ae:cf:be:d1:97:
                    e9:f7:7a:9c:d6:68:07:c3:d0:cf:cf:f9:e0:8f:34:
                    51:2a:26:03:15:c3:bd:b1:ed:29:43:f5:ec:a2:36:
                    7b:34:8a:ad:8e:bd:eb:e3:31:d8:81:bf:2f:b7:e2:
                    33:5b
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:TRUE
            Netscape Cert Type: 
                SSL Server
            X509v3 Key Usage: critical
                Digital Signature
            X509v3 Extended Key Usage: critical
                TLS Web Server Authentication, Code Signing
            X509v3 Subject Key Identifier: 
                B5:48:56:8F:3F:B0:55:94:CA:E4:35:E4:CA:2A:84:19:5A:78:B6:F8
            X509v3 Subject Alternative Name: 
                DNS:*.nokeyencipherment.mycdn.ciab.test
    Signature Algorithm: sha256WithRSAEncryption
         5a:90:a6:c3:42:97:2c:2c:10:d4:03:9b:5f:5a:59:cd:0a:56:
         61:38:0c:ce:c8:0a:54:9a:58:57:b5:ba:1e:76:e4:9e:65:9a:
         f5:01:2a:d3:f3:c1:fb:57:bf:84:b3:ca:c1:30:36:34:fa:12:
         9e:e2:57:f2:8c:4f:c8:72:a9:9a:0b:9e:31:19:ad:31:48:96:
         f0:71:14:db:19:a7:f3:dc:6c:0d:c9:c7:f5:c8:82:9e:dd:93:
         80:11:92:31:0b:ce:70:39:f0:30:4a:ed:0f:f0:b5:fe:5d:68:
         b0:de:f5:30:87:1c:70:dd:5e:2e:24:29:e3:5a:97:4c:c7:9c:
         1a:93:eb:dd:45:85:0a:63:ec:81:be:7c:c9:7b:d6:f0:a4:39:
         b5:6c:8c:a2:a0:60:f0:a4:48:6f:e3:94:8b:fe:d9:59:f3:18:
         10:eb:f2:95:c3:98:8f:ac:3e:1c:86:f1:d8:78:b9:7a:7d:86:
         51:86:38:a0:7a:52:ae:4f:85:17:a6:cc:86:3a:f2:f8:de:ab:
         86:4e:19:c1:24:34:cb:b5:89:6d:71:cd:a7:7a:6c:97:ed:36:
         ce:e7:65:74:d9:8b:13:81:bd:3f:74:e1:6e:c9:bc:14:a2:c1:
         d3:1d:3a:92:94:bb:dd:f7:fa:28:e2:33:5a:16:d0:9e:9c:9f:
         da:d1:54:77
-----BEGIN CERTIFICATE-----
MIIEojCCA4qgAwIBAgIUGwkaGopqaBDMcjWr1J/y+6woedEwDQYJKoZIhvcNAQEL
BQAwgbYxCzAJBgNVBAYTAlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMG
RGVudmVyMRUwEwYDVQQKEwxDRE4taW4tYS1Cb3gxFTATBgNVBAsTDENETi1pbi1h
LUJveDEsMCoGA1UEAxQjKi5ub2tleWVuY2lwaGVybWVudC5teWNkbi5jaWFiLnRl
c3QxJzAlBgkqhkiG9w0BCQEWGG5vLXJlcGx5QGluZnJhLmNpYWIudGVzdDAeFw0x
OTA1MDMxNzE0MjFaFw0zOTA0MjgxNzE0MjFaMIG2MQswCQYDVQQGEwJVUzERMA8G
A1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZlcjEVMBMGA1UEChMMQ0ROLWlu
LWEtQm94MRUwEwYDVQQLEwxDRE4taW4tYS1Cb3gxLDAqBgNVBAMUIyoubm9rZXll
bmNpcGhlcm1lbnQubXljZG4uY2lhYi50ZXN0MScwJQYJKoZIhvcNAQkBFhhuby1y
ZXBseUBpbmZyYS5jaWFiLnRlc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQCiOfGhT8EgpDoMLq088R5HR0MXfgwc9XoFZmB5oyZvTDcQkGJH48Z2tFFt
d3qoa5sbyiu5/I8j/J4+NXvYBLLaTNMTz4fFknFmlssyR3CXUjvh3wlzWVJ0SkA+
GLfnZdX1zJ+35Grf/bILeV0kl/WKoFEyI9CZW7VbTaRluoCxrINlGtLzi+K+kHRc
lxv4EZ6TcgJZ0lT47OSxUTTTDKH4L2rF6dOltmGiGdsMyhZsIswKbM+eEliEYA6Q
3CT5s/3FwihS27gWJq+ORAKuz77Rl+n3epzWaAfD0M/P+eCPNFEqJgMVw72x7SlD
9eyiNns0iq2OvevjMdiBvy+34jNbAgMBAAGjgaUwgaIwDAYDVR0TBAUwAwEB/zAR
BglghkgBhvhCAQEEBAMCBkAwDgYDVR0PAQH/BAQDAgeAMCAGA1UdJQEB/wQWMBQG
CCsGAQUFBwMBBggrBgEFBQcDAzAdBgNVHQ4EFgQUtUhWjz+wVZTK5DXkyiqEGVp4
tvgwLgYDVR0RBCcwJYIjKi5ub2tleWVuY2lwaGVybWVudC5teWNkbi5jaWFiLnRl
c3QwDQYJKoZIhvcNAQELBQADggEBAFqQpsNClywsENQDm19aWc0KVmE4DM7IClSa
WFe1uh525J5lmvUBKtPzwftXv4SzysEwNjT6Ep7iV/KMT8hyqZoLnjEZrTFIlvBx
FNsZp/PcbA3Jx/XIgp7dk4ARkjELznA58DBK7Q/wtf5daLDe9TCHHHDdXi4kKeNa
l0zHnBqT691FhQpj7IG+fMl71vCkObVsjKKgYPCkSG/jlIv+2VnzGBDr8pXDmI+s
PhyG8dh4uXp9hlGGOKB6Uq5PhRemzIY68vjeq4ZOGcEkNMu1iW1xzad6bJftNs7n
ZXTZixOBvT904W7JvBSiwdMdOpKUu933+ijiM1oW0J6cn9rRVHc=
-----END CERTIFICATE-----
`
	SelfSignedRSAPrivateKeyNoKeyEnciphermentKeyUsage = `
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:a2:39:f1:a1:4f:c1:20:a4:3a:0c:2e:ad:3c:f1:
    1e:47:47:43:17:7e:0c:1c:f5:7a:05:66:60:79:a3:
    26:6f:4c:37:10:90:62:47:e3:c6:76:b4:51:6d:77:
    7a:a8:6b:9b:1b:ca:2b:b9:fc:8f:23:fc:9e:3e:35:
    7b:d8:04:b2:da:4c:d3:13:cf:87:c5:92:71:66:96:
    cb:32:47:70:97:52:3b:e1:df:09:73:59:52:74:4a:
    40:3e:18:b7:e7:65:d5:f5:cc:9f:b7:e4:6a:df:fd:
    b2:0b:79:5d:24:97:f5:8a:a0:51:32:23:d0:99:5b:
    b5:5b:4d:a4:65:ba:80:b1:ac:83:65:1a:d2:f3:8b:
    e2:be:90:74:5c:97:1b:f8:11:9e:93:72:02:59:d2:
    54:f8:ec:e4:b1:51:34:d3:0c:a1:f8:2f:6a:c5:e9:
    d3:a5:b6:61:a2:19:db:0c:ca:16:6c:22:cc:0a:6c:
    cf:9e:12:58:84:60:0e:90:dc:24:f9:b3:fd:c5:c2:
    28:52:db:b8:16:26:af:8e:44:02:ae:cf:be:d1:97:
    e9:f7:7a:9c:d6:68:07:c3:d0:cf:cf:f9:e0:8f:34:
    51:2a:26:03:15:c3:bd:b1:ed:29:43:f5:ec:a2:36:
    7b:34:8a:ad:8e:bd:eb:e3:31:d8:81:bf:2f:b7:e2:
    33:5b
publicExponent: 65537 (0x10001)
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAojnxoU/BIKQ6DC6tPPEeR0dDF34MHPV6BWZgeaMmb0w3EJBi
R+PGdrRRbXd6qGubG8orufyPI/yePjV72ASy2kzTE8+HxZJxZpbLMkdwl1I74d8J
c1lSdEpAPhi352XV9cyft+Rq3/2yC3ldJJf1iqBRMiPQmVu1W02kZbqAsayDZRrS
84vivpB0XJcb+BGek3ICWdJU+OzksVE00wyh+C9qxenTpbZhohnbDMoWbCLMCmzP
nhJYhGAOkNwk+bP9xcIoUtu4FiavjkQCrs++0Zfp93qc1mgHw9DPz/ngjzRRKiYD
FcO9se0pQ/XsojZ7NIqtjr3r4zHYgb8vt+IzWwIDAQABAoIBAEm+Kz+XwIO1A4oM
IcXFGW1vUGk6bAkx8TDJM+u3JT6Ml69Y4sQpH0tQdn9bQ4+RsqV0RmI6E1tZdxly
OISexiqDp6Omv+IoypHG1EFbxiuTPxNSzrn3jYq9Qey4UcjHOvaL+MKf+5EsgqXC
mnuK9Bv6+k3fh/BehtclOSjhGaUpvBj4FFpVm5aQT1PlrfsSJsY8AhAsqyO3mO0Z
ndfSd2u3/T795h7XVyarJlZquaExlAe+rsrmdCciU+0l27Wh9HzY145sDtuzSPGm
A8hWredHQ6qVwz4EpVZUj3HDowHi3krqlQitc144CBA/PejTV7jby0piKAZ8agKK
pCiRvxECgYEAzvBKma1oJWVa3oM1sYsan8RcquqHRQoE2q0xDmKgQPa4mihyXfCV
iPDPc2jw+0/M6fNgFlhVppMVqJyvScIELKH0HsRtoZyDiaoVkmQY/EG4xkS9N3zK
6OXs9y/1nHHFwsKe3EDHuPKWJAHDAx+x9OrQ5o67Vc9MDwmwt4bJiDkCgYEAyK/u
3UNWrBVQuT4k8klI485HCJCV6AO2egONw/dGFeZJ9W20Uw/eWAAuphT6HKqwSlBt
0EwTpjL9dnuoTrS2TcTLq7UwTN9Gw/1VhE1of+nbsa0kWAqDqhpCFMZQ2/X9ehDv
cqFemW4HE1la3E9izwlKGU3ehpHwWjRhJbH8kDMCgYB+uS2l4EgHpoK4AneuCrY6
ImBxFf/SKmmAlFCXM5RZU/0GAkDPABZCbt1LGneAHoUouy4bYOrKgAXiZFj/fP1b
a6337WgJcLQoaGyfYgbe60xAtjV9NkF3z92GHet1a0KkmtP3ov/rZTrGQAHw9sbe
abGVjtBvous7xj5elP7zGQKBgQC2TtRwHi8LLmXhkemgTCCyCX6P8kCrv0uyNa5A
Gk6JsGT5VopcdmrmiGvYJfA7wHdbWwsXETU8Ys/MJXN05EdECIV426UgACjZ/DYG
dQd8Q+Z21rHQZOTMzwO+uZVU7Hcyv1W2TY+RU9mLoz2eK2O4bljo+csvdj3gw/qI
ctLb7wKBgFOwu2jZJE66DJ2AUz5XgeY0EOSH8ZIpdEtRQriE8R397KHb8Pmt75fJ
F+nJV58mV+vR/r3qlMViq8Cfb8cZz47Dn00Q2udMoOLKLiKgu1yP2pdCs24VbwLF
7FDXBcacw5kMQfw5DxWDqgXqbwl4cdS62m8kRA/idZNqdHtCmICA
-----END RSA PRIVATE KEY-----
`
	SelfSignedNOSKIAKIRSACertificate = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            8f:38:5b:9b:6c:47:7c:f8
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = CO, L = Apache Traffic Control, O = Traffic Ops, OU = Unit Testing, CN = *.invalid2.invalid
        Validity
            Not Before: Mar  6 17:11:16 2019 GMT
            Not After : Mar  1 17:11:16 2039 GMT
        Subject: C = US, ST = CO, L = Apache Traffic Control, O = Traffic Ops, OU = Unit Testing, CN = *.invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:bd:b2:ca:98:43:96:47:34:02:6a:71:3c:0a:74:
                    19:99:0f:88:63:39:3e:09:1a:86:d6:9f:9c:83:02:
                    59:70:b9:26:45:b3:d5:09:3b:75:28:0a:50:d9:7a:
                    85:78:37:23:57:41:27:f2:9e:8f:8e:b5:5f:04:52:
                    3c:d5:45:a7:57:a1:52:c9:7b:f8:3e:90:27:2e:fe:
                    58:4b:a4:e0:46:b7:b3:af:db:b4:c0:72:1f:af:bb:
                    4b:91:2f:24:91:a2:bd:ab:d7:95:a2:df:a9:79:c3:
                    da:f4:4a:5b:ee:84:b9:a6:cd:a5:e3:55:a6:56:2b:
                    65:06:10:d9:a1:d8:22:c6:3e:1d:82:27:ce:a7:69:
                    6f:fa:b1:e8:14:9b:ba:fb:95:ef:83:7a:29:c2:3b:
                    75:ba:42:e7:9d:72:f3:b8:6a:d9:4f:cc:7e:d9:1a:
                    d5:39:9f:4f:f9:65:ed:b2:9d:67:8a:b5:f7:63:04:
                    d7:8f:d4:48:6a:a6:cb:2f:07:1c:08:85:be:c8:11:
                    8e:6a:94:51:ed:cd:eb:e3:b4:54:e8:e3:4c:30:72:
                    2f:59:da:33:0e:4c:95:22:21:70:83:8c:e0:ef:ab:
                    92:50:22:29:5c:cb:f4:e3:5c:d6:bc:f5:53:7f:09:
                    ac:d7:94:8c:3e:a9:12:57:83:70:f3:c1:47:58:1e:
                    f0:37
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Key Usage: 
                Key Encipherment, Data Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
            X509v3 Subject Alternative Name: 
                DNS:*.invalid2.invalid
    Signature Algorithm: sha256WithRSAEncryption
         6d:90:0d:2b:b7:22:db:95:71:43:3a:7b:41:63:af:20:a5:c8:
         1a:6b:52:a3:52:8f:58:35:f9:2e:57:2c:02:55:a9:fb:66:8c:
         62:26:40:ac:8a:b2:28:9e:e8:09:0d:a9:f5:43:e9:9e:32:c0:
         36:b4:ef:54:5e:8c:52:85:74:b2:87:aa:44:72:a1:1a:e2:e4:
         9e:96:3d:4c:23:2b:05:15:db:42:c1:c5:0e:cc:bd:ce:c0:06:
         35:4d:e7:02:de:b4:5a:74:90:ff:85:63:35:17:78:87:3e:91:
         b7:04:60:75:2c:39:04:65:50:90:76:2f:86:b6:af:c2:50:21:
         db:ea:0d:65:96:e2:ed:fe:9c:47:ba:27:12:e1:eb:44:77:5c:
         d3:6f:5a:21:d8:26:dc:8a:33:0b:15:21:a2:98:fa:a4:aa:64:
         81:3a:48:0e:eb:fd:bc:df:10:37:30:e4:a6:ff:0e:f2:ac:bc:
         a5:ab:f3:7b:fe:0e:5c:ea:42:d9:49:69:18:a6:8c:40:e4:f8:
         08:98:91:2c:3b:0a:05:17:d5:6b:33:30:b8:24:f2:dd:47:7c:
         e3:2e:af:d3:c0:19:b2:23:37:15:14:9d:a0:e1:5d:22:4b:95:
         1f:8b:ab:87:c1:ce:c0:d3:73:d5:03:9d:20:ee:9e:c9:51:54:
         71:e6:af:c1
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
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:bd:b2:ca:98:43:96:47:34:02:6a:71:3c:0a:74:
    19:99:0f:88:63:39:3e:09:1a:86:d6:9f:9c:83:02:
    59:70:b9:26:45:b3:d5:09:3b:75:28:0a:50:d9:7a:
    85:78:37:23:57:41:27:f2:9e:8f:8e:b5:5f:04:52:
    3c:d5:45:a7:57:a1:52:c9:7b:f8:3e:90:27:2e:fe:
    58:4b:a4:e0:46:b7:b3:af:db:b4:c0:72:1f:af:bb:
    4b:91:2f:24:91:a2:bd:ab:d7:95:a2:df:a9:79:c3:
    da:f4:4a:5b:ee:84:b9:a6:cd:a5:e3:55:a6:56:2b:
    65:06:10:d9:a1:d8:22:c6:3e:1d:82:27:ce:a7:69:
    6f:fa:b1:e8:14:9b:ba:fb:95:ef:83:7a:29:c2:3b:
    75:ba:42:e7:9d:72:f3:b8:6a:d9:4f:cc:7e:d9:1a:
    d5:39:9f:4f:f9:65:ed:b2:9d:67:8a:b5:f7:63:04:
    d7:8f:d4:48:6a:a6:cb:2f:07:1c:08:85:be:c8:11:
    8e:6a:94:51:ed:cd:eb:e3:b4:54:e8:e3:4c:30:72:
    2f:59:da:33:0e:4c:95:22:21:70:83:8c:e0:ef:ab:
    92:50:22:29:5c:cb:f4:e3:5c:d6:bc:f5:53:7f:09:
    ac:d7:94:8c:3e:a9:12:57:83:70:f3:c1:47:58:1e:
    f0:37
publicExponent: 65537 (0x10001)
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
	SelfSignedECDSACertificateNoDigitalSignatureKeyUsage = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            a5:62:ef:bb:5e:6d:f0:4f
        Signature Algorithm: ecdsa-with-SHA256
        Issuer: C = US, ST = CO, L = Apache Traffic Control, O = Traffic Ops, OU = Unit Testing, CN = *.invalid2.invalid
        Validity
            Not Before: Mar  6 17:37:23 2019 GMT
            Not After : Mar  1 17:37:23 2039 GMT
        Subject: C = US, ST = CO, L = Apache Traffic Control, O = Traffic Ops, OU = Unit Testing, CN = *.invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub:
                    04:83:a9:9d:f6:0e:6b:6f:85:c6:1c:65:f7:de:80:
                    e5:cb:d4:39:4a:e9:63:5f:e8:cb:ab:94:74:19:a0:
                    a3:44:5f:6a:5f:65:bf:93:62:31:18:60:83:bb:3e:
                    fa:3a:f6:c9:ff:25:b3:6c:cb:22:8b:0d:3d:28:48:
                    54:b2:65:f5:27
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
            X509v3 Subject Key Identifier: 
                72:37:53:33:5D:B6:BF:B5:4A:EE:05:97:7B:10:4B:83:ED:5F:0A:ED
            X509v3 Authority Key Identifier: 
                keyid:72:37:53:33:5D:B6:BF:B5:4A:EE:05:97:7B:10:4B:83:ED:5F:0A:ED

            X509v3 Subject Alternative Name: 
                DNS:*.invalid2.invalid
    Signature Algorithm: ecdsa-with-SHA256
         30:45:02:20:3c:9b:5e:d5:92:4c:87:b2:11:38:0b:04:68:ad:
         bd:e5:80:62:1f:1d:98:12:bf:5a:dd:fe:25:d8:ff:e4:61:7a:
         02:21:00:f7:92:3c:f7:af:dc:5c:39:9d:b8:29:48:83:80:a7:
         60:b1:9b:90:33:ec:0b:ed:e7:10:d1:08:af:dd:96:aa:01
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
	SelfSignedECDSAPrivateKeyNoDigitalSignatureKeyUsage = `
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
Private-Key: (256 bit)
priv:
    fa:e3:75:29:92:a7:a0:ff:33:a2:81:3b:0c:33:ee:
    ef:cd:cc:ee:f4:9c:14:f4:53:be:10:48:20:56:9b:
    00:0e
pub:
    04:83:a9:9d:f6:0e:6b:6f:85:c6:1c:65:f7:de:80:
    e5:cb:d4:39:4a:e9:63:5f:e8:cb:ab:94:74:19:a0:
    a3:44:5f:6a:5f:65:bf:93:62:31:18:60:83:bb:3e:
    fa:3a:f6:c9:ff:25:b3:6c:cb:22:8b:0d:3d:28:48:
    54:b2:65:f5:27
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPrjdSmSp6D/M6KBOwwz7u/NzO70nBT0U74QSCBWmwAOoAoGCCqGSM49
AwEHoUQDQgAEg6md9g5rb4XGHGX33oDly9Q5SuljX+jLq5R0GaCjRF9qX2W/k2Ix
GGCDuz76OvbJ/yWzbMsiiw09KEhUsmX1Jw==
-----END EC PRIVATE KEY-----
`
	SelfSignedECDSACertificate = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            28:02:47:96:16:5c:d6:80:99:4a:11:92:22:04:cb:ff:a3:a8:f0:ed
        Signature Algorithm: ecdsa-with-SHA256
        Issuer: C = US, ST = Colorado, L = Denver, O = Apache Traffic Control, OU = Apache Traffic Control, CN = *.ecdsatest.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Validity
            Not Before: May  3 16:46:33 2019 GMT
            Not After : Apr 28 16:46:33 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Apache Traffic Control, OU = Apache Traffic Control, CN = *.ecdsatest.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.tes
t
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub:
                    04:1d:54:a5:e0:eb:e0:b0:a8:c5:d6:30:ce:20:40:
                    5a:93:49:19:e5:0a:de:29:1f:a8:79:fd:74:5a:ff:
                    a7:d2:37:8f:6e:32:a3:d9:49:72:f8:c6:07:b2:cf:
                    d1:eb:8a:af:44:23:9e:c8:26:d5:ca:ba:c5:1d:a6:
                    78:51:1f:ba:d1
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:TRUE
            Netscape Cert Type: 
                SSL Server
            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign
            X509v3 Extended Key Usage: critical
                TLS Web Server Authentication, Code Signing
            X509v3 Subject Key Identifier: 
                01:2B:30:40:BE:CF:04:31:00:59:5C:92:CB:FE:8A:DD:40:6B:3F:CA
            X509v3 Subject Alternative Name: 
                DNS:*.ecdsatest.mycdn.ciab.test
    Signature Algorithm: ecdsa-with-SHA256
         30:45:02:20:4d:a9:99:a1:af:79:03:65:d1:f2:8d:4d:ae:ca:
         01:1b:8f:1d:ef:36:7e:14:36:b6:ea:5c:ec:f5:93:64:7e:aa:
         02:21:00:dc:01:04:c4:b8:1e:e6:b8:34:62:79:48:71:c6:c8:
         c4:21:ff:e6:6f:3a:24:5e:90:c7:c9:97:91:47:92:cc:2d
-----BEGIN CERTIFICATE-----
MIIDJjCCAsygAwIBAgIUKAJHlhZc1oCZShGSIgTL/6Oo8O0wCgYIKoZIzj0EAwIw
gcIxCzAJBgNVBAYTAlVTMREwDwYDVQQIEwhDb2xvcmFkbzEPMA0GA1UEBxMGRGVu
dmVyMR8wHQYDVQQKExZBcGFjaGUgVHJhZmZpYyBDb250cm9sMR8wHQYDVQQLExZB
cGFjaGUgVHJhZmZpYyBDb250cm9sMSQwIgYDVQQDFBsqLmVjZHNhdGVzdC5teWNk
bi5jaWFiLnRlc3QxJzAlBgkqhkiG9w0BCQEWGG5vLXJlcGx5QGluZnJhLmNpYWIu
dGVzdDAeFw0xOTA1MDMxNjQ2MzNaFw0zOTA0MjgxNjQ2MzNaMIHCMQswCQYDVQQG
EwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZlcjEfMB0GA1UE
ChMWQXBhY2hlIFRyYWZmaWMgQ29udHJvbDEfMB0GA1UECxMWQXBhY2hlIFRyYWZm
aWMgQ29udHJvbDEkMCIGA1UEAxQbKi5lY2RzYXRlc3QubXljZG4uY2lhYi50ZXN0
MScwJQYJKoZIhvcNAQkBFhhuby1yZXBseUBpbmZyYS5jaWFiLnRlc3QwWTATBgcq
hkjOPQIBBggqhkjOPQMBBwNCAAQdVKXg6+CwqMXWMM4gQFqTSRnlCt4pH6h5/XRa
/6fSN49uMqPZSXL4xgeyz9Hriq9EI57IJtXKusUdpnhRH7rRo4GdMIGaMAwGA1Ud
EwQFMAMBAf8wEQYJYIZIAYb4QgEBBAQDAgZAMA4GA1UdDwEB/wQEAwIChDAgBgNV
HSUBAf8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwMwHQYDVR0OBBYEFAErMEC+zwQx
AFlcksv+it1Aaz/KMCYGA1UdEQQfMB2CGyouZWNkc2F0ZXN0Lm15Y2RuLmNpYWIu
dGVzdDAKBggqhkjOPQQDAgNIADBFAiBNqZmhr3kDZdHyjU2uygEbjx3vNn4UNrbq
XOz1k2R+qgIhANwBBMS4Hua4NGJ5SHHGyMQh/+ZvOiRekMfJl5FHkswt
-----END CERTIFICATE-----
`
	SelfSignedECDSAPrivateKey = `
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
Private-Key: (256 bit)
priv:
    d5:5b:f3:26:9a:09:3d:78:09:58:58:55:a3:af:dd:
    f9:86:a4:45:6e:5e:f4:dc:94:e4:18:21:9a:bb:f7:
    7b:2f
pub:
    04:1d:54:a5:e0:eb:e0:b0:a8:c5:d6:30:ce:20:40:
    5a:93:49:19:e5:0a:de:29:1f:a8:79:fd:74:5a:ff:
    a7:d2:37:8f:6e:32:a3:d9:49:72:f8:c6:07:b2:cf:
    d1:eb:8a:af:44:23:9e:c8:26:d5:ca:ba:c5:1d:a6:
    78:51:1f:ba:d1
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEINVb8yaaCT14CVhYVaOv3fmGpEVuXvTclOQYIZq793svoAoGCCqGSM49
AwEHoUQDQgAEHVSl4OvgsKjF1jDOIEBak0kZ5QreKR+oef10Wv+n0jePbjKj2Uly
+MYHss/R64qvRCOeyCbVyrrFHaZ4UR+60Q==
-----END EC PRIVATE KEY-----
`
	SelfSignedECDSAPrivateKeyWithoutParams = `
Private-Key: (256 bit)
priv:
    d5:5b:f3:26:9a:09:3d:78:09:58:58:55:a3:af:dd:
    f9:86:a4:45:6e:5e:f4:dc:94:e4:18:21:9a:bb:f7:
    7b:2f
pub:
    04:1d:54:a5:e0:eb:e0:b0:a8:c5:d6:30:ce:20:40:
    5a:93:49:19:e5:0a:de:29:1f:a8:79:fd:74:5a:ff:
    a7:d2:37:8f:6e:32:a3:d9:49:72:f8:c6:07:b2:cf:
    d1:eb:8a:af:44:23:9e:c8:26:d5:ca:ba:c5:1d:a6:
    78:51:1f:ba:d1
ASN1 OID: prime256v1
NIST CURVE: P-256
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEINVb8yaaCT14CVhYVaOv3fmGpEVuXvTclOQYIZq793svoAoGCCqGSM49
AwEHoUQDQgAEHVSl4OvgsKjF1jDOIEBak0kZ5QreKR+oef10Wv+n0jePbjKj2Uly
+MYHss/R64qvRCOeyCbVyrrFHaZ4UR+60Q==
-----END EC PRIVATE KEY-----
`
	SelfSignedDSACertificate = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            3f:73:c5:db:de:d5:a0:c2:27:6b:b4:22:35:5c:3d:d7:10:46:b4:80
        Signature Algorithm: dsa_with_SHA256
        Issuer: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.dsatest.mycdn.ciab.test, emailAddress = no-reply@i
nfra.ciab.test
        Validity
            Not Before: May  3 17:24:13 2019 GMT
            Not After : Apr 28 17:24:13 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.dsatest.mycdn.ciab.test, emailAddress = no-reply@
infra.ciab.test
        Subject Public Key Info:
            Public Key Algorithm: dsaEncryption
                pub: 
                    3c:fc:6a:a1:a4:a7:3a:01:7d:6c:b2:6a:35:11:07:
                    0c:59:0a:3e:a1:de:cb:6d:51:17:61:1c:22:fc:be:
                    81:53:70:8f:50:81:d2:cf:86:f0:02:dc:e5:c0:3b:
                    d2:2c:8f:13:6e:9d:50:ce:cd:95:06:47:13:e0:a6:
                    ce:88:84:70:e3:10:2c:b9:6f:a1:1c:12:5a:b2:dd:
                    9c:03:b8:2b:08:fd:9c:d1:a4:de:f3:51:54:2b:3e:
                    ca:32:81:f8:5e:ef:dc:62:07:7d:4c:e0:af:2d:44:
                    f9:7e:a2:73:8c:8f:12:93:9c:ff:78:9d:2a:ff:79:
                    0e:e7:b1:83:36:02:e2:fd:61:c5:63:d0:1b:9f:bc:
                    9b:e8:ca:64:35:ce:77:9a:2e:a0:94:c4:5c:d0:3d:
                    1d:fd:6e:4c:ad:22:2b:ab:fc:3e:ee:5c:0c:eb:ce:
                    69:43:d4:db:3d:1d:17:a2:29:a6:e2:bb:ab:e3:fc:
                    3d:d8:ff:59:c9:d3:29:f6:e0:de:8a:3c:25:fa:a1:
                    62:50:1c:e7:7b:ae:88:5b:1e:46:1f:cf:05:42:db:
                    96:2f:87:03:bb:27:f1:44:fb:38:13:6b:c2:de:99:
                    37:ef:38:b5:6f:aa:89:d6:99:9c:61:26:d4:28:c6:
                    c4:35:ba:8b:b9:6f:63:7c:d7:62:b1:55:10:07:23:
                    d0
                P:   
                    00:93:2a:40:af:44:b3:c0:d7:04:28:39:b6:92:20:
                    31:ee:3a:1a:cd:00:54:d6:d8:6b:f1:e2:14:64:08:
                    8c:a5:29:e9:50:7e:ac:bc:23:d4:8a:56:23:aa:14:
                    15:21:1d:85:e4:f7:a4:87:f3:5c:9e:e1:86:1e:46:
                    72:08:a1:9b:4b:dc:4a:66:79:dd:67:65:78:4c:0d:
                    c7:0f:2b:d3:06:ff:23:9d:f4:9f:6a:1e:d0:60:3e:
                    fb:13:98:ee:f4:0d:39:d0:a5:71:4d:9a:58:84:d7:
                    db:a1:1b:f4:3d:68:35:c7:f5:82:ce:a5:ad:61:63:
                    cf:b3:b0:c5:10:13:e1:29:97:0c:45:39:a7:10:bc:
                    1b:33:de:29:3c:8c:85:10:40:c9:2c:be:74:93:5b:
                    2a:6c:ad:f9:98:e7:3e:00:f0:e2:a2:5c:b5:13:25:
                    06:5e:3f:fc:aa:c0:a7:14:3b:fc:10:a1:fa:50:6e:
                    5b:3b:d4:35:7c:61:0f:4e:a8:02:71:02:96:16:b6:
                    dd:73:a5:18:64:ab:be:81:27:0c:13:26:b1:e1:89:
                    2d:8a:b4:c4:79:7e:3a:a3:75:67:18:a5:d4:c2:c9:
                    44:66:18:fe:1d:d0:27:84:9c:ca:39:4b:37:c5:dc:
                    33:a9:58:d9:e3:1c:61:ea:86:37:a0:4b:b1:6a:b7:
                    17:9b
                Q:   
                    00:c6:22:da:59:95:f9:be:1a:3f:a7:2e:ce:2d:dd:
                    1c:42:b3:9b:f6:7b:d3:24:22:e4:5f:b0:e5:2f:e0:
                    d8:8b:f1
                G:   
                    45:01:d7:ea:0f:73:9f:70:b0:d8:f1:ec:d8:66:97:
                    15:f1:21:1e:47:c9:af:ba:28:79:4e:3f:f9:4a:d7:
                    fb:3f:91:c0:b3:0e:f7:fb:b8:03:05:7a:98:0a:2f:
                    68:1e:03:5c:b7:e9:f9:65:18:f2:fd:c0:89:85:87:
                    db:99:a6:bb:4a:58:0d:67:e9:7b:75:6d:3b:2b:7e:
                    81:1d:ad:88:58:ca:8e:73:4a:58:86:56:e8:99:e9:
                    52:c0:dd:e8:1b:f1:60:28:c8:67:bc:8e:cf:10:77:
                    e3:56:90:9a:e3:a1:07:e7:43:29:8b:bf:64:3d:ac:
                    56:29:cd:ce:f7:39:53:2e:34:97:d8:99:b9:02:3b:
                    0a:b1:e0:dc:57:b9:1e:35:b5:5e:28:e1:61:69:a7:
                    ed:a1:e1:5b:44:4b:af:45:75:3f:f2:18:01:60:62:
                    6e:b6:aa:7a:a0:8d:3b:d1:eb:ff:0f:8b:49:a3:8e:
                    9b:3a:44:19:4a:e8:db:9c:2b:06:dd:d8:4d:f0:23:
                    fd:a5:21:f2:26:11:a4:40:34:c8:5b:a5:46:9a:82:
                    ba:ea:4e:c0:c1:b3:14:72:90:a4:b5:51:d8:f2:66:
                    78:95:e1:33:95:e1:78:dd:73:44:45:f3:0f:49:6e:
                    f4:72:29:b0:3e:28:08:2f:f6:9e:2a:1b:17:b1:ac:
                    37
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:TRUE
            Netscape Cert Type: 
                SSL Server
            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign
            X509v3 Extended Key Usage: critical
                TLS Web Server Authentication, Code Signing
            X509v3 Subject Key Identifier: 
                9E:54:E5:D4:AA:5F:47:94:3B:C5:6E:2E:A4:98:46:DB:D0:6C:CB:3F
            X509v3 Subject Alternative Name: 
                DNS:*.dsatest.mycdn.ciab.test
    Signature Algorithm: dsa_with_SHA256
         r:   
             47:70:0c:3a:7e:cb:82:51:69:4d:cd:a2:32:9d:17:
             91:b7:da:d3:30:65:25:c0:c5:04:53:42:07:67:49:
             81:1b
         s:   
             23:bb:51:04:46:5e:e3:64:32:93:09:b4:0a:24:11:
             28:59:d6:61:46:62:33:b9:18:32:35:a2:f3:17:2b:
             21:24
-----BEGIN CERTIFICATE-----
MIIF6DCCBY6gAwIBAgIUP3PF297VoMIna7QiNVw91xBGtIAwCwYJYIZIAWUDBAMC
MIGsMQswCQYDVQQGEwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRl
bnZlcjEVMBMGA1UEChMMQ0ROLWluLWEtQm94MRUwEwYDVQQLEwxDRE4taW4tYS1C
b3gxIjAgBgNVBAMUGSouZHNhdGVzdC5teWNkbi5jaWFiLnRlc3QxJzAlBgkqhkiG
9w0BCQEWGG5vLXJlcGx5QGluZnJhLmNpYWIudGVzdDAeFw0xOTA1MDMxNzI0MTNa
Fw0zOTA0MjgxNzI0MTNaMIGsMQswCQYDVQQGEwJVUzERMA8GA1UECBMIQ29sb3Jh
ZG8xDzANBgNVBAcTBkRlbnZlcjEVMBMGA1UEChMMQ0ROLWluLWEtQm94MRUwEwYD
VQQLEwxDRE4taW4tYS1Cb3gxIjAgBgNVBAMUGSouZHNhdGVzdC5teWNkbi5jaWFi
LnRlc3QxJzAlBgkqhkiG9w0BCQEWGG5vLXJlcGx5QGluZnJhLmNpYWIudGVzdDCC
A0YwggI5BgcqhkjOOAQBMIICLAKCAQEAkypAr0SzwNcEKDm2kiAx7joazQBU1thr
8eIUZAiMpSnpUH6svCPUilYjqhQVIR2F5Pekh/NcnuGGHkZyCKGbS9xKZnndZ2V4
TA3HDyvTBv8jnfSfah7QYD77E5ju9A050KVxTZpYhNfboRv0PWg1x/WCzqWtYWPP
s7DFEBPhKZcMRTmnELwbM94pPIyFEEDJLL50k1sqbK35mOc+APDioly1EyUGXj/8
qsCnFDv8EKH6UG5bO9Q1fGEPTqgCcQKWFrbdc6UYZKu+gScMEyax4YktirTEeX46
o3VnGKXUwslEZhj+HdAnhJzKOUs3xdwzqVjZ4xxh6oY3oEuxarcXmwIhAMYi2lmV
+b4aP6cuzi3dHEKzm/Z70yQi5F+w5S/g2IvxAoIBAEUB1+oPc59wsNjx7NhmlxXx
IR5Hya+6KHlOP/lK1/s/kcCzDvf7uAMFepgKL2geA1y36fllGPL9wImFh9uZprtK
WA1n6Xt1bTsrfoEdrYhYyo5zSliGVuiZ6VLA3egb8WAoyGe8js8Qd+NWkJrjoQfn
QymLv2Q9rFYpzc73OVMuNJfYmbkCOwqx4NxXuR41tV4o4WFpp+2h4VtES69FdT/y
GAFgYm62qnqgjTvR6/8Pi0mjjps6RBlK6NucKwbd2E3wI/2lIfImEaRANMhbpUaa
grrqTsDBsxRykKS1UdjyZniV4TOV4Xjdc0RF8w9JbvRyKbA+KAgv9p4qGxexrDcD
ggEFAAKCAQA8/GqhpKc6AX1ssmo1EQcMWQo+od7LbVEXYRwi/L6BU3CPUIHSz4bw
AtzlwDvSLI8Tbp1Qzs2VBkcT4KbOiIRw4xAsuW+hHBJast2cA7grCP2c0aTe81FU
Kz7KMoH4Xu/cYgd9TOCvLUT5fqJzjI8Sk5z/eJ0q/3kO57GDNgLi/WHFY9Abn7yb
6MpkNc53mi6glMRc0D0d/W5MrSIrq/w+7lwM685pQ9TbPR0Xoimm4rur4/w92P9Z
ydMp9uDeijwl+qFiUBzne66IWx5GH88FQtuWL4cDuyfxRPs4E2vC3pk37zi1b6qJ
1pmcYSbUKMbENbqLuW9jfNdisVUQByPQo4GbMIGYMAwGA1UdEwQFMAMBAf8wEQYJ
YIZIAYb4QgEBBAQDAgZAMA4GA1UdDwEB/wQEAwIChDAgBgNVHSUBAf8EFjAUBggr
BgEFBQcDAQYIKwYBBQUHAwMwHQYDVR0OBBYEFJ5U5dSqX0eUO8VuLqSYRtvQbMs/
MCQGA1UdEQQdMBuCGSouZHNhdGVzdC5teWNkbi5jaWFiLnRlc3QwCwYJYIZIAWUD
BAMCA0cAMEQCIEdwDDp+y4JRaU3NojKdF5G32tMwZSXAxQRTQgdnSYEbAiAju1EE
Rl7jZDKTCbQKJBEoWdZhRmIzuRgyNaLzFyshJA==
-----END CERTIFICATE-----
`
	SelfSignedDSAPrivateKey = `
    P:   
        00:93:2a:40:af:44:b3:c0:d7:04:28:39:b6:92:20:
        31:ee:3a:1a:cd:00:54:d6:d8:6b:f1:e2:14:64:08:
        8c:a5:29:e9:50:7e:ac:bc:23:d4:8a:56:23:aa:14:
        15:21:1d:85:e4:f7:a4:87:f3:5c:9e:e1:86:1e:46:
        72:08:a1:9b:4b:dc:4a:66:79:dd:67:65:78:4c:0d:
        c7:0f:2b:d3:06:ff:23:9d:f4:9f:6a:1e:d0:60:3e:
        fb:13:98:ee:f4:0d:39:d0:a5:71:4d:9a:58:84:d7:
        db:a1:1b:f4:3d:68:35:c7:f5:82:ce:a5:ad:61:63:
        cf:b3:b0:c5:10:13:e1:29:97:0c:45:39:a7:10:bc:
        1b:33:de:29:3c:8c:85:10:40:c9:2c:be:74:93:5b:
        2a:6c:ad:f9:98:e7:3e:00:f0:e2:a2:5c:b5:13:25:
        06:5e:3f:fc:aa:c0:a7:14:3b:fc:10:a1:fa:50:6e:
        5b:3b:d4:35:7c:61:0f:4e:a8:02:71:02:96:16:b6:
        dd:73:a5:18:64:ab:be:81:27:0c:13:26:b1:e1:89:
        2d:8a:b4:c4:79:7e:3a:a3:75:67:18:a5:d4:c2:c9:
        44:66:18:fe:1d:d0:27:84:9c:ca:39:4b:37:c5:dc:
        33:a9:58:d9:e3:1c:61:ea:86:37:a0:4b:b1:6a:b7:
        17:9b
    Q:   
        00:c6:22:da:59:95:f9:be:1a:3f:a7:2e:ce:2d:dd:
        1c:42:b3:9b:f6:7b:d3:24:22:e4:5f:b0:e5:2f:e0:
        d8:8b:f1
    G:   
        45:01:d7:ea:0f:73:9f:70:b0:d8:f1:ec:d8:66:97:
        15:f1:21:1e:47:c9:af:ba:28:79:4e:3f:f9:4a:d7:
        fb:3f:91:c0:b3:0e:f7:fb:b8:03:05:7a:98:0a:2f:
        68:1e:03:5c:b7:e9:f9:65:18:f2:fd:c0:89:85:87:
        db:99:a6:bb:4a:58:0d:67:e9:7b:75:6d:3b:2b:7e:
        81:1d:ad:88:58:ca:8e:73:4a:58:86:56:e8:99:e9:
        52:c0:dd:e8:1b:f1:60:28:c8:67:bc:8e:cf:10:77:
        e3:56:90:9a:e3:a1:07:e7:43:29:8b:bf:64:3d:ac:
        56:29:cd:ce:f7:39:53:2e:34:97:d8:99:b9:02:3b:
        0a:b1:e0:dc:57:b9:1e:35:b5:5e:28:e1:61:69:a7:
        ed:a1:e1:5b:44:4b:af:45:75:3f:f2:18:01:60:62:
        6e:b6:aa:7a:a0:8d:3b:d1:eb:ff:0f:8b:49:a3:8e:
        9b:3a:44:19:4a:e8:db:9c:2b:06:dd:d8:4d:f0:23:
        fd:a5:21:f2:26:11:a4:40:34:c8:5b:a5:46:9a:82:
        ba:ea:4e:c0:c1:b3:14:72:90:a4:b5:51:d8:f2:66:
        78:95:e1:33:95:e1:78:dd:73:44:45:f3:0f:49:6e:
        f4:72:29:b0:3e:28:08:2f:f6:9e:2a:1b:17:b1:ac:
        37
-----BEGIN DSA PARAMETERS-----
MIICLAKCAQEAkypAr0SzwNcEKDm2kiAx7joazQBU1thr8eIUZAiMpSnpUH6svCPU
ilYjqhQVIR2F5Pekh/NcnuGGHkZyCKGbS9xKZnndZ2V4TA3HDyvTBv8jnfSfah7Q
YD77E5ju9A050KVxTZpYhNfboRv0PWg1x/WCzqWtYWPPs7DFEBPhKZcMRTmnELwb
M94pPIyFEEDJLL50k1sqbK35mOc+APDioly1EyUGXj/8qsCnFDv8EKH6UG5bO9Q1
fGEPTqgCcQKWFrbdc6UYZKu+gScMEyax4YktirTEeX46o3VnGKXUwslEZhj+HdAn
hJzKOUs3xdwzqVjZ4xxh6oY3oEuxarcXmwIhAMYi2lmV+b4aP6cuzi3dHEKzm/Z7
0yQi5F+w5S/g2IvxAoIBAEUB1+oPc59wsNjx7NhmlxXxIR5Hya+6KHlOP/lK1/s/
kcCzDvf7uAMFepgKL2geA1y36fllGPL9wImFh9uZprtKWA1n6Xt1bTsrfoEdrYhY
yo5zSliGVuiZ6VLA3egb8WAoyGe8js8Qd+NWkJrjoQfnQymLv2Q9rFYpzc73OVMu
NJfYmbkCOwqx4NxXuR41tV4o4WFpp+2h4VtES69FdT/yGAFgYm62qnqgjTvR6/8P
i0mjjps6RBlK6NucKwbd2E3wI/2lIfImEaRANMhbpUaagrrqTsDBsxRykKS1Udjy
ZniV4TOV4Xjdc0RF8w9JbvRyKbA+KAgv9p4qGxexrDc=
-----END DSA PARAMETERS-----
Private-Key: (2048 bit)
priv:
    40:62:c4:5a:b2:83:e5:2c:27:15:f8:46:e5:03:4e:
    8c:53:ba:44:df:db:34:1f:12:6a:75:eb:9f:30:d4:
    9e:41
pub: 
    3c:fc:6a:a1:a4:a7:3a:01:7d:6c:b2:6a:35:11:07:
    0c:59:0a:3e:a1:de:cb:6d:51:17:61:1c:22:fc:be:
    81:53:70:8f:50:81:d2:cf:86:f0:02:dc:e5:c0:3b:
    d2:2c:8f:13:6e:9d:50:ce:cd:95:06:47:13:e0:a6:
    ce:88:84:70:e3:10:2c:b9:6f:a1:1c:12:5a:b2:dd:
    9c:03:b8:2b:08:fd:9c:d1:a4:de:f3:51:54:2b:3e:
    ca:32:81:f8:5e:ef:dc:62:07:7d:4c:e0:af:2d:44:
    f9:7e:a2:73:8c:8f:12:93:9c:ff:78:9d:2a:ff:79:
    0e:e7:b1:83:36:02:e2:fd:61:c5:63:d0:1b:9f:bc:
    9b:e8:ca:64:35:ce:77:9a:2e:a0:94:c4:5c:d0:3d:
    1d:fd:6e:4c:ad:22:2b:ab:fc:3e:ee:5c:0c:eb:ce:
    69:43:d4:db:3d:1d:17:a2:29:a6:e2:bb:ab:e3:fc:
    3d:d8:ff:59:c9:d3:29:f6:e0:de:8a:3c:25:fa:a1:
    62:50:1c:e7:7b:ae:88:5b:1e:46:1f:cf:05:42:db:
    96:2f:87:03:bb:27:f1:44:fb:38:13:6b:c2:de:99:
    37:ef:38:b5:6f:aa:89:d6:99:9c:61:26:d4:28:c6:
    c4:35:ba:8b:b9:6f:63:7c:d7:62:b1:55:10:07:23:
    d0
P:   
    00:93:2a:40:af:44:b3:c0:d7:04:28:39:b6:92:20:
    31:ee:3a:1a:cd:00:54:d6:d8:6b:f1:e2:14:64:08:
    8c:a5:29:e9:50:7e:ac:bc:23:d4:8a:56:23:aa:14:
    15:21:1d:85:e4:f7:a4:87:f3:5c:9e:e1:86:1e:46:
    72:08:a1:9b:4b:dc:4a:66:79:dd:67:65:78:4c:0d:
    c7:0f:2b:d3:06:ff:23:9d:f4:9f:6a:1e:d0:60:3e:
    fb:13:98:ee:f4:0d:39:d0:a5:71:4d:9a:58:84:d7:
    db:a1:1b:f4:3d:68:35:c7:f5:82:ce:a5:ad:61:63:
    cf:b3:b0:c5:10:13:e1:29:97:0c:45:39:a7:10:bc:
    1b:33:de:29:3c:8c:85:10:40:c9:2c:be:74:93:5b:
    2a:6c:ad:f9:98:e7:3e:00:f0:e2:a2:5c:b5:13:25:
    06:5e:3f:fc:aa:c0:a7:14:3b:fc:10:a1:fa:50:6e:
    5b:3b:d4:35:7c:61:0f:4e:a8:02:71:02:96:16:b6:
    dd:73:a5:18:64:ab:be:81:27:0c:13:26:b1:e1:89:
    2d:8a:b4:c4:79:7e:3a:a3:75:67:18:a5:d4:c2:c9:
    44:66:18:fe:1d:d0:27:84:9c:ca:39:4b:37:c5:dc:
    33:a9:58:d9:e3:1c:61:ea:86:37:a0:4b:b1:6a:b7:
    17:9b
Q:   
    00:c6:22:da:59:95:f9:be:1a:3f:a7:2e:ce:2d:dd:
    1c:42:b3:9b:f6:7b:d3:24:22:e4:5f:b0:e5:2f:e0:
    d8:8b:f1
G:   
    45:01:d7:ea:0f:73:9f:70:b0:d8:f1:ec:d8:66:97:
    15:f1:21:1e:47:c9:af:ba:28:79:4e:3f:f9:4a:d7:
    fb:3f:91:c0:b3:0e:f7:fb:b8:03:05:7a:98:0a:2f:
    68:1e:03:5c:b7:e9:f9:65:18:f2:fd:c0:89:85:87:
    db:99:a6:bb:4a:58:0d:67:e9:7b:75:6d:3b:2b:7e:
    81:1d:ad:88:58:ca:8e:73:4a:58:86:56:e8:99:e9:
    52:c0:dd:e8:1b:f1:60:28:c8:67:bc:8e:cf:10:77:
    e3:56:90:9a:e3:a1:07:e7:43:29:8b:bf:64:3d:ac:
    56:29:cd:ce:f7:39:53:2e:34:97:d8:99:b9:02:3b:
    0a:b1:e0:dc:57:b9:1e:35:b5:5e:28:e1:61:69:a7:
    ed:a1:e1:5b:44:4b:af:45:75:3f:f2:18:01:60:62:
    6e:b6:aa:7a:a0:8d:3b:d1:eb:ff:0f:8b:49:a3:8e:
    9b:3a:44:19:4a:e8:db:9c:2b:06:dd:d8:4d:f0:23:
    fd:a5:21:f2:26:11:a4:40:34:c8:5b:a5:46:9a:82:
    ba:ea:4e:c0:c1:b3:14:72:90:a4:b5:51:d8:f2:66:
    78:95:e1:33:95:e1:78:dd:73:44:45:f3:0f:49:6e:
    f4:72:29:b0:3e:28:08:2f:f6:9e:2a:1b:17:b1:ac:
    37
-----BEGIN DSA PRIVATE KEY-----
MIIDVQIBAAKCAQEAkypAr0SzwNcEKDm2kiAx7joazQBU1thr8eIUZAiMpSnpUH6s
vCPUilYjqhQVIR2F5Pekh/NcnuGGHkZyCKGbS9xKZnndZ2V4TA3HDyvTBv8jnfSf
ah7QYD77E5ju9A050KVxTZpYhNfboRv0PWg1x/WCzqWtYWPPs7DFEBPhKZcMRTmn
ELwbM94pPIyFEEDJLL50k1sqbK35mOc+APDioly1EyUGXj/8qsCnFDv8EKH6UG5b
O9Q1fGEPTqgCcQKWFrbdc6UYZKu+gScMEyax4YktirTEeX46o3VnGKXUwslEZhj+
HdAnhJzKOUs3xdwzqVjZ4xxh6oY3oEuxarcXmwIhAMYi2lmV+b4aP6cuzi3dHEKz
m/Z70yQi5F+w5S/g2IvxAoIBAEUB1+oPc59wsNjx7NhmlxXxIR5Hya+6KHlOP/lK
1/s/kcCzDvf7uAMFepgKL2geA1y36fllGPL9wImFh9uZprtKWA1n6Xt1bTsrfoEd
rYhYyo5zSliGVuiZ6VLA3egb8WAoyGe8js8Qd+NWkJrjoQfnQymLv2Q9rFYpzc73
OVMuNJfYmbkCOwqx4NxXuR41tV4o4WFpp+2h4VtES69FdT/yGAFgYm62qnqgjTvR
6/8Pi0mjjps6RBlK6NucKwbd2E3wI/2lIfImEaRANMhbpUaagrrqTsDBsxRykKS1
UdjyZniV4TOV4Xjdc0RF8w9JbvRyKbA+KAgv9p4qGxexrDcCggEAPPxqoaSnOgF9
bLJqNREHDFkKPqHey21RF2EcIvy+gVNwj1CB0s+G8ALc5cA70iyPE26dUM7NlQZH
E+CmzoiEcOMQLLlvoRwSWrLdnAO4Kwj9nNGk3vNRVCs+yjKB+F7v3GIHfUzgry1E
+X6ic4yPEpOc/3idKv95DuexgzYC4v1hxWPQG5+8m+jKZDXOd5ouoJTEXNA9Hf1u
TK0iK6v8Pu5cDOvOaUPU2z0dF6IppuK7q+P8Pdj/WcnTKfbg3oo8JfqhYlAc53uu
iFseRh/PBULbli+HA7sn8UT7OBNrwt6ZN+84tW+qidaZnGEm1CjGxDW6i7lvY3zX
YrFVEAcj0AIgQGLEWrKD5SwnFfhG5QNOjFO6RN/bNB8SanXrnzDUnkE=
-----END DSA PRIVATE KEY-----
`
	SelfSignedX509v1Certificate = `
Certificate:
    Data:
        Version: 1 (0x0)
        Serial Number:
            3d:e7:13:5a:25:e3:d1:0f:80:f8:f9:93:76:2e:93:2c:59:a3:51:ae
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.x509v1-test.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Validity
            Not Before: May  3 16:18:30 2019 GMT
            Not After : Apr 28 16:18:30 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = CDN-in-a-Box, OU = CDN-in-a-Box, CN = *.x509v1-test.mycdn.ciab.test, emailAddress = no-reply@infra.ciab.test
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:cd:73:2d:46:66:78:49:56:29:3b:8c:9f:d3:49:
                    0a:24:3e:cc:28:25:73:8d:33:8a:24:d8:8e:c8:99:
                    fe:f9:16:0d:5f:ae:20:65:4e:da:f0:b4:5e:7e:fa:
                    ee:bf:ac:0e:d4:87:38:21:1e:f1:bd:9e:7c:b7:fb:
                    20:1d:ee:80:a0:75:b7:53:3f:57:27:89:27:67:50:
                    47:34:51:d8:b1:a1:0b:48:35:24:3c:d4:68:c1:0a:
                    6b:cc:99:80:7f:7c:83:ae:7e:a4:5a:b2:d5:81:7d:
                    22:5a:b4:21:4f:95:09:c6:04:bb:2f:58:8d:89:ba:
                    80:28:51:a4:af:4b:51:32:4d:2b:c2:30:73:9b:90:
                    86:1b:2b:75:59:2f:c9:d5:6d:e1:fd:1d:39:05:a0:
                    d7:ba:44:82:a1:29:68:73:1d:57:82:69:5d:09:01:
                    36:ed:d7:be:59:7b:ac:21:36:6c:7b:98:7e:1a:27:
                    a2:fe:11:a4:1d:d2:49:af:56:91:38:fd:f2:f9:cc:
                    de:38:5d:58:3a:bd:b8:a7:fd:13:1c:fa:e2:af:48:
                    9e:b1:65:76:f9:b0:3f:4d:d8:f9:85:55:a2:a6:dd:
                    4e:24:1c:0e:bd:fb:a5:97:d7:d9:4e:f0:b5:45:17:
                    61:23:25:09:6a:e5:ab:e2:2c:03:da:0a:67:50:8f:
                    0e:cb
                Exponent: 65537 (0x10001)
-----BEGIN CERTIFICATE-----
MIID6TCCAtECFD3nE1ol49EPgPj5k3YukyxZo1GuMA0GCSqGSIb3DQEBCwUAMIGw
MQswCQYDVQQGEwJVUzERMA8GA1UECBMIQ29sb3JhZG8xDzANBgNVBAcTBkRlbnZl
cjEVMBMGA1UEChMMQ0ROLWluLWEtQm94MRUwEwYDVQQLEwxDRE4taW4tYS1Cb3gx
JjAkBgNVBAMUHSoueDUwOXYxLXRlc3QubXljZG4uY2lhYi50ZXN0MScwJQYJKoZI
hvcNAQkBFhhuby1yZXBseUBpbmZyYS5jaWFiLnRlc3QwHhcNMTkwNTAzMTYxODMw
WhcNMzkwNDI4MTYxODMwWjCBsDELMAkGA1UEBhMCVVMxETAPBgNVBAgTCENvbG9y
YWRvMQ8wDQYDVQQHEwZEZW52ZXIxFTATBgNVBAoTDENETi1pbi1hLUJveDEVMBMG
A1UECxMMQ0ROLWluLWEtQm94MSYwJAYDVQQDFB0qLng1MDl2MS10ZXN0Lm15Y2Ru
LmNpYWIudGVzdDEnMCUGCSqGSIb3DQEJARYYbm8tcmVwbHlAaW5mcmEuY2lhYi50
ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzXMtRmZ4SVYpO4yf
00kKJD7MKCVzjTOKJNiOyJn++RYNX64gZU7a8LRefvruv6wO1Ic4IR7xvZ58t/sg
He6AoHW3Uz9XJ4knZ1BHNFHYsaELSDUkPNRowQprzJmAf3yDrn6kWrLVgX0iWrQh
T5UJxgS7L1iNibqAKFGkr0tRMk0rwjBzm5CGGyt1WS/J1W3h/R05BaDXukSCoSlo
cx1XgmldCQE27de+WXusITZse5h+Giei/hGkHdJJr1aROP3y+czeOF1YOr24p/0T
HPrir0iesWV2+bA/Tdj5hVWipt1OJBwOvfull9fZTvC1RRdhIyUJauWr4iwD2gpn
UI8OywIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQATLszHvJ/yhxuPEmsZN9Loinkb
X3c9q8Gmsh0jdqu7IPXD0lBHAfIXOl+oQslDjKVKwOmxo/UyS8jmMdJJZS1PAFad
CmJ1kfK/VmAtSn72wkrSQps7pBZqJGqxNAsZ3hgQFg8DuonQrPIqcnX7HyzlCp60
3TVFifTbpBkpQmTFjBZ3GO0tSL9weqjf4mEY0bc5uXBmlHFWTkjAdDPrWHytllDy
FMc81jDoUnJyyoRqXXUMsKnwiW4NubfUGl1dN8Z2p7J+mwrHBMrfUUpkTPSoUjS/
PLY+k5ShdxWLv94vuDLWpETJm8766cyWBtkrOyHFuDP9BiLooMAs2paY7LoI
-----END CERTIFICATE-----
`
	SelfSignedX509v1PrivateKey = `
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:cd:73:2d:46:66:78:49:56:29:3b:8c:9f:d3:49:
    0a:24:3e:cc:28:25:73:8d:33:8a:24:d8:8e:c8:99:
    fe:f9:16:0d:5f:ae:20:65:4e:da:f0:b4:5e:7e:fa:
    ee:bf:ac:0e:d4:87:38:21:1e:f1:bd:9e:7c:b7:fb:
    20:1d:ee:80:a0:75:b7:53:3f:57:27:89:27:67:50:
    47:34:51:d8:b1:a1:0b:48:35:24:3c:d4:68:c1:0a:
    6b:cc:99:80:7f:7c:83:ae:7e:a4:5a:b2:d5:81:7d:
    22:5a:b4:21:4f:95:09:c6:04:bb:2f:58:8d:89:ba:
    80:28:51:a4:af:4b:51:32:4d:2b:c2:30:73:9b:90:
    86:1b:2b:75:59:2f:c9:d5:6d:e1:fd:1d:39:05:a0:
    d7:ba:44:82:a1:29:68:73:1d:57:82:69:5d:09:01:
    36:ed:d7:be:59:7b:ac:21:36:6c:7b:98:7e:1a:27:
    a2:fe:11:a4:1d:d2:49:af:56:91:38:fd:f2:f9:cc:
    de:38:5d:58:3a:bd:b8:a7:fd:13:1c:fa:e2:af:48:
    9e:b1:65:76:f9:b0:3f:4d:d8:f9:85:55:a2:a6:dd:
    4e:24:1c:0e:bd:fb:a5:97:d7:d9:4e:f0:b5:45:17:
    61:23:25:09:6a:e5:ab:e2:2c:03:da:0a:67:50:8f:
    0e:cb
publicExponent: 65537 (0x10001)
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAzXMtRmZ4SVYpO4yf00kKJD7MKCVzjTOKJNiOyJn++RYNX64g
ZU7a8LRefvruv6wO1Ic4IR7xvZ58t/sgHe6AoHW3Uz9XJ4knZ1BHNFHYsaELSDUk
PNRowQprzJmAf3yDrn6kWrLVgX0iWrQhT5UJxgS7L1iNibqAKFGkr0tRMk0rwjBz
m5CGGyt1WS/J1W3h/R05BaDXukSCoSlocx1XgmldCQE27de+WXusITZse5h+Giei
/hGkHdJJr1aROP3y+czeOF1YOr24p/0THPrir0iesWV2+bA/Tdj5hVWipt1OJBwO
vfull9fZTvC1RRdhIyUJauWr4iwD2gpnUI8OywIDAQABAoIBAEFLBWyGTFwiQeBn
BLRFVi/GtWNc46hQZOro2BfwuRO4am+qCymnMfWlnRKF9TJ9IAlzH+eGyhUVNVXT
PZXFoqNcRfLzAmPSNu+il76M9G0fXVKJcQbUCqavBSt07V2W2NKv9NPOWgRZHH3v
GVcNapnADy0w22qWFvy5VblQGnH9eiLJvqFZaRqLKxvKVXtp5MUS5APhQbBTXpYl
/GLnILznEFjI3CdVT2IHwP+2/KuG+3FcX7O5fkq7TDGDTTjOveW9wBFtk73uScDR
hL2jibB6kvlTR2fU+DQK/mQFOgtt+6zOnXzzKm/9e+2L899gaur22FaHm9og5zFW
44EUL2ECgYEA/J5Fy4LqvjJJIYVuHMT2CB8WQUPePHxZ9TtiSoJg0eW40D221wP9
CqP5f8idIpDSniwtP8F71HEIrrqyIYBRoDlPbs2l4E4WVLL9+UvX4k7k+FWpnECj
gONi028QpPTvNH8+Kkf5pDBc9SraA+HIIXxURv2VqWfBEkhzs/I85icCgYEA0DNB
6AV5D2VEY/dI6YPqDuEg9mU3gCH6TGIh8SP4zVfDafAyw53tGvZ3LBtXApj3pFaB
FeZJYnpGpOB9tdQoFiBEVTjkibph+VR7qiytgaFIhIjKQCkKxtwXOzDStQWPBuDF
Lv/0o+dPFic8lNTgTh0Vc3OnydG5iqPBBWAYPL0CgYArn1UkGH5ay6ovPLBQDX8C
1gNsz8Bvp3WNUGzfuvXnKQkqBI4vQQQQM1KhS04/Ks0D/VLvAIVWoRJDwf+Co3r0
9RCPbLmpKzLV+3a59uvXq5IEhB5e2hah6iIlqrcwFQ+9e/+LI5SrUqKqv3SYWQPL
LIINJDsU3tLLSnGYcEst3wKBgQDMbEed5SHEeA36iWbRwXAjQ/D1fNRNvw7fyMrC
1isIk8+PSQTPBVU1UCIa8I0yQ7eDaFw+gGo1gxGx+an0ymbBstTlSIM8qABiqwzx
PgTubsmhOB49eQ7XymoU+A8rJlYUzsVNLIusEwWYHtZg29ORXwUc4sYwZvfipH51
JLEnkQKBgH1cU89zwxN6IoM4d7jrbt7wdMYlFRi7MOmTdTCY/7ZMet9Xoz1eddrs
DDE4j9+s72EYM27TbBaJc/o537GY6SXm4AtW9JvPRpCDXuKDt1hoSIu4bBhhT/gr
1sEJwCy19Gk/tOThG3wYKFZtMsEwn/gEKVKVUP5iweltZsfQLjHm
-----END RSA PRIVATE KEY-----
`
	CASignedRSACertificateChain = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 91562250343 (0x1551891067)
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Intermediate CA, emailAddress = no-reply@invalid2.invalid
        Validity
            Not Before: Mar  6 16:51:13 2019 GMT
            Not After : Mar  1 16:51:13 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Traffic Ops, CN = *.test.invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:c4:a0:21:18:68:49:75:9a:b1:0e:27:2e:a3:3b:
                    cf:90:5d:79:a5:6f:c9:98:01:64:0b:fa:54:8e:3d:
                    84:6f:ce:e3:d6:3a:f9:df:bd:6f:c1:90:36:29:bf:
                    3a:2c:7c:36:60:1c:d4:b5:d3:a4:4e:b5:61:5e:fa:
                    44:19:4d:3a:e6:fb:13:d3:58:c9:7f:24:8c:a7:e6:
                    46:b5:66:72:4c:01:ec:b4:3b:1c:a7:d8:f8:12:9a:
                    5c:be:fb:a7:ac:3a:22:87:09:99:bd:66:1c:76:a7:
                    bf:dd:fd:f9:18:86:ee:86:d1:7d:3c:48:a5:70:65:
                    db:64:7a:b7:ca:1e:a3:c4:b8:41:1c:5a:6d:47:f8:
                    86:cc:57:07:8a:0e:38:80:aa:b1:1e:a5:c5:8a:af:
                    ff:a2:f0:d9:8e:11:10:98:bc:f1:4b:eb:f2:9f:f3:
                    f5:38:32:b0:cc:ff:7b:e1:ef:40:c5:29:66:77:c7:
                    9c:56:1b:b4:d6:97:7b:c2:4b:82:3e:cb:f9:bb:c4:
                    41:10:80:5d:af:a1:f7:10:3b:40:83:d3:db:38:89:
                    12:82:a9:b5:4c:da:fa:1f:59:95:61:62:3e:b7:06:
                    07:db:21:71:d8:08:53:f2:22:96:93:c1:cd:ee:e1:
                    f4:56:a6:75:88:52:9f:41:15:4c:bc:3e:a5:e3:87:
                    db:df
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:FALSE
            Netscape Cert Type: 
                SSL Server
            X509v3 Subject Key Identifier: 
                F1:1A:C5:EA:AC:11:B1:30:2D:B1:D1:D4:0D:47:BF:FD:69:E3:36:C9
            X509v3 Authority Key Identifier: 
                keyid:C1:54:7D:53:F7:FD:24:32:E2:D7:BF:49:73:5A:AE:77:38:D8:8E:17
                DirName:/C=US/ST=Colorado/L=Denver/O=Traffic Ops/OU=Traffic Ops/CN=Traffic Ops Root CA/emailAddress=no-reply@invalid2.invalid
                serial:15:51:89:10:66

            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
            X509v3 Subject Alternative Name: 
                DNS:*.test.invalid2.invalid
    Signature Algorithm: sha256WithRSAEncryption
         6d:c3:a2:80:74:75:d7:31:6b:3d:71:d2:ed:91:86:e8:3a:e8:
         ef:4f:10:56:b1:02:a1:f6:49:91:50:0a:39:46:d9:91:63:f6:
         87:84:fe:c5:24:ba:58:4a:a4:ff:58:4a:2a:8d:91:86:22:fb:
         52:51:b0:76:da:8d:67:a9:64:f5:cc:10:39:19:9e:b8:85:c4:
         83:2d:81:b1:b4:79:1f:8b:f2:2f:fd:8b:85:9c:c2:48:36:1f:
         79:26:c1:ca:11:5d:bc:b0:8d:24:c4:86:2b:ea:38:f4:d5:04:
         48:3d:73:10:91:4f:0c:3b:4c:8f:f3:4b:02:09:0b:09:21:a2:
         ff:c8:85:55:52:69:25:3d:f3:3f:52:3f:0b:d6:66:b7:c7:85:
         7f:ab:23:11:da:08:83:a3:cf:45:e2:73:bb:58:5c:10:92:f3:
         fe:3a:8b:bb:3a:06:6a:cc:5e:6b:f9:36:93:4e:a4:46:3d:c6:
         57:c7:f7:05:9a:f2:37:7f:93:29:ae:c8:88:db:13:ed:21:a2:
         64:c8:8f:59:4b:cc:2f:f7:69:f1:a4:e0:08:d0:3f:1a:58:46:
         ab:34:0f:9c:35:c1:be:04:cd:25:f0:d7:89:00:e7:47:5b:2f:
         25:93:47:11:f4:03:3c:81:1e:63:a1:02:04:01:1a:53:4a:e0:
         1f:df:e2:66
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
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 91562250342 (0x1551891066)
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Root CA, emailAddress = no-reply@invalid2.invalid
        Validity
            Not Before: Mar  6 16:51:10 2019 GMT
            Not After : Mar  1 16:51:10 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Intermediate CA, emailAddress = no-reply@invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:b2:ed:f1:57:0a:61:59:4c:ff:8e:8b:00:ad:38:
                    e5:76:77:15:80:06:47:55:af:8d:54:a8:09:a6:88:
                    68:c8:9a:da:ec:88:17:22:93:58:f5:d9:4c:77:8c:
                    32:74:f7:36:fd:8d:b2:73:a9:5c:67:f5:e5:7f:35:
                    90:4b:98:58:b0:ea:8c:a1:6b:32:94:2d:fc:65:e5:
                    fd:ca:68:a7:43:e4:34:4c:5b:33:10:c2:b7:bf:d6:
                    85:57:d9:22:14:82:8f:87:5e:70:17:52:1a:bb:51:
                    cc:ca:93:7a:7b:ad:6b:13:c4:e1:12:2c:4b:53:78:
                    06:2c:78:4c:bc:3b:83:f8:92:28:20:9c:e3:6f:d0:
                    ef:17:41:94:f2:42:6b:fe:ea:f8:a0:94:03:4a:71:
                    a2:5b:87:20:f8:81:ac:c9:21:f8:cb:d7:bd:34:c6:
                    66:c2:99:cd:a6:b5:d5:90:fe:91:af:de:14:28:07:
                    c1:28:90:db:eb:cd:d5:7a:eb:b2:ca:49:f0:0b:a6:
                    4d:99:46:2e:48:b2:14:c5:78:b9:50:f6:ca:42:2d:
                    b4:cf:02:64:1d:b5:ee:8b:c6:81:30:fb:ce:21:42:
                    97:4a:20:9d:cb:96:a9:7f:36:70:db:69:e9:51:46:
                    e4:6a:b8:3a:77:6e:81:53:68:55:8c:e7:7e:e9:5c:
                    ab:a9
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Subject Key Identifier: 
                C1:54:7D:53:F7:FD:24:32:E2:D7:BF:49:73:5A:AE:77:38:D8:8E:17
            X509v3 Authority Key Identifier: 
                keyid:C1:8A:A7:62:2C:5B:E9:76:6A:6C:3E:DD:A4:53:31:DB:CD:19:FB:96

            X509v3 Basic Constraints: critical
                CA:TRUE, pathlen:0
            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign
    Signature Algorithm: sha256WithRSAEncryption
         38:a2:4f:cb:e1:8d:a9:f8:de:62:43:bf:49:05:4b:bc:62:a5:
         e3:c9:79:76:cb:1a:00:2e:dc:fa:1c:a1:1f:6a:51:b8:7e:ed:
         0d:99:8f:9a:75:ac:20:da:a9:4a:14:07:c0:b9:a7:54:f7:26:
         33:f1:78:62:d0:dc:0f:6a:8d:63:f5:15:48:5b:b0:74:df:27:
         37:14:f7:cf:50:ef:72:28:9b:1d:e4:72:69:90:85:6b:7e:c6:
         f0:a9:c4:61:66:be:30:69:13:2c:88:84:29:20:a0:60:54:7c:
         74:05:8d:cc:ab:2d:66:da:2d:15:ce:47:57:85:f4:1d:18:39:
         16:c9:99:fd:74:d6:5c:60:f0:f4:37:2d:c9:42:b3:d9:b8:5d:
         22:36:4c:78:94:44:33:78:74:50:0d:45:8a:9c:5a:f8:34:9d:
         c7:a3:fd:23:97:25:0f:8e:67:92:b5:00:9f:73:fd:a7:01:19:
         3f:c4:a5:0e:98:7d:2c:b1:3a:c2:95:3c:0b:b7:9a:b2:31:7c:
         e0:2a:6a:12:85:25:ae:1b:81:5f:9a:e5:29:a2:10:d6:04:29:
         86:96:39:b0:20:32:93:be:51:de:45:90:6e:48:54:10:c9:89:
         c2:66:fe:46:b3:c8:0b:36:7f:83:66:5b:ef:29:f0:22:84:76:
         b0:7f:ea:02
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
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:c4:a0:21:18:68:49:75:9a:b1:0e:27:2e:a3:3b:
    cf:90:5d:79:a5:6f:c9:98:01:64:0b:fa:54:8e:3d:
    84:6f:ce:e3:d6:3a:f9:df:bd:6f:c1:90:36:29:bf:
    3a:2c:7c:36:60:1c:d4:b5:d3:a4:4e:b5:61:5e:fa:
    44:19:4d:3a:e6:fb:13:d3:58:c9:7f:24:8c:a7:e6:
    46:b5:66:72:4c:01:ec:b4:3b:1c:a7:d8:f8:12:9a:
    5c:be:fb:a7:ac:3a:22:87:09:99:bd:66:1c:76:a7:
    bf:dd:fd:f9:18:86:ee:86:d1:7d:3c:48:a5:70:65:
    db:64:7a:b7:ca:1e:a3:c4:b8:41:1c:5a:6d:47:f8:
    86:cc:57:07:8a:0e:38:80:aa:b1:1e:a5:c5:8a:af:
    ff:a2:f0:d9:8e:11:10:98:bc:f1:4b:eb:f2:9f:f3:
    f5:38:32:b0:cc:ff:7b:e1:ef:40:c5:29:66:77:c7:
    9c:56:1b:b4:d6:97:7b:c2:4b:82:3e:cb:f9:bb:c4:
    41:10:80:5d:af:a1:f7:10:3b:40:83:d3:db:38:89:
    12:82:a9:b5:4c:da:fa:1f:59:95:61:62:3e:b7:06:
    07:db:21:71:d8:08:53:f2:22:96:93:c1:cd:ee:e1:
    f4:56:a6:75:88:52:9f:41:15:4c:bc:3e:a5:e3:87:
    db:df
publicExponent: 65537 (0x10001)
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
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            ab:5f:90:e1:3a:e1:1c:a7
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Root CA, emailAddress = no-reply@invalid2.invalid
        Validity
            Not Before: Mar  6 16:51:07 2019 GMT
            Not After : Feb 24 16:51:07 2059 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Root CA, emailAddress = no-reply@invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:c8:b7:77:56:ed:7a:46:61:6b:bd:94:5a:47:49:
                    be:cb:99:7a:7a:7c:cc:d8:90:5b:b3:9a:2c:bc:99:
                    bc:b6:ef:0a:77:26:ea:c3:46:9d:94:d5:4e:4a:b9:
                    2d:99:eb:ab:6b:0e:73:4e:09:21:fe:e0:14:97:3c:
                    be:72:5d:04:f1:19:10:fe:46:45:b5:97:99:e9:78:
                    43:8a:a7:a0:e7:02:af:8f:fa:51:9e:48:d6:01:b5:
                    26:41:8c:b0:c5:aa:f3:76:6f:fe:3a:ad:8c:14:a9:
                    20:ef:5f:40:51:22:bf:a8:7d:8d:40:24:ce:59:fb:
                    95:fb:e9:33:cf:76:a8:22:99:08:52:e1:f3:d4:d9:
                    e5:f9:52:9b:09:8d:54:e5:b9:20:5b:d4:bb:db:85:
                    74:a0:c1:cc:e6:1a:5e:90:57:5e:1b:9a:dd:ec:21:
                    f5:8a:94:b4:27:f0:04:7f:f3:cf:c6:db:dd:28:66:
                    f6:9e:79:46:14:f4:46:9e:32:42:91:41:e0:5c:06:
                    d0:cc:4d:36:a8:d9:f0:12:ee:17:da:65:79:37:11:
                    95:3a:d5:c1:32:d0:da:c3:aa:ec:79:dd:bf:2c:9a:
                    87:1b:52:3f:70:45:a6:51:6b:9e:7c:e6:78:e6:63:
                    b6:be:d1:44:fb:c6:ff:6d:9f:1f:16:fa:94:6c:67:
                    5c:49
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Subject Key Identifier: 
                C1:8A:A7:62:2C:5B:E9:76:6A:6C:3E:DD:A4:53:31:DB:CD:19:FB:96
            X509v3 Authority Key Identifier: 
                keyid:C1:8A:A7:62:2C:5B:E9:76:6A:6C:3E:DD:A4:53:31:DB:CD:19:FB:96

            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign
    Signature Algorithm: sha256WithRSAEncryption
         91:b5:26:bc:1f:99:b0:34:a4:f5:4a:04:a3:11:ab:4b:30:b8:
         f3:62:4d:63:a5:f2:27:b3:a8:51:d9:7e:5f:f8:3f:e4:66:4a:
         60:3d:7f:dd:e1:d0:24:7b:50:e7:30:8c:b1:13:8d:ed:d5:12:
         49:ce:88:b5:06:11:60:f5:8b:c4:c8:d5:39:93:92:fb:10:2b:
         98:9a:05:39:b5:b9:9e:37:1d:9c:0e:ee:07:bb:2e:64:e9:99:
         90:e0:09:a1:d3:55:f7:25:8d:cc:14:04:49:c9:18:bf:f8:2b:
         e6:4f:e4:9b:c1:d2:98:b3:c7:f0:84:f6:c6:70:b4:61:95:40:
         b4:b6:f1:5e:e8:97:ef:71:e1:ff:d7:c4:d1:95:c0:cf:0b:9d:
         7e:7f:e5:a2:75:59:16:0b:c1:2b:ec:cc:9d:bf:e6:e6:35:8a:
         84:80:31:09:d1:44:16:b4:8f:e3:fe:1c:58:8f:5c:2e:aa:f5:
         1a:a0:dd:80:19:c0:ee:7f:41:f3:9c:6f:c9:30:eb:68:1a:63:
         04:2f:2c:49:94:81:c2:04:f1:b5:87:e9:b2:2c:e3:7d:fd:7b:
         8f:1e:58:4c:41:c2:21:9a:a3:44:04:11:b0:49:dc:26:d5:86:
         dc:dd:6f:2c:52:40:40:12:16:b5:36:6e:12:18:93:b3:d8:b7:
         e4:fe:b0:19
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
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 91562280034 (0x1551898462)
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Intermediate CA, emailAddress = no-reply@invalid2.invalid
        Validity
            Not Before: Mar  6 18:54:25 2019 GMT
            Not After : Mar  1 18:54:25 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Traffic Ops, CN = *.test.invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:d2:32:c9:22:d4:be:03:55:07:d8:77:05:c1:c4:
                    11:c0:58:3d:d9:76:9a:12:f8:c6:5a:de:52:57:3c:
                    c7:9c:a4:d3:90:7d:99:85:cc:8d:7a:97:7a:0a:69:
                    e3:ca:d2:dd:9a:d8:19:b1:3e:6b:9e:89:0e:14:0f:
                    89:f6:c7:76:78:6d:5f:ec:8f:3e:96:d7:70:13:66:
                    cf:8d:dc:46:68:a6:4a:f9:92:3a:f9:4f:83:71:88:
                    76:96:03:62:42:14:11:45:68:c8:49:c2:85:f8:83:
                    de:29:52:e1:ca:1e:6e:62:32:be:e6:af:a9:c9:2c:
                    ac:17:5b:54:6b:67:78:65:f2:ad:ef:06:8a:34:0f:
                    3c:99:d6:52:28:fa:cd:4c:a8:cc:c6:1b:ce:96:77:
                    0e:3d:94:dc:20:8b:87:82:30:13:e7:5f:e5:0e:10:
                    d8:a3:d0:a9:7c:d0:a7:e8:0a:ca:7d:14:ec:5b:02:
                    d6:48:a4:fe:5c:82:66:fc:11:09:e6:b9:aa:15:c1:
                    b7:ad:6b:3d:ac:67:e4:bf:a9:ab:12:af:59:e5:0f:
                    dc:2e:67:be:1d:09:3c:67:3f:3b:7e:86:0b:b2:86:
                    24:8b:ad:ad:2c:77:8f:22:7b:d4:1d:0a:a4:87:c9:
                    80:1c:ae:63:6d:10:b3:58:df:2b:5b:10:dd:8a:1b:
                    4c:67
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: 
                CA:FALSE
            Netscape Cert Type: 
                SSL Server
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage: 
                TLS Web Server Authentication
            X509v3 Subject Alternative Name: 
                DNS:*.test.invalid2.invalid
    Signature Algorithm: sha256WithRSAEncryption
         1e:48:d1:8e:a4:12:b0:b0:89:51:01:3a:40:d0:5c:c7:92:3e:
         86:57:b7:22:71:62:ca:fa:f9:bb:60:c0:b8:bf:79:1b:30:38:
         2d:06:e4:d6:fa:18:16:3d:a1:7c:28:00:13:d3:8d:d9:a5:62:
         3e:39:6a:41:c3:d8:97:ed:8a:a3:29:34:be:44:30:91:26:fc:
         a0:9b:04:b9:c0:1c:00:69:3a:e4:19:a5:05:50:db:33:d2:20:
         54:2f:7a:20:b6:2e:93:49:13:f9:8d:02:82:f1:16:05:fb:68:
         75:8b:3c:c4:a1:00:90:db:a8:a5:44:ca:64:7d:62:3c:b2:05:
         82:12:37:73:1b:94:e7:a1:ba:03:85:d4:14:46:d9:01:fe:6a:
         33:65:01:c4:59:ca:3c:f9:bb:1f:22:ac:ed:63:e8:b9:df:7f:
         57:87:7b:51:a6:ac:a5:7c:05:76:f0:3b:7f:c6:c6:26:f8:23:
         1c:66:04:f1:70:59:38:d8:10:60:95:8b:06:00:9f:20:a1:74:
         2e:a5:c0:17:05:6b:6d:bc:c9:12:c8:51:4f:e5:97:08:db:c2:
         0a:ae:7e:00:92:b8:80:16:a8:bf:39:3b:42:40:11:12:4c:13:
         7b:b0:ee:fb:e3:41:88:77:86:d1:3a:cf:37:ab:57:ec:a7:bb:
         11:46:1c:1e
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
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 91562280032 (0x1551898460)
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Root CA, emailAddress = no-reply@invalid2.invalid
        Validity
            Not Before: Mar  6 18:54:22 2019 GMT
            Not After : Mar  1 18:54:22 2039 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Intermediate CA, emailAddress = no-reply@invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:b6:40:e1:f6:0b:ac:7e:a8:2f:b9:33:7f:11:87:
                    e2:5d:2f:c9:6e:75:f7:e1:1f:4a:8d:ab:b7:25:1a:
                    80:f8:2b:6d:6d:8e:e2:e6:96:1b:6f:ec:5d:54:18:
                    94:48:43:b6:86:29:08:17:d2:db:b8:2a:c4:41:47:
                    70:3d:b0:5d:5f:98:2b:6a:d0:dc:3a:ab:aa:c1:88:
                    f8:ee:e8:7d:8a:2d:54:12:c1:fe:ed:26:b8:69:5d:
                    93:af:5c:ce:14:ba:dd:92:e2:9d:24:85:24:d4:6c:
                    29:eb:27:e8:2c:0b:17:0d:f9:92:70:86:3a:d4:43:
                    1f:11:c6:ae:f7:31:ff:52:50:8d:31:b7:d3:50:c6:
                    08:90:92:4a:f2:23:f2:16:54:bf:11:e8:61:16:e3:
                    dd:fb:43:23:77:fa:97:3c:20:41:0f:46:c2:70:92:
                    d5:91:af:e6:02:32:b2:5d:93:c9:e9:5c:a1:e3:f5:
                    03:e4:d7:be:6d:5a:7b:55:ab:62:59:a8:80:7f:bc:
                    02:17:4a:4e:af:cd:72:4d:d1:18:57:b4:a9:e7:77:
                    d7:31:41:2e:c0:00:82:95:5f:3c:62:51:13:dd:e1:
                    b0:7d:25:dd:8f:bb:a0:28:19:f2:0d:7d:d3:96:03:
                    99:8b:1f:3c:a9:18:6b:26:ba:3e:66:bf:62:e1:13:
                    25:f3
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: critical
                CA:TRUE, pathlen:0
            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign
    Signature Algorithm: sha256WithRSAEncryption
         58:83:3e:47:af:54:06:1b:a6:52:fc:cb:b3:0d:b1:5d:33:e6:
         1c:49:73:b8:cd:30:e4:14:0c:7c:97:2f:67:78:2c:76:8d:30:
         bc:25:e4:29:03:d2:bb:6c:8f:04:7d:01:ba:a4:6d:97:6f:4d:
         83:7f:9c:8a:eb:d1:aa:4f:b0:9c:3a:7c:04:7e:b9:d5:4e:94:
         ee:8a:92:c3:65:7b:6a:53:8e:cb:8b:a0:3d:d4:26:ec:d3:5b:
         26:e3:92:60:20:80:9e:01:21:a2:11:25:a1:6f:f8:08:62:f2:
         c9:8e:1f:7c:d7:e8:3f:24:a0:97:4c:57:2f:c0:7f:41:5b:93:
         bb:0f:f4:20:f2:b8:e3:cc:e2:50:9d:f8:ea:c6:5f:d6:80:50:
         70:cd:b2:fa:01:df:6c:da:7e:8d:03:3d:95:07:90:86:3b:f6:
         4c:a1:59:2a:b7:f2:ea:e4:b9:84:74:00:96:9a:35:ed:ed:28:
         3d:90:3e:63:0e:be:85:0d:67:f7:50:86:67:f5:f1:ad:69:be:
         a0:14:c0:c9:95:2e:2f:43:60:1c:10:b0:1b:36:4c:a6:a2:a4:
         da:1c:9e:be:f3:82:df:25:90:b3:18:5b:87:74:3d:da:e7:31:
         c7:78:9f:42:ff:96:b0:6b:b9:7b:c4:70:9a:bd:ff:76:88:7a:
         be:83:a7:19
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
RSA Private-Key: (2048 bit, 2 primes)
modulus:
    00:d2:32:c9:22:d4:be:03:55:07:d8:77:05:c1:c4:
    11:c0:58:3d:d9:76:9a:12:f8:c6:5a:de:52:57:3c:
    c7:9c:a4:d3:90:7d:99:85:cc:8d:7a:97:7a:0a:69:
    e3:ca:d2:dd:9a:d8:19:b1:3e:6b:9e:89:0e:14:0f:
    89:f6:c7:76:78:6d:5f:ec:8f:3e:96:d7:70:13:66:
    cf:8d:dc:46:68:a6:4a:f9:92:3a:f9:4f:83:71:88:
    76:96:03:62:42:14:11:45:68:c8:49:c2:85:f8:83:
    de:29:52:e1:ca:1e:6e:62:32:be:e6:af:a9:c9:2c:
    ac:17:5b:54:6b:67:78:65:f2:ad:ef:06:8a:34:0f:
    3c:99:d6:52:28:fa:cd:4c:a8:cc:c6:1b:ce:96:77:
    0e:3d:94:dc:20:8b:87:82:30:13:e7:5f:e5:0e:10:
    d8:a3:d0:a9:7c:d0:a7:e8:0a:ca:7d:14:ec:5b:02:
    d6:48:a4:fe:5c:82:66:fc:11:09:e6:b9:aa:15:c1:
    b7:ad:6b:3d:ac:67:e4:bf:a9:ab:12:af:59:e5:0f:
    dc:2e:67:be:1d:09:3c:67:3f:3b:7e:86:0b:b2:86:
    24:8b:ad:ad:2c:77:8f:22:7b:d4:1d:0a:a4:87:c9:
    80:1c:ae:63:6d:10:b3:58:df:2b:5b:10:dd:8a:1b:
    4c:67
publicExponent: 65537 (0x10001)
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
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            94:91:f7:32:32:9b:c8:d2
        Signature Algorithm: sha256WithRSAEncryption
        Issuer: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Root CA, emailAddress = no-reply@invalid2.invalid
        Validity
            Not Before: Mar  6 18:54:22 2019 GMT
            Not After : Feb 24 18:54:22 2059 GMT
        Subject: C = US, ST = Colorado, L = Denver, O = Traffic Ops, OU = Traffic Ops, CN = Traffic Ops Root CA, emailAddress = no-reply@invalid2.invalid
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:a4:8c:9e:da:0e:9c:d2:2a:c5:82:83:00:ce:d3:
                    de:12:e5:78:65:7f:0d:74:e4:51:b1:5d:83:76:2e:
                    24:57:d8:0a:62:20:e4:c5:0f:5e:39:a9:77:35:e6:
                    bd:90:31:1a:52:94:6f:93:69:ee:56:5e:63:6c:51:
                    b7:b0:ea:2b:7b:d5:6c:e9:85:2c:0a:f2:02:0c:a0:
                    94:3e:5e:4e:af:15:11:27:a0:52:88:2d:a5:d2:35:
                    44:e8:55:61:5d:ff:69:2d:7f:8e:47:c9:59:98:c6:
                    7d:18:a6:f0:d6:79:46:18:ac:1d:17:74:fb:ea:03:
                    99:15:21:d0:7d:3e:7b:bc:d1:6c:23:44:3e:f0:d8:
                    56:6c:37:25:36:8f:c0:9c:fa:50:b8:1b:3a:a1:c6:
                    a1:f3:70:40:55:09:37:81:34:4c:1c:ed:fe:ac:c2:
                    ee:bd:75:69:a4:10:6a:0f:e3:f9:39:4f:8b:45:13:
                    ab:8e:80:ee:96:e6:f6:41:43:e2:47:44:39:0d:cc:
                    ea:30:28:c3:21:00:7d:e8:b4:5e:af:23:78:77:1f:
                    e9:e3:1e:0f:eb:64:8b:40:1e:9d:77:6b:c7:bc:93:
                    66:a5:f9:7f:08:1c:c0:75:22:c1:46:76:bd:99:25:
                    7a:c7:0e:36:f6:db:b9:6f:d6:78:f0:36:b9:82:9f:
                    62:81
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign
    Signature Algorithm: sha256WithRSAEncryption
         03:b6:83:29:6b:ca:3a:b6:3a:ef:a6:9b:81:3c:3a:7f:c4:71:
         c8:5c:8b:51:54:95:cf:19:2c:0c:f1:c5:37:6d:24:36:50:77:
         86:fa:b8:41:0e:16:75:f2:20:3c:b5:0d:4d:c2:34:3d:e2:78:
         86:b4:10:bd:76:fd:db:06:cf:38:48:d2:62:44:7f:ea:1e:b3:
         86:9a:6b:89:74:ce:73:00:ed:64:39:f0:25:1d:05:6b:86:a0:
         6b:d4:cc:b9:04:dd:30:05:a8:bc:0a:0b:bf:61:0c:04:50:42:
         9d:15:68:cd:35:a5:94:1d:fd:28:7f:df:e5:77:8a:31:20:df:
         cd:7c:56:d4:ef:28:79:b9:b7:6e:a0:80:84:61:49:21:8c:84:
         47:aa:0a:81:be:f4:42:8b:8a:0c:f3:68:1a:6b:2b:36:d2:b8:
         6a:63:f5:02:8b:3b:83:f8:bc:8b:e0:b7:64:6e:b5:fa:98:04:
         d4:7a:b1:dc:ca:99:06:32:3c:26:2c:85:fb:28:41:55:a8:f6:
         69:83:fc:b2:e9:36:84:b0:c1:71:a0:06:12:3b:72:9a:e9:d0:
         2b:0d:e5:fe:7f:98:bd:9d:0d:24:2a:10:03:a5:f1:88:ad:55:
         ae:27:2f:8e:22:16:b2:fc:77:87:7c:38:ce:f3:8c:8f:ff:fd:
         7f:0c:50:9b
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
	// TODO: Add unit test(s) for multi-cert ECDSA x509v3 certificates
	CASignedECDSACertificateChain                   = ``
	CASignedECDSACertificateChainPrivateKey         = ``
	CASignedECDSARootCA                             = ``
	CASignedNoSkiAkiECDSACertificateChain           = ``
	CASignedNoSkiAkiECDSACertificateChainPrivateKey = ``
	CASignedNoSkiAkiECDSARootCA                     = ``
	SelfSignedRSAPrivateKeyInvalidChain             = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsARDggbjMcM3jaBnyJ67lTNEoJhQHnWrtvOPe69rp8meKXyL
ExeoSTWe2u6stjTFH3zLD716Ph/G1ZRHpq7Kp8zloSMcKI0NsBBV82OEo3viifIX
by5a1Lr10WMD9PW66+LpSx389HPQPuSP/NAeMy7p3xGQQQWDLLOas7sEWFZ1DPYs
D7eF7vD6I9/ZtO7sCyRn36Hm16GnenyoTVkA/f3zelcxYfE02EezufAZuXgOB5y7
xoQApD/gVp+XqouHZ329ZtCCHP5Jeyd6zngpG7ZZAbkVpY+jUkIzevkazHV5ire2
m4XfGO+tOMyqbmlWi7Ts4owaQCQkZMlwK2p36QIDAQABAoIBAEDS0Snp73I8OxFl
qdMw4lSodPXQInGVVJAkUwtyJ2u7zQvqWi3F4KxVmxN2IxVXieF2zDIXzhVjDo9J
9LlmVixGQat+irhEem4FFiJ03Dx5O40iI49GuxztXeqnVKW6egS1pMWNXcOJg4Am
HQE2hGjFNkx442+O4ChuXOMkVQ1S7Xz8Ejif2wheH5oMRLDqrmuXqyC9mkespct+
g9QjMPxEXGp1wXTUzWopbGGMkofr9oL2HXX2CJX6lUAJDE1doaq8QmSkrUwKu1/6
X8S3gqX9kAmy8ff4oAHuqmwv/uRDNGCLtmcfLvtXjDYdlC3ML1FT4B2E3ctc9Kjo
iQkNYIECgYEAy9jPQVA2pTQQ493z3ty1QS5Pyj55qEYSqhwyw3VF8g/hehKRPuJC
LC6PDhWgiaJCo3aiJ8Ac/ve/elSOsTlGx/oReyujk0rixvd9Sxoq2od5zJWpGySB
kRmyb4f/1DH+p+LTGEl7pJOfbfF2b3YnCS5iSP3KS7jeLKGKI4J7rM0CgYEA3Qyt
10RcRIOfxOXT5Sb5nF8DuPKbnJ0x2Tg6oK3pxALX5/sN3FMnld5Xl1rnPGDgaSKf
Gk2kM1rW/6XvKu9seBAFs9bODBZUHNUP7Vv6aaxq0GWU/QTNqxkYaqdt6lLTjlP5
FDxCKqBlon77KMaIFOEQnBKGSj2tUWO7rC/dd40CgYEAjmaY8hFs+x9SJTyp3ifk
XvJRPwFBz3GUHE2ykKReBmldo/9Qg9NfUqn7uWUWTs+RKcv4HziviNXdZ0GmpNtU
POLOT3L+xChuH3xIhKx0/0/goDB0f8eS06BV7F/fMYbzVKi5up+qxh9yIkWp7Ndn
EZzbgA36wccVPaxjacb/SokCgYEA1xbxSRgRl/Fj01m3F7EXDVs+6gXX+UrUKIOY
OKVBZCNIJ0iYshyP1jqljHc9rfiuJF815YhLEFWCAvxZfrO+Hg2pHtcTY5uOeQex
Gct4HL9SqDlQAetcnPIsWgtU3r99b26yXUhNMeElRDq+9WxJGdfuK4+y8CaXsSyU
fvWMUDkCgYA2HozEb1cMu676Sa4TaiF/LvOd9HFZe1DJRfk4jOKbiAgLgvNCfRrr
WMMitU3/DUNIDPsPQxdfPIUqk4PDp5ZuZGau5AsHU2a8byOVm8SASJHv9ha4Memr
ME8pyqZrmhowVZcsMmtVudk74bn9163d3fiiQzCKxokItxyB2MBCyA==
-----END RSA PRIVATE KEY-----
`
	SelfSignedInconsistentCertChain = `
-----BEGIN CERTIFICATE-----
MIID9DCCAtygAwIBAgIRALsBZHSe6BK11fYqQLaGPn0wDQYJKoZIhvcNAQELBQAw
gYoxFDASBgNVBAYTC1BsYWNlaG9sZGVyMRQwEgYDVQQIEwtQbGFjZWhvbGRlcjEU
MBIGA1UEBxMLUGxhY2Vob2xkZXIxFDASBgNVBAoTC1BsYWNlaG9sZGVyMRQwEgYD
VQQLEwtQbGFjZWhvbGRlcjEaMBgGA1UEAwwRKi5hc2RmMi5jaWFiLnRlc3QwHhcN
MjIxMDEzMTgxMDIyWhcNMjMxMDEzMTgxMDIyWjCBijEUMBIGA1UEBhMLUGxhY2Vo
b2xkZXIxFDASBgNVBAgTC1BsYWNlaG9sZGVyMRQwEgYDVQQHEwtQbGFjZWhvbGRl
cjEUMBIGA1UEChMLUGxhY2Vob2xkZXIxFDASBgNVBAsTC1BsYWNlaG9sZGVyMRow
GAYDVQQDDBEqLmFzZGYyLmNpYWIudGVzdDCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBALAEQ4IG4zHDN42gZ8ieu5UzRKCYUB51q7bzj3uva6fJnil8ixMX
qEk1ntrurLY0xR98yw+9ej4fxtWUR6auyqfM5aEjHCiNDbAQVfNjhKN74onyF28u
WtS69dFjA/T1uuvi6Usd/PRz0D7kj/zQHjMu6d8RkEEFgyyzmrO7BFhWdQz2LA+3
he7w+iPf2bTu7AskZ9+h5tehp3p8qE1ZAP3983pXMWHxNNhHs7nwGbl4Dgecu8aE
AKQ/4Fafl6qLh2d9vWbQghz+SXsnes54KRu2WQG5FaWPo1JCM3r5Gsx1eYq3tpuF
3xjvrTjMqm5pVou07OKMGkAkJGTJcCtqd+kCAwEAAaNTMFEwDgYDVR0PAQH/BAQD
AgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHAYDVR0RBBUw
E4IRKi5hc2RmMi5jaWFiLnRlc3QwDQYJKoZIhvcNAQELBQADggEBACFFXQrxyE3Y
tlrsHc5pV2hmqvi4wAj0iFtAfq06c5UbnnCIMRkw1hKIjO2TRN3qEJKxlEklPOIr
uwvn66sIFbB0HgEawjCxoQ+03k/ZSErxYukRmWzJ2knjFP46bxCHxRjBxUvNy47R
DXAlyG5jhAECuxyOBaevgOlk2Eue4Bk7yM7HGCPEhC/xOEBe66rNJtbh6T2NKe9A
q5dMWXDybkreALeRA8wnosXf5n24AMC9QKYqZ/baAvG/ASKvVqTaP9YgQ0tyJhjV
qb6viLd41tRPW9s9Ds5ECmKxmju3Bsh1se94V83jE5nUeHg99/h5IQuFjaCVi5Zt
6/0395yFKRA=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIID8zCCAtugAwIBAgIQDen+PunwC7HgIIy5puEugjANBgkqhkiG9w0BAQsFADCB
ijEUMBIGA1UEBhMLUGxhY2Vob2xkZXIxFDASBgNVBAgTC1BsYWNlaG9sZGVyMRQw
EgYDVQQHEwtQbGFjZWhvbGRlcjEUMBIGA1UEChMLUGxhY2Vob2xkZXIxFDASBgNV
BAsTC1BsYWNlaG9sZGVyMRowGAYDVQQDDBEqLmFzZGYzLmNpYWIudGVzdDAeFw0y
MjEwMTMxODEzNDFaFw0yMzEwMTMxODEzNDFaMIGKMRQwEgYDVQQGEwtQbGFjZWhv
bGRlcjEUMBIGA1UECBMLUGxhY2Vob2xkZXIxFDASBgNVBAcTC1BsYWNlaG9sZGVy
MRQwEgYDVQQKEwtQbGFjZWhvbGRlcjEUMBIGA1UECxMLUGxhY2Vob2xkZXIxGjAY
BgNVBAMMESouYXNkZjMuY2lhYi50ZXN0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAscFl+DbbcaJPxpxFVglVKW8P3n+SenHhexhPpBtSiPtPXX5NsAmM
zmDypmuyOZUNWMqyFnd5g7R/0qzCDJnB/c60QtcxSeN8TIFoUocYEU8GAe3GWnMs
75QhhA/2ps36+JfUs7+ZALwt169EhuzX4/R0jmpm2yOtqUEQd1Q6nWxENFip1Ya7
04+tJqtXfpS3kHjpVL6v24bW16TZbjnKSp5nYpha2KvV3fTxmas/1f+s3R2KPAvF
yoh73pXxbG7dinVhuNCiVD1ltRrZ+WWKJLg/IFI6aj7WkvZHukgzPSBNz0BKyeN2
cmGHTNPkyClaEmbI13TaR7fUymwGxqZBewIDAQABo1MwUTAOBgNVHQ8BAf8EBAMC
BaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADAcBgNVHREEFTAT
ghEqLmFzZGYzLmNpYWIudGVzdDANBgkqhkiG9w0BAQsFAAOCAQEAERS0gDP+gRPO
x9pCT4/yHC464N7XnQSBq8etd3JQqCQ1KhM4XxGacNRU6RbLw4haL3S5NEVJlEAw
/ScQi4RMU3z9c6mIFTP3Ha45VFcPjK8VQeemgMjiWueVF8tKpbWBmH/RZnEkzGsK
5O24QpIvKI37XkgBWJbkk4kDMkKve9tXoPAHMEsFXUdoUnZUh4DtP9uPWErtfcGv
fdFz/L9Qn4sLjVWVzchGSTNiD5I6Yf0gIaSfndUgqlOK0qHL+e5njBTW60AgDTBm
bIlRqYdgVaemmQPQWKPjMUKbnevs+FGdNk8pvPTy7/SNVjOmV0+11UR5dlLnqyzQ
3Uh7vjRNqw==
-----END CERTIFICATE-----
`
	SelfSignedCertChain = `
-----BEGIN CERTIFICATE-----
MIID9DCCAtygAwIBAgIRALsBZHSe6BK11fYqQLaGPn0wDQYJKoZIhvcNAQELBQAw
gYoxFDASBgNVBAYTC1BsYWNlaG9sZGVyMRQwEgYDVQQIEwtQbGFjZWhvbGRlcjEU
MBIGA1UEBxMLUGxhY2Vob2xkZXIxFDASBgNVBAoTC1BsYWNlaG9sZGVyMRQwEgYD
VQQLEwtQbGFjZWhvbGRlcjEaMBgGA1UEAwwRKi5hc2RmMi5jaWFiLnRlc3QwHhcN
MjIxMDEzMTgxMDIyWhcNMjMxMDEzMTgxMDIyWjCBijEUMBIGA1UEBhMLUGxhY2Vo
b2xkZXIxFDASBgNVBAgTC1BsYWNlaG9sZGVyMRQwEgYDVQQHEwtQbGFjZWhvbGRl
cjEUMBIGA1UEChMLUGxhY2Vob2xkZXIxFDASBgNVBAsTC1BsYWNlaG9sZGVyMRow
GAYDVQQDDBEqLmFzZGYyLmNpYWIudGVzdDCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBALAEQ4IG4zHDN42gZ8ieu5UzRKCYUB51q7bzj3uva6fJnil8ixMX
qEk1ntrurLY0xR98yw+9ej4fxtWUR6auyqfM5aEjHCiNDbAQVfNjhKN74onyF28u
WtS69dFjA/T1uuvi6Usd/PRz0D7kj/zQHjMu6d8RkEEFgyyzmrO7BFhWdQz2LA+3
he7w+iPf2bTu7AskZ9+h5tehp3p8qE1ZAP3983pXMWHxNNhHs7nwGbl4Dgecu8aE
AKQ/4Fafl6qLh2d9vWbQghz+SXsnes54KRu2WQG5FaWPo1JCM3r5Gsx1eYq3tpuF
3xjvrTjMqm5pVou07OKMGkAkJGTJcCtqd+kCAwEAAaNTMFEwDgYDVR0PAQH/BAQD
AgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHAYDVR0RBBUw
E4IRKi5hc2RmMi5jaWFiLnRlc3QwDQYJKoZIhvcNAQELBQADggEBACFFXQrxyE3Y
tlrsHc5pV2hmqvi4wAj0iFtAfq06c5UbnnCIMRkw1hKIjO2TRN3qEJKxlEklPOIr
uwvn66sIFbB0HgEawjCxoQ+03k/ZSErxYukRmWzJ2knjFP46bxCHxRjBxUvNy47R
DXAlyG5jhAECuxyOBaevgOlk2Eue4Bk7yM7HGCPEhC/xOEBe66rNJtbh6T2NKe9A
q5dMWXDybkreALeRA8wnosXf5n24AMC9QKYqZ/baAvG/ASKvVqTaP9YgQ0tyJhjV
qb6viLd41tRPW9s9Ds5ECmKxmju3Bsh1se94V83jE5nUeHg99/h5IQuFjaCVi5Zt
6/0395yFKRA=
-----END CERTIFICATE-----`
	ValidIntermediateCert = `
-----BEGIN CERTIFICATE-----
MIIEvjCCA6agAwIBAgIQBtjZBNVYQ0b2ii+nVCJ+xDANBgkqhkiG9w0BAQsFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0yMTA0MTQwMDAwMDBaFw0zMTA0MTMyMzU5NTlaME8xCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxKTAnBgNVBAMTIERpZ2lDZXJ0IFRMUyBS
U0EgU0hBMjU2IDIwMjAgQ0ExMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAwUuzZUdwvN1PWNvsnO3DZuUfMRNUrUpmRh8sCuxkB+Uu3Ny5CiDt3+PE0J6a
qXodgojlEVbbHp9YwlHnLDQNLtKS4VbL8Xlfs7uHyiUDe5pSQWYQYE9XE0nw6Ddn
g9/n00tnTCJRpt8OmRDtV1F0JuJ9x8piLhMbfyOIJVNvwTRYAIuE//i+p1hJInuW
raKImxW8oHzf6VGo1bDtN+I2tIJLYrVJmuzHZ9bjPvXj1hJeRPG/cUJ9WIQDgLGB
Afr5yjK7tI4nhyfFK3TUqNaX3sNk+crOU6JWvHgXjkkDKa77SU+kFbnO8lwZV21r
eacroicgE7XQPUDTITAHk+qZ9QIDAQABo4IBgjCCAX4wEgYDVR0TAQH/BAgwBgEB
/wIBADAdBgNVHQ4EFgQUt2ui6qiqhIx56rTaD5iyxZV2ufQwHwYDVR0jBBgwFoAU
A95QNVbRTLtm8KPiGxvDl7I90VUwDgYDVR0PAQH/BAQDAgGGMB0GA1UdJQQWMBQG
CCsGAQUFBwMBBggrBgEFBQcDAjB2BggrBgEFBQcBAQRqMGgwJAYIKwYBBQUHMAGG
GGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBABggrBgEFBQcwAoY0aHR0cDovL2Nh
Y2VydHMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNydDBCBgNV
HR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRH
bG9iYWxSb290Q0EuY3JsMD0GA1UdIAQ2MDQwCwYJYIZIAYb9bAIBMAcGBWeBDAEB
MAgGBmeBDAECATAIBgZngQwBAgIwCAYGZ4EMAQIDMA0GCSqGSIb3DQEBCwUAA4IB
AQCAMs5eC91uWg0Kr+HWhMvAjvqFcO3aXbMM9yt1QP6FCvrzMXi3cEsaiVi6gL3z
ax3pfs8LulicWdSQ0/1s/dCYbbdxglvPbQtaCdB73sRD2Cqk3p5BJl+7j5nL3a7h
qG+fh/50tx8bIKuxT8b1Z11dmzzp/2n3YWzW2fP9NsarA4h20ksudYbj/NhVfSbC
EXffPgK2fPOre3qGNm+499iTcc+G33Mw+nur7SpZyEKEOxEXGlLzyQ4UfaJbcme6
ce1XR2bFuAJKZTRei9AqPCCcUZlM51Ke92sRKw2Sfh3oius2FkOH6ipjv3U/697E
A7sKPPcw7+uvTPyLNhBzPvOk
-----END CERTIFICATE-----`
)

func TestDecodePrivateKeyPKCS8RSA2048(t *testing.T) {

	privateKey, cleanPemPrivateKey, err := decodeRSAPrivateKey(PrivateKeyPKCS8RSA2048)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	pBlock, remain := pem.Decode([]byte(cleanPemPrivateKey))

	if pBlock == nil {
		t.Fatal("can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("remaining bytes after decode > 0. expected: 0")
	}

	if privateKey == nil {
		t.Fatal("RSA private key is nil. expect: not nil")
	}
}

func TestDecodePrivateKeyPKCS1RSA2048(t *testing.T) {

	privateKey, cleanPemPrivateKey, err := decodeRSAPrivateKey(PrivateKeyPKCS1RSA2048)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	pBlock, remain := pem.Decode([]byte(cleanPemPrivateKey))

	if pBlock == nil {
		t.Fatal("can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("remaining bytes after decode > 0. expected: 0")
	}

	if privateKey == nil {
		t.Fatal("RSA private key is nil")
	}
}

func TestDecodePrivateKeyECDSANISTPrime256V1(t *testing.T) {

	privateKey, cleanPemPrivateKey, err := decodeECDSAPrivateKey(PrivateKeyECDSANISTPrime256V1)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	var pemData = []byte(cleanPemPrivateKey)
	var parsedBlocks = make([]*pem.Block, 0)

	for len(pemData) > 0 {
		var block *pem.Block = nil

		// Check for at least one END marker
		if strings.Count(string(pemData), "\n-----END") == 0 {
			break
		}

		block, pemData = pem.Decode(pemData)

		if block == nil {
			t.Fatal("can't decode cleaned ecdsa private-key/param pem block")
		}

		parsedBlocks = append(parsedBlocks, block)
	}

	expectedParsedBlocks := 2
	if len(parsedBlocks) != expectedParsedBlocks {
		t.Fatalf("incorrect number of parsed pem blocks - expected:%d actual:%d", expectedParsedBlocks, len(parsedBlocks))
	}

	if privateKey == nil {
		t.Fatal("ECDSA private key is nil")
	}
}

func TestDecodePrivateKeyECDSANISTPrime256V1WithoutParams(t *testing.T) {

	privateKey, cleanPemPrivateKey, err := decodeECDSAPrivateKey(PrivateKeyECDSANISTPrime256V1WithoutParams)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	pBlock, remain := pem.Decode([]byte(cleanPemPrivateKey))

	if pBlock == nil {
		t.Fatal("can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("remaining bytes after decode > 0. expected: 0")
	}

	if privateKey == nil {
		t.Fatal("ECDSA private key is nil")
	}
}

func TestDecodePrivateKeyECDSANISTSecP384R1(t *testing.T) {
	privateKey, cleanPemPrivateKey, err := decodeECDSAPrivateKey(PrivateKeyECDSANISTSecP384R1)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	var pemData = []byte(cleanPemPrivateKey)
	var parsedBlocks = make([]*pem.Block, 0)

	for len(pemData) > 0 {
		var block *pem.Block = nil

		// Check for at least one END marker
		if strings.Count(string(pemData), "\n-----END") == 0 {
			break
		}

		block, pemData = pem.Decode(pemData)

		if block == nil {
			t.Fatal("can't decode cleaned ecdsa private-key/param pem block")
		}

		parsedBlocks = append(parsedBlocks, block)
	}

	expectedParsedBlocks := 2
	if len(parsedBlocks) != expectedParsedBlocks {
		t.Fatalf("incorrect number of parsed pem blocks - expected:%d actual:%d", expectedParsedBlocks, len(parsedBlocks))
	}

	if privateKey == nil {
		t.Fatal("ECDSA private key is nil")
	}
}

func TestDecodePrivateKeyECDSANISTSecP384R1WithoutParams(t *testing.T) {
	privateKey, cleanPemPrivateKey, err := decodeECDSAPrivateKey(PrivateKeyECDSANISTSecP384R1WithoutParams)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	pBlock, remain := pem.Decode([]byte(cleanPemPrivateKey))

	if pBlock == nil {
		t.Fatal("can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("remaining bytes after decode > 0. expected: 0")
	}

	if privateKey == nil {
		t.Fatal("ECDSA private key is nil")
	}
}

func TestDecodeRSAPrivateKeyBadData(t *testing.T) {

	// Expected to fail.
	privateKey, _, err := decodeRSAPrivateKey(BadKeyData)
	if err == nil && privateKey != nil {
		t.Fatal("unexpected result: decoding of bad private key data should have returned an error")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestDecodeECDSAPrivateKeyBadData(t *testing.T) {

	// Expected to fail.
	privateKey, _, err := decodeECDSAPrivateKey(BadKeyData)
	if err == nil && privateKey != nil {
		t.Fatal("unexpected result: decoding of bad private key data should have returned an error")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestDecodePrivateKeyRSAEncrypted(t *testing.T) {

	// Expected to fail on decode of encrypted pem rsa private key
	privateKey, _, err := decodeRSAPrivateKey(PrivateKeyEncryptedRSA2048)
	if err == nil && privateKey != nil {
		t.Fatal("unexpected result: decoding of encrypted private key should have returned an error")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestDecodePrivateKeyECDSAEncrypted(t *testing.T) {

	// Expected to fail on decode of encrypted pem ecdsa private key
	privateKey, _, err := decodeECDSAPrivateKey(PrivateKeyECDSANISTSecP384R1Encrypted)
	if err == nil && privateKey != nil {
		t.Fatal("unexpected result: decoding of encrypted private key should have returned an error")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateBadData(t *testing.T) {
	// should fail bad base64 data
	_, _, _, _, _, err := verifyCertKeyPair(BadCertData, BadKeyData, "", true)
	if err == nil {
		t.Fatalf("unexpected result: there should have been a base64 decoding failure")
	}
}

func TestVerifyAndEncodeCertificateSelfSignedCertKeyPairDSA(t *testing.T) {
	// should fail due to x509 + DSA being unsupported
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedDSACertificate, SelfSignedDSAPrivateKey, "", true)

	if err == nil {
		t.Fatalf("unexpected result: the DSA PKI algorithm is unsupported")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateSelfSignedX509v1(t *testing.T) {
	// should successfully validate as x509v1 must remain supported
	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(SelfSignedX509v1Certificate, SelfSignedX509v1PrivateKey, "", true)

	if err != nil {
		t.Fatalf("unexpected result: the x509v1 cert/key pair is valid and should have passed validation: %v", err)
	}

	if !unknownAuth {
		t.Fatalf("unexpected result: certificate verification should have detected signature of unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: certChain should not be empty")
	}
}

func TestVerifyAndEncodeCertificateSelfSignedNoSkiAkiCertKeyPair(t *testing.T) {

	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(SelfSignedNOSKIAKIRSACertificate, SelfSignedNOSKIAKIRSAPrivateKey, "", true)

	if err != nil {
		t.Fatalf("unexpected result: a certificate verification error should have occured: %v", err)
	}

	if !unknownAuth {
		t.Fatalf("unexpected result: certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: certChain should not be empty")
	}
}

func TestVerifyAndEncodeCertificateSelfSignedCertKeyPair(t *testing.T) {

	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(SelfSignedRSACertificate, SelfSignedRSAPrivateKey, "", true)

	if err != nil {
		t.Fatalf("unexpected result, a certificate verification error should have occured: %v", err)
	}

	if !unknownAuth {
		t.Fatalf("unexpected result, certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateSelfSignedCertKeyPairMisMatchedPrivateKey(t *testing.T) {

	// Should fail on cert/private-key mismatch
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedRSACertificate, PrivateKeyPKCS1RSA2048, "", true)

	if err == nil {
		t.Fatalf("unexpected result: a certificate/key modulus mismatch error should have occurred")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateSelfSignedNoServerAuthExtKeyUsage(t *testing.T) {

	// Should fail due to not having the serverAuth extKeyUsage (x509v3 only)
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedRSACertificateNoServerAuthExtKeyUsage, SelfSignedRSAPrivateKeyNoServerAuthExtKeyUsage, "", true)

	if err == nil {
		t.Fatalf("unexpected result: x509v3 certificate should have been rejected since it doesn't have the server auth extKeyUsage")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateSelfSignedRSANoKeyEnciphermentKeyUsage(t *testing.T) {

	// Should fail due to not having the keyEncipherment keyUsage (x509v3 only)
	// keyUsage extension must include keyEncipherment if the PKI algorithm is RSA.
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedRSACertificateNoKeyEnciphermentKeyUsage, SelfSignedRSAPrivateKeyNoKeyEnciphermentKeyUsage, "", true)

	if err == nil {
		t.Fatalf("unexpected result: x509v3 certificate should have been rejected since it doesn't have the keyEncipherment keyUsage (required for x509v3+RSA certificates)")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateCASignedCertKeyPair(t *testing.T) {

	// Should succeed, but with unknown authority warning.
	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(CASignedRSACertificateChain, CASignedRSACertificateChainPrivateKey, "", true)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	if !unknownAuth {
		t.Fatalf("unexpected result, certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateCASignedCertKeyPairWithRootCA(t *testing.T) {

	// should succeed and be fully validated.
	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(CASignedRSACertificateChain, CASignedRSACertificateChainPrivateKey, CASignedRSARootCA, true)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	if unknownAuth {
		t.Fatalf("unexpected result: warning for unknown authority even though rootCA is in certChain")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateCASignedNoSkiAkiCertKeyPair(t *testing.T) {

	// Should succeed, but with unknown authority warning.
	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(CASignedNoSkiAkiRSACertificateChain, CASignedNoSkiAkiRSAPrivateKey, "", true)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	if !unknownAuth {
		t.Fatalf("unexpected result, certificate verification should have detected unknown authority")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: certchain should not empty")
	}
}

func TestVerifyAndEncodeCertificateCASignedNoSkiAkiCertKeyPairWithRootCA(t *testing.T) {

	// should succeed and be fully validated despite not having subject/authority key identifier(s).
	certChain, certPrivateKey, unknownAuth, _, _, err := verifyCertKeyPair(CASignedNoSkiAkiRSACertificateChain, CASignedNoSkiAkiRSAPrivateKey, CASignedNoSkiAkiRSARootCA, true)

	if err != nil {
		t.Fatalf("unexpected result: " + err.Error())
	}

	if unknownAuth {
		t.Fatalf("unexpected result: warning for unknown authority even though rootCA is in certChain")
	}

	// Decode the clean Private Key
	pBlock, remain := pem.Decode([]byte(certPrivateKey))

	if pBlock == nil {
		t.Fatal("unexpected result: can't decode cleaned private key pem block")
	} else if len(remain) > 0 {
		t.Fatal("unexpected result: remaining bytes after decode > 0. expected: 0")
	}

	if len(certChain) == 0 {
		t.Fatal("unexpected: cert chain should not empty")
	}
}

func TestVerifyAndEncodeCertificateECDSASelfSignedCertificateKeyPairWithoutDigitalSignatureKeyUsage(t *testing.T) {
	// Should fail due to unsupported private/public key algorithm
	// keyUsage must contain the digitalSignature usage if a DSA based PKI algorithm is indicated (unlike RSA).
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedECDSACertificateNoDigitalSignatureKeyUsage, SelfSignedECDSAPrivateKeyNoDigitalSignatureKeyUsage, "", true)

	if err == nil {
		t.Fatalf("unexpected result: x509v3 certificate should have been rejected since it doesn't have the digitalSignature keyUsage (required for x509v3+ECDSA certificates)")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateECDSASelfSignedCertificateKeyPair(t *testing.T) {
	// Should be successful as the certificate and key are valid with proper keyUsage/extendedKeyUsage.
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedECDSACertificate, SelfSignedECDSAPrivateKey, "", true)

	if err != nil {
		t.Fatalf("unexpected result - valid ECDSA cert/key pair should have been validated: " + err.Error())
	}
}

func TestVerifyAndEncodeCertificateECDSASelfSignedCertificateKeyPairECDisabled(t *testing.T) {
	// Should be rejected as allowEC flag has been set to false
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedECDSACertificate, SelfSignedECDSAPrivateKey, "", false)

	if err == nil {
		t.Fatalf("unexpected result - ECDSA cert/key pair should have been rejected due to allowEC being false")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyAndEncodeCertificateECDSASelfSignedCertificateKeyPairWithoutParams(t *testing.T) {
	// Test result should be successful as the certificate and key are valid with proper keyUsage/extendedKeyUsage.
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedECDSACertificate, SelfSignedECDSAPrivateKeyWithoutParams, "", true)

	if err != nil {
		t.Fatalf("unexpected result - valid ECDSA cert/key pair should have been validated: " + err.Error())
	}
}

func TestVerifyAndEncodeCertificateECDSASelfSignedCertificateKeyPairMisMatchedPrivateKey(t *testing.T) {
	// Should fail due to mismatched ECDSA cert/private key pair
	_, _, _, _, _, err := verifyCertKeyPair(SelfSignedECDSACertificate, PrivateKeyECDSANISTPrime256V1, "", true)

	if err == nil {
		t.Fatalf("unexpected Result: Mismatched ECDSA cert/key pair should have failed verification")
	} else {
		t.Logf("expected error message: %s", err.Error())
	}
}

func TestVerifyInconsistentCertChain(t *testing.T) {
	_, _, _, _, isInconsistent, err := verifyCertKeyPair(SelfSignedInconsistentCertChain, SelfSignedRSAPrivateKeyInvalidChain, "", true)
	if err != nil {
		t.Errorf("expected mismatched to return no error")
	}
	if !isInconsistent {
		t.Errorf("expected chain to be considered inconsistent")
	}

	_, _, _, _, isInconsistent, err = verifyCertKeyPair(SelfSignedCertChain+ValidIntermediateCert, SelfSignedRSAPrivateKeyInvalidChain, "", true)
	if err != nil {
		t.Errorf("expected mismatched to return no error")
	}
	if !isInconsistent {
		t.Errorf("expected chain to be considered inconsistent")
	}
}
