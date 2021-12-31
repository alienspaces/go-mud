import 'package:go_mud_client/widgets/common/character.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

const int maxAttributes = 36;

class CharacterPlayWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterPlayWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  _CharacterPlayWidgetState createState() => _CharacterPlayWidgetState();
}

class _CharacterPlayWidgetState extends State<CharacterPlayWidget> {
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
                const CharacterWidget(),
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
