part of 'character_create_cubit.dart';

@immutable
abstract class CharacterCreateState extends Equatable {
  const CharacterCreateState();
}

@immutable
class CharacterCreateStateInitial extends CharacterCreateState {
  const CharacterCreateStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCreateStateCreating extends CharacterCreateState {
  const CharacterCreateStateCreating();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCreateStateCreated extends CharacterCreateState {
  final CharacterRecord? characterRecord;

  const CharacterCreateStateCreated({
    this.characterRecord,
  });

  @override
  List<Object?> get props => [characterRecord];
}

@immutable
class CharacterCreateStateError extends CharacterCreateState {
  final CreateCharacterRecord characterRecord;
  final String message;
  const CharacterCreateStateError({
    required this.characterRecord,
    required this.message,
  });

  @override
  List<Object> get props => [message];
}
