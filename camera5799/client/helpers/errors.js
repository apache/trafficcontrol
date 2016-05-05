/**
 * @author Ramiro Arenivar
 * For CSCI 5799
 */

// Local (client-only) collection
Errors = new Mongo.Collection(null);

throwError = function(message) {
  Errors.insert({message: message})
}