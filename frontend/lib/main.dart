import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/config.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/theme.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/cubit/character_create/character_create_cubit.dart';
import 'package:go_mud_client/cubit/character_collection/character_collection_cubit.dart';

void main() {
  // Initialise logger
  initLogger();

  // When hostname is localhost and we are running in an emulator set backend to specific IP
  if (!kIsWeb && config['serverHost'].toString() == 'localhost') {
    config['serverHost'] = '10.0.3.2';
  }

  // Server API calls
  final API api = API(config: config);

  // Server resource request methods and record management
  final RepositoryCollection repositories =
      RepositoryCollection(config: config, api: api);

  // Global game state observer that *may* provide the ability to call on
  // cubit methods when the state of others cubits change?
  Bloc.observer = GameStateObserver();

  // Run application
  runApp(MainApp(config: config, repositories: repositories));
}

class MainApp extends StatelessWidget {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  MainApp({Key? key, required this.config, required this.repositories})
      : super(key: key);

  final log = getLogger('MainApp', null);

  @override
  Widget build(BuildContext context) {
    log.fine('Building..');
    return MaterialApp(
      title: 'Go MUD Client',
      theme: getTheme(context),
      home: MultiBlocProvider(
        providers: [
          BlocProvider<DungeonCubit>(
            create: (BuildContext context) => DungeonCubit(
              config: config,
              repositories: repositories,
            ),
          ),
          BlocProvider<DungeonActionCubit>(
            create: (BuildContext context) => DungeonActionCubit(
              config: config,
              repositories: repositories,
            ),
          ),
          BlocProvider<DungeonCommandCubit>(
            create: (BuildContext context) => DungeonCommandCubit(
              config: config,
              repositories: repositories,
            ),
          ),
          BlocProvider<CharacterCreateCubit>(
            create: (BuildContext context) => CharacterCreateCubit(
              config: config,
              repositories: repositories,
            ),
          ),
          BlocProvider<CharacterCollectionCubit>(
            create: (BuildContext context) => CharacterCollectionCubit(
              config: config,
              repositories: repositories,
              dungeonActionCubit: BlocProvider.of<DungeonActionCubit>(context),
              characterCreateCubit:
                  BlocProvider.of<CharacterCreateCubit>(context),
            ),
          ),
          BlocProvider<CharacterCubit>(
            create: (BuildContext context) => CharacterCubit(
              config: config,
              repositories: repositories,
              dungeonActionCubit: BlocProvider.of<DungeonActionCubit>(context),
            ),
          ),
        ],
        child: const Navigation(),
      ),
    );
  }
}

// Interesting functionality but not clear how this is going to help at the moment.
class GameStateObserver extends BlocObserver {
  @override
  void onCreate(BlocBase bloc) {
    super.onCreate(bloc);
    final log = getLogger('GameStateObserver', 'onCreate');
    log.fine('onCreate -- ${bloc.runtimeType}');
  }

  @override
  void onChange(BlocBase bloc, Change change) {
    super.onChange(bloc, change);
    final log = getLogger('GameStateObserver', 'onChange');
    log.fine('onChange -- ${bloc.runtimeType}, $change');
  }

  @override
  void onError(BlocBase bloc, Object error, StackTrace stackTrace) {
    final log = getLogger('GameStateObserver', 'onError');
    log.fine('onError -- ${bloc.runtimeType}, $error');
    super.onError(bloc, error, stackTrace);
  }

  @override
  void onClose(BlocBase bloc) {
    super.onClose(bloc);
    final log = getLogger('GameStateObserver', 'onClose');
    log.fine('onClose -- ${bloc.runtimeType}');
  }
}
