# Test writing conventions

Keep it Short, Precise, and Happily Sacrifice Some English Grammar

## BDD (Behavior Driven Development)

### Given (Replaced with `Test`)

- The initial state of the system
- The context in which the behavior occurs
- The state of the world before the event happens

### When (Replaced with `With`)

- The event that triggers the behavior
- The change that occurs in the system
- The user action that initiates the behavior

### Then (Replaced with `Expect`)

- The expected outcome
- The behavior that should be observed
- The state of the world after the event happens

## Example

- TestSum/With non-empty int slice/Expect no error
- TestSum/With non-empty int slice/Expect correct sum
- TestSum/With non-empty int slice/Expect input is unchanged
- TestSum/With empty int slice/Expect no error
- TestSum/With empty int slice/Expect zero
- TestSum/With empty int slice/Expect input is unchanged
