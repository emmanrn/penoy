<script lang="ts">
	import type { RoleActions } from '../types';

	interface RoleChoiceProps {
		roleActions: RoleActions;
		confirmedRoles: string[];
		myChosenRole: string | null;
	}

	let { roleActions, confirmedRoles, myChosenRole }: RoleChoiceProps = $props();
</script>

<div class="my-auto flex h-full w-full flex-col items-center justify-center">
	<h1>Choose your role</h1>

	<div class="flex flex-row items-center justify-center gap-50">
		<button
			class="my-20 btn preset-filled"
			disabled={confirmedRoles.includes('clue-giver')}
			onclick={() => roleActions.choose('clue-giver')}
		>
			Clue Giver
		</button>
		<button
			class="my-20 btn preset-filled"
			disabled={confirmedRoles.includes('guesser')}
			onclick={() => roleActions.choose('guesser')}
		>
			Guesser
		</button>
	</div>

	<div class="flex flex-row items-center justify-center gap-4">
		{#if myChosenRole}
			<button class="btn preset-filled-brand" onclick={roleActions.confirm}>Confirm</button>
			<button class="btn preset-outlined" onclick={roleActions.cancel}>Cancel</button>
		{/if}
	</div>

	<div class="my-5">
		{#if confirmedRoles.length > 0}
			<p>{confirmedRoles.length} of 2 confirmed</p>
		{/if}
	</div>
</div>
