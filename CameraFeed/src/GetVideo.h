#pragma once

#include <cstdint>
#include <memory>
#include <string>
#include <sstream>
#include <curl/curl.h>

class JPEGHandler;

class GetVideo
{
  public:

  GetVideo() = delete;

  GetVideo(const std::string &theURI, const std::string &theUserName,
           const std::string &thePassword, JPEGHandler &theJPEGCallback,
           int theDebug);

  GetVideo(const GetVideo&) = delete;
  GetVideo(GetVideo&&) = delete;

  ~GetVideo();

  void easyPerform();

  protected:

  static size_t curlHeaderCallback(char *ptr, size_t size, size_t nmemb,
                                  void *userdata);

  static size_t curlWriteCallback(char *ptr, size_t size, size_t nmemb,
                                  void *userdata);

  void headerCallback(const std::string &theHeader);
  size_t writeCallback(char *ptr, size_t size, size_t nmemb);

  private:

  void setEasyOptions();

  const std::string myPassword;
  const std::string myURI;
  const std::string myUserName;

  const int myDebug;

  JPEGHandler &myJPEGHandler;

  enum class ReadState {Header, JPEG, Trailer};

  ReadState myReadState = ReadState::Header;
  std::stringstream myHeaderData;
  uint32_t myLineTerm = 0;

  std::unique_ptr<char[]> myJPEGData;
  size_t myJPEGBytesSaved = 0;
  uint32_t myJPEGSize = 0;

  std::string myBoundary;

  CURL *myEasyHandle = nullptr;
};
