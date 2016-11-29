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
	var data GetDeliveryServiceResponse
	err := get(to, deliveryServicesEp(), &data)
	if err != nil {
		return nil, err
	}

	return data.Response, nil
}

// DeliveryService gets the DeliveryService for the ID it's passed
func (to *Session) DeliveryService(id string) (*DeliveryService, error) {
	var data GetDeliveryServiceResponse
	err := get(to, deliveryServiceEp(id), &data)
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
	err = post(to, deliveryServicesEp(), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdateDeliveryService updates the DeliveryService matching the ID it's passed with
// the DeliveryService it is passed
func (to *Session) UpdateDeliveryService(id string, ds *DeliveryService) (*DeliveryServiceResponse, error) {
	var data DeliveryServiceResponse
	jsonReq, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	err = put(to, deliveryServiceEp(id), jsonReq, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeleteDeliveryService deletes the DeliveryService matching the ID it's passed
func (to *Session) DeleteDeliveryService(id string) (*DeleteDeliveryServiceResponse, error) {
	var data DeleteDeliveryServiceResponse
	err := del(to, deliveryServiceEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// DeliveryServiceState gets the DeliveryServiceState for the ID it's passed
func (to *Session) DeliveryServiceState(id string) (*DeliveryServiceState, error) {
	var data DeliveryServiceStateResponse
	err := get(to, deliveryServiceStateEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceHealth gets the DeliveryServiceHealth for the ID it's passed
func (to *Session) DeliveryServiceHealth(id string) (*DeliveryServiceHealth, error) {
	var data DeliveryServiceHealthResponse
	err := get(to, deliveryServiceHealthEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceCapacity gets the DeliveryServiceCapacity for the ID it's passed
func (to *Session) DeliveryServiceCapacity(id string) (*DeliveryServiceCapacity, error) {
	var data DeliveryServiceCapacityResponse
	err := get(to, deliveryServiceCapacityEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceRouting gets the DeliveryServiceRouting for the ID it's passed
func (to *Session) DeliveryServiceRouting(id string) (*DeliveryServiceRouting, error) {
	var data DeliveryServiceRoutingResponse
	err := get(to, deliveryServiceRoutingEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceServer gets the DeliveryServiceServer
func (to *Session) DeliveryServiceServer(page, limit string) ([]DeliveryServiceServer, error) {
	var data DeliveryServiceServerResponse
	err := get(to, deliveryServiceServerEp(page, limit), &data)
	if err != nil {
		return nil, err
	}

	return data.Response, nil
}

// DeliveryServiceSSLKeysByID gets the DeliveryServiceSSLKeys by ID
func (to *Session) DeliveryServiceSSLKeysByID(id string) (*DeliveryServiceSSLKeys, error) {
	var data DeliveryServiceSSLKeysResponse
	err := get(to, deliveryServiceSSLKeysByIDEp(id), &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

// DeliveryServiceSSLKeysByHostname gets the DeliveryServiceSSLKeys by Hostname
func (to *Session) DeliveryServiceSSLKeysByHostname(hostname string) (*DeliveryServiceSSLKeys, error) {
	var data DeliveryServiceSSLKeysResponse
	err := get(to, deliveryServiceSSLKeysByHostnameEp(hostname), &data)
	if err != nil {
		return nil, err
	}

	return &data.Response, nil
}

func get(to *Session, endpoint string, respStruct interface{}) error {
	return makeReq(to, "GET", endpoint, nil, respStruct)
}

func post(to *Session, endpoint string, body []byte, respStruct interface{}) error {
	return makeReq(to, "POST", endpoint, body, respStruct)
}

func put(to *Session, endpoint string, body []byte, respStruct interface{}) error {
	return makeReq(to, "PUT", endpoint, body, respStruct)
}

func del(to *Session, endpoint string, respStruct interface{}) error {
	return makeReq(to, "DELETE", endpoint, nil, respStruct)
}

func makeReq(to *Session, method, endpoint string, body []byte, respStruct interface{}) error {
	resp, err := to.request(method, endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(respStruct); err != nil {
		return err
	}

	return nil
}
