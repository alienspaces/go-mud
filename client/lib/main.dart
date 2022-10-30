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
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

void main() {
  // Initialise logger
  initLogger();

  // When hostname is localhost and we are running in an emulator set backend to specific IP
  if (!kIsWeb && config['serverHost'].toString() == 'localhost') {
    config['serverHost'] = '10.0.3.2';
  }

  // Dependencies
  final API api = API(config: config);

  final repositories = RepositoryCollection(config: config, api: api);

  // Run application
  runApp(MainApp(config: config, repositories: repositories));
}

class MainApp extends StatelessWidget {
  final Map<String, String> config;
  final RepositoryCollection repositories;

  MainApp({Key? key, required this.config, required this.repositories})
      : super(key: key);

  final log = getLogger('MainApp');

  @override
  Widget build(BuildContext context) {
    log.fine('Building..');
    return MaterialApp(
      title: 'Go MUD Client',
      theme: getTheme(context),
      home: MultiBlocProvider(
        providers: [
          BlocProvider<DungeonCubit>(
            create: (BuildContext context) =>
                DungeonCubit(config: config, repositories: repositories),
          ),
          BlocProvider<CharacterCubit>(
            create: (BuildContext context) =>
                CharacterCubit(config: config, repositories: repositories),
          ),
          BlocProvider<DungeonCharacterCubit>(
            create: (BuildContext context) => DungeonCharacterCubit(
                config: config, repositories: repositories),
          ),
          BlocProvider<DungeonActionCubit>(
            create: (BuildContext context) =>
                DungeonActionCubit(config: config, repositories: repositories),
          ),
          BlocProvider<DungeonCommandCubit>(
            create: (BuildContext context) =>
                DungeonCommandCubit(config: config, repositories: repositories),
          ),
        ],
        child: const Navigation(),
      ),
    );
  }
}
