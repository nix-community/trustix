# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/api.proto](#api_api-proto)
    - [GetLogAuditProofRequest](#trustix_api-v1-GetLogAuditProofRequest)
    - [GetLogConsistencyProofRequest](#trustix_api-v1-GetLogConsistencyProofRequest)
    - [GetLogEntriesRequest](#trustix_api-v1-GetLogEntriesRequest)
    - [GetMapValueRequest](#trustix_api-v1-GetMapValueRequest)
    - [KeyValuePair](#trustix_api-v1-KeyValuePair)
    - [Log](#trustix_api-v1-Log)
    - [Log.MetaEntry](#trustix_api-v1-Log-MetaEntry)
    - [LogEntriesResponse](#trustix_api-v1-LogEntriesResponse)
    - [LogHeadRequest](#trustix_api-v1-LogHeadRequest)
    - [LogSigner](#trustix_api-v1-LogSigner)
    - [LogsRequest](#trustix_api-v1-LogsRequest)
    - [LogsResponse](#trustix_api-v1-LogsResponse)
    - [MapValueResponse](#trustix_api-v1-MapValueResponse)
    - [ProofResponse](#trustix_api-v1-ProofResponse)
    - [SparseCompactMerkleProof](#trustix_api-v1-SparseCompactMerkleProof)
    - [ValueRequest](#trustix_api-v1-ValueRequest)
    - [ValueResponse](#trustix_api-v1-ValueResponse)
  
    - [Log.LogModes](#trustix_api-v1-Log-LogModes)
    - [LogSigner.KeyTypes](#trustix_api-v1-LogSigner-KeyTypes)
  
    - [LogAPI](#trustix_api-v1-LogAPI)
    - [NodeAPI](#trustix_api-v1-NodeAPI)
  
- [rpc/rpc.proto](#rpc_rpc-proto)
    - [DecideRequest](#trustix_rpc-v1-DecideRequest)
    - [DecisionResponse](#trustix_rpc-v1-DecisionResponse)
    - [EntriesResponse](#trustix_rpc-v1-EntriesResponse)
    - [EntriesResponse.EntriesEntry](#trustix_rpc-v1-EntriesResponse-EntriesEntry)
    - [FlushRequest](#trustix_rpc-v1-FlushRequest)
    - [FlushResponse](#trustix_rpc-v1-FlushResponse)
    - [LogValueDecision](#trustix_rpc-v1-LogValueDecision)
    - [LogValueResponse](#trustix_rpc-v1-LogValueResponse)
    - [SubmitRequest](#trustix_rpc-v1-SubmitRequest)
    - [SubmitResponse](#trustix_rpc-v1-SubmitResponse)
  
    - [SubmitResponse.Status](#trustix_rpc-v1-SubmitResponse-Status)
  
    - [LogRPC](#trustix_rpc-v1-LogRPC)
    - [RPCApi](#trustix_rpc-v1-RPCApi)
  
- [schema/loghead.proto](#schema_loghead-proto)
    - [LogHead](#trustix_schema-v1-LogHead)
  
- [schema/logleaf.proto](#schema_logleaf-proto)
    - [LogLeaf](#trustix_schema-v1-LogLeaf)
  
- [schema/mapentry.proto](#schema_mapentry-proto)
    - [MapEntry](#trustix_schema-v1-MapEntry)
  
- [schema/queue.proto](#schema_queue-proto)
    - [SubmitQueue](#trustix_schema-v1-SubmitQueue)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_api-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/api.proto



<a name="trustix_api-v1-GetLogAuditProofRequest"></a>

### GetLogAuditProofRequest
Get log audit proof for a given tree


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| Index | [uint64](#uint64) | required | Tree node index |
| TreeSize | [uint64](#uint64) | required | Tree size (proof reference) |






<a name="trustix_api-v1-GetLogConsistencyProofRequest"></a>

### GetLogConsistencyProofRequest
Get a consistency proof between two given log sizes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| FirstSize | [uint64](#uint64) | required | From tree size |
| SecondSize | [uint64](#uint64) | required | To tree size |






<a name="trustix_api-v1-GetLogEntriesRequest"></a>

### GetLogEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| Start | [uint64](#uint64) | required | Get entries from |
| Finish | [uint64](#uint64) | required | Get entries to |






<a name="trustix_api-v1-GetMapValueRequest"></a>

### GetMapValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| Key | [bytes](#bytes) | required | Map key |
| MapRoot | [bytes](#bytes) | required | Map root hash to derive proof from |






<a name="trustix_api-v1-KeyValuePair"></a>

### KeyValuePair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required | Map key |
| Value | [bytes](#bytes) | required | Map value |






<a name="trustix_api-v1-Log"></a>

### Log



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Mode | [Log.LogModes](#trustix_api-v1-Log-LogModes) | required |  |
| Protocol | [string](#string) | required |  |
| Signer | [LogSigner](#trustix_api-v1-LogSigner) | required |  |
| Meta | [Log.MetaEntry](#trustix_api-v1-Log-MetaEntry) | repeated |  |






<a name="trustix_api-v1-Log-MetaEntry"></a>

### Log.MetaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | optional |  |
| value | [string](#string) | optional |  |






<a name="trustix_api-v1-LogEntriesResponse"></a>

### LogEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Leaves | [trustix_schema.v1.LogLeaf](#trustix_schema-v1-LogLeaf) | repeated |  |






<a name="trustix_api-v1-LogHeadRequest"></a>

### LogHeadRequest
Request a signed head for a given log


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |






<a name="trustix_api-v1-LogSigner"></a>

### LogSigner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| KeyType | [LogSigner.KeyTypes](#trustix_api-v1-LogSigner-KeyTypes) | required |  |
| Public | [string](#string) | required |  |






<a name="trustix_api-v1-LogsRequest"></a>

### LogsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Protocols | [string](#string) | repeated | Allow to filter logs response based on the protocol identifier |






<a name="trustix_api-v1-LogsResponse"></a>

### LogsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Logs | [Log](#trustix_api-v1-Log) | repeated |  |






<a name="trustix_api-v1-MapValueResponse"></a>

### MapValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Value | [bytes](#bytes) | required | Note that the Value field is actually a MapEntry but we need to return the marshaled version as that&#39;s what the proof is created from |
| Proof | [SparseCompactMerkleProof](#trustix_api-v1-SparseCompactMerkleProof) | required |  |






<a name="trustix_api-v1-ProofResponse"></a>

### ProofResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Proof | [bytes](#bytes) | repeated |  |






<a name="trustix_api-v1-SparseCompactMerkleProof"></a>

### SparseCompactMerkleProof
Sparse merkle tree proof


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SideNodes | [bytes](#bytes) | repeated |  |
| NonMembershipLeafData | [bytes](#bytes) | optional |  |
| BitMask | [bytes](#bytes) | required |  |
| NumSideNodes | [uint64](#uint64) | required |  |






<a name="trustix_api-v1-ValueRequest"></a>

### ValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Digest | [bytes](#bytes) | required |  |






<a name="trustix_api-v1-ValueResponse"></a>

### ValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Value | [bytes](#bytes) | required |  |





 


<a name="trustix_api-v1-Log-LogModes"></a>

### Log.LogModes


| Name | Number | Description |
| ---- | ------ | ----------- |
| Log | 0 |  |



<a name="trustix_api-v1-LogSigner-KeyTypes"></a>

### LogSigner.KeyTypes


| Name | Number | Description |
| ---- | ------ | ----------- |
| ed25519 | 0 |  |


 

 


<a name="trustix_api-v1-LogAPI"></a>

### LogAPI
LogAPI is a logical grouping for RPC methods that are specific to a given
log.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetHead | [LogHeadRequest](#trustix_api-v1-LogHeadRequest) | [.trustix_schema.v1.LogHead](#trustix_schema-v1-LogHead) | Get signed head |
| GetLogConsistencyProof | [GetLogConsistencyProofRequest](#trustix_api-v1-GetLogConsistencyProofRequest) | [ProofResponse](#trustix_api-v1-ProofResponse) |  |
| GetLogAuditProof | [GetLogAuditProofRequest](#trustix_api-v1-GetLogAuditProofRequest) | [ProofResponse](#trustix_api-v1-ProofResponse) |  |
| GetLogEntries | [GetLogEntriesRequest](#trustix_api-v1-GetLogEntriesRequest) | [LogEntriesResponse](#trustix_api-v1-LogEntriesResponse) |  |
| GetMapValue | [GetMapValueRequest](#trustix_api-v1-GetMapValueRequest) | [MapValueResponse](#trustix_api-v1-MapValueResponse) |  |
| GetMHLogConsistencyProof | [GetLogConsistencyProofRequest](#trustix_api-v1-GetLogConsistencyProofRequest) | [ProofResponse](#trustix_api-v1-ProofResponse) |  |
| GetMHLogAuditProof | [GetLogAuditProofRequest](#trustix_api-v1-GetLogAuditProofRequest) | [ProofResponse](#trustix_api-v1-ProofResponse) |  |
| GetMHLogEntries | [GetLogEntriesRequest](#trustix_api-v1-GetLogEntriesRequest) | [LogEntriesResponse](#trustix_api-v1-LogEntriesResponse) |  |


<a name="trustix_api-v1-NodeAPI"></a>

### NodeAPI
NodeAPI is a logical grouping for RPC methods that are for the entire node
rather than individual logs.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Logs | [LogsRequest](#trustix_api-v1-LogsRequest) | [LogsResponse](#trustix_api-v1-LogsResponse) | Get a list of all logs published by this node |
| GetValue | [ValueRequest](#trustix_api-v1-ValueRequest) | [ValueResponse](#trustix_api-v1-ValueResponse) | Get values by their content-address |

 



<a name="rpc_rpc-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## rpc/rpc.proto



<a name="trustix_rpc-v1-DecideRequest"></a>

### DecideRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |
| Protocol | [string](#string) | required |  |






<a name="trustix_rpc-v1-DecisionResponse"></a>

### DecisionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Decision | [LogValueDecision](#trustix_rpc-v1-LogValueDecision) | required |  |
| Mismatches | [LogValueResponse](#trustix_rpc-v1-LogValueResponse) | repeated | Non-matches (hash mismatch) |
| Misses | [string](#string) | repeated | Full misses (log ids missing log entry entirely) |






<a name="trustix_rpc-v1-EntriesResponse"></a>

### EntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |
| Entries | [EntriesResponse.EntriesEntry](#trustix_rpc-v1-EntriesResponse-EntriesEntry) | repeated |  |






<a name="trustix_rpc-v1-EntriesResponse-EntriesEntry"></a>

### EntriesResponse.EntriesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | optional |  |
| value | [trustix_schema.v1.MapEntry](#trustix_schema-v1-MapEntry) | optional |  |






<a name="trustix_rpc-v1-FlushRequest"></a>

### FlushRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |






<a name="trustix_rpc-v1-FlushResponse"></a>

### FlushResponse







<a name="trustix_rpc-v1-LogValueDecision"></a>

### LogValueDecision



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogIDs | [string](#string) | repeated |  |
| Digest | [bytes](#bytes) | required |  |
| Confidence | [int32](#int32) | required |  |
| Value | [bytes](#bytes) | required |  |






<a name="trustix_rpc-v1-LogValueResponse"></a>

### LogValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Digest | [bytes](#bytes) | required |  |






<a name="trustix_rpc-v1-SubmitRequest"></a>

### SubmitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Items | [trustix_api.v1.KeyValuePair](#trustix_api-v1-KeyValuePair) | repeated |  |






<a name="trustix_rpc-v1-SubmitResponse"></a>

### SubmitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [SubmitResponse.Status](#trustix_rpc-v1-SubmitResponse-Status) | required |  |





 


<a name="trustix_rpc-v1-SubmitResponse-Status"></a>

### SubmitResponse.Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| OK | 0 |  |


 

 


<a name="trustix_rpc-v1-LogRPC"></a>

### LogRPC
RPCApi are &#34;private&#34; rpc methods for an instance related to a specific log.
This should only be available to trusted parties.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetHead | [.trustix_api.v1.LogHeadRequest](#trustix_api-v1-LogHeadRequest) | [.trustix_schema.v1.LogHead](#trustix_schema-v1-LogHead) |  |
| GetLogEntries | [.trustix_api.v1.GetLogEntriesRequest](#trustix_api-v1-GetLogEntriesRequest) | [.trustix_api.v1.LogEntriesResponse](#trustix_api-v1-LogEntriesResponse) |  |
| Submit | [SubmitRequest](#trustix_rpc-v1-SubmitRequest) | [SubmitResponse](#trustix_rpc-v1-SubmitResponse) |  |
| Flush | [FlushRequest](#trustix_rpc-v1-FlushRequest) | [FlushResponse](#trustix_rpc-v1-FlushResponse) |  |


<a name="trustix_rpc-v1-RPCApi"></a>

### RPCApi
RPCApi are &#34;private&#34; rpc methods for an instance.
This should only be available to trusted parties.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Logs | [.trustix_api.v1.LogsRequest](#trustix_api-v1-LogsRequest) | [.trustix_api.v1.LogsResponse](#trustix_api-v1-LogsResponse) | Get a list of all logs published/subscribed by this node |
| Decide | [DecideRequest](#trustix_rpc-v1-DecideRequest) | [DecisionResponse](#trustix_rpc-v1-DecisionResponse) | Decide on an output for key based on the configured decision method |
| GetValue | [.trustix_api.v1.ValueRequest](#trustix_api-v1-ValueRequest) | [.trustix_api.v1.ValueResponse](#trustix_api-v1-ValueResponse) | Get values by their content-address |

 



<a name="schema_loghead-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/loghead.proto



<a name="trustix_schema-v1-LogHead"></a>

### LogHead
Log


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogRoot | [bytes](#bytes) | required |  |
| TreeSize | [uint64](#uint64) | required |  |
| MapRoot | [bytes](#bytes) | required |  |
| MHRoot | [bytes](#bytes) | required |  |
| MHTreeSize | [uint64](#uint64) | required |  |
| Signature | [bytes](#bytes) | required | Aggregate signature |





 

 

 

 



<a name="schema_logleaf-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/logleaf.proto



<a name="trustix_schema-v1-LogLeaf"></a>

### LogLeaf
Leaf value of a merkle tree


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | optional |  |
| ValueDigest | [bytes](#bytes) | optional |  |
| LeafDigest | [bytes](#bytes) | required |  |





 

 

 

 



<a name="schema_mapentry-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/mapentry.proto



<a name="trustix_schema-v1-MapEntry"></a>

### MapEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Digest | [bytes](#bytes) | required | Value digest of tree node |
| Index | [uint64](#uint64) | required | Index of value in log |





 

 

 

 



<a name="schema_queue-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/queue.proto



<a name="trustix_schema-v1-SubmitQueue"></a>

### SubmitQueue
This type is internal only and not guaranteed stable


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Min | [uint64](#uint64) | required | Min is the _current_ (last popped) ID |
| Max | [uint64](#uint64) | required | Max is the last written item |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

