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

Map<int, LocationContent> getLocationContents(LocationData locationData) {
  final log = getLogger('GetLocationContents');

  Map<int, LocationContent> populatedByIndex = locationPopulatedByIndex[locationData.name] ?? {};
  Map<String, int> populatedByName = locationPopulatedByName[locationData.name] ?? {};
  List<int> unpopulated =
      locationUnpopulated[locationData.name] ?? [for (var i = 0; i < roomLocationCount; i++) i];

  Map<String, ContentType> newLocationContents = {};
  List<String> roomContentNames = [];

  log.warning('*** Dungeon objects ${locationData.objects?.length}');
  if (locationData.objects != null) {
    for (var dungeonObject in locationData.objects!) {
      newLocationContents[dungeonObject.name] = ContentType.object;
      roomContentNames.add(dungeonObject.name);
    }
  }
  log.warning('*** Dungeon characters ${locationData.characters?.length}');
  if (locationData.characters != null) {
    for (var dungeonCharacter in locationData.characters!) {
      newLocationContents[dungeonCharacter.name] = ContentType.character;
      roomContentNames.add(dungeonCharacter.name);
    }
  }
  log.warning('*** Dungeon monsters ${locationData.monsters?.length}');
  if (locationData.monsters != null) {
    for (var dungeonMonster in locationData.monsters!) {
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

  locationPopulatedByIndex[locationData.name] = populatedByIndex;
  locationPopulatedByName[locationData.name] = populatedByName;
  locationUnpopulated[locationData.name] = unpopulated;

  return populatedByIndex;
}
