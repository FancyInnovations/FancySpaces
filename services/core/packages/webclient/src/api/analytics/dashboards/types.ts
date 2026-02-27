export interface Chart {
    type: string;
    name: string;
    options: Record<string, any>;
}

export interface Dashboard {
    dashboard_id: string;
    project_id: string;
    name: string;
    summary: string;
    created_at: string;
    public: boolean
    charts: Chart[];
}
