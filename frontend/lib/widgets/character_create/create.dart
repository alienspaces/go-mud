import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/cubit/character_create/character_create_cubit.dart';

const int maxAttributes = 36;

class CharacterCreateWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const CharacterCreateWidget({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterCreateWidget> createState() => _CharacterCreateWidgetState();
}

class _CharacterCreateWidgetState extends State<CharacterCreateWidget> {
  int strength = 8;
  int dexterity = 8;
  int intelligence = 8;

  // Global key that uniquely identifies the Form widget
  final _formKey = GlobalKey<FormState>();

  // Form field controllers
  final characterNameController = TextEditingController();

  @override
  void dispose() {
    // Clean up the controller when the widget is disposed.
    characterNameController.dispose();
    super.dispose();
  }

  void _createCharacter(BuildContext context) {
    final log = getLogger('CharacterCreateWidget', '_createCharacter');
    log.fine('Creating character name >${characterNameController.text}<');
    log.fine('Creating character strength >$strength<');
    log.fine('Creating character dexterity >$dexterity<');
    log.fine('Creating character intelligence >$intelligence<');

    final characterCubit = BlocProvider.of<CharacterCreateCubit>(context);

    CreateCharacterRecord createCharacterRecord = CreateCharacterRecord(
      characterName: characterNameController.text,
      characterStrength: strength,
      characterDexterity: dexterity,
      characterIntelligence: intelligence,
    );

    characterCubit.createCharacter(createCharacterRecord).then((result) {
      if (result.characterRecord != null && result.exception == null) {
        log.info(">>> Closing character create window");
        widget.callbacks.closeCharacterCreatePage(context);
      }
    });
  }

  void _incrementStrength() {
    if (strength + dexterity + intelligence < maxAttributes) {
      setState(() {
        strength++;
      });
    }
  }

  void _decrementStrength() {
    if (strength > 8) {
      setState(() {
        strength--;
      });
    }
  }

  void _incrementDexterity() {
    if (strength + dexterity + intelligence < maxAttributes) {
      setState(() {
        dexterity++;
      });
    }
  }

  void _decrementDexterity() {
    if (dexterity > 8) {
      setState(() {
        dexterity--;
      });
    }
  }

  void _incrementIntelligence() {
    if (strength + dexterity + intelligence < maxAttributes) {
      setState(() {
        intelligence++;
      });
    }
  }

  void _decrementIntelligence() {
    if (intelligence > 8) {
      setState(() {
        intelligence--;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterCreateWidget', 'build');
    log.fine('Building..');

    InputDecoration fieldDecoration(String hintText) {
      return InputDecoration(
        border: const OutlineInputBorder(),
        hintText: hintText,
      );
    }

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(20, 5, 20, 5),
      textStyle: Theme.of(context).textTheme.labelLarge,
    );

    const double fieldHeight = 50;

    return BlocConsumer<CharacterCreateCubit, CharacterCreateState>(
      listener: (BuildContext context, CharacterCreateState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, CharacterCreateState state) {
        // Build attribute row
        EdgeInsetsGeometry padding = const EdgeInsets.fromLTRB(10, 2, 10, 2);

        List<Widget> attributeRowWidgets(
          String attributeName,
          int attributeValue,
          void Function() attributeDecrementFunc,
          void Function() attributeIncrementFunc,
        ) {
          return <Widget>[
            Flexible(
              flex: 2,
              child: Container(
                padding: padding,
                alignment: Alignment.centerLeft,
                child: Text(attributeName),
              ),
            ),
            Flexible(
              child: Container(
                padding: padding,
                child: ElevatedButton(
                  onPressed: () {
                    attributeDecrementFunc();
                  },
                  child: const Text('<'),
                ),
              ),
            ),
            Flexible(
              child: Container(
                padding: padding,
                alignment: Alignment.center,
                child: Text('$attributeValue'),
              ),
            ),
            Flexible(
              child: Container(
                padding: padding,
                child: ElevatedButton(
                  onPressed: () {
                    attributeIncrementFunc();
                  },
                  child: const Text('>'),
                ),
              ),
            ),
          ];
        }

        List<Widget> formWidgets = [];

        if (state is CharacterCreateStateError) {
          formWidgets.add(Text(state.exception.message));
        }

        formWidgets.add(
          Container(
            height: fieldHeight,
            width: 300,
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: TextFormField(
              controller: characterNameController,
              autofocus: true,
              decoration: fieldDecoration('Character Name'),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter character name';
                }
                return null;
              },
            ),
          ),
        );

        formWidgets.add(
          Container(
            margin: const EdgeInsets.fromLTRB(0, 4, 0, 4),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: attributeRowWidgets(
                'Strength',
                strength,
                _decrementStrength,
                _incrementStrength,
              ),
            ),
          ),
        );

        formWidgets.add(
          Container(
            margin: const EdgeInsets.fromLTRB(0, 4, 0, 4),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: attributeRowWidgets(
                'Dexterity',
                dexterity,
                _decrementDexterity,
                _incrementDexterity,
              ),
            ),
          ),
        );

        formWidgets.add(
          Container(
            margin: const EdgeInsets.fromLTRB(0, 4, 0, 4),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: attributeRowWidgets(
                'Intelligence',
                intelligence,
                _decrementIntelligence,
                _incrementIntelligence,
              ),
            ),
          ),
        );

        formWidgets.add(
          Container(
            height: fieldHeight,
            width: 200,
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: ElevatedButton(
              onPressed: () {
                if (_formKey.currentState!.validate()) {
                  _createCharacter(context);
                }
              },
              style: buttonStyle,
              child: const Text('Create Character'),
            ),
          ),
        );

        return Container(
          margin: const EdgeInsets.fromLTRB(20, 10, 20, 10),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: formWidgets,
            ),
          ),
        );
      },
    );
  }
}
