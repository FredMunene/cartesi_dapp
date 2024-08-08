## GO- 


## Overview
The program consists of a `go` script that interacts with a rollup server to handle requests for state advancement and inspection.

## Dependencies

```go
import(
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"dapp/rollups"
    )
```
+ *dapp/rollups* packages contains functions forconverting between hex and UTF-8 strings and handling the interaction with the rollup server.
+ *fmt* handles converting errors to strings for logging.
+ *log*  contains functions to output errors, information, and other messages.
+ *strconv* handles converting integers to strings.
+ *os* provides an interface to the operating system; i.e  reading environment variables .
+ *io* handles reading and writing data streams.
+ *encoding/json* provides functions for encoding and decoding JSON data.

## Implementation

1. Clone this repository.
```sh
git clone https://github.com/FredMunene/cartesi_dapp 
```

2. You can look at creating your own templates [here](https://docs.cartesi.io/cartesi-rollups/1.3/quickstart/).


### Variables
```go
var (
    Users []string
    toUpperTotal int
)

```
+ `users` will store the users names.
+ `toUpperTotal` is a count of valid transactions.

### Handle Advance Requests
Function receives and logs the advance request data.
+ Extracts the metadata, sender and payload.
+ Payload converts from hex to a string.
+ Checks if sentence is numeric.
    - if true, an error is reported  using **`SendReport` function**.
+ User is added to `users` slice and `toUpperTotal` counter is incremented.
+ Sentence is converted to uppercase and a notice using **`SendNotice` function**.
+ Returns *nil* error.