Toy proof-of-concept repo using "command" pattern
---
(still under construction. Not much works yet... especially tests)

```
make run-dev
```
What is the "command" pattern?
---

The command pattern encapsulates all information necessary to perform a request
before it is executed.  This makes it easy to chain commands together, add
middleware, etc.

The project root
---
The types of "commands" for this project are simply interfaces stored in the
root directory under `commands.go`. Commands manipulate service objects,
defined in the same file.

Domain entities are saved in a file called `model.go`. The structs in this file
do not contain any struct tags as that would violate their purpose: to serve as
stable domain values. There are very few imports in this package and really
only to define aggregate data types (e.g. guregu/null package).

Error types are stored in `error.go`and wrap underlying errors, conforming to
Go's error convention (see https://pkg.go.dev/errors).

Dependency isolation
---
Sub-packages only import from the root package and 3rd-party packages (Never
from other sub-packages... utilities are an exception). Dependencies should be
specified as interfaces and defined as close as possible to the struct that
uses it. If code is thoughtfully organized by imports and function, this layout
effectively enforces dependency inversion.

> Note that per Dave Cheney's advice on fewer, larger packages, sub-folders
are not nested very deeply (Ideally only a single layer). And hold more, larger
files than in other languages.

Pipelining commands
---
Inside the `pipeline` package, there are custom pipeline types. Holding lists
of the command, a pipeline will implement the command interface, and terminate
execution early when an error is encountered.

Handlers
---
Handlers do three things: They transform bytes into service objects, get
results from a command and transform those results back into bytes.

Routers
---
Routers hand the request body to byte handlers for processing and return
response bodies, translating errors into http error codes.

Dependency Setup
---
Commands are stored in packages under the `cmd` folder but they hold very
little logic.  The heavy lifting of dependency setup and pipeline definition is
done in the `boot` package.

> This project follows the idiomatic go convention to define a struct literal
rather than a constructor when there is no constructor logic.

Utilities
---
Utility packages have meaningful names (config, logger, etc... never generic
names like base, common or util). They are allowed to define global, singleton
objects. A utility's sub-packages are allowed to register itself with its
parent at import.

Code generation
---
Mocks, enum stringer methods, etc. are all generated using the `go:generate`
comments. The test setup duplicates the logic in the `boot` package using mock
dependencies. Since we've defined our dependencies as interfaces, we can use
mock-generation libraries like `golang/mock` to generate that code.

Mocks are saved in the `test` directory.
