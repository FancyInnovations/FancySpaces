export interface LogRecord {
    id: string; // generated UUID
    project_id?: string;
    service: string;
    timestamp: Date;
    level: string;
    message: string;
    properties?: Record<string, string>;
}

export interface LogsQueryResult {
    time: number;
    rows_count: number;
    records: LogRecord[];
}

export interface LogsQueryVolumeResult {
    time: number;
    rows_count: number;
    volume: {
        timestamp: number;
        count: number;
    }[];
}
