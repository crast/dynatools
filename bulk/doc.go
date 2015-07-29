/*
package bulk implements parallel bulk loading of data into Dynamo.

The bulk package is mainly for use in writing thousands or more of records into
a DynamoDB table, such as the situation when you are doing a bulk conversion
from another datastore, or restore from backup, etc.

it does this with BulkWriter, which will simultaneously execute batch
operations in a number of goroutines, with automatic scale-back when the table
starts returning provisioned throughput errors, and back-pressure on writes.
*/
package bulk
