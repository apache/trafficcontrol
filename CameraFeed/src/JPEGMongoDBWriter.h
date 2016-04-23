#pragma once

#include <string>

#include <mongocxx/client.hpp>
#include <mongocxx/collection.hpp>
#include <mongocxx/instance.hpp>

#include "JPEGHandler.h"

class JPEGMongoDBWriter : public JPEGHandler
{
  public:
  
  JPEGMongoDBWriter() = delete;

  JPEGMongoDBWriter(const std::string &theLocation, const std::string &theUser,
                    const std::string &theCamera, int theDebug);

  ~JPEGMongoDBWriter() = default;

  virtual void handleJPEG(const char *theJPEG, size_t theSize) override;

  private:

  const std::string myCamera;
  mongocxx::instance myInstance;
  mongocxx::client myClient;
  mongocxx::collection myCollection;
  const int myDebug;
  const std::string myUser;
};

