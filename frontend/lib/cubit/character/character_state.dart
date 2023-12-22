part of 'character_cubit.dart';

@immutable
abstract class CharacterState extends Equatable {
  const CharacterState();
}

@immutable
class CharacterStateInitial extends CharacterState {
  const CharacterStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterStateError extends CharacterState {
  final CharacterRecord characterRecord;
  final String message;
  const CharacterStateError({
    required this.characterRecord,
    required this.message,
  });

  @override
  List<Object> get props => [message];
}

@immutable
class CharacterStateLoading extends CharacterState {
  const CharacterStateLoading();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterStateLoaded extends CharacterState {
  final List<CharacterRecord>? characterRecords;
  final CharacterRecord? currentCharacterRecord;

  const CharacterStateLoaded({
    required this.characterRecords,
    this.currentCharacterRecord,
  });

  @override
  List<Object?> get props => [characterRecords, currentCharacterRecord];
}

@immutable
class CharacterStateSelecting extends CharacterState {
  final CharacterRecord characterRecord;

  const CharacterStateSelecting({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}

@immutable
class CharacterStateSelected extends CharacterState {
  final CharacterRecord characterRecord;

  const CharacterStateSelected({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}

@immutable
class CharacterStateEntering extends CharacterState {
  final CharacterRecord characterRecord;

  const CharacterStateEntering({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}

@immutable
class CharacterStateEntered extends CharacterState {
  final CharacterRecord characterRecord;

  const CharacterStateEntered({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}

@immutable
class CharacterStateExiting extends CharacterState {
  final CharacterRecord characterRecord;

  const CharacterStateExiting({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}

@immutable
class CharacterStateExited extends CharacterState {
  final CharacterRecord characterRecord;

  const CharacterStateExited({
    required this.characterRecord,
  });

  @override
  List<Object> get props => [characterRecord];
}
