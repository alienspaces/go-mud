import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

class GameDungeonGridWidget extends StatefulWidget {
  const GameDungeonGridWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonGridWidgetState createState() => _GameDungeonGridWidgetState();
}

typedef DungeonGridMemberFunction = Widget Function(DungeonActionRecord record, String key);

// Room content location maps
const int roomLocationCount = 14;
Map<String, Map<int, String>> locationPopulatedByIndex = {};
Map<String, Map<String, int>> locationPopulatedByName = {};
List<int> roomLocationsUnpopulated = [for (var i = 0; i < roomLocationCount; i++) i];

int characterIdx = 0;
int monsterIdx = 0;
int objectIdx = 0;

double gridMemberWidth = 50;
double gridMemberHeight = 50;

enum MemberType { character, monster, object }

class _GameDungeonGridWidgetState extends State<GameDungeonGridWidget> {
  Map<String, String> directionLabelMap = {
    'north': 'N',
    'northeast': 'NE',
    'east': 'E',
    'southeast': 'SE',
    'south': 'S',
    'southwest': 'SW',
    'west': 'W',
    'northwest': 'NW',
    'up': 'U',
    'down': 'D',
  };

  List<Widget> generateGrid(BuildContext context) {
    final log = getLogger('CharacterCreateWidget');

    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    var dungeonActionRecord = dungeonActionCubit.dungeonActionRecord;
    if (dungeonActionRecord == null) {
      log.warning('Dungeon cubit dungeon record is null, cannot create character');
      return [];
    }

    Map<int, String> populatedByIndex =
        locationPopulatedByIndex[dungeonActionRecord.location.name] ?? {};
    Map<String, int> populatedByName =
        locationPopulatedByName[dungeonActionRecord.location.name] ?? {};

    Map<String, MemberType> newLocationContents = {};
    List<String> roomContentNames = [];

    log.warning('*** Dungeon objects ${dungeonActionRecord.objects}');
    if (dungeonActionRecord.objects != null) {
      for (var dungeonObject in dungeonActionRecord.objects!) {
        newLocationContents[dungeonObject.name] = MemberType.object;
        roomContentNames.add(dungeonObject.name);
      }
    }
    log.warning('*** Dungeon characters ${dungeonActionRecord.characters}');
    if (dungeonActionRecord.characters != null) {
      for (var dungeonCharacter in dungeonActionRecord.characters!) {
        newLocationContents[dungeonCharacter.name] = MemberType.character;
        roomContentNames.add(dungeonCharacter.name);
      }
    }
    log.warning('*** Dungeon monsters ${dungeonActionRecord.monsters}');
    if (dungeonActionRecord.monsters != null) {
      for (var dungeonMonster in dungeonActionRecord.monsters!) {
        newLocationContents[dungeonMonster.name] = MemberType.monster;
        roomContentNames.add(dungeonMonster.name);
      }
    }

    // Remove objects, characters or monsters already allocated and remove
    // object, characters and monsters that are no longer present
    var currentContentNames = populatedByIndex.values.toList();
    for (var contentName in currentContentNames) {
      if (newLocationContents[contentName] != null) {
        // Location object, character or monster already allocated
        newLocationContents.remove(contentName);
      } else {
        // Location object, character or monster no longer present
        var contentIdx = populatedByName[contentName];
        if (contentIdx != null) {
          populatedByName.remove(contentName);
          populatedByIndex.remove(contentIdx);
          roomLocationsUnpopulated.remove(contentIdx);
        }
      }
    }

    // Shuffle unpopulated room locations
    roomLocationsUnpopulated.shuffle();

    // Allocate remaining objects, characters and monsters
    newLocationContents.forEach((name, type) {
      var contentIdx = roomLocationsUnpopulated.removeAt(0);
      populatedByName[name] = contentIdx;
      populatedByIndex[contentIdx] = name;
    });

    locationPopulatedByIndex[dungeonActionRecord.location.name] = populatedByIndex;
    locationPopulatedByName[dungeonActionRecord.location.name] = populatedByName;

    int roomGridIdx = 0;
    List<Widget Function()> dunegonGridMemberFunctions = [
      // Top Row
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'northwest'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'north'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'northeast'),
      // Second Row
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'up'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      // Third Row
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'west'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'east'),
      // Fourth Row
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'down'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      // Bottom Row
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'southwest'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'south'),
      () => roomWidget(context, populatedByIndex, roomGridIdx++),
      () => directionWidget(context, dungeonActionCubit.dungeonActionRecord!, 'southeast'),
    ];

    List<Widget> gridWidgets = [];
    for (var gridMemberFunction in dunegonGridMemberFunctions) {
      gridWidgets.add(gridMemberFunction());
    }

    return gridWidgets;
  }

  // Direction widget
  Widget directionWidget(BuildContext context, DungeonActionRecord record, String direction) {
    if (!record.location.directions.contains(direction)) {
      return emptyWidget('${directionLabelMap[direction]}');
    }

    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonGridWidget');

          final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
          if (dungeonCubit.dungeonRecord == null) {
            log.warning(
              'onPressed - Dungeon cubit missing dungeon record, cannot initialise action',
            );
            return;
          }
          // TODO: Check action has been selected prior to allowing direction to be selected for consistent support of look actions
          _submitAction(context, 'move $direction');
        },
        child: Text('${directionLabelMap[direction]}'),
      ),
    );
  }

  // Object widget
  Widget objectWidget(BuildContext context, String objectName, String direction) {
    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonGridWidget');

          final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
          if (dungeonCubit.dungeonRecord == null) {
            log.warning(
              'onPressed - Dungeon cubit missing dungeon record, cannot initialise action',
            );
            return;
          }
          // TODO: Check action has been selected prior to allowing direction to be selected for consistent support of look actions
          _submitAction(context, 'look $objectName');
        },
        child: Text(objectName),
      ),
    );
  }

  // Room widget
  Widget roomWidget(BuildContext context, Map<int, String> populatedByIndex, int idx) {
    if (populatedByIndex[idx] != null) {
      return emptyWidget(populatedByIndex[idx]!);
    }

    return emptyWidget('E$idx');
  }

  // Empty widget
  Widget emptyWidget(String label) {
    return Container(
      width: gridMemberWidth,
      height: gridMemberHeight,
      alignment: Alignment.center,
      margin: const EdgeInsets.all(2),
      decoration: BoxDecoration(
        color: const Color(0xFFD4D4D4),
        border: Border.all(
          color: const Color(0xFFD4D4D4),
        ),
        borderRadius: const BorderRadius.all(Radius.circular(5)),
      ),
      child: Text(label),
    );
  }

  void _submitAction(BuildContext context, String action) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Submitting move action..');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    if (dungeonCubit.dungeonRecord == null) {
      log.warning('Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning('Character cubit missing character record, cannot initialise action');
      return;
    }

    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    dungeonActionCubit.createAction(
      dungeonCubit.dungeonRecord!.id,
      characterCubit.characterRecord!.id,
      action,
    );
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStateCreated) {
          return Container(
            decoration: BoxDecoration(
              color: const Color(0xFFDEDEDE),
              border: Border.all(
                color: const Color(0xFFDEDEDE),
              ),
              borderRadius: const BorderRadius.all(Radius.circular(5)),
            ),
            padding: const EdgeInsets.all(1),
            margin: const EdgeInsets.all(5),
            width: gridMemberWidth * 5,
            height: gridMemberHeight * 5,
            child: GridView.count(
              crossAxisCount: 5,
              children: generateGrid(context),
            ),
          );
        }

        // Empty
        return Container();
      },
    );
  }
}
