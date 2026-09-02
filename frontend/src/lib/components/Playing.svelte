<script lang="ts">
	import { handleKeyDown } from '../utils';

	interface PlayingProps {
		role: string | null;
		word: string | null;
		mode: string | null;
		round_started: boolean;
		timeUp: boolean;
		wordChoices: string[];
		remainingSeconds: number;
		customwordInput: string;
		guessInput: string;
		guessCorrect: boolean;
		handleSubmitWord: () => void;
		handleNextWord: () => void;
		handleLeaveRoom: () => void;
		handleChooseWord: (word: string) => void;
		handleGuess: () => void;
	}

	let {
		role,
		word,
		mode,
		round_started,
		timeUp,
		wordChoices,
		remainingSeconds,
		customwordInput = $bindable(),
		guessInput = $bindable(),
		guessCorrect,
		handleSubmitWord,
		handleNextWord,
		handleLeaveRoom,
		handleChooseWord,
		handleGuess
	}: PlayingProps = $props();
</script>

<div class="my-auto flex h-full w-full flex-col items-center justify-center">
	<div class="flex flex-col space-y-4">
		<h2>You are the {role === 'clue-giver' ? 'Clue Giver' : 'Guesser'}</h2>
		<div class="mx-auto flex w-full flex-col items-start justify-center space-y-4">
			{#if word}
				{#if role === 'guesser'}
					{#if guessCorrect}
						<div class="mx-auto flex flex-col items-center justify-center">
							<h1>Correct</h1>
						</div>
					{/if}
				{:else if role === 'clue-giver'}
					{#if guessCorrect}
						<div class="mx-auto flex flex-col items-center justify-center">
							<h1>Word: {word}</h1>
							<h1>Guesser got it correctly</h1>
						</div>
					{:else}
						<div class="mx-auto flex flex-col items-center justify-center">
							<h1>Word: {word}</h1>
						</div>
					{/if}
				{/if}
			{:else if role === 'clue-giver' && wordChoices.length > 0}
				<p>Pick a word:</p>
				<div class="flex space-x-2">
					{#each wordChoices as choice}
						<button class="btn preset-filled" onclick={() => handleChooseWord(choice)}
							>{choice}</button
						>
					{/each}
				</div>
			{:else if word === null && !round_started}
				<h1>Waiting for word...</h1>
			{/if}

			{#if role === 'clue-giver' && mode === 'custom'}
				<input
					class="input"
					onkeydown={(e) => handleKeyDown(e, handleSubmitWord)}
					bind:value={customwordInput}
					placeholder="Type a custom word"
				/>
				<button class="btn preset-filled-brand" onclick={handleSubmitWord}>Submit</button>
			{:else if round_started}
				<!-- Timer here -->
				{#if timeUp}
					<div class="mx-auto flex flex-col items-center justify-center">
						{#if role === 'guesser'}
							<h1>Time's up!</h1>
							<h1>The word was {word}</h1>
						{:else}
							<h1>Time's up!</h1>
						{/if}
					</div>
				{:else}
					<h2 class="self-center text-2xl font-bold {remainingSeconds <= 10 ? 'text-red-500' : ''}">
						{remainingSeconds}s
					</h2>
				{/if}

				{#if role === 'guesser' && !timeUp}
					<div class="mx-auto flex flex-col items-center justify-center">
						<p>Round Started</p>
						<div class="my-4 flex flex-row items-center justify-center space-x-4">
							<input
								class="input"
								bind:value={guessInput}
								onkeydown={(e) => handleKeyDown(e, handleGuess)}
								placeholder="Type your guess"
							/>
							<button class="btn preset-filled-brand" onclick={handleGuess}>Guess</button>
						</div>
					</div>
				{/if}
			{/if}
		</div>
	</div>
	{#if round_started}
		<div class="mt-5 flex items-center justify-center space-x-4">
			<!-- <button class="btn preset-outlined capitalize" onclick={handleNextWord}>Next Word</button> -->
			<button class="btn preset-filled-error-500 font-semibold" onclick={handleLeaveRoom}
				>Leave Room</button
			>
		</div>
	{/if}
</div>
