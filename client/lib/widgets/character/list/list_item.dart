import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/repository/character/character_repository.dart';

class CharacterListItemWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final CharacterRecord characterRecord;
  const CharacterListItemWidget(
      {Key? key, required this.callbacks, required this.characterRecord})
      : super(key: key);

  /// Sets the current character state to the provided character
  void _selectCharacter(BuildContext context, CharacterRecord characterRecord) {
    final log = getLogger('CharacterListItemWidget');
    log.fine('Select ${characterRecord.id} ${characterRecord.name}');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    characterCubit.selectCharacter(characterRecord);

    callbacks.openDungeonPage(context);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterListItemWidget');
    log.fine('Display ${characterRecord.id} ${characterRecord.name}');

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(30, 15, 30, 15),
      textStyle: Theme.of(context).textTheme.button!.copyWith(fontSize: 18),
    );

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (BuildContext context, CharacterState state) {
        //
      },
      builder: (BuildContext context, CharacterState state) {
        // ignore: avoid_unnecessary_containers
        return Container(
          margin: const EdgeInsets.all(5),
          decoration: BoxDecoration(
            border: Border.all(width: 2),
          ),
          child: Column(
            children: [
              Container(
                margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
                child: Text(
                  characterRecord.name,
                  style: Theme.of(context).textTheme.headline3,
                ),
              ),
              Container(
                margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: <Widget>[
                    Container(
                      margin: const EdgeInsets.all(5),
                      child: ElevatedButton(
                        onPressed: () => {null},
                        style: buttonStyle,
                        child: const Text('Delete'),
                      ),
                    ),
                    Container(
                      margin: const EdgeInsets.all(5),
                      child: ElevatedButton(
                        onPressed: () =>
                            _selectCharacter(context, characterRecord),
                        style: buttonStyle,
                        child: const Text('Play'),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
