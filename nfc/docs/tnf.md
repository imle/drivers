### Type Name Format
```text
  TNF Value    Record Type
  ---------    -----------------------------------------
  0x00         Empty Record
               Indicates no type, id, or payload is associated with this NDEF Record.
               This record type is useful on newly formatted cards since every NDEF tag
               must have at least one NDEF Record.
               
  0x01         Well-Known Record
               Indicates the type field uses the RTD type name format.  This type name is used
               to stored any record defined by a Record Type Definition (RTD), such as storing
               RTD Text, RTD URIs, etc., and is one of the mostly frequently used and useful
               record types.
               
  0x02         MIME Media Record
               Indicates the payload is an intermediate or final chunk of a chunked NDEF Record
               
  0x03         Absolute URI Record
               Indicates the type field contains a value that follows the absolute-URI BNF
               construct defined by RFC 3986
               
  0x04         External Record
               Indicates the type field contains a value that follows the RTD external
               name specification
               
  0x05         Unknown Record
               Indicates the payload type is unknown
               
  0x06         Unchanged Record
               Indicates the payload is an intermediate or final chunk of a chunked NDEF Record
```