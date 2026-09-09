<script lang="ts">
	interface WaitingProps {
		roomCode: string | null;
		handleLeaveRoom: () => void;
	}
	let { roomCode, handleLeaveRoom }: WaitingProps = $props();
</script>

<div class="my-auto flex h-full flex-col items-center justify-center wrap-break-word">
	<h2>Waiting for opponent..</h2>
	{#if roomCode}
		<p>Room code</p>
		<h1 class="tracking-[12px]">{roomCode}</h1>
	{/if}

	<button class="btn preset-filled-error-500 font-semibold" onclick={handleLeaveRoom}
		>Leave Room</button
	>
</div>
