export interface ComponentAudit {
    reconcileId: number;
    duration: number;
    status: string;
    startDateTime: string;
    approvedBy?: string;
}

export interface ApprovedByDto {
    email?: string;
    compName: string;
    envName: string;
    teamName: string;
    rid: number;
}
