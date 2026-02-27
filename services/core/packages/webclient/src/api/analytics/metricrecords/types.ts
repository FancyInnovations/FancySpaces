export interface RecordQueryResult {
    time: number;
    rows_count: number;
    records: MetricRecord[];
}

export interface MetricRecord {
    project_id: string
    metric_id: string
    label: string
    timestamp: Date
    value: number
}
