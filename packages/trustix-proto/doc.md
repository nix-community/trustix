# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/api.proto](#api_api-proto)
    - [GetLogAuditProofRequest](#trustix-GetLogAuditProofRequest)
    - [GetLogConsistencyProofRequest](#trustix-GetLogConsistencyProofRequest)
    - [GetLogEntriesRequest](#trustix-GetLogEntriesRequest)
    - [GetMapValueRequest](#trustix-GetMapValueRequest)
    - [KeyValuePair](#trustix-KeyValuePair)
    - [Log](#trustix-Log)
    - [Log.MetaEntry](#trustix-Log-MetaEntry)
    - [LogEntriesResponse](#trustix-LogEntriesResponse)
    - [LogHeadRequest](#trustix-LogHeadRequest)
    - [LogSigner](#trustix-LogSigner)
    - [LogsRequest](#trustix-LogsRequest)
    - [LogsResponse](#trustix-LogsResponse)
    - [MapValueResponse](#trustix-MapValueResponse)
    - [ProofResponse](#trustix-ProofResponse)
    - [SparseCompactMerkleProof](#trustix-SparseCompactMerkleProof)
    - [ValueRequest](#trustix-ValueRequest)
    - [ValueResponse](#trustix-ValueResponse)
  
    - [Log.LogModes](#trustix-Log-LogModes)
    - [LogSigner.KeyTypes](#trustix-LogSigner-KeyTypes)
  
    - [LogAPI](#trustix-LogAPI)
    - [NodeAPI](#trustix-NodeAPI)
  
- [rpc/rpc.proto](#rpc_rpc-proto)
    - [DecideRequest](#trustix-DecideRequest)
    - [DecisionResponse](#trustix-DecisionResponse)
    - [EntriesResponse](#trustix-EntriesResponse)
    - [EntriesResponse.EntriesEntry](#trustix-EntriesResponse-EntriesEntry)
    - [FlushRequest](#trustix-FlushRequest)
    - [FlushResponse](#trustix-FlushResponse)
    - [LogValueDecision](#trustix-LogValueDecision)
    - [LogValueResponse](#trustix-LogValueResponse)
    - [SubmitRequest](#trustix-SubmitRequest)
    - [SubmitResponse](#trustix-SubmitResponse)
  
    - [SubmitResponse.Status](#trustix-SubmitResponse-Status)
  
    - [LogRPC](#trustix-LogRPC)
    - [RPCApi](#trustix-RPCApi)
  
- [schema/loghead.proto](#schema_loghead-proto)
    - [LogHead](#-LogHead)
  
- [schema/logleaf.proto](#schema_logleaf-proto)
    - [LogLeaf](#-LogLeaf)
  
- [schema/mapentry.proto](#schema_mapentry-proto)
    - [MapEntry](#-MapEntry)
  
- [schema/queue.proto](#schema_queue-proto)
    - [SubmitQueue](#-SubmitQueue)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_api-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/api.proto



<a name="trustix-GetLogAuditProofRequest"></a>

### GetLogAuditProofRequest
Get log audit proof for a given tree


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| Index | [uint64](#uint64) | required | Tree node index |
| TreeSize | [uint64](#uint64) | required | Tree size (proof reference) |






<a name="trustix-GetLogConsistencyProofRequest"></a>

### GetLogConsistencyProofRequest
Get a consistency proof between two given log sizes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| FirstSize | [uint64](#uint64) | required | From tree size |
| SecondSize | [uint64](#uint64) | required | To tree size |






<a name="trustix-GetLogEntriesRequest"></a>

### GetLogEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| Start | [uint64](#uint64) | required | Get entries from |
| Finish | [uint64](#uint64) | required | Get entries to |






<a name="trustix-GetMapValueRequest"></a>

### GetMapValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |
| Key | [bytes](#bytes) | required | Map key |
| MapRoot | [bytes](#bytes) | required | Map root hash to derive proof from |






<a name="trustix-KeyValuePair"></a>

### KeyValuePair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required | Map key |
| Value | [bytes](#bytes) | required | Map value |






<a name="trustix-Log"></a>

### Log



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Mode | [Log.LogModes](#trustix-Log-LogModes) | required |  |
| Protocol | [string](#string) | required |  |
| Signer | [LogSigner](#trustix-LogSigner) | required |  |
| Meta | [Log.MetaEntry](#trustix-Log-MetaEntry) | repeated |  |






<a name="trustix-Log-MetaEntry"></a>

### Log.MetaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | optional |  |
| value | [string](#string) | optional |  |






<a name="trustix-LogEntriesResponse"></a>

### LogEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Leaves | [LogLeaf](#LogLeaf) | repeated |  |






<a name="trustix-LogHeadRequest"></a>

### LogHeadRequest
Request a signed head for a given log


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required | Log identifier |






<a name="trustix-LogSigner"></a>

### LogSigner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| KeyType | [LogSigner.KeyTypes](#trustix-LogSigner-KeyTypes) | required |  |
| Public | [string](#string) | required |  |






<a name="trustix-LogsRequest"></a>

### LogsRequest







<a name="trustix-LogsResponse"></a>

### LogsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Logs | [Log](#trustix-Log) | repeated |  |






<a name="trustix-MapValueResponse"></a>

### MapValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Value | [bytes](#bytes) | required | Note that the Value field is actually a MapEntry but we need to return the marshaled version as that&#39;s what the proof is created from |
| Proof | [SparseCompactMerkleProof](#trustix-SparseCompactMerkleProof) | required |  |






<a name="trustix-ProofResponse"></a>

### ProofResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Proof | [bytes](#bytes) | repeated |  |






<a name="trustix-SparseCompactMerkleProof"></a>

### SparseCompactMerkleProof
Sparse merkle tree proof


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SideNodes | [bytes](#bytes) | repeated |  |
| NonMembershipLeafData | [bytes](#bytes) | required |  |
| BitMask | [bytes](#bytes) | required |  |
| NumSideNodes | [uint64](#uint64) | required |  |






<a name="trustix-ValueRequest"></a>

### ValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Digest | [bytes](#bytes) | required |  |






<a name="trustix-ValueResponse"></a>

### ValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Value | [bytes](#bytes) | required |  |





 


<a name="trustix-Log-LogModes"></a>

### Log.LogModes


| Name | Number | Description |
| ---- | ------ | ----------- |
| Log | 0 |  |



<a name="trustix-LogSigner-KeyTypes"></a>

### LogSigner.KeyTypes


| Name | Number | Description |
| ---- | ------ | ----------- |
| ed25519 | 0 |  |


 

 


<a name="trustix-LogAPI"></a>

### LogAPI
LogAPI is a logical grouping for RPC methods that are specific to a given
log.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetHead | [LogHeadRequest](#trustix-LogHeadRequest) | [.LogHead](#LogHead) | Get signed head |
| GetLogConsistencyProof | [GetLogConsistencyProofRequest](#trustix-GetLogConsistencyProofRequest) | [ProofResponse](#trustix-ProofResponse) |  |
| GetLogAuditProof | [GetLogAuditProofRequest](#trustix-GetLogAuditProofRequest) | [ProofResponse](#trustix-ProofResponse) |  |
| GetLogEntries | [GetLogEntriesRequest](#trustix-GetLogEntriesRequest) | [LogEntriesResponse](#trustix-LogEntriesResponse) |  |
| GetMapValue | [GetMapValueRequest](#trustix-GetMapValueRequest) | [MapValueResponse](#trustix-MapValueResponse) |  |
| GetMHLogConsistencyProof | [GetLogConsistencyProofRequest](#trustix-GetLogConsistencyProofRequest) | [ProofResponse](#trustix-ProofResponse) |  |
| GetMHLogAuditProof | [GetLogAuditProofRequest](#trustix-GetLogAuditProofRequest) | [ProofResponse](#trustix-ProofResponse) |  |
| GetMHLogEntries | [GetLogEntriesRequest](#trustix-GetLogEntriesRequest) | [LogEntriesResponse](#trustix-LogEntriesResponse) |  |


<a name="trustix-NodeAPI"></a>

### NodeAPI
NodeAPI is a logical grouping for RPC methods that are for the entire node
rather than individual logs.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Logs | [LogsRequest](#trustix-LogsRequest) | [LogsResponse](#trustix-LogsResponse) | Get a list of all logs published by this node |
| GetValue | [ValueRequest](#trustix-ValueRequest) | [ValueResponse](#trustix-ValueResponse) | Get values by their content-address |

 



<a name="rpc_rpc-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## rpc/rpc.proto



<a name="trustix-DecideRequest"></a>

### DecideRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |
| Protocol | [string](#string) | required |  |






<a name="trustix-DecisionResponse"></a>

### DecisionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Decision | [LogValueDecision](#trustix-LogValueDecision) | required |  |
| Mismatches | [LogValueResponse](#trustix-LogValueResponse) | repeated | Non-matches (hash mismatch) |
| Misses | [string](#string) | repeated | Full misses (log ids missing log entry entirely) |






<a name="trustix-EntriesResponse"></a>

### EntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |
| Entries | [EntriesResponse.EntriesEntry](#trustix-EntriesResponse-EntriesEntry) | repeated |  |






<a name="trustix-EntriesResponse-EntriesEntry"></a>

### EntriesResponse.EntriesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | optional |  |
| value | [MapEntry](#MapEntry) | optional |  |






<a name="trustix-FlushRequest"></a>

### FlushRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |






<a name="trustix-FlushResponse"></a>

### FlushResponse







<a name="trustix-LogValueDecision"></a>

### LogValueDecision



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogIDs | [string](#string) | repeated |  |
| Digest | [bytes](#bytes) | required |  |
| Confidence | [int32](#int32) | required |  |
| Value | [bytes](#bytes) | required |  |






<a name="trustix-LogValueResponse"></a>

### LogValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Digest | [bytes](#bytes) | required |  |






<a name="trustix-SubmitRequest"></a>

### SubmitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Items | [KeyValuePair](#trustix-KeyValuePair) | repeated |  |






<a name="trustix-SubmitResponse"></a>

### SubmitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [SubmitResponse.Status](#trustix-SubmitResponse-Status) | required |  |





 


<a name="trustix-SubmitResponse-Status"></a>

### SubmitResponse.Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| OK | 0 |  |


 

 


<a name="trustix-LogRPC"></a>

### LogRPC
RPCApi are &#34;private&#34; rpc methods for an instance related to a specific log.
This should only be available to trusted parties.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetHead | [LogHeadRequest](#trustix-LogHeadRequest) | [.LogHead](#LogHead) |  |
| GetLogEntries | [GetLogEntriesRequest](#trustix-GetLogEntriesRequest) | [LogEntriesResponse](#trustix-LogEntriesResponse) |  |
| Submit | [SubmitRequest](#trustix-SubmitRequest) | [SubmitResponse](#trustix-SubmitResponse) |  |
| Flush | [FlushRequest](#trustix-FlushRequest) | [FlushResponse](#trustix-FlushResponse) |  |


<a name="trustix-RPCApi"></a>

### RPCApi
RPCApi are &#34;private&#34; rpc methods for an instance.
This should only be available to trusted parties.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Logs | [LogsRequest](#trustix-LogsRequest) | [LogsResponse](#trustix-LogsResponse) | Get a list of all logs published/subscribed by this node |
| Decide | [DecideRequest](#trustix-DecideRequest) | [DecisionResponse](#trustix-DecisionResponse) | Decide on an output for key based on the configured decision method |
| GetValue | [ValueRequest](#trustix-ValueRequest) | [ValueResponse](#trustix-ValueResponse) | Get values by their content-address |

 



<a name="schema_loghead-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/loghead.proto



<a name="-LogHead"></a>

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



<a name="-LogLeaf"></a>

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



<a name="-MapEntry"></a>

### MapEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Digest | [bytes](#bytes) | required | Value digest of tree node |
| Index | [uint64](#uint64) | required | Index of value in log |





 

 

 

 



<a name="schema_queue-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/queue.proto



<a name="-SubmitQueue"></a>

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

