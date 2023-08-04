// Application packages
import 'package:go_mud_client/logger.dart';
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
  final log = getLogger('getLocationContents', null);

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
  log.fine('*** Dungeon objects ${locationData.locationObjects?.length}');

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
  log.fine('*** Dungeon characters ${locationData.locationCharacters?.length}');

  if (locationData.locationCharacters != null) {
    for (var character in locationData.locationCharacters!) {
      log.info(
        'Character ${character.name} '
        'health ${character.health} '
        'currentHealth ${character.currentHealth} '
        'fatigue ${character.fatigue} '
        'currentFatigue ${character.currentFatigue}',
      );

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
  log.fine('*** Dungeon monsters ${locationData.locationMonsters?.length}');

  if (locationData.locationMonsters != null) {
    for (var monster in locationData.locationMonsters!) {
      log.info(
        'Monster ${monster.name} '
        'health ${monster.health} '
        'currentHealth ${monster.currentHealth} '
        'fatigue ${monster.fatigue} '
        'currentFatigue ${monster.currentFatigue}',
      );

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

  // TODO: 12-implement-death: "Replace" here

  // Remove objects, characters or monsters already allocated and remove
  // object, characters and monsters that are no longer present
  // var currentLocationContent = populatedByIndex.values.toList();

  for (var idx = 0; idx < roomLocationCount; idx++) {
    var contentData = populatedByIndex[idx];
    if (contentData == null) {
      continue;
    }
    if (allocateContent[contentData.name] != null) {
      log.info("Replacing idx $idx content name ${contentData.name}");
      // Replace existing object, character or monster
      populatedByIndex[idx] = allocateContent[contentData.name]!;
      populatedByName[contentData.name] = idx;
      // Remove from list that needs allocating
      allocateContent.remove(contentData.name);
    } else {
      populatedByName.remove(contentData.name);
      populatedByIndex.remove(idx);
      unpopulated.remove(idx);
    }
  }

  // Shuffle unpopulated room locations
  unpopulated.shuffle();

  log.fine('*** Unpopulated before ${unpopulated.length}');

  // Allocate remaining objects, characters and monsters
  allocateContent.forEach((name, contentData) {
    var contentIdx = unpopulated.removeAt(0);
    populatedByName[name] = contentIdx;
    populatedByIndex[contentIdx] = contentData;
  });

  log.fine('*** Unpopulated after ${unpopulated.length}');

  locationPopulatedByIndex[locationData.locationName] = populatedByIndex;
  locationPopulatedByName[locationData.locationName] = populatedByName;
  locationUnpopulated[locationData.locationName] = unpopulated;

  return populatedByIndex;
}
