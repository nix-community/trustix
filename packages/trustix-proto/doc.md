# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/api.proto](#api/api.proto)
    - [GetLogAuditProofRequest](#trustix.GetLogAuditProofRequest)
    - [GetLogConsistencyProofRequest](#trustix.GetLogConsistencyProofRequest)
    - [GetLogEntriesRequest](#trustix.GetLogEntriesRequest)
    - [GetMapValueRequest](#trustix.GetMapValueRequest)
    - [KeyValuePair](#trustix.KeyValuePair)
    - [Log](#trustix.Log)
    - [Log.MetaEntry](#trustix.Log.MetaEntry)
    - [LogEntriesResponse](#trustix.LogEntriesResponse)
    - [LogHeadRequest](#trustix.LogHeadRequest)
    - [LogSigner](#trustix.LogSigner)
    - [LogsRequest](#trustix.LogsRequest)
    - [LogsResponse](#trustix.LogsResponse)
    - [MapValueResponse](#trustix.MapValueResponse)
    - [ProofResponse](#trustix.ProofResponse)
    - [SparseCompactMerkleProof](#trustix.SparseCompactMerkleProof)
    - [ValueRequest](#trustix.ValueRequest)
    - [ValueResponse](#trustix.ValueResponse)
  
    - [LogSigner.KeyTypes](#trustix.LogSigner.KeyTypes)
  
  
    - [LogAPI](#trustix.LogAPI)
    - [NodeAPI](#trustix.NodeAPI)
  

- [rpc/rpc.proto](#rpc/rpc.proto)
    - [DecisionResponse](#trustix.DecisionResponse)
    - [EntriesResponse](#trustix.EntriesResponse)
    - [EntriesResponse.EntriesEntry](#trustix.EntriesResponse.EntriesEntry)
    - [FlushRequest](#trustix.FlushRequest)
    - [FlushResponse](#trustix.FlushResponse)
    - [KeyRequest](#trustix.KeyRequest)
    - [LogValueDecision](#trustix.LogValueDecision)
    - [LogValueResponse](#trustix.LogValueResponse)
    - [SubmitRequest](#trustix.SubmitRequest)
    - [SubmitResponse](#trustix.SubmitResponse)
  
    - [SubmitResponse.Status](#trustix.SubmitResponse.Status)
  
  
    - [LogRPC](#trustix.LogRPC)
    - [RPCApi](#trustix.RPCApi)
  

- [schema/loghead.proto](#schema/loghead.proto)
    - [LogHead](#.LogHead)
  
  
  
  

- [schema/logleaf.proto](#schema/logleaf.proto)
    - [LogLeaf](#.LogLeaf)
  
  
  
  

- [schema/mapentry.proto](#schema/mapentry.proto)
    - [MapEntry](#.MapEntry)
  
  
  
  

- [schema/queue.proto](#schema/queue.proto)
    - [SubmitQueue](#.SubmitQueue)
  
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="api/api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/api.proto



<a name="trustix.GetLogAuditProofRequest"></a>

### GetLogAuditProofRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Index | [uint64](#uint64) | required |  |
| TreeSize | [uint64](#uint64) | required |  |






<a name="trustix.GetLogConsistencyProofRequest"></a>

### GetLogConsistencyProofRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| FirstSize | [uint64](#uint64) | required |  |
| SecondSize | [uint64](#uint64) | required |  |






<a name="trustix.GetLogEntriesRequest"></a>

### GetLogEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Start | [uint64](#uint64) | required |  |
| Finish | [uint64](#uint64) | required |  |






<a name="trustix.GetMapValueRequest"></a>

### GetMapValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Key | [bytes](#bytes) | required |  |
| MapRoot | [bytes](#bytes) | required |  |






<a name="trustix.KeyValuePair"></a>

### KeyValuePair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |
| Value | [bytes](#bytes) | required |  |






<a name="trustix.Log"></a>

### Log



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Signer | [LogSigner](#trustix.LogSigner) | required | required string Mode = 2; |
| Meta | [Log.MetaEntry](#trustix.Log.MetaEntry) | repeated |  |






<a name="trustix.Log.MetaEntry"></a>

### Log.MetaEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | optional |  |
| value | [string](#string) | optional |  |






<a name="trustix.LogEntriesResponse"></a>

### LogEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Leaves | [LogLeaf](#LogLeaf) | repeated |  |






<a name="trustix.LogHeadRequest"></a>

### LogHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |






<a name="trustix.LogSigner"></a>

### LogSigner



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| KeyType | [LogSigner.KeyTypes](#trustix.LogSigner.KeyTypes) | required |  |
| Public | [string](#string) | required |  |






<a name="trustix.LogsRequest"></a>

### LogsRequest







<a name="trustix.LogsResponse"></a>

### LogsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Logs | [Log](#trustix.Log) | repeated |  |






<a name="trustix.MapValueResponse"></a>

### MapValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Value | [bytes](#bytes) | required | Note that the Value field is actually a MapEntry but we need to return the marshaled version as that&#39;s what the proof is created from |
| Proof | [SparseCompactMerkleProof](#trustix.SparseCompactMerkleProof) | required |  |






<a name="trustix.ProofResponse"></a>

### ProofResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Proof | [bytes](#bytes) | repeated |  |






<a name="trustix.SparseCompactMerkleProof"></a>

### SparseCompactMerkleProof



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| SideNodes | [bytes](#bytes) | repeated |  |
| NonMembershipLeafData | [bytes](#bytes) | required |  |
| BitMask | [bytes](#bytes) | required |  |
| NumSideNodes | [uint64](#uint64) | required |  |






<a name="trustix.ValueRequest"></a>

### ValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Digest | [bytes](#bytes) | required |  |






<a name="trustix.ValueResponse"></a>

### ValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Value | [bytes](#bytes) | required |  |





 


<a name="trustix.LogSigner.KeyTypes"></a>

### LogSigner.KeyTypes


| Name | Number | Description |
| ---- | ------ | ----------- |
| ed25519 | 0 |  |


 

 


<a name="trustix.LogAPI"></a>

### LogAPI


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetHead | [LogHeadRequest](#trustix.LogHeadRequest) | [.LogHead](#LogHead) |  |
| GetLogConsistencyProof | [GetLogConsistencyProofRequest](#trustix.GetLogConsistencyProofRequest) | [ProofResponse](#trustix.ProofResponse) |  |
| GetLogAuditProof | [GetLogAuditProofRequest](#trustix.GetLogAuditProofRequest) | [ProofResponse](#trustix.ProofResponse) |  |
| GetLogEntries | [GetLogEntriesRequest](#trustix.GetLogEntriesRequest) | [LogEntriesResponse](#trustix.LogEntriesResponse) |  |
| GetMapValue | [GetMapValueRequest](#trustix.GetMapValueRequest) | [MapValueResponse](#trustix.MapValueResponse) |  |
| GetMHLogConsistencyProof | [GetLogConsistencyProofRequest](#trustix.GetLogConsistencyProofRequest) | [ProofResponse](#trustix.ProofResponse) |  |
| GetMHLogAuditProof | [GetLogAuditProofRequest](#trustix.GetLogAuditProofRequest) | [ProofResponse](#trustix.ProofResponse) |  |
| GetMHLogEntries | [GetLogEntriesRequest](#trustix.GetLogEntriesRequest) | [LogEntriesResponse](#trustix.LogEntriesResponse) |  |


<a name="trustix.NodeAPI"></a>

### NodeAPI


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Logs | [LogsRequest](#trustix.LogsRequest) | [LogsResponse](#trustix.LogsResponse) | Get map[LogID]Log |
| GetValue | [ValueRequest](#trustix.ValueRequest) | [ValueResponse](#trustix.ValueResponse) |  |

 



<a name="rpc/rpc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## rpc/rpc.proto



<a name="trustix.DecisionResponse"></a>

### DecisionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Decision | [LogValueDecision](#trustix.LogValueDecision) | required |  |
| Mismatches | [LogValueResponse](#trustix.LogValueResponse) | repeated | Non-matches (hash mismatch) |
| Misses | [string](#string) | repeated | Full misses (log ids missing log entry entirely) |






<a name="trustix.EntriesResponse"></a>

### EntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |
| Entries | [EntriesResponse.EntriesEntry](#trustix.EntriesResponse.EntriesEntry) | repeated |  |






<a name="trustix.EntriesResponse.EntriesEntry"></a>

### EntriesResponse.EntriesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) | optional |  |
| value | [MapEntry](#MapEntry) | optional |  |






<a name="trustix.FlushRequest"></a>

### FlushRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |






<a name="trustix.FlushResponse"></a>

### FlushResponse







<a name="trustix.KeyRequest"></a>

### KeyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | required |  |






<a name="trustix.LogValueDecision"></a>

### LogValueDecision



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogIDs | [string](#string) | repeated |  |
| Digest | [bytes](#bytes) | required |  |
| Confidence | [int32](#int32) | required |  |
| Value | [bytes](#bytes) | required |  |






<a name="trustix.LogValueResponse"></a>

### LogValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Digest | [bytes](#bytes) | required |  |






<a name="trustix.SubmitRequest"></a>

### SubmitRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogID | [string](#string) | required |  |
| Items | [KeyValuePair](#trustix.KeyValuePair) | repeated |  |






<a name="trustix.SubmitResponse"></a>

### SubmitResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [SubmitResponse.Status](#trustix.SubmitResponse.Status) | required |  |





 


<a name="trustix.SubmitResponse.Status"></a>

### SubmitResponse.Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| OK | 0 |  |


 

 


<a name="trustix.LogRPC"></a>

### LogRPC


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetHead | [LogHeadRequest](#trustix.LogHeadRequest) | [.LogHead](#LogHead) |  |
| GetLogEntries | [GetLogEntriesRequest](#trustix.GetLogEntriesRequest) | [LogEntriesResponse](#trustix.LogEntriesResponse) |  |
| Submit | [SubmitRequest](#trustix.SubmitRequest) | [SubmitResponse](#trustix.SubmitResponse) |  |
| Flush | [FlushRequest](#trustix.FlushRequest) | [FlushResponse](#trustix.FlushResponse) |  |


<a name="trustix.RPCApi"></a>

### RPCApi
TrustixRPC

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Logs | [LogsRequest](#trustix.LogsRequest) | [LogsResponse](#trustix.LogsResponse) | Get map[LogID]Log (all local logs) |
| Decide | [KeyRequest](#trustix.KeyRequest) | [DecisionResponse](#trustix.DecisionResponse) | Compare(inputHash) |
| GetValue | [ValueRequest](#trustix.ValueRequest) | [ValueResponse](#trustix.ValueResponse) | Get stored value by digest |

 



<a name="schema/loghead.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/loghead.proto



<a name=".LogHead"></a>

### LogHead



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| LogRoot | [bytes](#bytes) | required | Log |
| TreeSize | [uint64](#uint64) | required |  |
| MapRoot | [bytes](#bytes) | required | Map |
| MHRoot | [bytes](#bytes) | required | Map head fields |
| MHTreeSize | [uint64](#uint64) | required |  |
| Signature | [bytes](#bytes) | required | Aggregate signature |





 

 

 

 



<a name="schema/logleaf.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/logleaf.proto



<a name=".LogLeaf"></a>

### LogLeaf



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Key | [bytes](#bytes) | optional |  |
| ValueDigest | [bytes](#bytes) | optional |  |
| LeafDigest | [bytes](#bytes) | required |  |





 

 

 

 



<a name="schema/mapentry.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/mapentry.proto



<a name=".MapEntry"></a>

### MapEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Digest | [bytes](#bytes) | required |  |
| Index | [uint64](#uint64) | required |  |





 

 

 

 



<a name="schema/queue.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## schema/queue.proto



<a name=".SubmitQueue"></a>

### SubmitQueue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Min | [uint64](#uint64) | required | Min is the _current_ (last popped) ID |
| Max | [uint64](#uint64) | required | Max is the last written item |





 

 

 

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

