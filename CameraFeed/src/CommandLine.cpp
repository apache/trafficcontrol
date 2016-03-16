#include <getopt.h>
#include <iostream>
#include <stdexcept>

#include "CommandLine.h"


CommandLine::CommandLine(int argc, char **argv)
{
  extern int optind, opterr;
  extern char* optarg;

  opterr = 1;

  while (true)
  {
    enum Option
    {
      Debug,
      Help,
      Password,
      UserName,
    };

    static struct option options[] = {
      {"debug", optional_argument, 0, Debug},
      {"username", required_argument, 0, UserName},
      {"password", required_argument, 0, Password},
      {"help", no_argument, 0, Help},
      {0, 0, 0,  0 }
    };

    int optionIndex = 0;
    auto c = ::getopt_long(argc, argv, "",
                           options, &optionIndex);
    if (c == -1)
      break;

    switch (c) {
      case Debug:
        if (optarg)
        {
          myDebug = std::atoi(optarg);
        }
        else
        {
          myDebug = 1;
        }
        break;

      case Help:
        usage(argv[0]);
        std::exit(0);
        break;

      case Password:
        myPassword = optarg;
        break;

      case UserName:
        myUserName = optarg;
        break;

      default:
        usage(argv[0]);
        std::exit(1);
    }
  }

  if (argc - optind != 1)
  {
    throw std::runtime_error("No URI given.");
  }

  myURI = argv[optind];
}

int CommandLine::getDebug() const noexcept
{
  return myDebug;
}

std::string CommandLine::getPassword() const noexcept
{
  return myPassword;
}

std::string CommandLine::getUserName() const noexcept
{
  return myUserName;
}

std::string CommandLine::getURI() const noexcept
{
  return myURI;
}

void CommandLine::usage(const char *argv0)
{
  std::cerr << argv0 << " [OPTIONS] camera_host_name" << std::endl
            << std::endl
            << "Records data from the Amcrest IPM-721S cammer at the "
            << "given host name" << std::endl
            << std::endl
            << "OPTIONS" << std::endl
            << "  --debug[=level]       debug level (1 if no level specified)"
            << std::endl
            << "  --help                print this help and exit"
            << std::endl
            << "  --password=password   password for camera access"
            << std::endl
            << "  --username=username   User name for camera access"
            << std::endl;
}

