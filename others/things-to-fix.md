# Project To-Do List

## High Priority
- [ ] **Fix Address Issue**: Validate "Main Address" logic.
    - *Problem*: Currently, a user can potentially have 1 or more addresses marked as "main" (IsDefault).
    - *Goal*: Ensure only one address can be set as default/main per user. When a new address is set as default, the previous default should be unset.
