// Application packages
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

// Room content location maps
const int roomLocationCount = 14;

enum ContentType { character, monster, object }

class LocationContent {
  final ContentType type;
  final String name;
  final int? health;
  final int? currentHealth;
  final int? fatigue;
  final int? currentFatigue;
  LocationContent({
    required this.type,
    required this.name,
    this.health,
    this.currentHealth,
    this.fatigue,
    this.currentFatigue,
  });
}

Map<String, Map<int, LocationContent>> locationPopulatedByIndex = {};
Map<String, Map<String, int>> locationPopulatedByName = {};
Map<String, List<int>> locationUnpopulated = {};

// getLocationContents will populate things that are already at a
// location in the same location index for each turn so they do not
// appear to jump around in the location.
Map<int, LocationContent> getLocationContents(LocationData locationData) {
  // Existing allocation content indexed by their current location index
  Map<int, LocationContent> populatedByIndex =
      locationPopulatedByIndex[locationData.locationName] ?? {};

  // Location indexes indexed by the contents "unique" name
  Map<String, int> populatedByName =
      locationPopulatedByName[locationData.locationName] ?? {};

  // List of unpopulated indexes
  List<int> unpopulated = locationUnpopulated[locationData.locationName] ??
      [for (var i = 0; i < roomLocationCount; i++) i];

  // Content to be allocated indexed by the contents "unique" name
  Map<String, LocationContent> allocateContent = {};

  List<String> roomContentNames = [];

  // Add location objects to content to be allocated
  if (locationData.locationObjects != null) {
    for (var dungeonObject in locationData.locationObjects!) {
      allocateContent[dungeonObject.name] = LocationContent(
        type: ContentType.object,
        name: dungeonObject.name,
      );
      roomContentNames.add(dungeonObject.name);
    }
  }

  // Add location characters to content to be allocated
  if (locationData.locationCharacters != null) {
    for (var character in locationData.locationCharacters!) {
      allocateContent[character.name] = LocationContent(
        type: ContentType.character,
        name: character.name,
        health: character.health,
        currentHealth: character.currentHealth,
        fatigue: character.fatigue,
        currentFatigue: character.currentFatigue,
      );
      roomContentNames.add(character.name);
    }
  }

  // Add location monsters to content to be allocated
  if (locationData.locationMonsters != null) {
    for (var monster in locationData.locationMonsters!) {
      allocateContent[monster.name] = LocationContent(
        type: ContentType.monster,
        name: monster.name,
        health: monster.health,
        currentHealth: monster.currentHealth,
        fatigue: monster.fatigue,
        currentFatigue: monster.currentFatigue,
      );
      roomContentNames.add(monster.name);
    }
  }

  // Replace objects, characters and monsters that are already allocated and
  // remove object, characters and monsters that are no longer present
  for (var idx = 0; idx < roomLocationCount; idx++) {
    var contentData = populatedByIndex[idx];
    if (contentData == null) {
      continue;
    }
    if (allocateContent[contentData.name] != null) {
      // Replace in populated indexes
      populatedByIndex[idx] = allocateContent[contentData.name]!;
      populatedByName[contentData.name] = idx;
      // Remove from list that needs allocating
      allocateContent.remove(contentData.name);
    } else {
      // Remove from populated indexes
      populatedByName.remove(contentData.name);
      populatedByIndex.remove(idx);
      unpopulated.remove(idx);
    }
  }

  // Shuffle unpopulated room locations
  unpopulated.shuffle();

  // Allocate remaining objects, characters and monsters
  allocateContent.forEach((name, contentData) {
    var contentIdx = unpopulated.removeAt(0);
    populatedByName[name] = contentIdx;
    populatedByIndex[contentIdx] = contentData;
  });

  locationPopulatedByIndex[locationData.locationName] = populatedByIndex;
  locationPopulatedByName[locationData.locationName] = populatedByName;
  locationUnpopulated[locationData.locationName] = unpopulated;

  return populatedByIndex;
}
