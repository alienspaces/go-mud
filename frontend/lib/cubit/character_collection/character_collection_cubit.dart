import 'dart:async';
import 'package:bloc/bloc.dart';
import 'package:meta/meta.dart';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/character_create/character_create_cubit.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/logger.dart';

part 'character_collection_state.dart';

class CharacterCollectionCubit extends Cubit<CharacterCollectionState> {
  final Map<String, String> config;
  final RepositoryCollection repositories;
  final DungeonActionCubit dungeonActionCubit;
  final CharacterCreateCubit characterCreateCubit;

  List<CharacterRecord>? characterRecords;

  StreamSubscription? dungeonActionSubscription;
  StreamSubscription? characterCreateSubscription;

  CharacterCollectionCubit({
    required this.config,
    required this.repositories,
    required this.dungeonActionCubit,
    required this.characterCreateCubit,
  }) : super(const CharacterCollectionStateInitial()) {
    dungeonActionSubscription?.cancel();
    dungeonActionSubscription = dungeonActionCubit.stream.listen((state) {
      final log =
          getLogger('CharacterCollectionCubit', 'dungeonActionCubit(listener)');
      log.warning('Dungeon action cubit emitted state');
      if (state is DungeonActionStateError) {
        loadCharacters();
      }
    });

    characterCreateSubscription?.cancel();
    characterCreateSubscription = characterCreateCubit.stream.listen((state) {
      final log = getLogger(
          'CharacterCollectionCubit', 'characterCreateCubit(listener)');
      log.warning('Character create cubit emitted state');
      if (state is CharacterCreateStateCreated) {
        loadCharacters();
      }
    });
  }

  Future<void> loadCharacters() async {
    final log = getLogger('CharacterCollectionCubit', 'loadCharacters');
    log.fine('Loading characters...');
    emit(const CharacterCollectionStateLoading());

    List<CharacterRecord>? characterRecords;

    try {
      characterRecords = await repositories.characterRepository.getMany();
    } catch (err) {
      emit(const CharacterCollectionStateError());
      return;
    }

    emit(CharacterCollectionStateLoaded(characterRecords: characterRecords));
  }
}
