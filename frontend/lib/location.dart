// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

// Room content location maps
const int roomLocationCount = 14;

enum ContentType { character, monster, object }

// class LocationContent {
//   // final int index;
//   final ContentType type;
//   final String name;
//   final int? health;
//   final int? fatigue;
//   LocationContent({
//     // required this.index,
//     required this.type,
//     required this.name,
//     this.health,
//     this.fatigue,
//   });
// }

class LocationContent {
  final ContentType type;
  final String name;
  final int? health;
  final int? fatigue;
  LocationContent({
    required this.type,
    required this.name,
    this.health,
    this.fatigue,
  });
}

Map<String, Map<int, LocationContent>> locationPopulatedByIndex = {};
Map<String, Map<String, int>> locationPopulatedByName = {};
Map<String, List<int>> locationUnpopulated = {};

Map<int, LocationContent> getLocationContents(LocationData locationData) {
  final log = getLogger('getLocationContents', null);

  Map<int, LocationContent> populatedByIndex =
      locationPopulatedByIndex[locationData.locationName] ?? {};
  Map<String, int> populatedByName =
      locationPopulatedByName[locationData.locationName] ?? {};
  List<int> unpopulated = locationUnpopulated[locationData.locationName] ??
      [for (var i = 0; i < roomLocationCount; i++) i];

  Map<String, LocationContent> newLocationContents = {};
  List<String> roomContentNames = [];

  log.fine('*** Dungeon objects ${locationData.locationObjects?.length}');
  if (locationData.locationObjects != null) {
    for (var dungeonObject in locationData.locationObjects!) {
      newLocationContents[dungeonObject.name] = LocationContent(
        type: ContentType.object,
        name: dungeonObject.name,
      );
      roomContentNames.add(dungeonObject.name);
    }
  }
  log.fine('*** Dungeon characters ${locationData.locationCharacters?.length}');
  if (locationData.locationCharacters != null) {
    for (var dungeonCharacter in locationData.locationCharacters!) {
      newLocationContents[dungeonCharacter.name] = LocationContent(
        type: ContentType.character,
        name: dungeonCharacter.name,
        health: dungeonCharacter.health,
        fatigue: dungeonCharacter.fatigue,
      );
      roomContentNames.add(dungeonCharacter.name);
    }
  }
  log.fine('*** Dungeon monsters ${locationData.locationMonsters?.length}');
  if (locationData.locationMonsters != null) {
    for (var dungeonMonster in locationData.locationMonsters!) {
      newLocationContents[dungeonMonster.name] = LocationContent(
        type: ContentType.monster,
        name: dungeonMonster.name,
        health: dungeonMonster.health,
        fatigue: dungeonMonster.fatigue,
      );
      roomContentNames.add(dungeonMonster.name);
    }
  }

  // Remove objects, characters or monsters already allocated and remove
  // object, characters and monsters that are no longer present
  var currentLocationContent = populatedByIndex.values.toList();
  for (var locationContent in currentLocationContent) {
    if (newLocationContents[locationContent.name] != null) {
      // Location object, character or monster already allocated
      newLocationContents.remove(locationContent.name);
    } else {
      // Location object, character or monster no longer present
      var contentIdx = populatedByName[locationContent.name];
      if (contentIdx != null) {
        populatedByName.remove(locationContent.name);
        populatedByIndex.remove(contentIdx);
        unpopulated.remove(contentIdx);
      }
    }
  }

  // Shuffle unpopulated room locations
  unpopulated.shuffle();

  log.fine('*** Unpopulated before ${unpopulated.length}');

  // Allocate remaining objects, characters and monsters
  newLocationContents.forEach((name, contentData) {
    var contentIdx = unpopulated.removeAt(0);
    populatedByName[name] = contentIdx;
    populatedByIndex[contentIdx] = LocationContent(
      // index: contentIdx,
      type: contentData.type,
      name: name,
      health: contentData.health,
      fatigue: contentData.fatigue,
    );
  });

  log.fine('*** Unpopulated after ${unpopulated.length}');

  locationPopulatedByIndex[locationData.locationName] = populatedByIndex;
  locationPopulatedByName[locationData.locationName] = populatedByName;
  locationUnpopulated[locationData.locationName] = unpopulated;

  return populatedByIndex;
}
