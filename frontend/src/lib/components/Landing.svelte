<script lang="ts">
	import { handleKeyDown } from '../utils';

	interface LandingProps {
		handleCreateRoom: () => void;
		roomCodeInput: string;
	}
	let { roomCodeInput = $bindable(), handleCreateRoom }: LandingProps = $props();
</script>

<div class="align-center mx-h-md my-auto flex h-full w-full flex-col justify-center space-y-4">
	<div>
		<h1>Penoy</h1>
	</div>

	<div class="max-h-md mx-auto my-auto h-full w-full max-w-md space-y-4">
		<label class="label">
			<span class="my-auto label-text">Room Code</span>
			<input
				class="input"
				bind:value={roomCodeInput}
				onkeydown={(e) => handleKeyDown(e, handleCreateRoom)}
				type="text"
				placeholder="Enter room code to join"
			/>
		</label>
		<button type="button" class="btn preset-filled-success-500" onclick={handleCreateRoom}
			>Join Room</button
		>
	</div>
</div>
