#pragma once

class CurlInit
{
  public:

  CurlInit(const CurlInit&) = delete;
  CurlInit(CurlInit&&) = delete;

  ~CurlInit();

  CurlInit& operator=(const CurlInit&) = delete;
  CurlInit& operator=(CurlInit&&) = delete;

  protected:

  CurlInit();

  private:

  static CurlInit myCurlInit;
};
