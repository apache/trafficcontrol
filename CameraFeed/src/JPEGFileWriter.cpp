#include <cerrno>
#include <cstring>
#include <fstream>
#include <iomanip>
#include <sstream>
#include <stdexcept>

#include "JPEGFileWriter.h"

JPEGFileWriter::JPEGFileWriter(const std::string &theDirectory,
                               const std::string &theFileBase) :
  myDirectory(theDirectory),
  myFileBase(theFileBase)
{
}

void JPEGFileWriter::handleJPEG(const char *theJPEG, size_t theSize)
{
  ++myCount;
  std::ostringstream fileName;
  fileName << myDirectory << "/" << myFileBase << std::setw(9)
           << std::setfill('0') << myCount << ".jpg";
  std::ofstream jpegFile(fileName.str(), std::ios::binary);
  auto thisErrno = errno;
  if (! jpegFile.is_open())
  {
    throw std::logic_error("Failed to open " + fileName.str() + ": " +
                           std::strerror(thisErrno));
  }

  jpegFile.write(theJPEG, theSize);
  jpegFile.close();
}
