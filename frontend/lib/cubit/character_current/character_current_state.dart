part of 'character_current_cubit.dart';

@immutable
abstract class CharacterCurrentState extends Equatable {
  const CharacterCurrentState();
}

@immutable
class CharacterCurrentStateInitial extends CharacterCurrentState {
  const CharacterCurrentStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCurrentStateError extends CharacterCurrentState {
  final CharacterRecord characterRecord;
  final String message;
  const CharacterCurrentStateError({
    required this.characterRecord,
    required this.message,
  });

  @override
  List<Object> get props => [message];
}

@immutable
class CharacterCurrentStateLoading extends CharacterCurrentState {
  const CharacterCurrentStateLoading();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCurrentStateLoaded extends CharacterCurrentState {
  final List<CharacterRecord>? characterRecords;
  final CharacterRecord? currentCharacterRecord;

  const CharacterCurrentStateLoaded({
    required this.characterRecords,
    this.currentCharacterRecord,
  });

  @override
  List<Object?> get props => [characterRecords, currentCharacterRecord];
}

@immutable
class CharacterCurrentStateSelected extends CharacterCurrentState {
  final CharacterRecord characterRecord;

  const CharacterCurrentStateSelected({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}
