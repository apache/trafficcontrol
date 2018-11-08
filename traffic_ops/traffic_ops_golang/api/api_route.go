package api

import "bytes"
import "net/http"
import "regexp"
import "strconv"
//import trops "github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang"
import log "github.com/apache/trafficcontrol/lib/go-log"

// todo - make the version string modular and accessible everywhere
const ServerString = "Traffic Operations/3.0.0";
const errorString = "Check the Traffic Ops log file(s) for details\n";


// CompiledRoute ...
type compiledRoute struct {
	Handler http.HandlerFunc
	Regex   *regexp.Regexp
	Params  []string
}


var AllRoutes *map[string][]compiledRoute;
// Writes a message indicating an internal server error back to the client (in plain text)
func errorResponse(writer http.ResponseWriter) {
	err := []byte(errorString);
	writer.Header().Set("Content-Length", strconv.Itoa(len(err)));
	writer.WriteHeader(http.StatusInternalServerError);
	writer.Write(err);
}

// Handles the disallowed request methods for this endpoint
func AvailableRoutesBadMethodHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Allow", "OPTIONS");
	writer.Header().Set("Server", ServerString);
	writer.WriteHeader(http.StatusMethodNotAllowed);
}

// Writes a list of all available API routes in a response to the client in plaintext
// (or writes an error message via errorResponse if something wicked happens)
func AvailableRoutesHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Server", ServerString);

	// Check for a CORS preflight request (all headers accepted, so immediate response given)
	if request.Header.Get("Origin") != "" &&
	   request.Header.Get("Access-Control-Request-Method") != "" &&
	   request.Header.Get("Access-Control-Request-Headers") != "" {

		writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS");
		writer.Header().Set("Access-Control-Allow-Headers", "*");
		writer.WriteHeader(http.StatusNoContent);
		return;
	}

	// ... otherwise, return a list of all supported methods/routes

	// if for some reason the variable didn't get set properly, return nothing
	if AllRoutes == nil {
		writer.Header().Set("Content-Length", "0");
		writer.WriteHeader(http.StatusNoContent);
		log.Warnln("API routes were requested, but weren't set!");
		return;
	}

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8");

	var body bytes.Buffer;
	contentLength := 0;
	for method, routes := range *AllRoutes {
		for _, route := range routes {
			n, err := body.WriteString(method);
			if err != nil {
				log.Errorf("Unable to append method to routes buffer: %s", err.Error());
				errorResponse(writer);
				return;
			}
			body.WriteRune(' ');
			contentLength += n + 1;

			n, err = body.WriteString(route.Regex.String());
			if err != nil {
				log.Errorf("Unable to append route to routes buffer: %s", err.Error());
				errorResponse(writer);
				return;
			}
			body.WriteRune('\n');
			contentLength += n +1;
		}
	}
}
