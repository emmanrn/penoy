<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import './app.css';
	import {
		cancelRole,
		chooseRole,
		chooseWord,
		confirmRole,
		connect,
		createRoom,
		joinRoom,
		leaveRoom,
		nextWord,
		sendGuess,
		submitWord,
		type IncomingWSMessage,
		type Mode
	} from './lib/ws';
	import Landing from './lib/components/Landing.svelte';
	import Waiting from './lib/components/Waiting.svelte';
	import Playing from './lib/components/Playing.svelte';
	import RoleChoice from './lib/components/RoleChoice.svelte';
	import type { RoleActions } from './lib/types';
	import { startCountdown } from './lib/utils';

	let gameState: 'landing' | 'waiting' | 'playing' | 'choosing_role' = 'landing';
	let role: string | null = null;
	let word: string | null = null;
	let mode: string | null = null;
	let round_started: boolean = false;
	let roomCodeInput: string = '';
	let roomCode: string | null = null;
	let customwordInput: string = '';
	let confirmedRoles: string[] = [];
	let myChosenRole: string | null = null;
	let wordChoices: string[] = [];
	let remainingSeconds: number = 0;
	let timeUp: boolean = false;
	let timerInterval: ReturnType<typeof setInterval> | null = null;
	let guessInput: string = '';
	let guessCorrect: boolean = false;

	let connected: boolean = false;
	let roomNotFound: boolean = false;

	function beginCountdown(endsAtMs: number): void {
		if (timerInterval) clearInterval(timerInterval);
		timerInterval = startCountdown(endsAtMs, (secs) => {
			remainingSeconds = secs;
			if (secs <= 0) stopCountdown();
		});
	}

	function stopCountdown(): void {
		if (timerInterval) {
			clearInterval(timerInterval);
			timerInterval = null;
		}
	}

	function handleMessage(msg: IncomingWSMessage) {
		switch (msg.type) {
			case 'room_joined':
				roomCode = msg.roomCode ?? null;
				word = null;
				gameState = 'waiting';
				roomNotFound = false;
				if (roomCode) {
					window.history.pushState({}, '', `/room/${roomCode}`);
				}
				break;
			case 'room_not_found':
				// show error message in landing page
				roomNotFound = true;
				break;
			case 'round_start':
				gameState = 'playing';
				word = msg.word ?? null;
				round_started = true;
				timeUp = false;
				guessCorrect = false;
				role = msg.role ?? role;
				wordChoices = [];
				if (msg.roundEndsAt) beginCountdown(msg.roundEndsAt);
				break;
			case 'awaiting_word':
				gameState = 'playing';
				word = null;
				mode = 'custom';
				role = 'clue-giver';
				break;
			case 'word_locked':
				word = msg.word ?? null;
				break;
			case 'round_end':
				gameState = 'playing';
				word = null;
				round_started = false;
				guessCorrect = false;
				role = msg.role ?? role;
				stopCountdown();
				break;
			case 'opponent_left':
				gameState = 'waiting';
				word = null;
				round_started = false;
				role = null;
				stopCountdown();
				break;
			case 'choose_role_phase':
				gameState = 'choosing_role';
				confirmedRoles = [];
				break;
			case 'role_state':
				confirmedRoles = msg.confirmedRoles ?? [];
				break;
			case 'confirm_rejected':
				confirmedRoles = msg.confirmedRoles ?? [];
				myChosenRole = null;
				break;
			case 'word_choices':
				gameState = 'playing';
				wordChoices = msg.words ?? [];
				word = null;
				role = msg.role ?? role;
				break;
			case 'waiting_word':
				gameState = 'playing';
				word = null;
				role = msg.role ?? role;
				break;
			case 'time_up':
				timeUp = true;
				word = msg.word ?? word;
				stopCountdown();
				break;
			case 'correct_guess':
				word = msg.word ?? null;
				guessCorrect = true;
				stopCountdown();
				break;
		}
	}

	onMount(() => {
		connect(handleMessage, () => {
			connected = true;
			const match = window.location.pathname.match(/^\/room\/([A-Za-z0-9]+)$/);
			if (match) {
				joinRoom(match[1]);
			}
		});
	});

	onDestroy(() => {
		stopCountdown();
	});

	function handleCreateRoom(mode: Mode): void {
		createRoom(mode); // can be 'custom' test it later
	}

	function handleJoinRoom(): void {
		joinRoom(roomCodeInput);
	}

	function handleSubmitWord(): void {
		submitWord(customwordInput);
		customwordInput = '';
	}

	function handleGuess(): void {
		sendGuess(guessInput);
		guessInput = '';
	}

	function handleLeaveRoom(): void {
		leaveRoom();
		gameState = 'landing';
		role = null;
		word = null;
		myChosenRole = null;
		mode = null;
		roomCodeInput = '';
		roomNotFound = false;
		round_started = false;
		window.history.pushState({}, '', '/');
	}

	function handleChooseRole(r: 'clue-giver' | 'guesser'): void {
		myChosenRole = r;
		chooseRole(r);
	}

	function handleCancelRole(): void {
		myChosenRole = null;
		cancelRole();
	}

	function handleConfirmRole(): void {
		confirmRole();
	}

	function handleChooseWord(chosen: string): void {
		chooseWord(chosen);
	}

	const roleActions: RoleActions = {
		choose: handleChooseRole,
		cancel: handleCancelRole,
		confirm: handleConfirmRole
	};
</script>

{#if gameState === 'landing'}
	<Landing {handleCreateRoom} {handleJoinRoom} bind:roomNotFound bind:roomCodeInput />
{:else if gameState === 'waiting'}
	<Waiting {roomCode} {handleLeaveRoom} />
{:else if gameState === 'playing'}
	<Playing
		{role}
		{word}
		{mode}
		{round_started}
		{timeUp}
		{guessCorrect}
		{wordChoices}
		{remainingSeconds}
		bind:customwordInput
		bind:guessInput
		{handleSubmitWord}
		{handleLeaveRoom}
		{handleChooseWord}
		{handleGuess}
	/>
{:else if gameState === 'choosing_role'}
	<!-- role component -->
	<RoleChoice {roleActions} {confirmedRoles} {myChosenRole} />
{/if}
