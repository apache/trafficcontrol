#include <chrono>
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

JPEGMongoDBWriter::JPEGMongoDBWriter(const std::string &theLocation,
                                     const std::string &theUser,
                                     const std::string &theCamera,
                                     int theDebug) :
  myCamera{theCamera},
  myInstance{},
  myClient{mongocxx::uri{"mongodb://" + theLocation}},
  // Database: CSCI5799, Collection: CameraFeed 
  myCollection{myClient["CSCI5799"]["CameraFeed"]},
  myDebug(theDebug),
  myUser{theUser}
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

  doc.append(bsoncxx::builder::basic::kvp("user", myUser));
  doc.append(bsoncxx::builder::basic::kvp("camera_id", myCamera));

  int64_t ms =
    std::chrono::duration_cast<std::chrono::milliseconds>(
      std::chrono::system_clock::now().time_since_epoch()).count();

  // The RetrieveVideo microservice works much better with raw milliseconds
  // since epoch (usual 1970 date).
  doc.append(bsoncxx::builder::basic::kvp(
               "msSinceEpoch", bsoncxx::types::b_int64{ms}));

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
