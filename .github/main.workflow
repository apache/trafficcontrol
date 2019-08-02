workflow "Daily Tests" {
  resolves = ["Go Test"]
  on = "schedule(0 0 * * *)"
}

action "Go Test" {
  uses = "./tests/traffic_ops/golang"
}
