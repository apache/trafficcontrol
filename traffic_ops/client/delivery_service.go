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

package client

import "encoding/json"

// DeliveryServices gets an array of DeliveryServices
func (to *Session) DeliveryServices() ([]DeliveryService, error) {
	var data DeliveryServiceResponse
	err := makeReq(to, deliveryServicesEp(), nil, &data)
	if err != nil {
		return nil, err
	}

	return data.Response, nil
}

// DeliveryService gets the DeliveryService for the ID it's passed
func (to *Session) DeliveryService(id string) (*DeliveryService, error) {
	var data DeliveryServiceResponse
	err := makeReq(to, deliveryServiceEp(id), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response[0], nil
}

// CreateDeliveryService creates the DeliveryService it's passed
func (to *Session) CreateDeliveryService(ds *DeliveryService) (*CreateDeliveryServiceResponse, error) {
	var data CreateDeliveryServiceResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	err = makeReq(to, deliveryServicesEp(), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeliveryServiceState gets the DeliveryServiceState for the ID it's passed
func (to *Session) DeliveryServiceState(id string) (*DeliveryServiceState, error) {
	var data DeliveryServiceStateResponse
	err := makeReq(to, deliveryServiceStateEp(id), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceHealth gets the DeliveryServiceHealth for the ID it's passed
func (to *Session) DeliveryServiceHealth(id string) (*DeliveryServiceHealth, error) {
	var data DeliveryServiceHealthResponse
	err := makeReq(to, deliveryServiceHealthEp(id), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceCapacity gets the DeliveryServiceCapacity for the ID it's passed
func (to *Session) DeliveryServiceCapacity(id string) (*DeliveryServiceCapacity, error) {
	var data DeliveryServiceCapacityResponse
	err := makeReq(to, deliveryServiceCapacityEp(id), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceRouting gets the DeliveryServiceRouting for the ID it's passed
func (to *Session) DeliveryServiceRouting(id string) (*DeliveryServiceRouting, error) {
	var data DeliveryServiceRoutingResponse
	err := makeReq(to, deliveryServiceRoutingEp(id), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceServer gets the DeliveryServiceServer
func (to *Session) DeliveryServiceServer(page, limit string) ([]DeliveryServiceServer, error) {
	var data DeliveryServiceServerResponse
	err := makeReq(to, deliveryServiceServerEp(page, limit), nil, &data)
	if err != nil {
		return nil, err
	}

	return data.Response, nil
}

// DeliveryServiceSSLKeysByID gets the DeliveryServiceSSLKeys by ID
func (to *Session) DeliveryServiceSSLKeysByID(id string) (*DeliveryServiceSSLKeys, error) {
	var data DeliveryServiceSSLKeysResponse
	err := makeReq(to, deliveryServiceSSLKeysByIDEp(id), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceSSLKeysByHostname gets the DeliveryServiceSSLKeys by Hostname
func (to *Session) DeliveryServiceSSLKeysByHostname(hostname string) (*DeliveryServiceSSLKeys, error) {
	var data DeliveryServiceSSLKeysResponse
	err := makeReq(to, deliveryServiceSSLKeysByHostnameEp(hostname), nil, &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

func makeReq(to *Session, endpoint string, body []byte, respStruct interface{}) error {
	resp, err := to.request(endpoint, body)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(respStruct); err != nil {
		return err
	}

	return nil
}
