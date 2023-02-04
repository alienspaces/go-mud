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
class CharacterStateCreate extends CharacterState {
  const CharacterStateCreate();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterStateCreateError extends CharacterState {
  final CreateCharacterRecord characterRecord;
  final String message;
  const CharacterStateCreateError({
    required this.characterRecord,
    required this.message,
  });

  @override
  List<Object> get props => [characterRecord, message];
}

@immutable
class CharacterStateError extends CharacterState {
  final CharacterRecord characterRecord;
  final String message;
  const CharacterStateError(
      {required this.characterRecord, required this.message});

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

  const CharacterStateLoaded(
      {required this.characterRecords, this.currentCharacterRecord});

  @override
  List<Object?> get props => [characterRecords, currentCharacterRecord];
}

@immutable
class CharacterStateLoadError extends CharacterState {
  const CharacterStateLoadError();

  @override
  List<Object?> get props => [];
}

@immutable
class CharacterStateSelected extends CharacterState {
  final CharacterRecord characterRecord;
  const CharacterStateSelected({required this.characterRecord});

  @override
  List<Object> get props => [characterRecord];
}
