#pragma once

#include <string>

class CommandLine
{
  public:
  CommandLine(int argc, char **argv);

  int getDebug() const noexcept;
  std::string getPassword() const noexcept;
  std::string getUserName() const noexcept;
  std::string getURI() const noexcept;

  void usage(const char *argv0);

  private:

  std::string myUserName;
  std::string myPassword;
  int myDebug = 0;
  std::string myURI;
};
