/*
Provides DynamoDB streams support for Dynago.

This provides access to the DynamoDB streams API with all the conveniences
that one would expect from using Dynago: document unmarshaling to go types,
clean API's and a way to write applications simply.

This package provides two interfaces: a low-level Client for making the
individual API requests, and a Streamer for those writing applications that
utilize streams.
*/
package streams
