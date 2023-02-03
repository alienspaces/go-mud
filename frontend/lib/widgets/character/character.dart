import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/widgets/character/create/create.dart';
import 'package:go_mud_client/widgets/character/create/create_button.dart';
import 'package:go_mud_client/widgets/character/list/list.dart';

class CharacterContainerWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterContainerWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterContainerWidget> createState() =>
      _CharacterContainerWidgetState();
}

class _CharacterContainerWidgetState extends State<CharacterContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterContainerWidget', 'build');
    log.fine('Building..');

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (BuildContext context, CharacterState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, CharacterState characterState) {
        List<Widget> pageWidgets = [];

        if (characterState is CharacterStateCreate ||
            characterState is CharacterStateCreateError) {
          log.info('Displaying character create');
          // ignore: avoid_unnecessary_containers
          pageWidgets.add(Container(
            child: CharacterCreateWidget(
              callbacks: widget.callbacks,
            ),
          ));
        } else {
          log.info('Displaying character list');
          // ignore: avoid_unnecessary_containers
          pageWidgets.add(Container(
            child: CharacterListWidget(
              callbacks: widget.callbacks,
            ),
          ));

          // ignore: avoid_unnecessary_containers
          pageWidgets.add(Container(
            child: CharacterCreateButtonWidget(
              callbacks: widget.callbacks,
            ),
          ));
        }

        return Column(
          children: pageWidgets,
        );
      },
    );
  }
}
