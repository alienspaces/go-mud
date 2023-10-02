import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_current_state.dart';

class CharacterCurrentCubit extends Cubit<CharacterCurrentState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;
  final DungeonActionCubit dungeonActionCubit;

  CharacterRecord? characterRecord;

  // streamSubscription is listening to events from the dungeon action
  // cubit, specifically events that might require this cubit to refresh
  // and emit an updated state.
  StreamSubscription? streamSubscription;

  CharacterCurrentCubit({
    required this.config,
    required this.repositories,
    required this.dungeonActionCubit,
  }) : super(const CharacterCurrentStateInitial()) {
    streamSubscription?.cancel();
    streamSubscription = dungeonActionCubit.stream.listen((state) {
      final log =
          getLogger('CharacterCurrentCubit', 'dungeonActionCubit(listener)');
      log.warning('Dungeon action cubit emitted state');
      if (state is DungeonActionStateError) {
        log.warning('Dungeon action cubit emitted error event');
        if (characterRecord != null) {
          log.warning('Clearing character record');
          clearCharacter();
        }
      }
    });
  }

  void clearCharacter() {
    final log = getLogger('CharacterCurrentCubit', 'clearCharacter');
    log.warning('Clearing character');
    characterRecord = null;
    emit(const CharacterCurrentStateInitial());
  }

  Future<void> loadCharacter(String characterID) async {
    final log = getLogger('CharacterCurrentCubit', 'loadCharacter');
    log.fine('Loading character ID $characterID');

    emit(const CharacterCurrentStateLoading());

    CharacterRecord? loadedCharacterRecord =
        await repositories.characterRepository.getOne(characterID);

    log.fine('Loaded character $loadedCharacterRecord');

    if (loadedCharacterRecord != null) {
      emit(CharacterCurrentStateSelected(
          characterRecord: loadedCharacterRecord));
    }
  }

  Future<void> refreshCharacter(String characterID) async {
    final log = getLogger('CharacterCurrentCubit', 'refreshCharacter');
    log.fine('Refreshing character ID $characterID');

    emit(const CharacterCurrentStateLoading());

    CharacterRecord? characterRecord =
        await repositories.characterRepository.getOne(characterID);

    if (characterRecord != null) {
      return selectCharacter(characterRecord);
    }
  }

  Future<void> selectCharacter(CharacterRecord characterRecord) async {
    this.characterRecord = characterRecord;
    emit(CharacterCurrentStateSelected(characterRecord: characterRecord));
  }
}
