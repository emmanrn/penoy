<script lang="ts">
	interface WaitingProps {
		handleLeaveRoom: () => void;
	}
	let { handleLeaveRoom }: WaitingProps = $props();
</script>

<div class="my-auto flex h-full flex-col items-center justify-center wrap-break-word">
	<h1>Waiting for opponent..</h1>

	<button class="btn preset-filled-error-500 font-semibold" onclick={handleLeaveRoom}
		>Leave Room</button
	>
</div>
