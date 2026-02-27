export interface Template {
    name: string;
    description: string;
    metrics: TemplateMetric[];
    dashboards: TemplateDashboard[]
}

export interface TemplateMetric {
    name: string;

    multi_sender: boolean;
    aggregation_interval?: number;
    apply_extra_aggregation?: boolean;

    pull_metric?: boolean;
    pull_interval?: number;
    pull_url?: string;
}

export interface TemplateDashboard {
    name: string;
    summary: string;
    public: boolean;
    charts: TemplateChart[];
}

export interface TemplateChart {
    type: string;
    name: string;
    options: Record<string, any>;
}
