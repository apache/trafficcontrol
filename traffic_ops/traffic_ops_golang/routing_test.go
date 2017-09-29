package main

import (
	"net/http"
	"testing"

	"fmt"

	"context"
	"net/http/httptest"
	"bytes"
)

func TestCreateRouteMap(t *testing.T) {
	authBase := AuthBase{false, "secret", nil, func(handlerFunc http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(),"authWasCalled","true")
			handlerFunc(w,r.WithContext(ctx))
		}
	}}


	//expected := make(map[string][]PathHandler)
	//expected["path1"] = []PathHandler{PathHandler{}}
	//expected["path2"] = []PathHandler{PathHandler{}}

	PathOneHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authWasCalled := getAuthWasCalled(ctx)


		fmt.Fprintf(w, "%s %s", "path1", authWasCalled)
	}

	PathTwoHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authWasCalled := getAuthWasCalled(ctx)

		fmt.Fprintf(w, "%s %s", "path2", authWasCalled)
	}

	routes := []Route{{Version: 1.2, Method: http.MethodGet, Path: `path1`, Handler: PathOneHandler, RequiredPrivLevel: ServersPrivLevel}, {Version: 1.2, Method: http.MethodGet, Path: `path2`, Handler: PathTwoHandler, RequiredPrivLevel: 0}}

	routeMap := CreateRouteMap(routes, authBase)

	route1Handler := routeMap["GET"][0].Handler

	w := httptest.NewRecorder()
	r, err := http.NewRequest("", "/", nil)
	if err != nil {
		t.Error("Error creating new request")
	}

	route1Handler(w,r)

	if bytes.Compare(w.Body.Bytes(), []byte("path1 true")) != 0 {
		t.Errorf("Got: %s \nExpected to receive path1 true\n",w.Body.Bytes())
	}

	route2Handler := routeMap["GET"][1].Handler

	w = httptest.NewRecorder()

	route2Handler(w,r)

	if bytes.Compare(w.Body.Bytes(), []byte("path2 false")) != 0 {
		t.Errorf("Got: %s \nExpected to receive path2 false\n", w.Body.Bytes())
	}
}


func getAuthWasCalled(ctx context.Context) string {
	val := ctx.Value("authWasCalled")
	if val != nil {
		return val.(string)
	}
	return "false"
}