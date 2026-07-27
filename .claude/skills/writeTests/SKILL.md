# Repository Testing Standards

## Purpose

This skill defines the testing philosophy, conventions, structure and quality standards for every test in this repository.

The goal is complete consistency across the entire codebase.

Every test should look as though it was written by the same engineer regardless of whether it tests:

* HTTP Handlers
* gRPC Handlers
* gRPC Client Wrappers
* Services
* Repositories
* Middleware
* Interceptors
* Validators
* Utilities
* Helper Functions
* Domain Logic
* Storage Logic

Only the code under test changes.

The testing philosophy never changes.

---

# Before Writing Any Test

Before writing a single test case, understand the complete execution flow.

Do **not** generate tests immediately.

Instead, first analyse:

* the function being tested
* every dependency it calls
* every interface involved
* every mock implementation
* every shared helper
* every domain error
* every storage error
* every gRPC mapping
* every validation rule
* every rollback path
* every middleware/helper involved

The generated tests should reflect the actual execution flow, not assumptions.

---

# Examples Are Only Examples

Every example in this document is illustrative.

Examples are **never** rules.

Examples are intended only to explain the testing philosophy.

The actual implementation should always be derived from the repository currently being analysed.

Do **not** blindly copy examples.

Instead, inspect:

* the function being tested
* existing mocks
* storage layer
* service layer
* repository layer
* grpc layer
* grpc utilities
* domain errors
* validation helpers
* helper functions
* shared packages

The repository is always the source of truth.

---

# Understand The Entire Flow First

Before creating tests, determine:

```
Entry point

↓

Validation

↓

Business Rules

↓

Storage

↓

Repository

↓

gRPC

↓

External Dependencies

↓

Rollback

↓

Response
```

The actual flow depends entirely on the implementation.

Never assume the flow.

Discover it.

---

# Mirror Control Flow

Tests must follow the exact control flow of the implementation.

If the implementation contains early returns, the tests should appear in exactly the same order.

Someone reading the implementation should immediately find the corresponding test.

Example

```
Validation

↓

Authentication

↓

Business Rules

↓

Database

↓

gRPC

↓

Rollback

↓

Success
```

The tests should appear in that exact order.

Never reorder tests based on preference.

The implementation determines the order.

---

# Test Structure

Every exported function should have exactly one table-driven test.

Examples

```go
func TestUpdateUserHandler(t *testing.T)
```

```go
func TestRenameFileGrpcHandler(t *testing.T)
```

```go
func TestCreateUser(t *testing.T)
```

Each test should contain

* Test table
* Arrange
* Act
* Assert

No multiple test functions unless absolutely necessary.

---

# Table Structure

Every table should contain the following whenever applicable.

```
name

expectedCode

expectedData

expectedError
```

Additional fields should only represent

* Inputs
* Dependency behaviour
* Mock behaviour
* Expected captured values

Examples

```
body

userID

email

password

fileID

grpcErr

dbErr

jwtErr

expectedPhone

expectedFileName
```

Do not introduce unnecessary fields.

---

# Naming Convention

Every test name should describe behaviour.

Good

```
user updates password successfully

file rename fails due to duplicate filename

user registration fails due to invalid email

grpc rename fails because metadata is missing
```

Bad

```
success

error

duplicate

validation
```

Each name should explain

* subject
* action
* reason
* outcome

---

# Test Ordering

Tests should always appear in execution order.

General ordering

```
Happy Paths

↓

Authentication

↓

Authorization

↓

Input Parsing

↓

Validation

↓

Business Rules

↓

Storage

↓

Repository

↓

gRPC

↓

Rollback

↓

Internal Errors

↓

Edge Cases
```

Not every function contains every stage.

The function itself determines the order.

---

# Happy Paths

Always begin with successful execution paths.

Every logical success path deserves its own test.

Examples

```
password updated successfully

phone updated successfully

both password and phone updated successfully

file uploaded successfully

file renamed successfully
```

---

# Validation Coverage

Validation should be exhaustive.

If validation exists, every meaningful invalid input should exist.

Examples

UUID

```
missing

whitespace

invalid
```

Password

```
missing

spaces

short

numbers only

letters only

missing uppercase

missing special character
```

JSON

```
invalid

malformed

empty
```

Multiple validation failures should also be tested whenever possible.

---

# Business Rules

Every business rule deserves a dedicated test.

Examples

```
duplicate email

duplicate filename

same password

quota exceeded

already deleted

already exists
```

---

# Dependency Failures

Every dependency should have explicit failure coverage.

Examples

Database

```
duplicate

not found

internal error
```

gRPC

```
missing metadata

missing user id

not found

permission denied

internal error
```

Filesystem

```
disk full

permission denied

missing file
```

JWT

```
expired

missing

invalid
```

Only include dependency failures that actually exist.

---

# Rollback

If rollback exists, every rollback path should be tested.

Examples

```
grpc fails and rollback succeeds

grpc fails and rollback fails

database update fails before rollback
```

Rollback is business logic.

Treat it as such.

---

# Edge Cases

Every exported function should include meaningful edge cases.

Examples

```
multiple validation failures

empty collection

nil response

duplicate retry

already deleted

maximum length

minimum length
```

Only include edge cases relevant to the implementation.

---

# AAA Pattern

Every test should follow

```
Arrange

↓

Act

↓

Assert
```

Arrange

* create mocks
* create request
* initialise dependencies

Act

* invoke exactly one function

Assert

* verify status/result
* verify returned data
* verify returned errors
* verify captured mock state

---

# Assertion Order

Assertions should always appear in the same order.

1.

Status / Return value

2.

Returned data

3.

Returned error

4.

Captured mock state

Maintain this ordering throughout the repository.

---

# Mock Philosophy

Mocks represent domain behaviour.

Mocks are not simple stubs.

Mocks should behave like the real implementation.

Tests should never manually inject behaviour.

Bad

```go
service.Update = func(...) error
```

Good

```go
MockUserService{
    UpdateUserErr: mocks.DbOpDuplicatePhone,
}
```

---

# Mock Analysis (Mandatory)

Before writing tests, inspect every related mock.

Determine whether the mock correctly represents the real implementation.

Do not assume existing mocks are correct.

Analyse:

* returned objects
* returned errors
* captured arguments
* rollback behaviour
* gRPC mappings
* storage mappings
* repository mappings

If the mock does not reflect the actual implementation:

1. Identify the inconsistency.
2. Correct the mock.
3. Use the corrected mock.
4. Explain why the mock required modification.

The mock implementation must always match the actual application behaviour.

Never work around an incorrect mock by modifying the tests.

Fix the mock first.

---

# Mock Design

Mocks should expose enumerated behaviour.

Example

```go
type DbOperationError int
```

```
DbOpSuccess

DbOpNotFound

DbOpDuplicate

DbOpInternalError
```

Never use

```
errors.New(...)
```

inside tests.

Mocks translate enumerations into actual domain errors.

---

# Mock Behaviour

Mocks should return realistic domain objects.

Examples

Instead of

```
&User{}
```

return

```
ID

Email

Password

Phone

CreatedAt

UpdatedAt
```

Populate realistic values whenever the real implementation would.

---

# Capturing State

Mocks should capture all meaningful inputs.

Examples

```
ID

UserID

Email

Password

Request

Filename

Metadata
```

This allows verification after execution.

---

# Assertions On Mock State

Whenever appropriate, verify captured values.

Examples

```
UserID

Phone

Email

Filename

Request Body

Metadata
```

Verify values.

Do not merely verify that no error occurred.

---

# Repository Is The Source Of Truth

When writing tests, inspect:

* implementation
* interfaces
* service
* repository
* storage
* grpc clients
* grpc utils
* validation
* helper packages
* shared handlers
* mocks

Do not infer behaviour.

Derive behaviour from the repository.

---

# Transport Independence

These standards apply equally to

* HTTP
* gRPC
* Services
* Storage
* Repositories
* Middleware
* Validators
* Utilities
* Helpers

Only Arrange and Assert differ.

The testing philosophy remains identical.

---

# Expected Output

Whenever asked to generate tests:

* Produce exactly one table-driven test per exported function.
* Follow the implementation's execution order.
* Cover every success path.
* Cover every validation path.
* Cover every business rule.
* Cover every dependency failure.
* Cover every rollback path.
* Cover meaningful edge cases.
* Analyse and verify all related mocks before writing tests.
* Fix incorrect mocks when necessary before using them.
* Use realistic domain objects.
* Use enumerated mock behaviours.
* Capture important inputs.
* Verify returned values.
* Verify captured state.
* Maintain identical naming, formatting and structure across the repository.

The finished test should be indistinguishable from every other test in the repository.