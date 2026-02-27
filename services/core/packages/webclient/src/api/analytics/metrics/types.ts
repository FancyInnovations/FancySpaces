export interface Metric {
    project_id: string;
    metric_id: string;
    name: string;

    multi_sender: boolean;
    aggregation_interval?: number;
    apply_extra_aggregation?: boolean;

    pull_metric?: boolean;
    pull_interval?: number;
    pull_url?: string;
}
