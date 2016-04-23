#pragma once

#include <string>

class CommandLine
{
  public:
  CommandLine(int argc, char **argv);

  std::string getCamera() const noexcept;
  int getDebug() const noexcept;
  std::string getMongoLocation() const noexcept;
  std::string getPassword() const noexcept;

  /**
   * Returns the user of the system
   */
  std::string getUser() const noexcept;

  /**
   * Returns the user name for logging into the camera
   */
  std::string getUserName() const noexcept;
  std::string getURI() const noexcept;

  void usage(const char *argv0);

  private:

  std::string myCamera;
  std::string myMongoLocation;
  std::string myUser;
  std::string myUserName;
  std::string myPassword;
  int myDebug = 0;
  std::string myURI;
};
