import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

const int maxAttributes = 36;

class CharacterTrainWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterTrainWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  _CharacterTrainWidgetState createState() => _CharacterTrainWidgetState();
}

class _CharacterTrainWidgetState extends State<CharacterTrainWidget> {
  Widget _buildAttribute(String name, int? value) {
    EdgeInsetsGeometry padding = const EdgeInsets.fromLTRB(10, 2, 10, 2);

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
    final log = getLogger('CharacterTrainWidget');
    log.fine('Building..');

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (BuildContext context, CharacterState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, CharacterState state) {
        if (state is CharacterStateSelected) {
          return Container(
            margin: const EdgeInsets.fromLTRB(20, 2, 20, 2),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: <Widget>[
                Column(
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
                      child: _buildAttribute(
                          'Strength', state.characterRecord.strength),
                    ),
                    Container(
                      child: _buildAttribute(
                          'Dexterity', state.characterRecord.dexterity),
                    ),
                    Container(
                      child: _buildAttribute(
                          'Intelligence', state.characterRecord.intelligence),
                    ),
                    Container(
                      child: _buildAttribute(
                          'Health', state.characterRecord.health),
                    ),
                    Container(
                      child: _buildAttribute(
                          'Fatigue', state.characterRecord.fatigue),
                    ),
                  ],
                ),
                // ignore: avoid_unnecessary_containers
                Container(
                  child: ElevatedButton(
                    onPressed: () {
                      widget.callbacks.openGamePage(context);
                    },
                    child: const Text(
                      'Play',
                    ),
                  ),
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
