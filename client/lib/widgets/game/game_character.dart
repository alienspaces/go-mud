import 'package:go_mud_client/widgets/common/character.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

class GameCharacterWidget extends StatefulWidget {
  const GameCharacterWidget({Key? key}) : super(key: key);

  @override
  _GameCharacterWidgetState createState() => _GameCharacterWidgetState();
}

class _GameCharacterWidgetState extends State<GameCharacterWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameCharacterWidget');
    log.info('Building..');

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (BuildContext context, CharacterState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, CharacterState state) {
        if (state is CharacterStateSelected) {
          // ignore: avoid_unnecessary_containers
          return Container(
            child: Column(
              children: const <Widget>[
                CharacterWidget(),
              ],
            ),
          );
        }

        // Empty
        return Container();
      },
    );
  }
}
