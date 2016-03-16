#pragma once

#include <cstdint>
#include <string>

#include "JPEGHandler.h"

class JPEGFileWriter : public JPEGHandler
{
  public:
  JPEGFileWriter() = delete;
  JPEGFileWriter(const std::string &theDirectory,
                 const std::string &theFileBase);

  virtual void handleJPEG(const char *theJPEG, size_t theSize) override;


  private:
  const std::string myDirectory;
  const std::string myFileBase;
  uint32_t myCount = 0;
};
