#include <cstdlib>
#include <cstring>
#include <iostream>
#include <string>

#include "CommandLine.h"
#include "GetVideo.h"
#include "JPEGFileWriter.h"

//http://stackoverflow.com/questions/12797647/ffmpeg-how-to-mux-mjpeg-encoded-data-into-mp4-or-avi-container-c
//http://video.stackexchange.com/questions/7903/how-to-losslessly-encode-a-jpg-image-sequence-to-a-video-in-ffmpeg
int main(int argc, char **argv)
{
  try
  {
    CommandLine commandLine(argc, argv);
    JPEGFileWriter jpegWriter("/tmp/jpegs", "");
    GetVideo getVideo(commandLine.getURI(), commandLine.getUserName(),
                      commandLine.getPassword(), jpegWriter,
                      commandLine.getDebug());
    getVideo.easyPerform();
  }
  catch (const std::exception &exception)
  {
    std::cerr << argv[0] << ": " << exception.what() << std::endl;
  }

  return 0;
}
