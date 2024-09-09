export class Snowflake {
    private sequence: number = 0
    private lastTimestamp: number = -1
    private workerId: number

    constructor(workerId: number) {
        this.workerId = workerId
    }

    public nextId(): string {
        let timestamp = Date.now()
        if (timestamp < this.lastTimestamp) {
            throw new Error('Invalid system clock')
        }

        if (this.lastTimestamp === timestamp) {
            this.sequence = (this.sequence + 1) & 0xfff
            if (this.sequence === 0) {
                // Sequence exhausted, wait for the next second
                timestamp = this.waitNextMillis(timestamp)
            }
        } else {
            this.sequence = 0
        }

        this.lastTimestamp = timestamp
        return (
            ((timestamp - 1288834974657) << 22) |
            (this.workerId << 12) |
            this.sequence
        ).toString()
    }

    private waitNextMillis(timestamp: number): number {
        while (timestamp === this.lastTimestamp) {
            timestamp = Date.now()
        }
        return timestamp
    }
}
