# Nest M.U.D - How To Play

- [Getting Started](README.md)
- [How to Play](README-HOWTOPLAY.md)
- [API](README-API.md)
- [Design](README-DESIGN.md)
- [TODO](README-TODO.md)

## Characters

Characters are created and begin their adventure at the start of a dungeon.

## Actions

Characters are controlled with actions that are simple sentences.

### Movement Actions

Character can `move` from one location to another location using the `move [direction]` action.

_Example_:

```text
move north
move southwest
move up
```

### Look Actions

A character can `look` to get a description of their current location, look in a specific direction, at an object, at another character or at a monster that is at the current location.

When a character looks at an object that is not in their possession a brief description of the object will be provided. To look more closely at an object, a character must have the object equipped or stashed in their backpack.

When a character looks a direction a brief description of the location will be provided. To get a more detailed view of a location a character must move to that location.

_Example_:

```text
look
look north
look chest
look goblin
look down
```

### Equip, Stash and Drop Actions

A character may `equip` an object that is stashed in their backpack or at the current location. Any item that is currently equipped where the newly equipped item would normally be equipped, will be automatically stashed if there is enough room in the character's backpack. If there is not enough room in a character's backback to stash an already equipped object, equipping the new object will fail.

A character may `stash` an object that is currently equipped or at the character's current location into their backpack. If there is not enough room in a characters backback to stash an object, stashing the object will fail.

A character may `drop` any item that is currently equipped or stashed in their backpack. If there is not enough room at the character's current location to drop the object, dropping the object will fail.

_Example_:

```text
equip sword
stash gold hammer
drop dragon tongue
```

### Attack and Defend Actions

📝 [Issue-2](https://gitlab.com/alienspaces/go-mud/-/issues/2)

A character equipped with a melee weapon can attack any target if the target is in the same room. A character equipped with two weapons, one in each hand, can specify which weapon to use.

_Example_:

```text
attack Grumpy Dwarf
attack Grumpy Dwarf with Rusted Dagger
```

A character equipped with a ranged weapon attacks any target if the target is in the same room or any adjacent room.

The character equipped with a melee weapon or shield spends the turn defending themselves which decreases the chance of getting hit by an opponent. A character equipped with a shield should be more effective in defending that an character who is not. A character equipped with a ranged weapon cannot defend at all.

_Example_:

```text
defend
```

### Use Actions

📝 [Issue-3](https://gitlab.com/alienspaces/go-mud/-/issues/3)

A character that has the item stashed or equipped uses the item. The result of using an item should be detailed in the action response. 

_Example_:

```text
use Potion of Healing
```

_Example Response_:

`You use the Potion of Healing, replenishing 10 health.`

Items should be able to be used multiple times if they support it.

_Example Response_:

- A potion may have 5 uses
- A wand may be able to be used forever but needs 5 rounds to recharge
- A door may toggle between open and closed

The result of not being able to use an item should be detailed in the action response.

`The Potion of Healing is empty.`

### Loot Actions

📝 [Issue-4](https://gitlab.com/alienspaces/go-mud/-/issues/4)

A character or monster may loot any other monster or character that has reached it's demise at the current location. All looted items will be automatically stashed. Objects that could not be looted due to the character or monster running out of room in their backpack will remain on the corpse. Once a monster or character has no more objects to loot, they disappear.

_Example_:

```text
loot Angry Kobold
```

### Speech Actions

📝 [Issue-5](https://gitlab.com/alienspaces/go-mud/-/issues/5)

Characters and monsters should be able to say things to each other in the current room. Saying something without specifying another character or monster can be heard by everyone in the room. Saying something to a character or monster can be only heard by the character or monster being spoken to.

_Example_:

```text
say hello
say Legislate look north
say Dirty Goblin King hello
```

The user interface should provide a `Say` button which then presents a set of words that can be used to form a sentence along with a button to say the sentence, a button to delete the most recent word and a button to cancel saying anything at all.

Pressing the say button once will then wait for the user to click a monster or character before displaying a set of words. Pressing the say button twice immediately displays a set of words. Pressing the say button three times or pressing any other action button will reset the say button back to its unpressed state.
