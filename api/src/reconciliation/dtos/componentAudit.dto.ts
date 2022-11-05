export interface ComponentAudit {
    reconcileId: number;
    duration: number;
    status: string;
    startDateTime: string;
    approvedBy?: string;
}