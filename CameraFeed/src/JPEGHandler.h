#pragma once

class JPEGHandler
{
  public:
  virtual ~JPEGHandler() = default;

  virtual void handleJPEG(const char *theJPEG, size_t theSize) = 0;
};
