workflow "Daily Tests" {
  resolves = ["Go Test"]
  on = "schedule(0 0 * * *)"
}

action "Go Test" {
  uses = "./traffic_ops/app/bin/tests/Dockerfile-golangtest"
}
