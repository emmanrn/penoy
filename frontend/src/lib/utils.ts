// this is for when pressing enter in an input field to call the submit handler
export function handleKeyDown(event: KeyboardEvent, onEnter: () => void): void {
	if (event.key === 'Enter') {
		onEnter();
	}
}

// timer function
export function startCountdown(
	endsAtMs: number,
	onTick: (secondsLeft: number) => void
): ReturnType<typeof setInterval> {
	const tick = () => {
		const secondsLeft = Math.max(0, Math.round((endsAtMs - Date.now()) / 1000));
		onTick(secondsLeft);
	};

	tick();
	return setInterval(tick, 1000);
}
