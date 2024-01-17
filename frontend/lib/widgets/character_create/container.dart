import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character_create/character_create_cubit.dart';
import 'package:go_mud_client/widgets/character_create/create.dart';

class CharacterCreateContainerWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterCreateContainerWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterCreateContainerWidget> createState() =>
      _CharacterCreateContainerWidgetState();
}

class _CharacterCreateContainerWidgetState
    extends State<CharacterCreateContainerWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterCreateContainerWidget', 'build');

    return BlocConsumer<CharacterCreateCubit, CharacterCreateState>(
      listener: (BuildContext context, CharacterCreateState state) {
        //
      },
      builder:
          (BuildContext context, CharacterCreateState characterCreateState) {
        List<Widget> pageWidgets = [];

        log.info('Displaying character create widget');

        //
        // Create character widget
        //
        pageWidgets.add(
          CharacterCreateWidget(
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
