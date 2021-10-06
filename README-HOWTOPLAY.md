# Nest M.U.D - How To Play

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)

## Characters

Characters are created and begin their adventure at the start of a dungeon.

## Actions

Characters are controlled with actions that are simple sentences.

### Movement Actions

Character can `move` from one location to another location using the `move [direction]` action.

**Syntax Examples:**

```text
move north
move southwest
move up
```

Possible directions:

- north
- northeast
- east
- southeast
- south
- southwest
- west
- northwest
- up
- down

### Look Actions

📝 _Unimplemented_

A character can simply `look` to get a description of their current location, look in a specific direction, at an object, at another character or at a monster that is in the current location.

When a character looks at an object that is not in their possession a brief description of the object will be provided. To look more closely at an object, a character must have the object equipped or stashed in their backpack.

When a character looks a direction a brief description of the location will be provided. To get a more detailed view of a location a character must move into that location.

**Syntax Examples:**

```text
look
look north
look chest
look goblin
```

### Equip, Stash and Drop Actions

📝 _Unimplemented_

A character may `equip` or `stash` an object into their backpack that is in their current location. A character may also `equip` and object that has been stashed in their backback. Any item that is currently equipped where that item would normally be equipped will be stashed if there is enough room in the characters backpack or dropped if there is not. A character may `drop` any item that is equipped or stashed in their backpack.

**Syntax Examples:**

```text
equip sword
stash gold hammer
drop dragon tongue
```

### Use Actions

📝 _Unimplemented_

A character may attempt `use` any object that is equipped, stashed or in their current location. Some items may only be used when equipped or stashed.

**Syntax Examples:**

```text
use sword
use door
```
