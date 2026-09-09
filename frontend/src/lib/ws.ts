// mirroring the OutcomeMessage type from the backend
export type IncomingWSMessage = {
	type: string;
	role: string | null;
	words: string[];
	word: string | null;
	outcome: string | null;
	confirmedRoles: string[];
	roundEndsAt?: number;
	roomCode: string | null;
};

// room mode
export type Mode = 'random' | 'custom';

let socket: WebSocket | null = null;
let messageHandler: ((msg: IncomingWSMessage) => void) | null = null;

// called once onMount of the application, so once a user opens the website, they automatically
// connect to a websocket
export function connect(onMessage: (msg: IncomingWSMessage) => void): void {
	messageHandler = onMessage;
	socket = new WebSocket('ws://localhost:8080/ws');

	socket.onopen = () => {
		console.log('WS connected');
	};

	// whenever the server sends a message
	socket.onmessage = (event: MessageEvent) => {
		const msg: IncomingWSMessage = JSON.parse(event.data);
		console.log('WS received:', msg);
		if (messageHandler) messageHandler(msg);
	};

	socket.onclose = () => {
		console.log('WS closed');
	};

	socket.onerror = (err: Event) => {
		console.log('WS error:', err);
	};
}

function send(msg: Record<string, unknown>): void {
	if (socket && socket.readyState === WebSocket.OPEN) {
		socket.send(JSON.stringify(msg));
	} else {
		console.warn('WS not open, cannot send:', msg);
	}
}

export function createRoom(mode: Mode): void {
	send({ type: 'create_room', mode });
}

// functions used that sends the jsons to the server
export function joinRoom(roomCode: string): void {
	send({ type: 'join_room', roomCode });
}

export function leaveRoom(): void {
	send({ type: 'leave_room' });
}

export function submitWord(word: string): void {
	send({ type: 'submit_word', word });
}

export function nextWord(): void {
	send({ type: 'next_word' });
}

export function chooseRole(role: 'clue-giver' | 'guesser'): void {
	send({ type: 'choose_role', role });
}

export function cancelRole(): void {
	send({ type: 'cancel_role' });
}

export function confirmRole(): void {
	send({ type: 'confirm_role' });
}

export function chooseWord(word: string): void {
	send({ type: 'choose_word', word });
}

export function sendGuess(guess: string): void {
	send({ type: 'guess', guess });
}

export function disconnect(): void {
	if (socket) {
		socket.close();
		socket = null;
	}
}
