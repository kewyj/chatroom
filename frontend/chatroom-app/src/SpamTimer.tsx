export class TimerExample {
    private timerId: ReturnType<typeof setTimeout> | null = null;

    // Function to set isVisible to true
    setIsVisibleTrue(setVisibility : (isVisible : boolean) => void) {
        console.log('isVisible set to true');

        // If there's an existing timer, clear it
        if (this.timerId) {
            clearTimeout(this.timerId);
        }

        // Set a new timer to set isVisible back to false after 5 seconds
        this.timerId = setTimeout(() => {
            setVisibility(false)
            console.log('isVisible set to false after 3 seconds');
        }, 3000);
    }
}