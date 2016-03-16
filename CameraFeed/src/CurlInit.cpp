#include <cstdlib>
#include <iostream>
#include <curl/curl.h>

#include "CurlInit.h"

CurlInit CurlInit::myCurlInit;

CurlInit::CurlInit()
{
  auto curlInitStatus = curl_global_init(CURL_GLOBAL_ALL);
  if (curlInitStatus)
  {
    std::cerr << "Critical curl initialization failure: "
              << curlInitStatus << std::endl;
    std::exit(1);
  }
}

CurlInit::~CurlInit()
{
  curl_global_cleanup();
}
