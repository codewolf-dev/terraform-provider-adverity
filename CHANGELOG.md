## 0.2.5

### FIXES:

- Fixed provider checksum mismatch

## 0.2.4

### FIXES:

- Fixed an inconsistent result when modifying the schedule(s) of a datastream caused by an Adverity API quirk that silently drops cron fields when datatype is included in the same request (0.2.2 only fixed this partially)
- Fixed the conversion of float numbers in parameters
- Validate instance URL format, schema and host

## 0.2.3

### FIXES:

- Fixed resource imports by using all required IDs for subsequent reads on the API

## 0.2.2

### FIXES:

- Fixed an inconsistent result when modifying the schedule(s) of a datastream (this time for real)
- Fixed broken deployments of datastreams with no schedule definition
- Validate conflicting attributes in schedule configurations

## 0.2.1

### FIXES:

- Fixed an inconsistent result when modifying the schedule(s) of a datastream

## 0.2.0

### FEATURES:

Resource:
- Authorization (reimplemented and renamed from Connection to align with UI terminology)

### DEPRECATIONS

Resource:
- Connection (deprecated and replaced by Authorization; will be removed in a future release)

## 0.1.2

### FIXES:

- Fixed an inconsistent result when modifying the datatype of a datastream

## 0.1.1

### FIXES:

- Fixed provider checksum mismatch

## 0.1.0

### FEATURES:

Resources:
- Workspace
- Connection
- Datastream
- Destination
- Destination Mapping

Data Sources:
- Connection Type
- Datastream Type
- Destination Type