#include <iostream>
#include <stdexcept>
#include <string>
#include <sys/time.h>

#include <bsoncxx/builder/basic/document.hpp>
#include <bsoncxx/builder/basic/kvp.hpp>
#include <bsoncxx/types.hpp>
#include <mongocxx/exception/exception.hpp>
#include <mongocxx/uri.hpp>

#include "JPEGMongoDBWriter.h"

JPEGMongoDBWriter::JPEGMongoDBWriter(const std::string &theURI, int theDebug) :
  myInstance{},
  myClient{mongocxx::uri{theURI}},
  myCollection{myClient["CSCI5799"]["CameraFeed"]},
  myDebug(theDebug)
{
}

void JPEGMongoDBWriter::handleJPEG(const char *theJPEG, size_t theSize)
{
  // From mongo shell
  // > use CSCI5799
  // > show collections
  // > db.CameraFeed.find()  # list all JPEGs
  // > db.CameraFeed.help()  # lists commands
  // > db.CameraFeed.deleteMany({}) # delete all documents in collection

  auto doc = bsoncxx::builder::basic::document{};
  // TODO: need username
  doc.append(bsoncxx::builder::basic::kvp("user",
                                          "dummyUser"));
  // TODO: need camera id
  doc.append(bsoncxx::builder::basic::kvp("camera_id",
                                          "Camera123"));

  struct timeval currentTime;
  gettimeofday(&currentTime, 0);

  doc.append(bsoncxx::builder::basic::kvp("unix_time",
                 [&](bsoncxx::builder::basic::sub_document subdoc)
                 {
                   subdoc.append(
                     bsoncxx::builder::basic::kvp(
                       "seconds", bsoncxx::types::b_int64{currentTime.tv_sec}));
                   subdoc.append(
                     bsoncxx::builder::basic::kvp(
                       "microseconds",
                       bsoncxx::types::b_int64{currentTime.tv_usec}));
                 }));

  doc.append(bsoncxx::builder::basic::kvp(
               "jpeg",
               bsoncxx::types::b_binary{
                 bsoncxx::binary_sub_type::k_binary,
                   static_cast<uint32_t>(theSize), // compiler warning
                   reinterpret_cast<const uint8_t*>(theJPEG)
                   }));

  try
  {
    auto insertStatus = myCollection.insert_one(std::move(doc.view()));
    if (myDebug > 2)
    {
      std::cout << "Inserted value" << std::endl;
      if (insertStatus->inserted_id().type() == bsoncxx::type::k_oid)
      {
        std::cout << "  id: "
                  << insertStatus->inserted_id().get_oid().value.to_string()
                  << std::endl;
      }
    }
  }
  catch (const mongocxx::exception &exception)
  {
    std::string error{"Error inserting JPEG into MongoDB: "};
    error += exception.what();
    throw std::runtime_error(error);
  }
}
