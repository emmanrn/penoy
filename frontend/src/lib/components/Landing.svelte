<script lang="ts">
	import { handleKeyDown } from '../utils';
	import type { Mode } from '../ws';

	interface LandingProps {
		handleCreateRoom: (mode: Mode) => void;
		handleJoinRoom: () => void;
		roomCodeInput: string;
	}
	let { roomCodeInput = $bindable(), handleCreateRoom, handleJoinRoom }: LandingProps = $props();

	let showInput: boolean = $state(false);
</script>

<div class="align-center mx-h-md my-auto flex h-full w-full flex-col justify-center space-y-4">
	<div>
		<h1>Penoy</h1>
	</div>

	{#if !showInput}
		<div class="max-h-md mx-auto my-auto h-full w-full max-w-md space-y-4">
			<button
				type="button"
				class="btn preset-filled-success-500"
				onclick={() => handleCreateRoom('random')}>Create Room</button
			>
			<button
				type="button"
				class="btn preset-filled-success-500"
				onclick={() => {
					showInput = true;
					handleJoinRoom();
				}}>Join Room</button
			>
		</div>
	{:else}
		<div class="max-h-md mx-auto my-auto h-full w-full max-w-md space-y-4">
			<label class="label">
				<span class="my-auto label-text">Room Code</span>
				<input
					class="input"
					bind:value={roomCodeInput}
					onkeydown={(e) => handleKeyDown(e, handleJoinRoom)}
					type="text"
					placeholder="Enter room code to join"
				/>
			</label>
			<div class="mx-auto flex items-center justify-center space-x-4">
				<button type="button" class="btn preset-filled-success-500" onclick={handleJoinRoom}
					>Join Room</button
				>
				<button
					type="button"
					class="btn preset-filled-error-500"
					onclick={() => (showInput = false)}>Cancel</button
				>
			</div>
		</div>
	{/if}
</div>
