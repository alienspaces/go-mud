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

    return BlocConsumer<CharacterCubit, CharacterState>(
      listener: (BuildContext context, CharacterState state) {
        //
      },
      builder: (BuildContext context, CharacterState characterState) {
        List<Widget> pageWidgets = [];

        if (characterState is CharacterStateCreate ||
            characterState is CharacterStateCreateError) {
          log.info('Displaying character create widget');

          //
          // Create character widget
          //
          pageWidgets.add(Container(
            color: Colors.brown,
            child: CharacterCreateWidget(
              callbacks: widget.callbacks,
            ),
          ));
        } else {
          log.info('Displaying character list widget');

          //
          // Character list widget
          //
          pageWidgets.add(Container(
            color: Colors.brown,
            child: CharacterListWidget(
              callbacks: widget.callbacks,
            ),
          ));

          log.info('Displaying character create button widget');

          //
          // Create character button
          //
          pageWidgets.add(Container(
            color: Colors.brown,
            child: CharacterCreateButtonWidget(
              callbacks: widget.callbacks,
            ),
          ));
        }

        return Container(
          color: Colors.brown,
          child: Column(
            children: pageWidgets,
          ),
        );
      },
    );
  }
}
