export interface RoleActions {
	choose: (role: 'clue-giver' | 'guesser') => void;
	cancel: () => void;
	confirm: () => void;
}
