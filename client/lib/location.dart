// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';

// Room content location maps
const int roomLocationCount = 14;
enum ContentType { character, monster, object }

class LocationContent {
  int index;
  String name;
  ContentType type;
  LocationContent({required this.index, required this.name, required this.type});
}

Map<String, Map<int, LocationContent>> locationPopulatedByIndex = {};
Map<String, Map<String, int>> locationPopulatedByName = {};
Map<String, List<int>> locationUnpopulated = {};

Map<int, LocationContent> getLocationContents(DungeonActionRecord dungeonActionRecord) {
  final log = getLogger('GetLocationContents');

  var location = dungeonActionRecord.location;

  Map<int, LocationContent> populatedByIndex = locationPopulatedByIndex[location.name] ?? {};
  Map<String, int> populatedByName = locationPopulatedByName[location.name] ?? {};
  List<int> unpopulated =
      locationUnpopulated[location.name] ?? [for (var i = 0; i < roomLocationCount; i++) i];

  Map<String, ContentType> newLocationContents = {};
  List<String> roomContentNames = [];

  log.warning('*** Dungeon objects ${location.objects?.length}');
  if (location.objects != null) {
    for (var dungeonObject in location.objects!) {
      newLocationContents[dungeonObject.name] = ContentType.object;
      roomContentNames.add(dungeonObject.name);
    }
  }
  log.warning('*** Dungeon characters ${location.characters?.length}');
  if (location.characters != null) {
    for (var dungeonCharacter in location.characters!) {
      newLocationContents[dungeonCharacter.name] = ContentType.character;
      roomContentNames.add(dungeonCharacter.name);
    }
  }
  log.warning('*** Dungeon monsters ${location.monsters?.length}');
  if (location.monsters != null) {
    for (var dungeonMonster in location.monsters!) {
      newLocationContents[dungeonMonster.name] = ContentType.monster;
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

  log.warning('*** Unpopulated before ${unpopulated.length}');

  // Allocate remaining objects, characters and monsters
  newLocationContents.forEach((name, type) {
    var contentIdx = unpopulated.removeAt(0);
    populatedByName[name] = contentIdx;
    populatedByIndex[contentIdx] = LocationContent(index: contentIdx, name: name, type: type);
  });

  log.warning('*** Unpopulated after ${unpopulated.length}');

  locationPopulatedByIndex[dungeonActionRecord.location.name] = populatedByIndex;
  locationPopulatedByName[dungeonActionRecord.location.name] = populatedByName;
  locationUnpopulated[dungeonActionRecord.location.name] = unpopulated;

  return populatedByIndex;
}
