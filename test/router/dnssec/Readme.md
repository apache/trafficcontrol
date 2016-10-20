DNSSEC Tests
============

Running the test

`ginkgo -- -ns=router-01.thecdn.example.com:53  -ds=ds-01.thecdn.example.com.`

Sample Output
```
Running Suite: Dnssec Suite
===========================
Random Seed: 1476984556
Will run 4 of 4 specs

2016/10/20 11:29:17 Nameserver router-01.thecdn.example.com:53
2016/10/20 11:29:17 DeliveryService ds-01.thecdn.example.com.
••••
Ran 4 of 4 Specs in 0.110 seconds
SUCCESS! -- 4 Passed | 0 Failed | 0 Pending | 0 Skipped PASS

Ginkgo ran 1 suite in 825.345359ms
Test Suite Passed
```