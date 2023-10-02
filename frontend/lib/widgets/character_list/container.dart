import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character_collection/character_collection_cubit.dart';
import 'package:go_mud_client/widgets/character_list/create_button.dart';
import 'package:go_mud_client/widgets/character_list/list.dart';

class CharacterListContainerWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterListContainerWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterListContainerWidget> createState() =>
      _CharacterListContainerWidgetState();
}

class _CharacterListContainerWidgetState
    extends State<CharacterListContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterListContainerWidget', 'build');

    return BlocConsumer<CharacterCollectionCubit, CharacterCollectionState>(
      listener: (BuildContext context, CharacterCollectionState state) {
        //
      },
      builder: (BuildContext context,
          CharacterCollectionState characterCollectionState) {
        List<Widget> pageWidgets = [];

        log.info('Displaying character list widget');

        //
        // Character list widget
        //
        pageWidgets.add(
          CharacterListWidget(
            callbacks: widget.callbacks,
          ),
        );

        log.info('Displaying character create button widget');

        //
        // Create character button
        //
        pageWidgets.add(
          CharacterCreateButtonWidget(
            callbacks: widget.callbacks,
          ),
        );

        return Column(
          children: pageWidgets,
        );
      },
    );
  }
}
