import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

class CharacterWidget extends StatelessWidget {
  const CharacterWidget({Key? key}) : super(key: key);
  Widget _buildAttribute(String name, int? value) {
    EdgeInsetsGeometry padding = const EdgeInsets.fromLTRB(10, 10, 10, 10);

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: <Widget>[
        const Spacer(flex: 1),
        Flexible(
          flex: 2,
          child: Container(
            padding: padding,
            alignment: Alignment.centerLeft,
            child: Text(name),
          ),
        ),
        const Spacer(flex: 1),
        Flexible(
          flex: 2,
          child: Container(
            padding: padding,
            alignment: Alignment.centerLeft,
            child: Text('$value'),
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterPlayWidget');
    log.info('Building..');

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (BuildContext context, CharacterState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, CharacterState state) {
        if (state is CharacterStateSelected) {
          return Container(
            margin: const EdgeInsets.fromLTRB(20, 2, 20, 2),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: <Widget>[
                // ignore: avoid_unnecessary_containers
                Container(
                  child: Text(
                    state.characterRecord.name,
                    style: Theme.of(context).textTheme.headline5,
                  ),
                ),
                Container(
                  child: _buildAttribute('Strength', state.characterRecord.strength),
                ),
                Container(
                  child: _buildAttribute('Dexterity', state.characterRecord.dexterity),
                ),
                Container(
                  child: _buildAttribute('Intelligence', state.characterRecord.intelligence),
                ),
                Container(
                  child: _buildAttribute('Health', state.characterRecord.health),
                ),
                Container(
                  child: _buildAttribute('Fatigue', state.characterRecord.fatigue),
                ),
              ],
            ),
          );
        }

        // Shouldn't get here..
        return Container();
      },
    );
  }
}
