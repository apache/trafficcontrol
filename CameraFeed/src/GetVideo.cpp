#include <cstdint>
#include <regex>
#include <stdexcept>
#include <iostream>

#include "GetVideo.h"
#include "JPEGHandler.h"

// This class reads data from the Amcrest IPM-721S per section 4.1.4 of the
// AMCREST_CGI_SDK_API.pdf (obtainable from Amcrest's website).

GetVideo::GetVideo(const std::string &theURI, const std::string &theUserName,
                   const std::string &thePassword,
                   JPEGHandler &theJPEGCallback,
                   int theDebug) :
  myPassword(thePassword),
  myURI(theURI + "/cgi-bin/mjpg/video.cgi?channel=0&subtype=1"),
  myUserName(theUserName),
  myDebug(theDebug),
  myJPEGHandler(theJPEGCallback)
{
  myEasyHandle = curl_easy_init();
  if (NULL == myEasyHandle)
  {
    throw std::logic_error("Failed to get curl easy handle.");
  }

  setEasyOptions();
}

GetVideo::~GetVideo()
{
  if (myEasyHandle)
  {
    curl_easy_cleanup(myEasyHandle);
  }
}

void GetVideo::easyPerform()
{
  auto status = curl_easy_perform(myEasyHandle);
  if (CURLE_OK != status)
  {
    throw std::logic_error{std::string("Error during transfer: ") +
        curl_easy_strerror(status)};
  }
}

void GetVideo::setEasyOptions()
{
  if (myDebug)
  {
    curl_easy_setopt(myEasyHandle, CURLOPT_VERBOSE, 1L);
    curl_easy_setopt(myEasyHandle, CURLOPT_HEADER, 1L);
  }

  curl_easy_setopt(myEasyHandle, CURLOPT_URL, myURI.c_str());
  curl_easy_setopt(myEasyHandle, CURLOPT_HTTPAUTH, CURLAUTH_BASIC);
  curl_easy_setopt(myEasyHandle, CURLOPT_USERNAME, myUserName.c_str());
  curl_easy_setopt(myEasyHandle, CURLOPT_PASSWORD, myPassword.c_str());

  // Certificates on the camera don't verify, so skip it.
  // TODO: see what it would take to get cert to verify
  // TODO: ssl seems to break the camera (it reboots itself)
  // curl_easy_setopt(myEasyHandle, CURLOPT_SSL_VERIFYPEER, 0L);
  // curl_easy_setopt(myEasyHandle, CURLOPT_SSL_VERIFYHOST, 0L);

  curl_easy_setopt(myEasyHandle, CURLOPT_HEADERFUNCTION, curlHeaderCallback);
  curl_easy_setopt(myEasyHandle, CURLOPT_HEADERDATA, this);

  curl_easy_setopt(myEasyHandle, CURLOPT_WRITEFUNCTION, curlWriteCallback);
  curl_easy_setopt(myEasyHandle, CURLOPT_WRITEDATA, this);
}

void GetVideo::headerCallback(const std::string &theHeader)
{
  if (myBoundary.empty())
  {
    std::regex boundaryRE{"Content-Type:.*boundary=([A-Za-z0-9]*)"};
    std::smatch boundaryMatch;
    if (std::regex_search(theHeader, boundaryMatch, boundaryRE))
    {
      myBoundary = boundaryMatch[1];
    }
  }
}

size_t GetVideo::writeCallback(char *ptr, size_t size, size_t nmemb)
{
  uint32_t ii = 0;
  while (ii < nmemb)
  {
    if (ReadState::Header == myReadState)
    {
      if (ptr[ii] == '\r' || ptr[ii] == '\n')
      {
        myLineTerm++;
      }
      else
      {
        myHeaderData << ptr[ii];
      }

      if (myLineTerm == 2)
      {
        myLineTerm = 0;

        if (myDebug)
        {
          std::cerr << "Header line:" << myHeaderData.str() << std::endl;
        }

        if (myHeaderData.str().empty())
        {
          // Done with multi-part header.
          myReadState = ReadState::JPEG;
        }
        else
        {
          static std::regex contentLengthRE{"Content-Length: ([0-9]*)"};
          std::smatch contentLengthMatch;
          if (std::regex_search(myHeaderData.str(), contentLengthMatch,
                                contentLengthRE))
          {
            myJPEGSize = std::stoi(contentLengthMatch[1]);
            myJPEGData.reset(new char[myJPEGSize]);
            myJPEGBytesSaved = 0;
            if (myDebug)
            {
              std::cerr << "jpg size: " << myJPEGSize << std::endl;
            }
          }
          myHeaderData.str("");
        }
      }

      ++ii;
    }
    else if (ReadState::JPEG == myReadState)
    {
      // TODO: passes non-jpeg data to callback when using debug
      auto remainingBytes = nmemb - ii;
      if (remainingBytes > myJPEGSize)
      {
        remainingBytes = myJPEGSize;
      }
      std::memcpy(&myJPEGData.get()[myJPEGBytesSaved], ptr+ii, remainingBytes);
      myJPEGBytesSaved += remainingBytes;
      myJPEGSize -= remainingBytes;
      ii += remainingBytes;
      if (myJPEGSize == 0)
      {
        try
        {
          myJPEGHandler.handleJPEG(myJPEGData.get(), myJPEGBytesSaved);
        }
        catch (const std::exception &exception)
        {
          std::cerr << exception.what() << std::endl;
          return size + 1; // Signal to curl that something went amiss.
        }
        myReadState = ReadState::Trailer;
      }
    }
    else
    {
      if (ptr[ii] == '\r' || ptr[ii] == '\n')
      {
        myLineTerm++;
      }
      if (myLineTerm == 2)
      {
        myReadState = ReadState::Header;
      }
      ++ii;
    }
  }
  
  return size * ii;
}

size_t GetVideo::curlHeaderCallback(char *ptr, size_t size, size_t nmemb,
                                    void *userdata)
{
  std::string header(ptr, size * nmemb);
  reinterpret_cast<GetVideo*>(userdata)->headerCallback(header);
  return size * nmemb;
}

size_t GetVideo::curlWriteCallback(char *ptr, size_t size, size_t nmemb,
                               void *userdata)
{
  GetVideo *getVideo = reinterpret_cast<GetVideo*>(userdata);
  return getVideo->writeCallback(ptr, size, nmemb);
}
