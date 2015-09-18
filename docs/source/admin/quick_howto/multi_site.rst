.. 
.. Copyright 2015 Comcast Cable Communications Management, LLC
.. 
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
.. 
..     http://www.apache.org/licenses/LICENSE-2.0
.. 
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
.. 

.. _rl-multi-site-origin-qht:

***************************
Configure Multi Site Origin
***************************

1) Create "cachegroups" for the origin locations, and assign the appropriate parent-child relationship between the mid cg's and org cgs (click the image to see full size):

.. image:: C5C4CD22-949A-48FD-8976-C673083E2177.png
	:scale: 100%
	:align: center

2) Create a profile to assign to each of the origins:

.. image:: 19BB6EC1-B6E8-4D22-BFA0-B7D6A9708B42.png
	:scale: 100%
	:align: center

3) Create server entries for the origination vips:

.. image:: D28614AA-9758-45ED-9EFD-3A284FC4218E.png
	:scale: 100%
	:align: center

4) assign the org servers to the delivery service that will have the multi site feature

.. image:: 066CEF4F-C1A3-4A89-8B52-4F72B0531367.png
	:scale: 100%
	:align: center

5) Check the multi-site check box in the delivery service screen and make sure that Content Routing Type is set to HTTP_LIVE_NATL:

.. image:: 71DA92BB-8E1E-4921-BC95-574E659812FF.png
	:scale: 100%
	:align: center

6) assign the parent.config location parameter to the MID profile

.. image:: D22DCAA3-18CC-48F4-965B-5312993F9820.png
	:scale: 100%
	:align: center


7) Configure the mid hdr_rewrite on the delivery service, example: ::

	cond %{REMAP_PSEUDO_HOOK} __RETURN__ set-config proxy.config.http.parent_origin.dead_server_retry_enabled 1 __RETURN__ set-config proxy.config.http.parent_origin.simple_retry_enabled 1 __RETURN__ set-config proxy.config.http.parent_origin.simple_retry_response_codes "400,404,412" __RETURN__ set-config proxy.config.http.parent_origin.dead_server_retry_response_codes "502,503" __RETURN__ set-config proxy.config.http.connect_attempts_timeout 2 __RETURN__ set-config proxy.config.http.connect_attempts_max_retries 2 __RETURN__ set-config proxy.config.http.connect_attempts_max_retries_dead_server 1 __RETURN__ set-config proxy.config.http.transaction_active_timeout_in 5 [L] __RETURN__

8) Turn on parent_proxy_routing in the MID profile.
