Testing
==
Integration tests for Traffic Control are here. For unit tests, see each individual component. Go code specifically has `_test.go` tests alongside the application code.

These tests and frameworks are currently experimental.

Frameworks
--
We have two frameworks for integrations tests: UI and API. Both are written in Go, use `go test`, and look at `environment.json` for test environment configuration.

Go Test
--
While designed for unit testing application code, `go test` does what we want, so it makes sense to reuse it for integration tests.

Test files are suffixed `_test.go` as required by `go test`. To run all integration tests, run `go test ./...` from the `infrastructure/test` directory.

Environment
--
All Traffic Control tests look at `test/environment.json` for service location and login information. There is a small helper Go library for loading the environment struct, which may be imported from `github.com/Comcast/traffic_control/infrastructure/test/environment`.

UI Tests
--
UI tests use Selenium, via `github.com/tebeka/selenium`.

A Selenium server must be running on `localhost` to run UI tests. The Selenium stanadlone server may be downloaded from `http://selenium-release.storage.googleapis.com/2.53/selenium-server-standalone-2.53.0.jar` and run with `java -jar selenium-server-standalone-2.53.0.jar`.

API Tests
--
API tests use the Go framework at `test/apitest`. We use a custom Go framework for API testing, because `github.com/tebeka/selenium` is not capable of HTTP methods other than GET.

Currently, the majority of functions in `apitest` are for comparing JSON endpoints. But it should be easy to extend for other formats.
