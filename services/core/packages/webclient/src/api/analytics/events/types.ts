export interface Event {
    id: string; // generated id
    name: string;
    timestamp: Date;
    properties: Record<string, string>;
}

export interface EventQueryResult {
    time: number;
    count: number;
    events: Event[];
}

export interface EventQueryVolumeResult {
    time: number;
    count: number;
    volume: {
        timestamp: Date;
        count: number;
    }[];
}
