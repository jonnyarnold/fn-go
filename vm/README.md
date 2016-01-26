# fnvm

That's pronounced "FunVum".

In an effort to learn how VMs work, I'm going to build one. I'm going to borrow heavily from [Lua's VM](http://luaforge.net/docman/83/98/ANoFrillsIntroToLua51VMInstructions.pdf), which seems to be both efficient and well-documented.

## The Machine

fnvm consists of the following machine:

- **Registers** hold integers. These integers can represent constant numbers, booleans or characters. They can also hold addresses to more complex data structures.
- **The Program Counter** keeps track of where in the instruction list we currently are.
- **The Constant List** stores any primitive constants - numbers, booleans or strings.

## Input

Input is provided as a bytecode list.

## Open Questions
