part of 'character_collection_cubit.dart';

@immutable
abstract class CharacterCollectionState extends Equatable {
  const CharacterCollectionState();
}

@immutable
class CharacterCollectionStateInitial extends CharacterCollectionState {
  const CharacterCollectionStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCollectionStateError extends CharacterCollectionState {
  const CharacterCollectionStateError();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCollectionStateLoading extends CharacterCollectionState {
  const CharacterCollectionStateLoading();

  @override
  List<Object> get props => [];
}

@immutable
class CharacterCollectionStateLoaded extends CharacterCollectionState {
  final List<CharacterRecord>? characterRecords;

  const CharacterCollectionStateLoaded({
    required this.characterRecords,
  });

  @override
  List<Object?> get props => [characterRecords];
}
