import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

class CharacterCreateButtonWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;

  const CharacterCreateButtonWidget({Key? key, required this.callbacks})
      : super(key: key);

  void _initCreateCharacter(BuildContext context) {
    final log = getLogger('CharacterCreateButtonWidget');
    log.fine('Initiating character creation');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    characterCubit.initCreateCharacter();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterCreateButtonWidget');
    log.fine('Building..');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);

    if (characterCubit.canCreateCharacter()) {
      // ignore: avoid_unnecessary_containers
      return Container(
        child: ElevatedButton(
          child: const Text('Create Character'),
          onPressed: () => _initCreateCharacter(context),
        ),
      );
    }

    return Container();
  }
}
