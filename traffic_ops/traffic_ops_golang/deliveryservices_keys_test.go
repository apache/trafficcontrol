package main

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

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
	"strings"
	"testing"
)

const (
	BadData = "This is bad data and it is not base64 encoded"

	SelfSigneCertOnly = `
LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURrakNDQW5vQ0NRQ2Z3ZDIxOUpLcFVEQU5C
Z2txaGtpRzl3MEJBUXNGQURDQmlqRUxNQWtHQTFVRUJoTUMKVlZNeEVUQVBCZ05WQkFnTUNFTnZi
Rzl5WVdSdk1ROHdEUVlEVlFRSERBWkVaVzUyWlhJeEVEQU9CZ05WQkFvTQpCME52YldOaGMzUXhE
akFNQmdOVkJBc01CWFpwY0dWeU1SVXdFd1lEVlFRRERBeDNkM2N1ZEdWemRDNXZjbWN4CkhqQWNC
Z2txaGtpRzl3MEJDUUVXRDIxcFkydGxlVUIwWlhOMExtOXlaekFlRncweE56RXhNVFl4TmpNNU1U
TmEKRncweU56RXhNVFF4TmpNNU1UTmFNSUdLTVFzd0NRWURWUVFHRXdKVlV6RVJNQThHQTFVRUNB
d0lRMjlzYjNKaApaRzh4RHpBTkJnTlZCQWNNQmtSbGJuWmxjakVRTUE0R0ExVUVDZ3dIUTI5dFky
RnpkREVPTUF3R0ExVUVDd3dGCmRtbHdaWEl4RlRBVEJnTlZCQU1NREhkM2R5NTBaWE4wTG05eVp6
RWVNQndHQ1NxR1NJYjNEUUVKQVJZUGJXbGoKYTJWNVFIUmxjM1F1YjNKbk1JSUJJakFOQmdrcWhr
aUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBNm02agppQU1KcExNeFpoSVVFOURyUHZuaTA2
TmViTFJhaDFSUzgrUjNKM1BOMVhDY2d3d1VaaWZaQnBQZ2FTYXBGWGJ2ClpkamxyZlpGcGtBZmdG
UVhseUdhZTZCVStuNEN3TmxEZGdvMVhiRWM4cW9HSWRzN2FZSEVrclFadFp5aCtYRlkKMkJSdlM2
Y2JHR0VjMngwdzVJa1hhMTM2V0NKY0x0QnpkVDdGQVZRSlZodTl4UFBuWGs4aWdUWmc2dEZ0MjdF
YgplYkhDVWVEQXJHVVJGaUZXZFhtOGRHQ1BVVkNaeFNDdnh1WGxJMTVkZGdmcDlHYkZYeENVUDVW
UjlQajZCdHFlCjdReW1GUk9zSWtscEN3NDZlYTBBODdpa1ZNYWRQQzIxVHViZmF5VUFNUTh4bDYw
aW94NVk4OWc0WEZyRWdreWQKV3hobVBXclZwVlFyUERXS2d3SURBUUFCTUEwR0NTcUdTSWIzRFFF
QkN3VUFBNElCQVFDWlpXNnJsbi9LYkdWbgpVWis3RFY5UFVCNHIyNUI5d3JSQUx5OHJRdzNVQWVQ
SFJxWmlBT1ZOV3hycjM5Njd6ZFBNYzhkMWtZY2t0K2NHCkcxUFU1QmNiVVZuUzJxRTBUa2Z2cWo2
citPaWNrcnU2dEtmWERhcU1PRXRmTmFud1Q5dURaQ1RlT2FpU0hCengKRnBQLzlURDA1Z0VZYmQx
VzRSUnpUNi9TN3lMWTB1WWhWQUhGZGwyZTd6T2podHk0aURQRjA0ZmQrWHlaKzNXUwpJeFNnYU1y
VHAwMG1hRnRicFBuaFRheElHek1ZdG9kaTZGYVVjT21CZk1scFJHS3lrc04wWjNlSjZSNm5oTWIz
ClJHSjludUdMS3ZxUzV1OUJnZ2NEd28xcXYyazhNY2RTZzZwbEs2WG9kTHlNZEJWVEI2Szdod1N6
MXVWcDNtWDEKSFZHRTQrb3UKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`
	GoodTLSKeys = `
LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUYzRENDQThTZ0F3SUJBZ0lDRUFBd0RRWUpL
b1pJaHZjTkFRRUxCUUF3Z1l3eEN6QUpCZ05WQkFZVEFsVlQKTVJFd0R3WURWUVFJREFoRGIyeHZj
bUZrYnpFa01DSUdBMVVFQ2d3YlNYQmpaRzRnUTJWeWRHbG1hV05oZEdVZwpRWFYwYUc5eWFYUjVN
U1F3SWdZRFZRUUxEQnRKY0dOa2JpQkRaWEowYVdacFkyRjBaU0JCZFhSb2IzSnBkSGt4CkhqQWNC
Z05WQkFNTUZVbHdZMlJ1SUVsdWRHVnliV1ZrYVdGMFpTQkRRVEFlRncweE56RXhNVFl5TURVd01U
bGEKRncweU5qQXlNREl5TURVd01UbGFNSEF4Q3pBSkJnTlZCQVlUQWxWVE1SRXdEd1lEVlFRSURB
aERiMnh2Y21GawpiekVQTUEwR0ExVUVCd3dHUkdWdWRtVnlNUTh3RFFZRFZRUUtEQVpKY0dOa2Jp
QXhFakFRQmdOVkJBc01DVWx3ClkyUnVJR1JsZGpFWU1CWUdBMVVFQXd3UGQzZDNMbVY0WVcxd2JH
VXVZMjl0TUlJQklqQU5CZ2txaGtpRzl3MEIKQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBMWg5Zlh3
SjE3UVlOTVhNNWpFZ0hlY1Z4V20rQ05QUWJpMk8wYUxtNQpsYUV1b2N0bC94bmNsRGdrT2NWNTBz
SWdVUlNuYVJYSENNYTkzemI5NXhRZWZUSXZSRi8xa1U5ZTY4bW1Gano4CmZyQWRPYU05MFRkVWYx
eXYyNlczN25UOUR4MjZDQWlTd0FMWFZkeCs0b1Bvck5IdjIweDNxUGJzKzNjRHVuQWYKL1hWWis2
dHhuOVZzZzAwb2RIWm9mcVpibUdUdkcveWJRY0dQalJQYVdZSGFMYWZqcERCM3dQc0REdTBMMGJY
VApITlE4Vld6S3drOFVRdHpndEt4WEFFMlp3TFNkUVdzV0VLSEZGTmthT3F3ZnEyMVVZdmN0NHEy
cnE5MXREbmJhCmZxTWpOTnNxVmtYZFIyNFdjUDhtZ1YxY3NWTkx1WmhtUWFCTlo2dytoR2l0aVFJ
REFRQUJvNElCWVRDQ0FWMHcKQ1FZRFZSMFRCQUl3QURBUkJnbGdoa2dCaHZoQ0FRRUVCQU1DQmtB
d013WUpZSVpJQVliNFFnRU5CQ1lXSkU5dwpaVzVUVTB3Z1IyVnVaWEpoZEdWa0lGTmxjblpsY2lC
RFpYSjBhV1pwWTJGMFpUQWRCZ05WSFE0RUZnUVVpWmVaCjVVbXdGRjFWQ2dwd1hmM2crcTFGUzZN
d2djTUdBMVVkSXdTQnV6Q0J1SUFVSVY3aXF0Myt4WWExRngvbll6eGQKMWh4dk53MmhnWnVrZ1pn
d2daVXhDekFKQmdOVkJBWVRBbFZUTVJFd0R3WURWUVFJREFoRGIyeHZjbUZrYnpFUApNQTBHQTFV
RUJ3d0dSR1Z1ZG1WeU1TUXdJZ1lEVlFRS0RCdEpjR05rYmlCRFpYSjBhV1pwWTJGMFpTQkJkWFJv
CmIzSnBkSGt4SkRBaUJnTlZCQXNNRzBsd1kyUnVJRU5sY25ScFptbGpZWFJsSUdGMWRHaHZjbWww
ZVRFV01CUUcKQTFVRUF3d05TWEJqWkc0Z1VtOXZkQ0JEUVlJQ0VBQXdEZ1lEVlIwUEFRSC9CQVFE
QWdXZ01CTUdBMVVkSlFRTQpNQW9HQ0NzR0FRVUZCd01CTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElD
QVFCSFgrWW1PNlRMYk1UWkJ5dXdyZ2lICnV5N3JQenBYVFhTSnZVTSt1cHhTZW5ybXo4SDVmQnp5
cExmS0UxTk9LNThScjI1T3lBZDFJcWdHMUlkKzh5Y2cKYXRldGpJMXBvYU5DWnZLUFJyNDh3amNU
VnNmZ3g3SFhXMGZ5a0kwS1QyMzlPazR1c0VNM3VxbWJrbjdBbm5mNwpkYnk4Qkt6NHA5bjF5ZTdl
SHBpcHUxQTVXeGpCUWVEeW1sUHcvQ005RTB2NUYvaGdpV3hXK2U4bDV1N2g0M2JQCjFCTk82UEpY
Si9HT05pUWg0YmNMMFI3NjN0TklLR1ZOWUhTZmZXNzF2eW9ZNnJlWGRLQk5uazkxT2dqRndIVE4K
QUxKdHI5ank1eHk5dktBVFk3Z1RoVldEaFcvdWl2OExVSUN4eFpXbkN1Y0t2aGdCMmJMV2FIeElP
cjlZeWphWApoRVZkSGdUVlhQbDVjTVhCdytVWkxuM3Z3NUVRSHVNRkFJV1IyWk5JaHJFdmphUG9u
NXFhRXZlUUNEc2RhMzl3CnJUcjNsaGV5QkdNaFRubUE0d05WK2xOU0xJbUFjMm83OGFldGUwY1Rz
SjYxOXNuR09EVW85UU8yWVovazd2UFMKbUhlSzBWclp2bkU4dWtsRkpVVERheUxvRCtTbnN4enRS
SVR2ay9FdVJheTJQMWZHMS9xWm43UVVJZytaTTRhVgpTOHZQcmZJOXJ3dThDUmd4bVRUdFg1Wk1l
ejg5Znk2RmxWRElUalZOaFJyVkljMm0zeEwwM0FTVGI2TmFua0pvCkVuK01MUXBDRlZLeVF5OU43
MlJaTno2OEtpUEgvZnJhQWRxNG1MNU5BY1I0UC82R2c3U24rcWpuZi80bDQzeloKKzd1enFqd2hk
dFUvWUhKOWo4L2Nvdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KLS0tLS1CRUdJTiBDRVJU
SUZJQ0FURS0tLS0tCk1JSUdCVENDQSsyZ0F3SUJBZ0lDRUFBd0RRWUpLb1pJaHZjTkFRRUxCUUF3
Z1pVeEN6QUpCZ05WQkFZVEFsVlQKTVJFd0R3WURWUVFJREFoRGIyeHZjbUZrYnpFUE1BMEdBMVVF
Qnd3R1JHVnVkbVZ5TVNRd0lnWURWUVFLREJ0SgpjR05rYmlCRFpYSjBhV1pwWTJGMFpTQkJkWFJv
YjNKcGRIa3hKREFpQmdOVkJBc01HMGx3WTJSdUlFTmxjblJwClptbGpZWFJsSUdGMWRHaHZjbWww
ZVRFV01CUUdBMVVFQXd3TlNYQmpaRzRnVW05dmRDQkRRVEFlRncweE56RXgKTVRZeU1ETTRNVFZh
RncweU56RXhNVFF5TURNNE1UVmFNSUdNTVFzd0NRWURWUVFHRXdKVlV6RVJNQThHQTFVRQpDQXdJ
UTI5c2IzSmhaRzh4SkRBaUJnTlZCQW9NRzBsd1kyUnVJRU5sY25ScFptbGpZWFJsSUVGMWRHaHZj
bWwwCmVURWtNQ0lHQTFVRUN3d2JTWEJqWkc0Z1EyVnlkR2xtYVdOaGRHVWdRWFYwYUc5eWFYUjVN
UjR3SEFZRFZRUUQKREJWSmNHTmtiaUJKYm5SbGNtMWxaR2xoZEdVZ1EwRXdnZ0lpTUEwR0NTcUdT
SWIzRFFFQkFRVUFBNElDRHdBdwpnZ0lLQW9JQ0FRREtreFE5aGRNdFN3Zk5FQnFHUEFLenU3Q1Bo
SCt4dTNrTDRlK3JKMDlFQisrZnhHOTA1bXdNClltOHMxR0MvaTFWMWhQYk9rK3pMV1hjY1podEM0
OUJ0dHNEQlZSbWFndDRxNmVlRDEyQXh6VkJveldqNFluRnkKWkUwNENpWmxBSU4zcU40VG5OVC9P
M2l5dm1qNThRRElGVk81MVlOU3JyN2oyZFFSTHBveVM2czg3a3B3OUE2VAoyNEwxcExrbUZ1QWdD
TE1GUEc1SFpXeVpTU3RwWE96T2M3TElUWlFYUXp1bndMYXpOOVo0QXo4ellDOWlsTzZWCnROTk5j
K1k3TXBGclRhRUZGU3NNMitSRWV4dVB0Q090VC9aRWNPd1A4ODRUNkFDY1VTVHY4NmlFL0VGT1Bn
WmgKdlRLMkQwTnphdDNaVHNjNU4ydnZOMGVabTZDT25WQXZZTndyVFdHNHYzWVV0THIvUEVvRm05
bXRQZFNBK1RzaQpMa0dGUmp3QW9BbmhWaWVGQUZ1bFFuc3diWkFhSlJjL3hTN0JKdnRzLzNKOWk3
bDFvcHF1MEVibTZMOGpMZUh2CnAvb3AxVEQ4TElRa2NwN0dkc1hrNExZSDZWb3BhTk9pOHlvYUVm
S1d3clhoeEJkQ2hISGxZQ2NmZWN5RDhPMzkKOGV1b0dRMHppL2VhQ3IzTUhZVTErYTIrNVRSUGR6
SHgvbC8zVjF3eFg5aG1PR2Nyd1RDdGNZUnpSanphMlVsRQpVYWNtd0JxR3lmYTFHbi9pT3RvOXlx
dFJ2V2xQR3E5ekpOUXpOQko2aVdSN2d5YjB0ODIwanBkdG1GeU5CNVB2CkpjZUh5VnlpeC83c0JV
QjVHd01PMGlHbkxQcFk5SXg4N25ZUDNpdE4zamtxYWhUclQ3dEh1d0lEQVFBQm8yWXcKWkRBZEJn
TlZIUTRFRmdRVUlWN2lxdDMreFlhMUZ4L25ZenhkMWh4dk53MHdId1lEVlIwakJCZ3dGb0FVV2Nk
SQpCUjZ0M044dzd5VXQvVE9mWHdtdWJBZ3dFZ1lEVlIwVEFRSC9CQWd3QmdFQi93SUJBREFPQmdO
VkhROEJBZjhFCkJBTUNBWVl3RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUtiZDhuS0ExaGx4Q1Rj
elRnOVBIU2RxSFZFWVdxUkQKZHpyb09uejBLY0lnSkY0bmNSYU4zbm5hNzR3SW0zbkpEL1d1d2RG
OGRWYnBENFlhZVpLaFFnbVhtWWdnVkgwbAo4SUlrZEMzcVlOVERqNHhoZzdxWkMzcHpLdU5JZzd5
SHYzQXhyL3JiOUJ2cFVKeGJLby9hbG4vRzZQcmJOb1F4CjR4YjZ4bEs2V2NZT1JCL1VnekY2aktN
NUNCWGFVdDJRUW1XdHM3TWdub056WjRreHNZYnV3ajNyT1B6SElnUlkKWXdKdCtraDV0bHloSFBr
Z2p3VVRiTWcvTXYvb21melozK1hNQkQyWjY5ZU5kOVJIRTlYUzZBZE1LdERpcHJoZwpUTHUrOXRP
Vmx5ZFFqWmszS3NxRnY1Ny8xVHdLcnQ5SkZqWGsvUFBRYWxYb1Azb3lrb2l2aWFjZDRNeWlCbWFX
Cmk1V2J0Y3hHUDV2d2FObjI2N09HbGliZlFtUWdBMVZJZ0RxZnRyUXdlNlQvWndPS1hxRnFiNTZn
cWRXNFZkT1UKNWJHTjB1TFhlbXJpdU5Rb2xUcFZtaE51MGJOSTVXQ2lZN2MzMFhOMk01ZWJjQldE
SE9yNERPQ2crbW1YK3RKVApHaTFqdmZtQTZMT1lueVE4eEFWMDgwdVZlMzVOUTEwZTlteTZMdVhN
aVE3dWJ2SFVTNU53YVFDRnBMREl2YmhmCjJYVmM0Nk5kMFYvYm9sTTIzQVFSbHp1bzBJU3NEU2VX
U1FWTSt1UVorcExlOUc1bjJNeHFGZDg0VnU0ODlhdWMKMDJJOWNIaU5haFd5WFB0dm52TVd3SW1J
QjMwQ1loaDVCeUdXbmVEbzNkNVpqMkhBQTZGd082N0xsSWhTZTdxegpheER0Y000bk8zUUgKLS0t
LS1FTkQgQ0VSVElGSUNBVEUtLS0tLQotLS0tLUJFR0lOIENFUlRJRklDQVRFLS0tLS0KTUlJR0Vq
Q0NBL3FnQXdJQkFnSUpBTDJWMWhuTFd1emdNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1JR1ZNUXN3Q1FZ
RApWUVFHRXdKVlV6RVJNQThHQTFVRUNBd0lRMjlzYjNKaFpHOHhEekFOQmdOVkJBY01Ca1JsYm5a
bGNqRWtNQ0lHCkExVUVDZ3diU1hCalpHNGdRMlZ5ZEdsbWFXTmhkR1VnUVhWMGFHOXlhWFI1TVNR
d0lnWURWUVFMREJ0SmNHTmsKYmlCRFpYSjBhV1pwWTJGMFpTQmhkWFJvYjNKcGRIa3hGakFVQmdO
VkJBTU1EVWx3WTJSdUlGSnZiM1FnUTBFdwpIaGNOTVRjeE1URTJNakF4T1RJMFdoY05NamN4TVRF
ME1qQXhPVEkwV2pDQmxURUxNQWtHQTFVRUJoTUNWVk14CkVUQVBCZ05WQkFnTUNFTnZiRzl5WVdS
dk1ROHdEUVlEVlFRSERBWkVaVzUyWlhJeEpEQWlCZ05WQkFvTUcwbHcKWTJSdUlFTmxjblJwWm1s
allYUmxJRUYxZEdodmNtbDBlVEVrTUNJR0ExVUVDd3diU1hCalpHNGdRMlZ5ZEdsbQphV05oZEdV
Z1lYVjBhRzl5YVhSNU1SWXdGQVlEVlFRRERBMUpjR05rYmlCU2IyOTBJRU5CTUlJQ0lqQU5CZ2tx
CmhraUc5dzBCQVFFRkFBT0NBZzhBTUlJQ0NnS0NBZ0VBMkhzTGtubTl3eDFWMzN6dHRNV24xc0ZF
SnlMTk1lZVcKaEEwZHdlS0NocURWdzVsQXlJaWQwbHFqbmpJRWl2bEFNOHpkMk5BUFpUSEZHOUJU
M3VuZHVOY0RjV3VSZ2gzMwptbHlzc0VoYXJXdEU2VTdsenc4Uk1HV2t5V1FrSjJFMFN4akUvVm1L
UWpIMy80QWs4U2hoVFpGS0VadXdlUmRnCmoyMThxVWVtc09WK0VOVHNuR1V4b0FQcHI1Y0dHbzRp
Z2ZPOXRwSTFnN1BXbmtZclZGdHdzUG95MkNLeHFoL0kKM0Y2N0VacTJ5Uk9CeTlnQmhDNUZ1Wmh1
dmdwdHZTYWtOUXkvdys5YnVSZmZzaGI2OXdMc2JxUWF5aGpsTHFjUQpYa2NoNHk2Y0c4WWxnK2hy
cXptRzN2RzRxMDNZRno4UHVNZC91TnhlSlFBYlpKL3NYaDdqVlFOSml0Z2k2b095CmoxeFZ6TktV
cWIxTXgyNURoQnkxMjUvbDh3R0Jkc0NSVVlSNnI3ME5BWlJuZGprUUtWNGwzdFRRNCtvMnJkamQK
Qy8yU2syL2JIVmJ5dE1xNEl6KzFzT3NJVVYzTGErUitEV2NMQTZVb25wRW9jTmI4YnQ1YnU0SmhS
RHJuSndQaQo0alJhSCtLZ3hwQlRJbFozYUZuU3MwSTN3a1BPZ0EzRHJJNkxXeGd0M0JmUHVLOURQ
bWExcmM0QjFjd1NRVE5rCnl5VVkyYm9Yc25MTU1QZWJGM2Rqd2J0WUVkcUFhaXh4enpRQ1ZZbGhm
Q0VRanFjZVNkNlRKMlNkekg0STNKN3gKT3NMejFMai85MmFEZWVxNzdxL2M5RmlkbmZOZTc1ekI3
Z1hOL2JUdGI1NVpjVFRPOFhBRXZlZ3hjcnc1Q25PMQpNbkZ3NnhEZisrVUNBd0VBQWFOak1HRXdI
UVlEVlIwT0JCWUVGRm5IU0FVZXJkemZNTzhsTGYwem4xOEpybXdJCk1COEdBMVVkSXdRWU1CYUFG
Rm5IU0FVZXJkemZNTzhsTGYwem4xOEpybXdJTUE4R0ExVWRFd0VCL3dRRk1BTUIKQWY4d0RnWURW
UjBQQVFIL0JBUURBZ0dHTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElDQVFCK2FZV1poKy9nQ21GeQpv
STU4MkFDNSs1NGM3eGVMUm10UVBtRnBBQTJZNGlUU2VHeWp2Mmc1TTlhZ1E0SXh3dWU1RGx1UWxz
RGNEc0p3CmxzWFZIeUFsQ2Y2bkRiSk9wZjVtNml4ZnZCRitRekVlWGpVRzc0aVNDV3JBcFlGbXNO
c0NybHNNL0VQSm9ncXUKN2NnR3ZId1dZTmpQenV3b1UwQVdITzlGZ0liRHozTWMxNVpsSElRMlQv
aXh4eUJmcElQYzVEOTFrNHUvWldTYgpoTzZsc00xaUVBdnY5VG9VWGtLdHRDUTBmTjM0QndZWDQ3
QXhDNzF0U3Q5L1ZLSURJRVhqenFKSnBiSmRCR1NtCld1aVpIV3dCelZ3N2s1eStTK3dNNXZaL09N
dzB0aEwrSTUzem1reEdFZjN0TUg0UE9lS3ZYL3UxOW4wVEZrLzEKOTZjenRWbVRCNlc2N21iOXpa
OHlTZUp0RDl1WEU0NHljeGs5M2dUbmFya0lHZFRmOFhvbkNWbGtIR3BiNlRXTQo1RTdhV0NvVmJx
UVJrRjFycVB6YThoUkpHSDFFR2ttZGJrdXlVTUxtUTFGSks5Zi94RkljcmRCR3BNbHh5T29MCk1T
ZUZlYWdvU0xOcTV4emp1Y0YvNm1rQVVKOFVaQXhyU0dtYURFYVN0MS9xSlJIb3ZIY3QyYjVGV1pM
dVJ0MGcKZkRzR29LVDhJRnVhcFNzbUZFSW52RDBXS1hEU1BsVThyVHkxWkJEbENxcDNlTkg0dENm
TGFoazd2TEEwSmJmNwpVTHJ0RG5yTmZqeFVZRmNGMFJjYXVSMWRsclFRUkFzMzFRT3pYK0pPcjRn
VUc3eVhhM2o2NG54Mkc1TXBMZ1NvCk9FVWpmYWtLNzErVi9IYlF0NDc3elI0azdjUmJpQT09Ci0t
LS0tRU5EIENFUlRJRklDQVRFLS0tLS0K`

	rootCA = `
-----BEGIN CERTIFICATE-----
MIIGEjCCA/qgAwIBAgIJAL2V1hnLWuzgMA0GCSqGSIb3DQEBCwUAMIGVMQswCQYD
VQQGEwJVUzERMA8GA1UECAwIQ29sb3JhZG8xDzANBgNVBAcMBkRlbnZlcjEkMCIG
A1UECgwbSXBjZG4gQ2VydGlmaWNhdGUgQXV0aG9yaXR5MSQwIgYDVQQLDBtJcGNk
biBDZXJ0aWZpY2F0ZSBhdXRob3JpdHkxFjAUBgNVBAMMDUlwY2RuIFJvb3QgQ0Ew
HhcNMTcxMTE2MjAxOTI0WhcNMjcxMTE0MjAxOTI0WjCBlTELMAkGA1UEBhMCVVMx
ETAPBgNVBAgMCENvbG9yYWRvMQ8wDQYDVQQHDAZEZW52ZXIxJDAiBgNVBAoMG0lw
Y2RuIENlcnRpZmljYXRlIEF1dGhvcml0eTEkMCIGA1UECwwbSXBjZG4gQ2VydGlm
aWNhdGUgYXV0aG9yaXR5MRYwFAYDVQQDDA1JcGNkbiBSb290IENBMIICIjANBgkq
hkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2HsLknm9wx1V33zttMWn1sFEJyLNMeeW
hA0dweKChqDVw5lAyIid0lqjnjIEivlAM8zd2NAPZTHFG9BT3unduNcDcWuRgh33
mlyssEharWtE6U7lzw8RMGWkyWQkJ2E0SxjE/VmKQjH3/4Ak8ShhTZFKEZuweRdg
j218qUemsOV+ENTsnGUxoAPpr5cGGo4igfO9tpI1g7PWnkYrVFtwsPoy2CKxqh/I
3F67EZq2yROBy9gBhC5FuZhuvgptvSakNQy/w+9buRffshb69wLsbqQayhjlLqcQ
Xkch4y6cG8Ylg+hrqzmG3vG4q03YFz8PuMd/uNxeJQAbZJ/sXh7jVQNJitgi6oOy
j1xVzNKUqb1Mx25DhBy125/l8wGBdsCRUYR6r70NAZRndjkQKV4l3tTQ4+o2rdjd
C/2Sk2/bHVbytMq4Iz+1sOsIUV3La+R+DWcLA6UonpEocNb8bt5bu4JhRDrnJwPi
4jRaH+KgxpBTIlZ3aFnSs0I3wkPOgA3DrI6LWxgt3BfPuK9DPma1rc4B1cwSQTNk
yyUY2boXsnLMMPebF3djwbtYEdqAaixxzzQCVYlhfCEQjqceSd6TJ2SdzH4I3J7x
OsLz1Lj/92aDeeq77q/c9FidnfNe75zB7gXN/bTtb55ZcTTO8XAEvegxcrw5CnO1
MnFw6xDf++UCAwEAAaNjMGEwHQYDVR0OBBYEFFnHSAUerdzfMO8lLf0zn18JrmwI
MB8GA1UdIwQYMBaAFFnHSAUerdzfMO8lLf0zn18JrmwIMA8GA1UdEwEB/wQFMAMB
Af8wDgYDVR0PAQH/BAQDAgGGMA0GCSqGSIb3DQEBCwUAA4ICAQB+aYWZh+/gCmFy
oI582AC5+54c7xeLRmtQPmFpAA2Y4iTSeGyjv2g5M9agQ4Ixwue5DluQlsDcDsJw
lsXVHyAlCf6nDbJOpf5m6ixfvBF+QzEeXjUG74iSCWrApYFmsNsCrlsM/EPJogqu
7cgGvHwWYNjPzuwoU0AWHO9FgIbDz3Mc15ZlHIQ2T/ixxyBfpIPc5D91k4u/ZWSb
hO6lsM1iEAvv9ToUXkKttCQ0fN34BwYX47AxC71tSt9/VKIDIEXjzqJJpbJdBGSm
WuiZHWwBzVw7k5y+S+wM5vZ/OMw0thL+I53zmkxGEf3tMH4POeKvX/u19n0TFk/1
96cztVmTB6W67mb9zZ8ySeJtD9uXE44ycxk93gTnarkIGdTf8XonCVlkHGpb6TWM
5E7aWCoVbqQRkF1rqPza8hRJGH1EGkmdbkuyUMLmQ1FJK9f/xFIcrdBGpMlxyOoL
MSeFeagoSLNq5xzjucF/6mkAUJ8UZAxrSGmaDEaSt1/qJRHovHct2b5FWZLuRt0g
fDsGoKT8IFuapSsmFEInvD0WKXDSPlU8rTy1ZBDlCqp3eNH4tCfLahk7vLA0Jbf7
ULrtDnrNfjxUYFcF0RcauR1dlrQQRAs31QOzX+JOr4gUG7yXa3j64nx2G5MpLgSo
OEUjfakK71+V/HbQt477zR4k7cRbiA==
-----END CERTIFICATE-----
`
)

func TestVerifyAndEncodeCertificate(t *testing.T) {

	// should fail bad base64 data
	dat, err := verifyAndEncodeCertificate(BadData, "")
	if err == nil {
		t.Errorf("Unexpected result, there should have been a base64 decoding failure")
	}

	// should fail, can't verify self signed cert
	dat, err = verifyAndEncodeCertificate(SelfSigneCertOnly, rootCA)
	if err == nil {
		t.Errorf("Unexpected result, a certificate verification error should have occured")
	}

	// should pass
	dat, err = verifyAndEncodeCertificate(GoodTLSKeys, rootCA)
	if err != nil {
		t.Errorf("Test failure: %s", err)
	}

	pemCerts := make([]byte, base64.StdEncoding.EncodedLen(len(dat)))
	_, err = base64.StdEncoding.Decode(pemCerts, []byte(dat))
	if err != nil {
		t.Errorf("Test failure: bad retrun value from verifyAndEncodeCertificate(): %v", err)
	}

	certs := strings.SplitAfter(string(pemCerts), "-----END CERTIFICATE-----")
	length := len(certs) - 1
	if length != 2 {
		t.Errorf("Test failure: expected 2 certs from verifyAndEncodeCertificate(), got: %d ", length)
	}
}

// tests the generateSSLCertificate() function.
// verifys the proper creation of a CSR and private key and that
// both are encoded properly.
func TestGenerateDeliveryServiceSSLKeysCertificate(t *testing.T) {
	// test data
	const cdn = "over-the-top"
	const deliveryservice = "test-ds"
	const businessUnit = "IPCDN"
	const city = "Denver"
	const organization = "Comcast"
	const hostName = "foobar.test-ds.com"
	const country = "US"
	const state = "CO"
	const version = 1

	dsSslKeys := tc.DeliveryServiceSSLKeys{
		CDN:             cdn,
		DeliveryService: deliveryservice,
		BusinessUnit:    businessUnit,
		City:            city,
		Organization:    organization,
		Hostname:        hostName,
		Country:         country,
		State:           state,
		Key:             deliveryservice,
		Version:         version,
	}

	// test generating a certificate request and privte key
	err := generateSSLCertificate(&dsSslKeys)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// The Crt should be base64 encoded, test that it can be decoded.
	crtPem := make([]byte, base64.StdEncoding.EncodedLen(len(dsSslKeys.Certificate.Crt)))
	_, err = base64.StdEncoding.Decode(crtPem, []byte(dsSslKeys.Certificate.Crt))
	if err != nil {
		t.Errorf("unexpected error while trying to base64 decode the Crt")
	}

	// test that the Crt may be pem decoded.
	crtBlock, _ := pem.Decode(crtPem)
	if crtBlock == nil {
		t.Errorf("failed to decode PEM block containing the Crt")
	}

	// verify that the Crt is valid by parsing and reading it.
	crt, err := x509.ParseCertificate(crtBlock.Bytes)
	if err != nil {
		t.Errorf("unexpected error parsing the Crt: %v", err)
	}
	if crt.Subject.CommonName != hostName {
		t.Errorf("expected '%s' in crt.DNSNames[0], got %v", hostName, crt.Subject.CommonName)
	}

	// CSR should be base64 encoded, test that it can be decoded.
	csrPem := make([]byte, base64.StdEncoding.EncodedLen(len(dsSslKeys.Certificate.CSR)))
	_, err = base64.StdEncoding.Decode(csrPem, []byte(dsSslKeys.Certificate.CSR))
	if err != nil {
		t.Errorf("unexpected error while trying to base64 decode the CSR")
	}

	// test that the CSR may be pem decoded.
	csrBlock, _ := pem.Decode(csrPem)
	if csrBlock == nil {
		t.Errorf("failed to decode PEM block containing the CSR")
	}

	// verify that the CSR is valid by parsing and reading it.
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		t.Errorf("unexpected error parsing the CSR: %v", err)
	}
	if csr.Subject.Country[0] != country {
		t.Errorf("expected '%s' in csr.Subject.Country, got %v", country, csr.Subject.Country[0])
	}
	if csr.Subject.Organization[0] != organization {
		t.Errorf("expected '%s' in csr.Subject.Organization, got %v", organization, csr.Subject.Organization[0])
	}
	if csr.Subject.OrganizationalUnit[0] != businessUnit {
		t.Errorf("expected '%s' in csr.Subject.OrganizationalUnit, got %v", businessUnit, csr.Subject.OrganizationalUnit[0])
	}
	if csr.Subject.Locality[0] != city {
		t.Errorf("expected '%s' in csr.Subject.Locality, got %v", city, csr.Subject.Locality[0])
	}
	if csr.Subject.Province[0] != state {
		t.Errorf("expected '%s' in csr.Subject.Province, got %v", state, csr.Subject.Province[0])
	}
	if csr.Subject.CommonName != hostName {
		t.Errorf("expected '%s' in csr.Subject.CommonName, got %v", hostName, csr.Subject.CommonName)
	}

	// The private key should be base64 encoded, test that it can be decoded.
	keyPem := make([]byte, base64.StdEncoding.EncodedLen(len(dsSslKeys.Certificate.Key)))
	_, err = base64.StdEncoding.Decode(keyPem, []byte(dsSslKeys.Certificate.Key))
	if err != nil {
		t.Errorf("unexpected error while trying to base64 decode the CSR")
	}

	// test that the private key may be pem decoded.
	keyBlock, _ := pem.Decode(keyPem)
	if keyBlock == nil {
		t.Errorf("failed to decode PEM block containing the CSR")
	}

	// test that the private key may be properly parsed.
	_, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		t.Errorf("failed to parse the generated private key: %v", err)
	}
}
