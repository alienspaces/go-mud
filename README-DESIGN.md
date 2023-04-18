# Go M.U.D - Design

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)
- [TODO](README-TODO.md)

## API Design

📝 _Incomplete_

The API is designed to emulate the experience of hand typed commands with responses that describe the result of that command.

Each time a command is submitted an action is created. Responses include the results of the command along with the names and descriptions of locations, characters, monsters and objects that are present at the location when the action occurred.

All other actions that occurred in the same location as the character since their previous command will also be returned. This way a character is given a running commentary of other actions that have occurred in that location.

All character commands, monster commands and other automatic location actions can occur at most, every two seconds.

A game client could choose to automatically submit a single `look` command for a character every two seconds when the player is detected as idle to retrieve any other actions that have occurred.

## Game Design

### Dungeons 

A dungeon is a set of interconnected locations. 

### Locations

A location may at most contain a total of 15 entities, encompassing monsters, characters and objects. A location may have 10 different entrances and exists including north, northeast, east, southeast, south, southwest, west, northwest, up and down.

**Location Diagram:**

```ascii
 --- --- --- --- --- 
|NW | 1 | N | 2 | NE|
 --- --- --- --- ---
| 3 | 4 | U | 5 | 6 |
 --- --- --- --- ---
| W | 7 | 8 | 9 | E |
 --- --- --- --- ---
|10 |11 | D | 12| 13|
 --- --- --- --- ---
|SW |14 | S | 15| SE|
 --- --- --- --- ---
```

Mundane objects that have no significant special abilities and aren't a quest object will be stacked at a location while other objects, monsters and characters will each contribute to the number of entities at that location.

When a location has reached the maximum number of 15 entities, no new monsters, characters or objects can enter, spawn or otherwise be added to that location. 

## Objects

An object is any thing in a dungeon that can be looked at, used, equipped or stashed. An object might be a weapon, a piece of armour, a potion, clothing, food, a key, a book, virtually anything. 

Not all objects are able to be stashed in a character backpack or equipped.

Quest objects are bound to a character when equipped or stashed and cannot be passed to another monster or character. When a quest object is dropped, it will disappear after 3 turns and must be collected again if needed to continue the quest. No other monster or character, other than the character the quest object is bound to, may pick up a dropped quest object.

## Monsters and Characters

Monsters and characters move around a dungeon, interacting with each other and objects they might discover.

Monsters come in many different shapes and sizes and demeanour. Characters are controlled by players while monsters are free willed entities that might stand guard at a specific location or move around the dungeon. Monsters will chase characters when provoked.

When a monster or character reaches it's demise, they fall to the ground and may be looted by other monsters or characters at that location. After 3 turns, the monster or character will disappear, along with all of their loot.
