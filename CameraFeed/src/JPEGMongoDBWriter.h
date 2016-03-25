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

  JPEGMongoDBWriter(const std::string &theURI, int theDebug);

  ~JPEGMongoDBWriter() = default;

  virtual void handleJPEG(const char *theJPEG, size_t theSize) override;

  private:

  mongocxx::instance myInstance;
  mongocxx::client myClient;
  mongocxx::collection myCollection;
  const int myDebug;
};

