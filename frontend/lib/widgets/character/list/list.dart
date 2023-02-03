import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/widgets/character/list/list_item.dart';

class CharacterListWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterListWidget({Key? key, required this.callbacks})
      : super(key: key);

  @override
  State<CharacterListWidget> createState() => _CharacterListWidgetState();
}

class _CharacterListWidgetState extends State<CharacterListWidget> {
  @override
  void initState() {
    final log = getLogger('CharacterListWidget', 'initState');
    log.fine('Initialising state..');

    super.initState();

    _loadCharacters(context);
  }

  void _loadCharacters(BuildContext context) {
    final log = getLogger('CharacterListWidget', '_loadCharacters');
    log.info('Loading characters');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    characterCubit.loadCharacters();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterListWidget', 'build');
    log.fine('Building..');

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (context, state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, CharacterState state) {
        log.fine('builder...');
        List<Widget> widgets = [];

        if (state is CharacterStateLoaded) {
          // Character list
          state.characterRecords?.forEach((characterRecord) {
            log.info([
              'Displaying character list item widget',
              characterRecord.characterName,
            ]);
            widgets.add(
              // ignore: avoid_unnecessary_containers
              Container(
                child: CharacterListItemWidget(
                  callbacks: widget.callbacks,
                  characterRecord: characterRecord,
                ),
              ),
            );
          });
        } else if (state is CharacterStateLoadError) {
          widgets.add(
            // ignore: avoid_unnecessary_containers
            Container(
              child: ElevatedButton(
                onPressed: () => _loadCharacters(context),
                child: const Text('Load Characters'),
              ),
            ),
          );
        } else {
          widgets.add(
            // ignore: avoid_unnecessary_containers
            Container(
              child: const Text('Loading characters...'),
            ),
          );
        }

        // ignore: avoid_unnecessary_containers
        return Container(
          child: Column(
            children: widgets,
          ),
        );
      },
    );
  }
}
