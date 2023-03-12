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
    final log =
        getLogger('CharacterCreateButtonWidget', '_initCreateCharacter');
    log.info('Initiating character creation');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    characterCubit.initCreateCharacter();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterCreateButtonWidget', 'build');
    log.info('Building..');

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(30, 15, 30, 15),
      textStyle: Theme.of(context).textTheme.labelLarge!.copyWith(fontSize: 18),
    );

    final characterCubit = BlocProvider.of<CharacterCubit>(context);

    if (characterCubit.canCreateCharacter()) {
      // ignore: avoid_unnecessary_containers
      return Container(
        margin: const EdgeInsets.all(5),
        // decoration: BoxDecoration(
        //   border: Border.all(width: 2),
        // ),
        alignment: Alignment.centerRight,
        child: ElevatedButton(
          onPressed: () => _initCreateCharacter(context),
          style: buttonStyle,
          child: const Text('Create Character'),
        ),
      );
    }

    return Container();
  }
}
